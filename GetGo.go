package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func Install_picture(url string, filepath string) error {
	req, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer req.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	_, err = io.Copy(file, req.Body)
	if err != nil {
		log.Println(err)
	}
	return nil
}
