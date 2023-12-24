package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
)

// ExecuteCommand runs a system command and returns its output
func ExecuteCommand(command string) string {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	return string(output)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		command := scanner.Text()

		switch command {
		case "cpu":
			fmt.Fprintln(conn, ExecuteCommand("top -bn1 | grep \"Cpu(s)\" | sed \"s/.*, *\\([0-9.]*\\)%* id.*/\\1/\" | sed \"s/^/100 - /\" | bc"))
		case "memory":
			fmt.Fprintln(conn, ExecuteCommand("free -m | awk '/Mem/{print $3\" \"$2}'"))
		case "network":
			fmt.Fprintln(conn, ExecuteCommand("cat /sys/class/net/[e]*/statistics/{r,t}x_bytes"))
		case "storage":
			fmt.Fprintln(conn, ExecuteCommand("df -Ph | grep mmcblk0p5 | awk '{print $2\" \"$3}' | sed 's/G//g'"))
		case "uptime":
			fmt.Fprintln(conn, ExecuteCommand("cut -f1 -d. /proc/uptime"))
		case "exit":
			return
		default:
			fmt.Fprintln(conn, "Invalid command")
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", ":12345")
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
