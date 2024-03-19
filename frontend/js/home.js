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

  // Show selected page
  const selectedPage = document.getElementById(page);
  selectedPage.style.display = 'block';
}

// Handle podcasts panel functionality
const podcastButtons = document.querySelectorAll('.podcast-button');
const textbox = document.getElementById('rss-textbox');
const generateButton = document.getElementById('generate-button');

podcastButtons.forEach(button => {
  button.addEventListener('click', () => {
    const rssUrl = button.dataset.rssUrl;
    textbox.value = rssUrl;
  });
});

// Handle generate button functionality
generateButton.addEventListener('click', () => {
  const rssUrl = textbox.value;
  // Perform generate action
});