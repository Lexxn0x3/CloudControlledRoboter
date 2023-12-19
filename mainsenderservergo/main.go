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

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
)

var logLevel = INFO // Set the default log level

var targetConnection *net.TCPConn
var healthCheckTicker *time.Ticker

func main() {
	targetPointer := flag.String("target", "false", "IP address of the target robot")
	listenPortPointer := flag.String("listenport", "4200", "port to listen for JSON objects")
	targetPortPointer := flag.String("targetport", "6969", "port to connect to")
	portPointer := flag.String("streamport", "false", "start port of streams")
	flag.IntVar(&logLevel, "loglevel", INFO, "log level (0=DEBUG, 1=INFO, 2=WARNING, 3=ERROR)")
	flag.Parse()

	healthCheckTicker = time.NewTicker(50 * time.Millisecond)
	reader := bufio.NewReader(os.Stdin)

	targetIP := strings.TrimSpace(*targetPointer)
	if targetIP == "false" {
		fmt.Print("Enter target IP address: ")
		targetIP, _ = reader.ReadString('\n')
		targetIP = strings.TrimSpace(targetIP)
	}

	portStr := strings.TrimSpace(*portPointer)
	if portStr == "false" {
		fmt.Print("Enter port for streams: ")
		portStr, _ = reader.ReadString('\n')
		portStr = strings.TrimSpace(portStr)
	}

	targetPort := strings.TrimSpace(*targetPortPointer)
	connected := false
	for !connected {
		currentConnection, err := connectToServer(targetIP, targetPort)
		if err != nil {
			log(ERROR, "Error connecting to server:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		log(INFO, "Connected to", targetIP+":"+targetPort)

		_, err = currentConnection.Write([]byte("startstreams " + portStr + "\n"))
		if err != nil {
			log(ERROR, "Error sending startstreams command:", err)
			time.Sleep(2 * time.Second)
			continue
		}
		connected = true
		defer currentConnection.Close()
		targetConnection = currentConnection
	}

	go startJSONServer(listenPortPointer)
	go runHealthcheck(targetIP, targetPort, portStr)

	for {
		time.Sleep(100 * time.Millisecond)
	}
}

func connectToServer(targetIP string, targetPort string) (*net.TCPConn, error) {
	addr := targetIP + ":" + targetPort
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	conn.SetNoDelay(true)
	return conn, nil
}

func runHealthcheck(targetIP string, targetPort string, portStr string) {
	for {
		<-healthCheckTicker.C
		msg := "healthcheck " + strconv.FormatInt(time.Now().UnixNano(), 10) + "\n"
		_, err := (*targetConnection).Write([]byte(msg))
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			log(WARNING, "Connection lost. Attempting to reconnect...")
			newConn, err := connectToServer(targetIP, targetPort)
			if err != nil {
				log(ERROR, "Reconnection failed:", err)
				continue
			}
			targetConnection = newConn
			_, err = targetConnection.Write([]byte("startstreams " + portStr + "\n"))
			if err != nil {
				log(ERROR, "Error sending startstreams command:", err)
				continue
			}
			log(INFO, "Reconnected to", targetIP+":"+targetPort)
		}
	}
}

func startJSONServer(portPointer *string) {
	port := strings.TrimSpace(*portPointer)
	ln, err := net.Listen("tcp4", "0.0.0.0:"+port)
	if err != nil {
		log(ERROR, "Error setting up JSON TCP server:", err)
		return
	}
	defer ln.Close()
	log(INFO, "JSON TCP server listening at 0.0.0.0:"+port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log(WARNING, "Error accepting connection on JSON server:", err)
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
				log(ERROR, "Error unmarshalling JSON: ", text)
				log(ERROR, "Invalid JSON:", err)
				continue
			}

			if _, ok := jsonData["motor1"]; ok {
				var motor Motor
				if err := mapstructure.Decode(jsonData, &motor); err != nil {
					log(ERROR, "Error decoding motor data:", err)
					continue
				}
				handleMotorData(motor)
			} else if _, ok := jsonData["mode"]; ok {
				var lightbar Lightbar
				if err := mapstructure.Decode(jsonData, &lightbar); err != nil {
					log(ERROR, "Error decoding lightbar data:", err)
					continue
				}
				handleLightbarData(lightbar)
			} else if _, ok := jsonData["buzzer"]; ok {
				var buzzer Buzzer
				if err := mapstructure.Decode(jsonData, &buzzer); err != nil {
					log(ERROR, "Error decoding buzzer data:", err)
					continue
				}
				handleBuzzerData(buzzer)
			} else if _, ok := jsonData["laser"]; ok {
				var laser Laser
				if err := mapstructure.Decode(jsonData, &laser); err != nil {
					log(ERROR, "Error decoding laser data:", err)
					continue
				}
				handleLaserdata(laser)
			} else {
				log(WARNING, "Unrecognized JSON object")
			}
		}

		if err := scanner.Err(); err != nil {
			netErr, ok := err.(net.Error)
			if !ok || !netErr.Timeout() {
				log(ERROR, "Error reading from connection:", err)
				break
			}
			scanner = bufio.NewScanner(conn)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func log(level int, v ...interface{}) {
	if level < logLevel {
		return
	}
	levelStr := []string{"DEBUG", "INFO", "WARNING", "ERROR"}[level]
	prefix := fmt.Sprintf("%s [%s] ", time.Now().Format("2006-01-02 15:04:05"), levelStr)
	fmt.Println(prefix, v)
}
