// Viki Dashboard - JavaScript Application

// State
let currentPhase = 'init';
let projectData = {
    name: '',
    spec: '',
    plan: '',
    tasks: []
};

// WebSocket connection
let ws = null;

// Initialize app
document.addEventListener('DOMContentLoaded', () => {
    initWebSocket();
    loadProjectState();
    setupKeyboardShortcuts();
});

// WebSocket connection to backend
function initWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;

    try {
        ws = new WebSocket(wsUrl);

        ws.onopen = () => {
            log('Connected to Viki server', 'success');
        };

        ws.onmessage = (event) => {
            handleServerMessage(JSON.parse(event.data));
        };

        ws.onclose = () => {
            log('Disconnected from server', 'warning');
            // Reconnect after 3 seconds
            setTimeout(initWebSocket, 3000);
        };

        ws.onerror = () => {
            log('WebSocket error - running in demo mode', 'warning');
        };
    } catch (e) {
        log('Running in demo mode (no server)', 'info');
    }
}

// Handle server messages
function handleServerMessage(data) {
    switch (data.type) {
        case 'output':
            log(data.message, data.level || 'info');
            break;
        case 'phase_update':
            updatePhase(data.phase);
            break;
        case 'chat_response':
            addMessage(data.message, 'assistant');
            break;
        case 'status':
            updateStatus(data);
            break;
    }
}

// Phase navigation
function selectPhase(phase) {
    // Update UI
    document.querySelectorAll('.pipeline-step').forEach(step => {
        step.classList.remove('active');
        if (step.dataset.phase === phase) {
            step.classList.add('active');
        }
    });

    // Show corresponding content
    document.querySelectorAll('.phase-content').forEach(content => {
        content.classList.add('hidden');
    });

    const phaseContent = document.getElementById(`phase-${phase}`);
    if (phaseContent) {
        phaseContent.classList.remove('hidden');
    }

    // Update header
    updatePhaseHeader(phase);
    currentPhase = phase;
}

function updatePhaseHeader(phase) {
    const headers = {
        init: { title: '💡 Start Your Project', desc: 'Give your project a name to begin' },
        specify: { title: '📝 Specify Requirements', desc: 'Describe what you want to build' },
        plan: { title: '📐 Architecture Plan', desc: 'AI will design the system architecture' },
        task: { title: '✅ Task Breakdown', desc: 'Break the plan into actionable tasks' },
        execute: { title: '💻 Implementation', desc: 'Generate and review code' },
        review: { title: '🔍 Code Review', desc: 'AI reviews code quality and suggests improvements' }
    };

    const header = headers[phase] || headers.init;
    document.getElementById('phaseTitle').textContent = header.title;
    document.getElementById('phaseDesc').textContent = header.desc;
}

// Run actions
async function runAction(action) {
    log(`Running: viki ${action}...`, 'info');

    // Show spinner
    const btn = event.target;
    const originalText = btn.innerHTML;
    btn.innerHTML = '<span class="spinner"></span> Running...';
    btn.disabled = true;

    try {
        // Collect input data
        let payload = { action };

        switch (action) {
            case 'init':
                payload.name = document.getElementById('projectName').value || 'my-project';
                break;
            case 'specify':
                payload.description = document.getElementById('ideaDesc').value;
                break;
        }

        // Send to server
        if (ws && ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify(payload));
        } else {
            // Demo mode - simulate response
            await simulateAction(action, payload);
        }

    } catch (error) {
        log(`Error: ${error.message}`, 'error');
    } finally {
        btn.innerHTML = originalText;
        btn.disabled = false;
    }
}

// Simulate actions in demo mode
async function simulateAction(action, payload) {
    await sleep(1500); // Simulate server processing

    switch (action) {
        case 'init':
            projectData.name = payload.name;
            log(`✅ Project "${payload.name}" initialized!`, 'success');
            markPhaseComplete('init');
            selectPhase('specify');
            break;

        case 'specify':
            projectData.spec = payload.description;
            log('✅ Specification generated!', 'success');
            document.getElementById('specPreview').textContent =
                payload.description.substring(0, 100) + '...';
            markPhaseComplete('specify');
            selectPhase('plan');
            break;

        case 'plan':
            log('🏗️ Creating architecture plan...', 'info');
            await sleep(1000);
            log('✅ Architecture plan created!', 'success');
            markPhaseComplete('plan');
            selectPhase('task');
            break;

        case 'approve':
            log('✅ Plan approved!', 'success');
            break;

        case 'task':
            log('📋 Breaking into tasks...', 'info');
            await sleep(1000);
            log('✅ Tasks created!', 'success');
            markPhaseComplete('task');
            selectPhase('execute');
            break;

        case 'execute':
            log('💻 Generating implementation...', 'info');
            await sleep(1500);
            log('✅ Code generated!', 'success');
            markPhaseComplete('execute');
            selectPhase('review');
            break;

        case 'review':
            log('🔍 Running code review...', 'info');
            await sleep(1000);
            log('✅ Review complete! Score: 8/10', 'success');
            markPhaseComplete('review');
            break;
    }
}

function markPhaseComplete(phase) {
    const step = document.querySelector(`.pipeline-step[data-phase="${phase}"]`);
    if (step) {
        step.classList.add('completed');
        step.querySelector('.step-status').textContent = '✓';
    }
}

// Chat functionality
function sendMessage() {
    const input = document.getElementById('chatInput');
    const message = input.value.trim();

    if (!message) return;

    // Add user message
    addMessage(message, 'user');
    input.value = '';

    // Send to server or simulate
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({
            type: 'chat',
            message: message,
            agent: document.getElementById('agentSelect').value
        }));
    } else {
        // Simulate AI response
        simulateChatResponse(message);
    }
}

function addMessage(content, role) {
    const container = document.getElementById('chatMessages');
    const avatar = role === 'user' ? '👤' : '🤖';

    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${role}`;
    messageDiv.innerHTML = `
        <div class="message-avatar">${avatar}</div>
        <div class="message-content">${content}</div>
    `;

    container.appendChild(messageDiv);
    container.scrollTop = container.scrollHeight;
}

async function simulateChatResponse(userMessage) {
    await sleep(1000);

    const agent = document.getElementById('agentSelect').value;
    const agentNames = {
        pm: 'Product Manager',
        architect: 'Architect',
        developer: 'Developer',
        qa: 'QA Engineer',
        devops: 'DevOps',
        security: 'Security',
        ux_designer: 'UX Designer'
    };

    const responses = [
        `As a ${agentNames[agent]}, I think that's a great approach!`,
        "Let me help you break that down into actionable steps.",
        "That's an interesting challenge. Here's how I'd approach it...",
        "Good question! The key considerations are quality, scalability, and maintainability.",
        "I'd recommend starting with a clear specification before diving into implementation."
    ];

    const response = responses[Math.floor(Math.random() * responses.length)];
    addMessage(response, 'assistant');
}

function switchAgent() {
    const agent = document.getElementById('agentSelect').value;
    const agentNames = {
        pm: 'Product Manager',
        architect: 'System Architect',
        developer: 'Developer',
        qa: 'QA Engineer',
        devops: 'DevOps',
        security: 'Security Analyst',
        ux_designer: 'UX Designer'
    };

    addMessage(`Switched to ${agentNames[agent]}. How can I help?`, 'assistant');
}

function quickAction(action) {
    const messages = {
        suggest: 'What should I do next with my project?',
        clarify: 'Can you help clarify the requirements?',
        example: 'Can you show me an example?'
    };

    document.getElementById('chatInput').value = messages[action];
    sendMessage();
}

// Console logging
function log(message, level = 'info') {
    const console = document.getElementById('consoleContent');
    const logDiv = document.createElement('div');
    logDiv.className = `log ${level}`;
    logDiv.textContent = `> ${message}`;
    console.appendChild(logDiv);
    console.scrollTop = console.scrollHeight;
}

function clearOutput() {
    document.getElementById('consoleContent').innerHTML =
        '<div class="log info">Console cleared.</div>';
}

// Help modal
function showHelp() {
    document.getElementById('helpModal').classList.remove('hidden');
}

function closeHelp() {
    document.getElementById('helpModal').classList.add('hidden');
}

// Theme toggle
document.getElementById('themeToggle')?.addEventListener('click', () => {
    document.body.classList.toggle('light-theme');
    const btn = document.getElementById('themeToggle');
    btn.textContent = document.body.classList.contains('light-theme') ? '🌙' : '☀️';
});

// Keyboard shortcuts
function setupKeyboardShortcuts() {
    document.addEventListener('keydown', (e) => {
        // Enter to send chat message
        if (e.key === 'Enter' && !e.shiftKey && document.activeElement.id === 'chatInput') {
            e.preventDefault();
            sendMessage();
        }

        // Escape to close modals
        if (e.key === 'Escape') {
            closeHelp();
        }

        // Ctrl+Enter to run action
        if (e.ctrlKey && e.key === 'Enter') {
            const btn = document.querySelector('.btn-primary:not([disabled])');
            if (btn) btn.click();
        }
    });
}

// Load saved project state
function loadProjectState() {
    const saved = localStorage.getItem('viki_project');
    if (saved) {
        try {
            projectData = JSON.parse(saved);
            // Restore UI state
        } catch (e) {
            // Ignore parse errors
        }
    }
}

// Save project state
function saveProjectState() {
    localStorage.setItem('viki_project', JSON.stringify(projectData));
}

// Utility
function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

// Update status from server
function updateStatus(data) {
    if (data.phases) {
        data.phases.forEach(phase => {
            if (phase.complete) {
                markPhaseComplete(phase.name);
            }
        });
    }
}

// Export for debugging
window.viki = {
    selectPhase,
    runAction,
    sendMessage,
    log,
    projectData
};
