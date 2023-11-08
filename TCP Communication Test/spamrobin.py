import rclpy
from rclpy.node import Node
from sensor_msgs.msg import Image
import socket
import json

class ImagePublisher(Node):
    def __init__(self):
        super().__init__('image_publisher')
        self.tcp_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.tcp_socket.connect(('192.168.8.225', 30001))  # IP-Adresse und Port entsprechend anpassen
        self.create_subscription(Image, '/camera/color/image_raw', self.image_callback, 10)

    def image_callback(self, msg):
        try:
            # Erstelle ein JSON-Objekt mit den Bilddaten und anderen Informationen
            image_json = {
                "header": {
                    "stamp": {
                        "sec": msg.header.stamp.sec,
                        "nanosec": msg.header.stamp.nanosec
                    },
                    "frame_id": msg.header.frame_id
                },
                "height": msg.height,
                "width": msg.width,
                "encoding": msg.encoding,
                "is_bigendian": msg.is_bigendian,
                "step": msg.step,
                "data": list(msg.data)  # Füge msg.data direkt im JSON-Objekt ein
            }
            # Konvertiere das JSON-Objekt in einen JSON-String, füge einen Carriage Return und einen Line Feed hinzu und sende ihn über TCP
            json_data = json.dumps(image_json) + '\r\n'
            self.tcp_socket.sendall(json_data.encode())
            self.get_logger().info("JSON-Nachricht über TCP versendet")
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

