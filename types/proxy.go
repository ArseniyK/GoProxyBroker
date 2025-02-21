package types

import "fmt"

type Proxy struct {
	IP    string
	Port  int
	Type  []ProxyType
	Level ProxyLevel
}

func (proxy Proxy) String() string {
	return fmt.Sprintf("%s:%d %s A: %s", proxy.IP, proxy.Port, proxy.Type, proxy.Level)
}
