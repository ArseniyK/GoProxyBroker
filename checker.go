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

type CheckStatus struct {
	isAlive bool
	level   types.ProxyLevel
}

func CheckProxy(proxy types.Proxy, publicIp string) types.Proxy {
	newProxy := proxy
	newProxy.Type = []types.ProxyType{}

	results := make([]bool, len(proxy.Type))
	for _, t := range proxy.Type {
		if t == types.HTTP {
			status, _ := checkHTTP(proxy, publicIp)
			results = append(results, status.isAlive)
			if status.isAlive {
				newProxy.Level = status.level
				newProxy.Type = append(newProxy.Type, t)
			}
		}
		if t == types.HTTPS {
			status, _ := checkHTTPS(proxy)
			results = append(results, status.isAlive)
			if status.isAlive {
				newProxy.Type = append(newProxy.Type, t)
			}
		}
	}

	newProxy.IsAlive = Any(results)
	return newProxy
}

func checkLevel(body string, publicIp string) types.ProxyLevel {
	// If the public IP is missing, it’s at least anonymous
	if !strings.Contains(body, publicIp) {
		// If none of the proxy headers exist, it's a high anonymity proxy
		if !strings.Contains(body, "via") && !strings.Contains(body, "proxy") && !strings.Contains(body, "X-Forwarded-For") {
			return types.HIGH
		}
		return types.ANONYMOUS
	}
	return types.TRANSPARENT
}

func checkHTTP(proxy types.Proxy, publicIp string) (CheckStatus, error) {
	testURL := "http://httpbin.org/get?show_env"

	proxyURL, err := url.Parse(fmt.Sprintf("http://%s:%d", proxy.IP, proxy.Port))
	if err != nil {
		return CheckStatus{isAlive: false}, fmt.Errorf("invalid proxy URL: %v", err)
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
		return CheckStatus{isAlive: false}, fmt.Errorf("proxy check failed: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return CheckStatus{isAlive: false}, fmt.Errorf("proxy check failed: %v", err)
		}
		bodyString := string(body)

		if !strings.Contains(bodyString, testURL) {
			return CheckStatus{isAlive: false}, fmt.Errorf("proxy check failed: %v", err)
		}

		return CheckStatus{isAlive: true, level: checkLevel(bodyString, publicIp)}, nil
	}

	return CheckStatus{isAlive: false}, fmt.Errorf("proxy returned non-200 status: %d", response.StatusCode)

}

func checkHTTPS(proxy types.Proxy) (CheckStatus, error) {
	testURL := "https://httpbin.org/get?show_env"

	proxyURL, err := url.Parse(fmt.Sprintf("https://%s:%d", proxy.IP, proxy.Port))
	if err != nil {
		return CheckStatus{isAlive: false}, fmt.Errorf("invalid proxy URL: %v", err)
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
		return CheckStatus{isAlive: false}, fmt.Errorf("proxy check failed: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		_, err := io.ReadAll(response.Body)
		if err != nil {
			return CheckStatus{isAlive: false}, fmt.Errorf("proxy check failed: %v", err)
		}

		return CheckStatus{isAlive: true}, nil
	}

	return CheckStatus{isAlive: false}, fmt.Errorf("proxy returned non-200 status: %d", response.StatusCode)
}
