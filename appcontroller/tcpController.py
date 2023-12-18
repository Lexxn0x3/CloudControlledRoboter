import time
import socket
from betterPrinting import BetterPrinting

 
#Class responsible for establishing and managing a TCP connection to the robot.

class TCPController:
    def __init__(self, robot_host, robot_port, bot_instance):
        self.robot_host = robot_host    # Store the robot host (IP address) for TCP communication
        self.robot_port = robot_port    # Store the port number for TCP communication
        self.robot_socket = None        # Initialize the robot_socket as None initially
        self.auto_reconnect = True      # Flag to enable automatic reconnection in case of connection issues
        self.connect()              # Establish the initial connection to the robot

        self.bot = bot_instance         # Reference to the main Bot instance
        self.bp = BetterPrinting(self.bot.i_print , self.bot.d_print, self.bot.e_print) # Setup Better Printing

    #Establishes a connection to the robot.
    def connect(self):
        try:
            # Create a socket and connect to the robot
            self.robot_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.robot_socket.setsockopt(socket.IPPROTO_TCP, socket.TCP_NODELAY, True)
            self.robot_socket.connect((self.robot_host, int(self.robot_port)))
            print(f"Connected to {self.robot_host}:{self.robot_port}")

        except Exception as e:
            self.bp.error_print(f"Error connecting to the robot: {e}")
            if self.auto_reconnect:
                self.bp.debug_print("Retrying in 2 seconds...")
                time.sleep(2)
                self.connect()

    #Sends JSON data to the robot.
    def send_json_data(self, json_data):
        try:
            if not self.robot_socket:
                self.connect()

            # Send the JSON data to the server
            json_data += '\n'   # Append a newline character to the JSON data string, because newline-delimited JSON is used
            self.robot_socket.sendall(json_data.encode('utf-8'))
            self.bp.debug_print(f"Sent JSON data: {json_data}")

        except Exception as e:
            self.bp.error_print(f"Error sending JSON data: {e}")
            if self.auto_reconnect:
                self.bp.debug_print("Retrying in 2 seconds...")
                time.sleep(2)
                self.connect()


    def close_connection(self):
        if self.robot_socket:
            self.robot_socket.close()
            print("Closed the connection.")