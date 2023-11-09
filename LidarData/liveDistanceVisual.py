import socket
import pandas as pd
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

# Create a list to store the data
data_list = []

# Set up the plot
fig, ax = plt.subplots()

def update(i):
    # Wait for a connection
    connection, client_address = sock.accept()

    try:
        # Receive the data in small chunks and add it to the list
        data = connection.recv(512)
        if data:
            data_list.append(float(data.decode('utf-8')))

            # Clear the current plot
            ax.clear()

            # Plot the new data
            ax.plot(data_list)

    finally:
        # Clean up the connection
        connection.close()

# Set up the animation
ani = FuncAnimation(fig, update, interval=1000)

# Show the plot
plt.show()
