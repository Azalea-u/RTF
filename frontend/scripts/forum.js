// forum.js
export function loadHomePage() {
    const app = document.getElementById('app');
    app.innerHTML = '<h2>Welcome to the Forum App!</h2>';
}

export async function loadForumPage() {
    const app = document.getElementById('app');
    app.innerHTML = '<h2>Forum</h2><div id="posts"></div>';
    // Post creation form
    app.innerHTML += `
        <form id="post-form">
            <input type="text" id="title" placeholder="Title" required>
            <textarea id="content" placeholder="Content" required></textarea>
            <div id="categories">
                <label><input type="radio" name="category" value="general"> General</label>
                <label><input type="radio" name="category" value="tech"> Tech</label>
                <label><input type="radio" name="category" value="music"> Music</label>
                <label><input type="radio" name="category" value="health"> Health</label>
            </div>

            <butto type="submit">Post</button>
        </form>
    `;

    const postForm = document.getElementById('post-form');
    postForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const title = document.getElementById('title').value;
        const content = document.getElementById('content').value;
        const category = document.querySelector('input[name="category"]:checked').value;
        const response = await fetch('/api/create-post', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ title, content, category }),
        });

        if (response.ok) {
            alert('Post created successfully.');
            window.location.reload();
        } else {
            alert('Failed to create post.');
        }
    });

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