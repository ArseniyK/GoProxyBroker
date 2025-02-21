package main

import "flag"

func main() {
	limit := flag.Int("timeout", 10, "The maximum number of proxies")
	flag.Parse()

	b := Broker{}
	b.Init()
	b.Find(*limit, true)
}
