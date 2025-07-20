// Подключение к веб-сокету
const socket = new WebSocket(`ws://${window.location.host}/chat`);

// Функция добавления сообщения в чат
function addMessageToChat(message, isNew = false) {
    const displayContainer = document.getElementById('displayContainer');

    // Если сообщение новое, добавляем в начало, иначе - как обычно
    if (isNew) {
        const messageDiv = document.createElement('div');
        messageDiv.classList.add('content-block');
        messageDiv.dataset.id = message.ID;

        const messageText = document.createElement('p');
        messageText.classList.add('display-text');
        messageText.textContent = message.Message;

        const deleteBtn = document.createElement('button');
        deleteBtn.classList.add('delete-btn');
        deleteBtn.textContent = '✖';
        deleteBtn.onclick = function() {
            deleteMessage(message.ID, this);
        };

        messageDiv.appendChild(messageText);
        messageDiv.appendChild(deleteBtn);

        // Добавляем новое сообщение в начало
        displayContainer.insertBefore(messageDiv, displayContainer.firstChild);
    }
}

// Обработчик входящих сообщений через веб-сокет
socket.onmessage = function(event) {
    const message = JSON.parse(event.data);
    addMessageToChat(message, true);
};

// Функции для работы с медиа
let selectedMedia = null;
let selectedMediaType = null;

function previewMedia(input) {
    const file = input.files[0];
    if (!file) return;
    const reader = new FileReader();
    reader.onload = function (e) {
        if (file.type.startsWith("image/")) {
            selectedMedia = e.target.result;
            selectedMediaType = "image";
            document.getElementById("photoPreview").innerHTML = `<img src="${selectedMedia}" alt="Preview" style="max-width:120px;max-height:90px;">`;
        } else if (file.type.startsWith("video/")) {
            selectedMedia = e.target.result;
            selectedMediaType = "video";
            document.getElementById("photoPreview").innerHTML = `<video src="${selectedMedia}" controls style="max-width:120px;max-height:90px;"></video>`;
        } else {
            selectedMedia = null;
            selectedMediaType = null;
            document.getElementById("photoPreview").innerHTML = '';
        }
    };
    reader.readAsDataURL(file);
}

// Обработчик отправки формы
document.addEventListener('DOMContentLoaded', function() {
    const chatForm = document.getElementById('chatForm');
    if (chatForm) {
        chatForm.addEventListener('submit', function(e) {
            e.preventDefault();
            const textInput = document.getElementById('textInput');
            const message = textInput.value.trim();

            if (!message && !selectedMedia) return;

            // Отправка сообщения на сервер
            fetch('/chat', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: 'message=' + encodeURIComponent(message)
            }).then(response => {
                if (response.ok) {
                    textInput.value = '';

                    // Очистка превью медиа
                    if (selectedMedia) {
                        selectedMedia = null;
                        selectedMediaType = null;
                        document.getElementById("photoPreview").innerHTML = '';
                        document.getElementById("fileInput").value = '';
                    }
                }
            });
        });
    }
});

// Функция удаления сообщения
function deleteMessage(id, btn) {
    fetch('/chat?id=' + id, { method: 'DELETE' })
        .then(res => {
            if (res.ok) {
                btn.closest('.content-block').remove();
            }
        });
}

// Обработка закрытия соединения
socket.onclose = function(event) {
    if (event.wasClean) {
        console.log(`Соединение закрыто чисто, код=${event.code} причина=${event.reason}`);
    } else {
        console.log('Соединение прервано');
        // Попытка переподключения через 5 секунд
        setTimeout(() => {
            window.location.reload();
        }, 5000);
    }
};

// Обработка ошибок соединения
socket.onerror = function(error) {
    console.log('Ошибка соединения:', error.message);
};