package main

import (
	"log"
	"net/http"
	// "github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://kr.indeed.com/jobs?q=springboot&limit=50"

func main() {
	getPages()
}

func getPages() int {
	res, err := http.Get(baseURL)
	checkError(err)
	checkCode(res)

	defer res.Body.Close()

	// doc, err := goquery.NewDocumentFromReader(res.Body)
	// checkError(err)

	// fmt.Println(doc)

	return 0
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("request failed with statusCode")
	}
}
