import socket
from collections import deque
import sys

# Function to create a socket to connect to server
def connect_to_server(ip_address, server_port):
    
    # Create a TCP/IP socket
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    
    # Connect to the server
    server_address = (ip_address, server_port)
    sock.connect(server_address)
    return sock

# prgenerate a buffer
def pregenerate_lidar_data(maxlen):
    lidar_data_buffer = deque(maxlen=maxlen)
    return lidar_data_buffer

def signal_handler(sig, frame):
    print("Ctrl+C detected. Exiting...")
    sys.exit(0)