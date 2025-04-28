// 页面加载完成后获取用户信息
window.onload = async function () {
    try {
        const response = await fetch('/get_user_info', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        });

        if (response.ok) {
            const data = await response.json();
            document.getElementById('username').textContent = data.username;
            document.getElementById('user-email').textContent = data.email;
            document.getElementById('useraccount').textContent = data.useraccount;
        } else {
            const errorData = await response.json();
            showCustomAlert(errorData.error);
        }
    } catch (error) {
        showCustomAlert('请求出错，请稍后再试');
        console.error(error);
    }
};

function goToChangeUsernamePage() {
    window.location.href = '/change_username_page';
}
function goToChangeEmailPage() {
    window.location.href = '/change_email_page';
}
// 返回翻译主界面的函数
function goBackToTranslate() {
    window.location.href = '/';
}

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

// 初始化
window.addEventListener('DOMContentLoaded', () => {
    initTheme(); // 这里调用了initTheme函数
});