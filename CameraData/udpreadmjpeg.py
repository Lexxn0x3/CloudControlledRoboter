import cv2
import numpy as np
import socket

def main():
    # Set up a UDP server
    udp_ip = "192.168.8.103"
    udp_port = 30001

    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.bind((udp_ip, udp_port))

    print("UDP server listening at {}:{}".format(udp_ip, udp_port))

    data = b''
    try:
        while True:
            # Receive data
            packet, _ = sock.recvfrom(65536)
            data += packet

            # Check for the start and end of the frame
            a = data.find(b'\xff\xd8')
            b = data.find(b'\xff\xd9')
            if a != -1 and b != -1:
                jpg = data[a:b+2]
                data = data[b+2:]

                # Decode the JPEG data and show the frame
                if len(jpg) > 0:
                    frame = cv2.imdecode(np.frombuffer(jpg, dtype=np.uint8), cv2.IMREAD_COLOR)
                    if frame is not None:
                        cv2.imshow('Video Stream', frame)

                # Exit on 'q' key
                if cv2.waitKey(1) & 0xFF == ord('q'):
                    break
    finally:
        cv2.destroyAllWindows()
        sock.close()

if __name__ == "__main__":
    main()
