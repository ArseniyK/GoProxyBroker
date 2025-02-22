package main

import (
	"ProxyBroker/providers"
	"ProxyBroker/types"
	"fmt"
	"sync"
)

type Broker struct {
	publicIP  string
	countries map[string]struct{}
}

func (broker *Broker) Init(countries []string) {
	publicIP, err := GetPublicIP()
	if err != nil {
		fmt.Println("Error fetching public IP:", err)
		return
	}
	broker.publicIP = publicIP
	broker.countries = MakeSet(countries)

	fmt.Println("Public IP:", broker.publicIP)
	fmt.Println("Countries:", countries)
}

func (broker *Broker) Find(limit int, check bool) {
	var wg sync.WaitGroup
	checkWg := sync.WaitGroup{}
	proxyChan := make(chan types.Proxy, 10)

	for _, provider := range providers.Providers {
		wg.Add(1)
		go func(p providers.Provider) {
			defer wg.Done()
			proxies := p.GetProxies()

			for _, proxy := range proxies {
				checkWg.Add(1)
				go func(px types.Proxy) {
					defer checkWg.Done()

					if check {
						px = CheckProxy(px, broker.publicIP)
					}
					px.CountryCode = GetGeoIP(px.IP)
					if px.IsAlive && broker.checkCountry(px) {
						proxyChan <- px
					}
				}(proxy)
			}
		}(provider)
	}

	go func() {
		wg.Wait()
		checkWg.Wait()
		close(proxyChan)
	}()

	count := 0
	for proxy := range proxyChan {
		fmt.Println(proxy.String())
		count++

		if limit > 0 && count >= limit {
			break
		}
	}

	fmt.Printf("Found %d proxy\n", count)
}

func (broker *Broker) checkCountry(proxy types.Proxy) bool {
	if broker.countries == nil && len(broker.countries) == 0 {
		return true
	}
	if _, ok := broker.countries[proxy.CountryCode]; ok {
		return true
	} else {
		return false
	}
}
