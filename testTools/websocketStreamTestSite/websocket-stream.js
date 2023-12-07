document.addEventListener('DOMContentLoaded', function () {
    const streamElement = document.getElementById('stream');
    const connectButton = document.getElementById('connect-button');
    let ws; // Declare WebSocket variable
    let buffer = [];
    let frameStart = -1;

    function connectWebSocket(websocketUrl) 
    {
        if (ws) {
            ws.close();
        }

        ws = new WebSocket(websocketUrl);
        ws.binaryType = 'arraybuffer';

        ws.onmessage = function (event) {
            const data = new Uint8Array(event.data);
            buffer = buffer.concat(Array.from(data));
    
            // Search for JPEG start and end markers
            for (let i = 0; i < buffer.length - 1; i++) {
                if (buffer[i] === 0xFF && buffer[i + 1] === 0xD8) { // JPEG start marker (FFD8)
                    frameStart = i;
                }
                if (buffer[i] === 0xFF && buffer[i + 1] === 0xD9 && frameStart >= 0) { // JPEG end marker (FFD9)
                    const frame = buffer.slice(frameStart, i + 2);
                    const blob = new Blob([new Uint8Array(frame)], { type: 'image/jpeg' });
                    const url = URL.createObjectURL(blob);
                    streamElement.src = url;
    
                    // Force the browser to recognize the new image source
                    streamElement.onload = () => URL.revokeObjectURL(url);
    
                    buffer = buffer.slice(i + 2); // Clear the buffer up to the end of the frame
                    frameStart = -1;
                    break;
                }
            }
        };

        ws.onerror = function (error) {
            console.error('WebSocket Error:', error);
        };

        ws.onclose = function () {
            console.log('WebSocket connection closed. Attempting to reconnect...');
            setTimeout(() => connectWebSocket(websocketUrl), 3000);
        };
    }

    connectButton.addEventListener('click', function() {
        const addressInput = document.getElementById('websocket-address').value;
        if (!addressInput) {
            alert('Please enter a WebSocket IP and port.');
            return;
        }
        const websocketUrl = `ws://${addressInput}`;
        connectWebSocket(websocketUrl);
    });
        

    //connectWebSocket(); // Initial connection

// ... (rest of your code) ...
});