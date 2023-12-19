import asyncio
import json
import websockets
from betterPrinting import BetterPrinting

class WebSocketController:
    def __init__(self, ip, port, datahandler, bot_instance):
        self.ip = ip            # Store the IP address for WebSocket communication
        self.port = port        # Store the port number for WebSocket communication
        self.bot = bot_instance # Reference to the main Bot instance
        self.dh = datahandler   # Reference to the DataHandler instance for processing incoming data
        self.loop = asyncio.get_event_loop()    # Get the event loop for managing asynchronous tasks

        self.bp = BetterPrinting(self.bot.i_print , self.bot.d_print, self.bot.e_print)

    #Handles WebSocket communication with clients.
    async def handle_websocket(self, websocket, _):
        try:
            async for message in websocket:
                try:
                    data = json.loads(message)
                    self.bp.debug_print(f"Received WebSocket data: {data}")
                    
                    # Check for a specific key in the data to determine the type
                    data_type = data.get("type")

                    if data_type == "stop_Flag":
                        await self.dh.handle_stopFlag_data(data)
                    elif data_type == "motor":
                        await self.dh.handle_direction_data(data)
                    elif data_type == "buzzer":
                        await self.dh.handle_buzzer_data(data)
                    elif data_type == "lightbar":
                        await self.dh.handle_lightbar_data(data)
                    elif data_type == "laser":
                        await self.dh.handle_laser_data(data)
                    elif data_type == "autonom":
                        await self.dh.handle_autonomous_driving(data)
                    elif data_type == "brakeassist":
                        await self.dh.handle_brake_assistant(data)
                    elif data_type == "detection":
                        await self.dh.handle_detection_system(data)
                    else:
                        self.bp.error_print(f"Unknown data type: {data_type}")
                    
                except json.JSONDecodeError as e:
                    self.bp.error_print(f"Error decoding JSON data: {e}")
        except websockets.exceptions.ConnectionClosed as e:
            self.bp.error_print(f"WebSocket connection closed: {e}")

    #Starts the WebSocket server.
    async def start_websocket_server(self):
        server = await websockets.serve(self.handle_websocket, self.ip, self.port)
        print(f"WebSocket server started on ws://{self.ip}:{self.port}")

    #Starts the WebSocket server and runs the event loop.
    def start(self):
        self.loop.run_until_complete(self.start_websocket_server())
        self.loop.run_forever()