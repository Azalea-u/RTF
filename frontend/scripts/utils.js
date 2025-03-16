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

export function TimeAgo(timestamp) {
    const now = new Date();
    const date = new Date(timestamp);
    const seconds = Math.floor((now - date) / 1000);
    if (seconds < 60) {
        return `${seconds}s ago`;
    } else if (seconds < 3600) {
        const minutes = Math.floor(seconds / 60);
        return `${minutes}m ago`;
    } else if (seconds < 86400) {
        const hours = Math.floor(seconds / 3600);
        return `${hours}h ago`;
    } else {
        const days = Math.floor(seconds / 86400);
        return `${days}d ago`;
    }
}
