# rss2podcast
A locally hosted, AI generated podcast from an rss feed.

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

## Contributing
Contributions are welcome. Please open a pull request with your changes.

## License
This project is licensed under the terms of the MIT License.