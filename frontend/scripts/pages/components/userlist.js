import { showAlert } from "../../utils.js";

let users = [];

async function fetchUsers() {
    try {
        const response = await fetch('/api/get-users', {
            method: 'GET',
            credentials: 'include',
        });
        if (!response.ok) {
            throw new Error('Network response was not ok');
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
            ${Array.isArray(users) && users.length > 0 ? users.map(user => `<p data-user-id="${user.id}" class="user">${user.username}</p>`).join('') : 'No users found'}
        </div>
    `;
    
    return container;
}