import { showAlert } from "./utils.js";
import { inChat, updateChat } from "./pages/components/chat.js";
import { getUser, updateUserList } from "./pages/components/userlist.js";
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
        case 'message':
            handleChatMessage(message.content);
            break;
        case 'user_connected':
        case 'user_disconnected':
            updateUserList();
            break;
        default:
            console.warn('Unknown message type:', message.type);
    }
}

function handleChatMessage(message) {
    const ids = message.split(',');
    let sender_id = ids[0];
    let receiver_id = ids[1];

    if (receiver_id === currentUserId) {
        showAlert(`New message from ${getUser(sender_id)}`, 'success');
        if (inChat) {
            updateChat(receiver_id);
            updateChat(sender_id);
        }
    }
    if (sender_id === currentUserId || receiver_id === currentUserId) {
        updateUserList();
    }
}