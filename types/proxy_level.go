package types

type ProxyLevel int

const (
	NONE        ProxyLevel = iota
	TRANSPARENT ProxyLevel = iota
	ANONYMOUS
	HIGH
)

func (p ProxyLevel) String() string {
	return [...]string{"None", "Transparent", "Anonymous", "High"}[p]
}
