import cv2
from flask import Flask, Response

app = Flask(__name__)
cap = cv2.VideoCapture(0, cv2.CAP_V4L2)

# Aktivieren der Hardwarebeschleunigung (sofern verfügbar)
if cap.isOpened():
    print("Hardwarebeschleunigung aktiviert")
    cap.set(cv2.CAP_PROP_FPS, 30)  # Setzen Sie die gewünschte Bildrate
    cap.set(cv2.CAP_PROP_FRAME_WIDTH, 1920)  # Setzen Sie die Breite auf die gewünschte Auflösung
    cap.set(cv2.CAP_PROP_FRAME_HEIGHT, 1080)  # Setzen Sie die Höhe auf die gewünschte Auflösung
    cap.set(cv2.CAP_PROP_BUFFERSIZE, 1)  # Begrenzen Sie den Frame-Buffer auf 1
    cap.set(cv2.CAP_PROP_FOURCC, cv2.VideoWriter_fourcc(*"MJPG"))

def generate_frames():
    while True:
        success, frame = cap.read()
        if not success:
            break
        else:
            _, buffer = cv2.imencode('.jpg', frame, [cv2.IMWRITE_JPEG_QUALITY, 30])  # Reduzieren Sie die Qualität auf 30
            frame = buffer.tobytes()
            yield (b'--frame\r\n'
                   b'Content-Type: image/jpeg\r\n\r\n' + frame + b'\r\n')

@app.route('/video_feed')
def video_feed():
    return Response(generate_frames(), mimetype='multipart/x-mixed-replace; boundary=frame')

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)

