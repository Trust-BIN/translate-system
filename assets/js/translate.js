// DOM 参数引用
const elements = {
    sourceText: document.getElementById('sourceText'),
    sourceLang: document.getElementById('sourceLang'),
    targetLang: document.getElementById('targetLang'),
    result: document.getElementById('result'),
    loading: document.getElementById('loading'),
    inputCount: document.getElementById('inputCount'),
};

let isTranslating = false;

// document.getElementById('history-btn').addEventListener('click', fetchHistory);

// 更新字符计数
function updateCharCount() {
    const text = elements.sourceText.value;
    const maxLength = parseInt(elements.sourceText.getAttribute('maxlength'));
    elements.inputCount.textContent = `${text.length}/${maxLength}`;
}

// 验证输入
function validateInput(text) {
    if (text.length === 0) {
        alert('请输入需要翻译的文本');
        return false;
    }
    return true;
}

// 设置加载状态
function setLoadingState(isLoading) {
    isTranslating = isLoading;
    elements.loading.style.display = isLoading ? 'flex' : 'none';
}

// 简单的语言检测函数（这里只是示例，实际需要使用更准确的方法）
function detectLanguage() {

    const text = elements.sourceText.value;
    const sourceLang = elements.sourceLang;
    if (sourceLang.value.includes("自动检测") ) {
        let detectedLang = '';
        if (/[\u4e00-\u9fa5]/.test(text)) {
            detectedLang = 'zh';
        } else if (/[a-zA-Z]/.test(text)) {
            detectedLang = 'en';
        }

        if (detectedLang) {
            const langName = getLanguageName(detectedLang);
            sourceLang.value = `自动检测（${langName}）`;
        } else {
            sourceLang.value = '自动检测';
        }
    }

}

// 根据语言代码获取语言名称
function getLanguageName(langCode) {
    switch (langCode) {
        case 'en':
            return '英语';
        // case '日语':
        //     return '日语';
        // case '韩语':
        //     return '韩语';
        // case '法语':
        //     return '法语';
        case 'zh':
            return '汉语';
        // case '印尼语':
        //     return '印尼语';
        // case '俄语':
        //     return '俄语';
        // case '阿拉伯语':
        //     return '阿拉伯语';
        // case '孟加拉语':
        //     return '孟加拉语'
        default:
            return '汉语';
    }
}

// 交换语言
function swapLanguages() {
    const sourceValue = elements.sourceLang.value;
    const targetValue = elements.targetLang.value;
    elements.sourceLang.value = targetValue;
    elements.targetLang.value = sourceValue;
}

// 向自己的后端服务器发送翻译请求
async function translateText(text) {
    const sourceLang = elements.sourceLang.value;
    // const targetLang = elements.targetLang.options[elements.targetLang.selectedIndex].text;
    const targetLang = elements.targetLang.value;

    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 60000);

    try {
        const messages = [
            {
                "role": "system",
                "content": "你是一个翻译专家，将用户输入的文本翻译成目标语种。用户可以向助手发送需要翻译的内容，助手会回答相应的翻译结果，并确保符合目标语种的语言习惯，你可以调整语气和风格，并考虑到某些词语的文化内涵和地区差异。同时作为翻译家，需将原文翻译成具有信达雅标准的译文。\"信\" 即忠实于原文的内容与意图；\"达\" 意味着译文应通顺易懂，表达清晰；\"雅\" 则追求译文的文化审美和语言的优美。目标是创作出既忠于原作精神，又符合目标语言文化和读者审美的翻译。只需输出翻译结果，无需任何注释内容。"
            },
            {
                "role": "user",
                "content": `请将以下 ${sourceLang === 'auto' ? '文本' : sourceLang} 翻译成 ${targetLang},并严格按照"翻译结果:"的格式输出，不要将"翻译结果"这四个字更换成别的语种。需要翻译的内容: ${text}`,
                "s_text":`${text}`,
            }
        ];

        const
            response = await fetch('/translate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                model: "deepseek-chat",
                messages: messages,
                stream: false
            }),
            signal: controller.signal
        });

        console.log(response);

        clearTimeout(timeoutId);

        if (!response.ok) {
            throw new Error(`HTTP 错误 ${response.status}`);
        }

        const data = await response.json();
        if (data.choices && data.choices.length > 0) {
            return data.choices[0].message.content;
        }
        return '未获取到有效翻译结果';

    } catch (error) {
        if (error.name === 'AbortError') {
            throw new Error('请求超时，请检查服务状态');
        }
        throw error;
    }
}

// 优化后的显示翻译结果函数
function displayTranslation(text) {
    // 定义可能的结果前缀（兼容不同格式）
    const resultPrefixes = ["翻译结果：", "翻译结果:", "Translation:", "译文："];

    let translatedText = text;

    // 检查所有可能的前缀
    for (const prefix of resultPrefixes) {
        const resultStart = text.indexOf(prefix);
        if (resultStart !== -1) {
            translatedText = text.substring(resultStart + prefix.length).trim();
            break; // 找到第一个匹配的前缀后就退出循环
        }
    }

    // 处理可能的多余换行和空格
    translatedText = translatedText
        .replace(/\n+/g, '\n')  // 合并连续换行
        .replace(/^\n+|\n+$/g, '')  // 去除开头和结尾的换行
        .trim();

    // 处理换行并显示
    elements.result.innerHTML = translatedText
        .split('\n')
        .filter(line => line.trim().length > 0)  // 过滤空行
        .map(line => `<p>${line}</p>`)
        .join('');
    return translatedText
}

// 处理翻译错误
function handleTranslationError(error) {
    const errorMessage = `翻译出错: ${error.message}`;
    elements.result.innerHTML = `<p class="error-message">${errorMessage}</p>`;
    console.error(error);
}

// 清空输入文本
function clearText() {
    elements.sourceText.value = '';
    elements.sourceLang.options[0].text = '自动检测';
    elements.sourceLang.value = 'auto';
    updateCharCount();
}

// 复制翻译结果
function copyResult() {
    const text = elements.result.innerText;
    navigator.clipboard.writeText(text)
        .then(() => alert('翻译结果已复制到剪贴板'))
        .catch(err => console.error('复制失败:', err));
}

// 核心翻译逻辑
async function handleTranslation() {
    if (isTranslating) return;

    const text = elements.sourceText.value.trim();
    if (!validateInput(text)) return;

    try {
        setLoadingState(true);
        const translation = await translateText(text);
        displayTranslation(translation);
    } catch (error) {
        handleTranslationError(error);
    } finally {
        setLoadingState(false);
    }
}

// 初始化字符计数
updateCharCount();
elements.sourceText.addEventListener('input', updateCharCount);

// source源文本监听
elements.sourceText.addEventListener('input', detectLanguage);


// 退出登录
async function logout() {
    try {
        const response = await fetch('/logout', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        });

        if (!response.ok) {
            throw new Error(`HTTP 错误 ${response.status}`);
        }

        // 手动重定向到登录页面
        window.location.href = '/login';
    } catch (error) {
        console.error('退出登录出错:', error);
        alert('退出登录出错，请稍后再试');
    }
}

// 获取 Cookie 中的用户名
function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

// 将用户名显示在页面上
document.addEventListener('DOMContentLoaded', function () {
    const username = getCookie('username');
    if (username) {
        const usernameDisplay = document.getElementById('username-display');
        usernameDisplay.textContent = username;
    }
});

// 主题切换功能
function toggleTheme() {
    const currentTheme = document.documentElement.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';

    document.documentElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);

    // 更新按钮图标状态
    updateThemeIcon(newTheme);

    // 强制重绘侧边栏（解决某些浏览器渐变背景切换问题）
    const sidebar = document.querySelector('.left_nav');
    sidebar.style.display = 'none';
    sidebar.offsetHeight; // 触发重绘
    sidebar.style.display = 'flex';
}

function updateThemeIcon(theme) {
    const themeSwitch = document.querySelector('.theme-switch');
    if (theme === 'dark') {
        themeSwitch.innerHTML = '<i class="fas fa-moon dark-icon"></i>';
    } else {
        themeSwitch.innerHTML = '<i class="fas fa-sun light-icon"></i>';
    }
}

// 初始化主题
function initTheme() {
    const savedTheme = localStorage.getItem('theme') || 'light';
    document.documentElement.setAttribute('data-theme', savedTheme);
    updateThemeIcon(savedTheme);
}

// 页面加载时初始化主题
window.addEventListener('DOMContentLoaded', initTheme);
/*-------------------------------------------------------------------------------------------------------*/
/*-------------------------------------------------------------------------------------------------------*/
/*-------------------------------------------------------------------------------------------------------*/
/*-------------------------------------------------------------------------------------------------------*/
/*-------------------------------------------------------------------------------------------------------*/
/*-------------------------------------------------------------------------------------------------------*/
/*-------------------------------------------------------------------------------------------------------*/
/*-------------------------------------------------------------------------------------------------------*/


