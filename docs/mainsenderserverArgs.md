---
title: Arguments
layout: default
parent: MainsenderServerGo
grand_parent: MicroServices
nav_order: 1
---

# Command-Line Arguments Documentation for the Client Program

## Overview

This document explains the command-line arguments available in the client program for connecting to the `mainsendergo` server.

## Arguments

1. **target** (`-target`): The IP address of the target robot. If not provided, the program will prompt for it.
   
   Usage: `-target="192.168.1.10"`

2. **listenport** (`-listenport`): The port on which the client will listen for incoming JSON objects. Default is `4200`.
   
   Usage: `-listenport="4200"`

3. **targetport** (`-targetport`): The port to connect to on the target robot. Default is `6969`.
   
   Usage: `-targetport="6969"`

4. **streamport** (`-streamport`): The starting port for streams. If not provided, the program will prompt for it.
   
   Usage: `-streamport="5000"`

## Usage

- The client program can be started with a combination of these arguments to configure its connection and behavior.
- Example command to start the client:
  
```bash
go run . -target="192.168.1.10" -listenport="4200" -targetport="6969" -streamport="5000"
```


## Notes

- Omitting an argument will either use its default value or prompt the user to enter it during runtime.
- Properly setting these arguments is crucial for the correct operation of the client program.

