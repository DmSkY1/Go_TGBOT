package RandomKey

import (
	"bufio"
	"encoding/json"
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
		log.Println(err)
		return "", err
	}
	defer file.Close()

	// Read all lines into a slice
	var keys []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		keys = append(keys, scanner.Text())
	}

	// Check for errors during the file read
	if err := scanner.Err(); err != nil {
		return "", err
	}
	// read json file
	data, err := os.ReadFile("C:/Users/DmSkY/Desktop/Go_Bot/using_API_keys.json")
	if err != nil {
		return "", err
	}

	//create instance
	var api_keys ApiKeys

	//unmarshal json data to the struct
	err = json.Unmarshal(data, &api_keys)
	if err != nil {
		return "", err
	}

	keyName := keys[Random_index(filename, keys)]
	Key := api_keys.Keys[keyName]

	if Key.Using < 25 {
		Key.Using++
		Key.LeftToUse--

		api_keys.Keys[keyName] = Key

		updateData, err := json.MarshalIndent(api_keys, "", " ")
		if err != nil {
			log.Println("[-] Error marshaling json")
			return "", err
		}
		err = os.WriteFile("C:/Users/DmSkY/Desktop/Go_Bot/using_API_keys.json", updateData, 0644)
		if err != nil {
			log.Println("[-] Error writing json:", err)
			return "", err
		}

		log.Println("[+] Successful key verification")
		return keyName, nil
	} else if Key.Using >= 25 {
		return "", nil
	} else {
		return "", nil
	}

}

func Random_index(filename string, keys []string) int {
	// Seed the random number generator
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Select a random key
	randomIndex := randGen.Intn(len(keys))

	return randomIndex

}
