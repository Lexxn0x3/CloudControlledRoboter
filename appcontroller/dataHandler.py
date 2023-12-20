from threading import Thread
from betterPrinting import BetterPrinting
from lidarDistance.LidarDistanceMain import LidarDistanceSystem
from detect.detectObject import RobotController


# Class responsible for processing incoming data and translating it into robot control commands.

class DataHandler():
    # Flags to indicate obstacles in different directions and stop robot movement if obstacles are detected.
    stopFront = False        # Flag to stop if there is an obstacle in front of the robot
    stopLeft = False         # Flag to stop if there is an obstacle on the left
    stopRight = False        # Flag to stop if there is an obstacle on the right
    stopFrontRight = False   # Flag to stop if there is an obstacle in the front-right
    stopFrontLeft = False    # Flag to stop if there is an obstacle in the front-left
    stopBack = False
        
    def __init__(self, bot_instance, lidar_thread):
        self.bot = bot_instance
        self.lidar_thread = lidar_thread
        self.bp = BetterPrinting(self.bot.i_print , self.bot.d_print, self.bot.e_print)

    #################
    # calc keyInput #
    #################
    
    #Process break assistant
    async def handle_brake_assistant(self, data):
        state = data.get("state", False)
        if not state:
            self.bot.stopLidarDistanceSystem()  # Stop LidarDistanceSystem gracefully
        else:
            if self.lidar_thread is not None and not self.lidar_thread.is_alive():
                self.bot.startLidarDistanceSystem(self.bot.web_ip, self.bot.web_port)  # Start LidarDistanceSystem    
    
    #Process detection system
    async def handle_detection_system(self, data):
        state = data.get("detect", False)
        if state:
                self.robot_controller = RobotController(self.bot)
                self.robot_controller.run()
        else:
            self.robot_controller = None
        
    
    
    # Processes incoming stop flag data and stops the robot if any stop condition is met.
    async def handle_stopFlag_data(self, data):
        self.stopFront = data.get("stop_front", False)
        self.stopLeft = data.get("stop_left", False)
        self.stopRight = data.get("stop_right", False)
        self.stopFrontLeft = data.get("stop_front_left", False)
        self.stopFrontRight = data.get("stop_front_right", False)

        if any([self.stopFront, self.stopFrontLeft, self.stopFrontRight, self.stopLeft, self.stopRight]) == True:
            self.bot.stop()


    #Processes incoming buzzer data and sends the corresponding command to the robot.
    async def handle_buzzer_data(self, data):
        state = data.get("state", False)
        
        self.bot.send_buzzer_data(state)


    # Processes incoming laser data
    async def handle_laser_data(self, data):
        state = data.get("state", False)
        
        self.bot.send_laser_data(state)
    
    
    #Processes incoming lightbar data and controls the lightbar on the robot.
    async def handle_lightbar_data(self, data):
        on = data.get("on", False)
        lb_speed = data.get("speed", 10)
        isEffect = data.get("isEffect", False)
        effect = data.get("effect", 0)
        red = data.get("r", 0)
        green = data.get("g", 0)
        blue = data.get("b", 255)
       
        if on:
            self.bp.info_print("light on")
            self.bot.send_lightbar_data(isEffect, red, green, blue, effect, lb_speed)
        else:
            self.bp.info_print("light off")
            self.bot.send_lightbar_data(False, 0, 0, 0, 0, 1)
    
    
    #Processes incoming direction data and translates it into corresponding robot control commands.
    async def handle_direction_data(self, data):
        w = data.get("w", False)
        a = data.get("a", False)
        s = data.get("s", False)
        d = data.get("d", False)
        q = data.get("q", False)
        e = data.get("e", False)
        speed = data.get("speed", 0)

        # Translate direction data into corresponding moving direction functions
        if not any([w, a, s, d, q, e]):
            # All keys are False
            self.bp.info_print("Stopping")
            self.bot.stop()
        elif w and not any([a, s, d, q, e, self.stopFront]):
            # Only w is True
            self.bp.info_print("Driving forward")
            self.bot.drive_forward(speed)
        elif a and not any([w, s, d, q, e, self.stopLeft]):
            # Only a is True
            self.bp.info_print("Driving left")
            self.bot.drive_left(speed)
        elif s and not any([w, a, d, q, e]):
            # Only s is True
            self.bp.info_print("Driving backward")
            self.bot.drive_backward(speed)
        elif d and not any([w, a, s, q, e, self.stopRight]):
            # Only d is True
            self.bp.info_print("Driving right")
            self.bot.drive_right(speed)
        elif q and not any([w, a, s, d, e]):
            # Only q is True
            self.bp.info_print("Spinning left")
            self.bot.spin_left(speed)
        elif e and not any([w, a, s, d, q]):
            # Only e is True
            self.bp.info_print("Spinning right")
            self.bot.spin_right(speed)
        elif w and d and not any([a, s, q, e, self.stopFrontRight]):
            # w and d are True
            self.bp.info_print("Driving right forward")
            self.bot.drive_right_forward(speed)
        elif w and a and not any([d, s, q, e, self.stopFrontLeft]):
            # w and a are True
            self.bp.info_print("Driving left forward")
            self.bot.drive_left_forward(speed)
        elif s and d and not any([w, a, q, e]):
            # s and d are True
            self.bp.info_print("Driving right backward")
            self.bot.drive_right_backward(speed)
        elif s and a and not any([w, d, q, e]):
            # s and a are True
            self.bp.info_print("Driving left backward")
            self.bot.drive_left_backward(speed)
        elif w and q and not any([a, s, d, e, self.stopFront, self.stopLeft, self.stopFrontLeft]):
            # w and q are True
            self.bp.info_print("Driving left-forward curve")
            self.bot.drive_curve_left_forward(speed)
        elif w and e and not any([a, s, d, q, self.stopFront, self.stopRight, self.stopFrontRight]):
            # w and e are True
            self.bp.info_print("Driving right-forward curve")
            self.bot.drive_curve_right_forward(speed)
        elif s and q and not any([w, a, d, e]):
            # s and q are True
            self.bp.info_print("Driving left-backward curve")
            self.bot.drive_curve_left_backward(speed)
        elif s and e and not any([q, w, a, d]):
            # s and e are True
            self.bp.info_print("Driving right-backward curve")
            self.bot.drive_curve_right_backward(speed)
        else:
            # else just stop
            self.bp.info_print("key combination not valid - stop")
            self.bot.stop()
