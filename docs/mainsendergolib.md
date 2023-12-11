---
title: Rosmasterlib
layout: default
parent: MainsenderGo
grand_parent: MicroServices
nav_order: 1
---

# Rosmasterlib Package Documentation for `mainsenderGo`

## Overview

The `rosmasterlib` package is a vital component of the `mainsenderGo` application, designed to interface with the robot's hardware. It facilitates communication with the robot, enabling control over motors, lights, beeps, and retrieving sensor data like battery voltage and gyroscope readings.

## Structure and Functions

### Type `Rosmaster`

- Represents the main interface for robot communication and control.
- Fields include serial port, device IDs, function codes for different operations, and debugging flags.

### Constructor `NewRosmaster`

- Initializes a new `Rosmaster` instance.
- Connects to the robot via a serial port.
- Parameters: `comPort` (string), `baudRate` (int).

### Methods

#### `Close`

- Closes the serial port and stops internal processes.

#### `readSerial`

- Internal method to continuously read from the serial port.
- Parses incoming data and handles it based on predefined protocols.

#### `writeSerial`

- Internal method to write data to the serial port.
- Ensures synchronization with the robot's communication protocol.

#### `parseData`

- Parses incoming data packets and updates the internal state of `Rosmaster`.
- Handles different types of data like battery voltage and sensor readings.

#### `sum`

- Calculates the checksum for data packets.

#### `limitMotorValue`

- Ensures motor speed values are within valid ranges.

#### `SetMotor`

- Controls the speed of the robot's motors.
- Parameters: speed values for each motor.

#### `SetCarMotion`

- Sets the motion parameters for the robot.
- Parameters: motion values (float64) in XYZ axes.

#### `SetBeep`

- Controls the beep function of the robot.
- Parameter: duration of the beep.

#### `SetColorfulLamps` and `SetColorfulEffect`

- Controls the robot's RGB lamps and effects.
- Parameters: color values and effect settings.

#### `GetBatteryVoltage`, `GetGyroscope`, `GetAcceleration`, `GetMagnitude`

- Retrieves various sensor data from the robot.

### Usage

1. **Initialization**: Create a `Rosmaster` instance with the appropriate serial port settings.
2. **Control Operations**: Use methods like `SetMotor`, `SetBeep`, and `SetColorfulLamps` to control the robot's hardware.
3. **Data Retrieval**: Use methods like `GetBatteryVoltage` and `GetGyroscope` to retrieve sensor data from the robot.
4. **Shutdown**: Call `Close` to properly disconnect and clean up resources.

## Notes

- Ensuring the correct serial port configuration is crucial for effective communication with the robot.
- The package provides a comprehensive API for controlling and monitoring the robot, which is integral to the `mainsenderGo` application.
- Adequate error handling and debugging capabilities are implemented to facilitate smooth operation and maintenance.
- For proper usage look into code for streamhandlers or mainsendergo
