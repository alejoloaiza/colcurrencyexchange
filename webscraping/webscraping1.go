package webscraping

import (
	"currencyexchange/datahandling"
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

func webscraping1() {
	//// TOMANDO LA INFO DE NUTIFINANZAS
	cNuti := colly.NewCollector(
	//colly.CacheDir("./cache"),
	)
	var Titles1 []string
	var Prices1 []string

	cNutiurlToVisit := "http://nutifinanzas.com/"
	cNuti.OnHTML("div.col-md-12.TituloDivisa", func(e *colly.HTMLElement) {
		//fmt.Println(strings.TrimSpace(e.Text))
		Titles1 = append(Titles1, strings.TrimSpace(e.Text))
	})
	cNuti.OnHTML("div.col-md-9", func(e *colly.HTMLElement) {
		re := regexp.MustCompile(`\r?\n`)
		TempPrice := strings.TrimSpace(e.Text)
		TempPrice = re.ReplaceAllString(TempPrice, " ")
		TempPrice = strings.Replace(TempPrice, " ", "", -1)
		TempPrice = strings.Replace(TempPrice, "Compramos:$", "", -1)
		TempPrice = strings.Replace(TempPrice, "Vendemos:$", "@", -1)
		Prices1 = append(Prices1, TempPrice)
	})
	cNuti.OnResponse(func(r *colly.Response) {
		//fmt.Printf(string(r.Body))
	})
	cNuti.OnRequest(func(r *colly.Request) {
		//r.Ctx.Put("url", r.URL.String())
		//	fmt.Println(r.Headers)
	})
	cNuti.OnScraped(func(r *colly.Response) {
		for i, curval := range Titles1 {
			var extData datahandling.ExchangeData
			extprice := strings.Split(Prices1[i], "@")
			extData.Name = "Nutifinanzas"
			extData.Currency = curval
			extData.PriceBuy = extprice[0]
			extData.PriceSell = extprice[1]
			//fmt.Printf("%s > %s > C: %s V: %s\r\n", extData.Name, extData.Currency, extData.PriceBuy, extData.PriceSell)
			datahandling.GeneralData = append(datahandling.GeneralData, extData)
		}
	})
	cNuti.OnError(func(r *colly.Response, err error) {
		fmt.Println("DETAILS: Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	cNuti.Visit(cNutiurlToVisit)
}
