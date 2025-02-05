export async function loadPosts() {
    const response = await fetch("/api/posts");
    if (!response.ok) return;

    const posts = await response.json();
    const postSection = document.getElementById("posts");
    postSection.innerHTML = posts.map(post => `
        <article>
            <h3>${post.title}</h3>
            <p>${post.content}</p>
            <small>Category: ${post.category}</small>
        </article>
    `).join("");
}

export async function createPost(title, content) {
    await fetch("/api/create-post", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title, content }),
    });
}
