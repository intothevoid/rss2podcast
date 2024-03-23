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
hideSpinner();

// Function to load page content
function loadPage(page) {
  // Hide all pages
  const pages = document.querySelectorAll('.page');
  pages.forEach(page => {
    page.style.display = 'none';
  });
}

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
  showSpinner();

  // Send GET request to the specified URL with cache disabled
  fetch(url, { cache: 'no-store' })
    .then(response => response.blob())
    .then(blob => {
      // Create a URL for the blob object
      const fileUrl = URL.createObjectURL(blob);

      // Create a hyperlink element
      const link = document.createElement('a');
      link.href = fileUrl;
      link.download = 'podcast.mp3';
      link.textContent = 'Download';

      // Append the hyperlink to the home page
      const downloadDiv = document.getElementById('download-url');
      downloadDiv.appendChild(link);

      // Enable UI elements
      enableUI();
      hideSpinner();
    })
    .catch(error => {
      // Handle any errors here
      console.error(error);
      enableUI();
      hideSpinner();
    });
}

// Function to disable all UI elements
function disableUI() {
  document.getElementById('rss-url').style.display = 'none';
  document.getElementById('generate-button').style.display = 'none';
  document.getElementById('download-url').style.display = 'none';
}

// Function to enable all UI elements
function enableUI() {
  document.getElementById('rss-url').style.display = 'block';
  document.getElementById('generate-button').style.display = 'block';
  document.getElementById('download-url').style.display = 'block';
}

// Function to show spinner
function showSpinner() {
  document.getElementById('spinner').style.display = 'block';
}

// Function to hide spinner
function hideSpinner() {
  document.getElementById('spinner').style.display = 'none';
}