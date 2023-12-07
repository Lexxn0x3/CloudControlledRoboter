import asyncio
import websockets

async def save_jpeg_frames(uri):
    buffer = bytearray()
    frame_count = 0

    async with websockets.connect(uri) as websocket:
        async for message in websocket:
            buffer.extend(message)

            # Find the start and end of the JPEG frame
            start = buffer.find(b'\xff\xd8')
            end = buffer.find(b'\xff\xd9', start + 2)

            # If a complete frame is found
            while start != -1 and end != -1 and start < end:
                frame = buffer[start:end + 2]
                frame_count += 1
                with open(f"frame_{frame_count}.jpg", "wb") as file:
                    file.write(frame)

                # Remove the processed frame from the buffer
                buffer = buffer[end + 2:]

                # Look for the next frame
                start = buffer.find(b'\xff\xd8')
                end = buffer.find(b'\xff\xd9', start + 2)

# Example usage
asyncio.run(save_jpeg_frames("ws://192.168.8.103:5001"))
