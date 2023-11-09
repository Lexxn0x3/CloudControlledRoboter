import socket
import json

# Configure TCP Connection
ip_address = '192.168.8.103'
server_port = 30002

# Create a TCP socket
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

# Connect to the server
sock.connect((ip_address, server_port))

while True:
    try:
        # Receive distance measurement data from the sender
        data = sock.recv(1024)
        if not data:
            break

        # Decode the JSON data and extract the distance and angle
        data = json.loads(data.decode('utf-8'))
        angle = data['angle']
        distance = data['distance']

        # Print the distance measurement with the corresponding angle
        print(f'Angle: {angle}, Distance: {distance}')

    except KeyboardInterrupt:
        break

# Close the TCP socket
sock.close()