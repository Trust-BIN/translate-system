// 显示自定义提示框
function showCustomAlert(message) {
    const overlay = document.getElementById('custom-alert-overlay');
    const messageElement = document.getElementById('custom-alert-message');
    messageElement.textContent = message;
    overlay.style.display = 'flex';
}

// 隐藏自定义提示框
function hideCustomAlert() {
    const overlay = document.getElementById('custom-alert-overlay');
    overlay.style.display = 'none';
}

// 初始化主题
function initTheme() {
    const savedTheme = localStorage.getItem('theme') || 'light';
    document.documentElement.setAttribute('data-theme', savedTheme);

    // 设置正确的图标
    const themeIcon = document.querySelector('.theme-switch i');
    themeIcon.className = savedTheme === 'dark' ? 'fas fa-moon' : 'fas fa-sun';
}

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', () => {
    initTheme();
});