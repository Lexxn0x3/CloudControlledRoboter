use clap::{App, Arg};

pub struct Config {
    pub single_connection_port: u16,
    pub multi_connection_port: u16,
    pub websocket_port: u16,
    pub debug_level: String,
    pub buffer_size: usize,
}

pub fn parse_arguments() -> Config {
    let matches = App::new("My Program")
        .arg(Arg::with_name("single-connection-port")
            .long("single-connection-port")
            .help("Sets the single connection port")
            .takes_value(true)
            .default_value("3001")
            .validator(validate_port))
        .arg(Arg::with_name("multi-connection-port")
            .long("multi-connection-port")
            .help("Sets the multi connection port")
            .takes_value(true)
            .default_value("4001")
            .validator(validate_port))
        .arg(Arg::with_name("websocket-port")
            .long("websocket-port")
            .help("Sets the multi websocket port")
            .takes_value(true)
            .default_value("5001")
            .validator(validate_port))
        .arg(Arg::with_name("debug-level")
            .long("debug-level")
            .help("Sets the debug level")
            .takes_value(true)
            .default_value("info")
            .possible_values(&["info", "error", "debug"]))
        .arg(Arg::with_name("buffer-size")
            .long("buffer-size")
            .help("Sets the buffer size in bytes: higher values would be 131072 or 512000 (standard)")
            .takes_value(true)
            .default_value("512000")
            .validator(|v| v.parse::<usize>()
                             .map(|_| ())
                             .map_err(|_| String::from("Buffer size must be an integer"))))
        .get_matches();

    Config
    {
        single_connection_port: matches.value_of("single-connection-port").unwrap().parse().unwrap(),
        multi_connection_port: matches.value_of("multi-connection-port").unwrap().parse().unwrap(),
        websocket_port: matches.value_of("websocket-port").unwrap().parse().unwrap(),
        debug_level: matches.value_of("debug-level").unwrap().to_string(),
        buffer_size: matches.value_of("buffer-size").unwrap().parse().unwrap(),
    }
}

fn validate_port(v: &str) -> Result<(), String> {
    v.parse::<u16>()
        .map_err(|_| String::from("The port must be a number"))
        .and_then(|val| {
            if val > 1024 && val < 65535 {
                Ok(())
            } else {
                Err(String::from("Port must be a number between 1025 and 65534"))
            }
        })
}
