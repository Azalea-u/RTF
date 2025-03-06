import { renderPage } from "../router.js";

export default function register() {
    const container = document.createElement('div');
    container.innerHTML = `
        <h1>Registration</h1>
        <form id="registration-form">
            <input type="text" id="username" placeholder="Username" required>
            <input type="email" id="email" placeholder="Email" required>
            <input type="password" id="password" placeholder="Password" required>
            <input type="text" id="first-name" placeholder="First Name" required>
            <input type="text" id="last-name" placeholder="Last Name" required>
            <input type="number" id="age" placeholder="Age" required>
            <select id="gender" required>
                <option value="">Select Gender</option>
                <option value="male">Male</option>
                <option value="female">Female</option>
                <option value="other">Other</option>
            </select>
            <button type="submit">Register</button>
            <p>Already have an account? <a href="#" id="login-link">Login</a></p>
        </form>
    `;

    container.querySelector('#registration-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData(e.target);
        const data = Object.fromEntries(formData.entries());
        const alert = document.getElementById('alert');

        const response = await fetch('/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });

        // Hide loading indicator
        alert.classList.remove('loading');

        if (response.ok) {
            alert.classList.remove('hidden');
            alert.textContent = 'Registration successful!';
            alert.classList.add('success');
            e.target.reset();
            setTimeout(() => {
                renderPage('/login');
            }, 2000);
        } else {
            const error = await response.json();
            alert.classList.remove('hidden');
            alert.textContent = error.message || 'Registration failed. Please try again.';
            alert.classList.add('error');
        }
    });

    container.querySelector('#login-link').addEventListener('click', (e) => {
        e.preventDefault();
        renderPage('/login');
    });

    return container;
}