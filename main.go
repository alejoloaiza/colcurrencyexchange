package main

import (
	"fmt"

	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

func main() {

	cNuti := colly.NewCollector(
		colly.CacheDir("./cache"),
	)
	var Titles []string
	var Prices []string

	//urlToVisit := fmt.Sprintf("http://usa.visa.com/support/consumer/travel-support/exchange-rate-calculator.html/?submitButton=Calculate+Exchange+Rates&fromCurr=%s&toCurr=%s&fee=0", myCardCurrency, myTransactionCurrency)
	//fmt.Printf("Visiting " + urlToVisit)
	cNutiurlToVisit := "http://nutifinanzas.com/"
	cNuti.OnHTML("div.col-md-12.TituloDivisa", func(e *colly.HTMLElement) {
		//fmt.Println(strings.TrimSpace(e.Text))
		Titles = append(Titles, strings.TrimSpace(e.Text))
	})
	cNuti.OnHTML("div.col-md-9", func(e *colly.HTMLElement) {
		re := regexp.MustCompile(`\r?\n`)
		TempPrice := strings.TrimSpace(e.Text)
		TempPrice = re.ReplaceAllString(TempPrice, " ")
		TempPrice = strings.Replace(TempPrice, " ", "", -1)
		TempPrice = strings.Replace(TempPrice, "Compramos:$", "", -1)
		TempPrice = strings.Replace(TempPrice, "Vendemos:$", "/", -1)
		Prices = append(Prices, TempPrice)
	})
	cNuti.OnResponse(func(r *colly.Response) {
		//fmt.Printf(string(r.Body))
	})
	cNuti.OnRequest(func(r *colly.Request) {
		//r.Ctx.Put("url", r.URL.String())
		//	fmt.Println(r.Headers)
	})
	cNuti.OnScraped(func(r *colly.Response) {
		for i, curval := range Titles {
			fmt.Printf("%s: %s\r\n", curval, Prices[i])
		}
	})
	cNuti.OnError(func(r *colly.Response, err error) {
		fmt.Println("DETAILS: Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	cNuti.Visit(cNutiurlToVisit)
	cNuti.Wait()

}
