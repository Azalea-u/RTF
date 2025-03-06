import { renderPage } from "../router.js";
import { showAlert } from "../utils.js";

export default function home() {
    const container = document.createElement('div');
    container.innerHTML = `
        <nav class="navbar">
            <h1>Welcome to the Forum</h1>
            <button id="logout-button">Logout</button>
        </nav>
        <div id="content">
            <h2>This is the home page</h2>
        </div>
    `;

    const alert = document.getElementById('alert');
    container.querySelector('#logout-button').addEventListener('click', async () => {
        const response = await fetch('/api/logout', {
            method: 'POST',
            credentials: 'include',
        });

        if (response.ok) {
            showAlert('Logout successful', 'success');
            setTimeout(() => {
                renderPage('/login');
            }, 1000);
        } else {
            showAlert('Logout failed', 'error');
        }
    });

    return container;
}