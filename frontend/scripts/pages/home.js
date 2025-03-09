import { renderPage } from "../router.js";
import { showAlert } from "../utils.js";
import UserList from "./components/userlist.js";
import Chat from "./components/chat.js";
import { initWebSocket } from "../websocket.js";
import Posts from "./components/posts.js";

export default async function home() {
    initWebSocket();
    const container = document.createElement("div");
    container.innerHTML = `
        <nav class="navbar">
            <h1>Welcome ${localStorage.getItem("username")}</h1>
            <button id="logout-button">Logout</button>
        </nav>
        <div id="content"></div>
    `;

    const userList = await UserList();
    container.appendChild(userList);

    const postsComponent = await Posts();
    container.querySelector("#content").appendChild(postsComponent);

    container.querySelector("#logout-button").addEventListener("click", async () => {
        const response = await fetch("/api/logout", {
            method: "POST",
            credentials: "include",
        });

        if (response.ok) {
            showAlert("Logout successful", "success");
            setTimeout(() => renderPage("/login"), 1000);
        } else {
            showAlert("Logout failed", "error");
        }
    });

    return container;
}

export async function onUserClick(userId, username) {
    const chatComponent = await Chat(userId, username);
    const contentDiv = document.querySelector("#content");
    contentDiv.innerHTML = "";
    contentDiv.appendChild(chatComponent);
}
