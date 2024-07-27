package checker

import (
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Represents a single local variable or local value
type local struct {
	typ              types.Type
	initialised      bool
	singleAssignment bool
	shadow           bool
}

func (l *local) copy() *local {
	return &local{
		typ:              l.typ,
		initialised:      l.initialised,
		singleAssignment: l.singleAssignment,
	}
}

func newLocal(typ types.Type, initialised, singleAssignment bool) *local {
	return &local{
		typ:              typ,
		initialised:      initialised,
		singleAssignment: singleAssignment,
	}
}

// Contains definitions of local variables and values
type localEnvironment struct {
	parent *localEnvironment
	locals map[value.Symbol]*local
}

// Get the local with the specified name from this local environment
func (l *localEnvironment) getLocal(name string) *local {
	local := l.locals[value.ToSymbol(name)]
	return local
}

// Resolve the local with the given name from this local environment or any parent environment
func (l *localEnvironment) resolveLocal(name string) (*local, bool) {
	nameSymbol := value.ToSymbol(name)
	currentEnv := l
	for {
		if currentEnv == nil {
			return nil, false
		}
		loc, ok := currentEnv.locals[nameSymbol]
		if ok {
			return loc, l == currentEnv
		}
		currentEnv = currentEnv.parent
	}
}

func newLocalEnvironment(parent *localEnvironment) *localEnvironment {
	return &localEnvironment{
		parent: parent,
		locals: make(map[value.Symbol]*local),
	}
}

func (c *Checker) popLocalEnv() {
	c.localEnvs = c.localEnvs[:len(c.localEnvs)-1]
}

func (c *Checker) pushNestedLocalEnv() {
	c.pushLocalEnv(newLocalEnvironment(c.currentLocalEnv()))
}

func (c *Checker) pushIsolatedLocalEnv() {
	c.pushLocalEnv(newLocalEnvironment(nil))
}

func (c *Checker) pushLocalEnv(env *localEnvironment) {
	c.localEnvs = append(c.localEnvs, env)
}

func (c *Checker) currentLocalEnv() *localEnvironment {
	return c.localEnvs[len(c.localEnvs)-1]
}
