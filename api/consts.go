package api

type Module int64
type Controller int64

const (
	Unbound Module = iota
)

const (
	Settings Controller = iota
	Service  Controller = iota
)

func (s Module) String() string {
	switch s {
	case Unbound:
		return "unbound"
	}
	return "unknown"
}

func (s Controller) String() string {
	switch s {
	case Settings:
		return "settings"
	case Service:
		return "service"
	}
	return "unknown"
}
