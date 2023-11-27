use std::sync::{Arc, mpsc, Mutex};
use std::thread;
use ws::{listen, Handler, Sender, Message, Result as WsResult, CloseCode};
use log::{info, error, debug};

// Define a WebSocket handler struct that will implement the Handler trait.
pub(crate) struct WebSocketServer {
    pub(crate) out: Sender,
    pub(crate) websocket_client_senders: Arc<Mutex<Vec<Sender>>>,
}

impl Handler for WebSocketServer
{
    fn on_open(&mut self, handshake: ws::Handshake) -> WsResult<()>
    {
        // Add the new WebSocket client's sender to the list
        self.websocket_client_senders.lock().unwrap().push(self.out.clone());
        info!("WebSocket client connected: {}", handshake.peer_addr.unwrap());
        Ok(())
    }

    fn on_close(&mut self, code: CloseCode, reason: &str)
    {
        // Remove the WebSocket client's sender from the list
        let mut clients = self.websocket_client_senders.lock().unwrap();
        clients.retain(|client| client.token() != self.out.token());

        info!("WebSocket client disconnected with code {:?}, reason: {}", code, reason);
    }

    fn on_error(&mut self, error: ws::Error)
    {
        error!("WebSocket encountered an error: {:?}", error);
    }
    fn on_message(&mut self, msg: Message) -> WsResult<()>
    {
        Ok(())
    }
}
pub(crate) fn broadcast_message_to_websocket_clients(websocket_client_senders: &Arc<Mutex<Vec<Sender>>>, message: Vec<u8>) {
    let clients = websocket_client_senders.lock().unwrap();
    for client in clients.iter() {
        let _ = client.send(Message::binary(message.clone()));
    }
}