import socket
import cv2
import numpy as np

def main():
    # Set up a TCP server
    tcp_ip = "192.168.8.182"
    tcp_port = 30001

    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.bind((tcp_ip, tcp_port))
    server_socket.listen(1)

    print(f"TCP server listening at {tcp_ip}:{tcp_port}")

    conn, addr = server_socket.accept()
    print(f"Connection from {addr}")

    data = b''
    try:
        while True:
            # Receive data
            packet = conn.recv(65536)
            if not packet: break
            data += packet

            # Check for the start and end of the frame
            a = data.find(b'\xff\xd8')
            b = data.find(b'\xff\xd9')
            if a != -1 and b != -1:
                jpg = data[a:b+2]
                data = data[b+2:]

                # Check if jpg buffer is not empty
                if len(jpg) > 0:
                    # Decode the JPEG data and show the frame
                    frame = cv2.imdecode(np.frombuffer(jpg, dtype=np.uint8), cv2.IMREAD_COLOR)
                    if frame is not None:
                        cv2.imshow('Video Stream', frame)

                # Exit on 'q' key
                if cv2.waitKey(1) & 0xFF == ord('q'):
                    break
    finally:
        cv2.destroyAllWindows()
        conn.close()
        server_socket.close()

if __name__ == "__main__":
    main()
import socket
import cv2
import numpy as np

def main():
    # Set up a TCP server
    tcp_ip = "192.168.8.182"
    tcp_port = 30001

    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.bind((tcp_ip, tcp_port))
    server_socket.listen(1)

    print(f"TCP server listening at {tcp_ip}:{tcp_port}")

    conn, addr = server_socket.accept()
    print(f"Connection from {addr}")

    data = b''
    try:
        while True:
            # Receive data
            packet = conn.recv(65536)
            if not packet: break
            data += packet

            # Check for the start and end of the frame
            a = data.find(b'\xff\xd8')
            b = data.find(b'\xff\xd9')
            if a != -1 and b != -1:
                jpg = data[a:b+2]
                data = data[b+2:]

                # Check if jpg buffer is not empty
                if len(jpg) > 0:
                    # Decode the JPEG data and show the frame
                    frame = cv2.imdecode(np.frombuffer(jpg, dtype=np.uint8), cv2.IMREAD_COLOR)
                    if frame is not None:
                        cv2.imshow('Video Stream', frame)

                # Exit on 'q' key
                if cv2.waitKey(1) & 0xFF == ord('q'):
                    break
    finally:
        cv2.destroyAllWindows()
        conn.close()
        server_socket.close()

if __name__ == "__main__":
    main()
