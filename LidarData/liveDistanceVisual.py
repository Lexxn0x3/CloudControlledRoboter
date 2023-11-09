import socket
import json
import matplotlib.pyplot as plt
from matplotlib.animation import FuncAnimation

# Configure TCP server
ip_address = '192.168.8.103'
server_port = 30002

# Create a TCP socket
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

# Bind the socket to the port
sock.bind((ip_address, server_port))

# Listen for incoming connections
sock.listen(1)

# Create lists to store the angles and distances
angles = []
distances = []

# Set up the plot
fig, ax = plt.subplots()

def update(i):
    # Wait for a connection
    connection, client_address = sock.accept()

    try:
        # Receive the data in small chunks and add it to the list
        data = connection.recv(512)
        if data:
            # Decode the data and convert it from JSON to a dictionary
            data_dict = json.loads(data.decode('utf-8'))

            # Extract the angle and distance and add them to the lists
            angles.append(data_dict['angle'])
            distances.append(data_dict['distance'])

            # Clear the current plot
            ax.clear()

            # Plot the new data
            ax.plot(angles, distances)

    finally:
        # Clean up the connection
        connection.close()

# Set up the animation
ani = FuncAnimation(fig, update, interval=1000)

# Show the plot
plt.show()
