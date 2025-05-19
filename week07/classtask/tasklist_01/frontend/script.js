document.addEventListener('DOMContentLoaded', () => {
    const taskInput = document.getElementById('task-input');
    const addTaskButton = document.getElementById('add-task-button');
    const taskListElement = document.getElementById('task-list');
    const themeToggleButton = document.getElementById('theme-toggle-button');
    const body = document.body;
    const filterButtonsContainer = document.querySelector('.filter-buttons');

    const API_BASE_URL = '/api/tasks'; 

    let localTasksCache = []; 

    const urlParams = new URLSearchParams(window.location.search);
    let currentFilter = urlParams.get('filter') || 'all';

    if (!['all', 'done', 'undone'].includes(currentFilter)) {
        currentFilter = 'all';
    }

    const applyTheme = (theme) => {
        body.classList.remove('light-theme', 'dark-theme');
        body.classList.add(theme);
        try {
            localStorage.setItem('todoTheme', theme);
        } catch (e) {
            console.warn("LocalStorage is not available. Theme won't be saved.");
        }
    };

    themeToggleButton.addEventListener('click', () => {
        const newTheme = body.classList.contains('dark-theme') ? 'light-theme' : 'dark-theme';
        applyTheme(newTheme);
    });

    const loadTheme = () => {
        let savedTheme = 'light-theme';
        try {
            savedTheme = localStorage.getItem('todoTheme') || 'light-theme';
        } catch (e) { /* ignore */ }
        applyTheme(savedTheme);
    };

    const renderTasks = (tasksToDisplay) => {
        taskListElement.innerHTML = ''; 

        updateFilterButtonsActiveState();

        if (!tasksToDisplay || tasksToDisplay.length === 0) {
            const emptyMessage = document.createElement('li');
            let messageText = '还没有任务！';
            if (currentFilter === 'done') messageText = '没有已完成的任务。';
            else if (currentFilter === 'undone') messageText = '所有任务都已完成！';

            emptyMessage.textContent = messageText;
            emptyMessage.style.textAlign = 'center';
            emptyMessage.style.padding = '10px';
            emptyMessage.style.color = body.classList.contains('dark-theme') ? '#aaa' : '#777';
            taskListElement.appendChild(emptyMessage);
            return;
        }

        tasksToDisplay.forEach(task => {
            const listItem = document.createElement('li');
            listItem.classList.add('task-item');
            listItem.dataset.id = task.id;
            if (task.completed) {
                listItem.classList.add('completed');
            }

            const checkbox = document.createElement('input');
            checkbox.type = 'checkbox';
            checkbox.classList.add('task-checkbox');
            checkbox.checked = task.completed;
            checkbox.addEventListener('change', () => toggleTaskCompleteAPI(task.id, !task.completed));

            const taskNameSpan = document.createElement('span');
            taskNameSpan.classList.add('task-name');
            taskNameSpan.textContent = task.name;

            const taskTimeSpan = document.createElement('span');
            taskTimeSpan.classList.add('task-time');
            taskTimeSpan.textContent = task.addedTime || "N/A";

            const deleteButton = document.createElement('button');
            deleteButton.classList.add('delete-task-button');
            deleteButton.textContent = '删除';
            deleteButton.addEventListener('click', () => deleteTaskAPI(task.id));

            listItem.appendChild(checkbox);
            listItem.appendChild(taskNameSpan);
            listItem.appendChild(taskTimeSpan);
            listItem.appendChild(deleteButton);
            taskListElement.appendChild(listItem);
        });
        updateFilterButtonsActiveState();
    };

    const formatFrontendTime = (dateObj) => {
        const month = (dateObj.getMonth() + 1).toString().padStart(2, '0');
        const day = dateObj.getDate().toString().padStart(2, '0');
        const hours = dateObj.getHours().toString().padStart(2, '0');
        const minutes = dateObj.getMinutes().toString().padStart(2, '0');
        const seconds = dateObj.getSeconds().toString().padStart(2, '0');
        return `${month}月${day}日 ${hours}:${minutes}:${seconds}`;
    };

    const fetchTasksAPI = async (filter = 'all') => {
        let url = `${API_BASE_URL}/`; 
        if (filter === 'done' || filter === 'undone' || filter === 'all') {
             url = `${API_BASE_URL}/${filter}`;
        }
        if (filter === 'all' && url.endsWith('/all') === false) {
            url = `${API_BASE_URL}/all`;
        }


        try {
            const response = await fetch(url);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const tasksFromServer = await response.json();
            localTasksCache = tasksFromServer || [];
            renderTasks(localTasksCache);
            saveTasksToLocalStorage(localTasksCache); 
        } catch (error) {
            console.error(`无法从服务器获取任务 (${filter}):`, error);
            alert(`加载任务失败 (${filter})。尝试加载本地缓存。`);
            loadTasksFromLocalStorageAndRender(filter); 
        }
    };

    const addTaskAPI = async () => {
        const taskName = taskInput.value.trim();
        if (taskName === '') {
            alert('任务内容不能为空！');
            return;
        }

        const newTaskData = {
            name: taskName,
            completed: false,
            addedTime: formatFrontendTime(new Date()) 
        };

        try {
            const response = await fetch(`${API_BASE_URL}/`, { 
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(newTaskData),
            });
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            taskInput.value = '';
            fetchTasksAPI(currentFilter);
        } catch (error) {
            console.error("添加任务失败:", error);
            alert("添加任务失败，请检查网络连接或服务器状态。");
        }
    };

    const toggleTaskCompleteAPI = async (taskId, newCompletedState) => {
        try {
            const response = await fetch(`${API_BASE_URL}/${taskId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ completed: newCompletedState }),
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            fetchTasksAPI(currentFilter);
        } catch (error) {
            console.error("更新任务状态失败:", error);
            alert("更新任务状态失败，请稍后再试。");
        }
    };

    const deleteTaskAPI = async (taskId) => {
        if (!confirm('确定要删除这个任务吗？')) {
            return;
        }
        try {
            const response = await fetch(`${API_BASE_URL}/${taskId}`, {
                method: 'DELETE',
            });
            if (!response.ok && response.status !== 204) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            fetchTasksAPI(currentFilter);
        } catch (error) {
            console.error("删除任务失败:", error);
            alert("删除任务失败，请稍后再试。");
        }
    };

    const saveTasksToLocalStorage = (tasksToSave) => {
        try {
            localStorage.setItem(`todoAppTasks_${currentFilter}`, JSON.stringify(tasksToSave));
            localStorage.setItem('todoAppTasks_lastFilter', currentFilter); 
        } catch (e) {
            console.warn("LocalStorage is not available. Tasks won't be saved locally.");
        }
    };

    const loadTasksFromLocalStorageAndRender = (filterToLoad) => {
        try {
            const storedTasks = localStorage.getItem(`todoAppTasks_${filterToLoad}`);
            if (storedTasks) {
                localTasksCache = JSON.parse(storedTasks);
                renderTasks(localTasksCache);
            } else {
                renderTasks([]); 
            }
        } catch (e) {
            console.warn("LocalStorage data is corrupted or unavailable. Displaying empty list.");
            renderTasks([]);
        }
    };    const loadLastFilter = () => {
        try {
            return localStorage.getItem('todoAppTasks_lastFilter') || 'all';
        } catch(e) {
            return 'all';
        }
    };

    filterButtonsContainer.addEventListener('click', (event) => {
        if (event.target.tagName === 'BUTTON') {
            const newFilter = event.target.dataset.filter;
            if (newFilter !== currentFilter) {
                currentFilter = newFilter;

                const newURL = new URL(window.location);
                newURL.searchParams.set('filter', newFilter);
                window.history.pushState({}, '', newURL);
                
                fetchTasksAPI(currentFilter); 

                const buttons = filterButtonsContainer.querySelectorAll('button');
                buttons.forEach(btn => {
                    btn.classList.remove('active');
                });
                event.target.classList.add('active');
            }
        }
    });

    const updateFilterButtonsActiveState = () => {
        const buttons = filterButtonsContainer.querySelectorAll('button');
        buttons.forEach(button => {
            if (button.dataset.filter === currentFilter) {
                button.classList.add('active');
            } else {
                button.classList.remove('active');
            }
        });
    };


    addTaskButton.addEventListener('click', addTaskAPI);
    taskInput.addEventListener('keypress', (event) => {
        if (event.key === 'Enter') {
            addTaskAPI();
        }
    });    loadTheme();

    if (!urlParams.has('filter')) {
        currentFilter = loadLastFilter();
        const newURL = new URL(window.location);
        newURL.searchParams.set('filter', currentFilter);
        window.history.replaceState({}, '', newURL);
    }
    
    updateFilterButtonsActiveState(); 
    fetchTasksAPI(currentFilter); 
});