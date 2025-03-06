import { renderPage } from "../router.js";

export default function register() {
    const container = document.createElement('div');
    container.innerHTML = `
        <h1>Login</h1>
        <form id="login-form">
            <input type="text" id="username-or-email" placeholder="Username or Email" required>
            <input type="password" id="password" placeholder="Password" required>
            <button type="submit">Login</button>
            <p>Don't have an account? <a href="#" id="register-link">Register</a></p>
        </form>
    `;

    container.querySelector('#login-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData(e.target);
        const data = Object.fromEntries(formData.entries());
        const alert = document.getElementById('alert');

        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify(data),
        });

        if (response.ok) {
            alert.classList.remove('hidden');
            alert.textContent = 'Login successful';
            alert.classList.add('success');
            e.target.reset();
            setTimeout(() => {
                renderPage('/');
            }, 2000);
        } else {
            const error = await response.json();
            alert.classList.remove('hidden');
            alert.textContent = error.message || 'Login failed, please try again later';
            alert.classList.add('error');
        }
    });

    container.querySelector('#register-link').addEventListener('click', (e) => {
        e.preventDefault();
        renderPage('/register');
    });

    return container;
}