use std::sync::mpsc;
use std::sync::{Arc, Mutex};
use std::thread;
use std::net::{TcpListener, TcpStream};
use std::io::{Read, Write, Error};
use std::sync::mpsc::{Receiver, Sender};
use tungstenite::accept;

fn main()
{
    let (tx, rx) = mpsc::channel::<Vec<u8>>();
    let rx = Arc::new(Mutex::new(rx));


    single_connection_listener(tx);

    loop
    {
        multi_connection_listener(&rx);
    }
}

fn multi_connection_listener(rx: &Arc<Mutex<Receiver<Vec<u8>>>>)
{
    let multi_listener = TcpListener::bind("0.0.0.0:3002").unwrap();
    println!("Multi-connection TCP listener started on 127.0.0.1:3002");
    for stream in multi_listener.incoming()
    {
        let rx_clone = Arc::clone(&rx);
        if let Ok(stream) = stream
        {
            println!("Accepted new client connection");
            thread::spawn(move ||
            {
                handle_client_connection(stream, rx_clone).unwrap();
            });
        }
    }
}

fn single_connection_listener(tx: Sender<Vec<u8>>)
{
    let single_listener = TcpListener::bind("0.0.0.0:3001").unwrap();
    println!("Single-connection TCP listener started on 127.0.0.1:3001");
    let tx_clone = tx.clone();
    let listener = thread::spawn(move ||
    {
        loop
        {
            let tx_clone = tx.clone();
            if let Ok((stream, addr)) = single_listener.accept()
            {
                println!("Accepted connection from {:?}", addr);
                handle_single_connection(stream, tx_clone).unwrap();
            }
        }
    });
}

fn handle_single_connection(mut stream: TcpStream, tx: mpsc::Sender<Vec<u8>>) -> Result<(), Error>
{
    let mut buffer = [0; 131072];
    loop
    {
        let nbytes = stream.read(&mut buffer)?;
        if nbytes == 0 {
            break;
        }
        println!("Received {} bytes from the single connection", nbytes);
        tx.send(buffer[..nbytes].to_vec()).unwrap();
    }
    Ok(())
}

fn handle_client_connection(mut stream: TcpStream, rx: Arc<Mutex<mpsc::Receiver<Vec<u8>>>>) -> Result<(), Error>
{
    loop
    {
        let data = rx.lock().unwrap().recv().unwrap();
        stream.write_all(&data)?;
        println!("Sent data to a client");
    }
    Ok(())
}
