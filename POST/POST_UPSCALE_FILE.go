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

type Post_file struct {
	Status string `json:"status"`
	Data   struct {
		Id  string `json:"id"`
		Url string `json:"url"`
	} `json:"data"`
}

func PostImage(apiKey string, filePath string, upscale_factor string) (string, error) {
	url := "https://api.picsart.io/tools/1.0/upscale"
	// Создание буфера для данных запроса
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Ошибка при открытии файла", err)
		return "", err
	}
	defer file.Close()
	log.Println(filePath)
	// Создаем буфер и многокомпонентный писатель
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Добавляем параметр upscale_factor
	writer.WriteField("upscale_factor", upscale_factor)

	// Добавляем файл к части image
	part, err := writer.CreateFormFile("image", filePath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	writer.Close()

	// Создаем новый запрос
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Добавляем заголовки
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-picsart-api-key", apiKey)

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			log.Println("Ошибка запроса --- Номер ошибки : 400")
		}
		return "", nil
	}

	// Читаем ответ
	var postFile Post_file
	if err := json.NewDecoder(resp.Body).Decode(&postFile); err != nil {
		fmt.Println(err)
		return "", err
	}
	return postFile.Data.Url, nil
}
