package Commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	installPhoto "main.go/INSTALL_PICTURE"
	post_file "main.go/POST"
	rand_key "main.go/RandomKey"
)

// Константы состоянийы
const (
	_ = iota
	WaitingForImageState
	WaitingForProcessing
)

// Состоние пользователей
var (
	mu sync.RWMutex
	wg sync.WaitGroup
)

// Функция для обработки событий upscale_image
func Upscale_image(bot *tgbotapi.BotAPI, update tgbotapi.Update, us_state *map[int64]int, us_active_commang *map[int64]string) {

	// Определение нужных переменных
	chatID := update.Message.Chat.ID
	messageID := update.Message.MessageID
	// Проверка, является ли сообщение сжатой фотографией
	if update.Message.Photo != nil {
		invalidFormatMessage(bot, chatID, messageID, us_active_commang)
		return

	} else if update.Message.Document != nil && isImageFile(update.Message.Document) {
		if getUserState(chatID, us_state) == WaitingForImageState {
			setUserState(chatID, us_state, WaitingForProcessing)
			informationMessage(bot, chatID, messageID)
			log.Printf("\033[32m[INFO]\033[0m Получена фотография от пользователя. ChatID [%d]", chatID)

			// Получаем API-ключ из файла
			api_key, err := rand_key.GetRandomAPIKey()
			if err != nil {
				log.Println("\033[31m[Error]\033[0m Ошибка при получении API-ключа:", err)
				return
			}
			log.Printf("\033[32m[INFO]\033[0m API ключ сгенерирован. ChatID [%d]", chatID)

			// Получение id документа
			fileID := update.Message.Document.FileID
			file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
			if err != nil {
				log.Println("\033[31m[Error]\033[0m Ошибка при получении файла:", err)
				return
			}
			log.Printf("\033[32m[INFO]\033[0m ID фотографии получено. ChatID [%d]", chatID)

			// Установка и отправка на обработку сообщения
			downloadURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)
			filepath := filepath.Join("picture", fmt.Sprintf("%s_%s", fileID, update.Message.Document.FileName))

			// Скачиваем файл из Telegram API
			err = installPhoto.InstallPhoto(filepath, downloadURL)
			if err != nil {
				log.Println(err)
			}
			log.Printf("\033[32m[INFO]\033[0m Фотография успешно скачена. ChatID [%d]", chatID)

			upscale_factor := "4"
			wg.Add(1)
			go func() {
				defer wg.Done()
				// Отправка запроса на обработку фотографии
				result, err := post_file.PostImage(api_key, filepath, upscale_factor)

				log.Printf("\033[32m[INFO]\033[0m Фотография от пользователя [%d] отправлена в обработку", chatID)

				if err != nil {
					log.Println(err)
					deleteMsg := tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID+1)
					if _, err := bot.Request(deleteMsg); err != nil {
						log.Println("\033[31m[Error]\033[0m Ошибка при удалении сообщения:", err)
					}
					photoMsg := tgbotapi.NewAnimation(chatID, tgbotapi.FileURL("https://c.tenor.com/ErB2RhcIXpwAAAAd/tenor.gif"))
					photoMsg.ReplyToMessageID = update.Message.MessageID
					photoMsg.Caption = "🚫 Упс, не удалось обработать вашу фотографию!\n😟Вероятно, файл повреждён или его формат неправильный 🛠️.\nПопробуйте отправить другой файл!📸"

					// Отправляем новое сообщение с гифом
					if _, err := bot.Send(photoMsg); err != nil {
						log.Println("\033[31m[Error]\033[0m Ошибка при отправке анимации:", err)
					}
					defer os.Remove(filepath)
					return
				}
				defer os.Remove(filepath)
				deleteMsg := tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID+1)
				if _, err := bot.Request(deleteMsg); err != nil {
					log.Println("\033[31m[Error]\033[0m Ошибка при удалении сообщения:", err)
				}
				log.Printf("\033[32m[INFO]\033[0m Информационное сообщение удалено успешно. ChatID [%d]", chatID)
				// https://s7.gifyu.com/images/SGWok.gif митсури
				// Отправляем новое сообщение с гифом и обновленной подписью

				document := tgbotapi.NewDocument(chatID, tgbotapi.FileURL(result)) // Отправка готового результата
				document.Caption = "🎯 *Готово!* 🎯\n\n" +
					"🎉 Ваше изображение теперь выглядит лучше! 📸🌈\n\n" +
					"🔍 Посмотрите внимательно и наслаждайтесь результатом! 😊"
				// Отправляем сообщение
				_, err = bot.Send(document)
				if err != nil {

					log.Println("\033[31m[Error]\033[0m Ошибка отправки документа121:", err)
				}
				log.Printf("\033[32m[INFO]\033[0m Документ успешно отправлен. ChatID [%d]", chatID)
			}()
			wg.Wait()
			setUserState(chatID, us_state, IdleState)
			return
		}

	} else if update.Message != nil && update.Message.Text != "" {
		if strings.HasPrefix(update.Message.Text, "https://") || strings.HasPrefix(update.Message.Text, "http://") {
			if getUserState(chatID, us_state) == WaitingForImageState {
				setUserState(chatID, us_state, WaitingForProcessing)
				informationMessage(bot, chatID, messageID)

				api_key, err := rand_key.GetRandomAPIKey()
				if err != nil {
					log.Println("\033[31m[Error]\033[0m Ошибка при получении API-ключа:", err)
					return
				}
				log.Printf("\033[32m[INFO]\033[0m API ключ сгенерирован. ChatID [%d]", chatID)

				upscale_factor := "4"
				url := update.Message.Text

				wg.Add(1)
				go func() {
					defer wg.Done()
					// Загрузка файла
					res, err := post_file.DownloadFileUrl(url, "URl_Image", strconv.Itoa(int(chatID)))
					if err != nil {
						log.Println("\033[31m[Error]\033[0m Ошибка при загрузке файла:", err)
						return
					}
					log.Printf("\033[32m[INFO]\033[0m Фотография успешно скачена. ChatID [%d]", chatID)

					// Отправка на обработку
					res_post, err := post_file.PostImage(api_key, res, upscale_factor)
					if err != nil {
						log.Println(err)
						setUserState(chatID, us_state, IdleState)
						return
					}
					log.Printf("\033[32m[INFO]\033[0m Документ успешно обработан. ChatID [%d]", chatID)

					// Удаление не обработанной фотографии
					defer os.Remove(res)

					// Отправление документа
					document := tgbotapi.NewDocument(chatID, tgbotapi.FileID(res_post))
					document.Caption = fmt.Sprintf("🎯 *Готово!* 🎯\n\n" +
						"🎉 Ваше изображение теперь выглядит лучше! 📸🌈\n\n" +
						"🔍 _Посмотрите внимательно и наслаждайтесь результатом!_ 😊")
					_, err = bot.Send(document)
					if err != nil {
						log.Println("\033[31m[Error]\033[0m Ошибка отправки документа:", err)
						return
					}
					log.Printf("\033[32m[INFO]\033[0m Документ успешно отправлен. ChatID [%d]", chatID)

					// Удаление сообщения последнего сообщения
					wg.Add(1)
					go func() {
						defer wg.Done()
						deleteMsg := tgbotapi.NewDeleteMessage(chatID, update.Message.MessageID+1)
						if _, err := bot.Request(deleteMsg); err != nil {
							log.Println("\033[31m[Error]\033[0m Ошибка при удалении сообщения:", err)
							return
						}
						log.Printf("\033[32m[INFO]\033[0m Информационное сообщение удалено успешно. ChatID [%d]", chatID)
					}()
					setUserState(chatID, us_state, IdleState)
				}()
				wg.Wait()
				return
			}
		} else {
			if (*us_state)[chatID] == WaitingForImageState {
				bot.Send(tgbotapi.NewMessage(chatID, "Пожалуйста, отправьте фотографию. Бот не воспринимает ваш текст. Я буду ждать ващей фотографии ))"))
				return
			} else {
				errorMessage(bot, chatID, messageID, us_active_commang)
				setUserState(chatID, us_state, IdleState)
				return
			}
		}
	}
}

func errorMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int, us_active_commang *map[int64]string) {
	errorMessage_url := tgbotapi.NewAnimation(chatID, tgbotapi.FileURL("https://gifs.obs.ru-moscow-1.hc.sbercloud.ru/1d06b49de1ac9de5cbc468d1d449d74658d39b7471c689cf5ec7570106908a9e.gif"))
	errorMessage_url.Caption = fmt.Sprintf("🚫 <i><strong>Ошибка! Неверный формат запроса.</strong></i> Пожалуйста, отправьте правильный URL-адрес. 🌐\n\n"+
		"Чтобы продолжить, повторите команду /%s с <strong>корректным URL.</strong>⚠️", (*us_active_commang)[chatID])
	errorMessage_url.ParseMode = "HTML"
	errorMessage_url.ReplyToMessageID = messageID // Указываем ID сообщения для ответа
	bot.Send(errorMessage_url)
}

func informationMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	// Формируем сообщение для отправки
	photoMsg := tgbotapi.NewAnimation(chatID, tgbotapi.FileURL("https://i.gifer.com/origin/38/3823ad20629c89b3dd4821b80eee79eb_w200.gif"))
	photoMsg.ReplyToMessageID = messageID
	photoMsg.Caption = "📸 Ваше фото в обработке! 🚀\n" +
		"Я занимаюсь улучшением и увеличением вашего изображения. Это займет примерно 10 секунд. ⏳✨\n\n" +
		"Пожалуйста, подождите немного, и ваше фото будет готово к просмотру. Спасибо за терпение! 😊"
	bot.Send(photoMsg)
}

func invalidFormatMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int, us_active_commang *map[int64]string) {
	error_photo := tgbotapi.NewAnimation(chatID, tgbotapi.FileURL("https://gifs.obs.ru-moscow-1.hc.sbercloud.ru/6e330e6add14701f3a495e17e01e17cccb377fba621adb0f6aeec8430cfc5153.gif"))
	error_photo.Caption = fmt.Sprintf("🚫 <i><strong>Неверный формат запроса!</strong></i> Вы отправили неправильный тип данных. Пожалуйста, загрузите фотографию как файл (например, JPEG или PNG), а не сжатое изображение. 📁\n\n"+
		"Чтобы продолжить, повторите команду /%s и отправьте фотографию в виде файла. ⚠️", (*us_active_commang)[chatID])
	error_photo.ParseMode = "HTML"
	error_photo.ReplyToMessageID = messageID // Указываем ID сообщения для ответа
	_, err := bot.Send(error_photo)
	if err != nil {
		log.Println("\033[31m[Error]\033[0m Ошибка отправки GIF:", err)
	}
}

// Проверка, является ли документ фотографией
func isImageFile(doc *tgbotapi.Document) bool {
	imageMimeTypes := []string{"image/jpeg", "image/png", "image/gif"}
	for _, mimeType := range imageMimeTypes {
		if doc.MimeType == mimeType {
			return true
		}
	}
	return false
}

func getUserState(chatID int64, us_state *map[int64]int) int {
	mu.RLock()
	defer mu.RUnlock()
	return (*us_state)[chatID]
}

func setUserState(chatID int64, us_state *map[int64]int, state_now int) {
	mu.Lock()
	defer mu.Unlock()
	(*us_state)[chatID] = state_now
}
