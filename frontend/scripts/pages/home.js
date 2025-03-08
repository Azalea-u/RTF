import { renderPage } from "../router.js";
import { showAlert } from "../utils.js";
import UserList from "./components/userlist.js";
import Chat from "./components/chat.js";

export default async function home() {
    const container = document.createElement('div');
    container.innerHTML = `
        <nav class="navbar">
            <h1>Welcome ${localStorage.getItem('username')}</h1>
            <button id="logout-button">Logout</button>
        </nav>
        <div id="content">
            <h2>Select a user to chat with</h2>
        </div>
    `;

    // Fetch and render the user list
    const userList = await UserList(onUserClick);
    container.appendChild(userList);

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

    async function onUserClick(userId, username) {
        const chatComponent = await Chat(userId, username);
        const contentDiv = container.querySelector('#content');
        contentDiv.innerHTML = ''; 
        contentDiv.appendChild(chatComponent);
    }

    return container;
}