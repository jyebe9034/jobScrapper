package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string
	title    string
	location string
	salary   string
	summary  string
}

// Scrape indeed by a term
func Scrape(term string) {
	var baseURL string = "https://kr.indeed.com/jobs?q=" + term + "&limit=50"
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totlaPages := getPages(baseURL)

	for i := 0; i < totlaPages; i++ {
		go getPage(i, baseURL, c)
	}

	for i := 0; i < totlaPages; i++ {
		extractJobs := <-c
		jobs = append(jobs, extractJobs...)
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted ", len(jobs))
}

func getPage(page int, baseURL string, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Request: " + pageURL)
	res, err := http.Get(pageURL)
	checkError(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)

	searchCards := doc.Find(".resultWithShelf")

	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	mainC <- jobs
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("data-jk")
	title := cleanString(card.Find(".jobTitle>span").Text())
	location := cleanString(card.Find(".company_location div.companyLocation").Text())
	salary := cleanString(card.Find(".salary-snippet").Text())
	summary := cleanString(card.Find(".job-snippet").Text())
	c <- extractedJob{
		id:       id,
		title:    title,
		location: location,
		salary:   salary,
		summary:  summary,
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages(baseURL string) int {
	pages := 0
	res, err := http.Get(baseURL)
	checkError(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkError(err)

	w := csv.NewWriter(file)
	defer w.Flush() // ????????? ???????????? ???????????? ??????

	headers := []string{"Link", "Title", "Location", "Salary", "Summary"}

	wErr := w.Write(headers)
	checkError(wErr)

	// ???????????? go routine?????? ????????? ??????
	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkError(jwErr)
	}
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
