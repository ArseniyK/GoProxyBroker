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
	levels    []types.ProxyLevel
}

func (broker *Broker) checkCountry(proxy types.Proxy) bool {
	if broker.countries == nil || len(broker.countries) == 0 {
		return true
	}
	_, ok := broker.countries[proxy.CountryCode]
	return ok
}

func (broker *Broker) checkLevels(proxy types.Proxy) bool {
	if broker.levels == nil || len(broker.levels) == 0 {
		return true
	}

	if proxy.Level == types.NONE {
		return true
	}

	for _, v := range broker.levels {
		if v == proxy.Level {
			return true
		}
	}
	return false
}

func (broker *Broker) Init(countries []string, levels []types.ProxyLevel) {
	publicIP, err := GetPublicIP()
	if err != nil {
		fmt.Println("Error fetching public IP:", err)
		return
	}
	broker.publicIP = publicIP
	broker.countries = MakeSet(countries)
	broker.levels = levels

	fmt.Println("Public IP:", broker.publicIP)
	fmt.Println("Countries:", countries)
	fmt.Println("Levels:", levels)
}

func (broker *Broker) grab(proxyChan chan types.Proxy) {
	var wg sync.WaitGroup
	checkWg := sync.WaitGroup{}
	input, output := distinct()

	for _, provider := range providers.Providers {
		wg.Add(1)
		go func(p providers.Provider) {
			defer wg.Done()
			for proxy := range p.GetProxies() {
				input <- proxy
			}
		}(provider)
	}

	go func() {
		for proxy := range output {
			checkWg.Add(1)
			go func(px types.Proxy) {
				defer checkWg.Done()
				px = CheckProxy(px, broker.publicIP)
				px.CountryCode = GetGeoIP(px.IP)
				if px.IsAlive && broker.checkCountry(px) && broker.checkLevels(px) {
					proxyChan <- px
				}
			}(proxy)
		}
	}()

	go func() {
		wg.Wait()
		checkWg.Wait()
		close(input)
		close(proxyChan)
	}()
}

func (broker *Broker) Find(limit int, check bool) {
	proxyChan := make(chan types.Proxy, 100)

	go broker.grab(proxyChan)

	count := 0
	for proxy := range proxyChan {
		fmt.Println(proxy)
		count++

		if limit > 0 && count >= limit {
			break
		}
	}

	fmt.Printf("Found %d proxy\n", count)
}

func (broker *Broker) Serve(host string, port int) {
	proxyChan := make(chan types.Proxy, 100)
	pool := NewProxyPool([]types.Proxy{})
	go broker.grab(proxyChan)

	go func() {
		for proxy := range proxyChan {
			pool.Put(proxy)
		}
	}()

	proxyServer := ProxyServer{}
	proxyServer.Init(pool)
	proxyServer.Start(host, port)
}
