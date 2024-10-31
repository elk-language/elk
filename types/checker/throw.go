package checker

import (
	"github.com/elk-language/elk/types"
)

type catchScope struct {
	typ types.Type
}

func makeCatchScope(typ types.Type) catchScope {
	return catchScope{
		typ: typ,
	}
}

func (c *Checker) popCatchScope() {
	c.catchScopes = c.catchScopes[:len(c.catchScopes)-1]
}

func (c *Checker) pushCatchScope(catchScope catchScope) {
	c.catchScopes = append(c.catchScopes, catchScope)
}

func (c *Checker) enclosingCatchScope() catchScope {
	if len(c.catchScopes) < 1 {
		panic("no catch scopes!")
	}
	return c.catchScopes[len(c.catchScopes)-1]
}
