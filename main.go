package main

import (
	"ProxyBroker/types"
	"flag"
	"fmt"
	"strings"
)

var levelMap = map[string]types.ProxyLevel{
	"transparent": types.TRANSPARENT,
	"anonymous":   types.ANONYMOUS,
	"high":        types.HIGH,
}

func main() {
	limit := flag.Int("limit", 10, "The maximum number of proxies")
	check := flag.Bool("check", true, "Check found proxies")
	var countriesString string
	flag.StringVar(&countriesString, "countries", "", "List of comma separated ISO country codes where should be located proxies")

	var levels []types.ProxyLevel
	flag.Func("lvl", "Comma-separated proxy levels (transparent, anonymous, high)", func(s string) error {
		values := strings.Split(s, ",")
		for _, v := range values {
			v = strings.TrimSpace(v) // Remove spaces
			level, exists := levelMap[v]
			if !exists {
				return fmt.Errorf("invalid level: %s (allowed: transparent, anonymous, high)", v)
			}
			levels = append(levels, level)
		}
		return nil
	})

	flag.Parse()
	countries := strings.Split(countriesString, ",")

	b := Broker{}
	b.Init(countries, levels)
	b.Find(*limit, *check)
}
