package main

import (
	"fmt"
	"log"
	"sync"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	cmd "main.go/Comand"
)

// Константы состояний
const (
	IdleState = iota
	WaitingForImageState
	WaitingForImageToRemoveBackground
)

var (
	userStates        = make(map[int64]int)
	userActiveCommand = make(map[int64]string)
	mu                sync.RWMutex
)

func main() {
	bot, err := tgbotapi.NewBotAPI("6798282567:AAEC8jxADvq9CTSaHBtmxYkflbP7pj72gvU")
	if err != nil {
		log.Fatal(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 20

	updatesChan, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println("Ошибка при получении канала обновлений:", err)
		return
	}

	log.Println("Бот запущен")
	for update := range updatesChan {
		if update.Message != nil {
			go processUpdate(bot, update)
		}
	}
}

func processUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	mu.Lock()
	state := userStates[chatID]
	mu.Unlock()

	// Обработка команд
	if update.Message.IsCommand() {
		command := update.Message.Command()
		if state != WaitingForImageToRemoveBackground && state != WaitingForImageState {
			messageForUserToUpscaleImage := tgbotapi.NewMessage(chatID, "📷 Отправьте мне изображение в виде файла или ссылку на него, и я постараюсь сделать его лучше! ✨")
			switch command {
			case "start":
				start_message := tgbotapi.NewAnimationShare(chatID, "https://gifs.obs.ru-moscow-1.hc.sbercloud.ru/4fdc711265b196861471c3b2a2ce51aa4ced8c09bfe12493a5edbc1e1a8e3700.gif")
				start_message.Caption = fmt.Sprintf("👋 Привет! Я бот, который поможет вам с обработкой фотографий и файлов. 📸📁\n\n" +
					"Вот что я умею:\n\n" +
					"🛠 Редактировать метаданные файлов\n" +
					"✨ Улучшать качество фотографий\n" +
					"🧹 Удалять водяные знаки с изображений\n" +
					"🔄 Преобразовывать файлы в различные форматы\n\n" +
					"Если хотите узнать больше о моих функциях, просто воспользуйтесь командой /help. Если у вас есть файл или фотография, которую нужно обработать, отправьте её мне, и я возьмусь за дело! 😊\n\n" +
					"Если возникнут вопросы или потребуется помощь, не стесняйтесь обращаться. Я всегда рад помочь!")

				start_message.ParseMode = "HTML"
				_, err := bot.Send(start_message)
				if err != nil {
					log.Println("Ошибка отправки GIF:", err)
				}
				mu.Lock()
				userStates[chatID] = IdleState
				userActiveCommand[chatID] = "start"
				mu.Unlock()
			case "upscale_image_x2":
				bot.Send(messageForUserToUpscaleImage)
				mu.Lock()
				userStates[chatID] = WaitingForImageState
				userActiveCommand[chatID] = "upscale_image_x2"
				mu.Unlock()
			case "upscale_image_x4":
				bot.Send(messageForUserToUpscaleImage)
				mu.Lock()
				userStates[chatID] = WaitingForImageState
				userActiveCommand[chatID] = "upscale_image_x4"
				mu.Unlock()
			case "upscale_image_x6":
				bot.Send(messageForUserToUpscaleImage)
				mu.Lock()
				userStates[chatID] = WaitingForImageState
				userActiveCommand[chatID] = "upscale_image_x6"
				mu.Unlock()
			case "upscale_image_x8":
				bot.Send(messageForUserToUpscaleImage)
				mu.Lock()
				userStates[chatID] = WaitingForImageState
				userActiveCommand[chatID] = "upscale_image_x8"
				mu.Unlock()
			case "remove_background":
				bot.Send(messageForUserToUpscaleImage)
				mu.Lock()
				userStates[chatID] = WaitingForImageToRemoveBackground
				userActiveCommand[chatID] = "remove_background"
				mu.Unlock()

			}
		}
	} else {
		switch state {
		case WaitingForImageState:
			go func() {
				cmd.Upscale_image(bot, update, &userStates, &userActiveCommand)
			}()
		case WaitingForImageToRemoveBackground:
			go func() {
				cmd.Remove_background_image(bot, update, &userStates, &userActiveCommand)
			}()
		}
	}

}
