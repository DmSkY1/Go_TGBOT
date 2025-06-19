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
)

type File_data struct {
	Status string `json:"status"`
	Data   struct {
		Id  string `json:"id"`
		Url string `json:"url"`
	} `json:"data"`
}

func RemoveBackground(apikey string, filepath string) (string, error) {
	log.Print("\033[32m[INFO]\033[0m Подготовка запроса к отправке.")
	url := "https://api.picsart.io/tools/1.0/removebg"

	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("\033[31m[Error]\033[0m Ошибка при открытии файла %v", err)
	}

	defer file.Close()
	var requestBody bytes.Buffer

	writer := multipart.NewWriter(&requestBody)

	if err := writer.WriteField("output_type", "cutout"); err != nil {
		return "", fmt.Errorf("\033[31m[Error]\033[0m Ошибка при создпнии поля№1 в запросе: %v", err)
	}

	if err := writer.WriteField("format", "PNG"); err != nil {
		return "", fmt.Errorf("\033[31m[Error]\033[0m Ошибка при создпнии поля№2 в запросе: %v", err)
	}

	form, err := writer.CreateFormFile("image", filepath)
	if err != nil {
		return "", fmt.Errorf("\033[31m[Error]\033[0m Ошибка при создании формы для файла: %v", err)
	}

	if _, err = io.Copy(form, file); err != nil {
		return "", fmt.Errorf("\033[31m[Error]\033[0m Ошибка при копировании файла: %v", err)
	}

	writer.Close()

	// Создаем новый запрос
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return "", fmt.Errorf("\033[31m[Error]\033[0m Ошибка при создании запроса: %v", err)
	}

	// Добавляем заголовки
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-picsart-api-key", apikey)

	// Выполняем запрос
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("\033[31m[Error]\033[0m  Ошибка выполнения завпроса %v", err)
	}
	log.Print("\033[32m[INFO]\033[0m Запрос отправлен.")

	defer resp.Body.Close()

	var file_data File_data
	if err := json.NewDecoder(resp.Body).Decode(&file_data); err != nil {
		return "", fmt.Errorf("\033[31m[Error]\033[0m Ошибка при декодировании ответа: %v", err)
	}
	log.Print("\033[32m[INFO]\033[0m Данные запроса успешно обратоны.")
	return file_data.Data.Url, nil
}
