# Share-Go

A modern, secure file-sharing web application built with Go for local network transfers. Share files between devices with QR codes, device notifications, and automatic cleanup - perfect for sharing photos, documents, and media across your home network.

## Features

* **Concurrent Transfers**: Multiple file transfers can happen simultaneously with unique access tokens
* **Multi-File Uploads**: Upload multiple files at once with categorization (audio, photos, videos, documents, contacts, or mixed)
* **Secure Storage**: Files stored privately outside the web root with token-based access control
* **Shareable Links**: Generate secure share links with optional PIN protection
* **QR Code Sharing**: Built-in QR code generation and scanner for easy mobile access
* **Device Notifications**: One-tap notifications to registered devices on your network
* **Auto-Discovery**: Automatic device registration and discovery across the local network
* **File Metadata**: Rich metadata display including file size, type, and transfer details
* **Automatic Cleanup**: Background process removes expired transfers to reclaim disk space
* **Responsive UI**: Modern web interface that works on desktop and mobile devices

## Architecture

Built with a clean, modular architecture:

- **HTTP Server**: Standard Go net/http with custom handlers
- **Storage Layer**: Concurrent-safe file storage with metadata tracking
- **Device Registry**: In-memory device management for notifications
- **Template System**: HTML templates with dynamic content rendering
- **Background Tasks**: Automatic cleanup with configurable TTL (Time To Live)

## Directory Structure

```
├── devices/
│   └── registry.go        # Device registry and notification system
├── handlers/
│   ├── device.go          # Device-related HTTP handlers
│   ├── file.go            # File serving handlers
│   ├── helpers.go         # Helper functions for handlers
│   ├── home.go            # Home page handler
│   ├── meta.go            # File metadata handlers
│   ├── receive.go         # File receiving handlers
│   ├── server.go          # Main server setup and routing
│   └── upload.go          # File upload handlers
├── storage/
│   └── store.go           # File storage abstraction and management
├── static/
│   ├── css/
│   │   └── style.css      # Application stylesheets
│   └── js/
│       ├── app.js         # Main application JavaScript
│       ├── device.js      # Device management scripts
│       ├── jsqr.js        # QR code scanning library
│       ├── qrcode.js      # QR code generation library
│       ├── qrscanner.js   # QR scanner interface
│       ├── receive.js     # File receiving interface
│       └── share.js       # File sharing interface
├── templates/
│   ├── device.html        # Device registration page
│   ├── main.html          # Main application page
│   ├── receive.html       # File receiving page
│   ├── send.html          # File sending page
│   └── share.html         # File sharing page
├── uploads/               # Temporary file storage directory
├── utils/
│   ├── ensure.go          # File system utilities
│   ├── ip.go              # IP address detection utilities
│   └── meta.go            # File metadata utilities
├── issues.md              # Known issues and improvements
├── main.go                # Application entry point
├── go.mod                 # Go module definition
├── README.md              # This file
└── LICENSE                # GPL v3 license
```

## Prerequisites

* **Go** (v1.24.2 or later)
* **Modern web browser** with JavaScript enabled
* **Camera access** (optional, for QR code scanning)

## Quick Start

1. **Clone and build**:
   ```bash
   git clone https://github.com/N-Jangra/share-go.git
   cd share-go
   go mod tidy
   go build -o server main.go
   ```

2. **Run the server**:
   ```bash
   ./server
   ```
   The server will display available URLs (typically `http://localhost:8080` and network IPs).

3. **Access the application**:
   Open your browser and navigate to the displayed URL.

## Usage Guide

### For Senders
1. **Upload Files**: Click "Send File", select file category, choose files, and optionally set a PIN
2. **Share Options**: After upload, you'll get:
   - A shareable link with unique token
   - QR code for easy scanning
   - List of registered devices for direct notifications

### For Receivers
1. **Receive Options**:
   - Click the share link
   - Scan the QR code
   - Accept device notifications
2. **Download**: Files download automatically and are removed from the server after transfer

### Device Registration
1. Open `/device` in a browser tab
2. Register a friendly name
3. Keep the tab open to receive notifications

### API Endpoints
- `GET /` - Main application page
- `POST /uploadFile` - Upload files
- `GET /incoming?id=<id>&token=<token>` - Access shared files
- `GET /meta?id=<id>&token=<token>` - Get file metadata
- `GET /device` - Device registration page
- `GET /api/devices` - List registered devices
- `POST /api/devices/register` - Register a device
- `POST /api/devices/notify` - Send notification to device

## Configuration

The application uses sensible defaults but can be customized:

- **Port**: Change `port` constant in `main.go` (default: 8080)
- **Upload Directory**: Modify `uploadsDir` constant (default: "uploads")
- **Transfer TTL**: Adjust `transferTTL` for expiration time (default: 2 hours)
- **Cleanup Interval**: Change `cleanupInterval` for cleanup frequency (default: 15 minutes)

## Development

### Building
```bash
go build -o server main.go
```

### Running
```bash
./server
```

### Customization
- **UI**: Modify templates in `templates/` and styles in `static/css/`
- **Logic**: Extend handlers in `handlers/` or storage in `storage/`
- **Features**: Add new functionality by implementing additional endpoints

## Security Notes

- Files are stored with access tokens for security
- Optional PIN protection for sensitive transfers
- Automatic cleanup prevents disk space issues
- Network-only access (no external internet exposure recommended)

## Technical Details

- **Concurrency**: Thread-safe storage handles multiple simultaneous transfers
- **Cleanup**: Background goroutine removes expired files every 15 minutes
- **Device Discovery**: Automatic registration with platform + browser detection
- **File Limits**: 10MB per file (configurable in upload handler)
- **Network**: Designed for local network use with automatic IP detection

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.
