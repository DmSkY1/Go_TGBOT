# Delta Bot

<p align="center">
  <img src="https://media1.tenor.com/m/ZQ06vsxdy1cAAAAd/bear-scream.gif" alt="Delta Bot" width="500"/>
</p>

<p align="center">
  Telegram-бот для удаления фона с изображений и улучшения их качества!
</p>

---

## 📸 Пример
Пример удаления фона с изображения:

<p align="center">
  <img src="https://i.postimg.cc/90h1y6ss/example.gif" alt="Delta Bot" width="320"/>
</p>


---

## 🔗 Попробовать бота
Отсканируйте QR-код или нажмите кнопку ниже, чтобы начать использовать Delta Bot в Telegram!

<p align="center">
  <img src="https://i.postimg.cc/QxDn773V/qr.png" alt="QR-код" width="300"/>
</p>

<p align="center">
  <a href="https://t.me/delta_ph_bot" target="_blank">
    <img src="https://img.shields.io/badge/🎉_Попробовать_бота-0088cc?style=for-the-badge&logo=telegram&logoColor=white" alt="Попробовать Delta Bot">
  </a>
</p>

---

## 🛠 Установка и настройка
Следуйте этим шагам, чтобы настроить Delta Bot:

1. **Создание бота**:
   - Перейдите в [BotFather](https://t.me/BotFather) в Telegram.
   - Создайте нового бота и получите **токен бота**.
   - Разрешите боту добавляться в группы.

2. **Клонирование репозитория**:
   - Скачайте репозиторий проекта на свой компьютер.

3. **Установка зависимостей**:
   - Установите необходимые библиотеки Go:
     ```bash
     go get github.com/go-telegram-bot-api/telegram-bot-api
     go get github.com/andybalholm/brotli 
     go get github.com/joho/godotenv 
     go get github.com/klauspost/compress 
     go get github.com/valyala/bytebufferpool 
     go get github.com/xlab/multipartstreamer 
     ```
4. Настройка переменных окружения:
   - Создайте файл .env в корне проекта.
   - Добавьте токен бота:
   ```bash
    TOKEN=your_bot_token_here
   ```
5. Получение API-ключей Picsart:
   - Перейдите на сайт [Picsart API](https://docs.picsart.io/docs/creative-apis-get-api-key) и получите API-ключи.
   - Создайте файл ApiKey.txt в корне проекта и добавьте в него полученные ключи.
6. Генерация JSON-файла:
   - Перейдите в папку generate_json.
   - Запустите скрипт для создания JSON-файла:
   ```bash
    go run json_generate_script.go
   ```
7. Запуск бота:
   - Запустите бота командой:
    ```bash
    go run main.go
   ```

🎉 Ваш Delta Bot готов к работе!

---

## 📝 Примечания
- Убедитесь, что файл ApiKey.txt назван точно так (с учетом регистра).
- Храните API-ключи и токен бота в безопасности, не публикуйте их.
- При возникновении проблем проверьте вывод в консоли или обратитесь к документации Picsart.

---

## 📧 Контакты
По каким либо вопросам писать мне в [Telegram](https://t.me/supchik_mmm)