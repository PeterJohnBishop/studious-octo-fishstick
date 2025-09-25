package main

import (
	"fmt"
	"log"
	"os"
	"studious-octo-fishstick/api"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hella World")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	url := "https://api.clickup.com/api/v2/user"
	headers := []api.Header{
		{Key: "Authorization", Value: os.Getenv("API_KEY")},
		{Key: "Content-Type", Value: "application/json"},
	}
	data, err := api.SendRequest("GET", headers, url, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Response:", string(data))
}
