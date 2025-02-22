package providers

import (
	"ProxyBroker/types"
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var ipPortPattern = regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}):(\d{2,5})`)

var globalHeaders = map[string]string{
	"User-Agent":    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 GPxBroker/1.0",
	"Accept":        "*/*",
	"Pragma":        "no-cache",
	"Cache-control": "no-cache",
	"Cookie":        "cookie=ok",
	"Referer":       "https://www.google.com/",
}

type Provider struct {
	URL       string
	ProxyType []types.ProxyType
	Timeout   time.Duration
}

func (provider *Provider) GetProxies() <-chan types.Proxy {
	client := &http.Client{Timeout: provider.Timeout, Transport: &types.TransportWrapper{Headers: globalHeaders}}

	ch := make(chan types.Proxy)
	go func() {
		defer close(ch) // Close channel when done
		proxies, err := provider.fetchProxies(client)

		if err != nil {
			log.Printf("Error fetching proxies from %s: %v", provider.URL, err)
			return
		}

		for _, proxy := range proxies {
			ch <- proxy
		}

	}()
	return ch

}

func (provider *Provider) fetchProxies(client *http.Client) ([]types.Proxy, error) {
	resp, err := client.Get(provider.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("bad status: " + resp.Status)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		return nil, err
	}

	body := buf.Bytes()

	return provider.findProxies(string(body)), nil
}

func (provider *Provider) findProxies(page string) []types.Proxy {
	matches := ipPortPattern.FindAllStringSubmatch(page, -1)
	var proxies []types.Proxy

	for _, match := range matches {

		if len(match) < 2 {
			continue
		}

		proxies = append(proxies, types.Proxy{
			IP:      match[1],
			Port:    toInt(match[2]),
			Type:    provider.ProxyType,
			IsAlive: true,
		})
	}

	return proxies
}

func toInt(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return num
}
