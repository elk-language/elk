package types

type SimpleLiteral interface {
	Type
	StringValue() string
}
