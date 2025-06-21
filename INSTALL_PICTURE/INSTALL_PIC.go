package INSTALL_PICTURE

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func InstallPhoto(filepath string, url string) error {
	// Создаем файл
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("\033[31m[Error]\033[0m Error when creating the file:", err)
	}
	defer out.Close()

	// Загружаем файл
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("\033[31m[Error]\033[0m Error when uploading a file:", err)
	}
	defer resp.Body.Close()

	// Копируем содержимое в файл
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("\033[31m[Error]\033[0m Error when copying content to a file:", err)
	}

	return nil
}
