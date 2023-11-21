use std::time::{Duration, Instant};

pub struct Statistics {
    bytes_received: usize,
    bytes_sent: usize,
    last_update: Instant,
    pub buffer_size: usize,
    pub buffer_usage: usize,
}

impl Statistics {
    pub fn new() -> Self {
        Self {
            bytes_received: 0,
            bytes_sent: 0,
            last_update: Instant::now(),
            buffer_size: 0,
            buffer_usage: 0,
        }
    }

    pub fn add_received(&mut self, bytes: usize) {
        self.bytes_received += bytes;
    }

    pub fn add_sent(&mut self, bytes: usize) {
        self.bytes_sent += bytes;
    }

    pub fn set_buffer_size(&mut self, bytes: usize)
    {
        self. buffer_size = bytes;
    }

    pub fn set_buffer_usage(&mut self, bytes: usize)
    {
        self.buffer_usage = bytes;
    }

    pub fn throughput(&mut self) -> (f64, f64) {
        let now = Instant::now();
        let elapsed = now.duration_since(self.last_update);

        if elapsed > Duration::ZERO {
            let received_throughput = self.bytes_received as f64 / elapsed.as_secs_f64();
            let sent_throughput = self.bytes_sent as f64 / elapsed.as_secs_f64();

            // Reset counters and update the last_update time
            self.bytes_received = 0;
            self.bytes_sent = 0;
            self.last_update = now;

            (received_throughput, sent_throughput)
        } else {
            (0.0, 0.0)
        }
    }
}