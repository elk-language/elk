package test

import (
	"fmt"

	"github.com/elk-language/elk/vm"
)

// Represents a single test case
type Case struct {
	Name   string
	Fn     *vm.Closure
	Parent *Suite
}

func NewCase(name string, fn *vm.Closure, parent *Suite) *Case {
	return &Case{
		Name:   name,
		Fn:     fn,
		Parent: parent,
	}
}

func (c *Case) FullName() string {
	if c.Parent == nil {
		return c.Name
	}

	return fmt.Sprintf("%s %s", c.Parent.FullName(), c.Name)
}
