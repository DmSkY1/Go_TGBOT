package POST

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type Post_file struct {
	Status string `json:"status"`
	Data   struct {
		Id  string `json:"id"`
		Url string `json:"url"`
	} `json:"data"`
}

const (
	url = "https://api.picsart.io/tools/1.0/upscale"
)

func PostImage(apiKey string, filePath string, upscale_factor string) (string, error) {
	log.Print("\033[32m[INFO]\033[0m Подготовка запроса к отправке.")

	// Создание буфера для данных запроса
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("\033[31m[Error]\033[0m Ошибка при открытии файла %v", err)
	}
	defer file.Close()

	// Создаем буфер и многокомпонентный писатель
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Добавляем параметр upscale_factor
	if err := writer.WriteField("upscale_factor", upscale_factor); err != nil {
		return "", fmt.Errorf("\033[31m[Error]\033[0m Ошибка при создпнии поля в запросе: %v", err)
	}

	// Добавляем файл к части image
	part, err := writer.CreateFormFile("image", filePath)
	if err != nil {
		return "", fmt.Errorf("\033[31m[Error]\033[0m Ошибка при создании формы для файла: %v", err)
	}
	if _, err = io.Copy(part, file); err != nil {
		return "", fmt.Errorf("\033[31m[Error]\033[0m Ошибка при копировании файла: %v", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("\033[31m[Error]\033[0m Ошибка при закрытии writer: %v", err)
	}

	// Создаем новый запрос
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return "", fmt.Errorf("\033[31m[Error]\033[0m Ошибка при создании запроса: %v", err)
	}

	// Добавляем заголовки
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-picsart-api-key", apiKey)

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("\033[31m[Error]\033[0m  Ошибка выполнения завпроса %v", err)
	}
	defer resp.Body.Close()
	log.Print("\033[32m[INFO]\033[0m Запрос отправлен.")

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return "", fmt.Errorf("\033[31m[Error]\033[0m Ошибка запроса: 400")
		}
	}

	// Читаем ответ
	var postFile Post_file
	if err := json.NewDecoder(resp.Body).Decode(&postFile); err != nil {
		return "", fmt.Errorf("\033[31m[Error]\033[0m Ошибка при декодировании ответа: %v", err)
	}
	result := strings.Replace(postFile.Data.Url, "?type=jpg&to=max&r=0", "", 1)

	log.Print("\033[32m[INFO]\033[0m Данные запроса успешно обратоны.")
	return result, nil
}
