package main

import "flag"

func main() {
	limit := flag.Int("limit", 10, "The maximum number of proxies")
	check := flag.Bool("check", true, "Check found proxies")
	flag.Parse()

	b := Broker{}
	b.Init()
	b.Find(*limit, *check)
}
