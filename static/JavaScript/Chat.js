// Подключение к веб-сокету
const socket = new WebSocket(`ws://${window.location.host}/chat`);

// Хранилище ID обработанных сообщений
const processedMessageIds = new Set();

// Функция добавления сообщения в чат
function addMessageToChat(message, isFromServer = false) {
    const displayContainer = document.getElementById('displayContainer');
    if (!displayContainer) return;

    // Для серверных сообщений проверяем дубликаты
    if (isFromServer && message.id && processedMessageIds.has(message.id)) {
        return; // Пропускаем уже обработанные
    }

    const messageDiv = document.createElement('div');
    messageDiv.classList.add('content-block');
    if (message.id) {
        messageDiv.dataset.id = message.id;
        processedMessageIds.add(message.id); // Регистрируем ID
    }

    const messageText = typeof message === 'object' ?
        (message.message || message.Message || JSON.stringify(message)) :
        String(message);

    const messageElement = document.createElement('p');
    messageElement.classList.add('display-text');
    messageElement.textContent = messageText;

    const deleteBtn = document.createElement('button');
    deleteBtn.classList.add('delete-btn');
    deleteBtn.textContent = '✖';

    if (message.id) {
        deleteBtn.onclick = function() {
            deleteMessage(message.id, this);
        };
    } else {
        deleteBtn.style.display = 'none';
    }

    messageDiv.appendChild(messageElement);
    messageDiv.appendChild(deleteBtn);

    // Все новые сообщения добавляем в начало
    displayContainer.insertBefore(messageDiv, displayContainer.firstChild);
}

// Обработчик веб-сокета
socket.onmessage = function(event) {
    try {
        const message = JSON.parse(event.data);
        addMessageToChat(message, true); // Только серверные сообщения
    } catch (e) {
        console.error('Ошибка парсинга:', e);
    }
};

// Обработчик отправки формы
document.getElementById('chatForm')?.addEventListener('submit', async function(e) {
    e.preventDefault();
    const textInput = document.getElementById('textInput');
    const message = textInput.value.trim();
    const fileInput = document.getElementById('fileInput');

    if (!message && !fileInput.files[0]) return;

    // Отправка на сервер
    const formData = new FormData();
    formData.append('message', message);
    if (fileInput.files[0]) formData.append('media', fileInput.files[0]);

    try {
        await fetch('/chat', {
            method: 'POST',
            body: formData
        });
        textInput.value = '';
        fileInput.value = '';
        document.getElementById('photoPreview').innerHTML = '';
    } catch (error) {
        console.error('Ошибка отправки:', error);
    }
});