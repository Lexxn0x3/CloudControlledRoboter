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

while True:
    # Wait for a connection
    print('Waiting for a connection...')
    connection, client_address = sock.accept()

    try:
        print('Connection from', client_address)

        # Create a pandas dataframe to store the distance data
        df = pd.DataFrame(columns=['Distance'])

        # Open a new window to show the live graph
        fig, ax = plt.subplots()

        # Plot the distance data live
        while True:

            # Receive a new data point
            data = connection.recv(512)

            # Decode the incoming stream
            data_decoded = data.decode('utf-8')

            # Convert the decoded stream to a float
            data_float = float(data_decoded)

            # Add the new data point to the pandas dataframe
            df.append({'Distance': data_float}, ignore_index=True)

            # Update the matplotlib plot
            ax.clear()
            ax.plot(df['Distance'])
            ax.set_title('Distance Measurements')
            ax.set_xlabel('Sample Number')
            ax.set_ylabel('Distance (m)')
            fig.canvas.draw()

            plt.pause(0.05)

    finally:
        # Close the new window
        plt.close(fig)

        # Clean up the connection
        connection.close()