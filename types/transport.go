package types

import "net/http"

type TransportWrapper struct {
	Transport http.RoundTripper
	Headers   map[string]string
}

func (t *TransportWrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original request
	newReq := req.Clone(req.Context())

	// Add custom headers
	for key, value := range t.Headers {
		newReq.Header.Set(key, value)
	}

	// Use the wrapped transport (default if nil)
	if t.Transport == nil {
		t.Transport = http.DefaultTransport
	}

	return t.Transport.RoundTrip(newReq)
}
