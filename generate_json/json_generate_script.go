package main

import (
	"bufio"
	"encoding/json"
	"os"
)

type ApiKeyData struct {
	Using int `json:"using"`
	Left  int `json:"left_to_use"`
}

type JsonTemplate struct {
	ApiKeys map[string]ApiKeyData `json:"api_keys"`
}

func main() {

	file, err := os.Open("NewApiKey.txt")
	if err != nil {
		return
	}
	defer file.Close()

	var keys []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		keys = append(keys, scanner.Text())
	}

	var jsonData JsonTemplate

	jsonData = JsonTemplate{
		ApiKeys: make(map[string]ApiKeyData),
	}

	for _, value := range keys {
		jsonData.ApiKeys[value] = ApiKeyData{Using: 0, Left: 25}
		fileData, err := json.MarshalIndent(jsonData, "", " ")
		if err != nil {
			return
		}
		if err := os.WriteFile("using_API_keys.json", fileData, 0644); err != nil {
			return
		}
	}
}
