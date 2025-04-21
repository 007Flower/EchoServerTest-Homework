package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxMessageSize   = 1024
	inactivityPeriod = 30 * time.Second
	logDir           = "client_logs"
)

func main() {
	port := flag.String("port", "4000", "Port for the server to listen on")
	flag.Parse()

	address := ":" + *port
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	logEvent(fmt.Sprintf("Server listening on %s", address))

	for {
		conn, err := listener.Accept()
		if err != nil {
			logEvent(fmt.Sprintf("Error accepting: %v", err))
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	logEvent(fmt.Sprintf("Client connected: %s", clientAddr))

	// Ensure log directory exists
	os.MkdirAll(logDir, 0755)

	// Open log file for client messages
	safeFileName := strings.ReplaceAll(clientAddr, ":", "_") + ".log"
	logFilePath := filepath.Join(logDir, safeFileName)
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(conn, "Server error: unable to open log file\n")
		return
	}
	defer logFile.Close()

	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, maxMessageSize), maxMessageSize)

	timer := time.NewTimer(inactivityPeriod)
	resetTimer := func() {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timer.Reset(inactivityPeriod)
	}

	done := make(chan bool)

	go func() {
		for scanner.Scan() {
			resetTimer()
			input := strings.TrimSpace(scanner.Text())

			if len(input) > maxMessageSize {
				conn.Write([]byte("Message too long. Truncated.\n"))
				input = input[:maxMessageSize]
			}

			logFile.WriteString(fmt.Sprintf("%s: %s\n", time.Now().Format(time.RFC3339), input))

			// Handle server personality and commands
			switch input {
			case "":
				conn.Write([]byte("Say something...\n"))
			case "hello":
				conn.Write([]byte("Hi there!\n"))
			case "bye":
				conn.Write([]byte("Goodbye!\n"))
				done <- true
				return
			case "/quit":
				conn.Write([]byte("Goodbye!\n"))
				done <- true
				return
			case "/time":
				conn.Write([]byte(time.Now().Format(time.RFC1123) + "\n"))
			default:
				if strings.HasPrefix(input, "/echo ") {
					conn.Write([]byte(strings.TrimPrefix(input, "/echo ") + "\n"))
				} else {
					conn.Write([]byte(input + "\n"))
				}
			}
		}
		if err := scanner.Err(); err != nil {
			logEvent(fmt.Sprintf("Error reading from client %s: %v", clientAddr, err))
		}
		done <- true
	}()

	select {
	case <-timer.C:
		conn.Write([]byte("Disconnected due to inactivity.\n"))
	case <-done:
	}

	logEvent(fmt.Sprintf("Client disconnected: %s", clientAddr))
}

func logEvent(message string) {
	timestamp := time.Now().Format(time.RFC3339)
	fmt.Printf("[%s] %s\n", timestamp, message)
}
