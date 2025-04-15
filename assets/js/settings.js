// 主题切换功能
function toggleTheme() {
    const currentTheme = document.documentElement.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';

    document.documentElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);

    // 更新按钮图标
    const themeIcon = document.querySelector('.theme-switch i');
    themeIcon.className = newTheme === 'dark' ? 'fas fa-moon' : 'fas fa-sun';
}

// 初始化主题
function initTheme() {
    const savedTheme = localStorage.getItem('theme') || 'light';
    document.documentElement.setAttribute('data-theme', savedTheme);

    // 设置正确的图标
    const themeIcon = document.querySelector('.theme-switch i');
    themeIcon.className = savedTheme === 'dark' ? 'fas fa-moon' : 'fas fa-sun';
}

function goToChangePasswordPage() {
    window.location.href = '/change_password';
}

function goToDeleteAccountPage(){
    window.location.href = '/delete_account'
}

// 返回翻译主界面的函数
function goBackToTranslate() {
    window.location.href = '/';
}

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', () => {
    initTheme();
});