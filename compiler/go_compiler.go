package compiler

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strings"
	"unicode"

	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/ds"
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
	closureGoCompilerMode
)

// represents a nativeElkLocal variable or value
type nativeElkLocal struct {
	name    string
	elkType types.Type
	goLocal *goLocal
}

func (n *nativeElkLocal) goIdent() string {
	return n.goLocal.name
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
	goType  string
	comment string
	free    bool
}

func (l *goLocal) markFree() {
	l.free = true
}

func (l *goLocal) markReserved() {
	l.free = false
}

func newGoLocal(name string, goType string, comment string) *goLocal {
	return &goLocal{
		name:    name,
		goType:  goType,
		comment: comment,
	}
}

type goValue struct {
	tmpLocal     *goLocal
	inline       string
	inlineGoType string
	elkType      types.Type
}

func (v *goValue) markFree() {
	if v.isTmp() {
		v.tmpLocal.markFree()
	}
}

func (v *goValue) goType() string {
	if v.isInline() {
		return v.inlineGoType
	}
	return v.tmpLocal.goType
}

func (v *goValue) isInline() bool {
	return v.tmpLocal == nil
}

func (v *goValue) isTmp() bool {
	return v.tmpLocal != nil
}

func (v *goValue) value() string {
	if v.inline != "" {
		return v.inline
	}

	return v.tmpLocal.name
}

func newInlineGoValue(v string, typ types.Type, goType string) *goValue {
	return &goValue{
		inline:       v,
		elkType:      typ,
		inlineGoType: goType,
	}
}

func newTmpGoValue(goLocal *goLocal, typ types.Type) *goValue {
	return &goValue{
		tmpLocal: goLocal,
		elkType:  typ,
	}
}

var errGoValue = newInlineGoValue("ERR", types.Untyped{}, goValueType)
var nilGoValue = newInlineGoValue("value.Nil", types.Nil{}, goValueType)

func CreateGoCompiler(parent *GoCompiler, checker types.Checker, loc *position.Location, errors *diagnostic.SyncDiagnosticList, output io.Writer) *GoCompiler {
	bigIntCache := concurrent.NewMap[string, *nativeBigInt]()
	symbolCache := concurrent.NewMap[string, *nativeSymbol]()
	compiler := NewGoCompiler("main", topLevelGoCompilerMode, loc, checker, bigIntCache, symbolCache, output)
	compiler.Errors = errors
	if parent != nil {
		compiler.SetParent(parent)
	}
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
	c.emitPackage("import \"github.com/elk-language/elk/value\"\n")
	c.emitPackage("import \"github.com/elk-language/elk/vm\"\n\n")
}

func (c *GoCompiler) InitGlobalEnv() Compiler {
	envCompiler := NewGoCompiler("initGlobalEnv", topLevelGoCompilerMode, c.loc, c.checker, c.bigIntCache, c.symbolCache, c.output)
	envCompiler.SetParent(c)
	envCompiler.Errors = c.Errors
	envCompiler.compileGlobalEnv()

	return envCompiler
}

func (c *GoCompiler) FinishGlobalEnvCompiler() {
	if c.buff.Len() == 0 {
		return
	}

	var funcBuffer bytes.Buffer
	fmt.Fprintf(&funcBuffer, "func initGlobalEnv() () {\n")
	c.compileLocalsTo(&funcBuffer)
	c.emitPrependBytes(funcBuffer.Bytes())
	c.emit("}")

	c.parent.emit("initGlobalEnv()\n")
}

func (c *GoCompiler) InitMethodCompiler(location *position.Location) (Compiler, int) {
	methodCompiler := NewGoCompiler("methodDefinitions", topLevelGoCompilerMode, c.loc, c.checker, c.bigIntCache, c.symbolCache, c.output)
	methodCompiler.Errors = c.Errors
	methodCompiler.SetParent(c)

	return methodCompiler, 0
}

func (c *GoCompiler) CompileMethods(location *position.Location, execOffset int) {
	c.registerGoLocal("class", "*value.Class")
	c.compileMethodsWithinModule(c.checker.Env().Root, location)

	if c.buff.Len() == 0 {
		return
	}

	c.parent.emit("methodDefinitions()\n")

	var funcBuffer bytes.Buffer
	fmt.Fprintf(&funcBuffer, "func methodDefinitions() {\n")
	c.compileLocalsTo(&funcBuffer)
	c.emitPrependBytes(funcBuffer.Bytes())
	c.emit("}\n")
}

func (c *GoCompiler) InitIvarIndicesCompiler(location *position.Location) (Compiler, int) {
	ivarCompiler := NewGoCompiler("ivarIndices", topLevelGoCompilerMode, c.loc, c.checker, c.bigIntCache, c.symbolCache, c.output)
	ivarCompiler.Errors = c.Errors
	ivarCompiler.SetParent(c)

	return ivarCompiler, 0
}

func (c *GoCompiler) FinishIvarIndicesCompiler(location *position.Location, execOffset int) Compiler {
	if c.buff.Len() == 0 {
		return c.parent
	}

	if c.parent != nil {
		c.parent.emit("ivarIndices(thread)\n")
	}

	var funcBuffer bytes.Buffer
	fmt.Fprintf(&funcBuffer, "func ivarIndices(thread *vm.Thread) {\n")
	c.compileLocalsTo(&funcBuffer)
	c.emitPrependBytes(funcBuffer.Bytes())
	c.emit("}\n")

	return c.parent
}

func (c *GoCompiler) CompileConstantDeclaration(node *ast.ConstantDeclarationNode, namespace types.Namespace, constName value.Symbol) {
	c.registerGoLocal("namespace", goValueType)

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
	c.emit("value.AddConstant(namespace, %s, %s)\n", constNameSymbol, c.convertToValue(init))
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
	methodCompiler.compileMethodBody(node.Parameters, node.Body, node.Location())

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
	goLocals          ds.OrderedMap[string, *goLocal]
	tmpLocalCounter   int
	lastElkLocalIndex int
	callCacheCounter  int
	bigIntCache       *concurrent.Map[string, *nativeBigInt]
	symbolCache       *concurrent.Map[string, *nativeSymbol]
	isGenerator       bool
	isAsync           bool
	unhygienic        bool
}

func NewGoCompiler(name string, mode goMode, loc *position.Location, checker types.Checker, bigIntCache *concurrent.Map[string, *nativeBigInt], symbolCache *concurrent.Map[string, *nativeSymbol], output io.Writer) *GoCompiler {
	return &GoCompiler{
		FuncName:          name,
		mode:              mode,
		Errors:            diagnostic.NewSyncDiagnosticList(),
		scopes:            nativeElkScopes{newNativeElkScope("", defaultNativeElkScopeType)}, // start with an empty set for the 0th scope
		goLocals:          ds.MakeOrderedMap[string, *goLocal](),
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
	if parent == nil {
		return
	}

	p := parent.(*GoCompiler)
	c.parent = p
	p.children.Append(c)
}

func (c *GoCompiler) registerGoLocal(name string, goType string) *goLocal {
	return c.registerGoLocalWithComment(name, goType, "")
}

func (c *GoCompiler) registerGoLocalWithComment(name string, goType string, comment string) *goLocal {
	if local, exists := c.goLocals.GetOk(name); exists {
		return local
	}
	local := newGoLocal(name, goType, comment)
	c.goLocals.Set(name, local)
	return local
}

func (c *GoCompiler) compileGlobalEnv() {
	env := c.checker.Env()
	c.compileModuleDefinition(env.Root, env.Root, value.ToSymbol("Root"))

	c.registerGoLocal("parentNamespace", goValueType)
	c.registerGoLocal("namespace", goValueType)

	c.registerGoLocal("class", "*value.Class")
	c.registerGoLocal("superclass", "*value.Class")
	c.registerGoLocal("mixin", "*value.Mixin")
}

// Entry point for compiling the body of a method.
func (c *GoCompiler) compileMethodBody(parameters []ast.ParameterNode, body []ast.StatementNode, loc *position.Location) {
	c.compileMethodFuncLiteralBody(parameters, body)

	var funcBuffer bytes.Buffer
	fmt.Fprintf(&funcBuffer, "func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) { // loc: %s\n", loc.String())
	c.compileLocalsTo(&funcBuffer)
	c.emitPrependBytes(funcBuffer.Bytes())
	c.emit("}")
}

func (c *GoCompiler) compileMethodFuncLiteralBody(parameters []ast.ParameterNode, body []ast.StatementNode) {
	c.registerGoLocal("self", goValueType)
	c.emit("self = args[0]\n")

	for i, param := range parameters {
		p := param.(*ast.MethodParameterNode)
		pSpan := p.Location()

		pName := identifierToName(p.Name)
		local := c.defineLocal(pName, c.typeOf(p.TypeNode), pSpan)
		if local == nil {
			return
		}

		localName := local.goIdent()
		c.emit("%s = args[%d]\n", localName, i+1)

		if p.Initialiser != nil {
			c.emit("if %s.IsUndefined() {\n", localName)
			val := c.compileExpression(p.Initialiser)
			c.emit("%s = %s", localName, c.convertToValue(val))
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

	c.emitReturn(c.convertToValue(val))
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
	c.registerGoLocal("err", goValueType)
	symbol := c.emitSymbol(name.String())
	c.emit("err = value.SetInstanceVariableByName(self, %s, %s)\n", symbol, val)
	c.emit("if err.IsNotUndefined() { panic(err) }\n")
}

func (c *GoCompiler) compileMethodsWithinModule(module *types.Module, location *position.Location) {
	if types.NamespaceHasAnyDefinableMethods(module) {
		nameSymbol := c.emitSymbol(module.Name())
		c.emit("class = (*value.Module)(value.GetConstant(%s).Pointer()).SingletonClass() // %s\n", nameSymbol, module.Name())

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
		c.emit("class = (*value.Class)(value.GetConstant(%s).Pointer()) // %s\n", namespaceSymbol, namespace.Name())

		for methodName, method := range types.SortedOwnMethods(namespace) {
			c.compileMethodDefinition(methodName, method, location)
		}

		if singletonHasCompiledMethods {
			c.emit("class = class.SingletonClass() // &%s\n", namespace.Name())

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
			c.registerGoLocal("aliasClass", "*value.Class")

			namespaceSymbol := c.emitSymbol(method.DefinedUnder.Name())
			switch namespace.(type) {
			case *value.Class:
				c.emit("aliasClass = (*value.Class)(value.GetConstant(%s).Pointer())\n", namespaceSymbol)
			case *value.Module:
				c.emit("aliasClass = value.GetConstant(%s).SingletonClass()\n", namespaceSymbol)
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

	c.emit("vm.Def(&class.MethodContainer, %q,\n", name.String())

	methodCompiler := (*GoCompiler)(method.Body.(*GoSourceMethod))
	c.emitBytes(methodCompiler.buff.Bytes())
	c.emitPackageBytes(methodCompiler.packageBuff.Bytes())
	methodCompiler.buff.Reset()
	methodCompiler.packageBuff.Reset()
	c.emit(",\n")

	if len(method.Params) > 0 {
		c.emit("vm.DefWithParameters(%d),\n", len(method.Params))
	}

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
	constNameSymbol := c.emitSymbol(name.String())

	return c.valueToNarrowerType(
		newInlineGoValue(
			fmt.Sprintf("value.GetConstant(%s)\n", constNameSymbol),
			typ,
			"value.Value",
		),
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
	c.emit("class = (*value.Class)(value.GetConstant(%s).Pointer())\n", classNameSymbol)
	c.emit("superclass = (*value.Class)(value.GetConstant(%s).Pointer())\n", superclassNameSymbol)
	c.emit("class.SetSuperclass(superclass)\n")
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

func mangleIdentifier(name string) string {
	var b strings.Builder

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

func mangleFileName(name string) string {
	return fmt.Sprintf("__file_%s", mangleIdentifier(name))
}

func (c *GoCompiler) InitExpressionCompiler(location *position.Location) Compiler {
	name := mangleFileName(location.FilePath)
	exprCompiler := NewGoCompiler(name, topLevelGoCompilerMode, location, c.checker, c.bigIntCache, c.symbolCache, c.output)
	exprCompiler.SetParent(c)
	exprCompiler.Errors = c.Errors

	return exprCompiler
}

func (c *GoCompiler) CompileExpressionsInFile(node *ast.ProgramNode) {
	c.compileProgram(node)

	if c.buff.Len() == 0 {
		return
	}

	if c.parent != nil {
		c.parent.emit("%s(thread)\n", c.FuncName)
	}

	var funcBuffer bytes.Buffer
	if c.FuncName == "main" {
		fmt.Fprintf(&funcBuffer, "func %s() { // loc: %s\n", c.FuncName, c.loc.FilePath)
		fmt.Fprintf(&funcBuffer, "thread := vm.New()\n")
	} else {
		fmt.Fprintf(&funcBuffer, "func %s(thread *vm.Thread) { // loc: %s\n", c.FuncName, c.loc.FilePath)
	}
	c.registerGoLocal("self", goValueType)
	c.compileLocalsTo(&funcBuffer)
	fmt.Fprintf(&funcBuffer, "self = value.Ref(value.GlobalObject)\n")
	c.emitPrependBytes(funcBuffer.Bytes())
	c.emit("}\n")
}

func (c *GoCompiler) compileLocalsTo(buff io.Writer) {
	for _, local := range c.goLocals.All() {
		fmt.Fprintf(buff, "var %s %s", local.name, local.goType)
		if local.comment != "" {
			fmt.Fprintf(buff, " // %s", local.comment)
		}
		fmt.Fprintf(buff, "\n_ = %s\n", local.name)
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
		return nilGoValue
	}
	return lastValue
}

func (c *GoCompiler) compileStatement(node ast.StatementNode) *goValue {
	switch node := node.(type) {
	case *ast.ExpressionStatementNode:
		return c.compileExpression(node.Expression)
	default:
		return nilGoValue
	}
}

func (c *GoCompiler) compileExpression(node ast.ExpressionNode) *goValue {
	switch node := node.(type) {
	case nil, *ast.AliasDeclarationNode, *ast.IncludeExpressionNode,
		*ast.MethodDefinitionNode, *ast.UsingExpressionNode,
		*ast.ConstantDeclarationNode, *ast.TypeDefinitionNode, *ast.GenericTypeDefinitionNode,
		*ast.MethodSignatureDefinitionNode, *ast.ImplementExpressionNode,
		*ast.StructDeclarationNode, *ast.GenericReceiverlessMethodCallNode,
		*ast.ReceiverlessMethodCallNode, *ast.AttrDeclarationNode,
		*ast.SetterDeclarationNode, *ast.GetterDeclarationNode, *ast.InitDefinitionNode,
		*ast.InstanceVariableDeclarationNode, *ast.MacroDefinitionNode,
		*ast.ReceiverlessMacroCallNode, *ast.MacroCallNode, *ast.ScopedMacroCallNode:
		return nilGoValue
	case *ast.ClassDeclarationNode:
		return c.compileClassDeclarationNode(node)
	case *ast.ModuleDeclarationNode:
		return c.compileModuleDeclarationNode(node)
	case *ast.MixinDeclarationNode:
		return c.compileMixinDeclarationNode(node)
	case *ast.InterfaceDeclarationNode:
		return c.compileInterfaceDeclarationNode(node)
	case *ast.VariableDeclarationNode:
		return c.compileVariableDeclarationNode(node)
	case *ast.MethodCallNode:
		return c.compileMethodCallNode(node)
	case *ast.GenericMethodCallNode:
		return c.compileGenericMethodCallNode(node)
	case *ast.PublicConstantNode:
		return c.compilePublicConstantNode(node)
	case *ast.PrivateConstantNode:
		return c.compilePrivateConstantNode(node)
	case *ast.GenericConstantNode:
		return c.compileExpression(node.Constant)
	case *ast.SelfLiteralNode:
		return newInlineGoValue("self", c.checker.SelfType(), "value.Value")
	case *ast.ReturnExpressionNode:
		return c.compileReturnExpressionNode(node)
	case *ast.IntLiteralNode:
		i, err := value.ParseBigInt(node.Value, 0)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		if i.IsSmallInt() {
			return newInlineGoValue(
				fmt.Sprintf("value.SmallInt(%d)", i.ToSmallInt()),
				c.typeOf(node),
				"value.SmallInt",
			)
		}
		bigIntVar := c.emitBigInt(node.Value)
		return newInlineGoValue(
			bigIntVar,
			c.typeOf(node),
			"*value.BigInt",
		)
	case *ast.Int8LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 8)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineGoValue(
			fmt.Sprintf("value.Int8(%d)", i),
			c.typeOf(node),
			"value.Int8",
		)
	case *ast.Int16LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 16)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineGoValue(
			fmt.Sprintf("value.Int16(%d)", i),
			c.typeOf(node),
			"value.Int16",
		)
	case *ast.Int32LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 32)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineGoValue(
			fmt.Sprintf("value.Int32(%d)", i),
			c.typeOf(node),
			"value.Int32",
		)
	case *ast.Int64LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 64)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineGoValue(
			fmt.Sprintf("value.Int64(%d)", i),
			c.typeOf(node),
			"value.Int64",
		)
	case *ast.UInt8LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 8)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineGoValue(
			fmt.Sprintf("value.UInt8(%d)", i),
			c.typeOf(node),
			"value.UInt8",
		)
	case *ast.UInt16LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 16)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineGoValue(
			fmt.Sprintf("value.UInt16(%d)", i),
			c.typeOf(node),
			"value.UInt16",
		)
	case *ast.UInt32LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 32)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineGoValue(
			fmt.Sprintf("value.UInt32(%d)", i),
			c.typeOf(node),
			"value.UInt32",
		)
	case *ast.UInt64LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 64)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineGoValue(
			fmt.Sprintf("value.UInt64(%d)", i),
			c.typeOf(node),
			"value.UInt64",
		)
	case *ast.BinaryExpressionNode:
		return c.compileBinaryExpressionNode(node)
	case *ast.AssignmentExpressionNode:
		return c.compileAssignmentExpressionNode(node)
	case *ast.PublicIdentifierNode:
		return c.compileLocalVariableAccess(node.Value)
	case *ast.PrivateIdentifierNode:
		return c.compileLocalVariableAccess(node.Value)
	case *ast.IfExpressionNode:
		return c.compileIfExpression(
			ifConditionType,
			node.Condition,
			node.ThenBody,
			node.ElseBody,
			c.typeOf(node),
		)
	case *ast.UnlessExpressionNode:
		return c.compileIfExpression(
			unlessConditionType,
			node.Condition,
			node.ThenBody,
			node.ElseBody,
			c.typeOf(node),
		)
	case *ast.ModifierIfElseNode:
		return c.compileModifierIfExpression(
			ifConditionType,
			node.Condition,
			node.ThenExpression,
			node.ElseExpression,
			c.typeOf(node),
		)
	default:
		panic(fmt.Sprintf("invalid expression node: %T", node))
	}
}

func (c *GoCompiler) compileVariableDeclarationNode(node *ast.VariableDeclarationNode) *goValue {
	initialised := node.Initialiser != nil

	local := c.defineLocal(identifierToName(node.Name), c.typeOf(node.TypeNode), node.Location())
	if local == nil {
		return errGoValue
	}

	if initialised {
		init := c.compileExpression(node.Initialiser)
		return c.emitSetLocal(local.name, c.convertToValue(init))
	}

	return nilGoValue
}

func (c *GoCompiler) compileReturnExpressionNode(node *ast.ReturnExpressionNode) *goValue {
	var val string
	if node.Value != nil {
		val = c.convertToValue(c.compileExpression(node.Value))
	} else {
		val = "value.Nil"
	}

	c.emitReturn(val)
	return nilGoValue
}

func (c *GoCompiler) compilePublicConstantNode(node *ast.PublicConstantNode) *goValue {
	return c.emitGetConst(value.ToSymbol(node.Value), c.typeOf(node))
}

func (c *GoCompiler) compilePrivateConstantNode(node *ast.PrivateConstantNode) *goValue {
	return c.emitGetConst(value.ToSymbol(node.Value), c.typeOf(node))
}

func (c *GoCompiler) compileMethodCallNode(node *ast.MethodCallNode) *goValue {
	return c.compileMethodCall(
		node.Receiver,
		node.Op,
		node.MethodName,
		node.PositionalArguments,
		c.typeOf(node),
		node.Location(),
	)
}

func (c *GoCompiler) compileGenericMethodCallNode(node *ast.GenericMethodCallNode) *goValue {
	return c.compileMethodCall(
		node.Receiver,
		node.Op,
		node.MethodName,
		node.PositionalArguments,
		c.typeOf(node),
		node.Location(),
	)
}

func (c *GoCompiler) compileMethodCall(receiver ast.ExpressionNode, op *token.Token, nameNode ast.IdentifierNode, args []ast.ExpressionNode, typ types.Type, location *position.Location) *goValue {
	name := identifierToName(nameNode)

	switch op.Type {
	case token.QUESTION_DOT:
		receiverVal := c.compileExpression(receiver)
		resultVar := c.defineTmpGoLocal(goValueType)

		receiverValString := c.convertToValue(receiverVal)
		c.emit("if value.IsNil(%s) {\n", receiverValString)
		c.emit("%s = value.Nil", resultVar)
		c.emit("} else {\n")
		callResult := c.compileInnerMethodCall(receiverValString, c.typeOf(receiver), name, op, args, typ, location)
		c.emit("%s = %s", resultVar, c.convertToValue(callResult))
		c.emit("}\n")

		return newTmpGoValue(resultVar, typ)
	case token.QUESTION_DOT_DOT:
		receiverVal := c.compileExpression(receiver)
		resultVar := c.defineTmpGoLocal(goValueType)

		receiverValString := c.convertToValue(receiverVal)
		c.emit("if value.IsNil(%s) {\n", receiverValString)
		c.emit("%s = value.Nil", resultVar)
		c.emit("} else {\n")
		c.compileInnerMethodCall(receiverValString, c.typeOf(receiver), name, op, args, typ, location)
		c.emit("%s = %s", resultVar, receiverVal)
		c.emit("}\n")

		return newTmpGoValue(resultVar, typ)
	case token.DOT_DOT:
		receiverVal := c.compileExpression(receiver)
		resultVar := c.defineTmpGoLocal(goValueType)

		receiverValString := c.convertToValue(receiverVal)
		c.compileInnerMethodCall(receiverValString, c.typeOf(receiver), name, op, args, typ, location)
		c.emit("%s = %s", resultVar, receiverVal)

		return newTmpGoValue(resultVar, typ)
	case token.DOT:
		receiverVal := c.compileExpression(receiver)
		receiverValString := c.convertToValue(receiverVal)
		return c.compileInnerMethodCall(receiverValString, c.typeOf(receiver), name, op, args, typ, location)
	default:
		panic(fmt.Sprintf("invalid method call operator: %#v", op))
	}
}

func (c *GoCompiler) compileInnerMethodCall(receiverVal string, receiverType types.Type, name string, op *token.Token, args []ast.ExpressionNode, typ types.Type, location *position.Location) *goValue {
	// TODO: implement closures and optimised increments/decrements
	// receiverType := c.typeOf(receiver)
	// switch name {
	// case "call":
	// 	return c.emitCall(callInfo, location)
	// case "++":
	// 	return c.compileIncrement(receiverType, location)
	// case "--":
	// 	return c.compileDecrement(receiverType, location)
	// }

	callArgsVar := c.defineTmpGoLocal("[]value.Value")
	c.emit("%s = make([]value.Value, %d)\n", callArgsVar, len(args)+1)

	c.emit("%s[0] = %s\n", callArgsVar, receiverVal)
	for i, posArg := range args {
		argVal := c.compileExpression(posArg)
		c.emit("%s[%d] = %s\n", callArgsVar, i+1, c.convertToValue(argVal))
	}

	nameSym := c.emitSymbol(name)
	callCache := c.emitCallCache()
	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitAddCallFrame(location)
	c.emit("%s, err = thread.CallMethodByNameWithCache(%s, &%s, %s...)\n", tmp, nameSym, callCache, callArgsVar)
	c.emit("%s = nil\n", callArgsVar)
	c.emitPopCallFrame()
	c.emitErrorPropagation()

	return newTmpGoValue(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileModuleDeclarationNode(node *ast.ModuleDeclarationNode) *goValue {
	typ := c.typeOf(node).(*types.Module)
	return c.compileNamespaceDeclarationNode(fmt.Sprintf("module_%s", mangleIdentifier(typ.Name())), node.Body, typ, node.Location())
}

func (c *GoCompiler) compileInterfaceDeclarationNode(node *ast.InterfaceDeclarationNode) *goValue {
	typ := c.typeOf(node).(*types.Interface)
	return c.compileNamespaceDeclarationNode(fmt.Sprintf("interface_%s", mangleIdentifier(typ.Name())), node.Body, typ, node.Location())
}

func (c *GoCompiler) compileMixinDeclarationNode(node *ast.MixinDeclarationNode) *goValue {
	typ := c.typeOf(node).(*types.Mixin)
	return c.compileNamespaceDeclarationNode(fmt.Sprintf("mixin_%s", mangleIdentifier(typ.Name())), node.Body, typ, node.Location())
}

func (c *GoCompiler) compileClassDeclarationNode(node *ast.ClassDeclarationNode) *goValue {
	typ := c.typeOf(node).(*types.Class)
	return c.compileNamespaceDeclarationNode(fmt.Sprintf("class_%s", mangleIdentifier(typ.Name())), node.Body, typ, node.Location())
}

func (c *GoCompiler) compileNamespaceDeclarationNode(name string, body []ast.StatementNode, typ types.Namespace, loc *position.Location) *goValue {
	if len(body) <= 0 {
		return nilGoValue
	}

	classCompiler := NewGoCompiler(name, topLevelGoCompilerMode, loc, c.checker, c.bigIntCache, c.symbolCache, c.output)
	classCompiler.SetParent(c)
	classCompiler.Errors = c.Errors

	classCompiler.compileNamespaceBody(body, typ)

	return nilGoValue
}

func (c *GoCompiler) compileNamespaceBody(body []ast.StatementNode, typ types.Namespace) {
	c.registerGoLocal("self", goValueType)
	c.compileStatements(body)
	if c.buff.Len() == 0 {
		return
	}

	c.parent.emit("%s(thread)\n", c.FuncName)

	var funcBuffer bytes.Buffer
	fmt.Fprintf(&funcBuffer, "func %s(thread *vm.Thread) { // namespace: %s, loc: %s\n", c.FuncName, typ.Name(), c.loc.String())
	c.compileLocalsTo(&funcBuffer)

	switch typ := typ.(type) {
	case *types.SingletonClass:
		nameSymbol := c.emitSymbol(typ.AttachedObject.Name())
		fmt.Fprintf(&funcBuffer, "self = value.Ref(value.GetConstant(%s).SingletonClass())\n", nameSymbol)
	case *types.Module:
		nameSymbol := c.emitSymbol(typ.Name())
		fmt.Fprintf(&funcBuffer, "self = value.Ref(value.GetConstant(%s).SingletonClass())\n", nameSymbol)
	default:
		nameSymbol := c.emitSymbol(typ.Name())
		fmt.Fprintf(&funcBuffer, "self = value.GetConstant(%s)\n", nameSymbol)
	}

	c.emitPrependBytes(funcBuffer.Bytes())
	c.emit("}\n")
}

func (c *GoCompiler) compileAssignmentExpressionNode(node *ast.AssignmentExpressionNode) *goValue {
	switch n := node.Left.(type) {
	case *ast.PublicIdentifierNode:
		return c.localVariableAssignment(n.Value, node.Op, node.Right, c.typeOf(node), node.Location())
	case *ast.PrivateIdentifierNode:
		return c.localVariableAssignment(n.Value, node.Op, node.Right, c.typeOf(node), node.Location())
		// TODO: Implement all assignment types
	// case *ast.SubscriptExpressionNode:
	// 	return c.subscriptAssignment(node, n, valueIsIgnored)
	// case *ast.PublicInstanceVariableNode:
	// 	return c.instanceVariableAssignment(node, n, valueIsIgnored)
	// case *ast.AttributeAccessNode:
	// 	return c.attributeAssignment(node, n)
	default:
		c.Errors.AddFailure(
			fmt.Sprintf("cannot assign to: %T", node.Left),
			node.Location(),
		)
		return errGoValue
	}
}

func (c *GoCompiler) localVariableAssignment(name string, operator *token.Token, right ast.ExpressionNode, typ types.Type, loc *position.Location) *goValue {
	switch operator.Type {
	case token.OR_OR_EQUAL:
		varIdent := c.compileLocalVariableAccess(name)
		c.emit("if value.Falsy(%s) {\n", c.convertToValue(varIdent))

		rightVal := c.compileExpression(right)
		c.emit("%s = %s\n", varIdent.value(), c.convertToValue(rightVal))

		c.emit("}\n")

		return varIdent
	case token.AND_AND_EQUAL:
		varIdent := c.compileLocalVariableAccess(name)
		c.emit("if value.Truthy(%s) {\n", c.convertToValue(varIdent))

		rightVal := c.compileExpression(right)
		c.emit("%s = %s\n", varIdent.value(), c.convertToValue(rightVal))

		c.emit("}\n")

		return varIdent
	case token.QUESTION_QUESTION_EQUAL:
		varIdent := c.compileLocalVariableAccess(name)
		c.emit("if value.IsNil(%s) {\n", c.convertToValue(varIdent))

		rightVal := c.compileExpression(right)
		c.emit("%s = %s\n", varIdent.value(), c.convertToValue(rightVal))

		c.emit("}\n")

		return varIdent
	case token.EQUAL_OP:
		return c.setLocal(name, right)
	case token.COLON_EQUAL:
		c.defineLocal(name, typ, loc)
		return c.setLocal(name, right)
	default:
		c.Errors.AddFailure(
			fmt.Sprintf("assignment using this operator has not been implemented: %s", operator.Type.Name()),
			loc,
		)
		return errGoValue
	}
}

func (c *GoCompiler) setLocal(name string, valueNode ast.ExpressionNode) *goValue {
	val := c.compileExpression(valueNode)
	return c.emitSetLocal(name, c.convertToValue(val))
}

func (c *GoCompiler) emitSetLocal(name string, val string) *goValue {
	var variable *nativeElkLocal
	if local, ok := c.resolveLocal(name); ok {
		variable = local
	} else if upvalue, ok := c.resolveUpvalue(name); ok {
		variable = upvalue
	} else {
		panic(fmt.Sprintf("undefined local: %s\n", name))
	}

	ident := variable.goIdent()
	c.emit("%s = %s\n", ident, val)
	return newTmpGoValue(
		variable.goLocal,
		variable.elkType,
	)
}

func (c *GoCompiler) compileModifierIfExpression(condType conditionType, condition, then, els ast.ExpressionNode, typ types.Type) *goValue {
	var elsFunc func() *goValue
	if els != nil {
		elsFunc = func() *goValue {
			return c.compileExpression(els)
		}
	}

	return c.compileIfWithConditionExpression(
		condType,
		condition,
		func() *goValue {
			return c.compileExpression(then)
		},
		elsFunc,
		typ,
	)
}

func (c *GoCompiler) compileIfExpression(condType conditionType, condition ast.ExpressionNode, then, els []ast.StatementNode, typ types.Type) *goValue {
	var elsFunc func() *goValue
	if els != nil {
		elsFunc = func() *goValue {
			return c.compileStatements(els)
		}
	}

	return c.compileIfWithConditionExpression(
		condType,
		condition,
		func() *goValue {
			return c.compileStatements(then)
		},
		elsFunc,
		typ,
	)
}

func (c *GoCompiler) compileIfWithConditionExpression(condType conditionType, condition ast.ExpressionNode, then, els func() *goValue, typ types.Type) *goValue {
	if result := resolve(condition); !result.IsUndefined() {
		// if gets optimised away
		c.enterScope("", defaultNativeElkScopeType)
		defer c.leaveScope()

		var checkFunc func(value.Value) bool
		switch condType {
		case ifConditionType:
			checkFunc = value.Truthy
		case unlessConditionType:
			checkFunc = value.Falsy
		case isNilConditionType:
			checkFunc = value.IsNil
		default:
			panic(fmt.Sprintf("invalid if condition type: %d", condType))
		}

		if checkFunc(result) {
			if then == nil {
				return nilGoValue
			}
			return then()
		}

		if els == nil {
			return nilGoValue
		}
		return els()
	}

	cond := func() *goValue {
		return c.compileExpression(condition)
	}

	return c.compileIf(
		condType,
		cond,
		then,
		els,
		typ,
	)
}

type conditionType uint8

const (
	ifConditionType conditionType = iota
	unlessConditionType
	isNilConditionType
)

func (c *GoCompiler) compileIf(condType conditionType, condition, then, els func() *goValue, typ types.Type) *goValue {
	c.enterScope("", defaultNativeElkScopeType)
	condVal := condition()

	ifResultVar := c.defineTmpGoLocal(goValueType)

	var condFunc string
	switch condType {
	case ifConditionType:
		condFunc = "value.Truthy"
	case unlessConditionType:
		condFunc = "value.Falsy"
	case isNilConditionType:
		condFunc = "value.IsNil"
	default:
		panic(fmt.Sprintf("invalid if condition type: %d", condType))
	}

	c.emit("if %s(%s) {\n", condFunc, c.convertToValue(condVal))
	thenVal := then()
	c.emit("%s = %s\n", ifResultVar, c.convertToValue(thenVal))
	c.emit("}")

	c.leaveScope()

	if els != nil {
		c.emit(" else {\n")
		elseVal := els()
		c.emit("%s = %s\n", ifResultVar, c.convertToValue(elseVal))
		c.emit("}")
	}

	c.emit("\n")

	return newTmpGoValue(
		ifResultVar,
		typ,
	)
}

func (c *GoCompiler) compileLocalVariableAccess(name string) *goValue {
	if local, ok := c.resolveLocal(name); ok {
		return newTmpGoValue(
			local.goLocal,
			local.elkType,
		)
	}

	if upvalue, ok := c.resolveUpvalue(name); ok {
		return newTmpGoValue(
			upvalue.goLocal,
			upvalue.elkType,
		)
	}

	panic(fmt.Sprintf("no such variable: %s", name))
}

// Resolve a local variable
func (c *GoCompiler) resolveLocal(name string) (*nativeElkLocal, bool) {
	var localVal *nativeElkLocal
	var found bool
	for i := len(c.scopes) - 1; i >= 0; i-- {
		varScope := c.scopes[i]
		local, ok := varScope.localTable[name]
		if ok {
			localVal = local
			found = true
			break
		}
		if !c.unhygienic && varScope.typ == macroBoundaryNativeElkScopeType {
			break
		}
	}

	if !found {
		return localVal, false
	}

	return localVal, true
}

// Resolve an upvalue from an outer context
func (c *GoCompiler) resolveUpvalue(name string) (*nativeElkLocal, bool) {
	if c.mode != closureGoCompilerMode {
		// don't look into parent compilers if we aren't in a closure
		return nil, false
	}

	parent := c.parent
	if parent == nil {
		return nil, false
	}
	local, ok := parent.resolveLocal(name)
	if ok {
		return local, true
	}

	local, ok = parent.resolveUpvalue(name)
	if ok {
		return local, true
	}

	return nil, false
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
	c.emitPackage("var %s = value.ParseBigIntPanic(%q, 0)\n", ident, val)
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

func (c *GoCompiler) defineTmpGoLocal(goType string) *goLocal {
	for local := range c.goLocals.Values() {
		if local.free && local.goType == goType {
			local.markReserved()
			return local
		}
	}

	tmp := c.getTmpIdent()
	return c.registerGoLocal(tmp, goType)
}

func (c *GoCompiler) emitErrorPropagation() {
	switch c.mode {
	case topLevelGoCompilerMode:
		c.emit("if err.IsNotUndefined() { thread.Panic(err) }\n")
	default:
		c.emit("if err.IsNotUndefined() { return value.Undefined, err }\n")
	}
}

func (c *GoCompiler) emitCallCache() string {
	c.callCacheCounter++
	callCacheName := fmt.Sprintf("cc_%s_%d", c.FuncName, c.callCacheCounter)
	c.emitPackage("var %s = &value.CallCache{}\n", callCacheName)

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
	c.registerGoLocal("err", goValueType)
}

func (c *GoCompiler) compileBinaryExpressionNode(node *ast.BinaryExpressionNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	switch node.Op.Type {
	case token.PLUS:
		return c.compileAdd(node)
	case token.MINUS:
		return c.compileSubtract(node)
	case token.STAR:
		return c.compileMultiply(node)
	// case token.SLASH:
	// 	c.emitAddCallFrame(node.Location())

	// 	if c.checker.IsSubtype(typ, c.checker.StdInt()) {
	// 		tmp := c.defineTmpGoLocal(goValueType)
	// 		c.registerErr()
	// 		c.emit("%s, err = value.DivideInt(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
	// 		c.emitPopCallFrame()
	// 		c.emitErrorPropagation()

	// 		return newTmpGoValue(
	// 			tmp,
	// 			c.checker.Std(symbol.Value),
	// 		)
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
	// 		tmp := c.defineTmpGoLocal(goValueType)
	// 		c.registerErr()
	// 		c.emit("%s, err = %s.DivideVal(%s)\n", tmp, c.goValueAs(left, "AsFloat"), c.convertToValue(right))
	// 		c.emitPopCallFrame()
	// 		c.emitErrorPropagation()

	// 		return newTmpGoValue(
	// 			tmp,
	// 			c.checker.Std(symbol.Value),
	// 		)
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.StdBigFloat()) {
	// 		tmp := c.defineTmpGoLocal(goValueType)
	// 		c.registerErr()
	// 		c.emit("%s, err = %s.DivideVal(%s)\n", tmp, c.goValueCast(left, "*value.BigFloat"), c.convertToValue(right))
	// 		c.emitErrorPropagation()

	// 		return newTmpGoValue(
	// 			tmp,
	// 			c.checker.Std(symbol.Value),
	// 		)
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinAddable)) {
	// 		tmp := c.defineTmpGoLocal(goValueType)
	// 		c.registerErr()
	// 		c.emit("%s, err = value.DivideVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
	// 		c.emitPopCallFrame()
	// 		c.emitErrorPropagation()

	// 		return newTmpGoValue(
	// 			tmp,
	// 			c.checker.Std(symbol.Value),
	// 		)
	// 	}

	// 	callCache := c.emitCallCache()
	// 	tmp := c.defineTmpGoLocal(goValueType)
	// 	c.registerErr()
	// 	c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpMultiply, &%s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
	// 	c.emitPopCallFrame()
	// 	c.emitErrorPropagation()
	// 	return newTmpGoValue(
	// 		tmp,
	// 		c.checker.Std(symbol.Value),
	// 	)
	// case token.STAR_STAR:
	// 	if c.checker.IsSubtype(typ, c.checker.StdInt()) {
	// 		tmp := c.defineTmpGoLocal(goValueType)
	// 		c.registerErr()
	// 		c.emit("%s, err = value.ExponentiateInt(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
	// 		c.emitErrorPropagation()

	// 		return newTmpGoValue(
	// 			tmp,
	// 			c.checker.Std(symbol.Value),
	// 		)
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
	// 		tmp := c.defineTmpGoLocal(goValueType)
	// 		c.registerErr()
	// 		c.emit("%s, err = %s.ExponentiateVal(%s)\n", tmp, c.goValueAs(left, "AsFloat"), c.convertToValue(right))
	// 		c.emitPopCallFrame()
	// 		c.emitErrorPropagation()

	// 		return newTmpGoValue(
	// 			tmp,
	// 			c.checker.Std(symbol.Value),
	// 		)
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinNumeric)) {
	// 		tmp := c.defineTmpGoLocal(goValueType)
	// 		c.registerErr()
	// 		c.emit("%s, err = value.ExponentiateVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
	// 		c.emitPopCallFrame()
	// 		c.emitErrorPropagation()

	// 		return newTmpGoValue(
	// 			tmp,
	// 			c.checker.Std(symbol.Value),
	// 		)
	// 	}

	// 	callCache := c.emitCallCache()
	// 	tmp := c.defineTmpGoLocal(goValueType)
	// 	c.registerErr()
	// 	c.emitAddCallFrame(node.Location())
	// 	c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpMultiply, &%s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
	// 	c.emitPopCallFrame()
	// 	c.emitErrorPropagation()
	// 	return newTmpGoValue(
	// 		tmp,
	// 		c.checker.Std(symbol.Value),
	// 	)
	// case token.LBITSHIFT:
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
	// 		tmp := c.defineTmpGoLocal(goValueType)
	// 		c.registerErr()
	// 		c.emit("%s, err = value.LeftBitshiftVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
	// 		c.emitPopCallFrame()
	// 		c.emitErrorPropagation()

	// 		return newTmpGoValue(
	// 			tmp,
	// 			c.checker.Std(symbol.Value),
	// 		)
	// 	}

	// 	callCache := c.emitCallCache()
	// 	tmp := c.defineTmpGoLocal(goValueType)
	// 	c.registerErr()
	// 	c.emitAddCallFrame(node.Location())
	// 	c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpLeftBitshift, &%s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
	// 	c.emitPopCallFrame()
	// 	c.emitErrorPropagation()
	// 	return newTmpGoValue(
	// 		tmp,
	// 		c.checker.Std(symbol.Value),
	// 	)
	// case token.LTRIPLE_BITSHIFT:
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinLogicBitshiftable)) {
	// 		tmp := c.defineTmpGoLocal(goValueType)
	// 		c.registerErr()
	// 		c.emit("%s, err = value.LogicalLeftBitshiftVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
	// 		c.emitPopCallFrame()
	// 		c.emitErrorPropagation()

	// 		return newTmpGoValue(
	// 			tmp,
	// 			c.checker.Std(symbol.Value),
	// 		)
	// 	}

	// 	callCache := c.emitCallCache()
	// 	tmp := c.defineTmpGoLocal(goValueType)
	// 	c.registerErr()
	// 	c.emitAddCallFrame(node.Location())
	// 	c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpLogicalLeftBitshift, &%s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
	// 	c.emitPopCallFrame()
	// 	c.emitErrorPropagation()
	// 	return newTmpGoValue(
	// 		tmp,
	// 		c.checker.Std(symbol.Value),
	// 	)
	// case token.RBITSHIFT:
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
	// 		tmp := c.defineTmpGoLocal(goValueType)
	// 		c.registerErr()
	// 		c.emit("%s, err = value.RightBitshiftVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
	// 		c.emitPopCallFrame()
	// 		c.emitErrorPropagation()

	// 		return newTmpGoValue(
	// 			tmp,
	// 			c.checker.Std(symbol.Value),
	// 		)
	// 	}

	// 	callCache := c.emitCallCache()
	// 	tmp := c.defineTmpGoLocal(goValueType)
	// 	c.registerErr()
	// 	c.emitAddCallFrame(node.Location())
	// 	c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpRightBitshift, &%s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
	// 	c.emitPopCallFrame()
	// 	c.emitErrorPropagation()
	// 	return newTmpGoValue(
	// 		tmp,
	// 		c.checker.Std(symbol.Value),
	// 	)
	// case token.RTRIPLE_BITSHIFT:
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinLogicBitshiftable)) {
	// 		tmp := c.defineTmpGoLocal(goValueType)
	// 		c.registerErr()
	// 		c.emit("%s, err = value.LogicalRightBitshiftVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
	// 		c.emitPopCallFrame()
	// 		c.emitErrorPropagation()

	// 		return newTmpGoValue(
	// 			tmp,
	// 			c.checker.Std(symbol.Value),
	// 		)
	// 	}

	// 	callCache := c.emitCallCache()
	// 	tmp := c.defineTmpGoLocal(goValueType)
	// 	c.registerErr()
	// 	c.emitAddCallFrame(node.Location())
	// 	c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpLogicalRightBitshift, &%s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
	// 	c.emitPopCallFrame()
	// 	c.emitErrorPropagation()
	// 	return newTmpGoValue(
	// 		tmp,
	// 		c.checker.Std(symbol.Value),
	// 	)
	// case token.AND:
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
	// 		tmp := c.defineTmpGoLocal(goValueType)
	// 		c.registerErr()
	// 		c.emit("%s, err = value.BitwiseAndVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
	// 		c.emitPopCallFrame()
	// 		c.emitErrorPropagation()

	// 		return newTmpGoValue(
	// 			tmp,
	// 			c.checker.Std(symbol.Value),
	// 		)
	// 	}

	// 	callCache := c.emitCallCache()
	// 	tmp := c.defineTmpGoLocal(goValueType)
	// 	c.registerErr()
	// 	c.emitAddCallFrame(node.Location())
	// 	c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpAnd, &%s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
	// 	c.emitPopCallFrame()
	// 	c.emitErrorPropagation()
	// 	return newTmpGoValue(
	// 		tmp,
	// 		c.checker.Std(symbol.Value),
	// 	)
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
	// 		tmp := c.defineTmpGoLocal(goValueType)
	// 		c.registerErr()
	// 		c.emit("%s, err = value.LessThanEqualValInt(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
	// 		c.emitErrorPropagation()

	// 		return newTmpGoValue(
	// 			tmp,
	// 			c.checker.Std(symbol.Value),
	// 		)
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
	// 		tmp := c.defineTmpGoLocal(goValueType)
	// 		c.registerErr()
	// 		c.emit("%s, err = %s.LessThanEqualVal(%s)\n", tmp, c.goValueAs(left, "AsFloat"), c.convertToValue(right))
	// 		c.emitErrorPropagation()

	// 		return newTmpGoValue(
	// 			tmp,
	// 			c.checker.Std(symbol.Value),
	// 		)
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.StdBigFloat()) {
	// 		tmp := c.defineTmpGoLocal(goValueType)
	// 		c.registerErr()
	// 		c.emit("%s, err = %s.LessThanEqualVal(%s)\n", tmp, c.goValueCast(left, "*value.BigFloat"), c.convertToValue(right))
	// 		c.emitErrorPropagation()

	// 		return newTmpGoValue(
	// 			tmp,
	// 			c.checker.Std(symbol.Value),
	// 		)
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinAddable)) {
	// 		tmp := c.defineTmpGoLocal(goValueType)
	// 		c.registerErr()
	// 		c.emit("%s, err = value.LessThanEqualVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
	// 		c.emitErrorPropagation()

	// 		return newTmpGoValue(
	// 			tmp,
	// 			c.checker.Std(symbol.Value),
	// 		)
	// 	}

	// 	callCache := c.emitCallCache()
	// 	tmp := c.defineTmpGoLocal(goValueType)
	// 	c.registerErr()
	// 	c.emitAddCallFrame(node.Location())
	// 	c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpLessThanEqual, &%s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
	// 	c.emitPopCallFrame()
	// 	c.emitErrorPropagation()
	// 	return newTmpGoValue(
	// 		tmp,
	// 		c.checker.Std(symbol.Value),
	// 	)
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
		return errGoValue
	}
}

func (c *GoCompiler) compileMultiply(node *ast.BinaryExpressionNode) *goValue {
	left := c.compileExpression(node.Left)
	right := c.compileExpression(node.Right)
	typ := c.typeOf(node)

	switch left.goType() {
	case "value.SmallInt":
		return c.compileMultiplySmallInt(left, right, typ)
	case "*value.BigInt":
		return c.compileMultiplyBigInt(left, right, typ)
	case "value.Int64":
		return c.compileMultiplyInt64(left, right)
	case "value.Int32":
		return c.compileMultiplyInt32(left, right)
	case "value.Int16":
		return c.compileMultiplyInt16(left, right)
	case "value.Int8":
		return c.compileMultiplyInt8(left, right)
	case "value.UInt64":
		return c.compileMultiplyUInt64(left, right)
	case "value.UInt32":
		return c.compileMultiplyUInt32(left, right)
	case "value.UInt16":
		return c.compileMultiplyUInt16(left, right)
	case "value.UInt8":
		return c.compileMultiplyUInt8(left, right)
	case "value.Float":
		return c.compileMultiplyFloat(left, right, typ)
	case "*value.BigFloat":
		return c.compileMultiplyBigFloat(left, right, typ)
	case "value.Float64":
		return c.compileMultiplyFloat64(left, right)
	case "value.Float32":
		return c.compileMultiplyFloat32(left, right)
	}

	if c.checker.IsSubtype(left.elkType, c.checker.Std(symbol.S_BuiltinMultipliable)) {
		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emit("%s, err = value.MultiplyVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
		c.emitErrorPropagation()

		return newTmpGoValue(
			tmp,
			typ,
		)
	}

	callCache := c.emitCallCache()
	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitAddCallFrame(node.Location())
	c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpMultiply, &%s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
	c.emitPopCallFrame()
	c.emitErrorPropagation()
	return newTmpGoValue(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileMultiplyBigInt(left, right *goValue, typ types.Type) *goValue {
	switch right.goType() {
	case "value.SmallInt":
		return newInlineGoValue(
			fmt.Sprintf("(%s).MultiplySmallInt(%s)\n", left.value(), right.value()),
			left.elkType,
			goValueType,
		)
	case "value.Float":
		return newInlineGoValue(
			fmt.Sprintf("(%s).MultiplyFloat(%s)\n", left.value(), right.value()),
			right.elkType,
			"value.Float",
		)
	case "*value.BigFloat":
		return newInlineGoValue(
			fmt.Sprintf("(%s).MultiplyBigFloat(%s)\n", left.value(), right.value()),
			right.elkType,
			"*value.BigFloat",
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newInlineGoValue(
			fmt.Sprintf("(%s).MultiplyInt(%s)\n", left.value(), c.convertToValue(right)),
			left.elkType,
			goValueType,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emit("%s, err = (%s).MultiplyVal(%s)\n", tmp, left.value(), c.convertToValue(right))
	c.emitErrorPropagation()

	return newTmpGoValue(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileMultiplySmallInt(left, right *goValue, typ types.Type) *goValue {
	switch right.goType() {
	case "value.SmallInt":
		return newInlineGoValue(
			fmt.Sprintf("(%s).MultiplySmallInt(%s)\n", left.value(), right.value()),
			left.elkType,
			"value.SmallInt",
		)
	case "value.Float":
		return newInlineGoValue(
			fmt.Sprintf("(%s).MultiplyFloat(%s)\n", left.value(), right.value()),
			right.elkType,
			"value.Float",
		)
	case "*value.BigFloat":
		return newInlineGoValue(
			fmt.Sprintf("(%s).MultiplyBigFloat(%s)\n", left.value(), right.value()),
			right.elkType,
			"*value.BigFloat",
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newInlineGoValue(
			fmt.Sprintf("(%s).MultiplyInt(%s)\n", left.value(), c.convertToValue(right)),
			left.elkType,
			goValueType,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emit("%s, err = (%s).MultiplyVal(%s)\n", tmp, left.value(), c.convertToValue(right))
	c.emitErrorPropagation()

	return newTmpGoValue(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileMultiplyFloat(left, right *goValue, typ types.Type) *goValue {
	switch right.goType() {
	case "value.SmallInt":
		return newInlineGoValue(
			fmt.Sprintf("(%s).MultiplySmallInt(%s)\n", left.value(), right.value()),
			left.elkType,
			"value.Float",
		)
	case "value.Float":
		return newInlineGoValue(
			fmt.Sprintf("(%s).MultiplyFloat(%s)\n", left.value(), right.value()),
			left.elkType,
			"value.Float",
		)
	case "*value.BigFloat":
		return newInlineGoValue(
			fmt.Sprintf("(%s).MultiplyBigFloat(%s)\n", left.value(), right.value()),
			right.elkType,
			"*value.BigFloat",
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newInlineGoValue(
			fmt.Sprintf("(%s).MultiplyInt(%s)\n", left.value(), c.convertToValue(right)),
			left.elkType,
			"value.Float",
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emit("%s, err = (%s).MultiplyVal(%s)\n", tmp, left.value(), c.convertToValue(right))
	c.emitErrorPropagation()

	return newTmpGoValue(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileMultiplyBigFloat(left, right *goValue, typ types.Type) *goValue {
	switch right.goType() {
	case "value.SmallInt":
		return newInlineGoValue(
			fmt.Sprintf("(%s).MultiplySmallInt(%s)\n", left.value(), right.value()),
			left.elkType,
			"*value.BigFloat",
		)
	case "value.Float":
		return newInlineGoValue(
			fmt.Sprintf("(%s).MultiplyFloat(%s)\n", left.value(), right.value()),
			left.elkType,
			"*value.BigFloat",
		)
	case "*value.BigFloat":
		return newInlineGoValue(
			fmt.Sprintf("(%s).MultiplyBigFloat(%s)\n", left.value(), right.value()),
			left.elkType,
			"*value.BigFloat",
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newInlineGoValue(
			fmt.Sprintf("(%s).MultiplyInt(%s)\n", left.value(), c.convertToValue(right)),
			left.elkType,
			"*value.BigFloat",
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emit("%s, err = (%s).MultiplyVal(%s)\n", tmp, left.value(), c.convertToValue(right))
	c.emitErrorPropagation()

	return newTmpGoValue(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileMultiplyFloat64(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.Float64":
		return newInlineGoValue(
			fmt.Sprintf("(%s * %s)\n", left.value(), right.value()),
			left.elkType,
			"value.Float64",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s * %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.Float64",
	)
}

func (c *GoCompiler) compileMultiplyFloat32(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.Float32":
		return newInlineGoValue(
			fmt.Sprintf("(%s * %s)\n", left.value(), right.value()),
			left.elkType,
			"value.Float32",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s * %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.Float32",
	)
}

func (c *GoCompiler) compileMultiplyInt64(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.Int64":
		return newInlineGoValue(
			fmt.Sprintf("(%s * %s)\n", left.value(), right.value()),
			left.elkType,
			"value.Int64",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s * %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.Int64",
	)
}

func (c *GoCompiler) compileMultiplyInt32(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.Int32":
		return newInlineGoValue(
			fmt.Sprintf("(%s * %s)\n", left.value(), right.value()),
			left.elkType,
			"value.Int32",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s * %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.Int32",
	)
}

func (c *GoCompiler) compileMultiplyInt16(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.Int16":
		return newInlineGoValue(
			fmt.Sprintf("(%s * %s)\n", left.value(), right.value()),
			left.elkType,
			"value.Int16",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s * %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.Int16",
	)
}

func (c *GoCompiler) compileMultiplyInt8(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.Int8":
		return newInlineGoValue(
			fmt.Sprintf("(%s * %s)\n", left.value(), right.value()),
			left.elkType,
			"value.Int8",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s * %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.Int8",
	)
}

func (c *GoCompiler) compileMultiplyUInt64(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.UInt64":
		return newInlineGoValue(
			fmt.Sprintf("(%s * %s)\n", left.value(), right.value()),
			left.elkType,
			"value.UInt64",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s * %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.UInt64",
	)
}

func (c *GoCompiler) compileMultiplyUInt32(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.UInt32":
		return newInlineGoValue(
			fmt.Sprintf("(%s * %s)\n", left.value(), right.value()),
			left.elkType,
			"value.UInt32",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s * %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.UInt32",
	)
}

func (c *GoCompiler) compileMultiplyUInt16(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.UInt16":
		return newInlineGoValue(
			fmt.Sprintf("(%s * %s)\n", left.value(), right.value()),
			left.elkType,
			"value.UInt16",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s * %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.UInt16",
	)
}

func (c *GoCompiler) compileMultiplyUInt8(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.UInt8":
		return newInlineGoValue(
			fmt.Sprintf("(%s * %s)\n", left.value(), right.value()),
			left.elkType,
			"value.UInt8",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s * %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.UInt8",
	)
}

func (c *GoCompiler) compileSubtract(node *ast.BinaryExpressionNode) *goValue {
	left := c.compileExpression(node.Left)
	right := c.compileExpression(node.Right)
	typ := c.typeOf(node)

	switch left.goType() {
	case "value.SmallInt":
		return c.compileSubtractSmallInt(left, right, typ)
	case "*value.BigInt":
		return c.compileSubtractBigInt(left, right, typ)
	case "value.Int64":
		return c.compileSubtractInt64(left, right)
	case "value.Int32":
		return c.compileSubtractInt32(left, right)
	case "value.Int16":
		return c.compileSubtractInt16(left, right)
	case "value.Int8":
		return c.compileSubtractInt8(left, right)
	case "value.UInt64":
		return c.compileSubtractUInt64(left, right)
	case "value.UInt32":
		return c.compileSubtractUInt32(left, right)
	case "value.UInt16":
		return c.compileSubtractUInt16(left, right)
	case "value.UInt8":
		return c.compileSubtractUInt8(left, right)
	case "value.Float":
		return c.compileSubtractFloat(left, right, typ)
	case "value.Float64":
		return c.compileSubtractFloat64(left, right)
	case "value.Float32":
		return c.compileSubtractFloat32(left, right)
	case "*value.BigFloat":
		return c.compileSubtractBigFloat(left, right, typ)
	}

	if c.checker.IsSubtype(left.elkType, c.checker.Std(symbol.S_BuiltinSubtractable)) {
		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emit("%s, err = value.SubtractVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
		c.emitErrorPropagation()

		return newTmpGoValue(
			tmp,
			typ,
		)
	}

	callCache := c.emitCallCache()
	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitAddCallFrame(node.Location())
	c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpSubtract, &%s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
	c.emitPopCallFrame()
	c.emitErrorPropagation()
	return newTmpGoValue(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileSubtractBigInt(left, right *goValue, typ types.Type) *goValue {
	switch right.goType() {
	case "value.SmallInt":
		return newInlineGoValue(
			fmt.Sprintf("(%s).SubtractSmallInt(%s)\n", left.value(), right.value()),
			left.elkType,
			goValueType,
		)
	case "value.Float":
		return newInlineGoValue(
			fmt.Sprintf("(%s).SubtractFloat(%s)\n", left.value(), right.value()),
			right.elkType,
			"value.Float",
		)
	case "*value.BigFloat":
		return newInlineGoValue(
			fmt.Sprintf("(%s).SubtractBigFloat(%s)\n", left.value(), right.value()),
			right.elkType,
			"*value.BigFloat",
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newInlineGoValue(
			fmt.Sprintf("(%s).SubtractInt(%s)\n", left.value(), c.convertToValue(right)),
			left.elkType,
			goValueType,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emit("%s, err = (%s).SubtractVal(%s)\n", tmp, left.value(), c.convertToValue(right))
	c.emitErrorPropagation()

	return newTmpGoValue(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileSubtractSmallInt(left, right *goValue, typ types.Type) *goValue {
	switch right.goType() {
	case "value.SmallInt":
		return newInlineGoValue(
			fmt.Sprintf("(%s).SubtractSmallInt(%s)\n", left.value(), right.value()),
			left.elkType,
			"value.SmallInt",
		)
	case "value.Float":
		return newInlineGoValue(
			fmt.Sprintf("(%s).SubtractFloat(%s)\n", left.value(), right.value()),
			right.elkType,
			"value.Float",
		)
	case "*value.BigFloat":
		return newInlineGoValue(
			fmt.Sprintf("(%s).SubtractBigFloat(%s)\n", left.value(), right.value()),
			right.elkType,
			"*value.BigFloat",
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newInlineGoValue(
			fmt.Sprintf("(%s).SubtractInt(%s)\n", left.value(), c.convertToValue(right)),
			left.elkType,
			goValueType,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emit("%s, err = (%s).SubtractVal(%s)\n", tmp, left.value(), c.convertToValue(right))
	c.emitErrorPropagation()

	return newTmpGoValue(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileSubtractFloat(left, right *goValue, typ types.Type) *goValue {
	switch right.goType() {
	case "value.SmallInt":
		return newInlineGoValue(
			fmt.Sprintf("(%s).SubtractSmallInt(%s)\n", left.value(), right.value()),
			left.elkType,
			"value.Float",
		)
	case "value.Float":
		return newInlineGoValue(
			fmt.Sprintf("(%s).SubtractFloat(%s)\n", left.value(), right.value()),
			left.elkType,
			"value.Float",
		)
	case "*value.BigFloat":
		return newInlineGoValue(
			fmt.Sprintf("(%s).SubtractBigFloat(%s)\n", left.value(), right.value()),
			right.elkType,
			"*value.BigFloat",
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newInlineGoValue(
			fmt.Sprintf("(%s).SubtractInt(%s)\n", left.value(), c.convertToValue(right)),
			left.elkType,
			"value.Float",
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emit("%s, err = (%s).SubtractVal(%s)\n", tmp, left.value(), c.convertToValue(right))
	c.emitErrorPropagation()

	return newTmpGoValue(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileSubtractBigFloat(left, right *goValue, typ types.Type) *goValue {
	switch right.goType() {
	case "value.SmallInt":
		return newInlineGoValue(
			fmt.Sprintf("(%s).SubtractSmallInt(%s)\n", left.value(), right.value()),
			left.elkType,
			"*value.BigFloat",
		)
	case "value.Float":
		return newInlineGoValue(
			fmt.Sprintf("(%s).SubtractFloat(%s)\n", left.value(), right.value()),
			left.elkType,
			"*value.BigFloat",
		)
	case "*value.BigFloat":
		return newInlineGoValue(
			fmt.Sprintf("(%s).SubtractBigFloat(%s)\n", left.value(), right.value()),
			left.elkType,
			"*value.BigFloat",
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newInlineGoValue(
			fmt.Sprintf("(%s).SubtractInt(%s)\n", left.value(), c.convertToValue(right)),
			left.elkType,
			"*value.BigFloat",
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emit("%s, err = (%s).SubtractVal(%s)\n", tmp, left.value(), c.convertToValue(right))
	c.emitErrorPropagation()

	return newTmpGoValue(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileSubtractFloat64(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.Float64":
		return newInlineGoValue(
			fmt.Sprintf("(%s - %s)\n", left.value(), right.value()),
			left.elkType,
			"value.Float64",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s - %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.Float64",
	)
}

func (c *GoCompiler) compileSubtractFloat32(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.Float32":
		return newInlineGoValue(
			fmt.Sprintf("(%s - %s)\n", left.value(), right.value()),
			left.elkType,
			"value.Float32",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s - %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.Float32",
	)
}

func (c *GoCompiler) compileSubtractInt64(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.Int64":
		return newInlineGoValue(
			fmt.Sprintf("(%s - %s)\n", left.value(), right.value()),
			left.elkType,
			"value.Int64",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s - %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.Int64",
	)
}

func (c *GoCompiler) compileSubtractInt32(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.Int32":
		return newInlineGoValue(
			fmt.Sprintf("(%s - %s)\n", left.value(), right.value()),
			left.elkType,
			"value.Int32",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s - %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.Int32",
	)
}

func (c *GoCompiler) compileSubtractInt16(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.Int16":
		return newInlineGoValue(
			fmt.Sprintf("(%s - %s)\n", left.value(), right.value()),
			left.elkType,
			"value.Int16",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s - %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.Int16",
	)
}

func (c *GoCompiler) compileSubtractInt8(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.Int8":
		return newInlineGoValue(
			fmt.Sprintf("(%s - %s)\n", left.value(), right.value()),
			left.elkType,
			"value.Int8",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s - %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.Int8",
	)
}

func (c *GoCompiler) compileSubtractUInt64(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.UInt64":
		return newInlineGoValue(
			fmt.Sprintf("(%s - %s)\n", left.value(), right.value()),
			left.elkType,
			"value.UInt64",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s - %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.UInt64",
	)
}

func (c *GoCompiler) compileSubtractUInt32(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.UInt32":
		return newInlineGoValue(
			fmt.Sprintf("(%s - %s)\n", left.value(), right.value()),
			left.elkType,
			"value.UInt32",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s - %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.UInt32",
	)
}

func (c *GoCompiler) compileSubtractUInt16(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.UInt16":
		return newInlineGoValue(
			fmt.Sprintf("(%s - %s)\n", left.value(), right.value()),
			left.elkType,
			"value.UInt16",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s - %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.UInt16",
	)
}

func (c *GoCompiler) compileSubtractUInt8(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.UInt8":
		return newInlineGoValue(
			fmt.Sprintf("(%s - %s)\n", left.value(), right.value()),
			left.elkType,
			"value.UInt8",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s - %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.UInt8",
	)
}

func (c *GoCompiler) compileAdd(node *ast.BinaryExpressionNode) *goValue {
	left := c.compileExpression(node.Left)
	right := c.compileExpression(node.Right)
	typ := c.typeOf(node)

	switch left.goType() {
	case "value.SmallInt":
		return c.compileAddSmallInt(left, right, typ)
	case "*value.BigInt":
		return c.compileAddBigInt(left, right, typ)
	case "value.Int64":
		return c.compileAddInt64(left, right)
	case "value.Int32":
		return c.compileAddInt32(left, right)
	case "value.Int16":
		return c.compileAddInt16(left, right)
	case "value.Int8":
		return c.compileAddInt8(left, right)
	case "value.UInt64":
		return c.compileAddUInt64(left, right)
	case "value.UInt32":
		return c.compileAddUInt32(left, right)
	case "value.UInt16":
		return c.compileAddUInt16(left, right)
	case "value.UInt8":
		return c.compileAddUInt8(left, right)
	case "value.Float":
		return c.compileAddFloat(left, right, typ)
	case "value.Float64":
		return c.compileAddFloat64(left, right)
	case "value.Float32":
		return c.compileAddFloat32(left, right)
	case "*value.BigFloat":
		return c.compileAddBigFloat(left, right, typ)
	}

	if c.checker.IsSubtype(left.elkType, c.checker.Std(symbol.S_BuiltinAddable)) {
		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emit("%s, err = value.AddVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
		c.emitErrorPropagation()

		return newTmpGoValue(
			tmp,
			typ,
		)
	}

	callCache := c.emitCallCache()
	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitAddCallFrame(node.Location())
	c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpAdd, &%s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
	c.emitPopCallFrame()
	c.emitErrorPropagation()
	return newTmpGoValue(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileAddBigInt(left, right *goValue, typ types.Type) *goValue {
	switch right.goType() {
	case "value.SmallInt":
		return newInlineGoValue(
			fmt.Sprintf("(%s).AddSmallInt(%s)\n", left.value(), right.value()),
			left.elkType,
			goValueType,
		)
	case "value.Float":
		return newInlineGoValue(
			fmt.Sprintf("(%s).AddFloat(%s)\n", left.value(), right.value()),
			right.elkType,
			"value.Float",
		)
	case "*value.BigFloat":
		return newInlineGoValue(
			fmt.Sprintf("(%s).AddBigFloat(%s)\n", left.value(), right.value()),
			right.elkType,
			"*value.BigFloat",
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newInlineGoValue(
			fmt.Sprintf("(%s).AddInt(%s)\n", left.value(), c.convertToValue(right)),
			left.elkType,
			goValueType,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emit("%s, err = (%s).AddVal(%s)\n", tmp, left.value(), c.convertToValue(right))
	c.emitErrorPropagation()

	return newTmpGoValue(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileAddSmallInt(left, right *goValue, typ types.Type) *goValue {
	switch right.goType() {
	case "value.SmallInt":
		return newInlineGoValue(
			fmt.Sprintf("(%s).AddSmallInt(%s)\n", left.value(), right.value()),
			left.elkType,
			"value.SmallInt",
		)
	case "value.Float":
		return newInlineGoValue(
			fmt.Sprintf("(%s).AddFloat(%s)\n", left.value(), right.value()),
			right.elkType,
			"value.Float",
		)
	case "*value.BigFloat":
		return newInlineGoValue(
			fmt.Sprintf("(%s).AddBigFloat(%s)\n", left.value(), right.value()),
			right.elkType,
			"*value.BigFloat",
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newInlineGoValue(
			fmt.Sprintf("(%s).AddInt(%s)\n", left.value(), c.convertToValue(right)),
			left.elkType,
			goValueType,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emit("%s, err = (%s).AddVal(%s)\n", tmp, left.value(), c.convertToValue(right))
	c.emitErrorPropagation()

	return newTmpGoValue(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileAddFloat(left, right *goValue, typ types.Type) *goValue {
	switch right.goType() {
	case "value.SmallInt":
		return newInlineGoValue(
			fmt.Sprintf("(%s).AddSmallInt(%s)\n", left.value(), right.value()),
			left.elkType,
			"value.Float",
		)
	case "value.Float":
		return newInlineGoValue(
			fmt.Sprintf("(%s).AddFloat(%s)\n", left.value(), right.value()),
			right.elkType,
			"value.Float",
		)
	case "*value.BigFloat":
		return newInlineGoValue(
			fmt.Sprintf("(%s).AddBigFloat(%s)\n", left.value(), right.value()),
			right.elkType,
			"*value.BigFloat",
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newInlineGoValue(
			fmt.Sprintf("(%s).AddInt(%s)\n", left.value(), c.convertToValue(right)),
			left.elkType,
			"value.Float",
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emit("%s, err = (%s).AddVal(%s)\n", tmp, left.value(), c.convertToValue(right))
	c.emitErrorPropagation()

	return newTmpGoValue(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileAddBigFloat(left, right *goValue, typ types.Type) *goValue {
	switch right.goType() {
	case "value.SmallInt":
		return newInlineGoValue(
			fmt.Sprintf("(%s).AddSmallInt(%s)\n", left.value(), right.value()),
			left.elkType,
			"*value.BigFloat",
		)
	case "value.Float":
		return newInlineGoValue(
			fmt.Sprintf("(%s).AddFloat(%s)\n", left.value(), right.value()),
			left.elkType,
			"*value.BigFloat",
		)
	case "*value.BigFloat":
		return newInlineGoValue(
			fmt.Sprintf("(%s).AddBigFloat(%s)\n", left.value(), right.value()),
			left.elkType,
			"*value.BigFloat",
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newInlineGoValue(
			fmt.Sprintf("(%s).AddInt(%s)\n", left.value(), c.convertToValue(right)),
			left.elkType,
			"*value.BigFloat",
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emit("%s, err = (%s).AddVal(%s)\n", tmp, left.value(), c.convertToValue(right))
	c.emitErrorPropagation()

	return newTmpGoValue(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileAddFloat64(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.Float64":
		return newInlineGoValue(
			fmt.Sprintf("(%s + %s)\n", left.value(), right.value()),
			left.elkType,
			"value.Float64",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s + %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.Float64",
	)
}

func (c *GoCompiler) compileAddFloat32(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.Float32":
		return newInlineGoValue(
			fmt.Sprintf("(%s + %s)\n", left.value(), right.value()),
			left.elkType,
			"value.Float32",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s + %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.Float32",
	)
}

func (c *GoCompiler) compileAddInt64(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.Int64":
		return newInlineGoValue(
			fmt.Sprintf("(%s + %s)\n", left.value(), right.value()),
			left.elkType,
			"value.Int64",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s + %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.Int64",
	)
}

func (c *GoCompiler) compileAddInt32(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.Int32":
		return newInlineGoValue(
			fmt.Sprintf("(%s + %s)\n", left.value(), right.value()),
			left.elkType,
			"value.Int32",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s + %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.Int32",
	)
}

func (c *GoCompiler) compileAddInt16(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.Int16":
		return newInlineGoValue(
			fmt.Sprintf("(%s + %s)\n", left.value(), right.value()),
			left.elkType,
			"value.Int16",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s + %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.Int16",
	)
}

func (c *GoCompiler) compileAddInt8(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.Int8":
		return newInlineGoValue(
			fmt.Sprintf("(%s + %s)\n", left.value(), right.value()),
			left.elkType,
			"value.Int8",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s + %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.Int8",
	)
}

func (c *GoCompiler) compileAddUInt64(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.UInt64":
		return newInlineGoValue(
			fmt.Sprintf("(%s + %s)\n", left.value(), right.value()),
			left.elkType,
			"value.UInt64",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s + %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.UInt64",
	)
}

func (c *GoCompiler) compileAddUInt32(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.UInt32":
		return newInlineGoValue(
			fmt.Sprintf("(%s + %s)\n", left.value(), right.value()),
			left.elkType,
			"value.UInt32",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s + %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.UInt32",
	)
}

func (c *GoCompiler) compileAddUInt16(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.UInt16":
		return newInlineGoValue(
			fmt.Sprintf("(%s + %s)\n", left.value(), right.value()),
			left.elkType,
			"value.UInt16",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s + %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.UInt16",
	)
}

func (c *GoCompiler) compileAddUInt8(left, right *goValue) *goValue {
	switch right.goType() {
	case "value.UInt8":
		return newInlineGoValue(
			fmt.Sprintf("(%s + %s)\n", left.value(), right.value()),
			left.elkType,
			"value.UInt8",
		)
	}

	return newInlineGoValue(
		fmt.Sprintf("(%s + %s)\n", left.value(), c.valueToNarrowerType(right)),
		left.elkType,
		"value.UInt8",
	)
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
			return newInlineGoValue(
				fmt.Sprintf("value.Int64(%d)", v),
				c.checker.Std(symbol.Int64),
				"value.Int64",
			)
		case value.UInt64:
			return newInlineGoValue(
				fmt.Sprintf("value.UInt64(%d)", v),
				c.checker.Std(symbol.UInt64),
				"value.UInt64",
			)
		default:
			return nil
		}
	}

	switch val.ValueFlag() {
	case value.TRUE_FLAG:
		return newInlineGoValue(
			"true",
			types.Bool{},
			"bool",
		)
	case value.FALSE_FLAG:
		return newInlineGoValue(
			"false",
			types.Bool{},
			"bool",
		)
	case value.NIL_FLAG:
		return nilGoValue
	case value.SMALL_INT_FLAG:
		return newInlineGoValue(
			fmt.Sprintf("value.SmallInt(%d)", val.AsSmallInt()),
			c.checker.StdInt(),
			"value.SmallInt",
		)
	case value.INT64_FLAG:
		return newInlineGoValue(
			fmt.Sprintf("value.Int64(%d)", val.AsInt64()),
			c.checker.Std(symbol.Int64),
			"value.Int64",
		)
	case value.UINT64_FLAG:
		return newInlineGoValue(
			fmt.Sprintf("value.UInt64(%d)", val.AsUInt64()),
			c.checker.Std(symbol.UInt64),
			"value.UInt64",
		)
	case value.INT32_FLAG:
		return newInlineGoValue(
			fmt.Sprintf("value.Int32(%d)", val.AsInt32()),
			c.checker.Std(symbol.Int32),
			"value.Int32",
		)
	case value.UINT32_FLAG:
		return newInlineGoValue(
			fmt.Sprintf("value.UInt32(%d)", val.AsUInt32()),
			c.checker.Std(symbol.UInt32),
			"value.UInt32",
		)
	case value.INT16_FLAG:
		return newInlineGoValue(
			fmt.Sprintf("value.Int16(%d)", val.AsInt16()),
			c.checker.Std(symbol.Int16),
			"value.Int16",
		)
	case value.UINT16_FLAG:
		return newInlineGoValue(
			fmt.Sprintf("value.UInt16(%d)", val.AsUInt16()),
			c.checker.Std(symbol.UInt16),
			"value.UInt16",
		)
	case value.INT8_FLAG:
		return newInlineGoValue(
			fmt.Sprintf("value.Int8(%d)", val.AsInt8()),
			c.checker.Std(symbol.Int8),
			"value.Int8",
		)
	case value.UINT8_FLAG:
		return newInlineGoValue(
			fmt.Sprintf("value.UInt8(%d)", val.AsUInt8()),
			c.checker.Std(symbol.UInt8),
			"value.UInt8",
		)
	case value.CHAR_FLAG:
		return newInlineGoValue(
			fmt.Sprintf("value.Char(%q)", val.AsChar()),
			c.checker.Std(symbol.Char),
			"value.Char",
		)
	case value.FLOAT_FLAG:
		return newInlineGoValue(
			fmt.Sprintf("value.Float(%f)", val.AsFloat()),
			c.checker.Std(symbol.Float),
			"value.Float",
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
	return newInlineGoValue(
		buff.String(),
		c.checker.Std(symbol.ArrayList),
		"*value.ArrayList",
	)
}

func (c *GoCompiler) convertToValue(v *goValue) string {
	switch v.goType() {
	case goValueType:
		return v.value()
	case "value.SmallInt":
		return fmt.Sprintf("%s.ToValue()", v.value())
	case "*value.BigInt":
		return fmt.Sprintf("value.Ref(%s)", v.value())
	}

	if c.checker.IsSubtype(v.elkType, c.checker.StdInt()) {
		return v.value()
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Object)) {
		return fmt.Sprintf("value.Ref(%s)", v.value())
	}

	return fmt.Sprintf("%s.ToValue()", v.value())
}

func (c *GoCompiler) valueToNarrowerType(v *goValue) *goValue {
	if v.goType() != goValueType {
		return v
	}

	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.String)) {
		return newInlineGoValue(
			fmt.Sprintf("%s.AsString()", v.value()),
			v.elkType,
			"value.String",
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Char)) {
		return newInlineGoValue(
			fmt.Sprintf("%s.AsChar()", v.value()),
			v.elkType,
			"value.Char",
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Float)) {
		return newInlineGoValue(
			fmt.Sprintf("%s.AsFloat()", v.value()),
			v.elkType,
			"value.Float",
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Float64)) {
		return newInlineGoValue(
			fmt.Sprintf("%s.AsFloat64()", v.value()),
			v.elkType,
			"value.Float64",
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Float32)) {
		return newInlineGoValue(
			fmt.Sprintf("%s.AsFloat32()", v.value()),
			v.elkType,
			"value.Float32",
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.BigFloat)) {
		return newInlineGoValue(
			fmt.Sprintf("(*value.BigFloat)(%s.Pointer())", v.value()),
			v.elkType,
			"*value.BigFloat",
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Int64)) {
		return newInlineGoValue(
			fmt.Sprintf("%s.AsInt64()", v.value()),
			v.elkType,
			"value.Int64",
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Int32)) {
		return newInlineGoValue(
			fmt.Sprintf("%s.AsInt32()", v.value()),
			v.elkType,
			"value.Int32",
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Int16)) {
		return newInlineGoValue(
			fmt.Sprintf("%s.AsInt16()", v.value()),
			v.elkType,
			"value.Int16",
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Int8)) {
		return newInlineGoValue(
			fmt.Sprintf("%s.AsInt8()", v.value()),
			v.elkType,
			"value.Int8",
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.UInt)) {
		return newInlineGoValue(
			fmt.Sprintf("%s.AsUInt()", v.value()),
			v.elkType,
			"value.UInt",
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.UInt64)) {
		return newInlineGoValue(
			fmt.Sprintf("%s.AsUInt64()", v.value()),
			v.elkType,
			"value.UInt64",
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.UInt32)) {
		return newInlineGoValue(
			fmt.Sprintf("%s.AsUInt32()", v.value()),
			v.elkType,
			"value.UInt32",
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.UInt16)) {
		return newInlineGoValue(
			fmt.Sprintf("%s.AsUInt16()", v.value()),
			v.elkType,
			"value.UInt16",
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.UInt8)) {
		return newInlineGoValue(
			fmt.Sprintf("%s.AsUInt8()", v.value()),
			v.elkType,
			"value.UInt8",
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.ArrayList)) {
		return newInlineGoValue(
			fmt.Sprintf("(*value.ArrayList)(%s.Pointer())", v.value()),
			v.elkType,
			"*value.ArrayList",
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.ArrayTuple)) {
		return newInlineGoValue(
			fmt.Sprintf("(*value.ArrayTuple)(%s.Pointer())", v.value()),
			v.elkType,
			"*value.ArrayTuple",
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.HashMap)) {
		return newInlineGoValue(
			fmt.Sprintf("(*value.HashMap)(%s.Pointer())", v.value()),
			v.elkType,
			"*value.HashMap",
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.HashRecord)) {
		return newInlineGoValue(
			fmt.Sprintf("(*value.HashRecord)(%s.Pointer())", v.value()),
			v.elkType,
			"*value.HashRecord",
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.HashSet)) {
		return newInlineGoValue(
			fmt.Sprintf("(*value.HashSet)(%s.Pointer())", v.value()),
			v.elkType,
			"*value.HashSet",
		)
	}

	return v
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
	return newInlineGoValue(
		buff.String(),
		c.checker.Std(symbol.ArrayTuple),
		"*value.ArrayTuple",
	)
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
	return newInlineGoValue(
		buff.String(),
		c.checker.Std(symbol.HashSet),
		"*value.HashSet",
	)
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
	return newInlineGoValue(
		buff.String(),
		c.checker.Std(symbol.HashMap),
		"*value.HashMap",
	)
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
	return newInlineGoValue(
		buff.String(),
		c.checker.Std(symbol.HashRecord),
		"*value.HashRecord",
	)
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

	return newInlineGoValue(
		fmt.Sprintf("Pair{Key: %s, Value: %s}", k.inline, v.inline),
		c.checker.Std(symbol.Pair),
		"value.Pair",
	)
}

func (c *GoCompiler) enterScope(label string, typ nativeElkScopeType) {
	c.scopes = append(c.scopes, newNativeElkScope(label, typ))
}

func (c *GoCompiler) leaveScope() {
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
	goLocal := c.registerGoLocalWithComment(
		fmt.Sprintf("l%d", c.lastElkLocalIndex),
		goValueType,
		fmt.Sprintf("var %s: %s", name, types.Inspect(typ)),
	)
	newVar := &nativeElkLocal{
		name:    name,
		elkType: typ,
		goLocal: goLocal,
	}
	scope.localTable[name] = newVar

	return newVar
}
