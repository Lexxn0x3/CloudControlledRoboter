---
title: LidarDistance
layout: default
parent: MicroServices
has_children: true
nav_order: 6
---

# Lidar Distance Microservice

## Introduction

This document provides documentation for the Lidar Distance Controller, a Python-based program designed to handle Lidar data streaming, distance sending, and WebSocket communication. The program consists of three main threads: `LidarDataThread`, `DistanceSenderThread`, and `WebSocketClientThread`. Additionally, utility functions and global variables are imported from `utils` and `globals` modules, respectively.

### Components

1. **Lidar Distance Main (`LidarDistanceMain.py`):**
   - This module initializes and manages the main control logic of the robot.
   - It establishes connections, handles motor control, and manages various functionalities such as driving, spinning, and lightbar control.

2. **LidarDataThread (`LidarDataThread.py`):**
   - Manages the Lidar data streaming from a specified IP and port.
   - Utilizes a Lidar data buffer to store and provide data to other threads.

3. **DistanceSenderThread (`DistanceSenderThread.py`):**
   - Uses the Lidar data buffer for accessing Lidar data
   - Checks distances in five directions: left, front-left, front, front-right and right and sets the associated stopping flag if the distance is shorter than the minimum distance
   - Sends a stream of distance data for the five directions over a TCP connection to a specified IP and port.

4. **WebSocketClientThread (`WebSocketClientThread.py`):**
   - Establishes a WebSocket client connection to a specified IP and port.
   - Sends a json package with the necessary information for stopping the robot to app control handler if one of the stop-flags, set by the DistanceSenderThread, is changing

5. **utils (`utils.py`):**
   - Provides utility functions, including `pregenerate_lidar_data` for generating initial Lidar data, `connect_to_server` for connecting to a tcp server and `signal_handler` for handling Ctrl+C signals.

6. **globals (`globals.py`):**
    - Contains global variables, including `stop_threads` to control thread termination and stop flags, e.g. `stop_front`, to stop the robot if to close to an obsticle.

## Usage

### Prerequisites

- Python 3.x

### Installation

1. Clone the repository:

   ```bash
   git clone [repository_url]
   ```


### Running the System

1. Run the LidarDistanceMain script:

   ```bash
   python3 LidarDistanceController.py [lidarstream_ip] [lidarstream_port] [appHandler_ip] [appHandler_port] [distanceSender_ip] [distanceSender_port]
   ```

   Example:

   ```bash
   python3 LidarDistanceController.py 192.168.8.20 9011 192.168.8.20 6942 192.168.8.20 3031
   ```

   Replace the arguments with the appropriate values.
2. To stop the program, use Ctrl+C. The program will gracefully terminate, stopping all threads.