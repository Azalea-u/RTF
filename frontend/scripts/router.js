import { checkAuth } from './utils.js';
import home from './pages/home.js';
import register from './pages/register.js';
import login from './pages/login.js';

const routes = {
    '/': home,
    '/register': register,
    '/login': login,
};

export async function renderPage(path) {
    const app = document.getElementById('app');
    const page = routes[path] || login;

    app.innerHTML = '';

    const pageNode = await page();
    if (pageNode instanceof Node) {
        app.appendChild(pageNode);
    } else {
        console.error('Rendered page is not a valid Node:', pageNode);
    }

    if (path === '/login' || path === '/register') {
        const authenticated = await checkAuth();
        if (authenticated) {
            renderPage('/');
        }
    } else {
        const authenticated = await checkAuth();
        if (!authenticated) {
            renderPage('/login');
        }
    }
}