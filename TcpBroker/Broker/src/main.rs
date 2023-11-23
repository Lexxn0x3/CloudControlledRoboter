use std::collections::HashMap;
use std::io::{self, Write};
use std::net::SocketAddr;
use std::sync::Arc;
use std::time::Instant;

use crossterm::terminal::{disable_raw_mode, enable_raw_mode};
use crossterm::{
    execute,
    terminal::{self, EnterAlternateScreen, LeaveAlternateScreen},
};
use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tokio::net::{TcpListener, TcpStream};
use tokio::sync::{mpsc, Mutex};
use tokio::time::{self, Duration};
use tokio::spawn;
use tokio::sync::mpsc::Receiver;
use tui::backend::CrosstermBackend;
use tui::layout::{Constraint, Direction, Layout};
use tui::style::{Color, Modifier, Style};
use tui::widgets::{Block, Borders, Gauge, List, ListItem, Paragraph};
use tui::Terminal;

mod config;
mod ui;
use ui::UI;
mod statistics;

use statistics::Statistics;

#[tokio::main]
async fn main() -> std::io::Result<()> {
    let config = config::parse_arguments();
    let stats = Arc::new(Mutex::new(Statistics::new()));

    // Enter the alternate screen and clear it
    if config.debug_level != "none"
    {
        let mut stdout = io::stdout();
        execute!(stdout, EnterAlternateScreen, terminal::Clear(terminal::ClearType::All))?;
    }

    // Wrap the UI in an Arc<Mutex<>> to share between contexts
    let ui = Arc::new(Mutex::new(UI::new()?));

    let (tx, mut rx) = mpsc::channel::<Vec<u8>>(100);


    if config.debug_level != "none"
    {
        // Periodically print statistics
        let stats_for_ui = Arc::clone(&stats);
        let ui_for_drawing = Arc::clone(&ui);
        tokio::spawn(async move {
            let mut interval = time::interval(Duration::from_millis(100));
            loop {
                interval.tick().await;
                let mut ui = ui_for_drawing.lock().await;
                let mut stats = stats_for_ui.lock().await;
                
                let (received_throughput, sent_throughput) = stats.throughput(); // Get current throughput
                
                ui.data_throughput = received_throughput;
                ui.buffer_size = config.buffer_size;
                ui.buffer_usage = stats.buffer_usage;

                if let Err(e) = ui.draw()
                {
                    eprintln!("Error drawing UI: {}", e);
                }
            }
        });
    }

    // Accept a single connection for receiving data
    let stats_for_reading = Arc::clone(&stats);
    let read_handle = tokio::spawn(async move {
        let listener = TcpListener::bind(format!("0.0.0.0:{}", config.server_port)).await.unwrap();
        println!("Server is running for receiving data on port: {}", config.server_port); // Print server start information
        if let Ok((mut socket, _)) = listener.accept().await {
            let mut buf = vec![0u8; config.buffer_size]; // Create a buffer with the configured size
            loop {
                match socket.read(&mut buf).await {
                    Ok(0) => break, // Connection was closed
                    Ok(n) => {
                        let mut stats = stats_for_reading.lock().await;
                        stats.add_received(n);
                        stats.set_buffer_usage(n);
                        tx.send(buf[..n].to_vec()).await.unwrap();
                    }
                    Err(e) => eprintln!("Failed to read from socket: {:?}", e),
                }
            }
        }
    });

    // Accept multiple connections for sending data
    let client_map: Arc<Mutex<HashMap<SocketAddr, TcpStream>>> = Arc::new(Mutex::new(HashMap::new()));

    // The listener task that accepts new connections
    let client_map_for_accepting = Arc::clone(&client_map);
    tokio::spawn(async move {
        let listener = TcpListener::bind(format!("0.0.0.0:{}", config.client_port)).await.unwrap();
        println!("Server is running for sending data on port: {}", config.client_port);

        loop {
            match listener.accept().await {
                Ok((socket, addr)) => {
                    client_map_for_accepting.lock().await.insert(addr, socket);
                    println!("New client: {}", addr);
                }
                Err(e) => eprintln!("Failed to accept client: {}", e),
            }
        }
    });


    // The task that handles sending data to clients
    let client_map_for_sending = Arc::clone(&client_map);
    let disconnected = Arc::new(Mutex::new(Vec::new())); // Declare outside the loop

    let sender = tokio::spawn(async move
    {
        while let Some(data) = rx.recv().await {
            let mut clients = client_map_for_sending.lock().await;

            // Collect the addresses to avoid holding the lock while spawning tasks
            let addrs: Vec<SocketAddr> = clients.keys().cloned().collect();

            for addr in addrs {
                if let Some(mut stream) = clients.remove(&addr) { // Remove the stream to pass ownership to the task
                    let data = data.clone(); // Clone the data for each task
                    let disconnected_clone = Arc::clone(&disconnected); // Clone the Arc to share ownership of the Mutex

                    tokio::spawn(async move {
                        match time::timeout(Duration::from_millis(200), stream.write_all(&data)).await {
                            Ok(result) => {
                                if let Err(e) = result {
                                    eprintln!("Failed to write to client {}: {}", addr, e);
                                    let mut disconnected = disconnected_clone.lock().await;
                                    disconnected.push(addr);
                                }
                            },
                            Err(_) => {
                                eprintln!("Write to client {} timed out", addr);
                                let mut disconnected = disconnected_clone.lock().await;
                                disconnected.push(addr);
                            },
                        }
                    });
                }
            }

            // Reinserting the streams back should be handled here,
            // after all tasks have been spawned, or by the tasks themselves if they succeed.

            // Handle disconnections here
            // ...
        }
        // Disconnection handling loop goes here
    });

    sender.await.ok();

    // Wait for both tasks to complete
    /*
    let _ = tokio::try_join!(read_handle, write_handle);

    {
        let mut ui = ui.lock().await;
        ui.cleanup()?;

        // Leave the alternate screen
        execute!(io::stdout(), LeaveAlternateScreen)?;
    }
    */

    Ok(())
}

async fn WriteToAllClients(mut rx: Receiver<Vec<u8>>, stats_for_writing: Arc<Mutex<Statistics>>, client_map: Arc<Mutex<HashMap<SocketAddr, TcpStream>>>) {
// Read from the channel and write to all connected clients
    while let Some(data) = rx.recv().await
    {
        let mut clients = client_map.lock().await;
        let mut disconnected = Vec::new(); // Track disconnected clients

        for (addr, socket) in clients.iter_mut()
        {
            if let Err(e) = socket.write_all(&data).await
            {
                eprintln!("Failed to write to client {}: {}", addr, e);
                disconnected.push(*addr); // Mark the client for removal
            } else {
                let mut stats = stats_for_writing.lock().await;
            }
        }

        // Remove disconnected clients
        for addr in disconnected {
            clients.remove(&addr);
        }
    }
}
