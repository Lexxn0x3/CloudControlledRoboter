package main

import (
	"bufio"
	"encoding/json"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mattzi/mainsendergo/rosmasterlib"
	"github.com/mattzi/mainsendergo/streamhandlers"
	"github.com/mitchellh/mapstructure"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var wg sync.WaitGroup
	cmdChan := make(chan string)
	doneReadingChan := make(chan bool)
	stopChan = make(chan struct{})

	go func() {
		scanner := bufio.NewScanner(conn)
		for {
			conn.SetReadDeadline(time.Now().Add(10 * time.Second))
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
				} else if strings.HasPrefix(text, "laser") {
					laserChannel <- text
				} else {
					cmdChan <- text
				}
			}

			if err := scanner.Err(); err != nil {
				netErr, ok := err.(net.Error)
				if !ok || !netErr.Timeout() {
					log(WARNING, "Error reading from connection:", err)
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
			log(DEBUG, "Received command:", cmd)
			if strings.HasPrefix(strings.ToLower(cmd), "startstreams") {
				port := strings.TrimSpace(cmd[len("startstreams"):])
				lidarPort, err := strconv.Atoi(port)
				if err != nil {
					log(WARNING, "Invalid port number for streams:", err)
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
				//rosmaster.SetColorfulEffect(6, 255, 255)

				go streamhandlers.HandleCameraStream(conn.RemoteAddr().(*net.TCPAddr).IP.String(), port, cameraDoneChan, &wg, logLevel)
				go streamhandlers.HandleLidarStream(conn.RemoteAddr().(*net.TCPAddr).IP.String(), lidarPortStr, lidarDoneChan, &wg, logLevel)
				go streamhandlers.HandleBatteryStream(conn.RemoteAddr().(*net.TCPAddr).IP.String(), batteryPortStr, batteryDoneChan, &wg, rosmaster, logLevel)
				go handleHealthcheck(&wg)
				healthCheckChan <- "healthcheck 0"
			} else if strings.HasPrefix(strings.ToLower(cmd), "stopstreams") {
				log(INFO, "Received stopstreams command")
				closeAllChannels(stopChan)

				wg.Wait()
				log(INFO, "All streams stopped")
				cameraDoneChan, lidarDoneChan, batteryDoneChan = nil, nil, nil
			}

		case <-doneReadingChan:
			rosmaster.SetMotor(0, 0, 0, 0)
			rosmaster.SetColorfulLamps(0xFF, 0, 0, 0)
			rosmaster.BlockedHealthcheck = true
			rosmaster.SetColorfulEffect(0, 255, 255)
			threeBeep()
			log(INFO, "Connection closed by client.")
			closeAllChannels(cameraDoneChan, lidarDoneChan, batteryDoneChan)
			return
		case <-stopChan:
			rosmaster.SetMotor(0, 0, 0, 0)
			rosmaster.SetColorfulLamps(0xFF, 0, 0, 0)
			rosmaster.BlockedHealthcheck = true
			rosmaster.SetColorfulEffect(0, 255, 255)
			threeBeep()
			log(INFO, "Connection closed due to stop chan receive.")
			closeAllChannels(cameraDoneChan, lidarDoneChan, batteryDoneChan)
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func handleIncomingJson() {
	for {
		select {
		case msg := <-motorChan:
			var jsonData map[string]interface{}
			unmarshalMsg := strings.Replace(msg, "motor ", "", 1)
			if err := json.Unmarshal([]byte(unmarshalMsg), &jsonData); err != nil {
				log(WARNING, "Invalid JSON:", err)
				continue
			}

			var motor Motor
			err := mapstructure.Decode(jsonData, &motor)
			if err != nil {
				log(ERROR, "Error decoding JSON:", err)
				continue
			}
			log(DEBUG, "Received motor:", motor)

			rosmaster.SetMotor(motor.Motor1, motor.Motor2, motor.Motor3, motor.Motor4)
		case msg := <-lightbarChan:
			var jsonData map[string]interface{}
			unmarshalMsg := strings.Replace(msg, "lightbar ", "", 1)
			if err := json.Unmarshal([]byte(unmarshalMsg), &jsonData); err != nil {
				log(WARNING, "Invalid JSON:", err)
				continue
			}
			var lightbar Lightbar
			err := mapstructure.Decode(jsonData, &lightbar)
			if err != nil {
				log(ERROR, "Error decoding JSON:", err)
				continue
			}
			log(DEBUG, "Received lightbar:", lightbar)
			if lightbar.Mode {
				rosmaster.SetColorfulEffect(lightbar.Effect, lightbar.Speed, lightbar.Parm)
			} else {
				rosmaster.SetColorfulLamps(lightbar.LedID, lightbar.R, lightbar.G, lightbar.B)
			}
		case msg := <-buzzerChan:
			var jsonData map[string]interface{}
			unmarshalMsg := strings.Replace(msg, "buzzer ", "", 1)
			if err := json.Unmarshal([]byte(unmarshalMsg), &jsonData); err != nil {
				log(WARNING, "Invalid JSON:", err)
				continue
			}
			var buzzer Buzzer
			err := mapstructure.Decode(jsonData, &buzzer)
			if err != nil {
				log(ERROR, "Error decoding JSON:", err)
				continue
			}
			log(DEBUG, "Received buzzer:", buzzer.Duration)

			rosmaster.SetBeep(buzzer.Duration)
		case msg := <-laserChannel:
			var jsonData map[string]interface{}
			unmarshalMsg := strings.Replace(msg, "laser ", "", 1)
			if err := json.Unmarshal([]byte(unmarshalMsg), &jsonData); err != nil {
				log(WARNING, "Invalid JSON:", err)
				continue
			}
			var laser Laser
			err := mapstructure.Decode(jsonData, &laser)
			if err != nil {
				log(ERROR, "Error decoding JSON:", err)
				continue
			}
			// Access a GPIO pin
			pin := gpioreg.ByName("GPIO14") // EXCUSE ME? GPIO14=GPIO13 https://jetsonhacks.com/nvidia-jetson-nano-j41-header-pinout/ <- Sysfs GPIO is this pin number with Pin being the pin on the board
			if pin == nil {
				return
			}

			if laser.Status {
				// Set the pin as output (for example)
				if err := pin.Out(gpio.High); err != nil {
					log(ERROR, "Error setting GPIO pin:", err)
				}
			} else {
				// Set the pin as output (for example)
				if err := pin.Out(gpio.Low); err != nil {
					log(ERROR, err)
				}
			}
		}
	}
}
