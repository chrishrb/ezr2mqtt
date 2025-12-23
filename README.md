# ezr2mqtt - Möhlenhoff Alpha 2 MQTT Gateway

A bridge service that connects Möhlenhoff EZR (Alpha 2) heating controllers to MQTT, enabling integration with home automation systems like Home Assistant.

## Features

- 🔄 **Bidirectional Communication**: Read temperature values and control heating zones via MQTT
- 📊 **Real-time Monitoring**: Periodic polling of EZR device status
- 🏠 **Home Assistant Ready**: MQTT integration for easy setup
- 🐳 **Docker Support**: Easy deployment with Docker and Docker Compose
- 🔧 **Flexible Configuration**: YAML-based configuration
- 🌡️ **Full Device Support**: Access to temperatures, heating modes, zones, and more

## What is Möhlenhoff Alpha 2?

Möhlenhoff Alpha 2 is a heating control system with EZR controllers that manage individual room temperatures. This gateway allows you to integrate these controllers into your smart home setup.

## Quick Start

### Docker Compose (Recommended)

1. Clone the repository:
   ```bash
   git clone https://github.com/chrishrb/ezr2mqtt.git
   cd ezr2mqtt
   ```

2. Copy the example configuration:
   ```bash
   cp ezr2mqtt.example.yaml ezr2mqtt.yaml
   ```

3. Edit `ezr2mqtt.yaml` with your settings:
   ```yaml
   api:
     type: mqtt
     mqtt:
       urls:
         - mqtt://myuser:mypass@mqtt.example.com:1883
   
   ezr:
     - host: EZR01A3AF.lan
       name: EG
   
   general:
     poll_every: 60s
   ```

4. Start the service:
   ```bash
   docker-compose up -d
   ```

### Binary Installation

1. Download the latest release for your platform from the [releases page](https://github.com/chrishrb/ezr2mqtt/releases)

2. Copy the example configuration:
   ```bash
   cp ezr2mqtt.example.yaml ezr2mqtt.yaml
   ```

3. Edit `ezr2mqtt.yaml` with your settings:
   ```yaml
   api:
     type: mqtt
     mqtt:
       urls:
         - mqtt://myuser:mypass@mqtt.example.com:1883
   
   ezr:
     - host: EZR01A3AF.lan
       name: EG
   
   general:
     poll_every: 60s
   ```

4. Run the service:
   ```bash
   ./ezr2mqtt start -c ezr2mqtt.yaml
   ```

## Configuration

### Complete Configuration Example

```yaml
api:
  type: mqtt
  mqtt:
    urls:
      - mqtt://username:password@mqtt-broker:1883
    prefix: ezr                    # MQTT topic prefix (default: ezr)
    group: ezr2mqtt                # MQTT group ID (default: ezr2mqtt)
    connect_timeout: 10s           # Connection timeout
    connect_retry_delay: 1s        # Retry delay on connection failure
    keepalive_interval: 60s        # Keep-alive interval

ezr:
  - name: ground_floor             # Friendly name for the device
    type: http                     # Type: http or mock
    http:
      host: EZR01A3AF.lan          # EZR device hostname or IP
  - name: first_floor
    type: http
    http:
      host: EZR01B2CD.lan

general:
  poll_every: 60s                  # How often to poll EZR devices
```

### Configuration Options

#### API Settings
- **type**: Communication type (currently only `mqtt` is supported)
- **mqtt.urls**: List of MQTT broker URLs
- **mqtt.prefix**: MQTT topic prefix (default: `ezr`)
- **mqtt.group**: MQTT consumer group (default: `ezr2mqtt`)
- **mqtt.connect_timeout**: Connection timeout duration
- **mqtt.connect_retry_delay**: Delay between connection retries
- **mqtt.keepalive_interval**: MQTT keep-alive interval

#### EZR Settings
- **name**: Unique identifier for the device
- **type**: Client type - `http` for real devices, `mock` for testing
- **http.host**: Hostname or IP address of the EZR controller

#### General Settings
- **poll_every**: Polling interval for fetching device status (e.g., `60s`, `5m`)

## MQTT Topics

### Published Topics (Device → MQTT)

The service publishes device state to MQTT with the following structure:

```
ezr/{device_name}/0/state/meta
ezr/{device_name}/+/state/temperature_target
ezr/{device_name}/+/state/temperature_actual
ezr/{device_name}/+/state/heatarea_mode
```

### Subscribed Topics (MQTT → Device)

Send commands to control your heating system:

#### Set Target Temperature

```
Topic: ezr/{device_name}/{room_id}/set/temperature_target
Payload: "22.20"
```

#### Set Heat Area Mode

```
Topic: ezr/{device_name}/{room_id}/set/heatarea_mode
Payload: "auto"
```

Modes:
- auto
- day
- night

## Development

### Prerequisites

- Go 1.23 or later
- Make (optional, for using Makefile commands)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/chrishrb/ezr2mqtt.git
cd ezr2mqtt

# Build
make build

# Or build manually
go build -o bin/ezr2mqtt main.go
```

### Running Tests

```bash
# Run unit tests
make test

# Run end-to-end tests
make e2e
```

### Available Make Targets

```bash
make run           # Run the application
make build         # Build the binary
make test          # Run unit tests
make e2e           # Run end-to-end tests
make compile       # Cross-compile for multiple platforms
make format        # Format code
make lint          # Run linter
make clean         # Clean build artifacts
```

### Project Structure

```
.
├── api/              # API interfaces and MQTT implementation
├── cmd/              # CLI commands (Cobra)
├── config/           # Configuration loading and validation
├── handlers/         # Message handlers for device control
├── polling/          # Periodic polling logic
├── store/            # In-memory state storage
├── transport/        # Device communication (HTTP, mock)
├── e2e/              # End-to-end tests
└── main.go           # Application entry point
```

## Docker

### Building the Docker Image

```bash
docker build -t ezr2mqtt .
```

### Running with Docker

```bash
docker run -d \
  --name ezr2mqtt \
  --network host \
  -v $(pwd)/ezr2mqtt.yaml:/config/ezr2mqtt.yaml:ro \
  ezr2mqtt:latest
```

### Docker Compose

The included `docker-compose.yml` provides a ready-to-use setup:

```yaml
version: '3.8'

services:
  ezr2mqtt:
    build: .
    container_name: ezr2mqtt
    restart: unless-stopped
    volumes:
      - ./ezr2mqtt.yaml:/config/ezr2mqtt.yaml:ro
    environment:
      - TZ=Europe/Berlin
    network_mode: host
```

## Troubleshooting

### Connection Issues

- Ensure your EZR device is reachable on the network
- Check that the hostname/IP in the configuration is correct
- Verify MQTT broker credentials and connectivity

### MQTT Not Receiving Messages

- Check MQTT broker logs
- Verify the topic prefix in your configuration
- Ensure the service has successfully connected (check logs)

### Viewing Logs

```bash
# Docker Compose
docker-compose logs -f ezr2mqtt

# Docker
docker logs -f ezr2mqtt
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under MIT.

## Support

If you encounter any issues or have questions, please [open an issue](https://github.com/chrishrb/ezr2mqtt/issues) on GitHub.
