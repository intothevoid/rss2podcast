// This file is intentionally left blank.
// Add your configuration UI code here.

// Handle save configuration button click
function saveConfig() {
    // Get the configuration values from the input field
    const subject = document.getElementById('subject').value;
    const podcaster = document.getElementById('podcaster').value;
    const rssMaxArticles = document.getElementById('rss_max_articles').value;
    const ollamaEndpoint = document.getElementById('ollama_end_point').value;
    const ollamaModel = document.getElementById('ollama_model').value;
    const ttsEngine = document.getElementById('tts_engine').value;
    const coquiUrl = document.getElementById('coqui_url').value;
    const kokoroUrl = document.getElementById('kokoro_url').value;
    const kokoroVoice = document.getElementById('kokoro_voice').value;
    const kokoroSpeed = document.getElementById('kokoro_speed').value;
    const kokoroFormat = document.getElementById('kokoro_format').value;

    // Add the configuration values to the object
    const config = {
        "subject": subject,
        "podcaster": podcaster,
        "rss_max_articles": rssMaxArticles,
        "ollama_endpoint": ollamaEndpoint,
        "ollama_model": ollamaModel,
        "tts": {
            "engine": ttsEngine,
            "coqui": {
                "url": coquiUrl
            },
            "kokoro": {
                "url": kokoroUrl,
                "voice": kokoroVoice,
                "speed": parseFloat(kokoroSpeed),
                "format": kokoroFormat
            }
        }
    };

    // Create a POST request to the server
    var url = "http://localhost:8080/configure/";

    // Send POST request to the specified URL with cache disabled
    fetch(url, { 
        method: "POST", 
        mode: 'no-cors',
        headers: {
            'Content-Type': 'text/plain',
        },
        body: JSON.stringify(config)
    })
        .then(response => response.json())
        .then(data => {
            console.log("Configuration saved")

            // Add div to the page with success message
            var div = document.createElement('div');
            div.innerHTML = "Configuration saved";
            document.body.appendChild(div);
        })
        .catch(error => {
            // Handle any errors here
            console.error(error);
        });
}

document.addEventListener('DOMContentLoaded', function() {
    // Add event listener for TTS engine selection
    document.getElementById('tts_engine').addEventListener('change', function() {
        const coquiContainer = document.getElementById('coqui_url_container');
        const kokoroContainer = document.getElementById('kokoro_url_container');
        
        if (this.value === 'coqui') {
            coquiContainer.style.display = 'block';
            kokoroContainer.style.display = 'none';
        } else if (this.value === 'kokoro') {
            coquiContainer.style.display = 'none';
            kokoroContainer.style.display = 'block';
        } else {
            coquiContainer.style.display = 'none';
            kokoroContainer.style.display = 'none';
        }
    });

    // Trigger the change event to set initial state
    document.getElementById('tts_engine').dispatchEvent(new Event('change'));
});
