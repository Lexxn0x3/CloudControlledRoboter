---
title: Lidar
layout: default
parent: MainsenderGo
grand_parent: MicroServices
nav_order: 1
---

# RPLidar Package Documentation for `mainsenderGo`

## Overview

The `rplidar` package, a component of the `mainsenderGo` application, is responsible for interfacing with an RPLidar device. It includes functionality for connecting to, controlling, and retrieving data from the lidar sensor.

## Structure and Functions

### Constants

- `SyncByte`, `SyncByte2`, etc.: Constants for various byte codes used in lidar communication.

### Types

#### `RPLidar`

- Represents an RPLidar device with serial port communication.
- Fields include `serialPort`, `portName`, `baudrate`, and `timeout`.

#### `Info`

- Contains information about the RPLidar device, such as model, firmware, hardware, and serial number.

#### `Measurement`

- Represents a single lidar measurement, including quality, angle, and distance.

### Functions

#### `NewRPLidar`

- Creates and returns a new `RPLidar` instance.
- Parameters: `portName`, `baudrate`, `timeout`.

#### `Connect` and `Disconnect`

- `Connect`: Establishes a serial connection to the RPLidar device.
- `Disconnect`: Closes the serial connection.

#### `sendCommand` and `readResponse`

- Internal functions to send commands to and read responses from the lidar device.

#### `GetInfo`

- Retrieves and returns the `Info` structure of the lidar device.

#### `StartScan` and `StopScan`

- `StartScan`: Initiates a scanning operation, returning the expected response size.
- `StopScan`: Ends the scanning operation.

#### `StartMotor` and `StopMotor`

- Control functions to start and stop the lidar motor.

#### `IterMeasurements`

- Initiates a scanning operation and returns a channel through which `Measurement` instances are sent.
- Parses raw scan data and provides it in a consumable format.

### Usage

1. **Initialization**: Create an `RPLidar` instance with appropriate serial port settings.
2. **Connection**: Call `Connect` to establish a connection to the RPLidar device.
3. **Operations**: Use functions like `GetInfo`, `StartScan`, `IterMeasurements`, etc., to control the lidar and retrieve data.
4. **Disconnection**: Ensure to call `Disconnect` to close the serial connection when done.

## Notes

- Proper error handling is crucial, especially in serial communication and data parsing.
- The package is designed to interact seamlessly with the RPLidar device, providing a straightforward API for higher-level functions in `mainsenderGo`.
- Ensure that the serial port and baud rate settings are correctly configured for your specific RPLidar model.
- For proper usage look into code for streamhandlers or mainsendergo
