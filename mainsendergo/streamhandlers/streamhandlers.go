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

// Define log levels
const (
	DEBUG   = 0
	INFO    = 1
	WARNING = 2
	ERROR   = 3
)

// Logging function with timestamp and log level.
func log(level int, logLevel int, v ...interface{}) {
	if level < logLevel {
		return
	}

	levelStr := []string{"DEBUG", "INFO", "WARNING", "ERROR"}[level]
	prefix := fmt.Sprintf("%s [%s] ", time.Now().Format("2006-01-02 15:04:05"), levelStr)
	fmt.Println(prefix, fmt.Sprint(v...))
}

func HandleCameraStream(addr string, port string, doneChan chan struct{}, wg *sync.WaitGroup, logLevel int) {
	defer wg.Done()

	for {
		log(INFO, logLevel, fmt.Sprintf("Starting ffmpeg on port %s", port))
		cmd := exec.Command("ffmpeg", "-input_format", "mjpeg", "-i", "/dev/video0", "-c:v", "copy", "-f", "mjpeg", fmt.Sprintf("tcp://%s:%s", addr, port))

		if logLevel == DEBUG {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}

		if err := cmd.Start(); err != nil {
			log(ERROR, logLevel, fmt.Sprintf("Error starting camera stream on port %s: %v", port, err))
			return
		}

		errChan := make(chan error)
		go func() {
			errChan <- cmd.Wait()
		}()

		select {
		case err := <-errChan:
			if err != nil {
				log(ERROR, logLevel, fmt.Sprintf("ffmpeg exited with error on port %s: %v", port, err))
				log(INFO, logLevel, "Attempting to restart ffmpeg...")
				continue
			} else {
				log(INFO, logLevel, fmt.Sprintf("ffmpeg exited normally on port %s", port))
				return
			}
		case <-doneChan:
			if _, ok := <-doneChan; ok {
				close(doneChan)
			}
			if cmd.Process != nil {
				log(INFO, logLevel, fmt.Sprintf("Stopping ffmpeg on port %s", port))
				cmd.Process.Kill()
			}
			return
		}
	}
}

func HandleLidarStream(addr string, port string, doneChan chan struct{}, wg *sync.WaitGroup, logLevel int) {
	defer wg.Done()

	for {
		select {
		case <-doneChan:
			return
		default:
			tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", addr, port))
			if err != nil {
				log(ERROR, logLevel, "Error resolving TCP address:", err)
				time.Sleep(time.Second * 2)
				continue
			}

			conn, err := net.DialTCP("tcp", nil, tcpAddr)
			if err != nil {
				log(ERROR, logLevel, fmt.Sprintf("Error connecting to LiDAR server at %s:%s: %v", addr, port, err))
				time.Sleep(time.Second * 2)
				continue
			}

			handleConnection(conn, doneChan, logLevel)

			conn.Close()
			time.Sleep(time.Second * 2)
		}
	}
}

func handleConnection(conn *net.TCPConn, doneChan chan struct{}, logLevel int) {
	conn.SetNoDelay(true)

	lidar := rplidar.NewRPLidar("/dev/rplidar", 115200, time.Second*3)
	err := lidar.Connect()
	if err != nil {
		log(ERROR, logLevel, "Error connecting to RPLidar:", err)
		return
	}
	defer lidar.Disconnect()

	info, err := lidar.GetInfo()
	if err != nil {
		log(ERROR, logLevel, "Error getting info:", err)
		return
	}
	log(INFO, logLevel, fmt.Sprintf("RPLidar Info: %+v", info))

	time.Sleep(time.Second * 1)

	measurements, err := lidar.IterMeasurements()
	if err != nil {
		log(ERROR, logLevel, "Error starting measurements:", err)
		return
	}

	for {
		select {
		case measurement, ok := <-measurements:
			if !ok {
				return
			}
			msg := fmt.Sprintf("{\"Angle\": %f, \"Distance\": %f}\n", measurement.Angle, measurement.Distance)
			_, err := conn.Write([]byte(msg))
			if err != nil {
				log(ERROR, logLevel, "Error sending data to server:", err)
				return
			}
		case <-doneChan:
			return
		}
	}
}

func HandleBatteryStream(addr string, port string, doneChan chan struct{}, wg *sync.WaitGroup, rm *rosmasterlib.Rosmaster, logLevel int) {
	defer wg.Done()

	for {
		select {
		case <-doneChan:
			return
		default:
			tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", addr, port))
			if err != nil {
				log(ERROR, logLevel, "Error resolving TCP address:", err)
				time.Sleep(time.Second * 2)
				continue
			}

			conn, err := net.DialTCP("tcp", nil, tcpAddr)
			if err != nil {
				log(ERROR, logLevel, fmt.Sprintf("Error connecting to Battery server at %s:%s: %v", addr, port, err))
				time.Sleep(time.Second * 2)
				continue
			}
			conn.SetNoDelay(true)

			handleBatteryConnection(conn, doneChan, rm, logLevel)

			conn.Close()
			time.Sleep(time.Second * 2)
		}
	}
}

func handleBatteryConnection(conn *net.TCPConn, doneChan chan struct{}, rm *rosmasterlib.Rosmaster, logLevel int) {
	for {
		select {
		case <-time.After(1 * time.Second):
			voltage := rm.GetBatteryVoltage()
			log(INFO, logLevel, fmt.Sprintf("Aktuelle Batteriespannung: %v V", voltage))
			_, err := conn.Write([]byte(fmt.Sprintf("%f\n", voltage)))
			if err != nil {
				log(ERROR, logLevel, "Error sending data to server:", err)
				return
			}
		case <-doneChan:
			return
		}
	}
}
