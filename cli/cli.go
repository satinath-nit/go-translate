package cli

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/Jeffail/gabs"
)

type RequestBody struct {
	SourceLang string
	TargetLang string
	SourceText string
}

const translateURL = "https://translate.googleapis.com/translate_a/single"

func RequestTranslate(reqBody *RequestBody, strChan chan string, wg *sync.WaitGroup) {
	fmt.Println("Requesting translation...")
	fmt.Println(reqBody.SourceLang, reqBody.TargetLang, reqBody.SourceText)
	client := &http.Client{}
	req, err := http.NewRequest("GET", translateURL, nil)
	if err != nil {
		fmt.Println("Error creating request.")
		fmt.Println(err)
		os.Exit(1)
	}

	query := req.URL.Query()
	query.Add("client", "gtx")
	query.Add("sl", reqBody.SourceLang)
	query.Add("tl", reqBody.TargetLang)
	query.Add("dt", "t")
	query.Add("q", reqBody.SourceText)
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println(resp)

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		strChan <- "Too many requests. Try again later."
		wg.Done()
		return
	}

	parsedJson, err := gabs.ParseJSONBuffer(resp.Body)
	fmt.Println("Parsed Response from Google Translator-->", parsedJson)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	translatedText := parsedJson.Index(0).Index(0).Index(0).String()

	fmt.Println(translatedText)
	strChan <- translatedText
	wg.Done()

}
