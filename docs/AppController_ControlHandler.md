---
title: AppControlHandler
layout: default
parent: App Controller
grand_parent: MicroServices
nav_order: 1
---

# AppControlHandler

The `AppControlHandler` module, contained in the `AppControlHandler.py` file, is a critical component of the Robot Controller system. It defines the `Bot` class, responsible for managing connections, processing incoming commands, and controlling the robot's movements. The class is designed to be instantiated with parameters related to TCP and WebSocket connections, as well as print preferences.

## Class: `Bot`

### Attributes

- **i_print**: A boolean flag indicating whether informational prints are enabled (default is `True`).
- **d_print**: A boolean flag indicating whether debugging prints are enabled (default is `True`).
- **e_print**: A boolean flag indicating whether error prints are enabled (default is `True`).
- **tcpc**: An instance of the `TCPController` class for managing TCP communication with the robot.
- **dh**: An instance of the `DataHandler` class for processing incoming data.
- **websocket_controller**: An instance of the `WebSocketController` class for managing WebSocket communication.

### Methods

#### `__init__(self, tcp_ip, tcp_port, web_ip, web_port, info_print=True, debug_print=True, error_print=True)`

Initializes the `Bot` instance.

- **Parameters**:
  - `info_print`: Flag for enabling/disabling informational prints (default is `True`).
  - `debug_print`: Flag for enabling/disabling debugging prints (default is `True`).
  - `error_print`: Flag for enabling/disabling error prints (default is `True`).
  - `tcp_ip`: IP address for TCP communication.
  - `tcp_port`: Port number for TCP communication.
  - `web_ip`: IP address for WebSocket communication.
  - `web_port`: Port number for WebSocket communication.

- **Usage**:
  ```python
  my_bot = Bot(tcp_ip, tcp_port, web_ip, web_port, info_print, debug_print, error_print)
  ```

#### `startTCPConection(self, ip, port)`

Starts the TCP connection.

- **Parameters**:
  - `ip`: IP address for TCP connection.
  - `port`: Port number for TCP connection.


#### `startDataHandler(self)`

Starts the `DataHandler` module.


#### `startWebsocket(self, web_ip, web_port)`

Starts the WebSocket connection.

- **Parameters**:
  - `web_ip`: IP address for WebSocket connection.
  - `web_port`: Port number for WebSocket connection.



#### `stop(self)`

Stops all robot movements.


### Movement Control Methods

These methods provide various movement commands for the robot.

#### `send_motor_data(self, motor1, motor2, motor3, motor4)`

Sends motor control data to the robot.

- **Parameters**:
  - `motor1`: Speed for motor 1.
  - `motor2`: Speed for motor 2.
  - `motor3`: Speed for motor 3.
  - `motor4`: Speed for motor 4.


#### `send_lightbar_data(self, isEffect, red, green, blue, effect, speed)`

Sends lightbar control data to the robot.

- **Parameters**:
  - `isEffect`: Flag indicating if it's an effect.
  - `red`, `green`, `blue`: RGB color values.
  - `effect`: Effect type.
  - `speed`: Speed of the lightbar.


#### `send_buzzer_data(self, onBuzzer)`

Sends buzzer control data to the robot.

- **Parameters**:
  - `onBuzzer`: Flag indicating if the buzzer should be turned on.


## Script Execution

The script can be executed from the command line with the following arguments:

```bash
python appControlHandler.py tcp_ip tcp_port web_ip web_port [info_print] [debug_print] [error_print]
```

- `tcp_ip`: IP address for TCP communication.
- `tcp_port`: Port number for TCP communication.
- `web_ip`: IP address for WebSocket communication.
- `web_port`: Port number for WebSocket communication.
- `[info_print]`, `[debug_print]`, `[error_print]`: Optional flags for enabling/disabling prints (default is `True`).

Example:

```bash
python appControlHandler.py 192.168.8.20 4200 192.168.8.105 5000 True True True 
