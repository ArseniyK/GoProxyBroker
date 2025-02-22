package types

import "fmt"

type Proxy struct {
	IP      string
	Port    int
	Type    []ProxyType
	IsAlive bool
	Level   ProxyLevel
}

func (proxy Proxy) String() string {
	var level = ""
	if proxy.Level > NONE {
		level = fmt.Sprintf(" A: %s", proxy.Level)
	}
	return fmt.Sprintf("%s:%d %s%s", proxy.IP, proxy.Port, proxy.Type, level)
}
