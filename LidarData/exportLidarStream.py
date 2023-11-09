import socket
import serial

# Configure serial port settings
port = '/dev/rplidar'  # Replace with the actual port name
baudrate = 115200

# Configure TCP Conection
ipAdress = '192.168.8.103'
serverport = 30002

# Open serial connection
ser = serial.Serial(port, baudrate)

# Create a TCP socket
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

# Connect to the server
sock.connect((ipAdress, serverport))

while True:
    # Read data packet from the lidar
    data_packet = ser.read(512)

    # Send data packet to the server
    sock.sendall(data_packet)
