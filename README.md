#
# File Sharing Web Application

This is a simple file-sharing web application built with Go. The application allows users to upload, download, and share files between users in a local network. Every upload generates a unique access token so multiple transfers can happen concurrently without overwriting each other. The refreshed web UI provides QR code sharing, a built-in QR scanner, and one-tap device notifications for trusted peers on the same network.

## Features

* **File Upload**: Users can upload multiple files in one transfer and categorize them (audio, photos, videos, documents, contacts, or mixed).
* **Private Storage**: Uploaded files are stored outside of the public `static/` directory and require an access token to download.
* **Share Links**: Each upload displays a share link (`/incoming?id=<id>&token=<token>`) that can be copied, sent, or encoded as a QR code, with optional PIN protection on downloads.
* **QR workflows**: The share screen renders a scannable QR code and ships with a built-in camera scanner so phones can join without typing URLs.
* **Device Notifications**: Register any device once and it will appear in the sender UI for one-tap delivery; receivers get pop-up prompts to accept or decline.
* **Automatic Device Discovery**: Every tab auto-registers itself (with a friendly name derived from the platform + browser) so the sender list updates instantly and ignores the current device.
* **File Metadata**: Displays metadata such as file name, size, and type via the UI or `/meta` endpoint.
* **Automatic Cleanup**: Expired transfers are removed on a background schedule to reclaim disk space.

## Directory Structure

```bash
├── devices/               # Device registry for push notifications
├── handlers/              # HTTP handlers and HTML rendering
├── storage/               # Concurrency-safe transfer store
├── static/
│   ├── css/
│   │   └── style.css
│   └── js/                # QR generator, scanner, and UI helpers
├── templates/             # HTML templates (main, send, share, receive, device)
├── utils/                 # Utility helpers
├── main.go                # Application entry point
├── go.mod
└── README.md
```

## Prerequisites

* **Go** (v1.20 or later)
* A modern web browser

## Installation

1. Clone the repository:

   ```bash
   git clone https://N-Jangra/share-go.git
   cd share-go
   ```

2. Build the Go server:

   ```bash
   go mod tidy
   go build -o server main.go
   ```

3. Run the server:

   ```bash
   ./server
   ```

   This will start the server on port `8080`.

4. Open your browser and go to `http://localhost:8080`.

## Usage

1. **Main Page**: Choose between sending or receiving files or opening the device listener.
2. **Upload File**: Click on "Send File", pick the content type, attach one or many files, and optionally require a PIN. When the upload finishes you will see a copyable share link, QR code, scanner overlay, and any registered devices you can notify instantly.
3. **Device Listener**: Open `/device`, register a friendly name once, and leave the page open. Senders will see your device in their list and you will receive pop-ups to accept or dismiss transfers.
4. **Receive File**: Recipients can use the share link, scan the QR code, or respond to the device pop-up. Accepting starts the download and removes the file from the server.
5. **Automation**: The `/meta?id=<id>&token=<token>` endpoint exposes file details if you need to script around transfers.

## Customization

You can customize the app by modifying the following:

* **CSS**: Customize the styles by editing `static/css/style.css`.
* **HTML**: Modify the templates located in the `templates` folder for changes to the UI.
* **File Handling**: Modify the handler logic under `handlers/` or the storage logic under `storage/` to add new behavior.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

#
