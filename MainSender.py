import socket
import subprocess
import cv2
import numpy as np

def main():
    # Set up a TCP server
    tcp_ip = "0.0.0.0"
    tcp_port = 6969

    cameraStreamProcess = None

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
            print(packet)

            if cameraStreamProcess is not None:
                cameraStreamProcess.poll()
                if cameraStreamProcess.returncode is not None:
                    print("FAILED STREAMER")

            keyword = b'startstreams'
            if packet.lower().startswith(keyword):
                port = packet[len(keyword):]
                print("Start streams")
                conn.sendall(b'ok')
                cameraStreamProcess = subprocess.Popen(["ffmpeg", "-input_format", "mjpeg", "-i", "/dev/video0", "-c:v", "copy", "-f", "mjpeg", "tcp://" + addr[0] + ":" + port.decode()], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            elif packet.lower().startswith(b'stopstreams'):
                print("Stop streams")
                if cameraStreamProcess is not None:
                    cameraStreamProcess.terminate()
                    cameraStreamProcess = None
    finally:
        cv2.destroyAllWindows()
        conn.close()
        server_socket.close()

if __name__ == "__main__":
    main()