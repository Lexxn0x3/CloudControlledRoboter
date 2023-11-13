---
title: VoiceModul
layout: default
parent: Testing
nav_order: 3
---
# Information to Voice Module
## Setup
The voice module can be set up using [this Documentation](http://www.yahboom.net/study/ROSMASTER-X3#!).

First you have to run this command
```yaml
sudo gedit /etc/udev/rules.d/usb.rules
```
before binding comment out the line that includes "myseria" with a "#", otherwise the voice board and the ROS extension board will be identified as the same device. It needs to be saved and quit.
![Alt text](image-2.png)
then run the following commands
```yaml
sudo udevadm trigger
sudo service udev reload
sudo service udev restart
```
then the ROS expansion board needs to be binded

you have to run this command
```yaml
ll /dev/ttyUSB*
```
it will show all the USB devices
if it shows two you can follow the documentation linked above
if it shows three you have to plug out the white cable of the voice module, run the "ll /dev/ttyUSB*" command again.
Now you know which one the voice module was, you can plug the cable in.
Now you have two devices left to identify.

when you run this command 
```yaml
ll /dev/rplidar
```
you will see which of those is the lidar.
That leaves the last one to be the ROS extension board.

Then run this command with the name of the extension board instead of "ttyUSB1"
```yaml
udevadm info --attribute-walk --name=/dev/ttyUSB1 |grep devpath
```
you will see all the devpaths. The first one is important

![Alt text](image-3.png)

Now run this command 
```yaml
sudo gedit /etc/udev/rules.d/myserial.rules
```
add the content shown below and modify the "ATTRS{devpath}=="1.4.3"" to your first devpath
```yaml
KERNEL=="ttyUSB*",ATTRS{devpath}=="1.4.3",ATTRS{idVendor}=="1a86",ATTRS{idProduct}=="7523",MODE:="0777",SYMLINK+="myserial"
```
then save and exit and reload with the following commands
```yaml
sudo udevadm trigger
sudo service udev reload
sudo service udev restart
```