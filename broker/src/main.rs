mod config;
mod websocket;

use std::sync::mpsc;
use std::sync::{Arc, Mutex};
use std::thread;
use std::net::{TcpListener, TcpStream};
use std::io::{Read, Write, Error};
use std::sync::mpsc::{Sender};
use env_logger;
use log::{info, error, debug};
use env_logger::{Builder, Env};

fn main()
{
    //test
    //Get config from Args
    let config = config::parse_arguments();

    // Initialize the logger with the specified level
    Builder::from_env(Env::default().default_filter_or(config.debug_level)).init();

    let client_senders = Arc::new(Mutex::new(vec![]));

    // Start single_connection_listener in a separate thread
    let client_senders_clone = Arc::clone(&client_senders);
    thread::spawn(move || {
        single_connection_listener(client_senders_clone, config.single_connection_port, config.buffer_size);
    });

    // Clone client_senders for the multi_connection_listener and start it in a separate thread
    let client_senders_clone_for_multi = Arc::clone(&client_senders);
    thread::spawn(move || {
        multi_conection_listener(&client_senders_clone_for_multi, config.multi_connection_port);
    });

    if  !config.no_websocket
    {
        // Delay the start of the WebSocket server
        thread::sleep(std::time::Duration::from_secs(1));
        websocket::start_websocket("127.0.0.1", config.multi_connection_port, config.buffer_size, config.websocket_connection_port);
    }

    // Keep the main thread alive
    loop {
        thread::sleep(std::time::Duration::from_secs(60)); // Sleep for a minute in each iteration
    }
}

fn multi_conection_listener(client_senders: &Arc<Mutex<Vec<Sender<Vec<u8>>>>>, port: u16) {
    let multi_listener = TcpListener::bind(format!("0.0.0.0:{}", port)).unwrap();
    info!("Multi-connection TCP listener started on {}:{}", multi_listener.local_addr().unwrap().ip(), multi_listener.local_addr().unwrap().port());

    for stream in multi_listener.incoming() {
        let client_senders_clone = Arc::clone(&client_senders);
        if let Ok(stream) = stream {
            stream.set_nodelay(true).expect("set_nodelay call failed");
            debug!("Accepted new multi connection from {}:{}", stream.peer_addr().unwrap().ip(), stream.peer_addr().unwrap().port());

            let (client_tx, client_rx) = mpsc::channel();
            client_senders_clone.lock().unwrap().push(client_tx);

            thread::spawn(move || {
                handle_client_connection(stream, client_rx).unwrap();
            });
        }
    }
}

fn single_connection_listener(client_senders: Arc<Mutex<Vec<Sender<Vec<u8>>>>>, port: u16, buffer_size: usize) {
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

fn handle_client_connection(mut stream: TcpStream, rx: mpsc::Receiver<Vec<u8>>) -> Result<(), Error> {
    debug!("Multi-connection client handler started.");
    loop {
        match rx.recv() {
            Ok(data) => {
                if stream.write_all(&data).is_err() {
                    error!("Error writing to multi-connection client. Client may have disconnected.");
                    break;
                }
            },
            Err(_) => {
                debug!("Multi-connection client disconnected. Ending handler thread.");
                break;
            }
        }
    }
    debug!("Multi-connection client handler ended.");
    Ok(())
}
