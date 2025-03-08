import { renderPage } from "../../router.js";
import { showAlert } from "../../utils.js";
import { onUserClick } from "../home.js";

let users = [];

async function fetchUsers() {
    try {
        const response = await fetch('/api/get-users', {
            method: 'GET',
            credentials: 'include',
        });
        if (!response.ok) {
            throw new Error('Failed to fetch users');
        }
        users = await response.json();
    } catch (error) {
        showAlert('Something went wrong, please try again', 'error');
        console.error('Fetch users error:', error);
        users = [];
    }
}

export default async function UserList() {
    await fetchUsers();

    const container = document.createElement('aside');
    container.id = 'user-list';
    container.setAttribute('role', 'complementary');
    container.innerHTML = `
        <h2>Users</h2>
        <div class="user-list">
            ${Array.isArray(users) && users.length > 0 ? 
                users.map(user => `
                    <p data-user-id="${user.id}" class="user">
                        ${user.username} 
                        <span class="${user.online ? 'online' : 'offline'}"></span>
                    </p>`).join('') : 
                'No users found'}
        </div>
    `;

    const userElements = container.querySelectorAll('.user');
    userElements.forEach(userElement => {
        userElement.addEventListener('click', () => {
            const userId = userElement.getAttribute('data-user-id');
            onUserClick(userId, userElement.textContent);
        });
    });

    return container;
}
// Reload the user list only for the specified user IDs
export function reloadUserList(userIds = []) {
    const userElements = document.querySelectorAll('.user');
    const currentUser = localStorage.getItem('userId');
    if (userIds.length === 0) {
        updateUserList();
        return;
    }
    userIds.forEach(userId => {
        if (userId !== currentUser) {
            updateUserList();
        }
    })
}

function updateUserList() {
    fetch ('/api/get-users')
        .then(response => response.json())
        .then(data => {
            const userList = document.getElementById('user-list');
            userList.innerHTML = `
                <h2>Users</h2>
                <div class="user-list">
                    ${Array.isArray(data) && data.length > 0 ? 
                        data.map(user => `
                            <p data-user-id="${user.id}" class="user">
                                ${user.username} 
                                <span class="${user.online ? 'online' : 'offline'}"></span>
                            </p>`).join('') : 
                        'No users found'}
                </div>
            `;
            const userElements = userList.querySelectorAll('.user');
            userElements.forEach(userElement => {
                userElement.addEventListener('click', () => {
                    const userId = userElement.getAttribute('data-user-id');
                    onUserClick(userId, userElement.textContent);
                });
            });
        })
        .catch(error => {
            console.error('Fetch users error:', error);
            renderPage('/login');
        });
}