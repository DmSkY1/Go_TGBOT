package POST

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func DownloadFileUrl(url string, fileid string, filename string) (string, error) {

	response, e := http.Get(url)
	if e != nil {
		log.Println("\033[31m[Error]\033[0m Ошибка в GET запросе. DownloadFileUrl")
		return "", nil
	}
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create(fmt.Sprintf("picture/%s_%s.%s", fileid, filename, "jpeg"))
	if err != nil {
		log.Println("\033[31m[Error]\033[0m Ошибка создания файла. DownloadFileUrl")
		return "", nil
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Println("\033[31m[Error]\033[0m Ошибка копирования данных в файл. DownloadFileUrl")
		return "", nil
	}

	return fmt.Sprintf("picture/%s_%s.%s", fileid, filename, "jpeg"), nil
}
