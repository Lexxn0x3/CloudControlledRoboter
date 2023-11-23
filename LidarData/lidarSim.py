import socket
import json
import random
import time

def generate_lidar_data():
    # Simulate realistic lidar data with one decimal for each angle
    angle = round(random.uniform(0, 359), 1)
    distance = round(random.uniform(101, 2500))

    return {'angle': angle, 'distance': distance}

def send_lidar_data(ip, port):
    # Create a TCP socket
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    try:
        # Connect to the server
        client_socket.connect((ip, port))
        print(f"Connected to {ip}:{port}")

        while True:
            # Generate lidar data
            lidar_data = generate_lidar_data()

            # Convert data to JSON
            data_json = json.dumps(lidar_data)

            # Send data to the server
            client_socket.sendall(data_json.encode('utf-8'))
            print(f"Sent data: {data_json}")

            # Wait for a short time before sending the next data
            time.sleep(0.1)

    except Exception as e:
        print(f"Error: {e}")

    finally:
        # Close the socket
        client_socket.close()

if __name__ == "__main__":
    # Set the IP address and port of the server (broker)
    server_ip = "192.168.8.20"
    server_port = 3011

    # Start sending lidar data to the server
    send_lidar_data(server_ip, server_port)
