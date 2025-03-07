import { showAlert } from "./utils";
import { inChat, updateChat } from "./pages/components/chat";
import { updateUserList, reloadUserList } from "./pages/components/userlist";

let socket;
let currentUserId = localStorage.getItem('userId');

export function initWebSocket() {
    if (socket && socket.readyState === WebSocket.OPEN) return;

    const protocol = location.protocol === 'https:' ? 'wss' : 'ws';
    socket = new WebSocket(`${protocol}://${location.host}/api/chat`);

    socket.onopen = () => {
        console.log('WebSocket connection established');
    };

    socket.onmessage = (event) => {
        const message = JSON.parse(event.data);

        switch (message.type) {
            case 'chat':
                if (message.sender_id !== currentUserId) {
                    showAlert('New message from ' + document.getElementById('chat-username').textContent, 'success');
                    if (inChat) {
                        updateChat(message);
                    }
                    reloadUserList(message.receiver_id, message.sender_id); // Reload the user list when a new message is received from only receiver and sender to reorders the list of users
                }
                break;
            case 'user_connect':
                updateUserList(message.sender_id, 'online');
                break;
            case 'user_disconnect':
                updateUserList(message.sender_id, 'offline');
                break;
            default:
                console.log('Unknown message type:', message.type);
                break;
        }
    };

    socket.onerror = (error) => {
        console.error('WebSocket error:', error);
        showAlert('WebSocket error occurred. Please try again.', 'error');
    };

    socket.onclose = () => {
        console.log('WebSocket connection closed');
        setTimeout(initWebSocket, 3000);
    };
}