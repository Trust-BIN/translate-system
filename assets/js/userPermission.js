// 页面加载完成后执行
document.addEventListener('DOMContentLoaded', function() {
    // 初始化主题
    fetchUserPermissions()

    initTheme();

    //获取所有用户信息
    fetchAllRoles();

    // 监听表单提交事件
    const editForm = document.getElementById('editPermissionForm');
    editForm.addEventListener('submit', handleEditPermissionSubmit);
});

// 从后端获取所有角色信息
async function fetchAllRoles() {
    const response = await fetch('/get_all_roles', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        }
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('网络响应不正常');
            }
            return response.json();
        })
        .then(data => {
            if (data.success) {
                renderRoleOptions(data.data);
            } else {
                console.error('获取角色信息失败:', data.error);
            }
        })
        .catch(error => {
            console.error('获取角色信息时出错:', error);
            alert('获取角色信息失败，请稍后重试');
        });
}

// 渲染角色选项
function renderRoleOptions(roles) {
    const roleSelect = document.getElementById('editRole');
    roleSelect.innerHTML = ''; // 清空现有选项

    roles.forEach(role => {
        const option = document.createElement('option');
        option.value = role;
        option.textContent = role;
        roleSelect.appendChild(option);
    });
}

// 从后端获取用户权限数据
async function fetchUserPermissions() {
    const response = await fetch('/get_userPermission_page', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        }
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('网络响应不正常');
            }
            return response.json();
        })
        .then(data => {
            if (data.success) {
                renderUserTable(data.data);
            } else {
                console.error('获取数据失败:', data.error);
            }
        })
        .catch(error => {
            console.error('获取用户权限时出错:', error);
            alert('获取用户权限数据失败，请稍后重试');
        });
}

// 渲染用户表格
function renderUserTable(userData) {
    const tableBody = document.getElementById('userTableBody');
    tableBody.innerHTML = ''; // 清空现有内容

    if (!userData || userData.length === 0) {
        tableBody.innerHTML = '<tr><td colspan="4" class="no-data">暂无用户数据</td></tr>';
        return;
    }

    userData.forEach(user => {
        const row = document.createElement('tr');

        // 用户名单元格
        const usernameCell = document.createElement('td');
        usernameCell.textContent = user.username || '未知用户';

        // 用户账号单元格
        const accountCell = document.createElement('td');
        accountCell.textContent = user.useraccount || 'N/A';

        // 权限单元格
        const roleCell = document.createElement('td');
        roleCell.textContent = user.role || '无权限';

        // 操作单元格
        const actionCell = document.createElement('td');
        const editButton = document.createElement('button');
        editButton.className = 'btn-edit';
        editButton.innerHTML = '<i class="fas fa-edit"></i> 编辑';
        editButton.onclick = function() {
            editUserPermission(user.useraccount, user.role);
        };

        actionCell.appendChild(editButton);

        // 将所有单元格添加到行中
        row.appendChild(usernameCell);
        row.appendChild(accountCell);
        row.appendChild(roleCell);
        row.appendChild(actionCell);

        // 将行添加到表格中
        tableBody.appendChild(row);
    });
}

// 搜索用户功能
function searchUsers() {
    const input = document.getElementById('searchInput');
    const filter = input.value.toUpperCase();
    const table = document.querySelector('.user-table');
    const tr = table.getElementsByTagName('tr');

    for (let i = 1; i < tr.length; i++) { // 从1开始跳过表头
        let match = false;
        const tds = tr[i].getElementsByTagName('td');

        // 检查用户名和账号是否匹配搜索条件
        for (let j = 0; j < 3; j++) { // 前3列(用户名、账号、权限)
            if (tds[j]) {
                const txtValue = tds[j].textContent || tds[j].innerText;
                if (txtValue.toUpperCase().indexOf(filter) > -1) {
                    match = true;
                    break;
                }
            }
        }

        tr[i].style.display = match ? '' : 'none';
    }
}

// 编辑用户权限
function editUserPermission(userId, currentRole) {
    document.getElementById('editUserAccount').value = userId;
    document.getElementById('editRole').value = currentRole;
    openEditModal();
}

// 打开模态框
function openEditModal() {
    const modal = document.getElementById('editPermissionModal');
    modal.style.display = 'flex';
}

// 关闭模态框
function closeEditModal() {
    const modal = document.getElementById('editPermissionModal');
    modal.style.display = 'none';
}

// 处理表单提交
async function handleEditPermissionSubmit(event) {
    event.preventDefault();
    const userAccount = document.getElementById('editUserAccount').value;
    const newRole = document.getElementById('editRole').value;
    console.log(userAccount)
    console.log(newRole)

    try {
        const response = await fetch('/update_user_permission', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                user_Account: userAccount,
                role: newRole
            })
        });

        if (response.ok) {
            const data = await response.json();
            if (data.success) {
                alert('权限修改成功');
                closeEditModal();
                // 重新获取用户权限数据并更新表格
                await fetchUserPermissions();
            } else {
                alert('修改失败: ' + (data.error || '未知错误'));
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

// 主题切换功能
function initTheme() {
    const savedTheme = localStorage.getItem('theme') || 'light';
    document.documentElement.setAttribute('data-theme', savedTheme);
    updateThemeIcon(savedTheme);
}

function toggleTheme() {
    const currentTheme = document.documentElement.getAttribute('data-theme');
    const newTheme = currentTheme === 'light' ? 'dark' : 'light';

    document.documentElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
    updateThemeIcon(newTheme);
}

function updateThemeIcon(theme) {
    const icon = document.querySelector('.theme-switch i');
    if (theme === 'light') {
        icon.className = 'fas fa-moon';
    } else {
        icon.className = 'fas fa-sun';
    }
}