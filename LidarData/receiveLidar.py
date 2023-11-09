import socket

#Confifure TCP Conection
host = '192.168.8.10'
port = 30002

# Create a TCP socket
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

# Bind the socket to an address and port
sock.bind(('host', port))

# Listen for incoming connections
sock.listen(1)

# Accept a connection from the client
client_socket, address = sock.accept()

while True:
    # Receive data from the client
    data = client_socket.recv(1024)

    # Print the received data to the console
    print(data)
