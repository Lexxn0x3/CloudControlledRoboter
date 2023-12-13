---
title: Tips and Tricks
layout: default
has_children: false
nav_order: 10
---

## Network and Communication Optimization
### Disable Nagle's Algorithm for TCP Packet Sending
#### Problem
The aggregation of TCP packets can lead to increased latency, parsing and buffer issues, impacting the speed and reliability of data transmission.

#### Solution
Disabling Nagle's Algorithm can result in faster TCP packet sending. This approach helps in reducing aggregation-related latency, making it beneficial for handling merged packets and alleviating buffer and delay issues.

### Wi-Fi Network
#### Problem
Stability and performance issues can arise when the robot is broadcasting its own Wi-Fi or is connected to a slow network.

#### Solution
Switching the robot's own Wi-Fi off and connecting to a different, more stable, and faster network can enhance its overall performance and reliability in data transmission and operational stability.

## Software and Development Tools
### Utilizing ChatGPT for Programming Assistance
ChatGPT with GPT-4 is effective in understanding, refactoring, and translating code. It's suitable for small functions but less efficient for larger segments. ChatGPT is great for brainstorming and initial code drafts, but the generated code should be reviewed and modified as needed.

### General Programming Advice
#### Multi-threading and CPU Usage
When using multiple threads, avoid infinite loops without any idle state. These can consume maximum CPU resources. Instead, consider using asynchronous calls or, at the very least, add a sleep timer within loops to manage CPU usage effectively. A good default sleep duration could be a few milliseconds, just enough to prevent CPU maxing out, such as `time.sleep(0.01)` in Python.

#### Cross-Platform and Cross-Architecture Compilation
- **Go and Rust:** Both languages can be easily compiled for Linux and Windows, and for both ARM and x86 architectures, offering great flexibility for cross-platform development.
- **Python:** As an interpreted language, Python works across different platforms without the need for specific compilation.
- **C#:** C# has become more versatile with .NET Core (now .NET 5 and onwards), supporting development for Windows, Linux, macOS, and other operating systems, including ARM architectures. This makes C# a viable option for applications that need to run on a variety of hardware, including ARM-based systems.

#### Python Specifics and Caveats
- **Performance:** Python can be slower, especially for computation-heavy tasks due to its interpreted nature. This is most noticeable in scenarios with extensive loops and data processing.
- **Type Checking:** Python uses dynamic typing, meaning types are checked at runtime. Although flexible, it can lead to runtime errors. Optional type hints available in Python 3.5+ offer a form of static type checking, but are not enforced by the runtime.
- **Runtime Crashes:** Due to dynamic typing, Python programs are more susceptible to runtime crashes, particularly with unexpected data types. Proper error handling and testing are essential to mitigate this risk.
- **Scope of Loop Variables:** Variables declared in loops remain accessible outside of the loop, unlike in some other languages where they are confined to the loop. This can lead to unanticipated behaviors if not carefully managed.
- **Mutable Default Parameters:** Using mutable objects (like lists or dictionaries) as default parameters in functions can lead to shared objects across function calls, resulting in unintended side effects. It's advisable to use `None` as a default and initialize the mutable object within the function.

## Robot Hardware and Peripheral Management
### Battery Voltage
A significant drop in the battery voltage can lead to partial robot malfunction, where some functionalities remain operational while others, like driving mechanisms, fail.

### USB Devices 
When USB devices, such as cameras, are passed to a ROS Docker container, they become inaccessible from the main operating system. This means that applications or commands outside the container, like `mainsendergo` or `ffmpeg`, cannot access these devices. It's important to plan device usage accordingly and set up necessary bridges or sharing mechanisms if simultaneous access is required.

### Limitations of Voice Module
The voice module comes with preprogrammed commands that are not easily modifiable. This rigidity can significantly limit its usability in dynamic or custom application scenarios, where flexibility in command customization is essential.

### OLED Screen 
The OLED screen of the robot consumes approximately 10% of the CPU resources, primarily due to its continual redrawing of each line. To optimize CPU usage, it is recommended to terminate the OLED screen process once a connection with the robot has been established. The `htop` command can be used to identify and manage this process. Alternatively, modifying the script that initiates the OLED screen process is also possible.

## User Interface and Interaction
### Command Line, VNC, and UI Limitations
It's generally more efficient to use the robot's command line for operations, as the VNC (Virtual Network Computing) connection tends to be slow and sluggish. Additionally, the user interface (UI), even when accessed through built-in HDMI or DisplayPort (DP), can be notably slow. Relying on command-line interactions can enhance responsiveness and efficiency in controlling and configuring the robot. 

Terminal software such as Termius is highly beneficial for managing the robot's operations. It allows users to have multiple terminal windows open simultaneously in an organized arrangement. This setup makes it easier to oversee and handle various tasks efficiently, enhancing the overall workflow and productivity when interacting with the robot.

## LiDAR System
- **Detection of Close Objects:** The LiDAR system struggles to detect objects that are very close to it. This limitation requires software-level handling to ensure close-proximity objects are appropriately accounted for in the robot's navigation and sensing algorithms.
- **Obstructed Rear View:** The LiDAR's ability to perceive the area behind the vehicle is obstructed by the sound module. This necessitates software compensation to mitigate potential blind spots in the robot's environmental awareness.
- **Latency in Data Acquisition:** The LiDAR system exhibits noticeable latency in data acquisition. The time it takes to complete a full rotation and update all points can be significant. This latency must be factored into the robot's movement and decision-making algorithms, especially considering the distance the robot can travel during one LiDAR cycle.

## Setting Up a Linux VM for Testing
For testing in a Linux environment, there are two recommended options:

1. **Using Windows Subsystem for Linux (WSL):** This allows you to run a Linux environment directly on Windows, without the overhead of a traditional virtual machine.

2. **Creating a New VM with Ubuntu Server Image:** Use the Ubuntu Server image, which doesn’t include a UI, for a lightweight setup. This is ideal for testing in a more controlled environment.

### Sharing Project Folders in VMware
In VMware Player and Workstation, you can easily share your project files between the host and the Linux VM. To do this, add a shared folder in the VM settings. Then, in the VM, execute the following command to mount the shared folder: `sudo vmhgfs-fuse .host:/ /mnt/hgfs/ -o allow_other -o uid=1000`
The mounted shared folder will be accessible in the `/mnt/hgfs` directory of your Linux VM, allowing for seamless file sharing and collaboration between your host and the VM.