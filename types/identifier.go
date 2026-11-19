package types

import (
	"iter"

	"github.com/elk-language/elk/position"
)

// Represents an identifier of a constant, local variable, method etc
type Identifier interface {
	Name() string
	DocComment() string
	AllDeclarations() iter.Seq[*position.Location]
	AllUses() iter.Seq[*position.Location]
}
