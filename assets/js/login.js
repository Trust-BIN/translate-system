// 登录处理函数
async function handleLogin(event) {
    event.preventDefault(); // 阻止表单默认提交行为
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    console.log('用户名:', username);
    console.log('密码:', password);

    // 确保用户名和密码不为空
    if (!username || !password) {
        alert('请输入用户名和密码');
        return;
    }

    const formData = new URLSearchParams();
    formData.append('username', username);
    formData.append('password', password);
    console.log("测试")
    console.log(formData.get('username'));

    try {
        const response = await fetch('/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            },
            body: formData
        });

        if (response.ok) {
            if (response.redirected) {
                window.location.href = response.url; // 登录成功，重定向到首页
            }
        } else {
            const data = await response.json();
            alert(data.error); // 显示错误信息
        }
    } catch (error) {
        alert('请求出错，请稍后再试');
        console.error(error);
    }
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

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', () => {
    initTheme();

    // 自动聚焦用户名输入框
    document.getElementById('username').focus();
});