// forum.js
export function loadHomePage() {
    const app = document.getElementById('app');
    app.innerHTML = '<h2>Welcome to the Forum App!</h2>';
}

export async function loadForumPage() {
    const app = document.getElementById('app');
    app.innerHTML = '<h2>Forum</h2><div id="posts"></div>';

    // Fetch posts from the backend
    const response = await fetch('/api/posts');
    if (response.ok) {
        const posts = await response.json();
        const postsContainer = document.getElementById('posts');
        posts.forEach(post => {
            const postElement = document.createElement('div');
            postElement.innerHTML = `
                <h3>${post.title}</h3>
                <p>${post.content}</p>
                <small>Posted by ${post.user_id} on ${new Date(post.created_at).toLocaleString()}</small>
            `;
            postsContainer.appendChild(postElement);
        });
    } else {
        app.innerHTML += '<p>Failed to load posts.</p>';
    }
}