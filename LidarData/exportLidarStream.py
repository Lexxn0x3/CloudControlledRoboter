import socket
from rplidar import RPLidar

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
    print('Recording measurments... Press Ctrl+C to stop.')
    for measurment in lidar.iter_measurments():
        # Send data packet to the server
        sock.sendall(str(measurment).encode('utf-8'))
finally:
    print('Stopping.')
    lidar.stop()
    lidar.disconnect()
