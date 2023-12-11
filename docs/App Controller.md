---
title: TCP Broker
layout: default
parent: MicroServices
has_children: true
nav_order: 1
---

# TcpBroker

The TcpBroker is a pivotal component within our microservice architecture that acts as a messaging hub for various services, including robotic operations. It utilizes TCP for reliable inter-service communication.

## Information flow

From Single connection to Multi Connection

## Architecture Overview

TcpBroker facilitates the exchange of messages and data streams between services such as pathfinding and depth analysis, and robotic systems. It is designed to manage and route data streams effectively, ensuring that all components within the architecture communicate synchronously.

### Connecting Robots and Services

Robots and other services connect to TcpBroker through TCP connections, enabling them to send and receive data. This centralized communication model allows for scalability and reliability in data handling.

### Testing with Mock ffmpeg Stream

The functionality of TcpBroker can be tested using a mock ffmpeg stream. This is useful for simulating video data transmission which the broker can then disseminate to other services for processing. Use the following command to create a test stream:

```bash
ffmpeg -re -i file_example_MP4_1920_18MG.mp4 -c:v mjpeg -q:v 5 -f mjpeg tcp://0.0.0.0:12345
```

### Python Program for Stream Decoding

VideoDecode.py is a Python program that demonstrates how to connect to TcpBroker and decode the transmitted video stream into images.

### Test setup

1. Start the Broker with the IP and ports you wish
2. Start the VideoDecode.py with the client port and ip of the Broker (it will wait with decoding till it actually has a completed frame)
3. Start the ffmpeg stream so it connects to the Broker on the receiving end (with the port)

Now the ffmpeg tool should send a TCP MJPEG stream to the Broker which in return broadcasts it to all its clients. 

The VideoDecode.py will accumilate one frame and then save it as a file in the same Directory.
