import json
import sys
from threading import Thread
from dataHandler import DataHandler
from tcpController import TCPController
from websocketController import WebSocketController
from lidarDistance.LidarDistanceMain import LidarDistanceSystem
from detect.detectObject import RobotController


# The main controller class for the robot.
# This module defines the `Bot` class, which acts as the main controller for a robot. It handles TCP communication,
# data processing, and motor control through various methods.

class Bot():
    
    def __init__(self, tcp_ip, tcp_port, web_ip, web_port, info_print, debug_print, error_print):
        print('starting AppControlHandler...')
        self.i_print = info_print   #Flag to enable/disable informational print statements.
        self.d_print = debug_print  #Flag to enable/disable debugging print statements.
        self.e_print = error_print  #Flag to enable/disable error print statements.
        self.startLidarDistanceSystem(web_ip, web_port)  # Start LidarDistanceSystem in a separate thread
        self.startTCPConection(tcp_ip, tcp_port)
        self.startDataHandler()
        self.startWebsocket(web_ip, web_port)
        
        print('------------------------------\nAppControlHandler ready\n------------------------------\n')
   

    ###################
    #     tcp con     #
    ###################
    
    def startTCPConection(self, ip, port):
        print("starting TCP Connection...")
        self.tcpc = TCPController(robot_host=ip, robot_port=port, bot_instance=self)
        
    
    #################
    #  dataHandler  #
    #################

    def startDataHandler(self):
        print("starting Datahandler...")
        self.dh = DataHandler(self, self.lidar_thread)
    
    
    ###################
    #  websocket con  #
    ###################
        
    def startWebsocket(self, web_ip, web_port):
        print("starting Websocket Connection...")
        self.websocket_controller = WebSocketController(web_ip, web_port, datahandler=self.dh, bot_instance=self)
        self.websocket_controller.start()

    ##########################
    #  LidarDistance Thread  #
    ##########################
        
    def startLidarDistanceSystem(self, web_ip, web_port):
        # Start LidarDistanceSystem in a separate thread
        self.lidar_thread = Thread(target=self.start_lidar_distance_system, args=(web_ip, web_port))
        self.lidar_thread.start()

    
    def start_lidar_distance_system(self, web_ip, web_port):
        self.LDC = LidarDistanceSystem("192.168.8.20", 9011, web_ip, web_port, "192.168.8.20", 3031)
    

    def stopLidarDistanceSystem(self):
        if self.lidar_thread is not None and self.lidar_thread.is_alive():
            self.LDC.stop_threads = True  # Set the stop_threads flag to True
            self.lidar_thread.join(timeout=2)  # Wait for LidarDistanceThread to finish for up to 2 seconds
            if self.lidar_thread.is_alive():
                # If the thread is still alive after the timeout, consider terminating it forcefully
                self.lidar_thread._stop()  # Note: Using _stop is generally not recommended, but it forcefully terminates the thread
        

    #################
    #     stop      #
    #################
    
    def stop(self):
        self.send_motor_data(0, 0, 0, 0)


    ##################
    #     drive      #
    ##################
    
    #speed goes from -100 up to 100
    
    #drive forwards
    def drive_forward(self, speed):                  
        self.send_motor_data(speed, speed, speed, speed)

    #drive right forward
    def drive_right_forward(self, speed):
        self.send_motor_data(speed, 0, 0, speed)
    
    #drive left forward
    def drive_left_forward(self, speed):
        self.send_motor_data(0, speed, speed, 0)
    
    #drive backwards
    def drive_backward(self, speed):                 
        self.send_motor_data(-1*speed, -1*speed, -1*speed, -1*speed)

    #drive right backward
    def drive_right_backward(self, speed):
        self.send_motor_data(0, -1*speed, -1*speed, 0)
        
    #drive left backward
    def drive_left_backward(self, speed):
        self.send_motor_data(-1*speed, 0, 0, -1*speed)
        
    #drive right
    def drive_right(self, speed):                     
        self.send_motor_data(speed, -1*speed, -1*speed, speed)

    #drive left
    def drive_left(self, speed):                      
        self.send_motor_data(-1*speed, speed, speed, -1*speed)
        
        
    ##################
    #      spin      #
    ##################
    
    #spin right
    def spin_right(self, speed):
        self.send_motor_data(0.5*speed, 0.5*speed, -0.5*speed, -0.5*speed)

    #spin left
    def spin_left(self, speed):
        self.send_motor_data(-0.5*speed, -0.5*speed, 0.5*speed, 0.5*speed)


    ##################
    #   drive curve  #
    ##################

    #drive curve right forward
    def drive_curve_right_forward(self, speed):
        self.send_motor_data(speed, speed, 0.5*speed, 0.5*speed)

    #drive curve left forward
    def drive_curve_left_forward(self, speed):
        self.send_motor_data(0.5*speed, 0.5*speed, speed, speed)

    #drive curve right backward
    def drive_curve_right_backward(self, speed):
        self.send_motor_data(-1*speed, -1*speed, -0.5*speed, -0.5*speed)

    #drive curve left backward
    def drive_curve_left_backward(self, speed):
        self.send_motor_data(-0.5*speed, -0.5*speed, -1*speed, -1*speed)


    ##################
    #  motor_control #
    ##################
    
    #send the given motor instructions to the robot
    def send_motor_data(self, motor1, motor2, motor3, motor4):
        try:
            #while True:
                # Prepare the JSON data for motor control
            motor_data = {
                "motor1": motor1,
                "motor2": motor2,
                "motor3": motor3,
                "motor4": motor4
            }

            # Convert motor data to JSON string
            json_data_str = json.dumps(motor_data)

            #time.sleep(4)
            # Send the JSON data to the robot
            self.tcpc.send_json_data(json_data_str)

        except Exception as e:
            self.pb.error_print(f"Error sending motor data to the robot: {e}")
        

    #################
    #    lightbar   #
    #################

    def send_lightbar_data(self, isEffect, red, green, blue, effect, speed):
        try:
            # Prepare the JSON data for motor control
            motor_data = {
                "mode": isEffect,    #true = effect, false = single light
                "ledid": 0xff,
                "red": red,
                "green": green,
                "blue": blue,
                "effect": effect,
                "speed": speed,      
                "parm": 255
            }

            # Convert motor data to JSON string
            json_data_str = json.dumps(motor_data)

            # Send the JSON data to the robot
            self.tcpc.send_json_data(json_data_str)

        except Exception as e:
            self.pb.error_print(f"Error sending lightbar data to the robot: {e}")
            
            
    #################
    #     buzzer    #
    #################
    
    def send_buzzer_data(self, onBuzzer):
        try:
            buzzer_data = {
                "buzzer": int(onBuzzer)
            }
            
            json_data_str = json.dumps(buzzer_data)
            
            self.tcpc.send_json_data(json_data_str)
            
        except Exception as e:
            self.pb.error_print(f"Error sending buzzer data to the robot: {e}")
      
            
    #################
    #     laser     #
    #################
    
    def send_laser_data(self, onLaser):
        try:
            laser_data = {
                "laser": onLaser
            }
            
            json_data_str = json.dumps(laser_data)
            
            self.tcpc.send_json_data(json_data_str)
            
        except Exception as e:
            self.pb.error_print(f"Error sending laser data to the robot: {e}")
            
    
    #################
    #   detection   #
    #################
    
    def send_detection_data(self, onDetection):
        try:
            if onDetection:
                self.robot_controller = RobotController()
                self.robot_controller.run()
            else:
                self.robot_controller = None
                
        except Exception as e:
            self.pb.error_print(f"Error starting object detection: {e}")
    
    


if __name__ == "__main__":
    # Check if the correct number of command-line arguments is provided
    if len(sys.argv) < 5:
        print("Usage: python your_script.py tcp_ip tcp_port web_ip web_port [info_print] [debug_print] [error_print]")
        sys.exit(1)

    # Extract variables from command-line arguments
    tcp_ip, tcp_port, web_ip, web_port = sys.argv[1:5]

    # Extract optional boolean variables with default values set to True
    info_print = sys.argv[5].lower() == 'true' if len(sys.argv) > 5 else True
    debug_print = sys.argv[6].lower() == 'true' if len(sys.argv) > 6 else True
    error_print = sys.argv[7].lower() == 'true' if len(sys.argv) > 7 else True

    # Create an instance of the bot class with the provided variables
    my_bot = Bot(tcp_ip, tcp_port, web_ip, web_port, info_print, debug_print, error_print)
    
