import { showAlert, TimeAgo } from '../../utils.js';

function messageBubble(message) {
    const container = document.createElement('div');
    container.classList.add('message-bubble');
    container.innerHTML = `
        <div class="message" sender-id="${message.sender_id}">
            <p>${message.content}</p>
            <span class="timestamp">${TimeAgo(message.created_at)}</span>
        </div>
    `;
    return container;
}

export default async function Chat(userId) {
    const container = document.createElement('div');
    container.id = 'chat';
    messages = await fetch(`/api/messages/${userId}`, {
        method: 'GET',
        credentials: 'include',
    });
    if (!messages.ok) {
        console.error('Failed to fetch messages');
        showAlert('Failed to fetch messages, please try again', 'error');
        return;
    }

    container.innerHTML = `
        <h2>Chating with ${userId}</h2>
        <div class="messages">
            ${Array.isArray(messages) && messages.length > 0 ? 
                messages.map(message => messageBubble(message)).join('') : 
                'No messages yet for this conversation'}
        </div>
        <form id="message-form">
            <input type="text" id="message-input" placeholder="Type your message here..." required>
            <button type="submit">Send</button>
        </form>
    `;

    container.querySelector('#message-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        const messageInput = container.querySelector('#message-input');
        const message = messageInput.value;
        if (!message) {
            console.error('Message is empty');
            return;
        }
        const response = await fetch(`/api/messages/${userId}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify({ content: message }),
        });
        if (!response.ok) {
            console.error('Failed to send message');
            showAlert('Failed to send message, please try again', 'error');
            return;
        }
        messageInput.value = '';
    });

    return container;
}