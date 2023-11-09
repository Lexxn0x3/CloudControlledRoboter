import socket
import json

# Configure TCP Connection
ip_address = '192.168.8.103'
server_port = 30002

# Create a TCP/IP socket
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

# Bind the socket to the port
server_address = (ip_address, server_port)
sock.bind(server_address)

# Listen for incoming connections
sock.listen(1)

while True:
    # Wait for a connection
    print('waiting for a connection')
    connection, client_address = sock.accept()

    try:
        print('connection from', client_address)

        # Receive the data in small chunks and retransmit it
        while True:
            data = connection.recv(1024)
            if data:
                # Decode the JSON data and extract the distance and angle
                data = json.loads(data.decode('utf-8'))
                angle = data['angle']
                distance = data['distance']

                # Print the distance measurement with the corresponding angle
                print(f'Angle: {angle}, Distance: {distance}')
            else:
                break
            
    except KeyboardInterrupt:
        break
    finally:
        # Clean up the connection
        connection.close()
