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

// 返回翻译主界面的函数
function goBackToTranslate() {
    window.location.href = '/';
}

// 从 Cookie 中获取用户名
function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

// 格式化时间
function formatTime(time) {
    return new Date(time).toLocaleString();
}

// 获取历史记录的函数
async function fetchHistory() {
    try {
        const response = await fetch('/trans_history', {
            method: 'GET',
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();


        if (data.success) {
            displayHistory(data.data);
            // console.log(data)
        } else {
            console.error('获取历史记录失败:', data.message || '未知错误');
            showCustomAlert('获取历史记录失败: ' + (data.message || '未知错误'));
        }
    } catch (error) {
        console.error('获取历史记录时出错:', error);
        if (error.message.includes('401')) {
            window.location.href = '/login';
        } else {
            showCustomAlert('获取历史记录时出错: ' + error.message);
        }
    }
}

// 显示历史记录
function displayHistory(data) {
    const table = document.getElementById('history-table');
    if (!table) {
        console.error('未找到历史记录表格元素');
        return;
    }

    const tableBody = table.getElementsByTagName('tbody')[0];
    tableBody.innerHTML = '';

    if (data == null) {
        console.log("fgh")
        const row = tableBody.insertRow();
        const cell = row.insertCell();
        cell.colSpan = 4;
        cell.textContent = '暂无历史记录';
        cell.style.textAlign = 'center';
        cell.style.padding = '20px';
        return;
    }

    data.forEach((record, index) => {
        const row = tableBody.insertRow();
        const originalTextCell = row.insertCell();
        const translatedTextCell = row.insertCell();
        const translationTimeCell = row.insertCell();
        const deleteCell = row.insertCell();

        originalTextCell.textContent = record.original_text;
        translatedTextCell.textContent = record.translated_text;
        translationTimeCell.textContent = formatTime(record.translation_time);

        // 添加删除按钮
        const deleteButton = document.createElement('button');
        deleteButton.innerHTML = '<i class="fas fa-trash-alt"></i> 删除';
        deleteButton.addEventListener('click', () => {
            if (confirm('确定要删除这条记录吗？')) {
                deleteHistoryRecord(record.original_text, record.translation_time);
            }
        });
        deleteCell.appendChild(deleteButton);
    });
}

// 搜索历史记录
function searchTable() {
    const input = document.getElementById("searchInput");
    const filter = input.value.toUpperCase();
    const table = document.getElementById("history-table");
    const tr = table.getElementsByTagName("tr");

    for (let i = 1; i < tr.length; i++) {
        let found = false;
        const tdList = tr[i].getElementsByTagName("td");

        for (let j = 0; j < tdList.length - 1; j++) { // 不搜索操作列
            const td = tdList[j];
            if (td) {
                const txtValue = td.textContent || td.innerText;
                if (txtValue.toUpperCase().includes(filter)) {
                    found = true;
                    break;
                }
            }
        }

        tr[i].style.display = found ? "" : "none";
    }
}

// 删除历史记录
async function deleteHistoryRecord(originalText, translationTime) {
    try {
        const username = getCookie('username');
        const formData = {
            username: username,
            original_text: originalText,
            translation_time: translationTime
        };

        const response = await fetch('/delete_history_record', {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        });

        if (response.ok) {
            showCustomAlert('历史记录删除成功');
            await fetchHistory();
        } else {
            const errorData = await response.json();
            showCustomAlert('删除失败: ' + (errorData.error || '未知错误'));
        }
    } catch (error) {
        showCustomAlert('请求出错，请稍后再试');
        console.error(error);
    }
}

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', () => {
    initTheme();
    fetchHistory();
});