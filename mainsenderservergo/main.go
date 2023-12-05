package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

// Structs for JSON objects
type Motor struct {
	Motor1 int8 `json:"motor1" mapstructure:"motor1"`
	Motor2 int8 `json:"motor2" mapstructure:"motor2"`
	Motor3 int8 `json:"motor3" mapstructure:"motor3"`
	Motor4 int8 `json:"motor4" mapstructure:"motor4"`
}

type Lightbar struct {
	Mode   bool   `json:"mode" mapstructure:"mode"`
	LedID  string `json:"ledid" mapstructure:"ledid"`
	RGB    string `json:"rgb" mapstructure:"rgb"`
	Effect string `json:"effect" mapstructure:"effect"`
	Speed  string `json:"speed" mapstructure:"speed"`
	Parm   string `json:"parm" mapstructure:"parm"`
}

type Buzzer struct {
	Duration int `json:"buzzer" mapstructure:"buzzer"`
}

var targetConnection *net.Conn

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter target IP address: ")
	ip, _ := reader.ReadString('\n')
	ip = strings.TrimSpace(ip)

	fmt.Print("Enter port for streams: ")
	portStr, _ := reader.ReadString('\n')
	portStr = strings.TrimSpace(portStr)

	addr := ip + ":6969"
	currentConnection, err := net.Dial("tcp", addr)
	if err != nil {
		logWithTimestamp("Error connecting to server:", err)
		return
	}
	defer currentConnection.Close()
	logWithTimestamp("Connected to", addr)

	_, err = currentConnection.Write([]byte("startstreams " + portStr + "\n"))
	if err != nil {
		logWithTimestamp("Error sending startstreams command:", err)
		return
	}

	targetConnection = &currentConnection
	// Start a new TCP server to handle JSON objects
	go startJSONServer()
	go runHealthcheck()

	select {}
}

func runHealthcheck() {
	for {
		msg := "healthcheck " + strconv.FormatInt(time.Now().Unix(), 10) + "\n"
		_, err := (*targetConnection).Write([]byte(msg))
		if err != nil {
			logWithTimestamp("Error sending healthcheck message:", err)
			return
		}
		logWithTimestamp("Sent healthcheck message:", strings.Split(msg, "\n")[0])
		time.Sleep(2 * time.Second)
	}
}

func startJSONServer() {
	ln, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		logWithTimestamp("Error setting up JSON TCP server:", err)
		return
	}
	defer ln.Close()
	logWithTimestamp("JSON TCP server listening at localhost:8081")

	for {
		conn, err := ln.Accept()
		if err != nil {
			logWithTimestamp("Error accepting connection on JSON server:", err)
			continue
		}

		go handleJSONConnection(conn)
	}
}

func handleJSONConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for {
		conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		if scanner.Scan() {
			text := scanner.Text()
			var jsonData map[string]interface{}
			if err := json.Unmarshal([]byte(text), &jsonData); err != nil {
				logWithTimestamp("Invalid JSON:", err)
				continue
			}

			if _, ok := jsonData["motor1"]; ok {
				var motor Motor
				if err := mapstructure.Decode(jsonData, &motor); err != nil {
					logWithTimestamp("Error decoding motor data:", err)
					continue
				}
				handleMotorData(motor)
			} else if _, ok := jsonData["mode"]; ok {
				var lightbar Lightbar
				if err := mapstructure.Decode(jsonData, &lightbar); err != nil {
					logWithTimestamp("Error decoding lightbar data:", err)
					continue
				}
				handleLightbarData(lightbar)
			} else if _, ok := jsonData["buzzer"]; ok {
				var buzzer Buzzer
				if err := mapstructure.Decode(jsonData, &buzzer); err != nil {
					logWithTimestamp("Error decoding buzzer data:", err)
					continue
				}
				handleBuzzerData(buzzer)
			} else {
				logWithTimestamp("Unrecognized JSON object")
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
}

func handleMotorData(motor Motor) {
	logWithTimestamp("Received motor data:", motor)
	// Process motor data
	byteSlice, err := json.Marshal(motor)
	if err != nil {
		logWithTimestamp("Error marshalling motor data:", err)
		return
	}
	msg := "motor " + string(byteSlice) + "\n"

	(*targetConnection).Write([]byte(msg))
}

func handleLightbarData(lightbar Lightbar) {
	logWithTimestamp("Received lightbar data:", lightbar)
	// Process lightbar data
	byteSlice, err := json.Marshal(lightbar)
	if err != nil {
		logWithTimestamp("Error marshalling lightbar data:", err)
		return
	}
	msg := "lightbar " + string(byteSlice) + "\n"

	(*targetConnection).Write([]byte(msg))
}

func handleBuzzerData(buzzer Buzzer) {
	logWithTimestamp("Received buzzer data:", buzzer)
	// Process buzzer data
	byteSlice, err := json.Marshal(buzzer)
	if err != nil {
		logWithTimestamp("Error marshalling buzzer data:", err)
		return
	}
	msg := "buzzer " + string(byteSlice) + "\n"

	(*targetConnection).Write([]byte(msg))
}

func logWithTimestamp(v ...interface{}) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), v)
}
