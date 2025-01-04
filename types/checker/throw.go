package checker

import (
	"fmt"

	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value/symbol"
)

type catchScope struct {
	typ        types.Type
	hasFinally bool
}

func makeCatchScope(typ types.Type, hasFinally bool) catchScope {
	return catchScope{
		typ:        typ,
		hasFinally: hasFinally,
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

func (c *Checker) checkThrowExpressionNode(node *ast.ThrowExpressionNode) *ast.ThrowExpressionNode {
	var thrownType types.Type
	if node.Value == nil {
		thrownType = c.Std(symbol.Error)
	} else {
		node.Value = c.checkExpression(node.Value)
		thrownType = c.TypeOf(node.Value)
	}

	if !node.Unchecked {
		c.checkThrowType(thrownType, node.Span())
	}

	return node
}

func (c *Checker) checkCalledMethodThrowType(method *types.Method, span *position.Span) {
	if types.IsNever(method.ThrowType) {
		return
	}

	if c.mode == tryMode {
		c.mode = validTryMode
		return
	}

	c.checkThrowType(method.ThrowType, span)
}

func (c *Checker) checkThrowType(throwType types.Type, span *position.Span) {
	for _, catchScope := range c.catchScopes {
		if c.isSubtype(throwType, catchScope.typ, nil) {
			return
		}
	}

	switch c.mode {
	case methodMode, closureInferReturnTypeMode, initMode:
		expectedThrowType := c.NewNormalisedUnion(c.throwType, throwType)
		c.addFailure(
			fmt.Sprintf(
				"thrown value of type `%s` must be caught or added to the signature of the function `%s`",
				types.InspectWithColor(throwType),
				lexer.Colorize(inspectThrow(expectedThrowType)),
			),
			span,
		)
	default:
		c.addFailure(
			fmt.Sprintf(
				"thrown value of type `%s` must be caught",
				types.InspectWithColor(throwType),
			),
			span,
		)
	}
}

func inspectThrow(throwType types.Type) string {
	return fmt.Sprintf(
		"! %s",
		types.Inspect(throwType),
	)
}
