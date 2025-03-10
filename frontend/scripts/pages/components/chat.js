import { renderPage } from '../../router.js';
import { showAlert, TimeAgo } from '../../utils.js';

export let inChat = false;

function messageBubble(message) {
    const container = document.createElement('div');
    container.classList.add('message-bubble');
    const currentUserId = Number(localStorage.getItem('user_id'));
    container.innerHTML = `
        <div class="message ${message.sender_id === currentUserId ? 'sent' : 'received'}" sender-id="${message.sender_id}">
            <p>${message.content}</p>
            <span class="timestamp">${TimeAgo(message.created_at)}</span>
        </div>
    `;
    return container;
}

export default async function Chat(userId, username) {
    const container = document.createElement('div');
    container.id = 'chat';

    const response = await fetch(`/api/messages/${userId}`, {
        method: 'GET',
        credentials: 'include',
    });

    if (!response.ok) {
        console.error('Failed to fetch messages');
        showAlert('Failed to fetch messages, please try again', 'error');
        return;
    }

    const messages = await response.json();
    inChat = true;

    container.innerHTML = `
        <div class="chat-header">
            <button id="exit-chat">Exit</button>
            <h2>Chatting with <span class="chat-username">${username}</span></h2>
        </div>
        <div class="messages"></div>
        <form id="message-form">
            <input type="text" id="message-input" placeholder="Type your message here..." required>
            <button type="submit" id="send-button">Send</button>
        </form>
    `;

    const messagesContainer = container.querySelector('.messages');

    function scrollToBottom() {
        setTimeout(() => {
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
        }, 100); // Small delay to ensure DOM updates before scrolling
    }

    if (Array.isArray(messages) && messages.length > 0) {
        messages.forEach(message => {
            messagesContainer.appendChild(messageBubble(message));
        });
    } else {
        messagesContainer.innerHTML = 'No messages yet for this conversation';
    }

    scrollToBottom(); // Scroll down after loading messages

    container.querySelector('#message-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        const messageInput = container.querySelector('#message-input');
        const message = messageInput.value.trim();

        if (!message) {
            console.error('Message is empty');
            showAlert('Please enter a message', 'error');
            return;
        }

        const sendResponse = await fetch(`/api/messages/${userId}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ content: message }),
        });

        if (!sendResponse.ok) {
            console.error('Failed to send message');
            showAlert('Failed to send message, please try again', 'error');
            return;
        }

        messageInput.value = '';

        updateChat(userId);
    });

    container.querySelector('#exit-chat').addEventListener('click', () => {
        inChat = false;
        renderPage('/');
    });

    return container;
}

export function updateChat(userId) {
    fetch(`/api/messages/${userId}`)
        .then(response => response.json())
        .then(messages => {
            const messagesContainer = document.querySelector('.messages');
            if (!messagesContainer) return;

            messagesContainer.innerHTML = '';

            if (messages.length === 0) {
                messagesContainer.innerHTML = 'No messages yet for this conversation';
                return;
            }

            messages.forEach(message => {
                messagesContainer.appendChild(messageBubble(message));
            });

            messagesContainer.scrollTop = messagesContainer.scrollHeight; // Scroll to bottom after update
        })
        .catch(error => {
            console.error('Error fetching updated messages:', error);
            showAlert('Failed to update chat messages', 'error');
        });
}
