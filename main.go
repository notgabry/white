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
	WhiteUtils "white/utils"

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
	if runtime.GOOS != "linux" {
		WhiteUtils.CreateError("Your OS is not supported.")
		return
	}

	if len(os.Args) < 2 {
		WhiteUtils.CreateError("Incorrect use. Go to github.com/NotGabry/white for more info.")
		return
	}

	switch os.Args[1] {
	case "--query", "-q":
		if len(os.Args) < 3 {
			WhiteUtils.CreateError("No args found.")
			return
		}

		current, err := user.Current()
		if err != nil {
			WhiteUtils.CreateError("Cannot find user.")
			return
		}

		data, err := os.ReadFile(fmt.Sprintf("/home/%s/.white/config.json", current.Username))
		if err != nil {
			WhiteUtils.CreateError("No file found.")
			return
		}

		stats, err := os.ReadFile(fmt.Sprintf("/home/%s/.white/stats.json", current.Username))
		if err != nil {
			statsTemplate := map[string]int{
				"SuccessQuery": 0,
				"ErrorQuery":   0,
			}

			os.WriteFile(fmt.Sprintf("/home/%s/.white/stats.json", current.Username), WhiteUtils.MapToString(statsTemplate), 0644)
		}

		var stat WhiteUtils.Stats
		json.Unmarshal(stats, &stat)

		var datas Config
		json.Unmarshal(data, &datas)

		if datas.Key == "" {
			WhiteUtils.CreateError("Invalid Key string.")
			return
		}
		if datas.MaxTokens == 0 {
			WhiteUtils.CreateError("Invalid MaxTokens int.")
			return
		}

		client := &http.Client{}
		url := "https://api.openai.com/v1/completions"

		tokens := datas.MaxTokens

		if len(os.Args) == 4 {
			value, err := strconv.Atoi(os.Args[3])
			if err != nil {
				WhiteUtils.CreateError("Invalid Max Token int, using default one.")
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
			WhiteUtils.CreateError("Error during the request.")
			return
		}
		defer resp.Body.Close()

		bodyBytes, _ := io.ReadAll(resp.Body)
		var respm Response
		json.Unmarshal(bodyBytes, &respm)

		if resp.StatusCode == 200 {
			t := stat
			t.SuccessQuery++
			os.WriteFile(fmt.Sprintf("/home/%s/.white/stats.json", current.Username), WhiteUtils.JSONToString(t), 0644)
			WhiteUtils.CreateReponse(os.Args[2], strings.Replace(respm.Choices[0].Text, "\n", "", 2))
		} else {
			t := stat
			t.ErrorQuery++
			os.WriteFile(fmt.Sprintf("/home/%s/.white/stats.json", current.Username), WhiteUtils.JSONToString(t), 0644)
			WhiteUtils.CreateError("Error during the request.")
		}
	case "--stats", "-s":
		current, err := user.Current()
		if err != nil {
			WhiteUtils.CreateError("Cannot find user.")
			return
		}
		stats, err := os.ReadFile(fmt.Sprintf("/home/%s/.white/stats.json", current.Username))
		if err != nil {
			WhiteUtils.CreateError("Stats file not found.")
			return
		}

		var stat WhiteUtils.Stats
		json.Unmarshal(stats, &stat)

		fmt.Printf(`
%s Success Queries%s%v
%s Error Queries%s%v

`, WhiteUtils.CreateTag("✔", "green"), color.Gray.Renderln("・"), stat.SuccessQuery, WhiteUtils.CreateTag("✖", "red"), color.Gray.Renderln("・"), stat.ErrorQuery)

	default:
		WhiteUtils.CreateError("Incorrect use. Go to github.com/NotGabry/white for more info.")
	}
}
