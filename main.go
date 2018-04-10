package main

import (
	"fmt"

	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	fmt.Println("Simple scrapper to find which credit card to use in a transaction (Visa or Mastercard)")
	//myTransactionCurrency := "EUR"
	//myCardCurrency := "USD"
	/*var myHeader = http.Header{}
	myHeader.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:59.0) Gecko/20100101 Firefox/59.0")
	myHeader.Add("Accept-Encoding", "gzip, deflate, br")
	myHeader.Add("Accept-Language", "en-US,en;q=0.5")
	myHeader.Add("Connection", "keep-alive")
	myHeader.Add("Content-Length", "0")
	myHeader.Add("Host", "chat.usa.visa.com")
	myHeader.Add("Content-Type", "text/plain;charset=UTF-8")
	myHeader.Add("Upgrade-Insecure-Requests", "1")
	*/
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0"),
		colly.CacheDir("./cache"),
	)
	var Titles []string
	var Prices []string

	//urlToVisit := fmt.Sprintf("http://usa.visa.com/support/consumer/travel-support/exchange-rate-calculator.html/?submitButton=Calculate+Exchange+Rates&fromCurr=%s&toCurr=%s&fee=0", myCardCurrency, myTransactionCurrency)
	//fmt.Printf("Visiting " + urlToVisit)
	urlToVisit := "http://nutifinanzas.com/"
	c.OnHTML("div.col-md-12.TituloDivisa", func(e *colly.HTMLElement) {
		//fmt.Println(strings.TrimSpace(e.Text))
		Titles = append(Titles, strings.TrimSpace(e.Text))
	})
	c.OnHTML("div.col-md-9", func(e *colly.HTMLElement) {
		re := regexp.MustCompile(`\r?\n`)
		TempPrice := strings.TrimSpace(e.Text)
		TempPrice = re.ReplaceAllString(TempPrice, " ")
		TempPrice = strings.Replace(TempPrice, " ", "", -1)
		TempPrice = strings.Replace(TempPrice, "Compramos:$", "", -1)
		TempPrice = strings.Replace(TempPrice, "Vendemos:$", "/", -1)
		Prices = append(Prices, TempPrice)
	})
	c.OnResponse(func(r *colly.Response) {
		//fmt.Printf(string(r.Body))
	})
	c.OnRequest(func(r *colly.Request) {
		//r.Ctx.Put("url", r.URL.String())
		//	fmt.Println(r.Headers)
	})
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Scrapped finished")
		for i, curval := range Titles {
			fmt.Printf("%s: %s\r\n", curval, Prices[i])
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("DETAILS: Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.Visit(urlToVisit)
	c.Wait()

}
