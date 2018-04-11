package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/schollz/closestmatch"
)

type ExchangeData struct {
	Name      string
	Currency  string
	PriceBuy  string
	PriceSell string
}

var generalData []ExchangeData

func main() {

	//// TOMANDO LA INFO DE NUTIFINANZAS
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
			//fmt.Printf("%s > %s > C: %s V: %s\r\n", extData.Name, extData.Currency, extData.PriceBuy, extData.PriceSell)
			generalData = append(generalData, extData)
		}
	})
	cNuti.OnError(func(r *colly.Response, err error) {
		fmt.Println("DETAILS: Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	cNuti.Visit(cNutiurlToVisit)

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
			var extData ExchangeData
			extData.Name = "Unicambios"
			extData.Currency = curval
			extData.PriceBuy = Prices2[i]
			extData.PriceSell = Prices3[i]
			//fmt.Printf("%s > %s > C: %s V: %s\r\n", extData.Name, extData.Currency, extData.PriceBuy, extData.PriceSell)
			generalData = append(generalData, extData)
		}
	})
	cUni.OnError(func(r *colly.Response, err error) {
		fmt.Println("DETAILS: Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	cUni.Visit(cUniurlToVisit)

	MergeCollideAndPrint()
}
func MergeCollideAndPrint() {

	// NORMALIZATION
	Currency1Stuff := []string{"Euro Baja EUR", "Colon Costarricense CRC", "Rupia India INR", "Quetzal GTQ", "Peso Uruguayo UYU", "Peso Dominicano DOP", "Peso Boliviano BOP", "Lira Turca TRY", "Florin Antillas AWG", "Dolar Nueva Zelanda NZD", "Corona Sueca SEK", "Corona Noruega NOK", "Colon CRC", "Bolivar Fuerte VEF", "Nuevo Sol PEN", "Peso Mexicano MXN", "Yen japones JPY", "Libra Esterlina GBP", "Yuan Chino CNY", "Peso Chileno CLP", "Franco Suizo CHF", "Dolar canadiense CAD", "US Dolar(Cheque Viajero)", "Dolar Americano USD", "Euro EUR 500 y 200", "Euro EUR", "Peso Argentino ARS", "Dolar Australiano AUD", "Real Brasil BRL"}
	bagSizes := []int{2, 3, 4, 5, 6, 7, 8}
	Currency1Type := closestmatch.New(Currency1Stuff, bagSizes)
	//fmt.Println(cmType.Closest(curAsset.Type))
	for i, curCurrency := range generalData {
		curCurrency.Currency = Currency1Type.Closest(curCurrency.Currency)
		r, _ := regexp.Compile(".([0-9])([0-9])([0-9])")
		generalData[i].Currency = curCurrency.Currency
		curCurrency.PriceBuy = strings.Replace(strings.Replace(strings.Replace(curCurrency.PriceBuy, ",", "", -1), ".00", "", -1), "$", "", -1)
		curCurrency.PriceSell = strings.Replace(strings.Replace(strings.Replace(curCurrency.PriceSell, ",", "", -1), ".00", "", -1), "$", "", -1)
		if r.MatchString(curCurrency.PriceBuy) {
			curCurrency.PriceBuy = strings.Replace(curCurrency.PriceBuy, ".", "", -1)
		}
		if r.MatchString(curCurrency.PriceBuy) {
			curCurrency.PriceSell = strings.Replace(curCurrency.PriceSell, ".", "", -1)
		}
		generalData[i].PriceBuy = curCurrency.PriceBuy
		generalData[i].PriceSell = curCurrency.PriceSell
		generalData[i].Currency = curCurrency.Currency

		//	fmt.Printf("%s > %s > C: %s V: %s\r\n", curCurrency.Name, curCurrency.Currency, curCurrency.PriceBuy, curCurrency.PriceSell)
	}
	Currency2Stuff := []string{"Colon Costarricense CRC", "Rupia India INR", "Quetzal GTQ", "Peso Uruguayo UYU", "Peso Dominicano DOP", "Peso Boliviano BOP", "Lira Turca TRY", "Florin Antillas AWG", "Dolar Nueva Zelanda NZD", "Corona Sueca SEK", "Corona Noruega NOK", "Colon CRC", "Bolivar Fuerte VEF", "Nuevo Sol PEN", "Peso Mexicano MXN", "Yen japones JPY", "Libra Esterlina GBP", "Yuan Chino CNY", "Peso Chileno CLP", "Franco Suizo CHF", "Dolar canadiense CAD", "US Dolar(Cheque Viajero)", "Dolar Americano USD", "Euro EUR 500 y 200", "Euro EUR", "Peso Argentino ARS", "Dolar Australiano AUD", "Real Brasil BRL"}
	Currency2Type := closestmatch.New(Currency2Stuff, bagSizes)
	for i, curCurrency := range generalData {
		curCurrency.Currency = Currency2Type.Closest(curCurrency.Currency)
		generalData[i].Currency = curCurrency.Currency
		fmt.Printf("%s > %s > C: %s V: %s\r\n", curCurrency.Name, curCurrency.Currency, curCurrency.PriceBuy, curCurrency.PriceSell)
	}

	// MERGE COLLIDE
	i := 0
	j := 0
	for i < len(generalData) {
		for j < len(generalData) {
			if generalData[i].Name != generalData[j].Name && generalData[i].Currency == generalData[j].Currency {

			}

			j = j + 1

		}
		i = i + 1
	}
}
