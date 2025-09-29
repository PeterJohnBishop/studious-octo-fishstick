package main

import (
	"fmt"
	"os"
	"studious-octo-fishstick/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if _, err := tea.NewProgram(tui.InitialRootModel()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}

// err := godotenv.Load()
// if err != nil {
// 	log.Fatal("Error loading .env file")
// }

// headers := []api.Header{
// 	{Key: "Authorization", Value: os.Getenv("API_KEY")},
// 	{Key: "Content-Type", Value: "application/json"},
// }
// url := "https://api.clickup.com/api/v2/list/901112072150/task"
// params := []api.Params{
// 	{Key: "Archived", Value: "false"},
// }
// resp, err := api.SendRequest("GET", headers, url, params, nil)
// if err != nil {
// 	fmt.Println("Error:", err)
// 	return
// }
// api.PrettyPrintJSON(resp)
