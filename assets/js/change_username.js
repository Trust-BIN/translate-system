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

// 处理修改用户名
async function handleChangeUsername(event) {
    event.preventDefault();

    const newUsername = document.getElementById('new-username').value.trim();

    if (!newUsername) {
        alert('请输入有效用户名');
        return;
    }

    try {
        const response = await fetch('/change_username', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                new_username: newUsername
            })
        });

        if (response.ok) {
            const data = await response.json();
            if (data.success) {
                alert('用户名修改成功');
                goBackToPersonCenter();
            } else {
                alert('修改失败: ' + (data.message || '未知错误'));
            }
        } else {
            const errorData = await response.json();
            alert('修改失败: ' + (errorData.error || '未知错误'));
        }
    } catch (error) {
        alert('请求出错，请稍后再试');
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
    document.getElementById('new-username').focus();
});