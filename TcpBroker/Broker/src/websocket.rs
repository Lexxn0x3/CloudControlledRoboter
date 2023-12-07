use std::net::TcpStream;
use std::thread;
use std::io::{Read, BufReader};
use std::sync::{Arc, Mutex};

use ws::{listen, Sender, Message};
use log::{info, debug, error};


pub fn start_websocket(broker_ip: &str, broker_port: u16, buffer_size: usize, websocket_port: u16) {
    debug!("Starting websocket");

    // Shared state among WebSocket connections
    let state = Arc::new(Mutex::new(Vec::new()));
    let clients = Arc::new(Mutex::new(Vec::new()));

    // Clone `clients` Arc for the TCP connection handling thread
    let clients_for_tcp_thread = clients.clone();

    // Clone the broker_ip and broker_port to move into the thread
    let broker_ip = broker_ip.to_string();
    let broker_port = broker_port.to_string();
    let websocket_port = websocket_port.to_string();

    thread::spawn(move || {
        // Construct the connection string using the IP and port
        let tcp_address = format!("{}:{}", broker_ip, broker_port);

        match TcpStream::connect(&tcp_address) {
            Ok(stream) => {
                stream.set_nodelay(true).expect("set_nodelay call failed");
                debug!("Connected to TCP server at {}", tcp_address);
                let mut reader = BufReader::new(stream);
                // Initialize the buffer with the configurable size
                let mut buffer = vec![0; buffer_size];
                loop {
                    match reader.read(&mut buffer) {
                        Ok(size) => {
                            if size == 0 {
                                debug!("End of TCP stream");
                                break;
                            }
                            let mut state = state.lock().unwrap();
                            state.clear();
                            state.extend_from_slice(&buffer[..size]);
                            debug!("Received {} bytes from TCP stream", size);
                            // Send the data to all connected WebSocket clients
                            broadcast_to_clients(&state, &clients_for_tcp_thread);
                        }
                        Err(e) => {
                            error!("Error reading from TCP stream: {}", e);
                            break;
                        }
                    }
                }
            }
            Err(e) => {
                error!("Failed to connect to TCP server at {}: {}", tcp_address, e);
            }
        }
    });


    // Clone `clients` Arc for the WebSocket server
    let clients_for_ws = clients.clone();

    // Construct the address string with the configurable port
    let address = format!("0.0.0.0:{}", websocket_port);

    // Start WebSocket server
    if let Err(e) = listen(&address, move |out| {
        clients_for_ws.lock().unwrap().push(out.clone());
        move |_msg| {
            // Handle incoming messages from WebSocket clients here
            Ok(())
        }
    }) {
        error!("Failed to start WebSocket server: {}", e);
    } else {
        info!("WebSocket server started on {}", address);
    }
}

// Function to broadcast data to all connected WebSocket clients
fn broadcast_to_clients(data: &[u8], clients: &Arc<Mutex<Vec<Sender>>>) {
    let clients = clients.lock().unwrap();
    for client in clients.iter() {
        if let Err(e) = client.send(Message::binary(data)) {
            error!("Error sending data to WebSocket client: {}", e);
        }
    }
}
