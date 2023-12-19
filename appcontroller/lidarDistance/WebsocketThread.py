import threading
import time
import websockets
import asyncio
import json
from lidarDistance.globals import stop_front, stop_left, stop_right, stop_front_right, stop_front_left, stop_threads, stop_back

class WebSocketClientThread(threading.Thread):
    def __init__(self, ip_address, app_handler_port):
        super(WebSocketClientThread, self).__init__()
        self.app_handler_port = app_handler_port
        self.ip_address = ip_address
    def run(self):
        try:
            last_flags = {"stop_front": False, "stop_left": False, "stop_right": False, "stop_front_left": False, "stop_front_right": False, "stop_back": False}
            
            while not stop_threads:
                time.sleep(0.01)  # Adjust sleep time as needed

                # Check if any of the stop flags have changed
                if (
                    stop_front != last_flags["stop_front"]
                    or stop_left != last_flags["stop_left"]
                    or stop_right != last_flags["stop_right"]
                    or stop_front_left != last_flags["stop_front_left"]
                    or stop_front_right != last_flags["stop_front_right"]
                    or stop_back != last_flags["stop_back"]
                ):
                    # Send data to the app handler when stop flags change
                    self.send_data_to_app_handler(self.ip_address, self.app_handler_port)

                    # Update the last_flags dictionary
                    last_flags = {
                        "stop_front": stop_front,
                        "stop_left": stop_left,
                        "stop_right": stop_right,
                        "stop_front_left": stop_front_left,
                        "stop_front_right": stop_front_right,
                        "stop_back": stop_back
                    }

        except KeyboardInterrupt:
            print("Stopping")
    
    def send_data_to_app_handler(ip_address, appHandler_port):
        data_to_send = {
            "type": "stop_Flag",
            "stop_front": stop_front,
            "stop_left": stop_left,
            "stop_right": stop_right,
            "stop_front_left": stop_front_left,
            "stop_front_right": stop_front_right,
            "stop_back": stop_back
        }
        uri = f"ws://{ip_address}:{appHandler_port}"
        async def send_data():
            async with websockets.connect(uri) as websocket:
                await websocket.send(json.dumps(data_to_send))
        print("send data")
        asyncio.run(send_data())
        print("send data")