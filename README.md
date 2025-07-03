# Simple TTS

This is a simple Text-to-Speech (TTS) server written in Go.

## Features

*   Provides a `/say` endpoint to synthesize text to speech.
*   Uses Google Cloud Text-to-Speech API.
*   Serves a simple HTML page to interact with the API.

## Usage

1.  **Run the server:**

    ```bash
    go run cmd/tts/main.go
    ```

2.  **Open your browser:**

    Navigate to `http://localhost:8080`

3.  **Synthesize speech:**

    Use the form on the page or send a GET request to the `/say` endpoint:

    ```
    http://localhost:8080/say?say=Hello, world!
    ```

## Configuration

*   `PORT`: The port the server listens on (default: `8080`).
