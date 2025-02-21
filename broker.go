package main

import (
	"ProxyBroker/providers"
	"ProxyBroker/types"
	"fmt"
	"sync"
)

type Broker struct {
	publicIP string
}

func (broker *Broker) Init() {
	publicIP, err := getPublicIP()
	if err != nil {
		fmt.Println("Error fetching public IP:", err)
		return
	}
	fmt.Println("Public IP:", publicIP)
	broker.publicIP = publicIP
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
					if !check || CheckProxy(px, broker.publicIP) {
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
