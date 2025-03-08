import { showAlert } from "./utils.js";
import { inChat, updateChat } from "./pages/components/chat.js";
import { reloadUserList } from "./pages/components/userlist.js";
import { renderPage } from "./router.js";

let socket;
const currentUserId = localStorage.getItem('userId');

export function initWebSocket() {
    if (socket && socket.readyState === WebSocket.OPEN) return;

    const protocol = location.protocol === 'https:' ? 'wss' : 'ws';
    socket = new WebSocket(`${protocol}://${location.host}/api/ws`);

    socket.onopen = () => {
        console.log('WebSocket connection established');
    };

    socket.onmessage = handleMessage;

    socket.onerror = (error) => {
        console.error('WebSocket error:', error);
        showAlert('WebSocket error occurred. Please try again.', 'error');
    };

    socket.onclose = () => {
        console.log('WebSocket connection closed. Attempting to reconnect...');
        renderPage('/login');
    };
}

function handleMessage(event) {
    const message = JSON.parse(event.data);

    switch (message.type) {
        case 'chat':
            handleChatMessage(message.content);
            break;
        case 'user_connected':
        case 'user_disconnected':
            reloadUserList();
            break;
        default:
            console.warn('Unknown message type:', message.type);
    }
}

function handleChatMessage(message) {
    if (message.sender_id !== currentUserId) {
        showAlert(`New message from ${document.getElementById('chat-username').textContent}`, 'success');
        if (inChat) {
            updateChat(message);
        }
        reloadUserList([message.sender_id, message.receiver_id]);
    }
}