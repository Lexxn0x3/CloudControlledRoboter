# serverLidar_json.py

import socket
import json

# Configure TCP Connection
ip_address = '192.168.8.103'
server_port = 30002

def initialize_socket():
    # Create a TCP/IP socket
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # Bind the socket to the port
    server_address = (ip_address, server_port)
    sock.bind(server_address)

    return sock

def pregenerate_lidar_data():
    # Pregenerate Dic for saving lidar data
    lidar_data = {}

    for i in range(3600):
        angle = i/10.0
        lidar_data[angle] = None

    return lidar_data

def receive_lidar_data(sock, lidar_data):
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
                    # Split the received data by newline and process each line as a separate JSON object
                    for line in data.splitlines():
                        try:
                            # Decode the JSON data and extract the distance and angle
                            data = json.loads(line.decode('utf-8'))
                            angle = data['angle']
                            distance = data['distance']

                            lidar_data[round(angle, 1)] = distance

                        except json.JSONDecodeError:
                            # print('Failed to decode JSON object: ', line)
                            continue
                else:
                    break

        except KeyboardInterrupt:
            break
        finally:
            # Clean up the connection
            connection.close()
            for zeile in lidar_data:
                print(zeile, ": ", lidar_data[zeile])

if __name__ == '__main__':
    sock = initialize_socket()
    lidar_data = pregenerate_lidar_data()
    receive_lidar_data(sock, lidar_data)
