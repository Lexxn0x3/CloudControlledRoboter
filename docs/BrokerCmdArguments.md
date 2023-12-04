---
title: Arguments
layout: default
parent: TCP Broker
grand_parent: MicroServices
nav_order: 1
---

# CMD Arguments

## Introduction
This document provides the necessary information to start the Rust program and details the command-line arguments it accepts. The program serves as a TCP server/client handler with UI for statistics and debugging.

## Information flow

From Single connection to Multi Connection

## Requirements
- Rust programming language (latest stable version recommended)
- Cargo (Rust's package manager and build system)

## Installation
Clone the repository and build the project using Cargo:
```bash
git clone https://github.com/your-repository/your-project.git
cd your-project
cargo build --release
```

## Starting the Program
To start the program, navigate to the target/release directory and run:
```bash
./my_program
```

Alternatively, you can run the program with Cargo, passing any arguments after `--` to ensure they are not interpreted by Cargo itself:
```bash
cargo run -- --single-connection-port 3001 --multi-connection-port 4001 --debug-level info --buffer-size 4096
```

## Command-Line Arguments
The program accepts the following arguments:
- `--single-connection-port <PORT>`: Sets the server port (default: 3001).
- `--multi-connection-port <PORT>`: Sets the client port (default: 4001).
- `--debug-level <LEVEL>`: Sets the debug level. Possible values are debug, info, error (default: info).
- `--buffer-size <SIZE>`: Sets the buffer size in bytes (default: 512000).
- `--websocket-port <SIZE>`: Sets the buffer size in bytes (default: 2097152) 2Mb.
- `--websocket-port <PORT>`: Sets the websocket port (default: 5001).

For example, to start the program with a custom server port and buffer size, you would run:
```bash
./my_program --single-connection-port 3020 --buffer-size 8192
```

Or with Cargo:
```bash
cargo run -- --single-connection-port 3020 --buffer-size 8192
```

Replace `<PORT>`, `<LEVEL>`, and `<SIZE>` with your desired configurations.
