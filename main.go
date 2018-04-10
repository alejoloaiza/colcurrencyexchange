package main

import (
	"fmt"

	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

type ExchangeData struct {
	Name      string
	Currency  string
	PriceBuy  string
	PriceSell string
}

func main() {
	var generalData []ExchangeData

	cNuti := colly.NewCollector(
		colly.CacheDir("./cache"),
	)
	var Titles []string
	var Prices []string

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
		TempPrice = strings.Replace(TempPrice, "Vendemos:$", "@", -1)
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
			var extData ExchangeData
			extprice := strings.Split(Prices[i], "@")
			extData.Name = "Nutifinanzas"
			extData.Currency = curval
			extData.PriceBuy = extprice[0]
			extData.PriceSell = extprice[1]
			fmt.Printf("%s: C: %s V: %s\r\n", extData.Currency, extData.PriceBuy, extData.PriceSell)
			generalData = append(generalData, extData)
		}
	})
	cNuti.OnError(func(r *colly.Response, err error) {
		fmt.Println("DETAILS: Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	cNuti.Visit(cNutiurlToVisit)
	cNuti.Wait()

}
