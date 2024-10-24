package checker

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
)

type loop struct {
	label                         string
	returnsValueFromLastIteration bool
	returnType                    types.Type
}

func (c *Checker) registerLoop(label string, returnsValueFromLastIteration bool) *loop {
	newLoop := &loop{
		label:                         label,
		returnsValueFromLastIteration: returnsValueFromLastIteration,
	}
	c.loops = append(c.loops, newLoop)
	return newLoop
}

func (c *Checker) popLoop() {
	c.loops = c.loops[:len(c.loops)-1]
}

func (c *Checker) findLoop(label string, span *position.Span) *loop {
	if len(c.loops) < 1 {
		c.addFailure(
			"cannot jump with `break` or `continue` outside of a loop",
			span,
		)
		return nil
	}

	if label == "" {
		// if there is no label, choose the closest enclosing loop
		return c.loops[len(c.loops)-1]
	}

	for _, loop := range c.loops {
		if loop.label == label {
			return loop
		}
	}

	c.addFailure(
		fmt.Sprintf("label $%s does not exist or is not attached to an enclosing loop", label),
		span,
	)
	return nil
}
