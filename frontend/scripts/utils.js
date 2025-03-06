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

export function showAlert(message, type) {
    const alert = document.getElementById('alert');
    alert.classList.remove('hidden');
    alert.textContent = message;
    alert.classList.add(type);
    setTimeout(() => {
        alert.classList.add('hidden');
    }, 2000);
}