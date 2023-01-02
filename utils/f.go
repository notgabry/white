package WhiteUtils

import (
	"encoding/json"
	"fmt"

	"github.com/gookit/color"
)

type Stats struct {
	SuccessQuery int `json:"SuccessQuery"`
	ErrorQuery   int `json:"ErrorQuery"`
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
		CreateError("Cannot convert map to string.")
		return nil
	}
	return d
}

func JSONToString(data Stats) []byte {
	d, err := json.Marshal(data)
	if err != nil {
		CreateError("Cannot convert json to string.")
		return nil
	}
	return d
}
