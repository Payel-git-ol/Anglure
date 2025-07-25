document.addEventListener('DOMContentLoaded', () => {
    const socket = new WebSocket(`ws://${window.location.host}/chat`);
    const processedMessageIds = new Set();

    const displayContainer = document.getElementById('displayContainer');
    const textInput = document.getElementById('textInput');
    const fileInput = document.getElementById('fileInput');
    const photoPreview = document.getElementById('photoPreview');
    const chatForm = document.getElementById('chatForm');

    socket.addEventListener('open', () => console.log('📡 WS connected'));

    socket.addEventListener('message', event => {
        try {
            const message = JSON.parse(event.data);

            // Приводим id к строке для безопасности
            const idStr = message.id !== undefined && message.id !== null ? String(message.id) : '';

            if (!idStr) {
                console.warn('Получено сообщение без id — игнорируем');
                return;
            }

            if (processedMessageIds.has(idStr)) {
                // Уже обработано
                return;
            }

            // Ищем временное сообщение с таким же текстом и именем, чтобы заменить его
            const tempMsgDivs = displayContainer.querySelectorAll('[data-id^="temp-"]');
            let replacedTemp = false;
            tempMsgDivs.forEach(div => {
                const nameDiv = div.querySelector('.nameUserTextMessage');
                const msgP = div.querySelector('.display-text');

                if (
                    nameDiv?.textContent === (message.name || 'Неизвестный') &&
                    msgP?.textContent === (message.message || '')
                ) {
                    // Заменяем временное сообщение на сообщение с сервера
                    div.dataset.id = idStr;

                    // Показываем кнопку удаления
                    const delBtn = div.querySelector('.delete-btn');
                    if (delBtn) {
                        delBtn.style.display = 'inline-block';
                        delBtn.onclick = () => deleteMessage(idStr, delBtn);
                    }

                    replacedTemp = true;
                    processedMessageIds.add(idStr);
                }
            });

            if (!replacedTemp) {
                // Если не нашли временное, добавляем новое сообщение
                addMessageToChat(message, true);
                processedMessageIds.add(idStr);
            }

        } catch (err) {
            console.error('❌ Ошибка парсинга сообщения:', err);
        }
    });

    socket.addEventListener('close', () => console.log('📡 WS disconnected'));
    socket.addEventListener('error', err => console.error('❌ WS ошибка:', err));

    chatForm.addEventListener('submit', async e => {
        e.preventDefault();

        const messageText = textInput.value.trim();
        const file = fileInput.files[0];

        if (!messageText && !file) return;

        const formData = new FormData();
        formData.append('message', messageText);
        if (file) formData.append('media', file);

        // Добавляем временное сообщение с уникальным temp id
        const tempId = 'temp-' + Date.now();

        addMessageToChat({
            id: tempId,
            name: getCookie('user_name') || 'Вы',
            message: messageText
        });

        try {
            await fetch('/chat', {
                method: 'POST',
                body: formData
            });

            textInput.value = '';
            fileInput.value = '';
            photoPreview.innerHTML = '';
        } catch (err) {
            console.error('❌ Ошибка отправки:', err);
        }
    });

    function addMessageToChat(message, isFromServer = false) {
        if (!displayContainer) return;

        const idStr = message.id !== undefined && message.id !== null ? String(message.id) : '';

        if (idStr && displayContainer.querySelector(`[data-id="${idStr}"]`)) return;

        const messageDiv = document.createElement('div');
        messageDiv.classList.add('content-block');
        if (idStr) messageDiv.dataset.id = idStr;

        const nameDiv = document.createElement('div');
        nameDiv.classList.add('nameUserTextMessage');
        nameDiv.textContent = message.name || 'Неизвестный';

        const messageText = document.createElement('p');
        messageText.classList.add('display-text');
        messageText.textContent = message.message || '';

        const deleteBtn = document.createElement('button');
        deleteBtn.classList.add('delete-btn');
        deleteBtn.textContent = '✖';

        if (idStr && isFromServer && !idStr.startsWith('temp-')) {
            deleteBtn.onclick = () => deleteMessage(idStr, deleteBtn);
        } else {
            deleteBtn.style.display = 'none';
        }

        messageDiv.appendChild(nameDiv);
        messageDiv.appendChild(messageText);
        messageDiv.appendChild(deleteBtn);
        displayContainer.insertBefore(messageDiv, displayContainer.firstChild);
    }

    window.deleteMessage = async function(id, buttonElement) {
        if (!confirm('Удалить сообщение?')) return;

        try {
            const response = await fetch(`/chat?id=${encodeURIComponent(id)}`, {
                method: 'DELETE'
            });

            if (response.ok) {
                const messageDiv = buttonElement.closest('.content-block');
                if (messageDiv) messageDiv.remove();
                processedMessageIds.delete(String(id));
            } else {
                console.error('❌ Ошибка удаления:', response.statusText);
            }
        } catch (err) {
            console.error('❌ Ошибка при удалении:', err);
        }
    };

    window.previewMedia = function(input) {
        photoPreview.innerHTML = '';
        const file = input.files[0];
        if (!file) return;

        const reader = new FileReader();
        reader.onload = () => {
            const media = document.createElement(file.type.startsWith('image/') ? 'img' : 'video');
            media.src = reader.result;
            media.style.maxWidth = '200px';
            media.controls = true;
            photoPreview.appendChild(media);
        };
        reader.readAsDataURL(file);
    };

    function getCookie(name) {
        const value = `; ${document.cookie}`;
        const parts = value.split(`; ${name}=`);
        return parts.length === 2 ? decodeURIComponent(parts.pop().split(';')[0]) : '';
    }

    const profileBtn = document.querySelector('.profile-button');
    const dropdown = document.querySelector('.dropdown-content');
    const profileModal = document.getElementById('profileModal');
    const closeModalBtn = profileModal?.querySelector('.close-modal');
    const profileMenuBtn = document.getElementById('profileButton');

    if (profileBtn && dropdown && profileModal && closeModalBtn && profileMenuBtn) {
        profileBtn.addEventListener('click', () => {
            dropdown.style.display = dropdown.style.display === 'block' ? 'none' : 'block';
        });

        window.addEventListener('click', (e) => {
            if (!profileBtn.contains(e.target) && !dropdown.contains(e.target)) {
                dropdown.style.display = 'none';
            }
        });

        profileMenuBtn.addEventListener('click', () => {
            dropdown.style.display = 'none';
            profileModal.classList.add('show');
        });

        closeModalBtn.addEventListener('click', () => {
            profileModal.classList.remove('show');
        });

        window.addEventListener('click', (e) => {
            if (e.target === profileModal) {
                profileModal.classList.remove('show');
            }
        });
    }
});
