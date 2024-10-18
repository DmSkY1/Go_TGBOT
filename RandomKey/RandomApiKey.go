package RandomKey

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
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

func GetRandomAPIKey(filename string) (string, error) {
	// open the file
	file, err := os.Open(filename)
	if err != nil {
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
	data, err := ioutil.ReadFile("../using_API_keys.json")
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
	Key := api_keys.Keys[keys[Random_index(filename, keys)]]

	if Key.Using < 25 {
		Key.Using++
		Key.LeftToUse--
		log.Println("[+] Successful key verification")
		go func() {
			updateData, err := json.MarshalIndent(api_keys, "", " ")
			if err != nil {
				log.Println("[-] Error marshaling json")
				return
			}
			err = ioutil.WriteFile("../using_API_keys.json", updateData, 0644)
			if err != nil {
				log.Println("[-] Error writing json:", err)
				return
			}
		}()
		return keys[Random_index(filename, keys)], nil
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
