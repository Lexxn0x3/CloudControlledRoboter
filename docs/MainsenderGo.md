---
title: MainsenderGo
layout: default
parent: MicroServices
has_children: true
nav_order: 4
---

# Documentation for TCP Server Application in Go

## Overview

This Go program is a crucial component of our infrastructure, being the only program that runs directly on the robot. It implements a TCP server that listens for incoming connections on a specified port and handles JSON messages for controlling motors, lightbars, and a buzzer. Additionally, it manages health checks to monitor the connection's status.

## Structure and Functionality

### Main Components

- **Motor**: Represents a motor with its control values.
- **Lightbar**: Represents a lightbar with various attributes like mode, color, effect, and speed.
- **Buzzer**: Represents a buzzer with a duration attribute.

### Global Variables

- `motorChan`, `lightbarChan`, `buzzerChan`: Channels for passing motor, lightbar, and buzzer data.
- `rosmaster`: Instance of `rosmasterlib.Rosmaster`, used to set motor and lightbar states and control the buzzer.

### Main Function

1. **TCP Server Setup**: Initializes a TCP server listening on a user-defined port.
2. **Signal Handling**: Captures OS signals for graceful shutdown.
3. **JSON Handler**: Runs a goroutine to handle incoming JSON messages for motor, lightbar, and buzzer.
4. **Connection Handling**: Accepts incoming connections and processes data sent by clients.

### JSON Message Handling (`handleIncomingJson`)

Processes incoming JSON messages from the channels and updates the state of motors, lightbars, and the buzzer accordingly.

### Health Check Routine (`handleHealthcheck`)

Monitors the health of the connection and triggers a shutdown if no health check message is received within a specified time frame.

### Connection Management (`handleConnection`)

1. **Command Processing**: Interprets and acts upon commands received from the client, such as starting or stopping streams.
2. **Stream Handling**: Manages streams for camera, lidar, and battery data.
3. **Health Check Integration**: Incorporates the health check routine to monitor connection integrity.

### Utility Functions

- `threeBeep()`: Triggers the buzzer to beep three times.
- `closeAllChannels(chans ...chan struct{})`: Closes all provided channels.
- `logWithTimestamp(v ...interface{})`: Logs messages with a timestamp.

## Usage

To run the server:

1. **Compile and Run**: Compile the Go code and run the resulting binary. Optionally, specify the port to listen on using the `-port` flag.
2. **Client Connection**: Connect a client to the server's IP and port.
3. **Data Transmission**: The client can send JSON-formatted messages to control motors, lightbars, and buzzers, and send commands to start/stop data streams.

## Compiling for Different Architectures and Operating Systems

To compile the Go program for a different architecture or operating system, use the `GOOS` and `GOARCH` environment variables. For example:

- To compile for Windows on an AMD64 architecture:
```bash
GOOS=windows GOARCH=amd64 go build
```

- To compile for Linux on an ARM64 architecture, which is what we need for the robot:
```bash
GOOS=linux GOARCH=arm64 go build
```

Replace the values of `GOOS` and `GOARCH` as needed to target the desired platform.

## Notes

- Ensure that the `rosmasterlib` and `streamhandlers` packages are correctly imported and configured.
- Proper error handling and logging are implemented for robustness.
- The program is designed to run indefinitely until an interrupt signal is received.