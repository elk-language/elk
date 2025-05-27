package checker

import (
	"fmt"

	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

// Expand top level macros
func (c *Checker) expandTopLevelMacrosInFile(filename string, node *ast.ProgramNode) {
	node.State = ast.EXPANDING_TOP_LEVEL_MACROS
	for _, importPath := range node.ImportPaths {
		importedAst, ok := c.ASTCache.GetUnsafe(importPath)
		if !ok {
			continue
		}
		switch importedAst.State {
		case ast.EXPANDED_TOP_LEVEL_MACROS, ast.EXPANDING_TOP_LEVEL_MACROS:
			continue
		}
		c.expandTopLevelMacrosInFile(importPath, importedAst)
	}

	prevFilename := c.Filename
	c.Filename = filename
	c.expandTopLevelMacros(node.Body)
	c.Filename = prevFilename
	node.State = ast.EXPANDED_TOP_LEVEL_MACROS
}

func (c *Checker) expandTopLevelMacros(statements []ast.StatementNode) {
	for _, statement := range statements {
		switch stmt := statement.(type) {
		case *ast.ExpressionStatementNode:
			stmt.Expression = c.expandTopLevelMacrosInExpression(stmt.Expression)
		}
	}
}

func (c *Checker) hoistSingletonDeclarationAndExpandMacros(node *ast.SingletonBlockExpressionNode) {
	c.hoistSingletonDeclarationWithFunc(node, func(body []ast.StatementNode) {
		c.expandTopLevelMacros(body)
	})
}

func (c *Checker) hoistInterfaceDeclarationAndExpandMacros(node *ast.InterfaceDeclarationNode) {
	c.hoistInterfaceDeclarationWithFunc(node, func(body []ast.StatementNode) {
		c.expandTopLevelMacros(body)
	})
}

func (c *Checker) hoistMixinDeclarationAndExpandMacros(node *ast.MixinDeclarationNode) {
	c.hoistMixinDeclarationWithFunc(node, func(body []ast.StatementNode) {
		c.expandTopLevelMacros(body)
	})
}

func (c *Checker) hoistModuleDeclarationAndExpandMacros(node *ast.ModuleDeclarationNode) {
	c.hoistModuleDeclarationWithFunc(node, func(body []ast.StatementNode) {
		c.expandTopLevelMacros(body)
	})
}

func (c *Checker) hoistClassDeclarationAndExpandMacros(node *ast.ClassDeclarationNode) {
	c.hoistClassDeclarationWithFunc(node, func(body []ast.StatementNode) {
		c.expandTopLevelMacros(body)
	})
}

func (c *Checker) expandTopLevelMacrosInExpression(expr ast.ExpressionNode) ast.ExpressionNode {
	switch expr := expr.(type) {
	case *ast.MacroBoundaryNode:
		c.expandTopLevelMacros(expr.Body)
	case *ast.ModuleDeclarationNode:
		c.hoistModuleDeclarationAndExpandMacros(expr)
	case *ast.ClassDeclarationNode:
		c.hoistClassDeclarationAndExpandMacros(expr)
	case *ast.MixinDeclarationNode:
		c.hoistMixinDeclarationAndExpandMacros(expr)
	case *ast.InterfaceDeclarationNode:
		c.hoistInterfaceDeclarationAndExpandMacros(expr)
	case *ast.SingletonBlockExpressionNode:
		c.hoistSingletonDeclaration(expr)
	case *ast.ExtendWhereBlockExpressionNode:
		c.expandTopLevelMacros(expr.Body)
	case *ast.StructDeclarationNode:
		c.hoistStructDeclaration(expr)
	case *ast.MacroCallNode:
		posArgs := []ast.ExpressionNode{expr.Receiver}

		result := c.expandMacroByName(
			expr.MacroName,
			append(posArgs, expr.PositionalArguments...),
			expr.NamedArguments,
			expr.Location(),
		)
		if result == nil {
			return expr
		}
		return c.expandTopLevelMacrosInExpression(result)
	case *ast.ReceiverlessMacroCallNode:
		result := c.expandMacroByName(
			expr.MacroName,
			expr.PositionalArguments,
			expr.NamedArguments,
			expr.Location(),
		)
		if result == nil {
			return expr
		}
		return c.expandTopLevelMacrosInExpression(result)
	}

	return expr
}

func (c *Checker) resolveMacro(name value.Symbol) *types.Method {
	for methodScope := range ds.ReverseSlice(c.methodScopes) {
		var namespace types.Namespace
		switch n := methodScope.container.(type) {
		case *types.Class:
			namespace = n.Singleton()
		case *types.Mixin:
			namespace = n.Singleton()
		case *types.Module:
			namespace = n
		case *types.SingletonClass:
			namespace = n
		default:
			continue
		}

		for parent := range types.Parents(namespace) {
			macro := parent.Method(name)
			if macro != nil {
				return macro
			}
		}
	}

	return nil
}

func (c *Checker) getMacro(name value.Symbol, loc *position.Location) *types.Method {
	macro := c.resolveMacro(name)
	if macro == nil {
		c.addFailure(
			fmt.Sprintf(
				"undefined macro `%s`",
				name.String(),
			),
			loc,
		)
	}

	return macro
}

func (c *Checker) expandMacroByName(name string, posArgs []ast.ExpressionNode, namedArgs []ast.NamedArgumentNode, loc *position.Location) ast.ExpressionNode {
	macroName := value.ToSymbol(name + "!")
	macro := c.getMacro(macroName, loc)
	if macro == nil {
		return nil
	}

	return c.expandMacro(macro, posArgs, namedArgs, loc)
}

func (c *Checker) expandMacro(macro *types.Method, posArgs []ast.ExpressionNode, namedArgs []ast.NamedArgumentNode, loc *position.Location) ast.ExpressionNode {
	checkedArgs := c.checkMacroArguments(macro, posArgs, namedArgs, loc)
	if c.Errors.IsFailure() {
		return nil
	}

	runtimeArgs := make([]value.Value, 0, len(checkedArgs)+1)
	runtimeArgs = append(
		runtimeArgs,
		value.Ref(value.NodeMixin),
	)

	for _, arg := range checkedArgs {
		runtimeArgs = append(runtimeArgs, value.Ref(arg))
	}

	promise := vm.NewPromiseForBytecode(c.threadPool, macro.Bytecode, runtimeArgs...)
	result, err := promise.AwaitSync()
	if !err.IsUndefined() {
		c.addFailure(
			fmt.Sprintf(
				"error while executing macro `%s`: %s",
				types.InspectWithColor(macro),
				lexer.Colorize(err.Inspect()),
			),
			loc,
		)
		return nil
	}

	resultNode := result.AsReference().(ast.ExpressionNode)
	// wrap in a macro boundary to make it hygienic
	resultBoundary, ok := resultNode.(*ast.MacroBoundaryNode)
	if ok && resultBoundary.Name == "" {
		resultNode = ast.NewMacroBoundaryNode(
			resultNode.Location(),
			resultBoundary.Body,
			types.Inspect(macro),
		)
	} else {
		resultNode = ast.NewMacroBoundaryNode(
			resultNode.Location(),
			ast.ExpressionToStatements(resultNode),
			types.Inspect(macro),
		)
	}
	// update location
	resultNode = ast.Splice(resultNode, loc, nil).(ast.ExpressionNode)

	return resultNode
}

func (c *Checker) checkMacroArguments(
	method *types.Method,
	positionalArguments []ast.ExpressionNode,
	namedArguments []ast.NamedArgumentNode,
	location *position.Location,
) (
	_posArgs []ast.ExpressionNode,
) {
	reqParamCount := method.RequiredParamCount()
	requiredPosParamCount := len(method.Params) - method.OptionalParamCount
	if method.PostParamCount != -1 {
		requiredPosParamCount -= method.PostParamCount + 1
	}
	if method.HasNamedRestParam() {
		requiredPosParamCount--
	}
	argCount := len(positionalArguments) + len(namedArguments)
	positionalRestParamIndex := method.PositionalRestParamIndex()
	var checkedPositionalArguments []ast.ExpressionNode

	// push `undefined` for every missing optional positional argument
	// before the rest parameter
	for range positionalRestParamIndex - len(positionalArguments) {
		checkedPositionalArguments = append(
			checkedPositionalArguments,
			ast.NewUndefinedLiteralNode(location),
		)
	}

	var currentParamIndex int
	// check all positional arguments before the rest parameter
	for ; currentParamIndex < len(positionalArguments); currentParamIndex++ {
		posArg := positionalArguments[currentParamIndex]
		if currentParamIndex == positionalRestParamIndex {
			break
		}
		if currentParamIndex >= len(method.Params) {
			c.addWrongArgumentCountError(
				len(positionalArguments)+len(namedArguments),
				method,
				location,
			)
			break
		}
		param := method.Params[currentParamIndex]

		posArgType := c.MacroTypeOf(posArg)
		checkedPositionalArguments = append(checkedPositionalArguments, posArg)

		if !c.isSubtype(posArgType, param.Type, posArg.Location()) {
			c.addFailure(
				fmt.Sprintf(
					"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
					types.InspectWithColor(param.Type),
					param.Name.String(),
					types.InspectWithColor(method),
					types.InspectWithColor(posArgType),
				),
				posArg.Location(),
			)
		}
	}

	if method.HasPositionalRestParam() {
		if len(positionalArguments) < requiredPosParamCount {
			c.addFailure(
				fmt.Sprintf(
					"expected %d... positional arguments in call to `%s`, got %d",
					requiredPosParamCount,
					types.InspectWithColor(method),
					len(positionalArguments),
				),
				location,
			)
			return nil
		}
		restPositionalArguments := ast.NewArrayTupleLiteralNode(
			location,
			nil,
		)
		posRestParam := method.Params[positionalRestParamIndex]

		currentArgIndex := currentParamIndex
		// check rest arguments
		for ; currentArgIndex < min(argCount-method.PostParamCount, len(positionalArguments)); currentArgIndex++ {
			posArg := positionalArguments[currentArgIndex]
			posArgType := c.MacroTypeOf(posArg)
			restPositionalArguments.Elements = append(restPositionalArguments.Elements, posArg)
			if !c.isSubtype(posArgType, posRestParam.Type, posArg.Location()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for rest parameter `*%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(posRestParam.Type),
						posRestParam.Name.String(),
						types.InspectWithColor(method),
						types.InspectWithColor(posArgType),
					),
					posArg.Location(),
				)
			}
		}
		checkedPositionalArguments = append(checkedPositionalArguments, restPositionalArguments)

		currentParamIndex = positionalRestParamIndex
		// check post arguments
		for ; currentArgIndex < len(positionalArguments); currentArgIndex++ {
			posArg := positionalArguments[currentArgIndex]
			currentParamIndex++
			param := method.Params[currentParamIndex]

			posArgType := c.MacroTypeOf(posArg)
			checkedPositionalArguments = append(checkedPositionalArguments, posArg)
			if !c.isSubtype(posArgType, param.Type, posArg.Location()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(param.Type),
						param.Name.String(),
						types.InspectWithColor(method),
						types.InspectWithColor(posArgType),
					),
					posArg.Location(),
				)
			}
		}
		currentParamIndex++

		if method.PostParamCount > 0 {
			reqParamCount++
		}
	}

	firstNamedParamIndex := currentParamIndex
	definedNamedArgumentsSlice := make([]bool, len(namedArguments))

	for i := range method.Params {
		param := method.Params[i]
		switch param.Kind {
		case types.PositionalRestParameterKind, types.NamedRestParameterKind:
			continue
		}
		paramName := param.Name.String()
		var found bool

		for namedArgIndex, namedArgI := range namedArguments {
			var namedArg *ast.NamedCallArgumentNode
			switch n := namedArgI.(type) {
			case *ast.NamedCallArgumentNode:
				namedArg = n
			case *ast.DoubleSplatExpressionNode:
				continue
			default:
				panic(fmt.Sprintf("invalid named argument node: %T", namedArgI))
			}

			if namedArg.Name != paramName {
				continue
			}
			if found || i < firstNamedParamIndex {
				c.addFailure(
					fmt.Sprintf(
						"duplicated argument `%s` in call to `%s`",
						paramName,
						types.InspectWithColor(method),
					),
					namedArg.Location(),
				)
			}
			found = true
			definedNamedArgumentsSlice[namedArgIndex] = true

			namedArgType := c.MacroTypeOf(namedArg.Value)
			checkedPositionalArguments = append(checkedPositionalArguments, namedArg.Value)
			if !c.isSubtype(namedArgType, param.Type, namedArg.Location()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(param.Type),
						param.Name.String(),
						types.InspectWithColor(method),
						types.InspectWithColor(namedArgType),
					),
					namedArg.Location(),
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
			c.addFailure(
				fmt.Sprintf(
					"argument `%s` is missing in call to `%s`",
					paramName,
					types.InspectWithColor(method),
				),
				location,
			)
		} else {
			// the parameter is missing and is optional
			// we push undefined as its value
			checkedPositionalArguments = append(
				checkedPositionalArguments,
				ast.NewUndefinedLiteralNode(location),
			)
		}
	}

	if method.HasNamedRestParam() {
		namedRestArgs := ast.NewHashRecordLiteralNode(
			location,
			nil,
		)
		namedRestParam := method.Params[len(method.Params)-1]
		for i, defined := range definedNamedArgumentsSlice {
			if defined {
				continue
			}

			namedArgI := namedArguments[i]
			switch namedArg := namedArgI.(type) {
			case *ast.NamedCallArgumentNode:
				namedRestArgs.Elements = append(
					namedRestArgs.Elements,
					ast.NewSymbolKeyValueExpressionNode(
						namedArg.Location(),
						namedArg.Name,
						namedArg.Value,
					),
				)
				namedArgType := c.MacroTypeOf(namedArg.Value)
				c.checkNamedRestArgumentType(
					method.Name.String(),
					namedArgType,
					namedRestParam,
					namedArg.Location(),
				)
			case *ast.DoubleSplatExpressionNode:
				c.addFailure(
					fmt.Sprintf(
						"double splat arguments cannot be used in macro call `%s`",
						types.InspectWithColor(method),
					),
					namedArg.Location(),
				)
			default:
				panic(fmt.Sprintf("invalid named argument node: %T", namedArgI))
			}
		}

		checkedPositionalArguments = append(checkedPositionalArguments, namedRestArgs)
	} else {
		for i, defined := range definedNamedArgumentsSlice {
			if defined {
				continue
			}

			namedArgI := namedArguments[i]
			switch namedArg := namedArgI.(type) {
			case *ast.NamedCallArgumentNode:
				c.addFailure(
					fmt.Sprintf(
						"nonexistent parameter `%s` given in call to `%s`",
						namedArg.Name,
						types.InspectWithColor(method),
					),
					namedArg.Location(),
				)
			case *ast.DoubleSplatExpressionNode:
				c.addFailure(
					fmt.Sprintf(
						"double splat arguments cannot be used in macro call `%s`",
						types.InspectWithColor(method),
					),
					namedArg.Location(),
				)
			}
		}
	}

	return checkedPositionalArguments
}

func (c *Checker) hoistMacroDefinition(node *ast.MacroDefinitionNode) {
	c.setDefinedMacros(true)
	definedUnder := c.currentMethodScope().container
	switch d := definedUnder.(type) {
	case *types.Module:
	case *types.Class, *types.Mixin:
		definedUnder = d.Singleton()
	default:
		c.addFailure(
			fmt.Sprintf(
				"cannot declare macro `%s` in this context",
				node.Name,
			),
			node.Location(),
		)
	}

	macro := c.declareMacro(
		definedUnder,
		node.DocComment(),
		value.ToSymbol(node.Name+"!"),
		node.Parameters,
		node.Location(),
	)
	macro.Node = node
	node.SetType(macro)
	c.registerMacroCheck(macro, node)
}

type macroCheckEntry struct {
	macro          *types.Method
	constantScopes []constantScope
	methodScopes   []methodScope
	node           *ast.MacroDefinitionNode
}

func (c *Checker) registerMacroCheck(macro *types.Method, node *ast.MacroDefinitionNode) {
	c.macroChecks = append(c.macroChecks, macroCheckEntry{
		macro:          macro,
		constantScopes: c.constantScopesCopy(),
		methodScopes:   c.methodScopesCopy(),
		node:           node,
	})
}

// Check macro definition bodies
func (c *Checker) checkMacros() {
	concurrent.Foreach(
		concurrencyLimit,
		c.macroChecks,
		func(macroCheck macroCheckEntry) {
			macro := macroCheck.macro
			node := macroCheck.node
			macroChecker := c.newMethodChecker(
				node.Location().FilePath,
				macroCheck.constantScopes,
				macroCheck.methodScopes,
				c.StdNode().Singleton(),
				macro.ReturnType,
				macro.ThrowType,
				false,
				c.threadPool,
			)
			macroChecker.checkMacroDefinition(node, macro)
		},
	)

	c.macroChecks = nil
	c.compiler = nil
}

func (c *Checker) checkMacroDefinition(node *ast.MacroDefinitionNode, macro *types.Method) {
	c.method = macro
	c.checkMethod(
		c.currentMethodScope().container,
		macro,
		node.Parameters,
		nil,
		nil,
		node.Body,
		node.Location(),
	)

	c.method = nil

	if c.shouldCompile() && macro.IsCompilable() {
		macro.Bytecode = c.compiler.CompileMacroBody(node, macro.Name)
	}
}

func (c *Checker) declareMacro(
	macroNamespace types.Namespace,
	docComment string,
	name value.Symbol,
	paramNodes []ast.ParameterNode,
	location *position.Location,
) *types.Method {
	exprNodeType := c.StdExpressionNode()

	var params []*types.Parameter
	for _, paramNode := range paramNodes {
		switch p := paramNode.(type) {
		case *ast.FormalParameterNode:
			var declaredType types.Type
			if p.TypeNode != nil {
				p.TypeNode = c.checkTypeNode(p.TypeNode)
				declaredType = c.TypeOf(p.TypeNode)
			} else {
				c.addFailure(
					fmt.Sprintf("cannot declare parameter `%s` without a type", p.Name),
					paramNode.Location(),
				)
			}

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

			if !c.isSubtype(declaredType, exprNodeType, p.Location()) {
				c.addFailure(
					fmt.Sprintf(
						"type `%s` does not inherit from `%s`, macro parameters must be expression nodes",
						types.InspectWithColor(declaredType),
						types.InspectWithColor(exprNodeType),
					),
					p.Location(),
				)
			}

			name := value.ToSymbol(p.Name)
			paramType := types.NewParameter(
				name,
				declaredType,
				kind,
				false,
			)
			p.SetType(paramType)
			params = append(params, paramType)
		default:
			c.addFailure(
				fmt.Sprintf("invalid param type %T", paramNode),
				paramNode.Location(),
			)
		}
	}

	newMacro := types.NewMethod(
		docComment,
		0,
		name,
		nil,
		params,
		exprNodeType,
		types.Never{},
		macroNamespace,
	)
	newMacro.SetLocation(location)
	macroNamespace.SetMethod(name, newMacro)
	return newMacro
}
