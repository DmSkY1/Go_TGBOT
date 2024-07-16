package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func InstallPhoto(fileID string, bot *tgbotapi.BotAPI) {
	log.Println("Открываем файл")
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		log.Println(err)
	}
	// Получаем прямую ссылку на файл фотографии
	fileURL := file.Link(bot.Token)
	response, err := http.Get(fileURL) // делаем запрос
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close() // закрываем процесс
	file_name := fmt.Sprintf("%s.jpg", fileID)
	file_path := filepath.Join("picture/", file_name)

	//созранение файла
	file_save, err := os.Create(file_path)
	if err != nil {
		log.Println(err)
	}
	defer file_save.Close()
	log.Println("Копирование содержимого файла")
	// Копируем содержимое ответа в файл
	_, err = io.Copy(file_save, response.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println("Фотография сохранена")
}
