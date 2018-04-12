package webscraping

import (
	"currencyexchange/datahandling"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func webscraping2() {
	//// TOMANDO LA INFO DE UNICAMBIOS
	var Titles2 []string
	var Prices2 []string
	var Prices3 []string
	cUni := colly.NewCollector(
	//colly.CacheDir("./cache"),
	)

	cUniurlToVisit := "http://www.unicambios.com.co/"
	cUni.OnHTML("[class^=row_][class$=col_1]", func(e *colly.HTMLElement) {
		if strings.TrimSpace(e.Text) != "" {
			Titles2 = append(Titles2, strings.TrimSpace(e.Text))
		}
	})
	cUni.OnHTML("[class^=row_][class$=col_2]", func(e *colly.HTMLElement) {
		//fmt.Println(e.Attr("class"))
		if strings.Contains(e.Attr("class"), "row_0 ") || strings.Contains(e.Attr("class"), "row_1 ") {
			return
		} else {
			if strings.TrimSpace(e.Text) != "" {
				Prices2 = append(Prices2, strings.TrimSpace(e.Text))
			}
		}
	})
	cUni.OnHTML("[class^=row_][class$=col_3]", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "row_0 ") || strings.Contains(e.Attr("class"), "row_1 ") {
			return
		} else {
			if strings.TrimSpace(e.Text) != "" {
				Prices3 = append(Prices3, strings.TrimSpace(e.Text))
			}
		}

	})
	cUni.OnResponse(func(r *colly.Response) {
		//fmt.Printf(string(r.Body))
	})
	cUni.OnRequest(func(r *colly.Request) {
		//r.Ctx.Put("url", r.URL.String())
		//	fmt.Println(r.Headers)
	})

	cUni.OnScraped(func(r *colly.Response) {
		for i, curval := range Titles2 {
			var extData datahandling.ExchangeData
			extData.Name = "Unicambios"
			extData.Currency = curval
			extData.PriceBuy = Prices2[i]
			extData.PriceSell = Prices3[i]
			//fmt.Printf("%s > %s > C: %s V: %s\r\n", extData.Name, extData.Currency, extData.PriceBuy, extData.PriceSell)
			datahandling.GeneralData = append(datahandling.GeneralData, extData)
		}
	})
	cUni.OnError(func(r *colly.Response, err error) {
		fmt.Println("DETAILS: Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	cUni.Visit(cUniurlToVisit)

}
