/**
 * Todo App - Viki Demo
 * A simple, elegant todo app built with vanilla JavaScript
 */

// ===== Data Model =====
class Task {
    constructor(description) {
        this.id = crypto.randomUUID();
        this.description = description.trim();
        this.isComplete = false;
        this.createdAt = new Date().toISOString();
    }
}

// ===== Storage Manager =====
const Storage = {
    KEY: 'viki-todo-tasks',

    load() {
        try {
            const data = localStorage.getItem(this.KEY);
            return data ? JSON.parse(data) : [];
        } catch (e) {
            console.error('Failed to load tasks:', e);
            return [];
        }
    },

    save(tasks) {
        try {
            localStorage.setItem(this.KEY, JSON.stringify(tasks));
        } catch (e) {
            console.error('Failed to save tasks:', e);
        }
    }
};

// ===== App State =====
let tasks = [];
let currentFilter = 'all';

// ===== DOM Elements =====
const taskInput = document.getElementById('taskInput');
const addBtn = document.getElementById('addBtn');
const taskList = document.getElementById('taskList');
const emptyState = document.getElementById('emptyState');
const taskCount = document.querySelector('.task-count');
const filterBtns = document.querySelectorAll('.filter-btn');

// ===== Task Operations =====
function addTask(description) {
    if (!description.trim()) return false;

    const task = new Task(description);
    tasks.unshift(task); // Add to beginning
    Storage.save(tasks);
    render();
    return true;
}

function toggleTask(id) {
    const task = tasks.find(t => t.id === id);
    if (task) {
        task.isComplete = !task.isComplete;
        Storage.save(tasks);
        render();
    }
}

function deleteTask(id) {
    tasks = tasks.filter(t => t.id !== id);
    Storage.save(tasks);
    render();
}

function updateTask(id, newDescription) {
    const task = tasks.find(t => t.id === id);
    if (task && newDescription.trim()) {
        task.description = newDescription.trim();
        Storage.save(tasks);
        render();
    }
}

// ===== Filtering =====
function getFilteredTasks() {
    switch (currentFilter) {
        case 'active':
            return tasks.filter(t => !t.isComplete);
        case 'completed':
            return tasks.filter(t => t.isComplete);
        default:
            return tasks;
    }
}

function setFilter(filter) {
    currentFilter = filter;
    filterBtns.forEach(btn => {
        btn.classList.toggle('active', btn.dataset.filter === filter);
    });
    render();
}

// ===== Rendering =====
function createTaskElement(task) {
    const li = document.createElement('li');
    li.className = `task-item${task.isComplete ? ' completed' : ''}`;
    li.dataset.id = task.id;

    li.innerHTML = `
        <div class="task-checkbox" title="Toggle complete"></div>
        <span class="task-text">${escapeHtml(task.description)}</span>
        <button class="delete-btn" title="Delete task">🗑️</button>
    `;

    // Event: Toggle completion
    li.querySelector('.task-checkbox').addEventListener('click', () => {
        toggleTask(task.id);
    });

    // Event: Delete
    li.querySelector('.delete-btn').addEventListener('click', () => {
        li.style.animation = 'fadeIn 0.3s ease-out reverse';
        setTimeout(() => deleteTask(task.id), 250);
    });

    // Event: Double-click to edit
    const textSpan = li.querySelector('.task-text');
    textSpan.addEventListener('dblclick', () => {
        const input = document.createElement('input');
        input.type = 'text';
        input.className = 'task-text-input';
        input.value = task.description;

        textSpan.replaceWith(input);
        input.focus();
        input.select();

        const save = () => {
            updateTask(task.id, input.value);
        };

        input.addEventListener('blur', save);
        input.addEventListener('keydown', (e) => {
            if (e.key === 'Enter') save();
            if (e.key === 'Escape') render();
        });
    });

    return li;
}

function render() {
    const filtered = getFilteredTasks();

    // Clear and re-render
    taskList.innerHTML = '';
    filtered.forEach(task => {
        taskList.appendChild(createTaskElement(task));
    });

    // Update empty state
    emptyState.classList.toggle('hidden', tasks.length > 0);

    // Update count
    const activeCount = tasks.filter(t => !t.isComplete).length;
    taskCount.textContent = `${activeCount} item${activeCount !== 1 ? 's' : ''} left`;
}

// ===== Utilities =====
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// ===== Event Handlers =====
function handleAddClick() {
    if (addTask(taskInput.value)) {
        taskInput.value = '';
        taskInput.focus();
    }
}

function handleKeyDown(e) {
    if (e.key === 'Enter') {
        handleAddClick();
    }
}

// ===== Initialize =====
function init() {
    // Load tasks from storage
    tasks = Storage.load();

    // Add event listeners
    addBtn.addEventListener('click', handleAddClick);
    taskInput.addEventListener('keydown', handleKeyDown);

    filterBtns.forEach(btn => {
        btn.addEventListener('click', () => setFilter(btn.dataset.filter));
    });

    // Initial render
    render();

    // Focus input
    taskInput.focus();

    console.log('✨ Todo App initialized - Built with Viki SDD Framework');
}

// Start the app
init();
