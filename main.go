package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

var systemdServiceUnit = strings.TrimSpace(`
[Unit]
Description=Stats Server
After=network.target

[Service]
Type=simple
Restart=always
RestartSec=15
ExecStart=/usr/local/bin/stats-server serve

[Install]
WantedBy=default.target
`)

var commands = map[string]string{
	"cpu":     `top -bn1 | grep "Cpu(s)" | sed "s/.*, *\([0-9.]*\)%* id.*/\1/" | sed "s/^/100 - /" | bc`,
	"memory":  `free -m | awk '/Mem/{print $3 " " $2}'`,
	"network": `cat /sys/class/net/[e]*/statistics/{r,t}x_bytes`,
	"storage": `df -Ph | grep mmcblk0p5 | awk '{print $2 " " $3}' | sed 's/G//g'`,
	"uptime":  `cut -f1 -d. /proc/uptime`,
}

func init() {
	log.SetFlags(0)
}

// runShellCommand runs a system command and returns its output
func runShellCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}

// handleConnection handles one command per connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		log.Printf("error reading from %s, %s", conn.RemoteAddr(), scanner.Err())
		return
	}

	command := scanner.Text()
	log.Printf("received command %s from %s", command, conn.RemoteAddr())

	shellCmd, valid := commands[strings.TrimSpace(string(command))]
	if !valid {
		fmt.Fprintln(conn, "invalid command")
		return
	}

	output, err := runShellCommand(shellCmd)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Fprintln(conn, output)
}

func runCommand(args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, output)
	}

	return nil
}

func setupSystemdService() error {
	if _, err := os.Stat("/usr/local/bin/stats-server"); err == nil {
		log.Println("binary already exists in /usr/local/bin, updating...")
	}

	log.Println("copying itself to /usr/local/bin/stats-server")
	if err := runCommand("cp", "-f", os.Args[0], "/usr/local/bin/stats-server"); err != nil {
		return err
	}

	log.Println("writing unit to /etc/systemd/system/stats-server.service")
	if err := os.WriteFile(
		"/etc/systemd/system/stats-server.service",
		[]byte(systemdServiceUnit),
		os.ModePerm,
	); err != nil {
		return err
	}

	log.Println("running systemctl daemon-reload")
	if err := runCommand("systemctl", "daemon-reload"); err != nil {
		return err
	}

	log.Println("running systemctl enable stats-server.service")
	if err := runCommand("systemctl", "enable", "stats-server.service"); err != nil {
		return err
	}

	log.Println("running systemctl (re)start stats-server.service")
	if err := runCommand("systemctl", "restart", "stats-server.service"); err != nil {
		return err
	}

	return nil
}

func startTCPServer() error {
	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = ":12345"
	}

	ln, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}
	defer ln.Close()

	log.Printf("listening on %s...", host)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("error accepting connection: %s", err)
			continue
		}

		log.Printf("connection from %s", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func showHelp() {
	fmt.Println("usage: stats-server [setup|serve]")
	fmt.Println("  setup: setup systemd service")
	fmt.Println("  serve: start tcp server")

	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		showHelp()
	}

	switch os.Args[1] {
	case "setup":
		if err := setupSystemdService(); err != nil {
			log.Fatal(err)
		}

		log.Println("setup complete")
	case "serve":
		if err := startTCPServer(); err != nil {
			log.Fatal(err)
		}

		log.Println("server exited")
	default:
		showHelp()
	}
}
