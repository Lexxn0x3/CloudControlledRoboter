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

# Create a DataFrame to store the data
df = pd.DataFrame()

# Set up the plot
fig, ax = plt.subplots()

def update(i):
    # Wait for a connection
    connection, client_address = sock.accept()

    try:
        # Receive the data in small chunks and add it to the DataFrame
        data = connection.recv(512)
        if data:
            df_new = pd.read_csv(data, sep=',')
            df = pd.concat([df, df_new])

            # Clear the current plot
            ax.clear()

            # Plot the new data
            df.plot(ax=ax)

    finally:
        # Clean up the connection
        connection.close()

# Set up the animation
ani = FuncAnimation(fig, update, interval=1000)

# Show the plot
plt.show()
