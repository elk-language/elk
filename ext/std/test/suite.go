package test

import (
	"fmt"

	"github.com/elk-language/elk/vm"
)

// Represents a test suite, a group of tests like `describe` or `context`
type Suite struct {
	Name       string
	Parent     *Suite
	SubSuites  []*Suite
	Cases      []*Case
	BeforeEach []*vm.Closure
	AfterEach  []*vm.Closure
	BeforeAll  []*vm.Closure
	AfterAll   []*vm.Closure
}

// Create a new tests suite
func NewSuite(name string, parent *Suite) *Suite {
	return &Suite{
		Name:   name,
		Parent: parent,
	}
}

func (s *Suite) NewSubSuite(name string) *Suite {
	newSuite := NewSuite(name, s)
	s.SubSuites = append(s.SubSuites, newSuite)
	return newSuite
}

func (s *Suite) NewCase(name string, fn *vm.Closure) *Case {
	newCase := NewCase(name, fn, s)
	s.Cases = append(s.Cases, newCase)
	return newCase
}

func (s *Suite) FullName() string {
	if s.Parent == nil {
		return s.Name
	}

	return fmt.Sprintf("%s %s", s.Parent.FullName(), s.Name)
}

func (s *Suite) RegisterBeforeEach(fn *vm.Closure) {
	s.BeforeEach = append(s.BeforeEach, fn)
}

func (s *Suite) RegisterBeforeAll(fn *vm.Closure) {
	s.BeforeAll = append(s.BeforeAll, fn)
}

func (s *Suite) RegisterAfterEach(fn *vm.Closure) {
	s.AfterEach = append(s.AfterEach, fn)
}

func (s *Suite) RegisterAfterAll(fn *vm.Closure) {
	s.AfterAll = append(s.AfterAll, fn)
}
