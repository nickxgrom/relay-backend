package sender

type Type byte

const (
	System   Type = iota
	Operator Type = iota
	Client   Type = iota
)
