package POST

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func DownloadFile(url string, fileid string, filename string) (string, error) {

	response, err := http.Get(url)
	if err != nil {
		log.Println("\033[31m[Error]\033[0m Ошибка в GET запросе. DownloadFile")
		return "", nil
	}
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create(fmt.Sprintf("picture/%s_%s.%s", fileid, filename[:strings.Index(filename, ".")], "jpeg"))
	if err != nil {
		log.Println("\033[31m[Error]\033[0m Ошибка создания файла. DownloadFile")
		return "", nil
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Println("\033[31m[Error]\033[0m Ошибка копирования данных в файл. DownloadFile")
		return "", nil
	}

	return fmt.Sprintf("picture/%s_%s.%s", fileid, filename[:strings.Index(filename, ".")], "jpeg"), nil
}
