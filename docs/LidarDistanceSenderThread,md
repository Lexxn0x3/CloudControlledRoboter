---
title: DistanceSenderThread
layout: default
parent: LidarDistance
grand_parent: MicroServices
nav_order: 3
---

# Distance Sender Thread

The `DistanceSenderThread` module, contained in the `DistanceSenderThread.py` file, is a critical component of the Lidar Distance system. This module handles the connection back to the broker, continuously analyzes Lidar data, and provides a stream of distance data for the five directions left, front-left, front, front-right and right to the broker.

## Class: `DistanceSenderThread`

### Attributes

- **lidar_data_buffer**: The Lidar data buffer shared with the Lidar Data Thread.
- **ip_address**: The IP address of the broker.
- **server_port**: The port number of the broker.

### Methods

#### `__init__(self, ip_address, server_port, lidar_data_buffer)`

Initializes the `DistanceSenderThread` instance.

- **Parameters:**
  - `ip_address`: IP address for the broker connection.
  - `server_port`: Port number for the broker connection.
  - `lidar_data_buffer`: Buffer for storing Lidar data.

- **Usage:**
  ```python
  sender_thread = DistanceSenderThread(ip_address, server_port, lidar_data_buffer)
  ```

#### `run(self)`

The main execution method of the thread. Manages the continuous analysis of Lidar data and sending distance information back to the broker.

#### `send_Distance_json(self, sock, distances)`

Sends a JSON package containing distance information back to the broker.

- **Parameters:**
  - `sock`: The socket connection to the Lidar TCP port
  - `distances`: Tuple containing left, front, right, front-left, and front-right distances.

#### `check_Distance(self, lidar_data)`

Analyzes Lidar data and checks distances in various directions.

- **Parameters:**
  - `lidar_data`: List of Lidar data containing angle-distance pairs.

- **Returns:**
  - Tuple containing left, front, right, front-left, and front-right distances.

## Script Execution
The Distance Sender Thread can be executed as part of the Lidar Distance Controller system. Ensure that the necessary dependencies are installed and run the main controller script.

#### Part in the Code
```python
sender_thread = DistanceSenderThread(distanceSenderIP, distanceSenderPort, lidar_data_buffer)
sender_thread.start()
```