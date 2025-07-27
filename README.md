# OSC2Hue

A Go application that bridges OSC (Open Sound Control) messages to Philips Hue lighting control.

## Description

This project provides a bridge between OSC messages and Philips Hue smart lights, allowing you to control your lighting setup through OSC commands.

## Features

- **🚀 Zero-config setup**: Automatic bridge discovery and authentication (just press the bridge button when prompted)
- **💡 Auto-discovery**: Finds lights from your Hue bridge at startup
- **🏷️ Dual addressing**: Supports both UUID and numeric light IDs
- **🌍 Global controls**: Commands to control all lights at once
- **🎨 CIE XY colors**: Use of Philips Hue colorimetry ([convert to RGB](https://viereck.ch/hue-xy-rgb/))
- **⚡ Transition support**: Smooth transitions with duration control
- **🎵 Tidal Cycles integration**: Ready-to-use examples for live coding

### Using with OSC Applications

You can send these messages to control your lights from any OSC-capable application such as:
- **[Max/MSP](https://cycling74.com/products/max)** - Visual programming for multimedia (built-in OSC objects)
- **[Pure Data (Pd)](https://puredata.info/)** - Real-time audio/visual programming (built-in or with [OSC externals](https://github.com/pd-externals/osc))
- **[SuperCollider](https://supercollider.github.io/)** - Audio synthesis platform (built-in NetAddr/sendMsg)
- **[TidalCycles](https://tidalcycles.org/)** - Live coding music pattern language (with [example implementation](examples/TIDAL_INTEGRATION.md))
- **[Processing](https://processing.org/)** - Creative coding environment (with [oscP5 library](https://sojamo.de/libraries/oscp5/))
- **[openFrameworks](https://openframeworks.cc/)** - Creative coding toolkit (with ofxOsc addon)

## Usage

### Prerequisites

- Philips Hue Bridge on your network
- That's it! No additional software required.

### Quick Start

1. **Download the latest release** for your operating system from the [Releases page](../../releases):
   - **Linux (x64)**: `osc2hue-linux-amd64`
   - **macOS (Intel)**: `osc2hue-darwin-amd64` 
   - **macOS (Apple Silicon)**: `osc2hue-darwin-arm64`
   - **Windows (x64)**: `osc2hue-windows-amd64.exe`

2. **Make it executable** (Linux/macOS):
   ```bash
   chmod +x osc2hue-*
   ```

3. **Run the application**:
   ```bash
   ./osc2hue-linux-amd64
   # or
   ./osc2hue-darwin-amd64
   # or on Windows:
   # osc2hue-windows-amd64.exe
   ```

4. **Press the link button on your Hue bridge when prompted**

5. **Start sending OSC messages to control your lights!**
   The server will be listening on port 8080 (default) for OSC messages.

That's it! The application will automatically:
- Discover your Hue bridge on the network
- Authenticate with your bridge (via button press)
- Save the configuration for future use
- Discover and map your lights

**Note:** If authentication times out or fails, simply run the application again - it will remember your bridge IP and only ask for authentication.

### OSC Message Format

The bridge accepts OSC messages in the following formats:

#### Light Control Commands
- **Turn light on/off:**
  ```
  /hue/light/{id}/on {0|1} [duration_ms]
  ```
  - `{id}`: Light ID (UUID or number 1,2,3... based on discovery order)
  - `{0|1}`: 0 = off, 1 = on
  - `[duration_ms]`: Optional transition duration in milliseconds

- **Set brightness:**
  ```
  /hue/light/{id}/brightness {value} [duration_ms]
  ```
  - `{id}`: Light ID (UUID or number)
  - `{value}`: 0.0-1.0 (float)
  - `[duration_ms]`: Optional transition duration in milliseconds

- **Set color:**
  ```
  /hue/light/{id}/color {x} {y} [duration_ms]
  ```
  - `{id}`: Light ID (UUID or number)
  - `{x}`, `{y}`: CIE XY color coordinates (0.0-1.0)
  - `[duration_ms]`: Optional transition duration in milliseconds

- **Unified set command (with optional parameters):**
  ```
  /hue/light/{id}/set {x|-1} [y|-1] [brightness|-1] [duration_ms|-1]
  ```
  - `{id}`: Light ID (UUID or number)
  - `{x|-1}`, `{y|-1}`: CIE XY color coordinates (0.0-1.0) or -1 to skip color change
  - `{brightness|-1}`: 0.0-1.0 (float) or -1 to skip brightness change
  - `{duration_ms|-1}`: Transition duration in milliseconds or -1 to skip
  - **Note:** When setting color, both X and Y must be provided (or both -1 to skip)

#### Global Commands
- **Control all lights:**
  ```
  /hue/all/on {0|1} [duration_ms]
  /hue/all/brightness {value} [duration_ms]
  /hue/all/color {x} {y} [duration_ms]
  /hue/all/set {x|-1} [y|-1] [brightness|-1] [duration_ms|-1]
  ```
  - Same parameters as individual light commands
  - `/set` command supports null values using -1 to skip parameters
  - `[duration_ms]`: Optional transition duration in milliseconds

#### Examples
```bash
# Turn light 1 on
/hue/1/on 1

# Turn light 1 on with 2 second transition
/hue/1/on 1 2000

# Set light 2 to 50% brightness with 1 second transition
/hue/2/brightness 0.5 1000

# Set light 3 to warm white color with 500ms transition
/hue/3/color 0.4 0.4 500

# Set light 1 to cool blue at 30% brightness with smooth 2 second transition
/hue/1/set 0.15 0.06 0.3 2000

# Set all lights to warm white at 80% brightness instantly
/hue/all/set 0.4 0.4 0.8

# Set all lights to green at full brightness with 1 second transition
/hue/all/set 0.3 0.6 1.0 1000

# Set all lights to cool blue color with 500ms transition
/hue/all/color 0.15 0.06 500

### Unified Set Commands (with null value support)

The `/set` commands allow you to modify only specific properties by using `-1` for null/skip values:

```bash
# Set only color (x=0.4, y=0.5), keep current brightness and no transition
/hue/1/set 0.4 0.5 -1 -1

# Set only brightness (60%), keep current color and no transition  
/hue/2/set -1 -1 0.6 -1

# Set color and brightness with transition, but skip duration
/hue/all/set 0.3 0.6 0.8 2000

# Change only transition duration, keep current color/brightness
/hue/all/set -1 -1 -1 1500

# Set color only with transition duration
/hue/3/set 0.2 0.7 -1 3000
```

**Note:** When setting color, both X and Y coordinates must be provided (or both set to -1 to skip color entirely).

#### More Examples
```bash
# Turn all lights off instantly
/hue/all/on 0

# Turn all lights off with 3 second transition
/hue/all/on 0 3000

# Set all lights to 75% brightness with smooth 1.5 second transition
/hue/all/brightness 0.75 1500

# Set all lights to warm red color instantly
/hue/all/color 0.6 0.3

# Set all lights to 30% brightness over 5 seconds using the unified command
/hue/all/set -1 -1 0.3 5000
```

**Note:** The application discovers actual lights from your bridge at startup and supports both UUID addressing (`/hue/abc-123-def/on`) and numeric addressing (`/hue/1/on`) for convenience.

## Configuration

### Configuration File

A `config.json` file will automatically be created in the project root with the following structure:

```json
{
  "osc": {
    "host": "0.0.0.0",
    "port": 8080
  },
  "hue": {
    "bridge_ip": "",
    "api_key": ""
  }
}
```

**Note:** Leave `bridge_ip` empty for automatic discovery, or set a specific IP address to skip discovery. Leave `api_key` empty for automatic authentication.

### Configuration Options

#### OSC Settings
- **`host`**: IP address to bind the OSC server to
  - `"0.0.0.0"` - Listen on all interfaces
  - `"127.0.0.1"` - Listen only on localhost
  - `"192.168.1.10"` - Listen on specific IP
- **`port`**: UDP port number for OSC messages (default: 8080)

#### Hue Settings
- **`bridge_ip`**: IP address of your Philips Hue Bridge
- **`api_key`**: Authorized API key for Hue Bridge API access

### Getting Hue Bridge Credentials

#### Automatic Bridge Discovery
The application automatically discovers Hue bridges on your network at startup! Simply leave the `bridge_ip` field empty or set to the default value:

```json
{
  "osc": {
    "host": "0.0.0.0",
    "port": 8080
  },
  "hue": {
    "bridge_ip": "",
    "api_key": ""
  }
}
```

#### Manual Bridge IP Setup
If automatic discovery fails, you can manually find your bridge IP:

Visit https://discovery.meethue.com/

#### Getting an API Key (Automatic Authentication)
The application can automatically authenticate with your bridge! Simply run the application and it will:

1. **Discover your bridge** automatically (if needed)
2. **Prompt you to press the link button** on your bridge
3. **Automatically obtain and save** your API key

```bash
./osc2hue
# Output:
# Discovering Hue bridges...
# Found Hue bridge at 192.168.1.74
# 🔗 Press the link button on your Hue bridge now...
# ✅ Authentication successful!
# Configuration saved with new API key
```

#### Manual Authentication (Alternative)
If you prefer manual setup, you can still get an API key manually:

1. **Press the link button** on your Hue Bridge (you have 30 seconds)
2. **Send a registration request:**
   ```bash
   curl -X POST http://{bridge_ip}/api \
     -H "Content-Type: application/json" \
     -d '{"devicetype":"osc2hue#mydevice"}'
   ```
3. **Copy the API key** from the response and add it to your config.json

## Development

### Prerequisites for Development

- Go 1.23 or later
- Git

### Building from Source

- **Clone this repository:**
   ```bash
   git clone <repository-url>
   cd osc2hue
   ```

- **Install dependencies:**
   ```bash
   go mod tidy
   ```

- **Run tests:**
   ```bash
   go test ./...
   ```

- **Build the application:**
   ```bash
   go build -o osc2hue .
   ```

- **Run the application:**
   ```bash
   ./osc2hue
   ```

- **Cross-compile for different platforms:**
   ```bash
   GOOS=linux GOARCH=amd64 go build -o osc2hue-linux-amd64 .
   GOOS=darwin GOARCH=amd64 go build -o osc2hue-darwin-amd64 .
   GOOS=darwin GOARCH=arm64 go build -o osc2hue-darwin-arm64 .
   GOOS=windows GOARCH=amd64 go build -o osc2hue-windows-amd64.exe .
   ```

### Project Structure
```
osc2hue/
├── internal/
│   ├── config/           # Configuration management
│   ├── hue/             # Hue bridge integration
│   └── osc/             # OSC server implementation
├── examples/            # Example code and integrations
│   ├── TIDAL_INTEGRATION.md     # Tidal Cycles guide
│   ├── tidal-simple-osc.tidal   # Tidal examples
│   └── *.go             # Test clients
├── handlers.go          # OSC message handlers
├── main.go             # Main application entry point
├── go.mod              # Go module definition
└── README.md           # This file
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
