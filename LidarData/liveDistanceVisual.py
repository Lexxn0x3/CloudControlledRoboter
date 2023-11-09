import socket
import pandas as pd
import matplotlib.pyplot as plt

# Configure TCP server
ip_address = '192.168.8.103'
server_port = 30002

# Create a TCP socket
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

# Bind the socket to the port
sock.bind((ip_address, server_port))

# Listen for incoming connections
sock.listen(1)

# Create a pandas dataframe to store the distance data
df = pd.DataFrame(columns=['Distance'])

# Create a matplotlib figure
fig, ax = plt.subplots()

# Plot the distance data live
while True:

    # Wait for a new data point
    data = sock.recv(512)

    # Add the new data point to the pandas dataframe
    df.append({'Distance': float(data)}, ignore_index=True)

    # Update the matplotlib plot
    ax.clear()
    ax.plot(df['Distance'])
    ax.set_title('Distance Measurements')
    ax.set_xlabel('Sample Number')
    ax.set_ylabel('Distance (m)')
    fig.canvas.draw()