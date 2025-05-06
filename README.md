# rss2podcast

<img src="resources/ui.jpg" width=600px></img>

A locally hosted, AI generated podcast from an rss feed. 

![workflow](https://github.com/intothevoid/rss2podcast/actions/workflows/go.yml/badge.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/intothevoid/rss2podcast)
![GitHub](https://img.shields.io/github/license/intothevoid/rss2podcast)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/intothevoid/rss2podcast)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/intothevoid/rss2podcast)

Powered by -

<img src="resources/rss.png" width=80px height=80px></img>
<img src="resources/ollama.png" width=80px height=80px></img>
<img src="resources/coqui.png" width=80px height=80px></img>
<img src="resources/kokoro.jpg" width=80px height=80px></img>

## Features

- RSS feed parsing and article extraction
- Article summarization using Ollama
- Text-to-speech conversion using multiple engines:
  - Kokoro TTS (Recommended)
  - MLX Audio TTS
  - Coqui TTS
- Podcast generation with customizable settings
- Web interface for configuration and control

## Requirements

- Go 1.21 or later
- Ollama (for article summarization)
- One of the following TTS engines:
  - Kokoro TTS (recommended)
  - MLX Audio TTS
  - Coqui TTS

## Installation

1. Clone the repository:
```bash
git clone https://github.com/intothevoid/rss2podcast.git
cd rss2podcast
```

2. Install dependencies:
```bash
go mod download
```

3. Configure the application by editing `config.yaml` or using the web interface.

## Configuration

The application can be configured using the web interface or by editing the `config.yaml` file. The following settings are available:

### RSS Settings
- `url`: The RSS feed URL to parse
- `max_articles`: Maximum number of articles to process
- `filters`: List of filters to apply to articles

### Ollama Settings
- `end_point`: The Ollama API endpoint
- `model`: The Ollama model to use for summarization

### Podcast Settings
- `subject`: The podcast subject
- `podcaster`: The podcaster name

### TTS Settings
- `engine`: The TTS engine to use ("kokoro", "mlx", or "coqui")
- `kokoro`: Kokoro TTS settings
  - `url`: The Kokoro TTS API endpoint
  - `voice`: The voice to use
  - `speed`: The speech speed (0.25 to 4.0)
  - `format`: The audio format (mp3, opus, flac, wav, pcm)
- `mlx`: MLX Audio TTS settings
  - `url`: The MLX Audio TTS API endpoint
  - `voice`: The voice to use
  - `speed`: The speech speed (0.5 to 2.0)
  - `format`: The audio format (mp3, wav)
- `coqui`: Coqui TTS settings
  - `url`: The Coqui TTS API endpoint

## Usage

1. Start the application:
```bash
go run cmd/rss2podcast/main.go
```

2. Access the web interface at `http://localhost:8080`

3. Configure the application using the web interface or edit `config.yaml`

4. The application will:
   - Parse the RSS feed
   - Extract and summarize articles
   - Convert the summary to audio using the selected TTS engine
   - Generate a podcast file

## TTS Engines

### Kokoro TTS (Recommended)
Kokoro TTS offers OpenAI-compatible speech synthesis with support for multiple voices and formats. It provides excellent quality with low latency.

### MLX Audio TTS
MLX Audio TTS is a powerful text-to-speech engine that provides high-quality speech synthesis with support for multiple voices and formats. It offers additional features like direct audio playback and output folder management.

### Coqui TTS
Coqui TTS provides high-quality speech synthesis with support for multiple voices and formats.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Ollama](https://ollama.ai/) for the LLM API
- [Kokoro TTS](https://github.com/kokoro-tts/kokoro) for the TTS engine
- [MLX Audio TTS](https://github.com/mlx-audio/mlx-tts) for the TTS engine
- [Coqui TTS](https://github.com/coqui-ai/TTS) for the TTS engine

## How it works
The application reads an rss feed, extracts the articles and summarises them. 

RSS + Ollama + TTS = Podcast

### RSS
The application reads an rss feed and extracts the articles. Each of these articles are then processed by scraping the article content.

### Ollama
The application uses a locally hosted version of Ollama. The Ollama API is used to summarise the article content. Default model used is mistral:7b

### TTS
The summarised article content is then converted into an audio podcast using the Coqui TTS API.

## Dependencies

This project requires the following dependencies to be installed on your system. 

### Ollama

You can install the Ollama server by following the instructions on the [official website](https://ollama.com).

Ollama needs to be running on your local machine for the application to work. The application is configured to use the default Ollama server URL `http://localhost:11434/api/generate`. This can be changed via the config.yaml file.

### ffmpeg

`ffmpeg` is a command-line tool for handling multimedia files. It is used to convert the generated audio files to the MP3 format.

#### macOS

You can use Homebrew to install `ffmpeg` on macOS:

```bash
brew install ffmpeg
```

#### Windows

1. Download the `ffmpeg` build for Windows from the [official website](https://ffmpeg.org/download.html).
2. Extract the downloaded ZIP file.
3. Add the `bin` directory from the extracted folder to your system's PATH.

#### Linux

The installation command depends on your Linux distribution.

##### Ubuntu/Debian

```bash
sudo apt update
sudo apt install ffmpeg
```

### Kokoro TTS (Recommended)

Kokoro TTS is a text-to-speech synthesis system that uses deep learning to create human-like speech from text. You can install the Kokoro TTS server by following the instructions on the [official website](https://github.com/nazdridoy/kokoro-tts).

#### Docker:
  
Create a docker-compose.yml file and add the following:

```yaml
services:
kokoro-fastapi-cpu:
    ports:
        - 8880:8880
    image: ghcr.io/remsky/kokoro-fastapi-cpu:latest # or v0.2.3 for last stable version
```

Start the server by running the following command:

```bash
docker compose up -d
```

This will start the Kokoro TTS server on port 8880. The server provides a REST API for text-to-speech conversion.

### Coqui TTS

Coqui TTS is a text-to-speech synthesis system that uses deep learning to create human-like speech from text. You can install the Coqui TTS server by following the instructions on the [official website](https://coqui.ai/tts).

#### Docker

Start the container by using the following command:

```bash
docker run -d -p 5002:5002 --platform linux/amd64 --entrypoint /usr/local/bin/tts-server ghcr.io/coqui-ai/tts-cpu --model_name tts_models/en/ljspeech/vits
```

### MLX Audio TTS

MLX Audio TTS is a text-to-speech synthesis system that uses deep learning to create human-like speech from text. You can install the MLX Audio TTS server by following the instructions on the [official website](https://github.com/mlx-audio/mlx-tts).

#### Docker

As of this writing, MLX Audio TTS needs to be run locally as Docker does not allow GPU access on Apple Silicon.

```bash
# Install the package
pip install mlx-audio

# Create a virtual environment
python -m venv venv

# Activate the virtual environment
source venv/bin/activate

# Install the dependencies
pip install -r requirements.txt

# Run the server
mlx_audio.server
```
rss2podcast will automatically request the MLX Audio TTS server to generate the audio file.

## Testing
To run the tests, use the following command:
```bash
go test ./...
```