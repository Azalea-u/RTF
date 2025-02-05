import { checkAuth, loginUser, registerUser, logoutUser } from "./auth.js";
import { loadPosts, createPost } from "./forum.js";

document.addEventListener("DOMContentLoaded", async () => {
    const user = await checkAuth();

    if (!user) {
        loadLoginPage();
    } else {
        loadForumPage(user);
    }
});

function loadLoginPage() {
    document.getElementById("app").innerHTML = `
        <h2>Login</h2>
        <form id="login-form">
            <input type="text" id="username" placeholder="Username" required>
            <input type="password" id="password" placeholder="Password" required>
            <button type="submit">Login</button>
        </form>
        <p>No account? <a href="#" id="register-link">Register</a></p>
    `;

    document.getElementById("login-form").addEventListener("submit", async (e) => {
        e.preventDefault();
        const username = document.getElementById("username").value;
        const password = document.getElementById("password").value;
        const success = await loginUser(username, password);

        if (success) {
            const user = await checkAuth();
            loadForumPage(user);
        }
    });

    document.getElementById("register-link").addEventListener("click", (e) => {
        e.preventDefault();
        loadRegisterPage();
    });
}

function loadRegisterPage() {
    document.getElementById("app").innerHTML = `
        <h2>Register</h2>
        <form id="register-form">
            <input type="text" id="username" placeholder="Username" required>
            <input type="email" id="email" placeholder="Email" required>
            <input type="password" id="password" placeholder="Password" required>
            <input type="text" id="first-name" placeholder="First Name" required>
            <input type="text" id="last-name" placeholder="Last Name" required>
            <select id="gender" required>
                <option value="">Select Gender</option>
                <option value="male">Male</option>
                <option value="female">Female</option>
            </select>
            <button type="submit">Register</button>
        </form>
        <p>Already have an account? <a href="#" id="login-link">Login</a></p>
    `;

    document.getElementById("register-form").addEventListener("submit", async (e) => {
        e.preventDefault();
        const username = document.getElementById("username").value;
        const email = document.getElementById("email").value;
        const password = document.getElementById("password").value;
        const firstName = document.getElementById("first-name").value;
        const lastName = document.getElementById("last-name").value;
        const gender = document.getElementById("gender").value;

        const success = await registerUser(username, email, password, firstName, lastName, gender);
        if (success) {
            loadLoginPage(); // Redirect to login after successful registration
        }
    });

    document.getElementById("login-link").addEventListener("click", (e) => {
        e.preventDefault();
        loadLoginPage();
    });
}

async function loadForumPage(user) {
    document.getElementById("app").innerHTML = `
        <h2>Welcome, ${user.username}</h2>
        <button id="logout">Logout</button>
        <section id="posts"></section>
        <form id="post-form">
            <input type="text" id="post-title" placeholder="Title" required>
            <textarea id="post-content" placeholder="Write something..." required></textarea>
            <div id="category">
                <label><input type="radio" name="category" value="general" required>General</label>
                <label><input type="radio" name="category" value="tech">Tech</label>
                <label><input type="radio" name="category" value="lifestyle">Lifestyle</label>
            </div>
            <button type="submit">Post</button>
        </form>
    `;

    document.getElementById("logout").addEventListener("click", async () => {
        await logoutUser();
        loadLoginPage();
    });

    document.getElementById("post-form").addEventListener("submit", async (e) => {
        e.preventDefault();
        const title = document.getElementById("post-title").value;
        const content = document.getElementById("post-content").value;
        await createPost(title, content);
        loadForumPage(user);
    });

    await loadPosts();
}
