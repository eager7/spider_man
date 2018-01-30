package middleware

type StopSign interface {
	Sign() bool
	Signed() bool
	Reset()
	Deal(code string)
	DealCount() uint32
	Summary() string
}
