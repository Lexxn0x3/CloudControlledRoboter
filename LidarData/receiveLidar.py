import socket

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

        # Receive the data in small chunks and print it
        while True:
            data = connection.recv(512)
            if data:
                print('Received "%s"' % data)
            else:
                print('No more data from', client_address)
                break

    finally:
        # Clean up the connection
        connection.close()
