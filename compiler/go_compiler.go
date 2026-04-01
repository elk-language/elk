package compiler

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/diagnostic"
	reflag "github.com/elk-language/elk/regex/flag"
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
func (m *GoSourceMethod) ToValue() value.Value                      { return value.Ref(m) }

type nativeBigInt struct {
	id  int
	val string
}

func (n *nativeBigInt) goIdent() string {
	return fmt.Sprintf("bi%d", n.id)
}

type nativeBigFloat struct {
	id  int
	val string
}

func (n *nativeBigFloat) goIdent() string {
	return fmt.Sprintf("bf%d", n.id)
}

type nativeSymbol struct {
	id  int
	val string
}

func (s *nativeSymbol) goIdent() string {
	return fmt.Sprintf("sym%d", s.id)
}

type nativeMethod struct {
	ident string
	init  string
}

// Returns true if the method takes a slice of args,
// false if it takes individual arguments with narrower types
func (n *nativeMethod) hasArgsSlice() bool {
	return n.init != ""
}

func (m *nativeMethod) goIdent() string {
	return m.ident
}

type nativeConstant struct {
	ident   string
	elkType types.Type
	goType  *value.GoType
	init    string
}

func (m *nativeConstant) goIdent() string {
	return m.ident
}

type nativeValue struct {
	ident   string
	val     value.Value
	goType  *value.GoType
	elkType types.Type
}

func (s *nativeValue) goIdent() string {
	return s.ident
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

type goLoopInfo struct {
	label                         string
	labelIsUsed                   bool
	resultVar                     *goLocal
	returnsValueFromLastIteration bool
}

var goValueType = value.FetchGoType("value.Value")

type goLocal struct {
	name       string
	goType     *value.GoType
	comment    string
	predefined bool
	free       bool
	elkLocal   bool
}

func (l *goLocal) markFree() {
	if l.elkLocal {
		return
	}

	l.free = true
}

func (l *goLocal) markReserved() {
	l.free = false
}

func newGoLocal(name string, goType *value.GoType, comment string) *goLocal {
	return &goLocal{
		name:    name,
		goType:  goType,
		comment: comment,
	}
}

type goValuePair struct {
	key   *goValue
	value *goValue
}

func (p goValuePair) markFree() {
	p.key.markFree()
	p.value.markFree()
}

type goValue struct {
	value   string
	goType  *value.GoType
	elkType types.Type
	locals  []*goLocal // temporary variables used in this expression
}

// Mark any temporary variables used in this expression as free
// so they can get used to hold other values
func (v *goValue) markFree() {
	for _, local := range v.locals {
		local.markFree()
	}
}

// Get the value and mark it as free.
// Should be used when emitting statements
// that use the value for the last time (it wont be used any longer).
//
// If the value will be used again or you're using it to create
// the expression of another `*goValue` use `value` property instead.
func (v *goValue) fetchValue() string {
	v.markFree()
	return v.value
}

// Create a new go value while inheriting locals of this one
func (v *goValue) newGoValue(value string, typ types.Type, goType *value.GoType) *goValue {
	return &goValue{
		value:   value,
		elkType: typ,
		goType:  goType,
		locals:  v.locals,
	}
}

func newGoValue(v string, typ types.Type, goType *value.GoType) *goValue {
	return &goValue{
		value:   v,
		elkType: typ,
		goType:  goType,
	}
}

// Create a new go value that contains the given dependencies.
// Inherits locals (temporary variables) from dependencies.
func newGoValueWithDependencies(v string, typ types.Type, goType *value.GoType, dependencies ...*goValue) *goValue {
	var locals []*goLocal
	for _, dependency := range dependencies {
		locals = append(locals, dependency.locals...)
	}

	return newGoValueWithLocals(v, typ, goType, locals...)
}

// Create a new go value that contains the local temporary variables.
func newGoValueWithLocals(v string, typ types.Type, goType *value.GoType, locals ...*goLocal) *goValue {
	return &goValue{
		value:   v,
		elkType: typ,
		goType:  goType,
		locals:  locals,
	}
}

// Create a new go value based on a Go local.
func newGoValueWithLocal(local *goLocal, typ types.Type) *goValue {
	return &goValue{
		value:   local.name,
		elkType: typ,
		goType:  local.goType,
		locals:  []*goLocal{local},
	}
}

var errGoValue = newGoValue("ERR", types.Untyped{}, goValueType)
var nilGoValue = newGoValue("value.Nil", types.Nil{}, goValueType)

type goImportEntry struct {
	path string
	name string
}

func newGoImportEntry(path, name string) *goImportEntry {
	return &goImportEntry{
		path: path,
		name: name,
	}
}

func CreateGoCompiler(parent *GoCompiler, checker types.Checker, loc *position.Location, errors *diagnostic.SyncDiagnosticList, output io.Writer) *GoCompiler {
	compiler := NewGoCompiler("main", topLevelGoCompilerMode, loc, checker, newGlobalData(), output)
	compiler.Errors = errors
	if parent != nil {
		compiler.SetParent(parent)
	}
	return compiler
}

func (c *GoCompiler) CreateMainCompiler(checker types.Checker, loc *position.Location, errors *diagnostic.SyncDiagnosticList, output io.Writer) Compiler {
	compiler := NewGoCompiler("main", topLevelGoCompilerMode, loc, checker, newGlobalData(), output)
	compiler.Errors = errors
	return compiler
}

func (c *GoCompiler) InitMainCompiler() {
	c.registerGoPackageClause("main")

	c.registerGoImport("github.com/elk-language/elk/value", "")
	c.registerGoImport("github.com/elk-language/elk/vm", "")
	c.registerGoImport("github.com/elk-language/elk/value/symbol", "")

	// noops to stop Go from complaining about unused imports
	c.emitPackage("var _ = symbol.Value\n")
	c.emitPackage("var _ = vm.New\n")
	c.emitPackage("var _ = value.Truthy\n\n")
}

func (c *GoCompiler) InitGlobalEnv() Compiler {
	envCompiler := NewGoCompiler("initGlobalEnv", topLevelGoCompilerMode, c.loc, c.checker, c.globalData, c.output)
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

	c.parent.emit("\ninitGlobalEnv()\n")
}

func (c *GoCompiler) InitMethodCompiler(location *position.Location) (Compiler, int) {
	methodCompiler := NewGoCompiler("methodDefinitions", topLevelGoCompilerMode, c.loc, c.checker, c.globalData, c.output)
	methodCompiler.Errors = c.Errors
	methodCompiler.SetParent(c)

	return methodCompiler, 0
}

func (c *GoCompiler) CompileMethods(location *position.Location, execOffset int) {
	c.registerGoLocal("class", value.FetchGoType("*value.Class"))
	c.compileMethodsWithinModule(c.checker.Env().Root, location)

	if c.buff.Len() == 0 {
		return
	}

	c.parent.emit("\nmethodDefinitions()\n")

	var funcBuffer bytes.Buffer
	fmt.Fprintf(&funcBuffer, "func methodDefinitions() {\n")
	c.compileLocalsTo(&funcBuffer)
	c.emitPrependBytes(funcBuffer.Bytes())
	c.emit("}\n")
}

func (c *GoCompiler) InitIvarIndicesCompiler(location *position.Location) (Compiler, int) {
	ivarCompiler := NewGoCompiler("ivarIndices", topLevelGoCompilerMode, c.loc, c.checker, c.globalData, c.output)
	ivarCompiler.Errors = c.Errors
	ivarCompiler.SetParent(c)

	return ivarCompiler, 0
}

func (c *GoCompiler) FinishIvarIndicesCompiler(location *position.Location, execOffset int) Compiler {
	if c.buff.Len() == 0 {
		return c.parent
	}

	if c.parent != nil {
		c.parent.emit("\nivarIndices(thread)\n")
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

	c.emit("\n")
	switch n := namespace.(type) {
	case *types.SingletonClass:
		attachedObjectConst := c.emitGetConst(value.ToSymbol(n.AttachedObject.Name()), types.Any{})
		c.emit("namespace = value.Ref((%s).SingletonClass())\n", attachedObjectConst.fetchValue())
	default:
		namespaceConst := c.emitGetConst(value.ToSymbol(n.Name()), types.Any{})
		c.emit("namespace = %s\n", c.convertToValue(namespaceConst).fetchValue())
	}

	init := c.valueToNarrowerType(
		c.compileExpression(node.Initialiser, false),
	)

	fullConstName := c.getFullConstName(namespace.Name(), constName.String())
	elkType := c.typeOf(node)
	goType := init.goType
	goIdent := mangleGoIdentifier(fullConstName)
	c.globalData.constantCache.SetUnsafe(
		fullConstName,
		&nativeConstant{
			ident:   goIdent,
			elkType: elkType,
			goType:  goType,
		},
	)

	constNameSymbol := c.emitSymbol(constName.String())
	c.emitPackage("var %s %s\n", goIdent, goType)
	c.emit("value.AddConstant(namespace, %s, %s)\n", constNameSymbol, c.convertToValue(init).fetchValue())
	c.emit("%s = %s\n", goIdent, init.fetchValue())
}

func (c *GoCompiler) getFullConstName(namespaceName, constName string) string {
	if namespaceName == "Root" {
		return constName
	}

	return fmt.Sprintf("%s::%s", namespaceName, constName)
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

	method := c.typeOf(node).(*types.Method)

	goName := c.registerElkMethodName(method.NamespacedName())

	methodCompiler := NewGoCompiler(goName, mode, node.Location(), c.checker, c.globalData, c.output)
	methodCompiler.isGenerator = node.IsGenerator()
	methodCompiler.isAsync = node.IsAsync()
	methodCompiler.Errors = c.Errors
	methodCompiler.method = method
	methodCompiler.compileMethodBody(node.Parameters, node.Body, method.ReturnType, node.Location())

	return methodCompiler
}

type globalData struct {
	bigFloatCache *concurrent.Map[string, *nativeBigFloat]
	bigIntCache   *concurrent.Map[string, *nativeBigInt]
	symbolCache   *concurrent.Map[string, *nativeSymbol]
	valueCache    *concurrent.Map[string, *nativeValue]
	methodCache   *concurrent.OrderedMap[string, *nativeMethod]
	constantCache *concurrent.Map[string, *nativeConstant]
	goImports     *concurrent.OrderedMap[string, *goImportEntry]
}

func newGlobalData() *globalData {
	return &globalData{
		bigFloatCache: concurrent.NewMap[string, *nativeBigFloat](),
		bigIntCache:   concurrent.NewMap[string, *nativeBigInt](),
		symbolCache:   concurrent.NewMap[string, *nativeSymbol](),
		valueCache:    concurrent.NewMap[string, *nativeValue](),
		methodCache:   concurrent.NewOrderedMap[string, *nativeMethod](),
		constantCache: concurrent.NewMap[string, *nativeConstant](),
		goImports:     concurrent.NewOrderedMap[string, *goImportEntry](),
	}
}

// Compiles Elk source code to Go source code.
type GoCompiler struct {
	Errors            *diagnostic.SyncDiagnosticList
	FuncName          string
	method            *types.Method
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
	loopInfo          []*goLoopInfo
	loopCounter       int // increments when a loop is entered, does not decrement ever
	tmpLocalCounter   int
	lastElkLocalIndex int
	callCacheCounter  int
	currentLineNumber int
	globalData        *globalData
	isGenerator       bool
	isAsync           bool
	unhygienic        bool
}

func NewGoCompiler(name string, mode goMode, loc *position.Location, checker types.Checker, globalData *globalData, output io.Writer) *GoCompiler {
	return &GoCompiler{
		FuncName:          name,
		mode:              mode,
		Errors:            diagnostic.NewSyncDiagnosticList(),
		scopes:            nativeElkScopes{newNativeElkScope("", defaultNativeElkScopeType)}, // start with an empty set for the 0th scope
		goLocals:          ds.MakeOrderedMap[string, *goLocal](),
		lastElkLocalIndex: -1,
		globalData:        globalData,
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

// Write the accumulated source code to the main buffer.
func (c *GoCompiler) Flush() {
	c.flushImport()
	c.output.Write([]byte("\n"))
	c.flushPackage()
	c.output.Write([]byte("\n"))
	c.flushInner()
}

func (c *GoCompiler) flushImport() {
	imports := c.globalData.goImports
	if imports.Len() == 0 {
		return
	}

	if packageEntry, ok := imports.GetUnsafe(""); ok {
		// special case, package name
		fmt.Fprintf(c.output, "package %s\n\n", packageEntry.name)
	}

	fmt.Fprint(c.output, "import (")
	for _, importEntry := range imports.Map.All() {
		if importEntry.path == "" {
			continue
		}

		if importEntry.name != "" {
			fmt.Fprintf(c.output, "%s ", importEntry.name)
		}
		fmt.Fprintf(c.output, "%q\n", importEntry.path)
	}
	fmt.Fprint(c.output, ")\n")

	imports.ClearUnsafe()
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

func (c *GoCompiler) addLoopInfo(label string, resultVar *goLocal, returnsValFromLastIteration bool) *goLoopInfo {
	info := &goLoopInfo{
		label:                         label,
		resultVar:                     resultVar,
		returnsValueFromLastIteration: returnsValFromLastIteration,
	}
	c.loopInfo = append(
		c.loopInfo,
		info,
	)

	return info
}

func (c *GoCompiler) findLoopInfo(label string, location *position.Location) *goLoopInfo {
	if len(c.loopInfo) < 1 {
		c.Errors.AddFailure(
			"cannot jump with `break` or `continue` outside of a loop",
			location,
		)
		return nil
	}

	if label == "" {
		// if there is no label, choose the closest enclosing loop
		return c.loopInfo[len(c.loopInfo)-1]
	}

	for _, currentJumpSet := range c.loopInfo {
		if currentJumpSet.label == label {
			return currentJumpSet
		}
	}

	c.Errors.AddFailure(
		fmt.Sprintf("label $%s does not exist or is not attached to an enclosing loop", label),
		location,
	)
	return nil
}

func (c *GoCompiler) popLoopInfo() {
	c.loopInfo = c.loopInfo[:len(c.loopInfo)-1]
}

func (c *GoCompiler) registerGoLocal(name string, goType *value.GoType) *goLocal {
	return c.registerGoLocalWithComment(name, goType, "")
}

func (c *GoCompiler) registerGoLocalWithComment(name string, goType *value.GoType, comment string) *goLocal {
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

	c.registerGoLocal("class", value.FetchGoType("*value.Class"))
	c.registerGoLocal("superclass", value.FetchGoType("*value.Class"))
	c.registerGoLocal("mixin", value.FetchGoType("*value.Mixin"))
}

// Entry point for compiling the body of a method.
func (c *GoCompiler) compileMethodBody(parameters []ast.ParameterNode, body []ast.StatementNode, returnType types.Type, loc *position.Location) {
	c.compileMethodFuncLiteralWithNativeArgsBody(parameters, body, returnType, loc)
	c.emitPackageBytes(c.buff.Bytes())
	c.buff.Reset()

	c.emit("func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {\n")

	c.emit("result, err := %s(thread, args[0]", c.FuncName)

	for i, param := range parameters {
		p := param.(*ast.MethodParameterNode)
		paramType := c.typeOf(p).(*types.Parameter)
		typ := paramType.Type

		argVal := newGoValue(
			fmt.Sprintf("args[%d]", i+1),
			typ,
			goValueType,
		)
		if p.Initialiser != nil {
			c.emit(", %s", argVal.value)
		} else {
			c.emit(", %s", c.valueToNarrowerType(argVal).value)
		}
	}

	c.emit(")\n")

	returnVal := newGoValue(
		"result",
		returnType,
		c.elkTypeToGoType(returnType, false),
	)
	c.emit("return %s, err", c.convertToValue(returnVal).fetchValue())

	c.emit("}")
}

func (c *GoCompiler) compileMethodFuncLiteralWithNativeArgsBody(parameters []ast.ParameterNode, body []ast.StatementNode, returnType types.Type, loc *position.Location) {
	var funcBuffer bytes.Buffer
	fmt.Fprintf(&funcBuffer, "func %s(thread *vm.Thread, self value.Value", c.FuncName)

	self := c.registerGoLocal("self", goValueType)
	self.predefined = true

	goReturnType := c.elkTypeToGoType(returnType, false)
	errVal := c.registerGoLocal("err", goValueType)
	errVal.predefined = true

	result := c.registerGoLocal("result", goReturnType)
	result.predefined = true

	for _, param := range parameters {
		p := param.(*ast.MethodParameterNode)
		pSpan := p.Location()

		pName := identifierToName(p.Name)
		paramType := c.typeOf(p).(*types.Parameter)
		typ := paramType.Type
		local := c.defineLocal(pName, typ, c.elkTypeToGoType(typ, false), pSpan)
		if local == nil {
			return
		}

		localName := local.goIdent()
		if p.Initialiser != nil {
			fmt.Fprintf(&funcBuffer, ", arg_%s value.Value", localName)

			c.emit("if (%s).IsUndefined() {\n", localName)
			val := c.compileExpression(p.Initialiser, false)
			c.emit("%s = %s\n", localName, c.valueToNarrowerType(val).fetchValue())
			c.emit("} else {\n")

			argVal := newGoValue(
				localName,
				typ,
				goValueType,
			)
			c.emit("%s = %s\n", localName, c.valueToNarrowerType(argVal).fetchValue())
			c.emit("}\n")
		} else {
			local.goLocal.predefined = true
			fmt.Fprintf(&funcBuffer, ", %s %s", localName, local.goLocal.goType.String())
		}

		if p.SetInstanceVariable {
			val := c.convertToValue(
				newGoValue(
					localName,
					local.elkType,
					local.goLocal.goType,
				),
			)
			c.emitSetInstanceVariable(value.ToSymbol(pName), val.fetchValue())
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

	c.emitAddCallFrame(loc)
	returnVal := c.compileStatements(body)

	c.emitReturn(c.valueToNarrowerType(returnVal).fetchValue())

	fmt.Fprintf(&funcBuffer, ") (result %s, err value.Value) { // method: %s, loc: %s\n", goReturnType, types.Inspect(c.method), loc.String())
	c.compileLocalsTo(&funcBuffer)
	c.emitPrependBytes(funcBuffer.Bytes())
	c.emit("\n}\n")

	// TODO: implement generators
	// c.emitFinalReturn(location, nil)
}

func (c *GoCompiler) compileMethodFuncLiteralBody(parameters []ast.ParameterNode, body []ast.StatementNode) {
	c.registerGoLocal("self", goValueType)
	c.emit("self = args[0]\n")

	for i, param := range parameters {
		p := param.(*ast.MethodParameterNode)
		pSpan := p.Location()

		pName := identifierToName(p.Name)
		paramType := c.typeOf(p).(*types.Parameter)
		typ := paramType.Type
		local := c.defineLocal(pName, typ, c.elkTypeToGoType(typ, false), pSpan)
		if local == nil {
			return
		}

		localName := local.goIdent()
		if p.Initialiser != nil {
			argVal := newGoValue(
				fmt.Sprintf("args[%d]", i+1),
				typ,
				goValueType,
			)
			c.emit("if (%s).IsUndefined() {\n", argVal.value)
			val := c.compileExpression(p.Initialiser, false)
			c.emit("%s = %s\n", localName, c.valueToNarrowerType(val).fetchValue())
			c.emit("} else {\n")

			c.emit("%s = %s\n", localName, c.valueToNarrowerType(argVal).fetchValue())

			c.emit("}\n")
		} else {
			argVal := newGoValue(
				fmt.Sprintf("args[%d]", i+1),
				typ,
				goValueType,
			)
			c.emit("%s = %s\n", localName, c.valueToNarrowerType(argVal).fetchValue())
		}

		if p.SetInstanceVariable {
			val := c.convertToValue(
				newGoValue(
					localName,
					local.elkType,
					local.goLocal.goType,
				),
			)
			c.emitSetInstanceVariable(value.ToSymbol(pName), val.fetchValue())
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

	c.emitReturn(c.convertToValue(val).fetchValue())
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

func (c *GoCompiler) compileMethodsWithinInterface(iface *types.Interface, location *position.Location) {
	singleton := iface.Singleton()
	if types.NamespaceHasAnyDefinableMethods(singleton) {
		ifaceVal := c.emitGetConst(value.ToSymbol(iface.Name()), c.checker.Std(symbol.Interface))
		c.emit("class = (%s).SingletonClass() // %s\n", ifaceVal.fetchValue(), iface.Name())

		for methodName, method := range types.SortedOwnMethods(singleton) {
			c.compileMethodDefinition(methodName, method, location)

			for i, overload := range method.Overloads {
				overloadName := value.ToSymbol(
					fmt.Sprintf("%s@%d", methodName.String(), i+1),
				)
				c.compileMethodDefinition(overloadName, overload, location)
			}
		}
	}

	for _, subtype := range types.SortedSubtypes(iface) {
		if subtype.Type == iface {
			continue
		}
		c.compileMethodsWithinType(subtype.Type, location)
	}
}

func (c *GoCompiler) compileMethodsWithinModule(module *types.Module, location *position.Location) {
	if types.NamespaceHasAnyDefinableMethods(module) {
		moduleVal := c.emitGetConst(value.ToSymbol(module.Name()), c.checker.Std(symbol.Module))
		c.emit("class = (%s).SingletonClass() // %s\n", moduleVal.fetchValue(), module.Name())

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

func (c *GoCompiler) compileMethodsWithinClassOrMixin(namespace types.Namespace, location *position.Location) {
	namespaceHasCompiledMethods := types.NamespaceHasAnyDefinableMethods(namespace)

	singleton := namespace.Singleton()
	singletonHasCompiledMethods := types.NamespaceHasAnyDefinableMethods(singleton)

	if namespaceHasCompiledMethods || singletonHasCompiledMethods {
		namespaceVal := c.emitGetConst(value.ToSymbol(namespace.Name()), c.checker.Std(symbol.Class))
		c.emit("class = %s // %s\n", c.valueToNarrowerType(namespaceVal).fetchValue(), namespace.Name())

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
		c.compileMethodsWithinClassOrMixin(t, location)
	case *types.Mixin:
		c.compileMethodsWithinClassOrMixin(t, location)
	case *types.Interface:
		c.compileMethodsWithinInterface(t, location)
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
			c.registerGoLocal("aliasClass", value.FetchGoType("*value.Class"))

			switch namespace.(type) {
			case *value.Class:
				classVal := c.emitGetConst(value.ToSymbol(method.DefinedUnder.Name()), c.checker.Std(symbol.Class))
				c.emit("aliasClass = %s\n", c.valueToNarrowerType(classVal).fetchValue())
			case *value.Module:
				moduleVal := c.emitGetConst(value.ToSymbol(method.DefinedUnder.Name()), c.checker.Std(symbol.Module))
				c.emit("aliasClass = (%s).SingletonClass()\n", moduleVal.fetchValue())
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

	methodCompiler := (*GoCompiler)(method.Body.(*GoSourceMethod))

	c.emit("vm.Def(&class.MethodContainer, %q, ", name.String())
	c.emitBytes(methodCompiler.buff.Bytes())
	methodCompiler.buff.Reset()

	c.emitPackageBytes(methodCompiler.packageBuff.Bytes())
	methodCompiler.packageBuff.Reset()
	c.emit(",")

	if len(method.Params) > 0 {
		c.emit("vm.DefWithParameters(%d), ", len(method.Params))
	}

	c.emit(")\n")

	method.SetCompiled(true)
	method.Body = nil
}

func (c *GoCompiler) compileNamespaceDefinition(parentNamespace, namespace types.Namespace, constName value.Symbol) {
	if !namespace.IsDefined() && !namespace.IsNative() {
		switch p := parentNamespace.(type) {
		case *types.SingletonClass:
			namespaceVal := c.emitGetConst(value.ToSymbol(p.AttachedObject.Name()), types.Any{})
			c.emit("parentNamespace = (%s).SingletonClass()\n", namespaceVal.fetchValue())
		default:
			namespaceVal := c.emitGetConst(value.ToSymbol(p.Name()), types.Any{})
			c.emit("parentNamespace = %s\n", c.convertToValue(namespaceVal).fetchValue())
		}

		goIdent := mangleGoIdentifier(constName.String())
		var elkType types.Type
		var goType *value.GoType

		switch namespace.(type) {
		case *types.Module:
			elkType = c.checker.Std(symbol.Module)
			goType = value.FetchGoType("*value.Module")
			c.emit("%s = value.NewModule()\n", goIdent)
			c.emit("namespace = value.Ref(%s)\n", goIdent)
		case *types.Class:
			elkType = c.checker.Std(symbol.Class)
			goType = value.FetchGoType("*value.Class")
			c.emit("%s = value.NewClassWithOptions(value.ClassWithSuperclass(nil))\n", goIdent)
			c.emit("namespace = value.Ref(%s)\n", goIdent)
		case *types.Mixin:
			elkType = c.checker.Std(symbol.Mixin)
			goType = value.FetchGoType("*value.Mixin")
			c.emit("%s = value.NewMixin()\n", goIdent)
			c.emit("%s = value.NewMixin()\n", goIdent)
		case *types.Interface:
			elkType = c.checker.Std(symbol.Interface)
			goType = value.FetchGoType("*value.Interface")
			c.emit("%s = value.NewInterface()\n", goIdent)
			c.emit("namespace = value.Ref(%s)\n", goIdent)
		}

		c.globalData.constantCache.SetUnsafe(
			constName.String(),
			&nativeConstant{
				ident:   goIdent,
				elkType: elkType,
				goType:  goType,
			},
		)
		c.emitPackage("var %s %s // %s\n", goIdent, goType, constName.String())
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

// Emit go package clause line
func (c *GoCompiler) registerGoPackageClause(name string) {
	c.registerGoImport("", name)
}

// Emit import level code
func (c *GoCompiler) registerGoImport(path, name string) {
	imports := c.globalData.goImports
	imports.Lock()
	defer imports.Unlock()

	if _, ok := imports.GetUnsafe(path); ok {
		return
	}

	imports.SetUnsafe(path, newGoImportEntry(path, name))
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

func (c *GoCompiler) emitGetConst(fullName value.Symbol, elkType types.Type) *goValue {
	if nativeConst, ok := value.NativeConstantMap[fullName.String()]; ok {
		return newGoValue(
			nativeConst.GoExpr,
			elkType,
			nativeConst.GoType,
		)
	}

	c.globalData.constantCache.Lock()
	defer c.globalData.constantCache.Unlock()

	fullNameString := fullName.String()
	if constant, ok := c.globalData.constantCache.GetUnsafe(fullNameString); ok {
		return newGoValue(
			constant.goIdent(),
			elkType,
			constant.goType,
		)
	}

	val := c.emitDynamicGetConst(fullName, elkType)
	goIdent := mangleGoIdentifier(fullNameString)
	c.emitPackage("var %s %s // %s\n", goIdent, val.goType, fullNameString)
	c.globalData.constantCache.SetUnsafe(
		fullNameString,
		&nativeConstant{
			ident:   goIdent,
			elkType: elkType,
			goType:  val.goType,
			init:    val.fetchValue(),
		},
	)

	return newGoValue(
		goIdent,
		elkType,
		val.goType,
	)
}

func (c *GoCompiler) emitDynamicGetConst(fullName value.Symbol, elkType types.Type) *goValue {
	constNameSymbol := c.emitSymbol(fullName.String())

	return c.valueToNarrowerType(
		newGoValue(
			fmt.Sprintf("value.GetConstant(%s)", constNameSymbol),
			elkType,
			goValueType,
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

	classVal := c.emitGetConst(value.ToSymbol(class.Name()), c.checker.Std(symbol.Class))
	c.emit("class = %s\n", c.valueToNarrowerType(classVal).fetchValue())

	superclassVal := c.emitGetConst(value.ToSymbol(superclass.Name()), c.checker.Std(symbol.Class))
	c.emit("superclass = %s\n", c.valueToNarrowerType(superclassVal).fetchValue())

	c.emit("class.SetSuperclass(superclass)\n")
}

func (c *GoCompiler) CompileIvarIndices(target types.NamespaceWithIvarIndices, location *position.Location) {
	c.registerGoLocal("class", value.FetchGoType("*value.Class"))

	switch target := target.(type) {
	case *types.SingletonClass:
		namespaceVal := c.emitGetConst(value.ToSymbol(target.AttachedObject.Name()), types.Any{})
		c.emit("class = (%s).SingletonClass()\n", namespaceVal.fetchValue())
	case *types.Module:
		namespaceVal := c.emitGetConst(value.ToSymbol(target.Name()), types.Any{})
		c.emit("class = (%s).SingletonClass()\n", namespaceVal.fetchValue())
	default:
		namespaceVal := c.emitGetConst(value.ToSymbol(target.Name()), c.checker.Std(symbol.Class))
		c.emit("class = %s\n", c.valueToNarrowerType(namespaceVal).fetchValue())
	}

	c.emit("class.IvarIndices = %s\n", target.IvarIndices().ToGoSource())
}

func (c *GoCompiler) CompileInclude(target types.Namespace, mixin *types.Mixin, location *position.Location) {
	switch t := target.(type) {
	case *types.SingletonClass:
		namespaceVal := c.emitGetConst(value.ToSymbol(t.AttachedObject.Name()), types.Any{})
		c.emit("class = (%s).SingletonClass()\n", namespaceVal.fetchValue())
	default:
		namespaceVal := c.emitGetConst(value.ToSymbol(target.Name()), c.checker.Std(symbol.Class))
		c.emit("class = %s\n", c.valueToNarrowerType(namespaceVal).fetchValue())
	}

	mixinVal := c.emitGetConst(value.ToSymbol(mixin.Name()), c.checker.Std(symbol.Mixin))
	c.emit("mixin = %s\n", c.valueToNarrowerType(mixinVal).fetchValue())

	c.emit("class.IncludeMixin(mixin)\n")
}

func isValidGoIdentRune(r rune, first bool) bool {
	return r == '_' || unicode.IsLetter(r) || !first && unicode.IsDigit(r)
}

type goMangleMapping struct {
	from string
	to   string
}

var goMangleMap = []goMangleMapping{
	{"===", "_eqq_"},
	{"<=>", "_cmp_"},
	{"<<<", "_llsh_"},
	{">>>", "_rrsh_"},
	{".:", "_im_"},
	{"::", "_ns_"},
	{">=", "_gte_"},
	{"<=", "_lte_"},
	{"<<", "_lsh_"},
	{">>", "_rsh_"},
	{"==", "_eq_"},
	{"=~", "_leq_"},
	{"&~", "_andnot_"},
	{"++", "_inc_"},
	{"--", "_dec_"},
	{"**", "_exp_"},
	{"+@", "_plus_"},
	{"-@", "_neg_"},
	{":", "_cln_"},
	{".", "_dot_"},
	{"@", "_at_"},
	{">", "_gt_"},
	{"<", "_lt_"},
	{"~", "_tld_"},
	{"&", "_and_"},
	{"|", "_or_"},
	{"^", "_xor_"},
	{"+", "_add_"},
	{"-", "_sub_"},
	{"*", "_mul_"},
	{"/", "_div_"},
	{"%", "_mod_"},
}

func mangleGoIdentifier(name string) string {
	for _, mapping := range goMangleMap {
		name = strings.ReplaceAll(name, mapping.from, mapping.to)
	}

	var b strings.Builder

	for i, r := range name {
		if isValidGoIdentRune(r, i == 0) {
			b.WriteRune(r)
			continue
		}

		if r>>8 == 0 {
			fmt.Fprintf(&b, `_x%02x_`, r)
		} else if r>>16 == 0 {
			fmt.Fprintf(&b, `_u%04x_`, r)
		} else {
			fmt.Fprintf(&b, `_U%08X_`, r)
		}
	}

	return b.String()
}

func mangleFileName(name string) string {
	return fmt.Sprintf("__file_%s", mangleGoIdentifier(name))
}

func (c *GoCompiler) InitExpressionCompiler(location *position.Location) Compiler {
	name := mangleFileName(location.FilePath)
	exprCompiler := NewGoCompiler(name, topLevelGoCompilerMode, location, c.checker, c.globalData, c.output)
	exprCompiler.SetParent(c)
	exprCompiler.Errors = c.Errors

	return exprCompiler
}

func (c *GoCompiler) CompileExpressionsInFile(node *ast.ProgramNode) {
	var initCode string
	if c.FuncName == "main" {
		initCode = c.buff.String()
		c.buff.Reset()
	}

	c.emitAddCallFrame(node.Location())
	c.compileProgram(node)

	if c.buff.Len() == 0 && c.FuncName != "main" {
		return
	}

	if c.parent != nil {
		c.parent.emit("%s(thread)\n", c.FuncName)
	}

	var funcBuffer bytes.Buffer
	if c.FuncName == "main" {
		var methodVarsBuff bytes.Buffer
		for _, nativeMethod := range c.globalData.methodCache.Map.All() {
			if nativeMethod.init == "" {
				continue
			}

			fmt.Fprintf(&methodVarsBuff, "%s = %s\n", nativeMethod.goIdent(), nativeMethod.init)
		}
		if methodVarsBuff.Len() > 0 {
			fmt.Fprintln(&methodVarsBuff)
			c.emitPrependBytes(methodVarsBuff.Bytes())
		}

		c.emitPrependBytes([]byte(initCode))
		fmt.Fprintf(&funcBuffer, "func %s() { // loc: %s\n", c.FuncName, c.loc.FilePath)
		fmt.Fprintf(&funcBuffer, "thread := vm.New()\n_ = thread\n")
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
		if local.predefined {
			continue
		}

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
	for i, stmt := range nodes {
		lastValue = c.compileStatement(stmt, i != len(nodes)-1)
		if types.IsNever(lastValue.elkType) {
			break
		}
	}

	if lastValue == nil {
		return nilGoValue
	}
	return lastValue
}

func (c *GoCompiler) compileStatement(node ast.StatementNode, valueIsIgnored bool) *goValue {
	switch node := node.(type) {
	case *ast.ExpressionStatementNode:
		return c.compileExpression(node.Expression, valueIsIgnored)
	default:
		return nilGoValue
	}
}

func (c *GoCompiler) compileExpression(node ast.ExpressionNode, valueIsIgnored bool) *goValue {
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
	case *ast.TypeofExpressionNode:
		return c.compileExpression(node.Value, valueIsIgnored)
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
	case *ast.ValueDeclarationNode:
		return c.compileValueDeclarationNode(node)
	case *ast.MethodCallNode:
		return c.compileMethodCallNode(node, valueIsIgnored)
	case *ast.GenericMethodCallNode:
		return c.compileGenericMethodCallNode(node, valueIsIgnored)
	case *ast.PublicConstantNode:
		return c.compilePublicConstantNode(node)
	case *ast.PrivateConstantNode:
		return c.compilePrivateConstantNode(node)
	case *ast.GenericConstantNode:
		return c.compileExpression(node.Constant, valueIsIgnored)
	case *ast.SelfLiteralNode:
		return newGoValue("self", c.checker.SelfType(), goValueType)
	case *ast.ReturnExpressionNode:
		return c.compileReturnExpressionNode(node)
	case *ast.IntLiteralNode:
		return c.compileIntLiteralNode(node)
	case *ast.Int8LiteralNode:
		return c.compileInt8LiteralNode(node)
	case *ast.Int16LiteralNode:
		return c.compileInt16LiteralNode(node)
	case *ast.Int32LiteralNode:
		return c.compileInt32LiteralNode(node)
	case *ast.Int64LiteralNode:
		return c.compileInt64LiteralNode(node)
	case *ast.UInt8LiteralNode:
		return c.compileUInt8LiteralNode(node)
	case *ast.UInt16LiteralNode:
		return c.compileUInt16LiteralNode(node)
	case *ast.UInt32LiteralNode:
		return c.compileUInt32LiteralNode(node)
	case *ast.UInt64LiteralNode:
		return c.compileUInt64LiteralNode(node)
	case *ast.UIntLiteralNode:
		return c.compileUIntLiteralNode(node)
	case *ast.BigFloatLiteralNode:
		return c.compileBigFloatLiteralNode(node)
	case *ast.FloatLiteralNode:
		return c.compileFloatLiteralNode(node)
	case *ast.Float64LiteralNode:
		return c.compileFloat64LiteralNode(node)
	case *ast.Float32LiteralNode:
		return c.compileFloat32LiteralNode(node)
	case *ast.RawStringLiteralNode:
		return c.compileRawStringLiteralNode(node)
	case *ast.DoubleQuotedStringLiteralNode:
		return c.compileDoubleQuotedStringLiteralNode(node)
	case *ast.InterpolatedStringLiteralNode:
		return c.compileInterpolatedStringLiteralNode(node)
	case *ast.InterpolatedSymbolLiteralNode:
		return c.compileInterpolatedSymbolLiteralNode(node)
	case *ast.UninterpolatedRegexLiteralNode:
		return c.compileUninterpolatedRegexLiteralNode(node)
	case *ast.InterpolatedRegexLiteralNode:
		return c.compileInterpolatedRegexLiteralNode(node)
	case *ast.RawCharLiteralNode:
		return c.compileRawCharLiteralNode(node)
	case *ast.CharLiteralNode:
		return c.compileCharLiteralNode(node)
	case *ast.SimpleSymbolLiteralNode:
		return c.compileSimpleSymbolLiteralNode(node)
	case *ast.NilLiteralNode:
		return nilGoValue
	case *ast.TrueLiteralNode:
		return newGoValue("value.True", types.Bool{}, value.FetchGoType("value.Bool"))
	case *ast.FalseLiteralNode:
		return newGoValue("value.False", types.Bool{}, value.FetchGoType("value.Bool"))
	case *ast.BinaryExpressionNode:
		return c.compileBinaryExpressionNode(node, valueIsIgnored)
	case *ast.ArrayTupleLiteralNode:
		return c.compileArrayTupleLiteralNode(node)
	case *ast.WordArrayTupleLiteralNode:
		return c.compileWordArrayTupleLiteralNode(node)
	case *ast.SymbolArrayTupleLiteralNode:
		return c.compileSymbolArrayTupleLiteralNode(node)
	case *ast.BinArrayTupleLiteralNode:
		return c.compileBinArrayTupleLiteralNode(node)
	case *ast.HexArrayTupleLiteralNode:
		return c.compileHexArrayTupleLiteralNode(node)
	case *ast.ArrayListLiteralNode:
		return c.compileArrayListLiteralNode(node)
	case *ast.WordArrayListLiteralNode:
		return c.compileWordArrayListLiteralNode(node)
	case *ast.SymbolArrayListLiteralNode:
		return c.compileSymbolArrayListLiteralNode(node)
	case *ast.BinArrayListLiteralNode:
		return c.compileBinArrayListLiteralNode(node)
	case *ast.HexArrayListLiteralNode:
		return c.compileHexArrayListLiteralNode(node)
	case *ast.HashSetLiteralNode:
		return c.compileHashSetLiteralNode(node)
	case *ast.WordHashSetLiteralNode:
		return c.compileWordHashSetLiteralNode(node)
	case *ast.SymbolHashSetLiteralNode:
		return c.compileSymbolHashSetLiteralNode(node)
	case *ast.BinHashSetLiteralNode:
		return c.compileBinHashSetLiteralNode(node)
	case *ast.HexHashSetLiteralNode:
		return c.compileHexHashSetLiteralNode(node)
	case *ast.HashMapLiteralNode:
		return c.compileHashMapLiteralNode(node)
	case *ast.HashRecordLiteralNode:
		return c.compileHashRecordLiteralNode(node)
	case *ast.AssignmentExpressionNode:
		return c.compileAssignmentExpressionNode(node)
	case *ast.PublicIdentifierNode:
		return c.compileLocalVariableAccess(node.Value, c.typeOf(node))
	case *ast.PrivateIdentifierNode:
		return c.compileLocalVariableAccess(node.Value, c.typeOf(node))
	case *ast.RangeLiteralNode:
		return c.compileRangeLiteralNode(node)
	case *ast.AsExpressionNode:
		return c.compileAsExpressionNode(node, valueIsIgnored)
	case *ast.MustExpressionNode:
		return c.compileMustExpressionNode(node, valueIsIgnored)
	case *ast.WhileExpressionNode:
		return c.compileWhileExpressionNode("", node, valueIsIgnored)
	case *ast.LoopExpressionNode:
		return c.compileLoopExpressionNode("", node, valueIsIgnored)
	case *ast.IfExpressionNode:
		return c.compileIfExpression(
			ifConditionType,
			node.Condition,
			node.ThenBody,
			node.ElseBody,
			c.typeOf(node),
			valueIsIgnored,
		)
	case *ast.UnlessExpressionNode:
		return c.compileIfExpression(
			unlessConditionType,
			node.Condition,
			node.ThenBody,
			node.ElseBody,
			c.typeOf(node),
			valueIsIgnored,
		)
	case *ast.ModifierIfElseNode:
		return c.compileModifierIfExpression(
			ifConditionType,
			node.Condition,
			node.ThenExpression,
			node.ElseExpression,
			c.typeOf(node),
			valueIsIgnored,
		)
	// case *ast.ForInExpressionNode:
	// return c.compileForInExpressionNode("", node)
	default:
		panic(fmt.Sprintf("invalid expression node: %T", node))
	}
}

// func (c *GoCompiler) compileForInExpressionNode(label string, node *ast.ForInExpressionNode) {
// 	return c.compileForIn()
// }

// func (c *GoCompiler) compileForIn(label string, param ast.PatternNode, inExpression ast.ExpressionNode, then func() *goValue, location *position.Location, collectionLiteral bool) *goValue {
// 	if result := c.compileForInAsNumericFor(label, param, inExpression, then, location, collectionLiteral); result != nil {
// 		return result
// 	}
// }

// func (c *GoCompiler) compileForInAsNumericFor(
// 	label string,
// 	param ast.PatternNode,
// 	inExpression ast.ExpressionNode,
// 	then func() *goValue,
// 	location *position.Location,
// 	collectionLiteral bool,
// ) bool {
// 	var paramExpr ast.ExpressionNode
// 	var paramName string
// 	switch p := param.(type) {
// 	case *ast.PublicIdentifierNode:
// 		paramExpr = p
// 		paramName = p.Value
// 	case *ast.PrivateIdentifierNode:
// 		paramExpr = p
// 		paramName = p.Value
// 	default:
// 		return false
// 	}

// 	switch in := inExpression.(type) {
// 	case *ast.RangeLiteralNode:
// 		return c.compileForInRangeLiteralAsNumericFor(label, in, then, paramExpr, paramName, collectionLiteral, location)
// 	case *ast.IntLiteralNode:
// 		return c.compileForInIntLiteralAsNumericFor(label, in, then, paramExpr, paramName, collectionLiteral, location)
// 	}

// 	inExpressionType := c.typeOf(inExpression)
// 	if c.checker.IsSubtype(inExpressionType, c.checker.Std(symbol.Range)) {
// 		return c.compileForInRangeAsNumericFor(label, inExpression, then, paramExpr, paramName, collectionLiteral, location)
// 	}
// 	if c.checker.IsSubtype(inExpressionType, c.checker.Std(symbol.Int)) {
// 		return c.compileForInIntAsNumericFor(label, inExpression, then, paramExpr, paramName, collectionLiteral, location)
// 	}

// 	return false
// }

func (c *GoCompiler) compileSimpleSymbolLiteralNode(node *ast.SimpleSymbolLiteralNode) *goValue {
	varName := c.emitSymbol(node.Content)
	return newGoValue(
		varName,
		c.typeOf(node),
		value.FetchGoType("value.Symbol"),
	)
}

func (c *GoCompiler) compileCharLiteralNode(node *ast.CharLiteralNode) *goValue {
	return newGoValue(
		fmt.Sprintf("value.Char(%q)", node.Value),
		c.typeOf(node),
		value.FetchGoType("value.Char"),
	)
}

func (c *GoCompiler) compileRawCharLiteralNode(node *ast.RawCharLiteralNode) *goValue {
	return newGoValue(
		fmt.Sprintf("value.Char(%q)", node.Value),
		c.typeOf(node),
		value.FetchGoType("value.Char"),
	)
}

func (c *GoCompiler) compileRangeLiteralNode(node *ast.RangeLiteralNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	if node.Start == nil {
		end := c.compileExpression(node.End, false)

		switch node.Op.Type {
		case token.CLOSED_RANGE_OP, token.LEFT_OPEN_RANGE_OP:
			end := c.convertToValue(end)
			return newGoValueWithDependencies(
				fmt.Sprintf("value.NewBeginlessClosedRange(%s)", end.value),
				c.typeOf(node),
				value.FetchGoType("*value.BeginlessClosedRange"),
				end,
			)
		case token.RIGHT_OPEN_RANGE_OP, token.OPEN_RANGE_OP:
			end := c.convertToValue(end)
			return newGoValueWithDependencies(
				fmt.Sprintf("value.NewBeginlessOpenRange(%s)", end.value),
				c.typeOf(node),
				value.FetchGoType("*value.BeginlessOpenRange"),
				end,
			)
		default:
			panic(fmt.Sprintf("invalid range operator: %#v", node.Op))
		}
	}
	if node.End == nil {
		start := c.compileExpression(node.Start, false)

		switch node.Op.Type {
		case token.CLOSED_RANGE_OP, token.RIGHT_OPEN_RANGE_OP:
			start := c.convertToValue(start)
			return newGoValueWithDependencies(
				fmt.Sprintf("value.NewEndlessClosedRange(%s)", start.value),
				c.typeOf(node),
				value.FetchGoType("*value.EndlessClosedRange"),
				start,
			)
		case token.LEFT_OPEN_RANGE_OP, token.OPEN_RANGE_OP:
			start := c.convertToValue(start)
			return newGoValueWithDependencies(
				fmt.Sprintf("value.NewEndlessOpenRange(%s)", start.value),
				c.typeOf(node),
				value.FetchGoType("*value.EndlessOpenRange"),
				start,
			)
		default:
			panic(fmt.Sprintf("invalid range operator: %#v", node.Op))
		}
	}

	start := c.compileExpression(node.Start, false)
	end := c.compileExpression(node.End, false)
	switch node.Op.Type {
	case token.CLOSED_RANGE_OP:
		start := c.convertToValue(start)
		end := c.convertToValue(end)
		return newGoValueWithDependencies(
			fmt.Sprintf("value.NewClosedRange(%s, %s)", start.value, end.value),
			c.typeOf(node),
			value.FetchGoType("*value.ClosedRange"),
			start, end,
		)
	case token.OPEN_RANGE_OP:
		start := c.convertToValue(start)
		end := c.convertToValue(end)
		return newGoValueWithDependencies(
			fmt.Sprintf("value.NewOpenRange(%s, %s)", start.value, end.value),
			c.typeOf(node),
			value.FetchGoType("*value.OpenRange"),
			start, end,
		)
	case token.LEFT_OPEN_RANGE_OP:
		start := c.convertToValue(start)
		end := c.convertToValue(end)
		return newGoValueWithDependencies(
			fmt.Sprintf("value.NewLeftOpenRange(%s, %s)", start.value, end.value),
			c.typeOf(node),
			value.FetchGoType("*value.LeftOpenRange"),
			start, end,
		)
	case token.RIGHT_OPEN_RANGE_OP:
		start := c.convertToValue(start)
		end := c.convertToValue(end)
		return newGoValueWithDependencies(
			fmt.Sprintf("value.NewRightOpenRange(%s, %s)", start.value, end.value),
			c.typeOf(node),
			value.FetchGoType("*value.RightOpenRange"),
			start, end,
		)
	default:
		panic(fmt.Sprintf("invalid range operator: %#v", node.Op))
	}
}

func (c *GoCompiler) compileLoopExpressionNode(label string, node *ast.LoopExpressionNode, valueIsIgnored bool) *goValue {
	return c.compileLoop(label, node.ThenBody, c.typeOf(node), valueIsIgnored)
}

func (c *GoCompiler) compileLoop(label string, body []ast.StatementNode, elkType types.Type, valueIsIgnored bool) *goValue {
	var result *goValue
	var tmpVar *goLocal

	if valueIsIgnored {
		result = nilGoValue
	} else {
		goType := c.elkTypeToGoType(elkType, false)
		tmpVar = c.defineTmpGoLocal(goType)
		result = newGoValueWithLocal(tmpVar, elkType)

		if c.checker.IsSubtype(types.Nil{}, elkType) {
			c.emit("%s = value.Nil\n", tmpVar.name)
		}
	}

	c.enterScope(label, loopNativeElkScopeType)
	loopInfo := c.addLoopInfo(label, tmpVar, false)
	c.loopCounter++

	prevBuff := c.switchBuffer(bytes.Buffer{})
	c.emit("for {\n")
	then := c.compileStatements(body)
	if !valueIsIgnored {
		c.emit("%s = %s\n", tmpVar.name, then.fetchValue())
	}

	c.emit("}\n")

	newBuff := c.switchBuffer(prevBuff)
	if loopInfo.labelIsUsed {
		c.emit("%s: ", label)
	}
	c.emitBytes(newBuff.Bytes())

	c.leaveScope()
	c.popLoopInfo()
	return result
}

// Switch to a new buffer and return the previous one
func (c *GoCompiler) switchBuffer(buff bytes.Buffer) bytes.Buffer {
	prevBuff := c.buff
	c.buff = buff
	return prevBuff
}

func (c *GoCompiler) compileWhileExpressionNode(label string, node *ast.WhileExpressionNode, valueIsIgnored bool) *goValue {
	if resolved := resolve(node.Condition, c.checker); resolved.IsNotUndefined() {
		if value.Falsy(resolved) {
			// the loop won't run at all
			// it can be optimised into a simple NIL operation
			return nilGoValue
		}

		// the loop is endless
		return c.compileLoop(label, node.ThenBody, c.typeOf(node), valueIsIgnored)
	}

	var result *goValue
	var tmpVar *goLocal

	if valueIsIgnored {
		result = nilGoValue
	} else {
		elkType := c.typeOf(node)
		goType := c.elkTypeToGoType(elkType, false)
		tmpVar = c.defineTmpGoLocal(goType)
		result = newGoValueWithLocal(tmpVar, elkType)

		if c.checker.IsSubtype(types.Nil{}, elkType) {
			c.emit("%s = value.Nil\n", tmpVar.name)
		}
	}

	c.enterScope(label, loopNativeElkScopeType)
	loopInfo := c.addLoopInfo(label, tmpVar, true)
	c.loopCounter++

	prevBuff := c.switchBuffer(bytes.Buffer{})
	// loop start
	c.emit("for {\n")

	// loop condition eg. `i < 5`
	cond := c.valueToNarrowerType(c.compileExpression(node.Condition, false))

	switch cond.goType.Name {
	case "value.Bool", "bool":
		c.emit("if !(%s) { break }\n", cond.fetchValue())
	default:
		c.emit("if value.Falsy(%s) { break }\n", cond.fetchValue())
	}

	// loop body
	then := c.valueToNarrowerType(c.compileStatements(node.ThenBody))
	if !valueIsIgnored {
		c.emit("%s = %s\n", tmpVar.name, then.fetchValue())
	}

	// after loop
	c.emit("}\n")

	newBuff := c.switchBuffer(prevBuff)
	if loopInfo.labelIsUsed {
		c.emit("%s: ", label)
	}
	c.emitBytes(newBuff.Bytes())

	c.leaveScope()
	c.popLoopInfo()

	return result
}

func (c *GoCompiler) compileMustExpressionNode(node *ast.MustExpressionNode, valueIsIgnored bool) *goValue {
	val := c.compileExpression(node.Value, false)
	defer val.markFree()

	narrowVal := c.valueToNarrowerType(val)

	var result *goValue

	if valueIsIgnored {
		result = val
	} else {
		tmpVar := c.defineTmpGoLocal(narrowVal.goType)
		result = newGoValueWithLocal(tmpVar, narrowVal.elkType)
		c.emit("%s = %s\n", tmpVar.name, narrowVal.fetchValue())
	}

	c.registerErr()
	c.emitSetCallFrameLineNumber(node.Location())
	c.emit("err = value.Must(%s)\n", c.convertToValue(result).fetchValue())
	c.emitErrorPropagation()

	return result
}

func (c *GoCompiler) compileAsExpressionNode(node *ast.AsExpressionNode, valueIsIgnored bool) *goValue {
	val := c.compileExpression(node.Value, false)
	class := c.compileExpression(node.RuntimeType, false)
	defer val.markFree()
	defer class.markFree()

	narrowVal := c.valueToNarrowerType(val)
	narrowClass := c.valueToNarrowerType(class)

	var result *goValue

	if valueIsIgnored {
		result = val
	} else {
		tmpVar := c.defineTmpGoLocal(narrowVal.goType)
		result = newGoValueWithLocal(tmpVar, narrowVal.elkType)
		c.emit("%s = %s\n", tmpVar.name, narrowVal.fetchValue())
	}

	c.registerErr()
	c.emitSetCallFrameLineNumber(node.Location())
	switch narrowClass.goType.Name {
	case "*value.Class", "*value.Mixin":
		c.emit("err = value.As(%s, %s)\n", c.convertToValue(result).fetchValue(), narrowClass.fetchValue())
	default:
		c.emit("err = value.AsUnsafe(%s, %s)\n", c.convertToValue(result).fetchValue(), c.convertToValue(class).fetchValue())
	}
	c.emitErrorPropagation()

	return result
}

func (c *GoCompiler) compileWordArrayTupleLiteralNode(node *ast.WordArrayTupleLiteralNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	c.Errors.AddFailure("invalid word arrayTuple literal", node.Location())
	return nilGoValue
}

func (c *GoCompiler) compileSymbolArrayTupleLiteralNode(node *ast.SymbolArrayTupleLiteralNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	c.Errors.AddFailure("invalid symbol arrayTuple literal", node.Location())
	return nilGoValue
}

func (c *GoCompiler) compileBinArrayTupleLiteralNode(node *ast.BinArrayTupleLiteralNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	c.Errors.AddFailure("invalid binary arrayTuple literal", node.Location())
	return nilGoValue
}

func (c *GoCompiler) compileHexArrayTupleLiteralNode(node *ast.HexArrayTupleLiteralNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	c.Errors.AddFailure("invalid hex arrayTuple literal", node.Location())
	return nilGoValue
}

func (c *GoCompiler) compileWordArrayListLiteralNode(node *ast.WordArrayListLiteralNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	c.Errors.AddFailure("invalid word arrayList literal", node.Location())
	return nilGoValue
}

func (c *GoCompiler) compileSymbolArrayListLiteralNode(node *ast.SymbolArrayListLiteralNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	c.Errors.AddFailure("invalid symbol arrayList literal", node.Location())
	return nilGoValue
}

func (c *GoCompiler) compileBinArrayListLiteralNode(node *ast.BinArrayListLiteralNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	c.Errors.AddFailure("invalid binary arrayList literal", node.Location())
	return nilGoValue
}

func (c *GoCompiler) compileHexArrayListLiteralNode(node *ast.HexArrayListLiteralNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	c.Errors.AddFailure("invalid hex arrayList literal", node.Location())
	return nilGoValue
}

func (c *GoCompiler) compileArrayTupleLiteralNode(node *ast.ArrayTupleLiteralNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	var buff bytes.Buffer

	elementValues := make([]*goValue, 0, len(node.Elements))
	var tmp *goLocal

	typ := c.typeOf(node)
	elementType, _ := c.checker.GetIteratorElementType(typ)
	goElementType := c.elkTypeToGoType(elementType, false)
	var goType *value.GoType
	if goElementType.Name == "value.Value" {
		goType = value.FetchGoType("*value.ArrayTupleOfValue")
	} else {
		goType = value.FetchGenericGoType(
			"*value.NativeArrayTuple",
			[]*value.GoType{
				goElementType,
			},
		)
	}

	finalizeStaticElements := func() {
		if tmp != nil {
			return
		}

		tmp = c.defineTmpGoLocal(goType)
		for _, elementValue := range elementValues {
			if goElementType.Name == "value.Value" {
				buff.WriteString(c.convertToValue(elementValue).value)
			} else {
				buff.WriteString(c.valueToNarrowerType(elementValue).value)
			}
			buff.WriteRune(',')
		}
		buff.WriteString(")")
		c.emit("%s = %s\n", tmp.name, buff.String())
		buff.Reset()
	}

	if goType.Name == "*value.ArrayTupleOfValue" {
		fmt.Fprintf(&buff, "value.NewArrayTupleOfValueWithElementsAndTotalCapacity(%d,", len(node.Elements))
	} else {
		fmt.Fprintf(&buff, "value.NewNativeArrayTupleWithElementsAndTotalCapacity[%s](%d,", goElementType.String(), len(node.Elements))
	}

	for i := 0; i < len(node.Elements); i++ {
		elementNode := node.Elements[i]

		switch elementNode := elementNode.(type) {
		case *ast.KeyValueExpressionNode:
			if tmp != nil {
				key := c.compileExpression(elementNode.Key, false)
				value := c.compileExpression(elementNode.Value, false)

				c.registerErr()
				c.emitSetCallFrameLineNumber(elementNode.Location())
				c.emit("err = %s.AppendAt(%s, %s)\n", tmp.name, c.convertToValue(key).fetchValue(), c.convertToValue(value).fetchValue())
				c.emitErrorPropagation()
				continue
			}
			index, ok := c.parseArrayIndex(elementNode.Key)
			if !ok {
				return errGoValue
			}
			if index == -1 {
				finalizeStaticElements()
				i--
				continue
			}

			value := c.compileExpression(elementNode.Value, false)
			if index >= len(elementValues) {
				newElementsCount := (index + 1) - len(elementValues)
				c.expandValueSlice(&elementValues, newElementsCount)
			}

			elementValues[index] = value
		case *ast.ModifierNode:
			finalizeStaticElements()

			var condType conditionType
			switch elementNode.Modifier.Type {
			case token.IF:
				condType = ifConditionType
			case token.UNLESS:
				condType = unlessConditionType
			default:
				panic(fmt.Sprintf("invalid collection modifier: %#v", elementNode.Modifier))
			}

			c.compileIfWithConditionExpression(
				condType,
				elementNode.Right,
				func() *goValue {
					c.compileArrayAppend(tmp, elementNode.Left)
					return nilGoValue
				},
				nil,
				c.typeOf(elementNode),
				true,
			)
		case *ast.ModifierForInNode:
			finalizeStaticElements()
			// TODO: compile for in
		case *ast.ModifierIfElseNode:
			finalizeStaticElements()

			c.compileIfWithConditionExpression(
				ifConditionType,
				elementNode.Condition,
				func() *goValue {
					c.compileArrayAppend(tmp, elementNode.ThenExpression)
					return nilGoValue
				},
				func() *goValue {
					c.compileArrayAppend(tmp, elementNode.ElseExpression)
					return nilGoValue
				},
				c.typeOf(elementNode),
				true,
			)
		default:
			if tmp != nil {
				c.compileArrayAppend(tmp, elementNode)
			} else {
				element := c.compileExpression(elementNode, false)
				elementValues = append(elementValues, element)
			}
		case *ast.SymbolKeyValueExpressionNode:
			panic(fmt.Sprintf("invalid arraytuple literal element node: %T", elementNode))
		}
	}

	if tmp == nil {
		for _, elementValue := range elementValues {
			if goElementType.Name == "value.Value" {
				buff.WriteString(c.convertToValue(elementValue).value)
			} else {
				buff.WriteString(c.valueToNarrowerType(elementValue).value)
			}
			buff.WriteRune(',')
		}
		buff.WriteString(")")

		return newGoValueWithDependencies(
			buff.String(),
			c.typeOf(node),
			goType,
			elementValues...,
		)
	}

	for _, dependency := range elementValues {
		dependency.markFree()
	}
	return newGoValueWithLocal(
		tmp,
		c.typeOf(node),
	)
}

func (c *GoCompiler) compileArrayAppend(tmp *goLocal, expr ast.ExpressionNode) {
	switch expr := expr.(type) {
	case *ast.KeyValueExpressionNode:
		key := c.compileExpression(expr.Key, false)
		value := c.compileExpression(expr.Value, false)

		c.registerErr()
		c.emitSetCallFrameLineNumber(expr.Location())
		c.emit("err = %s.AppendAt(%s, %s)\n", tmp.name, c.convertToNativeInt(key).fetchValue(), c.convertToValue(value).fetchValue())
		c.emitErrorPropagation()
	default:
		c.compileCollectionAppendExpr(tmp, expr)
	}
}

func (c *GoCompiler) compileCollectionAppendExpr(tmp *goLocal, expr ast.ExpressionNode) {
	c.compileCollectionAppend(tmp, c.compileExpression(expr, false))
}

func (c *GoCompiler) compileCollectionAppend(tmp *goLocal, val *goValue) {
	c.emit("%s.Append(%s)\n", tmp.name, c.convertToValue(val).fetchValue())
}

func (c *GoCompiler) compileHashSetAppendExpr(tmp *goLocal, expr ast.ExpressionNode) {
	c.compileHashSetAppend(tmp, c.compileExpression(expr, false), expr.Location())
}

func (c *GoCompiler) compileHashSetAppend(tmp *goLocal, val *goValue, loc *position.Location) {
	switch tmp.goType.Name {
	case "*vm.HashSetOfValue", "vm.HashSet":
		c.registerErr()
		c.emitSetCallFrameLineNumber(loc)
		c.emit("_, err = %s.AppendVal(thread, %s)\n", tmp.name, c.convertToValue(val).fetchValue())
		c.emitErrorPropagation()
	default:
		c.emit("%s.Append(%s)\n", tmp.name, c.valueToNarrowerType(val).fetchValue())
	}
}

func (c *GoCompiler) compileMapSet(tmp *goLocal, key, val *goValue, loc *position.Location) {
	switch tmp.goType.Name {
	case "*vm.HashMapOfValue", "vm.HashMap", "*vm.HashRecordOfValue", "vm.HashRecord":
		c.registerErr()
		c.emitSetCallFrameLineNumber(loc)
		c.emit("err = %s.SetVal(thread, %s, %s)\n", tmp.name, c.convertToValue(key).fetchValue(), c.convertToValue(val).fetchValue())
		c.emitErrorPropagation()
	case "vm.NativeKeyHashRecord", "*vm.NativeKeyHashMap":
		c.emit("%s.Set(%s, %s)\n", tmp.name, c.valueToNarrowerType(key).fetchValue(), c.convertToValue(val).fetchValue())
	default:
		c.emit("%s.Set(%s, %s)\n", tmp.name, c.valueToNarrowerType(key).fetchValue(), c.valueToNarrowerType(val).fetchValue())
	}
}

func (c *GoCompiler) parseArrayIndex(node ast.ExpressionNode) (int, bool) {
	var index int

	switch n := node.(type) {
	case *ast.IntLiteralNode:
		i, err := value.ParseBigInt(n.Value, 0)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return 0, false
		}
		index = int(i.ToSmallInt())
	case *ast.Int8LiteralNode:
		i, err := value.StrictParseInt(n.Value, 0, 8)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return 0, false
		}
		index = int(i)
	case *ast.Int16LiteralNode:
		i, err := value.StrictParseInt(n.Value, 0, 16)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return 0, false
		}
		index = int(i)
	case *ast.Int32LiteralNode:
		i, err := value.StrictParseInt(n.Value, 0, 32)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return 0, false
		}
		index = int(i)
	case *ast.Int64LiteralNode:
		i, err := value.StrictParseInt(n.Value, 0, 64)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return 0, false
		}
		index = int(i)
	case *ast.UInt8LiteralNode:
		i, err := value.StrictParseUint(n.Value, 0, 8)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return 0, false
		}
		index = int(i)
	case *ast.UInt16LiteralNode:
		i, err := value.StrictParseUint(n.Value, 0, 16)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return 0, false
		}
		index = int(i)
	case *ast.UInt32LiteralNode:
		i, err := value.StrictParseUint(n.Value, 0, 32)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return 0, false
		}
		index = int(i)
	case *ast.UInt64LiteralNode:
		i, err := value.StrictParseUint(n.Value, 0, 64)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return 0, false
		}
		index = int(i)
	case *ast.UIntLiteralNode:
		i, err := value.StrictParseUint(n.Value, 0, value.SmallIntBits)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return 0, false
		}
		index = int(i)
	default:
		return -1, true
	}

	if index < 0 {
		c.Errors.AddFailure(
			fmt.Sprintf("negative array indices are invalid: %d", index),
			node.Location(),
		)
	}

	return index, true
}

func (c *GoCompiler) expandValueSlice(slice *[]*goValue, newElements int) {
	if newElements < 1 {
		return
	}

	newCollection := slices.Grow(*slice, newElements)
	for range newElements {
		newCollection = append(newCollection, nilGoValue)
	}
	*slice = newCollection
}

const invalidCapacityErrMessage = "capacity cannot be specified in collection literals with conditional elements or loops"

func (c *GoCompiler) compileArrayListLiteralNode(node *ast.ArrayListLiteralNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	var buff bytes.Buffer
	var capacity *goValue
	if node.Capacity != nil {
		capacity = c.compileExpression(node.Capacity, false)
	} else {
		capacity = newGoValue(
			"0",
			c.checker.Std(symbol.Int),
			value.FetchGoType("int"),
		)
	}

	elementValues := make([]*goValue, 0, len(node.Elements))
	dependencies := make([]*goValue, 0, len(node.Elements)+1)
	var tmp *goLocal

	typ := c.typeOf(node)
	elementType, _ := c.checker.GetIteratorElementType(typ)
	goElementType := c.elkTypeToGoType(elementType, false)
	var goType *value.GoType
	if goElementType.Name == "value.Value" {
		goType = value.FetchGoType("*value.ArrayListOfValue")
	} else {
		goType = value.FetchGenericGoType(
			"*value.NativeArrayList",
			[]*value.GoType{
				goElementType,
			},
		)
	}

	finalizeStaticElements := func() {
		if tmp != nil {
			return
		}

		tmp = c.defineTmpGoLocal(goType)
		for _, elementValue := range elementValues {
			if goElementType.Name == "value.Value" {
				buff.WriteString(c.convertToValue(elementValue).value)
			} else {
				buff.WriteString(c.valueToNarrowerType(elementValue).value)
			}
			buff.WriteRune(',')
		}
		buff.WriteString(")")
		c.emit("%s = %s\n", tmp.name, buff.String())
		buff.Reset()
		elementValues = nil
	}

	dependencies = append(dependencies, capacity)
	if goType.Name == "*value.ArrayListOfValue" {
		fmt.Fprintf(&buff, "value.NewArrayListOfValueWithElementsAndTotalCapacity(%d + %s,", len(node.Elements), c.convertToNativeInt(capacity).value)
	} else {
		fmt.Fprintf(&buff, "value.NewNativeArrayListWithElementsAndTotalCapacity[%s](%d + %s,", goElementType.String(), len(node.Elements), c.convertToNativeInt(capacity).value)
	}

	for i := 0; i < len(node.Elements); i++ {
		elementNode := node.Elements[i]

		switch elementNode := elementNode.(type) {
		case *ast.KeyValueExpressionNode:
			if tmp != nil {
				key := c.compileExpression(elementNode.Key, false)
				value := c.compileExpression(elementNode.Value, false)

				c.registerErr()
				c.emitSetCallFrameLineNumber(elementNode.Location())
				c.emit("err = %s.AppendAt(%s, %s)\n", tmp.name, c.convertToValue(key).fetchValue(), c.convertToValue(value).fetchValue())
				c.emitErrorPropagation()
				continue
			}
			index, ok := c.parseArrayIndex(elementNode.Key)
			if !ok {
				return errGoValue
			}
			if index == -1 {
				finalizeStaticElements()
				i--
				continue
			}

			value := c.compileExpression(elementNode.Value, false)
			if index >= len(elementValues) {
				newElementsCount := (index + 1) - len(elementValues)
				c.expandValueSlice(&elementValues, newElementsCount)
			}

			elementValues[index] = value
			dependencies = append(dependencies, value)
		case *ast.ModifierNode:
			if node.Capacity != nil {
				c.Errors.AddFailure(
					invalidCapacityErrMessage,
					node.Capacity.Location(),
				)
				return nilGoValue
			}

			finalizeStaticElements()

			var condType conditionType
			switch elementNode.Modifier.Type {
			case token.IF:
				condType = ifConditionType
			case token.UNLESS:
				condType = unlessConditionType
			default:
				panic(fmt.Sprintf("invalid collection modifier: %#v", elementNode.Modifier))
			}

			c.compileIfWithConditionExpression(
				condType,
				elementNode.Right,
				func() *goValue {
					c.compileArrayAppend(tmp, elementNode.Left)
					return nilGoValue
				},
				nil,
				c.typeOf(elementNode),
				true,
			)
		case *ast.ModifierForInNode:
			if node.Capacity != nil {
				c.Errors.AddFailure(
					invalidCapacityErrMessage,
					node.Capacity.Location(),
				)
				return nilGoValue
			}

			// TODO: compile for in
			finalizeStaticElements()
		case *ast.ModifierIfElseNode:
			if node.Capacity != nil {
				c.Errors.AddFailure(
					invalidCapacityErrMessage,
					node.Capacity.Location(),
				)
				return nilGoValue
			}

			finalizeStaticElements()

			c.compileIfWithConditionExpression(
				ifConditionType,
				elementNode.Condition,
				func() *goValue {
					c.compileArrayAppend(tmp, elementNode.ThenExpression)
					return nilGoValue
				},
				func() *goValue {
					c.compileArrayAppend(tmp, elementNode.ElseExpression)
					return nilGoValue
				},
				c.typeOf(elementNode),
				true,
			)
		default:
			if tmp != nil {
				c.compileArrayAppend(tmp, elementNode)
			} else {
				element := c.compileExpression(elementNode, false)
				elementValues = append(elementValues, element)
				dependencies = append(dependencies, dependencies...)
			}
		case *ast.SymbolKeyValueExpressionNode:
			panic(fmt.Sprintf("invalid arraylist literal element node: %T", elementNode))
		}
	}

	if tmp == nil {
		for _, elementValue := range elementValues {
			if goElementType.Name == "value.Value" {
				buff.WriteString(c.convertToValue(elementValue).value)
			} else {
				buff.WriteString(c.valueToNarrowerType(elementValue).value)
			}
			buff.WriteRune(',')
		}
		buff.WriteString(")")

		return newGoValueWithDependencies(
			buff.String(),
			c.typeOf(node),
			goType,
			dependencies...,
		)
	}

	for _, dependency := range dependencies {
		dependency.markFree()
	}
	return newGoValueWithLocal(
		tmp,
		c.typeOf(node),
	)
}

func (c *GoCompiler) compileWordHashSetLiteralNode(node *ast.WordHashSetLiteralNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	c.Errors.AddFailure("invalid word hashset literal", node.Location())
	return nilGoValue
}

func (c *GoCompiler) compileSymbolHashSetLiteralNode(node *ast.SymbolHashSetLiteralNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	c.Errors.AddFailure("invalid symbol hashset literal", node.Location())
	return nilGoValue
}

func (c *GoCompiler) compileBinHashSetLiteralNode(node *ast.BinHashSetLiteralNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	c.Errors.AddFailure("invalid bin hashset literal", node.Location())
	return nilGoValue
}

func (c *GoCompiler) compileHexHashSetLiteralNode(node *ast.HexHashSetLiteralNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	c.Errors.AddFailure("invalid hex hashset literal", node.Location())
	return nilGoValue
}

func (c *GoCompiler) compileHashSetLiteralNode(node *ast.HashSetLiteralNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	var buff bytes.Buffer
	var capacity *goValue
	if node.Capacity != nil {
		capacity = c.compileExpression(node.Capacity, false)
	} else {
		capacity = newGoValue(
			"0",
			c.checker.Std(symbol.Int),
			value.FetchGoType("int"),
		)
	}

	elementValues := make([]*goValue, 0, len(node.Elements))
	dependencies := make([]*goValue, 0, len(node.Elements)+1)
	var tmp *goLocal

	typ := c.typeOf(node)
	elementType, _ := c.checker.GetIteratorElementType(typ)
	goElementType := c.elkTypeToGoType(elementType, false)
	var goType *value.GoType
	if goElementType.Name == "value.Value" {
		goType = value.FetchGoType("*vm.HashSetOfValue")
	} else {
		goType = value.FetchGenericGoType(
			"*vm.NativeHashSet",
			[]*value.GoType{
				goElementType,
			},
		)
	}

	finalizeStaticElements := func() {
		if tmp != nil {
			return
		}

		tmp = c.defineTmpGoLocal(goType)
		for _, elementValue := range elementValues {
			if goType.Name == "*vm.HashSetOfValue" {
				buff.WriteString(c.convertToValue(elementValue).value)
			} else {
				buff.WriteString(c.valueToNarrowerType(elementValue).value)
			}
			buff.WriteRune(',')
			elementValue.markFree()
		}

		buff.WriteString(")")

		for _, dependency := range dependencies {
			dependency.markFree()
		}
		if goType.Name == "*vm.HashSetOfValue" {
			c.registerErr()
			c.emitSetCallFrameLineNumber(node.Location())
			c.emit("%s, err = %s\n", tmp.name, buff.String())
			c.emitErrorPropagation()

		} else {
			c.emit("%s = %s\n", tmp.name, buff.String())
		}
		buff.Reset()
		elementValues = nil
	}

	dependencies = append(dependencies, capacity)
	if goType.Name == "*vm.HashSetOfValue" {
		fmt.Fprintf(&buff, "vm.NewHashSetOfValueWithCapacityAndElements(thread, %d + %s,", len(node.Elements), c.convertToNativeInt(capacity).value)
	} else {
		fmt.Fprintf(&buff, "vm.NewNativeHashSetWithElementsAndTotalCapacity[%s](%d + %s,", goElementType.String(), len(node.Elements), c.convertToNativeInt(capacity).value)
	}

	for i := 0; i < len(node.Elements); i++ {
		elementNode := node.Elements[i]

		switch elementNode := elementNode.(type) {
		case *ast.ModifierNode:
			if node.Capacity != nil {
				c.Errors.AddFailure(
					invalidCapacityErrMessage,
					node.Capacity.Location(),
				)
				return nilGoValue
			}

			finalizeStaticElements()

			var condType conditionType
			switch elementNode.Modifier.Type {
			case token.IF:
				condType = ifConditionType
			case token.UNLESS:
				condType = unlessConditionType
			default:
				panic(fmt.Sprintf("invalid collection modifier: %#v", elementNode.Modifier))
			}

			c.compileIfWithConditionExpression(
				condType,
				elementNode.Right,
				func() *goValue {
					c.compileHashSetAppendExpr(tmp, elementNode.Left)
					return nilGoValue
				},
				nil,
				c.typeOf(elementNode),
				true,
			)
		case *ast.ModifierForInNode:
			if node.Capacity != nil {
				c.Errors.AddFailure(
					invalidCapacityErrMessage,
					node.Capacity.Location(),
				)
				return nilGoValue
			}

			// TODO: compile for in
			finalizeStaticElements()
		case *ast.ModifierIfElseNode:
			if node.Capacity != nil {
				c.Errors.AddFailure(
					invalidCapacityErrMessage,
					node.Capacity.Location(),
				)
				return nilGoValue
			}

			finalizeStaticElements()

			c.compileIfWithConditionExpression(
				ifConditionType,
				elementNode.Condition,
				func() *goValue {
					c.compileHashSetAppendExpr(tmp, elementNode.ThenExpression)
					return nilGoValue
				},
				func() *goValue {
					c.compileHashSetAppendExpr(tmp, elementNode.ElseExpression)
					return nilGoValue
				},
				c.typeOf(elementNode),
				true,
			)
		default:
			if tmp != nil {
				c.compileHashSetAppendExpr(tmp, elementNode)
			} else {
				element := c.compileExpression(elementNode, false)
				elementValues = append(elementValues, element)
				dependencies = append(dependencies, element)
			}
		case *ast.SymbolKeyValueExpressionNode:
			panic(fmt.Sprintf("invalid hashset literal element node: %T", elementNode))
		}
	}

	if tmp == nil {
		if goType.Name != "*vm.HashSetOfValue" {
			for _, elementValue := range elementValues {
				if goElementType.Name == "value.Value" {
					buff.WriteString(c.convertToValue(elementValue).value)
				} else {
					buff.WriteString(c.valueToNarrowerType(elementValue).value)
				}
				buff.WriteRune(',')
			}
			buff.WriteString(")")

			return newGoValueWithDependencies(
				buff.String(),
				c.typeOf(node),
				goType,
				elementValues...,
			)
		}

		finalizeStaticElements()
	} else {
		for _, dependency := range dependencies {
			dependency.markFree()
		}
	}

	return newGoValueWithLocal(
		tmp,
		c.typeOf(node),
	)
}

func (c *GoCompiler) compileHashMapLiteralNode(node *ast.HashMapLiteralNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	var buff bytes.Buffer
	var capacity *goValue
	if node.Capacity != nil {
		capacity = c.compileExpression(node.Capacity, false)
	} else {
		capacity = newGoValue(
			"0",
			c.checker.Std(symbol.Int),
			value.FetchGoType("int"),
		)
	}

	pairValues := make([]goValuePair, 0, len(node.Elements))
	dependencies := make([]*goValue, 0, len(node.Elements)*2+1)
	var tmp *goLocal

	var keyType, valType types.Type
	var goKeyType, goValType *value.GoType

	typ := c.typeOf(node)
	elementType, _ := c.checker.GetIteratorElementType(typ)
	if g, ok := elementType.(*types.Generic); ok {
		if c.checker.IsTheSameNamespace(g.Namespace, c.checker.Std(symbol.Pair).(*types.Class)) {
			keyType = types.Normalise(g.Get(0).Type)
			goKeyType = c.elkTypeToGoKeyType(keyType)

			valType = types.Normalise(g.Get(1).Type)
			goValType = c.elkTypeToGoType(valType, false)
		}
	}
	var goType *value.GoType
	if goKeyType.Name == "value.Value" {
		goType = value.FetchGoType("*vm.HashMapOfValue")
	} else if goValType.Name == "value.Value" {
		goType = value.FetchGenericGoType(
			"*vm.NativeKeyHashMap",
			[]*value.GoType{
				goKeyType,
			},
		)
	} else {
		goType = value.FetchGenericGoType(
			"*vm.NativeHashMap",
			[]*value.GoType{
				goKeyType,
				goValType,
			},
		)
	}

	finalizeStaticElements := func() {
		if tmp != nil {
			return
		}

		tmp = c.defineTmpGoLocal(goType)
		for _, pairValue := range pairValues {
			switch goType.Name {
			case "*vm.HashMapOfValue":
				fmt.Fprintf(
					&buff,
					"value.MakePairOfValue(%s, %s)",
					c.convertToValue(pairValue.key).value,
					c.convertToValue(pairValue.value).value,
				)
			case "*vm.NativeKeyHashMap":
				fmt.Fprintf(
					&buff,
					"value.MakeNativePair(%s, %s)",
					c.valueToNarrowerType(pairValue.key).value,
					c.convertToValue(pairValue.value).value,
				)
			case "*vm.NativeHashMap":
				fmt.Fprintf(
					&buff,
					"value.MakeNativePair(%s, %s)",
					c.valueToNarrowerType(pairValue.key).value,
					c.valueToNarrowerType(pairValue.value).value,
				)
			default:
				panic(fmt.Sprintf("invalid hash map go type: %s", goType.String()))
			}

			buff.WriteRune(',')
		}

		for _, dependency := range dependencies {
			dependency.markFree()
		}
		buff.WriteString(")")
		if goType.Name == "*vm.HashMapOfValue" {
			c.registerErr()
			c.emitSetCallFrameLineNumber(node.Location())
			c.emit("%s, err = %s\n", tmp.name, buff.String())
			c.emitErrorPropagation()
		} else {
			c.emit("%s = %s\n", tmp.name, buff.String())
		}
		buff.Reset()
		pairValues = nil
	}

	dependencies = append(dependencies, capacity)
	switch goType.Name {
	case "*vm.HashMapOfValue":
		fmt.Fprintf(
			&buff,
			"vm.NewHashMapOfValueWithCapacityAndElements(thread, %d + %s,",
			len(node.Elements),
			c.convertToNativeInt(capacity).value,
		)
	case "*vm.NativeKeyHashMap":
		fmt.Fprintf(
			&buff,
			"vm.NewNativeKeyHashMapWithElementsAndTotalCapacity[%s](%d + %s,",
			goKeyType.String(),
			len(node.Elements),
			c.convertToNativeInt(capacity).value,
		)
	case "*vm.NativeHashMap":
		fmt.Fprintf(
			&buff,
			"vm.NewNativeHashMapWithElementsAndTotalCapacity[%s, %s](%d + %s,",
			goKeyType.String(),
			goValType.String(),
			len(node.Elements),
			c.convertToNativeInt(capacity).value,
		)
	default:
		panic(fmt.Sprintf("invalid hash map go type: %s", goType.String()))
	}

	for i := 0; i < len(node.Elements); i++ {
		elementNode := node.Elements[i]

		switch elementNode := elementNode.(type) {
		case *ast.ModifierNode:
			if node.Capacity != nil {
				c.Errors.AddFailure(
					invalidCapacityErrMessage,
					node.Capacity.Location(),
				)
				return nilGoValue
			}

			finalizeStaticElements()

			var condType conditionType
			switch elementNode.Modifier.Type {
			case token.IF:
				condType = ifConditionType
			case token.UNLESS:
				condType = unlessConditionType
			default:
				panic(fmt.Sprintf("invalid collection modifier: %#v", elementNode.Modifier))
			}

			c.compileIfWithConditionExpression(
				condType,
				elementNode.Right,
				func() *goValue {
					switch then := elementNode.Left.(type) {
					case *ast.KeyValueExpressionNode:
						key := c.compileExpression(then.Key, false)
						val := c.compileExpression(then.Value, false)
						c.compileMapSet(tmp, key, val, then.Location())
					case *ast.SymbolKeyValueExpressionNode:
						key := c.valueToGoSource(value.ToSymbol(identifierToName(then.Key)).ToValue(), c.checker.Std(symbol.Symbol), false)
						val := c.compileExpression(then.Value, false)
						c.compileMapSet(tmp, key, val, then.Location())
					default:
						panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
					}

					return nilGoValue
				},
				nil,
				c.typeOf(elementNode),
				true,
			)
		case *ast.ModifierForInNode:
			if node.Capacity != nil {
				c.Errors.AddFailure(
					invalidCapacityErrMessage,
					node.Capacity.Location(),
				)
				return nilGoValue
			}

			// TODO: compile for in
			finalizeStaticElements()
		case *ast.ModifierIfElseNode:
			if node.Capacity != nil {
				c.Errors.AddFailure(
					invalidCapacityErrMessage,
					node.Capacity.Location(),
				)
				return nilGoValue
			}

			finalizeStaticElements()

			c.compileIfWithConditionExpression(
				ifConditionType,
				elementNode.Condition,
				func() *goValue {
					switch then := elementNode.ThenExpression.(type) {
					case *ast.KeyValueExpressionNode:
						key := c.compileExpression(then.Key, false)
						val := c.compileExpression(then.Value, false)
						c.compileMapSet(tmp, key, val, then.Location())
					case *ast.SymbolKeyValueExpressionNode:
						key := c.valueToGoSource(value.ToSymbol(identifierToName(then.Key)).ToValue(), c.checker.Std(symbol.Symbol), false)
						val := c.compileExpression(then.Value, false)
						c.compileMapSet(tmp, key, val, then.Location())
					default:
						panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
					}

					return nilGoValue
				},
				func() *goValue {
					switch then := elementNode.ElseExpression.(type) {
					case *ast.KeyValueExpressionNode:
						key := c.compileExpression(then.Key, false)
						val := c.compileExpression(then.Value, false)
						c.compileMapSet(tmp, key, val, then.Location())
					case *ast.SymbolKeyValueExpressionNode:
						key := c.valueToGoSource(value.ToSymbol(identifierToName(then.Key)).ToValue(), c.checker.Std(symbol.Symbol), false)
						val := c.compileExpression(then.Value, false)
						c.compileMapSet(tmp, key, val, then.Location())
					default:
						panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
					}

					return nilGoValue
				},
				c.typeOf(elementNode),
				true,
			)
		case *ast.PublicIdentifierNode:
			key := c.valueToGoSource(value.ToSymbol(elementNode.Value).ToValue(), c.checker.Std(symbol.Symbol), false)
			val := c.compileLocalVariableAccess(elementNode.Value, c.typeOf(elementNode))
			if tmp != nil {
				c.compileMapSet(tmp, key, val, elementNode.Location())
			} else {
				pairValues = append(pairValues, goValuePair{key: key, value: val})
				dependencies = append(dependencies, key, val)
			}
		case *ast.PrivateIdentifierNode:
			key := c.valueToGoSource(value.ToSymbol(elementNode.Value).ToValue(), c.checker.Std(symbol.Symbol), false)
			val := c.compileLocalVariableAccess(elementNode.Value, c.typeOf(elementNode))
			if tmp != nil {
				c.compileMapSet(tmp, key, val, elementNode.Location())
			} else {
				pairValues = append(pairValues, goValuePair{key: key, value: val})
				dependencies = append(dependencies, key, val)
			}
		case *ast.KeyValueExpressionNode:
			key := c.compileExpression(elementNode.Key, false)
			val := c.compileExpression(elementNode.Value, false)
			if tmp != nil {
				c.compileMapSet(tmp, key, val, elementNode.Location())
			} else {
				pairValues = append(pairValues, goValuePair{key: key, value: val})
				dependencies = append(dependencies, key, val)
			}
		case *ast.SymbolKeyValueExpressionNode:
			key := c.valueToGoSource(value.ToSymbol(identifierToName(elementNode.Key)).ToValue(), c.checker.Std(symbol.Symbol), false)
			val := c.compileExpression(elementNode.Value, false)
			if tmp != nil {
				c.compileMapSet(tmp, key, val, elementNode.Location())
			} else {
				pairValues = append(pairValues, goValuePair{key: key, value: val})
				dependencies = append(dependencies, key, val)
			}
		default:
			panic(fmt.Sprintf("invalid element in hashmap literal: %#v", elementNode))
		}
	}

	if tmp == nil {
		switch goType.Name {
		case "*vm.NativeKeyHashMap":
			for _, pairValue := range pairValues {
				fmt.Fprintf(
					&buff,
					"value.MakeNativePair(%s, %s),",
					c.valueToNarrowerType(pairValue.key).value,
					c.convertToValue(pairValue.value).value,
				)
			}
			buff.WriteString(")")

			return newGoValueWithDependencies(
				buff.String(),
				c.typeOf(node),
				goType,
				dependencies...,
			)
		case "*vm.NativeHashMap":
			for _, pairValue := range pairValues {
				fmt.Fprintf(
					&buff,
					"value.MakeNativePair(%s, %s),",
					c.valueToNarrowerType(pairValue.key).value,
					c.valueToNarrowerType(pairValue.value).value,
				)
			}
			buff.WriteString(")")

			return newGoValueWithDependencies(
				buff.String(),
				c.typeOf(node),
				goType,
				dependencies...,
			)
		default:
			finalizeStaticElements()
		}
	} else {
		for _, dependency := range dependencies {
			dependency.markFree()
		}
	}

	return newGoValueWithLocal(
		tmp,
		c.typeOf(node),
	)
}

func (c *GoCompiler) compileHashRecordLiteralNode(node *ast.HashRecordLiteralNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	var buff bytes.Buffer

	pairValues := make([]goValuePair, 0, len(node.Elements))
	dependencies := make([]*goValue, 0, len(node.Elements)*2)
	var tmp *goLocal

	var keyType, valType types.Type
	var goKeyType, goValType *value.GoType

	typ := c.typeOf(node)
	elementType, _ := c.checker.GetIteratorElementType(typ)
	if g, ok := elementType.(*types.Generic); ok {
		if c.checker.IsTheSameNamespace(g.Namespace, c.checker.Std(symbol.Pair).(*types.Class)) {
			keyType = types.Normalise(g.Get(0).Type)
			goKeyType = c.elkTypeToGoKeyType(keyType)

			valType = types.Normalise(g.Get(1).Type)
			goValType = c.elkTypeToGoType(valType, false)
		}
	}
	var goType *value.GoType
	if goKeyType.Name == "value.Value" {
		goType = value.FetchGoType("*vm.HashRecordOfValue")
	} else if goValType.Name == "value.Value" {
		goType = value.FetchGenericGoType(
			"vm.NativeKeyHashRecord",
			[]*value.GoType{
				goKeyType,
			},
		)
	} else {
		goType = value.FetchGenericGoType(
			"vm.NativeHashRecord",
			[]*value.GoType{
				goKeyType,
				goValType,
			},
		)
	}

	finalizeStaticElements := func() {
		if tmp != nil {
			return
		}

		tmp = c.defineTmpGoLocal(goType)
		switch goType.Name {
		case "*vm.HashRecordOfValue":
			for _, pairValue := range pairValues {
				fmt.Fprintf(
					&buff,
					"value.MakePairOfValue(%s, %s),",
					c.convertToValue(pairValue.key).value,
					c.convertToValue(pairValue.value).value,
				)
			}
			buff.WriteString(")")
		case "vm.NativeKeyHashRecord":
			for _, pairValue := range pairValues {
				fmt.Fprintf(
					&buff,
					"%s: %s,",
					c.valueToNarrowerType(pairValue.key).value,
					c.convertToValue(pairValue.value).value,
				)
			}
			buff.WriteString("})")
		case "vm.NativeHashRecord":
			for _, pairValue := range pairValues {
				fmt.Fprintf(
					&buff, "%s: %s,",
					c.valueToNarrowerType(pairValue.key).value,
					c.valueToNarrowerType(pairValue.value).value,
				)
			}
			buff.WriteString("})")
		default:
			panic(fmt.Sprintf("invalid hash record go type: %s", goType.String()))
		}

		for _, dependency := range dependencies {
			dependency.markFree()
		}
		if goType.Name == "*vm.HashRecordOfValue" {
			c.registerErr()
			c.emitSetCallFrameLineNumber(node.Location())
			c.emit("%s, err = %s\n", tmp.name, buff.String())
			c.emitErrorPropagation()
		} else {
			c.emit("%s = %s\n", tmp.name, buff.String())
		}
		buff.Reset()
		pairValues = nil
	}

	switch goType.Name {
	case "*vm.HashRecordOfValue":
		buff.WriteString("vm.NewHashRecordOfValueWithElements(thread,")
	case "vm.NativeKeyHashRecord":
		fmt.Fprintf(&buff, "vm.MakeNativeKeyHashRecordFromMap(map[%s]value.Value{", goKeyType.String())
	case "vm.NativeHashRecord":
		fmt.Fprintf(&buff, "vm.MakeNativeHashRecordFromMap(map[%s]%s{", goKeyType.String(), goValType.String())
	default:
		panic(fmt.Sprintf("invalid hash record go type: %s", goType.String()))
	}

	for i := 0; i < len(node.Elements); i++ {
		elementNode := node.Elements[i]

		switch elementNode := elementNode.(type) {
		case *ast.ModifierNode:
			finalizeStaticElements()

			var condType conditionType
			switch elementNode.Modifier.Type {
			case token.IF:
				condType = ifConditionType
			case token.UNLESS:
				condType = unlessConditionType
			default:
				panic(fmt.Sprintf("invalid collection modifier: %#v", elementNode.Modifier))
			}

			c.compileIfWithConditionExpression(
				condType,
				elementNode.Right,
				func() *goValue {
					switch then := elementNode.Left.(type) {
					case *ast.KeyValueExpressionNode:
						key := c.compileExpression(then.Key, false)
						val := c.compileExpression(then.Value, false)
						c.compileMapSet(tmp, key, val, then.Location())
					case *ast.SymbolKeyValueExpressionNode:
						key := c.valueToGoSource(value.ToSymbol(identifierToName(then.Key)).ToValue(), c.checker.Std(symbol.Symbol), false)
						val := c.compileExpression(then.Value, false)
						c.compileMapSet(tmp, key, val, then.Location())
					default:
						panic(fmt.Sprintf("invalid hash record element: %#v", elementNode))
					}

					return nilGoValue
				},
				nil,
				c.typeOf(elementNode),
				true,
			)
		case *ast.ModifierForInNode:
			// TODO: compile for in
			finalizeStaticElements()
		case *ast.ModifierIfElseNode:
			finalizeStaticElements()

			c.compileIfWithConditionExpression(
				ifConditionType,
				elementNode.Condition,
				func() *goValue {
					switch then := elementNode.ThenExpression.(type) {
					case *ast.KeyValueExpressionNode:
						key := c.compileExpression(then.Key, false)
						val := c.compileExpression(then.Value, false)
						c.compileMapSet(tmp, key, val, then.Location())
					case *ast.SymbolKeyValueExpressionNode:
						key := c.valueToGoSource(value.ToSymbol(identifierToName(then.Key)).ToValue(), c.checker.Std(symbol.Symbol), false)
						val := c.compileExpression(then.Value, false)
						c.compileMapSet(tmp, key, val, then.Location())
					default:
						panic(fmt.Sprintf("invalid hash record element: %#v", elementNode))
					}

					return nilGoValue
				},
				func() *goValue {
					switch then := elementNode.ElseExpression.(type) {
					case *ast.KeyValueExpressionNode:
						key := c.compileExpression(then.Key, false)
						val := c.compileExpression(then.Value, false)
						c.compileMapSet(tmp, key, val, then.Location())
					case *ast.SymbolKeyValueExpressionNode:
						key := c.valueToGoSource(value.ToSymbol(identifierToName(then.Key)).ToValue(), c.checker.Std(symbol.Symbol), false)
						val := c.compileExpression(then.Value, false)
						c.compileMapSet(tmp, key, val, then.Location())
					default:
						panic(fmt.Sprintf("invalid hash record element: %#v", elementNode))
					}

					return nilGoValue
				},
				c.typeOf(elementNode),
				true,
			)
		case *ast.PublicIdentifierNode:
			key := c.valueToGoSource(value.ToSymbol(elementNode.Value).ToValue(), c.checker.Std(symbol.Symbol), false)
			val := c.compileLocalVariableAccess(elementNode.Value, c.typeOf(elementNode))
			if tmp != nil {
				c.compileMapSet(tmp, key, val, elementNode.Location())
			} else {
				pairValues = append(pairValues, goValuePair{key: key, value: val})
				dependencies = append(dependencies, key, val)
			}
		case *ast.PrivateIdentifierNode:
			key := c.valueToGoSource(value.ToSymbol(elementNode.Value).ToValue(), c.checker.Std(symbol.Symbol), false)
			val := c.compileLocalVariableAccess(elementNode.Value, c.typeOf(elementNode))
			if tmp != nil {
				c.compileMapSet(tmp, key, val, elementNode.Location())
			} else {
				pairValues = append(pairValues, goValuePair{key: key, value: val})
				dependencies = append(dependencies, key, val)
			}
		case *ast.KeyValueExpressionNode:
			key := c.compileExpression(elementNode.Key, false)
			val := c.compileExpression(elementNode.Value, false)
			if tmp != nil {
				c.compileMapSet(tmp, key, val, elementNode.Location())
			} else {
				pairValues = append(pairValues, goValuePair{key: key, value: val})
				dependencies = append(dependencies, key, val)
			}
		case *ast.SymbolKeyValueExpressionNode:
			key := c.valueToGoSource(value.ToSymbol(identifierToName(elementNode.Key)).ToValue(), c.checker.Std(symbol.Symbol), false)
			val := c.compileExpression(elementNode.Value, false)
			if tmp != nil {
				c.compileMapSet(tmp, key, val, elementNode.Location())
			} else {
				pairValues = append(pairValues, goValuePair{key: key, value: val})
				dependencies = append(dependencies, key, val)
			}
		default:
			panic(fmt.Sprintf("invalid element in hashrecord literal: %#v", elementNode))
		}
	}

	if tmp == nil {
		switch goType.Name {
		case "vm.NativeKeyHashRecord":
			for _, pairValue := range pairValues {
				fmt.Fprintf(
					&buff,
					"%s: %s,",
					c.valueToNarrowerType(pairValue.key).value,
					c.convertToValue(pairValue.value).value,
				)
			}
			buff.WriteString("})")

			return newGoValueWithDependencies(
				buff.String(),
				c.typeOf(node),
				goType,
				dependencies...,
			)
		case "vm.NativeHashRecord":
			for _, pairValue := range pairValues {
				fmt.Fprintf(
					&buff,
					"%s: %s,",
					c.valueToNarrowerType(pairValue.key).value,
					c.valueToNarrowerType(pairValue.value).value,
				)
			}
			buff.WriteString("})")

			return newGoValueWithDependencies(
				buff.String(),
				c.typeOf(node),
				goType,
				dependencies...,
			)
		default:
			finalizeStaticElements()
		}
	} else {
		for _, dependency := range dependencies {
			dependency.markFree()
		}
	}

	return newGoValueWithLocal(
		tmp,
		c.typeOf(node),
	)
}

func (c *GoCompiler) compileStringInterpolationNode(node *ast.StringInterpolationNode) *goValue {
	return c.compileInterpolationNode(node.Expression, node.Location())
}

func (c *GoCompiler) compileRegexInterpolationNode(node *ast.RegexInterpolationNode) *goValue {
	return c.compileInterpolationNode(node.Expression, node.Location())
}

func (c *GoCompiler) compileInterpolationNode(expr ast.ExpressionNode, loc *position.Location) *goValue {
	exprVal := c.compileExpression(expr, false)
	exprVal = c.valueToNarrowerType(exprVal)

	typ := exprVal.goType
	switch typ.Name {
	case "value.String":
		return exprVal
	case "value.Char", "value.Float64", "value.Float32",
		"value.Float", "value.SmallInt",
		"value.Int64", "value.Int32", "value.Int16", "value.Int8",
		"value.UInt64", "value.UInt32", "value.UInt16", "value.UInt8",
		"value.Symbol", "*value.BigInt", "*value.Regex":
		return exprVal.newGoValue(
			fmt.Sprintf("(%s).ToString()", exprVal.value),
			c.checker.Std(symbol.String),
			value.FetchGoType("value.String"),
		)
	}

	return c.compileMethodCallWithLiteralArgValuesAndName(
		exprVal.elkType,
		c.checker.StdString(),
		"symbol.L_to_string",
		"to_string",
		[]*goValue{exprVal},
		loc,
		false,
	)
}

func (c *GoCompiler) compileStringInspectInterpolationNode(node *ast.StringInspectInterpolationNode) *goValue {
	expr := c.compileExpression(node.Expression, false)
	expr = c.valueToNarrowerType(expr)

	typ := expr.goType
	switch typ.Name {
	case "value.String", "value.Char", "value.Float64", "value.Float32",
		"value.Float", "value.SmallInt",
		"value.Int64", "value.Int32", "value.Int16", "value.Int8",
		"value.UInt64", "value.UInt32", "value.UInt16", "value.UInt8",
		"value.Symbol", "*value.BigInt", "*value.Regex":
		return expr.newGoValue(
			fmt.Sprintf("value.String((%s).Inspect())", expr.value),
			c.checker.Std(symbol.String),
			value.FetchGoType("value.String"),
		)
	}

	return c.compileMethodCallWithLiteralArgValuesAndName(
		expr.elkType,
		c.checker.StdString(),
		"symbol.L_inspect",
		"inspect",
		[]*goValue{expr},
		node.Location(),
		false,
	)
}

func (c *GoCompiler) compileStringLiteralContentNode(node ast.StringLiteralContentNode) *goValue {
	switch node := node.(type) {
	case *ast.StringLiteralContentSectionNode:
		return newGoValue(
			fmt.Sprintf("value.String(%q)", node.Value),
			c.checker.StdString(),
			value.FetchGoType("value.String"),
		)
	case *ast.StringInspectInterpolationNode:
		return c.compileStringInspectInterpolationNode(node)
	case *ast.StringInterpolationNode:
		return c.compileStringInterpolationNode(node)
	default:
		panic(fmt.Sprintf("invalid string content node: %T", node))
	}
}

func (c *GoCompiler) compileRegexLiteralContentNode(node ast.RegexLiteralContentNode) *goValue {
	switch node := node.(type) {
	case *ast.RegexLiteralContentSectionNode:
		return newGoValue(
			fmt.Sprintf("value.String(%q)", node.Value),
			c.checker.StdString(),
			value.FetchGoType("value.String"),
		)
	case *ast.RegexInterpolationNode:
		return c.compileRegexInterpolationNode(node)
	default:
		panic(fmt.Sprintf("invalid regex content node: %T", node))
	}
}

func (c *GoCompiler) compileUninterpolatedRegexLiteralNode(node *ast.UninterpolatedRegexLiteralNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	re, err := value.CompileRegex(node.Content, node.Flags)
	if mergeRegexDiagnostics(c.Errors, err, node.Location()) {
		return nilGoValue
	}

	return c.valueToGoSource(re.ToValue(), c.typeOf(node), false)
}

func (c *GoCompiler) compileInterpolatedRegexLiteralNode(node *ast.InterpolatedRegexLiteralNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	elkType := c.typeOf(node)
	goType := c.elkTypeToGoType(elkType, false)
	resultVar := c.defineTmpGoLocal(goType)
	c.registerErr()

	var buff bytes.Buffer
	fmt.Fprintf(&buff, "%s, err = value.CompileRegexVal(", resultVar.name)

	for i, element := range node.Content {
		if i != 0 {
			buff.WriteString(" + ")
		}
		section := c.compileRegexLiteralContentNode(element)
		buff.WriteString(section.fetchValue())
	}

	buff.WriteByte(',')
	c.compileRegexFlags(&buff, node.Flags)

	buff.WriteString(")\n")
	c.emitBytes(buff.Bytes())

	return newGoValueWithLocal(
		resultVar,
		elkType,
	)
}

func (c *GoCompiler) compileInterpolatedStringLiteralNode(node *ast.InterpolatedStringLiteralNode) *goValue {
	var buff strings.Builder
	var dependencies []*goValue

	for i, element := range node.Content {
		if i != 0 {
			buff.WriteString(" + ")
		}
		section := c.compileStringLiteralContentNode(element)
		dependencies = append(dependencies, section)
		buff.WriteString(section.value)
	}

	return newGoValueWithDependencies(
		buff.String(),
		c.checker.StdString(),
		value.FetchGoType("value.String"),
		dependencies...,
	)
}

func (c *GoCompiler) compileInterpolatedSymbolLiteralNode(node *ast.InterpolatedSymbolLiteralNode) *goValue {
	result := c.compileInterpolatedStringLiteralNode(node.Content)
	return result.newGoValue(
		fmt.Sprintf("(%s).ToSymbol()", result.value),
		c.typeOf(node),
		value.FetchGoType("value.Symbol"),
	)
}

func (c *GoCompiler) compileDoubleQuotedStringLiteralNode(node *ast.DoubleQuotedStringLiteralNode) *goValue {
	return newGoValue(
		fmt.Sprintf("value.String(%q)", node.Value),
		c.typeOf(node),
		value.FetchGoType("value.String"),
	)
}

func (c *GoCompiler) compileRawStringLiteralNode(node *ast.RawStringLiteralNode) *goValue {
	return newGoValue(
		fmt.Sprintf("value.String(%q)", node.Value),
		c.typeOf(node),
		value.FetchGoType("value.String"),
	)
}

func (c *GoCompiler) compileIntLiteralNode(node *ast.IntLiteralNode) *goValue {
	i, err := value.ParseBigInt(node.Value, 0)
	if !err.IsUndefined() {
		c.Errors.AddFailure(err.Error(), node.Location())
		return errGoValue
	}
	if i.IsSmallInt() {
		return newGoValue(
			fmt.Sprintf("value.SmallInt(%d)", i.ToSmallInt()),
			c.typeOf(node),
			value.FetchGoType("value.SmallInt"),
		)
	}
	bigIntVar := c.emitBigInt(node.Value)
	return newGoValue(
		bigIntVar,
		c.typeOf(node),
		value.FetchGoType("*value.BigInt"),
	)
}

func (c *GoCompiler) compileInt8LiteralNode(node *ast.Int8LiteralNode) *goValue {
	i, err := value.StrictParseInt(node.Value, 0, 8)
	if !err.IsUndefined() {
		c.Errors.AddFailure(err.Error(), node.Location())
		return errGoValue
	}
	return newGoValue(
		fmt.Sprintf("value.Int8(%d)", i),
		c.typeOf(node),
		value.FetchGoType("value.Int8"),
	)
}

func (c *GoCompiler) compileInt16LiteralNode(node *ast.Int16LiteralNode) *goValue {
	i, err := value.StrictParseInt(node.Value, 0, 16)
	if !err.IsUndefined() {
		c.Errors.AddFailure(err.Error(), node.Location())
		return errGoValue
	}
	return newGoValue(
		fmt.Sprintf("value.Int16(%d)", i),
		c.typeOf(node),
		value.FetchGoType("value.Int16"),
	)
}

func (c *GoCompiler) compileInt32LiteralNode(node *ast.Int32LiteralNode) *goValue {
	i, err := value.StrictParseInt(node.Value, 0, 32)
	if !err.IsUndefined() {
		c.Errors.AddFailure(err.Error(), node.Location())
		return errGoValue
	}
	return newGoValue(
		fmt.Sprintf("value.Int32(%d)", i),
		c.typeOf(node),
		value.FetchGoType("value.Int32"),
	)
}

func (c *GoCompiler) compileInt64LiteralNode(node *ast.Int64LiteralNode) *goValue {
	i, err := value.StrictParseInt(node.Value, 0, 64)
	if !err.IsUndefined() {
		c.Errors.AddFailure(err.Error(), node.Location())
		return errGoValue
	}
	return newGoValue(
		fmt.Sprintf("value.Int64(%d)", i),
		c.typeOf(node),
		value.FetchGoType("value.Int64"),
	)
}

func (c *GoCompiler) compileUInt8LiteralNode(node *ast.UInt8LiteralNode) *goValue {
	i, err := value.StrictParseUint(node.Value, 0, 8)
	if !err.IsUndefined() {
		c.Errors.AddFailure(err.Error(), node.Location())
		return errGoValue
	}
	return newGoValue(
		fmt.Sprintf("value.UInt8(%d)", i),
		c.typeOf(node),
		value.FetchGoType("value.UInt8"),
	)
}

func (c *GoCompiler) compileUInt16LiteralNode(node *ast.UInt16LiteralNode) *goValue {
	i, err := value.StrictParseUint(node.Value, 0, 16)
	if !err.IsUndefined() {
		c.Errors.AddFailure(err.Error(), node.Location())
		return errGoValue
	}
	return newGoValue(
		fmt.Sprintf("value.UInt16(%d)", i),
		c.typeOf(node),
		value.FetchGoType("value.UInt16"),
	)
}

func (c *GoCompiler) compileUInt32LiteralNode(node *ast.UInt32LiteralNode) *goValue {
	i, err := value.StrictParseUint(node.Value, 0, 32)
	if !err.IsUndefined() {
		c.Errors.AddFailure(err.Error(), node.Location())
		return errGoValue
	}
	return newGoValue(
		fmt.Sprintf("value.UInt32(%d)", i),
		c.typeOf(node),
		value.FetchGoType("value.UInt32"),
	)
}

func (c *GoCompiler) compileUInt64LiteralNode(node *ast.UInt64LiteralNode) *goValue {
	i, err := value.StrictParseUint(node.Value, 0, 64)
	if !err.IsUndefined() {
		c.Errors.AddFailure(err.Error(), node.Location())
		return errGoValue
	}
	return newGoValue(
		fmt.Sprintf("value.UInt64(%d)", i),
		c.typeOf(node),
		value.FetchGoType("value.UInt64"),
	)
}

func (c *GoCompiler) compileUIntLiteralNode(node *ast.UIntLiteralNode) *goValue {
	i, err := value.StrictParseUint(node.Value, 0, value.SmallIntBits)
	if !err.IsUndefined() {
		c.Errors.AddFailure(err.Error(), node.Location())
		return errGoValue
	}
	return newGoValue(
		fmt.Sprintf("value.UInt(%d)", i),
		c.typeOf(node),
		value.FetchGoType("value.UInt"),
	)
}

func (c *GoCompiler) compileBigFloatLiteralNode(node *ast.BigFloatLiteralNode) *goValue {
	v := c.emitBigFloat(node.Value)
	return newGoValue(
		v,
		c.typeOf(node),
		value.FetchGoType("*value.BigFloat"),
	)
}

func (c *GoCompiler) compileFloatLiteralNode(node *ast.FloatLiteralNode) *goValue {
	f, err := strconv.ParseFloat(node.Value, value.SmallIntBits)
	if err != nil {
		c.Errors.AddFailure(err.Error(), node.Location())
		return errGoValue
	}
	return newGoValue(
		fmt.Sprintf("value.Float(%g)", f),
		c.typeOf(node),
		value.FetchGoType("value.Float"),
	)
}

func (c *GoCompiler) compileFloat64LiteralNode(node *ast.Float64LiteralNode) *goValue {
	f, err := strconv.ParseFloat(node.Value, 64)
	if err != nil {
		c.Errors.AddFailure(err.Error(), node.Location())
		return errGoValue
	}
	return newGoValue(
		fmt.Sprintf("value.Float64(%g)", f),
		c.typeOf(node),
		value.FetchGoType("value.Float64"),
	)
}

func (c *GoCompiler) compileFloat32LiteralNode(node *ast.Float32LiteralNode) *goValue {
	f, err := strconv.ParseFloat(node.Value, 32)
	if err != nil {
		c.Errors.AddFailure(err.Error(), node.Location())
		return errGoValue
	}
	return newGoValue(
		fmt.Sprintf("value.Float32(%g)", f),
		c.typeOf(node),
		value.FetchGoType("value.Float32"),
	)
}

func (c *GoCompiler) compileVariableDeclarationNode(node *ast.VariableDeclarationNode) *goValue {
	initialised := node.Initialiser != nil

	var elkType types.Type
	if node.TypeNode != nil {
		elkType = c.typeOf(node.TypeNode)
	} else {
		elkType = c.typeOf(node)
	}

	goType := c.elkTypeToGoType(elkType, false)
	local := c.defineLocal(
		identifierToName(node.Name),
		elkType,
		goType,
		node.Location(),
	)
	if local == nil {
		return errGoValue
	}

	if initialised {
		init := c.compileExpression(node.Initialiser, false)
		return c.emitSetLocal(local.name, c.valueToNarrowerType(init))
	}

	return nilGoValue
}

func (c *GoCompiler) compileValueDeclarationNode(node *ast.ValueDeclarationNode) *goValue {
	initialised := node.Initialiser != nil

	var elkType types.Type
	if node.TypeNode != nil {
		elkType = c.typeOf(node.TypeNode)
	} else {
		elkType = c.typeOf(node)
	}

	if initialised {
		init := c.valueToNarrowerType(c.compileExpression(node.Initialiser, false))
		local := c.defineLocal(
			identifierToName(node.Name),
			elkType,
			init.goType,
			node.Location(),
		)
		if local == nil {
			return errGoValue
		}

		return c.emitSetLocal(local.name, init)
	}

	local := c.defineLocal(
		identifierToName(node.Name),
		elkType,
		c.elkTypeToGoType(elkType, false),
		node.Location(),
	)
	if local == nil {
		return errGoValue
	}

	return nilGoValue
}

func (c *GoCompiler) compileReturnExpressionNode(node *ast.ReturnExpressionNode) *goValue {
	var val string

	if node.Value != nil {
		expr := c.compileExpression(node.Value, false)
		if c.method == nil {
			val = c.convertToValue(expr).fetchValue()
		} else {
			goReturnType := c.elkTypeToGoType(c.method.ReturnType, false)
			if goReturnType.Name == "value.Value" {
				val = c.convertToValue(expr).fetchValue()
			} else {
				val = c.valueToNarrowerType(expr).fetchValue()
			}
		}
	} else {
		val = "value.Nil"
	}

	c.emitReturn(val)
	return newGoValue(
		"value.Nil",
		types.Never{},
		goValueType,
	)
}

func (c *GoCompiler) compilePublicConstantNode(node *ast.PublicConstantNode) *goValue {
	return c.emitGetConst(value.ToSymbol(node.Value), c.typeOf(node))
}

func (c *GoCompiler) compilePrivateConstantNode(node *ast.PrivateConstantNode) *goValue {
	return c.emitGetConst(value.ToSymbol(node.Value), c.typeOf(node))
}

func (c *GoCompiler) compileMethodCallNode(node *ast.MethodCallNode, valueIsIgnored bool) *goValue {
	return c.compileMethodCall(
		node.Receiver,
		node.Op,
		node.MethodName,
		node.PositionalArguments,
		c.typeOf(node),
		node.Location(),
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileGenericMethodCallNode(node *ast.GenericMethodCallNode, valueIsIgnored bool) *goValue {
	return c.compileMethodCall(
		node.Receiver,
		node.Op,
		node.MethodName,
		node.PositionalArguments,
		c.typeOf(node),
		node.Location(),
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileMethodCall(receiver ast.ExpressionNode, op *token.Token, nameNode ast.IdentifierNode, args []ast.ExpressionNode, typ types.Type, location *position.Location, valueIsIgnored bool) *goValue {
	name := identifierToName(nameNode)

	switch op.Type {
	case token.QUESTION_DOT:
		receiverVal := c.compileExpression(receiver, false)
		resultVar := c.defineTmpGoLocal(goValueType)

		c.emit("if value.IsNil(%s) {\n", c.convertToValue(receiverVal).fetchValue())
		c.emit("%s = value.Nil\n", resultVar.name)
		c.emit("} else {\n")
		callResult := c.compileInnerMethodCall(receiverVal, c.typeOf(receiver), name, op, args, typ, location, valueIsIgnored)
		c.emit("%s = %s\n", resultVar.name, c.convertToValue(callResult).fetchValue())
		c.emit("}\n")

		return newGoValueWithLocal(resultVar, typ)
	case token.QUESTION_DOT_DOT:
		receiverVal := c.compileExpression(receiver, false)
		resultVar := c.defineTmpGoLocal(goValueType)

		c.emit("if value.IsNil(%s) {\n", c.convertToValue(receiverVal).fetchValue())
		c.emit("%s = value.Nil\n", resultVar.name)
		c.emit("} else {\n")
		c.compileInnerMethodCall(receiverVal, c.typeOf(receiver), name, op, args, typ, location, valueIsIgnored)
		c.emit("%s = %s\n", resultVar.name, receiverVal.fetchValue())
		c.emit("}\n")

		return newGoValueWithLocal(resultVar, typ)
	case token.DOT_DOT:
		receiverVal := c.compileExpression(receiver, false)
		resultVar := c.defineTmpGoLocal(goValueType)

		c.compileInnerMethodCall(receiverVal, c.typeOf(receiver), name, op, args, typ, location, valueIsIgnored)
		c.emit("%s = %s\n", resultVar.name, receiverVal.fetchValue())

		return newGoValueWithLocal(resultVar, typ)
	case token.DOT:
		receiverVal := c.compileExpression(receiver, false)

		return c.compileInnerMethodCall(receiverVal, c.typeOf(receiver), name, op, args, typ, location, valueIsIgnored)
	default:
		panic(fmt.Sprintf("invalid method call operator: %#v", op))
	}
}

func (c *GoCompiler) compileInnerMethodCall(receiver *goValue, receiverType types.Type, name string, op *token.Token, args []ast.ExpressionNode, typ types.Type, location *position.Location, valueIsIgnored bool) *goValue {
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

	return c.compileMethodCallWithArgNodes(receiver, receiverType, typ, name, args, location, valueIsIgnored)
}

func (c *GoCompiler) compileMethodCallWithArgNodes(receiver *goValue, receiverType, returnType types.Type, name string, args []ast.ExpressionNode, loc *position.Location, valueIsIgnored bool) *goValue {
	argsSlice := make([]*goValue, len(args)+1)

	argsSlice[0] = receiver
	for i, posArg := range args {
		argsSlice[i+1] = c.compileExpression(posArg, false)
	}

	return c.compileMethodCallWithLiteralArgValues(receiverType, returnType, name, argsSlice, loc, valueIsIgnored)
}

func (c *GoCompiler) compileMethodCallWithLiteralArgValues(receiverType, returnType types.Type, name string, args []*goValue, loc *position.Location, valueIsIgnored bool) *goValue {
	nameSym := c.emitSymbol(name)
	return c.compileMethodCallWithLiteralArgValuesAndName(
		receiverType,
		returnType,
		nameSym,
		name,
		args,
		loc,
		valueIsIgnored,
	)
}

func (c *GoCompiler) registerElkMethodName(methodName string) string {
	c.globalData.methodCache.Lock()

	var goName string
	if entry, ok := c.globalData.methodCache.GetUnsafe(methodName); ok {
		goName = entry.ident
	} else {
		goName = mangleGoIdentifier(methodName)
		c.globalData.methodCache.SetUnsafe(
			methodName,
			&nativeMethod{
				ident: goName,
			},
		)
	}

	c.globalData.methodCache.Unlock()

	return goName
}

func (c *GoCompiler) RegisterMethod(node *ast.MethodDefinitionNode) {
	method := c.typeOf(node).(*types.Method)
	c.registerElkMethodName(method.NamespacedName())
}

func (c *GoCompiler) compileOptimizedNativeMethodCall(receiverType, returnType types.Type, args []*goValue, name string, loc *position.Location, valueIsIgnored bool) *goValue {
	switch receiverType := receiverType.(type) {
	case types.Self:
		return c.compileOptimizedNativeMethodCall(
			c.checker.SelfType(),
			returnType,
			args,
			name,
			loc,
			valueIsIgnored,
		)
	case *types.SingletonClass:
		if receiverType.Children.Len() != 0 {
			break
		}

		// class has no children so method lookup can be static
		return c._compileOptimizedNativeMethodCall(
			receiverType,
			returnType,
			args,
			name,
			loc,
			valueIsIgnored,
		)
	case *types.Class:
		if receiverType.Children.Len() != 0 {
			break
		}

		// class has no children so method lookup can be static
		return c._compileOptimizedNativeMethodCall(
			receiverType,
			returnType,
			args,
			name,
			loc,
			valueIsIgnored,
		)
	case *types.Module:
		return c._compileOptimizedNativeMethodCall(
			receiverType,
			returnType,
			args,
			name,
			loc,
			valueIsIgnored,
		)
	}

	return nil
}

func (c *GoCompiler) generateGetNamespace(typ types.Namespace) string {
	switch typ := typ.(type) {
	case *types.SingletonClass:
		namespaceVal := c.emitGetConst(value.ToSymbol(typ.AttachedObject.Name()), types.Any{})
		return fmt.Sprintf("(%s).SingletonClass()", namespaceVal.value)
	case *types.Module:
		namespaceVal := c.emitGetConst(value.ToSymbol(typ.Name()), c.checker.Std(symbol.Module))
		return fmt.Sprintf("(%s).SingletonClass()", namespaceVal.value)
	case *types.Class:
		namespaceVal := c.emitGetConst(value.ToSymbol(typ.Name()), c.checker.Std(symbol.Class))
		return c.valueToNarrowerType(namespaceVal).value
	default:
		panic(fmt.Sprintf("invalid namespace: %T", typ))
	}
}

func (c *GoCompiler) _compileOptimizedNativeMethodCall(receiverType, returnType types.Type, args []*goValue, name string, loc *position.Location, valueIsIgnored bool) *goValue {
	method := c.checker.GetMethod(receiverType, value.ToSymbol(name), nil)
	c.globalData.methodCache.Lock()

	namespacedMethodName := method.NamespacedName()
	goMethodName, ok := c.globalData.methodCache.GetUnsafe(namespacedMethodName)
	if !ok {
		goIdent := mangleGoIdentifier(namespacedMethodName)
		nameSym := c.emitSymbol(name)
		goMethodName = &nativeMethod{
			ident: goIdent,
			init: fmt.Sprintf(
				"vm.MethodToFunc((%s).LookupMethod(%s))",
				c.generateGetNamespace(method.DefinedUnder),
				nameSym,
			),
		}
		c.globalData.methodCache.SetUnsafe(
			namespacedMethodName,
			goMethodName,
		)

		c.emitPackage(
			"var %s vm.NativeFunction // %s\n",
			goIdent,
			namespacedMethodName,
		)
	}

	c.globalData.methodCache.Unlock()

	var tmp *goLocal
	var tmpName string
	if valueIsIgnored {
		tmpName = "_"
	} else {
		if goMethodName.hasArgsSlice() {
			tmp = c.defineTmpGoLocal(goValueType)
		} else {
			tmp = c.defineTmpGoLocal(c.elkTypeToGoType(method.ReturnType, false))
		}
		tmpName = tmp.name
	}

	if goMethodName.hasArgsSlice() {
		callArgsVar := c.defineTmpGoLocal(value.FetchGoType("[]value.Value"))
		c.emit("%[1]s = value.ResizeNativeArgs(%[1]s, %d)\n", callArgsVar.name, len(args)+1)

		for i, posArg := range args {
			c.emit("%s[%d] = %s\n", callArgsVar.name, i, c.convertToValue(posArg).fetchValue())
		}

		c.registerErr()
		c.emitSetCallFrameLineNumber(loc)
		c.emit(
			"%s, err = %s(thread, %s) // receiver: %s, name: %s\n",
			tmpName,
			goMethodName.goIdent(),
			callArgsVar.name,
			types.Inspect(receiverType),
			name,
		)

		c.emitErrorPropagation()
	} else {
		c.registerErr()
		c.emitSetCallFrameLineNumber(loc)
		c.emit(
			"%s, err = %s(thread, %s",
			tmpName,
			goMethodName.goIdent(),
			c.convertToValue(args[0]).fetchValue(),
		)

		for i, arg := range args[1:] {
			param := method.Params[i]
			goParamType := c.elkTypeToGoType(param.Type, false)
			if goMethodName.hasArgsSlice() || param.IsOptional() || goParamType.Name == "value.Value" {
				c.emit(", %s", c.convertToValue(arg).fetchValue())
				continue
			}

			c.emit(", %s", c.valueToNarrowerType(arg).fetchValue())
		}

		c.emit(
			") // receiver: %s, name: %s\n",
			types.Inspect(receiverType),
			name,
		)
		c.emitErrorPropagation()
	}

	if valueIsIgnored {
		return nilGoValue
	}

	result := newGoValueWithLocal(
		tmp,
		returnType,
	)
	return c.narrowDownValue(result)
}

func (c *GoCompiler) compileMethodCallWithLiteralArgValuesAndName(receiverType, returnType types.Type, nameSym, name string, args []*goValue, loc *position.Location, valueIsIgnored bool) *goValue {
	result := c.compileOptimizedNativeMethodCall(
		receiverType,
		returnType,
		args,
		name,
		loc,
		valueIsIgnored,
	)
	if result != nil {
		return result
	}

	callCache := c.emitCallCache()

	var tmp *goLocal
	var tmpName string
	if valueIsIgnored {
		tmpName = "_"
	} else {
		tmp = c.defineTmpGoLocal(goValueType)
		tmpName = tmp.name
	}

	callArgsVar := c.defineTmpGoLocal(value.FetchGoType("[]value.Value"))
	c.emit("%[1]s = value.ResizeNativeArgs(%[1]s, %d)\n", callArgsVar.name, len(args)+1)

	for i, posArg := range args {
		c.emit("%s[%d] = %s\n", callArgsVar.name, i, c.convertToValue(posArg).fetchValue())
	}

	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit(
		"%s, err = thread.CallMethodByNameWithCache(%s, &%s, %s...) // receiver: %s, name: %s\n",
		tmpName,
		nameSym,
		callCache,
		callArgsVar.name,
		types.Inspect(receiverType),
		name,
	)
	c.emitErrorPropagation()

	if valueIsIgnored {
		return nilGoValue
	}

	result = newGoValueWithLocal(
		tmp,
		returnType,
	)
	return c.narrowDownValue(result)
}

func (c *GoCompiler) narrowDownValue(val *goValue) *goValue {
	if val.goType.Name != "value.Value" {
		return val
	}

	result := c.valueToNarrowerType(val)
	if result.goType.Name == "value.Value" {
		return result
	}

	narrowTmp := c.defineTmpGoLocal(result.goType)
	c.emit("%s = %s\n", narrowTmp.name, result.fetchValue())
	val.markFree()
	return newGoValueWithLocal(narrowTmp, result.elkType)
}

func (c *GoCompiler) compileModuleDeclarationNode(node *ast.ModuleDeclarationNode) *goValue {
	typ := c.typeOf(node).(*types.Module)
	return c.compileNamespaceDeclarationNode(fmt.Sprintf("module_%s", mangleGoIdentifier(typ.Name())), node.Body, typ, node.Location())
}

func (c *GoCompiler) compileInterfaceDeclarationNode(node *ast.InterfaceDeclarationNode) *goValue {
	typ := c.typeOf(node).(*types.Interface)
	return c.compileNamespaceDeclarationNode(fmt.Sprintf("interface_%s", mangleGoIdentifier(typ.Name())), node.Body, typ, node.Location())
}

func (c *GoCompiler) compileMixinDeclarationNode(node *ast.MixinDeclarationNode) *goValue {
	typ := c.typeOf(node).(*types.Mixin)
	return c.compileNamespaceDeclarationNode(fmt.Sprintf("mixin_%s", mangleGoIdentifier(typ.Name())), node.Body, typ, node.Location())
}

func (c *GoCompiler) compileClassDeclarationNode(node *ast.ClassDeclarationNode) *goValue {
	typ := c.typeOf(node).(*types.Class)
	return c.compileNamespaceDeclarationNode(fmt.Sprintf("class_%s", mangleGoIdentifier(typ.Name())), node.Body, typ, node.Location())
}

func (c *GoCompiler) compileNamespaceDeclarationNode(name string, body []ast.StatementNode, typ types.Namespace, loc *position.Location) *goValue {
	if len(body) <= 0 {
		return nilGoValue
	}

	classCompiler := NewGoCompiler(name, topLevelGoCompilerMode, loc, c.checker, c.globalData, c.output)
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
		namespaceVal := c.emitGetConst(value.ToSymbol(typ.AttachedObject.Name()), types.Any{})
		fmt.Fprintf(&funcBuffer, "self = value.Ref((%s).SingletonClass())\n", namespaceVal.fetchValue())
	case *types.Module:
		namespaceVal := c.emitGetConst(value.ToSymbol(typ.Name()), types.Any{})
		fmt.Fprintf(&funcBuffer, "self = value.Ref((%s).SingletonClass())\n", namespaceVal.fetchValue())
	default:
		namespaceVal := c.emitGetConst(value.ToSymbol(typ.Name()), types.Any{})
		fmt.Fprintf(&funcBuffer, "self = %s\n", c.convertToValue(namespaceVal).fetchValue())
	}

	c.emitPrependBytes(funcBuffer.Bytes())
	c.emit("}\n")
}

func (c *GoCompiler) compileAssignmentExpressionNode(node *ast.AssignmentExpressionNode) *goValue {
	switch n := node.Left.(type) {
	case *ast.PublicIdentifierNode:
		return c.localVariableAssignment(n.Value, node.Op, node.Right, c.typeOf(node.Left), c.typeOf(node), node.Location())
	case *ast.PrivateIdentifierNode:
		return c.localVariableAssignment(n.Value, node.Op, node.Right, c.typeOf(node.Left), c.typeOf(node), node.Location())
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

func (c *GoCompiler) localVariableAssignment(name string, operator *token.Token, right ast.ExpressionNode, varType, assignmentType types.Type, loc *position.Location) *goValue {
	switch operator.Type {
	case token.OR_OR_EQUAL:
		varIdent := c.compileLocalVariableAccess(name, varType)
		switch varIdent.goType.Name {
		case "value.Bool", "bool":
			c.emit("if !(%s) {\n", varIdent.value)
		default:
			c.emit("if value.Falsy(%s) {\n", c.convertToValue(varIdent).value)
		}

		rightVal := c.compileExpression(right, false)
		c.emit("%s = %s\n", varIdent.value, c.convertToValue(rightVal).fetchValue())

		c.emit("}\n")

		return varIdent
	case token.AND_AND_EQUAL:
		varIdent := c.compileLocalVariableAccess(name, varType)
		switch varIdent.goType.Name {
		case "value.Bool", "bool":
			c.emit("if %s {\n", varIdent.value)
		default:
			c.emit("if value.Truthy(%s) {\n", c.convertToValue(varIdent).value)
		}

		rightVal := c.compileExpression(right, false)
		c.emit("%s = %s\n", varIdent.value, c.convertToValue(rightVal).fetchValue())

		c.emit("}\n")

		return varIdent
	case token.QUESTION_QUESTION_EQUAL:
		varIdent := c.compileLocalVariableAccess(name, varType)
		c.emit("if value.IsNil(%s) {\n", c.convertToValue(varIdent).value)

		rightVal := c.compileExpression(right, false)
		c.emit("%s = %s\n", varIdent.value, c.convertToValue(rightVal).fetchValue())

		c.emit("}\n")

		return varIdent
	case token.EQUAL_OP:
		return c.setLocal(name, right)
	case token.COLON_EQUAL:
		c.defineLocal(name, assignmentType, c.elkTypeToGoType(assignmentType, false), loc)
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
	val := c.compileExpression(valueNode, false)
	return c.emitSetLocal(name, val)
}

func (c *GoCompiler) emitSetLocal(name string, val *goValue) *goValue {
	var variable *nativeElkLocal
	if local, ok := c.resolveLocal(name); ok {
		variable = local
	} else if upvalue, ok := c.resolveUpvalue(name); ok {
		variable = upvalue
	} else {
		panic(fmt.Sprintf("undefined local: %s\n", name))
	}

	ident := variable.goIdent()
	if variable.goLocal.goType.Name == "value.Value" {
		c.emit("%s = %s\n", ident, c.convertToValue(val).fetchValue())
	} else {
		c.emit("%s = %s\n", ident, c.valueToNarrowerType(val).fetchValue())
	}

	return newGoValueWithLocal(
		variable.goLocal,
		variable.elkType,
	)
}

func (c *GoCompiler) compileModifierIfExpression(condType conditionType, condition, then, els ast.ExpressionNode, typ types.Type, valueIsIgnored bool) *goValue {
	var elsFunc func() *goValue
	if els != nil {
		elsFunc = func() *goValue {
			return c.compileExpression(els, valueIsIgnored)
		}
	}

	return c.compileIfWithConditionExpression(
		condType,
		condition,
		func() *goValue {
			return c.compileExpression(then, valueIsIgnored)
		},
		elsFunc,
		typ,
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileIfExpression(condType conditionType, condition ast.ExpressionNode, then, els []ast.StatementNode, typ types.Type, valueIsIgnored bool) *goValue {
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
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileIfWithConditionExpression(condType conditionType, condition ast.ExpressionNode, then, els func() *goValue, typ types.Type, valueIsIgnored bool) *goValue {
	if result := resolve(condition, c.checker); !result.IsUndefined() {
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
		return c.compileExpression(condition, false)
	}

	return c.compileIf(
		condType,
		cond,
		then,
		els,
		typ,
		valueIsIgnored,
	)
}

type conditionType uint8

const (
	ifConditionType conditionType = iota
	unlessConditionType
	isNilConditionType
)

func (c *GoCompiler) compileIf(condType conditionType, condition, then, els func() *goValue, typ types.Type, valueIsIgnored bool) *goValue {
	c.enterScope("", defaultNativeElkScopeType)
	condVal := condition()

	var ifResultVar *goLocal
	if !valueIsIgnored {
		ifResultVar = c.defineTmpGoLocal(goValueType)
	}

	switch condVal.goType.Name {
	case "bool", "value.Bool":
		switch condType {
		case ifConditionType:
			c.emit("if %s {\n", condVal.fetchValue())
		case unlessConditionType:
			c.emit("if !(%s) {\n", condVal.fetchValue())
		default:
			panic(fmt.Sprintf("invalid if condition type: %d", condType))
		}

		thenVal := then()
		if !valueIsIgnored && !types.IsNever(thenVal.elkType) {
			c.emit("%s = %s\n", ifResultVar.name, c.convertToValue(thenVal).fetchValue())
		}
		c.emit("}")
	default:
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

		c.emit("if %s(%s) {\n", condFunc, c.convertToValue(condVal).fetchValue())

		thenVal := then()
		if !valueIsIgnored && !types.IsNever(thenVal.elkType) {
			c.emit("%s = %s\n", ifResultVar.name, c.convertToValue(thenVal).fetchValue())
		}
		thenVal.markFree()

		c.emit("}")
	}

	c.leaveScope()

	if els != nil {
		c.emit(" else {\n")
		elseVal := els()
		if !valueIsIgnored && !types.IsNever(elseVal.elkType) {
			c.emit("%s = %s\n", ifResultVar.name, c.convertToValue(elseVal).fetchValue())
		}
		elseVal.markFree()

		c.emit("}")
	} else if !valueIsIgnored {
		c.emit(" else {\n")
		c.emit("%s = value.Nil\n", ifResultVar.name)
		c.emit("}")
	}

	c.emit("\n")

	if valueIsIgnored {
		return nilGoValue
	}

	return newGoValueWithLocal(
		ifResultVar,
		typ,
	)
}

func (c *GoCompiler) compileLocalVariableAccess(name string, elkType types.Type) *goValue {
	if local, ok := c.resolveLocal(name); ok {
		return newGoValueWithLocal(
			local.goLocal,
			elkType,
		)
	}

	if upvalue, ok := c.resolveUpvalue(name); ok {
		return newGoValueWithLocal(
			upvalue.goLocal,
			elkType,
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

func (c *GoCompiler) emitBigFloat(val string) string {
	c.globalData.bigFloatCache.Lock()
	defer c.globalData.bigFloatCache.Unlock()

	bigFloat, ok := c.globalData.bigFloatCache.GetUnsafe(val)
	if ok {
		return bigFloat.goIdent()
	}

	bigFloat = &nativeBigFloat{
		id:  c.globalData.bigFloatCache.Len(),
		val: val,
	}
	c.globalData.bigFloatCache.SetUnsafe(val, bigFloat)
	ident := bigFloat.goIdent()
	c.emitPackage("var %s = value.ParseBigFloatPanic(%q)\n", ident, val)
	return ident
}

func (c *GoCompiler) emitBigInt(val string) string {
	c.globalData.bigIntCache.Lock()
	defer c.globalData.bigIntCache.Unlock()

	bigInt, ok := c.globalData.bigIntCache.GetUnsafe(val)
	if ok {
		return bigInt.goIdent()
	}

	bigInt = &nativeBigInt{
		id:  c.globalData.bigIntCache.Len(),
		val: val,
	}
	c.globalData.bigIntCache.SetUnsafe(val, bigInt)
	ident := bigInt.goIdent()
	c.emitPackage("var %s = value.ParseBigIntPanic(%q, 0)\n", ident, val)
	return ident
}

func (c *GoCompiler) emitSymbol(val string) string {
	c.globalData.symbolCache.Lock()
	defer c.globalData.symbolCache.Unlock()

	symbol, ok := c.globalData.symbolCache.GetUnsafe(val)
	if ok {
		return symbol.goIdent()
	}

	symbol = &nativeSymbol{
		id:  c.globalData.symbolCache.Len(),
		val: val,
	}
	c.globalData.symbolCache.SetUnsafe(val, symbol)
	ident := symbol.goIdent()
	c.emitPackage("var %s = value.ToSymbol(%q)\n", ident, val)
	return ident
}

func (c *GoCompiler) emitCachedRange(val value.Value, typ types.Type) *goValue {
	rangeSource := c.rangeToGoSource(val, typ, false)
	return c.emitCachedValue("range", val, rangeSource, typ)
}

func (c *GoCompiler) emitCachedValue(prefix string, val value.Value, goValue *goValue, elkType types.Type) *goValue {
	if goValue == nil {
		return nil
	}

	c.globalData.valueCache.Lock()
	defer c.globalData.valueCache.Unlock()

	inspect := val.Inspect()
	if nativeVal, ok := c.globalData.valueCache.GetUnsafe(inspect); ok {
		return newGoValue(
			nativeVal.goIdent(),
			nativeVal.elkType,
			nativeVal.goType,
		)
	}

	nativeVal := &nativeValue{
		ident: fmt.Sprintf("%s%d", prefix, c.globalData.valueCache.Len()),
		val:   val,
	}
	c.globalData.valueCache.SetUnsafe(inspect, nativeVal)
	ident := nativeVal.goIdent()

	c.emitPackage("var %s = %s\n", ident, goValue.fetchValue())

	nativeVal.goType = goValue.goType
	if elkType == nil {
		elkType = nativeVal.elkType
	}
	nativeVal.elkType = elkType

	return newGoValue(
		ident,
		elkType,
		nativeVal.goType,
	)
}

func (c *GoCompiler) emitCachedRegex(val *value.Regex, typ types.Type) *goValue {
	source := c.regexToGoSource(val, typ)
	return c.emitCachedValue("regex", val.ToValue(), source, typ)
}

func (c *GoCompiler) emitCachedArrayTuple(tuple value.ArrayTuple, typ types.Type) *goValue {
	tupleSource := c.arrayTupleToGoSource(tuple, typ, false)
	return c.emitCachedValue("arrtuple", tuple.ToValue(), tupleSource, typ)
}

func (c *GoCompiler) emitCachedHashRecord(record vm.HashRecord, typ types.Type) *goValue {
	recordSource := c.hashRecordToGoSource(record, typ, false)
	return c.emitCachedValue("hshrec", record.ToValue(), recordSource, typ)
}

func (c *GoCompiler) getTmpIdent() string {
	c.tmpLocalCounter++
	return fmt.Sprintf("t%d", c.tmpLocalCounter)
}

func (c *GoCompiler) defineTmpGoLocal(goType *value.GoType) *goLocal {
	for local := range c.goLocals.Values() {
		if !local.elkLocal && local.free && local.goType.Equal(goType) {
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
		c.emit("if err.IsNotUndefined() {\n")
		c.emitCaptureStackTrace()
		c.emit("thread.Panic(err)\n")
		c.emit("}\n")
	case methodGoCompilerMode, setterMethodGoCompilerMode, initMethodGoCompilerMode:
		c.emit("if err.IsNotUndefined() {\n")
		c.emitCaptureStackTrace()
		c.emit("return result, err\n")
		c.emit("}\n")
	default:
		c.emit("if err.IsNotUndefined() {\n")
		c.emitCaptureStackTrace()
		c.emit("return value.Undefined, err\n")
		c.emit("}\n")
	}
}

func (c *GoCompiler) emitCaptureStackTrace() {
	c.emit("thread.CaptureStackTrace()\n")
}

func (c *GoCompiler) emitCallCache() string {
	c.callCacheCounter++
	callCacheName := fmt.Sprintf("cc_%s_%d", c.FuncName, c.callCacheCounter)
	c.emitPackage("var %s = &value.CallCache{}\n", callCacheName)

	return callCacheName
}

func (c *GoCompiler) emitAddCallFrame(loc *position.Location) {
	funcNameSym := c.emitSymbol(c.FuncName)
	fileNameSym := c.emitSymbol(loc.FilePath)
	lineNumber := loc.StartPos.Line

	c.currentLineNumber = lineNumber
	c.registerGoLocal("callFrame", value.FetchGoType("*vm.CallFrame"))
	c.emit(
		"callFrame = thread.AddNativeCallFrame(%s, %s, %d)\n",
		funcNameSym,
		fileNameSym,
		lineNumber,
	)
	c.emit("defer thread.PopNativeCallFrame()\n")
}

func (c *GoCompiler) emitSetCallFrameLineNumber(loc *position.Location) {
	newLineNumber := loc.StartPos.Line
	if c.currentLineNumber == newLineNumber {
		return
	}

	c.currentLineNumber = newLineNumber
	c.emit(
		"callFrame.SetNativeLineNumber(%d)\n",
		newLineNumber,
	)
}

func (c *GoCompiler) registerErr() {
	c.registerGoLocal("err", goValueType)
}

func (c *GoCompiler) compileBinaryExpressionNode(node *ast.BinaryExpressionNode, valueIsIgnored bool) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	switch node.Op.Type {
	case token.PLUS:
		return c.compileAdd(node, valueIsIgnored)
	case token.MINUS:
		return c.compileSubtract(node, valueIsIgnored)
	case token.STAR:
		return c.compileMultiply(node, valueIsIgnored)
	case token.SLASH:
		return c.compileDivide(node, valueIsIgnored)
	case token.STAR_STAR:
		return c.compileExponentiate(node, valueIsIgnored)
	case token.LBITSHIFT:
		return c.compileLeftBitshift(node, valueIsIgnored)
	case token.LTRIPLE_BITSHIFT:
		return c.compileLogicalLeftBitshift(node, valueIsIgnored)
	case token.RBITSHIFT:
		return c.compileRightBitshift(node, valueIsIgnored)
	case token.RTRIPLE_BITSHIFT:
		return c.compileLogicalRightBitshift(node, valueIsIgnored)
	case token.AND:
		return c.compileBitwiseAnd(node, valueIsIgnored)
	case token.OR:
		return c.compileBitwiseOr(node, valueIsIgnored)
	case token.AND_TILDE:
		return c.compileBitwiseAndNot(node, valueIsIgnored)
	case token.XOR:
		return c.compileBitwiseXor(node, valueIsIgnored)
	case token.PERCENT:
		return c.compileModulo(node, valueIsIgnored)
	case token.LAX_EQUAL:
		return c.compileLaxEqual(node, valueIsIgnored)
	case token.LAX_NOT_EQUAL:
		return c.compileLaxNotEqual(node, valueIsIgnored)
	case token.EQUAL_EQUAL:
		return c.compileEqual(node, valueIsIgnored)
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
	case token.LESS_EQUAL:
		return c.compileLessEqual(node, valueIsIgnored)
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
	//  c.emitSetCallFrameLineNumber(node.Location())
	// 	c.emit("%s, err = thread.CallMethodByNameWithCache(symbol.OpLessThanEqual, &%s, %s, %s)", tmp.name, callCache, c.convertToValue(left), c.convertToValue(right))
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

func (c *GoCompiler) compileLessEqual(node *ast.BinaryExpressionNode, valueIsIgnored bool) *goValue {
	left := c.compileExpression(node.Left, false)
	narrowLeft := c.valueToNarrowerType(left)
	right := c.compileExpression(node.Right, false)
	typ := c.typeOf(node)

	switch narrowLeft.goType.Name {
	case "value.SmallInt", "*value.BigInt", "value.Float", "*value.BigFloat":
		return c.compileLessEqualCoercibleNumeric(narrowLeft, right, node.Location())
	case "value.Int64", "value.Int32", "value.Int16", "value.Int8",
		"value.UInt64", "value.UInt32", "value.UInt16", "value.UInt8",
		"value.Float64", "value.Float32":
		return c.compileLessEqualStrictNumeric(narrowLeft, right)
	}

	if c.checker.IsSubtype(left.elkType, c.checker.Std(symbol.S_BuiltinNumeric)) {
		if valueIsIgnored {
			c.registerErr()
			c.emitSetCallFrameLineNumber(node.Location())
			c.emit(
				"_, err = value.LessThanEqual(%s, %s)\n",
				c.convertToValue(left).fetchValue(),
				c.convertToValue(right).fetchValue(),
			)
			c.emitErrorPropagation()
			left.markFree()
			right.markFree()

			return nilGoValue
		}

		tmp := c.defineTmpGoLocal(value.FetchGoType("bool"))
		c.registerErr()
		c.emitSetCallFrameLineNumber(node.Location())
		c.emit(
			"%s, err = value.LessThanEqual(%s, %s)\n",
			tmp.name,
			c.convertToValue(left).fetchValue(),
			c.convertToValue(right).fetchValue(),
		)
		c.emitErrorPropagation()
		left.markFree()
		right.markFree()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	return c.compileMethodCallWithLiteralArgValuesAndName(
		left.elkType,
		typ,
		"symbol.OpLessThanEqual",
		"<=",
		[]*goValue{
			left,
			right,
		},
		node.Location(),
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileLessEqualCoercibleNumeric(left, right *goValue, loc *position.Location) *goValue {
	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("value.Bool((%s).LessThanEqualSmallInt(%s))", left.value, narrowRight.value),
			types.Bool{},
			value.FetchGoType("value.Bool"),
			left,
			narrowRight,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("value.Bool((%s).LessThanEqualFloat(%s))", left.value, narrowRight.value),
			types.Bool{},
			value.FetchGoType("value.Bool"),
			left,
			narrowRight,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("value.Bool((%s).LessThanEqualBigFloat(%s))", left.value, narrowRight.value),
			types.Bool{},
			value.FetchGoType("value.Bool"),
			left,
			narrowRight,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.Std(symbol.S_BuiltinNumeric)) {
		return newGoValueWithDependencies(
			fmt.Sprintf("value.Bool((%s).LessThanEqualInt(%s))", left.value, c.convertToValue(right).value),
			types.Bool{},
			value.FetchGoType("value.Bool"),
			left,
			right,
		)
	}

	tmp := c.defineTmpGoLocal(value.FetchGoType("value.Bool"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).LessThanEqual(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocals(
		fmt.Sprintf("value.Bool(%s)", tmp.name),
		types.Bool{},
		value.FetchGoType("value.Bool"),
		tmp,
	)
}

func (c *GoCompiler) compileLessEqualStrictNumeric(left, right *goValue) *goValue {
	return newGoValueWithDependencies(
		fmt.Sprintf("value.Bool((%s) <= (%s))", left.value, c.valueToNarrowerType(right).value),
		types.Bool{},
		value.FetchGoType("bool"),
		left,
		right,
	)
}

func (c *GoCompiler) compileDivide(node *ast.BinaryExpressionNode, valueIsIgnored bool) *goValue {
	left := c.compileExpression(node.Left, false)
	narrowLeft := c.valueToNarrowerType(left)
	right := c.compileExpression(node.Right, false)
	typ := c.typeOf(node)

	switch narrowLeft.goType.Name {
	case "value.SmallInt":
		return c.compileDivideSmallInt(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "*value.BigInt":
		return c.compileDivideBigInt(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Int64":
		return c.compileDivideInt64(narrowLeft, right, node.Location())
	case "value.Int32":
		return c.compileDivideInt32(narrowLeft, right, node.Location())
	case "value.Int16":
		return c.compileDivideInt16(narrowLeft, right, node.Location())
	case "value.Int8":
		return c.compileDivideInt8(narrowLeft, right, node.Location())
	case "value.UInt64":
		return c.compileDivideUInt64(narrowLeft, right, node.Location())
	case "value.UInt32":
		return c.compileDivideUInt32(narrowLeft, right, node.Location())
	case "value.UInt16":
		return c.compileDivideUInt16(narrowLeft, right, node.Location())
	case "value.UInt8":
		return c.compileDivideUInt8(narrowLeft, right, node.Location())
	case "value.Float":
		return c.compileDivideFloat(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "*value.BigFloat":
		return c.compileDivideBigFloat(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Float64":
		return c.compileDivideFloat64(narrowLeft, right)
	case "value.Float32":
		return c.compileDivideFloat32(narrowLeft, right)
	}

	if c.checker.IsSubtype(left.elkType, c.checker.Std(symbol.S_BuiltinDividable)) {
		if valueIsIgnored {
			c.registerErr()
			c.emitSetCallFrameLineNumber(node.Location())
			c.emit("_, err = value.DivideVal(%s, %s)\n", c.convertToValue(left).fetchValue(), c.convertToValue(right).fetchValue())
			c.emitErrorPropagation()
			return nilGoValue
		}

		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(node.Location())
		c.emit("%s, err = value.DivideVal(%s, %s)\n", tmp.name, c.convertToValue(left).fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	return c.compileMethodCallWithLiteralArgValuesAndName(
		left.elkType,
		typ,
		"symbol.OpDivide",
		"/",
		[]*goValue{
			left,
			right,
		},
		node.Location(),
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileDivideBigInt(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(loc)
		c.emit("%s, err = (%s).DivideSmallInt(%s)\n", tmp.name, left.fetchValue(), narrowRight.fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).DivideFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("value.Float"),
			left,
			right,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).DivideBigFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("*value.BigFloat"),
			left,
			right,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		if valueIsIgnored {
			c.registerErr()
			c.emitSetCallFrameLineNumber(loc)
			c.emit("_, err = (%s).DivideInt(%s)\n", left.fetchValue(), c.convertToValue(right).fetchValue())
			c.emitErrorPropagation()
			return nilGoValue
		}

		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(loc)
		c.emit("%s, err = (%s).DivideInt(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	if valueIsIgnored {
		c.registerErr()
		c.emitSetCallFrameLineNumber(loc)
		c.emit("_, err = (%s).DivideVal(%s)\n", left.fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).DivideVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileDivideSmallInt(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(loc)
		c.emit("%s, err = (%s).DivideSmallInt(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).DivideFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("value.Float"),
			left,
			right,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).DivideBigFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("*value.BigFloat"),
			left,
			right,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		if valueIsIgnored {
			c.registerErr()
			c.emitSetCallFrameLineNumber(loc)
			c.emit("_, err = (%s).DivideInt(%s)\n", left.fetchValue(), c.convertToValue(right).fetchValue())
			c.emitErrorPropagation()
			return nilGoValue
		}

		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(loc)
		c.emit("%s, err = (%s).DivideInt(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	if valueIsIgnored {
		c.registerErr()
		c.emitSetCallFrameLineNumber(loc)
		c.emit("_, err = (%s).DivideVal(%s)\n", left.fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).DivideVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileDivideFloat(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).DivideSmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("value.Float"),
			left,
			right,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).DivideFloat(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("value.Float"),
			left,
			right,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).DivideBigFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("*value.BigFloat"),
			left,
			right,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).DivideInt(%s)", left.value, c.convertToValue(right).value),
			left.elkType,
			value.FetchGoType("value.Float"),
			left,
			right,
		)
	}

	if valueIsIgnored {
		c.registerErr()
		c.emitSetCallFrameLineNumber(loc)
		c.emit("_, err = (%s).DivideVal(%s)\n", left.fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).DivideVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileDivideBigFloat(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).DivideSmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left,
			right,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).DivideFloat(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left,
			right,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).DivideBigFloat(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left,
			right,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).DivideInt(%s)", left.value, c.convertToValue(right).value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left,
			right,
		)
	}

	if valueIsIgnored {
		c.registerErr()
		c.emitSetCallFrameLineNumber(loc)
		c.emit("_, err = (%s).DivideVal(%s)\n", left.fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).DivideVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileDivideFloat64(left, right *goValue) *goValue {
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) / (%s)", left.value, c.valueToNarrowerType(right).value),
		left.elkType,
		value.FetchGoType("value.Float64"),
		left,
		right,
	)
}

func (c *GoCompiler) compileDivideFloat32(left, right *goValue) *goValue {
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) / (%s)", left.value, c.valueToNarrowerType(right).value),
		left.elkType,
		value.FetchGoType("value.Float32"),
		left,
		right,
	)
}

func (c *GoCompiler) compileDivideInt64(left, right *goValue, loc *position.Location) *goValue {
	tmp := c.defineTmpGoLocal(value.FetchGoType("value.Int64"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).DivideInt64(%s)\n", tmp.name, left.fetchValue(), c.valueToNarrowerType(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		left.elkType,
	)
}

func (c *GoCompiler) compileDivideInt32(left, right *goValue, loc *position.Location) *goValue {
	tmp := c.defineTmpGoLocal(value.FetchGoType("value.Int32"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).DivideInt32(%s)\n", tmp.name, left.fetchValue(), c.valueToNarrowerType(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		left.elkType,
	)
}

func (c *GoCompiler) compileDivideInt16(left, right *goValue, loc *position.Location) *goValue {
	tmp := c.defineTmpGoLocal(value.FetchGoType("value.Int16"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).DivideInt16(%s)\n", tmp.name, left.fetchValue(), c.valueToNarrowerType(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		left.elkType,
	)
}

func (c *GoCompiler) compileDivideInt8(left, right *goValue, loc *position.Location) *goValue {
	tmp := c.defineTmpGoLocal(value.FetchGoType("value.Int8"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).DivideInt8(%s)\n", tmp.name, left.fetchValue(), c.valueToNarrowerType(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		left.elkType,
	)
}

func (c *GoCompiler) compileDivideUInt64(left, right *goValue, loc *position.Location) *goValue {
	tmp := c.defineTmpGoLocal(value.FetchGoType("value.UInt64"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).DivideUInt64(%s)\n", tmp.name, left.fetchValue(), c.valueToNarrowerType(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		left.elkType,
	)
}

func (c *GoCompiler) compileDivideUInt32(left, right *goValue, loc *position.Location) *goValue {
	tmp := c.defineTmpGoLocal(value.FetchGoType("value.UInt32"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).DivideUInt32(%s)\n", tmp.name, left.fetchValue(), c.valueToNarrowerType(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		left.elkType,
	)
}

func (c *GoCompiler) compileDivideUInt16(left, right *goValue, loc *position.Location) *goValue {
	tmp := c.defineTmpGoLocal(value.FetchGoType("value.UInt16"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).DivideUInt16(%s)\n", tmp.name, left.fetchValue(), c.valueToNarrowerType(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		left.elkType,
	)
}

func (c *GoCompiler) compileDivideUInt8(left, right *goValue, loc *position.Location) *goValue {
	tmp := c.defineTmpGoLocal(value.FetchGoType("value.UInt8"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).DivideUInt8(%s)\n", tmp.name, left.fetchValue(), c.valueToNarrowerType(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		left.elkType,
	)
}

func (c *GoCompiler) compileExponentiate(node *ast.BinaryExpressionNode, valueIsIgnored bool) *goValue {
	left := c.compileExpression(node.Left, false)
	narrowLeft := c.valueToNarrowerType(left)
	right := c.compileExpression(node.Right, false)
	typ := c.typeOf(node)

	switch narrowLeft.goType.Name {
	case "value.SmallInt":
		return c.compileExponentiateSmallInt(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "*value.BigInt":
		return c.compileExponentiateBigInt(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Int64":
		return c.compileExponentiateInt64(narrowLeft, right)
	case "value.Int32":
		return c.compileExponentiateInt32(narrowLeft, right)
	case "value.Int16":
		return c.compileExponentiateInt16(narrowLeft, right)
	case "value.Int8":
		return c.compileExponentiateInt8(narrowLeft, right)
	case "value.UInt64":
		return c.compileExponentiateUInt64(narrowLeft, right)
	case "value.UInt32":
		return c.compileExponentiateUInt32(narrowLeft, right)
	case "value.UInt16":
		return c.compileExponentiateUInt16(narrowLeft, right)
	case "value.UInt8":
		return c.compileExponentiateUInt8(narrowLeft, right)
	case "value.Float":
		return c.compileExponentiateFloat(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "*value.BigFloat":
		return c.compileExponentiateBigFloat(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Float64":
		return c.compileExponentiateFloat64(narrowLeft, right)
	case "value.Float32":
		return c.compileExponentiateFloat32(narrowLeft, right)
	}

	return c.compileMethodCallWithLiteralArgValuesAndName(
		left.elkType,
		typ,
		"symbol.OpExponentiate",
		"**",
		[]*goValue{
			left,
			right,
		},
		node.Location(),
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileExponentiateBigInt(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ExponentiateSmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ExponentiateFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("value.Float"),
			left,
			right,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ExponentiateBigFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("*value.BigFloat"),
			left,
			right,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ExponentiateInt(%s)", left.value, c.convertToValue(right).value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).ExponentiateVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileExponentiateSmallInt(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ExponentiateSmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("value.SmallInt"),
			left, right,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ExponentiateFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("value.Float"),
			left, right,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ExponentiateBigFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, right,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ExponentiateInt(%s)", left.value, c.convertToValue(right).value),
			left.elkType,
			goValueType,
			left, right,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).ExponentiateVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileExponentiateFloat(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ExponentiateSmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("value.Float"),
			left, right,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ExponentiateFloat(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("value.Float"),
			left, right,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ExponentiateBigFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, right,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ExponentiateInt(%s)", left.value, c.convertToValue(right).value),
			left.elkType,
			value.FetchGoType("value.Float"),
			left, right,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).ExponentiateVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileExponentiateBigFloat(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ExponentiateSmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, right,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ExponentiateFloat(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, right,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ExponentiateBigFloat(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, right,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ExponentiateInt(%s)", left.value, c.convertToValue(right).value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, right,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).ExponentiateVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileExponentiateFloat64(left, right *goValue) *goValue {
	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s).ExponentiateFloat64(%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.Float64"),
		left,
		right,
	)
}

func (c *GoCompiler) compileExponentiateFloat32(left, right *goValue) *goValue {
	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s).ExponentiateFloat32(%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.Float32"),
		left,
		right,
	)
}

func (c *GoCompiler) compileExponentiateInt64(left, right *goValue) *goValue {
	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s).ExponentiateInt64(%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.Int64"),
		left,
		right,
	)
}

func (c *GoCompiler) compileExponentiateInt32(left, right *goValue) *goValue {
	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s).ExponentiateInt32(%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.Int32"),
		left,
		right,
	)
}

func (c *GoCompiler) compileExponentiateInt16(left, right *goValue) *goValue {
	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s).ExponentiateInt16(%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.Int16"),
		left,
		right,
	)
}

func (c *GoCompiler) compileExponentiateInt8(left, right *goValue) *goValue {
	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s).ExponentiateInt8(%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.Int8"),
		left,
		right,
	)
}

func (c *GoCompiler) compileExponentiateUInt64(left, right *goValue) *goValue {
	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s).ExponentiateUInt64(%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.UInt64"),
		left,
		right,
	)
}

func (c *GoCompiler) compileExponentiateUInt32(left, right *goValue) *goValue {
	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s).ExponentiateUInt32(%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.UInt32"),
		left,
		right,
	)
}

func (c *GoCompiler) compileExponentiateUInt16(left, right *goValue) *goValue {
	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s).ExponentiateUInt16(%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.UInt16"),
		left,
		right,
	)
}

func (c *GoCompiler) compileExponentiateUInt8(left, right *goValue) *goValue {
	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s).ExponentiateUInt8(%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.UInt8"),
		left,
		right,
	)
}

func (c *GoCompiler) compileBitwiseAnd(node *ast.BinaryExpressionNode, valueIsIgnored bool) *goValue {
	left := c.compileExpression(node.Left, false)
	narrowLeft := c.valueToNarrowerType(left)
	right := c.compileExpression(node.Right, false)
	typ := c.typeOf(node)

	switch narrowLeft.goType.Name {
	case "value.Int64", "value.Int32", "value.Int16", "value.Int8",
		"value.UInt64", "value.UInt32", "value.UInt16", "value.UInt8":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s) & (%s)", narrowLeft.value, c.valueToNarrowerType(right).value),
			typ,
			narrowLeft.goType,
			left,
			right,
		)
	}

	if c.checker.IsSubtype(left.elkType, c.checker.Std(symbol.S_BuiltinInt)) {
		if valueIsIgnored {
			return nilGoValue
		}

		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(node.Location())
		c.emit("%s, err = value.BitwiseAndVal(%s, %s)\n", tmp.name, c.convertToValue(left).fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	return c.compileMethodCallWithLiteralArgValuesAndName(
		left.elkType,
		typ,
		"symbol.OpAnd",
		"&",
		[]*goValue{
			left,
			right,
		},
		node.Location(),
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileBitwiseAndNot(node *ast.BinaryExpressionNode, valueIsIgnored bool) *goValue {
	left := c.compileExpression(node.Left, false)
	narrowLeft := c.valueToNarrowerType(left)
	right := c.compileExpression(node.Right, false)
	typ := c.typeOf(node)

	switch narrowLeft.goType.Name {
	case "value.Int64", "value.Int32", "value.Int16", "value.Int8",
		"value.UInt64", "value.UInt32", "value.UInt16", "value.UInt8":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s) &^ (%s)", narrowLeft.value, c.valueToNarrowerType(right).value),
			typ,
			narrowLeft.goType,
			left,
			right,
		)
	}

	if c.checker.IsSubtype(left.elkType, c.checker.Std(symbol.S_BuiltinInt)) {
		if valueIsIgnored {
			return nilGoValue
		}

		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(node.Location())
		c.emit("%s, err = value.BitwiseAndNotVal(%s, %s)\n", tmp.name, c.convertToValue(left).fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	return c.compileMethodCallWithLiteralArgValuesAndName(
		left.elkType,
		typ,
		"symbol.OpAndNot",
		"&~",
		[]*goValue{
			left,
			right,
		},
		node.Location(),
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileBitwiseOr(node *ast.BinaryExpressionNode, valueIsIgnored bool) *goValue {
	left := c.compileExpression(node.Left, false)
	narrowLeft := c.valueToNarrowerType(left)
	right := c.compileExpression(node.Right, false)
	typ := c.typeOf(node)

	switch narrowLeft.goType.Name {
	case "value.Int64", "value.Int32", "value.Int16", "value.Int8",
		"value.UInt64", "value.UInt32", "value.UInt16", "value.UInt8":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s) | (%s)", narrowLeft.value, c.valueToNarrowerType(right).value),
			typ,
			narrowLeft.goType,
			left,
			right,
		)
	}

	if c.checker.IsSubtype(left.elkType, c.checker.Std(symbol.S_BuiltinInt)) {
		if valueIsIgnored {
			return nilGoValue
		}

		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(node.Location())
		c.emit("%s, err = value.BitwiseOrVal(%s, %s)\n", tmp.name, c.convertToValue(left).fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	return c.compileMethodCallWithLiteralArgValuesAndName(
		left.elkType,
		typ,
		"symbol.OpOr",
		"|",
		[]*goValue{
			left,
			right,
		},
		node.Location(),
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileBitwiseXor(node *ast.BinaryExpressionNode, valueIsIgnored bool) *goValue {
	left := c.compileExpression(node.Left, false)
	narrowLeft := c.valueToNarrowerType(left)
	right := c.compileExpression(node.Right, false)
	typ := c.typeOf(node)

	switch narrowLeft.goType.Name {
	case "value.Int64", "value.Int32", "value.Int16", "value.Int8",
		"value.UInt64", "value.UInt32", "value.UInt16", "value.UInt8":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s) ^ (%s)", narrowLeft.value, c.valueToNarrowerType(right).value),
			typ,
			narrowLeft.goType,
			left,
			right,
		)
	}

	if c.checker.IsSubtype(left.elkType, c.checker.Std(symbol.S_BuiltinInt)) {
		if valueIsIgnored {
			return nilGoValue
		}

		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(node.Location())
		c.emit("%s, err = value.BitwiseXorVal(%s, %s)\n", tmp.name, c.convertToValue(left).fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	return c.compileMethodCallWithLiteralArgValuesAndName(
		left.elkType,
		typ,
		"symbol.OpXor",
		"^",
		[]*goValue{
			left,
			right,
		},
		node.Location(),
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileLeftBitshift(node *ast.BinaryExpressionNode, valueIsIgnored bool) *goValue {
	left := c.compileExpression(node.Left, false)
	narrowLeft := c.valueToNarrowerType(left)
	right := c.compileExpression(node.Right, false)
	typ := c.typeOf(node)

	switch narrowLeft.goType.Name {
	case "value.SmallInt":
		return c.compileLeftBitshiftSmallInt(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "*value.BigInt":
		return c.compileLeftBitshiftBigInt(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Int64":
		return c.compileLeftBitshiftInt64(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Int32":
		return c.compileLeftBitshiftInt32(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Int16":
		return c.compileLeftBitshiftInt16(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Int8":
		return c.compileLeftBitshiftInt8(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.UInt64":
		return c.compileLeftBitshiftUInt64(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.UInt32":
		return c.compileLeftBitshiftUInt32(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.UInt16":
		return c.compileLeftBitshiftUInt16(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.UInt8":
		return c.compileLeftBitshiftUInt8(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	}

	if c.checker.IsSubtype(left.elkType, c.checker.Std(symbol.S_BuiltinLogicBitshiftable)) {
		if valueIsIgnored {
			return nilGoValue
		}

		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(node.Location())
		c.emit("%s, err = value.LeftBitshiftVal(%s, %s)\n", tmp.name, c.convertToValue(left).fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	return c.compileMethodCallWithLiteralArgValuesAndName(
		left.elkType,
		typ,
		"symbol.OpLeftBitshift",
		"<<",
		[]*goValue{
			left,
			right,
		},
		node.Location(),
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileLogicalLeftBitshift(node *ast.BinaryExpressionNode, valueIsIgnored bool) *goValue {
	left := c.compileExpression(node.Left, false)
	narrowLeft := c.valueToNarrowerType(left)
	right := c.compileExpression(node.Right, false)
	typ := c.typeOf(node)

	switch narrowLeft.goType.Name {
	case "value.Int64":
		return c.compileLeftBitshiftInt64(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Int32":
		return c.compileLeftBitshiftInt32(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Int16":
		return c.compileLeftBitshiftInt16(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Int8":
		return c.compileLeftBitshiftInt8(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.UInt64":
		return c.compileLeftBitshiftUInt64(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.UInt32":
		return c.compileLeftBitshiftUInt32(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.UInt16":
		return c.compileLeftBitshiftUInt16(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.UInt8":
		return c.compileLeftBitshiftUInt8(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	}

	if c.checker.IsSubtype(left.elkType, c.checker.Std(symbol.S_BuiltinLogicBitshiftable)) {
		if valueIsIgnored {
			return nilGoValue
		}

		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(node.Location())
		c.emit("%s, err = value.LogicalLeftBitshiftVal(%s, %s)\n", tmp.name, c.convertToValue(left).fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	return c.compileMethodCallWithLiteralArgValuesAndName(
		left.elkType,
		typ,
		"symbol.OpLogicalLeftBitshift",
		"<<<",
		[]*goValue{
			left,
			right,
		},
		node.Location(),
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileRightBitshift(node *ast.BinaryExpressionNode, valueIsIgnored bool) *goValue {
	left := c.compileExpression(node.Left, false)
	narrowLeft := c.valueToNarrowerType(left)
	right := c.compileExpression(node.Right, false)
	typ := c.typeOf(node)

	switch narrowLeft.goType.Name {
	case "value.SmallInt":
		return c.compileRightBitshiftSmallInt(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "*value.BigInt":
		return c.compileRightBitshiftBigInt(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Int64":
		return c.compileRightBitshiftInt64(narrowLeft, right, typ, node.Location())
	case "value.Int32":
		return c.compileRightBitshiftInt32(narrowLeft, right, typ, node.Location())
	case "value.Int16":
		return c.compileRightBitshiftInt16(narrowLeft, right, typ, node.Location())
	case "value.Int8":
		return c.compileRightBitshiftInt8(narrowLeft, right, typ, node.Location())
	case "value.UInt64":
		return c.compileRightBitshiftUInt64(narrowLeft, right, typ, node.Location())
	case "value.UInt32":
		return c.compileRightBitshiftUInt32(narrowLeft, right, typ, node.Location())
	case "value.UInt16":
		return c.compileRightBitshiftUInt16(narrowLeft, right, typ, node.Location())
	case "value.UInt8":
		return c.compileRightBitshiftUInt8(narrowLeft, right, typ, node.Location())
	}

	if c.checker.IsSubtype(left.elkType, c.checker.Std(symbol.S_BuiltinLogicBitshiftable)) {
		if valueIsIgnored {
			return nilGoValue
		}

		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(node.Location())
		c.emit("%s, err = value.RightBitshiftVal(%s, %s)\n", tmp.name, c.convertToValue(left).fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	return c.compileMethodCallWithLiteralArgValuesAndName(
		left.elkType,
		typ,
		"symbol.OpRightBitshift",
		">>",
		[]*goValue{
			left,
			right,
		},
		node.Location(),
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileLogicalRightBitshift(node *ast.BinaryExpressionNode, valueIsIgnored bool) *goValue {
	left := c.compileExpression(node.Left, false)
	narrowLeft := c.valueToNarrowerType(left)
	right := c.compileExpression(node.Right, false)
	typ := c.typeOf(node)

	switch narrowLeft.goType.Name {
	case "value.Int64":
		return c.compileLogicalRightBitshiftInt64(narrowLeft, right, typ, node.Location())
	case "value.Int32":
		return c.compileLogicalRightBitshiftInt32(narrowLeft, right, typ, node.Location())
	case "value.Int16":
		return c.compileLogicalRightBitshiftInt16(narrowLeft, right, typ, node.Location())
	case "value.Int8":
		return c.compileLogicalRightBitshiftInt8(narrowLeft, right, typ, node.Location())
	case "value.UInt64":
		return c.compileRightBitshiftUInt64(narrowLeft, right, typ, node.Location())
	case "value.UInt32":
		return c.compileRightBitshiftUInt32(narrowLeft, right, typ, node.Location())
	case "value.UInt16":
		return c.compileRightBitshiftUInt16(narrowLeft, right, typ, node.Location())
	case "value.UInt8":
		return c.compileRightBitshiftUInt8(narrowLeft, right, typ, node.Location())
	}

	if c.checker.IsSubtype(left.elkType, c.checker.Std(symbol.S_BuiltinLogicBitshiftable)) {
		if valueIsIgnored {
			return nilGoValue
		}

		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(node.Location())
		c.emit("%s, err = value.LogicalRightBitshiftVal(%s, %s)\n", tmp.name, c.convertToValue(left).fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	return c.compileMethodCallWithLiteralArgValuesAndName(
		left.elkType,
		typ,
		"symbol.OpLogicalRightBitshift",
		">>>",
		[]*goValue{
			left,
			right,
		},
		node.Location(),
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileLeftBitshiftSmallInt(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).LeftBitshiftSmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "*value.BigInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).LeftBitshiftBigInt(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.Int64":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).LeftBitshiftInt64(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.Int32":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).LeftBitshiftInt32(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.Int16":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).LeftBitshiftInt16(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.Int8":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).LeftBitshiftInt8(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.UInt64":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).LeftBitshiftUInt64(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.UInt32":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).LeftBitshiftUInt32(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.UInt16":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).LeftBitshiftUInt16(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.UInt8":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).LeftBitshiftUInt8(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).LeftBitshiftVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileRightBitshiftSmallInt(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).RightBitshiftSmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "*value.BigInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).RightBitshiftBigInt(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.Int64":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).RightBitshiftInt64(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.Int32":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).RightBitshiftInt32(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.Int16":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).RightBitshiftInt16(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.Int8":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).RightBitshiftInt8(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.UInt64":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).RightBitshiftUInt64(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.UInt32":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).RightBitshiftUInt32(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.UInt16":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).RightBitshiftUInt16(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.UInt8":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).RightBitshiftUInt8(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).RightBitshiftVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileLeftBitshiftBigInt(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).LeftBitshiftSmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "*value.BigInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).LeftBitshiftBigInt(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.Int64":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).LeftBitshiftInt64(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.Int32":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).LeftBitshiftInt32(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.Int16":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).LeftBitshiftInt16(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.Int8":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).LeftBitshiftInt8(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.UInt64":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).LeftBitshiftUInt64(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.UInt32":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).LeftBitshiftUInt32(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.UInt16":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).LeftBitshiftUInt16(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.UInt8":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).LeftBitshiftUInt8(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).LeftBitshiftVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileLeftBitshiftInt64(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(value.FetchGoType("value.Int64"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = value.StrictIntLeftBitshift(%s, %s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileLeftBitshiftInt32(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(value.FetchGoType("value.Int32"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = value.StrictIntLeftBitshift(%s, %s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileLeftBitshiftInt16(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(value.FetchGoType("value.Int16"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = value.StrictIntLeftBitshift(%s, %s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileLeftBitshiftInt8(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(value.FetchGoType("value.Int8"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = value.StrictIntLeftBitshift(%s, %s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileLeftBitshiftUInt64(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(value.FetchGoType("value.UInt64"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = value.StrictIntLeftBitshift(%s, %s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileLeftBitshiftUInt32(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(value.FetchGoType("value.UInt32"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = value.StrictIntLeftBitshift(%s, %s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileLeftBitshiftUInt16(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(value.FetchGoType("value.UInt16"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = value.StrictIntLeftBitshift(%s, %s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileLeftBitshiftUInt8(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(value.FetchGoType("value.UInt8"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = value.StrictIntLeftBitshift(%s, %s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileRightBitshiftBigInt(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).RightBitshiftSmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "*value.BigInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).RightBitshiftBigInt(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.Int64":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).RightBitshiftInt64(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.Int32":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).RightBitshiftInt32(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.Int16":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).RightBitshiftInt16(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.Int8":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).RightBitshiftInt8(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.UInt64":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).RightBitshiftUInt64(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.UInt32":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).RightBitshiftUInt32(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.UInt16":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).RightBitshiftUInt16(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.UInt8":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).RightBitshiftUInt8(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).RightBitshiftVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileLogicalRightBitshiftInt64(left, right *goValue, typ types.Type, loc *position.Location) *goValue {
	tmp := c.defineTmpGoLocal(value.FetchGoType("value.Int64"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = value.StrictIntLogicalRightBitshift(%s, %s, value.LogicalRightShift64)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileLogicalRightBitshiftInt32(left, right *goValue, typ types.Type, loc *position.Location) *goValue {
	tmp := c.defineTmpGoLocal(value.FetchGoType("value.Int32"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = value.StrictIntLogicalRightBitshift(%s, %s, value.LogicalRightShift32)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileLogicalRightBitshiftInt16(left, right *goValue, typ types.Type, loc *position.Location) *goValue {
	tmp := c.defineTmpGoLocal(value.FetchGoType("value.Int16"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = value.StrictIntLogicalRightBitshift(%s, %s, value.LogicalRightShift16)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileLogicalRightBitshiftInt8(left, right *goValue, typ types.Type, loc *position.Location) *goValue {
	tmp := c.defineTmpGoLocal(value.FetchGoType("value.Int8"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = value.StrictIntLogicalRightBitshift(%s, %s, value.LogicalRightShift8)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileRightBitshiftInt64(left, right *goValue, typ types.Type, loc *position.Location) *goValue {
	tmp := c.defineTmpGoLocal(value.FetchGoType("value.Int64"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = value.StrictIntRightBitshift(%s, %s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileRightBitshiftInt32(left, right *goValue, typ types.Type, loc *position.Location) *goValue {
	tmp := c.defineTmpGoLocal(value.FetchGoType("value.Int32"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = value.StrictIntRightBitshift(%s, %s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileRightBitshiftInt16(left, right *goValue, typ types.Type, loc *position.Location) *goValue {
	tmp := c.defineTmpGoLocal(value.FetchGoType("value.Int16"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = value.StrictIntRightBitshift(%s, %s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileRightBitshiftInt8(left, right *goValue, typ types.Type, loc *position.Location) *goValue {
	tmp := c.defineTmpGoLocal(value.FetchGoType("value.Int8"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = value.StrictIntRightBitshift(%s, %s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileRightBitshiftUInt64(left, right *goValue, typ types.Type, loc *position.Location) *goValue {
	tmp := c.defineTmpGoLocal(value.FetchGoType("value.UInt64"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = value.StrictIntRightBitshift(%s, %s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileRightBitshiftUInt32(left, right *goValue, typ types.Type, loc *position.Location) *goValue {
	tmp := c.defineTmpGoLocal(value.FetchGoType("value.UInt32"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = value.StrictIntRightBitshift(%s, %s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileRightBitshiftUInt16(left, right *goValue, typ types.Type, loc *position.Location) *goValue {
	tmp := c.defineTmpGoLocal(value.FetchGoType("value.UInt16"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = value.StrictIntRightBitshift(%s, %s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileRightBitshiftUInt8(left, right *goValue, typ types.Type, loc *position.Location) *goValue {
	tmp := c.defineTmpGoLocal(value.FetchGoType("value.UInt8"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = value.StrictIntRightBitshift(%s, %s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileMultiply(node *ast.BinaryExpressionNode, valueIsIgnored bool) *goValue {
	left := c.compileExpression(node.Left, false)
	narrowLeft := c.valueToNarrowerType(left)
	right := c.valueToNarrowerType(c.compileExpression(node.Right, false))
	typ := c.typeOf(node)

	switch narrowLeft.goType.Name {
	case "value.SmallInt":
		return c.compileMultiplySmallInt(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "*value.BigInt":
		return c.compileMultiplyBigInt(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Int64":
		return c.compileMultiplyInt64(narrowLeft, right, valueIsIgnored)
	case "value.Int32":
		return c.compileMultiplyInt32(narrowLeft, right, valueIsIgnored)
	case "value.Int16":
		return c.compileMultiplyInt16(narrowLeft, right, valueIsIgnored)
	case "value.Int8":
		return c.compileMultiplyInt8(narrowLeft, right, valueIsIgnored)
	case "value.UInt64":
		return c.compileMultiplyUInt64(narrowLeft, right, valueIsIgnored)
	case "value.UInt32":
		return c.compileMultiplyUInt32(narrowLeft, right, valueIsIgnored)
	case "value.UInt16":
		return c.compileMultiplyUInt16(narrowLeft, right, valueIsIgnored)
	case "value.UInt8":
		return c.compileMultiplyUInt8(narrowLeft, right, valueIsIgnored)
	case "value.Float":
		return c.compileMultiplyFloat(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "*value.BigFloat":
		return c.compileMultiplyBigFloat(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Float64":
		return c.compileMultiplyFloat64(narrowLeft, right, valueIsIgnored)
	case "value.Float32":
		return c.compileMultiplyFloat32(narrowLeft, right, valueIsIgnored)
	case "value.String":
		return c.compileMultiplyString(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Char":
		return c.compileMultiplyChar(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	}

	if c.checker.IsSubtype(left.elkType, c.checker.Std(symbol.S_BuiltinMultipliable)) {
		if valueIsIgnored {
			return nilGoValue
		}

		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(node.Location())
		c.emit(
			"%s, err = value.MultiplyVal(%s, %s)\n",
			tmp.name,
			c.convertToValue(left).fetchValue(),
			c.convertToValue(right).fetchValue(),
		)
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	return c.compileMethodCallWithLiteralArgValuesAndName(
		left.elkType,
		typ,
		"symbol.OpMultiply",
		"*",
		[]*goValue{
			left,
			right,
		},
		node.Location(),
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileMultiplyBigInt(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).MultiplySmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).MultiplyFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("value.Float"),
			left,
			right,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).MultiplyBigFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("*value.BigFloat"),
			left,
			right,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).MultiplyInt(%s)", left.value, c.convertToValue(right).value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).MultiplyVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileMultiplySmallInt(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).MultiplySmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("value.SmallInt"),
			left,
			right,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).MultiplyFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("value.Float"),
			left,
			right,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).MultiplyBigFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("*value.BigFloat"),
			left,
			right,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).MultiplyInt(%s)", left.value, c.convertToValue(right).value),
			left.elkType,
			goValueType,
			left,
			right,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).MultiplyVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileMultiplyFloat(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).MultiplySmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("value.Float"),
			left,
			right,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).MultiplyFloat(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("value.Float"),
			left,
			right,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).MultiplyBigFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("*value.BigFloat"),
			left,
			right,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).MultiplyInt(%s)", left.value, c.convertToValue(right).value),
			left.elkType,
			value.FetchGoType("value.Float"),
			left,
			right,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).MultiplyVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileMultiplyBigFloat(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).MultiplySmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left,
			right,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).MultiplyFloat(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left,
			right,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).MultiplyBigFloat(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left,
			right,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).MultiplyInt(%s)", left.value, c.convertToValue(right).value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left,
			right,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).MultiplyVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileMultiplyString(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		tmp := c.defineTmpGoLocal(value.FetchGoType("value.String"))
		c.registerErr()
		c.emitSetCallFrameLineNumber(loc)
		c.emit("%s, err = (%s).RepeatSmallInt(%s)\n", tmp.name, left.fetchValue(), narrowRight.fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	tmp := c.defineTmpGoLocal(value.FetchGoType("value.String"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).Repeat(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileMultiplyChar(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		tmp := c.defineTmpGoLocal(value.FetchGoType("value.String"))
		c.registerErr()
		c.emitSetCallFrameLineNumber(loc)
		c.emit("%s, err = (%s).RepeatSmallInt(%s)\n", tmp.name, left.fetchValue(), narrowRight.fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	tmp := c.defineTmpGoLocal(value.FetchGoType("value.String"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).Repeat(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileMultiplyFloat64(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) * (%s)", left.value, c.valueToNarrowerType(right).value),
		left.elkType,
		value.FetchGoType("value.Float64"),
		left, right,
	)
}

func (c *GoCompiler) compileMultiplyFloat32(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) * (%s)", left.value, c.valueToNarrowerType(right).value),
		left.elkType,
		value.FetchGoType("value.Float32"),
		left, right,
	)
}

func (c *GoCompiler) compileMultiplyInt64(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) * (%s)", left.value, c.valueToNarrowerType(right).value),
		left.elkType,
		value.FetchGoType("value.Int64"),
		left, right,
	)
}

func (c *GoCompiler) compileMultiplyInt32(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) * (%s)", left.value, c.valueToNarrowerType(right).value),
		left.elkType,
		value.FetchGoType("value.Int32"),
		left, right,
	)
}

func (c *GoCompiler) compileMultiplyInt16(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) * (%s)", left.value, c.valueToNarrowerType(right).value),
		left.elkType,
		value.FetchGoType("value.Int16"),
		left, right,
	)
}

func (c *GoCompiler) compileMultiplyInt8(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) * (%s)", left.value, c.valueToNarrowerType(right).value),
		left.elkType,
		value.FetchGoType("value.Int8"),
		left, right,
	)
}

func (c *GoCompiler) compileMultiplyUInt64(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) * (%s)", left.value, c.valueToNarrowerType(right).value),
		left.elkType,
		value.FetchGoType("value.UInt64"),
		left, right,
	)
}

func (c *GoCompiler) compileMultiplyUInt32(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) * (%s)", left.value, c.valueToNarrowerType(right).value),
		left.elkType,
		value.FetchGoType("value.UInt32"),
		left, right,
	)
}

func (c *GoCompiler) compileMultiplyUInt16(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) * (%s)", left.value, c.valueToNarrowerType(right).value),
		left.elkType,
		value.FetchGoType("value.UInt16"),
		left, right,
	)
}

func (c *GoCompiler) compileMultiplyUInt8(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) * (%s)", left.value, c.valueToNarrowerType(right).value),
		left.elkType,
		value.FetchGoType("value.UInt8"),
		left, right,
	)
}

func (c *GoCompiler) compileSubtract(node *ast.BinaryExpressionNode, valueIsIgnored bool) *goValue {
	left := c.compileExpression(node.Left, false)
	narrowLeft := c.valueToNarrowerType(left)
	right := c.valueToNarrowerType(c.compileExpression(node.Right, false))
	typ := c.typeOf(node)

	switch narrowLeft.goType.Name {
	case "value.SmallInt":
		return c.compileSubtractSmallInt(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "*value.BigInt":
		return c.compileSubtractBigInt(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Int64":
		return c.compileSubtractInt64(narrowLeft, right, valueIsIgnored)
	case "value.Int32":
		return c.compileSubtractInt32(narrowLeft, right, valueIsIgnored)
	case "value.Int16":
		return c.compileSubtractInt16(narrowLeft, right, valueIsIgnored)
	case "value.Int8":
		return c.compileSubtractInt8(narrowLeft, right, valueIsIgnored)
	case "value.UInt64":
		return c.compileSubtractUInt64(narrowLeft, right, valueIsIgnored)
	case "value.UInt32":
		return c.compileSubtractUInt32(narrowLeft, right, valueIsIgnored)
	case "value.UInt16":
		return c.compileSubtractUInt16(narrowLeft, right, valueIsIgnored)
	case "value.UInt8":
		return c.compileSubtractUInt8(narrowLeft, right, valueIsIgnored)
	case "value.Float":
		return c.compileSubtractFloat(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Float64":
		return c.compileSubtractFloat64(narrowLeft, right, valueIsIgnored)
	case "value.Float32":
		return c.compileSubtractFloat32(narrowLeft, right, valueIsIgnored)
	case "*value.BigFloat":
		return c.compileSubtractBigFloat(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	}

	if c.checker.IsSubtype(left.elkType, c.checker.Std(symbol.S_BuiltinSubtractable)) {
		if valueIsIgnored {
			return nilGoValue
		}

		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(node.Location())
		c.emit("%s, err = value.SubtractVal(%s, %s)\n", tmp.name, c.convertToValue(left).fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	return c.compileMethodCallWithLiteralArgValuesAndName(
		left.elkType,
		typ,
		"symbol.OpSubtract",
		"-",
		[]*goValue{
			left,
			right,
		},
		node.Location(),
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileSubtractBigInt(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).SubtractSmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left, narrowRight,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).SubtractFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("value.Float"),
			left, narrowRight,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).SubtractBigFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, narrowRight,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		rightVal := c.convertToValue(right)
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).SubtractInt(%s)", left.value, rightVal.value),
			left.elkType,
			goValueType,
			left, rightVal,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).SubtractVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileSubtractSmallInt(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).SubtractSmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("value.SmallInt"),
			left, narrowRight,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).SubtractFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("value.Float"),
			left, narrowRight,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).SubtractBigFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, narrowRight,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		rightVal := c.convertToValue(right)
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).SubtractInt(%s)", left.value, rightVal.value),
			left.elkType,
			goValueType,
			left, rightVal,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).SubtractVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileSubtractFloat(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).SubtractSmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("value.Float"),
			left, narrowRight,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).SubtractFloat(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("value.Float"),
			left, narrowRight,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).SubtractBigFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, narrowRight,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		rightVal := c.convertToValue(right)
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).SubtractInt(%s)", left.value, rightVal.value),
			left.elkType,
			value.FetchGoType("value.Float"),
			left, rightVal,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).SubtractVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileSubtractBigFloat(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).SubtractSmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, narrowRight,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).SubtractFloat(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, narrowRight,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).SubtractBigFloat(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, narrowRight,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		rightVal := c.convertToValue(right)
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).SubtractInt(%s)", left.value, rightVal.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, rightVal,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).SubtractVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileSubtractFloat64(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) - (%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.Float64"),
		left, narrowRight,
	)
}

func (c *GoCompiler) compileSubtractFloat32(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) - (%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.Float32"),
		left, narrowRight,
	)
}

func (c *GoCompiler) compileSubtractInt64(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) - (%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.Int64"),
		left, narrowRight,
	)
}

func (c *GoCompiler) compileSubtractInt32(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) - (%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.Int32"),
		left, narrowRight,
	)
}

func (c *GoCompiler) compileSubtractInt16(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) - (%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.Int16"),
		left, narrowRight,
	)
}

func (c *GoCompiler) compileSubtractInt8(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) - (%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.Int8"),
		left, narrowRight,
	)
}

func (c *GoCompiler) compileSubtractUInt64(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) - (%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.UInt64"),
		left, narrowRight,
	)
}

func (c *GoCompiler) compileSubtractUInt32(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) - (%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.UInt32"),
		left, narrowRight,
	)
}

func (c *GoCompiler) compileSubtractUInt16(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) - (%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.UInt16"),
		left, narrowRight,
	)
}

func (c *GoCompiler) compileSubtractUInt8(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s) - (%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.UInt8"),
		left, narrowRight,
	)
}

func (c *GoCompiler) compileLaxEqual(node *ast.BinaryExpressionNode, valueIsIgnored bool) *goValue {
	left := c.valueToNarrowerType(c.compileExpression(node.Left, false))
	narrowLeft := c.valueToNarrowerType(left)
	right := c.compileExpression(node.Right, false)
	typ := c.typeOf(node)

	switch narrowLeft.goType.Name {
	case "value.String", "*value.Regex", "value.Symbol", "value.Char",
		"value.SmallInt", "*value.BigInt", "value.Float", "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("value.Bool((%s).LaxEqual(%s))", narrowLeft.value, c.convertToValue(right).value),
			typ,
			value.FetchGoType("value.Bool"),
			narrowLeft, right,
		)
	case "value.Int64", "value.Int32", "value.Int16", "value.Int8":
		return c.compileLaxEqualStrictSignedInt(narrowLeft, right)
	case "value.UInt64", "value.UInt32", "value.UInt16", "value.UInt8":
		return c.compileLaxEqualStrictUnsignedInt(narrowLeft, right)
	case "value.Float64", "value.Float32":
		return c.compileLaxEqualStrictFloat(narrowLeft, right)
	}

	if c.checker.IsSubtype(left.elkType, c.checker.Std(symbol.S_BuiltinEquatable)) {
		if valueIsIgnored {
			return nilGoValue
		}

		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(node.Location())
		c.emit("%s, err = value.LaxEqualVal(%s, %s)\n", tmp.name, c.convertToValue(left).fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	return c.compileMethodCallWithLiteralArgValuesAndName(
		left.elkType,
		typ,
		"symbol.OpLaxEqual",
		"=~",
		[]*goValue{
			left,
			right,
		},
		node.Location(),
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileLaxNotEqual(node *ast.BinaryExpressionNode, valueIsIgnored bool) *goValue {
	left := c.valueToNarrowerType(c.compileExpression(node.Left, false))
	defer left.markFree()

	narrowLeft := c.valueToNarrowerType(left)

	right := c.compileExpression(node.Right, false)
	defer right.markFree()

	typ := c.typeOf(node)

	switch narrowLeft.goType.Name {
	case "value.String", "*value.Regex", "value.Symbol", "value.Char",
		"value.SmallInt", "*value.BigInt", "value.Float", "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("value.Bool(!(%s).LaxEqual(%s))", narrowLeft.value, c.convertToValue(right).value),
			typ,
			value.FetchGoType("value.Bool"),
			left, right,
		)
	case "value.Int64", "value.Int32", "value.Int16", "value.Int8":
		return c.compileLaxNotEqualStrictSignedInt(narrowLeft, right)
	case "value.UInt64", "value.UInt32", "value.UInt16", "value.UInt8":
		return c.compileLaxNotEqualStrictUnsignedInt(narrowLeft, right)
	case "value.Float64", "value.Float32":
		return c.compileLaxNotEqualStrictFloat(narrowLeft, right)
	}

	if c.checker.IsSubtype(left.elkType, c.checker.Std(symbol.S_BuiltinEquatable)) {
		return newGoValueWithDependencies(
			fmt.Sprintf("value.ToNotBool(value.LaxEqual(%s, %s))", c.convertToValue(left).value, c.convertToValue(right).value),
			typ,
			value.FetchGoType("value.Bool"),
			left, right,
		)
	}

	result := c.compileMethodCallWithLiteralArgValuesAndName(
		left.elkType,
		typ,
		"symbol.OpLaxEqual",
		"=~",
		[]*goValue{
			left,
			right,
		},
		node.Location(),
		valueIsIgnored,
	)

	return result.newGoValue(
		fmt.Sprintf("value.ToNotBool(%s)", result.value),
		types.Bool{},
		value.FetchGoType("value.Bool"),
	)
}

func (c *GoCompiler) compileLaxEqualStrictSignedInt(left, right *goValue) *goValue {
	return newGoValueWithDependencies(
		fmt.Sprintf("value.Bool(value.StrictSignedIntLaxEqual(%s, %s))", left.value, c.valueToNarrowerType(right).value),
		types.Bool{},
		value.FetchGoType("value.Bool"),
		left, right,
	)
}

func (c *GoCompiler) compileLaxNotEqualStrictSignedInt(left, right *goValue) *goValue {
	return newGoValueWithDependencies(
		fmt.Sprintf("value.Bool(!value.StrictSignedIntLaxEqual(%s, %s))", left.value, c.valueToNarrowerType(right).value),
		types.Bool{},
		value.FetchGoType("value.Bool"),
		left, right,
	)
}

func (c *GoCompiler) compileLaxEqualStrictUnsignedInt(left, right *goValue) *goValue {
	return newGoValueWithDependencies(
		fmt.Sprintf("value.Bool(value.StrictUnsignedIntLaxEqual(%s, %s))", left.value, c.valueToNarrowerType(right).value),
		types.Bool{},
		value.FetchGoType("value.Bool"),
		left, right,
	)
}

func (c *GoCompiler) compileLaxNotEqualStrictUnsignedInt(left, right *goValue) *goValue {
	return newGoValueWithDependencies(
		fmt.Sprintf("value.Bool(!value.StrictUnsignedIntLaxEqual(%s, %s))", left.value, c.valueToNarrowerType(right).value),
		types.Bool{},
		value.FetchGoType("value.Bool"),
		left, right,
	)
}

func (c *GoCompiler) compileLaxEqualStrictFloat(left, right *goValue) *goValue {
	return newGoValueWithDependencies(
		fmt.Sprintf("value.Bool(value.StrictFloatLaxEqual(%s, %s))", left.value, c.valueToNarrowerType(right).value),
		types.Bool{},
		value.FetchGoType("value.Bool"),
		left, right,
	)
}

func (c *GoCompiler) compileLaxNotEqualStrictFloat(left, right *goValue) *goValue {
	return newGoValueWithDependencies(
		fmt.Sprintf("value.Bool(!value.StrictFloatLaxEqual(%s, %s))", left.value, c.valueToNarrowerType(right).value),
		types.Bool{},
		value.FetchGoType("value.Bool"),
		left, right,
	)
}

func (c *GoCompiler) compileEqual(node *ast.BinaryExpressionNode, valueIsIgnored bool) *goValue {
	left := c.valueToNarrowerType(c.compileExpression(node.Left, false))
	narrowLeft := c.valueToNarrowerType(left)
	right := c.compileExpression(node.Right, false)
	typ := c.typeOf(node)

	switch narrowLeft.goType.Name {
	case "value.SmallInt":
		return c.compileEqualSmallInt(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "*value.BigInt":
		return c.compileEqualBigInt(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Int64":
		return c.compileAddInt64(narrowLeft, right, valueIsIgnored)
	case "value.Int32":
		return c.compileAddInt32(narrowLeft, right, valueIsIgnored)
	case "value.Int16":
		return c.compileAddInt16(narrowLeft, right, valueIsIgnored)
	case "value.Int8":
		return c.compileAddInt8(narrowLeft, right, valueIsIgnored)
	case "value.UInt64":
		return c.compileAddUInt64(narrowLeft, right, valueIsIgnored)
	case "value.UInt32":
		return c.compileAddUInt32(narrowLeft, right, valueIsIgnored)
	case "value.UInt16":
		return c.compileAddUInt16(narrowLeft, right, valueIsIgnored)
	case "value.UInt8":
		return c.compileAddUInt8(narrowLeft, right, valueIsIgnored)
	case "value.Float":
		return c.compileAddFloat(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Float64":
		return c.compileAddFloat64(narrowLeft, right, valueIsIgnored)
	case "value.Float32":
		return c.compileAddFloat32(narrowLeft, right, valueIsIgnored)
	case "*value.BigFloat":
		return c.compileAddBigFloat(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.String":
		return c.compileAddString(narrowLeft, right, node.Location(), valueIsIgnored)
	case "value.Char":
		return c.compileAddChar(narrowLeft, right, node.Location(), valueIsIgnored)
	}

	if c.checker.IsSubtype(left.elkType, c.checker.Std(symbol.S_BuiltinAddable)) {
		if valueIsIgnored {
			return nilGoValue
		}

		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(node.Location())
		c.emit("%s, err = value.EqualVal(%s, %s)\n", tmp.name, c.convertToValue(left).fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	return c.compileMethodCallWithLiteralArgValuesAndName(
		left.elkType,
		typ,
		"symbol.OpEqual",
		"==",
		[]*goValue{
			left,
			right,
		},
		node.Location(),
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileEqualBigInt(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("value.Bool((%s).EqualSmallInt(%s))", left.value, narrowRight.value),
			typ,
			value.FetchGoType("value.Bool"),
			left, narrowRight,
		)
	case "*value.BigInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("value.Bool((%s).EqualBigInt(%s))", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left, narrowRight,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newGoValueWithDependencies(
			fmt.Sprintf("value.Bool((%s).EqualInt(%s))", left.value, c.convertToValue(right).value),
			left.elkType,
			goValueType,
			left, right,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).EqualVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileEqualSmallInt(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("value.Bool((%s).EqualSmallInt(%s))", left.value, narrowRight.value),
			typ,
			goValueType,
			left, narrowRight,
		)
	case "*value.BigInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("value.Bool((%s).EqualBigInt(%s))", left.value, narrowRight.value),
			typ,
			value.FetchGoType("value.Bool"),
			left, narrowRight,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		return newGoValueWithDependencies(
			fmt.Sprintf("value.Bool((%s).EqualInt(%s))", left.value, c.convertToValue(right).value),
			typ,
			value.FetchGoType("value.Bool"),
			left, right,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).EqualVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileModulo(node *ast.BinaryExpressionNode, valueIsIgnored bool) *goValue {
	left := c.valueToNarrowerType(c.compileExpression(node.Left, false))
	narrowLeft := c.valueToNarrowerType(left)
	right := c.compileExpression(node.Right, false)
	typ := c.typeOf(node)

	switch narrowLeft.goType.Name {
	case "value.SmallInt":
		return c.compileModuloSmallInt(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "*value.BigInt":
		return c.compileModuloBigInt(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Int64":
		return c.compileModuloInt64(narrowLeft, right, node.Location(), valueIsIgnored)
	case "value.Int32":
		return c.compileModuloInt32(narrowLeft, right, node.Location(), valueIsIgnored)
	case "value.Int16":
		return c.compileModuloInt16(narrowLeft, right, node.Location(), valueIsIgnored)
	case "value.Int8":
		return c.compileModuloInt8(narrowLeft, right, node.Location(), valueIsIgnored)
	case "value.UInt":
		return c.compileModuloUInt(narrowLeft, right, node.Location(), valueIsIgnored)
	case "value.UInt64":
		return c.compileModuloUInt64(narrowLeft, right, node.Location(), valueIsIgnored)
	case "value.UInt32":
		return c.compileModuloUInt32(narrowLeft, right, node.Location(), valueIsIgnored)
	case "value.UInt16":
		return c.compileModuloUInt16(narrowLeft, right, node.Location(), valueIsIgnored)
	case "value.UInt8":
		return c.compileModuloUInt8(narrowLeft, right, node.Location(), valueIsIgnored)
	case "value.Float":
		return c.compileModuloFloat(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Float64":
		return c.compileModuloFloat64(narrowLeft, right, valueIsIgnored)
	case "value.Float32":
		return c.compileModuloFloat32(narrowLeft, right, valueIsIgnored)
	case "*value.BigFloat":
		return c.compileModuloBigFloat(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	}

	if c.checker.IsSubtype(left.elkType, c.checker.Std(symbol.S_BuiltinNumeric)) {
		if valueIsIgnored {
			return nilGoValue
		}

		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(node.Location())
		c.emit("%s, err = value.ModuloVal(%s, %s)\n", tmp.name, c.convertToValue(left).fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	return c.compileMethodCallWithLiteralArgValuesAndName(
		left.elkType,
		typ,
		"symbol.OpModulo",
		"%",
		[]*goValue{
			left,
			right,
		},
		node.Location(),
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileModuloFloat(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ModuloSmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("value.Float"),
			left, narrowRight,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ModuloFloat(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("value.Float"),
			left, narrowRight,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ModuloBigFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, narrowRight,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		valRight := c.convertToValue(right)
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ModuloInt(%s)", left.value, valRight.value),
			left.elkType,
			goValueType,
			left, valRight,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).ModuloVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileModuloBigFloat(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ModuloSmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, narrowRight,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ModuloFloat(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, narrowRight,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ModuloBigFloat(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, narrowRight,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		valRight := c.convertToValue(right)
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ModuloInt(%s)", left.value, valRight.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, valRight,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).ModuloVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileModuloSmallInt(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(loc)
		c.emit("%s, err = (%s).ModuloSmallInt(%s)\n", tmp.name, left.fetchValue(), narrowRight.fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	case "*value.BigInt":
		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(loc)
		c.emit("%s, err = (%s).ModuloBigInt(%s)\n", tmp.name, left.fetchValue(), narrowRight.fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ModuloFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("value.Float"),
			left, narrowRight,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ModuloBigFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, narrowRight,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		valRight := c.convertToValue(right)
		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(loc)
		c.emit("%s, err = (%s).ModuloInt(%s)\n", tmp.name, left.fetchValue(), valRight.fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	valRight := c.convertToValue(right)
	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).ModuloVal(%s)\n", tmp.name, left.fetchValue(), valRight.fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileModuloBigInt(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(loc)
		c.emit("%s, err = (%s).ModuloSmallInt(%s)\n", tmp.name, left.fetchValue(), narrowRight.fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ModuloFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("value.Float"),
			left, narrowRight,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).ModuloBigFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, narrowRight,
		)
	case "*value.BigInt":
		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(loc)
		c.emit("%s, err = (%s).ModuloBigInt(%s)\n", tmp.name, left.fetchValue(), narrowRight.fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		valRight := c.convertToValue(right)
		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(loc)
		c.emit("%s, err = (%s).ModuloInt(%s)\n", tmp.name, left.fetchValue(), valRight.fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	valRight := c.convertToValue(right)
	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).ModuloVal(%s)\n", tmp.name, left.fetchValue(), valRight.fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileModuloInt64(left, right *goValue, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(left.goType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).ModuloInt64(%s)\n", tmp.name, left.fetchValue(), c.valueToNarrowerType(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		left.elkType,
	)
}

func (c *GoCompiler) compileModuloInt32(left, right *goValue, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(left.goType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).ModuloInt32(%s)\n", tmp.name, left.fetchValue(), c.valueToNarrowerType(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		left.elkType,
	)
}

func (c *GoCompiler) compileModuloInt16(left, right *goValue, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(left.goType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).ModuloInt16(%s)\n", tmp.name, left.fetchValue(), c.valueToNarrowerType(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		left.elkType,
	)
}

func (c *GoCompiler) compileModuloInt8(left, right *goValue, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(left.goType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).ModuloInt8(%s)\n", tmp.name, left.fetchValue(), c.valueToNarrowerType(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		left.elkType,
	)
}

func (c *GoCompiler) compileModuloUInt64(left, right *goValue, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(left.goType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).ModuloUInt64(%s)\n", tmp.name, left.fetchValue(), c.valueToNarrowerType(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		left.elkType,
	)
}

func (c *GoCompiler) compileModuloUInt(left, right *goValue, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(left.goType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).ModuloUInt(%s)\n", tmp.name, left.fetchValue(), c.valueToNarrowerType(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		left.elkType,
	)
}

func (c *GoCompiler) compileModuloUInt32(left, right *goValue, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(left.goType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).ModuloUInt32(%s)\n", tmp.name, left.fetchValue(), c.valueToNarrowerType(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		left.elkType,
	)
}

func (c *GoCompiler) compileModuloUInt16(left, right *goValue, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(left.goType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).ModuloUInt16(%s)\n", tmp.name, left.fetchValue(), c.valueToNarrowerType(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		left.elkType,
	)
}

func (c *GoCompiler) compileModuloUInt8(left, right *goValue, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	tmp := c.defineTmpGoLocal(left.goType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).ModuloUInt8(%s)\n", tmp.name, left.fetchValue(), c.valueToNarrowerType(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		left.elkType,
	)
}

func (c *GoCompiler) compileModuloFloat64(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s).ModuloFloat64(%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.Float64"),
		left, narrowRight,
	)
}

func (c *GoCompiler) compileModuloFloat32(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	return newGoValueWithDependencies(
		fmt.Sprintf("(%s).ModuloFloat32(%s)", left.value, narrowRight.value),
		left.elkType,
		value.FetchGoType("value.Float32"),
		left, narrowRight,
	)
}

func (c *GoCompiler) compileAdd(node *ast.BinaryExpressionNode, valueIsIgnored bool) *goValue {
	left := c.valueToNarrowerType(c.compileExpression(node.Left, false))
	narrowLeft := c.valueToNarrowerType(left)
	right := c.compileExpression(node.Right, false)
	typ := c.typeOf(node)

	switch narrowLeft.goType.Name {
	case "value.SmallInt":
		return c.compileAddSmallInt(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "*value.BigInt":
		return c.compileAddBigInt(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Int64":
		return c.compileAddInt64(narrowLeft, right, valueIsIgnored)
	case "value.Int32":
		return c.compileAddInt32(narrowLeft, right, valueIsIgnored)
	case "value.Int16":
		return c.compileAddInt16(narrowLeft, right, valueIsIgnored)
	case "value.Int8":
		return c.compileAddInt8(narrowLeft, right, valueIsIgnored)
	case "value.UInt64":
		return c.compileAddUInt64(narrowLeft, right, valueIsIgnored)
	case "value.UInt32":
		return c.compileAddUInt32(narrowLeft, right, valueIsIgnored)
	case "value.UInt16":
		return c.compileAddUInt16(narrowLeft, right, valueIsIgnored)
	case "value.UInt8":
		return c.compileAddUInt8(narrowLeft, right, valueIsIgnored)
	case "value.Float":
		return c.compileAddFloat(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.Float64":
		return c.compileAddFloat64(narrowLeft, right, valueIsIgnored)
	case "value.Float32":
		return c.compileAddFloat32(narrowLeft, right, valueIsIgnored)
	case "*value.BigFloat":
		return c.compileAddBigFloat(narrowLeft, right, typ, node.Location(), valueIsIgnored)
	case "value.String":
		return c.compileAddString(narrowLeft, right, node.Location(), valueIsIgnored)
	case "value.Char":
		return c.compileAddChar(narrowLeft, right, node.Location(), valueIsIgnored)
	}

	if c.checker.IsSubtype(left.elkType, c.checker.Std(symbol.S_BuiltinAddable)) {
		if valueIsIgnored {
			return nilGoValue
		}

		tmp := c.defineTmpGoLocal(goValueType)
		c.registerErr()
		c.emitSetCallFrameLineNumber(node.Location())
		c.emit("%s, err = value.AddVal(%s, %s)\n", tmp.name, c.convertToValue(left).fetchValue(), c.convertToValue(right).fetchValue())
		c.emitErrorPropagation()

		return newGoValueWithLocal(
			tmp,
			typ,
		)
	}

	return c.compileMethodCallWithLiteralArgValuesAndName(
		left.elkType,
		typ,
		"symbol.OpAdd",
		"+",
		[]*goValue{
			left,
			right,
		},
		node.Location(),
		valueIsIgnored,
	)
}

func (c *GoCompiler) compileAddBigInt(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).AddSmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left, narrowRight,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).AddFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("value.Float"),
			left, narrowRight,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).AddBigFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, narrowRight,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		converted := c.convertToValue(right)
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).AddInt(%s)", left.value, converted.value),
			left.elkType,
			goValueType,
			left, converted,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).AddVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileAddSmallInt(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).AddSmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			goValueType,
			left, narrowRight,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).AddFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("value.Float"),
			left, narrowRight,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).AddBigFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, narrowRight,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		converted := c.convertToValue(right)
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).AddInt(%s)", left.value, converted.value),
			left.elkType,
			goValueType,
			left, converted,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).AddVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileAddFloat(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).AddSmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("value.Float"),
			left, narrowRight,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).AddFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("value.Float"),
			left, narrowRight,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).AddBigFloat(%s)", left.value, narrowRight.value),
			narrowRight.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, narrowRight,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		converted := c.convertToValue(right)
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).AddInt(%s)", left.value, converted.value),
			left.elkType,
			value.FetchGoType("value.Float"),
			left, converted,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).AddVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileAddBigFloat(left, right *goValue, typ types.Type, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.SmallInt":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).AddSmallInt(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, narrowRight,
		)
	case "value.Float":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).AddFloat(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, narrowRight,
		)
	case "*value.BigFloat":
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).AddBigFloat(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, narrowRight,
		)
	}

	if c.checker.IsSubtype(right.elkType, c.checker.StdInt()) {
		converted := c.convertToValue(right)
		return newGoValueWithDependencies(
			fmt.Sprintf("(%s).AddInt(%s)", left.value, converted.value),
			left.elkType,
			value.FetchGoType("*value.BigFloat"),
			left, converted,
		)
	}

	tmp := c.defineTmpGoLocal(goValueType)
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = (%s).AddVal(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		typ,
	)
}

func (c *GoCompiler) compileAddString(left, right *goValue, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.String":
		return newGoValueWithDependencies(
			fmt.Sprintf("%s.ConcatString(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("value.String"),
			left, narrowRight,
		)
	case "value.Char":
		return newGoValueWithDependencies(
			fmt.Sprintf("%s.ConcatChar(%s)", left.value, narrowRight.value),
			left.elkType,
			value.FetchGoType("value.String"),
			left, narrowRight,
		)
	}

	tmp := c.defineTmpGoLocal(value.FetchGoType("value.String"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = %s.Concat(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		left.elkType,
	)
}

func (c *GoCompiler) compileAddChar(left, right *goValue, loc *position.Location, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	switch narrowRight.goType.Name {
	case "value.String":
		return newGoValueWithDependencies(
			fmt.Sprintf("%s.ConcatString(%s)", left.value, narrowRight.value),
			right.elkType,
			value.FetchGoType("value.String"),
			left, narrowRight,
		)
	case "value.Char":
		return newGoValueWithDependencies(
			fmt.Sprintf("%s.ConcatChar(%s)", left.value, narrowRight.value),
			c.checker.Std(symbol.String),
			value.FetchGoType("value.String"),
			left, narrowRight,
		)
	}

	tmp := c.defineTmpGoLocal(value.FetchGoType("value.String"))
	c.registerErr()
	c.emitSetCallFrameLineNumber(loc)
	c.emit("%s, err = %s.Concat(%s)\n", tmp.name, left.fetchValue(), c.convertToValue(right).fetchValue())
	c.emitErrorPropagation()

	return newGoValueWithLocal(
		tmp,
		c.checker.Std(symbol.String),
	)
}

func (c *GoCompiler) compileAddFloat64(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	return newGoValue(
		fmt.Sprintf("(%s) + (%s)", left.fetchValue(), c.valueToNarrowerType(right).fetchValue()),
		left.elkType,
		value.FetchGoType("value.Float64"),
	)
}

func (c *GoCompiler) compileAddFloat32(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	return newGoValue(
		fmt.Sprintf("(%s) + (%s)", left.fetchValue(), c.valueToNarrowerType(right).fetchValue()),
		left.elkType,
		value.FetchGoType("value.Float32"),
	)
}

func (c *GoCompiler) compileAddInt64(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	return newGoValue(
		fmt.Sprintf("(%s) + (%s)", left.fetchValue(), narrowRight.fetchValue()),
		left.elkType,
		value.FetchGoType("value.Int64"),
	)
}

func (c *GoCompiler) compileAddInt32(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	narrowRight := c.valueToNarrowerType(right)
	return newGoValue(
		fmt.Sprintf("(%s) + (%s)", left.fetchValue(), narrowRight.fetchValue()),
		left.elkType,
		value.FetchGoType("value.Int32"),
	)
}

func (c *GoCompiler) compileAddInt16(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	return newGoValue(
		fmt.Sprintf("(%s) + (%s)", left.fetchValue(), c.valueToNarrowerType(right).fetchValue()),
		left.elkType,
		value.FetchGoType("value.Int16"),
	)
}

func (c *GoCompiler) compileAddInt8(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	return newGoValue(
		fmt.Sprintf("(%s) + (%s)", left.fetchValue(), c.valueToNarrowerType(right).fetchValue()),
		left.elkType,
		value.FetchGoType("value.Int8"),
	)
}

func (c *GoCompiler) compileAddUInt64(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	return newGoValue(
		fmt.Sprintf("(%s) + (%s)", left.fetchValue(), c.valueToNarrowerType(right).fetchValue()),
		left.elkType,
		value.FetchGoType("value.UInt64"),
	)
}

func (c *GoCompiler) compileAddUInt32(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	return newGoValue(
		fmt.Sprintf("(%s) + (%s)", left.fetchValue(), c.valueToNarrowerType(right).fetchValue()),
		left.elkType,
		value.FetchGoType("value.UInt32"),
	)
}

func (c *GoCompiler) compileAddUInt16(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	return newGoValue(
		fmt.Sprintf("(%s) + (%s)", left.fetchValue(), c.valueToNarrowerType(right).fetchValue()),
		left.elkType,
		value.FetchGoType("value.UInt16"),
	)
}

func (c *GoCompiler) compileAddUInt8(left, right *goValue, valueIsIgnored bool) *goValue {
	if valueIsIgnored {
		return nilGoValue
	}

	return newGoValue(
		fmt.Sprintf("(%s) + (%s)", left.fetchValue(), c.valueToNarrowerType(right).fetchValue()),
		left.elkType,
		value.FetchGoType("value.UInt8"),
	)
}

func (c *GoCompiler) resolve(node ast.ExpressionNode) *goValue {
	result := resolve(node, c.checker)
	if result.IsUndefined() {
		return nil
	}

	return c.valueToGoSource(result, c.typeOf(node), true)
}

func (c *GoCompiler) valueToGoSource(val value.Value, typ types.Type, allowMutable bool) *goValue {
	if val.IsReference() {
		switch v := val.AsReference().(type) {
		case value.ArrayList:
			if !allowMutable {
				return nil
			}
			return c.arrayListToGoSource(v, typ)
		case value.ArrayTuple:
			cached := c.emitCachedArrayTuple(v, typ)
			if cached != nil {
				return cached
			}
			if !allowMutable {
				return nil
			}
			return c.arrayTupleToGoSource(v, typ, true)
		case vm.HashSet:
			if !allowMutable {
				return nil
			}
			return c.hashSetToGoSource(v, typ)
		case vm.HashMap:
			if !allowMutable {
				return nil
			}
			return c.hashMapToGoSource(v, typ)
		case vm.HashRecord:
			cached := c.emitCachedHashRecord(v, typ)
			if cached != nil {
				return cached
			}
			if !allowMutable {
				return nil
			}
			return c.hashRecordToGoSource(v, typ, allowMutable)
		case value.String:
			return newGoValue(
				fmt.Sprintf("value.String(%q)", v.String()),
				c.checker.Std(symbol.String),
				value.FetchGoType("value.String"),
			)
		case value.Int64:
			return newGoValue(
				fmt.Sprintf("value.Int64(%d)", v),
				c.checker.Std(symbol.Int64),
				value.FetchGoType("value.Int64"),
			)
		case value.UInt64:
			return newGoValue(
				fmt.Sprintf("value.UInt64(%d)", v),
				c.checker.Std(symbol.UInt64),
				value.FetchGoType("value.UInt64"),
			)
		case *value.BigInt:
			return newGoValue(
				c.emitBigInt(string(v.ToString())),
				c.checker.Std(symbol.Int),
				value.FetchGoType("*value.BigInt"),
			)
		case *value.BeginlessClosedRange, *value.BeginlessOpenRange,
			*value.EndlessClosedRange, *value.EndlessOpenRange,
			*value.ClosedRange, *value.OpenRange,
			*value.LeftOpenRange, *value.RightOpenRange:
			return c.emitCachedRange(value.Ref(v), typ)
		case *value.Regex:
			return c.emitCachedRegex(v, typ)
		default:
			panic(fmt.Sprintf("cannot convert elk value to Go source: %T, %s", val, val.Inspect()))
		}
	}

	switch val.ValueFlag() {
	case value.BOOL_FLAG:
		if val.AsBool() {
			return newGoValue(
				"value.True",
				types.Bool{},
				value.FetchGoType("value.Bool"),
			)
		} else {
			return newGoValue(
				"value.False",
				types.Bool{},
				value.FetchGoType("value.Bool"),
			)
		}
	case value.NIL_FLAG:
		return nilGoValue
	case value.SMALL_INT_FLAG:
		return newGoValue(
			fmt.Sprintf("value.SmallInt(%d)", val.AsSmallInt()),
			c.checker.StdInt(),
			value.FetchGoType("value.SmallInt"),
		)
	case value.INT64_FLAG:
		return newGoValue(
			fmt.Sprintf("value.Int64(%d)", val.AsInt64()),
			c.checker.Std(symbol.Int64),
			value.FetchGoType("value.Int64"),
		)
	case value.UINT_FLAG:
		return newGoValue(
			fmt.Sprintf("value.UInt(%d)", val.AsUInt()),
			c.checker.Std(symbol.UInt),
			value.FetchGoType("value.UInt"),
		)
	case value.UINT64_FLAG:
		return newGoValue(
			fmt.Sprintf("value.UInt64(%d)", val.AsUInt64()),
			c.checker.Std(symbol.UInt64),
			value.FetchGoType("value.UInt64"),
		)
	case value.INT32_FLAG:
		return newGoValue(
			fmt.Sprintf("value.Int32(%d)", val.AsInt32()),
			c.checker.Std(symbol.Int32),
			value.FetchGoType("value.Int32"),
		)
	case value.UINT32_FLAG:
		return newGoValue(
			fmt.Sprintf("value.UInt32(%d)", val.AsUInt32()),
			c.checker.Std(symbol.UInt32),
			value.FetchGoType("value.UInt32"),
		)
	case value.INT16_FLAG:
		return newGoValue(
			fmt.Sprintf("value.Int16(%d)", val.AsInt16()),
			c.checker.Std(symbol.Int16),
			value.FetchGoType("value.Int16"),
		)
	case value.UINT16_FLAG:
		return newGoValue(
			fmt.Sprintf("value.UInt16(%d)", val.AsUInt16()),
			c.checker.Std(symbol.UInt16),
			value.FetchGoType("value.UInt16"),
		)
	case value.INT8_FLAG:
		return newGoValue(
			fmt.Sprintf("value.Int8(%d)", val.AsInt8()),
			c.checker.Std(symbol.Int8),
			value.FetchGoType("value.Int8"),
		)
	case value.UINT8_FLAG:
		return newGoValue(
			fmt.Sprintf("value.UInt8(%d)", val.AsUInt8()),
			c.checker.Std(symbol.UInt8),
			value.FetchGoType("value.UInt8"),
		)
	case value.CHAR_FLAG:
		return newGoValue(
			fmt.Sprintf("value.Char(%q)", rune(val.AsChar())),
			c.checker.Std(symbol.Char),
			value.FetchGoType("value.Char"),
		)
	case value.SYMBOL_FLAG:
		return newGoValue(
			c.emitSymbol(val.AsInlineSymbol().String()),
			c.checker.Std(symbol.Symbol),
			value.FetchGoType("value.Symbol"),
		)
	case value.FLOAT_FLAG:
		return newGoValue(
			fmt.Sprintf("value.Float(%g)", val.AsFloat()),
			c.checker.Std(symbol.Float),
			value.FetchGoType("value.Float"),
		)
	case value.FLOAT64_FLAG:
		return newGoValue(
			fmt.Sprintf("value.Float64(%g)", val.AsFloat64()),
			c.checker.Std(symbol.Float64),
			value.FetchGoType("value.Float64"),
		)
	case value.FLOAT32_FLAG:
		return newGoValue(
			fmt.Sprintf("value.Float32(%g)", val.AsFloat32()),
			c.checker.Std(symbol.Float32),
			value.FetchGoType("value.Float32"),
		)
	}

	panic(fmt.Sprintf("cannot convert elk value to Go source: %T, %s", val, val.Inspect()))
}

func (c *GoCompiler) arrayListToGoSource(v value.ArrayList, typ types.Type) *goValue {
	elementType, _ := c.checker.GetIteratorElementType(typ)
	if types.IsUntyped(elementType) {
		return c.arrayListOfValueToGoSource(v)
	}

	goElementType := c.elkTypeToGoType(elementType, false)
	if goElementType.Name == "value.Value" {
		return c.arrayListOfValueToGoSource(v)
	}

	var buff strings.Builder

	fmt.Fprintf(
		&buff,
		"value.NewNativeArrayListWithElements[%s](%d, ",
		goElementType.String(),
		v.LeftCapacity(),
	)

	var dependencies []*goValue
	for _, element := range v.Elements() {
		el := c.valueToGoSource(element, elementType, true)
		if el == nil {
			return nil
		}

		fmt.Fprintf(&buff, "%s, ", c.valueToNarrowerType(el).value)
		dependencies = append(dependencies, el)
	}

	buff.WriteString(")")
	return newGoValueWithDependencies(
		buff.String(),
		c.checker.Std(symbol.ArrayList),
		value.FetchGenericGoType(
			"*value.NativeArrayList",
			[]*value.GoType{
				goElementType,
			},
		),
		dependencies...,
	)
}

func (c *GoCompiler) arrayListOfValueToGoSource(v value.ArrayList) *goValue {
	var buff strings.Builder

	fmt.Fprintf(
		&buff,
		"value.NewArrayListOfValueWithElements(%d, ",
		v.LeftCapacity(),
	)

	var dependencies []*goValue
	for _, element := range v.Elements() {
		el := c.valueToGoSource(element, nil, true)
		if el == nil {
			return nil
		}

		fmt.Fprintf(&buff, "%s, ", c.convertToValue(el).value)
		dependencies = append(dependencies, el)
	}

	buff.WriteString(")")
	return newGoValueWithDependencies(
		buff.String(),
		c.checker.Std(symbol.ArrayList),
		value.FetchGoType("*value.ArrayListOfValue"),
		dependencies...,
	)
}

func (c *GoCompiler) convertToValue(v *goValue) *goValue {
	switch v.goType.Name {
	case "value.Value":
		return v
	case "bool":
		return v.newGoValue(
			fmt.Sprintf("value.ToBoolVal(%s)", v.value),
			v.elkType,
			goValueType,
		)
	}

	return v.newGoValue(
		fmt.Sprintf("(%s).ToValue()", v.value),
		v.elkType,
		goValueType,
	)
}

func (c *GoCompiler) convertToNativeInt(v *goValue) *goValue {
	if v == nil {
		newGoValue(
			"0",
			types.Any{},
			value.FetchGoType("int"),
		)
	}

	switch v.goType.Name {
	case "value.SmallInt", "value.Float",
		"value.Int64", "value.Int32", "value.Int16", "value.Int8",
		"value.UInt", "value.UInt64", "value.UInt32", "value.UInt16", "value.UInt8":
		v.newGoValue(
			fmt.Sprintf("int(%s)", v.value),
			v.elkType,
			value.FetchGoType("int"),
		)
	case "*value.BigInt":
		v.newGoValue(
			fmt.Sprintf("int((%s).ToSmallInt())", v.value),
			v.elkType,
			value.FetchGoType("int"),
		)
	}

	return v.newGoValue(
		fmt.Sprintf("(%s).AsAnyInt()", c.convertToValue(v).value),
		v.elkType,
		value.FetchGoType("int"),
	)
}

func (c *GoCompiler) valueToNarrowerType(v *goValue) *goValue {
	if v.goType.Name != "value.Value" {
		return v
	}

	if c.checker.IsSubtype(v.elkType, types.Bool{}) {
		return v.newGoValue(
			fmt.Sprintf("value.ToBool(%s)", v.value),
			v.elkType,
			value.FetchGoType("value.Bool"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Symbol)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsSymbol()", v.value),
			v.elkType,
			value.FetchGoType("value.Symbol"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.String)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsString()", v.value),
			v.elkType,
			value.FetchGoType("value.String"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Char)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsChar()", v.value),
			v.elkType,
			value.FetchGoType("value.Char"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Float)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsFloat()", v.value),
			v.elkType,
			value.FetchGoType("value.Float"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Float64)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsFloat64()", v.value),
			v.elkType,
			value.FetchGoType("value.Float64"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Float32)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsFloat32()", v.value),
			v.elkType,
			value.FetchGoType("value.Float32"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.BigFloat)) {
		return v.newGoValue(
			fmt.Sprintf("(*value.BigFloat)((%s).Pointer())", v.value),
			v.elkType,
			value.FetchGoType("*value.BigFloat"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Int64)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsInt64()", v.value),
			v.elkType,
			value.FetchGoType("value.Int64"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Int32)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsInt32()", v.value),
			v.elkType,
			value.FetchGoType("value.Int32"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Int16)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsInt16()", v.value),
			v.elkType,
			value.FetchGoType("value.Int16"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Int8)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsInt8()", v.value),
			v.elkType,
			value.FetchGoType("value.Int8"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.UInt)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsUInt()", v.value),
			v.elkType,
			value.FetchGoType("value.UInt"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.UInt64)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsUInt64()", v.value),
			v.elkType,
			value.FetchGoType("value.UInt64"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.UInt32)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsUInt32()", v.value),
			v.elkType,
			value.FetchGoType("value.UInt32"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.UInt16)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsUInt16()", v.value),
			v.elkType,
			value.FetchGoType("value.UInt16"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.UInt8)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsUInt8()", v.value),
			v.elkType,
			value.FetchGoType("value.UInt8"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.ArrayList)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsReference().(value.ArrayList)", v.value),
			v.elkType,
			value.FetchGoType("value.ArrayList"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.ArrayTuple)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsReference().(value.ArrayTuple)", v.value),
			v.elkType,
			value.FetchGoType("value.ArrayTuple"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.HashMap)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsReference().(vm.HashMap)", v.value),
			v.elkType,
			value.FetchGoType("vm.HashMap"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.HashRecord)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsReference().(vm.HashRecord)", v.value),
			v.elkType,
			value.FetchGoType("vm.HashRecord"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.HashSet)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsReference().(vm.HashSet)", v.value),
			v.elkType,
			value.FetchGoType("vm.HashSet"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.BeginlessClosedRange)) {
		return v.newGoValue(
			fmt.Sprintf("(*value.BeginlessClosedRange)((%s).Pointer())", v.value),
			v.elkType,
			value.FetchGoType("*value.BeginlessClosedRange"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.BeginlessOpenRange)) {
		return v.newGoValue(
			fmt.Sprintf("(*value.BeginlessOpenRange)((%s).Pointer())", v.value),
			v.elkType,
			value.FetchGoType("*value.BeginlessOpenRange"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.EndlessClosedRange)) {
		return v.newGoValue(
			fmt.Sprintf("(*value.EndlessClosedRange)((%s).Pointer())", v.value),
			v.elkType,
			value.FetchGoType("*value.EndlessClosedRange"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.EndlessOpenRange)) {
		return v.newGoValue(
			fmt.Sprintf("(*value.EndlessOpenRange)((%s).Pointer())", v.value),
			v.elkType,
			value.FetchGoType("*value.EndlessOpenRange"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.ClosedRange)) {
		return v.newGoValue(
			fmt.Sprintf("(*value.ClosedRange)((%s).Pointer())", v.value),
			v.elkType,
			value.FetchGoType("*value.ClosedRange"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.OpenRange)) {
		return v.newGoValue(
			fmt.Sprintf("(*value.OpenRange)((%s).Pointer())", v.value),
			v.elkType,
			value.FetchGoType("*value.OpenRange"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.LeftOpenRange)) {
		return v.newGoValue(
			fmt.Sprintf("(*value.LeftOpenRange)((%s).Pointer())", v.value),
			v.elkType,
			value.FetchGoType("*value.LeftOpenRange"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.RightOpenRange)) {
		return v.newGoValue(
			fmt.Sprintf("(*value.RightOpenRange)((%s).Pointer())", v.value),
			v.elkType,
			value.FetchGoType("*value.RightOpenRange"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Class)) {
		return v.newGoValue(
			fmt.Sprintf("(*value.Class)((%s).Pointer())", v.value),
			v.elkType,
			value.FetchGoType("*value.Class"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Module)) {
		return v.newGoValue(
			fmt.Sprintf("(*value.Module)((%s).Pointer())", v.value),
			v.elkType,
			value.FetchGoType("*value.Module"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Mixin)) {
		return v.newGoValue(
			fmt.Sprintf("(*value.Mixin)((%s).Pointer())", v.value),
			v.elkType,
			value.FetchGoType("*value.Mixin"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Interface)) {
		return v.newGoValue(
			fmt.Sprintf("(*value.Interface)((%s).Pointer())", v.value),
			v.elkType,
			value.FetchGoType("*value.Interface"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Time)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsTime()", v.value),
			v.elkType,
			value.FetchGoType("value.Time"),
		)
	}
	if c.checker.IsSubtype(v.elkType, c.checker.Std(symbol.Date)) {
		return v.newGoValue(
			fmt.Sprintf("(%s).AsDate()", v.value),
			v.elkType,
			value.FetchGoType("value.Date"),
		)
	}

	return v
}

func (c *GoCompiler) elkTypeToGoType(elkType types.Type, specialized bool) *value.GoType {
	if c.checker.IsSubtype(elkType, types.Bool{}) {
		return value.FetchGoType("value.Bool")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Symbol)) {
		return value.FetchGoType("value.Symbol")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.String)) {
		return value.FetchGoType("value.String")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Char)) {
		return value.FetchGoType("value.Char")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Regex)) {
		return value.FetchGoType("*value.Regex")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Float)) {
		return value.FetchGoType("value.Float")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Float64)) {
		return value.FetchGoType("value.Float64")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Float32)) {
		return value.FetchGoType("value.Float32")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.BigFloat)) {
		return value.FetchGoType("*value.BigFloat")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Int64)) {
		return value.FetchGoType("value.Int64")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Int32)) {
		return value.FetchGoType("value.Int32")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Int16)) {
		return value.FetchGoType("value.Int16")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Int8)) {
		return value.FetchGoType("value.Int8")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.UInt)) {
		return value.FetchGoType("value.UInt")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.UInt64)) {
		return value.FetchGoType("value.UInt64")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.UInt32)) {
		return value.FetchGoType("value.UInt32")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.UInt16)) {
		return value.FetchGoType("value.UInt16")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.UInt8)) {
		return value.FetchGoType("value.UInt8")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.ArrayList)) {
		if !specialized {
			return value.FetchGoType("value.ArrayList")
		}

		elkType, ok := elkType.(*types.Generic)
		if !ok {
			return value.FetchGoType("*value.ArrayListOfValue")
		}

		elementType := elkType.Get(0).Type

		goElementType := c.elkTypeToGoType(elementType, true)
		if goElementType.Equal(goValueType) {
			return value.FetchGoType("*value.ArrayListOfValue")
		}

		return value.FetchGenericGoType(
			"*value.NativeArrayList",
			[]*value.GoType{
				goElementType,
			},
		)
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.ArrayTuple)) {
		if !specialized {
			return value.FetchGoType("value.ArrayTuple")
		}

		elkType, ok := elkType.(*types.Generic)
		if !ok {
			return value.FetchGoType("*value.ArrayTupleOfValue")
		}

		elementType := elkType.Get(0).Type

		goElementType := c.elkTypeToGoType(elementType, true)
		if goElementType.Equal(goValueType) {
			return value.FetchGoType("*value.ArrayTupleOfValue")
		}

		return value.FetchGenericGoType(
			"*value.NativeArrayTuple",
			[]*value.GoType{
				goElementType,
			},
		)
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.HashMap)) {
		if !specialized {
			return value.FetchGoType("vm.HashMap")
		}

		elkType, ok := elkType.(*types.Generic)
		if !ok {
			return value.FetchGoType("*vm.HashMapOfValue")
		}

		keyType := elkType.Get(0).Type
		valType := elkType.Get(1).Type

		nativeKeyType := c.elkTypeToGoKeyType(keyType)
		if nativeKeyType.Equal(goValueType) {
			return value.FetchGoType("*vm.HashMapOfValue")
		}
		nativeValType := c.elkTypeToGoType(valType, true)
		if nativeValType.Equal(goValueType) {
			return value.FetchGenericGoType(
				"*vm.NativeKeyHashMap",
				[]*value.GoType{
					nativeKeyType,
				},
			)
		}

		return value.FetchGenericGoType(
			"*vm.NativeHashMap",
			[]*value.GoType{
				nativeKeyType,
				nativeValType,
			},
		)
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.HashRecord)) {
		if !specialized {
			return value.FetchGoType("vm.HashRecord")
		}

		elkType, ok := elkType.(*types.Generic)
		if !ok {
			return value.FetchGoType("vm.HashRecordOfValue")
		}

		keyType := elkType.Get(0).Type
		valType := elkType.Get(1).Type

		nativeKeyType := c.elkTypeToGoKeyType(keyType)
		if nativeKeyType.Equal(goValueType) {
			return value.FetchGoType("vm.HashRecordOfValue")
		}
		nativeValType := c.elkTypeToGoType(valType, true)
		if nativeValType.Equal(goValueType) {
			return value.FetchGenericGoType(
				"vm.NativeKeyHashRecord",
				[]*value.GoType{
					nativeKeyType,
				},
			)
		}

		return value.FetchGenericGoType(
			"vm.NativeHashRecord",
			[]*value.GoType{
				nativeKeyType,
				nativeValType,
			},
		)
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.HashSet)) {
		if !specialized {
			return value.FetchGoType("vm.HashSet")
		}

		elkType, ok := elkType.(*types.Generic)
		if !ok {
			return value.FetchGoType("*vm.HashSetOfValue")
		}

		elementType := elkType.Get(0).Type

		goElementType := c.elkTypeToGoKeyType(elementType)
		if goElementType.Equal(goValueType) {
			return value.FetchGoType("*vm.HashSetOfValue")
		}

		return value.FetchGenericGoType(
			"*vm.NativeHashSet",
			[]*value.GoType{
				goElementType,
			},
		)
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.BeginlessClosedRange)) {
		return value.FetchGoType("*value.BeginlessClosedRange")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.BeginlessOpenRange)) {
		return value.FetchGoType("*value.BeginlessOpenRange")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.EndlessClosedRange)) {
		return value.FetchGoType("*value.EndlessClosedRange")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.EndlessOpenRange)) {
		return value.FetchGoType("*value.EndlessOpenRange")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.ClosedRange)) {
		return value.FetchGoType("*value.ClosedRange")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.OpenRange)) {
		return value.FetchGoType("*value.OpenRange")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.LeftOpenRange)) {
		return value.FetchGoType("*value.LeftOpenRange")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.RightOpenRange)) {
		return value.FetchGoType("*value.RightOpenRange")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Date)) {
		return value.FetchGoType("value.Date")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Time)) {
		return value.FetchGoType("value.Time")
	}

	return goValueType
}

// Convert an elk type to a native go type that can be used as a key
// in a hash map or an element in a hash set
func (c *GoCompiler) elkTypeToGoKeyType(elkType types.Type) *value.GoType {
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Symbol)) {
		return value.FetchGoType("value.Symbol")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.String)) {
		return value.FetchGoType("value.String")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Char)) {
		return value.FetchGoType("value.Char")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Float)) {
		return value.FetchGoType("value.Float")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Float64)) {
		return value.FetchGoType("value.Float64")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Float32)) {
		return value.FetchGoType("value.Float32")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Int64)) {
		return value.FetchGoType("value.Int64")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Int32)) {
		return value.FetchGoType("value.Int32")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Int16)) {
		return value.FetchGoType("value.Int16")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Int8)) {
		return value.FetchGoType("value.Int8")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.UInt)) {
		return value.FetchGoType("value.UInt")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.UInt64)) {
		return value.FetchGoType("value.UInt64")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.UInt32)) {
		return value.FetchGoType("value.UInt32")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.UInt16)) {
		return value.FetchGoType("value.UInt16")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.UInt8)) {
		return value.FetchGoType("value.UInt8")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Date)) {
		return value.FetchGoType("value.Date")
	}
	if c.checker.IsSubtype(elkType, c.checker.Std(symbol.Time)) {
		return value.FetchGoType("value.Time")
	}

	return goValueType
}

func (c *GoCompiler) rangeToGoSource(v value.Value, typ types.Type, mutable bool) *goValue {
	switch val := v.AsReference().(type) {
	case *value.BeginlessClosedRange:
		end := c.convertToValue(c.valueToGoSource(val.End, typ, mutable))
		return newGoValueWithDependencies(
			fmt.Sprintf(
				"value.NewBeginlessClosedRange(%s)",
				end.value,
			),
			c.checker.Std(symbol.BeginlessClosedRange),
			value.FetchGoType("*value.BeginlessClosedRange"),
			end,
		)
	case *value.BeginlessOpenRange:
		end := c.convertToValue(c.valueToGoSource(val.End, typ, mutable))
		return newGoValueWithDependencies(
			fmt.Sprintf(
				"value.NewBeginlessOpenRange(%s)",
				end.value,
			),
			c.checker.Std(symbol.BeginlessOpenRange),
			value.FetchGoType("*value.BeginlessOpenRange"),
			end,
		)
	case *value.EndlessClosedRange:
		start := c.convertToValue(c.valueToGoSource(val.Start, typ, mutable))
		return newGoValueWithDependencies(
			fmt.Sprintf(
				"value.NewEndlessClosedRange(%s)",
				start.value,
			),
			c.checker.Std(symbol.EndlessClosedRange),
			value.FetchGoType("*value.EndlessClosedRange"),
			start,
		)
	case *value.EndlessOpenRange:
		start := c.convertToValue(c.valueToGoSource(val.Start, typ, mutable))
		return newGoValueWithDependencies(
			fmt.Sprintf(
				"value.NewEndlessOpenRange(%s)",
				start.value,
			),
			c.checker.Std(symbol.EndlessOpenRange),
			value.FetchGoType("*value.EndlessOpenRange"),
			start,
		)
	case *value.ClosedRange:
		start := c.convertToValue(c.valueToGoSource(val.Start, typ, mutable))
		end := c.convertToValue(c.valueToGoSource(val.End, typ, mutable))

		return newGoValueWithDependencies(
			fmt.Sprintf(
				"value.NewClosedRange(%s, %s)",
				start.value,
				end.value,
			),
			c.checker.Std(symbol.ClosedRange),
			value.FetchGoType("*value.ClosedRange"),
			start,
			end,
		)
	case *value.OpenRange:
		start := c.convertToValue(c.valueToGoSource(val.Start, typ, mutable))
		end := c.convertToValue(c.valueToGoSource(val.End, typ, mutable))

		return newGoValueWithDependencies(
			fmt.Sprintf(
				"value.NewOpenRange(%s, %s)",
				start.value,
				end.value,
			),
			c.checker.Std(symbol.OpenRange),
			value.FetchGoType("*value.OpenRange"),
			start,
			end,
		)
	case *value.LeftOpenRange:
		start := c.convertToValue(c.valueToGoSource(val.Start, typ, mutable))
		end := c.convertToValue(c.valueToGoSource(val.End, typ, mutable))

		return newGoValueWithDependencies(
			fmt.Sprintf(
				"value.NewLeftOpenRange(%s, %s)",
				start.value,
				end.value,
			),
			c.checker.Std(symbol.LeftOpenRange),
			value.FetchGoType("*value.LeftOpenRange"),
			start,
			end,
		)
	case *value.RightOpenRange:
		start := c.convertToValue(c.valueToGoSource(val.Start, typ, mutable))
		end := c.convertToValue(c.valueToGoSource(val.End, typ, mutable))

		return newGoValueWithDependencies(
			fmt.Sprintf(
				"value.NewRightOpenRange(%s, %s)",
				start.value,
				end.value,
			),
			c.checker.Std(symbol.RightOpenRange),
			value.FetchGoType("*value.RightOpenRange"),
			start,
			end,
		)
	default:
		panic(fmt.Sprintf("invalid range value: %#v", val))
	}
}

func (c *GoCompiler) arrayTupleOfValueToGoSource(v value.ArrayTuple, mutable bool) *goValue {
	var buff strings.Builder

	buff.WriteString("value.NewArrayTupleOfValueWithElements(0, ")

	var dependencies []*goValue
	for _, element := range v.Elements() {
		el := c.valueToGoSource(element, nil, mutable)
		if el == nil {
			return nil
		}

		fmt.Fprintf(&buff, "%s, ", c.convertToValue(el).value)
		dependencies = append(dependencies, el)
	}

	buff.WriteString(")")
	return newGoValueWithDependencies(
		buff.String(),
		c.checker.Std(symbol.ArrayTuple),
		value.FetchGoType("*value.ArrayTupleOfValue"),
		dependencies...,
	)
}

func (c *GoCompiler) compileRegexFlags(buff io.Writer, flags bitfield.BitField8) {
	c.registerGoImport("github.com/elk-language/elk/bitfield", "")

	if !flags.IsAnyFlagSet() {
		fmt.Fprint(buff, "bitfield.BitField8FromBitFlag(0)")
		return
	}

	c.registerGoImport("github.com/elk-language/elk/regex/flag", "reflag")

	fmt.Fprint(buff, "bitfield.BitField8FromBitFlag(")
	var i int
	for _, flag := range flags.AllSetFlags() {
		if i != 0 {
			fmt.Fprint(buff, " | ")
		}

		switch flag {
		case reflag.CaseInsensitiveFlag:
			fmt.Fprint(buff, "reflag.CaseInsensitiveFlag")
		case reflag.MultilineFlag:
			fmt.Fprint(buff, "reflag.MultilineFlag")
		case reflag.DotAllFlag:
			fmt.Fprint(buff, "reflag.DotAllFlag")
		case reflag.UngreedyFlag:
			fmt.Fprint(buff, "reflag.UngreedyFlag")
		case reflag.ExtendedFlag:
			fmt.Fprint(buff, "reflag.ExtendedFlag")
		case reflag.ASCIIFlag:
			fmt.Fprint(buff, "reflag.ASCIIFlag")
		default:
			panic(fmt.Sprintf("invalid regex flag: %d", flag))
		}

		i++
	}

	fmt.Fprint(buff, ")")
}

func (c *GoCompiler) regexToGoSource(v *value.Regex, typ types.Type) *goValue {
	var buff bytes.Buffer

	fmt.Fprintf(&buff, "value.MustCompileRegex(%q,", v.Source)
	c.compileRegexFlags(&buff, v.Flags)
	buff.WriteString(")")

	return newGoValue(
		buff.String(),
		c.checker.Std(symbol.Regex),
		value.FetchGoType("*value.Regex"),
	)
}

func (c *GoCompiler) arrayTupleToGoSource(v value.ArrayTuple, typ types.Type, mutable bool) *goValue {
	elementType, _ := c.checker.GetIteratorElementType(typ)
	if types.IsUntyped(elementType) {
		return c.arrayTupleOfValueToGoSource(v, mutable)
	}

	goElementType := c.elkTypeToGoType(elementType, true)
	if goElementType.Name == "value.Value" {
		return c.arrayTupleOfValueToGoSource(v, mutable)
	}

	var buff strings.Builder

	fmt.Fprintf(&buff, "value.NewNativeArrayTupleWithElements[%s](0, ", goElementType.String())

	var dependencies []*goValue
	for _, element := range v.Elements() {
		el := c.valueToGoSource(element, elementType, mutable)
		if el == nil {
			return nil
		}

		fmt.Fprintf(&buff, "%s, ", c.valueToNarrowerType(el).value)
		dependencies = append(dependencies, el)
	}

	buff.WriteString(")")
	return newGoValueWithDependencies(
		buff.String(),
		c.checker.Std(symbol.ArrayTuple),
		value.FetchGenericGoType(
			"*value.NativeArrayTuple",
			[]*value.GoType{
				goElementType,
			},
		),
		dependencies...,
	)
}

func (c *GoCompiler) hashSetToGoSource(v vm.HashSet, typ types.Type) *goValue {
	elementType, _ := c.checker.GetIteratorElementType(typ)
	if types.IsUntyped(elementType) {
		return c.hashSetOfValueToGoSource(v)
	}

	goElementType := c.elkTypeToGoKeyType(elementType)
	if goElementType.Name == "value.Value" {
		return c.hashSetOfValueToGoSource(v)
	}

	var buff strings.Builder

	fmt.Fprintf(&buff, "vm.NewNativeHashSetWithElements[%s](", goElementType.String())

	var dependencies []*goValue
	for _, element := range inspectSort(slices.Collect(v.All())) {
		el := c.valueToGoSource(element, elementType, true)
		if el == nil {
			return nil
		}

		fmt.Fprintf(&buff, "%s, ", c.valueToNarrowerType(el).value)
		dependencies = append(dependencies, el)
	}

	buff.WriteString(")")
	return newGoValueWithDependencies(
		buff.String(),
		c.checker.Std(symbol.HashSet),
		value.FetchGenericGoType(
			"*vm.NativeHashSet",
			[]*value.GoType{
				goElementType,
			},
		),
		dependencies...,
	)
}

func (c *GoCompiler) hashSetOfValueToGoSource(v vm.HashSet) *goValue {
	var buff strings.Builder

	fmt.Fprintf(&buff, "vm.MustNewHashSetOfValueWithCapacityAndElements(nil, 0, ")

	var dependencies []*goValue
	for _, element := range inspectSort(slices.Collect(v.All())) {
		el := c.valueToGoSource(element, nil, true)
		if el == nil {
			return nil
		}

		fmt.Fprintf(&buff, "%s, ", c.convertToValue(el).value)
		dependencies = append(dependencies, el)
	}

	buff.WriteString(")")
	return newGoValueWithDependencies(
		buff.String(),
		c.checker.Std(symbol.HashSet),
		value.FetchGoType("*vm.HashSetOfValue"),
		dependencies...,
	)
}

func (c *GoCompiler) hashMapOfValueToGoSource(v vm.HashMap) *goValue {
	var buff strings.Builder

	fmt.Fprintf(&buff, "vm.MustNewHashMapOfValueWithCapacityAndElements(nil, 0, ")

	var dependencies []*goValue
	for _, pair := range inspectSort(slices.Collect(v.All())) {
		p := c.valuePairToGoSource(pair, true)
		if p == nil {
			return nil
		}

		fmt.Fprintf(&buff, "%s, ", p.value)
		dependencies = append(dependencies, p)
	}

	buff.WriteString(")")
	return newGoValueWithDependencies(
		buff.String(),
		c.checker.Std(symbol.HashMap),
		value.FetchGoType("*vm.HashMapOfValue"),
		dependencies...,
	)
}

func (c *GoCompiler) hashMapToGoSource(v vm.HashMap, typ types.Type) *goValue {
	elementType, _ := c.checker.GetIteratorElementType(typ)
	if types.IsUntyped(elementType) {
		return c.hashMapOfValueToGoSource(v)
	}
	if !c.checker.IsSubtype(elementType, c.checker.Std(symbol.Pair)) {
		return c.hashMapOfValueToGoSource(v)
	}

	pairType, ok := elementType.(*types.Generic)
	if !ok {
		return c.hashMapOfValueToGoSource(v)
	}

	keyType := pairType.Get(0).Type
	valType := pairType.Get(1).Type

	goKeyType := c.elkTypeToGoKeyType(keyType)
	if goKeyType.Name == "value.Value" {
		return c.hashMapOfValueToGoSource(v)
	}
	goValType := c.elkTypeToGoKeyType(valType)

	if goValType.Name == "value.Value" {
		var buff strings.Builder

		fmt.Fprintf(&buff, "vm.NewNativeKeyHashMapFromMap(map[%s]value.Value{", goKeyType.String())

		var dependencies []*goValue
		for _, pair := range inspectSort(slices.Collect(v.All())) {
			keySource := c.valueToGoSource(pair.Key(), keyType, true)
			if keySource == nil {
				return nil
			}
			valSource := c.valueToGoSource(pair.Value(), valType, true)
			if valSource == nil {
				return nil
			}

			fmt.Fprintf(
				&buff,
				"%s: %s, ",
				c.valueToNarrowerType(keySource).value,
				c.convertToValue(valSource).value,
			)
			dependencies = append(dependencies, keySource, valSource)
		}

		buff.WriteString("})")
		return newGoValueWithDependencies(
			buff.String(),
			c.checker.Std(symbol.HashMap),
			value.FetchGenericGoType(
				"*vm.NativeHashMap",
				[]*value.GoType{
					goKeyType,
					goValType,
				},
			),
			dependencies...,
		)
	}

	var buff strings.Builder

	fmt.Fprintf(&buff, "vm.NewNativeHashMapFromMap(map[%s]%s{", goKeyType.String(), goValType.String())

	var dependencies []*goValue
	for _, pair := range inspectSort(slices.Collect(v.All())) {
		keySource := c.valueToGoSource(pair.Key(), keyType, true)
		if keySource == nil {
			return nil
		}
		valSource := c.valueToGoSource(pair.Value(), valType, true)
		if valSource == nil {
			return nil
		}

		fmt.Fprintf(
			&buff,
			"%s: %s, ",
			c.valueToNarrowerType(keySource).value,
			c.valueToNarrowerType(valSource).value,
		)
		dependencies = append(dependencies, keySource, valSource)
	}

	buff.WriteString("})")
	return newGoValueWithDependencies(
		buff.String(),
		c.checker.Std(symbol.HashMap),
		value.FetchGenericGoType(
			"*vm.NativeHashMap",
			[]*value.GoType{
				goKeyType,
				goValType,
			},
		),
		dependencies...,
	)
}

func (c *GoCompiler) hashRecordToGoSource(v vm.HashRecord, typ types.Type, allowMutable bool) *goValue {
	elementType, _ := c.checker.GetIteratorElementType(typ)
	if types.IsUntyped(elementType) {
		return c.hashRecordOfValueToGoSource(v, allowMutable)
	}
	if !c.checker.IsSubtype(elementType, c.checker.Std(symbol.Pair)) {
		return c.hashRecordOfValueToGoSource(v, allowMutable)
	}

	pairType, ok := elementType.(*types.Generic)
	if !ok {
		return c.hashRecordOfValueToGoSource(v, allowMutable)
	}

	keyType := pairType.Get(0).Type
	valType := pairType.Get(1).Type

	goKeyType := c.elkTypeToGoKeyType(keyType)
	if goKeyType.Name == "value.Value" {
		return c.hashRecordOfValueToGoSource(v, allowMutable)
	}
	goValType := c.elkTypeToGoType(valType, true)

	if goValType.Name == "value.Value" {
		var buff strings.Builder

		fmt.Fprintf(&buff, "vm.MakeNativeKeyHashRecordFromMap(map[%s]value.Value{", goKeyType.String())

		var dependencies []*goValue
		for _, pair := range inspectSort(slices.Collect(v.All())) {
			keySource := c.valueToGoSource(pair.Key(), keyType, true)
			if keySource == nil {
				return nil
			}
			valSource := c.valueToGoSource(pair.Value(), valType, true)
			if valSource == nil {
				return nil
			}

			fmt.Fprintf(
				&buff,
				"%s: %s, ",
				c.valueToNarrowerType(keySource).value,
				c.convertToValue(valSource).value,
			)
			dependencies = append(dependencies, keySource, valSource)
		}

		buff.WriteString("})")
		return newGoValueWithDependencies(
			buff.String(),
			c.checker.Std(symbol.HashMap),
			value.FetchGenericGoType(
				"vm.NativeHashRecord",
				[]*value.GoType{
					goKeyType,
					goValType,
				},
			),
			dependencies...,
		)
	}

	var buff strings.Builder

	fmt.Fprintf(&buff, "vm.MakeNativeHashRecordFromMap(map[%s]%s{", goKeyType.String(), goValType.String())

	var dependencies []*goValue
	for _, pair := range inspectSort(slices.Collect(v.All())) {
		keySource := c.valueToGoSource(pair.Key(), keyType, true)
		if keySource == nil {
			return nil
		}
		valSource := c.valueToGoSource(pair.Value(), valType, true)
		if valSource == nil {
			return nil
		}

		fmt.Fprintf(
			&buff,
			"%s: %s, ",
			c.valueToNarrowerType(keySource).value,
			c.valueToNarrowerType(valSource).value,
		)
		dependencies = append(dependencies, keySource, valSource)
	}

	buff.WriteString("})")
	return newGoValueWithDependencies(
		buff.String(),
		c.checker.Std(symbol.HashMap),
		value.FetchGenericGoType(
			"vm.NativeHashRecord",
			[]*value.GoType{
				goKeyType,
				goValType,
			},
		),
		dependencies...,
	)
}

func (c *GoCompiler) hashRecordOfValueToGoSource(v vm.HashRecord, allowMutable bool) *goValue {
	var buff strings.Builder

	buff.WriteString("vm.MustNewHashRecordOfValueWithElements(nil, ")

	var dependencies []*goValue
	for _, pair := range inspectSort(slices.Collect(v.All())) {
		p := c.valuePairToGoSource(pair, allowMutable)
		if p == nil {
			return nil
		}

		fmt.Fprintf(&buff, "%s, ", p.value)
		dependencies = append(dependencies, p)
	}

	buff.WriteString(")")
	return newGoValueWithDependencies(
		buff.String(),
		c.checker.Std(symbol.HashRecord),
		value.FetchGoType("*vm.HashRecordOfValue"),
		dependencies...,
	)
}

func (c *GoCompiler) valuePairToGoSource(p value.PairOfValue, allowMutable bool) *goValue {
	k := c.valueToGoSource(p.Key(), nil, allowMutable)
	if k == nil {
		return nil
	}
	v := c.valueToGoSource(p.Value(), nil, allowMutable)
	if v == nil {
		return nil
	}

	return newGoValueWithDependencies(
		fmt.Sprintf(
			"value.MakePairOfValue(%s, %s)",
			c.convertToValue(k).value,
			c.convertToValue(v).value,
		),
		c.checker.Std(symbol.Pair),
		value.FetchGoType("value.PairOfValue"),
		k, v,
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
func (c *GoCompiler) defineLocal(name string, elkType types.Type, goType *value.GoType, location *position.Location) *nativeElkLocal {
	varScope := c.scopes.last()
	_, ok := varScope.localTable[name]
	if ok {
		c.Errors.AddFailure(
			fmt.Sprintf("a variable with this name has already been declared in this scope `%s`", name),
			location,
		)
		return nil
	}
	return c.defineVariableInScope(varScope, name, elkType, goType, location)
}

func (c *GoCompiler) defineVariableInScope(scope *nativeElkScope, name string, elkType types.Type, goType *value.GoType, location *position.Location) *nativeElkLocal {
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
		goType,
		fmt.Sprintf("var %s: %s", name, types.Inspect(elkType)),
	)
	goLocal.elkLocal = true

	newVar := &nativeElkLocal{
		name:    name,
		elkType: elkType,
		goLocal: goLocal,
	}
	scope.localTable[name] = newVar

	return newVar
}
