package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/gookit/color"
)

type Response struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Text string `json:"text"`
}

type Config struct {
	Key       string `json:"Key"`
	MaxTokens int    `json:"MaxTokens"`
}

func main() {
	if len(os.Args) < 2 {
		CreateError("No args found.")
		return
	}
	current, err := user.Current()
	if err != nil {
		CreateError("Cannot find user.")
		return
	}
	data, err := os.ReadFile(fmt.Sprintf("/home/%v/.white/config.json", current.Username))
	if err != nil {
		CreateError("No file found.")
		return
	}
	var datas Config
	json.Unmarshal(data, &datas)

	if datas.Key == "" {
		CreateError("Invalid Key string.")
		return
	}
	if datas.MaxTokens == 0 {
		CreateError("Invalid MaxTokens int.")
		return
	}

	client := &http.Client{}
	url := "https://api.openai.com/v1/completions"

	tokens := datas.MaxTokens

	if len(os.Args) == 3 {
		value, err := strconv.Atoi(os.Args[2])
		if err != nil {
			CreateError("Invalid Max Token int, using default one.")
		} else {
			tokens = value
		}
	}
	MapRequest := map[string]interface{}{
		"model":      "text-davinci-003",
		"prompt":     os.Args[1],
		"max_tokens": tokens,
	}
	requestBody, _ := json.Marshal(MapRequest)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", datas.Key))

	resp, _ := client.Do(req)
	if err != nil {
		CreateError("Error during the request.")
		os.Exit(1)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var respm Response
	json.Unmarshal(bodyBytes, &respm)

	if resp.StatusCode == 200 {
		CreateReponse(os.Args[1], strings.Replace(respm.Choices[0].Text, "\n", "", 2))
	} else {
		CreateError("Error during the request.")
	}
}

func CreateError(a string) {
	fmt.Printf("%s%v%s %s\n", color.Gray.Renderln("["), color.Red.Renderln("✖"), color.Gray.Renderln("]"), a)
}

func CreateTag(code string, colors string) string {
	if colors == "yellow" {
		code = color.Yellow.Renderln(code)
	}
	if colors == "blue" {
		code = color.LightBlue.Renderln(code)
	}
	c := fmt.Sprintf("%s%v%s", color.Gray.Renderln("["), code, color.Gray.Renderln("]"))
	return c
}

func CreateReponse(question string, response string) {
	fmt.Printf(`
%v %s
%v
%s %s

`, CreateTag("?", "yellow"), question, color.Gray.Renderln("・"), CreateTag("&", "blue"), response)
}
