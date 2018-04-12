package main

import "currencyexchange/datahandling"
import "currencyexchange/webscraping"

func main() {
	webscraping.webscraping1()
	webscraping.webscraping2()
	datahandling.MergeCollideAndPrint()
}
