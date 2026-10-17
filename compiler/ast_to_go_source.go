package compiler

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

func (c *GoCompiler) positionToGoSource(pos *position.Position) string {
	if pos == nil {
		return "nil"
	}

	return fmt.Sprintf(
		"position.New(%d, %d, %d)",
		pos.ByteOffset,
		pos.Line,
		pos.Column,
	)
}

func (c *GoCompiler) spanToGoSource(span *position.Span) string {
	if span == nil {
		return "nil"
	}

	return fmt.Sprintf(
		"position.NewSpan(%s, %s)",
		c.positionToGoSource(span.StartPos),
		c.positionToGoSource(span.EndPos),
	)
}

func (c *GoCompiler) registerPositionImport() {
	c.registerGoImport("github.com/elk-language/elk/position", "")
}

func (c *GoCompiler) registerTokenImport() {
	c.registerGoImport("github.com/elk-language/elk/token", "")
}

func (c *GoCompiler) registerBitfieldImport() {
	c.registerGoImport("github.com/elk-language/elk/bitfield", "")
}

func (c *GoCompiler) registerAstImport() {
	c.registerGoImport("github.com/elk-language/elk/parser/ast", "")
}

func (c *GoCompiler) locationToGoSource(loc *position.Location) string {
	if loc == nil || loc.Span == nil {
		return "nil"
	}

	c.registerPositionImport()
	return fmt.Sprintf(
		"position.NewLocation(%q, %s)",
		loc.FilePath,
		c.spanToGoSource(loc.Span),
	)
}

func (c *GoCompiler) tokenToGoSource(tok *token.Token) string {
	if tok == nil {
		return "nil"
	}

	c.registerTokenImport()
	loc := c.locationToGoSource(tok.Location())
	if tok.Value == "" {
		return fmt.Sprintf("token.New(%s, token.%s)", loc, tok.Type.String())
	}
	return fmt.Sprintf("token.NewWithValue(%s, token.%s, %q)", loc, tok.Type.String(), tok.Value)
}

func (c *GoCompiler) bitField8ToGoSource(b bitfield.BitField8) string {
	c.registerBitfieldImport()
	return fmt.Sprintf("bitfield.BitField8FromInt[byte](%d)", b.Byte())
}

func (c *GoCompiler) bitFlag8ToGoSource(b bitfield.BitFlag8) string {
	c.registerBitfieldImport()
	return fmt.Sprintf("bitfield.BitFlag8(%d)", b)
}

func (c *GoCompiler) astEnum(typeName string, v int) string {
	return fmt.Sprintf("ast.%s(%d)", typeName, v)
}

func isNilNode(node ast.Node) bool {
	if node == nil {
		return true
	}
	v := reflect.ValueOf(node)
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Slice, reflect.Map, reflect.Func:
		return v.IsNil()
	default:
		return false
	}
}

func astSliceToGoSource[T ast.Node](c *GoCompiler, typeName string, nodes []T, spliceInfo *goAstSpliceInfo) string {
	if len(nodes) == 0 {
		return "nil"
	}

	var buff strings.Builder
	fmt.Fprintf(&buff, "%s{", typeName)
	for i, node := range nodes {
		if i > 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(c.astNodeToGoSource(node, spliceInfo).value)
	}
	buff.WriteString("}")
	return buff.String()
}

type goAstSpliceInfo struct {
	values []*goValue
	cursor int
}

func (s *goAstSpliceInfo) currentValue() *goValue {
	value := s.values[s.cursor]
	s.cursor++
	return value
}

func newGoAstSpliceInfo(values []*goValue) *goAstSpliceInfo {
	return &goAstSpliceInfo{
		values: values,
	}
}

func (c *GoCompiler) astNodeToGoSource(node ast.Node, spliceInfo *goAstSpliceInfo) *goValue {
	if isNilNode(node) {
		return newGoValue(
			"nil",
			c.AstNodeType(),
			value.FetchGoType("ast.Node"),
		)
	}

	c.registerAstImport()

	switch n := node.(type) {
	case *ast.AliasDeclarationEntry:
		return newGoValue(
			fmt.Sprintf("ast.NewAliasDeclarationEntry(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.NewName, spliceInfo).value,
				c.astNodeToGoSource(n.OldName, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.AliasDeclarationEntry"),
		)
	case *ast.AliasDeclarationNode:
		return newGoValue(
			fmt.Sprintf("ast.NewAliasDeclarationNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]*ast.AliasDeclarationEntry", n.Entries, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.AliasDeclarationNode"),
		)
	case *ast.AnyTypeNode:
		return newGoValue(
			fmt.Sprintf("ast.NewAnyTypeNode(%s)",
				c.locationToGoSource(n.Location()),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.AnyTypeNode"),
		)
	case *ast.ArrayListLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewArrayListLiteralNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.ExpressionNode", n.Elements, spliceInfo),
				c.astNodeToGoSource(n.Capacity, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ArrayListLiteralNode"),
		)
	case *ast.ArrayTupleLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewArrayTupleLiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.ExpressionNode", n.Elements, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ArrayTupleLiteralNode"),
		)
	case *ast.AsExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewAsExpressionNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Value, spliceInfo).value,
				c.astNodeToGoSource(n.RuntimeType, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.AsExpressionNode"),
		)
	case *ast.AsPatternNode:
		return newGoValue(
			fmt.Sprintf("ast.NewAsPatternNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Pattern, spliceInfo).value,
				c.astNodeToGoSource(n.Name, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.AsPatternNode"),
		)
	case *ast.AssignmentExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewAssignmentExpressionNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.tokenToGoSource(n.Op),
				c.astNodeToGoSource(n.Left, spliceInfo).value,
				c.astNodeToGoSource(n.Right, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.AssignmentExpressionNode"),
		)
	case *ast.AttrDeclarationNode:
		return newGoValue(
			fmt.Sprintf("ast.NewAttrDeclarationNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.DocComment()),
				c.bitFlag8ToGoSource(n.Flags.ToBitFlag()),
				astSliceToGoSource(c, "[]ast.ParameterNode", n.Entries, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.AttrDeclarationNode"),
		)
	case *ast.AttributeAccessNode:
		return newGoValue(
			fmt.Sprintf("ast.NewAttributeAccessNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Receiver, spliceInfo).value,
				c.astNodeToGoSource(n.AttributeName, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.AttributeAccessNode"),
		)
	case *ast.AttributeParameterNode:
		return newGoValue(
			fmt.Sprintf("ast.NewAttributeParameterNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Name, spliceInfo).value,
				c.astNodeToGoSource(n.TypeNode, spliceInfo).value,
				c.astNodeToGoSource(n.Initialiser, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.AttributeParameterNode"),
		)
	case *ast.AwaitExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewAwaitExpressionNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Value, spliceInfo).value,
				fmt.Sprintf("%t", n.Sync),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.AwaitExpressionNode"),
		)
	case *ast.BigFloatLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewBigFloatLiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.BigFloatLiteralNode"),
		)
	case *ast.BinArrayListLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewBinArrayListLiteralNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.IntCollectionContentNode", n.Elements, spliceInfo),
				c.astNodeToGoSource(n.Capacity, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.BinArrayListLiteralNode"),
		)
	case *ast.BinArrayTupleLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewBinArrayTupleLiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.IntCollectionContentNode", n.Elements, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.BinArrayTupleLiteralNode"),
		)
	case *ast.BinHashSetLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewBinHashSetLiteralNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.IntCollectionContentNode", n.Elements, spliceInfo),
				c.astNodeToGoSource(n.Capacity, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.BinHashSetLiteralNode"),
		)
	case *ast.BinaryExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewBinaryExpressionNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.tokenToGoSource(n.Op),
				c.astNodeToGoSource(n.Left, spliceInfo).value,
				c.astNodeToGoSource(n.Right, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.BinaryExpressionNode"),
		)
	case *ast.BinaryPatternNode:
		return newGoValue(
			fmt.Sprintf("ast.NewBinaryPatternNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.tokenToGoSource(n.Op),
				c.astNodeToGoSource(n.Left, spliceInfo).value,
				c.astNodeToGoSource(n.Right, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.BinaryPatternNode"),
		)
	case *ast.BinaryTypeNode:
		return newGoValue(
			fmt.Sprintf("ast.NewBinaryTypeNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.tokenToGoSource(n.Op),
				c.astNodeToGoSource(n.Left, spliceInfo).value,
				c.astNodeToGoSource(n.Right, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.BinaryTypeNode"),
		)
	case *ast.BoolLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewBoolLiteralNode(%s)",
				c.locationToGoSource(n.Location()),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.BoolLiteralNode"),
		)
	case *ast.BoxOfExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewBoxOfExpressionNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Expression, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.BoxOfExpressionNode"),
		)
	case *ast.BoxTypeNode:
		return newGoValue(
			fmt.Sprintf("ast.NewBoxTypeNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.TypeNode, spliceInfo).value,
				fmt.Sprintf("%t", n.Immutable),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.BoxTypeNode"),
		)
	case *ast.BreakExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewBreakExpressionNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Label, spliceInfo).value,
				c.astNodeToGoSource(n.Value, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.BreakExpressionNode"),
		)
	case *ast.BreakpointNode:
		return newGoValue(
			fmt.Sprintf("ast.NewBreakpointNode(%s)",
				c.locationToGoSource(n.Location()),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.BreakpointNode"),
		)
	case *ast.CallNode:
		return newGoValue(
			fmt.Sprintf("ast.NewCallNode(%s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Receiver, spliceInfo).value,
				fmt.Sprintf("%t", n.NilSafe),
				astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments, spliceInfo),
				astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.CallNode"),
		)
	case *ast.CallableTypeNode:
		return newGoValue(
			fmt.Sprintf("ast.NewCallableTypeNode(%s, %s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.ParameterNode", n.Parameters, spliceInfo),
				c.astNodeToGoSource(n.ReturnType, spliceInfo).value,
				c.astNodeToGoSource(n.ThrowType, spliceInfo).value,
				fmt.Sprintf("%t", n.IsClosure),
				fmt.Sprintf("%t", n.IsPure),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.CallableTypeNode"),
		)
	case *ast.CatchNode:
		return newGoValue(
			fmt.Sprintf("ast.NewCatchNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Pattern, spliceInfo).value,
				c.astNodeToGoSource(n.StackTraceVar, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.StatementNode", n.Body, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.CatchNode"),
		)
	case *ast.CharLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewCharLiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.CharLiteralNode"),
		)
	case *ast.ClassDeclarationNode:
		return newGoValue(
			fmt.Sprintf("ast.NewClassDeclarationNode(%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.DocComment()),
				fmt.Sprintf("%t", n.Abstract),
				fmt.Sprintf("%t", n.Sealed),
				fmt.Sprintf("%t", n.Primitive),
				fmt.Sprintf("%t", n.NoInit),
				fmt.Sprintf("%t", n.Immutable),
				c.astNodeToGoSource(n.Constant, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.TypeParameterNode", n.TypeParameters, spliceInfo),
				c.astNodeToGoSource(n.Superclass, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.StatementNode", n.Body, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ClassDeclarationNode"),
		)
	case *ast.ClosureLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewClosureLiteralNode(%s, %s, %s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.ParameterNode", n.Parameters, spliceInfo),
				c.astNodeToGoSource(n.ReturnType, spliceInfo).value,
				c.astNodeToGoSource(n.ThrowType, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.StatementNode", n.Body, spliceInfo),
				fmt.Sprintf("%t", n.Lambda),
				fmt.Sprintf("%t", n.Pure),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ClosureLiteralNode"),
		)
	case *ast.ConstantAsNode:
		return newGoValue(
			fmt.Sprintf("ast.NewConstantAsNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Constant, spliceInfo).value,
				fmt.Sprintf("%q", n.AsName),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ConstantAsNode"),
		)
	case *ast.ConstantDeclarationNode:
		return newGoValue(
			fmt.Sprintf("ast.NewConstantDeclarationNode(%s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.DocComment()),
				c.astNodeToGoSource(n.Constant, spliceInfo).value,
				c.astNodeToGoSource(n.TypeNode, spliceInfo).value,
				c.astNodeToGoSource(n.Initialiser, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ConstantDeclarationNode"),
		)
	case *ast.ConstantLookupNode:
		return newGoValue(
			fmt.Sprintf("ast.NewConstantLookupNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Left, spliceInfo).value,
				c.astNodeToGoSource(n.Right, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ConstantLookupNode"),
		)
	case *ast.ConstructorCallNode:
		return newGoValue(
			fmt.Sprintf("ast.NewConstructorCallNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.ClassNode, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments, spliceInfo),
				astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ConstructorCallNode"),
		)
	case *ast.ContinueExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewContinueExpressionNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Label, spliceInfo).value,
				c.astNodeToGoSource(n.Value, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ContinueExpressionNode"),
		)
	case *ast.DeferExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewDeferExpressionNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Expression, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.DeferExpressionNode"),
		)
	case *ast.DoExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewDoExpressionNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.StatementNode", n.Body, spliceInfo),
				astSliceToGoSource(c, "[]*ast.CatchNode", n.Catches, spliceInfo),
				astSliceToGoSource(c, "[]ast.StatementNode", n.Finally, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.DoExpressionNode"),
		)
	case *ast.DoubleQuotedStringLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewDoubleQuotedStringLiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.DoubleQuotedStringLiteralNode"),
		)
	case *ast.DoubleSplatExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewDoubleSplatExpressionNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Value, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.DoubleSplatExpressionNode"),
		)
	case *ast.EmptyStatementNode:
		return newGoValue(
			fmt.Sprintf("ast.NewEmptyStatementNode(%s)",
				c.locationToGoSource(n.Location()),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.EmptyStatementNode"),
		)
	case *ast.ExactTypeNode:
		return newGoValue(
			fmt.Sprintf("ast.NewExactTypeNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.TypeNode, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ExactTypeNode"),
		)
	case *ast.ExpressionStatementNode:
		return newGoValue(
			fmt.Sprintf("ast.NewExpressionStatementNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Expression, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ExpressionStatementNode"),
		)
	case *ast.ExtendWhereBlockExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewExtendWhereBlockExpressionNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.StatementNode", n.Body, spliceInfo),
				astSliceToGoSource(c, "[]ast.TypeParameterNode", n.Where, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ExtendWhereBlockExpressionNode"),
		)
	case *ast.FalseLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewFalseLiteralNode(%s)",
				c.locationToGoSource(n.Location()),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.FalseLiteralNode"),
		)
	case *ast.Float32LiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewFloat32LiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.Float32LiteralNode"),
		)
	case *ast.Float64LiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewFloat64LiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.Float64LiteralNode"),
		)
	case *ast.FloatLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewFloatLiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.FloatLiteralNode"),
		)
	case *ast.ForInExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewForInExpressionNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Pattern, spliceInfo).value,
				c.astNodeToGoSource(n.InExpression, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.StatementNode", n.ThenBody, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ForInExpressionNode"),
		)
	case *ast.FormalParameterNode:
		return newGoValue(
			fmt.Sprintf("ast.NewFormalParameterNode(%s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Name, spliceInfo).value,
				c.astNodeToGoSource(n.TypeNode, spliceInfo).value,
				c.astNodeToGoSource(n.Initialiser, spliceInfo).value,
				c.astEnum("ParameterKind", int(n.Kind)),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.FormalParameterNode"),
		)
	case *ast.GenericConstantNode:
		return newGoValue(
			fmt.Sprintf("ast.NewGenericConstantNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Constant, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.TypeNode", n.TypeArguments, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.GenericConstantNode"),
		)
	case *ast.GenericConstructorCallNode:
		return newGoValue(
			fmt.Sprintf("ast.NewGenericConstructorCallNode(%s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.ClassNode, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.TypeNode", n.TypeArguments, spliceInfo),
				astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments, spliceInfo),
				astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.GenericConstructorCallNode"),
		)
	case *ast.GenericMethodCallNode:
		return newGoValue(
			fmt.Sprintf("ast.NewGenericMethodCallNode(%s, %s, %s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Receiver, spliceInfo).value,
				c.tokenToGoSource(n.Op),
				c.astNodeToGoSource(n.MethodName, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.TypeNode", n.TypeArguments, spliceInfo),
				astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments, spliceInfo),
				astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.GenericMethodCallNode"),
		)
	case *ast.GenericReceiverlessMethodCallNode:
		return newGoValue(
			fmt.Sprintf("ast.NewGenericReceiverlessMethodCallNode(%s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.MethodName, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.TypeNode", n.TypeArguments, spliceInfo),
				astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments, spliceInfo),
				astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.GenericReceiverlessMethodCallNode"),
		)
	case *ast.GenericTypeDefinitionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewGenericTypeDefinitionNode(%s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.DocComment()),
				c.astNodeToGoSource(n.Constant, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.TypeParameterNode", n.TypeParameters, spliceInfo),
				c.astNodeToGoSource(n.TypeNode, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.GenericTypeDefinitionNode"),
		)
	case *ast.GetterDeclarationNode:
		return newGoValue(
			fmt.Sprintf("ast.NewGetterDeclarationNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.DocComment()),
				c.bitFlag8ToGoSource(n.Flags.ToBitFlag()),
				astSliceToGoSource(c, "[]ast.ParameterNode", n.Entries, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.GetterDeclarationNode"),
		)
	case *ast.GoExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewGoExpressionNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.StatementNode", n.Body, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.GoExpressionNode"),
		)
	case *ast.HashMapLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewHashMapLiteralNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.ExpressionNode", n.Elements, spliceInfo),
				c.astNodeToGoSource(n.Capacity, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.HashMapLiteralNode"),
		)
	case *ast.HashRecordLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewHashRecordLiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.ExpressionNode", n.Elements, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.HashRecordLiteralNode"),
		)
	case *ast.HashSetLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewHashSetLiteralNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.ExpressionNode", n.Elements, spliceInfo),
				c.astNodeToGoSource(n.Capacity, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.HashSetLiteralNode"),
		)
	case *ast.HexArrayListLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewHexArrayListLiteralNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.IntCollectionContentNode", n.Elements, spliceInfo),
				c.astNodeToGoSource(n.Capacity, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.HexArrayListLiteralNode"),
		)
	case *ast.HexArrayTupleLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewHexArrayTupleLiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.IntCollectionContentNode", n.Elements, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.HexArrayTupleLiteralNode"),
		)
	case *ast.HexHashSetLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewHexHashSetLiteralNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.IntCollectionContentNode", n.Elements, spliceInfo),
				c.astNodeToGoSource(n.Capacity, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.HexHashSetLiteralNode"),
		)
	case *ast.IfExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewIfExpressionNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Condition, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.StatementNode", n.ThenBody, spliceInfo),
				astSliceToGoSource(c, "[]ast.StatementNode", n.ElseBody, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.IfExpressionNode"),
		)
	case *ast.ImplementExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewImplementExpressionNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.ComplexConstantNode", n.Constants, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ImplementExpressionNode"),
		)
	case *ast.ImportStatementNode:
		return newGoValue(
			fmt.Sprintf("ast.NewImportStatementNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Path, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ImportStatementNode"),
		)
	case *ast.IncludeExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewIncludeExpressionNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.ComplexConstantNode", n.Constants, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.IncludeExpressionNode"),
		)
	case *ast.InferredObjectPatternNode:
		return newGoValue(
			fmt.Sprintf("ast.NewInferredObjectPatternNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.PatternNode", n.Attributes, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.InferredObjectPatternNode"),
		)
	case *ast.InitDefinitionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewInitDefinitionNode(%s, %s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.DocComment()),
				c.bitFlag8ToGoSource(n.Flags.ToBitFlag()),
				astSliceToGoSource(c, "[]ast.ParameterNode", n.Parameters, spliceInfo),
				c.astNodeToGoSource(n.ThrowType, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.StatementNode", n.Body, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.InitDefinitionNode"),
		)
	case *ast.InstanceMethodLookupNode:
		return newGoValue(
			fmt.Sprintf("ast.NewInstanceMethodLookupNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Receiver, spliceInfo).value,
				c.astNodeToGoSource(n.Name, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.InstanceMethodLookupNode"),
		)
	case *ast.InstanceOfTypeNode:
		return newGoValue(
			fmt.Sprintf("ast.NewInstanceOfTypeNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.TypeNode, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.InstanceOfTypeNode"),
		)
	case *ast.InstanceValueDeclarationNode:
		return newGoValue(
			fmt.Sprintf("ast.NewInstanceValueDeclarationNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.DocComment()),
				c.astNodeToGoSource(n.Name, spliceInfo).value,
				c.astNodeToGoSource(n.TypeNode, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.InstanceValueDeclarationNode"),
		)
	case *ast.InstanceVariableDeclarationNode:
		return newGoValue(
			fmt.Sprintf("ast.NewInstanceVariableDeclarationNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.DocComment()),
				c.astNodeToGoSource(n.Name, spliceInfo).value,
				c.astNodeToGoSource(n.TypeNode, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.InstanceVariableDeclarationNode"),
		)
	case *ast.Int16LiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewInt16LiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.Int16LiteralNode"),
		)
	case *ast.Int32LiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewInt32LiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.Int32LiteralNode"),
		)
	case *ast.Int64LiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewInt64LiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.Int64LiteralNode"),
		)
	case *ast.Int8LiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewInt8LiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.Int8LiteralNode"),
		)
	case *ast.IntLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewIntLiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.IntLiteralNode"),
		)
	case *ast.InterfaceDeclarationNode:
		return newGoValue(
			fmt.Sprintf("ast.NewInterfaceDeclarationNode(%s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.DocComment()),
				c.astNodeToGoSource(n.Constant, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.TypeParameterNode", n.TypeParameters, spliceInfo),
				astSliceToGoSource(c, "[]ast.StatementNode", n.Body, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.InterfaceDeclarationNode"),
		)
	case *ast.InterpolatedRegexLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewInterpolatedRegexLiteralNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.RegexLiteralContentNode", n.Content, spliceInfo),
				c.bitField8ToGoSource(n.Flags),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.InterpolatedRegexLiteralNode"),
		)
	case *ast.InterpolatedStringLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewInterpolatedStringLiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.StringLiteralContentNode", n.Content, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.InterpolatedStringLiteralNode"),
		)
	case *ast.InterpolatedSymbolLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewInterpolatedSymbolLiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Content, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.InterpolatedSymbolLiteralNode"),
		)
	case *ast.IntersectionTypeNode:
		return newGoValue(
			fmt.Sprintf("ast.NewIntersectionTypeNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.TypeNode", n.Elements, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.IntersectionTypeNode"),
		)
	case *ast.InvalidNode:
		return newGoValue(
			fmt.Sprintf("ast.NewInvalidNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.tokenToGoSource(n.Token),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.InvalidNode"),
		)
	case *ast.KeyValueExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewKeyValueExpressionNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Key, spliceInfo).value,
				c.astNodeToGoSource(n.Value, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.KeyValueExpressionNode"),
		)
	case *ast.KeyValuePatternNode:
		return newGoValue(
			fmt.Sprintf("ast.NewKeyValuePatternNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Key, spliceInfo).value,
				c.astNodeToGoSource(n.Value, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.KeyValuePatternNode"),
		)
	case *ast.LabeledExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewLabeledExpressionNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Label),
				c.astNodeToGoSource(n.Expression, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.LabeledExpressionNode"),
		)
	case *ast.ListPatternNode:
		return newGoValue(
			fmt.Sprintf("ast.NewListPatternNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.PatternNode", n.Elements, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ListPatternNode"),
		)
	case *ast.LogicalExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewLogicalExpressionNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.tokenToGoSource(n.Op),
				c.astNodeToGoSource(n.Left, spliceInfo).value,
				c.astNodeToGoSource(n.Right, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.LogicalExpressionNode"),
		)
	case *ast.LoopExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewLoopExpressionNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.StatementNode", n.ThenBody, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.LoopExpressionNode"),
		)
	case *ast.MacroBoundaryNode:
		return newGoValue(
			fmt.Sprintf("ast.NewMacroBoundaryNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.StatementNode", n.Body, spliceInfo),
				fmt.Sprintf("%q", n.Name),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.MacroBoundaryNode"),
		)
	case *ast.MacroCallNode:
		return newGoValue(
			fmt.Sprintf("ast.NewMacroCallNode(%s, %s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astEnum("MacroKind", int(n.Kind)),
				c.astNodeToGoSource(n.Receiver, spliceInfo).value,
				c.astNodeToGoSource(n.MacroName, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments, spliceInfo),
				astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.MacroCallNode"),
		)
	case *ast.MacroDefinitionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewMacroDefinitionNode(%s, %s, %s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.DocComment()),
				fmt.Sprintf("%t", n.IsSealed()),
				c.astNodeToGoSource(n.Name, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.ParameterNode", n.Parameters, spliceInfo),
				c.astNodeToGoSource(n.ReturnType, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.StatementNode", n.Body, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.MacroDefinitionNode"),
		)
	case *ast.MacroNameNode:
		return newGoValue(
			fmt.Sprintf("ast.NewMacroNameNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.MacroNameNode"),
		)
	case *ast.MapPatternNode:
		return newGoValue(
			fmt.Sprintf("ast.NewMapPatternNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.PatternNode", n.Elements, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.MapPatternNode"),
		)
	case *ast.MatchExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewMatchExpressionNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Expression, spliceInfo).value,
				c.astNodeToGoSource(n.Pattern, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.MatchExpressionNode"),
		)
	case *ast.MethodCallNode:
		return newGoValue(
			fmt.Sprintf("ast.NewMethodCallNode(%s, %s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Receiver, spliceInfo).value,
				c.tokenToGoSource(n.Op),
				c.astNodeToGoSource(n.MethodName, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments, spliceInfo),
				astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.MethodCallNode"),
		)
	case *ast.MethodDefinitionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewMethodDefinitionNode(%s, %s, %s, %s, %s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.DocComment()),
				c.bitFlag8ToGoSource(n.Flags.ToBitFlag()),
				c.astNodeToGoSource(n.Name, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.TypeParameterNode", n.TypeParameters, spliceInfo),
				astSliceToGoSource(c, "[]ast.ParameterNode", n.Parameters, spliceInfo),
				c.astNodeToGoSource(n.ReturnType, spliceInfo).value,
				c.astNodeToGoSource(n.ThrowType, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.StatementNode", n.Body, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.MethodDefinitionNode"),
		)
	case *ast.MethodLookupAsNode:
		return newGoValue(
			fmt.Sprintf("ast.NewMethodLookupAsNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.MethodLookup, spliceInfo).value,
				c.astNodeToGoSource(n.AsName, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.MethodLookupAsNode"),
		)
	case *ast.MethodLookupNode:
		return newGoValue(
			fmt.Sprintf("ast.NewMethodLookupNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Receiver, spliceInfo).value,
				c.astNodeToGoSource(n.Name, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.MethodLookupNode"),
		)
	case *ast.MethodParameterNode:
		return newGoValue(
			fmt.Sprintf("ast.NewMethodParameterNode(%s, %s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Name, spliceInfo).value,
				fmt.Sprintf("%t", n.SetInstanceVariable),
				c.astNodeToGoSource(n.TypeNode, spliceInfo).value,
				c.astNodeToGoSource(n.Initialiser, spliceInfo).value,
				c.astEnum("ParameterKind", int(n.Kind)),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.MethodParameterNode"),
		)
	case *ast.MethodSignatureDefinitionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewMethodSignatureDefinitionNode(%s, %s, %s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.DocComment()),
				c.astNodeToGoSource(n.Name, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.TypeParameterNode", n.TypeParameters, spliceInfo),
				astSliceToGoSource(c, "[]ast.ParameterNode", n.Parameters, spliceInfo),
				c.astNodeToGoSource(n.ReturnType, spliceInfo).value,
				c.astNodeToGoSource(n.ThrowType, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.MethodSignatureDefinitionNode"),
		)
	case *ast.MixinDeclarationNode:
		return newGoValue(
			fmt.Sprintf("ast.NewMixinDeclarationNode(%s, %s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.DocComment()),
				fmt.Sprintf("%t", n.Abstract),
				c.astNodeToGoSource(n.Constant, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.TypeParameterNode", n.TypeParameters, spliceInfo),
				astSliceToGoSource(c, "[]ast.StatementNode", n.Body, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.MixinDeclarationNode"),
		)
	case *ast.ModifierForInNode:
		return newGoValue(
			fmt.Sprintf("ast.NewModifierForInNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.ThenExpression, spliceInfo).value,
				c.astNodeToGoSource(n.Pattern, spliceInfo).value,
				c.astNodeToGoSource(n.InExpression, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ModifierForInNode"),
		)
	case *ast.ModifierIfElseNode:
		return newGoValue(
			fmt.Sprintf("ast.NewModifierIfElseNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.ThenExpression, spliceInfo).value,
				c.astNodeToGoSource(n.Condition, spliceInfo).value,
				c.astNodeToGoSource(n.ElseExpression, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ModifierIfElseNode"),
		)
	case *ast.ModifierNode:
		return newGoValue(
			fmt.Sprintf("ast.NewModifierNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.tokenToGoSource(n.Modifier),
				c.astNodeToGoSource(n.Left, spliceInfo).value,
				c.astNodeToGoSource(n.Right, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ModifierNode"),
		)
	case *ast.ModuleDeclarationNode:
		return newGoValue(
			fmt.Sprintf("ast.NewModuleDeclarationNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.DocComment()),
				c.astNodeToGoSource(n.Constant, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.StatementNode", n.Body, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ModuleDeclarationNode"),
		)
	case *ast.MustExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewMustExpressionNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Value, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.MustExpressionNode"),
		)
	case *ast.MustPatternNode:
		return newGoValue(
			fmt.Sprintf("ast.NewMustPatternNode(%s)",
				c.locationToGoSource(n.Location()),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.MustPatternNode"),
		)
	case *ast.NamedCallArgumentNode:
		return newGoValue(
			fmt.Sprintf("ast.NewNamedCallArgumentNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Name, spliceInfo).value,
				c.astNodeToGoSource(n.Value, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.NamedCallArgumentNode"),
		)
	case *ast.NeverTypeNode:
		return newGoValue(
			fmt.Sprintf("ast.NewNeverTypeNode(%s)",
				c.locationToGoSource(n.Location()),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.NeverTypeNode"),
		)
	case *ast.NewExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewNewExpressionNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments, spliceInfo),
				astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.NewExpressionNode"),
		)
	case *ast.NilLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewNilLiteralNode(%s)",
				c.locationToGoSource(n.Location()),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.NilLiteralNode"),
		)
	case *ast.NilSafeSubscriptExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewNilSafeSubscriptExpressionNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Receiver, spliceInfo).value,
				c.astNodeToGoSource(n.Key, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.NilSafeSubscriptExpressionNode"),
		)
	case *ast.NilablePatternNode:
		return newGoValue(
			fmt.Sprintf("ast.NewNilablePatternNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Pattern, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.NilablePatternNode"),
		)
	case *ast.NilableTypeNode:
		return newGoValue(
			fmt.Sprintf("ast.NewNilableTypeNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.TypeNode, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.NilableTypeNode"),
		)
	case *ast.NotTypeNode:
		return newGoValue(
			fmt.Sprintf("ast.NewNotTypeNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.TypeNode, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.NotTypeNode"),
		)
	case *ast.NumericForExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewNumericForExpressionNode(%s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Initialiser, spliceInfo).value,
				c.astNodeToGoSource(n.Condition, spliceInfo).value,
				c.astNodeToGoSource(n.Increment, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.StatementNode", n.ThenBody, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.NumericForExpressionNode"),
		)
	case *ast.ObjectPatternNode:
		return newGoValue(
			fmt.Sprintf("ast.NewObjectPatternNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.ObjectType, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.PatternNode", n.Attributes, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ObjectPatternNode"),
		)
	case *ast.ParameterStatementNode:
		return newGoValue(
			fmt.Sprintf("ast.NewParameterStatementNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Parameter, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ParameterStatementNode"),
		)
	case *ast.PatternExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewPatternExpressionNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.PatternNode, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.PatternExpressionNode"),
		)
	case *ast.PatternStatementNode:
		return newGoValue(
			fmt.Sprintf("ast.NewPatternStatementNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Pattern, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.PatternStatementNode"),
		)
	case *ast.PostfixExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewPostfixExpressionNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.tokenToGoSource(n.Op),
				c.astNodeToGoSource(n.Expression, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.PostfixExpressionNode"),
		)
	case *ast.PrivateConstantNode:
		return newGoValue(
			fmt.Sprintf("ast.NewPrivateConstantNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.PrivateConstantNode"),
		)
	case *ast.PrivateIdentifierNode:
		return newGoValue(
			fmt.Sprintf("ast.NewPrivateIdentifierNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.PrivateIdentifierNode"),
		)
	case *ast.ProgramNode:
		return newGoValue(
			fmt.Sprintf("ast.NewProgramNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.StatementNode", n.Body, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ProgramNode"),
		)
	case *ast.PublicConstantAsNode:
		return newGoValue(
			fmt.Sprintf("ast.NewPublicConstantAsNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Target, spliceInfo).value,
				fmt.Sprintf("%q", n.AsName),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.PublicConstantAsNode"),
		)
	case *ast.PublicConstantNode:
		return newGoValue(
			fmt.Sprintf("ast.NewPublicConstantNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.PublicConstantNode"),
		)
	case *ast.PublicIdentifierNode:
		return newGoValue(
			fmt.Sprintf("ast.NewPublicIdentifierNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.PublicIdentifierNode"),
		)
	case *ast.PublicInstanceVariableNode:
		return newGoValue(
			fmt.Sprintf("ast.NewPublicInstanceVariableNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.PublicInstanceVariableNode"),
		)
	case *ast.QuoteExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewQuoteExpressionNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astEnum("QuoteKind", int(n.Kind)),
				astSliceToGoSource(c, "[]ast.StatementNode", n.Body, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.QuoteExpressionNode"),
		)
	case *ast.RangeLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewRangeLiteralNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.tokenToGoSource(n.Op),
				c.astNodeToGoSource(n.Start, spliceInfo).value,
				c.astNodeToGoSource(n.End, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.RangeLiteralNode"),
		)
	case *ast.RawCharLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewRawCharLiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.RawCharLiteralNode"),
		)
	case *ast.RawStringLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewRawStringLiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.RawStringLiteralNode"),
		)
	case *ast.ReceiverlessMacroCallNode:
		return newGoValue(
			fmt.Sprintf("ast.NewReceiverlessMacroCallNode(%s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astEnum("MacroKind", int(n.Kind)),
				c.astNodeToGoSource(n.MacroName, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments, spliceInfo),
				astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ReceiverlessMacroCallNode"),
		)
	case *ast.ReceiverlessMethodCallNode:
		return newGoValue(
			fmt.Sprintf("ast.NewReceiverlessMethodCallNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.MethodName, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments, spliceInfo),
				astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ReceiverlessMethodCallNode"),
		)
	case *ast.RecordPatternNode:
		return newGoValue(
			fmt.Sprintf("ast.NewRecordPatternNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.PatternNode", n.Elements, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.RecordPatternNode"),
		)
	case *ast.RegexInterpolationNode:
		return newGoValue(
			fmt.Sprintf("ast.NewRegexInterpolationNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Expression, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.RegexInterpolationNode"),
		)
	case *ast.RegexLiteralContentSectionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewRegexLiteralContentSectionNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.RegexLiteralContentSectionNode"),
		)
	case *ast.RestPatternNode:
		return newGoValue(
			fmt.Sprintf("ast.NewRestPatternNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Identifier, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.RestPatternNode"),
		)
	case *ast.ReturnExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewReturnExpressionNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Value, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ReturnExpressionNode"),
		)
	case *ast.ScopedMacroCallNode:
		return newGoValue(
			fmt.Sprintf("ast.NewScopedMacroCallNode(%s, %s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astEnum("MacroKind", int(n.Kind)),
				c.astNodeToGoSource(n.Receiver, spliceInfo).value,
				c.astNodeToGoSource(n.MacroName, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments, spliceInfo),
				astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ScopedMacroCallNode"),
		)
	case *ast.SelectCaseNode:
		return newGoValue(
			fmt.Sprintf("ast.NewSelectCaseNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Expression, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.StatementNode", n.Body, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.SelectCaseNode"),
		)
	case *ast.SelectExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewSelectExpressionNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]*ast.SelectCaseNode", n.Cases, spliceInfo),
				astSliceToGoSource(c, "[]ast.StatementNode", n.ElseBody, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.SelectExpressionNode"),
		)
	case *ast.SelfLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewSelfLiteralNode(%s)",
				c.locationToGoSource(n.Location()),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.SelfLiteralNode"),
		)
	case *ast.SetPatternNode:
		return newGoValue(
			fmt.Sprintf("ast.NewSetPatternNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.PatternNode", n.Elements, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.SetPatternNode"),
		)
	case *ast.SetterDeclarationNode:
		return newGoValue(
			fmt.Sprintf("ast.NewSetterDeclarationNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.DocComment()),
				astSliceToGoSource(c, "[]ast.ParameterNode", n.Entries, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.SetterDeclarationNode"),
		)
	case *ast.SignatureParameterNode:
		return newGoValue(
			fmt.Sprintf("ast.NewSignatureParameterNode(%s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Name, spliceInfo).value,
				c.astNodeToGoSource(n.TypeNode, spliceInfo).value,
				fmt.Sprintf("%t", n.Optional),
				c.astEnum("ParameterKind", int(n.Kind)),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.SignatureParameterNode"),
		)
	case *ast.SimpleSymbolLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewSimpleSymbolLiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Content),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.SimpleSymbolLiteralNode"),
		)
	case *ast.SingletonBlockExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewSingletonBlockExpressionNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.StatementNode", n.Body, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.SingletonBlockExpressionNode"),
		)
	case *ast.SingletonTypeNode:
		return newGoValue(
			fmt.Sprintf("ast.NewSingletonTypeNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.TypeNode, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.SingletonTypeNode"),
		)
	case *ast.SplatExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewSplatExpressionNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Value, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.SplatExpressionNode"),
		)
	case *ast.StringInspectInterpolationNode:
		return newGoValue(
			fmt.Sprintf("ast.NewStringInspectInterpolationNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Expression, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.StringInspectInterpolationNode"),
		)
	case *ast.StringInterpolationNode:
		return newGoValue(
			fmt.Sprintf("ast.NewStringInterpolationNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Expression, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.StringInterpolationNode"),
		)
	case *ast.StringLiteralContentSectionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewStringLiteralContentSectionNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.StringLiteralContentSectionNode"),
		)
	case *ast.StructDeclarationNode:
		return newGoValue(
			fmt.Sprintf("ast.NewStructDeclarationNode(%s, %s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.DocComment()),
				fmt.Sprintf("%t", n.Immutable),
				c.astNodeToGoSource(n.Constant, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.TypeParameterNode", n.TypeParameters, spliceInfo),
				astSliceToGoSource(c, "[]ast.StructBodyStatementNode", n.Body, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.StructDeclarationNode"),
		)
	case *ast.SubscriptExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewSubscriptExpressionNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Receiver, spliceInfo).value,
				c.astNodeToGoSource(n.Key, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.SubscriptExpressionNode"),
		)
	case *ast.SwitchCaseNode:
		return newGoValue(
			fmt.Sprintf("ast.NewSwitchCaseNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Pattern, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.StatementNode", n.Body, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.SwitchCaseNode"),
		)
	case *ast.SwitchExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewSwitchExpressionNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Value, spliceInfo).value,
				astSliceToGoSource(c, "[]*ast.SwitchCaseNode", n.Cases, spliceInfo),
				astSliceToGoSource(c, "[]ast.StatementNode", n.ElseBody, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.SwitchExpressionNode"),
		)
	case *ast.SymbolArrayListLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewSymbolArrayListLiteralNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.SymbolCollectionContentNode", n.Elements, spliceInfo),
				c.astNodeToGoSource(n.Capacity, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.SymbolArrayListLiteralNode"),
		)
	case *ast.SymbolArrayTupleLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewSymbolArrayTupleLiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.SymbolCollectionContentNode", n.Elements, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.SymbolArrayTupleLiteralNode"),
		)
	case *ast.SymbolHashSetLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewSymbolHashSetLiteralNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.SymbolCollectionContentNode", n.Elements, spliceInfo),
				c.astNodeToGoSource(n.Capacity, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.SymbolHashSetLiteralNode"),
		)
	case *ast.SymbolKeyValueExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewSymbolKeyValueExpressionNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Key, spliceInfo).value,
				c.astNodeToGoSource(n.Value, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.SymbolKeyValueExpressionNode"),
		)
	case *ast.SymbolKeyValuePatternNode:
		return newGoValue(
			fmt.Sprintf("ast.NewSymbolKeyValuePatternNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Key, spliceInfo).value,
				c.astNodeToGoSource(n.Value, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.SymbolKeyValuePatternNode"),
		)
	case *ast.ThrowExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewThrowExpressionNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%t", n.Unchecked),
				c.astNodeToGoSource(n.Value, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ThrowExpressionNode"),
		)
	case *ast.TrueLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewTrueLiteralNode(%s)",
				c.locationToGoSource(n.Location()),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.TrueLiteralNode"),
		)
	case *ast.TryExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewTryExpressionNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Value, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.TryExpressionNode"),
		)
	case *ast.TuplePatternNode:
		return newGoValue(
			fmt.Sprintf("ast.NewTuplePatternNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.PatternNode", n.Elements, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.TuplePatternNode"),
		)
	case *ast.TypeDefinitionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewTypeDefinitionNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.DocComment()),
				c.astNodeToGoSource(n.Constant, spliceInfo).value,
				c.astNodeToGoSource(n.TypeNode, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.TypeDefinitionNode"),
		)
	case *ast.TypeExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewTypeExpressionNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.TypeNode, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.TypeExpressionNode"),
		)
	case *ast.TypeStatementNode:
		return newGoValue(
			fmt.Sprintf("ast.NewTypeStatementNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.TypeNode, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.TypeStatementNode"),
		)
	case *ast.TypeofExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewTypeofExpressionNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Value, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.TypeofExpressionNode"),
		)
	case *ast.UInt16LiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewUInt16LiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.UInt16LiteralNode"),
		)
	case *ast.UInt32LiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewUInt32LiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.UInt32LiteralNode"),
		)
	case *ast.UInt64LiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewUInt64LiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.UInt64LiteralNode"),
		)
	case *ast.UInt8LiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewUInt8LiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.UInt8LiteralNode"),
		)
	case *ast.UIntLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewUIntLiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Value),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.UIntLiteralNode"),
		)
	case *ast.UnaryExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewUnaryExpressionNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.tokenToGoSource(n.Op),
				c.astNodeToGoSource(n.Right, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.UnaryExpressionNode"),
		)
	case *ast.UnaryTypeNode:
		return newGoValue(
			fmt.Sprintf("ast.NewUnaryTypeNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.tokenToGoSource(n.Op),
				c.astNodeToGoSource(n.TypeNode, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.UnaryTypeNode"),
		)
	case *ast.UndefinedLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewUndefinedLiteralNode(%s)",
				c.locationToGoSource(n.Location()),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.UndefinedLiteralNode"),
		)
	case *ast.UnhygienicNode:
		return newGoValue(
			fmt.Sprintf("ast.NewUnhygienicNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Node, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.UnhygienicNode"),
		)
	case *ast.UninterpolatedRegexLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewUninterpolatedRegexLiteralNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.Content),
				c.bitField8ToGoSource(n.Flags),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.UninterpolatedRegexLiteralNode"),
		)
	case *ast.UnionTypeNode:
		return newGoValue(
			fmt.Sprintf("ast.NewUnionTypeNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.TypeNode", n.Elements, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.UnionTypeNode"),
		)
	case *ast.UnlessExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewUnlessExpressionNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Condition, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.StatementNode", n.ThenBody, spliceInfo),
				astSliceToGoSource(c, "[]ast.StatementNode", n.ElseBody, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.UnlessExpressionNode"),
		)
	case *ast.UnquoteForInExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewUnquoteForInExpressionNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Pattern, spliceInfo).value,
				c.astNodeToGoSource(n.InExpression, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.StatementNode", n.ThenBody, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.UnquoteForInExpressionNode"),
		)
	case *ast.UnquoteIfExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewUnquoteIfExpressionNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Condition, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.StatementNode", n.ThenBody, spliceInfo),
				astSliceToGoSource(c, "[]ast.StatementNode", n.ElseBody, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.UnquoteIfExpressionNode"),
		)
	case *ast.UnquoteNode:
		if spliceInfo == nil {
			return newGoValue(
				fmt.Sprintf("ast.NewUnquoteNode(%s, %s, %s)",
					c.locationToGoSource(n.Location()),
					c.astEnum("UnquoteKind", int(n.Kind)),
					c.astNodeToGoSource(n.Expression, spliceInfo).value,
				),
				node.MacroType(c.checker.Env()),
				value.FetchGoType("*ast.UnquoteNode"),
			)
		}

		return c.convertValueToNarrowerType(spliceInfo.currentValue())
	case *ast.UntilExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewUntilExpressionNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Condition, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.StatementNode", n.ThenBody, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.UntilExpressionNode"),
		)
	case *ast.UsingAllEntryNode:
		return newGoValue(
			fmt.Sprintf("ast.NewUsingAllEntryNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Namespace, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.UsingAllEntryNode"),
		)
	case *ast.UsingEntryWithSubentriesNode:
		return newGoValue(
			fmt.Sprintf("ast.NewUsingEntryWithSubentriesNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Namespace, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.UsingSubentryNode", n.Subentries, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.UsingEntryWithSubentriesNode"),
		)
	case *ast.UsingExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewUsingExpressionNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.UsingEntryNode", n.Entries, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.UsingExpressionNode"),
		)
	case *ast.UsingSubentryAsNode:
		return newGoValue(
			fmt.Sprintf("ast.NewUsingSubentryAsNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Target, spliceInfo).value,
				c.astNodeToGoSource(n.AsName, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.UsingSubentryAsNode"),
		)
	case *ast.ValueDeclarationNode:
		return newGoValue(
			fmt.Sprintf("ast.NewValueDeclarationNode(%s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Name, spliceInfo).value,
				c.astNodeToGoSource(n.TypeNode, spliceInfo).value,
				c.astNodeToGoSource(n.Initialiser, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ValueDeclarationNode"),
		)
	case *ast.ValuePatternDeclarationNode:
		return newGoValue(
			fmt.Sprintf("ast.NewValuePatternDeclarationNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Pattern, spliceInfo).value,
				c.astNodeToGoSource(n.Initialiser, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.ValuePatternDeclarationNode"),
		)
	case *ast.VariableDeclarationNode:
		return newGoValue(
			fmt.Sprintf("ast.NewVariableDeclarationNode(%s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%q", n.DocComment()),
				c.astNodeToGoSource(n.Name, spliceInfo).value,
				c.astNodeToGoSource(n.TypeNode, spliceInfo).value,
				c.astNodeToGoSource(n.Initialiser, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.VariableDeclarationNode"),
		)
	case *ast.VariablePatternDeclarationNode:
		return newGoValue(
			fmt.Sprintf("ast.NewVariablePatternDeclarationNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Pattern, spliceInfo).value,
				c.astNodeToGoSource(n.Initialiser, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.VariablePatternDeclarationNode"),
		)
	case *ast.VariantTypeParameterNode:
		return newGoValue(
			fmt.Sprintf("ast.NewVariantTypeParameterNode(%s, %s, %s, %s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astEnum("Variance", int(n.Variance)),
				fmt.Sprintf("%q", n.Name),
				c.astNodeToGoSource(n.LowerBound, spliceInfo).value,
				c.astNodeToGoSource(n.UpperBound, spliceInfo).value,
				c.astNodeToGoSource(n.Default, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.VariantTypeParameterNode"),
		)
	case *ast.VoidTypeNode:
		return newGoValue(
			fmt.Sprintf("ast.NewVoidTypeNode(%s)",
				c.locationToGoSource(n.Location()),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.VoidTypeNode"),
		)
	case *ast.WhileExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewWhileExpressionNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				c.astNodeToGoSource(n.Condition, spliceInfo).value,
				astSliceToGoSource(c, "[]ast.StatementNode", n.ThenBody, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.WhileExpressionNode"),
		)
	case *ast.WordArrayListLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewWordArrayListLiteralNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.WordCollectionContentNode", n.Elements, spliceInfo),
				c.astNodeToGoSource(n.Capacity, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.WordArrayListLiteralNode"),
		)
	case *ast.WordArrayTupleLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewWordArrayTupleLiteralNode(%s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.WordCollectionContentNode", n.Elements, spliceInfo),
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.WordArrayTupleLiteralNode"),
		)
	case *ast.WordHashSetLiteralNode:
		return newGoValue(
			fmt.Sprintf("ast.NewWordHashSetLiteralNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				astSliceToGoSource(c, "[]ast.WordCollectionContentNode", n.Elements, spliceInfo),
				c.astNodeToGoSource(n.Capacity, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.WordHashSetLiteralNode"),
		)
	case *ast.YieldExpressionNode:
		return newGoValue(
			fmt.Sprintf("ast.NewYieldExpressionNode(%s, %s, %s)",
				c.locationToGoSource(n.Location()),
				fmt.Sprintf("%t", n.Forward),
				c.astNodeToGoSource(n.Value, spliceInfo).value,
			),
			node.MacroType(c.checker.Env()),
			value.FetchGoType("*ast.YieldExpressionNode"),
		)
	default:
		panic(fmt.Sprintf("invalid ast node when converting to Go source: %T", node))
	}
}
