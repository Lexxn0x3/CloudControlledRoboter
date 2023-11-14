use tokio::net::{TcpListener, TcpStream};
use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tokio::sync::{mpsc, Mutex};
use tokio::time::{self, Duration, Instant};
use std::sync::Arc;
use std::collections::HashMap;
use std::net::SocketAddr;
use std::io::{self, Write};

struct TrafficStats {
    bytes_received: usize,
    bytes_sent: HashMap<SocketAddr, usize>,
    last_print_time: Instant,
}

impl TrafficStats {
    fn new() -> Self {
        Self {
            bytes_received: 0,
            bytes_sent: HashMap::new(),
            last_print_time: Instant::now(),
        }
    }

    fn add_received(&mut self, bytes: usize) {
        self.bytes_received += bytes;
    }

    fn add_sent(&mut self, addr: SocketAddr, bytes: usize) {
        *self.bytes_sent.entry(addr).or_insert(0) += bytes;
    }

    fn print_and_reset(&mut self) {
        let time_elapsed = self.last_print_time.elapsed().as_secs_f64();
        let received_speed = self.bytes_received as f64 / time_elapsed;
        
        println!("Received: {} bytes ({} bytes/sec)", self.bytes_received, received_speed);
        for (addr, &bytes) in &self.bytes_sent {
            let sent_speed = bytes as f64 / time_elapsed;
            println!("Sent to {}: {} bytes ({} bytes/sec)", addr, bytes, sent_speed);
        }
        
        // Reset statistics
        self.bytes_received = 0;
        self.bytes_sent.clear();
        self.last_print_time = Instant::now();
    }
}

#[tokio::main]
async fn main() -> std::io::Result<()> {
    let (tx, mut rx) = mpsc::channel::<Vec<u8>>(100);
    let stats = Arc::new(Mutex::new(TrafficStats::new()));

    // Periodically print statistics
    let stats_for_printing = Arc::clone(&stats);
    tokio::spawn(async move {
        let mut interval = time::interval(Duration::from_secs(5));
        loop {
            interval.tick().await;
            print!("{esc}[2J{esc}[1;1H", esc = 27 as char); // Clear the console
            let mut stats = stats_for_printing.lock().await;
            stats.print_and_reset();
            io::stdout().flush().unwrap();
        }
    });

    // Accept a single connection for receiving data
    let stats_for_reading = Arc::clone(&stats);
    let read_handle = tokio::spawn(async move {
        let listener = TcpListener::bind("0.0.0.0:12345").await.unwrap();
        if let Ok((mut socket, _)) = listener.accept().await {
            let mut buf = [0u8; 4096];
            loop {
                match socket.read(&mut buf).await {
                    Ok(0) => break, // Connection was closed
                    Ok(n) => {
                        let mut stats = stats_for_reading.lock().await;
                        stats.add_received(n);
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
        let listener = TcpListener::bind("0.0.0.0:54321").await.unwrap();
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
                    stats.add_sent(*addr, data.len());
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

    Ok(())
}
