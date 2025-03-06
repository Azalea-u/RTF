import { renderPage } from "./router.js";

addEventListener('load', () => {
    renderPage(window.location.pathname);
});
document.addEventListener('DOMContentLoaded', () => {
    const alert = document.getElementById('alert');
    alert.classList.add('hidden');
    alert.querySelector('.close').addEventListener('click', () => {
        alert.classList.add('hidden');
    });
});