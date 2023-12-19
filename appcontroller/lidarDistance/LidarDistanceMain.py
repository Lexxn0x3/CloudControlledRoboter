import sys
from lidarDistance.LidarDataThread import LidarDataThread
from lidarDistance.DistanceSenderThread import DistanceSenderThread
from lidarDistance.WebsocketThread import WebSocketClientThread
from lidarDistance.utils import pregenerate_lidar_data
from lidarDistance.globals import stop_threads, maxLenBuffer
from lidarDistance.SpeedDataThread import SpeedThread


class LidarDistanceSystem:
    def __init__(self, lidarstreamIP="192.168.8.20", lidarstreamPort=9011, appHandlerIP="192.168.8.20", appHandlerPort=6942, distanceSenderIP="192.168.8.20", distanceSenderPort=3031):
        self.lidarstreamIP = lidarstreamIP
        self.lidarstreamPort = lidarstreamPort
        self.appHandlerIP = appHandlerIP
        self.appHandlerPort = appHandlerPort
        self.distanceSenderIP = distanceSenderIP
        self.distanceSenderPort = distanceSenderPort

        # Create a LidarDataThread instance
        self.lidar_data_buffer = pregenerate_lidar_data(maxLenBuffer)
        self.lidar_data_thread = LidarDataThread(self.lidarstreamIP, self.lidarstreamPort, self.lidar_data_buffer)
        
        # Create a DistanceSenderThread instance if distanceSenderIP and distanceSenderPort are provided
        self.sender_thread = None
        if self.distanceSenderIP is not None and self.distanceSenderPort is not None:
            self.sender_thread = DistanceSenderThread(self.distanceSenderIP, self.distanceSenderPort, self.lidar_data_buffer)
        
        # Create a WebSocket client thread
        self.websocket_client_thread = WebSocketClientThread(self.appHandlerIP, self.appHandlerPort)
        
        print("everything startet")
        self.start_processing()

    def start_processing(self):
        # Start LidarDataThread
        self.lidar_data_thread.start()

        # Start DistanceSenderThread if it was created
        if self.sender_thread is not None:
            self.sender_thread.start()

        # Start WebSocket client thread
        self.websocket_client_thread.start()

        try:
            while not stop_threads:
                # You can add any other code you want to run continuously here
                pass
        except KeyboardInterrupt:
            # Handle Ctrl+C interruption
            print("Ctrl+C detected. Stopping threads...")
            self.lidar_data_thread.join()  # Wait for LidarDataThread to finish
            if self.sender_thread is not None:
                self.sender_thread.join()  # Wait for DistanceSenderThread to finish
            self.websocket_client_thread.join()
            print("Threads stopped. Exiting...")


if __name__ == '__main__':
    #lidarstreamIP = '192.168.8.20'
    #lidarstreamPort = 9011
    #appHandlerIP = '192.168.8.20'
    #appHandlerPort = 6942
    #distanceSenderIP = '192.168.8.20'
    #distanceSenderPort = 3031

    # Check if the correct number of command-line arguments is provided
    if len(sys.argv) < 5:
        print("Usage: python your_script.py  lidarstream_ip lidarstream_port appHandler_ip appHandler_port (distancestream_ip distancestream_port)")
    
    elif len(sys.argv) == 5:
        # Extract variables from command-line arguments
        lidarstreamIP, lidarstreamPort, appHandlerIP, appHandlerPort = sys.argv[1:5]

    elif len(sys.argv) <= 7:
        # Extract variables from command-line arguments with distance sending
        lidarstreamIP, lidarstreamPort, appHandlerIP, appHandlerPort, distanceSenderIP, distanceSenderPort = sys.argv[1:7]

    # Create an instance of LidarProcessingSystem and start processing
    lidar_system = LidarDistanceSystem()
