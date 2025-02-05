export async function checkAuth() {
    const response = await fetch("/api/user");
    if (!response.ok) return null;
    return await response.json();
}

export async function loginUser(username, password) {
    const response = await fetch("/api/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
    });

    return response.ok;
}

export async function registerUser(username, email, password, firstName, lastName, gender) {
    const response = await fetch("/api/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ 
            username, 
            email, 
            password, 
            first_name: firstName, 
            last_name: lastName, 
            gender 
        }),
    });

    return response.ok;
}

export async function logoutUser() {
    await fetch("/api/logout", { method: "POST" });
}
