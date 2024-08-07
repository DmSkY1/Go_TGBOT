package main

import (
	"fmt"
	"log"
	"path/filepath"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func main() {
	// create bot and conenct
	bot, err := tgbotapi.NewBotAPI("6798282567:AAEC8jxADvq9CTSaHBtmxYkflbP7pj72gvU")
	if err != nil {
		log.Fatal()
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 20

	updatesChan, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println("Ошибка при апдейте канала")
	}

	log.Println("Бот запущен")

	for update := range updatesChan {
		if update.Message != nil {
			go AsyncProcess(bot, update)
		}
	}
}

func AsyncProcess(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message.IsCommand() {
		command := update.Message.Command()
		switch command {
		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, я наконец то зработал!!!")
			bot.Send(msg)

		}
	} else {
		switch {
		case update.Message.Photo != nil:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ваше фото принято в работу. Пожалуйста ожидайте!"))
			//Полуение информации о фото
			photo := *update.Message.Photo
			// Выбираем фотографию наивысшего разрешения
			fileID := photo[len(photo)-1].FileID
			// Получаем информацию о файле
			file_name := fmt.Sprintf("%s.jpg", fileID)
			file_path := filepath.Join("picture/", file_name)
			InstallPhoto(fileID, bot)
			api_token, err := random_token()
			if err != nil {
				log.Println(err)
			}
			image_url := Post(file_path, api_token)
			Remove_Image(file_path)
			Install_picture(image_url, file_path)
			send_photo := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, file_path)
			send_photoDocument := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, file_path)
			send_photoDocument.Caption = "Фото обработано, наслаждайтесь работой!!!!"
			_, err = bot.Send(send_photo)
			if err != nil {
				log.Println("Ошибка при отправке фото")
			}
			_, err = bot.Send(send_photoDocument)
			if err != nil {
				log.Println("Ошибка при отправке Документа")
			}
			Remove_Image(file_path)
		}
	}
}
