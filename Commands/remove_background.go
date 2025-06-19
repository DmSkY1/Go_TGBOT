package Commands

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	installPhoto "main.go/INSTALL_PICTURE"
	post_file "main.go/POST"
	rand_key "main.go/RandomKey"
)

const (
	IdleState = iota
	WaitingForImageToRemoveBackground
	WaitingForProcessingRemoveBg
)

func Remove_background_image(bot *tgbotapi.BotAPI, update tgbotapi.Update, user_state *map[int64]int, user_active_command *map[int64]string) {
	chatID := update.Message.Chat.ID
	messageID := update.Message.MessageID

	if update.Message.Photo != nil {
		log.Printf("\033[32m[INFO]\033[0m Получена фотография от пользователя. ChatID [%d]", chatID)
		setUserState(chatID, user_state, WaitingForProcessingRemoveBg)
		photo := update.Message.Photo
		fileID := photo[len(photo)-1].FileID
		file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
		if err != nil {
			log.Println("\033[31m[Error]\033[0m Ошибка при получении id фотографии:", err)
			return
		}
		log.Printf("\033[32m[INFO]\033[0m ID фотографии получено. ChatID [%d]", chatID)

		infoMessage(bot, chatID, messageID)

		// Получаем прямую ссылку на файл фотографии
		fileURL := file.Link(bot.Token)
		if err := installPhoto.InstallPhoto(fmt.Sprintf("picture/%s.jpeg", fileID), fileURL); err != nil {
			log.Println(err)
			return
		} // устанавливаем фотографию в папку picture
		log.Printf("\033[32m[INFO]\033[0m Фотография успешно скачена. ChatID [%d]", chatID)

		api_key, err := rand_key.GetRandomAPIKey()
		if err != nil {
			log.Println("\033[31m[Error]\033[0m Ошибка при получении API-ключа:", err)
			return
		}
		log.Printf("\033[32m[INFO]\033[0m API ключ сгенерирован. ChatID [%d]", chatID)

		result, err := post_file.RemoveBackground(api_key, fmt.Sprintf("picture/%s.jpeg", fileID))
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("\033[32m[INFO]\033[0m Документ успешно обработан. ChatID [%d]", chatID)

		if err = os.Remove(fmt.Sprintf("picture/%s.jpeg", fileID)); err != nil {
			log.Println("\033[31m[Error]\033[0m Ошибка при удалении фотографии:", err)
			return
		}

		document := tgbotapi.NewDocument(chatID, tgbotapi.FileURL(result))
		document.Caption = fmt.Sprintf("🎯 *Готово!* 🎯\n\n" +
			"🎉 Ваше изображение теперь выглядит лучше! 📸🌈\n\n" +
			"🔍 _Посмотрите внимательно и наслаждайтесь результатом!_ 😊")
		_, err = bot.Send(document)
		if err != nil {
			log.Println("\033[31m[Error]\033[0m Ошибка отправки документа:", err)
			return
		}
		log.Printf("\033[32m[INFO]\033[0m Документ успешно отправлен. ChatID [%d]", chatID)

		go func() {
			deleteMsg := tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID+1)
			if _, err := bot.Request(deleteMsg); err != nil {
				log.Println("\033[31m[Error]\033[0m Ошибка при удалении сообщения:", err)
				return
			}
			log.Printf("\033[32m[INFO]\033[0m Информационное сообщение удалено успешно. ChatID [%d]", chatID)
		}()
		setUserState(chatID, user_state, IdleState)

	} else if update.Message.Document != nil && isImageFile(update.Message.Document) {
		log.Printf("\033[32m[INFO]\033[0m Получена фотография от пользователя. ChatID [%d]", chatID)
		setUserState(chatID, user_state, WaitingForProcessingRemoveBg)
		infoMessage(bot, chatID, messageID)

		fileID := update.Message.Document.FileID
		file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
		if err != nil {
			log.Println("\033[31m[Error]\033[0m Ошибка при получении id фотографии:", err)
			return
		}
		log.Printf("\033[32m[INFO]\033[0m ID фотографии получено. ChatID [%d]", chatID)

		download_url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)
		filepath := fmt.Sprintf("picture/%s_%s", fileID, update.Message.Document.FileName)

		if err = installPhoto.InstallPhoto(filepath, download_url); err != nil {
			log.Println(err)
			return
		}
		log.Printf("\033[32m[INFO]\033[0m Фотография успешно скачена. ChatID [%d]", chatID)

		api_key, err := rand_key.GetRandomAPIKey()
		if err != nil {
			log.Println("\033[31m[Error]\033[0m Ошибка при получении API-ключа:", err)
			return
		}
		log.Printf("\033[32m[INFO]\033[0m API ключ сгенерирован. ChatID [%d]", chatID)

		result, err := post_file.RemoveBackground(api_key, filepath)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("\033[32m[INFO]\033[0m Документ успешно обработан. ChatID [%d]", chatID)

		if err = os.Remove(filepath); err != nil {
			log.Println("\033[31m[Error]\033[0m Ошибка при удалении фотографии:", err)
			return
		}

		document := tgbotapi.NewDocument(chatID, tgbotapi.FileURL(result))
		document.Caption = fmt.Sprintf("🎯 *Готово!* 🎯\n\n" +
			"🎉 Ваше изображение теперь выглядит лучше! 📸🌈\n\n" +
			"🔍 _Посмотрите внимательно и наслаждайтесь результатом!_ 😊")
		_, err = bot.Send(document)
		if err != nil {
			log.Println("\033[31m[Error]\033[0m Ошибка отправки документа:", err)
			return
		}
		log.Printf("\033[32m[INFO]\033[0m Документ успешно отправлен. ChatID [%d]", chatID)
		go func() {
			deleteMsg := tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID+1)
			_, err := bot.Request(deleteMsg)
			if err != nil {
				log.Println("\033[31m[Error]\033[0m Ошибка при удалении сообщения:", err)
				return
			}
			log.Printf("\033[32m[INFO]\033[0m Информационное сообщение удалено успешно. ChatID [%d]", chatID)
		}()
		setUserState(chatID, user_state, IdleState)

	} else if update.Message != nil && update.Message.Text != "" {
		msg := tgbotapi.NewMessage(chatID, "Пожалуйста, отправьте фотографию. Бот не воспринимает ваш текст. Я буду ждать ващей фотографии ))")
		bot.Send(msg)
		return
	} else {
		errorMessage_url := tgbotapi.NewAnimation(chatID, tgbotapi.FileURL("https://gifs.obs.ru-moscow-1.hc.sbercloud.ru/1d06b49de1ac9de5cbc468d1d449d74658d39b7471c689cf5ec7570106908a9e.gif"))
		errorMessage_url.Caption = fmt.Sprintf("🚫 <i><strong>Ошибка! Неверный формат запроса.</strong></i> Пожалуйста, отправьте фотографию URL-адрес на нее. 🌐\n\n"+
			"Чтобы продолжить, повторите команду /%s с <strong>корректным URL.</strong>⚠️", (*user_active_command)[chatID])
		errorMessage_url.ParseMode = "HTML"
		errorMessage_url.ReplyToMessageID = messageID // Указываем ID сообщения для ответа
		_, err := bot.Send(errorMessage_url)
		if err != nil {
			log.Println("\033[31m[Error]\033[0m Ошибка отправки GIF:", err)
			return
		}
		return

	}
}

func infoMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	photoMsg := tgbotapi.NewAnimation(chatID, tgbotapi.FileURL("https://i.gifer.com/origin/38/3823ad20629c89b3dd4821b80eee79eb_w200.gif"))
	photoMsg.ReplyToMessageID = messageID
	photoMsg.Caption = fmt.Sprintf("📸 Ваше фото в обработке! 🚀\n" +
		"Я занимаюсь улучшением и увеличением вашего изображения. Это займет примерно 10 секунд. ⏳✨\n\n" +
		"Пожалуйста, подождите немного, и ваше фото будет готово к просмотру. Спасибо за терпение! 😊")
	_, err := bot.Send(photoMsg)
	if err != nil {
		log.Println("\033[31m[Error]\033[0m Ошибка при отправке сообщения", err)
		return
	}
}
