package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/satinath-nit/google-translate/cli"
)

var sourceLang string
var targetLang string
var sourceText string
var wg sync.WaitGroup

func init() {
	flag.StringVar(&sourceLang, "s", "en", "source language")
	flag.StringVar(&targetLang, "t", "zh", "target language")
	flag.StringVar(&sourceText, "str", "", "source text")
}

func main() {
	flag.Parse()
	if flag.NFlag() == 0 {
		fmt.Println("No arguments provided.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	strChan := make(chan string)
	wg.Add(1)

	reqBody := &cli.RequestBody{
		SourceLang: sourceLang,
		TargetLang: targetLang,
		SourceText: sourceText,
	}

	go cli.RequestTranslate(reqBody, strChan, &wg)

	processedStr := strings.ReplaceAll(<-strChan, "+", " ")
	unquotedStr, err := json.RawUnquote(processedStr)
	if err != nil {

	}
	fmt.Print(" Translated Text is --->", processedStr)
	close(strChan)
	wg.Wait()

}
