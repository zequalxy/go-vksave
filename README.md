# go-vksave
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

Простое приложение, которое позволяет скачивать изображения из чата или диалога ВК.

---

### Использование

1. Перейти по [адресу](https://oauth.vk.com/authorize?client_id=2685278&display=popup&redirect_uri=https://oauth.vk.com/blank.html&scope=messages,offline&response_type=token&v=5.131&state=123456) для получения ссылки с токеном.
2. Скопировать полученную ссылку из адресной строки браузера.
3. Скопировать ссылку на чат или диалог ВК, откуда необходимо загрузить изображения.
4. Запустить программу, использую аргументы `-t` для ссылки с токеном, `-c` для ссылки с чатом.
5. Запуск
```
make build
make run
```
