---
title: Services
layout: default
nav_order: 8
---

# Overview over Services

Put .sh files into path: /opt/roboterteam1/downloadservice.sh 


#### FOR MAINSENDER:
Move .service file to /etc/systemd/system/YOURNAME.service

```yaml
sudo systemctl daemon-reload 
```

```yaml
sudo systemctl start YOURNAME.service
```

```yaml
sudo systemctl enable YOURNAME.service
```
This now always starts the service on system startup

#### FOR BROKER:
Modify instance count in .sh

Move .service file to /etc/systemd/system/YOURNAME@.service

```yaml
sudo systemctl daemon-reload 
```

```yaml
sudo systemctl start YOURNAME@
```

```yaml
sudo systemctl enable YOURNAME@
```

This now always starts the service on system startup with port and port + 6000:
For the multiple instances the next ports are port + 10 and port + 10 + 6000