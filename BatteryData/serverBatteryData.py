import socket
import json


# Define the server address and port
server_address = ('192.168.8.103', 30003)

# Create a TCP socket
server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

# Bind the socket to the server address
server_socket.bind(server_address)

# Listen for incoming connections
server_socket.listen(1)
print(f"Waiting for a connection on {server_address}")

# Accept the connection
client_socket, client_address = server_socket.accept()
print(f"Connection from {client_address}")

try:
    while True:
        # Receive motor data from the client
        data = client_socket.recv(1024)
        if not data:
            break

        # Split the received data by newline and process each line as a separate JSON object
        for line in data.splitlines():
            try:
                # Decode the JSON data
                motor_data = json.loads(line.decode('utf-8'))

                # Process the received motor data (replace with your logic)
                print("Received motor data:", motor_data)
            except json.JSONDecodeError:
                # Handle JSON decoding errors
                continue

finally:
    # Clean up the connection
    client_socket.close()
    server_socket.close()