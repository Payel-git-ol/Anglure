// Функции для работы с медиа и сообщениями
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

function sendMessage() {
    const textInput = document.getElementById('textInput');
    const displayContainer = document.getElementById('displayContainer');
    if (!textInput.value.trim() && !selectedMedia) return;

    const newDiv = document.createElement('div');
    newDiv.classList.add('content-block');

    // Медиа
    if (selectedMedia && selectedMediaType === "image") {
        const img = document.createElement('img');
        img.src = selectedMedia;
        img.style.maxWidth = "150px";
        img.style.display = "block";
        newDiv.appendChild(img);
    }
    if (selectedMedia && selectedMediaType === "video") {
        const video = document.createElement('video');
        video.src = selectedMedia;
        video.controls = true;
        video.style.maxWidth = "150px";
        video.style.display = "block";
        newDiv.appendChild(video);
    }

    // Текст
    if (textInput.value.trim()) {
        const p = document.createElement('p');
        p.className = "display-text";
        p.textContent = textInput.value;

        // Кнопка удаления
        const delBtn = document.createElement('button');
        delBtn.textContent = "✖";
        delBtn.className = "delete-btn";
        delBtn.onclick = function(e) {
            e.stopPropagation();
            newDiv.remove();
        };

        // Обёртка для текста и кнопки
        const msgWrap = document.createElement('span');
        msgWrap.style.display = "flex";
        msgWrap.style.alignItems = "center";
        msgWrap.appendChild(p);
        msgWrap.appendChild(delBtn);

        newDiv.appendChild(msgWrap);
    }

    displayContainer.appendChild(newDiv);

    // Очистка
    selectedMedia = null;
    selectedMediaType = null;
    document.getElementById("photoPreview").innerHTML = '';
    document.getElementById("fileInput").value = '';
}

// Enter в поле ввода
document.getElementById('textInput').addEventListener('keydown', function(e) {
    if (e.key === 'Enter' && !e.shiftKey) {
        // Обработка будет в обработчике формы
    }
});

// Обработчик отправки формы
document.addEventListener('DOMContentLoaded', function() {
    const chatForm = document.getElementById('chatForm');
    if (chatForm) {
        chatForm.addEventListener('submit', function(e) {
            e.preventDefault();
            const textInput = document.getElementById('textInput');
            const message = textInput.value.trim();
            if (!message) return;

            fetch('/chat', {
                method: 'POST',
                headers: {'Content-Type': 'application/x-www-form-urlencoded'},
                body: 'message=' + encodeURIComponent(message)
            }).then(() => {
                sendMessage(); // Визуальный вывод сообщения
                textInput.value = ''; // Очистка поля ввода
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



function submitForm(e) {
    e.preventDefault(); // Предотвращаем стандартную отправку формы

    // Отправляем форму через Fetch API
    fetch('/chat', {
        method: 'POST',
        body: new FormData(document.getElementById('chatForm'))
    })
        .then(response => {
            // После успешной отправки обновляем страницу
            location.reload();
        })
        .catch(error => {
            console.error('Ошибка:', error);
        });
}
setInterval(() => location.reload(), 10000);