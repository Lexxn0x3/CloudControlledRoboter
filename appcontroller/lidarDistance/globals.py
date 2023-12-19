class LidarGlobals:
    def __init__(self):
        self.stop_front = False
        self.stop_right = False
        self.stop_left = False
        self.stop_front_right = False
        self.stop_front_left = False
        self.stop_back = False
        self.stop_threads = False
        self.minDist = 400
        self.maxLenBuffer = 150
        self.speed = 50