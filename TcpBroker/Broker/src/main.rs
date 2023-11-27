mod config;
mod websocket_server;

use std::sync::mpsc;
use std::sync::{Arc, Mutex};
use std::thread;
use std::net::{TcpListener, TcpStream};
use std::io::{Read, Write, Error};
use std::sync::mpsc::{Receiver, Sender};
use env_logger;
use log::{info, error, debug};
use env_logger::{Builder, Env};
use ws;

fn main() {
    let config = config::parse_arguments();

    // Initialize the logger with the specified level
    Builder::from_env(Env::default().default_filter_or(config.debug_level)).init();

    let client_senders = Arc::new(Mutex::new(vec![]));

    // Start single_connection_listener in a separate thread
    let client_senders_clone = Arc::clone(&client_senders);
    thread::spawn(move || {
        single_connection_listener(client_senders_clone, config.single_connection_port, config.buffer_size);
    });

    // Start multi_conection_listener in a separate thread
    let client_senders_clone = Arc::clone(&client_senders);
    thread::spawn(move || {
        multi_conection_listener(&client_senders_clone, config.multi_connection_port);
    });

    let websocket_client_senders = Arc::new(Mutex::new(vec![]));

    // Start the WebSocket listener in a new thread.
    // This thread will broadcast messages from the shared channel to WebSocket clients.
    let websocket_client_senders_clone = Arc::clone(&websocket_client_senders);
    let client_senders_clone = Arc::clone(&client_senders);
    thread::spawn(move || {
        multi_connection_websocket_listener(client_senders_clone, websocket_client_senders_clone, config.websocket_port);
    });

    // Main loop to keep the main thread alive
    loop {
        thread::sleep(std::time::Duration::from_secs(1));
    }
}

fn multi_conection_listener(client_senders: &Arc<Mutex<Vec<Sender<Vec<u8>>>>>, port: u16) {
    let multi_listener = TcpListener::bind(format!("0.0.0.0:{}", port)).unwrap();
    info!("Multi-connection TCP listener started on {}:{}", multi_listener.local_addr().unwrap().ip(), multi_listener.local_addr().unwrap().port());

    for stream in multi_listener.incoming() {
        let client_senders_clone = Arc::clone(&client_senders);
        if let Ok(stream) = stream {
            debug!("Accepted new multi connection from {}:{}", stream.peer_addr().unwrap().ip(), stream.peer_addr().unwrap().port());

            let (client_tx, client_rx) = mpsc::channel();
            client_senders_clone.lock().unwrap().push(client_tx);

            thread::spawn(move || {
                handle_client_connection(stream, client_rx).unwrap();
            });
        }
    }
}

// Starts a new thread that listens for messages on the channel and broadcasts them.
fn start_message_broadcaster_thread(rx: mpsc::Receiver<Vec<u8>>, websocket_client_senders: Arc<Mutex<Vec<ws::Sender>>>) {
    thread::spawn(move || {
        for message in rx {
           websocket_server::broadcast_message_to_websocket_clients(&websocket_client_senders, message);
        }
    });
}

// Thread that listens for messages from the shared channel and broadcasts them to WebSocket clients
fn start_websocket_broadcast_thread(rx: Receiver<Vec<u8>>, websocket_client_senders: Arc<Mutex<Vec<ws::Sender>>>) {
    thread::spawn(move || {
        for message in rx {
            let message = message.clone();
            let clients = websocket_client_senders.lock().unwrap();
            for client in clients.iter() {
                let _ = client.send(ws::Message::binary(message.clone()));
            }
        }
    });
}

// Your existing WebSocket server start function, now starts the broadcaster thread.
fn multi_connection_websocket_listener(client_senders: Arc<Mutex<Vec<mpsc::Sender<Vec<u8>>>>>, websocket_client_senders: Arc<Mutex<Vec<ws::Sender>>>, port: u16) {
    // Start a thread dedicated to broadcasting messages to WebSocket clients
    let (tx, rx) = mpsc::channel::<Vec<u8>>();
    start_websocket_broadcast_thread(rx, websocket_client_senders.clone());

    // Add the new broadcaster's sender to the shared channel vector
    {
        let mut senders = client_senders.lock().unwrap();
        senders.push(tx);
    }

    // Start the WebSocket server
    if let Err(e) = ws::listen(format!("0.0.0.0:{}", port), move |out| {
        // Add the new WebSocket client's sender to the list
        websocket_client_senders.lock().unwrap().push(out.clone());

        websocket_server::WebSocketServer {
            out: out,
            websocket_client_senders: websocket_client_senders.clone(),
        }
    }) {
        error!("Failed to start WebSocket server: {:?}", e);
    }
}

// This function sets up the WebSocket listener.

fn single_connection_listener(client_senders: Arc<Mutex<Vec<Sender<Vec<u8>>>>>, port: u16, buffer_size: usize)
{
    let single_listener = TcpListener::bind(format!("0.0.0.0:{}", port)).unwrap();
    info!("Single-connection TCP listener started on {}:{}", single_listener.local_addr().unwrap().ip(), single_listener.local_addr().unwrap().port());

    loop {
        let client_senders_clone = Arc::clone(&client_senders);
        match single_listener.accept() {
            Ok((stream, _)) => {
                debug!("Accepted single connection from {}:{}", stream.peer_addr().unwrap().ip(), stream.peer_addr().unwrap().port());
                if let Err(e) = handle_single_connection(stream, client_senders_clone, buffer_size) {
                    error!("Error handling single connection: {}", e);
                    // Here you can handle any cleanup or reset actions needed
                }
            },
            Err(e) => error!("Failed to accept connection: {}", e),
        }
    }
}

fn handle_single_connection(mut stream: TcpStream, clients: Arc<Mutex<Vec<mpsc::Sender<Vec<u8>>>>>, buffer_size: usize) -> Result<(), Error> {
    debug!("Single connection handler started.");
    let mut buffer = vec![0; buffer_size];  // Dynamic buffer based on buffer_size
    loop {
        match stream.read(&mut buffer) {
            Ok(nbytes) => {
                if nbytes == 0 {
                    debug!("Single connection closed by the client.");
                    break;
                }
                let data = buffer[..nbytes].to_vec();
                let mut clients_to_remove = vec![];
                for (i, client) in clients.lock().unwrap().iter().enumerate() {
                    if client.send(data.clone()).is_err() {
                        error!("Error sending data to a multi-connection client. Removing client.");
                        clients_to_remove.push(i);
                    }
                }
                let mut clients = clients.lock().unwrap();
                for i in clients_to_remove.iter().rev() {
                    clients.remove(*i);
                }
            },
            Err(e) => {
                error!("Error reading from single connection: {}", e);
                break;
            }
        }
    }
    debug!("Single connection handler ended.");
    Ok(())
}

fn handle_client_connection(mut stream: TcpStream, rx: mpsc::Receiver<Vec<u8>>) -> Result<(), Error>
{
    debug!("Multi-connection client handler started.");
    loop
    {
        match rx.recv()
        {
            Ok(data) => {
                if stream.write_all(&data).is_err()
                {
                    error!("Error writing to multi-connection client. Client may have disconnected.");
                    break;
                }
            },
            Err(_) =>
                {
                debug!("Multi-connection client disconnected. Ending handler thread.");
                break;
            }
        }
    }
    debug!("Multi-connection client handler ended.");
    Ok(())
}
