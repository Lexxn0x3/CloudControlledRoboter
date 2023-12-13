---
title: Tips and Tricks
layout: default
has_children: false
nav_order: 10
---

## Disable Nagle's Algorithm for TCP Packet Sending
### Problem
The aggregation of TCP packets can lead to increased latency, parsing and buffer issues, impacting the speed and reliability of data transmission.

### Solution
Disabling Nagle's Algorithm can result in faster TCP packet sending. This approach helps in reducing aggregation-related latency, making it beneficial for handling merged packets and alleviating buffer and delay issues.

## Battery Voltage
A significant drop in the battery voltage can lead to partial robot malfunction, where some functionalities remain operational while others, like driving mechanisms, fail.

## Wi-Fi Network
### Problem
Stability and performance issues can arise when the robot is broadcasting its own Wi-Fi or is connected to a slow network.

### Solution
Switching the robot's own Wi-Fi off and connecting to a different, more stable, and faster network can enhance its overall performance and reliability in data transmission and operational stability.

## MainSenderGo, MainSenderServerGo, Broker
Start the mainsendergo on the robot, and initiate the mainsenderservergo and broker on an external server. Alternatively, for testing purposes, these can also be run on the robot itself, though this might defy the intended use. Once set up, the app can connect to this configuration for optimal operation.

## Utilizing ChatGPT for Programming Assistance
ChatGPT, particularly with GPT-4, is highly proficient in understanding, refactoring, and translating code from one language to another. It is very effective for working with smaller functions. However, for larger code segments, it may become less efficient and sluggish. ChatGPT is an excellent tool for brainstorming ideas and providing initial code drafts. Nonetheless, it is advisable to thoroughly review, understand, and modify the generated code to ensure it meets the specific requirements and standards of the project.

## USB Device Accessibility in ROS Docker Container
When USB devices, such as cameras, are passed to a ROS Docker container, they become inaccessible from the main operating system. This means that applications or commands outside the container, like `mainsendergo` or `ffmpeg`, cannot access these devices. It's important to plan device usage accordingly and set up necessary bridges or sharing mechanisms if simultaneous access is required.
