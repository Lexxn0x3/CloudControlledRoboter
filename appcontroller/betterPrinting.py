# This module defines the `BetterPrinting` class, which provides enhanced printing functionalities
# with colored output for information, debugging, and error messages.

class BetterPrinting():
    class Colors:
        RESET = '\033[0m'
        RED = '\033[91m'
        GREEN = '\033[92m'
        YELLOW = '\033[93m'

    def __init__(self, onInfo=True, onDebug=True, onError=True):
        # Activate/deactivate debugging
        self.debug = onDebug
        self.info = onInfo
        self.error = onError

    def info_print(self, message):
        if self.info:
            print(f"{self.Colors.GREEN}[INFO]:{self.Colors.RESET} {message}")

    def debug_print(self, message):
        if self.debug:
            print(f"{self.Colors.YELLOW}[DEBUG]:{self.Colors.RESET} {message}")

    def error_print(self, message):
        if self.error:
            print(f"{self.Colors.RED}[ERROR]:{self.Colors.RESET} {message}")
