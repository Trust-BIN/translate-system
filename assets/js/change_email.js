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

// 处理修改邮箱
async function handleChangeEmail(event) {
    event.preventDefault();

    const newEmail = document.getElementById('new-email').value.trim();

    // 简单的邮箱格式验证
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(newEmail)) {
        showCustomAlert('请输入有效的邮箱地址');
        return;
    }

    try {
        const response = await fetch('/change_email', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                new_email: newEmail
            })
        });

        if (response.ok) {
            const data = await response.json();
            if (data.success) {
                showCustomAlert('邮箱修改成功');
                goBackToPersonCenter();
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

// 返回个人中心
function goBackToPersonCenter() {
    window.location.href = '/person_center';
}

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', () => {
    initTheme();

    // 自动聚焦输入框
    document.getElementById('new-email').focus();
});