package test

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

// Represents a test suite, a group of tests like `describe` or `context`
type Suite struct {
	Name       string
	Parent     *Suite
	SubSuites  []*Suite
	Cases      []*Case
	BeforeEach []value.Value
	AfterEach  []value.Value
	BeforeAll  []value.Value
	AfterAll   []value.Value
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

func (s *Suite) NewCase(name string, fn value.Value) *Case {
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

func (s *Suite) RegisterBeforeEach(fn value.Value) {
	s.BeforeEach = append(s.BeforeEach, fn)
}

func (s *Suite) RegisterBeforeAll(fn value.Value) {
	s.BeforeAll = append(s.BeforeAll, fn)
}

func (s *Suite) RegisterAfterEach(fn value.Value) {
	s.AfterEach = append(s.AfterEach, fn)
}

func (s *Suite) RegisterAfterAll(fn value.Value) {
	s.AfterAll = append(s.AfterAll, fn)
}
