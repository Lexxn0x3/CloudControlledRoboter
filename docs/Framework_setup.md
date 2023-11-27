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
   - [VScode](#vscode)
   - [Node.js](#nodejs)
   - [Git](#git)
3. [Create-React-App](#create-react-app)
4. [Using Electron](#using-electron)
   - [Project Initialization](#project-initialization)
   - [Installing Electron](#installing-electron)
5. [Using MUI - Joy UI](#using-mui-joy-ui)
   - [Installation](#installation)
6. [Installing Extras](#installing-extras)
   - [VScode Extensions](#vscode-extensions)
   - [Font Installation](#font-installation)
7. [React with Electron](#react-with-electron)
   - [Setting up React](#setting-up-react)
   - [Integration Tutorial](#integration-tutorial)

## Introduction

This documentation outlines the steps to set up a project using Create React App with Electron and MUI - Joy UI.

## Needed Tools

Make sure you have the following tools installed:

- **VScode**
- **Node.js**
  - Install Node.js using the command: `winget install -e --id OpenJS.NodeJS`
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
4. Create a CRACO Configuration File
    ```craco.config.js
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
    ```
5. Install webpack-node-externals
    ```bash
    npm install webpack-node-externals --save-dev
    ```
6. Install Electron
    ```bash
    npm install electron --save-dev
    ```
7. Create Your Electron Main Process File
    -   Add the following code to a new file called electron.js in the public directory:
        ```electron.js
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
        ```

    -   Add the following to your package.json file:
        ```package.json
        "main": "public/electron.js",
        "homepage": "./",
        ```
        This is the entry point for the Electron App

    -   Add custom start and build scripts
        ```package.json
        "scripts": {
        "build": "craco build",
        "start": "electron ."
        },
        ```
8. Run your App
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
2. Download the required font.

Additional installations:

```bash
npm install react react-dom
npm install react-scripts
npm install electron-builder cross-env --save-dev
```

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






