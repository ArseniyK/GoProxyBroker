package main

import (
	"flag"
	"strings"
)

func main() {
	limit := flag.Int("limit", 10, "The maximum number of proxies")
	check := flag.Bool("check", true, "Check found proxies")
	var countriesString string
	flag.StringVar(&countriesString, "countries", "", "List of comma separated ISO country codes where should be located proxies")
	flag.Parse()
	countries := strings.Split(countriesString, ",")
	b := Broker{}
	b.Init(countries)
	b.Find(*limit, *check)
}
