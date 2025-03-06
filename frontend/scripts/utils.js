export async function checkAuth() {
    try {
        const response = await fetch('/api/check-auth', {
            method: 'GET',
            credentials: 'include',
        });
        return response.ok;
    } catch (error) {
        console.error("Error checking authentication:", error);
        return false;
    }
}