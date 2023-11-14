package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func handleCameraStream(addr string, port string, doneChan chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Starting ffmpeg on port %s\n", port)
	cmd := exec.Command("ffmpeg", "-input_format", "mjpeg", "-i", "/dev/video0", "-c:v", "copy", "-f", "mjpeg", fmt.Sprintf("tcp://%s:%s", addr, port))

	// Pipe the stdout and stderr of the ffmpeg command to the terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting camera stream on port %s: %v\n", port, err)
		return
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			fmt.Printf("ffmpeg exited with error on port %s: %v\n", port, err)
		} else {
			fmt.Printf("ffmpeg exited normally on port %s\n", port)
		}
		close(doneChan) // Signal that the command has finished
	}()

	// Wait for the stop signal or command completion
	<-doneChan
	if cmd.Process != nil {
		fmt.Printf("Stopping ffmpeg on port %s\n", port)
		cmd.Process.Kill() // Attempt to kill the process if it's still running
	}
}

func main() {
	ln, err := net.Listen("tcp", "0.0.0.0:6969")
	if err != nil {
		fmt.Println("Error setting up TCP server:", err)
		return
	}
	defer ln.Close()
	fmt.Println("TCP server listening at 0.0.0.0:6969")

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connection accepted")

	var wg sync.WaitGroup
	cmdChan := make(chan string)
	doneReadingChan := make(chan bool)

	// Goroutine to read commands from the connection
	go func() {
		scanner := bufio.NewScanner(conn)
		for {
			// Set a deadline for the read operation
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			if scanner.Scan() {
				cmdChan <- scanner.Text()
			}

			if err := scanner.Err(); err != nil {
				netErr, ok := err.(net.Error)
				if !ok || !netErr.Timeout() {
					fmt.Println("Error reading from connection:", err)
					break
				}
			}
			// Reset the error for the next iteration
			scanner = bufio.NewScanner(conn)
		}
		close(doneReadingChan)
	}()

	var doneChan chan struct{}

	for {
		select {
		case cmd := <-cmdChan:
			if strings.HasPrefix(strings.ToLower(cmd), "startstreams") {
				port := strings.TrimSpace(cmd[len("startstreams"):])
				doneChan = make(chan struct{})
				wg.Add(1)
				go handleCameraStream(conn.RemoteAddr().(*net.TCPAddr).IP.String(), port, doneChan, &wg)
			} else if strings.HasPrefix(strings.ToLower(cmd), "stopstreams") {
				fmt.Println("Received stopstreams command")
				if doneChan != nil {
					close(doneChan) // Signal to stop the camera stream
					wg.Wait()       // Wait for the stream handler to finish
					fmt.Println("All streams stopped")
					doneChan = make(chan struct{}) // Reinitialize the channel for future use
				}
			}
		case <-doneReadingChan:
			return
		}
	}

	wg.Wait()
}
