package types

type Variance uint8

const (
	INVARIANT Variance = iota
	COVARIANT
	CONTRAVARIANT
)

type TypeVariable struct {
	Name       string
	LowerBound Type
	UpperBound Type
	Variance   Variance
}

func NewTypeVariable(name string, lowerBound, upperBound Type, variance Variance) *TypeVariable {
	return &TypeVariable{
		Name:       name,
		LowerBound: lowerBound,
		UpperBound: upperBound,
		Variance:   variance,
	}
}

func (t *TypeVariable) ToNonLiteral(env *GlobalEnvironment) Type {
	return t
}

func (*TypeVariable) IsLiteral() bool {
	return false
}

func (t *TypeVariable) inspect() string {
	return t.Name
}
