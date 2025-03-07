import { renderPage } from "../router.js";
import { showAlert } from "../utils.js";

export default function register() {
    const container = document.createElement('div');
    container.innerHTML = `
        <h1>Login</h1>
        <form id="login-form">
            <input type="text" id="username-or-email" name="email_or_username" placeholder="Username or Email" required>
            <input type="password" id="password" name="password" placeholder="Password" required>
            <button type="submit">Login</button>
            <p>Don't have an account? <a href="#" id="register-link">Register</a></p>
        </form>
    `;

    container.querySelector('#login-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData(e.target);
        const data = Object.fromEntries(formData.entries());

        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify(data),
        });

        if (response.ok) {
            const user = await response.json();
            localStorage.setItem('userId', user.id);
            localStorage.setItem('username', user.username);
            showAlert('Welcome, ' + user.username + '!', 'success', 'success');
            e.target.reset();
            initWebSocket();
            setTimeout(() => {
                renderPage('/');
            }, 1000);
        } else {
            const error = await response.json();
            showAlert(error.message || 'Login failed, please try again', 'error');
        }
    });

    container.querySelector('#register-link').addEventListener('click', (e) => {
        e.preventDefault();
        renderPage('/register');
    });

    return container;
}