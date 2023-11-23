package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mattzi/mainsendergo/rplidar"
)

func handleCameraStream(addr string, port string, doneChan chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Starting ffmpeg on port %s\n", port)
	cmd := exec.Command("ffmpeg", "-input_format", "mjpeg", "-i", "/dev/video0", "-c:v", "copy", "-f", "mjpeg", fmt.Sprintf("tcp://%s:%s", addr, port))

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
		close(doneChan)
	}()

	<-doneChan
	if cmd.Process != nil {
		fmt.Printf("Stopping ffmpeg on port %s\n", port)
		cmd.Process.Kill()
	}
}

func handleLidarStream(addr string, port string, doneChan chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", addr, port))
	if err != nil {
		fmt.Printf("Error connecting to LiDAR server at %s:%s: %v\n", addr, port, err)
		return
	}
	defer conn.Close()

	lidar := rplidar.NewRPLidar("/dev/rplidar", 115200, time.Second*3)
	err = lidar.Connect()
	if err != nil {
		fmt.Println("Error connecting to RPLidar:", err)
		return
	}
	defer lidar.Disconnect()

	// Retrieve and print device information
	info, err := lidar.GetInfo()
	if err != nil {
		fmt.Println("Error getting info: %v", err)
		return
	}
	fmt.Printf("RPLidar Info: %+v\n", info)

	time.Sleep(time.Second * 1)

	measurements, err := lidar.IterMeasurements()
	if err != nil {
		fmt.Println("Error starting measurements:", err)
		return
	}

	for {
		select {
		case measurement, ok := <-measurements:
			if !ok {
				return // Channel closed, end the loop
			}
			msg := fmt.Sprintf("{\"Angle\": %f, \"Distance\": %f}\n", measurement.Angle, measurement.Distance)
			_, err := conn.Write([]byte(msg))
			if err != nil {
				fmt.Println("Error sending data to server:", err)
				return
			}
		case <-doneChan:
			return
		}
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

	go func() {
		scanner := bufio.NewScanner(conn)
		for {
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
			scanner = bufio.NewScanner(conn)
		}
		close(doneReadingChan)
	}()

	var cameraDoneChan, lidarDoneChan chan struct{}

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

				cameraDoneChan = make(chan struct{})
				lidarDoneChan = make(chan struct{})

				wg.Add(2)
				go handleCameraStream(conn.RemoteAddr().(*net.TCPAddr).IP.String(), port, cameraDoneChan, &wg)
				go handleLidarStream(conn.RemoteAddr().(*net.TCPAddr).IP.String(), lidarPortStr, lidarDoneChan, &wg)
			} else if strings.HasPrefix(strings.ToLower(cmd), "stopstreams") {
				fmt.Println("Received stopstreams command")
				if cameraDoneChan != nil {
					close(cameraDoneChan)
				}
				if lidarDoneChan != nil {
					close(lidarDoneChan)
				}
				wg.Wait()
				fmt.Println("All streams stopped")
				cameraDoneChan, lidarDoneChan = nil, nil
			}
		case <-doneReadingChan:
			return
		}
	}

	wg.Wait()
}
