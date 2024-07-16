package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type VisualScaleResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		TaskId      string  `json:"task_id"`
		Image       string  `json:"image"`
		ReturnType  uint    `json:"return_type"`
		Type        string  `json:"type"`
		Progress    uint    `json:"progress"`
		State       int     `json:"state"`
		TimeElapsed float64 `json:"time_elapsed"`
	} `json:"data"`
}

func Post(nameFile string, api_token string) string {
	log.Println("Выполнение запроса на обработку")
	imgF, err := os.Open(nameFile)
	if err != nil {
		log.Println(err)
	}
	defer imgF.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("sync", "1")
	part, err := writer.CreateFormFile("image_file", filepath.Base(nameFile))

	if err != nil {
		log.Println(err)
	}
	_, err = io.Copy(part, imgF)
	if err != nil {
		log.Println(err)
	}
	err = writer.Close()
	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest("POST", "https://techhk.aoscdn.com/api/tasks/visual/scale", body)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-API-KEY", api_token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	var taskResp VisualScaleResponse
	err = json.Unmarshal(data, &taskResp)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}
	log.Println("Запрос выполнен, переход к установке")
	return taskResp.Data.Image
}
