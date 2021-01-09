package main

import (
	"fmt"
	"log"
	"net/http"
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

var baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

func main() {
	totalPages := getPages()
	fmt.Println("pages:", totalPages)

	var jobs []extractedJob

	for i := 0; i < totalPages; i++ {
		extractedJobs := getPage(i)
		jobs = append(jobs, extractedJobs...)
	}
}

func getPage(page int) []extractedJob {
	var jobs []extractedJob

	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")

	searchCards.Each(func(i int, card *goquery.Selection) {
		id, _ := card.Attr("data-jk")
		title := cleanString(card.Find(".title>a").Text())
		location := cleanString(card.Find(".sjcl").Text())
		salary := cleanString(card.Find(".salaryText").Text())
		summary := cleanString(card.Find(".sammary").Text())
		job := extractedJob{
			id:       id,
			title:    title,
			location: location,
			salary:   salary,
			summary:  summary}
		jobs = append(jobs, job)
	})
	return jobs
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages() int {
	var pages int

	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}
