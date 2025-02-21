package types

type ProxyLevel int

const (
	TRANSPARENT ProxyLevel = iota
	ANONYMOUS
	HIGH
)

func (p ProxyLevel) String() string {
	return [...]string{"Transparent", "Anonymous", "High"}[p]
}
