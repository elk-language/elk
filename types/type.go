package types

type Type interface {
	IsSupertypeOf(Type) bool
}
