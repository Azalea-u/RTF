import { renderPage } from "../router.js";
import { showAlert } from "../utils.js";

export default function register() {
    const container = document.createElement('div');
    container.innerHTML = `
        <h1>Registration</h1>
        <form id="registration-form">
            <input type="text" name="username" placeholder="Username" required>
            <input type="email" name="email" placeholder="Email" required>
            <input type="password" name="password" placeholder="Password" required>
            <input type="text" name="first_name" placeholder="First Name" required>
            <input type="text" name="last_name" placeholder="Last Name" required>            <input type="number" name="age" placeholder="Age" required>
            <select name="gender" required>
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
        if (!data.username || !data.email || !data.password || !data.first_name || !data.last_name || !data.age || !data.gender) {
            console.error('All fields must be filled out');
            return;
        }

        // change age to int
        data.age = parseInt(data.age);

        const response = await fetch('/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });

        if (response.ok) {
            showAlert('Registration successful', 'success');
            e.target.reset();
            setTimeout(() => {
                renderPage('/login');
            }, 1000);
        } else {
            const error = await response.json();
            showAlert(error.message || 'Registration failed, please try again', 'error');
        }
    });

    container.querySelector('#login-link').addEventListener('click', (e) => {
        e.preventDefault();
        renderPage('/login');
    });

    return container;
}