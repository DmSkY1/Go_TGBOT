package POST

import (
	"bytes"
	"encoding/json"
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
	url := "https://api.picsart.io/tools/1.0/removebg"

	file, err := os.Open(filepath)
	if err != nil {
		log.Println("Ошибкуа при открытии файла !!!!")
		return "", err
	}

	defer file.Close()
	var requestBody bytes.Buffer

	writer := multipart.NewWriter(&requestBody)

	writer.WriteField("output_type", "cutout")
	writer.WriteField("format", "PNG")

	form, err := writer.CreateFormFile("image", filepath)
	if err != nil {
		log.Println("Ошибка при создании multipart/form-data!!!!")
		return "", err
	}

	_, err = io.Copy(form, file)
	if err != nil {
		log.Println("Ошибка при создании формы")
		return "", err
	}

	writer.Close()

	//
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		log.Println("Ошибка при создании запроса")
	}

	//
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-picsart-api-key", apikey)

	//
	client := &http.Client{}

	//
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Ошибка при выполнении запроса -", err)
		return "", err
	}
	//
	defer resp.Body.Close()

	var file_data File_data
	if err := json.NewDecoder(resp.Body).Decode(&file_data); err != nil {
		log.Println(err)
		return "", err
	}
	return file_data.Data.Url, nil

}
