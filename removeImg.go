package main

import (
	"log"
	"os"
	"time"
)

func Remove_Image(file_name string) error {
	time.Sleep(100 * time.Millisecond)
	err := os.Remove(file_name)
	if err != nil {
		log.Println(err)
	}
	return nil
}
