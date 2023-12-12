---
title: Websocket Controller
layout: default
parent: App Controller
grand_parent: MicroServices
nav_order: 4
---

# WebSocketController 

The `WebSocketController` module, contained in the `WebSocketController.py` file, is a key component of the Robot Controller system. It defines the `WebSocketController` class, responsible for managing WebSocket communication with the robot.

## Class: `WebSocketController`

### Attributes

- **ip**: The IP address for WebSocket communication.
- **port**: The port number for WebSocket communication.
- **bot**: A reference to the main `Bot` instance for communication and control.
- **dh**: An instance of the `DataHandler` class for processing incoming data.
- **loop**: The asyncio event loop for managing asynchronous tasks.
- **bp**: An instance of the `BetterPrinting` class for improved console printing.

### Methods

#### `__init__(self, ip, port, datahandler, bot_instance)`

Initializes the `WebSocketController` instance.

- **Parameters**:
  - `ip`: The IP address for WebSocket communication.
  - `port`: The port number for WebSocket communication.
  - `datahandler`: An instance of the `DataHandler` class.
  - `bot_instance`: Reference to the main `Bot` instance.

#### `handle_websocket(self, websocket, _)`

Handles incoming WebSocket messages, processes the data, and triggers corresponding actions.

- **Parameters**:
  - `websocket`: The WebSocket connection.
  - `_`: Placeholder for additional data (not used).

#### `start_websocket_server(self)`

Starts the WebSocket server and awaits incoming connections.

#### `start(self)`

Initiates the WebSocket server and enters the event loop to handle WebSocket connections.

## Note

The `WebSocketController` class works alongside the `DataHandler` and `Bot` classes, managing WebSocket communication. It processes incoming data and triggers appropriate actions within the Robot Controller system.