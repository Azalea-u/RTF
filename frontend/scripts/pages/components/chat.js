import { renderPage } from '../../router.js';
import { showAlert, TimeAgo } from '../../utils.js';

export let inChat = false;

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

export default async function Chat(userId, username) {
    const container = document.createElement('div');
    container.id = 'chat';

    // Pagination state
    let offset = 0;
    let hasMore = true;
    let isLoading = false;

    // Initial messages load
    const response = await fetch(`/api/messages/${userId}?limit=10&offset=0`, {
        method: 'GET',
        credentials: 'include',
    });

    if (!response.ok) {
        console.error('Failed to fetch messages');
        showAlert('Failed to fetch messages, please try again', 'error');
        return;
    }

    const messages = await response.json();
    hasMore = messages.length === 10;
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

    messagesContainer.addEventListener('scroll', async () => {
        if (isLoading || !hasMore) return;
        
        // Load more when 80px from top
        if (messagesContainer.scrollTop < 80) {
            isLoading = true;
            const newOffset = offset + 10;
            
            try {
                const response = await fetch(`/api/messages/${userId}?limit=10&offset=${newOffset}`, {
                    credentials: 'include',
                });
                
                if (!response.ok) throw new Error('Failed to load messages');
                
                const newMessages = await response.json();
                hasMore = newMessages.length === 10;
                
                if (newMessages.length > 0) {
                    const oldHeight = messagesContainer.scrollHeight;
                    const fragment = document.createDocumentFragment();
                    
                    newMessages.reverse().forEach(message => {
                        fragment.appendChild(messageBubble(message));
                    });
                    
                    messagesContainer.prepend(fragment);
                    offset = newOffset;
                    
                    // Maintain scroll position
                    messagesContainer.scrollTop = messagesContainer.scrollHeight - oldHeight;
                }
            } catch (error) {
                console.error('Error loading more messages:', error);
                showAlert('Failed to load older messages', 'error');
            } finally {
                isLoading = false;
            }
        }
    });

    function scrollToBottom() {
        setTimeout(() => {
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
        }, 100);
    }

    // Initial messages rendering
    if (messages.length > 0) {
        messages.reverse().forEach(message => {
            messagesContainer.appendChild(messageBubble(message));
        });
    } else {
        messagesContainer.innerHTML = 'No messages yet for this conversation';
    }

    scrollToBottom();

    // Message submission handler
    container.querySelector('#message-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        const messageInput = container.querySelector('#message-input');
        const message = messageInput.value.trim();

        if (!message) {
            showAlert('Please enter a message', 'error');
            return;
        }

        try {
            const sendResponse = await fetch(`/api/messages/${userId}`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ content: message }),
            });

            if (!sendResponse.ok) throw new Error('Failed to send message');
            
            // Optimistically add new message
            const newMessage = {
                content: message,
                sender_id: Number(localStorage.getItem('user_id')),
                created_at: new Date().toISOString(), // Use current time for the new message
            };
            messagesContainer.appendChild(messageBubble(newMessage));
            scrollToBottom();
            messageInput.value = '';
            updateChat(userId);
            hasMore = true;
            offset = 0;
        } catch (error) {
            console.error('Message send failed:', error);
            showAlert('Failed to send message, please try again', 'error');
        }
    });

    container.querySelector('#exit-chat').addEventListener('click', () => {
        inChat = false;
        renderPage('/');
    });

    return container;
}

export function updateChat(userId) {
    fetch(`/api/messages/${userId}?limit=10&offset=0`)
        .then(response => response.json())
        .then(messages => {
            const messagesContainer = document.querySelector('.messages');
            if (!messagesContainer) return;

            messagesContainer.innerHTML = '';

            if (messages.length === 0) {
                messagesContainer.innerHTML = 'No messages yet for this conversation';
                return;
            }

            messages.reverse().forEach(message => {
                messagesContainer.appendChild(messageBubble(message));
            });

            messagesContainer.scrollTop = messagesContainer.scrollHeight;
        })
        .catch(error => {
            console.error('Error fetching updated messages:', error);
            showAlert('Failed to update chat messages', 'error');
        });
}