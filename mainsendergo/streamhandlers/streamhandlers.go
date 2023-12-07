package streamhandlers

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/mattzi/mainsendergo/rosmasterlib"
	"github.com/mattzi/mainsendergo/rplidar"
)

func HandleCameraStream(addr string, port string, doneChan chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		fmt.Printf("Starting ffmpeg on port %s\n", port)
		cmd := exec.Command("ffmpeg", "-input_format", "mjpeg", "-i", "/dev/video0", "-c:v", "copy", "-f", "mjpeg", fmt.Sprintf("tcp://%s:%s", addr, port))

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			fmt.Printf("Error starting camera stream on port %s: %v\n", port, err)
			return
		}

		errChan := make(chan error)
		go func() {
			errChan <- cmd.Wait()
		}()

		select {
		case err := <-errChan:
			if err != nil {
				fmt.Printf("ffmpeg exited with error on port %s: %v\n", port, err)
				fmt.Println("Attempting to restart ffmpeg...")
				continue
			} else {
				fmt.Printf("ffmpeg exited normally on port %s\n", port)
				return
			}
		case <-doneChan:
			if _, ok := <-doneChan; ok {
				close(doneChan)
			}
			if cmd.Process != nil {
				fmt.Printf("Stopping ffmpeg on port %s\n", port)
				cmd.Process.Kill()
			}
			return
		}
	}
}

func HandleLidarStream(addr string, port string, doneChan chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-doneChan:
			return
		default:
			// Attempt to establish TCP connection
			tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", addr, port))
			if err != nil {
				fmt.Println("Error resolving TCP address:", err)
				time.Sleep(time.Second * 2) // Wait before retrying
				continue
			}

			conn, err := net.DialTCP("tcp", nil, tcpAddr)
			if err != nil {
				fmt.Printf("Error connecting to LiDAR server at %s:%s: %v\n", addr, port, err)
				time.Sleep(time.Second * 2) // Wait before retrying
				continue
			}

			// If connected, handle the connection
			handleConnection(conn, doneChan)

			conn.Close()                // Close the connection when handleConnection returns
			time.Sleep(time.Second * 2) // Wait before trying to reconnect
		}
	}
}

func handleConnection(conn *net.TCPConn, doneChan chan struct{}) {
	conn.SetNoDelay(true)

	lidar := rplidar.NewRPLidar("/dev/rplidar", 115200, time.Second*3)
	err := lidar.Connect()
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

func HandleBatteryStream(addr string, port string, doneChan chan struct{}, wg *sync.WaitGroup, rm *rosmasterlib.Rosmaster) {
	defer wg.Done()

	for {
		select {
		case <-doneChan:
			return // Exit if a done signal is received
		default:
			// Attempt to establish TCP connection
			tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", addr, port))
			if err != nil {
				fmt.Println("Error resolving TCP address:", err)
				time.Sleep(time.Second * 2) // Wait before retrying
				continue
			}

			conn, err := net.DialTCP("tcp", nil, tcpAddr)
			if err != nil {
				fmt.Printf("Error connecting to Battery server at %s:%s: %v\n", addr, port, err)
				time.Sleep(time.Second * 2) // Wait before retrying
				continue
			}
			conn.SetNoDelay(true)

			// If connected, handle the connection
			handleBatteryConnection(conn, doneChan, rm)

			conn.Close()                // Close the connection when handleBatteryConnection returns
			time.Sleep(time.Second * 2) // Wait before trying to reconnect
		}
	}
}

func handleBatteryConnection(conn *net.TCPConn, doneChan chan struct{}, rm *rosmasterlib.Rosmaster) {
	for {
		select {
		case <-time.After(1 * time.Second): // Timer for sending battery voltage
			voltage := rm.GetBatteryVoltage()
			fmt.Printf("Aktuelle Batteriespannung: %v V\n", voltage)
			_, err := conn.Write([]byte(fmt.Sprintf("%f\n", voltage)))
			if err != nil {
				fmt.Println("Error sending data to server:", err)
				return
			}
		case <-doneChan: // Listen for the done signal
			return // Exit the function when done signal is received
		}
	}
}
