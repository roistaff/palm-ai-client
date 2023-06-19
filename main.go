package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	markdown "github.com/MichaelMure/go-term-markdown"
)

type Candidate struct {
	Output string `json:"output"`
}

func InputPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}
func main() {
	API_KEY := os.Getenv("PALM_KEY")
	if API_KEY == "" {
		fmt.Println("api key not found.Please type in terminal `export PALM_KEY=your_api_key`")
		os.Exit(0)
	}
	posturl := "https://generativelanguage.googleapis.com/v1beta2/models/text-bison-001:generateText?key=" + API_KEY
	for {
		client := new(http.Client)
		text := strings.Replace(InputPrompt(">>> "), "\"", "'", -1)
		var promptData = []byte(`{"prompt":{"text":"` + text + `"}}`)
		request, error := http.NewRequest("POST", posturl, bytes.NewBuffer(promptData))
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")
		response, error := client.Do(request)
		if error != nil {
			panic(error)
		}
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		var data map[string][]Candidate
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println("JSON parse erorr:", err)
			os.Exit(0)
		}
		if len(data["candidates"]) > 0 {
			output := data["candidates"][0].Output
			result := markdown.Render(output, 80, 4)
			fmt.Println("PaLM:\n\n", string(result)[1:])
		} else {
			fmt.Println("Erorr:No response.Try agin.")
		}
	}
}
