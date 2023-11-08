import socket
import os

for i in range(4, 9999):
    # Set the server's port and IP address
    SERVER_HOST = '192.168.8.225'
    SERVER_PORT = 30001

    # Create a socket object with IPv4 addressing and TCP protocol
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # Bind the socket to the address and port
    server_socket.bind((SERVER_HOST, SERVER_PORT))

    # Enable the server to accept connections, with the backlog parameter set to 5
    server_socket.listen(5)

    print(f"Listening as {SERVER_HOST}:{SERVER_PORT} ...")


    # Accept a connection
    client_socket, client_address = server_socket.accept()
    print(f"{client_address[0]}:{client_address[1]} Connected!")

    # Define the file path
    file_name = 'received_data' + str(i) + '.dat'
    file_path = os.path.abspath(file_name)
    print(f"Data will be written to: {file_path}")

    # Open the file to write the incoming packets
    with open(file_name, 'ab') as f:
        while True:
            # Receive data from the client
            data = client_socket.recv(4096)
            if not data:
                # No more data from client, close the connection
                print(f"{client_address[0]}:{client_address[1]} Disconnected!")
                break
            # Write the received data to the file
            f.write(data)
            # Optionally, send some acknowledgment to the client (not necessary)
            # client_socket.send("Data received\n".encode())

    # Close the client socket
    client_socket.close()

    # Close the server socket (unreachable in this snippet)
    server_socket.close()
