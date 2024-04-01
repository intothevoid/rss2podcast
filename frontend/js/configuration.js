// This file is intentionally left blank.
// Add your configuration UI code here.

// Handle save configuration button click
function saveConfig() {
    // Get the configuration values from the input field
    var subject = document.getElementById('subject').value;
    var podcaster = document.getElementById('podcaster').value;
    var rssMaxArticles = document.getElementById('rss_max_articles').value;
    var ollamaEndpoint = document.getElementById('ollama_end_point').value;
    var ollamaModel = document.getElementById('ollama_model').value;
    var ttsUrl = document.getElementById('tts_url').value;

    // Add the configuration values to the object
    var config = {
        "subject": subject,
        "podcaster": podcaster,
        "rss_max_articles": rssMaxArticles,
        "ollama_endpoint": ollamaEndpoint,
        "ollama_model": ollamaModel,
        "tts_url": ttsUrl
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
        body: JSON.stringify(config) // Directly stringify the config object
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
