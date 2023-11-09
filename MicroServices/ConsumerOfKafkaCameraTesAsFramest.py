import cv2
import numpy as np
from confluent_kafka import Consumer, KafkaException, KafkaError

# Kafka configuration for consumer
conf = {
    'bootstrap.servers': 'localhost:9092',
    'group.id': 'camera-consumer-group',
    'auto.offset.reset': 'earliest'
}

consumer = Consumer(conf)
consumer.subscribe(['camera-topic'])

def save_frame(frame_data, frame_number):
    """Save frame as JPEG file"""
    filename = f'frame_{frame_number:08d}.jpg'
    with open(filename, 'wb') as f:
        f.write(frame_data)
    print(f"Saved {filename}")

def main():
    try:
        frame_number = 0
        buffer = bytearray()

        while True:
            msg = consumer.poll(1.0)
            if msg is None:
                continue
            if msg.error():
                if msg.error().code() == KafkaError._PARTITION_EOF:
                    # End of partition event
                    continue
                else:
                    print(msg.error())
                    break

            # Append the message to the buffer
            buffer.extend(msg.value())

            # Check for the JPEG frame end marker (0xff 0xd9)
            while True:
                start = buffer.find(b'\xff\xd8')
                end = buffer.find(b'\xff\xd9', start)
                if start != -1 and end != -1:
                    # Extract the JPEG frame
                    frame_data = buffer[start:end + 2]
                    save_frame(frame_data, frame_number)
                    frame_number += 1
                    # Remove the processed frame from the buffer
                    buffer = buffer[end + 2:]
                else:
                    break

    except KeyboardInterrupt:
        pass
    except KafkaException as e:
        print(f"Kafka exception: {e}")
    finally:
        # Close down consumer to commit final offsets.
        consumer.close()

if __name__ == "__main__":
    main()
