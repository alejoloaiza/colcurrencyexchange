package datahandling

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/schollz/closestmatch"
)

type ExchangeData struct {
	Name      string
	Currency  string
	PriceBuy  string
	PriceSell string
}

var GeneralData []ExchangeData

func MergeCollideAndPrint() {

	// NORMALIZATION
	Currency1Stuff := []string{"Euro Baja EUR", "Colon Costarricense CRC", "Rupia India INR", "Quetzal GTQ", "Peso Uruguayo UYU", "Peso Dominicano DOP", "Peso Boliviano BOP", "Lira Turca TRY", "Florin Antillas AWG", "Dolar Nueva Zelanda NZD", "Corona Sueca SEK", "Corona Noruega NOK", "Colon CRC", "Bolivar Fuerte VEF", "Nuevo Sol PEN", "Peso Mexicano MXN", "Yen japones JPY", "Libra Esterlina GBP", "Yuan Chino CNY", "Peso Chileno CLP", "Franco Suizo CHF", "Dolar canadiense CAD", "US Dolar(Cheque Viajero)", "Dolar Americano USD", "Euro EUR 500 y 200", "Euro EUR", "Peso Argentino ARS", "Dolar Australiano AUD", "Real Brasil BRL"}
	bagSizes := []int{2, 3, 4, 5, 6, 7, 8}
	Currency1Type := closestmatch.New(Currency1Stuff, bagSizes)
	//fmt.Println(cmType.Closest(curAsset.Type))
	for i, curCurrency := range GeneralData {
		curCurrency.Currency = Currency1Type.Closest(curCurrency.Currency)
		r, _ := regexp.Compile(".([0-9])([0-9])([0-9])")
		GeneralData[i].Currency = curCurrency.Currency
		curCurrency.PriceBuy = strings.Replace(strings.Replace(strings.Replace(curCurrency.PriceBuy, ",", "", -1), ".00", "", -1), "$", "", -1)
		curCurrency.PriceSell = strings.Replace(strings.Replace(strings.Replace(curCurrency.PriceSell, ",", "", -1), ".00", "", -1), "$", "", -1)
		if r.MatchString(curCurrency.PriceBuy) {
			curCurrency.PriceBuy = strings.Replace(curCurrency.PriceBuy, ".", "", -1)
		}
		if r.MatchString(curCurrency.PriceBuy) {
			curCurrency.PriceSell = strings.Replace(curCurrency.PriceSell, ".", "", -1)
		}
		GeneralData[i].PriceBuy = curCurrency.PriceBuy
		GeneralData[i].PriceSell = curCurrency.PriceSell
		GeneralData[i].Currency = curCurrency.Currency

		//	fmt.Printf("%s > %s > C: %s V: %s\r\n", curCurrency.Name, curCurrency.Currency, curCurrency.PriceBuy, curCurrency.PriceSell)
	}
	Currency2Stuff := []string{"Colon Costarricense CRC", "Rupia India INR", "Quetzal GTQ", "Peso Uruguayo UYU", "Peso Dominicano DOP", "Peso Boliviano BOP", "Lira Turca TRY", "Florin Antillas AWG", "Dolar Nueva Zelanda NZD", "Corona Sueca SEK", "Corona Noruega NOK", "Colon CRC", "Bolivar Fuerte VEF", "Nuevo Sol PEN", "Peso Mexicano MXN", "Yen japones JPY", "Libra Esterlina GBP", "Yuan Chino CNY", "Peso Chileno CLP", "Franco Suizo CHF", "Dolar canadiense CAD", "US Dolar(Cheque Viajero)", "Dolar Americano USD", "Euro EUR 500 y 200", "Euro EUR", "Peso Argentino ARS", "Dolar Australiano AUD", "Real Brasil BRL"}
	Currency2Type := closestmatch.New(Currency2Stuff, bagSizes)
	for i, curCurrency := range GeneralData {
		curCurrency.Currency = Currency2Type.Closest(curCurrency.Currency)
		GeneralData[i].Currency = curCurrency.Currency
		//fmt.Printf("%s > %s > C: %s V: %s\r\n", curCurrency.Name, curCurrency.Currency, curCurrency.PriceBuy, curCurrency.PriceSell)
	}

	// MERGE COLLIDE
	i := 0
	for i < len(GeneralData) {
		j := 0
		for j < len(GeneralData) {
			if GeneralData[i].Name != GeneralData[j].Name && GeneralData[i].Currency == GeneralData[j].Currency {
				tmpBuy1, _ := strconv.ParseFloat(GeneralData[i].PriceBuy, 64)
				tmpBuy2, _ := strconv.ParseFloat(GeneralData[j].PriceBuy, 64)
				fmt.Printf("Moneda: " + GeneralData[i].Currency + " \n")
				if tmpBuy1 > tmpBuy2 {
					fmt.Printf(" >COMPRA: Mejor precio de compra: " + GeneralData[i].Name + "\n")
				} else if tmpBuy1 == tmpBuy2 {
					fmt.Printf(" >COMPRA: Los precios son iguales \n")
				} else {
					fmt.Printf(" >COMPRA: Mejor precio de compra: " + GeneralData[j].Name + "\n")
				}
				tmpSell1, _ := strconv.ParseFloat(GeneralData[i].PriceSell, 64)
				tmpSell2, _ := strconv.ParseFloat(GeneralData[j].PriceSell, 64)
				if tmpSell1 < tmpSell2 {
					fmt.Printf(" >VENTA: Mejor precio de venta: " + GeneralData[i].Name + "\n")
				} else if tmpSell1 == tmpSell2 {
					fmt.Printf(" >VENTA: Los precios son iguales \n")
				} else {
					fmt.Printf(" >VENTA: Mejor precio de venta: " + GeneralData[j].Name + "\n")
				}
			}

			j = j + 1

		}
		i = i + 1
	}
}
