import rclpy
from rclpy.node import Node
from sensor_msgs.msg import Image
import socket
import cv2
import numpy as np

class ImagePublisher(Node):
    def __init__(self):
        super().__init__('image_publisher')
        self.image_subscriber = self.create_subscription(Image, '/camera/color/image_raw', self.image_callback, 10)
        self.tcp_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.tcp_socket.connect(('192.168.8.225', 30001))  # IP-Adresse und Port entsprechend anpassen

    def image_callback(self, msg):
        try:
            # Konvertiere die Rohdaten in ein OpenCV-Bild
            height = msg.height
            width = msg.width
            encoding = msg.encoding
            image_data = np.array(msg.data, dtype=np.uint8).reshape((height, width, -1))
            cv_image = cv2.cvtColor(image_data, cv2.COLOR_BGR2RGB)  # Falls die Daten als BGR vorliegen
            # Kodiere das Bild als .jpg und sende es über TCP
            _, encoded_image = cv2.imencode('.jpg', cv_image)
            self.tcp_socket.sendall(encoded_image.tobytes())
            self.get_logger().info("Bild als .jpg über TCP versendet")
            self.image_subscriber.destroy()  # Beende den Subscriber nach dem Versenden des Bildes
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

