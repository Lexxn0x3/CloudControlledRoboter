title: Camera
layout: default
parent: TCP Broker
nav_order: 1
---
# CMD Arguments

## Introduction
This document outlines the steps required to start the Rust program and describes the command-line arguments it accepts.

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

## Command-Line Arguments
The program accepts the following arguments:
- `--config <PATH>`: Specify the path to the configuration file.
- `--verbose`: Enable verbose output for debugging purposes.
- `--help`: Display help information about the command-line arguments.

For detailed usage and more options, run:
```bash
./my_program --help
```
