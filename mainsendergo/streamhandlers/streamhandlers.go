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

func handleBatteryStream(addr string, port string, doneChan chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", addr, port))
	if err != nil {
		fmt.Printf("Error connecting to Battery server at %s:%s: %v\n", addr, port, err)
		return
	}
	defer conn.Close()

	rm := rosmasterlib.NewRosmaster("/dev/myserial", 115200)
	defer rm.Close()

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
