package base

type Item map[string]interface{}
type Data interface {
	Valid() bool
}
