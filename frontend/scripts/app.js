import { renderPage } from "./router.js";

addEventListener('load', () => {
    renderPage("/login");
});

const alert = document.getElementById('alert');
if (!alert.classList.contains('hidden')) {
    setTimeout(() => {
        alert.classList.add('hidden');
    }, 3000);
}