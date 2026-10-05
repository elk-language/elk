package compiler

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
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

func astSliceToGoSource[T ast.Node](c *GoCompiler, typeName string, nodes []T) string {
	if len(nodes) == 0 {
		return "nil"
	}

	var buff strings.Builder
	fmt.Fprintf(&buff, "%s{", typeName)
	for i, node := range nodes {
		if i > 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(c.astNodeToGoSource(node))
	}
	buff.WriteString("}")
	return buff.String()
}

func (c *GoCompiler) astNodeToGoSourceVal(node ast.Node) *goValue {
	source := c.astNodeToGoSource(node)

	return newGoValue(
		source,
		types.GetType(c.checker.Env().Root, symbol.C_Std, symbol.C_Elk, symbol.C_AST, symbol.C_Node),
		value.FetchGoType("ast.Node"),
	)
}

func (c *GoCompiler) astNodeToGoSource(node ast.Node) string {
	if isNilNode(node) {
		return "nil"
	}

	c.registerAstImport()

	switch n := node.(type) {
	case *ast.AliasDeclarationEntry:
		return fmt.Sprintf("ast.NewAliasDeclarationEntry(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.NewName),
			c.astNodeToGoSource(n.OldName),
		)
	case *ast.AliasDeclarationNode:
		return fmt.Sprintf("ast.NewAliasDeclarationNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]*ast.AliasDeclarationEntry", n.Entries),
		)
	case *ast.AnyTypeNode:
		return fmt.Sprintf("ast.NewAnyTypeNode(%s)",
			c.locationToGoSource(n.Location()),
		)
	case *ast.ArrayListLiteralNode:
		return fmt.Sprintf("ast.NewArrayListLiteralNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.ExpressionNode", n.Elements),
			c.astNodeToGoSource(n.Capacity),
		)
	case *ast.ArrayTupleLiteralNode:
		return fmt.Sprintf("ast.NewArrayTupleLiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.ExpressionNode", n.Elements),
		)
	case *ast.AsExpressionNode:
		return fmt.Sprintf("ast.NewAsExpressionNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Value),
			c.astNodeToGoSource(n.RuntimeType),
		)
	case *ast.AsPatternNode:
		return fmt.Sprintf("ast.NewAsPatternNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Pattern),
			c.astNodeToGoSource(n.Name),
		)
	case *ast.AssignmentExpressionNode:
		return fmt.Sprintf("ast.NewAssignmentExpressionNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.tokenToGoSource(n.Op),
			c.astNodeToGoSource(n.Left),
			c.astNodeToGoSource(n.Right),
		)
	case *ast.AttrDeclarationNode:
		return fmt.Sprintf("ast.NewAttrDeclarationNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.DocComment()),
			c.bitFlag8ToGoSource(n.Flags.ToBitFlag()),
			astSliceToGoSource(c, "[]ast.ParameterNode", n.Entries),
		)
	case *ast.AttributeAccessNode:
		return fmt.Sprintf("ast.NewAttributeAccessNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Receiver),
			c.astNodeToGoSource(n.AttributeName),
		)
	case *ast.AttributeParameterNode:
		return fmt.Sprintf("ast.NewAttributeParameterNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Name),
			c.astNodeToGoSource(n.TypeNode),
			c.astNodeToGoSource(n.Initialiser),
		)
	case *ast.AwaitExpressionNode:
		return fmt.Sprintf("ast.NewAwaitExpressionNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Value),
			fmt.Sprintf("%t", n.Sync),
		)
	case *ast.BigFloatLiteralNode:
		return fmt.Sprintf("ast.NewBigFloatLiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.BinArrayListLiteralNode:
		return fmt.Sprintf("ast.NewBinArrayListLiteralNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.IntCollectionContentNode", n.Elements),
			c.astNodeToGoSource(n.Capacity),
		)
	case *ast.BinArrayTupleLiteralNode:
		return fmt.Sprintf("ast.NewBinArrayTupleLiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.IntCollectionContentNode", n.Elements),
		)
	case *ast.BinHashSetLiteralNode:
		return fmt.Sprintf("ast.NewBinHashSetLiteralNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.IntCollectionContentNode", n.Elements),
			c.astNodeToGoSource(n.Capacity),
		)
	case *ast.BinaryExpressionNode:
		return fmt.Sprintf("ast.NewBinaryExpressionNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.tokenToGoSource(n.Op),
			c.astNodeToGoSource(n.Left),
			c.astNodeToGoSource(n.Right),
		)
	case *ast.BinaryPatternNode:
		return fmt.Sprintf("ast.NewBinaryPatternNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.tokenToGoSource(n.Op),
			c.astNodeToGoSource(n.Left),
			c.astNodeToGoSource(n.Right),
		)
	case *ast.BinaryTypeNode:
		return fmt.Sprintf("ast.NewBinaryTypeNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.tokenToGoSource(n.Op),
			c.astNodeToGoSource(n.Left),
			c.astNodeToGoSource(n.Right),
		)
	case *ast.BoolLiteralNode:
		return fmt.Sprintf("ast.NewBoolLiteralNode(%s)",
			c.locationToGoSource(n.Location()),
		)
	case *ast.BoxOfExpressionNode:
		return fmt.Sprintf("ast.NewBoxOfExpressionNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Expression),
		)
	case *ast.BoxTypeNode:
		return fmt.Sprintf("ast.NewBoxTypeNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.TypeNode),
			fmt.Sprintf("%t", n.Immutable),
		)
	case *ast.BreakExpressionNode:
		return fmt.Sprintf("ast.NewBreakExpressionNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Label),
			c.astNodeToGoSource(n.Value),
		)
	case *ast.BreakpointNode:
		return fmt.Sprintf("ast.NewBreakpointNode(%s)",
			c.locationToGoSource(n.Location()),
		)
	case *ast.CallNode:
		return fmt.Sprintf("ast.NewCallNode(%s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Receiver),
			fmt.Sprintf("%t", n.NilSafe),
			astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments),
			astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments),
		)
	case *ast.CallableTypeNode:
		return fmt.Sprintf("ast.NewCallableTypeNode(%s, %s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.ParameterNode", n.Parameters),
			c.astNodeToGoSource(n.ReturnType),
			c.astNodeToGoSource(n.ThrowType),
			fmt.Sprintf("%t", n.IsClosure),
			fmt.Sprintf("%t", n.IsPure),
		)
	case *ast.CatchNode:
		return fmt.Sprintf("ast.NewCatchNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Pattern),
			c.astNodeToGoSource(n.StackTraceVar),
			astSliceToGoSource(c, "[]ast.StatementNode", n.Body),
		)
	case *ast.CharLiteralNode:
		return fmt.Sprintf("ast.NewCharLiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.ClassDeclarationNode:
		return fmt.Sprintf("ast.NewClassDeclarationNode(%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.DocComment()),
			fmt.Sprintf("%t", n.Abstract),
			fmt.Sprintf("%t", n.Sealed),
			fmt.Sprintf("%t", n.Primitive),
			fmt.Sprintf("%t", n.NoInit),
			fmt.Sprintf("%t", n.Immutable),
			c.astNodeToGoSource(n.Constant),
			astSliceToGoSource(c, "[]ast.TypeParameterNode", n.TypeParameters),
			c.astNodeToGoSource(n.Superclass),
			astSliceToGoSource(c, "[]ast.StatementNode", n.Body),
		)
	case *ast.ClosureLiteralNode:
		return fmt.Sprintf("ast.NewClosureLiteralNode(%s, %s, %s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.ParameterNode", n.Parameters),
			c.astNodeToGoSource(n.ReturnType),
			c.astNodeToGoSource(n.ThrowType),
			astSliceToGoSource(c, "[]ast.StatementNode", n.Body),
			fmt.Sprintf("%t", n.Lambda),
			fmt.Sprintf("%t", n.Pure),
		)
	case *ast.ConstantAsNode:
		return fmt.Sprintf("ast.NewConstantAsNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Constant),
			fmt.Sprintf("%q", n.AsName),
		)
	case *ast.ConstantDeclarationNode:
		return fmt.Sprintf("ast.NewConstantDeclarationNode(%s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.DocComment()),
			c.astNodeToGoSource(n.Constant),
			c.astNodeToGoSource(n.TypeNode),
			c.astNodeToGoSource(n.Initialiser),
		)
	case *ast.ConstantLookupNode:
		return fmt.Sprintf("ast.NewConstantLookupNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Left),
			c.astNodeToGoSource(n.Right),
		)
	case *ast.ConstructorCallNode:
		return fmt.Sprintf("ast.NewConstructorCallNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.ClassNode),
			astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments),
			astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments),
		)
	case *ast.ContinueExpressionNode:
		return fmt.Sprintf("ast.NewContinueExpressionNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Label),
			c.astNodeToGoSource(n.Value),
		)
	case *ast.DeferExpressionNode:
		return fmt.Sprintf("ast.NewDeferExpressionNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Expression),
		)
	case *ast.DoExpressionNode:
		return fmt.Sprintf("ast.NewDoExpressionNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.StatementNode", n.Body),
			astSliceToGoSource(c, "[]*ast.CatchNode", n.Catches),
			astSliceToGoSource(c, "[]ast.StatementNode", n.Finally),
		)
	case *ast.DoubleQuotedStringLiteralNode:
		return fmt.Sprintf("ast.NewDoubleQuotedStringLiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.DoubleSplatExpressionNode:
		return fmt.Sprintf("ast.NewDoubleSplatExpressionNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Value),
		)
	case *ast.EmptyStatementNode:
		return fmt.Sprintf("ast.NewEmptyStatementNode(%s)",
			c.locationToGoSource(n.Location()),
		)
	case *ast.ExactTypeNode:
		return fmt.Sprintf("ast.NewExactTypeNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.TypeNode),
		)
	case *ast.ExpressionStatementNode:
		return fmt.Sprintf("ast.NewExpressionStatementNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Expression),
		)
	case *ast.ExtendWhereBlockExpressionNode:
		return fmt.Sprintf("ast.NewExtendWhereBlockExpressionNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.StatementNode", n.Body),
			astSliceToGoSource(c, "[]ast.TypeParameterNode", n.Where),
		)
	case *ast.FalseLiteralNode:
		return fmt.Sprintf("ast.NewFalseLiteralNode(%s)",
			c.locationToGoSource(n.Location()),
		)
	case *ast.Float32LiteralNode:
		return fmt.Sprintf("ast.NewFloat32LiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.Float64LiteralNode:
		return fmt.Sprintf("ast.NewFloat64LiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.FloatLiteralNode:
		return fmt.Sprintf("ast.NewFloatLiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.ForInExpressionNode:
		return fmt.Sprintf("ast.NewForInExpressionNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Pattern),
			c.astNodeToGoSource(n.InExpression),
			astSliceToGoSource(c, "[]ast.StatementNode", n.ThenBody),
		)
	case *ast.FormalParameterNode:
		return fmt.Sprintf("ast.NewFormalParameterNode(%s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Name),
			c.astNodeToGoSource(n.TypeNode),
			c.astNodeToGoSource(n.Initialiser),
			c.astEnum("ParameterKind", int(n.Kind)),
		)
	case *ast.GenericConstantNode:
		return fmt.Sprintf("ast.NewGenericConstantNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Constant),
			astSliceToGoSource(c, "[]ast.TypeNode", n.TypeArguments),
		)
	case *ast.GenericConstructorCallNode:
		return fmt.Sprintf("ast.NewGenericConstructorCallNode(%s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.ClassNode),
			astSliceToGoSource(c, "[]ast.TypeNode", n.TypeArguments),
			astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments),
			astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments),
		)
	case *ast.GenericMethodCallNode:
		return fmt.Sprintf("ast.NewGenericMethodCallNode(%s, %s, %s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Receiver),
			c.tokenToGoSource(n.Op),
			c.astNodeToGoSource(n.MethodName),
			astSliceToGoSource(c, "[]ast.TypeNode", n.TypeArguments),
			astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments),
			astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments),
		)
	case *ast.GenericReceiverlessMethodCallNode:
		return fmt.Sprintf("ast.NewGenericReceiverlessMethodCallNode(%s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.MethodName),
			astSliceToGoSource(c, "[]ast.TypeNode", n.TypeArguments),
			astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments),
			astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments),
		)
	case *ast.GenericTypeDefinitionNode:
		return fmt.Sprintf("ast.NewGenericTypeDefinitionNode(%s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.DocComment()),
			c.astNodeToGoSource(n.Constant),
			astSliceToGoSource(c, "[]ast.TypeParameterNode", n.TypeParameters),
			c.astNodeToGoSource(n.TypeNode),
		)
	case *ast.GetterDeclarationNode:
		return fmt.Sprintf("ast.NewGetterDeclarationNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.DocComment()),
			c.bitFlag8ToGoSource(n.Flags.ToBitFlag()),
			astSliceToGoSource(c, "[]ast.ParameterNode", n.Entries),
		)
	case *ast.GoExpressionNode:
		return fmt.Sprintf("ast.NewGoExpressionNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.StatementNode", n.Body),
		)
	case *ast.HashMapLiteralNode:
		return fmt.Sprintf("ast.NewHashMapLiteralNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.ExpressionNode", n.Elements),
			c.astNodeToGoSource(n.Capacity),
		)
	case *ast.HashRecordLiteralNode:
		return fmt.Sprintf("ast.NewHashRecordLiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.ExpressionNode", n.Elements),
		)
	case *ast.HashSetLiteralNode:
		return fmt.Sprintf("ast.NewHashSetLiteralNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.ExpressionNode", n.Elements),
			c.astNodeToGoSource(n.Capacity),
		)
	case *ast.HexArrayListLiteralNode:
		return fmt.Sprintf("ast.NewHexArrayListLiteralNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.IntCollectionContentNode", n.Elements),
			c.astNodeToGoSource(n.Capacity),
		)
	case *ast.HexArrayTupleLiteralNode:
		return fmt.Sprintf("ast.NewHexArrayTupleLiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.IntCollectionContentNode", n.Elements),
		)
	case *ast.HexHashSetLiteralNode:
		return fmt.Sprintf("ast.NewHexHashSetLiteralNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.IntCollectionContentNode", n.Elements),
			c.astNodeToGoSource(n.Capacity),
		)
	case *ast.IfExpressionNode:
		return fmt.Sprintf("ast.NewIfExpressionNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Condition),
			astSliceToGoSource(c, "[]ast.StatementNode", n.ThenBody),
			astSliceToGoSource(c, "[]ast.StatementNode", n.ElseBody),
		)
	case *ast.ImplementExpressionNode:
		return fmt.Sprintf("ast.NewImplementExpressionNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.ComplexConstantNode", n.Constants),
		)
	case *ast.ImportStatementNode:
		return fmt.Sprintf("ast.NewImportStatementNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Path),
		)
	case *ast.IncludeExpressionNode:
		return fmt.Sprintf("ast.NewIncludeExpressionNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.ComplexConstantNode", n.Constants),
		)
	case *ast.InferredObjectPatternNode:
		return fmt.Sprintf("ast.NewInferredObjectPatternNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.PatternNode", n.Attributes),
		)
	case *ast.InitDefinitionNode:
		return fmt.Sprintf("ast.NewInitDefinitionNode(%s, %s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.DocComment()),
			c.bitFlag8ToGoSource(n.Flags.ToBitFlag()),
			astSliceToGoSource(c, "[]ast.ParameterNode", n.Parameters),
			c.astNodeToGoSource(n.ThrowType),
			astSliceToGoSource(c, "[]ast.StatementNode", n.Body),
		)
	case *ast.InstanceMethodLookupNode:
		return fmt.Sprintf("ast.NewInstanceMethodLookupNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Receiver),
			c.astNodeToGoSource(n.Name),
		)
	case *ast.InstanceOfTypeNode:
		return fmt.Sprintf("ast.NewInstanceOfTypeNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.TypeNode),
		)
	case *ast.InstanceValueDeclarationNode:
		return fmt.Sprintf("ast.NewInstanceValueDeclarationNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.DocComment()),
			c.astNodeToGoSource(n.Name),
			c.astNodeToGoSource(n.TypeNode),
		)
	case *ast.InstanceVariableDeclarationNode:
		return fmt.Sprintf("ast.NewInstanceVariableDeclarationNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.DocComment()),
			c.astNodeToGoSource(n.Name),
			c.astNodeToGoSource(n.TypeNode),
		)
	case *ast.Int16LiteralNode:
		return fmt.Sprintf("ast.NewInt16LiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.Int32LiteralNode:
		return fmt.Sprintf("ast.NewInt32LiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.Int64LiteralNode:
		return fmt.Sprintf("ast.NewInt64LiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.Int8LiteralNode:
		return fmt.Sprintf("ast.NewInt8LiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.IntLiteralNode:
		return fmt.Sprintf("ast.NewIntLiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.InterfaceDeclarationNode:
		return fmt.Sprintf("ast.NewInterfaceDeclarationNode(%s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.DocComment()),
			c.astNodeToGoSource(n.Constant),
			astSliceToGoSource(c, "[]ast.TypeParameterNode", n.TypeParameters),
			astSliceToGoSource(c, "[]ast.StatementNode", n.Body),
		)
	case *ast.InterpolatedRegexLiteralNode:
		return fmt.Sprintf("ast.NewInterpolatedRegexLiteralNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.RegexLiteralContentNode", n.Content),
			c.bitField8ToGoSource(n.Flags),
		)
	case *ast.InterpolatedStringLiteralNode:
		return fmt.Sprintf("ast.NewInterpolatedStringLiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.StringLiteralContentNode", n.Content),
		)
	case *ast.InterpolatedSymbolLiteralNode:
		return fmt.Sprintf("ast.NewInterpolatedSymbolLiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Content),
		)
	case *ast.IntersectionTypeNode:
		return fmt.Sprintf("ast.NewIntersectionTypeNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.TypeNode", n.Elements),
		)
	case *ast.InvalidNode:
		return fmt.Sprintf("ast.NewInvalidNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.tokenToGoSource(n.Token),
		)
	case *ast.KeyValueExpressionNode:
		return fmt.Sprintf("ast.NewKeyValueExpressionNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Key),
			c.astNodeToGoSource(n.Value),
		)
	case *ast.KeyValuePatternNode:
		return fmt.Sprintf("ast.NewKeyValuePatternNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Key),
			c.astNodeToGoSource(n.Value),
		)
	case *ast.LabeledExpressionNode:
		return fmt.Sprintf("ast.NewLabeledExpressionNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Label),
			c.astNodeToGoSource(n.Expression),
		)
	case *ast.ListPatternNode:
		return fmt.Sprintf("ast.NewListPatternNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.PatternNode", n.Elements),
		)
	case *ast.LogicalExpressionNode:
		return fmt.Sprintf("ast.NewLogicalExpressionNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.tokenToGoSource(n.Op),
			c.astNodeToGoSource(n.Left),
			c.astNodeToGoSource(n.Right),
		)
	case *ast.LoopExpressionNode:
		return fmt.Sprintf("ast.NewLoopExpressionNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.StatementNode", n.ThenBody),
		)
	case *ast.MacroBoundaryNode:
		return fmt.Sprintf("ast.NewMacroBoundaryNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.StatementNode", n.Body),
			fmt.Sprintf("%q", n.Name),
		)
	case *ast.MacroCallNode:
		return fmt.Sprintf("ast.NewMacroCallNode(%s, %s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astEnum("MacroKind", int(n.Kind)),
			c.astNodeToGoSource(n.Receiver),
			c.astNodeToGoSource(n.MacroName),
			astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments),
			astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments),
		)
	case *ast.MacroDefinitionNode:
		return fmt.Sprintf("ast.NewMacroDefinitionNode(%s, %s, %s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.DocComment()),
			fmt.Sprintf("%t", n.IsSealed()),
			c.astNodeToGoSource(n.Name),
			astSliceToGoSource(c, "[]ast.ParameterNode", n.Parameters),
			c.astNodeToGoSource(n.ReturnType),
			astSliceToGoSource(c, "[]ast.StatementNode", n.Body),
		)
	case *ast.MacroNameNode:
		return fmt.Sprintf("ast.NewMacroNameNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.MapPatternNode:
		return fmt.Sprintf("ast.NewMapPatternNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.PatternNode", n.Elements),
		)
	case *ast.MatchExpressionNode:
		return fmt.Sprintf("ast.NewMatchExpressionNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Expression),
			c.astNodeToGoSource(n.Pattern),
		)
	case *ast.MethodCallNode:
		return fmt.Sprintf("ast.NewMethodCallNode(%s, %s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Receiver),
			c.tokenToGoSource(n.Op),
			c.astNodeToGoSource(n.MethodName),
			astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments),
			astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments),
		)
	case *ast.MethodDefinitionNode:
		return fmt.Sprintf("ast.NewMethodDefinitionNode(%s, %s, %s, %s, %s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.DocComment()),
			c.bitFlag8ToGoSource(n.Flags.ToBitFlag()),
			c.astNodeToGoSource(n.Name),
			astSliceToGoSource(c, "[]ast.TypeParameterNode", n.TypeParameters),
			astSliceToGoSource(c, "[]ast.ParameterNode", n.Parameters),
			c.astNodeToGoSource(n.ReturnType),
			c.astNodeToGoSource(n.ThrowType),
			astSliceToGoSource(c, "[]ast.StatementNode", n.Body),
		)
	case *ast.MethodLookupAsNode:
		return fmt.Sprintf("ast.NewMethodLookupAsNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.MethodLookup),
			c.astNodeToGoSource(n.AsName),
		)
	case *ast.MethodLookupNode:
		return fmt.Sprintf("ast.NewMethodLookupNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Receiver),
			c.astNodeToGoSource(n.Name),
		)
	case *ast.MethodParameterNode:
		return fmt.Sprintf("ast.NewMethodParameterNode(%s, %s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Name),
			fmt.Sprintf("%t", n.SetInstanceVariable),
			c.astNodeToGoSource(n.TypeNode),
			c.astNodeToGoSource(n.Initialiser),
			c.astEnum("ParameterKind", int(n.Kind)),
		)
	case *ast.MethodSignatureDefinitionNode:
		return fmt.Sprintf("ast.NewMethodSignatureDefinitionNode(%s, %s, %s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.DocComment()),
			c.astNodeToGoSource(n.Name),
			astSliceToGoSource(c, "[]ast.TypeParameterNode", n.TypeParameters),
			astSliceToGoSource(c, "[]ast.ParameterNode", n.Parameters),
			c.astNodeToGoSource(n.ReturnType),
			c.astNodeToGoSource(n.ThrowType),
		)
	case *ast.MixinDeclarationNode:
		return fmt.Sprintf("ast.NewMixinDeclarationNode(%s, %s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.DocComment()),
			fmt.Sprintf("%t", n.Abstract),
			c.astNodeToGoSource(n.Constant),
			astSliceToGoSource(c, "[]ast.TypeParameterNode", n.TypeParameters),
			astSliceToGoSource(c, "[]ast.StatementNode", n.Body),
		)
	case *ast.ModifierForInNode:
		return fmt.Sprintf("ast.NewModifierForInNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.ThenExpression),
			c.astNodeToGoSource(n.Pattern),
			c.astNodeToGoSource(n.InExpression),
		)
	case *ast.ModifierIfElseNode:
		return fmt.Sprintf("ast.NewModifierIfElseNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.ThenExpression),
			c.astNodeToGoSource(n.Condition),
			c.astNodeToGoSource(n.ElseExpression),
		)
	case *ast.ModifierNode:
		return fmt.Sprintf("ast.NewModifierNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.tokenToGoSource(n.Modifier),
			c.astNodeToGoSource(n.Left),
			c.astNodeToGoSource(n.Right),
		)
	case *ast.ModuleDeclarationNode:
		return fmt.Sprintf("ast.NewModuleDeclarationNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.DocComment()),
			c.astNodeToGoSource(n.Constant),
			astSliceToGoSource(c, "[]ast.StatementNode", n.Body),
		)
	case *ast.MustExpressionNode:
		return fmt.Sprintf("ast.NewMustExpressionNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Value),
		)
	case *ast.MustPatternNode:
		return fmt.Sprintf("ast.NewMustPatternNode(%s)",
			c.locationToGoSource(n.Location()),
		)
	case *ast.NamedCallArgumentNode:
		return fmt.Sprintf("ast.NewNamedCallArgumentNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Name),
			c.astNodeToGoSource(n.Value),
		)
	case *ast.NeverTypeNode:
		return fmt.Sprintf("ast.NewNeverTypeNode(%s)",
			c.locationToGoSource(n.Location()),
		)
	case *ast.NewExpressionNode:
		return fmt.Sprintf("ast.NewNewExpressionNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments),
			astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments),
		)
	case *ast.NilLiteralNode:
		return fmt.Sprintf("ast.NewNilLiteralNode(%s)",
			c.locationToGoSource(n.Location()),
		)
	case *ast.NilSafeSubscriptExpressionNode:
		return fmt.Sprintf("ast.NewNilSafeSubscriptExpressionNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Receiver),
			c.astNodeToGoSource(n.Key),
		)
	case *ast.NilablePatternNode:
		return fmt.Sprintf("ast.NewNilablePatternNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Pattern),
		)
	case *ast.NilableTypeNode:
		return fmt.Sprintf("ast.NewNilableTypeNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.TypeNode),
		)
	case *ast.NotTypeNode:
		return fmt.Sprintf("ast.NewNotTypeNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.TypeNode),
		)
	case *ast.NumericForExpressionNode:
		return fmt.Sprintf("ast.NewNumericForExpressionNode(%s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Initialiser),
			c.astNodeToGoSource(n.Condition),
			c.astNodeToGoSource(n.Increment),
			astSliceToGoSource(c, "[]ast.StatementNode", n.ThenBody),
		)
	case *ast.ObjectPatternNode:
		return fmt.Sprintf("ast.NewObjectPatternNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.ObjectType),
			astSliceToGoSource(c, "[]ast.PatternNode", n.Attributes),
		)
	case *ast.ParameterStatementNode:
		return fmt.Sprintf("ast.NewParameterStatementNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Parameter),
		)
	case *ast.PatternExpressionNode:
		return fmt.Sprintf("ast.NewPatternExpressionNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.PatternNode),
		)
	case *ast.PatternStatementNode:
		return fmt.Sprintf("ast.NewPatternStatementNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Pattern),
		)
	case *ast.PostfixExpressionNode:
		return fmt.Sprintf("ast.NewPostfixExpressionNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.tokenToGoSource(n.Op),
			c.astNodeToGoSource(n.Expression),
		)
	case *ast.PrivateConstantNode:
		return fmt.Sprintf("ast.NewPrivateConstantNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.PrivateIdentifierNode:
		return fmt.Sprintf("ast.NewPrivateIdentifierNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.ProgramNode:
		return fmt.Sprintf("ast.NewProgramNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.StatementNode", n.Body),
		)
	case *ast.PublicConstantAsNode:
		return fmt.Sprintf("ast.NewPublicConstantAsNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Target),
			fmt.Sprintf("%q", n.AsName),
		)
	case *ast.PublicConstantNode:
		return fmt.Sprintf("ast.NewPublicConstantNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.PublicIdentifierNode:
		return fmt.Sprintf("ast.NewPublicIdentifierNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.PublicInstanceVariableNode:
		return fmt.Sprintf("ast.NewPublicInstanceVariableNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.QuoteExpressionNode:
		return fmt.Sprintf("ast.NewQuoteExpressionNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astEnum("QuoteKind", int(n.Kind)),
			astSliceToGoSource(c, "[]ast.StatementNode", n.Body),
		)
	case *ast.RangeLiteralNode:
		return fmt.Sprintf("ast.NewRangeLiteralNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.tokenToGoSource(n.Op),
			c.astNodeToGoSource(n.Start),
			c.astNodeToGoSource(n.End),
		)
	case *ast.RawCharLiteralNode:
		return fmt.Sprintf("ast.NewRawCharLiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.RawStringLiteralNode:
		return fmt.Sprintf("ast.NewRawStringLiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.ReceiverlessMacroCallNode:
		return fmt.Sprintf("ast.NewReceiverlessMacroCallNode(%s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astEnum("MacroKind", int(n.Kind)),
			c.astNodeToGoSource(n.MacroName),
			astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments),
			astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments),
		)
	case *ast.ReceiverlessMethodCallNode:
		return fmt.Sprintf("ast.NewReceiverlessMethodCallNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.MethodName),
			astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments),
			astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments),
		)
	case *ast.RecordPatternNode:
		return fmt.Sprintf("ast.NewRecordPatternNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.PatternNode", n.Elements),
		)
	case *ast.RegexInterpolationNode:
		return fmt.Sprintf("ast.NewRegexInterpolationNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Expression),
		)
	case *ast.RegexLiteralContentSectionNode:
		return fmt.Sprintf("ast.NewRegexLiteralContentSectionNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.RestPatternNode:
		return fmt.Sprintf("ast.NewRestPatternNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Identifier),
		)
	case *ast.ReturnExpressionNode:
		return fmt.Sprintf("ast.NewReturnExpressionNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Value),
		)
	case *ast.ScopedMacroCallNode:
		return fmt.Sprintf("ast.NewScopedMacroCallNode(%s, %s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astEnum("MacroKind", int(n.Kind)),
			c.astNodeToGoSource(n.Receiver),
			c.astNodeToGoSource(n.MacroName),
			astSliceToGoSource(c, "[]ast.ExpressionNode", n.PositionalArguments),
			astSliceToGoSource(c, "[]ast.NamedArgumentNode", n.NamedArguments),
		)
	case *ast.SelectCaseNode:
		return fmt.Sprintf("ast.NewSelectCaseNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Expression),
			astSliceToGoSource(c, "[]ast.StatementNode", n.Body),
		)
	case *ast.SelectExpressionNode:
		return fmt.Sprintf("ast.NewSelectExpressionNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]*ast.SelectCaseNode", n.Cases),
			astSliceToGoSource(c, "[]ast.StatementNode", n.ElseBody),
		)
	case *ast.SelfLiteralNode:
		return fmt.Sprintf("ast.NewSelfLiteralNode(%s)",
			c.locationToGoSource(n.Location()),
		)
	case *ast.SetPatternNode:
		return fmt.Sprintf("ast.NewSetPatternNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.PatternNode", n.Elements),
		)
	case *ast.SetterDeclarationNode:
		return fmt.Sprintf("ast.NewSetterDeclarationNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.DocComment()),
			astSliceToGoSource(c, "[]ast.ParameterNode", n.Entries),
		)
	case *ast.SignatureParameterNode:
		return fmt.Sprintf("ast.NewSignatureParameterNode(%s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Name),
			c.astNodeToGoSource(n.TypeNode),
			fmt.Sprintf("%t", n.Optional),
			c.astEnum("ParameterKind", int(n.Kind)),
		)
	case *ast.SimpleSymbolLiteralNode:
		return fmt.Sprintf("ast.NewSimpleSymbolLiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Content),
		)
	case *ast.SingletonBlockExpressionNode:
		return fmt.Sprintf("ast.NewSingletonBlockExpressionNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.StatementNode", n.Body),
		)
	case *ast.SingletonTypeNode:
		return fmt.Sprintf("ast.NewSingletonTypeNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.TypeNode),
		)
	case *ast.SplatExpressionNode:
		return fmt.Sprintf("ast.NewSplatExpressionNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Value),
		)
	case *ast.StringInspectInterpolationNode:
		return fmt.Sprintf("ast.NewStringInspectInterpolationNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Expression),
		)
	case *ast.StringInterpolationNode:
		return fmt.Sprintf("ast.NewStringInterpolationNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Expression),
		)
	case *ast.StringLiteralContentSectionNode:
		return fmt.Sprintf("ast.NewStringLiteralContentSectionNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.StructDeclarationNode:
		return fmt.Sprintf("ast.NewStructDeclarationNode(%s, %s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.DocComment()),
			fmt.Sprintf("%t", n.Immutable),
			c.astNodeToGoSource(n.Constant),
			astSliceToGoSource(c, "[]ast.TypeParameterNode", n.TypeParameters),
			astSliceToGoSource(c, "[]ast.StructBodyStatementNode", n.Body),
		)
	case *ast.SubscriptExpressionNode:
		return fmt.Sprintf("ast.NewSubscriptExpressionNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Receiver),
			c.astNodeToGoSource(n.Key),
		)
	case *ast.SwitchCaseNode:
		return fmt.Sprintf("ast.NewSwitchCaseNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Pattern),
			astSliceToGoSource(c, "[]ast.StatementNode", n.Body),
		)
	case *ast.SwitchExpressionNode:
		return fmt.Sprintf("ast.NewSwitchExpressionNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Value),
			astSliceToGoSource(c, "[]*ast.SwitchCaseNode", n.Cases),
			astSliceToGoSource(c, "[]ast.StatementNode", n.ElseBody),
		)
	case *ast.SymbolArrayListLiteralNode:
		return fmt.Sprintf("ast.NewSymbolArrayListLiteralNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.SymbolCollectionContentNode", n.Elements),
			c.astNodeToGoSource(n.Capacity),
		)
	case *ast.SymbolArrayTupleLiteralNode:
		return fmt.Sprintf("ast.NewSymbolArrayTupleLiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.SymbolCollectionContentNode", n.Elements),
		)
	case *ast.SymbolHashSetLiteralNode:
		return fmt.Sprintf("ast.NewSymbolHashSetLiteralNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.SymbolCollectionContentNode", n.Elements),
			c.astNodeToGoSource(n.Capacity),
		)
	case *ast.SymbolKeyValueExpressionNode:
		return fmt.Sprintf("ast.NewSymbolKeyValueExpressionNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Key),
			c.astNodeToGoSource(n.Value),
		)
	case *ast.SymbolKeyValuePatternNode:
		return fmt.Sprintf("ast.NewSymbolKeyValuePatternNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Key),
			c.astNodeToGoSource(n.Value),
		)
	case *ast.ThrowExpressionNode:
		return fmt.Sprintf("ast.NewThrowExpressionNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%t", n.Unchecked),
			c.astNodeToGoSource(n.Value),
		)
	case *ast.TrueLiteralNode:
		return fmt.Sprintf("ast.NewTrueLiteralNode(%s)",
			c.locationToGoSource(n.Location()),
		)
	case *ast.TryExpressionNode:
		return fmt.Sprintf("ast.NewTryExpressionNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Value),
		)
	case *ast.TuplePatternNode:
		return fmt.Sprintf("ast.NewTuplePatternNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.PatternNode", n.Elements),
		)
	case *ast.TypeDefinitionNode:
		return fmt.Sprintf("ast.NewTypeDefinitionNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.DocComment()),
			c.astNodeToGoSource(n.Constant),
			c.astNodeToGoSource(n.TypeNode),
		)
	case *ast.TypeExpressionNode:
		return fmt.Sprintf("ast.NewTypeExpressionNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.TypeNode),
		)
	case *ast.TypeStatementNode:
		return fmt.Sprintf("ast.NewTypeStatementNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.TypeNode),
		)
	case *ast.TypeofExpressionNode:
		return fmt.Sprintf("ast.NewTypeofExpressionNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Value),
		)
	case *ast.UInt16LiteralNode:
		return fmt.Sprintf("ast.NewUInt16LiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.UInt32LiteralNode:
		return fmt.Sprintf("ast.NewUInt32LiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.UInt64LiteralNode:
		return fmt.Sprintf("ast.NewUInt64LiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.UInt8LiteralNode:
		return fmt.Sprintf("ast.NewUInt8LiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.UIntLiteralNode:
		return fmt.Sprintf("ast.NewUIntLiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Value),
		)
	case *ast.UnaryExpressionNode:
		return fmt.Sprintf("ast.NewUnaryExpressionNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.tokenToGoSource(n.Op),
			c.astNodeToGoSource(n.Right),
		)
	case *ast.UnaryTypeNode:
		return fmt.Sprintf("ast.NewUnaryTypeNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.tokenToGoSource(n.Op),
			c.astNodeToGoSource(n.TypeNode),
		)
	case *ast.UndefinedLiteralNode:
		return fmt.Sprintf("ast.NewUndefinedLiteralNode(%s)",
			c.locationToGoSource(n.Location()),
		)
	case *ast.UnhygienicNode:
		return fmt.Sprintf("ast.NewUnhygienicNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Node),
		)
	case *ast.UninterpolatedRegexLiteralNode:
		return fmt.Sprintf("ast.NewUninterpolatedRegexLiteralNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.Content),
			c.bitField8ToGoSource(n.Flags),
		)
	case *ast.UnionTypeNode:
		return fmt.Sprintf("ast.NewUnionTypeNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.TypeNode", n.Elements),
		)
	case *ast.UnlessExpressionNode:
		return fmt.Sprintf("ast.NewUnlessExpressionNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Condition),
			astSliceToGoSource(c, "[]ast.StatementNode", n.ThenBody),
			astSliceToGoSource(c, "[]ast.StatementNode", n.ElseBody),
		)
	case *ast.UnquoteForInExpressionNode:
		return fmt.Sprintf("ast.NewUnquoteForInExpressionNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Pattern),
			c.astNodeToGoSource(n.InExpression),
			astSliceToGoSource(c, "[]ast.StatementNode", n.ThenBody),
		)
	case *ast.UnquoteIfExpressionNode:
		return fmt.Sprintf("ast.NewUnquoteIfExpressionNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Condition),
			astSliceToGoSource(c, "[]ast.StatementNode", n.ThenBody),
			astSliceToGoSource(c, "[]ast.StatementNode", n.ElseBody),
		)
	case *ast.UnquoteNode:
		return fmt.Sprintf("ast.NewUnquoteNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astEnum("UnquoteKind", int(n.Kind)),
			c.astNodeToGoSource(n.Expression),
		)
	case *ast.UntilExpressionNode:
		return fmt.Sprintf("ast.NewUntilExpressionNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Condition),
			astSliceToGoSource(c, "[]ast.StatementNode", n.ThenBody),
		)
	case *ast.UsingAllEntryNode:
		return fmt.Sprintf("ast.NewUsingAllEntryNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Namespace),
		)
	case *ast.UsingEntryWithSubentriesNode:
		return fmt.Sprintf("ast.NewUsingEntryWithSubentriesNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Namespace),
			astSliceToGoSource(c, "[]ast.UsingSubentryNode", n.Subentries),
		)
	case *ast.UsingExpressionNode:
		return fmt.Sprintf("ast.NewUsingExpressionNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.UsingEntryNode", n.Entries),
		)
	case *ast.UsingSubentryAsNode:
		return fmt.Sprintf("ast.NewUsingSubentryAsNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Target),
			c.astNodeToGoSource(n.AsName),
		)
	case *ast.ValueDeclarationNode:
		return fmt.Sprintf("ast.NewValueDeclarationNode(%s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Name),
			c.astNodeToGoSource(n.TypeNode),
			c.astNodeToGoSource(n.Initialiser),
		)
	case *ast.ValuePatternDeclarationNode:
		return fmt.Sprintf("ast.NewValuePatternDeclarationNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Pattern),
			c.astNodeToGoSource(n.Initialiser),
		)
	case *ast.VariableDeclarationNode:
		return fmt.Sprintf("ast.NewVariableDeclarationNode(%s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%q", n.DocComment()),
			c.astNodeToGoSource(n.Name),
			c.astNodeToGoSource(n.TypeNode),
			c.astNodeToGoSource(n.Initialiser),
		)
	case *ast.VariablePatternDeclarationNode:
		return fmt.Sprintf("ast.NewVariablePatternDeclarationNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Pattern),
			c.astNodeToGoSource(n.Initialiser),
		)
	case *ast.VariantTypeParameterNode:
		return fmt.Sprintf("ast.NewVariantTypeParameterNode(%s, %s, %s, %s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astEnum("Variance", int(n.Variance)),
			fmt.Sprintf("%q", n.Name),
			c.astNodeToGoSource(n.LowerBound),
			c.astNodeToGoSource(n.UpperBound),
			c.astNodeToGoSource(n.Default),
		)
	case *ast.VoidTypeNode:
		return fmt.Sprintf("ast.NewVoidTypeNode(%s)",
			c.locationToGoSource(n.Location()),
		)
	case *ast.WhileExpressionNode:
		return fmt.Sprintf("ast.NewWhileExpressionNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			c.astNodeToGoSource(n.Condition),
			astSliceToGoSource(c, "[]ast.StatementNode", n.ThenBody),
		)
	case *ast.WordArrayListLiteralNode:
		return fmt.Sprintf("ast.NewWordArrayListLiteralNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.WordCollectionContentNode", n.Elements),
			c.astNodeToGoSource(n.Capacity),
		)
	case *ast.WordArrayTupleLiteralNode:
		return fmt.Sprintf("ast.NewWordArrayTupleLiteralNode(%s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.WordCollectionContentNode", n.Elements),
		)
	case *ast.WordHashSetLiteralNode:
		return fmt.Sprintf("ast.NewWordHashSetLiteralNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			astSliceToGoSource(c, "[]ast.WordCollectionContentNode", n.Elements),
			c.astNodeToGoSource(n.Capacity),
		)
	case *ast.YieldExpressionNode:
		return fmt.Sprintf("ast.NewYieldExpressionNode(%s, %s, %s)",
			c.locationToGoSource(n.Location()),
			fmt.Sprintf("%t", n.Forward),
			c.astNodeToGoSource(n.Value),
		)
	default:
		panic(fmt.Sprintf("invalid ast node when converting to Go source: %T", node))
	}
}
