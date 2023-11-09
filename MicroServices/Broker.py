import asyncio
import cv2
import json
import numpy as np
from confluent_kafka import Producer

# Kafka configuration
conf = {
    'bootstrap.servers': 'localhost:9092',
    'client.id': 'camera-lidar-service'
}

producer = Producer(conf)

def delivery_report(err, msg):
    if err is not None:
        #print('Message delivery failed: {}'.format(err))
        return
    else:
        #print('Message delivered to {} [{}]'.format(msg.topic(), msg.partition()))
        return

async def handle_camera_stream(reader, writer):
    while True:
        frame_data = await reader.read(4096)
        # Assume frame_data is in the correct format
        # Produce the frame data to Kafka topic 'camera-topic'
        producer.produce('camera-topic', frame_data, callback=delivery_report)
        producer.poll(0)

async def handle_lidar_stream(reader, writer):
    buffer = ""
    while True:
        data = await reader.read(4096)
        buffer += data.decode('utf-8')
        while '\n' in buffer:
            line, buffer = buffer.split('\n', 1)
            try:
                lidar_data = json.loads(line)
                # Produce the LiDAR data to Kafka topic 'lidar-topic'
                producer.produce('lidar-topic', json.dumps(lidar_data), callback=delivery_report)
                producer.poll(0)
            except json.JSONDecodeError:
                print("JSON Decode Error: Could not parse data")

async def main():
    camera_server = await asyncio.start_server(
        handle_camera_stream, '0.0.0.0', 8000
    )
    lidar_server = await asyncio.start_server(
        handle_lidar_stream, '0.0.0.0', 8001
    )

    await asyncio.gather(
        camera_server.serve_forever(),
        lidar_server.serve_forever()
    )

# Run the main coroutine
asyncio.run(main())
