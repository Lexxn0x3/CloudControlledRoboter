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
    let mut stdout = io::stdout();
    execute!(stdout, EnterAlternateScreen, terminal::Clear(terminal::ClearType::All))?;

    // Wrap the UI in an Arc<Mutex<>> to share between contexts
    let ui = Arc::new(Mutex::new(UI::new()?));

    let (tx, mut rx) = mpsc::channel::<Vec<u8>>(100);

    // Periodically print statistics
    let stats_for_ui = Arc::clone(&stats);
    let ui_for_drawing = Arc::clone(&ui);
    tokio::spawn(async move {
        let mut interval = time::interval(Duration::from_secs(5));
        loop {
            interval.tick().await;
            let mut ui = ui_for_drawing.lock().await;
            let mut stats = stats_for_ui.lock().await;
            
            let (received_throughput, sent_throughput) = stats.throughput(); // Get current throughput
            
            ui.data_throughput = received_throughput;

            if let Err(e) = ui.draw()
            {
                eprintln!("Error drawing UI: {}", e);
            }
        }
    });

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
                        stats.set_buffer_size(buf.len());
                        stats.set_buffer_usage(n);
                        tx.send(buf[..n].to_vec()).await.unwrap();
                    }
                    Err(e) => eprintln!("Failed to read from socket: {:?}", e),
                }
            }
        }
    });

    // Accept multiple connections for sending data
    let stats_for_writing = Arc::clone(&stats);
    let write_handle = tokio::spawn(async move {
        let listener = TcpListener::bind(format!("0.0.0.0:{}", config.client_port)).await.unwrap();
        println!("Server is running for sending data on port: {}", config.client_port); // Print server start information
        let client_map = Arc::new(Mutex::new(HashMap::<SocketAddr, TcpStream>::new()));

        // Accept new clients and store them in a HashMap
        let client_map_clone = Arc::clone(&client_map);
        tokio::spawn(async move {
            loop {
                if let Ok((socket, addr)) = listener.accept().await {
                    client_map_clone.lock().await.insert(addr, socket);
                }
            }
        });

        // Read from the channel and write to all connected clients
        while let Some(data) = rx.recv().await {
            let mut clients = client_map.lock().await;
            let mut disconnected = Vec::new(); // Track disconnected clients

            for (addr, socket) in clients.iter_mut() {
                if let Err(e) = socket.write_all(&data).await {
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
    });

    // Wait for both tasks to complete
    let _ = tokio::try_join!(read_handle, write_handle);

    {
        let mut ui = ui.lock().await;
        ui.cleanup()?;

        // Leave the alternate screen
        execute!(io::stdout(), LeaveAlternateScreen)?;
    }

    Ok(())
}
