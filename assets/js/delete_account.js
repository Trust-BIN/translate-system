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

// 切换密码可见性
function togglePasswordVisibility(inputId, icon) {
    const input = document.getElementById(inputId);
    if (input.type === 'password') {
        input.type = 'text';
        icon.classList.replace('fa-eye', 'fa-eye-slash');
    } else {
        input.type = 'password';
        icon.classList.replace('fa-eye-slash', 'fa-eye');
    }
}

// 处理注销账号
async function handleDeleteAccount(event) {
    event.preventDefault();

    const password = document.getElementById('password').value;

    try {
        const response = await fetch('/delete_my_account', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                password: password
            })
        });

        const data = await response.json();

        if (data.success) {
            showCustomAlert('账号已成功注销');
            window.location.href = '/login';
        } else {
            showCustomAlert('注销失败: ' + (data.error || '未知错误'));
        }
    } catch (error) {
        showCustomAlert('请求出错，请稍后再试');
        console.error(error);
    }
}

// 返回设置页面的函数
function goBackToSettings() {
    window.location.href = '/settings';
}

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', () => {
    initTheme();
});