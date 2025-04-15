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

function validateForm() {
    const username = document.getElementById('username').value;
    const useraccount = document.getElementById('useraccount').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirmPassword').value;

    console.log('用户名:', username); // 添加日志输出

    const usernameRegex = /^[a-zA-Z0-9_]{3,20}$/;
    if (!usernameRegex.test(username)) {
        alert('用户名只能包含字母、数字和下划线，且长度在 3 到 20 个字符之间');
        return false;
    }

    const useraccountRegex = /^[a-zA-Z0-9]{6,15}$/;
    if (!useraccountRegex.test(useraccount)) {
        alert('账号只能是由数字和字母组成的6~12位')
        return false;
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
        alert('请输入有效的邮箱地址');
        return false;
    }

    if (password.length < 6) {
        alert('密码长度至少为 6 个字符');
        return false;
    }

    if (password !== confirmPassword) {
        alert('两次输入的密码不一致');
        return false;
    }
    return true;
}

async function handleRegister(event) {
    event.preventDefault(); // 阻止表单默认提交行为
    if (!validateForm()) return;

    const username = document.getElementById('username').value;
    const useraccount = document.getElementById('useraccount').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    const formData = new URLSearchParams();
    formData.append('username', username);
    formData.append('useraccount',useraccount);
    formData.append('email', email);
    formData.append('password', password);

    try {
        const response = await fetch('/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
                'X-Content-Type-Options':'nosniff',
                'X-Frame-Options':'DENY',
            },
            body: formData
        });

        // 4. 解析响应数据
        const data = await response.json();

        console.log(data)

        if (data.success) {
            window.location.href = data.redirected; // 注册成功，重定向到登录页面
        } else {
            alert(data.error); // 显示错误信息
        }
    } catch (error) {
        alert('请求出错，请稍后再试');
        console.error(error);
    }
}

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', () => {
    initTheme();

    // 自动聚焦用户名输入框
    document.getElementById('username').focus();
});
