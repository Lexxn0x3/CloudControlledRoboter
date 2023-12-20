import threading
import time
import websockets
import asyncio
import json

class WebSocketClientThread(threading.Thread):
    def __init__(self, globals, ip_address, app_handler_port):
        super(WebSocketClientThread, self).__init__()
        self.app_handler_port = app_handler_port
        self.ip_address = ip_address
        self.globals = globals
        print(f"websocket ip {self.ip_address}", flush= True)
        print(f"websocket port {self.app_handler_port}", flush= True)

    def run(self):
        
        try:
            last_flags = {"stop_front": False, "stop_left": False, "stop_right": False, "stop_front_left": False, "stop_front_right": False, "stop_back": False}
            
            while not self.globals.stop_threads:
                time.sleep(0.01)  # Adjust sleep time as needed
                # Check if any of the stop flags have changed
                if (
                    self.globals.stop_front != last_flags["stop_front"]
                    or self.globals.stop_left != last_flags["stop_left"]
                    or self.globals.stop_right != last_flags["stop_right"]
                    or self.globals.stop_front_left != last_flags["stop_front_left"]
                    or self.globals.stop_front_right != last_flags["stop_front_right"]
                    or self.globals.stop_back != last_flags["stop_back"]
                ):
                    # Send data to the app handler when stop flags change
                    self.send_data_to_app_handler(self)

                    # Update the last_flags dictionary
                    last_flags = {
                        "stop_front": self.globals.stop_front,
                        "stop_left": self.globals.stop_left,
                        "stop_right": self.globals.stop_right,
                        "stop_front_left": self.globals.stop_front_left,
                        "stop_front_right": self.globals.stop_front_right,
                        "stop_back": self.globals.stop_back
                    }

        except KeyboardInterrupt:
            print("Stopping")
    
    def send_data_to_app_handler():

        data_to_send = {
            "type": "stop_Flag",
            "stop_front": self.globals.stop_front,
            "stop_left": self.globals.stop_left,
            "stop_right": self.globals.stop_right,
            "stop_front_left": self.globals.stop_front_left,
            "stop_front_right": self.globals.stop_front_right,
            "stop_back": self.globals.stop_back
        }
        print(data_to_send)
        uri = f"ws://{self.ip_address}:{self.app_handler_port}"
        print("sending", flush= True)
        async def send_data():
            async with websockets.connect(uri) as websocket:
                await websocket.send(json.dumps(data_to_send))
        print("send data", flush= True)
        asyncio.run(send_data())
        print("send data")