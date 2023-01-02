package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"runtime"
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

type Stats struct {
	SuccessQuery int `json:"SuccessQuery"`
	ErrorQuery   int `json:"ErrorQuery"`
}

func main() {
	if runtime.GOOS != "linux" {
		CreateError("Your OS is not supported.")
		return
	}

	if len(os.Args) < 2 {
		CreateError("Incorrect use. Go to github.com/NotGabry/white for more info.")
		return
	}

	switch os.Args[1] {
	case "--query", "-q":
		if len(os.Args) < 3 {
			CreateError("No args found.")
			return
		}

		current, err := user.Current()
		if err != nil {
			CreateError("Cannot find user.")
			return
		}

		data, err := os.ReadFile(fmt.Sprintf("/home/%s/.white/config.json", current.Username))
		if err != nil {
			CreateError("No file found.")
			return
		}

		stats, err := os.ReadFile(fmt.Sprintf("/home/%s/.white/stats.json", current.Username))
		if err != nil {
			statsTemplate := map[string]int{
				"SuccessQuery": 0,
				"ErrorQuery":   0,
			}

			os.WriteFile(fmt.Sprintf("/home/%s/.white/stats.json", current.Username), MapToString(statsTemplate), 0644)
		}

		var stat Stats
		json.Unmarshal(stats, &stat)

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

		if len(os.Args) == 4 {
			value, err := strconv.Atoi(os.Args[3])
			if err != nil {
				CreateError("Invalid Max Token int, using default one.")
			} else {
				tokens = value
			}
		}
		MapRequest := map[string]interface{}{
			"model":      "text-davinci-003",
			"prompt":     os.Args[2],
			"max_tokens": tokens,
		}
		requestBody, _ := json.Marshal(MapRequest)

		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", datas.Key))

		resp, err := client.Do(req)
		if err != nil {
			CreateError("Error during the request.")
			return
		}
		defer resp.Body.Close()

		bodyBytes, _ := io.ReadAll(resp.Body)
		var respm Response
		json.Unmarshal(bodyBytes, &respm)

		if resp.StatusCode == 200 {
			t := stat
			t.SuccessQuery++
			os.WriteFile(fmt.Sprintf("/home/%s/.white/stats.json", current.Username), JSONToString(t), 0644)
			CreateReponse(os.Args[2], strings.Replace(respm.Choices[0].Text, "\n", "", 2))
		} else {
			t := stat
			t.ErrorQuery++
			os.WriteFile(fmt.Sprintf("/home/%s/.white/stats.json", current.Username), JSONToString(t), 0644)
			CreateError("Error during the request.")
		}
	case "--stats", "-s":
		current, err := user.Current()
		if err != nil {
			CreateError("Cannot find user.")
			return
		}
		stats, err := os.ReadFile(fmt.Sprintf("/home/%s/.white/stats.json", current.Username))
		if err != nil {
			CreateError("Stats file not found.")
			return
		}

		var stat Stats
		json.Unmarshal(stats, &stat)

		fmt.Printf(`
%s Success Queries%s%v
%s Error Queries%s%v

`, CreateTag("✔", "green"), color.Gray.Renderln("・"), stat.SuccessQuery, CreateTag("✖", "red"), color.Gray.Renderln("・"), stat.ErrorQuery)

	default:
		CreateError("Incorrect use. Go to github.com/NotGabry/white for more info.")
	}
}

func CreateError(a string) {
	fmt.Printf("%s%s%s %s\n", color.Gray.Renderln("["), color.Red.Renderln("✖"), color.Gray.Renderln("]"), a)
}

func CreateTag(code string, colors string) string {
	if colors == "yellow" {
		code = color.Yellow.Renderln(code)
	}
	if colors == "blue" {
		code = color.LightBlue.Renderln(code)
	}
	if colors == "green" {
		code = color.Green.Renderln(code)
	}
	if colors == "red" {
		code = color.Red.Renderln(code)
	}
	c := fmt.Sprintf("%s%s%s", color.Gray.Renderln("["), code, color.Gray.Renderln("]"))
	return c
}

func CreateReponse(question string, response string) {
	fmt.Printf(`
%s %s
%s
%s %s

`, CreateTag("?", "yellow"), question, color.Gray.Renderln("・"), CreateTag("&", "blue"), response)
}

func MapToString(data map[string]int) []byte {
	d, err := json.Marshal(data)
	if err != nil {
		CreateError("Cannot converting to map.")
		return nil
	}
	return d
}

func JSONToString(data Stats) []byte {
	d, err := json.Marshal(data)
	if err != nil {
		CreateError("Cannot converting to map.")
		return nil
	}
	return d
}
