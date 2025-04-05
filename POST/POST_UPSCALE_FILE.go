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

	// Создание буфера для данных запроса
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Ошибка при открытии файла", err)
		return "", fmt.Errorf("ошибка при открытии файла: %w", err)
	}
	defer file.Close()

	// Создаем буфер и многокомпонентный писатель
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Добавляем параметр upscale_factor
	if err := writer.WriteField("upscale_factor", upscale_factor); err != nil {
		log.Printf("Ошибка при открытии файла: %v", err)
		return "", fmt.Errorf("ошибка при открытии файла: %w", err)
	}

	// Добавляем файл к части image
	part, err := writer.CreateFormFile("image", filePath)
	if err != nil {
		log.Printf("Ошибка при создании формы для файла: %v", err)
		return "", fmt.Errorf("ошибка при создании формы для файла: %w", err)
	}
	if _, err = io.Copy(part, file); err != nil {
		log.Printf("Ошибка при копировании файла: %v", err)
		return "", fmt.Errorf("ошибка при копировании файла: %w", err)
	}
	if err := writer.Close(); err != nil {
		log.Printf("Ошибка при закрытии writer: %v", err)
		return "", fmt.Errorf("ошибка при закрытии writer: %w", err)
	}

	// Создаем новый запрос
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		log.Printf("Ошибка при создании запроса: %v", err)
		return "", fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	// Добавляем заголовки
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-picsart-api-key", apiKey)

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		return "", fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			log.Println("[-] request failed -- Error number: 400")
			return "", fmt.Errorf("Request failed with status %d", http.StatusBadRequest)
		}
	}

	// Читаем ответ
	var postFile Post_file
	if err := json.NewDecoder(resp.Body).Decode(&postFile); err != nil {
		log.Printf("Ошибка при декодировании ответа: %v", err)
		return "", fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	result := strings.Replace(postFile.Data.Url, "?type=jpg&to=max&r=0", "", 1)

	return result, nil
}
