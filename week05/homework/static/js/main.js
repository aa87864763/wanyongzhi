// 等待页面加载完成
document.addEventListener('DOMContentLoaded', function() {
    // 获取DOM元素
    const questionForm = document.getElementById('questionForm');
    const generateBtn = document.getElementById('generateBtn');
    const resultCard = document.getElementById('resultCard');
    const loadingIndicator = document.getElementById('loadingIndicator');
    const resultContent = document.getElementById('resultContent');
    const questionTitle = document.getElementById('questionTitle');
    const optionsList = document.getElementById('optionsList');
    const correctAnswer = document.getElementById('correctAnswer');
    const errorMessage = document.getElementById('errorMessage');

    // 监听表单提交事件
    questionForm.addEventListener('submit', function(e) {
        e.preventDefault();
        generateQuestion();
    });

    // 生成题目
    async function generateQuestion() {
        // 显示加载指示器
        resultCard.style.display = 'block';
        loadingIndicator.style.display = 'block';
        resultContent.style.display = 'none';
        errorMessage.style.display = 'none';
        errorMessage.textContent = '';

        // 获取表单数据
        const formData = new FormData(questionForm);
        const requestData = {
            model: formData.get('model'),
            language: formData.get('language'),
            type: parseInt(formData.get('type')),
            keyword: formData.get('keyword')
        };

        try {
            // 发送请求
            const response = await fetch('/api/questions/create', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(requestData)
            });

            // 解析响应
            const data = await response.json();

            // 处理响应
            if (data.code === 0) {
                displayResult(data.aiRes);
            } else {
                displayError(data.msg || '生成题目失败');
            }
        } catch (error) {
            displayError('请求失败: ' + error.message);
        } finally {
            loadingIndicator.style.display = 'none';
        }
    }

    // 显示结果
    function displayResult(result) {
        // 设置题目标题
        questionTitle.textContent = result.title;

        // 清空选项列表
        optionsList.innerHTML = '';

        // 添加选项
        const rightAnswers = result.right || [];
        const optionLabels = ['A', 'B', 'C', 'D'];

        result.answer.forEach((option, index) => {
            const li = document.createElement('li');
            const isCorrect = rightAnswers.includes(index);
            
            if (isCorrect) {
                li.classList.add('option-correct');
            }
            
            li.textContent = `${optionLabels[index]}. ${option.replace(/^[A-D]\.\s*/, '')}`;
            optionsList.appendChild(li);
        });

        // 显示正确答案
        const answerLetters = rightAnswers.map(index => optionLabels[index]).join(', ');
        correctAnswer.textContent = answerLetters;

        // 显示结果
        resultContent.style.display = 'block';
    }

    // 显示错误
    function displayError(message) {
        errorMessage.textContent = message;
        errorMessage.style.display = 'block';
        resultContent.style.display = 'none';
    }
}); 