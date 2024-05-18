package types

// All types are subtypes of any (other than `void` and `never`).
// Any is not a subtype of anything other than itself.
type Any struct{}

func (n Any) ToNonLiteral(env *GlobalEnvironment) Type {
	return n
}

func IsAny(t Type) bool {
	_, ok := t.(Any)
	return ok
}

func (Any) inspect() string {
	return "any"
}
