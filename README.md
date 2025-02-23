# Proxy Broker

A Go-based project for managing and verifying proxy servers. The tool provides functionalities for checking the validity of proxies, ensuring they are alive, and optionally verifying their geographical location against predefined country restrictions.

## Features

- **Proxy Verification:**
    - Validate if proxies are functional (alive).
- **Geographical Filtering:**
    - Supports country-based proxy validation
- **Concurrency:**
    - Uses goroutines for asynchronous proxy verification for improved performance.
- **Extensibility:**
    - Designed with Go's strong typing and modular architecture for easy addition of new features.

## Installation

To get started, make sure you have Go installed (v1.24 or higher is recommended). Clone the repository and build the application:

```bash
# Clone the repository
git clone https://github.com/<your_username>/proxy-broker.git

# Navigate to the project directory
cd proxy-broker

# Build the application
go build -o proxy-broker main.go
```

## Usage

Run the application using the following command after building it:

```bash
./proxy-broker
```

