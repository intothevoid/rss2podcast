// This file is intentionally left blank.
// Add your configuration UI code here.

// Handle save configuration button click
function saveConfig() {
    console.log('Starting configuration save...');
    
    const config = {
        subject: document.getElementById('subject').value,
        podcaster: document.getElementById('podcaster').value,
        rss_max_articles: document.getElementById('rss_max_articles').value,
        ollama_endpoint: document.getElementById('ollama_end_point').value,
        ollama_model: document.getElementById('ollama_model').value,
        tts_engine: document.getElementById('tts_engine').value,
        coqui_url: document.getElementById('coqui_url').value,
        kokoro_url: document.getElementById('kokoro_url').value,
        kokoro_voice: document.getElementById('kokoro_voice').value,
        kokoro_speed: document.getElementById('kokoro_speed').value,
        kokoro_format: document.getElementById('kokoro_format').value,
        mlx_url: document.getElementById('mlx_url').value,
        mlx_voice: document.getElementById('mlx_voice').value,
        mlx_speed: document.getElementById('mlx_speed').value,
        mlx_format: document.getElementById('mlx_format').value
    };

    console.log('Sending configuration:', config);

    // Send the configuration to the server
    fetch('http://localhost:8080/configure/', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(config)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => {
        console.log('Configuration saved successfully:', data);
        alert('Configuration saved successfully!');
    })
    .catch(error => {
        console.error('Error saving configuration:', error);
        alert('Error saving configuration: ' + error.message);
    });
}

// Add event listener for page load
document.addEventListener('DOMContentLoaded', function() {
    console.log('Configuration page loaded');
    
    // Test server connectivity
    fetch('http://localhost:8080/configure/', {
        method: 'OPTIONS',
        mode: 'cors',
        credentials: 'omit'
    })
    .then(response => {
        console.log('Server is reachable, OPTIONS response:', response);
    })
    .catch(error => {
        console.error('Server connectivity test failed:', error);
    });

    // Set kokoro as default TTS engine
    const ttsEngineSelect = document.getElementById('tts_engine');
    ttsEngineSelect.value = 'kokoro';

    // Add event listener for TTS engine selection
    ttsEngineSelect.addEventListener('change', function() {
        const coquiContainer = document.getElementById('coqui_url_container');
        const kokoroContainer = document.getElementById('kokoro_url_container');
        const mlxContainer = document.getElementById('mlx_url_container');
        
        if (this.value === 'coqui') {
            coquiContainer.style.display = 'block';
            kokoroContainer.style.display = 'none';
            mlxContainer.style.display = 'none';
        } else if (this.value === 'kokoro') {
            coquiContainer.style.display = 'none';
            kokoroContainer.style.display = 'block';
            mlxContainer.style.display = 'none';
        } else if (this.value === 'mlx') {
            coquiContainer.style.display = 'none';
            kokoroContainer.style.display = 'none';
            mlxContainer.style.display = 'block';
        } else {
            coquiContainer.style.display = 'none';
            kokoroContainer.style.display = 'none';
            mlxContainer.style.display = 'none';
        }
    });

    // Trigger the change event to set initial state
    ttsEngineSelect.dispatchEvent(new Event('change'));
});
