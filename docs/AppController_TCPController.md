---
title: TCP Controller
layout: default
parent: App Controller
grand_parent: MicroServices
nav_order: 3
---

# TCPController

The `TCPController` module, contained in the `TCPController.py` file, is a crucial component of the Robot Controller system. It defines the `TCPController` class, responsible for managing TCP communication with the robot.

## Class: `TCPController`

### Attributes

- **robot_host**: The IP address of the robot for TCP communication.
- **robot_port**: The port number for TCP communication.
- **robot_socket**: The TCP socket for communication with the robot.
- **auto_reconnect**: A boolean flag indicating whether automatic reconnection is enabled (default is `True`).
- **bot**: A reference to the main `Bot` instance for communication and control.
- **bp**: An instance of the `BetterPrinting` class for improved console printing.

### Methods

#### `__init__(self, robot_host, robot_port, bot_instance)`

Initializes the `TCPController` instance.

- **Parameters**:
  - `robot_host`: The IP address of the robot for TCP communication.
  - `robot_port`: The port number for TCP communication.
  - `bot_instance`: Reference to the main `Bot` instance.

#### `connect(self)`

Establishes a connection to the robot via TCP.

#### `send_json_data(self, json_data)`

Sends JSON-formatted data to the robot over the established TCP connection.

- **Parameters**:
  - `json_data`: The JSON-formatted data to be sent.

#### `close_connection(self)`

Closes the TCP connection with the robot.


## Note

The `TCPController` class is designed to work in conjunction with the `Bot` class and manages the TCP communication with the robot. It provides methods for connecting, sending data, and closing the connection.