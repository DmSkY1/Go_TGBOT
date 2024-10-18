package Comand

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	installPhoto "main.go/INSTALL_PICTURE"
	post_file "main.go/POST"
	rand_key "main.go/RandomKey"
)

const (
	IdleState = iota
	WaitingForImageToRemoveBackground
)

func Remove_background_image(bot *tgbotapi.BotAPI, update tgbotapi.Update, user_state *map[int64]int, user_active_command *map[int64]string) {
	chatID := update.Message.Chat.ID
	state := (*user_state)[chatID]
	messageID := update.Message.MessageID

	if update.Message.Photo != nil {
		photo := *update.Message.Photo
		fileID := photo[len(photo)-1].FileID
		file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
		if err != nil {
			log.Println(err)
		}

		photoMsg := tgbotapi.NewAnimationShare(chatID, "https://i.gifer.com/origin/38/3823ad20629c89b3dd4821b80eee79eb_w200.gif")
		photoMsg.ReplyToMessageID = messageID
		photoMsg.Caption = fmt.Sprintf("📸 Ваше фото в обработке! 🚀\n" +
			"Я занимаюсь улучшением и увеличением вашего изображения. Это займет примерно 10 секунд. ⏳✨\n\n" +
			"Пожалуйста, подождите немного, и ваше фото будет готово к просмотру. Спасибо за терпение! 😊")
		bot.Send(photoMsg)

		// Получаем прямую ссылку на файл фотографии
		fileURL := file.Link(bot.Token)
		installPhoto.InstallPhoto(fmt.Sprintf("picture/%s.jpeg", fileID), fileURL) // устанавливаем фотографию в папку picture

		api_key, err := rand_key.GetRandomAPIKey("NewApiKey.txt")
		if err != nil {
			log.Println("Ошибка при получении API-ключа:", err)
			return
		}

		result, err := post_file.RemoveBackground(api_key, fmt.Sprintf("picture/%s.jpeg", fileID))
		if err != nil {
			log.Println("Ошибка при удалении фона:", err)
			return
		}
		err = os.Remove(fmt.Sprintf("picture/%s.jpeg", fileID))
		if err != nil {
			fmt.Println("Error deleting file 101!!:", err)
		}
		res, err := post_file.DownloadFileUrl(result, fileID, "1")
		if err != nil {
			log.Println("Ошибка при загрузке файла:", err)
			return
		}
		document := tgbotapi.NewDocumentUpload(chatID, res)
		document.Caption = fmt.Sprintf("🎯 *Готово!* 🎯\n\n" +
			"🎉 Ваше изображение теперь выглядит лучше! 📸🌈\n\n" +
			"🔍 _Посмотрите внимательно и наслаждайтесь результатом!_ 😊")
		_, err = bot.Send(document)
		if err != nil {
			log.Println("Ошибка отправки документа:", err)
		}

		go func() {
			deleteMsg := tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID+1)
			if _, err := bot.DeleteMessage(deleteMsg); err != nil {
				log.Println("Ошибка при удалении сообщения:", err)
			}
		}()
		err = os.Remove(fmt.Sprintf("picture/%s_1.jpeg", fileID))
		if err != nil {
			fmt.Println("Error deleting file: 92!!", err)
		}
		mu.Lock()
		(*user_state)[chatID] = IdleState
		mu.Unlock()

	} else if update.Message.Document != nil && isImageFile(update.Message.Document) {
		photoMsg := tgbotapi.NewAnimationShare(chatID, "https://i.gifer.com/origin/38/3823ad20629c89b3dd4821b80eee79eb_w200.gif")
		photoMsg.ReplyToMessageID = messageID
		photoMsg.Caption = fmt.Sprintf("📸 Ваше фото в обработке! 🚀\n" +
			"Я занимаюсь улучшением и увеличением вашего изображения. Это займет примерно 10 секунд. ⏳✨\n\n" +
			"Пожалуйста, подождите немного, и ваше фото будет готово к просмотру. Спасибо за терпение! 😊")
		bot.Send(photoMsg)

		fileID := update.Message.Document.FileID
		file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
		if err != nil {
			log.Println("Ошибка при получении файла:", err)
			return
		}

		download_url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)
		filepath := fmt.Sprintf("picture/%s_%s", fileID, update.Message.Document.FileName)

		err = installPhoto.InstallPhoto(filepath, download_url)
		if err != nil {
			log.Println("Ошибка при скачивании файла:", err)
		} else {
			log.Printf("Файл успешно скачан с именем: %s", filepath)
		}

		api_key, err := rand_key.GetRandomAPIKey("NewApiKey.txt")
		if err != nil {
			log.Println("Ошибка при получении API-ключа:", err)
			return
		}
		result, err := post_file.RemoveBackground(api_key, filepath)
		if err != nil {
			log.Println("Ошибка при отправке пост Запроса")
			return
		}
		err = os.Remove(filepath)
		if err != nil {
			fmt.Println("Error deleting file 101!!:", err)
		}

		res, err := post_file.DownloadFile(result, fileID, update.Message.Document.FileName)
		if err != nil {
			log.Println("Ошибка при загрузке файла:", err)
		}
		document := tgbotapi.NewDocumentUpload(chatID, res)
		document.Caption = fmt.Sprintf("🎯 *Готово!* 🎯\n\n" +
			"🎉 Ваше изображение теперь выглядит лучше! 📸🌈\n\n" +
			"🔍 _Посмотрите внимательно и наслаждайтесь результатом!_ 😊")
		_, err = bot.Send(document)
		if err != nil {
			log.Println("Ошибка отправки документа:", err)
		}
		go func() {
			deleteMsg := tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID+1)
			if _, err := bot.DeleteMessage(deleteMsg); err != nil {
				log.Println("Ошибка при удалении сообщения:", err)
			}
		}()
		err = os.Remove(filepath[:strings.Index(filepath, ".")+1] + "jpeg")
		if err != nil {
			fmt.Println("Error deleting file: 92!", err)
		}
		mu.Lock()
		(*user_state)[chatID] = IdleState
		mu.Unlock()

	} else if update.Message != nil && update.Message.Text != "" {

	} else {
		if state == WaitingForImageToRemoveBackground {
			msg := tgbotapi.NewMessage(chatID, "Пожалуйста, подождите немного, и ваше фото будет готово к просмотру. Спасибо за терпение!")
			bot.Send(msg)
			return
		} else {
			errorMessage_url := tgbotapi.NewAnimationShare(chatID, "https://gifs.obs.ru-moscow-1.hc.sbercloud.ru/1d06b49de1ac9de5cbc468d1d449d74658d39b7471c689cf5ec7570106908a9e.gif")
			errorMessage_url.Caption = fmt.Sprintf("🚫 <i><strong>Ошибка! Неверный формат запроса.</strong></i> Пожалуйста, отправьте фотографию URL-адрес на нее. 🌐\n\n"+
				"Чтобы продолжить, повторите команду /%s с <strong>корректным URL.</strong>⚠️", (*user_active_command)[chatID])
			errorMessage_url.ParseMode = "HTML"
			errorMessage_url.ReplyToMessageID = messageID // Указываем ID сообщения для ответа
			_, err := bot.Send(errorMessage_url)
			if err != nil {
				log.Println("Ошибка отправки GIF:", err)
			}
			mu.Lock()
			state = IdleState
			mu.Unlock()
			return
		}

	}
}
