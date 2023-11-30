package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mattzi/mainsendergo/streamhandlers"
)

func main() {
	ln, err := net.Listen("tcp", "0.0.0.0:6969")
	if err != nil {
		logWithTimestamp("Error setting up TCP server:", err)
		return
	}
	defer ln.Close()
	logWithTimestamp("TCP server listening at 0.0.0.0:6969")

	for {
		conn, err := ln.Accept()
		if err != nil {
			logWithTimestamp("Error accepting connection:", err)
			continue
		}
		logWithTimestamp("Connection accepted")

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var wg sync.WaitGroup
	cmdChan := make(chan string)
	doneReadingChan := make(chan bool)
	healthCheckChan := make(chan string)
	stopChan := make(chan struct{})
	healthCheckActive := false
	lastHealthCheckActive := healthCheckActive

	// Health check routine
	go func() {
		var lastTimestamp int64
		timer := time.NewTimer(3 * time.Second)
		logWithTimestamp("Health check routine started")

		for {

			select {
			case <-timer.C:
				if healthCheckActive {
					logWithTimestamp("Health check failed. No message received in time.")
					close(stopChan)
					return
				}
			case msg := <-healthCheckChan:
				lastHealthCheckActive = healthCheckActive
				healthCheckActive = true
				if lastHealthCheckActive != healthCheckActive {
					timer.Reset(3 * time.Second)
				}

				if strings.HasPrefix(msg, "healthcheck") {
					fields := strings.Fields(msg)
					if len(fields) < 2 {
						logWithTimestamp("Invalid healthcheck message format")
						continue
					}
					timestamp, err := strconv.ParseInt(fields[1], 10, 64)
					if err != nil {
						logWithTimestamp("Invalid timestamp in healthcheck message")
						continue
					}
					if timestamp <= lastTimestamp {
						logWithTimestamp("Received an old timestamp in healthcheck message")
						continue
					}
					lastTimestamp = timestamp
					logWithTimestamp("Health check passed")
					timer.Reset(3 * time.Second)
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
					logWithTimestamp("Error reading from connection:", err)
					break
				}
				scanner = bufio.NewScanner(conn)
			}
		}
		close(doneReadingChan)
	}()

	var cameraDoneChan, lidarDoneChan, batteryDoneChan chan struct{}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	for {
		select {
		case cmd := <-cmdChan:
			logWithTimestamp("Received command:", cmd)
			if strings.HasPrefix(strings.ToLower(cmd), "startstreams") {
				port := strings.TrimSpace(cmd[len("startstreams"):])
				lidarPort, err := strconv.Atoi(port)
				if err != nil {
					logWithTimestamp("Invalid port number for streams:", err)
					continue
				}
				lidarPortStr := strconv.Itoa(lidarPort + 10)
				batteryPortStr := strconv.Itoa(lidarPort + 20)

				cameraDoneChan = make(chan struct{})
				lidarDoneChan = make(chan struct{})
				batteryDoneChan = make(chan struct{})

				wg.Add(3)
				go streamhandlers.HandleCameraStream(conn.RemoteAddr().(*net.TCPAddr).IP.String(), port, cameraDoneChan, &wg)
				go streamhandlers.HandleLidarStream(conn.RemoteAddr().(*net.TCPAddr).IP.String(), lidarPortStr, lidarDoneChan, &wg)
				go streamhandlers.HandleBatteryStream(conn.RemoteAddr().(*net.TCPAddr).IP.String(), batteryPortStr, batteryDoneChan, &wg)
				healthCheckChan <- "healthcheck 0"
			} else if strings.HasPrefix(strings.ToLower(cmd), "stopstreams") {
				logWithTimestamp("Received stopstreams command")
				closeAllChannels(cameraDoneChan, lidarDoneChan, batteryDoneChan)

				wg.Wait()
				logWithTimestamp("All streams stopped")
				cameraDoneChan, lidarDoneChan, batteryDoneChan = nil, nil, nil
			}

		case <-doneReadingChan:
			logWithTimestamp("Connection closed by client.")
			closeAllChannels(cameraDoneChan, lidarDoneChan, batteryDoneChan)
			return
		case <-stopChan:
			logWithTimestamp("Connection closed due to health check failure.")
			closeAllChannels(cameraDoneChan, lidarDoneChan, batteryDoneChan)
			return
		case <-sigChan:
			logWithTimestamp("Connection closed due to SIGINT.")
			closeAllChannels(cameraDoneChan, lidarDoneChan, batteryDoneChan)
			return
		}

	}
}

func closeAllChannels(chans ...chan struct{}) {
	for _, ch := range chans {
		if ch != nil {
			close(ch)
		}
	}
}

func logWithTimestamp(v ...interface{}) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), v)
}
