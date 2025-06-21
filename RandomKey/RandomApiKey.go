package RandomKey

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

type ApiKey struct {
	Using     int `json:"using"`
	LeftToUse int `json:"left_to_use"`
}

type ApiKeys struct {
	Keys map[string]ApiKey `json:"api_keys"`
}

func GetRandomAPIKey() (string, error) {
	filename := "ApiKey.txt"
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("\033[31m[Error]\033[0m Error when opening a file %v.", err)
	}
	defer file.Close()

	var keys []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		keys = append(keys, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("\033[31m[Error]\033[0m Scanner error %v.", err)
	}
	data, err := os.ReadFile("using_api_keys.json")
	if err != nil {
		return "", fmt.Errorf("\033[31m[Error]\033[0m Error reading the JSON file %v.", err)
	}

	var api_keys ApiKeys

	err = json.Unmarshal(data, &api_keys)
	if err != nil {
		return "", fmt.Errorf("\033[31m[Error]\033[0m Error converting a JSON file to a structure %v.", err)
	}

	keyName := keys[Random_index(filename, keys)]
	Key := api_keys.Keys[keyName]
	for {
		if Key.Using < 25 {
			Key.Using++
			Key.LeftToUse--

			api_keys.Keys[keyName] = Key

			updateData, err := json.MarshalIndent(api_keys, "", " ")
			if err != nil {
				return "", fmt.Errorf("\033[31m[Error]\033[0m JSON processing error. %v", err)
			}
			err = os.WriteFile("using_api_keys.json", updateData, 0644)
			if err != nil {
				return "", fmt.Errorf("\033[31m[Error]\033[0m Error reading JSON. %v", err)
			}

			log.Println("\033[32m[INFO]\033[0m The key has been successfully generated.")
			return keyName, nil
		} else {
			keyName = keys[Random_index(filename, keys)]
			continue
		}
	}

}

func Random_index(filename string, keys []string) int {
	// Seed the random number generator
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Select a random key
	randomIndex := randGen.Intn(len(keys))

	return randomIndex

}
