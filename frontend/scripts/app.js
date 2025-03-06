import { renderPage } from "./router.js";

addEventListener('load', () => {
    renderPage(window.location.pathname);
});

const alert = document.getElementById('alert');
// check if alert class includes hidden
if (!alert.classList.contains('hidden')) {
    setTimeout(() => {
        alert.classList.add('hidden');
    }, 3000);
}