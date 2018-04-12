package main

import "currencyexchange/datahandling"
import "currencyexchange/webscraping"

func main() {
	webscraping.Webscraping1()
	webscraping.Webscraping2()
	datahandling.MergeCollideAndPrint()
}
