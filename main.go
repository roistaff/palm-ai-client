package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

type Candidate struct {
	Output string `json:"output"`
}

type Error struct {
	Code    int
	Message string
	Status  string
}

var Log string

func AddLog(text string) {
	Log += text + "\n"
}

func SaveLog() error {
	dir, _ := os.UserHomeDir()
	file_name := dir + "/palm-log" + time.Now().Format("2006,01,02,15:04") + ".txt"
	err := ioutil.WriteFile(file_name, []byte(Log), 0644)
	return err
}

func main() {
	savelog := flag.Bool("p", false, "bool flag")
	flag.Parse()
	color.Cyan("PaLM2 AI client\n\n")

	API_KEY := os.Getenv("PALM_KEY")
	if API_KEY == "" {
		fmt.Println("API key not found. Please set your API key using `export PALM_KEY=your_api_key`")
		os.Exit(0)
	}

	posturl := "https://generativelanguage.googleapis.com/v1beta2/models/text-bison-001:generateText?key=" + API_KEY

	AddLog("PaLM2 talk log")
	for {
		client := &http.Client{}

		validate := func(input string) error {
			if input == "" {
				return errors.New("No values.")
			}
			return nil
		}

		prompt := promptui.Prompt{
			Label:    "You  ",
			Validate: validate,
		}

		result, err := prompt.Run()

		if err == promptui.ErrInterrupt {
			if *savelog == true {
				SaveLog()
				fmt.Println("Log has saved.")
			}
			fmt.Println(color.HiGreenString("\U0001F5F8"), "Program has run successfully.")
			os.Exit(1)
		}

		text := strings.Replace(result, "\"", "'", -1)
		AddLog("You  :" + text)
		promptData := []byte(`{"prompt":{"text":"` + text + `"}}`)

		request, err := http.NewRequest("POST", posturl, bytes.NewBuffer(promptData))
		if err != nil {
			panic(err)
		}
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")

		response, err := client.Do(request)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		var data map[string][]Candidate
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println("JSON parse error:", err)
			os.Exit(0)
		}

		if len(data["candidates"]) > 0 {
			output := data["candidates"][0].Output
			result := markdown.Render(output, 80, 7)
			view := string(result)
			if view != "" {
				view = view[7:]
			}
			fmt.Println("PaLM2:", view)
			AddLog("PaLM2:" + output)
		} else {
			fmt.Println("Error: No response. Please try again.")
			AddLog("No response")
		}
	}
}
