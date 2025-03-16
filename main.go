package main

import (
	"fmt"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	cmd "main.go/Comand"
)

// Константы состояний
const (
	IdleState = iota
	WaitingForImageState
	WaitingForImageToRemoveBackground
	WaitingForChangeBackground
)

var (
	userStates        = make(map[int64]int)
	userActiveCommand = make(map[int64]string)
	mu                sync.RWMutex
)

const (
	CommandStart            = "start"
	CommandUpscaleImageX2   = "upscale_image_x2"
	CommandUpscaleImageX4   = "upscale_image_x4"
	CommandUpscaleImageX6   = "upscale_image_x6"
	CommandUpscaleImageX8   = "upscale_image_x8"
	CommandRemoveBackground = "remove_background"
	CommandHelp             = "help"
	ChangeBackground        = "change_background"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("6798282567:AAEC8jxADvq9CTSaHBtmxYkflbP7pj72gvU")
	if err != nil {
		log.Fatal(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 20

	updatesChan := bot.GetUpdatesChan(u)
	log.Println("Бот запущен")

	// Обработка всех новых сообщений в отдельном горутине, чтобы не заблокировать основной поток обработки сообщений
	var wg sync.WaitGroup

	// Инициализация состояний пользователей
	for update := range updatesChan {
		if update.Message != nil {
			wg.Add(1)
			go func(update tgbotapi.Update) {
				defer wg.Done()
				processUpdate(bot, update)
			}(update)
		}
	}
	wg.Wait()
}

// Обработка нового сообщения
func processUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	mu.Lock()
	state := userStates[chatID]
	mu.Unlock()

	if update.Message.IsCommand() {
		handleIdleState(bot, update)
	} else {
		switch state {
		case WaitingForImageState:
			handleWaitingForImageState(bot, update)
		case WaitingForImageToRemoveBackground:
			handleWaitingForImageToRemoveBackground(bot, update)
		case WaitingForChangeBackground:
			handleChangeBackround(bot, update)
		}
	}

}

func handleIdleState(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	command := update.Message.Command()
	messageForUserToUpscaleImage := tgbotapi.NewMessage(chatID, "🌟✨ Пришлите мне изображение в виде файла или ссылку на него, и я превращу его в настоящий шедевр! 🎨💫 Ваше фото засияет новыми красками! 📸🌠")
	switch command {
	case CommandStart:
		sendStartMessage(bot, chatID)
		setUserState(chatID, IdleState, CommandStart)
	case CommandUpscaleImageX2, CommandUpscaleImageX4, CommandUpscaleImageX6, CommandUpscaleImageX8:
		bot.Send(messageForUserToUpscaleImage)
		setUserState(chatID, WaitingForImageState, command)
	case CommandRemoveBackground:
		bot.Send(messageForUserToUpscaleImage)
		setUserState(chatID, WaitingForImageToRemoveBackground, command)
	case CommandHelp:
		sendHelpMessage(bot, chatID)
	case ChangeBackground:
		setUserState(chatID, WaitingForChangeBackground, command)
	}
}

func handleChangeBackround(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	go cmd.Change_background(bot, update)
}
func handleWaitingForImageState(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	go cmd.Upscale_image(bot, update, &userStates, &userActiveCommand)
}

func handleWaitingForImageToRemoveBackground(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	go cmd.Remove_background_image(bot, update, &userStates, &userActiveCommand)
}

func sendHelpMessage(bot *tgbotapi.BotAPI, chatID int64) {
	startMessage := tgbotapi.NewAnimation(chatID, tgbotapi.FileURL("https://c.tenor.com/h9Izn8ibp4AAAAAd/tenor.gif"))
	startMessage.Caption = fmt.Sprintf("🌟 Вот список моих команд: 🚀✨\n\n" +
		"🔹 /start - Начать работу со мной. Я расскажу, как всё устроено!\n" +
		"🔹 /help - Показать список всех доступных команд.\n" +
		"🔹 /upscale_image_x2 - Увеличить качество изображения в 2 раза.\n" +
		"🔹 /upscale_image_x4 - Увеличить качество изображения в 4 раза.\n" +
		"🔹 /upscale_image_x6 - Увеличить качество изображения в 6 раза.\n" +
		"🔹 /upscale_image_x8 - Увеличить качество изображения в 8 раза.\n" +
		"🔹 /remove_background - Удалить фон с изображения.\n\n" +
		"Выбирай нужную команду, и я помогу сделать твои фото ещё лучше! 🎨💫\n" +
		"Готов к работе!😊🚀")
	startMessage.ParseMode = "HTML"
	_, err := bot.Send(startMessage)
	if err != nil {
		log.Println("Ошибка отправки GIF:", err)
	}
}

func sendStartMessage(bot *tgbotapi.BotAPI, chatID int64) {
	startMessage := tgbotapi.NewAnimation(chatID, tgbotapi.FileURL("https://gifs.obs.ru-moscow-1.hc.sbercloud.ru/4fdc711265b196861471c3b2a2ce51aa4ced8c09bfe12493a5edbc1e1a8e3700.gif"))
	startMessage.Caption = fmt.Sprintf("👋 Привет! Я бот, который поможет вам с обработкой фотографий и файлов. 📸📁\n\n" +
		"Вот что я умею:\n\n" +
		"🛠 Редактировать метаданные файлов\n" +
		"✨ Улучшать качество фотографий\n" +
		"🧹 Удалять водяные знаки с изображений\n" +
		"🔄 Преобразовывать файлы в различные форматы\n\n" +
		"Если хотите узнать больше о моих функциях, просто воспользуйтесь командой /help. Если у вас есть файл или фотография, которую нужно обработать, отправьте её мне, и я возьмусь за дело! 😊\n\n" +
		"Если возникнут вопросы или потребуется помощь, не стесняйтесь обращаться. Я всегда рад помочь!")
	startMessage.ParseMode = "HTML"
	_, err := bot.Send(startMessage)
	if err != nil {
		log.Println("Ошибка отправки GIF:", err)
	}
}

func setUserState(chatID int64, state int, command string) {
	mu.Lock()
	defer mu.Unlock()
	userStates[chatID] = state
	userActiveCommand[chatID] = command

}
