from flask import Flask, Response, stream_with_context
from confluent_kafka import Consumer, KafkaException, KafkaError
from threading import Thread
import time

app = Flask(__name__)

# Kafka configuration for consumer
conf = {
    'bootstrap.servers': 'localhost:9092',
    'group.id': 'camera-consumer-group',
    'auto.offset.reset': 'earliest'
}

consumer = Consumer(conf)
consumer.subscribe(['camera-topic'])

buffer = bytearray()  # Shared buffer for Kafka consumer and Flask


def kafka_consumer_thread():
    global buffer
    try:
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

    except KafkaException as e:
        print(f"Kafka exception: {e}")
    finally:
        # Close down consumer to commit final offsets.
        consumer.close()


@app.route('/video')
def video_feed():
    # Define the generator function that will stream the MJPEG frames
    @stream_with_context
    def generate():
        global buffer
        while True:
            start = buffer.find(b'\xff\xd8')
            end = buffer.find(b'\xff\xd9', start)
            if start != -1 and end != -1:
                # Extract the JPEG frame
                frame_data = buffer[start:end + 2]
                # Serve the image in multipart/x-mixed-replace format
                yield (b'--frame\r\nContent-Type: image/jpeg\r\n\r\n' + frame_data + b'\r\n')
                # Remove the processed frame from the buffer
                buffer = buffer[end + 2:]
            time.sleep(0.04)  # Wait roughly the duration of one frame at 25 fps

    return Response(generate(), mimetype='multipart/x-mixed-replace; boundary=frame')


if __name__ == "__main__":
    # Start the Kafka consumer thread
    Thread(target=kafka_consumer_thread).start()
    # Start the Flask app
    app.run(host='0.0.0.0', port=5000, threaded=True)
