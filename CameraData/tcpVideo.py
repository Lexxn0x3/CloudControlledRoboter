import socket
import cv2
import numpy as np

# Define the host and port for receiving the video stream
HOST = "0.0.0.0"  # Change this to the actual IP address if needed
PORT = 30001

# Create a TCP socket
server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server_socket.bind((HOST, PORT))
server_socket.listen(0)

print(f"Server listening for connections on {HOST}:{PORT}")

# Accept an incoming connection
connection, address = server_socket.accept()
print(f"Connected to client at {address}")

# Continuously receive and decode video frames
while True:
    # Receive the frame size
    data = connection.recv(4)
    frame_size = np.frombuffer(data, dtype=np.int32)
    frame_size = frame_size[0]

    # Receive the frame data
    data = connection.recv(frame_size)
    frame_data = np.frombuffer(data, dtype=np.uint8)
    frame = cv2.imdecode(frame_data, cv2.IMREAD_COLOR)

    # Display the received frame
    cv2.imshow('Video Stream', frame)

    # Check if the user wants to quit
    if cv2.waitKey(1) & 0xFF == ord('q'):
        break

# Close the connection and release resources
connection.close()
server_socket.close()
cv2.destroyAllWindows()