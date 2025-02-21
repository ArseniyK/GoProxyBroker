package main

import (
	"ProxyBroker/types"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func CheckProxy(proxy types.Proxy, publicIp string) bool {
	results := make([]bool, len(proxy.Type))

	for _, t := range proxy.Type {
		if t == types.HTTP {
			isAlive, _ := checkHTTP(proxy, publicIp)
			results = append(results, isAlive)
		}
	}

	return Any(results)
}

func checkHTTP(proxy types.Proxy, publicIp string) (bool, error) {
	testURL := "http://httpbin.org/get?show_env"

	proxyURL, err := url.Parse(fmt.Sprintf("http://%s:%d", proxy.IP, proxy.Port))
	if err != nil {
		return false, fmt.Errorf("invalid proxy URL: %v", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		},
		Timeout: 60 * time.Second,
	}

	response, err := client.Get(testURL)
	if err != nil {
		return false, fmt.Errorf("proxy check failed: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err == nil {
			fmt.Println(string(body))
			if !strings.Contains(string(body), publicIp) {
				proxy.Level = types.ANONYMOUS
			}
		}

		return true, nil
	}

	return false, fmt.Errorf("proxy returned non-200 status: %d", response.StatusCode)

}

func checkHTTPS(proxy types.Proxy) (bool, error) {
	return false, nil
}
