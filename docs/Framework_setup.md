---
title: Framework Setup
layout: default
parent: App
nav_order: 1
---

# Framework Setup

## Table of Contents
1. [Introduction](#introduction)
2. [Needed Tools](#needed-tools)
3. [Create-React-App](#create-react-app)
4. [Using Electron](#using-electron)
   - [Project Initialization](#project-initialization)
   - [Run yor App](#run-your-app)
5. [Using MUI - Joy UI](#using-mui---joy-ui)
6. [Installing Extras](#installing-extras)
   - [VScode Extensions](#vscode-extensions)
   - [Font Installation](#font-installation)
   - [Electron hot-reload](#electron-hot-reload)

## Introduction

This documentation outlines the steps to set up a project using Create React App with Electron and MUI - Joy UI.

## Needed Tools

Make sure you have the following tools installed:

- **VScode**
- **Node.js**
  - Install Node.js using the command:

      ```bash
       winget install -e --id OpenJS.NodeJS
      ```
    
  - Verify installation with:
  
    ```bash
    $ node -v
    $ npm -v
    ```
    
- **Git**

*LAR (Local Admin Rights) is recommended.*

## Create-React-App

### Initiating Your App

Initiate your app by running the following commands in your project folder:

```bash
npm init react-app my-app
cd my-app
```


Refer to the [Create React App documentation](https://create-react-app.dev/docs/getting-started/). for more instructions.

## Using Electron

### Project Initialization

All of the following commands are run within the my-app folder

1. Create Your React Application
2. Remove the web-vitals dependency
    -   Uninstall the web-vitals package
      
        ```bash
        npm uninstall web-vitals
        ```
        
    -   Delete the reportWebVitals.js file:
      
        ```bash
        rm src/reportWebVitals.js
        ```
        
    -   Remove the following lines from the src/index.js file:
      
        ```bash
        import reportWebVitals from './reportWebVitals';
        reportWebVitals();
        ```
        
3. Install CRACO to Alter Your Webpack Configuration

    ```bash
    npm install @craco/craco
    ```
    
5. Create a CRACO Configuration File named craco.config.js

   
       const nodeExternals = require("webpack-node-externals");
   
       module.exports = {
       webpack: {
           configure: {
           target: "electron-renderer",
           externals: [
               nodeExternals({
               allowlist: [/webpack(\/.*)?/, "electron-devtools-installer"],
               }),
           ],
           },
       },
       };
    
6. Install webpack-node-externals

    ```bash
    npm install webpack-node-externals --save-dev
    ```
   
8. Install Electron

    ```bash
    npm install electron --save-dev
    ```
    
10. Create Your Electron Main Process File
    -   Add the following code to a new file called electron.js in the public directory:
    
           
              const electron = require("electron");
              const path = require("path");
      
              const app = electron.app;
              const BrowserWindow = electron.BrowserWindow;
      
              let mainWindow;
      
              function createWindow() {
              // Create the browser window.
              mainWindow = new BrowserWindow({
                  width: 800,
                  height: 600,
                  webPreferences: { nodeIntegration: true, contextIsolation: false },
              });
              // and load the index.html of the app.
              console.log(__dirname);
              mainWindow.loadFile(path.join(__dirname, "../build/index.html"));
              }
      
              // This method will be called when Electron has finished
              // initialization and is ready to create browser windows.
              // Some APIs can only be used after this event occurs.
              app.on("ready", createWindow);
              
        

    -   Add the following to your package.json file:
        
              "main": "public/electron.js",
              "homepage": "./",
        
        
        This is the entry point for the Electron App

    -   Add custom start and build scripts
      
              "scripts": {
              "build": "craco build",
              "start": "electron ."
              },
        
### Run your App

To run your app, use te following commands
 
 ```bash
 npm run build
 npm run start
 ```
 
 You might want to add a custom script calling both commands:

 ```package.json
 "scripts": {
     "build-start": "craco build && electron ."
     "build": "craco build",
     "start": "electron ."
     },
 ```
 
 run this command with:
 
 ```bash 
 npm run build-start
 ```


Full tutorial on [MongoDB Realm SDK for Electron + CRA integration](https://www.mongodb.com/docs/realm/sdk/node/integrations/electron-cra/).

Refer to the [Electron documentation](https://www.electronjs.org/de/docs/latest/tutorial/tutorial-first-app) for more instructions.

## Using MUI - Joy UI

Follow the [MUI - Joy UI Tutorial](https://mui.com/joy-ui/getting-started/installation/) for installation.

1. Install MUI - Joy UI:

   ```bash
   npm install @mui/joy @emotion/react @emotion/styled
   ```
   
   Don't forget to add dependencies to `package.json`.
3. Download the required font.

## Installing Extras

### VScode Extensions

- Add Node.js Module extensions to VScode.
- Recommended extensions: Git push & pull and GitLens.
- Disable all deprecated extensions (uninstall if needed).

### Font Installation

MUI & Joy UI use inter as default Font

```bash
npm install @fontsource-variable/inter
```

### Electron hot-reload

install:

```bash
npm install -D electron-reloader
```

Add this line to the electron.js file:

```bash
try {
  require('electron-reloader')(module)
} catch (_) {}
```

Refer to the [Flavicopes Tutorial](https://flaviocopes.com/electron-hot-reload/) for more information.






