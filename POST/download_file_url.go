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
		log.Println("Ошмбка в respons", e)
		return "", nil
	}
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create(fmt.Sprintf("picture/%s_%s.%s", fileid, filename, "jpeg"))
	if err != nil {
		log.Println(err)
		return "", nil
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Println(err)
		return "", nil
	}

	return fmt.Sprintf("picture/%s_%s.%s", fileid, filename, "jpeg"), nil
}
