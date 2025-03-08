import { TimeAgo } from "../../utils.js";
import { getUser  } from "./userlist.js";

let postOffset = 0;
let commentOffsets = {};
const limit = 10;

function createPostCard(post) {
    const container = document.createElement('div');
    container.classList.add('post-card');
    container.innerHTML = `
        <h3 class="post-title">${post.title}</h3>
        <p class="post-username">${getUser (post.user_id)}</p>
        <p class="post-content">${post.content}</p>
        <p class="post-categories">Categories: ${post.categories.join(', ')}</p>
        <p class="post-timestamp">${TimeAgo(post.created_at)}</p>
        <button data-post-id="${post.id}" class="comment-button">Comments</button>
        <div class="comments-dropdown" data-post-id="${post.id}"></div>
        <button class="load-more-comments" data-post-id="${post.id}" style="display:none;">Load More Comments</button>
    `;

    const commentButton = container.querySelector('.comment-button');
    const loadMoreButton = container.querySelector('.load-more-comments');

    commentButton.addEventListener('click', () => {
        const commentsDropdown = container.querySelector('.comments-dropdown');
        if (commentsDropdown.innerHTML === '') {
            loadComments(post.id, commentsDropdown, loadMoreButton);
        } else {
            commentsDropdown.innerHTML = '';
            loadMoreButton.style.display = 'none';
        }
    });

    return container;
}

async function loadComments(postId, commentsDropdown, loadMoreButton) {
    if (!commentOffsets[postId]) {
        commentOffsets[postId] = 0;
    }

    try {
        const response = await fetch(`/api/get-comments?post_id=${postId}&limit=${limit}&offset=${commentOffsets[postId]}`, {
            method: 'GET',
            credentials: 'include',
        });
        if (!response.ok) throw new Error('Failed to fetch comments');

        const comments = await response.json();
        comments.forEach(comment => {
            commentsDropdown.appendChild(createCommentBubble(comment));
        });
        commentOffsets[postId] += comments.length;

        if (comments.length === limit) {
            loadMoreButton.style.display = 'block';
        } else {
            loadMoreButton.style.display = 'none';
        }

        loadMoreButton.onclick = () => {
            loadMoreComments(postId, commentsDropdown, loadMoreButton);
        };
    } catch (error) {
        console.error('Fetch comments error:', error);
    }
}

async function loadMoreComments(postId, commentsDropdown, loadMoreButton) {
    try {
        const response = await fetch(`/api/get-comments?post_id=${postId}&limit=${limit}&offset=${commentOffsets[postId]}`, {
            method: 'GET',
            credentials: 'include',
        });
        if (!response.ok) throw new Error('Failed to fetch more comments');

        const comments = await response.json();
        comments.forEach(comment => {
            commentsDropdown.appendChild(createCommentBubble(comment));
        });
        commentOffsets[postId] += comments.length;

        if (comments.length < limit) {
            loadMoreButton.style.display = 'none';
        }
    } catch (error) {
        console.error('Fetch more comments error:', error);
    }
}

function createCommentBubble(comment) {
    const container = document.createElement('div');
    container.classList.add('comment-bubble');
    container.innerHTML = `
        <div class="comment" sender-id="${comment.sender_id}">
            <p class="comment-username">${getUser (comment.user_id)}</p>
            <p class="comment-content">${comment.content}</p>
            <span class="comment-timestamp">${TimeAgo(comment.created_at)}</span>
        </div>
    `;
    return container;
}

export default async function Posts() {
    const container = document.createElement('div');
    container.innerHTML = `
        <h2>Posts</h2>
        <form id="post-form" class="post-form">
            <input type="text" id="title" name="title" placeholder="Title..." required>
            <textarea id="content" name="content" placeholder="What's on your mind?" required></textarea>
            <input type="text" id="categories" name="categories" placeholder="Categories (separated by commas)" required>
            <button type="submit">Post</button>
        </form>
        <div id="post-list" class="post-list"></div>
    `;

    const postForm = container.querySelector('#post-form');
    postForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData(postForm);
        formData.append('user_id', localStorage.getItem('userId'));

        try {
            const response = await fetch('/api/create-post', {
                method: 'POST',
                body: formData,
                credentials: 'include',
            });
            if (!response.ok) throw new Error('Failed to create post');

            postForm.reset();
            postOffset = 0;
            loadPosts();
        } catch (error) {
            console.error('Error creating post:', error);
        }
    });

    await loadPosts();

    return container;
}

async function loadPosts() {
    const postList = document.getElementById('post-list');
    if (!postList) {
        console.error('Post list element not found');
        return;
    }
    postList.innerHTML = '';

    try {
        const response = await fetch(`/api/get-posts?limit=${limit}&offset=${postOffset}`, {
            method: 'GET',
            credentials: 'include',
        });
        if (!response.ok) throw new Error('Failed to fetch posts');

        const posts = await response.json();
        posts.forEach(post => {
            postList.appendChild(createPostCard(post));
        });
        postOffset += posts.length;

        if (posts.length === limit) {
            const loadMoreButton = document.createElement('button');
            loadMoreButton.textContent = 'Load More Posts';
            loadMoreButton.onclick = loadMorePosts;
            postList.appendChild(loadMoreButton);
        }
    } catch (error) {
        console.error('Fetch posts error:', error);
    }
}

async function loadMorePosts() {
    try {
        const response = await fetch(`/api/get-posts?limit=${limit}&offset=${postOffset}`, {
            method: 'GET',
            credentials: 'include',
        });
        if (!response.ok) throw new Error('Failed to fetch more posts');

        const posts = await response.json();
        if (posts.length === 0) return;

        const postList = document.getElementById('post-list');
        posts.forEach(post => {
            postList.appendChild(createPostCard(post));
        });
        postOffset += posts.length;
    } catch (error) {
        console.error('Fetch more posts error:', error);
    }
}