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

// 处理注销账号
async function handleDeleteAccount() {
    if (confirm('确定要注销账号吗？此操作将永久删除您的所有数据且不可恢复！')) {
        try {
            const response = await fetch('/delete_account', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            const data = await response.json()

            if (data.success) {
                alert('账号已成功注销');
                window.location.href = '/login';
            } else {
                const errorData = await response.json();
                alert('注销失败: ' + (errorData.error || '未知错误'));
            }
        } catch (error) {
            alert('请求出错，请稍后再试');
            console.error(error);
        }
    }
}

function goToChangePasswordPage() {
    window.location.href = '/change_password';
}

// 返回翻译主界面的函数
function goBackToTranslate() {
    window.location.href = '/';
}

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', () => {
    initTheme();
});