package providers

import (
	"ProxyBroker/types"
	"time"
)

var Providers = []Provider{
	{
		URL: "https://www.proxy-list.download/api/v1/get?type=http",
		ProxyType: []types.ProxyType{
			types.HTTP,
		},
		Timeout: 30 * time.Second,
	},
	{
		URL: "https://api.proxyscrape.com/?request=getproxies&proxytype=http",
		ProxyType: []types.ProxyType{
			types.HTTP,
		},
		Timeout: 30 * time.Second,
	},
	{
		URL: "https://proxy-daily.com/",
		ProxyType: []types.ProxyType{
			types.HTTP, types.HTTPS,
		},
		Timeout: 30 * time.Second,
	},
	{
		URL: "http://pubproxy.com/api/proxy?limit=20&format=txt",
		ProxyType: []types.ProxyType{
			types.HTTP, types.HTTPS,
		},
		Timeout: 30 * time.Second,
	},
	{
		URL: "https://api.proxyscrape.com/?request=getproxies&proxytype=http",
		ProxyType: []types.ProxyType{
			types.HTTP, types.HTTPS,
		},
		Timeout: 30 * time.Second,
	},
	{
		URL: "https://raw.githubusercontent.com/fate0/proxylist/master/proxy.list",
		ProxyType: []types.ProxyType{
			types.HTTP, types.HTTPS,
		},
		Timeout: 30 * time.Second,
	},
	{
		URL: "https://raw.githubusercontent.com/sunny9577/proxy-scraper/master/proxies.json",
		ProxyType: []types.ProxyType{
			types.HTTP, types.HTTPS,
		},
		Timeout: 30 * time.Second,
	},
	{
		URL: "https://github.com/TheSpeedX/PROXY-List/blob/master/http.txt",
		ProxyType: []types.ProxyType{
			types.HTTP, types.HTTPS,
		},
		Timeout: 30 * time.Second,
	},
	{
		URL: "https://raw.githubusercontent.com/clarketm/proxy-list/master/proxy-list.txt",
		ProxyType: []types.ProxyType{
			types.HTTP, types.HTTPS,
		},
		Timeout: 30 * time.Second,
	},
}
