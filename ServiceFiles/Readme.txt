Put .sh files into path: /opt/roboterteam1/downloadservice.sh
For Broker: Modify instance count in .sh


Move .service file to /etc/systemd/system/YOURNAME.service

sudo systemctl daemon-reload

FOR MAINSENDER:
sudo systemctl enable YOURNAME.service
this now always starts the service on system startup

FOR BROKER:
sudo systemctl enable YOURNAME.service@PORT
this now always starts the service on system startup with port and port + 6000
for the multiple instances the next ports are port + 10 and port + 10 + 6000