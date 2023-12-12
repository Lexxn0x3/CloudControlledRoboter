---
title: Better Printing
layout: default
parent: App Controller
grand_parent: MicroServices
nav_order: 5
---

# BetterPrinting

The `BetterPrinting` module, contained in the `BetterPrinting.py` file, is a utility module that enhances console printing by adding color-coding to different types of messages.

## Class: `BetterPrinting`

### Colors

- **RESET**: Resets the text color to the default.
- **RED**: Sets the text color to red for error messages.
- **GREEN**: Sets the text color to green for informational messages.
- **YELLOW**: Sets the text color to yellow for debugging messages.

### Methods

#### `__init__(self, onInfo=True, onDebug=True, onError=True)`

Initializes the `BetterPrinting` instance.

- **Parameters**:
  - `onInfo` (optional, default is `True`): Boolean flag indicating whether informational messages should be printed.
  - `onDebug` (optional, default is `True`): Boolean flag indicating whether debugging messages should be printed.
  - `onError` (optional, default is `True`): Boolean flag indicating whether error messages should be printed.

#### `info_print(self, message)`

Prints an informational message to the console with green color-coding.

- **Parameters**:
  - `message`: The message to be printed.

#### `debug_print(self, message)`

Prints a debugging message to the console with yellow color-coding.

- **Parameters**:
  - `message`: The message to be printed.

#### `error_print(self, message)`

Prints an error message to the console with red color-coding.

- **Parameters**:
  - `message`: The message to be printed.

### Usage

To use the `BetterPrinting` class in your code, create an instance and use its methods for improved console printing. Example:

```python
# Create an instance with default settings
bp = BetterPrinting()

# Print an informational message
bp.info_print("This is an informational message")

# Print a debugging message
bp.debug_print("This is a debugging message")

# Print an error message
bp.error_print("This is an error message")
```

## Note

The `BetterPrinting` class provides a simple and effective way to enhance console output by adding color-coded messages. It can be integrated into other modules, such as the `DataHandler`, `TCPController`, and `WebSocketController` classes, to improve the readability of log messages.