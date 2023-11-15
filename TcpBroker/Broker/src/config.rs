use clap::{App, Arg};

struct Config {
    server_port: u16,
    client_port: u16,
    debug_level: String,
    buffer_size: usize,
}

fn parse_arguments() -> Config {
    let matches = App::new("My Program")
        .arg(Arg::with_name("server-port")
            .long("server-port")
            .help("Sets the server port")
            .takes_value(true)
            .default_value("8080")
            .validator(validate_port))
        .arg(Arg::with_name("client-port")
            .long("client-port")
            .help("Sets the client port")
            .takes_value(true)
            .default_value("8081")
            .validator(validate_port))
        .arg(Arg::with_name("debug-level")
            .long("debug-level")
            .help("Sets the debug level")
            .takes_value(true)
            .default_value("info")
            .possible_values(&["trace", "debug", "info", "warn", "error"]))
        .arg(Arg::with_name("buffer-size")
            .long("buffer-size")
            .help("Sets the buffer size in bytes")
            .takes_value(true)
            .default_value("4096")
            .validator(|v| v.parse::<usize>()
                             .map(|_| ())
                             .map_err(|_| String::from("Buffer size must be an integer"))))
        .get_matches();

    Config {
        server_port: matches.value_of("server-port").unwrap().parse().unwrap(),
        client_port: matches.value_of("client-port").unwrap().parse().unwrap(),
        debug_level: matches.value_of("debug-level").unwrap().to_string(),
        buffer_size: matches.value_of("buffer-size").unwrap().parse().unwrap(),
    }
}

fn validate_port(v: String) -> Result<(), String> {
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