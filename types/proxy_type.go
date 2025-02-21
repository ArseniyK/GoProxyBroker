package types

type ProxyType int

const (
	HTTP ProxyType = iota
	HTTPS
	SOCKS
)

func (p ProxyType) String() string {
	return [...]string{"HTTP", "HTTPS", "SOCKS"}[p]
}
