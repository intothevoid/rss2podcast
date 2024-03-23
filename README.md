# rss2podcast

<img src="resources/ui.jpg" width=600px></img>

A locally hosted, AI generated podcast from an rss feed. 

![workflow](https://github.com/intothevoid/rss2podcast/actions/workflows/go.yml/badge.svg)

Powered by -

<img src="resources/rss.png" width=80px height=80px></img>
<img src="resources/ollama.png" width=80px height=80px></img>
<img src="resources/coqui.png" width=80px height=80px></img>

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

### Coqui TTS

Coqui TTS is a text-to-speech synthesis system that uses deep learning to create human-like speech from text. You can install the Coqui TTS server by following the instructions on the [official website](https://coqui.ai/tts).

#### Docker

Start the container by using the following command:

```bash
docker run -d -p 5002:5002 --platform linux/amd64 --entrypoint /usr/local/bin/tts-server ghcr.io/coqui-ai/tts-cpu --model_name tts_models/en/ljspeech/vits
```

## Installation 

Clone the repository and navigate into the directory:

```bash 
    git clone https://github.com/yourusername/your-repo.git
    cd your-repo
```

Then, install the dependencies:
```bash
go mod download
```

## Usage
To run the application, navigate to the cmd/rss2podcast directory and run:
```bash
go run main.go
```

## Testing
To run the tests, use the following command:
```bash
go test ./...
```

## Configuration

The application's configuration is stored in a `config.yaml` file. Here's what each section does:

### Podcast

This section contains information about the podcast.

```yaml
podcast:
  subject: "News" # The subject of the podcast
  podcaster: "Cody" # The name of the podcaster
```

### RSS

This section contains information about the RSS feed.

```yaml
rss:
  url: "https://www.reutersagency.com/feed/?taxonomy=best-topics&post_type=best" # The URL of the RSS feed
  max_articles: 10 # The maximum number of articles to fetch from the RSS feed
  filters: # Keywords to filter articles by
    - "Daily"
    - "Weekly"
```

### Ollama

This section contains information about the Ollama service.

```yaml
ollama:
  end_point: "http://localhost:11434/api/generate" # The URL of the Ollama service
  model: "mistral:7b" # The model used by the Ollama service
```

### TTS

This section contains information about the Text-to-Speech (TTS) service.

```yaml
tts:
  url: "http://localhost:5002/api/tts" # The URL of the TTS service
```

You can modify these values to suit your needs. Remember to restart the application after making changes to the configuration file.

## Contributing
Contributions are welcome. Please open a pull request with your changes.

## License
This project is licensed under the terms of the MIT License.