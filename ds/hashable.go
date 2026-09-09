package ds

type Hashable interface {
	HashUint64() uint64
	EqualAny(any) bool
}
