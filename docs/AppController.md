---
title: App Controller
layout: default
parent: MicroServices
has_children: true
nav_order: 2
---

# App Controller

## Introduction

This document provides documentation for the Robot Controller, a Python-based system that translates incoming WebSocket commands into robot-readable motor data, subsequently sent using TCP communication. The system is modularized into three main components: `appControlHandler`, `betterPrinting`, `tcpController`, `dataHandler` and `websocketController`.

### Components

1. **appControlHandler (`appControlHandler.py`):**
   - This module initializes and manages the main control logic of the robot.
   - It establishes connections, handles motor control, and manages various functionalities such as driving, spinning, and lightbar control.

2. **betterPrinting (`betterPrinting.py`):**
   - A utility module providing enhanced printing with colored output for information, debugging, and error messages.
   - It is used to improve the readability of log messages in the console.

3. **tcpController (`tcpController.py`):**
   - Manages the TCP connection to the robot.
   - Sends JSON data to the robot for motor control, lightbar control, and buzzer control.
   - Handles connection errors and supports automatic reconnection.

4. **dataHandler (`dataHandler.py`):**
   - Processes incoming data from various sources (stop flags, buzzer commands, lightbar commands, and direction commands).
   - Translates incoming data into corresponding robot control commands.
   - Utilizes the `BetterPrinting` class for improved console output.

5. **websocketController (`websocketController.py`):**
   - Manages a WebSocket server for real-time communication with clients.
   - Handles incoming data from clients and delegates the processing to the `dataHandler`.
   - Utilizes asynchronous features for concurrent handling of WebSocket connections.

## Usage

### Prerequisites

- Python 3.x

### Installation

1. Clone the repository:

   ```bash
   git clone [repository_url]
   ```

2. Install dependencies:

   ```bash
   pip install -r requirements.txt
   ```

### Running the System

1. Run the main control handler script:

   ```bash
   python3 appControlHandler.py [tcp_ip] [tcp_port] [web_ip] [web_port] [info_print] [debug_print] [error_print]
   ```

   Example:

   ```bash
   python3 appControlHandler.py 192.168.1.1 5000 192.168.1.1 8080 True True True
   ```

   Replace the arguments with the appropriate values.

