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

function goToUserPermissionPage(){
    window.location.href = '/userPermission_page'
}

// 返回翻译主界面的函数
function goBackToTranslate() {
    window.location.href = '/';
}

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', () => {
    initTheme();
});

// 从后端获取用户的权限
async function fetchUser(){
    await fetch('/check_userPermission',{
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        }
    })
        .then(response => {
            return response.json();
        })
        .then(data => {
                console.log(data);
                if (data.success){
                    window.location.href = data.redirected
                }else{
                    console.error("错误信息：", data.message)
                    alert("用户权限不足")
                }
        })
        .catch(error => {
            console.error('获取用户权限时出错:', error);
            alert('获取用户权限数据失败1，请稍后重试');
        });
}