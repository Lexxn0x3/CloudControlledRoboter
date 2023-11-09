import socket
import matplotlib.pyplot as plt
import matplotlib.animation as animation
import json

def animate(i):
    # Connect to the sender socket
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect((ip_address, port))

    # Receive the data and decode it from JSON
    data = s.recv(1024)
    data_dict = json.loads(data.decode('utf-8'))

    # Close the connection
    s.close()

    # Add the received data to the lists
    angles.append(data_dict['angle'])
    distances.append(data_dict['distance'])

    # Plot the new data
    line1.set_data(angles, distances)

    # Limit the number of data points displayed on the graph
    max_data_points = 100
    if len(angles) > max_data_points:
        angles.pop(0)
        distances.pop(0)

    return line1,

# Replace with the sender's IP address and port
ip_address = '192.168.8.103'
port = 30002

fig, ax = plt.subplots()
ax.set_xlim(0, 360)
ax.set_ylim(0, 5000)

angles = []
distances = []

line1, = ax.plot([], [], 'r-')

ani = animation.FuncAnimation(fig, animate, interval=1000, blit=True)
plt.show()