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
        kokoro_format: document.getElementById('kokoro_format').value
    };

    console.log('Sending configuration:', config);

    fetch('http://localhost:8080/configure/', {
        method: 'POST',
        mode: 'cors',
        credentials: 'omit',
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        },
        body: JSON.stringify(config)
    })
    .then(response => {
        console.log('Response received:', response);
        console.log('Response status:', response.status);
        console.log('Response headers:', response.headers);
        
        if (!response.ok) {
            return response.text().then(text => {
                console.error('Server error response:', text);
                throw new Error(text || `HTTP error! status: ${response.status}`);
            });
        }
        return response.text();
    })
    .then(text => {
        console.log('Success response text:', text);
        try {
            const data = JSON.parse(text);
            console.log('Parsed response:', data);
            alert(data.message || 'Configuration saved successfully!');
        } catch (e) {
            console.log('Raw response (not JSON):', text);
            alert('Configuration saved successfully!');
        }
    })
    .catch(error => {
        console.error('Detailed error:', error);
        console.error('Error stack:', error.stack);
        
        if (!window.navigator.onLine) {
            alert('You are offline. Please check your internet connection.');
            return;
        }

        if (error.message.includes('Failed to fetch')) {
            alert('Cannot connect to the server. Please ensure the server is running at http://localhost:8080');
        } else {
            alert('Error saving configuration: ' + error.message);
        }
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
