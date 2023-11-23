import socket
import json
import base64
import threading
import time  # Make sure to import the time module

# Set the server's port for the JSON TCP server and the MJPEG HTTP server
TCP_PORT = 30001
HTTP_PORT = 8080

# Frame queue to hold the frames received
frame_queue = []

# Lock for thread-safe operations on the frame queue
queue_lock = threading.Lock()

# Create a socket object for the JSON TCP server with IPv4 addressing and TCP protocol
tcp_server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
tcp_server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
tcp_server_socket.bind(('192.168.8.225', TCP_PORT))
tcp_server_socket.listen(5)

# Function to handle the TCP connection and receive JSON objects
def tcp_connection_handler(client_socket):
    try:
        while True:
            # Receive JSON object from the client
            json_data = b''
            while True:
                chunk = client_socket.recv(4096)
                if not chunk:  # No more data from the client
                    return
                json_data += chunk
                # Check if we have a complete JSON object
                try:
                    json_object = json.loads(json_data.decode('utf-8'))
                    # Extract and decode the base64 encoded image
                    if 'data' in json_object:
                        image_data = base64.b64decode(json_object['data'])
                        with queue_lock:
                            frame_queue.append(image_data)
                    json_data = b''  # Clear to receive the next JSON object
                except json.JSONDecodeError:
                    continue  # Continue receiving data if JSON is incomplete
    finally:
        client_socket.close()

# Function to serve the MJPEG stream over HTTP
def mjpeg_stream_handler(client_connection):
    try:
        client_connection.send(b'HTTP/1.1 200 OK\r\n')
        client_connection.send(b'Content-Type: multipart/x-mixed-replace; boundary=frame\r\n\r\n')
        while True:
            with queue_lock:
                if frame_queue:
                    frame = frame_queue.pop(0)
                else:
                    frame = None
            if frame:
                message = b'--frame\r\nContent-Type: image/jpeg\r\n\r\n' + frame + b'\r\n'
                client_connection.send(message)
            else:
                # If no frame is available, wait before trying again
                time.sleep(0.05)
    except Exception as e:
        print(f"Stream Client Disconnected: {e}")
    finally:
        client_connection.close()


# Start the TCP server in a new thread
def start_tcp_server():
    print(f"TCP Server listening on port {TCP_PORT}...")
    while True:
        client_socket, addr = tcp_server_socket.accept()
        print(f"TCP Client {addr[0]}:{addr[1]} Connected!")
        client_thread = threading.Thread(target=tcp_connection_handler, args=(client_socket,))
        client_thread.start()

# Start the HTTP server in a new thread
def start_http_server():
    http_server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    http_server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    http_server_socket.bind(('0.0.0.0', HTTP_PORT))
    http_server_socket.listen(5)
    print(f"HTTP Server listening on port {HTTP_PORT}...")
    while True:
        client_connection, addr = http_server_socket.accept()
        print(f"HTTP Client {addr[0]}:{addr[1]} Connected!")
        stream_thread = threading.Thread(target=mjpeg_stream_handler, args=(client_connection,))
        stream_thread.start()

# Run the TCP and HTTP servers
tcp_server_thread = threading.Thread(target=start_tcp_server)
tcp_server_thread.start()

http_server_thread = threading.Thread(target=start_http_server)
http_server_thread.start()

tcp_server_thread.join()
http_server_thread.join()
