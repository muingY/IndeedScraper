package main

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

//var baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"
var baseURL string = "https://kr.indeed.com/jobs?q="

func main() {
	var searchKey string
	fmt.Println("---< Indeed Job Scraper >---")
	fmt.Print("Search keyword: ")
	fmt.Scan(&searchKey)

	baseURL = baseURL + searchKey + "&limit=50"
	fmt.Println("---Start search---")

	totalPages := getPages()
	fmt.Println("Total pages:", totalPages)

	var jobs []extractedJob

	for i := 0; i < totalPages; i++ {
		extractedJobs := getPage(i)
		jobs = append(jobs, extractedJobs...)
		fmt.Printf("Page %d complete!\n", i+1)
	}

	writeJobs(jobs)
	fmt.Println("All Complete!")
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "TITLE", "LOCATION", "SALARY", "SUMMARY"}

	werr := w.Write(headers)
	checkErr(werr)

	for _, job := range jobs {
		jobSilce := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
		jobwerr := w.Write(jobSilce)
		checkErr(jobwerr)
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
		summary := cleanString(card.Find(".summary").Text())
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
	if pages == 0 {
		return 1
	}
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
