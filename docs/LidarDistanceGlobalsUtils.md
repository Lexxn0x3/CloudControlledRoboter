---
title: Globals and Utils
layout: default
parent: LidarDistance
grand_parent: MicroServices
nav_order: 5
---

# Globals and Utils

The `globals` and `utils` modules, containing shared global variables and utility functions, play a crucial role in the Lidar Distance system.

## Global Variables (`globals.py`)

### Variables

- **stop_front**: Flag indicating whether there is an obstacle in the front.
- **stop_right**: Flag indicating whether there is an obstacle on the right.
- **stop_left**: Flag indicating whether there is an obstacle on the left.
- **stop_front_right**: Flag indicating whether there is an obstacle in the front-right.
- **stop_front_left**: Flag indicating whether there is an obstacle in the front-left.
- **stop_threads**: Flag indicating whether all threads should stop.
- **minDist**: The minimum distance considered for obstacle detection.
- **maxLenBuffer**: The maximum length of the Lidar data buffer.

### Usage

These global variables are used across different threads to manage the system's state and obstacle detection.

## Utility Functions (`utils.py`)

### Functions

#### `connect_to_server(ip_address, server_port)`

Creates a TCP/IP socket and connects to the server.

- **Parameters**:
  - `ip_address`: IP address of the server.
  - `server_port`: Port number of the server.

- **Returns**:
  A connected socket.

#### `pregenerate_lidar_data(maxlen)`

Generates a Lidar data buffer with a specified maximum length.

- **Parameters**:
  - `maxlen`: Maximum length of the Lidar data buffer.

- **Returns**:
  A deque object representing the Lidar data buffer.

#### `signal_handler(sig, frame)`

Handles the Ctrl+C signal and exits the program gracefully.

- **Parameters**:
  - `sig`: Signal number.
  - `frame`: Current stack frame.

- **Usage**:
  This function is used as a signal handler for Ctrl+C interruptions.

### Script Execution

Globals and utils are imported and utilized in various modules and threads across the Lidar Distance system.

```python
from globals import stop_front, stop_right, stop_left, stop_front_right, stop_front_left, stop_threads, minDist, maxLenBuffer
from utils import connect_to_server, pregenerate_lidar_data, signal_handler
```
