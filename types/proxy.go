package types

import "fmt"

type Proxy struct {
	IP          string
	Port        int
	Type        []ProxyType
	IsAlive     bool
	Level       ProxyLevel
	CountryCode string
}

func (proxy Proxy) String() string {
	var level = ""
	if proxy.Level > NONE {
		level = fmt.Sprintf(" A: %s", proxy.Level)
	}
	var code = " "
	if proxy.CountryCode != "" {
		code = fmt.Sprintf(" %s ", proxy.CountryCode)
	}
	return fmt.Sprintf("%s:%d%s%s%s", proxy.IP, proxy.Port, code, proxy.Type, level)
}
