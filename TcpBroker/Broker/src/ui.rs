// ui.rs

use crossterm::terminal::{disable_raw_mode, enable_raw_mode};
use std::io;
use tui::backend::CrosstermBackend;
use tui::layout::{Constraint, Direction, Layout};
use tui::style::{Color, Modifier, Style};
use tui::widgets::{Block, Borders, Gauge, List, ListItem, Paragraph};
use tui::Terminal;

pub struct UI {
    terminal: Terminal<CrosstermBackend<io::Stdout>>,
    pub data_throughput: f64,
    pub client_statuses: Vec<String>,
    pub buffer_size: usize,
    pub buffer_usage: usize,
    pub debug_messages: Vec<String>,
}

impl UI {
    pub fn new() -> io::Result<Self> {
        enable_raw_mode()?;
        let stdout = io::stdout();
        let backend = CrosstermBackend::new(stdout);
        let terminal = Terminal::new(backend)?;
        Ok(Self {
            terminal,
            data_throughput: 0.0,
            client_statuses: vec![],
            buffer_size: 0,
            buffer_usage: 0,
            debug_messages: vec![],
        })
    }

    pub fn draw(&mut self) -> io::Result<()> {
        self.terminal.draw(|f| {
            let chunks = Layout::default()
                .direction(Direction::Vertical)
                .margin(1)
                .constraints(
                    [
                        Constraint::Length(3),
                        Constraint::Length(3 + self.client_statuses.len() as u16),
                        Constraint::Length(3),
                        Constraint::Min(1),
                    ]
                    .as_ref(),
                )
                .split(f.size());

            let throughput = format!("{:.2} bytes/sec", self.data_throughput);
            let throughput_widget = Paragraph::new(throughput)
                .block(Block::default().borders(Borders::ALL).title("Throughput"));
            f.render_widget(throughput_widget, chunks[0]);

            let clients: Vec<ListItem> = self
                .client_statuses
                .iter()
                .map(|s| ListItem::new(s.clone()))
                .collect();
            let clients_widget = List::new(clients)
                .block(Block::default().borders(Borders::ALL).title("Clients"));
            f.render_widget(clients_widget, chunks[1]);

            let buffer_gauge = Gauge::default()
                .block(Block::default().borders(Borders::ALL).title("Buffer Usage"))
                .gauge_style(Style::default().fg(Color::Cyan).bg(Color::Black).add_modifier(Modifier::ITALIC))
                .percent((self.buffer_usage as f64 / self.buffer_size as f64 * 100.0) as u16);
            f.render_widget(buffer_gauge, chunks[2]);

            let debug: Vec<ListItem> = self
                .debug_messages
                .iter()
                .map(|m| ListItem::new(m.clone()))
                .collect();
            let debug_widget = List::new(debug)
                .block(Block::default().borders(Borders::ALL).title("Debug"));
            f.render_widget(debug_widget, chunks[3]);
        })?;
        Ok(())
    }

    pub fn cleanup(&mut self) -> io::Result<()> {
        disable_raw_mode()?;
        self.terminal.show_cursor()?;
        Ok(())
    }
}
