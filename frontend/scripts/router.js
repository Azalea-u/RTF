
import { checkAuth } from './utils.js';
import home from './pages/home.js';
import register from './pages/register.js';
import login from './pages/login.js';

const routes = {
    '/': home,
    '/register': register,
    '/login': login,
};

export function renderPage(path) {
    const app = document.getElementById('app');
    const page = routes[path] || login;
    app.innerHTML = '';
    app.appendChild(page());
    if (path === '/login' || path === '/register') {
        checkAuth().then(authenticated => {
            if (authenticated) {
                renderPage('/');
            }
        });
    } else {
        checkAuth().then(authenticated => {
            if (!authenticated) {
                renderPage('/login');
            }
        });
    }
}