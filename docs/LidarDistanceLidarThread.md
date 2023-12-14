---
title: Lidar Data Thread
layout: default
parent: LidarDistance
grand_parent: MicroServices
nav_order: 1
---

# Lidar Data Thread

The `LidarDataThread` module, contained in the `LidarDataThread.py` file, is a crucial component of the Lidar Distance system. This module is responsible for managing the connection to the Lidar TCP port, continuously receiving Lidar data, and updating the Lidar data buffer.

## Class: `LidarDataThread`

### Attributes

- **ip_address**: The IP address of the Lidar TCP port.
- **server_port**: The port number of the Lidar TCP port.
- **lidar_data**: The Lidar data buffer shared with other threads.

### Methods

#### `__init__(self, ip_address, server_port, lidar_data_buffer`

Initializes the `LidarDataThread` instance.

- **Parameters**:
  - `ip_address`: IP address for the Lidar TCP connection.
  - `server_port`: Port number for the Lidar TCP connection.
  - `lidar_data_buffer`: Buffer for storing Lidar data.

- **Usage**:
  ```python
  lidar_data_thread = LidarDataThread(ip_address, server_port, lidar_data_buffer)

#### `run(self)`

The main execution method of the thread. Manages the continuous reception of Lidar data from the specified TCP port through calling the `receive_lidar_data` function and checking and connecting to the Lidar TCP port if connection is lost.

#### `receive_lidar_data(self, sock, lidar_data_buffer)`

Receives Lidar data from the Lidar TCP port and updates the Lidar data buffer.

- **Parameters**:
  - `sock`: The socket connection to the Lidar TCP port
  - `lidar_data_buffer`: Buffer for storing Lidar data.

- **Usage**:
  This method is called internally during the thread execution

## Script Execution
The Lidar Data Thread can be executed as part of the Lidar-Distance system. Ensure that the necessary dependencies are installed and run the main Lidar-Distance script.

#### Part in the Code
```python
# Example usage in LidarDistanceController.py
lidar_data_buffer = pregenerate_lidar_data(maxLenBuffer)
lidar_data_thread = LidarDataThread(lidarstreamIP, lidarstreamPort, lidar_data_buffer)
lidar_data_thread.start()