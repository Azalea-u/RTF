import { renderPage } from "../router.js";

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
            alert.classList.remove('hidden');
            alert.textContent = 'Logout successful';
            alert.classList.add('success');
            setTimeout(() => {
                renderPage('/login');
            }, 2000);
        } else {
            alert.classList.remove('hidden');
            alert.textContent = 'Logout failed';
            alert.classList.add('error');
        }
    });

    return container;
}