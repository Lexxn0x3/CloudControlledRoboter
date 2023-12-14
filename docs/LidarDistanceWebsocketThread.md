---
title: WebsocketThread
layout: default
parent: LidarDistance
grand_parent: MicroServices
nav_order: 4
---

# WebSocket Client Thread

The `WebSocketClientThread` module, contained in the `WebSocketClientThread.py` file, is a crucial component of the Lidar Distance system. This module manages a WebSocket client connection, continuously monitors changes in stop flags, and sends relevant data to the app handler.

## Class: `WebSocketClientThread`

### Attributes

- **ip_address**: The IP address of the app handler.
- **app_handler_port**: The port number of the app handler.

### Methods

#### `__init__(self, ip_address, app_handler_port)`

Initializes the `WebSocketClientThread` instance.

- **Parameters**:
  - `ip_address`: IP address for the app handler connection.
  - `app_handler_port`: Port number for the app handler connection.

- **Usage**:
  ```python
  websocket_client_thread = WebSocketClientThread(ip_address, app_handler_port)
  ```


#### `run(self)`

The main execution method of the thread. Monitors changes in stop flags and sends relevant data to the app handler.

#### `send_data_to_app_handler(self, ip_address, app_handler_port)`

Sends stop flag data as a json package to the app handler using a WebSocket connection.

- **Parameters:**
  - `ip_address`: IP address for the app handler connection.
  - `app_handler_port`: Port number for the app handler connection.

- **Usage:**
  - This method is called internally during the thread execution.

## Script Execution
The WebSocket Client Thread can be executed as part of the Lidar Distance Controller system. Ensure that the necessary dependencies are installed and run the main controller script.

#### Part in the Code
```python
websocket_client_thread = WebSocketClientThread(appHandlerIP, appHandlerPort)
websocket_client_thread.start()
```