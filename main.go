package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func check(error error) {
	if error != nil {
		fmt.Println(error)
	}
}

func getHtml(url string) *http.Response {
	response, error := http.Get(url)
	check(error)

	if response.StatusCode > 400 {
		fmt.Println("Status code:", response.StatusCode)
	}

	return response
}

func writeCsv(scrapedData []string) {
	filename := "data.csv"

	file, error := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	check(error)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	error = writer.Write(scrapedData)
	check(error)
}

func scrapePageData(doc *goquery.Document) {
	doc.Find("ul.srp-results>li.s-item").Each(func(index int, item *goquery.Selection) {
		a := item.Find("a.s-item_link")

		title := strings.TrimSpace(a.Text())
		url, _ := a.Attr("href")

		price_span := item.Find("span.s-item_price").Text()
		price := strings.Trim(price_span, "py6.")

		scrapedData := []string{title, price, url}

		writeCsv(scrapedData)
	})
}

func main() {
	url := "https://www.ebay.com/globaldeals?_trkparms=pageci%3A2e1adb5a-93bc-11ed-9a3b-6adecd69cf1b%7Cparentrq%3Aae56b2211850ab965cbff6f8fffe29e8%7Ciid%3A1"

	var previousUrl string

	for {
		response := getHtml(url)
		defer response.Body.Close()

		doc, error := goquery.NewDocumentFromReader(response.Body)
		check(error)

		scrapePageData(doc)

		href, _ := doc.Find("nav.pagination>a.pagination_next").Attr("href")

		if href == previousUrl {
			break
		} else {
			url = href
			previousUrl = href
		}
	}

}
