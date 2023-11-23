import rclpy
from rclpy.node import Node
from sensor_msgs.msg import Image
import socket
import json
import cv2
import numpy as np
import base64

class ImagePublisher(Node):
    def __init__(self):
        super().__init__('image_publisher')
        self.tcp_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.tcp_socket.connect(('192.168.8.225', 30001))  # IP-Adresse und Port entsprechend anpassen
        self.create_subscription(Image, '/camera/color/image_raw', self.image_callback, 10)

    def image_callback(self, msg):
        try:
            height = msg.height
            width = msg.width
            encoding = msg.encoding
            image_data = np.array(msg.data, dtype=np.uint8).reshape((height, width, -1))
            cv_image = cv2.cvtColor(image_data, cv2.COLOR_BGR2RGB)  # Falls die Daten als BGR vorliegen

            # Kodiere das Bild als .jpg
            _, encoded_image = cv2.imencode('.jpg', cv_image)

            # Konvertiere das encoded_image-Bytes-Objekt in Base64
            encoded_image_base64 = base64.b64encode(encoded_image.tobytes()).decode('utf-8')

            # Erstelle ein JSON-Objekt mit den Bilddaten (JPEG) und anderen Informationen
            image_json = {
                "header": {
                    "stamp": {
                        "sec": msg.header.stamp.sec,
                        "nanosec": msg.header.stamp.nanosec
                    },
                    "frame_id": msg.header.frame_id
                },
                "height": height,
                "width": width,
                "encoding": encoding,
                "is_bigendian": msg.is_bigendian,
                "step": msg.step,
                "data": encoded_image_base64
            }

            # Konvertiere das JSON-Objekt in einen JSON-String und sende ihn über TCP
            json_data = json.dumps(image_json) + '\n'
            self.tcp_socket.sendall(json_data.encode())
            self.get_logger().info("JSON-Nachricht (mit JPEG-Daten) über TCP versendet")
        except Exception as e:
            self.get_logger().error(f"Fehler in der Callback-Funktion: {str(e)}")

def main(args=None):
    rclpy.init(args=args)
    node = ImagePublisher()
    rclpy.spin(node)
    node.tcp_socket.close()
    rclpy.shutdown()

if __name__ == '__main__':
    main()

