package INSTALL_PICTURE

import (
	"os"
	"io"
	"net/http"
)

func InstallPhoto(filepath string, url string) error {
	// Создаем файл
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Загружаем файл
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Копируем содержимое в файл
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}