import socket
import cv2
import numpy as np
from threading import Thread
import time

class MJPEGStreamDecoder:
    def __init__(self, ip, port):
        self.ip = ip
        self.port = port
        self.socket = None
        self.connected = False
        self.running = False
        self.buffer = bytearray()
        self.frame_count = 0

    def connect(self):
        self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        try:
            self.socket.connect((self.ip, self.port))
            self.connected = True
            print("Connection established")
        except Exception as e:
            print(f"Failed to connect to TCP server: {e}")
            self.connected = False

    def start_decoding(self):
        if not self.connected:
            print("Not connected to server")
            return
        self.running = True
        Thread(target=self.decode_stream, daemon=True).start()

    def decode_stream(self):
        try:
            while self.running:
                # Adjust the buffer size based on the frame size
                data = self.socket.recv(65536)  # Assuming a single frame doesn't exceed this size
                if not data:
                    break
                self.buffer += data
                self.process_buffer()

        except Exception as e:
            print(f"Error while decoding stream: {e}")
        finally:
            self.running = False
            self.socket.close()

    def process_buffer(self):
        while True:
            start = self.buffer.find(b'\xff\xd8')
            end = self.buffer.find(b'\xff\xd9', start + 2)
            if start != -1 and end != -1:
                jpg = self.buffer[start:end+2]
                self.buffer = self.buffer[end+2:]
                frame = cv2.imdecode(np.frombuffer(jpg, dtype=np.uint8), cv2.IMREAD_COLOR)
                if frame is not None:
                    self.frame_count += 1
                    filename = f'frame_{self.frame_count}.jpg'
                    cv2.imwrite(filename, frame)
                    print(f"Saved {filename}")
                else:
                    print(f"Decoding failed for the frame between {start} and {end}.")
            else:
                # If no complete frame is found, keep the current partial data and try to read more
                # This handles the case where a frame could be split across multiple recv calls
                if start != -1 and (end == -1 or end < start):
                    # Keep the data from the start marker to the end of the buffer
                    self.buffer = self.buffer[start:]
                break

    def stop(self):
        self.running = False

# Usage
if __name__ == "__main__":
    decoder = MJPEGStreamDecoder("192.168.178.41", 4001)
    decoder.connect()
    decoder.start_decoding()
    try:
        while decoder.running:
            time.sleep(1)  # Sleep to simulate doing other tasks
    except KeyboardInterrupt:
        decoder.stop()
    cv2.destroyAllWindows()
