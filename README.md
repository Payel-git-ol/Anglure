# Anglure 🌐

Русский | [English version](#english-version)

---

## 📌 О проекте

Anglure — это социальная сеть, построенная на Go, с акцентом на модульность, масштабируемость и качественную архитектуру.  
Проект создан с расчётом на командную разработку: структура прозрачна, фронтенд окружение подготовлено, а архитектура вдохновлена международными стандартами.

---

## ⚙️ Стек технологий

- Backend: Go, net/http ,Gorilla/WebSocket
- Frontend: HTML/CSS/JS, шаблоны, заготовка под Angular
- Базы данных: PostgresSQL
- DevOps: Docker, Docker Compose
- Архитектура: [ISO/IEC/IEEE 42010](https://en.wikipedia.org/wiki/ISO/IEC_42010) — стандарт описания архитектуры программных систем[43dcd9a7-70db-4a1f-b0ae-981daa162054](https://en.wikipedia.org/wiki/ISO/IEC_42010?citationMarker=43dcd9a7-70db-4a1f-b0ae-981daa162054 "1")

---
## 🚀 Запуск програмы 

- Запустите Docker контейнер "docker run --hostname=0920f0a628e1 --env=POSTGRES_PASSWORD=mysecretpassword -env=PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/lib/postgresql/17/bin --env=GOSU_VERSION=1.17 --env=LANG=en_US.utf8 --env=PG_MAJOR=17 --env=PG_VERSION=17.5-1.pgdg120+1 --env=PGDATA=/var/lib/postgresql/data --volume=/var/lib/postgresql/data --network=bridge -p 5432:5432 --restart=no --runtime=runc -d postgres"
- После запустите уже само приложениие нужно прейти к дериктории куда вы загрузили своё приложение после запустите команду "go run cmd/app/main.go" и перейдите по ссылки которую вам выдало в терминале 
