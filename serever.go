package main

import (
	"ProxyBroker/types"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type ProxyPoolInterface interface {
	Get() (*types.Proxy, error) // Retrieve a proxy in round-robin fashion.
	Put(proxy types.Proxy)      // Add a proxy back to the pool.
}

type ProxyPool struct {
	mu           sync.Mutex
	proxies      []types.Proxy
	currentIndex int
}

func NewProxyPool(initialProxies []types.Proxy) *ProxyPool {
	return &ProxyPool{
		proxies:      initialProxies,
		currentIndex: 0,
	}
}

func (p *ProxyPool) Get() (*types.Proxy, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.proxies) == 0 {
		return nil, errors.New("proxy pool is empty")
	}

	// Retrieve the current proxy in round-robin order.
	proxy := &p.proxies[p.currentIndex]

	// Increment the index and wrap around if necessary (round-robin).
	p.currentIndex = (p.currentIndex + 1) % len(p.proxies)

	return proxy, nil
}

// Put adds a new proxy to the pool.
func (p *ProxyPool) Put(proxy types.Proxy) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Append the proxy to the pool.
	p.proxies = append(p.proxies, proxy)
}

type ProxyServer struct {
	pool *ProxyPool
}

func (ps *ProxyServer) Init(pool *ProxyPool) {
	ps.pool = pool
}

func (ps *ProxyServer) Start(host string, port int) error {
	fmt.Printf("Serve %s:%d\n", host, port)

	handler := func(w http.ResponseWriter, r *http.Request) {
		// Fetch a proxy via the broker
		proxy, err := ps.pool.Get()
		if err != nil {
			http.Error(w, "No available proxies", http.StatusServiceUnavailable)
			return
		}

		// Forward traffic through the selected proxy
		ps.forwardRequest(w, r, *proxy)
	}

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: http.HandlerFunc(handler),
	}

	return server.ListenAndServe()
}

func (ps *ProxyServer) forwardRequest(w http.ResponseWriter, r *http.Request, proxy types.Proxy) {
	// Set up the proxy URL
	proxyURL, err := url.Parse(fmt.Sprintf("http://%s:%d", proxy.IP, proxy.Port))
	if err != nil {
		http.Error(w, "Invalid proxy", http.StatusInternalServerError)
		return
	}

	// Use the ReverseProxy to forward the request
	proxyHandler := httputil.NewSingleHostReverseProxy(proxyURL)
	proxyHandler.ServeHTTP(w, r)
}
