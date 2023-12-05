package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/mattzi/mainsendergo/rosmasterlib"
	"github.com/mattzi/mainsendergo/streamhandlers"
	"github.com/mitchellh/mapstructure"
)

type Motor struct {
	Motor1 int8 `json:"motor1" mapstructure:"motor1"`
	Motor2 int8 `json:"motor2" mapstructure:"motor2"`
	Motor3 int8 `json:"motor3" mapstructure:"motor3"`
	Motor4 int8 `json:"motor4" mapstructure:"motor4"`
}

var motorChan = make(chan string)

type Lightbar struct {
	Mode   bool   `json:"mode" mapstructure:"mode"`
	LedID  string `json:"ledid" mapstructure:"ledid"`
	RGB    string `json:"rgb" mapstructure:"rgb"`
	Effect string `json:"effect" mapstructure:"effect"`
	Speed  string `json:"speed" mapstructure:"speed"`
	Parm   string `json:"parm" mapstructure:"parm"`
}

var lightbarChan = make(chan string)

type Buzzer struct {
	Duration int `json:"buzzer" mapstructure:"buzzer"`
}

var buzzerChan = make(chan string)

var rosmaster *rosmasterlib.Rosmaster

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ln, err := net.Listen("tcp", "0.0.0.0:6969")
	if err != nil {
		logWithTimestamp("Error setting up TCP server:", err)
		return
	}
	defer ln.Close()
	logWithTimestamp("TCP server listening at 0.0.0.0:6969")

	go handleIncomingJson()

	go func() {
		<-sigChan
		logWithTimestamp("SIGINT received, shutting down.")
		ln.Close()
		_, ok := <-stopChan
		if ok {
			close(stopChan)
			time.Sleep(500 * time.Millisecond)
		}
		return
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			logWithTimestamp("Error accepting connection:", err)
			return
		}
		logWithTimestamp("Connection accepted")

		go handleConnection(conn)
	}
}

func handleIncomingJson() {
	for {
		select {
		case msg := <-motorChan:
			var jsonData map[string]interface{}
			unmarshalMsg := strings.Replace(msg, "motor ", "", 1)
			if err := json.Unmarshal([]byte(unmarshalMsg), &jsonData); err != nil {
				logWithTimestamp("Invalid JSON:", err)
				continue
			}

			var motor Motor
			err := mapstructure.Decode(jsonData, &motor)
			if err != nil {
				logWithTimestamp("Error decoding JSON:", err)
				continue
			}
			logWithTimestamp("Received motor:", motor)

			rosmaster.SetMotor(motor.Motor1, motor.Motor2, motor.Motor3, motor.Motor4)
		case msg := <-lightbarChan:
			var jsonData map[string]interface{}
			unmarshalMsg := strings.Replace(msg, "lightbar ", "", 1)
			if err := json.Unmarshal([]byte(unmarshalMsg), &jsonData); err != nil {
				logWithTimestamp("Invalid JSON:", err)
				continue
			}
			var lightbar Lightbar
			err := mapstructure.Decode(jsonData, &lightbar)
			if err != nil {
				logWithTimestamp("Error decoding JSON:", err)
				continue
			}
			logWithTimestamp("Received lightbar:", lightbar)
		case msg := <-buzzerChan:
			var jsonData map[string]interface{}
			unmarshalMsg := strings.Replace(msg, "buzzer ", "", 1)
			if err := json.Unmarshal([]byte(unmarshalMsg), &jsonData); err != nil {
				logWithTimestamp("Invalid JSON:", err)
				continue
			}
			var buzzer Buzzer
			err := mapstructure.Decode(jsonData, &buzzer)
			if err != nil {
				logWithTimestamp("Error decoding JSON:", err)
				continue
			}
			logWithTimestamp("Received buzzer:", buzzer)
		}
	}
}

var stopChan chan struct{}
var healthCheckChan = make(chan string)

// Health check routine
func handleHealthcheck(wg *sync.WaitGroup) {
	defer wg.Done()
	healthCheckActive := false
	lastHealthCheckActive := healthCheckActive

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
				rosmaster.BlockedHealthcheck = false
				timer.Reset(3 * time.Second)
			}
		case <-stopChan:
			logWithTimestamp("Health check routine stopped")
			return
		}
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var wg sync.WaitGroup
	cmdChan := make(chan string)
	doneReadingChan := make(chan bool)
	stopChan = make(chan struct{})

	go func() {
		scanner := bufio.NewScanner(conn)
		for {
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			if scanner.Scan() {
				text := scanner.Text()
				if strings.HasPrefix(text, "healthcheck") {
					healthCheckChan <- text
				} else if strings.HasPrefix(text, "buzzer") {
					buzzerChan <- text
				} else if strings.HasPrefix(text, "lightbar") {
					lightbarChan <- text
				} else if strings.HasPrefix(text, "motor") {
					motorChan <- text
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

				wg.Add(4)

				rosmaster = rosmasterlib.NewRosmaster("/dev/myserial", 115200)
				defer rosmaster.Close()
				rosmaster.SetBeep(100)
				rosmaster.SetColorfulLamps(0xFF, 0, 0, 0)
				rosmaster.SetColorfulEffect(0, 255, 255)
				rosmaster.SetColorfulEffect(6, 255, 255)

				go streamhandlers.HandleCameraStream(conn.RemoteAddr().(*net.TCPAddr).IP.String(), port, cameraDoneChan, &wg)
				go streamhandlers.HandleLidarStream(conn.RemoteAddr().(*net.TCPAddr).IP.String(), lidarPortStr, lidarDoneChan, &wg)
				go streamhandlers.HandleBatteryStream(conn.RemoteAddr().(*net.TCPAddr).IP.String(), batteryPortStr, batteryDoneChan, &wg, rosmaster)
				go handleHealthcheck(&wg)
				healthCheckChan <- "healthcheck 0"
			} else if strings.HasPrefix(strings.ToLower(cmd), "stopstreams") {
				logWithTimestamp("Received stopstreams command")
				closeAllChannels(cameraDoneChan, lidarDoneChan, batteryDoneChan)

				wg.Wait()
				logWithTimestamp("All streams stopped")
				cameraDoneChan, lidarDoneChan, batteryDoneChan = nil, nil, nil
			}

		case <-doneReadingChan:
			rosmaster.SetMotor(0, 0, 0, 0)
			rosmaster.SetColorfulLamps(0xFF, 0, 0, 0)
			rosmaster.BlockedHealthcheck = true
			threeBeep()
			logWithTimestamp("Connection closed by client.")
			closeAllChannels(cameraDoneChan, lidarDoneChan, batteryDoneChan)
			return
		case <-stopChan:
			rosmaster.SetMotor(0, 0, 0, 0)
			rosmaster.SetColorfulLamps(0xFF, 0, 0, 0)
			rosmaster.BlockedHealthcheck = true
			threeBeep()
			logWithTimestamp("Connection closed due to stop chan receive.")
			closeAllChannels(cameraDoneChan, lidarDoneChan, batteryDoneChan)
			return
		}

	}
}
func threeBeep() {
	rosmaster.SetBeep(100)
	time.Sleep(150 * time.Millisecond)
	rosmaster.SetBeep(100)
	time.Sleep(150 * time.Millisecond)
	rosmaster.SetBeep(100)
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
