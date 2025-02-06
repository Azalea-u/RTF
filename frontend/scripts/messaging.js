// messaging.js
let socket;
let currentUser = null;
let currentConversationUserID = null;

export async function loadMessagingSidebar(user) {
  currentUser = user;
  const sidebar = document.getElementById("messaging-sidebar");
  sidebar.innerHTML = `
    <h3>Messages</h3>
    <div id="user-list"></div>
    <div id="chat-window" style="margin-top:20px;">
      <div id="chat-messages"></div>
      <form id="message-form">
        <input type="text" id="message-input" placeholder="Type a message" required>
        <button type="submit">Send</button>
      </form>
    </div>
  `;
  await loadUserList();
  initWebSocket();
  attachFormListener();
}

async function loadUserList() {
  const res = await fetch("/api/users");
  if (!res.ok) return;
  const users = await res.json();
  const userListDiv = document.getElementById("user-list");
  userListDiv.innerHTML = "";
  users.forEach(u => {
    if (u.id === currentUser.id) return;
    const div = document.createElement("div");
    div.textContent = u.username;
    div.dataset.userid = u.id;
    div.style.cursor = "pointer";
    div.addEventListener("click", () => {
      currentConversationUserID = u.id;
      loadConversation(u.id);
    });
    userListDiv.appendChild(div);
  });
}

async function loadConversation(otherUserID) {
  const res = await fetch("/api/messages?user_id=" + otherUserID);
  if (!res.ok) return;
  const messages = await res.json();
  const chatMessagesDiv = document.getElementById("chat-messages");
  chatMessagesDiv.innerHTML = "";
  messages.forEach(m => {
    const p = document.createElement("p");
    p.textContent = (m.sender_id === currentUser.id ? "You: " : "Them: ") + m.content;
    chatMessagesDiv.appendChild(p);
  });
}

function initWebSocket() {
  const protocol = location.protocol === "https:" ? "wss" : "ws";
  socket = new WebSocket(`${protocol}://${location.host}/api/chat`);

  socket.onopen = () => {
    console.log("WebSocket connected");
  };

  socket.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    // Only update if the message belongs to the current conversation
    if (currentConversationUserID &&
        (msg.sender_id === currentConversationUserID || msg.receiver_id === currentConversationUserID)) {
      appendMessage(msg);
    }
  };

  socket.onclose = () => {
    console.log("WebSocket disconnected, retrying in 3 seconds...");
    setTimeout(initWebSocket, 3000);
  };
}

function appendMessage(msg) {
  const chatMessagesDiv = document.getElementById("chat-messages");
  const p = document.createElement("p");
  p.textContent = (msg.sender_id === currentUser.id ? "You: " : "Them: ") + msg.content;
  chatMessagesDiv.appendChild(p);
}

function attachFormListener() {
  const form = document.getElementById("message-form");
  form.onsubmit = async (e) => {
    e.preventDefault();
    const input = document.getElementById("message-input");
    const content = input.value.trim();
    if (!content || !currentConversationUserID) return;
    const message = {
      receiver_id: currentConversationUserID,
      content,
      sender_id: currentUser.id
    };
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify(message));
    }
    input.value = "";
  };
}
