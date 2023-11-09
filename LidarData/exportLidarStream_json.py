import socket
from rplidar import RPLidar
import json

# Configure TCP Connection
ip_address = '192.168.8.103'
server_port = 30002

# Create a TCP socket
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

# Connect to the server
sock.connect((ip_address, server_port))

# Initialize RPLidar
lidar = RPLidar('/dev/rplidar')

try:
    print('Recording distance measurements... Press Ctrl+C to stop.')
    for scan in lidar.iter_scans():
        for (_, angle, distance) in scan:
            # Round the distance to the nearest whole number
            distance = round(distance)

            # Create a dictionary with the angle and distance
            data_dict = {'angle': angle, 'distance': distance}

            # Convert the dictionary to a JSON string
            data_json = json.dumps(data_dict)

            # Send the JSON string to the server
            sock.sendall(data_json.encode('utf-8'))
finally:
    print('Stopping.')
    lidar.stop()
    lidar.disconnect()
