#
# File Sharing Web Application

This is a simple file-sharing web application built with Go. The application allows users to upload, download, and share files between users in a local network. The user interface is built using HTML templates and styled with CSS.

## Features

* **File Upload**: Users can upload files to the server.
* **File Metadata**: Displays metadata such as file name, size, and type.
* **File Download**: Users can download files shared by others.
* **Incoming Files**: Users can accept or decline files sent by others.
* **Responsive Interface**: The app is designed to be simple and easy to use.

## Directory Structure

```bash
├── cmd/
│   └── server/
│       └── main.go           # Main entry point of the server
├── handlers/
│   ├── home.go               # Handlers for home page
│   ├── upload.go             # Handlers for file upload
│   ├── meta.go               # Handler for file metadata
│   ├── file.go               # Handler for file download
│   ├── receive.go            # Handlers for incoming, accept, and decline file
│   └── home.go               # Home page handler
├── static/
│   └── css/
│       └── style.css         # Styles for the application
│   └── uploaded/             # Directory where uploaded files are stored
├── templates/
│   ├── upload.html           # Upload file template
│   ├── incoming.html         # Incoming file accept/decline template
│   └── main.html             # Main page template
├── utils/
│   └── fs.go                 # Utility for file system operations
│   └── ip.go                 # Utility to get local IP address
└── README.md                 # Project documentation
```

## Prerequisites

* **Go** (v1.16 or later)
* A modern web browser

## Installation

1. Clone the repository:

   ```bash
   git clone https://N-Jangra/share-go.git
   cd share-go
   ```

2. Build the Go server:

   ```bash
   go build -o server main.go
   ```

3. Run the server:

   ```bash
   ./server
   ```

   This will start the server on port `8080`.

4. Open your browser and go to `http://localhost:8080`.

## Usage

1. **Main Page**: Choose between sending or receiving files.
2. **Upload File**: Click on "Send File" and select a file to upload.
3. **Accept/Decline File**: If you're receiving a file, you'll see metadata and the option to either accept or decline the file.
4. **File Metadata**: After uploading, you can view the file's name, size, and type.

## Customization

You can customize the app by modifying the following:

* **CSS**: Customize the styles by editing `static/css/style.css`.
* **HTML**: Modify the templates located in the `templates` folder for changes to the UI.
* **File Handling**: Modify the file-handling logic in `handlers/` to support additional functionality or custom file management.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

#
