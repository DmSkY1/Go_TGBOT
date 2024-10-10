package RandomKey

import (
	"bufio"
	"math/rand"
	"os"
	"time"
)

func GetRandomAPIKey(filename string) (string, error) {
	// Open the file
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

	// Seed the random number generator
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Select a random key
	randomIndex := randGen.Intn(len(keys))
	return keys[randomIndex], nil
}
