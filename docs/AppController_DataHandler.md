---
title: DataHandler
layout: default
parent: App Controller
grand_parent: MicroServices
nav_order: 2
---


# DataHandler
The `DataHandler` module, contained in the `DataHandler.py` file, is a crucial component of the Robot Controller system. It defines the `DataHandler` class, responsible for processing incoming data and translating it into appropriate commands for the robot.

## Class: `DataHandler`

### Attributes

- **stopFront**: A flag indicating if there is an obstacle in front of the robot.
- **stopLeft**: A flag indicating if there is an obstacle on the left of the robot.
- **stopRight**: A flag indicating if there is an obstacle on the right of the robot.
- **stopFrontRight**: A flag indicating if there is an obstacle in the front-right of the robot.
- **stopFrontLeft**: A flag indicating if there is an obstacle in the front-left of the robot.
- **bot**: A reference to the main `Bot` instance for communication and control.
- **bp**: An instance of the `BetterPrinting` class for improved console printing.

### Methods

#### `__init__(self, bot_instance)`

Initializes the `DataHandler` instance.

- **Parameters**:
  - `bot_instance`: Reference to the main `Bot` instance.

#### `handle_stopFlag_data(self, data)`

Processes incoming stop flag data and takes appropriate actions.

- **Parameters**:
  - `data`: Dictionary containing stop flag data.

#### `handle_buzzer_data(self, data)`

Processes incoming buzzer data and sends corresponding commands to the robot.

- **Parameters**:
  - `data`: Dictionary containing buzzer data.

#### `handle_lightbar_data(self, data)`

Processes incoming lightbar data and sends corresponding commands to the robot.

- **Parameters**:
  - `data`: Dictionary containing lightbar data.

#### `handle_direction_data(self, data)`

Processes incoming direction data and translates it into motor control commands for the robot.

- **Parameters**:
  - `data`: Dictionary containing direction data.

#### Movement Control Methods

These methods process specific directional commands and send corresponding motor control commands to the robot.

- `drive_forward(speed)`
- `drive_right_forward(speed)`
- `drive_left_forward(speed)`
- `drive_backward(speed)`
- `drive_right_backward(speed)`
- `drive_left_backward(speed)`
- `drive_right(speed)`
- `drive_left(speed)`
- `spin_right(speed)`
- `spin_left(speed)`
- `drive_curve_right_forward(speed)`
- `drive_curve_left_forward(speed)`
- `drive_curve_right_backward(speed)`
- `drive_curve_left_backward(speed)`

  Each of these methods takes a `speed` parameter and sends the appropriate motor control commands based on the specific movement.


## Note

The `DataHandler` class is designed to work in conjunction with the `Bot` class and handles the interpretation of incoming data, ensuring appropriate actions are taken to control the robot's movements.