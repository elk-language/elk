// Package checker implements the Elk type checker
package checker

import (
	"fmt"

	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	typed "github.com/elk-language/elk/types/ast" // typed AST
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// Check the types of Elk source code.
func CheckSource(sourceName string, source string, globalEnv *types.GlobalEnvironment, headerMode bool) (typed.Node, errors.ErrorList) {
	ast, err := parser.Parse(sourceName, source)
	if err != nil {
		return nil, err
	}

	return CheckAST(sourceName, ast, globalEnv, headerMode)
}

// Check the types of an Elk AST.
func CheckAST(sourceName string, ast *ast.ProgramNode, globalEnv *types.GlobalEnvironment, headerMode bool) (typed.Node, errors.ErrorList) {
	checker := newChecker(position.NewLocationWithSpan(sourceName, ast.Span()), globalEnv, headerMode)
	typedAst := checker.checkProgram(ast)
	return typedAst, checker.Errors
}

// Represents a single local variable or local value
type local struct {
	typ              types.Type
	initialised      bool
	singleAssignment bool
}

// Contains definitions of local variables and values
type localEnvironment struct {
	parent *localEnvironment
	locals map[value.Symbol]local
}

// Get the local with the specified name from this local environment
func (l *localEnvironment) getLocal(name string) (local, bool) {
	local, ok := l.locals[value.ToSymbol(name)]
	return local, ok
}

// Resolve the local with the given name from this local environment or any parent environment
func (l *localEnvironment) resolveLocal(name string) (local, bool) {
	nameSymbol := value.ToSymbol(name)
	currentEnv := l
	for {
		if currentEnv == nil {
			return local{}, false
		}
		loc, ok := currentEnv.locals[nameSymbol]
		if ok {
			return loc, true
		}
		currentEnv = currentEnv.parent
	}
}

func newLocalEnvironment(parent *localEnvironment) *localEnvironment {
	return &localEnvironment{
		parent: parent,
		locals: make(map[value.Symbol]local),
	}
}

type constantScope struct {
	container types.ConstantContainer
	local     bool
}

func makeLocalConstantScope(container types.ConstantContainer) constantScope {
	return constantScope{
		container: container,
		local:     true,
	}
}

func makeConstantScope(container types.ConstantContainer) constantScope {
	return constantScope{
		container: container,
		local:     false,
	}
}

// Holds the state of the type checking process
type Checker struct {
	Location       *position.Location
	Errors         errors.ErrorList
	GlobalEnv      *types.GlobalEnvironment
	HeaderMode     bool
	constantScopes []constantScope
	localEnvs      []*localEnvironment
	returnType     types.Type
	throwType      types.Type
	selfType       types.Type
}

// Instantiate a new Checker instance.
func newChecker(loc *position.Location, globalEnv *types.GlobalEnvironment, headerMode bool) *Checker {
	if globalEnv == nil {
		globalEnv = types.NewGlobalEnvironment()
	}
	return &Checker{
		Location:   loc,
		GlobalEnv:  globalEnv,
		HeaderMode: headerMode,
		selfType:   globalEnv.StdSubtype(symbol.Object),
		returnType: types.Any{},
		constantScopes: []constantScope{
			makeConstantScope(globalEnv.Std()),
			makeLocalConstantScope(globalEnv.Root),
		},
		localEnvs: []*localEnvironment{
			newLocalEnvironment(nil),
		},
	}
}

// Instantiate a new Checker instance.
func New() *Checker {
	globalEnv := types.NewGlobalEnvironment()
	return &Checker{
		GlobalEnv: globalEnv,
		constantScopes: []constantScope{
			makeConstantScope(globalEnv.Std()),
			makeLocalConstantScope(globalEnv.Root),
		},
		localEnvs: []*localEnvironment{
			newLocalEnvironment(nil),
		},
	}
}

func (c *Checker) CheckSource(sourceName string, source string) (typed.Node, errors.ErrorList) {
	ast, err := parser.Parse(sourceName, source)
	if err != nil {
		return nil, err
	}

	loc := position.NewLocationWithSpan(sourceName, ast.Span())
	c.Location = loc
	typedAst := c.checkProgram(ast)
	return typedAst, c.Errors
}

func (c *Checker) ClearErrors() {
	c.Errors = nil
}

func (c *Checker) popLocalEnv() {
	c.localEnvs = c.localEnvs[:len(c.localEnvs)-1]
}

func (c *Checker) pushLocalEnv(env *localEnvironment) {
	c.localEnvs = append(c.localEnvs, env)
}

func (c *Checker) currentLocalEnv() *localEnvironment {
	return c.localEnvs[len(c.localEnvs)-1]
}

func (c *Checker) popConstScope() {
	c.constantScopes = c.constantScopes[:len(c.constantScopes)-1]
}

func (c *Checker) pushConstScope(constScope constantScope) {
	c.constantScopes = append(c.constantScopes, constScope)
}

func (c *Checker) currentConstScope() constantScope {
	for i := range len(c.constantScopes) {
		constScope := c.constantScopes[len(c.constantScopes)-i-1]
		if constScope.local {
			return constScope
		}
	}

	panic("no local constant scopes!")
}

// Create a new location struct with the given position.
func (c *Checker) newLocation(span *position.Span) *position.Location {
	return position.NewLocationWithSpan(c.Location.Filename, span)
}

func (c *Checker) checkProgram(node *ast.ProgramNode) *typed.ProgramNode {
	newStatements := c.checkStatements(node.Body)
	return typed.NewProgramNode(
		node.Span(),
		newStatements,
	)
}

func (c *Checker) checkStatements(stmts []ast.StatementNode) []typed.StatementNode {
	var newStatements []typed.StatementNode
	for _, statement := range stmts {
		switch statement.(type) {
		case *ast.EmptyStatementNode:
			continue
		}
		newStatements = append(newStatements, c.checkStatement(statement))
	}

	return newStatements
}

func (c *Checker) checkStatement(node ast.Node) typed.StatementNode {
	switch node := node.(type) {
	case *ast.ExpressionStatementNode:
		expr := c.checkExpression(node.Expression)
		return typed.NewExpressionStatementNode(
			node.Span(),
			expr,
		)
	default:
		c.addError(
			fmt.Sprintf("incorrect statement type %#v", node),
			node.Span(),
		)
		return typed.NewInvalidNode(node.Span(), nil)
	}
}

func (c *Checker) isSubtype(a, b types.Type) bool {
	if types.IsVoid(a) || types.IsNever(a) || types.IsVoid(b) || types.IsNever(b) {
		return false
	}
	if types.IsAny(b) {
		return true
	}
	aNonLiteral := a.ToNonLiteral(c.GlobalEnv)
	if a != aNonLiteral && c.isSubtype(aNonLiteral, b) {
		return true
	}

	if aNilable, aIsNilable := a.(*types.Nilable); aIsNilable {
		return c.isSubtype(aNilable.Type, b) && c.isSubtype(c.GlobalEnv.StdSubtype(symbol.Nil), b)
	}

	if bNilable, bIsNilable := b.(*types.Nilable); bIsNilable {
		return c.isSubtype(a, bNilable.Type) || c.isSubtype(a, c.GlobalEnv.StdSubtype(symbol.Nil))
	}

	if aUnion, aIsUnion := a.(*types.Union); aIsUnion {
		for _, aElement := range aUnion.Elements {
			if !c.isSubtype(aElement, b) {
				return false
			}
		}
		return true
	}

	if bUnion, bIsUnion := b.(*types.Union); bIsUnion {
		for _, bElement := range bUnion.Elements {
			if c.isSubtype(a, bElement) {
				return true
			}
		}
		return false
	}

	originalA := a
	switch a := a.(type) {
	case types.Any:
		return types.IsAny(b)
	case *types.SingletonClass:
		b, ok := b.(*types.SingletonClass)
		if !ok {
			return false
		}
		return a.AttachedObject == b.AttachedObject
	case *types.Class:
		b, ok := b.(*types.Class)
		if !ok {
			return false
		}

		var currentClass types.ConstantContainer = a
		for {
			if currentClass == nil {
				return false
			}
			if currentClass == b {
				return true
			}

			currentClass = currentClass.Parent()
		}
	case *types.Module:
		b, ok := b.(*types.Module)
		if !ok {
			return false
		}
		return a == b
	case *types.Method:
		b, ok := b.(*types.Method)
		if !ok {
			return false
		}
		return a == b
	case *types.CharLiteral:
		b, ok := b.(*types.CharLiteral)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.StringLiteral:
		b, ok := b.(*types.StringLiteral)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.SymbolLiteral:
		b, ok := b.(*types.SymbolLiteral)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.FloatLiteral:
		b, ok := b.(*types.FloatLiteral)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.Float64Literal:
		b, ok := b.(*types.Float64Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.Float32Literal:
		b, ok := b.(*types.Float32Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.BigFloatLiteral:
		b, ok := b.(*types.BigFloatLiteral)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.IntLiteral:
		b, ok := b.(*types.IntLiteral)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.Int64Literal:
		b, ok := b.(*types.Int64Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.Int32Literal:
		b, ok := b.(*types.Int32Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.Int16Literal:
		b, ok := b.(*types.Int16Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.Int8Literal:
		b, ok := b.(*types.Int8Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.UInt64Literal:
		b, ok := b.(*types.UInt64Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.UInt32Literal:
		b, ok := b.(*types.UInt32Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.UInt16Literal:
		b, ok := b.(*types.UInt16Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.UInt8Literal:
		b, ok := b.(*types.UInt8Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	default:
		panic(fmt.Sprintf("invalid type: %T", originalA))
	}
}

func (c *Checker) checkExpressions(exprs []ast.ExpressionNode) []typed.ExpressionNode {
	var newExpressions []typed.ExpressionNode
	for _, expr := range exprs {
		newExpressions = append(newExpressions, c.checkExpression(expr))
	}

	return newExpressions
}

func (c *Checker) checkExpression(node ast.ExpressionNode) typed.ExpressionNode {
	switch n := node.(type) {
	case *ast.FalseLiteralNode:
		return typed.NewFalseLiteralNode(n.Span())
	case *ast.TrueLiteralNode:
		return typed.NewTrueLiteralNode(n.Span())
	case *ast.NilLiteralNode:
		return typed.NewNilLiteralNode(n.Span())
	case *ast.TypeExpressionNode:
		typeNode := c.checkTypeNode(n.TypeNode)
		return typed.NewTypeExpressionNode(n.Span(), typeNode)
	case *ast.IntLiteralNode:
		return typed.NewIntLiteralNode(
			n.Span(),
			n.Value,
			types.NewIntLiteral(n.Value),
		)
	case *ast.Int64LiteralNode:
		return typed.NewInt64LiteralNode(
			n.Span(),
			n.Value,
			types.NewInt64Literal(n.Value),
		)
	case *ast.Int32LiteralNode:
		return typed.NewInt32LiteralNode(
			n.Span(),
			n.Value,
			types.NewInt32Literal(n.Value),
		)
	case *ast.Int16LiteralNode:
		return typed.NewInt16LiteralNode(
			n.Span(),
			n.Value,
			types.NewInt16Literal(n.Value),
		)
	case *ast.Int8LiteralNode:
		return typed.NewInt8LiteralNode(
			n.Span(),
			n.Value,
			types.NewInt8Literal(n.Value),
		)
	case *ast.UInt64LiteralNode:
		return typed.NewUInt64LiteralNode(
			n.Span(),
			n.Value,
			types.NewUInt64Literal(n.Value),
		)
	case *ast.UInt32LiteralNode:
		return typed.NewUInt32LiteralNode(
			n.Span(),
			n.Value,
			types.NewUInt32Literal(n.Value),
		)
	case *ast.UInt16LiteralNode:
		return typed.NewUInt16LiteralNode(
			n.Span(),
			n.Value,
			types.NewUInt16Literal(n.Value),
		)
	case *ast.UInt8LiteralNode:
		return typed.NewUInt8LiteralNode(
			n.Span(),
			n.Value,
			types.NewUInt8Literal(n.Value),
		)
	case *ast.FloatLiteralNode:
		return typed.NewFloatLiteralNode(
			n.Span(),
			n.Value,
			types.NewFloatLiteral(n.Value),
		)
	case *ast.Float64LiteralNode:
		return typed.NewFloat64LiteralNode(
			n.Span(),
			n.Value,
			types.NewFloat64Literal(n.Value),
		)
	case *ast.Float32LiteralNode:
		return typed.NewFloat32LiteralNode(
			n.Span(),
			n.Value,
			types.NewFloat32Literal(n.Value),
		)
	case *ast.BigFloatLiteralNode:
		return typed.NewBigFloatLiteralNode(
			n.Span(),
			n.Value,
			types.NewBigFloatLiteral(n.Value),
		)
	case *ast.DoubleQuotedStringLiteralNode:
		return typed.NewDoubleQuotedStringLiteralNode(
			n.Span(),
			n.Value,
			types.NewStringLiteral(n.Value),
		)
	case *ast.RawStringLiteralNode:
		return typed.NewRawStringLiteralNode(
			n.Span(),
			n.Value,
			types.NewStringLiteral(n.Value),
		)
	case *ast.RawCharLiteralNode:
		return typed.NewRawCharLiteralNode(
			n.Span(),
			n.Value,
			types.NewCharLiteral(n.Value),
		)
	case *ast.CharLiteralNode:
		return typed.NewCharLiteralNode(
			n.Span(),
			n.Value,
			types.NewCharLiteral(n.Value),
		)
	case *ast.InterpolatedStringLiteralNode:
		return c.interpolatedStringLiteral(n)
	case *ast.SimpleSymbolLiteralNode:
		return typed.NewSimpleSymbolLiteralNode(
			n.Span(),
			n.Content,
			types.NewSymbolLiteral(n.Content),
		)
	case *ast.InterpolatedSymbolLiteralNode:
		return typed.NewInterpolatedSymbolLiteralNode(
			n.Span(),
			c.interpolatedStringLiteral(n.Content),
		)
	case *ast.VariableDeclarationNode:
		return c.variableDeclaration(n)
	case *ast.ValueDeclarationNode:
		return c.valueDeclaration(n)
	case *ast.ConstantDeclarationNode:
		return c.constantDeclaration(n)
	case *ast.PublicIdentifierNode:
		return c.publicIdentifier(n)
	case *ast.PrivateIdentifierNode:
		return c.privateIdentifier(n)
	case *ast.PublicConstantNode:
		return c.publicConstant(n)
	case *ast.PrivateConstantNode:
		return c.privateConstant(n)
	case *ast.ConstantLookupNode:
		return c.constantLookup(n)
	case *ast.ModuleDeclarationNode:
		return c.module(n)
	case *ast.MethodDefinitionNode:
		return c.methodDefinition(n)
	case *ast.AssignmentExpressionNode:
		return c.assignmentExpression(n)
	case *ast.ReceiverlessMethodCallNode:
		return c.receiverlessMethodCall(n)
	case *ast.MethodCallNode:
		return c.methodCall(n)
	case *ast.AttributeAccessNode:
		return c.attributeAccess(n)
	default:
		c.addError(
			fmt.Sprintf("invalid expression type %T", node),
			node.Span(),
		)
		return typed.NewInvalidNode(node.Span(), nil)
	}
}

func (c *Checker) typeOf(node typed.Node) types.Type {
	return typed.TypeOf(node, c.GlobalEnv)
}

func (c *Checker) toNilable(typ types.Type) types.Type {
	return types.ToNilable(typ, c.GlobalEnv)
}

func (c *Checker) checkMethodArguments(method *types.Method, positionalArguments []ast.ExpressionNode, namedArguments []ast.NamedArgumentNode, span *position.Span) []typed.ExpressionNode {
	reqParamCount := method.RequiredParamCount()
	requiredPosParamCount := len(method.Params) - method.OptionalParamCount - method.PostParamCount
	positionalRestParamIndex := method.PositionalRestParamIndex()
	var typedPositionalArguments []typed.ExpressionNode

	var currentParamIndex int
	var posArg ast.ExpressionNode
	for currentParamIndex, posArg = range positionalArguments {
		if currentParamIndex == positionalRestParamIndex {
			break
		}
		if currentParamIndex >= len(method.Params) {
			c.addWrongArgumentCountError(
				len(positionalArguments)+len(namedArguments),
				method,
				span,
			)
			break
		}
		param := method.Params[currentParamIndex]

		typedPosArg := c.checkExpression(posArg)
		typedPositionalArguments = append(typedPositionalArguments, typedPosArg)
		posArgType := c.typeOf(typedPosArg)
		if !c.isSubtype(posArgType, param.Type) {
			c.addError(
				fmt.Sprintf(
					"expected type `%s` for parameter `%s`, got type `%s`",
					types.Inspect(param.Type),
					param.Name,
					types.Inspect(posArgType),
				),
				posArg.Span(),
			)
		}
	}

	if method.HasPositionalRestParam() {
		if len(positionalArguments) < requiredPosParamCount {
			c.addError(
				fmt.Sprintf(
					"expected %d... positional arguments, got %d",
					requiredPosParamCount,
					len(positionalArguments),
				),
				span,
			)
			return nil
		}
		restPositionalArguments := typed.NewArrayTupleLiteralNode(
			span,
			nil,
		)
		posRestParam := method.Params[positionalRestParamIndex]

		for j := currentParamIndex; j < len(positionalArguments)-method.PostParamCount; j++ {
			posArg := positionalArguments[j]
			typedPosArg := c.checkExpression(posArg)
			restPositionalArguments.Elements = append(restPositionalArguments.Elements, typedPosArg)
			posArgType := c.typeOf(typedPosArg)
			if !c.isSubtype(posArgType, posRestParam.Type) {
				c.addError(
					fmt.Sprintf(
						"expected type `%s` for rest parameter `*%s`, got type `%s`",
						types.Inspect(posRestParam.Type),
						posRestParam.Name,
						types.Inspect(posArgType),
					),
					posArg.Span(),
				)
			}
		}
		typedPositionalArguments = append(typedPositionalArguments, restPositionalArguments)

		currentParamIndex = positionalRestParamIndex
		for i := len(positionalArguments) - method.PostParamCount; i < len(positionalArguments); i++ {
			posArg := positionalArguments[i]
			currentParamIndex++
			param := method.Params[currentParamIndex]

			typedPosArg := c.checkExpression(posArg)
			typedPositionalArguments = append(typedPositionalArguments, typedPosArg)
			posArgType := c.typeOf(typedPosArg)
			if !c.isSubtype(posArgType, param.Type) {
				c.addError(
					fmt.Sprintf(
						"expected type `%s` for parameter `%s`, got type `%s`",
						types.Inspect(param.Type),
						param.Name,
						types.Inspect(posArgType),
					),
					posArg.Span(),
				)
			}
		}
	}

	firstNamedParamIndex := currentParamIndex + 1
	definedNamedArgumentsSlice := make([]bool, len(namedArguments))

	for i := 0; i < len(method.Params); i++ {
		param := method.Params[i]
		paramName := param.Name.String()
		var found bool

		for namedArgIndex, namedArgI := range namedArguments {
			namedArg := namedArgI.(*ast.NamedCallArgumentNode)
			if namedArg.Name != paramName {
				continue
			}
			found = true
			if found || i < firstNamedParamIndex {
				c.addError(
					fmt.Sprintf(
						"duplicated argument `%s` in call to `%s`",
						paramName,
						method.Name,
					),
					namedArg.Span(),
				)
			}
			definedNamedArgumentsSlice[namedArgIndex] = true
			typedNamedArgValue := c.checkExpression(namedArg.Value)
			namedArgType := c.typeOf(typedNamedArgValue)
			typedPositionalArguments = append(typedPositionalArguments, typedNamedArgValue)
			if !c.isSubtype(namedArgType, param.Type) {
				c.addError(
					fmt.Sprintf(
						"expected type `%s` for parameter `%s`, got type `%s`",
						types.Inspect(param.Type),
						param.Name,
						types.Inspect(namedArgType),
					),
					namedArg.Span(),
				)
			}
		}

		if i < firstNamedParamIndex {
			continue
		}
		if found {
			continue
		}

		if i < reqParamCount {
			// the parameter is required
			// but is not present in the call
			c.addError(
				fmt.Sprintf(
					"argument `%s` is missing in call to `%s`",
					paramName,
					method.Name,
				),
				span,
			)
		} else {
			// the parameter is missing and is optional
			// we push undefined as its value
			typedPositionalArguments = append(
				typedPositionalArguments,
				typed.NewUndefinedLiteralNode(span),
			)
		}
	}

	if method.HasNamedRestParam {
		namedRestArgs := typed.NewHashRecordLiteralNode(
			span,
			nil,
		)
		namedRestParam := method.Params[len(method.Params)-1]
		for i, defined := range definedNamedArgumentsSlice {
			if defined {
				continue
			}

			namedArgI := namedArguments[i]
			namedArg := namedArgI.(*ast.NamedCallArgumentNode)
			typedNamedArgValue := c.checkExpression(namedArg.Value)
			namedRestArgs.Elements = append(
				namedRestArgs.Elements,
				typed.NewSymbolKeyValueExpressionNode(
					namedArg.Span(),
					namedArg.Name,
					typedNamedArgValue,
				),
			)
			namedArgType := c.typeOf(typedNamedArgValue)
			if !c.isSubtype(namedArgType, namedRestParam.Type) {
				c.addError(
					fmt.Sprintf(
						"expected type `%s` for named rest parameter `**%s`, got type `%s`",
						types.Inspect(namedRestParam.Type),
						namedRestParam.Name,
						types.Inspect(namedArgType),
					),
					namedArg.Span(),
				)
			}
		}

		typedPositionalArguments = append(typedPositionalArguments, namedRestArgs)
	} else {
		for i, defined := range definedNamedArgumentsSlice {
			if defined {
				continue
			}

			namedArgI := namedArguments[i]
			namedArg := namedArgI.(*ast.NamedCallArgumentNode)
			c.addError(
				fmt.Sprintf(
					"nonexistent parameter `%s` given in call to `%s`",
					namedArg.Name,
					method.Name,
				),
				namedArg.Span(),
			)
		}
	}

	return typedPositionalArguments
}

func (c *Checker) receiverlessMethodCall(node *ast.ReceiverlessMethodCallNode) *typed.ReceiverlessMethodCallNode {
	method := types.GetMethod(c.selfType, node.MethodName, c.GlobalEnv)
	if method == nil {
		c.addError(
			fmt.Sprintf("method `%s` is not defined in type `%s`", node.MethodName, types.Inspect(c.selfType)),
			node.Span(),
		)
		return typed.NewReceiverlessMethodCallNode(
			node.Span(),
			node.MethodName,
			c.checkExpressions(node.PositionalArguments),
			types.Void{},
		)
	}

	typedPositionalArguments := c.checkMethodArguments(method, node.PositionalArguments, node.NamedArguments, node.Span())

	return typed.NewReceiverlessMethodCallNode(
		node.Span(),
		node.MethodName,
		typedPositionalArguments,
		method.ReturnType,
	)
}

func (c *Checker) methodCall(node *ast.MethodCallNode) *typed.MethodCallNode {
	receiver := c.checkExpression(node.Receiver)
	receiverType := c.typeOf(receiver)
	method := types.GetMethod(receiverType, node.MethodName, c.GlobalEnv)
	if method == nil {
		c.addError(
			fmt.Sprintf("method `%s` is not defined in type `%s`", node.MethodName, types.Inspect(receiverType)),
			node.Span(),
		)
		return typed.NewMethodCallNode(
			node.Span(),
			receiver,
			node.NilSafe,
			node.MethodName,
			c.checkExpressions(node.PositionalArguments),
			types.Void{},
		)
	}

	typedPositionalArguments := c.checkMethodArguments(method, node.PositionalArguments, node.NamedArguments, node.Span())
	returnType := method.ReturnType
	if node.NilSafe {
		returnType = c.toNilable(returnType)
	}

	return typed.NewMethodCallNode(
		node.Span(),
		receiver,
		node.NilSafe,
		node.MethodName,
		typedPositionalArguments,
		returnType,
	)
}

func (c *Checker) attributeAccess(node *ast.AttributeAccessNode) typed.ExpressionNode {
	receiver := c.checkExpression(node.Receiver)
	receiverType := c.typeOf(receiver)
	method := types.GetMethod(receiverType, node.AttributeName, c.GlobalEnv)
	if method == nil {
		c.addError(
			fmt.Sprintf("method `%s` is not defined in type `%s`", node.AttributeName, types.Inspect(receiverType)),
			node.Span(),
		)
		return typed.NewAttributeAccessNode(
			node.Span(),
			receiver,
			node.AttributeName,
			types.Void{},
		)
	}

	typedPositionalArguments := c.checkMethodArguments(method, nil, nil, node.Span())

	return typed.NewMethodCallNode(
		node.Span(),
		receiver,
		false,
		node.AttributeName,
		typedPositionalArguments,
		method.ReturnType,
	)
}

func (c *Checker) addWrongArgumentCountError(got int, method *types.Method, span *position.Span) {
	c.addError(
		fmt.Sprintf("expected %s arguments, got %d", method.ExpectedParamCountString(), got),
		span,
	)
}

func (c *Checker) methodDefinition(node *ast.MethodDefinitionNode) *typed.MethodDefinitionNode {
	constScope := c.currentConstScope()
	env := newLocalEnvironment(nil)
	c.pushLocalEnv(env)
	defer c.popLocalEnv()

	oldMethod := constScope.container.MethodString(node.Name)

	var paramNodes []typed.ParameterNode
	var params []*types.Parameter
	for _, param := range node.Parameters {
		p, ok := param.(*ast.MethodParameterNode)
		if !ok {
			c.addError(
				fmt.Sprintf("invalid param type %T", node),
				param.Span(),
			)
		}
		var declaredType types.Type
		var declaredTypeNode typed.TypeNode
		if p.Type == nil {
			c.addError(
				fmt.Sprintf("cannot declare parameter `%s` without a type", p.Name),
				param.Span(),
			)
		} else {
			declaredTypeNode = c.checkTypeNode(p.Type)
			declaredType = c.typeOf(declaredTypeNode)
		}
		var initNode typed.ExpressionNode
		if p.Initialiser != nil {
			initNode = c.checkExpression(p.Initialiser)
			initType := c.typeOf(initNode)
			if !c.isSubtype(initType, declaredType) {
				c.addError(
					fmt.Sprintf("type `%s` cannot be assigned to type `%s`", types.Inspect(initType), types.Inspect(declaredType)),
					initNode.Span(),
				)
			}
		}
		c.addLocal(p.Name, local{typ: declaredType, initialised: true})
		var kind types.ParameterKind
		switch p.Kind {
		case ast.NormalParameterKind:
			kind = types.NormalParameterKind
		case ast.PositionalRestParameterKind:
			kind = types.PositionalRestParameterKind
		case ast.NamedRestParameterKind:
			kind = types.NamedRestParameterKind
		}
		if p.Initialiser != nil {
			kind = types.DefaultValueParameterKind
		}
		name := value.ToSymbol(p.Name)
		params = append(params, types.NewParameter(
			name,
			declaredType,
			kind,
			false,
		))
		paramNodes = append(paramNodes, typed.NewMethodParameterNode(
			p.Span(),
			p.Name,
			p.SetInstanceVariable,
			declaredTypeNode,
			initNode,
			typed.ParameterKind(p.Kind),
		))
	}

	var returnType types.Type
	var returnTypeNode typed.TypeNode
	if node.ReturnType != nil {
		returnTypeNode = c.checkTypeNode(node.ReturnType)
		returnType = c.typeOf(returnTypeNode)
	} else {
		returnType = types.Void{}
	}

	var throwType types.Type
	var throwTypeNode typed.TypeNode
	if node.ThrowType != nil {
		throwTypeNode = c.checkTypeNode(node.ThrowType)
		throwType = c.typeOf(throwTypeNode)
	}
	newMethod := constScope.container.NewMethod(
		node.Name,
		params,
		returnType,
		throwType,
	)

	setMethod := true
	if oldMethod != nil {
		if returnTypeNode != nil && oldMethod.ReturnType != nil && newMethod.ReturnType != oldMethod.ReturnType {
			c.addError(
				fmt.Sprintf("cannot redeclare method `%s` with a different return type, is `%s`, should be `%s`", node.Name, types.Inspect(newMethod.ReturnType), types.Inspect(oldMethod.ReturnType)),
				returnTypeNode.Span(),
			)
			setMethod = false
		}
		if newMethod.ThrowType != oldMethod.ThrowType {
			var span *position.Span
			if throwTypeNode != nil {
				span = throwTypeNode.Span()
			} else {
				span = node.Span()
			}
			c.addError(
				fmt.Sprintf("cannot redeclare method `%s` with a different throw type, is `%s`, should be `%s`", node.Name, types.Inspect(newMethod.ThrowType), types.Inspect(oldMethod.ThrowType)),
				span,
			)
			setMethod = false
		}

		if len(oldMethod.Params) > len(newMethod.Params) {
			c.addError(
				fmt.Sprintf("cannot redeclare method `%s` with less parameters", node.Name),
				position.JoinSpanOfCollection(node.Parameters),
			)
			setMethod = false
		} else {
			for i := range len(oldMethod.Params) {
				oldParam := oldMethod.Params[i]
				newParam := newMethod.Params[i]
				if oldParam.Name != newParam.Name {
					c.addError(
						fmt.Sprintf("cannot redeclare method `%s` with invalid parameter name, is `%s`, should be `%s`", node.Name, newParam.Name, oldParam.Name),
						paramNodes[i].Span(),
					)
					setMethod = false
					continue
				}
				if oldParam.Kind != newParam.Kind {
					c.addError(
						fmt.Sprintf("cannot redeclare method `%s` with invalid parameter kind, is `%s`, should be `%s`", node.Name, newParam.NameWithKind(), oldParam.NameWithKind()),
						paramNodes[i].Span(),
					)
					setMethod = false
					continue
				}
				if oldParam.Type != newParam.Type {
					c.addError(
						fmt.Sprintf("cannot redeclare method `%s` with invalid parameter type, is `%s`, should be `%s`", node.Name, types.Inspect(newParam.Type), types.Inspect(oldParam.Type)),
						paramNodes[i].Span(),
					)
					setMethod = false
					continue
				}
			}

			for i := len(oldMethod.Params); i < len(newMethod.Params); i++ {
				param := newMethod.Params[i]
				if !param.IsOptional() {
					c.addError(
						fmt.Sprintf("cannot redeclare method `%s` with additional required parameter `%s`", node.Name, param.Name),
						paramNodes[i].Span(),
					)
					setMethod = false
				}
			}
		}

	}

	if setMethod {
		constScope.container.SetMethod(node.Name, newMethod)
	}

	c.returnType = returnType
	c.throwType = throwType
	body := c.checkStatements(node.Body)
	c.returnType = nil
	c.throwType = nil

	return typed.NewMethodDefinitionNode(
		node.Span(),
		node.Name,
		paramNodes,
		returnTypeNode,
		throwTypeNode,
		body,
	)
}

func (c *Checker) assignmentExpression(node *ast.AssignmentExpressionNode) *typed.AssignmentExpressionNode {
	switch n := node.Left.(type) {
	case *ast.PublicIdentifierNode:
		return c.localVariableAssignment(n.Value, node.Op, node.Right, node.Span())
	case *ast.PrivateIdentifierNode:
		return c.localVariableAssignment(n.Value, node.Op, node.Right, node.Span())
	// case *ast.SubscriptExpressionNode:
	// case *ast.ConstantLookupNode:
	// case *ast.PublicConstantNode:
	// case *ast.PrivateConstantNode:
	// case *ast.InstanceVariableNode:
	// case *ast.AttributeAccessNode:
	default:
		c.Errors.Add(
			fmt.Sprintf("cannot assign to: %T", node.Left),
			c.newLocation(node.Span()),
		)
	}

	return nil
}

func (c *Checker) localVariableAssignment(name string, operator *token.Token, right ast.ExpressionNode, span *position.Span) *typed.AssignmentExpressionNode {
	switch operator.Type {
	case token.EQUAL_OP:
	// case token.COLON_EQUAL:
	// case token.OR_OR_EQUAL:
	// case token.AND_AND_EQUAL:
	// case token.QUESTION_QUESTION_EQUAL:
	// case token.PLUS_EQUAL:
	// case token.MINUS_EQUAL:
	// case token.STAR_EQUAL:
	// case token.SLASH_EQUAL:
	// case token.STAR_STAR_EQUAL:
	// case token.PERCENT_EQUAL:
	// case token.AND_EQUAL:
	// case token.OR_EQUAL:
	// case token.XOR_EQUAL:
	// case token.LBITSHIFT_EQUAL:
	// case token.LTRIPLE_BITSHIFT_EQUAL:
	// case token.RBITSHIFT_EQUAL:
	// case token.RTRIPLE_BITSHIFT_EQUAL:
	default:
		c.Errors.Add(
			fmt.Sprintf("assignment using this operator has not been implemented: %s", operator.Type.String()),
			c.newLocation(span),
		)
	}
	return nil
}

func (c *Checker) interpolatedStringLiteral(node *ast.InterpolatedStringLiteralNode) *typed.InterpolatedStringLiteralNode {
	var newContent []typed.StringLiteralContentNode
	for _, contentSection := range node.Content {
		newContent = append(newContent, c.checkStringContent(contentSection))
	}
	return typed.NewInterpolatedStringLiteralNode(
		node.Span(),
		newContent,
	)
}

func (c *Checker) checkStringContent(node ast.StringLiteralContentNode) typed.StringLiteralContentNode {
	switch n := node.(type) {
	case *ast.StringInspectInterpolationNode:
		return typed.NewStringInspectInterpolationNode(
			n.Span(),
			c.checkExpression(n.Expression),
		)
	case *ast.StringInterpolationNode:
		return typed.NewStringInterpolationNode(
			n.Span(),
			c.checkExpression(n.Expression),
		)
	case *ast.StringLiteralContentSectionNode:
		return typed.NewStringLiteralContentSectionNode(
			n.Span(),
			n.Value,
		)
	default:
		c.addError(
			fmt.Sprintf("invalid string content %T", node),
			node.Span(),
		)
		return typed.NewInvalidNode(node.Span(), nil)
	}
}

func (c *Checker) addError(message string, span *position.Span) {
	c.Errors.Add(
		message,
		c.newLocation(span),
	)
}

// Get the type of the constant with the given name
func (c *Checker) resolveConstantForSetter(name string) (types.Type, string) {
	constScope := c.currentConstScope()
	constant := constScope.container.ConstantString(name)
	fullName := types.MakeFullConstantName(constScope.container.Name(), name)
	if constant != nil {
		return constant, fullName
	}
	return nil, fullName
}

// Get the type of the public constant with the given name
func (c *Checker) resolvePublicConstant(name string, span *position.Span) (types.Type, string) {
	for i := range len(c.constantScopes) {
		constScope := c.constantScopes[len(c.constantScopes)-i-1]
		constant := constScope.container.ConstantString(name)
		if constant != nil {
			return constant, types.MakeFullConstantName(constScope.container.Name(), name)
		}
	}

	c.addError(
		fmt.Sprintf("undefined constant `%s`", name),
		span,
	)
	return nil, name
}

// Get the type of the private constant with the given name
func (c *Checker) resolvePrivateConstant(name string, span *position.Span) (types.Type, string) {
	for i := range len(c.constantScopes) {
		constScope := c.constantScopes[len(c.constantScopes)-i-1]
		if !constScope.local {
			continue
		}
		constant := constScope.container.ConstantString(name)
		if constant != nil {
			return constant, types.MakeFullConstantName(constScope.container.Name(), name)
		}
	}

	c.addError(
		fmt.Sprintf("undefined constant `%s`", name),
		span,
	)
	return nil, name
}

// Get the type with the given name
func (c *Checker) resolveType(name string, span *position.Span) (types.Type, string) {
	for i := range len(c.constantScopes) {
		constScope := c.constantScopes[len(c.constantScopes)-i-1]
		constant := constScope.container.SubtypeString(name)
		if constant != nil {
			return constant, types.MakeFullConstantName(constScope.container.Name(), name)
		}
	}

	c.addError(
		fmt.Sprintf("undefined type `%s`", name),
		span,
	)
	return nil, name
}

// Add the local with the given name to the current local environment
func (c *Checker) addLocal(name string, l local) {
	env := c.currentLocalEnv()
	env.locals[value.ToSymbol(name)] = l
}

// Get the local with the specified name from the current local environment
func (c *Checker) getLocal(name string) (local, bool) {
	env := c.currentLocalEnv()
	return env.getLocal(name)
}

// Resolve the local with the given name from the current local environment or any parent environment
func (c *Checker) resolveLocal(name string, span *position.Span) (local, bool) {
	env := c.currentLocalEnv()
	local, ok := env.resolveLocal(name)
	if !ok {
		c.addError(
			fmt.Sprintf("undefined local `%s`", name),
			span,
		)
	}
	return local, ok
}

func (c *Checker) resolveConstantLookupType(node *ast.ConstantLookupNode) (types.Type, string) {
	var leftContainerType types.Type
	var leftContainerName string

	switch l := node.Left.(type) {
	case *ast.PublicConstantNode:
		leftContainerType, leftContainerName = c.resolveType(l.Value, l.Span())
	case *ast.PrivateConstantNode:
		leftContainerType, leftContainerName = c.resolveType(l.Value, l.Span())
	case nil:
		leftContainerType = c.GlobalEnv.Root
	case *ast.ConstantLookupNode:
		leftContainerType, leftContainerName = c.resolveConstantLookupType(l)
	default:
		c.addError(
			fmt.Sprintf("invalid type node %T", node),
			node.Span(),
		)
		return nil, ""
	}

	var rightName string
	switch r := node.Right.(type) {
	case *ast.PublicConstantNode:
		rightName = r.Value
	case *ast.PrivateConstantNode:
		rightName = r.Value
		c.addError(
			fmt.Sprintf("cannot use private type `%s`", rightName),
			node.Span(),
		)
	default:
		c.addError(
			fmt.Sprintf("invalid type node %T", node),
			node.Span(),
		)
		return nil, ""
	}

	typeName := types.MakeFullConstantName(leftContainerName, rightName)
	if leftContainerType == nil {
		return nil, typeName
	}
	leftContainer, ok := leftContainerType.(types.ConstantContainer)
	if !ok {
		c.addError(
			fmt.Sprintf("cannot read subtypes from `%s`, it is not a type container", leftContainerName),
			node.Span(),
		)
		return nil, typeName
	}

	constant := leftContainer.SubtypeString(rightName)
	if constant == nil {
		c.addError(
			fmt.Sprintf("undefined type `%s`", typeName),
			node.Right.Span(),
		)
		return nil, typeName
	}

	return constant, typeName
}

func (c *Checker) checkTypeNode(node ast.TypeNode) typed.TypeNode {
	switch n := node.(type) {
	case *ast.PublicConstantNode:
		typ, _ := c.resolveType(n.Value, n.Span())
		if typ == nil {
			typ = types.Void{}
		}
		return typed.NewPublicConstantNode(
			n.Span(),
			n.Value,
			typ,
		)
	case *ast.PrivateConstantNode:
		typ, _ := c.resolveType(n.Value, n.Span())
		if typ == nil {
			typ = types.Void{}
		}
		return typed.NewPrivateConstantNode(
			n.Span(),
			n.Value,
			typ,
		)
	case *ast.ConstantLookupNode:
		return c.constantLookupType(n)
	case *ast.RawStringLiteralNode:
		return typed.NewRawStringLiteralNode(
			n.Span(),
			n.Value,
			types.NewStringLiteral(n.Value),
		)
	case *ast.DoubleQuotedStringLiteralNode:
		return typed.NewDoubleQuotedStringLiteralNode(
			n.Span(),
			n.Value,
			types.NewStringLiteral(n.Value),
		)
	case *ast.RawCharLiteralNode:
		return typed.NewRawCharLiteralNode(
			n.Span(),
			n.Value,
			types.NewCharLiteral(n.Value),
		)
	case *ast.CharLiteralNode:
		return typed.NewCharLiteralNode(
			n.Span(),
			n.Value,
			types.NewCharLiteral(n.Value),
		)
	case *ast.NilLiteralNode:
		return typed.NewNilLiteralNode(n.Span())
	case *ast.SimpleSymbolLiteralNode:
		return typed.NewSimpleSymbolLiteralNode(
			n.Span(),
			n.Content,
			types.NewSymbolLiteral(n.Content),
		)
	case *ast.BinaryTypeExpressionNode:
		switch n.Op.Type {
		case token.OR:
			return c.constructUnionType(n)
		case token.AND:
			return c.constructIntersectionType(n)
		default:
			panic("invalid binary type expression operator")
		}
	case *ast.IntLiteralNode:
		return typed.NewIntLiteralNode(
			n.Span(),
			n.Value,
			types.NewIntLiteral(n.Value),
		)
	case *ast.Int64LiteralNode:
		return typed.NewInt64LiteralNode(
			n.Span(),
			n.Value,
			types.NewInt64Literal(n.Value),
		)
	case *ast.Int32LiteralNode:
		return typed.NewInt32LiteralNode(
			n.Span(),
			n.Value,
			types.NewInt32Literal(n.Value),
		)
	case *ast.Int16LiteralNode:
		return typed.NewInt16LiteralNode(
			n.Span(),
			n.Value,
			types.NewInt16Literal(n.Value),
		)
	case *ast.Int8LiteralNode:
		return typed.NewInt8LiteralNode(
			n.Span(),
			n.Value,
			types.NewInt8Literal(n.Value),
		)
	case *ast.UInt64LiteralNode:
		return typed.NewUInt64LiteralNode(
			n.Span(),
			n.Value,
			types.NewUInt64Literal(n.Value),
		)
	case *ast.UInt32LiteralNode:
		return typed.NewUInt32LiteralNode(
			n.Span(),
			n.Value,
			types.NewUInt32Literal(n.Value),
		)
	case *ast.UInt16LiteralNode:
		return typed.NewUInt16LiteralNode(
			n.Span(),
			n.Value,
			types.NewUInt16Literal(n.Value),
		)
	case *ast.UInt8LiteralNode:
		return typed.NewUInt8LiteralNode(
			n.Span(),
			n.Value,
			types.NewUInt8Literal(n.Value),
		)
	case *ast.FloatLiteralNode:
		return typed.NewFloatLiteralNode(
			n.Span(),
			n.Value,
			types.NewFloatLiteral(n.Value),
		)
	case *ast.Float64LiteralNode:
		return typed.NewFloat64LiteralNode(
			n.Span(),
			n.Value,
			types.NewFloat64Literal(n.Value),
		)
	case *ast.Float32LiteralNode:
		return typed.NewFloat32LiteralNode(
			n.Span(),
			n.Value,
			types.NewFloat32Literal(n.Value),
		)
	case *ast.TrueLiteralNode:
		return typed.NewTrueLiteralNode(n.Span())
	case *ast.FalseLiteralNode:
		return typed.NewFalseLiteralNode(n.Span())
	case *ast.VoidTypeNode:
		return typed.NewVoidTypeNode(n.Span())
	case *ast.NeverTypeNode:
		return typed.NewNeverTypeNode(n.Span())
	case *ast.AnyTypeNode:
		return typed.NewAnyTypeNode(n.Span())
	case *ast.NilableTypeNode:
		typeNode := c.checkTypeNode(n.Type)
		typ := c.toNilable(c.typeOf(typeNode))
		return typed.NewNilableTypeNode(
			n.Span(),
			c.checkTypeNode(n.Type),
			typ,
		)
	default:
		c.addError(
			fmt.Sprintf("invalid type node %T", node),
			node.Span(),
		)
		return typed.NewInvalidNode(node.Span(), nil)
	}
}

func (c *Checker) constructUnionType(node *ast.BinaryTypeExpressionNode) *typed.UnionTypeNode {
	union := types.NewUnion()
	elements := new([]typed.TypeNode)
	c._constructUnionType(node, elements, union)

	return typed.NewUnionTypeNode(
		node.Span(),
		*elements,
		union,
	)
}

func (c *Checker) _constructUnionType(node *ast.BinaryTypeExpressionNode, elements *[]typed.TypeNode, union *types.Union) {
	leftBinaryType, leftIsBinaryType := node.Left.(*ast.BinaryTypeExpressionNode)
	if leftIsBinaryType && leftBinaryType.Op.Type == token.OR {
		c._constructUnionType(leftBinaryType, elements, union)
	} else {
		leftTypeNode := c.checkTypeNode(node.Left)
		*elements = append(*elements, leftTypeNode)

		leftType := c.typeOf(leftTypeNode)
		union.Elements = append(union.Elements, leftType)
	}

	rightBinaryType, rightIsBinaryType := node.Right.(*ast.BinaryTypeExpressionNode)
	if rightIsBinaryType && rightBinaryType.Op.Type == token.OR {
		c._constructUnionType(rightBinaryType, elements, union)
	} else {
		rightTypeNode := c.checkTypeNode(node.Right)
		*elements = append(*elements, rightTypeNode)

		rightType := c.typeOf(rightTypeNode)
		union.Elements = append(union.Elements, rightType)
	}
}

func (c *Checker) constructIntersectionType(node *ast.BinaryTypeExpressionNode) *typed.IntersectionTypeNode {
	intersection := types.NewIntersection()
	elements := new([]typed.TypeNode)
	c._constructIntersectionType(node, elements, intersection)

	return typed.NewIntersectionTypeNode(
		node.Span(),
		*elements,
		intersection,
	)
}

func (c *Checker) _constructIntersectionType(node *ast.BinaryTypeExpressionNode, elements *[]typed.TypeNode, intersection *types.Intersection) {
	leftBinaryType, leftIsBinaryType := node.Left.(*ast.BinaryTypeExpressionNode)
	if leftIsBinaryType && leftBinaryType.Op.Type == token.AND {
		c._constructIntersectionType(leftBinaryType, elements, intersection)
	} else {
		leftTypeNode := c.checkTypeNode(node.Left)
		*elements = append(*elements, leftTypeNode)

		leftType := c.typeOf(leftTypeNode)
		intersection.Elements = append(intersection.Elements, leftType)
	}

	rightBinaryType, rightIsBinaryType := node.Right.(*ast.BinaryTypeExpressionNode)
	if rightIsBinaryType && rightBinaryType.Op.Type == token.AND {
		c._constructIntersectionType(rightBinaryType, elements, intersection)
	} else {
		rightTypeNode := c.checkTypeNode(node.Right)
		*elements = append(*elements, rightTypeNode)

		rightType := c.typeOf(rightTypeNode)
		intersection.Elements = append(intersection.Elements, rightType)
	}
}

func (c *Checker) constantLookupType(node *ast.ConstantLookupNode) *typed.PublicConstantNode {
	typ, name := c.resolveConstantLookupType(node)
	if typ == nil {
		typ = types.Void{}
	}

	return typed.NewPublicConstantNode(
		node.Span(),
		name,
		typ,
	)
}

func (c *Checker) resolveConstantLookup(node *ast.ConstantLookupNode) (types.Type, string) {
	var leftContainerType types.Type
	var leftContainerName string

	switch l := node.Left.(type) {
	case *ast.PublicConstantNode:
		leftContainerType, leftContainerName = c.resolvePublicConstant(l.Value, l.Span())
	case *ast.PrivateConstantNode:
		leftContainerType, leftContainerName = c.resolvePrivateConstant(l.Value, l.Span())
	case nil:
		leftContainerType = c.GlobalEnv.Root
	case *ast.ConstantLookupNode:
		leftContainerType, leftContainerName = c.resolveConstantLookup(l)
	default:
		c.addError(
			fmt.Sprintf("invalid constant node %T", node),
			node.Span(),
		)
		return nil, ""
	}

	var rightName string
	switch r := node.Right.(type) {
	case *ast.PublicConstantNode:
		rightName = r.Value
	case *ast.PrivateConstantNode:
		rightName = r.Value
		c.addError(
			fmt.Sprintf("cannot read private constant `%s`", rightName),
			node.Span(),
		)
	default:
		c.addError(
			fmt.Sprintf("invalid constant node %T", node),
			node.Span(),
		)
		return nil, ""
	}

	constantName := types.MakeFullConstantName(leftContainerName, rightName)
	if leftContainerType == nil {
		return nil, constantName
	}

	var leftContainer types.ConstantContainer
	switch l := leftContainerType.(type) {
	case *types.Module:
		leftContainer = l
	case *types.SingletonClass:
		leftContainer = l.AttachedObject
	default:
		c.addError(
			fmt.Sprintf("cannot read constants from `%s`, it is not a constant container", leftContainerName),
			node.Span(),
		)
		return nil, constantName
	}

	constant := leftContainer.ConstantString(rightName)
	if constant == nil {
		c.addError(
			fmt.Sprintf("undefined constant `%s`", constantName),
			node.Right.Span(),
		)
		return nil, constantName
	}

	return constant, constantName
}

func (c *Checker) _resolveConstantLookupForSetter(node *ast.ConstantLookupNode, firstCall bool) (types.Type, string) {
	var leftContainerType types.Type
	var leftContainerName string

	switch l := node.Left.(type) {
	case *ast.PublicConstantNode:
		leftContainerType, leftContainerName = c.resolvePublicConstant(l.Value, l.Span())
	case *ast.PrivateConstantNode:
		leftContainerType, leftContainerName = c.resolvePrivateConstant(l.Value, l.Span())
	case nil:
		leftContainerType = c.GlobalEnv.Root
	case *ast.ConstantLookupNode:
		leftContainerType, leftContainerName = c._resolveConstantLookupForSetter(l, false)
	default:
		c.addError(
			fmt.Sprintf("invalid constant node %T", node),
			node.Span(),
		)
		return nil, ""
	}

	var rightName string
	switch r := node.Right.(type) {
	case *ast.PublicConstantNode:
		rightName = r.Value
	case *ast.PrivateConstantNode:
		rightName = r.Value
		c.addError(
			fmt.Sprintf("cannot read private constant `%s`", rightName),
			node.Span(),
		)
	default:
		c.addError(
			fmt.Sprintf("invalid constant node %T", node),
			node.Span(),
		)
		return nil, ""
	}

	constantName := types.MakeFullConstantName(leftContainerName, rightName)
	if leftContainerType == nil {
		return nil, constantName
	}
	var leftContainer types.ConstantContainer
	switch l := leftContainerType.(type) {
	case *types.Module:
		leftContainer = l
	case *types.SingletonClass:
		leftContainer = l.AttachedObject
	default:
		c.addError(
			fmt.Sprintf("cannot read constants from `%s`, it is not a constant container", leftContainerName),
			node.Span(),
		)
		return nil, constantName
	}

	constant := leftContainer.ConstantString(rightName)
	if constant == nil {
		if !firstCall {
			c.addError(
				fmt.Sprintf("undefined constant `%s`", constantName),
				node.Right.Span(),
			)
		}
		return nil, constantName
	}

	return constant, constantName
}

func (c *Checker) resolveConstantLookupForSetter(node *ast.ConstantLookupNode) (types.Type, string) {
	return c._resolveConstantLookupForSetter(node, true)
}

func (c *Checker) constantLookup(node *ast.ConstantLookupNode) *typed.PublicConstantNode {
	typ, name := c.resolveConstantLookup(node)
	if typ == nil {
		typ = types.Void{}
	}

	return typed.NewPublicConstantNode(
		node.Span(),
		name,
		typ,
	)
}

func (c *Checker) publicConstant(node *ast.PublicConstantNode) *typed.PublicConstantNode {
	typ, name := c.resolvePublicConstant(node.Value, node.Span())
	if typ == nil {
		typ = types.Void{}
	}

	return typed.NewPublicConstantNode(
		node.Span(),
		name,
		typ,
	)
}

func (c *Checker) privateConstant(node *ast.PrivateConstantNode) *typed.PrivateConstantNode {
	typ, name := c.resolvePrivateConstant(node.Value, node.Span())
	if typ == nil {
		typ = types.Void{}
	}

	return typed.NewPrivateConstantNode(
		node.Span(),
		name,
		typ,
	)
}

func (c *Checker) publicIdentifier(node *ast.PublicIdentifierNode) *typed.PublicIdentifierNode {
	local, ok := c.resolveLocal(node.Value, node.Span())
	if ok && !local.initialised {
		c.addError(
			fmt.Sprintf("cannot access uninitialised local `%s`", node.Value),
			node.Span(),
		)
	}
	return typed.NewPublicIdentifierNode(
		node.Span(),
		node.Value,
		local.typ,
	)
}

func (c *Checker) privateIdentifier(node *ast.PrivateIdentifierNode) *typed.PrivateIdentifierNode {
	local, ok := c.resolveLocal(node.Value, node.Span())
	if ok && !local.initialised {
		c.addError(
			fmt.Sprintf("cannot access uninitialised local `%s`", node.Value),
			node.Span(),
		)
	}
	return typed.NewPrivateIdentifierNode(
		node.Span(),
		node.Value,
		local.typ,
	)
}

func (c *Checker) variableDeclaration(node *ast.VariableDeclarationNode) *typed.VariableDeclarationNode {
	if _, ok := c.getLocal(node.Name.Value); ok {
		c.addError(
			fmt.Sprintf("cannot redeclare local `%s`", node.Name.Value),
			node.Span(),
		)
	}
	if node.Initialiser == nil {
		if node.Type == nil {
			c.addError(
				fmt.Sprintf("cannot declare a variable without a type `%s`", node.Name.Value),
				node.Span(),
			)
			c.addLocal(node.Name.Value, local{typ: types.Void{}})
			return typed.NewVariableDeclarationNode(
				node.Span(),
				node.Name,
				nil,
				nil,
				types.Void{},
			)
		}

		// without an initialiser but with a type
		declaredTypeNode := c.checkTypeNode(node.Type)
		declaredType := c.typeOf(declaredTypeNode)
		c.addLocal(node.Name.Value, local{typ: declaredType})
		return typed.NewVariableDeclarationNode(
			node.Span(),
			node.Name,
			declaredTypeNode,
			nil,
			types.Void{},
		)
	}

	// with an initialiser
	if node.Type == nil {
		// without a type, inference
		init := c.checkExpression(node.Initialiser)
		actualType := c.typeOf(init).ToNonLiteral(c.GlobalEnv)
		c.addLocal(node.Name.Value, local{typ: actualType, initialised: true})
		return typed.NewVariableDeclarationNode(
			node.Span(),
			node.Name,
			nil,
			init,
			actualType,
		)
	}

	// with a type and an initializer

	declaredTypeNode := c.checkTypeNode(node.Type)
	declaredType := c.typeOf(declaredTypeNode)
	init := c.checkExpression(node.Initialiser)
	actualType := c.typeOf(init)
	c.addLocal(node.Name.Value, local{typ: declaredType, initialised: true})
	if !c.isSubtype(actualType, declaredType) {
		c.addError(
			fmt.Sprintf("type `%s` cannot be assigned to type `%s`", types.Inspect(actualType), types.Inspect(declaredType)),
			init.Span(),
		)
	}

	return typed.NewVariableDeclarationNode(
		node.Span(),
		node.Name,
		declaredTypeNode,
		init,
		declaredType,
	)
}

func (c *Checker) valueDeclaration(node *ast.ValueDeclarationNode) *typed.ValueDeclarationNode {
	if _, ok := c.getLocal(node.Name.Value); ok {
		c.addError(
			fmt.Sprintf("cannot redeclare local `%s`", node.Name.Value),
			node.Span(),
		)
	}
	if node.Initialiser == nil {
		if node.Type == nil {
			c.addError(
				fmt.Sprintf("cannot declare a value without a type `%s`", node.Name.Value),
				node.Span(),
			)
			c.addLocal(node.Name.Value, local{typ: types.Void{}})
			return typed.NewValueDeclarationNode(
				node.Span(),
				node.Name,
				nil,
				nil,
				types.Void{},
			)
		}

		// without an initialiser but with a type
		declaredTypeNode := c.checkTypeNode(node.Type)
		declaredType := c.typeOf(declaredTypeNode)
		c.addLocal(node.Name.Value, local{typ: declaredType, singleAssignment: true})
		return typed.NewValueDeclarationNode(
			node.Span(),
			node.Name,
			declaredTypeNode,
			nil,
			types.Void{},
		)
	}

	// with an initialiser
	if node.Type == nil {
		// without a type, inference
		init := c.checkExpression(node.Initialiser)
		actualType := c.typeOf(init)
		c.addLocal(node.Name.Value, local{typ: actualType, initialised: true, singleAssignment: true})
		return typed.NewValueDeclarationNode(
			node.Span(),
			node.Name,
			nil,
			init,
			actualType,
		)
	}

	// with a type and an initializer

	declaredTypeNode := c.checkTypeNode(node.Type)
	declaredType := c.typeOf(declaredTypeNode)
	init := c.checkExpression(node.Initialiser)
	actualType := c.typeOf(init)
	c.addLocal(node.Name.Value, local{typ: declaredType, initialised: true, singleAssignment: true})
	if !c.isSubtype(actualType, declaredType) {
		c.addError(
			fmt.Sprintf("type `%s` cannot be assigned to type `%s`", types.Inspect(actualType), types.Inspect(declaredType)),
			init.Span(),
		)
	}

	return typed.NewValueDeclarationNode(
		node.Span(),
		node.Name,
		declaredTypeNode,
		init,
		declaredType,
	)
}

func (c *Checker) constantDeclaration(node *ast.ConstantDeclarationNode) *typed.ConstantDeclarationNode {
	if constant, constantName := c.resolveConstantForSetter(node.Name.Value); constant != nil {
		c.addError(
			fmt.Sprintf("cannot redeclare constant `%s`", constantName),
			node.Span(),
		)
	}

	scope := c.currentConstScope()
	if node.Type == nil {
		// without a type, inference
		init := c.checkExpression(node.Initialiser)
		actualType := c.typeOf(init)
		scope.container.DefineConstant(node.Name.Value, actualType)
		return typed.NewConstantDeclarationNode(
			node.Span(),
			node.Name,
			nil,
			init,
			actualType,
		)
	}

	// with a type

	declaredTypeNode := c.checkTypeNode(node.Type)
	declaredType := c.typeOf(declaredTypeNode)
	init := c.checkExpression(node.Initialiser)
	actualType := c.typeOf(init)
	scope.container.DefineConstant(node.Name.Value, actualType)
	if !c.isSubtype(actualType, declaredType) {
		c.addError(
			fmt.Sprintf("type `%s` cannot be assigned to type `%s`", types.Inspect(actualType), types.Inspect(declaredType)),
			init.Span(),
		)
	}

	return typed.NewConstantDeclarationNode(
		node.Span(),
		node.Name,
		declaredTypeNode,
		init,
		declaredType,
	)
}

func (c *Checker) module(node *ast.ModuleDeclarationNode) *typed.ModuleDeclarationNode {
	constScope := c.currentConstScope()
	var typedConstantNode typed.ExpressionNode
	var module *types.Module

	switch constant := node.Constant.(type) {
	case *ast.PublicConstantNode:
		constantType, constantName := c.resolveConstantForSetter(constant.Value)
		constantModule, constantIsModule := constantType.(*types.Module)
		if constantType != nil {
			if !constantIsModule {
				c.addError(
					fmt.Sprintf("cannot redeclare constant `%s`", constantName),
					node.Constant.Span(),
				)
			}
			module = constantModule
		} else {
			module = constScope.container.DefineModule(constant.Value, nil, nil)
		}
		typedConstantNode = typed.NewPublicConstantNode(
			constant.Span(),
			constantName,
			module,
		)
	case *ast.PrivateConstantNode:
		constantType, constantName := c.resolveConstantForSetter(constant.Value)
		constantModule, constantIsModule := constantType.(*types.Module)
		if constantType != nil {
			if !constantIsModule {
				c.addError(
					fmt.Sprintf("cannot redeclare constant `%s`", constantName),
					node.Constant.Span(),
				)
			}
			module = constantModule
		} else {
			module = constScope.container.DefineModule(constant.Value, nil, nil)
		}
		typedConstantNode = typed.NewPublicConstantNode(
			constant.Span(),
			constantName,
			module,
		)
	case *ast.ConstantLookupNode:
		constantType, constantName := c.resolveConstantLookupForSetter(constant)
		constantModule, constantIsModule := constantType.(*types.Module)
		if constantType != nil {
			if !constantIsModule {
				c.addError(
					fmt.Sprintf("cannot redeclare constant `%s`", constantName),
					node.Constant.Span(),
				)
			}
			module = constantModule
		} else {
			module = constScope.container.DefineModule(constantName, nil, nil)
		}
		typedConstantNode = typed.NewPublicConstantNode(
			constant.Span(),
			constantName,
			module,
		)
	case nil:
		module = types.NewModule("", nil, nil)
	default:
		c.addError(
			fmt.Sprintf("invalid module name node %T", node.Constant),
			node.Constant.Span(),
		)
	}

	c.pushConstScope(makeLocalConstantScope(module))
	prevSelfType := c.selfType
	c.selfType = module
	newBody := c.checkStatements(node.Body)
	c.selfType = prevSelfType
	c.popConstScope()

	return typed.NewModuleDeclarationNode(
		node.Span(),
		typedConstantNode,
		newBody,
		module,
	)
}
