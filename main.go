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
			var extData ExchangeData
			extprice := strings.Split(Prices1[i], "@")
			extData.Name = "Nutifinanzas"
			extData.Currency = curval
			extData.PriceBuy = extprice[0]
			extData.PriceSell = extprice[1]
			fmt.Printf("%s > %s > C: %s V: %s\r\n", extData.Name, extData.Currency, extData.PriceBuy, extData.PriceSell)
			generalData = append(generalData, extData)
		}
	})
	cNuti.OnError(func(r *colly.Response, err error) {
		fmt.Println("DETAILS: Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	cNuti.Visit(cNutiurlToVisit)
	cNuti.Wait()

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
			var extData ExchangeData
			extData.Name = "Unicambios"
			extData.Currency = curval
			extData.PriceBuy = Prices2[i]
			extData.PriceSell = Prices3[i]
			fmt.Printf("%s > %s > C: %s V: %s\r\n", extData.Name, extData.Currency, extData.PriceBuy, extData.PriceSell)
			generalData = append(generalData, extData)
		}
	})
	cUni.OnError(func(r *colly.Response, err error) {
		fmt.Println("DETAILS: Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	cUni.Visit(cUniurlToVisit)
	cUni.Wait()

}
