---
title: WebsocketClient
layout: default
parent: App
nav_order: 2
---

# WebsocketClient

## Table of Contents

1. [Introduction](#introduction)
2. [Getting Started](#getting-started)
   - [Prerequisites](#prerequisites)
   - [Installation](#installation)
3. [Usage](#usage)
   - [Creating a WebSocket Connection](#creating-a-websocket-connection)
   - [Sending Messages](#sending-messages)
   - [Receiving Messages](#receiving-messages)
   - [Closing the Connection](#closing-the-connection)
4. [Optional Features](#optional-features)
4. [Troubleshooting](#troubleshooting)
5. [Examples](#examples)
6. [References](#references)

## Introduction

WebSocket is a communication protocol that enables real-time bidirectional communication between a client and a server. This documentation provides guidelines for implementing and using WebSocket on the client side.

## Getting Started

### Prerequisites

- Basic knowledge of web development (HTML, CSS, JavaScript).
- A server environment with WebSocket support.

### Installation

No specific installation is required for WebSocket on the client side. It is natively supported in modern web browsers.

## Usage

### Creating a WebSocket Connection



      // Create a WebSocket connection
      const socket = new WebSocket('ws://127.0.0.1:5002');

      // Event listener for when the connection is open
      socket.onopen = () => {
          console.log('WebSocket connection opened');
      };


### Error handling


       //Event listener for when error occurs
       socket.onerror = (error) => {
      console.error('Camera WebSocket error:', error);
    };


### Sending Messages

   
      // Send a message to the server
      const messageToSend = 'Hello, server!';
      socket.send(messageToSend);


### Receiving Messages

      
      // Event listener for when a message is received
      socket.onmessage = (event) => {
          const receivedMessage = event.data;
          console.log('Server says: ${receivedMessage}');
      };
      

### Closing the Connection

      
      // Close the WebSocket connection
      socket.close();
            

## Optional Features

These features could prove useful

### Reconnecting Websockets

Install this [library](https://github.com/pladaria/reconnecting-websocket)

      npm install --save reconnecting-websocket

import it to the relevant file 

      import ReconnectingWebSocket from 'reconnecting-websocket';

and replace the method call with

      const socket = new ReconnectingWebSocket(wsServerUrl);

## Troubleshooting

Check the browser console for errors.
Use browser developer tools to inspect WebSocket connections.

## Examples


      // Example: Connect to a WebSocket server and send/receive messages
      const socket = new WebSocket('ws://127.0.0.1:5002');
      
      socket.onopen = () => {
          console.log('WebSocket connection opened');
      
          // Send a message to the server
          const messageToSend = 'Hello, server!';
          socket.send(messageToSend);
      };
      
      socket.onmessage = (event) => {
          const receivedMessage = event.data;
          console.log(`Server says: ${receivedMessage}`);
      };

      //Event listener for when error occurs
      socket.onerror = (error) => {
         console.error('Camera WebSocket error:', error);
      };
      
      // Close the WebSocket connection after 5 seconds
      setTimeout(() => {
          socket.close();
      }, 5000);



## References

- [WebSocket API MDN Web Docs](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)
