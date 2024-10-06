package types

type NumericLiteral interface {
	Type
	SimpleLiteral
	IsNegative() bool
	SetNegative(bool)
	CopyNumeric() NumericLiteral
}
