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

var commands = map[string]string{
	"cpu":     `top -bn1 | grep "Cpu(s)" | sed "s/.*, *\([0-9.]*\)%* id.*/\1/" | sed "s/^/100 - /" | bc`,
	"memory":  `free -m | awk '/Mem/{print $3 " " $2}'`,
	"network": `cat /sys/class/net/[e]*/statistics/{r,t}x_bytes`,
	"storage": `df -Ph | grep mmcblk0p5 | awk '{print $2 " " $3}' | sed 's/G//g'`,
	"uptime":  `cut -f1 -d. /proc/uptime`,
}

// executeCommand runs a system command and returns its output
func executeCommand(command string) string {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err.Error()
	}

	return string(output)
}

// handleConnection handles one command per connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		log.Printf("Error reading from %s, %s", conn.RemoteAddr(), scanner.Err())
		return
	}

	command := scanner.Text()

	cmd, valid := commands[strings.TrimSpace(string(command))]
	if !valid {
		fmt.Fprintln(conn, "invalid command")
		return
	}

	output := executeCommand(cmd)
	fmt.Fprintln(conn, output)
}

func main() {
	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = ":12345"
	}

	ln, err := net.Listen("tcp", host)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer ln.Close()

	log.Printf("Listening on %s...", host)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		log.Printf("Connection from %s", conn.RemoteAddr())
		go handleConnection(conn)
	}
}
