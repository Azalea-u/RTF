// app.js
import { loadHomePage } from './forum.js';
import { loadLoginPage } from './auth.js';
import { loadRegisterPage } from './auth.js';
import { loadForumPage } from './forum.js';

const routes = {
    '/': loadHomePage,
    '/login': loadLoginPage,
    '/register': loadRegisterPage,
    '/forum': loadForumPage,
};

function router() {
    const path = window.location.pathname;
    const app = document.getElementById('app');

    // Clear existing content
    app.innerHTML = '';

    // Load the appropriate page based on the route
    if (routes[path]) {
        routes[path]();
    } else {
        // Handle 404
        app.innerHTML = '<h2>404 - Page Not Found</h2>';
    }
}

// Handle navigation
document.addEventListener('click', (e) => {
    if (e.target.matches('[data-link]')) {
        e.preventDefault();
        const href = e.target.getAttribute('href');
        history.pushState(null, null, href);
        router();
    }
});

// Handle back/forward navigation
window.addEventListener('popstate', router);

// Initial load
router();