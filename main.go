package main

import (
	"fmt"
	"log"
	"os"
	"studious-octo-fishstick/api"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	headers := []api.Header{
		{Key: "Authorization", Value: os.Getenv("API_KEY")},
		{Key: "Content-Type", Value: "application/json"},
	}
	url := "https://api.clickup.com/api/v2/list/901112072150/task"
	params := []api.Params{
		{Key: "Archived", Value: "false"},
	}
	resp, err := api.SendRequest("GET", headers, url, params, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	api.PrettyPrintJSON(resp)
}
