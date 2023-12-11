---
title: MainsenderServerGo
layout: default
parent: MicroServices
has_children: true
nav_order: 5
---
# Client Program Documentation for Connecting to `mainsendergo`

## Overview

This Go program serves as a client that establishes a connection to the `mainsendergo` server, which runs on the robot. This is called mainsenderservergo because it runs on the server where the brokers are also running. It is designed to send control commands in the form of JSON objects for motors, lightbars, and buzzers. Additionally, it includes a health check mechanism to maintain the connection's stability.

## Structure and Functions

### Types

#### `Motor`, `Lightbar`, `Buzzer`

- Structs that represent JSON objects for controlling motors, lightbars, and buzzers.

### Global Variables

- `targetConnection`: A `*net.TCPConn` object to manage the connection to the server.
- `healthCheckTicker`: A ticker for periodically sending health check messages.

### Main Function (`main`)

- Parses command-line arguments to set up connection parameters.
- Establishes a connection to the `mainsendergo` server.
- Initiates a TCP server to listen for incoming JSON objects and handle them.
- Runs a health check loop to ensure continuous connection.

### Helper Functions

#### `connectToServer`

- Establishes a TCP connection to the given IP and port.
- Returns a `*net.TCPConn` object representing the connection.

#### `runHealthcheck`

- Periodically sends health check messages to the server.
- Reconnects to the server if the connection is lost.

#### `startJSONServer`

- Starts a TCP server to listen for incoming JSON objects.
- For each incoming connection, it initiates `handleJSONConnection`.

#### `handleJSONConnection`

- Reads and processes JSON objects from the connection.
- Depending on the JSON type (`Motor`, `Lightbar`, `Buzzer`), it handles the data appropriately.

#### `handleMotorData`, `handleLightbarData`, `handleBuzzerData`

- Handle the respective data types and send control commands to the `mainsendergo` server.

### Usage

1. **Start the Client**: Run the client program with necessary flags (`target`, `listenport`, `targetport`, `streamport`).
2. **JSON Control**: Send JSON formatted control messages to the specified listen port.
3. **Server Interaction**: The client forwards these messages to the `mainsendergo` server running on the robot.
4. **Health Check**: The client maintains the connection with periodic health checks.

## Notes

- Ensure proper network configurations for successful connection establishment.
- The client plays a crucial role in remotely controlling the robot via the `mainsendergo` server.
- Adequate error handling and logging mechanisms are implemented for reliability and troubleshooting.
- For details on the command-line arguments, refer to the separate documentation page.
