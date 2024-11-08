use clap::{App, Arg};

pub struct Config {
    pub single_connection_port: u16,
    pub multi_connection_port: u16,
    pub websocket_connection_port: u16,
    pub single_connection_timeout: u64,
    pub no_websocket: bool,
    pub debug_level: String,
    pub buffer_size: usize,
}

pub fn parse_arguments() -> Config {
    let matches = App::new("My Program")
        .arg(Arg::with_name("single-connection-port")
            .long("single-connection-port")
            .short('s')  // Changed to a char
            .help("Sets the port for single connections. Default is 3001.")
            .takes_value(true)
            .default_value("3001")
            .validator(validate_port))
        .arg(Arg::with_name("multi-connection-port")
            .long("multi-connection-port")
            .short('m')  // Changed to a char
            .help("Sets the port for multiple connections. Default is 9001.")
            .takes_value(true)
            .default_value("9001")
            .validator(validate_port))
        .arg(Arg::with_name("websocket-connection-port")
            .long("websocket-connection-port")
            .short('w')  // Changed to a char
            .help("Sets the port for websocket connections. Default is 5001.")
            .takes_value(true)
            .default_value("5001")
            .validator(validate_port))
        .arg(Arg::with_name("no-websocket")
            .long("no-websocket")
            .short('n')  // Changed to a char
            .help("Disables the websocket connection functionality.")
            .takes_value(false))
        .arg(Arg::with_name("debug-level")
            .long("debug-level")
            .short('d')  // Changed to a char
            .help("Sets the logging debug level. Possible values are 'info', 'error', 'debug'. Default is 'info'.")
            .takes_value(true)
            .default_value("info")
            .possible_values(&["info", "error", "debug"]))
        .arg(Arg::with_name("buffer-size")
            .long("buffer-size")
            .short('b')  // Changed to a char
            .help("Sets the buffer size in bytes. Higher values might be 131072 or 512000 (standard). Default is 512000.")
            .takes_value(true)
            .default_value("512000")
            .validator(|v| v.parse::<usize>()
                .map(|_| ())
                .map_err(|_| String::from("Buffer size must be an integer"))))
        .arg(Arg::with_name("single-connection-timeout")
            .long("single-connection-timeout")
            .short('t')  // Changed to a char
            .help("Sets the timeout for the single connection. If the single connection doesnt send data for default 5 seconds it gets disconnected")
            .takes_value(true)
            .default_value("5")
            .validator(|v| v.parse::<usize>()
                .map(|_| ())
                .map_err(|_| String::from("Timeout must be an integer"))))
        .get_matches();

    Config
    {
        single_connection_port: matches.value_of("single-connection-port").unwrap().parse().unwrap(),
        multi_connection_port: matches.value_of("multi-connection-port").unwrap().parse().unwrap(),
        websocket_connection_port: matches.value_of("websocket-connection-port").unwrap().parse().unwrap(),
        no_websocket: matches.is_present("no-websocket"),
        single_connection_timeout: matches.value_of("single-connection-timeout").unwrap().parse().unwrap(),
        debug_level: matches.value_of("debug-level").unwrap().to_string(),
        buffer_size: matches.value_of("buffer-size").unwrap().parse().unwrap(),
    }
}

fn validate_port(v: &str) -> Result<(), String> {
    v.parse::<u16>()
        .map_err(|_| String::from("The port must be a number"))
        .and_then(|val| {
            if val > 1024 && val < 65535
            {
                Ok(())
            } else {
                Err(String::from("Port must be a number between 1025 and 65534"))
            }
        })
}
