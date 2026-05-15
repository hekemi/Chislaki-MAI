// Этот файл отвечает только за интерфейс: он не содержит численных методов.
// Его задача — получать данные из Go-сервера и рисовать удобную навигацию.

const taskNav = document.getElementById('taskNav');
const taskTitle = document.getElementById('taskTitle');
const taskDescription = document.getElementById('taskDescription');
const taskMeta = document.getElementById('taskMeta');
const serverStatus = document.getElementById('serverStatus');

let tasksCache = [];

// renderTaskButtons строит левую панель из данных, пришедших с backend.
function renderTaskButtons(tasks) {
  taskNav.innerHTML = '';

  tasks.forEach((task) => {
    const button = document.createElement('button');
    button.className = 'task-button';
    button.type = 'button';
    button.dataset.id = String(task.id);
    button.innerHTML = `
      <span class="task-button__title">${task.title}</span>
      <span class="task-button__meta">Модуль №${task.id}</span>
    `;

    button.addEventListener('click', () => {
      // Обновляем hash, чтобы можно было быстро открывать конкретный модуль.
      window.location.hash = `task-${task.id}`;
      activateTask(task.id);
    });

    taskNav.appendChild(button);
  });
}

// activateTask обновляет центральную область интерфейса.
function activateTask(taskId) {
  const task = tasksCache.find((item) => item.id === taskId);
  if (!task) {
    return;
  }

  document.querySelectorAll('.task-button').forEach((button) => {
    button.classList.toggle('active', Number(button.dataset.id) === taskId);
  });

  taskTitle.textContent = task.title;
  taskDescription.textContent = task.description;
  taskMeta.textContent = JSON.stringify(task, null, 2);
}

// loadTasks получает данные из Go API и создает интерфейс.
async function loadTasks() {
  try {
    const response = await fetch('/api/tasks');
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }

    tasksCache = await response.json();
    renderTaskButtons(tasksCache);

    const hash = window.location.hash.replace('#task-', '');
    const taskIdFromHash = Number(hash) || tasksCache[0].id;
    activateTask(taskIdFromHash);

    serverStatus.textContent = 'API подключен';
  } catch (error) {
    console.error(error);
    taskTitle.textContent = 'Ошибка загрузки';
    taskDescription.textContent = 'Сервер не ответил или API недоступно.';
    taskMeta.textContent = String(error);
    serverStatus.textContent = 'Ошибка';
  }
}

// Если пользователь меняет hash вручную, интерфейс тоже должен обновиться.
window.addEventListener('hashchange', () => {
  const hash = window.location.hash.replace('#task-', '');
  const taskId = Number(hash);
  if (!Number.isNaN(taskId)) {
    activateTask(taskId);
  }
});

loadTasks();
