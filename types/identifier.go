package types

import (
	"iter"

	"github.com/elk-language/elk/position"
)

// Represents an identifier of a constant, local variable, method etc
type Identifier interface {
	IdentifierName() string
	DocComment() string
}

// Represents an identifier local to a particular file
type LocalIdentifier interface {
	Identifier
	AllDeclarations() iter.Seq[*position.Location]
	AllUses() iter.Seq[*position.Location]
}

// Represents a global identifier available in the entire program
type GlobalIdentifier interface {
	Identifier
	AllDeclarationFiles() iter.Seq[string] // file paths where the identifier is declared
	AllUseFiles() iter.Seq[string]         // file paths where the identifier is used
}
