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

// 处理修改密码
async function handleChangePassword(event) {
    event.preventDefault();

    const oldPassword = document.getElementById('old-password').value;
    const newPassword = document.getElementById('new-password').value;
    const confirmPassword = document.getElementById('confirm-new-password').value;

    // 验证新密码和确认密码是否一致
    if (newPassword !== confirmPassword) {
        showCustomAlert('新密码和确认密码不一致');
        return;
    }

    // 验证密码强度（可选）
    if (newPassword.length < 6) {
        showCustomAlert('密码长度至少为6位');
        return;
    }

    try {
        const response = await fetch('/change_my_password', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                old_password: oldPassword,
                new_password: newPassword
            })
        });

        if (response.ok) {
            const data = await response.json();
            if (data.success) {
                showCustomAlert('密码修改成功');
                goBackToLogin();
            } else {
                showCustomAlert('修改失败: ' + (data.message || '未知错误'));
            }
        } else {
            const errorData = await response.json();
            showCustomAlert('修改失败: ' + (errorData.error || '未知错误'));
        }
    } catch (error) {
        showCustomAlert('请求出错，请稍后再试');
        console.error(error);
    }
}

// 返回设置页面
function goBackToSettings() {
    window.location.href = '/settings';
}

function goBackToLogin() {
    window.location.href = '/login';
}

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', () => {
    initTheme();

    // 自动聚焦第一个输入框
    document.getElementById('old-password').focus();
});