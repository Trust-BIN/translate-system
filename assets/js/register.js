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
        showCustomAlert('用户名只能包含字母、数字和下划线，且长度在 3 到 20 个字符之间');
        return false;
    }

    const useraccountRegex = /^[a-zA-Z0-9]{6,15}$/;
    if (!useraccountRegex.test(useraccount)) {
        showCustomAlert('账号只能是由数字和字母组成的6~12位')
        return false;
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
        showCustomAlert('请输入有效的邮箱地址');
        return false;
    }

    if (password.length < 6) {
        showCustomAlert('密码长度至少为 6 个字符');
        return false;
    }

    if (password !== confirmPassword) {
        showCustomAlert('两次输入的密码不一致');
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
            showCustomAlert(data.error); // 显示错误信息
        }
    } catch (error) {
        showCustomAlert('请求出错，请稍后再试');
        console.error(error);
    }
}

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', () => {
    initTheme();

    // 自动聚焦用户名输入框
    document.getElementById('username').focus();
});


// // 用户头像上传
// document.addEventListener('DOMContentLoaded', function() {
//     const fileInput = document.getElementById('fileInput');
//     const uploadBtn = document.getElementById('uploadBtn');
//     const avatarPreview = document.getElementById('avatarPreview');
//     const previewImage = document.getElementById('previewImage');
//     const progressContainer = document.getElementById('progressContainer');
//     const progressBar = document.getElementById('progressBar');
//     const errorMessage = document.getElementById('errorMessage');
//
//     // 点击上传按钮触发文件选择
//     uploadBtn.addEventListener('click', function() {
//         fileInput.click();
//     });
//
//     // 文件选择变化时处理
//     fileInput.addEventListener('change', function(e) {
//         const file = e.target.files[0];
//         if (!file) return;
//
//         // 验证文件类型和大小
//         if (!validateFile(file)) {
//             return;
//         }
//
//         // 显示预览
//         showPreview(file);
//
//         // 模拟上传过程
//         simulateUpload(file);
//     });
//
//     // 验证文件
//     function validateFile(file) {
//         const validTypes = ['image/jpeg', 'image/png', 'image/gif'];
//         const maxSize = 2 * 1024 * 1024; // 2MB
//
//         errorMessage.style.display = 'none';
//
//         if (!validTypes.includes(file.type)) {
//             showError('请上传 JPG、PNG 或 GIF 格式的图片');
//             return false;
//         }
//
//         if (file.size > maxSize) {
//             showError('图片大小不能超过 2MB');
//             return false;
//         }
//
//         return true;
//     }
//
//     // 显示错误信息
//     function showError(message) {
//         errorMessage.textContent = message;
//         errorMessage.style.display = 'block';
//     }
//
//     // 显示图片预览
//     function showPreview(file) {
//         const reader = new FileReader();
//
//         reader.onload = function(e) {
//             previewImage.src = e.target.result;
//             previewImage.style.display = 'block';
//             avatarPreview.querySelector('i').style.display = 'none';
//         };
//
//         reader.readAsDataURL(file);
//     }
//
//     // 模拟上传过程
//     function simulateUpload(file) {
//         progressContainer.style.display = 'block';
//         progressBar.style.width = '0%';
//
//         // 模拟上传进度
//         let progress = 0;
//         const interval = setInterval(() => {
//             progress += Math.random() * 10;
//             if (progress >= 100) {
//                 progress = 100;
//                 clearInterval(interval);
//
//                 // 上传完成后的操作
//                 setTimeout(() => {
//                     alert('头像上传成功！');
//                     progressContainer.style.display = 'none';
//                 }, 300);
//             }
//             progressBar.style.width = progress + '%';
//         }, 200);
//     }
//
//     // 拖放上传功能
//     avatarPreview.addEventListener('dragover', function(e) {
//         e.preventDefault();
//         this.style.borderColor = '#409eff';
//         this.style.backgroundColor = '#ecf5ff';
//     });
//
//     avatarPreview.addEventListener('dragleave', function(e) {
//         e.preventDefault();
//         this.style.borderColor = '#ccc';
//         this.style.backgroundColor = '#f9f9f9';
//     });
//
//     avatarPreview.addEventListener('drop', function(e) {
//         e.preventDefault();
//         this.style.borderColor = '#ccc';
//         this.style.backgroundColor = '#f9f9f9';
//
//         const file = e.dataTransfer.files[0];
//         if (!file) return;
//
//         fileInput.files = e.dataTransfer.files;
//
//         if (!validateFile(file)) {
//             return;
//         }
//
//         showPreview(file);
//         simulateUpload(file);
//     });
// });