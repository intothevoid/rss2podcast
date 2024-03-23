// Handle navbar functionality
const navbarItems = document.querySelectorAll('.navbar-item');
navbarItems.forEach(item => {
  item.addEventListener('click', () => {
    const page = item.dataset.page;
    loadPage(page);
  });
});

// Load initial page
loadPage('home');

// Function to load page content
function loadPage(page) {
  // Hide all pages
  const pages = document.querySelectorAll('.page');
  pages.forEach(page => {
    page.style.display = 'none';
  });
}


// Handle generate button functionality
const generateButton = document.getElementById('generate-button');
generateButton.addEventListener('click', () => {
  const rssUrl = textbox.value;
  // Perform generate action
});

// function to load rss url in textbox
function loadTopic(rssUrl) {
  const textbox = document.getElementById('rss-url');
  textbox.value = rssUrl;
}

// Handle generate button functionality
function generateRSS() {
  var topic = document.getElementById("rss-url").value;
  var url = "http://localhost:8080/generate/" + topic;
  
  // Disable UI elements
  disableUI();

  // Send GET request to the specified URL
  fetch(url)
    .then(response => response.json())
    .then(data => {
      // Handle the response data here
      console.log(data);

      // Enable UI elements
      enableUI();
    })
    .catch(error => {
      // Handle any errors here
      console.error(error);
      enableUI();
    });
}

// Function to disable all UI elements
function disableUI() {
  document.getElementById('rss-url').disabled = true;
  document.getElementById('generate-button').disabled = true;
  document.getElementById('loading').style.display = 'block';
  document.getElementById('rss-feed').style.display = 'none';
}

// Function to enable all UI elements
function enableUI() {
  document.getElementById('rss-url').disabled = false;
  document.getElementById('generate-button').disabled = false;
  document.getElementById('loading').style.display = 'none';
  document.getElementById('rss-feed').style.display = 'block';
}

