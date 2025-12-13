package compiler

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strings"
	"unicode"

	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

type GoSourceMethod GoCompiler

func (c *GoSourceMethod) Name() value.Symbol                        { return value.ToSymbol(c.FuncName) }
func (*GoSourceMethod) Class() *value.Class                         { return nil }
func (*GoSourceMethod) Copy() value.Reference                       { return nil }
func (*GoSourceMethod) DirectClass() *value.Class                   { return nil }
func (*GoSourceMethod) Error() string                               { return "" }
func (*GoSourceMethod) Inspect() string                             { return "" }
func (*GoSourceMethod) InstanceVariables() *value.InstanceVariables { return nil }
func (*GoSourceMethod) OptionalParameterCount() int                 { return 0 }
func (*GoSourceMethod) ParameterCount() int                         { return 0 }
func (*GoSourceMethod) SingletonClass() *value.Class                { return nil }

type nativeBigInt struct {
	id  int
	val string
}

func (s *nativeBigInt) goIdent() string {
	return fmt.Sprintf("bi%d", s.id)
}

type nativeSymbol struct {
	id  int
	val string
}

func (s *nativeSymbol) goIdent() string {
	return fmt.Sprintf("sym%d", s.id)
}

// Compiler mode
type goMode uint8

const (
	topLevelGoCompilerMode goMode = iota
	methodGoCompilerMode
	initMethodGoCompilerMode
	setterMethodGoCompilerMode
)

// represents a nativeElkLocal variable or value
type nativeElkLocal struct {
	index int
	name  string
	typ   types.Type
}

func (n *nativeElkLocal) goIdent() string {
	return fmt.Sprintf("l%d", n.index)
}

type nativeElkLocalTable map[string]*nativeElkLocal

type nativeElkScopeType uint8

const (
	defaultNativeElkScopeType       nativeElkScopeType = iota
	loopNativeElkScopeType                             // this scope is a loop
	doFinallyNativeElkScopeType                        // this scope is inside do with a finally block
	macroBoundaryNativeElkScopeType                    // this scope is a macro boundary, locals from the outer scopes should be ignored
)

// set of local variables
type nativeElkScope struct {
	localTable nativeElkLocalTable
	label      string
	typ        nativeElkScopeType
}

func newNativeElkScope(label string, typ nativeElkScopeType) *nativeElkScope {
	return &nativeElkScope{
		localTable: nativeElkLocalTable{},
		label:      label,
		typ:        typ,
	}
}

// indices represent scope depths
// and elements are sets of local variable names in a particular scope
type nativeElkScopes []*nativeElkScope

// Get the last local variable scope.
func (s nativeElkScopes) last() *nativeElkScope {
	return s[len(s)-1]
}

const goValueType = "value.Value"

type goLocal struct {
	name    string
	elkType types.Type
	goType  string
}

func newGoLocal(name string, elkType types.Type, goType string) *goLocal {
	return &goLocal{
		name:    name,
		elkType: elkType,
		goType:  goType,
	}
}

type goValue struct {
	tmpLocal string
	inline   string
	typ      types.Type
}

func (v *goValue) IsInline() bool {
	return v.inline != ""
}

func (v *goValue) IsTmp() bool {
	return v.tmpLocal != ""
}

func (v *goValue) Value() string {
	if v.inline != "" {
		return v.inline
	}

	return v.tmpLocal
}

func newInlineValue(v string, typ types.Type) *goValue {
	return &goValue{
		inline: v,
		typ:    typ,
	}
}

func newTmpValue(v string, typ types.Type) *goValue {
	return &goValue{
		tmpLocal: v,
		typ:      typ,
	}
}

func CreateGoCompiler(parent *GoCompiler, checker types.Checker, loc *position.Location, errors *diagnostic.SyncDiagnosticList, output io.Writer) *GoCompiler {
	bigIntCache := concurrent.NewMap[string, *nativeBigInt]()
	symbolCache := concurrent.NewMap[string, *nativeSymbol]()
	compiler := NewGoCompiler("main", topLevelGoCompilerMode, loc, checker, bigIntCache, symbolCache, output)
	compiler.Errors = errors
	compiler.SetParent(parent)
	return compiler
}

func (c *GoCompiler) CreateMainCompiler(checker types.Checker, loc *position.Location, errors *diagnostic.SyncDiagnosticList, output io.Writer) Compiler {
	bigIntCache := concurrent.NewMap[string, *nativeBigInt]()
	symbolCache := concurrent.NewMap[string, *nativeSymbol]()
	compiler := NewGoCompiler("main", topLevelGoCompilerMode, loc, checker, bigIntCache, symbolCache, output)
	compiler.Errors = errors
	return compiler
}

func (c *GoCompiler) InitMainCompiler() {
	c.emitPackage("package main\n\n")
}

func (c *GoCompiler) InitGlobalEnv() Compiler {
	envCompiler := NewGoCompiler("initGlobalEnv", topLevelGoCompilerMode, c.loc, c.checker, c.bigIntCache, c.symbolCache, c.output)
	envCompiler.SetParent(c)
	envCompiler.Errors = c.Errors
	envCompiler.compileGlobalEnv()
	return envCompiler
}

func (c *GoCompiler) FinishGlobalEnvCompiler() {
	parent := c.parent
	parent.emit("%s()\n", c.FuncName)
}

func (c *GoCompiler) InitMethodCompiler(location *position.Location) (Compiler, int) {
	methodCompiler := NewGoCompiler("methodDefinitions", topLevelGoCompilerMode, c.loc, c.checker, c.bigIntCache, c.symbolCache, c.output)
	methodCompiler.Errors = c.Errors
	methodCompiler.SetParent(c)

	c.emit("methodDefinitions()\n")

	return methodCompiler, 0
}

func (c *GoCompiler) CompileMethods(location *position.Location, execOffset int) {
	c.registerGoLocal("class", c.checker.Std(symbol.Class), "*value.Class")
	c.compileMethodsWithinModule(c.checker.Env().Root, location)
}

func (c *GoCompiler) InitIvarIndicesCompiler(location *position.Location) (Compiler, int) {
	ivarCompiler := NewGoCompiler("ivarIndices", topLevelGoCompilerMode, c.loc, c.checker, c.bigIntCache, c.symbolCache, c.output)
	ivarCompiler.Errors = c.Errors
	ivarCompiler.SetParent(c)

	c.emit("ivarIndices()\n")

	return ivarCompiler, 0
}

func (c *GoCompiler) FinishIvarIndicesCompiler(location *position.Location, execOffset int) Compiler {
	return c.parent
}

func (c *GoCompiler) CompileConstantDeclaration(node *ast.ConstantDeclarationNode, namespace types.Namespace, constName value.Symbol) {
	c.registerGoLocal("namespace", c.checker.Std(symbol.Value), goValueType)

	switch n := namespace.(type) {
	case *types.SingletonClass:
		namespaceSymbol := c.emitSymbol(n.AttachedObject.Name())
		c.emit("namespace = value.Ref(value.GetConstant(%s).SingletonClass())\n", namespaceSymbol)
	default:
		namespaceSymbol := c.emitSymbol(n.Name())
		c.emit("namespace = value.GetConstant(%s)\n", namespaceSymbol)
	}

	init := c.compileExpression(node.Initialiser)
	constNameSymbol := c.emitSymbol(constName.String())
	c.emit("value.AddConstant(namespace, %s, %s)\n", constNameSymbol, init.Value())
}

func (c *GoCompiler) CompileMethodBody(node *ast.MethodDefinitionNode, name value.Symbol) Compiler {
	var mode goMode
	if node.IsSetter() {
		mode = setterMethodGoCompilerMode
	} else if identifierToName(node.Name) == "#init" {
		mode = initMethodGoCompilerMode
	} else {
		mode = methodGoCompilerMode
	}

	methodCompiler := NewGoCompiler(name.String(), mode, node.Location(), c.checker, c.bigIntCache, c.symbolCache, c.output)
	methodCompiler.isGenerator = node.IsGenerator()
	methodCompiler.isAsync = node.IsAsync()
	methodCompiler.Errors = c.Errors
	methodCompiler.compileMethodBody(node.Parameters, node.Body)

	return methodCompiler
}

// Compiles Elk source code to Go source code.
type GoCompiler struct {
	Errors            *diagnostic.SyncDiagnosticList
	FuncName          string
	scopes            nativeElkScopes
	output            io.Writer
	parent            *GoCompiler
	buff              bytes.Buffer // inner function code
	packageBuff       bytes.Buffer // package level code
	checker           types.Checker
	loc               *position.Location
	mode              goMode
	children          concurrent.Slice[*GoCompiler]
	goLocals          map[string]*goLocal
	tmpLocalCounter   int
	lastElkLocalIndex int
	callCacheCounter  int
	bigIntCache       *concurrent.Map[string, *nativeBigInt]
	symbolCache       *concurrent.Map[string, *nativeSymbol]
	isGenerator       bool
	isAsync           bool
}

func NewGoCompiler(name string, mode goMode, loc *position.Location, checker types.Checker, bigIntCache *concurrent.Map[string, *nativeBigInt], symbolCache *concurrent.Map[string, *nativeSymbol], output io.Writer) *GoCompiler {
	return &GoCompiler{
		FuncName:          name,
		mode:              mode,
		Errors:            diagnostic.NewSyncDiagnosticList(),
		scopes:            nativeElkScopes{newNativeElkScope("", defaultNativeElkScopeType)}, // start with an empty set for the 0th scope
		lastElkLocalIndex: -1,
		bigIntCache:       bigIntCache,
		symbolCache:       symbolCache,
		checker:           checker,
		output:            output,
		loc:               loc,
	}
}

func (c *GoCompiler) Parent() Compiler {
	if c.parent == nil {
		return nil
	}
	return c.parent
}

func (c *GoCompiler) Bytecode() *vm.BytecodeFunction {
	panic("cannot get bytecode from a GoCompiler")
}

func (c *GoCompiler) Method() value.Method {
	return (*GoSourceMethod)(c)
}

func (c *GoCompiler) Flush() {
	c.flushPackage()
	c.output.Write([]byte("\n"))
	c.flushInner()
}

func (c *GoCompiler) flushPackage() {
	c.output.Write(c.packageBuff.Bytes())
	c.output.Write([]byte("\n"))
	c.packageBuff.Reset()

	for _, child := range c.children.Slice {
		child.flushPackage()
	}
}

func (c *GoCompiler) flushInner() {
	c.output.Write(c.buff.Bytes())
	c.output.Write([]byte("\n"))
	c.packageBuff.Reset()

	for _, child := range c.children.Slice {
		child.flushInner()
	}
}

func (c *GoCompiler) SetParent(parent Compiler) {
	p := parent.(*GoCompiler)
	c.parent = p
	p.children.Append(c)
}

func (c *GoCompiler) registerGoLocal(name string, elkType types.Type, goType string) {
	_, exists := c.goLocals[name]
	if exists {
		return
	}
	c.goLocals[name] = newGoLocal(name, elkType, goType)
}

func (c *GoCompiler) compileGlobalEnv() {
	env := c.checker.Env()
	c.emit("func initGlobalEnv() {\n")

	c.emit("var parentNamespace, namespace value.Value\n")
	c.emit("var class, superclass, mixin *value.Class\n")
	c.compileModuleDefinition(env.Root, env.Root, value.ToSymbol("Root"))

	c.emit("}\n")
}

// Entry point for compiling the body of a method.
func (c *GoCompiler) compileMethodBody(parameters []ast.ParameterNode, body []ast.StatementNode) {
	c.compileMethodFuncLiteralBody(parameters, body)

	var funcBuffer bytes.Buffer
	fmt.Fprintf(&funcBuffer, "func(thread *VM, args []value.Value) (value.Value, value.Value) {\n")
	c.compileLocalsTo(&funcBuffer)
	c.emitPrependBytes(funcBuffer.Bytes())
	c.emit("}")
}

func (c *GoCompiler) compileMethodFuncLiteralBody(parameters []ast.ParameterNode, body []ast.StatementNode) {
	c.registerGoLocal("self", c.checker.SelfType(), goValueType)
	c.emit("self = args[0]\n")

	for i, param := range parameters {
		p := param.(*ast.MethodParameterNode)
		pSpan := p.Location()

		pName := identifierToName(p.Name)
		local := c.defineLocal(pName, c.typeOf(p), pSpan)
		if local == nil {
			return
		}

		localName := local.goIdent()
		c.emit("%s = args[%d]\n", localName, i+1)

		if p.Initialiser != nil {
			c.emit("if %s.IsUndefined() {\n", localName)
			val := c.compileExpression(p.Initialiser)
			c.emit("%s = %s", localName, val.Value())
			c.emit("}\n")
		}

		if p.SetInstanceVariable {
			c.emitSetInstanceVariable(value.ToSymbol(pName), localName)
		}
	}

	// TODO: implement async and generators

	// paramCount := len(parameters)
	// if c.isGenerator {
	// 	c.emit(location.StartPos.Line, bytecode.GENERATOR)
	// 	c.emit(location.EndPos.Line, bytecode.RETURN)
	// 	c.registerCatch(-1, -1, c.nextInstructionOffset(), false)
	// } else if c.isAsync {
	// 	poolVar := c.defineLocal("_pool", location)
	// 	paramCount++
	// 	c.predefinedLocals++
	// 	c.bytecode.IncrementOptionalParameterCount()

	// 	c.emitGetLocal(location.StartPos.Line, poolVar.index)
	// 	c.emit(location.StartPos.Line, bytecode.PROMISE)
	// 	c.emit(location.EndPos.Line, bytecode.RETURN)
	// }

	val := c.compileStatements(body)

	c.emitReturn(val.Value())
	// TODO: implement generators
	// c.emitFinalReturn(location, nil)
}

func (c *GoCompiler) emitReturn(val string) {
	// TODO: implement generators
	// if c.isGenerator {
	// 	c.emitYield(location, value)
	// 	c.emit(location.EndPos.Line, bytecode.STOP_ITERATION)
	// 	return
	// }

	switch c.mode {
	case setterMethodGoCompilerMode:
		// TODO: implement finally
		// if c.isNestedInFinally() {
		// 	c.emit(location.EndPos.Line, bytecode.GET_LOCAL8, 1)
		// 	c.emit(location.EndPos.Line, bytecode.RETURN_FINALLY)
		// } else {
		c.emit("return args[1], value.Undefined\n")
		// }
	case initMethodGoCompilerMode:
		// TODO: implement finally
		// if c.isNestedInFinally() {
		// 	c.emit(location.EndPos.Line, bytecode.SELF)
		// 	c.emit(location.EndPos.Line, bytecode.RETURN_FINALLY)
		// } else {
		c.emit("return self, value.Undefined\n")
		// }
	// TODO: implement namespaces
	// case namespaceBytecodeCompilerMode:
	// 	if value != nil {
	// 		c.compileNodeWithResult(value)
	// 	}
	// 	if c.lastOpCode != bytecode.NIL {
	// 		c.emit(location.EndPos.Line, bytecode.POP)
	// 		c.emit(location.EndPos.Line, bytecode.NIL)
	// 	}
	// 	if c.isNestedInFinally() {
	// 		c.emit(location.EndPos.Line, bytecode.RETURN_FINALLY)
	// 	} else {
	// 		c.emit(location.EndPos.Line, bytecode.RETURN)
	// 	}
	default:
		// TODO: implement finally
		// if c.isNestedInFinally() {
		// 	c.emit(location.EndPos.Line, bytecode.RETURN_FINALLY)
		// } else {
		c.emit("return %s, value.Undefined\n", val)
		// }
	}
}

func (c *GoCompiler) emitSetInstanceVariable(name value.Symbol, val string) {
	self := c.checker.SelfType()

	switch self := self.(type) {
	case types.NamespaceWithIvarIndices:
		index := self.IvarIndices().GetIndex(name)
		c.emitSetInstanceVariableByIndex(index, val)
	default:
		c.emitSetInstanceVariableByName(name, val)
	}
}

// Emit an instruction that sets the value of an instance variable by its index
func (c *GoCompiler) emitSetInstanceVariableByIndex(index int, val string) {
	c.emit("value.SetInstanceVariable(self, %d, %s)\n", index, val)
}

// Emit an instruction that sets the value of an instance variable by name
func (c *GoCompiler) emitSetInstanceVariableByName(name value.Symbol, val string) {
	c.registerGoLocal("err", c.checker.Std(symbol.Value), goValueType)
	symbol := c.emitSymbol(name.String())
	c.emit("err = value.SetInstanceVariableByName(self, %s, %s)\n", symbol, val)
	c.emit("if err.IsNotUndefined() { panic(err) }\n")
}

func (c *GoCompiler) compileMethodsWithinModule(module *types.Module, location *position.Location) {
	if types.NamespaceHasAnyDefinableMethods(module) {
		nameSymbol := c.emitSymbol(module.Name())
		c.emit("class = (*value.Module)(value.GetConstant(%s).Pointer()).SingletonClass()\n", nameSymbol)

		for methodName, method := range types.SortedOwnMethods(module) {
			c.compileMethodDefinition(methodName, method, location)

			for i, overload := range method.Overloads {
				overloadName := value.ToSymbol(
					fmt.Sprintf("%s@%d", methodName.String(), i+1),
				)
				c.compileMethodDefinition(overloadName, overload, location)
			}
		}
	}

	for _, subtype := range types.SortedSubtypes(module) {
		if subtype.Type == module {
			continue
		}
		c.compileMethodsWithinType(subtype.Type, location)
	}
}

func (c *GoCompiler) compileMethodsWithinNamespace(namespace types.Namespace, location *position.Location) {
	namespaceHasCompiledMethods := types.NamespaceHasAnyDefinableMethods(namespace)

	singleton := namespace.Singleton()
	singletonHasCompiledMethods := types.NamespaceHasAnyDefinableMethods(singleton)

	if namespaceHasCompiledMethods || singletonHasCompiledMethods {
		namespaceSymbol := c.emitSymbol(namespace.Name())
		c.emit("class = (*value.Class)(value.GetConstant(%s).Pointer())\n", namespaceSymbol)

		for methodName, method := range types.SortedOwnMethods(namespace) {
			c.compileMethodDefinition(methodName, method, location)
		}

		if singletonHasCompiledMethods {
			c.emit("class = class.SingletonClass()\n")

			for methodName, method := range types.SortedOwnMethods(singleton) {
				c.compileMethodDefinition(methodName, method, location)
			}
		}
	}

	for _, subtype := range types.SortedSubtypes(namespace) {
		if subtype.Type == namespace {
			continue
		}
		c.compileMethodsWithinType(subtype.Type, location)
	}
}

func (c *GoCompiler) compileMethodsWithinType(typ types.Type, location *position.Location) {
	switch t := typ.(type) {
	case *types.Module:
		c.compileMethodsWithinModule(t, location)
	case *types.Class:
		c.compileMethodsWithinNamespace(t, location)
	case *types.Mixin:
		c.compileMethodsWithinNamespace(t, location)
	case *types.Interface:
		c.compileMethodsWithinNamespace(t, location)
	}
}

func (c *GoCompiler) compileMethodDefinition(name value.Symbol, method *types.Method, location *position.Location) {
	if !method.IsDefinable() {
		return
	}

	if method.Base != nil {
		// handle aliases
		method = method.Base

		if method.IsNative() {
			namespace := value.RootModule.Constants.GetString(method.DefinedUnder.Name()).AsReference()
			c.registerGoLocal("aliasClass", c.checker.Std(symbol.Class), "*value.Class")

			namespaceSymbol := c.emitSymbol(method.DefinedUnder.Name())
			switch namespace.(type) {
			case *value.Class:
				c.emit("aliasClass = (*value.Class)(value.GetConstant(%s).Pointer())\n", namespaceSymbol)
			case *value.Module:
				c.emit("aliasClass = (*value.Module)(value.GetConstant(%s).Pointer()).SingletonClass()\n", namespaceSymbol)
			default:
				panic(fmt.Sprintf("invalid namespace %T", namespace))
			}

			oldNameSymbol := c.emitSymbol(method.Name.String())
			newNameSymbol := c.emitSymbol(name.String())
			c.emit("class.Methods[%s] = aliasClass.Methods[%s]\n", newNameSymbol, oldNameSymbol)

			method.SetCompiled(true)
			method.Body = nil
			return
		}
	}

	if method.IsAttribute() {
		if method.IsSetter() {
			nameStr := name.String()
			ivarName := value.ToSymbol(nameStr[:len(nameStr)-1])
			namespace := method.DefinedUnder

			var index int
			var ok bool

			switch n := namespace.(type) {
			case *types.Class:
				index, ok = n.IvarIndices().GetIndexOk(ivarName)
			case *types.SingletonClass:
				index, ok = n.IvarIndices().GetIndexOk(ivarName)
			case *types.Module:
				index, ok = n.IvarIndices().GetIndexOk(ivarName)
			default:
				index = -1
				ok = true
			}

			if !ok {
				panic(fmt.Sprintf("cannot get index of ivar `%s` in `%s`", ivarName.String(), namespace.Name()))
			}

			ivarNameSymbol := c.emitSymbol(ivarName.String())
			c.emit("vm.DefineSetter(&class.MethodContainer, %s, %d)\n", ivarNameSymbol, index)

			method.SetCompiled(true)
			method.Body = nil
			return
		}

		namespace := method.DefinedUnder

		var index int
		var ok bool

		switch n := namespace.(type) {
		case *types.Class:
			index, ok = n.IvarIndices().GetIndexOk(name)
		case *types.SingletonClass:
			index, ok = n.IvarIndices().GetIndexOk(name)
		case *types.Module:
			index, ok = n.IvarIndices().GetIndexOk(name)
		default:
			index = -1
			ok = true
		}

		if !ok {
			panic(fmt.Sprintf("cannot get index of ivar `%s` in `%s`", name.String(), namespace.Name()))
		}

		nameSymbol := c.emitSymbol(name.String())
		c.emit("vm.DefineGetter(&class.MethodContainer, %s, %d)\n", nameSymbol, index)

		method.SetCompiled(true)
		method.Body = nil
		return
	}

	c.emit("vm.Def(&class.MethodContainer, %s,\n", name.String())

	methodCompiler := (*GoCompiler)(method.Body.(*GoSourceMethod))
	c.emitBytes(methodCompiler.buff.Bytes())
	c.emitPackageBytes(methodCompiler.packageBuff.Bytes())
	methodCompiler.buff.Reset()
	methodCompiler.packageBuff.Reset()

	c.emit(")\n")

	method.SetCompiled(true)
	method.Body = nil
}

func (c *GoCompiler) compileNamespaceDefinition(parentNamespace, namespace types.Namespace, constName value.Symbol) {
	if !namespace.IsDefined() && !namespace.IsNative() {
		switch p := parentNamespace.(type) {
		case *types.SingletonClass:
			parentSymbol := c.emitSymbol(p.AttachedObject.Name())
			c.emit("parentNamespace = value.GetConstant(%s).SingletonClass()\n", parentSymbol)
		default:
			parentSymbol := c.emitSymbol(p.Name())
			c.emit("parentNamespace = value.GetConstant(%s)\n", parentSymbol)
		}

		switch namespace.(type) {
		case *types.Module:
			c.emit("namespace = value.Ref(value.NewModule())\n")
		case *types.Class:
			c.emit("namespace = value.Ref(value.NewClassWithOptions(value.ClassWithSuperclass(nil)))\n")
		case *types.Mixin:
			c.emit("namespace = value.Ref(value.NewMixin())\n")
		case *types.Interface:
			c.emit("namespace = value.Ref(value.NewInterface())\n")
		}
		constNameSymbol := c.emitSymbol(constName.String())
		c.emit("value.AddConstant(parentNamespace, %s, namespace)\n\n", constNameSymbol)
		namespace.SetDefined(true)
	}

	for name, subtype := range types.SortedSubtypes(namespace) {
		if subtype.Type == namespace {
			continue
		}
		c.compileSubtypeDefinition(namespace, subtype.Type, name)
	}
}

func (c *GoCompiler) compileSubtypeDefinition(parentNamespace types.Namespace, typ types.Type, constName value.Symbol) {
	switch t := typ.(type) {
	case *types.Module:
		c.compileModuleDefinition(parentNamespace, t, constName)
	case *types.Class:
		c.compileClassDefinition(parentNamespace, t, constName)
	case *types.Mixin:
		c.compileMixinDefinition(parentNamespace, t, constName)
	case *types.Interface:
		c.compileInterfaceDefinition(parentNamespace, t, constName)
	}
}

// Prepend before the main block of source code in a func
func (c *GoCompiler) emitPrepend(format string, a ...any) {
	prevBuff := c.buff.Bytes()
	c.buff = bytes.Buffer{}
	c.emit(format, a...)
	c.emitBytes(prevBuff)
}

// Prepend before the main block of source code in a func
func (c *GoCompiler) emitPrependBytes(byt []byte) {
	prevBuff := c.buff.Bytes()
	c.buff = bytes.Buffer{}
	c.emitBytes(byt)
	c.emitBytes(prevBuff)
}

// Emit code inside of a func
func (c *GoCompiler) emit(format string, a ...any) {
	fmt.Fprintf(&c.buff, format, a...)
}

// Emit code inside of a func
func (c *GoCompiler) emitBytes(byt []byte) {
	_, err := c.buff.Write(byt)
	if err != nil {
		panic(fmt.Sprintf("cannot emit bytes: %s", err))
	}
}

// Emit package level code
func (c *GoCompiler) emitPackage(format string, a ...any) {
	fmt.Fprintf(&c.packageBuff, format, a...)
}

// Emit package level code
func (c *GoCompiler) emitPackageBytes(byt []byte) {
	_, err := c.packageBuff.Write(byt)
	if err != nil {
		panic(fmt.Sprintf("cannot emit bytes: %s", err))
	}
}

func (c *GoCompiler) typeOf(node ast.Node) types.Type {
	return node.Type(c.checker.Env())
}

func (c *GoCompiler) compileModuleDefinition(parentNamespace types.Namespace, module *types.Module, constName value.Symbol) {
	c.compileNamespaceDefinition(parentNamespace, module, constName)
}

func (c *GoCompiler) compileClassDefinition(parentNamespace types.Namespace, class *types.Class, constName value.Symbol) {
	c.compileNamespaceDefinition(parentNamespace, class, constName)
}

func (c *GoCompiler) compileMixinDefinition(parentNamespace types.Namespace, mixin *types.Mixin, constName value.Symbol) {
	c.compileNamespaceDefinition(parentNamespace, mixin, constName)
}

func (c *GoCompiler) compileInterfaceDefinition(parentNamespace types.Namespace, iface *types.Interface, constName value.Symbol) {
	c.compileNamespaceDefinition(parentNamespace, iface, constName)
}

func (c *GoCompiler) emitGetConstValue(name value.Symbol) *goValue {
	return c.emitGetConst(name, c.checker.Std(symbol.Value))
}

func (c *GoCompiler) emitGetConst(name value.Symbol, typ types.Type) *goValue {
	tmp := c.getTmpIdent()
	c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
	constNameSymbol := c.emitSymbol(name.String())
	c.emit("%s = value.GetConstant(%s)\n", tmp, constNameSymbol)

	return newTmpValue(
		tmp,
		typ,
	)
}

func (c *GoCompiler) CompileClassInheritance(class *types.Class, location *position.Location) {
	if class.IsCompiled() {
		return
	}
	superclass := class.Superclass()
	if superclass == nil {
		return
	}

	class.SetCompiled(true)

	classNameSymbol := c.emitSymbol(class.Name())
	superclassNameSymbol := c.emitSymbol(superclass.Name())
	c.emit("class = (*value.Class)(value.GetConstant(%s).Pointer())", classNameSymbol)
	c.emit("superclass = (*value.Class)(value.GetConstant(%s).Pointer())", superclassNameSymbol)
	c.emit("class.SetSuperclass(superclass)")
}

func (c *GoCompiler) CompileIvarIndices(target types.NamespaceWithIvarIndices, location *position.Location) {
	switch target := target.(type) {
	case *types.SingletonClass:
		targetSymbol := c.emitSymbol(target.AttachedObject.Name())
		c.emit("class = value.GetConstant(%s).SingletonClass()\n", targetSymbol)
	case *types.Module:
		targetSymbol := c.emitSymbol(target.Name())
		c.emit("class = value.GetConstant(%s).SingletonClass()\n", targetSymbol)
	default:
		targetSymbol := c.emitSymbol(target.Name())
		c.emit("class = (*value.Class)(value.GetConstant(%s).Pointer())\n", targetSymbol)
	}

	c.emit("class.IvarIndices = %s\n", target.IvarIndices().ToGoSource())
}

func (c *GoCompiler) CompileInclude(target types.Namespace, mixin *types.Mixin, location *position.Location) {
	switch t := target.(type) {
	case *types.SingletonClass:
		targetSymbol := c.emitSymbol(t.AttachedObject.Name())
		c.emit("class = value.GetConstant(%s).SingletonClass()\n", targetSymbol)
	default:
		targetSymbol := c.emitSymbol(target.Name())
		c.emit("class = (*value.Class)(value.GetConstant(%s).Pointer())\n", targetSymbol)
	}

	mixinNameSymbol := c.emitSymbol(mixin.Name())
	c.emit("mixin = (*value.Mixin)(value.GetConstant(%s).Pointer())\n", mixinNameSymbol)
	c.emit("class.IncludeMixin(mixin)\n")
}

func mangleFileName(name string) string {
	var b strings.Builder

	b.WriteString("__file_")

	for i, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			// Identifiers cannot start with a digit
			if i == 0 && unicode.IsDigit(r) {
				b.WriteByte('_')
			}
			b.WriteRune(r)
			continue
		}

		// Everything else becomes underscore
		b.WriteByte('_')
	}

	return b.String()
}

func (c *GoCompiler) InitExpressionCompiler(location *position.Location) Compiler {
	name := mangleFileName(location.FilePath)
	exprCompiler := NewGoCompiler(name, topLevelGoCompilerMode, location, c.checker, c.bigIntCache, c.symbolCache, c.output)
	exprCompiler.Errors = c.Errors

	c.emit("%s()\n", name)

	return exprCompiler
}

func (c *GoCompiler) CompileExpressionsInFile(node *ast.ProgramNode) {
	c.compileProgram(node)

	var funcBuffer bytes.Buffer
	fmt.Fprintf(&funcBuffer, "func %s() {\n", c.FuncName)
	c.compileLocalsTo(&funcBuffer)
	c.emitPrependBytes(funcBuffer.Bytes())
	c.emit("}\n")
}

func (c *GoCompiler) compileLocalsTo(buff io.Writer) {
	for _, local := range c.goLocals {
		fmt.Fprintf(buff, "var %s %s\n", local.name, local.goType)
	}
	buff.Write([]byte("\n"))
}

func (c *GoCompiler) compileLocals() {
	c.compileLocalsTo(&c.buff)
}

// Entry point to the compilation process
func (c *GoCompiler) compileProgram(node *ast.ProgramNode) *goValue {
	return c.compileStatements(node.Body)
}

func (c *GoCompiler) compileStatements(nodes []ast.StatementNode) *goValue {
	var lastValue *goValue
	for _, stmt := range nodes {
		lastValue = c.compileStatement(stmt)
	}

	if lastValue == nil {
		return newInlineValue("value.Nil", c.checker.Std(symbol.Nil))
	}
	return lastValue
}

func (c *GoCompiler) compileStatement(node ast.StatementNode) *goValue {
	switch node := node.(type) {
	case *ast.ExpressionStatementNode:
		return c.compileExpression(node.Expression)
	default:
		return newInlineValue("value.Nil", c.checker.Std(symbol.Nil))
	}
}

func (c *GoCompiler) compileExpression(node ast.ExpressionNode) *goValue {
	switch node := node.(type) {
	case *ast.IntLiteralNode:
		i, err := value.ParseBigInt(node.Value, 0)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		if i.IsSmallInt() {
			return newInlineValue(
				fmt.Sprintf("value.SmallInt(%d)", i.ToSmallInt()),
				c.typeOf(node),
			)
		}
		bigIntVar := c.emitBigInt(node.Value)
		return newTmpValue(
			bigIntVar,
			c.typeOf(node),
		)
	case *ast.Int8LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 8)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineValue(
			fmt.Sprintf("value.Int8(%d)", i),
			c.typeOf(node),
		)
	case *ast.Int16LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 16)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineValue(
			fmt.Sprintf("value.Int16(%d)", i),
			c.typeOf(node),
		)
	case *ast.Int32LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 32)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineValue(
			fmt.Sprintf("value.Int32(%d)", i),
			c.typeOf(node),
		)
	case *ast.Int64LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 64)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineValue(
			fmt.Sprintf("value.Int64(%d)", i),
			c.typeOf(node),
		)
	case *ast.UInt8LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 8)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineValue(
			fmt.Sprintf("value.UInt8(%d)", i),
			c.typeOf(node),
		)
	case *ast.UInt16LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 16)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineValue(
			fmt.Sprintf("value.UInt16(%d)", i),
			c.typeOf(node),
		)
	case *ast.UInt32LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 32)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineValue(
			fmt.Sprintf("value.UInt32(%d)", i),
			c.typeOf(node),
		)
	case *ast.UInt64LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 64)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineValue(
			fmt.Sprintf("value.UInt64(%d)", i),
			c.typeOf(node),
		)
	case *ast.BinaryExpressionNode:
		return c.compileBinaryExpressionNode(node)
	default:
		panic(fmt.Sprintf("invalid expression node: %T", node))
	}
}

func (c *GoCompiler) emitBigInt(val string) string {
	c.bigIntCache.Lock()
	defer c.bigIntCache.Unlock()

	bigInt, ok := c.bigIntCache.GetUnsafe(val)
	if ok {
		return bigInt.goIdent()
	}

	bigInt = &nativeBigInt{
		id:  c.bigIntCache.Len(),
		val: val,
	}
	c.bigIntCache.SetUnsafe(val, bigInt)
	ident := bigInt.goIdent()
	c.emitPackage("var %s = value.MustParseInt(%q, 0)\n", ident, val)
	return ident
}

func (c *GoCompiler) emitSymbol(val string) string {
	c.symbolCache.Lock()
	defer c.symbolCache.Unlock()

	symbol, ok := c.symbolCache.GetUnsafe(val)
	if ok {
		return symbol.goIdent()
	}

	symbol = &nativeSymbol{
		id:  c.symbolCache.Len(),
		val: val,
	}
	c.symbolCache.SetUnsafe(val, symbol)
	ident := symbol.goIdent()
	c.emitPackage("var %s = value.ToSymbol(%q)\n", ident, val)
	return ident
}

func (c *GoCompiler) getTmpIdent() string {
	c.tmpLocalCounter++
	return fmt.Sprintf("t%d", c.tmpLocalCounter)
}

func (c *GoCompiler) emitErrorPropagation() {
	switch c.mode {
	case topLevelGoCompilerMode:
		c.emit("if err.IsNotUndefined() { thread.PrintErrorValue(err); os.Exit(1) }\n")
	default:
		c.emit("if err.IsNotUndefined() { return value.Undefined, err }\n")
	}
}

func (c *GoCompiler) emitCallCache() string {
	c.callCacheCounter++
	callCacheName := fmt.Sprintf("cc_%s_%d", c.FuncName, c.callCacheCounter)
	c.emitPackage("var %s *value.CallCache", callCacheName)

	return callCacheName
}

func (c *GoCompiler) emitAddCallFrame(loc *position.Location) {
	c.emit(
		"thread.AddCallFrame(value.CallFrame{FuncName: %q, FileName: %q, LineNumber: %d})\n",
		c.FuncName,
		loc.FilePath,
		loc.StartPos.Line,
	)
}
func (c *GoCompiler) emitPopCallFrame() {
	c.emit("thread.PopCallFrame()\n")
}

func (c *GoCompiler) registerErr() {
	c.registerGoLocal("err", types.Any{}, goValueType)
}

func (c *GoCompiler) compileBinaryExpressionNode(node *ast.BinaryExpressionNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	left := c.compileExpression(node.Left)
	right := c.compileExpression(node.Right)
	typ := left.typ

	switch node.Op.Type {
	case token.PLUS:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = value.AddInt(%s, %s)\n", tmp, left.Value(), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = %s.AddVal(%s)\n", tmp, c.goValueAs(left, "AsFloat"), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.StdBigFloat()) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = %s.AddVal(%s)\n", tmp, c.goValueCast(left, "*value.BigFloat"), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinAddable)) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = value.AddVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}

		callCache := c.emitCallCache()
		tmp := c.getTmpIdent()
		c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
		c.registerErr()
		c.emitAddCallFrame(node.Location())
		c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpAdd, %s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
		c.emitPopCallFrame()
		c.emitErrorPropagation()
		return newTmpValue(
			tmp,
			c.checker.Std(symbol.Value),
		)
	case token.MINUS:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = value.SubtractInt(%s, %s)\n", tmp, left.Value(), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = %s.SubtractVal(%s)\n", tmp, c.goValueAs(left, "AsFloat"), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.StdBigFloat()) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = %s.SubtractVal(%s)\n", tmp, c.goValueCast(left, "*value.BigFloat"), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinAddable)) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = value.SubtractVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}

		callCache := c.emitCallCache()
		tmp := c.getTmpIdent()
		c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
		c.registerErr()
		c.emitAddCallFrame(node.Location())
		c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpSubtract, %s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
		c.emitPopCallFrame()
		c.emitErrorPropagation()
		return newTmpValue(
			tmp,
			c.checker.Std(symbol.Value),
		)
	case token.STAR:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = value.MultiplyInt(%s, %s)\n", tmp, left.Value(), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = %s.MultiplyVal(%s)\n", tmp, c.goValueAs(left, "AsFloat"), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.StdBigFloat()) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = %s.MultiplyVal(%s)\n", tmp, c.goValueCast(left, "*value.BigFloat"), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinAddable)) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = value.MultiplyVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}

		callCache := c.emitCallCache()
		tmp := c.getTmpIdent()
		c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
		c.registerErr()
		c.emitAddCallFrame(node.Location())
		c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpMultiply, %s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
		c.emitPopCallFrame()
		c.emitErrorPropagation()
		return newTmpValue(
			tmp,
			c.checker.Std(symbol.Value),
		)
	case token.SLASH:
		c.emitAddCallFrame(node.Location())

		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = value.DivideInt(%s, %s)\n", tmp, left.Value(), c.convertToValue(right))
			c.emitPopCallFrame()
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = %s.DivideVal(%s)\n", tmp, c.goValueAs(left, "AsFloat"), c.convertToValue(right))
			c.emitPopCallFrame()
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.StdBigFloat()) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = %s.DivideVal(%s)\n", tmp, c.goValueCast(left, "*value.BigFloat"), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinAddable)) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = value.DivideVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
			c.emitPopCallFrame()
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}

		callCache := c.emitCallCache()
		tmp := c.getTmpIdent()
		c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
		c.registerErr()
		c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpMultiply, %s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
		c.emitPopCallFrame()
		c.emitErrorPropagation()
		return newTmpValue(
			tmp,
			c.checker.Std(symbol.Value),
		)
	case token.STAR_STAR:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = value.ExponentiateInt(%s, %s)\n", tmp, left.Value(), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = %s.ExponentiateVal(%s)\n", tmp, c.goValueAs(left, "AsFloat"), c.convertToValue(right))
			c.emitPopCallFrame()
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinNumeric)) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = value.ExponentiateVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
			c.emitPopCallFrame()
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}

		callCache := c.emitCallCache()
		tmp := c.getTmpIdent()
		c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
		c.registerErr()
		c.emitAddCallFrame(node.Location())
		c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpMultiply, %s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
		c.emitPopCallFrame()
		c.emitErrorPropagation()
		return newTmpValue(
			tmp,
			c.checker.Std(symbol.Value),
		)
	case token.LBITSHIFT:
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = value.LeftBitshiftVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
			c.emitPopCallFrame()
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}

		callCache := c.emitCallCache()
		tmp := c.getTmpIdent()
		c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
		c.registerErr()
		c.emitAddCallFrame(node.Location())
		c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpLeftBitshift, %s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
		c.emitPopCallFrame()
		c.emitErrorPropagation()
		return newTmpValue(
			tmp,
			c.checker.Std(symbol.Value),
		)
	case token.LTRIPLE_BITSHIFT:
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinLogicBitshiftable)) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = value.LogicalLeftBitshiftVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
			c.emitPopCallFrame()
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}

		callCache := c.emitCallCache()
		tmp := c.getTmpIdent()
		c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
		c.registerErr()
		c.emitAddCallFrame(node.Location())
		c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpLogicalLeftBitshift, %s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
		c.emitPopCallFrame()
		c.emitErrorPropagation()
		return newTmpValue(
			tmp,
			c.checker.Std(symbol.Value),
		)
	case token.RBITSHIFT:
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = value.RightBitshiftVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
			c.emitPopCallFrame()
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}

		callCache := c.emitCallCache()
		tmp := c.getTmpIdent()
		c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
		c.registerErr()
		c.emitAddCallFrame(node.Location())
		c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpRightBitshift, %s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
		c.emitPopCallFrame()
		c.emitErrorPropagation()
		return newTmpValue(
			tmp,
			c.checker.Std(symbol.Value),
		)
	case token.RTRIPLE_BITSHIFT:
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinLogicBitshiftable)) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = value.LogicalRightBitshiftVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
			c.emitPopCallFrame()
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}

		callCache := c.emitCallCache()
		tmp := c.getTmpIdent()
		c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
		c.registerErr()
		c.emitAddCallFrame(node.Location())
		c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpLogicalRightBitshift, %s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
		c.emitPopCallFrame()
		c.emitErrorPropagation()
		return newTmpValue(
			tmp,
			c.checker.Std(symbol.Value),
		)
	case token.AND:
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
			tmp := c.getTmpIdent()
			c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
			c.registerErr()
			c.emit("%s, err = value.BitwiseAndVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
			c.emitPopCallFrame()
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}

		callCache := c.emitCallCache()
		tmp := c.getTmpIdent()
		c.registerGoLocal(tmp, c.checker.Std(symbol.Value), goValueType)
		c.registerErr()
		c.emitAddCallFrame(node.Location())
		c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpAnd, %s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
		c.emitPopCallFrame()
		c.emitErrorPropagation()
		return newTmpValue(
			tmp,
			c.checker.Std(symbol.Value),
		)
	// case token.AND_TILDE:
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
	// 		c.emit(line, bytecode.BITWISE_AND_NOT)
	// 		return
	// 	}
	// 	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpAndNot, 1), location, false)
	// case token.OR:
	// 	if c.checker.IsSubtype(typ, c.checker.StdInt()) {
	// 		c.emit(line, bytecode.BITWISE_OR_INT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
	// 		c.emit(line, bytecode.BITWISE_OR)
	// 		return
	// 	}
	// 	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpOr, 1), location, false)
	// case token.XOR:
	// 	if c.checker.IsSubtype(typ, c.checker.StdInt()) {
	// 		c.emit(line, bytecode.BITWISE_XOR_INT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
	// 		c.emit(line, bytecode.BITWISE_XOR)
	// 		return
	// 	}
	// 	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpXor, 1), location, false)
	// case token.PERCENT:
	// 	if c.checker.IsSubtype(typ, c.checker.StdInt()) {
	// 		c.emit(line, bytecode.MODULO_INT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
	// 		c.emit(line, bytecode.MODULO_FLOAT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinNumeric)) {
	// 		c.emit(line, bytecode.MODULO)
	// 		return
	// 	}
	// 	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpModulo, 1), location, false)
	// case token.LAX_EQUAL:
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinEquatable)) {
	// 		c.emit(line, bytecode.LAX_EQUAL)
	// 		return
	// 	}
	// 	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpLaxEqual, 1), location, false)
	// case token.LAX_NOT_EQUAL:
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinEquatable)) {
	// 		c.emit(line, bytecode.LAX_NOT_EQUAL)
	// 		return
	// 	}
	// 	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpLaxEqual, 1), location, false)
	// 	c.emit(line, bytecode.NOT)
	// case token.EQUAL_EQUAL:
	// 	if c.checker.IsSubtype(typ, c.checker.StdInt()) {
	// 		c.emit(line, bytecode.EQUAL_INT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
	// 		c.emit(line, bytecode.EQUAL_INT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinEquatable)) {
	// 		c.emit(line, bytecode.EQUAL)
	// 		return
	// 	}
	// 	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpEqual, 1), location, false)
	// case token.NOT_EQUAL:
	// 	if c.checker.IsSubtype(typ, c.checker.StdInt()) {
	// 		c.emit(line, bytecode.NOT_EQUAL_INT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.StdInt()) {
	// 		c.emit(line, bytecode.NOT_EQUAL_FLOAT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinEquatable)) {
	// 		c.emit(line, bytecode.NOT_EQUAL)
	// 		return
	// 	}
	// 	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpEqual, 1), location, false)
	// 	c.emit(line, bytecode.NOT)
	// case token.STRICT_EQUAL:
	// 	c.emit(line, bytecode.STRICT_EQUAL)
	// case token.STRICT_NOT_EQUAL:
	// 	c.emit(line, bytecode.STRICT_NOT_EQUAL)
	// case token.GREATER:
	// 	if c.checker.IsSubtype(typ, c.checker.StdInt()) {
	// 		c.emit(line, bytecode.GREATER_INT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
	// 		c.emit(line, bytecode.GREATER_FLOAT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinNumeric)) {
	// 		c.emit(line, bytecode.GREATER)
	// 		return
	// 	}
	// 	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpGreaterThan, 1), location, false)
	// case token.GREATER_EQUAL:
	// 	if c.checker.IsSubtype(typ, c.checker.StdInt()) {
	// 		c.emit(line, bytecode.GREATER_EQUAL_I)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
	// 		c.emit(line, bytecode.GREATER_EQUAL_F)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinNumeric)) {
	// 		c.emit(line, bytecode.GREATER_EQUAL)
	// 		return
	// 	}
	// 	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpGreaterThanEqual, 1), location, false)
	// case token.LESS:
	// 	if c.checker.IsSubtype(typ, c.checker.StdInt()) {
	// 		c.emit(line, bytecode.LESS_INT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
	// 		c.emit(line, bytecode.LESS_FLOAT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinNumeric)) {
	// 		c.emit(line, bytecode.LESS)
	// 		return
	// 	}
	// 	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpLessThan, 1), location, false)
	// case token.LESS_EQUAL:
	// 	if c.checker.IsSubtype(typ, c.checker.StdInt()) {
	// 		c.emit(line, bytecode.LESS_EQUAL_INT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
	// 		c.emit(line, bytecode.LESS_EQUAL_FLOAT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinNumeric)) {
	// 		c.emit(line, bytecode.LESS_EQUAL)
	// 		return
	// 	}
	// 	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpLessThanEqual, 1), location, false)
	// case token.SPACESHIP_OP:
	// 	c.emit(line, bytecode.COMPARE)
	// case token.INSTANCE_OF_OP:
	// 	c.emit(line, bytecode.INSTANCE_OF)
	// case token.REVERSE_INSTANCE_OF_OP:
	// 	c.emit(line, bytecode.INSTANCE_OF)
	// 	c.emit(line, bytecode.NOT)
	// case token.ISA_OP:
	// 	c.emit(line, bytecode.IS_A)
	// case token.REVERSE_ISA_OP:
	// 	c.emit(line, bytecode.IS_A)
	// 	c.emit(line, bytecode.NOT)
	default:
		c.Errors.AddFailure(fmt.Sprintf("unknown binary operator: %s", node.Op.String()), node.Location())
		return nil
	}
}

func (c *GoCompiler) resolve(node ast.ExpressionNode) *goValue {
	result := resolve(node)
	if result.IsUndefined() {
		return nil
	}

	return c.valueToGoSource(result)
}

func (c *GoCompiler) valueToGoSource(val value.Value) *goValue {
	if val.IsReference() {
		switch v := val.AsReference().(type) {
		case *value.ArrayList:
			return c.arrayListToGoSource(v)
		case *value.ArrayTuple:
			return c.arrayTupleToGoSource(v)
		case *value.HashSet:
			return c.hashSetToGoSource(v)
		case *value.HashMap:
			return c.hashMapToGoSource(v)
		case *value.HashRecord:
			return c.hashRecordToGoSource(v)
		case value.Int64:
			return newInlineValue(
				fmt.Sprintf("value.Int64(%d)", v),
				c.checker.Std(symbol.Int64),
			)
		case value.UInt64:
			return newInlineValue(
				fmt.Sprintf("value.UInt64(%d)", v),
				c.checker.Std(symbol.UInt64),
			)
		default:
			return nil
		}
	}

	switch val.ValueFlag() {
	case value.TRUE_FLAG:
		return newInlineValue(
			"value.True",
			c.checker.StdValue(),
		)
	case value.FALSE_FLAG:
		return newInlineValue(
			"value.False",
			c.checker.StdValue(),
		)
	case value.NIL_FLAG:
		return newInlineValue(
			"value.Nil",
			c.checker.StdValue(),
		)
	case value.SMALL_INT_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.SmallInt(%d)", val.AsSmallInt()),
			c.checker.StdValue(),
		)
	case value.INT64_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.Int64(%d)", val.AsInt64()),
			c.checker.Std(symbol.Int64),
		)
	case value.UINT64_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.UInt64(%d)", val.AsUInt64()),
			c.checker.Std(symbol.UInt64),
		)
	case value.INT32_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.Int32(%d)", val.AsInt32()),
			c.checker.Std(symbol.Int32),
		)
	case value.UINT32_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.UInt32(%d)", val.AsUInt32()),
			c.checker.Std(symbol.UInt32),
		)
	case value.INT16_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.Int16(%d)", val.AsInt16()),
			c.checker.Std(symbol.Int16),
		)
	case value.UINT16_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.UInt16(%d)", val.AsUInt16()),
			c.checker.Std(symbol.UInt16),
		)
	case value.INT8_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.Int8(%d)", val.AsInt8()),
			c.checker.Std(symbol.Int8),
		)
	case value.UINT8_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.UInt8(%d)", val.AsUInt8()),
			c.checker.Std(symbol.UInt8),
		)
	case value.CHAR_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.Char(%q)", val.AsChar()),
			c.checker.Std(symbol.Char),
		)
	case value.FLOAT_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.Float(%f)", val.AsFloat()),
			c.checker.Std(symbol.Float),
		)
	}

	return nil
}

func (c *GoCompiler) arrayListToGoSource(v *value.ArrayList) *goValue {
	var buff strings.Builder

	fmt.Fprintf(&buff, "value.NewArrayListWithElements(%d, ", v.LeftCapacity())

	for _, element := range *v {
		el := c.valueToGoSource(element)
		if el == nil {
			return nil
		}

		fmt.Fprintf(&buff, "%s, ", c.convertToValue(el))
	}

	buff.WriteString(")")
	return newInlineValue(buff.String(), c.checker.Std(symbol.ArrayList))
}

func (c *GoCompiler) goValueAs(v *goValue, as string) string {
	if v.IsTmp() {
		return fmt.Sprintf("%s.%s()", v.tmpLocal, as)
	}
	if c.checker.IsTheSameType(v.typ, c.checker.Std(symbol.Value)) {
		return fmt.Sprintf("%s.%s()", v.inline, as)
	}

	return v.inline
}

func (c *GoCompiler) goValueCast(v *goValue, cast string) string {
	if v.IsTmp() {
		return fmt.Sprintf("(%s)(%s.Pointer())", cast, v.tmpLocal)
	}
	if c.checker.IsTheSameType(v.typ, c.checker.Std(symbol.Value)) {
		return fmt.Sprintf("(%s)(%s.Pointer())", cast, v.inline)
	}

	return v.inline
}

func (c *GoCompiler) convertToValue(v *goValue) string {
	if v.IsTmp() {
		return v.tmpLocal
	}

	if c.checker.IsTheSameType(v.typ, c.checker.Std(symbol.Value)) {
		return v.inline
	}
	if c.checker.IsSubtype(v.typ, c.checker.Std(symbol.Object)) {
		return fmt.Sprintf("value.Ref(%s)", v.inline)
	}

	return fmt.Sprintf("%s.ToValue()", v.inline)
}

func (c *GoCompiler) arrayTupleToGoSource(v *value.ArrayTuple) *goValue {
	var buff strings.Builder

	buff.WriteString("value.NewArrayTupleWithElements(0, ")

	for _, element := range *v {
		el := c.valueToGoSource(element)
		if el == nil {
			return nil
		}

		fmt.Fprintf(&buff, "%s, ", c.convertToValue(el))
	}

	buff.WriteString(")")
	return newInlineValue(buff.String(), c.checker.Std(symbol.ArrayTuple))
}

func (c *GoCompiler) hashSetToGoSource(v *value.HashSet) *goValue {
	var buff strings.Builder

	fmt.Fprintf(&buff, "vm.MustNewHashSetWithCapacityAndElements(nil, %d, ", v.LeftCapacity())

	for element := range v.All() {
		el := c.valueToGoSource(element)
		if el == nil {
			return nil
		}

		fmt.Fprintf(&buff, "%s, ", c.convertToValue(el))
	}

	buff.WriteString(")")
	return newInlineValue(buff.String(), c.checker.Std(symbol.HashSet))
}

func (c *GoCompiler) hashMapToGoSource(v *value.HashMap) *goValue {
	var buff strings.Builder

	fmt.Fprintf(&buff, "vm.MustNewHashMapWithCapacityAndElements(nil, %d, ", v.LeftCapacity())

	for pair := range v.All() {
		p := c.valuePairToGoSource(pair)
		if p == nil {
			return nil
		}

		fmt.Fprintf(&buff, "%s, ", p)
	}

	buff.WriteString(")")
	return newInlineValue(buff.String(), c.checker.Std(symbol.HashMap))
}

func (c *GoCompiler) hashRecordToGoSource(v *value.HashRecord) *goValue {
	var buff strings.Builder

	buff.WriteString("vm.MustNewHashRecordWithElements(nil, ")

	for pair := range v.All() {
		p := c.valuePairToGoSource(pair)
		if p == nil {
			return nil
		}

		fmt.Fprintf(&buff, "%s, ", p)
	}

	buff.WriteString(")")
	return newInlineValue(buff.String(), c.checker.Std(symbol.HashRecord))
}

func (c *GoCompiler) valuePairToGoSource(p value.Pair) *goValue {
	k := c.valueToGoSource(p.Key)
	if k == nil {
		return nil
	}
	v := c.valueToGoSource(p.Key)
	if v == nil {
		return nil
	}

	return newInlineValue(
		fmt.Sprintf("Pair{Key: %s, Value: %s}", k.inline, v.inline),
		c.checker.Std(symbol.Pair),
	)
}

func (c *GoCompiler) enterScope(label string, typ nativeElkScopeType) {
	c.scopes = append(c.scopes, newNativeElkScope(label, typ))
}

func (c *GoCompiler) leaveScope(line int) {
	currentDepth := len(c.scopes) - 1
	c.scopes[currentDepth] = nil
	c.scopes = c.scopes[:currentDepth]
}

// Register a local variable.
func (c *GoCompiler) defineLocal(name string, typ types.Type, location *position.Location) *nativeElkLocal {
	varScope := c.scopes.last()
	_, ok := varScope.localTable[name]
	if ok {
		c.Errors.AddFailure(
			fmt.Sprintf("a variable with this name has already been declared in this scope `%s`", name),
			location,
		)
		return nil
	}
	return c.defineVariableInScope(varScope, name, typ, location)
}

func (c *GoCompiler) defineVariableInScope(scope *nativeElkScope, name string, typ types.Type, location *position.Location) *nativeElkLocal {
	if c.lastElkLocalIndex == math.MaxInt {
		c.Errors.AddFailure(
			fmt.Sprintf("exceeded the maximum number of local variables (%d): %s", math.MaxInt, name),
			location,
		)
		return nil
	}

	c.lastElkLocalIndex++
	c.registerGoLocal(name, typ, goValueType)

	newVar := &nativeElkLocal{
		index: c.lastElkLocalIndex,
		name:  name,
		typ:   typ,
	}
	scope.localTable[name] = newVar
	return newVar
}
