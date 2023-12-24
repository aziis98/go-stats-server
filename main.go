package main

import (
	"fmt"
	"io"
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

// ExecuteCommand runs a system command and returns its output
func ExecuteCommand(command string) string {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	return string(output)
}

// handleConnection handles one command per connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	command, err := io.ReadAll(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd, valid := commands[strings.TrimSpace(string(command))]
	if !valid {
		fmt.Fprintln(conn, "Invalid command")
		return
	}

	stdout := ExecuteCommand(cmd)
	fmt.Fprintln(conn, stdout)
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

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		go handleConnection(conn)
	}
}
