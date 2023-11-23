package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

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
	healthCheckChan := make(chan string)
	stopChan := make(chan struct{})

	// Health check routine
	go func() {
		var lastTimestamp int64
		for {
			select {
			case <-time.After(12 * time.Second): // 10 seconds + 2 seconds grace period
				fmt.Println("Health check failed. No message received in time.")
				close(stopChan) // Signal to stop the program
				return
			case msg := <-healthCheckChan:
				if strings.HasPrefix(msg, "healthcheck") {
					fields := strings.Fields(msg)
					if len(fields) < 2 {
						fmt.Println("Invalid healthcheck message format")
						continue
					}
					timestamp, err := strconv.ParseInt(fields[1], 10, 64)
					if err != nil {
						fmt.Println("Invalid timestamp in healthcheck message")
						continue
					}
					if timestamp <= lastTimestamp {
						fmt.Println("Received an old timestamp in healthcheck message")
						continue
					}
					lastTimestamp = timestamp
					fmt.Println("Health check passed")
					// Reset the timer
					time.After(12 * time.Second)
				}
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(conn)
		for {
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			if scanner.Scan() {
				text := scanner.Text()
				if strings.HasPrefix(text, "healthcheck") {
					healthCheckChan <- text
				} else {
					cmdChan <- text
				}
			}

			if err := scanner.Err(); err != nil {
				netErr, ok := err.(net.Error)
				if !ok || !netErr.Timeout() {
					fmt.Println("Error reading from connection:", err)
					break
				}
				scanner = bufio.NewScanner(conn) // Reset the scanner
			}
		}
		close(doneReadingChan)
	}()

	var cameraDoneChan, lidarDoneChan, batteryDoneChan chan struct{}

	for {
		select {
		case cmd := <-cmdChan:
			if strings.HasPrefix(strings.ToLower(cmd), "startstreams") {
				port := strings.TrimSpace(cmd[len("startstreams"):])
				lidarPort, err := strconv.Atoi(port)
				if err != nil {
					fmt.Println("Invalid port number for streams:", err)
					continue
				}
				lidarPortStr := strconv.Itoa(lidarPort + 10) // Increment port number by 1 for LiDAR
				batteryPortStr := strconv.Itoa(lidarPort + 20)

				cameraDoneChan = make(chan struct{})
				lidarDoneChan = make(chan struct{})
				batteryDoneChan = make(chan struct{})

				wg.Add(3)
				go streamhandlers.handleCameraStream(conn.RemoteAddr().(*net.TCPAddr).IP.String(), port, cameraDoneChan, &wg)
				go streamhandlers.handleLidarStream(conn.RemoteAddr().(*net.TCPAddr).IP.String(), lidarPortStr, lidarDoneChan, &wg)
				go streamhandlers.handleBatteryStream(conn.RemoteAddr().(*net.TCPAddr).IP.String(), batteryPortStr, cameraDoneChan, &wg)
			} else if strings.HasPrefix(strings.ToLower(cmd), "stopstreams") {
				fmt.Println("Received stopstreams command")
				if cameraDoneChan != nil {
					close(cameraDoneChan)
				}
				if lidarDoneChan != nil {
					close(lidarDoneChan)
				}
				if batteryDoneChan != nil {
					close(lidarDoneChan)
				}
				wg.Wait()
				fmt.Println("All streams stopped")
				cameraDoneChan, lidarDoneChan, batteryDoneChan = nil, nil, nil
			}
		case <-doneReadingChan:
			if cameraDoneChan != nil {
				close(cameraDoneChan)
			}
			if lidarDoneChan != nil {
				close(lidarDoneChan)
			}
			if batteryDoneChan != nil {
				close(batteryDoneChan)
			}

			fmt.Println("Stopping program due to read done.")
			return
		case <-stopChan:
			// Close the stream channels if they are not nil
			if cameraDoneChan != nil {
				close(cameraDoneChan)
			}
			if lidarDoneChan != nil {
				close(lidarDoneChan)
			}
			if batteryDoneChan != nil {
				close(batteryDoneChan)
			}

			fmt.Println("Stopping program due to health check failure.")
			return // Exit the main function and thus the program
		}
	}

	wg.Wait()
}
