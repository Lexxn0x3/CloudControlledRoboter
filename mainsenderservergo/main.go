package main

import (
	"bufio"
	"encoding/json"
	"flag"
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
	Mode   bool `json:"mode" mapstructure:"mode"`
	LedID  byte `json:"ledid" mapstructure:"ledid"`
	R      byte `json:"red" mapstructure:"red"`
	G      byte `json:"green" mapstructure:"green"`
	B      byte `json:"blue" mapstructure:"blue"`
	Effect byte `json:"effect" mapstructure:"effect"`
	Speed  byte `json:"speed" mapstructure:"speed"`
	Parm   byte `json:"parm" mapstructure:"parm"`
}

type Buzzer struct {
	Duration int `json:"buzzer" mapstructure:"buzzer"`
}

var targetConnection *net.TCPConn

var healthCheckTicker *time.Ticker

func main() {
	targetPointer := flag.String("target", "false", "ip adress of the target robot")
	listenPortPointer := flag.String("listenport", "4200", "port to listen for JSON objects")
	targetPortPointer := flag.String("targetport", "6969", "port to connect to")
	portPointer := flag.String("streamport", "false", "start port of streams")
	flag.Parse()

	healthCheckTicker = time.NewTicker(50 * time.Millisecond)
	reader := bufio.NewReader(os.Stdin)

	targetIP := strings.TrimSpace(*targetPointer)
	if targetIP == "false" {
		fmt.Print("Enter target IP address: ")
		targetIp, _ := reader.ReadString('\n')
		targetIp = strings.TrimSpace(targetIp)
	}

	portStr := strings.TrimSpace(*portPointer)
	if portStr == "false" {
		fmt.Print("Enter port for streams: ")
		portStr, _ := reader.ReadString('\n')
		portStr = strings.TrimSpace(portStr)
	}

	targetPort := strings.TrimSpace(*targetPortPointer)

	addr := targetIP + ":" + targetPort
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		logWithTimestamp("Error resolving TCP address:", err)
		return
	}
	currentConnection, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		logWithTimestamp("Error connecting to server:", err)
		return
	}
	defer currentConnection.Close()
	currentConnection.SetNoDelay(true)
	logWithTimestamp("Connected to", addr)

	_, err = currentConnection.Write([]byte("startstreams " + portStr + "\n"))
	if err != nil {
		logWithTimestamp("Error sending startstreams command:", err)
		return
	}

	targetConnection = currentConnection
	// Start a new TCP server to handle JSON objects
	listenPort := strings.TrimSpace(*listenPortPointer)
	go startJSONServer(listenPort)
	go runHealthcheck()

	for {
		time.Sleep(100 * time.Millisecond)
	}
}

func runHealthcheck() {
	for {
		<-healthCheckTicker.C
		msg := "healthcheck " + strconv.FormatInt(time.Now().UnixNano(), 10) + "\n"
		_, err := (*targetConnection).Write([]byte(msg))
		if err != nil {
			logWithTimestamp("Error sending healthcheck message:", err)
			return
		}
		//logWithTimestamp("Sent healthcheck message:", strings.Split(msg, "\n")[0])
		//time.Sleep(50 * time.Millisecond)
	}
}

func startJSONServer(port string) {

	ln, err := net.Listen("tcp4", "0.0.0.0:"+port)
	if err != nil {
		logWithTimestamp("Error setting up JSON TCP server:", err)
		return
	}
	defer ln.Close()
	logWithTimestamp("JSON TCP server listening at 0.0.0.0:4200")

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
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		if scanner.Scan() {
			text := scanner.Text()
			var jsonData map[string]interface{}
			if err := json.Unmarshal([]byte(text), &jsonData); err != nil {
				logWithTimestamp("Error unmarshalling JSON: ", text)
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
		time.Sleep(10 * time.Millisecond)
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
