# Gemini Guidelines

This document provides instructions for interacting with the `simple-tts` project using the Gemini CLI.

## Project Overview

`simple-tts` is a web server written in Go that provides a Text-to-Speech (TTS) service. It uses the Google Cloud Text-to-Speech API to convert text into spoken audio.

The server has a single API endpoint, `/say`, which takes a `say` query parameter. The value of this parameter is synthesized into speech and returned as an MP3 audio file.

## Development

To run the server:

```bash
go run cmd/tts/main.go
```

The server will start on port 8080 by default, but this can be configured using the `PORT` environment variable.
