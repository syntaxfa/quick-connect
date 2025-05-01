package main

import (
	"fmt"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"golang.org/x/text/language"
)

func main() {
	languages := []translation.Language{
		{language.English, "en.json"},
		{language.Persian, "fa.json"},
		{language.French, "fr.json"},
	}

	t, tErr := translation.New(languages...)
	if tErr != nil {
		panic(tErr.Error())
	}

	fmt.Println(t.TranslateMessage("Hello", "ru", "en", "fr", "fa"))
}
