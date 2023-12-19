import threading
import time
import socket
import json
from lidarDistance.utils import connect_to_server

# LidarDataThread class to handle the connection to the LidarTCPPort and update the lidarDataBuffer
class LidarDataThread(threading.Thread):
    def __init__(self, globals, ip_address, server_port, lidar_data_buffer):
        super(LidarDataThread, self).__init__()
        self.ip_address = ip_address
        self.server_port = server_port
        self.lidar_data = lidar_data_buffer
        self.globals = globals

    def run(self):
        while not self.globals.stop_threads:
            try:
                sock = connect_to_server(self.ip_address, self.server_port)
                self.receive_lidar_data(sock, self.lidar_data)

            except socket.error as e:
                print(f"Socket error: {e}")
                print("Attempting to reconnect...")
                time.sleep(1)  #Add a delay before attempt to reconnect
            except KeyboardInterrupt:
                break
            finally:
                sock.close()

    def receive_lidar_data(self, sock, lidar_data_buffer):
        try:
            # Receive the data in small chunks and retransmit it
            while True:
                data = sock.recv(1024)
                if data:
                    # Split the received data by newline and process each line as a separate JSON object
                    for line in data.splitlines():
                        try:
                            # Decode the JSON data and extract the distance and angle
                            data = json.loads(line.decode('utf-8'))
                            angle = data['Angle']
                            distance = data['Distance']
                            
                            rounded_angle = round(angle)
                            if rounded_angle == 360:
                                rounded_angle = 0
                            
                            # Put the data into the buffer
                            lidar_data_buffer.append((rounded_angle, distance))

                            # Ensure the buffer size is limited
                            if len(lidar_data_buffer) > 150:
                                lidar_data_buffer.popleft()  # Remove the oldest element
                                
                        except json.JSONDecodeError:
                            continue
                else:
                    break
        except KeyboardInterrupt:
            pass
        finally:
            # Clean up the connection
            sock.close()