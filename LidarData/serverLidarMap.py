import matplotlib.pyplot as plt
import numpy as np
from serverLidarData_json import lidar_data  # Assuming lidar_data is defined in serverLidarData_json

# Initialize the 2D map
plt.ion()  # Activate interactive mode
fig, ax = plt.subplots()
ax.set_xlim(-10, 10)  # Map width
ax.set_ylim(0, 10)    # Map height

try:
    while True:
        # Clear the existing scatter plot
        ax.clear()

        # Extract valid distances and angles for plotting
        valid_angles = [a for a, d in lidar_data.items() if d is not None]
        valid_distances = [d for d in lidar_data.values() if d is not None]

        # Convert polar coordinates to Cartesian coordinates
        x = valid_distances * np.cos(np.radians(valid_angles))
        y = valid_distances * np.sin(np.radians(valid_angles))

        # Plot the Lidar data on the 2D map
        ax.scatter(x, y, marker='o', color='blue')

        # Optionally, you can add additional features to the plot

        # Update the plot
        plt.draw()
        plt.pause(0.001)

except KeyboardInterrupt:
    pass

finally:
    plt.ioff()  # Deactivate interactive mode
    plt.show()
