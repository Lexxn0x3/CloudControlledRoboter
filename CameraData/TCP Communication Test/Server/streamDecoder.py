import socket
import json
import base64

# Set the server's port and IP address
SERVER_HOST = '192.168.8.225'
SERVER_PORT = 30001

# Create a socket object with IPv4 addressing and TCP protocol
server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

# Set the socket option to reuse the address
server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)

# Bind the socket to the address and port
server_socket.bind((SERVER_HOST, SERVER_PORT))

# Enable the server to accept connections
server_socket.listen(5)

print(f"Listening as {SERVER_HOST}:{SERVER_PORT} ...")

# Function to receive the entire JSON object
def receive_json_object(sock):
    json_data = ''
    while True:
        chunk = sock.recv(4096).decode('utf-8')  # Receive chunk of data from the client
        if not chunk:  # No more data from the client
            return None
        json_data += chunk  # Append the chunk to the JSON data
        # Check if we have a complete JSON object
        try:
            # Attempt to parse the JSON data
            json_object = json.loads(json_data)
            return json_object  # Return the parsed JSON object
        except json.JSONDecodeError as e:
            # If JSON is incomplete, continue to accumulate data
            continue

# Main loop to accept connections
while True:
    # Accept a connection
    client_socket, client_address = server_socket.accept()
    print(f"{client_address[0]}:{client_address[1]} Connected!")

    # Initialize a counter for the image file name
    i = 0

    # This inner loop will keep listening for data from the client
    while True:
        json_object = receive_json_object(client_socket)
        if json_object is None:
            break

        # Process the JSON object
        try:
            if 'data' in json_object:
                # Decode the base64 encoded image
                raw_image = base64.b64decode(json_object['data'])
                # Write the image to a file
                image_file_name = f'image_{i}.jpg'
                with open(image_file_name, 'wb') as image_file:
                    image_file.write(raw_image)
                    print(f"Saved image to {image_file_name}")
                i += 1
            else:
                print("No 'data' key in JSON object")

        except KeyError as e:
            print(f"Key error: {e}")
        except Exception as e:
            print(f"An unexpected error occurred: {e}")

    # After the inner loop ends, the client has disconnected.
    print(f"{client_address[0]}:{client_address[1]} Disconnected.")
    client_socket.close()

# Close the server socket (unreachable in this snippet)
server_socket.close()