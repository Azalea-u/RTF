import { TimeAgo } from "../../utils.js";
import { getUser } from "./userlist.js";

let postOffset = 0;
let commentOffsets = {};
const limit = 10;

function createPostCard(post) {
    if (!post || !post.id || !post.title || !post.content) return document.createElement("div");

    const container = document.createElement("div");
    container.classList.add("post-card");
    container.innerHTML = `
        <h3 class="post-title">${post.title}</h3>
        <p class="post-username">${getUser(post.user_id) || "Unknown User"}</p>
        <p class="post-content">${post.content}</p>
        <p class="post-categories">Category: ${post.category || "Uncategorized"}</p>
        <p class="post-timestamp">${TimeAgo(post.created_at)}</p>
        <button data-post-id="${post.id}" class="show-comments">Comments</button>
        <form class="comment-form" data-post-id="${post.id}" style="display:none;">
            <input type="text" name="content" placeholder="Comment..." required>
            <button type="submit">Comment</button>
        </form>
        <div class="comments-dropdown" data-post-id="${post.id}"></div>
        <button class="load-more-comments" data-post-id="${post.id}" style="display:none;">Show more comments</button>
    `;

    const commentButton = container.querySelector(".show-comments");
    const commentForm = container.querySelector(".comment-form");
    const loadMoreButton = container.querySelector(".load-more-comments");

    commentButton.addEventListener("click", () => {
        const commentsDropdown = container.querySelector(".comments-dropdown");
        const isVisible = commentsDropdown.innerHTML !== "";

        if (isVisible) {
            commentsDropdown.innerHTML = "";
            loadMoreButton.style.display = "none";
            commentForm.style.display = "none";
            commentButton.style.display = "inline-block";
        } else {
            loadComments(post.id, commentsDropdown, loadMoreButton, commentForm);
            commentButton.style.display = "none";
        }
    });

    commentForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        const content = commentForm.querySelector('input[name="content"]').value;
        const commentsDropdown = container.querySelector(`.comments-dropdown[data-post-id="${post.id}"]`);

        if (content) {
            try {
                await postComment(post.id, content);
                commentForm.reset();
                commentOffsets[post.id] = 0;
                commentsDropdown.innerHTML = "";
                loadComments(post.id, commentsDropdown, loadMoreButton, commentForm);
            } catch (error) {
                showAlert(error.message || "An error occurred while posting the comment", "error");
            }
        }
    });

    return container;
}

async function loadComments(postId, commentsDropdown, loadMoreButton, commentForm) {
    if (commentOffsets[postId] === undefined) commentOffsets[postId] = 0;

    try {
        const response = await fetch(`/api/get-comments?post_id=${postId}&limit=${limit}&offset=${commentOffsets[postId]}`, {
            method: "GET",
            credentials: "include",
        });

        if (!response.ok) throw new Error("Failed to fetch comments");

        const comments = await response.json();
        if (comments.length === 0) {
            loadMoreButton.style.display = "none";
            commentsDropdown.innerHTML = "No comments yet";
            commentForm.style.display = "block";
            return;
        }

        comments.forEach(comment => {
            commentsDropdown.appendChild(createCommentBubble(comment));
        });

        commentOffsets[postId] += comments.length;
        loadMoreButton.style.display = comments.length === limit ? "block" : "none";
        loadMoreButton.onclick = () => loadMoreComments(postId, commentsDropdown, loadMoreButton);
        commentForm.style.display = "block";
    } catch (error) {
        console.error("Fetch comments error:", error);
    }
}

async function loadMoreComments(postId, commentsDropdown, loadMoreButton) {
    try {
        const response = await fetch(`/api/get-comments?post_id=${postId}&limit=${limit}&offset=${commentOffsets[postId]}`, {
            method: "GET",
            credentials: "include",
        });

        if (!response.ok) throw new Error("Failed to fetch more comments");

        const comments = await response.json();
        comments.forEach(comment => {
            commentsDropdown.appendChild(createCommentBubble(comment));
        });

        commentOffsets[postId] += comments.length;
        loadMoreButton.style.display = comments.length < limit ? "none" : "block";
    } catch (error) {
        console.error("Fetch more comments error:", error);
    }
}

function createCommentBubble(comment) {
    const container = document.createElement("div");
    container.classList.add("comment-bubble");
    container.innerHTML = `
        <div class="comment" sender-id="${comment.user_id}">
            <p class="comment-username">${getUser(comment.user_id)}</p>
            <p class="comment-content">${comment.content}</p>
            <span class="comment-timestamp">${TimeAgo(comment.created_at)}</span>
        </div>
    `;
    return container;
}

async function postComment(postId, content) {
    try {
        const response = await fetch(`/api/create-comment`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            credentials: "include",
            body: JSON.stringify({ post_id: postId, content: content }),
        });

        if (!response.ok) {
            const text = await response.text();
            const errorData = text ? JSON.parse(text) : {};
            throw new Error(errorData.message || "Failed to post comment");
        }

        const text = await response.text();
        return text ? JSON.parse(text) : {};
    } catch (error) {
        console.error("Post comment error:", error);
        throw error;
    }
}

function showAlert(message, type) {
    const alertBox = document.createElement("div");
    alertBox.className = `alert alert-${type}`;
    alertBox.textContent = message;
    document.body.appendChild(alertBox);
    setTimeout(() => alertBox.remove(), 3000);
}

export default async function Posts() {
    const container = document.createElement("div");
    container.innerHTML = `
        <form id="post-form" class="post-form">
            <h2>Posts</h2>
            <input type="text" id="title" name="title" placeholder="Title..." required>
            <textarea name="content" placeholder="Content..." required></textarea>
            <input type="text" id="categories" name="category" placeholder="Categories (separated by commas)" required>
            <button type="submit">Post</button>
        </form>
        <div id="post-list" class="post-list"></div>
    `;

    async function loadPosts(reset = false) {
        if (reset) postOffset = 0;

        const postList = container.querySelector("#post-list");
        if (!postList) return;

        postList.innerHTML = "";

        try {
            const response = await fetch(`/api/get-posts?limit=${limit}&offset=${postOffset}`, {
                method: "GET",
                credentials: "include",
            });

            if (!response.ok) throw new Error("Failed to fetch posts");

            const posts = await response.json();
            if (posts.length === 0) return;

            posts.forEach(post => {
                postList.appendChild(createPostCard(post));
            });

            postOffset += posts.length;
            addLoadMoreButton(posts.length);
        } catch (error) {
            console.error("Fetch posts error:", error);
        }
    }

    function addLoadMoreButton(postsLength) {
        const postList = container.querySelector("#post-list");
        const existingButton = postList.querySelector("button.show-more-posts");
        if (existingButton) existingButton.remove();

        if (postsLength === limit) {
            const loadMoreButton = document.createElement("button");
            loadMoreButton.textContent = "Show more posts";
            loadMoreButton.classList.add("show-more-posts");
            loadMoreButton.onclick = loadMorePosts;
            postList.appendChild(loadMoreButton);
        }
    }

    async function loadMorePosts() {
        try {
            const response = await fetch(`/api/get-posts?limit=${limit}&offset=${postOffset}`, {
                method: "GET",
                credentials: "include",
            });

            if (!response.ok) throw new Error("Failed to fetch more posts");

            const posts = await response.json();
            if (posts.length === 0) return;

            const postList = container.querySelector("#post-list");
            posts.forEach(post => {
                postList.appendChild(createPostCard(post));
            });

            postOffset += posts.length;
            addLoadMoreButton(posts.length);
        } catch (error) {
            console.error("Fetch more posts error:", error);
        }
    }

    const postForm = container.querySelector("#post-form");
    postForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        const formData = new FormData(postForm);
        const data = Object.fromEntries(formData.entries());

        if (!data.title || !data.content || !data.category) return;

        const response = await fetch("/api/create-post", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(data),
        });

        if (response.ok) {
            postForm.reset();
            await loadPosts(true);
        } else {
            const errorData = await response.json();
            console.error(errorData.message || "Failed to create post");
        }
    });

    setTimeout(() => loadPosts(true), 0);
    return container;
}
