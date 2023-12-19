---
title: XboxController Integration
layout: default
parent: App
nav_order: 4
---


# Xbox Controller Integration

## Introduction

This documentation provides guidance on integrating an Xbox controller into a React app. The integration includes handling joystick inputs, button presses, and updating application state accordingly.

### Prerequisites

- Xbox controller connected to your computer.
- An electron project integrating react

## Getting Started

### Installation

You shouldn't need to install anything if you set up the project correctly.

### Usage

1. **Create a Gamepad Component:**

   Create a React component (e.g., `GamepadComponent.js`) to handle Xbox controller events:

   ```jsx
   // GamepadComponent.js
   import React, { useEffect, useState } from 'react';

   const GamepadComponent = () => {
     // State to store the index of the connected controller
     const [gamepadIndex, setGamepadIndex] = useState(-1);

     // Handle connection and disconnection of the controller
     const handleGamepads = () => {
       const controller = navigator.getGamepads().findIndex(gamepad => gamepad != null);
       if(controller >= 0) {
         setGamepadIndex(controller);
         handleControllerButtons();
         handleControlerAxes();
       }
       else{
         setGamepadIndex(null);
       }
     };

     // Handle joystick inputs
     const handleControlerAxes = () => {
       // Logic to handle joystick inputs and update application state
       // ...
     };

     // Handle button presses
     const handleControllerButtons = () => {
       // Logic to handle button presses and update application state
       // ...
     };

      //Set up an interval to check the inputs
      useEffect(() =>{
        const interval = setInterval(handleGamepads, 150);
    
        return () => {
          clearInterval(interval);
        }
      });

     return (
       <div>
         {/* Display information or controls based on gamepad state */}
       </div>
     );
   };

   export default GamepadComponent;
   ```


## Features

### Joystick Inputs

- The app detects joystick movements and calculates direction and speed.
- The `handleControlerAxes` function updates the application state accordingly.

### Button Presses

- The app handles button presses on the Xbox controller.
- The `handleControllerButtons` function updates the application state based on button presses.

### Additional Features

- Customize button actions based on your application requirements.

## Tips and tricks

To see the gamepad objects use this command in the developer tools console


      navigator.getGamepads()


Here you can learn what inputs you can use
## References
[Youtube Video I learned from](https://www.youtube.com/watch?v=UPaKoTfqk8k)
