package Commands

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	mediaGroups = make(map[string][]tgbotapi.Message)
)

func Change_background(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Println("OK")
	if update.Message.MediaGroupID != "" {
		mediaGroupID := update.Message.MediaGroupID

		// Добавляем сообщение в медиагруппу
		mediaGroups[mediaGroupID] = append(mediaGroups[mediaGroupID], *update.Message)

		// Если медиагруппа собрана (2 фото)
		if len(mediaGroups[mediaGroupID]) == 2 {
			log.Printf("Медиагруппа %s собрана:\n", mediaGroupID)

			// Создаем срез для хранения ссылок на фотографии
			var photoLinks []string

			// Обрабатываем каждое сообщение в медиагруппе
			for i, msg := range mediaGroups[mediaGroupID] {
				if len(msg.Photo) > 0 {
					// Обработка фотографий
					photo := msg.Photo[len(msg.Photo)-1] // Берем фото с самым высоким разрешением
					log.Printf("Фото %d: FileID = %s\n", i+1, photo.FileID)

					// Получаем информацию о файле
					file, err := bot.GetFile(tgbotapi.FileConfig{FileID: photo.FileID})
					if err != nil {
						log.Println("Ошибка при получении файла:", err)
						continue
					}

					// Получаем ссылку на файл
					fileURL := file.Link(bot.Token)
					log.Printf("Ссылка на фото %d: %s\n", i+1, fileURL)

					// Сохраняем ссылку в срез
					photoLinks = append(photoLinks, fileURL)
				}
			}

			// Выводим все ссылки на фотографии
			log.Println("Ссылки на фотографии:", photoLinks)

			// Удаляем медиагруппу из хранилища
			delete(mediaGroups, mediaGroupID)
		}
	} else {
		// Обработка одиночных сообщений
		log.Println("Одиночное сообщение:", update.Message.Text)
	}

}
