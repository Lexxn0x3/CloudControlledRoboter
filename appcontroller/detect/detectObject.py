from ultralytics import YOLO
import time


class RobotController:
    def __init__(self, bot_instance ,model_path='werner.pt', ip_address='127.0.0.1', port=9001):
        self.bot = bot_instance 
        self.model = YOLO(model_path)
        self.address = f'tcp://{ip_address}:{port}'
        

    def calculate_midpoint(self, points):
        points = points.squeeze()
        x = int((points[1] + points[3]) / 2)
        y = int((points[0] + points[2]) / 2)
        return x, y

    
    def adjust_robot_position(self, object_midpoint, target_position, tolerance):
        object_y, object_x = object_midpoint
        target_x, target_y = target_position
        tolerance_x, tolerance_y = tolerance
        final_x = False
        final_y = False

        if object_x >= abs(target_x - tolerance_x)  and object_x <= abs(target_x + tolerance_x) and not final_x:
            self.bot.stop()
            print("stop")
            final_x = True

        if object_x < abs(target_x - tolerance_x)  or object_x > abs(target_x + tolerance_x):
            final_x = False  # Reset final_x if the x-coordinate is outside tolerance
            if object_x < (target_x - tolerance_x):
                self.bot.spin_left(30)
                print("spin left")
            else:
                self.bot.spin_right(30)
                print("spin right")
        
        if object_y >= abs(target_y - tolerance_y)  and object_y <= abs(target_y + tolerance_y) and final_x is True and not final_y:
            self.bot.stop()
            print("stop")
            final_y = True

        if object_y < abs(target_y - tolerance_y)  or object_y > abs(target_y + tolerance_y) and final_x is True:
            final_y = False  # Reset final_y if the y-coordinate is outside tolerance
            if object_y < (target_y - tolerance_y):
                self.bot.drive_forward(30)
                print("forward")
            else:
                self.bot.drive_backward(30)
                print("backward")
                
        if final_x and final_y:
            print("ready to shoot")
            self.bot.send_laser_data(True)
            time.sleep(3)
            self.bot.send_laser_data(False)
            

    def process_results(self, results):
        for result in results:
            #print(result.boxes[0].xyxy)
            #print(result.boxes[0].orig_shape[0])
            #print(result.boxes[0].orig_shape[1])
            if result.boxes and result.boxes[0] and result.boxes[0].xyxy is not None:
                target_x = int(result.boxes[0].orig_shape[1] * 50 / 100)
                target_y = int(result.boxes[0].orig_shape[0] * 70 / 100)

                tolerance_x = int(result.boxes[0].orig_shape[1] * 2 / 100)
                tolerance_y = int(result.boxes[0].orig_shape[0] * 2 / 100)

                object_midpoint = self.calculate_midpoint(result.boxes[0].xyxy)
                print(object_midpoint)
                target_position = (target_x, target_y)
                tolerance = (tolerance_x, tolerance_y)

                self.adjust_robot_position(object_midpoint, target_position, tolerance)

    def run(self):
        results = self.model.predict(self.address, classes=[0], conf=0.6, stream=True, save=False, show=False, save_txt=False)
        self.process_results(results)




