import threading
import time
import socket
import json
from lidarDistance.utils import connect_to_server

# DistanceSenderThread class to handle the connection back to the broker and analyze and provide a stream of distance data
class DistanceSenderThread(threading.Thread):
    def __init__(self, globals, ip_address, server_port, lidar_data_buffer):
        super(DistanceSenderThread, self).__init__()
        self.lidar_data_buffer = lidar_data_buffer
        self.ip_address = ip_address
        self.server_port = server_port
        self.globals = globals
    
    def run(self):
        while not self.globals.stop_threads:
            sender_sock = None
            try:
                sender_sock = connect_to_server(self.ip_address, self.server_port)
                last_distances = (10000, 10000, 10000, 10000, 10000)
                while True:
                    lidar_data = list(self.lidar_data_buffer)
                    distances = self.check_Distance(lidar_data)
                    if distances != last_distances:
                        self.send_Distance_json(sender_sock, distances)

                    distances = last_distances
                    time.sleep(0.05)

            except socket.error as e:
                print(f"Socket error: {e}")
                print("Attempting to reconnect...")
                time.sleep(1) #Add a delay before attempting to reconnect
            except KeyboardInterrupt:
                print("Stopping")
                break
            finally:
                if sender_sock is not None:
                    sender_sock.close()

    #Sends the json package back to broker
    def send_Distance_json(self, sock, distances):
        distance_data = {
            "left_distance": distances[0],
            "front_distance": distances[1], 
            "right_distance": distances[2],
            "front_left_distance": distances[3],
            "front_right_disctance": distances[4]
        }
        json_data = json.dumps(distance_data)
        sock.sendall(json_data.encode())
    
    #Function to check distance of LidarData
    def check_Distance(self, lidar_data):
            
        # Initialize variables outside the function if they are not already initialized
        left_D = left_D if 'left_D' in locals() else None
        front_D = front_D if 'front_D' in locals() else None
        right_D = right_D if 'right_D' in locals() else None
        back_D = back_D if 'back_D' in locals() else None
        front_right_D = front_right_D if 'front_right_D' in locals() else None
        front_left_D = front_left_D if 'front_left_D' in locals() else None
                
        for angle, distance in lidar_data:
            
            #check_front
            if 350 <= angle <= 360 or 0 <= angle <= 10:
                if distance <= self.globals.minDist:
                    self.globals.stop_front = True
                else:
                    self.globals.stop_front = False
                front_D = distance
            #check_left
            if 260 <= angle <= 280:
                if distance <= self.globals.minDist:
                    self.globals.stop_left = True
                else:
                    self.globals.stop_left = False
                left_D = distance
            #check_right
            if 80 <= angle <= 100:
                if distance <= self.globals.minDist:
                    self.globals.stop_right = True
                else:
                    self.globals.stop_right = False
                right_D = distance
            #check_front_left
            if 305 <= angle <= 325:
                if distance <= self.globals.minDist:
                    self.globals.stop_front_left = True
                else:
                    self.globals.stop_front_left = False
                front_left_D = distance
            
            #check_front_right
            if 35 <= angle <= 55:
                if distance <= self.globals.minDist:
                    self.globals.stop_front_right = True
                else:
                    self.globals.stop_front_right = False
                front_right_D = distance

            #check_back
            #if 160 <= angle <= 200:
            #    if distance <= self.globals.minDist:
            #        self.globals.stop_back = True
            #    else:
            #        self.globals.stop_back = False
            #    back_D = distance
                    
        
        return left_D or 10000, front_D or 10000, right_D or 10000, front_left_D or 10000, front_right_D or 10000, back_D or 10000