package checker

import (
	"fmt"

	"github.com/elk-language/elk/compiler"
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
	prevIsHeader := c.IsHeader()

	c.Filename = filename
	c.setIsHeaderForPath(filename)

	c.expandTopLevelMacros(node.Body)

	c.Filename = prevFilename
	c.SetHeader(prevIsHeader)

	node.State = ast.EXPANDED_TOP_LEVEL_MACROS
}

func (c *Checker) checkUsingExpressionForMacros(node *ast.UsingExpressionNode) {
	for _, entry := range node.Entries {
		c.resolveUsingEntry(entry, false)
		switch e := entry.(type) {
		case *ast.MethodLookupNode:
			if !e.IsMacro() {
				continue
			}
			c.checkUsingMethodLookupEntryNode(
				e.Receiver,
				c.identifierToName(e.Name),
				"",
				e.Location(),
			)
		case *ast.MethodLookupAsNode:
			if !e.IsMacro() {
				continue
			}
			c.checkUsingMethodLookupEntryNode(
				e.MethodLookup.Receiver,
				c.identifierToName(e.MethodLookup.Name),
				c.identifierToName(e.AsName),
				e.Location(),
			)
		case *ast.UsingEntryWithSubentriesNode:
			c.checkUsingEntryWithSubentriesForMacros(e)
		}
	}
}

func (c *Checker) checkUsingEntryWithSubentriesForMacros(node *ast.UsingEntryWithSubentriesNode) {
	for _, subentry := range node.Subentries {
		switch s := subentry.(type) {
		case *ast.MacroNameNode:
			c.checkUsingMethodLookupEntryNode(node.Namespace, c.identifierToName(s), "", s.Location())
		case *ast.UsingSubentryAsNode:
			if !s.IsMacro() {
				continue
			}
			value := c.identifierToName(s.Target)
			asName := c.identifierToName(s.AsName)
			c.checkUsingMethodLookupEntryNode(node.Namespace, value, asName, s.Location())
		case *ast.PublicConstantNode, *ast.PublicConstantAsNode, *ast.PublicIdentifierNode:
		default:
			panic(fmt.Sprintf("invalid using subentry node: %T", subentry))
		}
	}
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
	case *ast.UsingExpressionNode:
		c.checkUsingExpressionForMacros(expr)
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
		c.hoistSingletonDeclarationAndExpandMacros(expr)
	case *ast.ExtendWhereBlockExpressionNode:
		c.expandTopLevelMacros(expr.Body)
	case *ast.StructDeclarationNode:
		c.hoistStructDeclaration(expr)
	case *ast.MacroCallNode:
		posArgs := []ast.ExpressionNode{expr.Receiver}

		result := c.expandMacroByName(
			c.identifierToName(expr.MacroName),
			expr.Kind,
			append(posArgs, expr.PositionalArguments...),
			expr.NamedArguments,
			expr.Location(),
		)
		expr.SetType(types.Untyped{})
		if result == nil {
			return expr
		}
		return c.expandTopLevelMacrosInExpression(result.(ast.ExpressionNode))
	case *ast.ReceiverlessMacroCallNode:
		result := c.expandMacroByName(
			c.identifierToName(expr.MacroName),
			expr.Kind,
			expr.PositionalArguments,
			expr.NamedArguments,
			expr.Location(),
		)
		expr.SetType(types.Untyped{})
		if result == nil {
			return expr
		}
		return c.expandTopLevelMacrosInExpression(result.(ast.ExpressionNode))
	}

	return expr
}

func (c *Checker) resolveMacro(name value.Symbol) *types.Method {
	for methodScope := range ds.ReverseSlice(c.methodScopes) {
		macro := c.resolveMacroForNamespace(methodScope.container, name)
		if macro != nil {
			return macro
		}
	}

	return nil
}

func (c *Checker) resolveMacroForNamespace(namespace types.Namespace, name value.Symbol) *types.Method {
	switch n := namespace.(type) {
	case *types.Class:
		namespace = n.Singleton()
	case *types.Mixin:
		namespace = n.Singleton()
	case *types.Module, *types.SingletonClass, *types.UsingBufferNamespace:
		namespace = n
	default:
		return nil
	}

	for parent := range types.Parents(namespace) {
		macro := parent.Method(name)
		if macro != nil {
			return macro
		}
	}

	return nil
}

func macroName(namespace types.Namespace, name string) string {
	var namespaceName string
	switch n := namespace.(type) {
	case *types.SingletonClass:
		namespaceName = n.AttachedObject.Name()
	default:
		namespaceName = namespace.Name()
	}
	return fmt.Sprintf("%s::%s", namespaceName, name)
}

func (c *Checker) getMacroForNamespace(namespace types.Namespace, name value.Symbol, loc *position.Location) *types.Method {
	macro := c.resolveMacroForNamespace(namespace, name)
	if macro == nil {
		c.addUndefinedMacroError(macroName(namespace, name.String()), loc)
	}

	return macro
}

func (c *Checker) addUndefinedMacroError(name string, loc *position.Location) {
	c.addFailure(
		fmt.Sprintf(
			"undefined macro `%s`",
			lexer.Colorize(name),
		),
		loc,
	)
}

func (c *Checker) getMacro(name value.Symbol, loc *position.Location) *types.Method {
	macro := c.resolveMacro(name)
	if macro == nil {
		c.addUndefinedMacroError(name.String(), loc)
	}

	return macro
}

func (c *Checker) macroMethodName(macroName string) value.Symbol {
	return value.ToSymbol(macroName + "!")
}

func (c *Checker) expandMacroByName(name string, kind ast.MacroKind, posArgs []ast.ExpressionNode, namedArgs []ast.NamedArgumentNode, loc *position.Location) ast.Node {
	macroName := c.macroMethodName(name)
	macro := c.getMacro(macroName, loc)
	if macro == nil {
		return nil
	}
	if macro.IsPlaceholder() {
		return nil
	}

	return c.expandMacro(macro, kind, posArgs, namedArgs, loc)
}

func (c *Checker) expandMacro(macro *types.Method, kind ast.MacroKind, posArgs []ast.ExpressionNode, namedArgs []ast.NamedArgumentNode, loc *position.Location) ast.Node {
	exprNodeType := c.StdExpressionNode()
	patternNodeType := c.StdPatternNode()
	typeNodeType := c.StdTypeNode()

	checkedArgs := c.checkMacroArguments(macro, posArgs, namedArgs, loc)

	var expectedReturnType types.Type
	switch kind {
	case ast.MACRO_EXPRESSION_KIND:
		expectedReturnType = exprNodeType
	case ast.MACRO_PATTERN_KIND:
		expectedReturnType = patternNodeType
	case ast.MACRO_TYPE_KIND:
		expectedReturnType = typeNodeType
	}
	c.checkCanAssign(macro.ReturnType, expectedReturnType, loc)

	if c.Errors.IsFailure() {
		return nil
	}

	runtimeArgs := make([]value.Value, 0, len(checkedArgs)+1)
	runtimeArgs = append(
		runtimeArgs,
		value.Ref(value.NodeMixin),
	)

	for _, arg := range checkedArgs {
		switch arg := arg.(type) {
		case *ast.UndefinedLiteralNode:
			runtimeArgs = append(runtimeArgs, value.Undefined)
		default:
			runtimeArgs = append(runtimeArgs, value.Ref(arg))
		}
	}

	var promise *vm.Promise
	switch body := macro.Body.(type) {
	case *vm.NativeMethod:
		promise = vm.NewNativePromise(c.threadPool, body.Function, runtimeArgs...)
	case *vm.BytecodeFunction:
		promise = vm.NewBytecodePromise(c.threadPool, body, runtimeArgs...)
	default:
		panic(fmt.Sprintf("invalid compiled macro body: %T", body))
	}

	result, stackTrace, err := promise.AwaitSync()
	if !err.IsUndefined() {
		c.addFailure(
			fmt.Sprintf(
				"error while executing macro `%s`: %s\n%s",
				types.InspectWithColor(macro),
				lexer.Colorize(err.Inspect()),
				stackTrace.String(),
			),
			loc,
		)
		return nil
	}

	resultNode := result.AsReference().(ast.Node)

	switch macro.ReturnType {
	case exprNodeType:
		switch r := resultNode.(type) {
		case *ast.DoExpressionNode:
			if r.HasSingleScope() {
				resultNode = ast.NewMacroBoundaryNode(
					resultNode.Location(),
					r.Body,
					types.Inspect(macro),
				)
			} else {
				resultNode = ast.NewMacroBoundaryNode(
					resultNode.Location(),
					ast.ExpressionToStatements(r),
					types.Inspect(macro),
				)
			}
		case ast.ExpressionNode:
			resultNode = ast.NewMacroBoundaryNode(
				resultNode.Location(),
				ast.ExpressionToStatements(r),
				types.Inspect(macro),
			)
		default:
			panic(fmt.Sprintf("invalid expression macro result type: %T", resultNode))
		}
	case patternNodeType:
		r := resultNode.(ast.PatternNode)
		resultNode = ast.NewMacroBoundaryNode(
			resultNode.Location(),
			ast.PatternToStatements(r),
			types.Inspect(macro),
		)
	case typeNodeType:
		r := resultNode.(ast.TypeNode)
		resultNode = ast.NewMacroBoundaryNode(
			resultNode.Location(),
			ast.TypeToStatements(r),
			types.Inspect(macro),
		)
	}

	// update location
	resultNode = ast.Splice(resultNode, loc, nil)

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

			if c.identifierToName(namedArg.Name) != paramName {
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
						namedArg.Name.String(),
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
		c.macroMethodName(c.identifierToName(node.Name)),
		node.Parameters,
		node.ReturnType,
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
		MethodCheckConcurrencyLimit,
		c.macroChecks,
		func(macroCheck macroCheckEntry) {
			macro := macroCheck.macro
			node := macroCheck.node
			macroChecker := c.newMacroChecker(
				node.Location().FilePath,
				macroCheck.constantScopes,
				macroCheck.methodScopes,
				c.StdNode().Singleton(),
				macro.ReturnType,
				macro.ThrowType,
				macroMode,
				c.threadPool,
				node.Location(),
			)
			macroChecker.checkMacroDefinition(node, macro)
		},
	)

	c.macroChecks = nil
}

func (c *Checker) newMacroChecker(
	filename string,
	constScopes []constantScope,
	methodScopes []methodScope,
	selfType,
	returnType,
	throwType types.Type,
	mode mode,
	threadPool *vm.ThreadPool,
	loc *position.Location,
) *Checker {
	checker := &Checker{
		env:            c.env,
		Filename:       filename,
		mode:           mode,
		phase:          methodCheckPhase,
		selfType:       selfType,
		returnType:     returnType,
		throwType:      throwType,
		constantScopes: constScopes,
		methodScopes:   methodScopes,
		Errors:         c.Errors,
		flags:          c.flags,
		localEnvs: []*localEnvironment{
			newLocalEnvironment(nil, false),
		},
		typeDefinitionChecks: newTypeDefinitionChecks(),
		methodCache:          concurrent.NewSlice[*types.Method](),
		threadPool:           threadPool,
	}
	checker.macroCompiler = compiler.CreateBytecodeCompiler(nil, checker, loc, c.Errors)

	return checker
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

	if c.shouldCompileMacro() && macro.IsCompilable() {
		macro.Body = c.macroCompiler.CompileMacroBody(node, macro.Name)
	}
}

func (c *Checker) declareMacro(
	macroNamespace types.Namespace,
	docComment string,
	name value.Symbol,
	paramNodes []ast.ParameterNode,
	returnTypeNode ast.TypeNode,
	location *position.Location,
) *types.Method {
	nodeType := c.StdNode()
	exprNodeType := c.StdExpressionNode()
	typeNodeType := c.StdTypeNode()
	patternNodeType := c.StdPatternNode()

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

			if !c.isSubtype(declaredType, nodeType, p.Location()) {
				c.addFailure(
					fmt.Sprintf(
						"type `%s` does not inherit from `%s`, macro parameters must be nodes",
						types.InspectWithColor(declaredType),
						types.InspectWithColor(nodeType),
					),
					p.Location(),
				)
			}

			name := value.ToSymbol(c.identifierToName(p.Name))
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

	prevMode := c.mode
	c.setOutputPositionTypeMode()

	var returnType types.Type
	var typedReturnTypeNode ast.TypeNode
	if returnTypeNode != nil {
		typedReturnTypeNode = c.checkTypeNode(returnTypeNode)
		returnType = c.TypeOf(typedReturnTypeNode)
		switch returnType {
		case exprNodeType, patternNodeType, typeNodeType:
		default:
			c.addFailure(
				fmt.Sprintf(
					"invalid macro return type, got %s, should be %s, %s or %s",
					types.InspectWithColor(returnType),
					types.InspectWithColor(exprNodeType),
					types.InspectWithColor(patternNodeType),
					types.InspectWithColor(typeNodeType),
				),
				returnTypeNode.Location(),
			)
		}
	} else {
		returnType = exprNodeType
	}

	newMacro := types.NewMethod(
		docComment,
		0,
		name,
		nil,
		params,
		returnType,
		types.Never{},
		macroNamespace,
	)
	newMacro.SetMacro(true)
	newMacro.SetLocation(location)
	macroNamespace.SetMethod(name, newMacro)

	c.mode = prevMode

	return newMacro
}

func (c *Checker) checkMacroCallNode(node *ast.MacroCallNode) ast.Node {
	node.SetType(types.Untyped{})

	if c.mode == macroMode {
		c.addMacroInMacroError(node.MacroName.Location())
		return nil
	}

	posArgs := []ast.ExpressionNode{node.Receiver}
	result := c.expandMacroByName(
		c.identifierToName(node.MacroName),
		node.Kind,
		append(posArgs, node.PositionalArguments...),
		node.NamedArguments,
		node.Location(),
	)
	if result == nil {
		return nil
	}
	return result
}

func (c *Checker) checkMacroCallNodeForExpression(node *ast.MacroCallNode) ast.ExpressionNode {
	result := c.checkMacroCallNode(node)
	if result == nil {
		return node
	}
	return c.checkExpression(result.(ast.ExpressionNode))
}

func (c *Checker) checkMacroCallNodeForType(node *ast.MacroCallNode) ast.TypeNode {
	result := c.checkMacroCallNode(node)
	if result == nil {
		return node
	}
	return c.checkTypeNode(result.(ast.TypeNode))
}

func (c *Checker) checkMacroCallNodeForPattern(node *ast.MacroCallNode, matchedType types.Type) (ast.PatternNode, types.Type) {
	result := c.checkMacroCallNode(node)
	if result == nil {
		return node, types.Never{}
	}
	return c.checkPattern(result.(ast.PatternNode), matchedType)
}

func (c *Checker) addMacroInMacroError(loc *position.Location) {
	c.addFailure(
		"macros cannot be used in macro definitions",
		loc,
	)
}

func (c *Checker) checkScopedMacroCallNode(node *ast.ScopedMacroCallNode) ast.Node {
	node.SetType(types.Untyped{})

	if c.mode == macroMode {
		c.addMacroInMacroError(node.MacroName.Location())
		return nil
	}

	node.Receiver = c.checkExpression(node.Receiver)
	receiverType := c.TypeOf(node.Receiver)

	var receiverNamespace types.Namespace
	switch r := receiverType.(type) {
	case types.Namespace:
		receiverNamespace = r
	case types.Untyped:
		return nil
	default:
		c.addFailure(
			fmt.Sprintf("invalid macro scope %T", receiverType),
			node.Receiver.Location(),
		)
		return nil
	}

	macroName := c.macroMethodName(c.identifierToName(node.MacroName))
	macro := c.getMacroForNamespace(receiverNamespace, macroName, node.MacroName.Location())
	if macro == nil {
		return nil
	}

	result := c.expandMacro(
		macro,
		node.Kind,
		node.PositionalArguments,
		node.NamedArguments,
		node.Location(),
	)
	if result == nil {
		return nil
	}
	return result
}

func (c *Checker) checkScopedMacroCallNodeForExpression(node *ast.ScopedMacroCallNode) ast.ExpressionNode {
	result := c.checkScopedMacroCallNode(node)
	if result == nil {
		return node
	}
	return c.checkExpression(result.(ast.ExpressionNode))
}

func (c *Checker) checkScopedMacroCallNodeForType(node *ast.ScopedMacroCallNode) ast.TypeNode {
	result := c.checkScopedMacroCallNode(node)
	if result == nil {
		return node
	}
	return c.checkTypeNode(result.(ast.TypeNode))
}

func (c *Checker) checkScopedMacroCallNodeForPattern(node *ast.ScopedMacroCallNode, matchedType types.Type) (ast.PatternNode, types.Type) {
	result := c.checkScopedMacroCallNode(node)
	if result == nil {
		return node, types.Never{}
	}
	return c.checkPattern(result.(ast.PatternNode), matchedType)
}

func (c *Checker) checkReceiverlessMacroCallNode(node *ast.ReceiverlessMacroCallNode) ast.Node {
	node.SetType(types.Untyped{})

	if c.mode == macroMode {
		c.addMacroInMacroError(node.MacroName.Location())
		return nil
	}

	result := c.expandMacroByName(
		c.identifierToName(node.MacroName),
		node.Kind,
		node.PositionalArguments,
		node.NamedArguments,
		node.Location(),
	)
	if result == nil {
		return nil
	}
	return result
}

func (c *Checker) checkReceiverlessMacroCallNodeForExpression(node *ast.ReceiverlessMacroCallNode) ast.ExpressionNode {
	result := c.checkReceiverlessMacroCallNode(node)

	if result == nil {
		return node
	}
	return c.checkExpression(result.(ast.ExpressionNode))
}

func (c *Checker) checkReceiverlessMacroCallNodeForType(node *ast.ReceiverlessMacroCallNode) ast.TypeNode {
	result := c.checkReceiverlessMacroCallNode(node)
	if result == nil {
		return node
	}
	return c.checkTypeNode(result.(ast.TypeNode))
}

func (c *Checker) checkReceiverlessMacroCallNodeForPattern(node *ast.ReceiverlessMacroCallNode, matchedType types.Type) (ast.PatternNode, types.Type) {
	result := c.checkReceiverlessMacroCallNode(node)
	if result == nil {
		return node, types.Never{}
	}
	return c.checkPattern(result.(ast.PatternNode), matchedType)
}
