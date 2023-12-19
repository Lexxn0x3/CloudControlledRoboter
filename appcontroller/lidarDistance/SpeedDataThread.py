import threading
import time
import json
from lidarDistance.utils import connect_to_server

class SpeedThread(threading.Thread):
    def __init__(self, ip_address, speed_server_port):
        super(SpeedThread, self).__init__()
        self.ip_address = ip_address
        self.speed_server_port = speed_server_port
    def run(self):
        global stop_threads, speed
        try:
            speed_sock = connect_to_server(self.ip_address, self.speed_server_port)
            while not stop_threads:
                data = speed_sock.recv(1024)
            if data:
                # Split the received data by newline and process each line as a separate JSON object
                for line in data.splitlines():
                    try:
                        # Decode the JSON data and extract the distance and angle
                        data = json.loads(line.decode('utf-8'))
                        speed = data['Speed']
                        time.sleep(1)
                            
                    except json.JSONDecodeError:
                        continue
            else:
                pass
        except KeyboardInterrupt:
            print("Stopping")