---
title: Streamhandler
layout: default
parent: MainsenderGo
grand_parent: MicroServices
nav_order: 1
---

# Stream Handlers Documentation for `mainsenderGo`

## Overview

The `streamhandlers` package is part of the `mainsenderGo` application and is responsible for managing various streams - camera, LiDAR, and battery - on the robot. This package ensures continuous data flow and robust connection management for these specific streams.

## Structure and Functions

### Camera Stream Handling (`HandleCameraStream`)

- **Purpose**: Manages the camera stream using `ffmpeg` to send video data to a specified address and port.
- **Operations**:
  - Initiates `ffmpeg` with the right parameters for video streaming.
  - Monitors the process and restarts `ffmpeg` if it exits unexpectedly.
  - Stops `ffmpeg` on receiving a done signal through `doneChan`.

### LiDAR Stream Handling (`HandleLidarStream`)

- **Purpose**: Handles the LiDAR stream by establishing and maintaining a TCP connection for data transmission.
- **Operations**:
  - Attempts to continuously establish a TCP connection to the target address and port.
  - Manages LiDAR data transmission via `handleConnection`.
  - Retries connection after disconnection.

#### `handleConnection`

- Manages LiDAR data transmission over TCP.
- Connects to RPLidar, retrieves measurements, and sends them via TCP.
- Terminates on channel closure or a done signal.

### Battery Stream Handling (`HandleBatteryStream`)

- **Purpose**: Manages battery stream to send battery voltage data over TCP.
- **Operations**:
  - Similar to LiDAR handling, establishes a TCP connection for data transmission.
  - Manages battery data via `handleBatteryConnection`.

#### `handleBatteryConnection`

- Regularly transmits battery voltage (every second).
- Retrieves voltage from `rosmasterlib.Rosmaster`.
- Sends voltage data over TCP and terminates on a done signal.

## Usage

Integrate `streamhandlers` with `mainsenderGo`:

1. **Integrate**: Include the `streamhandlers` package in `mainsenderGo`.
2. **Initialize**: Use the handlers with necessary parameters (address, port, channels, wait group, and `rosmasterlib.Rosmaster` for the battery stream).
3. **Manage Streams**: Incorporate stream management into the robot's operational flow.

## Notes

- Ensure robust error handling for each stream.
- Design facilitates continuous connection attempts for reliability.
- These streams are an integral part of `mainsenderGo`, contributing to the overall functionality of the robot.
- For proper usage have a look at the mainsendergo