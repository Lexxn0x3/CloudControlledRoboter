package main

import (
	"flag"
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
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

// Global channels for different components.
var motorChan = make(chan string)
var lightbarChan = make(chan string)
var buzzerChan = make(chan string)
var laserChannel = make(chan string)
var healthCheckChan = make(chan string)

// Channel for stopping goroutines.
var stopChan chan struct{}

// GPIO pin interface.
var gpioPinOut gpio.PinIO

// ROS master library instance.
var rosmaster *rosmasterlib.Rosmaster

// Log levels.
const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
)

// Current log level.
var logLevel = INFO

func main() {
	// Initialize hardware-related features.
	if _, err := host.Init(); err != nil {
		log(DEBUG, "Error initializing periph host:", err)
		return
	}

	// Access a GPIO pin
	pin := gpioreg.ByName("GPIO14") // EXCUSE ME? GPIO14=GPIO13 https://jetsonhacks.com/nvidia-jetson-nano-j41-header-pinout/ <- Sysfs GPIO is this pin number with Pin being the pin on the board
	if pin == nil {
		return
	}

	if err := pin.Out(gpio.Low); err != nil {
		log(ERROR, err)
	}

	// Command line flag for setting the server port.
	listenPort := flag.String("port", "6969", "port to listen on")
	flag.IntVar(&logLevel, "loglevel", INFO, "log level (0=DEBUG, 1=INFO, 2=WARNING, 3=ERROR)")
	flag.Parse()

	// Channel for OS signals.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Setup TCP server.
	ln, err := net.Listen("tcp", "0.0.0.0:"+*listenPort)
	if err != nil {
		log(ERROR, "Error setting up TCP server:", err)
		return
	}
	defer ln.Close()
	log(INFO, "TCP server listening at 0.0.0.0:", *listenPort)

	// Handle incoming JSON data.
	go handleIncomingJson()

	// Goroutine to handle OS signals.
	go func() {
		<-sigChan
		log(INFO, "SIGINT received, shutting down.")
		ln.Close()
		_, ok := <-stopChan
		if ok {
			close(stopChan)
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	// Accepting and handling new connections.
	for {
		conn, err := ln.Accept()
		if err != nil {
			log(WARNING, "Error accepting connection:", err)
			return
		}
		log(INFO, "Connection accepted")

		go handleConnection(conn)
		time.Sleep(100 * time.Millisecond)
	}
}

// handleHealthcheck manages periodic health checks.
// This function runs as a goroutine to perform regular health checks.
// It ensures that the system is functioning correctly and takes appropriate
// actions (like shutting down) if a health check fails.
func handleHealthcheck(wg *sync.WaitGroup) {
	defer wg.Done()
	healthCheckActive := false
	lastHealthCheckActive := healthCheckActive

	var lastTimestamp int64
	timer := time.NewTimer(300 * time.Millisecond)
	log(INFO, "Health check routine started")

	for {
		select {
		case <-timer.C:
			if healthCheckActive {
				log(WARNING, "Health check failed. No message received in time.")
				close(stopChan)
				return
			}
		case msg := <-healthCheckChan:
			lastHealthCheckActive = healthCheckActive
			healthCheckActive = true
			if lastHealthCheckActive != healthCheckActive {
				timer.Reset(300 * time.Millisecond)
			}

			if strings.HasPrefix(msg, "healthcheck") {
				fields := strings.Fields(msg)
				if len(fields) < 2 {
					log(WARNING, "Invalid healthcheck message format")
					continue
				}
				timestamp, err := strconv.ParseInt(fields[1], 10, 64)
				if err != nil {
					log(WARNING, "Invalid timestamp in healthcheck message")
					continue
				}
				if timestamp <= lastTimestamp {
					log(WARNING, "Received an old timestamp in healthcheck message")
					continue
				}
				lastTimestamp = timestamp
				//logWithTimestamp("Health check passed")
				rosmaster.BlockedHealthcheck = false
				timer.Reset(300 * time.Millisecond)
			}
		case <-stopChan:
			log(INFO, "Health check routine stopped")
			return
		}
	}
}

// Helper function to beep three times.
func threeBeep() {
	rosmaster.SetBeep(100)
	time.Sleep(150 * time.Millisecond)
	rosmaster.SetBeep(100)
	time.Sleep(150 * time.Millisecond)
	rosmaster.SetBeep(100)
}

// Helper function to close all provided channels.
func closeAllChannels(chans ...chan struct{}) {
	for _, ch := range chans {
		if ch != nil {
			close(ch)
		}
	}
}

// Logging function with timestamp and log level.
func log(level int, v ...interface{}) {
	if level < logLevel {
		return
	}

	levelStr := []string{"DEBUG", "INFO", "WARNING", "ERROR"}[level]
	prefix := fmt.Sprintf("%s [%s] ", time.Now().Format("2006-01-02 15:04:05"), levelStr)
	fmt.Println(prefix, v)
}
