import socket
import subprocess
import threading
import time
from rplidar import RPLidar
import json

def handle_camera_stream(addr, port):
    cameraStreamProcess = subprocess.Popen(["ffmpeg", "-input_format", "mjpeg", "-i", "/dev/video0", "-c:v", "copy", "-f", "mjpeg", f"tcp://{addr}:{port}"], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    
    while True:
        if cameraStreamProcess.poll() is not None:
            print("Camera stream process has exited.")
            break
        time.sleep(1)  # Check every second

def handle_lidar_stream(addr, port):
    lidar = RPLidar('/dev/rplidar')
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect((addr, port))

    try:
        for new_scan, quality, angle, distance in lidar.iter_measures():
            distance = round(distance)
            data = {'new_scan': new_scan, 'angle': angle, 'distance': distance}
            sock.sendall((json.dumps(data) + '\n').encode('utf-8'))
    finally:
        lidar.stop()
        lidar.disconnect()
        sock.close()

def main():
    tcp_ip = "0.0.0.0"
    tcp_port = 6969

    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.bind((tcp_ip, tcp_port))
    server_socket.listen(1)

    print(f"TCP server listening at {tcp_ip}:{tcp_port}")

    conn, addr = server_socket.accept()
    print(f"Connection from {addr}")

    data = b''
    camera_thread = None
    lidar_thread = None

    try:
        while True:
            packet = conn.recv(65536)
            if not packet: 
                break
            data += packet

            if packet.lower().startswith(b'startstreams'):
                port = int(packet[len(b'startstreams'):].strip().decode())
                if camera_thread is None or not camera_thread.is_alive():
                    camera_thread = threading.Thread(target=handle_camera_stream, args=(addr[0], port))
                    camera_thread.start()
                if lidar_thread is None or not lidar_thread.is_alive():
                    lidar_thread = threading.Thread(target=handle_lidar_stream, args=(addr[0], port + 1))
                    lidar_thread.start()
                conn.sendall(b'ok')

            elif packet.lower().startswith(b'stopstreams'):
                if camera_thread is not None and camera_thread.is_alive():
                    camera_thread.terminate()  # Terminate the camera thread
                    camera_thread = None
                if lidar_thread is not None and lidar_thread.is_alive():
                    lidar_thread.terminate()  # Terminate the lidar thread
                    lidar_thread = None
    except Exception as e:
        print(f"An error occurred: {e}")
    finally:
        conn.close()
        server_socket.close()

if __name__ == "__main__":
    main()
