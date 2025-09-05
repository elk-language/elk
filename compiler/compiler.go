// Package Compiler implements
// the Elk Bytecode Compiler.
// It takes in Elk source code and outputs
// Elk Bytecode that can be run the Elk VM.
package compiler

import (
	"encoding/binary"
	"fmt"
	"math"
	"slices"
	"strconv"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"

	"github.com/elk-language/elk/token"
)

const MainName = "<main>"

func CreateCompiler(parent *Compiler, checker types.Checker, loc *position.Location, errors *diagnostic.SyncDiagnosticList) *Compiler {
	compiler := New(loc.FilePath, topLevelMode, loc, checker)
	compiler.Errors = errors
	compiler.Parent = parent
	return compiler
}

func (c *Compiler) CreateMainCompiler(checker types.Checker, loc *position.Location, errors *diagnostic.SyncDiagnosticList) *Compiler {
	compiler := New(loc.FilePath, topLevelMode, loc, checker)
	compiler.predefinedLocals = c.maxLocalIndex + 1
	compiler.scopes = c.scopes
	compiler.lastLocalIndex = c.lastLocalIndex
	compiler.maxLocalIndex = c.maxLocalIndex
	compiler.Errors = errors
	return compiler
}

func (c *Compiler) InitGlobalEnv() *Compiler {
	envCompiler := New("<namespaceDefinitions>", topLevelMode, c.Bytecode.Location, c.checker)
	envCompiler.Parent = c
	envCompiler.Errors = c.Errors
	envCompiler.compileGlobalEnv()
	return envCompiler
}

func (c *Compiler) EmitExecInParent() {
	parent := c.Parent
	location := parent.Bytecode.Location
	parent.emitValue(value.Ref(c.Bytecode), location)
	parent.emit(location.StartPos.Line, bytecode.EXEC)
	parent.emit(location.StartPos.Line, bytecode.POP)
}

// Compiler mode
type mode uint8

const (
	topLevelMode mode = iota
	namespaceMode
	methodMode
	macroMode
	setterMethodMode
	initMethodMode
	valuePatternDeclarationNode
)

// represents a local variable or value
type local struct {
	index      uint16
	hasUpvalue bool // is captured by some upvalue in a closure
}

type localTable map[string]*local

type scopeType uint8

const (
	defaultScopeType       scopeType = iota
	loopScopeType                    // this scope is a loop
	doFinallyScopeType               // this scope is inside do with a finally block
	macroBoundaryScopeType           // this scope is a macro boundary, locals from the outer scopes should be ignored
)

// set of local variables
type scope struct {
	localTable map[string]*local
	label      string
	typ        scopeType
}

func newScope(label string, typ scopeType) *scope {
	return &scope{
		localTable: localTable{},
		label:      label,
		typ:        typ,
	}
}

// indices represent scope depths
// and elements are sets of local variable names in a particular scope
type scopes []*scope

// Get the last local variable scope.
func (s scopes) last() *scope {
	return s[len(s)-1]
}

type loopJumpInfoType uint8

const (
	breakLoopJump           loopJumpInfoType = iota // break
	breakFinallyLoopJump                            // break inside of finally
	continueLoopJump                                // continue
	continueFinallyLoopJump                         // continue inside of finally
)

type loopJumpInfo struct {
	typ      loopJumpInfoType
	offset   int
	location *position.Location
}

type loopJumpSet struct {
	label                         string
	returnsValueFromLastIteration bool
	loopJumps                     []*loopJumpInfo
}

// Represents an upvalue, a captured variable from an outer context
type upvalue struct {
	index uint16 // index of the upvalue
	// index of the captured local if `isLocal` is true,
	// otherwise the index of the captured upvalue from the outer context
	upIndex uint16
	isLocal bool   // whether the captured variable is a local or an upvalue
	local   *local // the local that is captured through this upvalue
}

// Holds the state of the Compiler.
type Compiler struct {
	Name               string
	Bytecode           *vm.BytecodeFunction
	Errors             *diagnostic.SyncDiagnosticList
	scopes             scopes
	loopJumpSets       []*loopJumpSet
	offsetValueIds     []int // ids of integers in the value pool that represent bytecode offsets
	lastLocalIndex     int   // index of the last local variable
	maxLocalIndex      int   // max index of a local variable
	predefinedLocals   int
	mode               mode
	isGenerator        bool
	isAsync            bool
	unhygienic         bool
	secondToLastOpCode bytecode.OpCode
	lastOpCode         bytecode.OpCode
	patternNesting     int
	Parent             *Compiler
	upvalues           []*upvalue
	checker            types.Checker
}

// Instantiate a New Compiler instance.
func New(name string, mode mode, loc *position.Location, checker types.Checker) *Compiler {
	c := &Compiler{
		Bytecode: vm.NewBytecodeFunctionSimple(
			value.ToSymbol(name),
			[]byte{},
			loc,
		),
		scopes:         scopes{newScope("", defaultScopeType)}, // start with an empty set for the 0th scope
		lastLocalIndex: -1,
		maxLocalIndex:  -1,
		Name:           name,
		mode:           mode,
		checker:        checker,
		Errors:         diagnostic.NewSyncDiagnosticList(),
	}
	// reserve the first slot on the stack for `self`
	c.defineLocal("$self", position.DefaultLocation)
	switch mode {
	case topLevelMode, namespaceMode,
		methodMode, setterMethodMode, initMethodMode, macroMode:
		c.predefinedLocals = 1
	}
	return c
}

func (c *Compiler) EmitReturnNil() {
	location := c.Bytecode.Location
	c.emit(location.EndPos.Line, bytecode.NIL)
	c.emit(location.EndPos.Line, bytecode.RETURN)
}

func (c *Compiler) EmitReturn() {
	location := c.Bytecode.Location
	if c.lastOpCode != bytecode.RETURN {
		c.emit(location.EndPos.Line, bytecode.RETURN)
	}
	c.prepLocals()
}

func (c *Compiler) typeOf(node ast.Node) types.Type {
	return node.Type(c.checker.Env())
}

func (c *Compiler) compileGlobalEnv() {
	location := c.Bytecode.Location
	env := c.checker.Env()
	c.compileModuleDefinition(env.Root, env.Root, value.ToSymbol("Root"), location)
}

func (c *Compiler) compileNamespaceDefinition(parentNamespace, namespace types.Namespace, namespaceType byte, constName value.Symbol, location *position.Location) {
	if !namespace.IsDefined() && !namespace.IsNative() {
		switch p := parentNamespace.(type) {
		case *types.SingletonClass:
			c.emitGetConst(value.ToSymbol(p.AttachedObject.Name()), location)
			c.emit(location.StartPos.Line, bytecode.GET_SINGLETON)
		default:
			c.emitGetConst(value.ToSymbol(p.Name()), location)
		}
		c.emitValue(constName.ToValue(), location)
		c.emit(location.StartPos.Line, bytecode.DEF_NAMESPACE, namespaceType)
		namespace.SetDefined(true)
	}

	for name, subtype := range types.SortedSubtypes(namespace) {
		if subtype.Type == namespace {
			continue
		}
		c.compileSubtypeDefinition(namespace, subtype.Type, name, location)
	}
}

func (c *Compiler) compileModuleDefinition(parentNamespace types.Namespace, module *types.Module, constName value.Symbol, location *position.Location) {
	c.compileNamespaceDefinition(parentNamespace, module, bytecode.DEF_MODULE_FLAG, constName, location)
}

func (c *Compiler) compileClassDefinition(parentNamespace types.Namespace, class *types.Class, constName value.Symbol, location *position.Location) {
	c.compileNamespaceDefinition(parentNamespace, class, bytecode.DEF_CLASS_FLAG, constName, location)
}

func (c *Compiler) compileMixinDefinition(parentNamespace types.Namespace, mixin *types.Mixin, constName value.Symbol, location *position.Location) {
	c.compileNamespaceDefinition(parentNamespace, mixin, bytecode.DEF_MIXIN_FLAG, constName, location)
}

func (c *Compiler) compileInterfaceDefinition(parentNamespace types.Namespace, iface *types.Interface, constName value.Symbol, location *position.Location) {
	c.compileNamespaceDefinition(parentNamespace, iface, bytecode.DEF_INTERFACE_FLAG, constName, location)
}

func (c *Compiler) compileSubtypeDefinition(parentNamespace types.Namespace, typ types.Type, constName value.Symbol, location *position.Location) {
	switch t := typ.(type) {
	case *types.Module:
		c.compileModuleDefinition(parentNamespace, t, constName, location)
	case *types.Class:
		c.compileClassDefinition(parentNamespace, t, constName, location)
	case *types.Mixin:
		c.compileMixinDefinition(parentNamespace, t, constName, location)
	case *types.Interface:
		c.compileInterfaceDefinition(parentNamespace, t, constName, location)
	}
}

func (c *Compiler) CompileClassInheritance(class *types.Class, location *position.Location) {
	if class.IsCompiled() {
		return
	}
	superclass := class.Superclass()
	if superclass == nil {
		return
	}

	class.SetCompiled(true)
	name := value.ToSymbol(class.Name())
	// get the class
	c.emitGetConst(name, location)

	superclassName := value.ToSymbol(superclass.Name())
	c.emitGetConst(superclassName, location)

	c.emit(location.StartPos.Line, bytecode.SET_SUPERCLASS)
}

func (c *Compiler) CompileIvarIndices(target types.NamespaceWithIvarIndices, location *position.Location) {
	switch target := target.(type) {
	case *types.SingletonClass:
		c.emitGetConst(value.ToSymbol(target.AttachedObject.Name()), location)
		c.emit(location.StartPos.Line, bytecode.GET_SINGLETON)
	case *types.Module:
		c.emitGetConst(value.ToSymbol(target.Name()), location)
		c.emit(location.StartPos.Line, bytecode.GET_SINGLETON)
	default:
		c.emitGetConst(value.ToSymbol(target.Name()), location)
	}

	c.emitValue(value.Ref(target.IvarIndices()), location)
	c.emit(location.StartPos.Line, bytecode.DEF_IVARS)
}

func (c *Compiler) CompileInclude(target types.Namespace, mixin *types.Mixin, location *position.Location) {
	switch t := target.(type) {
	case *types.SingletonClass:
		targetName := value.ToSymbol(t.AttachedObject.Name())
		c.emitGetConst(targetName, location)
		c.emit(location.StartPos.Line, bytecode.GET_SINGLETON)
	default:
		targetName := value.ToSymbol(t.Name())
		c.emitGetConst(targetName, location)
	}

	mixinName := value.ToSymbol(mixin.Name())
	c.emitGetConst(mixinName, location)

	c.emit(location.StartPos.Line, bytecode.INCLUDE)
}

func (c *Compiler) InitExpressionCompiler(location *position.Location) *Compiler {
	exprCompiler := New("<file>", topLevelMode, location, c.checker)
	exprCompiler.Errors = c.Errors

	c.emitValue(value.Ref(exprCompiler.Bytecode), location)
	c.emit(location.StartPos.Line, bytecode.EXEC)
	c.emit(location.StartPos.Line, bytecode.POP)

	return exprCompiler
}

func (c *Compiler) CompileExpressionsInFile(node *ast.ProgramNode) {
	c.compileProgram(node)
}

// Entry point to the compilation process
func (c *Compiler) compileProgram(node ast.Node) {
	c.compileNode(node, false)
	c.emitReturn(node.Location(), nil)
	c.prepLocals()
}

func identifierToName(node ast.IdentifierNode) string {
	switch node := node.(type) {
	case *ast.PrivateIdentifierNode:
		return node.Value
	case *ast.PublicIdentifierNode:
		return node.Value
	default:
		return ""
	}
}

// Entry point for compiling the body of a function.
func (c *Compiler) compileFunction(location *position.Location, parameters []ast.ParameterNode, body []ast.StatementNode) {
	c.Bytecode.SetParameterCount(len(parameters))

	for _, param := range parameters {
		p := param.(*ast.FormalParameterNode)
		pSpan := p.Location()

		pName := identifierToName(p.Name)
		local := c.defineLocal(pName, pSpan)
		if local == nil {
			return
		}
		c.predefinedLocals++

		if p.Initialiser != nil {
			c.Bytecode.IncrementOptionalParameterCount()

			c.emitGetLocal(location.StartPos.Line, local.index)
			jump := c.emitJump(pSpan.StartPos.Line, bytecode.JUMP_UNLESS_UNDEF)

			c.compileNode(p.Initialiser, false)
			c.emitSetLocalPop(pSpan.StartPos.Line, local.index)

			c.patchJump(jump, pSpan)
		}
	}
	c.compileStatements(body, location, false)

	c.emitReturn(location, nil)
	c.prepLocals()
}

func (c *Compiler) InitMethodCompiler(location *position.Location) (*Compiler, int) {
	methodCompiler := New("<methodDefinitions>", topLevelMode, c.Bytecode.Location, c.checker)
	methodCompiler.Errors = c.Errors
	methodCompiler.Parent = c

	offset := c.nextInstructionOffset()
	c.emitValue(value.Ref(methodCompiler.Bytecode), location)
	c.emit(location.StartPos.Line, bytecode.EXEC)
	c.emit(location.StartPos.Line, bytecode.POP)

	return methodCompiler, offset
}

var ivarIndicesSymbol = value.ToSymbol("<ivarIndices>")

func (c *Compiler) InitIvarIndicesCompiler(location *position.Location) (*Compiler, int) {
	ivarCompiler := New(ivarIndicesSymbol.String(), topLevelMode, c.Bytecode.Location, c.checker)
	ivarCompiler.Errors = c.Errors
	ivarCompiler.Parent = c

	offset := c.nextInstructionOffset()
	c.emitValue(value.Ref(ivarCompiler.Bytecode), location)
	c.emit(location.StartPos.Line, bytecode.EXEC)
	c.emit(location.StartPos.Line, bytecode.POP)

	return ivarCompiler, offset
}

func (c *Compiler) FinishIvarIndicesCompiler(location *position.Location, execOffset int) *Compiler {
	if len(c.Bytecode.Instructions) > 0 {
		c.emit(location.EndPos.Line, bytecode.NIL)
		c.emit(location.EndPos.Line, bytecode.RETURN)
		return c.Parent
	}

	// If no instructions were emitted, remove the EXEC instruction block
	c.Parent.removeLastBytes(execOffset)
	c.Parent.removeBytecodeFunction(ivarIndicesSymbol)
	return c.Parent
}

func (c *Compiler) CompileMethods(location *position.Location, execOffset int) {
	c.compileMethodsWithinModule(c.checker.Env().Root, location)
	if len(c.Bytecode.Instructions) > 0 {
		c.emit(location.EndPos.Line, bytecode.NIL)
		c.emit(location.EndPos.Line, bytecode.RETURN)
		return
	}

	// If no instructions were emitted, remove the EXEC instruction block
	c.Parent.removeLastBytes(execOffset)
	c.Parent.removeBytecodeFunction(methodDefinitionsSymbol)
}

func (c *Compiler) removeLastBytes(offset int) {
	c.removeBytes(offset, c.nextInstructionOffset()-offset)
}

func (c *Compiler) removeBytes(offset int, count int) {
	c.Bytecode.Instructions = slices.Concat(c.Bytecode.Instructions[:offset], c.Bytecode.Instructions[offset+count:])
	lineInfo := c.Bytecode.LineInfoList.GetLineInfo(offset)
	lineInfo.InstructionCount -= count
}

var methodDefinitionsSymbol = value.ToSymbol("<methodDefinitions>")

func (c *Compiler) removeBytecodeFunction(name value.Symbol) {
	for i, val := range c.Bytecode.Values {
		val, ok := val.SafeAsReference().(*vm.BytecodeFunction)
		if !ok {
			continue
		}

		if val.Name() == name {
			c.Bytecode.Values[i] = value.Undefined
			if i == len(c.Bytecode.Values)-1 {
				c.Bytecode.Values = c.Bytecode.Values[:len(c.Bytecode.Values)-1]
			}
			break
		}
	}
}

func (c *Compiler) compileMethodsWithinModule(module *types.Module, location *position.Location) {
	if types.NamespaceHasAnyDefinableMethods(module) {
		c.emitGetConst(value.ToSymbol(module.Name()), location)
		c.emit(location.StartPos.Line, bytecode.GET_SINGLETON)

		for methodName, method := range types.SortedOwnMethods(module) {
			c.compileMethodDefinition(methodName, method, location)

			for i, overload := range method.Overloads {
				overloadName := value.ToSymbol(
					fmt.Sprintf("%s@%d", methodName.String(), i+1),
				)
				c.compileMethodDefinition(overloadName, overload, location)
			}
		}

		c.emit(location.StartPos.Line, bytecode.POP)
	}

	for _, subtype := range types.SortedSubtypes(module) {
		if subtype.Type == module {
			continue
		}
		c.compileMethodsWithinType(subtype.Type, location)
	}
}

func (c *Compiler) compileMethodDefinition(name value.Symbol, method *types.Method, location *position.Location) {
	if !method.IsDefinable() {
		return
	}

	if method.Base != nil {
		method = method.Base
	}

	if method.IsAttribute() {
		if method.IsSetter() {
			nameStr := name.String()
			ivarName := value.ToSymbol(nameStr[:len(nameStr)-1])
			c.emitValue(ivarName.ToValue(), location)
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
			c.emitSmallInt(value.SmallInt(index), location)
			c.emit(location.StartPos.Line, bytecode.DEF_SETTER)
			method.SetCompiled(true)
			return
		}

		c.emitValue(name.ToValue(), location)
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
		c.emitSmallInt(value.SmallInt(index), location)
		c.emit(location.StartPos.Line, bytecode.DEF_GETTER)
		method.SetCompiled(true)
		return
	}

	c.emitValue(value.Ref(method.Body), location)
	c.emitValue(name.ToValue(), location)
	c.emit(location.StartPos.Line, bytecode.DEF_METHOD)
	method.SetCompiled(true)
}

func (c *Compiler) compileMethodsWithinNamespace(namespace types.Namespace, location *position.Location) {
	namespaceHasCompiledMethods := types.NamespaceHasAnyDefinableMethods(namespace)

	singleton := namespace.Singleton()
	singletonHasCompiledMethods := types.NamespaceHasAnyDefinableMethods(singleton)

	if namespaceHasCompiledMethods || singletonHasCompiledMethods {
		c.emitGetConst(value.ToSymbol(namespace.Name()), location)

		for methodName, method := range types.SortedOwnMethods(namespace) {
			c.compileMethodDefinition(methodName, method, location)
		}

		if singletonHasCompiledMethods {
			c.emit(location.StartPos.Line, bytecode.GET_SINGLETON)

			for methodName, method := range types.SortedOwnMethods(singleton) {
				c.compileMethodDefinition(methodName, method, location)
			}
		}

		c.emit(location.StartPos.Line, bytecode.POP)
	}

	for _, subtype := range types.SortedSubtypes(namespace) {
		if subtype.Type == namespace {
			continue
		}
		c.compileMethodsWithinType(subtype.Type, location)
	}
}

func (c *Compiler) compileMethodsWithinType(typ types.Type, location *position.Location) {
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

func (c *Compiler) CompileMethodBody(node *ast.MethodDefinitionNode, name value.Symbol) *vm.BytecodeFunction {
	var mode mode
	if node.IsSetter() {
		mode = setterMethodMode
	} else if identifierToName(node.Name) == "#init" {
		mode = initMethodMode
	} else {
		mode = methodMode
	}

	methodCompiler := New(name.String(), mode, node.Location(), c.checker)
	methodCompiler.isGenerator = node.IsGenerator()
	methodCompiler.isAsync = node.IsAsync()
	methodCompiler.Errors = c.Errors
	methodCompiler.compileMethodBody(node.Location(), node.Parameters, node.Body)

	return methodCompiler.Bytecode
}

func (c *Compiler) CompileMacroBody(node *ast.MacroDefinitionNode, name value.Symbol) *vm.BytecodeFunction {
	methodCompiler := New(name.String(), macroMode, node.Location(), c.checker)
	methodCompiler.Errors = c.Errors
	methodCompiler.compileMacroBody(node.Location(), node.Parameters, node.Body)

	return methodCompiler.Bytecode
}

const macroLocationParamName = "_location"

// Entry point for compiling the body of a macro.
func (c *Compiler) compileMacroBody(location *position.Location, parameters []ast.ParameterNode, body []ast.StatementNode) {
	paramCount := len(parameters)

	for _, param := range parameters {
		p := param.(*ast.FormalParameterNode)
		pSpan := p.Location()

		pName := identifierToName(p.Name)
		local := c.defineLocal(pName, pSpan)
		if local == nil {
			return
		}
		c.predefinedLocals++

		if p.Initialiser != nil {
			c.Bytecode.IncrementOptionalParameterCount()

			c.emitGetLocal(location.StartPos.Line, local.index)
			jump := c.emitJump(pSpan.StartPos.Line, bytecode.JUMP_UNLESS_UNDEF)

			c.compileNode(p.Initialiser, false)
			c.emitSetLocalPop(pSpan.StartPos.Line, local.index)

			c.patchJump(jump, pSpan)
		}
	}

	c.Bytecode.SetParameterCount(paramCount)

	c.compileStatements(body, location, false)

	c.emitReturn(location, nil)
	c.emitFinalReturn(location, nil)
	c.prepLocals()
}

// Entry point for compiling the body of a method.
func (c *Compiler) compileMethodBody(location *position.Location, parameters []ast.ParameterNode, body []ast.StatementNode) {
	for _, param := range parameters {
		p := param.(*ast.MethodParameterNode)
		pSpan := p.Location()

		pName := identifierToName(p.Name)
		local := c.defineLocal(pName, pSpan)
		if local == nil {
			return
		}
		c.predefinedLocals++

		if p.Initialiser != nil {
			c.Bytecode.IncrementOptionalParameterCount()

			c.emitGetLocal(location.StartPos.Line, local.index)
			jump := c.emitJump(pSpan.StartPos.Line, bytecode.JUMP_UNLESS_UNDEF)

			c.compileNode(p.Initialiser, false)
			c.emitSetLocalPop(pSpan.StartPos.Line, local.index)

			c.patchJump(jump, pSpan)
		}

		if p.SetInstanceVariable {
			c.emitGetLocal(location.StartPos.Line, local.index)
			c.emitSetInstanceVariableNoPop(value.ToSymbol(pName), pSpan)
			// pop the value after setting it
			c.emit(pSpan.StartPos.Line, bytecode.POP)
		}
	}

	paramCount := len(parameters)
	if c.isGenerator {
		c.emit(location.StartPos.Line, bytecode.GENERATOR)
		c.emit(location.EndPos.Line, bytecode.RETURN)
		c.registerCatch(-1, -1, c.nextInstructionOffset(), false)
	} else if c.isAsync {
		poolVar := c.defineLocal("_pool", location)
		paramCount++
		c.predefinedLocals++
		c.Bytecode.IncrementOptionalParameterCount()

		c.emitGetLocal(location.StartPos.Line, poolVar.index)
		c.emit(location.StartPos.Line, bytecode.PROMISE)
		c.emit(location.EndPos.Line, bytecode.RETURN)
	}
	c.Bytecode.SetParameterCount(paramCount)

	c.compileStatements(body, location, false)

	c.emitReturn(location, nil)
	c.emitFinalReturn(location, nil)
	c.prepLocals()
}

// Entry point for compiling the body of a namespace eg. Module, Class, Mixin, Struct, Interface.
func (c *Compiler) compileNamespace(node ast.Node) bool {
	location := node.Location()
	switch n := node.(type) {
	case *ast.ClassDeclarationNode:
		if !c.compileStatementsOk(n.Body) {
			return false
		}
	case *ast.InterfaceDeclarationNode:
		if !c.compileStatementsOk(n.Body) {
			return false
		}
	case *ast.ModuleDeclarationNode:
		if !c.compileStatementsOk(n.Body) {
			return false
		}
	case *ast.MixinDeclarationNode:
		if !c.compileStatementsOk(n.Body) {
			return false
		}
	case *ast.SingletonBlockExpressionNode:
		if !c.compileStatementsOk(n.Body) {
			return false
		}
	default:
		c.Errors.AddFailure(
			fmt.Sprintf("incorrect namespace type %#v", n),
			location,
		)
		return false
	}

	c.emitReturn(location, nil)
	c.prepLocals()
	return true
}

func (c *Compiler) prepLocals() {
	localCount := c.maxLocalIndex + 1 - c.predefinedLocals
	if localCount == 0 {
		return
	}

	var newInstructions []byte
	var newBytes int
	if c.maxLocalIndex >= math.MaxUint8 {
		newBytes = 3
		newInstructions = make([]byte, 0, len(c.Bytecode.Instructions)+newBytes)
		newInstructions = append(newInstructions, byte(bytecode.PREP_LOCALS16))
		newInstructions = binary.BigEndian.AppendUint16(newInstructions, uint16(localCount))
	} else {
		newBytes = 2
		newInstructions = make([]byte, 0, len(c.Bytecode.Instructions)+newBytes)
		newInstructions = append(newInstructions, byte(bytecode.PREP_LOCALS8), byte(localCount))
	}

	c.Bytecode.Instructions = append(
		newInstructions,
		c.Bytecode.Instructions...,
	)
	lineInfo := c.Bytecode.LineInfoList.First()
	if lineInfo != nil {
		lineInfo.InstructionCount += newBytes
	}
	for _, catchEntry := range c.Bytecode.CatchEntries {
		catchEntry.From += len(newInstructions)
		catchEntry.To += len(newInstructions)
		catchEntry.JumpAddress += len(newInstructions)
	}

	for _, id := range c.offsetValueIds {
		currentValue := c.Bytecode.Values[id].MustSmallInt()
		c.Bytecode.Values[id] = (currentValue + value.SmallInt(len(newInstructions))).ToValue()
	}
}

func (c *Compiler) initLoopJumpSet(label string, returnsValFromLastIteration bool) {
	c.loopJumpSets = append(
		c.loopJumpSets,
		&loopJumpSet{
			label:                         label,
			returnsValueFromLastIteration: returnsValFromLastIteration,
		},
	)
}

func (c *Compiler) findLoopJumpSet(label string, location *position.Location) *loopJumpSet {
	if len(c.loopJumpSets) < 1 {
		c.Errors.AddFailure(
			"cannot jump with `break` or `continue` outside of a loop",
			location,
		)
		return nil
	}

	if label == "" {
		// if there is no label, choose the closest enclosing loop
		return c.loopJumpSets[len(c.loopJumpSets)-1]
	}

	for _, currentJumpSet := range c.loopJumpSets {
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

func (c *Compiler) addLoopJumpTo(jumpSet *loopJumpSet, typ loopJumpInfoType, offset int) {
	jumpSet.loopJumps = append(
		jumpSet.loopJumps,
		&loopJumpInfo{
			typ:    typ,
			offset: offset,
		},
	)
}

func (c *Compiler) addLoopJump(label string, typ loopJumpInfoType, offset int, location *position.Location) {
	jumpSet := c.findLoopJumpSet(label, location)
	if jumpSet == nil {
		return
	}

	c.addLoopJumpTo(jumpSet, typ, offset)
}

func (c *Compiler) compilePublicConstantNode(node *ast.PublicConstantNode) {
	c.emitGetConst(value.ToSymbol(node.Value), node.Location())
}

func (c *Compiler) compilePrivateConstantNode(node *ast.PrivateConstantNode) {
	c.emitGetConst(value.ToSymbol(node.Value), node.Location())
}

func (c *Compiler) nodeIsCompilable(node ast.Node) bool {
	switch node := node.(type) {
	case nil, *ast.AliasDeclarationNode, *ast.IncludeExpressionNode,
		*ast.EmptyStatementNode, *ast.MethodDefinitionNode, *ast.UsingExpressionNode,
		*ast.ConstantDeclarationNode, *ast.TypeDefinitionNode, *ast.GenericTypeDefinitionNode,
		*ast.MethodSignatureDefinitionNode, *ast.ImplementExpressionNode,
		*ast.StructDeclarationNode, *ast.GenericReceiverlessMethodCallNode,
		*ast.ReceiverlessMethodCallNode, *ast.AttrDeclarationNode,
		*ast.SetterDeclarationNode, *ast.GetterDeclarationNode, *ast.InitDefinitionNode,
		*ast.InstanceVariableDeclarationNode, *ast.MacroDefinitionNode,
		*ast.ReceiverlessMacroCallNode, *ast.MacroCallNode, *ast.ScopedMacroCallNode,
		*ast.ImportStatementNode:
		return false
	case *ast.ExpressionStatementNode:
		return c.nodeIsCompilable(node.Expression)
	case *ast.InterfaceDeclarationNode:
		return c.interfaceIsCompilable(node)
	case *ast.ClassDeclarationNode:
		return c.classIsCompilable(node)
	case *ast.ModuleDeclarationNode:
		return c.moduleIsCompilable(node)
	case *ast.MixinDeclarationNode:
		return c.mixinIsCompilable(node)
	case *ast.SingletonBlockExpressionNode:
		return c.singletonBlockIsCompilable(node)
	default:
		return true
	}
}

type expressionResult uint8

const (
	expressionCompiled              expressionResult = iota // expression has been compiled and can be popped
	expressionIgnored                                       // expression was ignored
	expressionCompiledWithoutResult                         // expression has been successfully compiled but should not be popped
)

func (c *Compiler) compileNodeWithoutResult(node ast.Node) {
	if c.compileNode(node, true) == expressionCompiled {
		c.emit(node.Location().EndPos.Line, bytecode.POP)
	}
}

func (c *Compiler) compileNodeWithResult(node ast.Node) {
	switch c.compileNode(node, false) {
	case expressionCompiledWithoutResult, expressionIgnored:
		c.emit(node.Location().EndPos.Line, bytecode.NIL)
	}
}

func (c *Compiler) mustCompileNode(node ast.Node, valueIsIgnored bool) {
	if valueIsIgnored {
		c.compileNodeWithoutResult(node)
	} else {
		c.compileNodeWithResult(node)
	}
}

func (c *Compiler) compileNode(node ast.Node, valueIsIgnored bool) expressionResult {
	switch node := node.(type) {
	case nil, *ast.AliasDeclarationNode, *ast.IncludeExpressionNode,
		*ast.EmptyStatementNode, *ast.MethodDefinitionNode, *ast.UsingExpressionNode,
		*ast.ConstantDeclarationNode, *ast.TypeDefinitionNode, *ast.GenericTypeDefinitionNode,
		*ast.MethodSignatureDefinitionNode, *ast.ImplementExpressionNode,
		*ast.StructDeclarationNode, *ast.GenericReceiverlessMethodCallNode,
		*ast.ReceiverlessMethodCallNode, *ast.AttrDeclarationNode,
		*ast.SetterDeclarationNode, *ast.GetterDeclarationNode, *ast.InitDefinitionNode,
		*ast.InstanceVariableDeclarationNode, *ast.MacroDefinitionNode,
		*ast.ReceiverlessMacroCallNode, *ast.MacroCallNode, *ast.ScopedMacroCallNode,
		*ast.ImportStatementNode:
		return expressionIgnored
	case *ast.ProgramNode:
		return c.compileStatements(node.Body, node.Location(), valueIsIgnored)
	case *ast.ExtendWhereBlockExpressionNode:
		c.compileStatements(node.Body, node.Location(), false)
	case *ast.ExpressionStatementNode:
		return c.compileNode(node.Expression, valueIsIgnored)
	case *ast.LabeledExpressionNode:
		return c.compileLabeledExpressionNode(node, valueIsIgnored)
	case *ast.UndefinedLiteralNode:
		c.emit(node.Location().StartPos.Line, bytecode.UNDEFINED)
	case *ast.PublicConstantNode:
		c.compilePublicConstantNode(node)
	case *ast.PrivateConstantNode:
		c.compilePrivateConstantNode(node)
	case *ast.GenericConstantNode:
		return c.compileNode(node.Constant, valueIsIgnored)
	case *ast.SelfLiteralNode:
		c.emit(node.Location().StartPos.Line, bytecode.SELF)
	case *ast.AssignmentExpressionNode:
		return c.compileAssignmentExpressionNode(node, valueIsIgnored)
	case *ast.InterfaceDeclarationNode:
		return c.compileInterfaceDeclarationNode(node)
	case *ast.ClassDeclarationNode:
		return c.compileClassDeclarationNode(node)
	case *ast.ModuleDeclarationNode:
		return c.compileModuleDeclarationNode(node)
	case *ast.MixinDeclarationNode:
		return c.compileMixinDeclarationNode(node)
	case *ast.SingletonBlockExpressionNode:
		return c.compileSingletonBlockExpressionNode(node)
	case *ast.GoExpressionNode:
		c.compileGoExpressionNode(node)
	case *ast.ClosureLiteralNode:
		c.compileClosureLiteralNode(node)
	case *ast.SwitchExpressionNode:
		return c.compileSwitchExpressionNode(node, valueIsIgnored)
	case *ast.SubscriptExpressionNode:
		c.compileSubscriptExpressionNode(node)
	case *ast.NilSafeSubscriptExpressionNode:
		return c.compileNilSafeSubscriptExpressionNode(node)
	case *ast.AttributeAccessNode:
		c.compileAttributeAccessNode(node)
	case *ast.NewExpressionNode:
		c.compileNewExpressionNode(node)
	case *ast.ConstructorCallNode:
		c.compileConstructorCallNode(node)
	case *ast.GenericConstructorCallNode:
		c.compileGenericConstructorCallNode(node)
	case *ast.MethodCallNode:
		c.compileMethodCallNode(node)
	case *ast.GenericMethodCallNode:
		c.compileGenericMethodCallNode(node)
	case *ast.CallNode:
		c.compileCallNode(node)
	case *ast.ReturnExpressionNode:
		c.compileReturnExpressionNode(node)
		return expressionCompiledWithoutResult
	case *ast.AwaitExpressionNode:
		c.compileAwaitExpressionNode(node)
	case *ast.YieldExpressionNode:
		c.compileYieldExpressionNode(node)
		return expressionCompiledWithoutResult
	case *ast.VariablePatternDeclarationNode:
		c.compilerVariablePatternDeclarationNode(node)
	case *ast.VariableDeclarationNode:
		return c.compileVariableDeclarationNode(node, valueIsIgnored)
	case *ast.ValuePatternDeclarationNode:
		c.compileValuePatternDeclarationNode(node)
	case *ast.ValueDeclarationNode:
		return c.compileValueDeclarationNode(node, valueIsIgnored)
	case *ast.PublicIdentifierNode:
		c.compileLocalVariableAccess(node.Value, node.Location())
	case *ast.PrivateIdentifierNode:
		c.compileLocalVariableAccess(node.Value, node.Location())
	case *ast.PublicInstanceVariableNode:
		c.compileInstanceVariableAccess(node.Value, node.Location())
	case *ast.BinaryExpressionNode:
		c.compileBinaryExpressionNode(node)
	case *ast.LogicalExpressionNode:
		return c.compileLogicalExpressionNode(node, valueIsIgnored)
	case *ast.UnaryExpressionNode:
		c.compileUnaryExpressionNode(node)
	case *ast.RangeLiteralNode:
		c.compileRangeLiteralNode(node)
	case *ast.HashSetLiteralNode:
		c.compileHashSetLiteralNode(node)
	case *ast.HashMapLiteralNode:
		c.compileHashMapLiteralNode(node)
	case *ast.HashRecordLiteralNode:
		c.compileHashRecordLiteralNode(node)
	case *ast.ArrayTupleLiteralNode:
		c.compileArrayTupleLiteralNode(node)
	case *ast.WordArrayTupleLiteralNode:
		c.compileWordArrayTupleLiteralNode(node)
	case *ast.SymbolArrayTupleLiteralNode:
		c.compileSymbolArrayTupleLiteralNode(node)
	case *ast.BinArrayTupleLiteralNode:
		c.compileBinArrayTupleLiteralNode(node)
	case *ast.HexArrayTupleLiteralNode:
		c.compileHexArrayTupleLiteralNode(node)
	case *ast.ArrayListLiteralNode:
		c.compileArrayListLiteralNode(node)
	case *ast.WordArrayListLiteralNode:
		c.compileWordArrayListLiteralNode(node)
	case *ast.SymbolArrayListLiteralNode:
		c.compileSymbolArrayListLiteralNode(node)
	case *ast.BinArrayListLiteralNode:
		c.compileBinArrayListLiteralNode(node)
	case *ast.HexArrayListLiteralNode:
		c.compileHexArrayListLiteralNode(node)
	case *ast.WordHashSetLiteralNode:
		c.compileWordHashSetLiteralNode(node)
	case *ast.SymbolHashSetLiteralNode:
		c.compileSymbolHashSetLiteralNode(node)
	case *ast.BinHashSetLiteralNode:
		c.compileBinHashSetLiteralNode(node)
	case *ast.HexHashSetLiteralNode:
		c.compileHexHashSetLiteralNode(node)
	case *ast.UninterpolatedRegexLiteralNode:
		c.compileUninterpolatedRegexLiteralNode(node)
	case *ast.InterpolatedRegexLiteralNode:
		c.compileInterpolatedRegexLiteralNode(node)
	case *ast.RawStringLiteralNode:
		c.emitValue(value.Ref(value.String(node.Value)), node.Location())
	case *ast.DoubleQuotedStringLiteralNode:
		c.emitValue(value.Ref(value.String(node.Value)), node.Location())
	case *ast.InterpolatedStringLiteralNode:
		c.compileInterpolatedStringLiteralNode(node)
	case *ast.InterpolatedSymbolLiteralNode:
		c.compileInterpolatedSymbolLiteralNode(node)
	case *ast.CharLiteralNode:
		c.emitValue(value.Char(node.Value).ToValue(), node.Location())
	case *ast.RawCharLiteralNode:
		c.emitValue(value.Char(node.Value).ToValue(), node.Location())
	case *ast.FalseLiteralNode:
		c.emit(node.Location().StartPos.Line, bytecode.FALSE)
	case *ast.TrueLiteralNode:
		c.emit(node.Location().StartPos.Line, bytecode.TRUE)
	case *ast.NilLiteralNode:
		c.emit(node.Location().StartPos.Line, bytecode.NIL)
	case *ast.ThrowExpressionNode:
		c.compileThrowExpressionNode(node)
	case *ast.MustExpressionNode:
		c.compileMustExpressionNode(node)
	case *ast.TryExpressionNode:
		return c.compileTryExpressionNode(node, valueIsIgnored)
	case *ast.AsExpressionNode:
		c.compileAsExpressionNode(node)
	case *ast.TypeofExpressionNode:
		return c.compileTypeofExpressionNode(node, valueIsIgnored)
	case *ast.DoExpressionNode:
		c.compileDoExpressionNode(node)
	case *ast.MacroBoundaryNode:
		c.compileMacroBoundaryNode(node)
	case *ast.UnhygienicNode:
		c.compileUnhygienicNode(node, valueIsIgnored)
	case *ast.IfExpressionNode:
		return c.compileIfExpression(
			false,
			node.Condition,
			node.ThenBody,
			node.ElseBody,
			node.Location(),
			valueIsIgnored,
		)
	case *ast.UnlessExpressionNode:
		return c.compileIfExpression(
			true,
			node.Condition,
			node.ThenBody,
			node.ElseBody,
			node.Location(),
			valueIsIgnored,
		)
	case *ast.ModifierIfElseNode:
		return c.compileModifierIfExpression(
			false,
			node.Condition,
			node.ThenExpression,
			node.ElseExpression,
			node.Location(),
			valueIsIgnored,
		)
	case *ast.ModifierNode:
		return c.compileModifierExpressionNode("", node, valueIsIgnored)
	case *ast.BreakExpressionNode:
		c.compileBreakExpressionNode(node)
	case *ast.ContinueExpressionNode:
		c.compileContinueExpressionNode(node)
	case *ast.LoopExpressionNode:
		c.compileLoopExpressionNode("", node.ThenBody, node.Location())
	case *ast.WhileExpressionNode:
		c.compileWhileExpressionNode("", node)
	case *ast.UntilExpressionNode:
		c.compileUntilExpressionNode("", node)
	case *ast.NumericForExpressionNode:
		c.compileNumericForExpressionNode("", node)
	case *ast.ForInExpressionNode:
		c.compileForInExpressionNode("", node)
	case *ast.ModifierForInNode:
		c.compileModifierForInNode("", node)
	case *ast.PostfixExpressionNode:
		return c.compilePostfixExpressionNode(node, valueIsIgnored)
	case *ast.SimpleSymbolLiteralNode:
		c.emitValue(value.ToSymbol(node.Content).ToValue(), node.Location())
	case *ast.IntLiteralNode:
		c.compileIntLiteralNode(node)
	case *ast.Int8LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 8)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return expressionCompiled
		}
		c.emitValue(value.Int8(i).ToValue(), node.Location())
	case *ast.Int16LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 16)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return expressionCompiled
		}
		c.emitValue(value.Int16(i).ToValue(), node.Location())
	case *ast.Int32LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 32)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return expressionCompiled
		}
		c.emitValue(value.Int32(i).ToValue(), node.Location())
	case *ast.Int64LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 64)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return expressionCompiled
		}
		c.emitValue(value.Int64(i).ToValue(), node.Location())
	case *ast.UInt8LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 8)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return expressionCompiled
		}
		c.emitValue(value.UInt8(i).ToValue(), node.Location())
	case *ast.UInt16LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 16)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return expressionCompiled
		}
		c.emitValue(value.UInt16(i).ToValue(), node.Location())
	case *ast.UInt32LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 32)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return expressionCompiled
		}
		c.emitValue(value.UInt32(i).ToValue(), node.Location())
	case *ast.UInt64LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 64)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return expressionCompiled
		}
		c.emitValue(value.UInt64(i).ToValue(), node.Location())
	case *ast.FloatLiteralNode:
		f, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			c.Errors.AddFailure(err.Error(), node.Location())
			return expressionCompiled
		}
		c.emitValue(value.Float(f).ToValue(), node.Location())
	case *ast.BigFloatLiteralNode:
		f, err := value.ParseBigFloat(node.Value)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return expressionCompiled
		}
		c.emitValue(value.Ref(f), node.Location())
	case *ast.Float64LiteralNode:
		f, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			c.Errors.AddFailure(err.Error(), node.Location())
			return expressionCompiled
		}
		c.emitValue(value.Float64(f).ToValue(), node.Location())
	case *ast.Float32LiteralNode:
		f, err := strconv.ParseFloat(node.Value, 32)
		if err != nil {
			c.Errors.AddFailure(err.Error(), node.Location())
			return expressionCompiled
		}
		c.emitValue(value.Float32(f).ToValue(), node.Location())
	case *ast.QuoteExpressionNode:
		c.compileQuoteExpressionNode(node)
	case *ast.UnquoteNode:
		c.compileUnquoteNode(node)
	default:
		c.Errors.AddFailure(
			fmt.Sprintf("compilation of this node has not been implemented: %T", node),
			node.Location(),
		)
	}

	return expressionCompiled
}

func (c *Compiler) compileUnquoteNode(node *ast.UnquoteNode) {
	c.compileNodeWithResult(node.Expression)
	var methodName value.Symbol

	switch node.Kind {
	case ast.UNQUOTE_EXPRESSION_KIND:
		methodName = symbol.L_to_ast_expr_node
	case ast.UNQUOTE_CONSTANT_KIND:
		methodName = symbol.L_to_ast_const_node
	case ast.UNQUOTE_PATTERN_KIND:
		methodName = symbol.L_to_ast_pattern_node
	case ast.UNQUOTE_PATTERN_EXPRESSION_KIND:
		methodName = symbol.L_to_ast_pattern_expr_node
	case ast.UNQUOTE_TYPE_KIND:
		methodName = symbol.L_to_ast_type_node
	case ast.UNQUOTE_IDENTIFIER_KIND:
		methodName = symbol.L_to_ast_ident_node
	case ast.UNQUOTE_INSTANCE_VARIABLE_KIND:
		methodName = symbol.L_to_ast_ivar_node
	}

	c.emitCallMethod(
		value.NewCallSiteInfo(
			methodName,
			0,
		),
		node.Location(),
		false,
	)
}

func (c *Compiler) compileQuoteExpressionNode(node *ast.QuoteExpressionNode) {
	location := node.Location()

	var newNode ast.ExpressionNode
	if expr := node.SingleExpression(); expr != nil {
		newNode = expr
	} else {
		newNode = ast.NewDoExpressionNode(location, node.Body, nil, nil)
	}

	c.emitGetConst(value.ToSymbol("Std::Kernel"), location)
	c.emitValue(value.Ref(newNode), location)

	tupleBaseOffset := c.nextInstructionOffset()
	// ArrayTuple base
	c.emit(location.StartPos.Line, bytecode.UNDEFINED)

	var unquoteCount int
	ast.Traverse(
		node,
		func(node, parent ast.Node) ast.TraverseOption {
			switch node := node.(type) {
			case *ast.UnquoteNode:
				unquoteCount++
				c.compileNodeWithResult(node)

				return ast.TraverseSkip
			}

			return ast.TraverseContinue
		},
		nil,
	)
	var argCount int

	if unquoteCount != 0 {
		argCount = 2
		c.emitNewArrayTuple(unquoteCount, location)
	} else {
		argCount = 1
		c.removeBytes(tupleBaseOffset, 1)
	}

	c.emitCallMethod(
		value.NewCallSiteInfo(
			symbol.S_splice,
			argCount,
		),
		location,
		false,
	)
}

func (c *Compiler) compileTypeofExpressionNode(node *ast.TypeofExpressionNode, valueIsIgnored bool) expressionResult {
	return c.compileNode(node.Value, valueIsIgnored)
}

func (c *Compiler) compileTryExpressionNode(node *ast.TryExpressionNode, valueIsIgnored bool) expressionResult {
	return c.compileNode(node.Value, valueIsIgnored)
}

func (c *Compiler) compileMustExpressionNode(node *ast.MustExpressionNode) {
	location := node.Location()
	c.compileNodeWithResult(node.Value)
	c.emit(location.StartPos.Line, bytecode.MUST)
}

func (c *Compiler) compileAsExpressionNode(node *ast.AsExpressionNode) {
	location := node.Location()
	c.compileNode(node.Value, false)
	c.compileNode(node.RuntimeType, false)
	c.emit(location.StartPos.Line, bytecode.AS)
}

func (c *Compiler) compileThrowExpressionNode(node *ast.ThrowExpressionNode) {
	location := node.Location()
	if node.Value != nil {
		c.compileNode(node.Value, false)
	} else {
		c.emitValue(value.Ref(value.NewError(value.ErrorClass, "error")), location)
	}

	c.emit(location.StartPos.Line, bytecode.THROW)
}

func (c *Compiler) isNestedInFinally() bool {
	for _, scope := range c.scopes {
		if scope.typ == doFinallyScopeType {
			return true
		}
	}

	return false
}

func (c *Compiler) registerCatch(from, to, jumpAddress int, finally bool) {
	doCatchEntry := vm.NewCatchEntry(
		from,
		to,
		jumpAddress,
		finally,
	)
	c.Bytecode.CatchEntries = append(
		c.Bytecode.CatchEntries,
		doCatchEntry,
	)
}

func (c *Compiler) CompileConstantDeclaration(node *ast.ConstantDeclarationNode, namespace types.Namespace, constName value.Symbol) {
	location := node.Location()
	switch n := namespace.(type) {
	case *types.SingletonClass:
		namespaceName := value.ToSymbol(n.AttachedObject.Name())
		c.emitGetConst(namespaceName, node.Constant.Location())
		c.emit(location.StartPos.Line, bytecode.GET_SINGLETON)
	default:
		namespaceName := value.ToSymbol(n.Name())
		c.emitGetConst(namespaceName, node.Constant.Location())
	}
	c.emitValue(constName.ToValue(), location)
	c.compileNode(node.Initialiser, false)
	c.emit(location.StartPos.Line, bytecode.DEF_CONST)
}

func (c *Compiler) compileDoExpressionNode(node *ast.DoExpressionNode) {
	location := node.Location()

	doStartOffset := c.nextInstructionOffset()

	var scopeType scopeType
	if len(node.Finally) > 0 {
		scopeType = doFinallyScopeType
	} else {
		scopeType = defaultScopeType
	}

	c.enterScope("", scopeType)
	c.compileStatementsWithResult(node.Body, location)
	c.leaveScope(location.EndPos.Line)

	doEndOffset := c.nextInstructionOffset()

	if len(node.Finally) > 0 {
		c.enterScope("", defaultScopeType)
		// pop the return value of finally leaving the return value of do
		c.compileStatementsWithoutResult(node.Finally)
		c.leaveScope(location.EndPos.Line)
	}

	if len(node.Catches) <= 0 && len(node.Finally) <= 0 {
		return
	}

	jumpOverCatchOffset := c.emitJump(location.StartPos.Line, bytecode.JUMP)

	var jumpsToEndOfCatch []int
	catchStartOffset := c.nextInstructionOffset()

	c.registerCatch(doStartOffset, doEndOffset, catchStartOffset, false)

	c.enterScope("", defaultScopeType)

	for _, catchNode := range node.Catches {
		location := catchNode.Location()

		if catchNode.StackTraceVar != nil {
			c.emit(location.EndPos.Line, bytecode.DUP_SECOND)

			var stackTraceVarName string
			switch s := catchNode.StackTraceVar.(type) {
			case *ast.PublicIdentifierNode:
				stackTraceVarName = s.Value
			case *ast.PrivateIdentifierNode:
				stackTraceVarName = s.Value
			default:
				panic(fmt.Sprintf("invalid stack trace variable name in catch: %T", catchNode.StackTraceVar))
			}

			stackTraceVar := c.defineLocal(stackTraceVarName, location)
			if stackTraceVar != nil {
				c.emitSetLocalPop(location.StartPos.Line, stackTraceVar.index)
			}
		}

		c.pattern(catchNode.Pattern)
		jumpOverCatchBody := c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS)

		c.compileStatementsWithResult(catchNode.Body, catchNode.Location())

		if len(node.Finally) < 1 {
			// pop the thrown value and the stack trace, leaving the return value of the catch
			c.emit(location.EndPos.Line, bytecode.POP_2_SKIP_ONE)
		}
		jump := c.emitJump(location.EndPos.Line, bytecode.JUMP)
		jumpsToEndOfCatch = append(jumpsToEndOfCatch, jump)

		c.patchJump(jumpOverCatchBody, location)
	}

	if len(node.Finally) > 0 {
		c.emit(location.EndPos.Line, bytecode.TRUE)
	} else {
		c.emit(location.EndPos.Line, bytecode.RETHROW)
	}

	var jumpOverFalseOffset int
	if len(node.Finally) > 0 {

		jumpOverFalseOffset = c.emitJump(location.EndPos.Line, bytecode.JUMP)
	}
	for _, jump := range jumpsToEndOfCatch {
		c.patchJump(jump, location)
	}
	if len(node.Finally) > 0 {
		c.emit(location.EndPos.Line, bytecode.FALSE)
		c.patchJump(jumpOverFalseOffset, location)

		jumpOverReturnBreakOrContinueEntryOffset := c.emitJump(location.EndPos.Line, bytecode.JUMP)
		finallyEntryOffset := c.nextInstructionOffset()
		c.registerCatch(doStartOffset, doEndOffset, finallyEntryOffset, true)
		// entry point for return when executing finally
		c.emit(location.EndPos.Line, bytecode.NIL)

		jumpOverBreakOrContinueEntryOffset := c.emitJump(location.EndPos.Line, bytecode.JUMP)
		// entry point for break or continue when executing finally
		c.emit(location.EndPos.Line, bytecode.UNDEFINED)

		c.patchJump(jumpOverBreakOrContinueEntryOffset, location)
		c.patchJump(jumpOverReturnBreakOrContinueEntryOffset, location)

		c.compileStatementsWithResult(node.Finally, location)

		c.emit(location.EndPos.Line, bytecode.SWAP)
		jumpOverFinallyBreakOrContinueOffset := c.emitJump(location.EndPos.Line, bytecode.JUMP_UNLESS_UNP)
		c.emit(location.EndPos.Line, bytecode.POP_2)
		c.emit(location.EndPos.Line, bytecode.JUMP_TO_FINALLY)
		c.patchJump(jumpOverFinallyBreakOrContinueOffset, location)

		jumpToRethrowOffset := c.emitJump(location.EndPos.Line, bytecode.JUMP_IF_NP)
		jumpToFinallyReturnOffset := c.emitJump(location.EndPos.Line, bytecode.JUMP_IF_NIL_NP)
		// FALSE
		c.emit(location.EndPos.Line, bytecode.POP_2)          // pop the flag and return value of finally
		c.emit(location.EndPos.Line, bytecode.POP_2_SKIP_ONE) // pop the thrown value and the stack trace leaving the return value of catch
		jumpToEndOffset := c.emitJump(location.EndPos.Line, bytecode.JUMP)

		c.patchJump(jumpToFinallyReturnOffset, location)
		// return with finally
		c.emit(location.EndPos.Line, bytecode.POP_2) // pop the flag and return value of finally
		c.emit(location.EndPos.Line, bytecode.RETURN_FINALLY)

		c.patchJump(jumpToRethrowOffset, location)
		// pop the flag and the return value of finally
		c.emit(location.EndPos.Line, bytecode.POP_2)
		c.emit(location.EndPos.Line, bytecode.RETHROW)

		c.patchJump(jumpToEndOffset, location)
	}

	c.leaveScope(location.EndPos.Line)

	c.patchJump(jumpOverCatchOffset, location)
}

func (c *Compiler) compileMacroBoundaryNode(node *ast.MacroBoundaryNode) {
	location := node.Location()

	c.enterScope("", defaultScopeType)
	c.compileStatementsWithResult(node.Body, location)
	c.leaveScope(location.EndPos.Line)
}

func (c *Compiler) compileUnhygienicNode(node *ast.UnhygienicNode, valueIsIgnored bool) {
	prevUnhygienic := c.unhygienic
	c.unhygienic = true

	c.compileNode(node.Node, valueIsIgnored)

	c.unhygienic = prevUnhygienic
}

// Count `finally` blocks we are currently nested in under
// the nearest enclosing loop or
// under the loop with the specified label.
func (c *Compiler) countFinallyInLoop(label string) int {
	var finallyCount int
	for i := range c.scopes {
		scope := c.scopes[len(c.scopes)-i-1]
		if scope.typ == doFinallyScopeType {
			finallyCount++
		}
		if label == "" {
			if scope.typ == loopScopeType {
				break
			}
			continue
		}

		if scope.label == label {
			break
		}
	}

	return finallyCount
}

func (c *Compiler) leaveScopeOnBreak(line int, label string) {
	var varsToPop int
	for i := range c.scopes {
		scope := c.scopes[len(c.scopes)-i-1]
		varsToPop += len(scope.localTable)
		c.closeUpvaluesInScope(line, scope)

		if label == "" {
			if scope.typ == loopScopeType {
				break
			}
			continue
		}

		if scope.label == label {
			break
		}
	}
	c.emitLeaveScope(line, c.lastLocalIndex, varsToPop)
}

func (c *Compiler) compileBreakExpressionNode(node *ast.BreakExpressionNode) {
	location := node.Location()
	if node.Value == nil {
		c.emit(location.StartPos.Line, bytecode.NIL)
	} else {
		c.compileNode(node.Value, false)
	}

	labelName := identifierToName(node.Label)
	finallyCount := c.countFinallyInLoop(labelName)
	if finallyCount <= 0 {
		c.leaveScopeOnBreak(location.StartPos.Line, labelName)

		breakJumpOffset := c.emitJump(location.StartPos.Line, bytecode.JUMP)
		c.addLoopJump(labelName, breakLoopJump, breakJumpOffset, location)
		return
	}

	jumpOffsetId := c.emitLoadValue(value.Undefined, location)
	c.offsetValueIds = append(c.offsetValueIds, jumpOffsetId)
	c.addLoopJump(labelName, breakFinallyLoopJump, jumpOffsetId, location)

	c.emitValue(value.SmallInt(finallyCount).ToValue(), location)
	c.emit(location.StartPos.Line, bytecode.JUMP_TO_FINALLY)
}

func (c *Compiler) leaveScopeOnContinue(line int, label string) {
	var varsToPop int

	if label == "" {
		for i := range c.scopes {
			scope := c.scopes[len(c.scopes)-i-1]
			if scope.typ == loopScopeType {
				break
			}
			c.closeUpvaluesInScope(line, scope)
			varsToPop += len(scope.localTable)
		}
	} else {
		for i := range c.scopes {
			scope := c.scopes[len(c.scopes)-i-1]
			if scope.label == label {
				break
			}
			c.closeUpvaluesInScope(line, scope)
			varsToPop += len(scope.localTable)
		}
	}
	c.emitLeaveScope(line, c.lastLocalIndex, varsToPop)
}

func (c *Compiler) compileContinueExpressionNode(node *ast.ContinueExpressionNode) {
	location := node.Location()

	labelName := identifierToName(node.Label)
	loop := c.findLoopJumpSet(labelName, location)
	if loop == nil {
		return
	}

	if !loop.returnsValueFromLastIteration {
		if node.Value != nil {
			c.compileNode(node.Value, false)
			c.emit(location.StartPos.Line, bytecode.POP)
		}
	} else {
		if node.Value == nil {
			c.emit(location.StartPos.Line, bytecode.NIL)
		} else {
			c.compileNode(node.Value, false)
		}
	}

	finallyCount := c.countFinallyInLoop(labelName)
	if finallyCount <= 0 {
		c.leaveScopeOnContinue(location.StartPos.Line, labelName)

		continueJumpOffset := c.emitJump(location.StartPos.Line, bytecode.LOOP)
		c.addLoopJumpTo(loop, continueLoopJump, continueJumpOffset)
		return
	}

	jumpOffsetId := c.emitLoadValue(value.Undefined, location)
	c.offsetValueIds = append(c.offsetValueIds, jumpOffsetId)
	c.addLoopJump(labelName, continueFinallyLoopJump, jumpOffsetId, location)

	c.emitValue(value.SmallInt(finallyCount).ToValue(), location)
	c.emit(location.StartPos.Line, bytecode.JUMP_TO_FINALLY)
}

// Patch loop jump addresses for `break` and `continue` expressions.
func (c *Compiler) patchLoopJumps(continueOffset int) {
	lastLoopJumpSet := c.loopJumpSets[len(c.loopJumpSets)-1]
	for _, loopJump := range lastLoopJumpSet.loopJumps {
		switch loopJump.typ {
		case breakFinallyLoopJump:
			c.Bytecode.Values[loopJump.offset] = value.SmallInt(c.nextInstructionOffset()).ToValue()
		case continueFinallyLoopJump:
			c.Bytecode.Values[loopJump.offset] = value.SmallInt(continueOffset).ToValue()
		case breakLoopJump:
			c.patchJump(loopJump.offset, loopJump.location)
		case continueLoopJump:
			target := continueOffset - loopJump.offset
			if target >= 0 {
				// jump forward
				// override the opcode to JUMP
				c.Bytecode.Instructions[loopJump.offset-1] = byte(bytecode.JUMP)
				c.patchJumpWithTarget(target-2, loopJump.offset, loopJump.location)
			} else {
				// jump backward
				// override the opcode to LOOP
				c.Bytecode.Instructions[loopJump.offset-1] = byte(bytecode.LOOP)
				c.patchJumpWithTarget((-target)+2, loopJump.offset, loopJump.location)
			}
		default:
			panic(fmt.Sprintf("invalid loop jump info: %#v", loopJump))
		}
	}
	c.loopJumpSets = c.loopJumpSets[:len(c.loopJumpSets)-1]
}

func (c *Compiler) compileLoopExpressionNode(label string, body []ast.StatementNode, location *position.Location) {
	c.enterScope(label, loopScopeType)
	c.initLoopJumpSet(label, false)

	start := c.nextInstructionOffset()
	c.enterScope("", defaultScopeType)
	if c.compileStatementsOk(body) {
		c.emit(location.EndPos.Line, bytecode.POP)
	}
	c.leaveScope(location.EndPos.Line)
	c.emitLoop(location, start)

	c.leaveScope(location.EndPos.Line)
	c.patchLoopJumps(start)
}

func (c *Compiler) compileWhileExpressionNode(label string, node *ast.WhileExpressionNode) {
	location := node.Location()

	if result := resolve(node.Condition); !result.IsUndefined() {
		if value.Falsy(result) {
			// the loop won't run at all
			// it can be optimised into a simple NIL operation
			c.emit(location.StartPos.Line, bytecode.NIL)
			return
		}

		// the loop is endless
		c.compileLoopExpressionNode(label, node.ThenBody, location)
		return
	}

	c.enterScope(label, loopScopeType)
	c.initLoopJumpSet(label, true)

	c.emit(location.StartPos.Line, bytecode.NIL)
	// loop start
	start := c.nextInstructionOffset()
	var loopBodyOffset int

	if optimisedJumpOp, optimisedCond := c.optimiseCondition(bytecode.JUMP_UNLESS, node.Condition, location); optimisedCond != nil {
		optimisedCond()
		loopBodyOffset = c.emitJump(location.StartPos.Line, optimisedJumpOp)
	} else {
		// loop condition eg. `i < 5`
		c.compileNodeWithResult(node.Condition)
		// jump past the loop if the condition is falsy
		loopBodyOffset = c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS)
	}
	// pop the return value of the last iteration
	c.emit(location.StartPos.Line, bytecode.POP)

	// loop body
	c.compileStatementsWithResult(node.ThenBody, location)

	c.closeUpvaluesInCurrentScope(location.EndPos.Line)
	// jump to loop condition
	c.emitLoop(location, start)

	// after loop
	c.patchJump(loopBodyOffset, location)

	c.leaveScope(location.EndPos.Line)
	c.patchLoopJumps(start)
}

func (c *Compiler) modifierWhileExpression(label string, node *ast.ModifierNode) {
	location := node.Location()

	body := node.Left
	condition := node.Right

	var conditionIsStaticFalsy bool

	if result := resolve(condition); !result.IsUndefined() {
		if value.Truthy(result) {
			// the loop is endless
			c.compileLoopExpressionNode(label, ast.ExpressionToStatements(body), location)
			return
		}

		// the loop will only iterate once
		conditionIsStaticFalsy = true
	}

	c.enterScope(label, loopScopeType)
	c.initLoopJumpSet(label, true)

	// loop start
	start := c.nextInstructionOffset()
	c.enterScope("", defaultScopeType)
	var loopBodyOffset int

	// loop body
	c.compileNodeWithResult(body)
	// continue
	continueOffset := c.nextInstructionOffset()
	if conditionIsStaticFalsy {
		// the loop has a static falsy condition
		// it will only finish one iteration
		c.leaveScope(location.EndPos.Line)
		c.patchLoopJumps(continueOffset)
		return
	}

	if optimisedJumpOp, optimisedCond := c.optimiseCondition(bytecode.JUMP_UNLESS, condition, location); optimisedCond != nil {
		optimisedCond()
		loopBodyOffset = c.emitJump(location.StartPos.Line, optimisedJumpOp)
	} else {
		// loop condition eg. `i < 5`
		c.compileNodeWithResult(condition)
		// jump past the loop if the condition is falsy
		loopBodyOffset = c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS)
	}
	// pop the return value of the last iteration
	c.emit(location.StartPos.Line, bytecode.POP)

	c.leaveScope(location.EndPos.Line)
	// jump to loop start
	c.emitLoop(location, start)

	// after loop
	c.patchJump(loopBodyOffset, location)

	c.leaveScope(location.EndPos.Line)
	c.patchLoopJumps(continueOffset)
}

func (c *Compiler) modifierUntilExpression(label string, node *ast.ModifierNode) {
	location := node.Location()

	body := node.Left
	condition := node.Right

	var conditionIsStaticTruthy bool

	if result := resolve(condition); !result.IsUndefined() {
		if value.Falsy(result) {
			// the loop is endless
			c.compileLoopExpressionNode(label, ast.ExpressionToStatements(body), location)
			return
		}

		// the loop will only iterate once
		conditionIsStaticTruthy = true
	}

	c.enterScope(label, loopScopeType)
	c.initLoopJumpSet(label, true)

	// loop start
	start := c.nextInstructionOffset()
	c.enterScope("", defaultScopeType)
	var loopBodyOffset int

	// loop body
	c.compileNodeWithResult(body)
	// continue
	continueOffset := c.nextInstructionOffset()
	if conditionIsStaticTruthy {
		// the loop has a static truthy condition
		// it will only finish one iteration
		c.leaveScope(location.EndPos.Line)
		c.patchLoopJumps(continueOffset)
		return
	}

	if optimisedJumpOp, optimisedCond := c.optimiseCondition(bytecode.JUMP_IF, condition, location); optimisedCond != nil {
		optimisedCond()
		loopBodyOffset = c.emitJump(location.StartPos.Line, optimisedJumpOp)
	} else {
		// loop condition eg. `i > 5`
		c.compileNodeWithResult(condition)
		// jump past the loop if the condition is truthy
		loopBodyOffset = c.emitJump(location.StartPos.Line, bytecode.JUMP_IF)
	}
	// pop the return value of the last iteration
	c.emit(location.StartPos.Line, bytecode.POP)

	c.leaveScope(location.EndPos.Line)
	// jump to loop start
	c.emitLoop(location, start)

	// after loop
	c.patchJump(loopBodyOffset, location)

	c.leaveScope(location.EndPos.Line)
	c.patchLoopJumps(continueOffset)
}

func (c *Compiler) compileUntilExpressionNode(label string, node *ast.UntilExpressionNode) {
	location := node.Location()

	if result := resolve(node.Condition); !result.IsUndefined() {
		if value.Falsy(result) {
			// the loop is endless
			c.compileLoopExpressionNode(label, node.ThenBody, location)
			return
		}

		// the loop won't run at all
		// it can be optimised into a simple NIL operation
		c.emit(location.StartPos.Line, bytecode.NIL)
		return
	}

	c.enterScope(label, loopScopeType)
	c.initLoopJumpSet(label, true)

	c.emit(location.StartPos.Line, bytecode.NIL)
	// loop start
	start := c.nextInstructionOffset()
	c.enterScope("", defaultScopeType)
	var loopBodyOffset int

	if optimisedJumpOp, optimisedCond := c.optimiseCondition(bytecode.JUMP_IF, node.Condition, location); optimisedCond != nil {
		optimisedCond()
		loopBodyOffset = c.emitJump(location.StartPos.Line, optimisedJumpOp)
	} else {
		// loop condition eg. `i > 5`
		c.compileNodeWithResult(node.Condition)
		// jump past the loop if the condition is truthy
		loopBodyOffset = c.emitJump(location.StartPos.Line, bytecode.JUMP_IF)
	}
	// pop the return value of the last iteration
	c.emit(location.StartPos.Line, bytecode.POP)

	// loop body
	c.compileStatementsWithResult(node.ThenBody, location)

	c.leaveScope(location.EndPos.Line)
	// jump to loop condition
	c.emitLoop(location, start)

	// after loop
	c.patchJump(loopBodyOffset, location)

	c.leaveScope(location.EndPos.Line)
	c.patchLoopJumps(start)
}

// Compile a labeled expression eg. `$foo: println("bar")`
func (c *Compiler) compileLabeledExpressionNode(node *ast.LabeledExpressionNode, valueIsIgnored bool) expressionResult {
	switch expr := node.Expression.(type) {
	case *ast.WhileExpressionNode:
		c.compileWhileExpressionNode(node.Label, expr)
	case *ast.UntilExpressionNode:
		c.compileUntilExpressionNode(node.Label, expr)
	case *ast.LoopExpressionNode:
		c.compileLoopExpressionNode(node.Label, expr.ThenBody, expr.Location())
	case *ast.NumericForExpressionNode:
		c.compileNumericForExpressionNode(node.Label, expr)
	case *ast.ForInExpressionNode:
		c.compileForInExpressionNode(node.Label, expr)
	case *ast.ModifierForInNode:
		c.compileModifierForInNode(node.Label, expr)
	case *ast.ModifierNode:
		return c.compileModifierExpressionNode(node.Label, expr, valueIsIgnored)
	default:
		return c.compileNode(node.Expression, valueIsIgnored)
	}

	return expressionCompiled
}

// Compile a for in loop eg. `for i in [1, 2] then println(i)`
func (c *Compiler) compileForInExpressionNode(label string, node *ast.ForInExpressionNode) {
	c.compileForIn(
		label,
		node.Pattern,
		node.InExpression,
		func() {
			c.compileStatements(node.ThenBody, node.Location(), false)
		},
		node.Location(),
		false,
	)
}

// Compile a for in loop eg. `println(i) for i in [1, 2]`
func (c *Compiler) compileModifierForInNode(label string, node *ast.ModifierForInNode) {
	c.compileForIn(
		label,
		node.Pattern,
		node.InExpression,
		func() {
			result := c.compileNode(node.ThenExpression, false)
			switch result {
			case expressionIgnored, expressionCompiledWithoutResult:
				c.emit(node.ThenExpression.Location().StartPos.Line, bytecode.NIL)
			}
		},
		node.Location(),
		false,
	)
}

func (c *Compiler) compileForInAsNumericFor(
	label string,
	param ast.PatternNode,
	inExpression ast.ExpressionNode,
	then func(),
	location *position.Location,
	collectionLiteral bool,
) bool {
	var paramExpr ast.ExpressionNode
	var paramName string
	switch p := param.(type) {
	case *ast.PublicIdentifierNode:
		paramExpr = p
		paramName = p.Value
	case *ast.PrivateIdentifierNode:
		paramExpr = p
		paramName = p.Value
	default:
		return false
	}

	switch in := inExpression.(type) {
	case *ast.RangeLiteralNode:
		return c.compileForInRangeLiteralAsNumericFor(label, in, then, paramExpr, paramName, collectionLiteral, location)
	case *ast.IntLiteralNode:
		return c.compileForInIntLiteralAsNumericFor(label, in, then, paramExpr, paramName, collectionLiteral, location)
	}

	inExpressionType := c.typeOf(inExpression)
	if c.checker.IsSubtype(inExpressionType, c.checker.Std(symbol.Range)) {
		return c.compileForInRangeAsNumericFor(label, inExpression, then, paramExpr, paramName, collectionLiteral, location)
	}
	if c.checker.IsSubtype(inExpressionType, c.checker.Std(symbol.Int)) {
		return c.compileForInIntAsNumericFor(label, inExpression, then, paramExpr, paramName, collectionLiteral, location)
	}

	return false
}

func (c *Compiler) compileForInIntAsNumericFor(label string, inExpression ast.ExpressionNode, then func(), paramExpr ast.ExpressionNode, paramName string, collectionLiteral bool, location *position.Location) bool {
	c.enterScope("", defaultScopeType)

	intVarName := fmt.Sprintf("#!forInt%d", len(c.scopes))

	c.compileNodeWithResult(inExpression)
	intVar := c.defineLocal(intVarName, location)

	c.emitSetLocalPop(location.StartPos.Line, intVar.index)

	init := ast.NewVariableDeclarationNode(
		paramExpr.Location(),
		"",
		ast.NewPublicIdentifierNode(paramExpr.Location(), paramName),
		nil,
		ast.NewIntLiteralNode(location, "0"),
	)
	increment := ast.NewPostfixExpressionNode(
		location,
		token.New(location, token.PLUS_PLUS),
		paramExpr,
	)

	cond := ast.NewBinaryExpressionNode(
		location,
		token.New(location, token.LESS),
		paramExpr,
		ast.NewPublicIdentifierNode(inExpression.Location(), intVarName),
	)

	c.compileNumericFor(
		label,
		init,
		cond,
		increment,
		then,
		collectionLiteral,
		location,
	)
	if !collectionLiteral {
		c.emit(location.EndPos.Line, bytecode.POP)
		c.emit(location.EndPos.Line, bytecode.NIL)
	}

	c.leaveScope(location.EndPos.Line)
	return true
}

func (c *Compiler) compileForInRangeAsNumericFor(label string, inExpression ast.ExpressionNode, then func(), paramExpr ast.ExpressionNode, paramName string, collectionLiteral bool, location *position.Location) bool {
	inExpressionType := c.typeOf(inExpression).(*types.Generic)

	if c.checker.IsSubtype(inExpressionType, c.checker.Std(symbol.BeginlessClosedRange)) ||
		c.checker.IsSubtype(inExpressionType, c.checker.Std(symbol.BeginlessOpenRange)) {
		return false
	}

	c.enterScope("", defaultScopeType)

	rangeVarName := fmt.Sprintf("#!forRange%d", len(c.scopes))
	rangeEndVarName := fmt.Sprintf("#!forRangeEnd%d", len(c.scopes))

	initVal := ast.NewMethodCallNode(
		location,
		ast.NewPublicIdentifierNode(inExpression.Location(), rangeVarName),
		token.New(location, token.DOT),
		ast.NewPublicIdentifierNode(inExpression.Location(), "start"),
		nil,
		nil,
	)
	rangeElementType := inExpressionType.Get(0).Type
	initVal.SetType(rangeElementType)

	var cmpOp token.Type
	if c.checker.IsSubtype(inExpressionType, c.checker.Std(symbol.ClosedRange)) {
		cmpOp = token.LESS_EQUAL
	} else if c.checker.IsSubtype(inExpressionType, c.checker.Std(symbol.EndlessClosedRange)) {
	} else if c.checker.IsSubtype(inExpressionType, c.checker.Std(symbol.RightOpenRange)) {
		cmpOp = token.LESS
	} else if c.checker.IsSubtype(inExpressionType, c.checker.Std(symbol.OpenRange)) {
		cmpOp = token.LESS
		initVal = ast.NewMethodCallNode(
			location,
			initVal,
			token.New(location, token.DOT),
			ast.NewPublicIdentifierNode(location, "++"),
			nil,
			nil,
		)
	} else if c.checker.IsSubtype(inExpressionType, c.checker.Std(symbol.LeftOpenRange)) {
		cmpOp = token.LESS_EQUAL
		initVal = ast.NewMethodCallNode(
			location,
			initVal,
			token.New(location, token.DOT),
			ast.NewPublicIdentifierNode(location, "++"),
			nil,
			nil,
		)
	} else if c.checker.IsSubtype(inExpressionType, c.checker.Std(symbol.EndlessOpenRange)) {
		initVal = ast.NewMethodCallNode(
			location,
			initVal,
			token.New(location, token.DOT),
			ast.NewPublicIdentifierNode(location, "++"),
			nil,
			nil,
		)
	}

	c.compileNodeWithResult(inExpression)
	rangeVar := c.defineLocal(rangeVarName, location)
	if cmpOp != token.ZERO_VALUE {
		c.emitSetLocalNoPop(location.StartPos.Line, rangeVar.index)

		c.emitCallMethod(value.NewCallSiteInfo(value.ToSymbol("end"), 0), location, false)
		rangeEndVar := c.defineLocal(rangeEndVarName, location)
		c.emitSetLocalPop(location.StartPos.Line, rangeEndVar.index)
	} else {
		c.emitSetLocalPop(location.StartPos.Line, rangeVar.index)
	}

	init := ast.NewVariableDeclarationNode(
		paramExpr.Location(),
		"",
		ast.NewPublicIdentifierNode(paramExpr.Location(), paramName),
		nil,
		initVal,
	)
	increment := ast.NewPostfixExpressionNode(
		location,
		token.New(location, token.PLUS_PLUS),
		paramExpr,
	)
	var cond ast.ExpressionNode
	if cmpOp != token.ZERO_VALUE {
		cond = ast.NewBinaryExpressionNode(
			location,
			token.New(location, cmpOp),
			paramExpr,
			ast.NewPublicIdentifierNode(inExpression.Location(), rangeEndVarName),
		)
	}
	c.compileNumericFor(
		label,
		init,
		cond,
		increment,
		then,
		collectionLiteral,
		location,
	)
	if !collectionLiteral {
		c.emit(location.EndPos.Line, bytecode.POP)
		c.emit(location.EndPos.Line, bytecode.NIL)
	}

	c.leaveScope(location.EndPos.Line)
	return true
}

func (c *Compiler) compileForInIntLiteralAsNumericFor(label string, inInt *ast.IntLiteralNode, then func(), paramExpr ast.ExpressionNode, paramName string, collectionLiteral bool, location *position.Location) bool {
	init := ast.NewVariableDeclarationNode(
		paramExpr.Location(),
		"",
		ast.NewPublicIdentifierNode(paramExpr.Location(), paramName),
		nil,
		ast.NewIntLiteralNode(paramExpr.Location(), "0"),
	)
	increment := ast.NewPostfixExpressionNode(
		paramExpr.Location(),
		token.New(paramExpr.Location(), token.PLUS_PLUS),
		paramExpr,
	)

	cond := ast.NewBinaryExpressionNode(
		inInt.Location(),
		token.New(inInt.Location(), token.LESS),
		paramExpr,
		inInt,
	)
	c.compileNumericFor(
		label,
		init,
		cond,
		increment,
		then,
		collectionLiteral,
		location,
	)
	if !collectionLiteral {
		c.emit(location.EndPos.Line, bytecode.POP)
		c.emit(location.EndPos.Line, bytecode.NIL)
	}

	return true
}

func (c *Compiler) compileForInRangeLiteralAsNumericFor(label string, inRange *ast.RangeLiteralNode, then func(), paramExpr ast.ExpressionNode, paramName string, collectionLiteral bool, location *position.Location) bool {
	if inRange.Start == nil {
		return false
	}

	var op token.Type
	var initVal ast.ExpressionNode

	switch inRange.Op.Type {
	case token.CLOSED_RANGE_OP:
		op = token.LESS_EQUAL
		initVal = inRange.Start
	case token.RIGHT_OPEN_RANGE_OP:
		op = token.LESS
		initVal = inRange.Start
	case token.LEFT_OPEN_RANGE_OP:
		op = token.LESS_EQUAL
		initVal = ast.NewMethodCallNode(
			inRange.Location(),
			inRange.Start,
			token.New(inRange.Op.Location(), token.DOT),
			ast.NewPublicIdentifierNode(location, "++"),
			nil,
			nil,
		)
	case token.OPEN_RANGE_OP:
		op = token.LESS
		initVal = ast.NewMethodCallNode(
			inRange.Location(),
			inRange.Start,
			token.New(inRange.Op.Location(), token.DOT),
			ast.NewPublicIdentifierNode(location, "++"),
			nil,
			nil,
		)
	default:
		return false
	}

	init := ast.NewVariableDeclarationNode(
		paramExpr.Location(),
		"",
		ast.NewPublicIdentifierNode(paramExpr.Location(), paramName),
		nil,
		initVal,
	)
	increment := ast.NewPostfixExpressionNode(
		inRange.Op.Location(),
		token.New(inRange.Op.Location(), token.PLUS_PLUS),
		paramExpr,
	)

	var cond ast.ExpressionNode
	if inRange.End != nil {
		cond = ast.NewBinaryExpressionNode(
			inRange.End.Location(),
			token.New(inRange.Op.Location(), op),
			paramExpr,
			inRange.End,
		)
	}
	c.compileNumericFor(
		label,
		init,
		cond,
		increment,
		then,
		collectionLiteral,
		location,
	)
	if !collectionLiteral {
		c.emit(location.EndPos.Line, bytecode.POP)
		c.emit(location.EndPos.Line, bytecode.NIL)
	}

	return true
}

func (c *Compiler) compileForIn(
	label string,
	param ast.PatternNode,
	inExpression ast.ExpressionNode,
	then func(),
	location *position.Location,
	collectionLiteral bool,
) {
	if c.compileForInAsNumericFor(label, param, inExpression, then, location, collectionLiteral) {
		return
	}

	c.enterScope(label, loopScopeType)
	c.initLoopJumpSet(label, false)

	c.compileNode(inExpression, false)
	inExpressionType := c.typeOf(inExpression)
	if c.checker.IsSubtype(inExpressionType, c.checker.Std(symbol.S_BuiltinIterable)) {
		c.emit(location.StartPos.Line, bytecode.GET_ITERATOR)
	} else {
		c.emitCallMethod(value.NewCallSiteInfo(symbol.L_iter, 0), inExpression.Location(), false)
	}

	iteratorVarName := fmt.Sprintf("#!forIn%d", len(c.scopes))
	iteratorVar := c.defineLocal(iteratorVarName, location)
	c.emitSetLocalPop(location.StartPos.Line, iteratorVar.index)

	// loop start
	start := c.nextInstructionOffset()
	continueOffset := start

	c.emitGetLocal(location.StartPos.Line, iteratorVar.index)

	var loopBodyOffset int

	nextType := c.checker.GetIteratorType(inExpressionType)
	if c.checker.IsSubtype(nextType, c.checker.Std(symbol.S_BuiltinIterator)) {
		loopBodyOffset = c.emitJump(location.StartPos.Line, bytecode.FOR_IN_BUILTIN)
	} else {
		c.emitCallNext(value.NewCallSiteInfo(symbol.L_next, 0), inExpression.Location())
		loopBodyOffset = c.emitJump(location.StartPos.Line, bytecode.FOR_IN)
	}

	switch p := param.(type) {
	case *ast.PrivateIdentifierNode:
		paramVar := c.defineLocal(p.Value, param.Location())
		c.emitSetLocalPop(param.Location().StartPos.Line, paramVar.index)
	case *ast.PublicIdentifierNode:
		paramVar := c.defineLocal(p.Value, param.Location())
		c.emitSetLocalPop(param.Location().StartPos.Line, paramVar.index)
	default:
		c.pattern(param)
		jumpOverErrorOffset := c.emitJump(location.StartPos.Line, bytecode.JUMP_IF)

		c.emitValue(
			value.Ref(value.NewError(
				value.PatternNotMatchedErrorClass,
				"assigned value does not match the pattern defined in for in loop",
			)),
			location,
		)
		c.emit(location.EndPos.Line, bytecode.THROW)

		c.patchJump(jumpOverErrorOffset, location)
		c.emit(param.Location().EndPos.Line, bytecode.POP)
	}

	// loop body
	then()

	// pop the return value of the block
	if !collectionLiteral {
		c.emit(location.EndPos.Line, bytecode.POP)
	}
	// jump to loop condition
	c.emitLoop(location, start)

	// after loop
	c.patchJump(loopBodyOffset, location)
	if !collectionLiteral {
		c.emit(location.EndPos.Line, bytecode.NIL)
	}

	c.leaveScope(location.EndPos.Line)
	c.patchLoopJumps(continueOffset)
}

// Compile a numeric for loop eg. `for i := 0; i < 5; i += 1 then println(i)`
func (c *Compiler) compileNumericForExpressionNode(label string, node *ast.NumericForExpressionNode) {
	location := node.Location()

	if node.Initialiser == nil && node.Condition == nil && node.Increment == nil {
		// the loop is endless
		c.compileLoopExpressionNode(label, node.ThenBody, location)
		return
	}

	c.compileNumericFor(
		label,
		node.Initialiser,
		node.Condition,
		node.Increment,
		func() {
			c.compileStatementsWithResult(node.ThenBody, location)
		},
		false,
		location,
	)
}

func (c *Compiler) compileNumericFor(label string, init, cond, increment ast.ExpressionNode, then func(), collectionLiteral bool, location *position.Location) {
	c.enterScope(label, loopScopeType)
	c.initLoopJumpSet(label, true)

	// loop initialiser eg. `i := 0`
	if init != nil {
		c.compileNodeWithoutResult(init)
	}

	if !collectionLiteral {
		c.emit(location.StartPos.Line, bytecode.NIL)
	}
	// loop start
	start := c.nextInstructionOffset()
	continueOffset := start

	var loopBodyOffset int
	// loop condition eg. `i < 5`
	if cond != nil {
		// jump past the loop if the condition is falsy
		if optimisedJumpOp, optimisedCond := c.optimiseCondition(bytecode.JUMP_UNLESS, cond, location); optimisedCond != nil {
			optimisedCond()
			loopBodyOffset = c.emitJump(location.StartPos.Line, optimisedJumpOp)
		} else {
			c.compileNodeWithResult(cond)
			loopBodyOffset = c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS)
		}
	}
	if !collectionLiteral {
		// pop the return value of the last iteration
		c.emit(location.EndPos.Line, bytecode.POP)
	}

	// loop body
	then()

	c.closeUpvaluesInCurrentScope(location.EndPos.Line)
	if increment != nil {
		continueOffset = c.nextInstructionOffset()
		// increment step eg. `i += 1`
		c.compileNodeWithoutResult(increment)
	}

	// jump to loop condition
	c.emitLoop(location, start)

	// after loop
	if cond != nil {
		c.patchJump(loopBodyOffset, location)
	}

	c.leaveScope(location.EndPos.Line)
	c.patchLoopJumps(continueOffset)
}

func (c *Compiler) emitSetterCall(name string, location *position.Location) {
	nameSymbol := value.ToSymbol(name + "=")
	callInfo := value.NewCallSiteInfo(nameSymbol, 1)
	c.emitCallMethod(callInfo, location, false)
}

func (c *Compiler) emitGetterCall(name string, location *position.Location) {
	nameSymbol := value.ToSymbol(name)
	callInfo := value.NewCallSiteInfo(nameSymbol, 0)
	c.emitCallMethod(callInfo, location, false)
}

func (c *Compiler) compileIncrement(typ types.Type, location *position.Location) {
	if c.checker.IsSubtype(typ, c.checker.StdInt()) {
		c.emit(location.EndPos.Line, bytecode.INCREMENT_INT)
		return
	}
	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinIncrementable)) {
		c.emit(location.EndPos.Line, bytecode.INCREMENT)
		return
	}

	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpIncrement, 0), location, false)
}

func (c *Compiler) compileDecrement(typ types.Type, location *position.Location) {
	if c.checker.IsSubtype(typ, c.checker.StdInt()) {
		c.emit(location.EndPos.Line, bytecode.DECREMENT_INT)
		return
	}
	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinIncrementable)) {
		c.emit(location.EndPos.Line, bytecode.DECREMENT)
		return
	}

	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpDecrement, 0), location, false)
}

func (c *Compiler) compileSubscript(typ types.Type, location *position.Location) {
	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinSubscriptable)) {
		c.emit(location.EndPos.Line, bytecode.SUBSCRIPT)
		return
	}

	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpSubscript, 1), location, false)
}

func (c *Compiler) compileSubscriptSet(typ types.Type, location *position.Location) {
	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinSubscriptable)) {
		c.emit(location.EndPos.Line, bytecode.SUBSCRIPT_SET)
		return
	}

	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpSubscriptSet, 2), location, false)
}

func (c *Compiler) compilePostfixExpressionNode(node *ast.PostfixExpressionNode, valueIsIgnored bool) expressionResult {
	switch n := node.Expression.(type) {
	case *ast.PublicIdentifierNode:
		// get variable value
		c.compileLocalVariableAccess(n.Value, n.Location())

		switch node.Op.Type {
		case token.PLUS_PLUS:
			c.compileIncrement(c.typeOf(n), node.Location())
		case token.MINUS_MINUS:
			c.compileDecrement(c.typeOf(n), node.Location())
		default:
			panic(fmt.Sprintf("invalid postfix operator: %#v", node.Op))
		}

		// set variable
		return c.setLocalWithoutValue(n.Value, n.Location(), valueIsIgnored)
	case *ast.PrivateIdentifierNode:
		// get variable value
		c.compileLocalVariableAccess(n.Value, n.Location())

		switch node.Op.Type {
		case token.PLUS_PLUS:
			c.compileIncrement(c.typeOf(n), node.Location())
		case token.MINUS_MINUS:
			c.compileDecrement(c.typeOf(n), node.Location())
		default:
			panic(fmt.Sprintf("invalid postfix operator: %#v", node.Op))
		}

		// set variable
		return c.setLocalWithoutValue(n.Value, n.Location(), valueIsIgnored)
	case *ast.SubscriptExpressionNode:
		// get value
		c.compileNodeWithResult(n.Receiver)
		c.compileNodeWithResult(n.Key)
		c.emit(node.Location().EndPos.Line, bytecode.DUP_2)

		receiverType := c.typeOf(n.Receiver)
		c.compileSubscript(receiverType, node.Location())

		switch node.Op.Type {
		case token.PLUS_PLUS:
			c.compileIncrement(c.typeOf(n), node.Location())
		case token.MINUS_MINUS:
			c.compileDecrement(c.typeOf(n), node.Location())
		default:
			panic(fmt.Sprintf("invalid postfix operator: %#v", node.Op))
		}

		// set value
		c.compileSubscriptSet(receiverType, node.Location())
	case *ast.PublicInstanceVariableNode:
		switch c.mode {
		case topLevelMode:
			c.Errors.AddFailure(
				"instance variables cannot be set in the top level",
				node.Location(),
			)
		}
		// get value
		ivarSymbol := value.ToSymbol(n.Value)
		c.emitGetInstanceVariableByName(ivarSymbol, n.Location())

		switch node.Op.Type {
		case token.PLUS_PLUS:
			c.compileIncrement(c.typeOf(n), node.Location())
		case token.MINUS_MINUS:
			c.compileDecrement(c.typeOf(n), node.Location())
		default:
			panic(fmt.Sprintf("invalid postfix operator: %#v", node.Op))
		}

		// set instance variable
		c.emitSetInstanceVariable(ivarSymbol, node.Location(), valueIsIgnored)
		return valueIgnoredToResult(valueIsIgnored)
	case *ast.AttributeAccessNode:
		// get value
		c.compileNodeWithResult(n.Receiver)
		name := identifierToName(n.AttributeName)
		nameSymbol := value.ToSymbol(name)
		callInfo := value.NewCallSiteInfo(nameSymbol, 0)
		c.emitCallMethod(callInfo, node.Location(), false)

		switch node.Op.Type {
		case token.PLUS_PLUS:
			c.compileIncrement(c.typeOf(n), node.Location())
		case token.MINUS_MINUS:
			c.compileDecrement(c.typeOf(n), node.Location())
		default:
			panic(fmt.Sprintf("invalid postfix operator: %#v", node.Op))
		}

		// set attribute
		c.emitSetterCall(name, node.Location())
	default:
		c.Errors.AddFailure(
			fmt.Sprintf("cannot assign to: %T", node.Expression),
			node.Location(),
		)
	}

	return expressionCompiled
}

func (c *Compiler) attributeAssignment(node *ast.AssignmentExpressionNode, attr *ast.AttributeAccessNode) {
	// compile the argument
	switch node.Op.Type {
	case token.EQUAL_OP:
		c.compileNodeWithResult(attr.Receiver)
		c.compileNodeWithResult(node.Right)
		c.emitSetterCall(identifierToName(attr.AttributeName), node.Location())
	case token.OR_OR_EQUAL:
		name := identifierToName(attr.AttributeName)
		location := node.Location()
		// Read the current value
		c.compileNodeWithResult(attr.Receiver)
		c.emitGetterCall(name, location)

		jump := c.emitJump(location.StartPos.Line, bytecode.JUMP_IF_NP)

		// if falsy
		c.emit(location.StartPos.Line, bytecode.POP)
		c.compileNodeWithResult(node.Right)
		c.emitSetterCall(name, location)

		// if truthy
		c.patchJump(jump, location)
	case token.AND_AND_EQUAL:
		name := identifierToName(attr.AttributeName)
		location := node.Location()
		// Read the current value
		c.compileNodeWithResult(attr.Receiver)
		c.emitGetterCall(name, location)

		jump := c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NP)

		// if truthy
		c.emit(location.StartPos.Line, bytecode.POP)
		c.compileNodeWithResult(node.Right)
		c.emitSetterCall(name, location)

		// if falsy
		c.patchJump(jump, location)
	case token.QUESTION_QUESTION_EQUAL:
		name := identifierToName(attr.AttributeName)
		location := node.Location()
		// Read the current value
		c.compileNodeWithResult(attr.Receiver)
		c.emitGetterCall(name, location)

		nilJump := c.emitJump(location.StartPos.Line, bytecode.JUMP_IF_NIL_NP)
		nonNilJump := c.emitJump(location.StartPos.Line, bytecode.JUMP)

		// if nil
		c.patchJump(nilJump, location)
		c.emit(location.StartPos.Line, bytecode.POP)
		c.compileNodeWithResult(node.Right)
		c.emitSetterCall(name, location)

		// if not nil
		c.patchJump(nonNilJump, location)
	default:
		c.Errors.AddFailure(fmt.Sprintf("unknown binary operator: %s", node.Op.String()), node.Location())
	}
}

func (c *Compiler) instanceVariableAssignment(node *ast.AssignmentExpressionNode, ivar *ast.PublicInstanceVariableNode, valueIsIgnored bool) expressionResult {
	switch c.mode {
	case topLevelMode:
		c.Errors.AddFailure(
			"instance variables cannot be set in the top level",
			node.Location(),
		)
	}

	ivarSymbol := value.ToSymbol(ivar.Value)
	switch node.Op.Type {
	case token.EQUAL_OP:
		c.compileNodeWithResult(node.Right)
		c.emitSetInstanceVariable(ivarSymbol, ivar.Location(), valueIsIgnored)
		return valueIgnoredToResult(valueIsIgnored)
	case token.OR_OR_EQUAL:
		location := node.Location()
		// Read the current value
		c.emitGetInstanceVariable(ivarSymbol, location)

		var jump int
		if valueIsIgnored {
			jump = c.emitJump(location.StartPos.Line, bytecode.JUMP_IF)
		} else {
			jump = c.emitJump(location.StartPos.Line, bytecode.JUMP_IF_NP)
		}

		// if falsy
		if !valueIsIgnored {
			c.emit(location.StartPos.Line, bytecode.POP)
		}
		c.compileNodeWithResult(node.Right)
		c.emitSetInstanceVariable(ivarSymbol, ivar.Location(), valueIsIgnored)

		// if truthy
		c.patchJump(jump, location)
		return valueIgnoredToResult(valueIsIgnored)
	case token.AND_AND_EQUAL:
		location := node.Location()
		// Read the current value
		c.emitGetInstanceVariable(ivarSymbol, location)

		var jump int
		if valueIsIgnored {
			jump = c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS)
		} else {
			jump = c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NP)
		}

		// if truthy
		if !valueIsIgnored {
			c.emit(location.StartPos.Line, bytecode.POP)
		}
		c.compileNodeWithResult(node.Right)
		c.emitSetInstanceVariable(ivarSymbol, ivar.Location(), valueIsIgnored)

		// if falsy
		c.patchJump(jump, location)
		return valueIgnoredToResult(valueIsIgnored)
	case token.QUESTION_QUESTION_EQUAL:
		location := node.Location()
		// Read the current value
		c.emitGetInstanceVariable(ivarSymbol, location)

		var jump int
		if valueIsIgnored {
			jump = c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NIL)
		} else {
			jump = c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NNP)
		}

		// if nil
		if !valueIsIgnored {
			c.emit(location.StartPos.Line, bytecode.POP)
		}
		c.compileNodeWithResult(node.Right)
		c.emitSetInstanceVariable(ivarSymbol, ivar.Location(), valueIsIgnored)

		// if not nil
		c.patchJump(jump, location)
		return valueIgnoredToResult(valueIsIgnored)
	default:
		c.Errors.AddFailure(fmt.Sprintf("unknown binary operator: %s", node.Op.String()), node.Location())
	}

	return expressionCompiled
}

func (c *Compiler) subscriptAssignment(node *ast.AssignmentExpressionNode, subscript *ast.SubscriptExpressionNode, valueIsIgnored bool) expressionResult {
	receiverType := c.typeOf(subscript.Receiver)

	switch node.Op.Type {
	case token.EQUAL_OP:
		c.compileNodeWithResult(subscript.Receiver)
		c.compileNodeWithResult(subscript.Key)
		c.compileNodeWithResult(node.Right)

		c.compileSubscriptSet(receiverType, node.Location())
	case token.OR_OR_EQUAL:
		location := node.Location()
		// Read the current value
		c.compileNodeWithResult(subscript.Receiver)
		c.compileNodeWithResult(subscript.Key)
		c.emit(node.Location().EndPos.Line, bytecode.DUP_2)
		c.compileSubscript(receiverType, node.Location())

		var jump int
		if valueIsIgnored {
			jump = c.emitJump(location.StartPos.Line, bytecode.JUMP_IF)
		} else {
			jump = c.emitJump(location.StartPos.Line, bytecode.JUMP_IF_NP)
		}

		// if falsy
		if !valueIsIgnored {
			c.emit(location.StartPos.Line, bytecode.POP)
		}
		c.compileNodeWithResult(node.Right)
		c.compileSubscriptSet(receiverType, node.Location())

		if valueIsIgnored {
			c.emit(location.StartPos.Line, bytecode.POP)
		}

		// if truthy
		c.patchJump(jump, location)
		if valueIsIgnored {
			c.emit(location.StartPos.Line, bytecode.POP_2)
		} else {
			c.emit(location.StartPos.Line, bytecode.POP_2_SKIP_ONE)
		}
		return valueIgnoredToResult(valueIsIgnored)
	case token.AND_AND_EQUAL:
		location := node.Location()
		// Read the current value
		c.compileNodeWithResult(subscript.Receiver)
		c.compileNodeWithResult(subscript.Key)
		c.emit(node.Location().EndPos.Line, bytecode.DUP_2)
		c.compileSubscript(receiverType, node.Location())

		var jump int
		if valueIsIgnored {
			jump = c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS)
		} else {
			jump = c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NP)
		}

		// if truthy
		if !valueIsIgnored {
			c.emit(location.StartPos.Line, bytecode.POP)
		}
		c.compileNodeWithResult(node.Right)
		c.compileSubscriptSet(receiverType, node.Location())

		// if falsy
		c.patchJump(jump, location)
		if valueIsIgnored {
			c.emit(location.StartPos.Line, bytecode.POP_2)
		} else {
			c.emit(location.StartPos.Line, bytecode.POP_2_SKIP_ONE)
		}
		return valueIgnoredToResult(valueIsIgnored)
	case token.QUESTION_QUESTION_EQUAL:
		location := node.Location()
		// Read the current value
		c.compileNodeWithResult(subscript.Receiver)
		c.compileNodeWithResult(subscript.Key)
		c.emit(node.Location().EndPos.Line, bytecode.DUP_2)
		c.compileSubscript(receiverType, node.Location())

		var jump int
		if valueIsIgnored {
			jump = c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NIL)
		} else {
			jump = c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NNP)
		}

		// if nil
		if !valueIsIgnored {
			c.emit(location.StartPos.Line, bytecode.POP)
		}
		c.compileNodeWithResult(node.Right)
		c.compileSubscriptSet(receiverType, node.Location())

		// if not nil
		c.patchJump(jump, location)
		if valueIsIgnored {
			c.emit(location.StartPos.Line, bytecode.POP_2)
		} else {
			c.emit(location.StartPos.Line, bytecode.POP_2_SKIP_ONE)
		}
		return valueIgnoredToResult(valueIsIgnored)
	default:
		c.Errors.AddFailure(fmt.Sprintf("unknown binary operator: %s", node.Op.String()), node.Location())
	}

	return expressionCompiled
}

func (c *Compiler) compileAssignmentExpressionNode(node *ast.AssignmentExpressionNode, valueIsIgnored bool) expressionResult {
	switch n := node.Left.(type) {
	case *ast.PublicIdentifierNode:
		return c.localVariableAssignment(n.Value, node.Op, node.Right, node.Location(), valueIsIgnored)
	case *ast.PrivateIdentifierNode:
		return c.localVariableAssignment(n.Value, node.Op, node.Right, node.Location(), valueIsIgnored)
	case *ast.SubscriptExpressionNode:
		return c.subscriptAssignment(node, n, valueIsIgnored)
	case *ast.PublicInstanceVariableNode:
		return c.instanceVariableAssignment(node, n, valueIsIgnored)
	case *ast.AttributeAccessNode:
		c.attributeAssignment(node, n)
	default:
		c.Errors.AddFailure(
			fmt.Sprintf("cannot assign to: %T", node.Left),
			node.Location(),
		)
	}

	return expressionCompiled
}

// Return the offset of the next instruction.
func (c *Compiler) nextInstructionOffset() int {
	return len(c.Bytecode.Instructions)
}

func (c *Compiler) setLocalWithoutValue(name string, location *position.Location, valueIsIgnored bool) expressionResult {
	if local, ok := c.resolveLocal(name); ok {
		return c.emitSetLocal(location.StartPos.Line, local.index, valueIsIgnored)
	} else if upvalue, ok := c.resolveUpvalue(name, location); ok {
		return c.emitSetUpvalue(location.StartPos.Line, upvalue.index, valueIsIgnored)
	}

	return valueIgnoredToResult(valueIsIgnored)
}

func (c *Compiler) setLocal(name string, valueNode ast.ExpressionNode, location *position.Location, valueIsIgnored bool) expressionResult {
	c.compileNodeWithResult(valueNode)
	return c.setLocalWithoutValue(name, location, valueIsIgnored)
}

func (c *Compiler) localVariableAssignment(name string, operator *token.Token, right ast.ExpressionNode, location *position.Location, valueIsIgnored bool) expressionResult {
	switch operator.Type {
	case token.OR_OR_EQUAL:
		c.compileLocalVariableAccess(name, location)
		var jump int
		if valueIsIgnored {
			jump = c.emitJump(location.StartPos.Line, bytecode.JUMP_IF)
		} else {
			jump = c.emitJump(location.StartPos.Line, bytecode.JUMP_IF_NP)
		}

		// if falsy
		if !valueIsIgnored {
			c.emit(location.StartPos.Line, bytecode.POP)
		}
		c.setLocal(name, right, location, valueIsIgnored)

		// if truthy
		c.patchJump(jump, location)
		return valueIgnoredToResult(valueIsIgnored)
	case token.AND_AND_EQUAL:
		c.compileLocalVariableAccess(name, location)
		var jump int
		if valueIsIgnored {
			jump = c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS)
		} else {
			jump = c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NP)
		}

		// if truthy
		if !valueIsIgnored {
			c.emit(location.StartPos.Line, bytecode.POP)
		}
		c.setLocal(name, right, location, valueIsIgnored)

		// if falsy
		c.patchJump(jump, location)
		return valueIgnoredToResult(valueIsIgnored)
	case token.QUESTION_QUESTION_EQUAL:
		c.compileLocalVariableAccess(name, location)
		var jump int
		if valueIsIgnored {
			jump = c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NIL)
		} else {
			jump = c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NNP)
		}

		// if nil
		if !valueIsIgnored {
			c.emit(location.StartPos.Line, bytecode.POP)
		}
		c.setLocal(name, right, location, valueIsIgnored)

		// if not nil
		c.patchJump(jump, location)
		return valueIgnoredToResult(valueIsIgnored)
	case token.EQUAL_OP:
		return c.setLocal(name, right, location, valueIsIgnored)
	case token.COLON_EQUAL:
		local := c.defineLocal(name, location)
		c.compileNodeWithResult(right)
		if local == nil {
			return valueIgnoredToResult(valueIsIgnored)
		}
		return c.emitSetLocal(location.StartPos.Line, local.index, valueIsIgnored)
	default:
		c.Errors.AddFailure(
			fmt.Sprintf("assignment using this operator has not been implemented: %s", operator.Type.Name()),
			location,
		)
	}

	return expressionCompiled
}

func (c *Compiler) compileInstanceVariableAccess(name string, location *position.Location) {
	c.emitGetInstanceVariable(value.ToSymbol(name), location)
}

func (c *Compiler) compileLocalVariableAccess(name string, location *position.Location) (*local, *upvalue, bool) {
	if local, ok := c.resolveLocal(name); ok {
		c.emitGetLocal(location.StartPos.Line, local.index)
		return local, nil, true
	} else if upvalue, ok := c.resolveUpvalue(name, location); ok {
		local := upvalue.local
		c.emitGetUpvalue(location.StartPos.Line, upvalue.index)
		return local, upvalue, true
	}

	return nil, nil, false
}

// Resolve an upvalue from an outer context and get its index.
func (c *Compiler) resolveUpvalue(name string, location *position.Location) (*upvalue, bool) {
	parent := c.Parent
	if parent == nil {
		return nil, false
	}
	local, ok := parent.resolveLocal(name)
	if ok {
		return c.addUpvalue(local, local.index, true, location), true
	}

	upvalue, ok := parent.resolveUpvalue(name, location)
	if ok {
		return c.addUpvalue(upvalue.local, upvalue.index, false, location), true
	}

	return nil, false
}

func (c *Compiler) addUpvalue(local *local, upIndex uint16, isLocal bool, location *position.Location) *upvalue {
	for _, upvalue := range c.upvalues {
		if upvalue.upIndex == upIndex && upvalue.isLocal == isLocal {
			return upvalue
		}
	}

	if len(c.upvalues) > math.MaxUint16 {
		c.Errors.AddFailure(
			fmt.Sprintf("upvalue limit reached: %d", math.MaxUint16),
			location,
		)
	}

	upvalue := &upvalue{
		index:   uint16(len(c.upvalues)),
		upIndex: upIndex,
		local:   local,
		isLocal: isLocal,
	}
	c.upvalues = append(c.upvalues, upvalue)
	c.Bytecode.UpvalueCount++
	local.hasUpvalue = true
	return upvalue
}

func (c *Compiler) compileModifierExpressionNode(label string, node *ast.ModifierNode, valueIsIgnored bool) expressionResult {
	switch node.Modifier.Type {
	case token.IF:
		return c.compileModifierIfExpression(false, node.Right, node.Left, nil, node.Location(), valueIsIgnored)
	case token.UNLESS:
		return c.compileModifierIfExpression(true, node.Right, node.Left, nil, node.Location(), valueIsIgnored)
	case token.WHILE:
		c.modifierWhileExpression(label, node)
	case token.UNTIL:
		c.modifierUntilExpression(label, node)
	default:
		c.Errors.AddFailure(
			fmt.Sprintf("illegal modifier: %s", node.Modifier.FetchValue()),
			node.Location(),
		)
	}

	return expressionCompiled
}

func (c *Compiler) compileModifierIfExpression(unless bool, condition, then, els ast.ExpressionNode, location *position.Location, valueIsIgnored bool) expressionResult {
	var elsFunc func()
	if els != nil {
		elsFunc = func() {
			c.mustCompileNode(els, valueIsIgnored)
		}
	}
	var jumpOp bytecode.OpCode
	if unless {
		jumpOp = bytecode.JUMP_IF
	} else {
		jumpOp = bytecode.JUMP_UNLESS
	}
	return c.compileIfWithConditionExpression(
		jumpOp,
		condition,
		func() {
			c.mustCompileNode(then, valueIsIgnored)
		},
		elsFunc,
		location,
		valueIsIgnored,
	)
}

func (c *Compiler) compileIfExpression(unless bool, condition ast.ExpressionNode, then, els []ast.StatementNode, location *position.Location, valueIsIgnored bool) expressionResult {
	var elsFunc func()
	if els != nil {
		elsFunc = func() {
			c.compileStatements(els, location, valueIsIgnored)
		}
	}

	var jumpOp bytecode.OpCode
	if unless {
		jumpOp = bytecode.JUMP_IF
	} else {
		jumpOp = bytecode.JUMP_UNLESS
	}

	return c.compileIfWithConditionExpression(
		jumpOp,
		condition,
		func() {
			c.compileStatements(then, location, valueIsIgnored)
		},
		elsFunc,
		location,
		valueIsIgnored,
	)
}

func (c *Compiler) compileIf(jumpOp bytecode.OpCode, condition, then, els func(), location *position.Location, valueIsIgnored bool) expressionResult {
	c.enterScope("", defaultScopeType)
	condition()

	thenJumpOffset := c.emitJump(location.StartPos.Line, jumpOp)

	then()
	c.leaveScope(location.StartPos.Line)

	compileElse := !valueIsIgnored || els != nil
	var elseJumpOffset int
	if compileElse {
		elseJumpOffset = c.emitJump(location.StartPos.Line, bytecode.JUMP)
	}

	c.patchJump(thenJumpOffset, location)

	if compileElse {
		if els != nil {
			c.enterScope("", defaultScopeType)
			els()
			c.leaveScope(location.StartPos.Line)
		} else if !valueIsIgnored {
			c.emit(location.StartPos.Line, bytecode.NIL)
		}
		c.patchJump(elseJumpOffset, location)
	}

	return valueIgnoredToResult(valueIsIgnored)
}

func valueIgnoredToResult(valueIsIgnored bool) expressionResult {
	if valueIsIgnored {
		return expressionCompiledWithoutResult
	}
	return expressionCompiled
}

func (c *Compiler) compileIfWithConditionExpression(jumpOp bytecode.OpCode, condition ast.ExpressionNode, then, els func(), location *position.Location, valueIsIgnored bool) expressionResult {
	if result := resolve(condition); !result.IsUndefined() {
		// if gets optimised away
		c.enterScope("", defaultScopeType)
		defer c.leaveScope(location.StartPos.Line)

		var checkFunc func(value.Value) bool
		switch jumpOp {
		case bytecode.JUMP_UNLESS:
			checkFunc = value.Truthy
		case bytecode.JUMP_IF:
			checkFunc = value.Falsy
		case bytecode.JUMP_IF_NIL:
			checkFunc = value.IsNil
		}

		if checkFunc(result) {
			if then == nil {
				if valueIsIgnored {
					return expressionCompiledWithoutResult
				}
				c.emit(location.StartPos.Line, bytecode.NIL)
				return expressionCompiled
			}
			then()
			return valueIgnoredToResult(valueIsIgnored)
		}

		if els == nil {
			if valueIsIgnored {
				return expressionCompiledWithoutResult
			}
			c.emit(location.StartPos.Line, bytecode.NIL)
			return expressionCompiled
		}
		els()
		return valueIgnoredToResult(valueIsIgnored)
	}

	cond := func() {
		c.compileNodeWithResult(condition)
	}
	if optimisedJumpOp, optimisedCond := c.optimiseCondition(jumpOp, condition, location); optimisedCond != nil {
		jumpOp = optimisedJumpOp
		cond = optimisedCond
	}

	return c.compileIf(
		jumpOp,
		cond,
		then,
		els,
		location,
		valueIsIgnored,
	)
}

func (c *Compiler) optimiseCondition(jumpOp bytecode.OpCode, condition ast.ExpressionNode, location *position.Location) (bytecode.OpCode, func()) {
	cond, ok := condition.(*ast.BinaryExpressionNode)
	if !ok {
		return 0, nil
	}

	switch cond.Op.Type {
	case token.LESS_EQUAL:
		return c.optimiseIfLessEqual(jumpOp, cond, location)
	case token.LESS:
		return c.optimiseIfLess(jumpOp, cond, location)
	case token.GREATER_EQUAL:
		return c.optimiseIfGreaterEqual(jumpOp, cond, location)
	case token.GREATER:
		return c.optimiseIfGreater(jumpOp, cond, location)
	case token.EQUAL_EQUAL:
		return c.optimiseIfEqual(jumpOp, cond, location)
	case token.NOT_EQUAL:
		return c.optimiseIfNotEqual(jumpOp, cond, location)
	}

	return 0, nil
}

func (c *Compiler) optimiseIfNotEqual(jumpOp bytecode.OpCode, condition *ast.BinaryExpressionNode, location *position.Location) (bytecode.OpCode, func()) {
	leftType := c.typeOf(condition.Left)
	rightType := c.typeOf(condition.Right)

	if jumpOp == bytecode.JUMP_UNLESS {
		if c.checker.IsSubtype(leftType, c.checker.StdNil()) {
			return bytecode.JUMP_IF_NIL, func() { c.compileNodeWithResult(condition.Right) }
		}
		if c.checker.IsSubtype(rightType, c.checker.StdNil()) {
			return bytecode.JUMP_IF_NIL, func() { c.compileNodeWithResult(condition.Left) }
		}
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			return bytecode.JUMP_IF_IEQ, func() {
				c.compileNodeWithResult(condition.Left)
				c.compileNodeWithResult(condition.Right)
			}
		}
		if c.checker.IsSubtype(rightType, c.checker.StdInt()) {
			return bytecode.JUMP_IF_IEQ, func() {
				c.compileNodeWithResult(condition.Right)
				c.compileNodeWithResult(condition.Left)
			}
		}
	}

	if jumpOp == bytecode.JUMP_IF {
		if c.checker.IsSubtype(leftType, c.checker.StdNil()) {
			return bytecode.JUMP_UNLESS_NIL, func() { c.compileNodeWithResult(condition.Right) }
		}
		if c.checker.IsSubtype(rightType, c.checker.StdNil()) {
			return bytecode.JUMP_UNLESS_NIL, func() { c.compileNodeWithResult(condition.Left) }
		}
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			return bytecode.JUMP_UNLESS_IEQ, func() {
				c.compileNodeWithResult(condition.Left)
				c.compileNodeWithResult(condition.Right)
			}
		}
		if c.checker.IsSubtype(rightType, c.checker.StdInt()) {
			return bytecode.JUMP_UNLESS_IEQ, func() {
				c.compileNodeWithResult(condition.Right)
				c.compileNodeWithResult(condition.Left)
			}
		}
	}

	return 0, nil
}

func (c *Compiler) optimiseIfEqual(jumpOp bytecode.OpCode, condition *ast.BinaryExpressionNode, location *position.Location) (bytecode.OpCode, func()) {
	leftType := c.typeOf(condition.Left)
	rightType := c.typeOf(condition.Right)

	if jumpOp == bytecode.JUMP_UNLESS {
		if c.checker.IsSubtype(leftType, c.checker.StdNil()) {
			return bytecode.JUMP_UNLESS_NIL, func() { c.compileNodeWithResult(condition.Right) }
		}
		if c.checker.IsSubtype(rightType, c.checker.StdNil()) {
			return bytecode.JUMP_UNLESS_NIL, func() { c.compileNodeWithResult(condition.Left) }
		}
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			return bytecode.JUMP_UNLESS_IEQ, func() {
				c.compileNodeWithResult(condition.Left)
				c.compileNodeWithResult(condition.Right)
			}
		}
		if c.checker.IsSubtype(rightType, c.checker.StdInt()) {
			return bytecode.JUMP_UNLESS_IEQ, func() {
				c.compileNodeWithResult(condition.Right)
				c.compileNodeWithResult(condition.Left)
			}
		}
	}
	if jumpOp == bytecode.JUMP_IF {
		if c.checker.IsSubtype(leftType, c.checker.StdNil()) {
			return bytecode.JUMP_IF_NIL, func() { c.compileNodeWithResult(condition.Right) }
		}
		if c.checker.IsSubtype(rightType, c.checker.StdNil()) {
			return bytecode.JUMP_IF_NIL, func() { c.compileNodeWithResult(condition.Left) }
		}
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			return bytecode.JUMP_IF_IEQ, func() {
				c.compileNodeWithResult(condition.Left)
				c.compileNodeWithResult(condition.Right)
			}
		}
		if c.checker.IsSubtype(rightType, c.checker.StdInt()) {
			return bytecode.JUMP_IF_IEQ, func() {
				c.compileNodeWithResult(condition.Right)
				c.compileNodeWithResult(condition.Left)
			}
		}
	}

	return 0, nil
}

func (c *Compiler) optimiseIfGreater(jumpOp bytecode.OpCode, condition *ast.BinaryExpressionNode, location *position.Location) (bytecode.OpCode, func()) {
	leftType := c.typeOf(condition.Left)
	rightType := c.typeOf(condition.Right)

	if jumpOp == bytecode.JUMP_UNLESS {
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			return bytecode.JUMP_UNLESS_IGT, func() {
				c.compileNodeWithResult(condition.Left)
				c.compileNodeWithResult(condition.Right)
			}
		}
		if c.checker.IsSubtype(rightType, c.checker.StdInt()) {
			return bytecode.JUMP_UNLESS_ILT, func() {
				c.compileNodeWithResult(condition.Right)
				c.compileNodeWithResult(condition.Left)
			}
		}
	}
	if jumpOp == bytecode.JUMP_IF {
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			return bytecode.JUMP_UNLESS_ILE, func() {
				c.compileNodeWithResult(condition.Left)
				c.compileNodeWithResult(condition.Right)
			}
		}
	}

	return 0, nil
}

func (c *Compiler) optimiseIfGreaterEqual(jumpOp bytecode.OpCode, condition *ast.BinaryExpressionNode, location *position.Location) (bytecode.OpCode, func()) {
	leftType := c.typeOf(condition.Left)
	rightType := c.typeOf(condition.Right)

	if jumpOp == bytecode.JUMP_UNLESS {
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			return bytecode.JUMP_UNLESS_IGE, func() {
				c.compileNodeWithResult(condition.Left)
				c.compileNodeWithResult(condition.Right)
			}
		}
		// Reverse only when leftType is subtype of BuiltinComparable
		if c.checker.IsSubtype(rightType, c.checker.StdInt()) {
			return bytecode.JUMP_UNLESS_ILE, func() {
				c.compileNodeWithResult(condition.Right)
				c.compileNodeWithResult(condition.Left)
			}
		}
	}
	if jumpOp == bytecode.JUMP_IF {
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			return bytecode.JUMP_UNLESS_ILT, func() {
				c.compileNodeWithResult(condition.Left)
				c.compileNodeWithResult(condition.Right)
			}
		}
	}

	return 0, nil
}

func (c *Compiler) optimiseIfLess(jumpOp bytecode.OpCode, condition *ast.BinaryExpressionNode, location *position.Location) (bytecode.OpCode, func()) {
	leftType := c.typeOf(condition.Left)
	rightType := c.typeOf(condition.Right)

	if jumpOp == bytecode.JUMP_UNLESS {
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			return bytecode.JUMP_UNLESS_ILT, func() {
				c.compileNodeWithResult(condition.Left)
				c.compileNodeWithResult(condition.Right)
			}
		}
		if c.checker.IsSubtype(rightType, c.checker.StdInt()) {
			return bytecode.JUMP_UNLESS_IGT, func() {
				c.compileNodeWithResult(condition.Right)
				c.compileNodeWithResult(condition.Left)
			}
		}
	}
	if jumpOp == bytecode.JUMP_IF {
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			return bytecode.JUMP_UNLESS_IGE, func() {
				c.compileNodeWithResult(condition.Left)
				c.compileNodeWithResult(condition.Right)
			}
		}
	}

	return 0, nil
}

func (c *Compiler) optimiseIfLessEqual(jumpOp bytecode.OpCode, condition *ast.BinaryExpressionNode, location *position.Location) (bytecode.OpCode, func()) {
	leftType := c.typeOf(condition.Left)
	rightType := c.typeOf(condition.Right)

	if jumpOp == bytecode.JUMP_UNLESS {
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			return bytecode.JUMP_UNLESS_ILE, func() {
				c.compileNodeWithResult(condition.Left)
				c.compileNodeWithResult(condition.Right)
			}
		}
		if c.checker.IsSubtype(rightType, c.checker.StdInt()) {
			return bytecode.JUMP_UNLESS_IGE, func() {
				c.compileNodeWithResult(condition.Right)
				c.compileNodeWithResult(condition.Left)
			}
		}
	}
	if jumpOp == bytecode.JUMP_IF {
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			return bytecode.JUMP_UNLESS_IGT, func() {
				c.compileNodeWithResult(condition.Left)
				c.compileNodeWithResult(condition.Right)
			}
		}
	}

	return 0, nil
}

func (c *Compiler) compileValueDeclarationNode(node *ast.ValueDeclarationNode, valueIsIgnored bool) expressionResult {
	initialised := node.Initialiser != nil

	local := c.defineLocal(identifierToName(node.Name), node.Location())
	if local == nil {
		return valueIgnoredToResult(valueIsIgnored)
	}
	if initialised {
		c.compileNodeWithResult(node.Initialiser)
	}

	if initialised {
		return c.emitSetLocal(node.Location().StartPos.Line, local.index, valueIsIgnored)
	}

	if !valueIsIgnored {
		c.emit(node.Location().StartPos.Line, bytecode.NIL)
		return expressionCompiled
	}

	return expressionCompiledWithoutResult
}

func (c *Compiler) compileAwaitExpressionNode(node *ast.AwaitExpressionNode) {
	location := node.Location()
	c.compileNodeWithResult(node.Value)
	if !c.isAsync || node.Sync {
		c.emit(location.StartPos.Line, bytecode.AWAIT_SYNC)
		return
	}

	c.emit(location.StartPos.Line, bytecode.AWAIT)
	c.emit(location.StartPos.Line, bytecode.AWAIT_RESULT)
}

func (c *Compiler) compileReturnExpressionNode(node *ast.ReturnExpressionNode) {
	location := node.Location()
	if node.Value != nil {
		c.emitReturn(location, node.Value)
	} else {
		c.emit(location.StartPos.Line, bytecode.NIL)
		c.emitReturn(location, nil)
	}
}

func (c *Compiler) compileYieldExpressionNode(node *ast.YieldExpressionNode) {
	location := node.Location()
	if node.Value != nil {
		c.emitYield(location, node.Value)
	} else {
		c.emit(location.StartPos.Line, bytecode.NIL)
		c.emitYield(location, nil)
	}
}

func (c *Compiler) compileNilSafeSubscriptExpressionNode(node *ast.NilSafeSubscriptExpressionNode) expressionResult {
	if c.resolveAndEmit(node) {
		return expressionCompiled
	}

	return c.compileIfWithConditionExpression(
		bytecode.JUMP_IF_NIL,
		node.Receiver,
		func() {
			c.compileNodeWithResult(node.Receiver)
			c.compileNodeWithResult(node.Key)

			receiverType := c.typeOf(node.Receiver)
			c.compileSubscript(receiverType, node.Location())
		},
		func() {},
		node.Location(),
		false,
	)
}

func (c *Compiler) relationalPattern(pattern ast.Node, opcode bytecode.OpCode) {
	location := pattern.Location()

	c.compileIf(
		bytecode.JUMP_UNLESS,
		func() {
			c.emit(location.StartPos.Line, bytecode.DUP)
			c.compileNodeWithResult(pattern)
			c.emit(location.StartPos.Line, bytecode.DUP_2)
			c.emit(location.StartPos.Line, bytecode.SWAP)
			c.emit(location.StartPos.Line, bytecode.GET_CLASS)
			c.emit(location.StartPos.Line, bytecode.IS_A)
		},
		func() {
			c.emit(location.StartPos.Line, opcode)
		},
		func() {
			c.emit(location.StartPos.Line, bytecode.POP_2)
			c.emit(location.StartPos.Line, bytecode.FALSE)
		},
		location,
		false,
	)
}

func (c *Compiler) literalPattern(pattern ast.Node, opcode bytecode.OpCode) {
	location := pattern.Location()
	c.emit(location.StartPos.Line, bytecode.DUP)
	c.compileNodeWithResult(pattern)
	c.emit(location.StartPos.Line, opcode)
}

func (c *Compiler) pattern(pattern ast.PatternNode) {
	location := pattern.Location()
	switch pat := pattern.(type) {
	case *ast.TrueLiteralNode, *ast.FalseLiteralNode, *ast.NilLiteralNode,
		*ast.CharLiteralNode, *ast.RawCharLiteralNode, *ast.DoubleQuotedStringLiteralNode,
		*ast.InterpolatedStringLiteralNode, *ast.RawStringLiteralNode,
		*ast.SimpleSymbolLiteralNode, *ast.InterpolatedSymbolLiteralNode,
		*ast.IntLiteralNode, *ast.Int64LiteralNode, *ast.UInt64LiteralNode,
		*ast.Int32LiteralNode, *ast.UInt32LiteralNode, *ast.Int16LiteralNode, *ast.UInt16LiteralNode,
		*ast.Int8LiteralNode, *ast.UInt8LiteralNode, *ast.FloatLiteralNode,
		*ast.Float64LiteralNode, *ast.Float32LiteralNode, *ast.BigFloatLiteralNode,
		*ast.PublicConstantNode, *ast.PrivateConstantNode, *ast.ConstantLookupNode:
		c.literalPattern(
			pat,
			bytecode.EQUAL,
		)
	case *ast.RangeLiteralNode:
		c.emit(location.StartPos.Line, bytecode.DUP)
		c.compileRangeLiteralNode(pat)
		c.emit(location.StartPos.Line, bytecode.SWAP)
		callInfo := value.NewCallSiteInfo(symbol.S_contains, 1)
		c.emitCallMethod(callInfo, location, false)
	case *ast.PublicIdentifierNode:
		switch c.mode {
		case valuePatternDeclarationNode:
			c.defineLocal(pat.Value, location)
		default:
			c.defineLocalOverrideCurrentScope(pat.Value, location)
		}
		c.setLocalWithoutValue(pat.Value, location, false)
		c.emit(location.StartPos.Line, bytecode.TRUE)
	case *ast.PrivateIdentifierNode:
		switch c.mode {
		case valuePatternDeclarationNode:
			c.defineLocal(pat.Value, location)
		default:
			c.defineLocalOverrideCurrentScope(pat.Value, location)
		}
		c.setLocalWithoutValue(pat.Value, location, false)
		c.emit(location.StartPos.Line, bytecode.TRUE)
	case *ast.ObjectPatternNode:
		c.objectPattern(pat)
	case *ast.AsPatternNode:
		c.asPattern(pat)
	case *ast.UninterpolatedRegexLiteralNode, *ast.InterpolatedRegexLiteralNode:
		c.emit(location.StartPos.Line, bytecode.DUP)
		c.compileNode(pat, false)
		c.emit(location.StartPos.Line, bytecode.SWAP)
		callInfo := value.NewCallSiteInfo(matchesSymbol, 1)
		c.emitCallMethod(callInfo, location, false)
	case *ast.UnaryExpressionNode:
		c.unaryPattern(pat)
	case *ast.BinaryPatternNode:
		c.binaryPattern(pat)
	case *ast.MapPatternNode:
		c.mapOrRecordPattern(c.typeOf(pat), pat.Location(), pat.Elements, true)
	case *ast.RecordPatternNode:
		c.mapOrRecordPattern(c.typeOf(pat), pat.Location(), pat.Elements, false)
	case *ast.SetPatternNode:
		c.setPattern(pat.Location(), pat.Elements)
	case *ast.ListPatternNode:
		c.listOrTuplePattern(c.typeOf(pat), pat.Location(), pat.Elements, true)
	case *ast.TuplePatternNode:
		c.listOrTuplePattern(c.typeOf(pat), pat.Location(), pat.Elements, false)
	case *ast.WordArrayListLiteralNode, *ast.SymbolArrayListLiteralNode, *ast.BinArrayListLiteralNode, *ast.HexArrayListLiteralNode,
		*ast.WordArrayTupleLiteralNode, *ast.SymbolArrayTupleLiteralNode, *ast.BinArrayTupleLiteralNode, *ast.HexArrayTupleLiteralNode,
		*ast.WordHashSetLiteralNode, *ast.SymbolHashSetLiteralNode, *ast.BinHashSetLiteralNode, *ast.HexHashSetLiteralNode:
		c.specialCollectionPattern(pat)
	case *ast.MacroBoundaryNode:
		stmt := pat.Body[0].(*ast.PatternStatementNode)
		c.pattern(stmt.Pattern)
	default:
		c.Errors.AddFailure(
			fmt.Sprintf("compilation of this pattern has not been implemented: %T", pattern),
			location,
		)
	}
}

func (c *Compiler) unaryPattern(pat *ast.UnaryExpressionNode) {
	switch pat.Op.Type {
	case token.EQUAL_EQUAL:
		c.literalPattern(
			pat.Right,
			bytecode.EQUAL,
		)
	case token.NOT_EQUAL:
		c.literalPattern(
			pat.Right,
			bytecode.NOT_EQUAL,
		)
	case token.LAX_EQUAL:
		c.literalPattern(
			pat.Right,
			bytecode.LAX_EQUAL,
		)
	case token.LAX_NOT_EQUAL:
		c.literalPattern(
			pat.Right,
			bytecode.LAX_NOT_EQUAL,
		)
	case token.STRICT_EQUAL:
		c.literalPattern(
			pat.Right,
			bytecode.STRICT_EQUAL,
		)
	case token.STRICT_NOT_EQUAL:
		c.literalPattern(
			pat.Right,
			bytecode.STRICT_NOT_EQUAL,
		)
	case token.LESS:
		c.relationalPattern(
			pat.Right,
			bytecode.LESS,
		)
	case token.LESS_EQUAL:
		c.relationalPattern(
			pat.Right,
			bytecode.LESS_EQUAL,
		)
	case token.GREATER:
		c.relationalPattern(
			pat.Right,
			bytecode.GREATER,
		)
	case token.GREATER_EQUAL:
		c.relationalPattern(
			pat.Right,
			bytecode.GREATER_EQUAL,
		)
	default:
		c.literalPattern(
			pat,
			bytecode.EQUAL,
		)
	}
}

func (c *Compiler) binaryPattern(pat *ast.BinaryPatternNode) {
	location := pat.Location()
	var op bytecode.OpCode
	switch pat.Op.Type {
	case token.OR_OR:
		op = bytecode.JUMP_IF_NP
	case token.AND_AND:
		op = bytecode.JUMP_UNLESS_NP
	default:
		panic(fmt.Sprintf("invalid binary pattern operator: %s", pat.Op.Type.Name()))
	}

	c.pattern(pat.Left)
	jump := c.emitJump(location.StartPos.Line, op)

	// branch one
	c.emit(location.StartPos.Line, bytecode.POP)
	c.pattern(pat.Right)

	// branch two
	c.patchJump(jump, location)
}

func (c *Compiler) asPattern(node *ast.AsPatternNode) {
	location := node.Location()
	var varName string
	switch n := node.Name.(type) {
	case *ast.PrivateIdentifierNode:
		varName = n.Value
	case *ast.PublicIdentifierNode:
		varName = n.Value
	default:
		panic(fmt.Sprintf("invalid as pattern name: %#v", node.Name))
	}

	switch c.mode {
	case valuePatternDeclarationNode:
		c.defineLocal(varName, location)
	default:
		c.defineLocalOverrideCurrentScope(varName, location)
	}
	c.setLocalWithoutValue(varName, location, false)
	c.pattern(node.Pattern)
}

func (c *Compiler) identifierObjectPatternAttribute(name string, location *position.Location) {
	c.emit(location.StartPos.Line, bytecode.DUP)
	callInfo := value.NewCallSiteInfo(value.ToSymbol(name), 0)
	c.emitCallMethod(callInfo, location, false)

	var identVar *local
	switch c.mode {
	case valuePatternDeclarationNode:
		identVar = c.defineLocal(name, location)
	default:
		identVar = c.defineLocalOverrideCurrentScope(name, location)
	}
	c.emitSetLocalPop(location.StartPos.Line, identVar.index)
}

func (c *Compiler) objectPattern(node *ast.ObjectPatternNode) {
	var jumpsToPatch []int
	c.enterPattern()

	location := node.Location()
	c.emit(node.ObjectType.Location().StartPos.Line, bytecode.DUP)
	c.compileNodeWithResult(node.ObjectType)
	c.emit(node.ObjectType.Location().StartPos.Line, bytecode.IS_A)

	jmp := c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NP)
	jumpsToPatch = append(jumpsToPatch, jmp)
	c.emit(location.StartPos.Line, bytecode.POP)

	for _, attr := range node.Attributes {
		location := attr.Location()
		switch e := attr.(type) {
		case *ast.SymbolKeyValuePatternNode:
			c.emit(location.StartPos.Line, bytecode.DUP)
			callInfo := value.NewCallSiteInfo(value.ToSymbol(e.Key), 0)
			c.emitCallMethod(callInfo, location, false)

			c.pattern(e.Value)
			c.emit(location.StartPos.Line, bytecode.POP_SKIP_ONE)
			jmp := c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NP)
			jumpsToPatch = append(jumpsToPatch, jmp)
			c.emit(location.StartPos.Line, bytecode.POP)
		case *ast.PublicIdentifierNode:
			c.identifierObjectPatternAttribute(e.Value, location)
		case *ast.PrivateIdentifierNode:
			c.identifierObjectPatternAttribute(e.Value, location)
		default:
			c.Errors.AddFailure(
				fmt.Sprintf("invalid object pattern attribute: %T", attr),
				location,
			)
		}
	}

	// leave true as the result of the happy path
	c.emit(location.StartPos.Line, bytecode.TRUE)

	// leave false on the stack from the falsy if that jumped here
	for _, jmp := range jumpsToPatch {
		c.patchJump(jmp, location)
	}
	c.leavePattern()
}

func (c *Compiler) specialCollectionPattern(node ast.PatternNode) {
	location := node.Location()
	c.emit(location.StartPos.Line, bytecode.DUP)
	switch node.(type) {
	case *ast.WordArrayListLiteralNode, *ast.SymbolArrayListLiteralNode, *ast.BinArrayListLiteralNode, *ast.HexArrayListLiteralNode:
		c.emitValue(value.Ref(value.ListMixin), location)
	case *ast.WordArrayTupleLiteralNode, *ast.SymbolArrayTupleLiteralNode, *ast.BinArrayTupleLiteralNode, *ast.HexArrayTupleLiteralNode:
		c.emitValue(value.Ref(value.TupleMixin), location)
	case *ast.WordHashSetLiteralNode, *ast.SymbolHashSetLiteralNode, *ast.BinHashSetLiteralNode, *ast.HexHashSetLiteralNode:
		c.emitValue(value.Ref(value.SetMixin), location)
	default:
		panic(fmt.Sprintf("invalid special collection pattern node: %#v", node))
	}
	c.emit(location.StartPos.Line, bytecode.IS_A)

	jmp := c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NP)
	c.emit(location.StartPos.Line, bytecode.POP)

	c.emit(location.StartPos.Line, bytecode.DUP)
	c.compileNodeWithResult(node)
	c.emit(location.StartPos.Line, bytecode.LAX_EQUAL)

	// leave false on the stack from the falsy if that jumped here
	c.patchJump(jmp, location)
}

func (c *Compiler) identifierMapPatternElement(name string, collectionType types.Type, location *position.Location) {
	c.emit(location.StartPos.Line, bytecode.DUP)
	c.emitValue(value.ToSymbol(name).ToValue(), location)
	c.compileSubscript(collectionType, location)

	var identVar *local
	switch c.mode {
	case valuePatternDeclarationNode:
		identVar = c.defineLocal(name, location)
	default:
		identVar = c.defineLocalOverrideCurrentScope(name, location)
	}
	if identVar == nil {
		return
	}
	c.emitSetLocalNoPop(location.StartPos.Line, identVar.index)
	c.emit(location.StartPos.Line, bytecode.POP)
}

func (c *Compiler) mapOrRecordPattern(typ types.Type, location *position.Location, elements []ast.PatternNode, isMap bool) {
	var jumpsToPatch []int
	c.enterPattern()

	c.emit(location.StartPos.Line, bytecode.DUP)
	if isMap {
		c.emitValue(value.Ref(value.MapMixin), location)
	} else {
		c.emitValue(value.Ref(value.RecordMixin), location)
	}
	c.emit(location.StartPos.Line, bytecode.IS_A)

	jmp := c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NP)
	jumpsToPatch = append(jumpsToPatch, jmp)
	c.emit(location.StartPos.Line, bytecode.POP)

	for _, element := range elements {
		location := element.Location()
		switch e := element.(type) {
		case *ast.SymbolKeyValuePatternNode:
			c.emit(location.StartPos.Line, bytecode.DUP)
			c.emitValue(value.ToSymbol(e.Key).ToValue(), location)
			c.compileSubscript(typ, location)

			c.pattern(e.Value)
			c.emit(location.StartPos.Line, bytecode.POP_SKIP_ONE)
			jmp := c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NP)
			jumpsToPatch = append(jumpsToPatch, jmp)
			c.emit(location.StartPos.Line, bytecode.POP)
		case *ast.KeyValuePatternNode:
			c.emit(location.StartPos.Line, bytecode.DUP)
			c.compileNodeWithResult(e.Key)
			c.compileSubscript(typ, location)

			c.pattern(e.Value)
			c.emit(location.StartPos.Line, bytecode.POP_SKIP_ONE)
			jmp := c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NP)
			jumpsToPatch = append(jumpsToPatch, jmp)
			c.emit(location.StartPos.Line, bytecode.POP)
		case *ast.PublicIdentifierNode:
			c.identifierMapPatternElement(e.Value, typ, location)
		case *ast.PrivateIdentifierNode:
			c.identifierMapPatternElement(e.Value, typ, location)
		default:
			c.Errors.AddFailure(
				fmt.Sprintf("invalid map pattern element: %T", element),
				location,
			)
		}
	}

	// leave true as the result of the happy path
	c.emit(location.StartPos.Line, bytecode.TRUE)

	// leave false on the stack from the falsy if that jumped here
	for _, jmp := range jumpsToPatch {
		c.patchJump(jmp, location)
	}
	c.leavePattern()
}

func (c *Compiler) setPattern(location *position.Location, elements []ast.PatternNode) {
	var jumpsToPatch []int
	var subPatternElements []ast.PatternNode

	var restElementIsPresent bool
	for _, element := range elements {
		switch e := element.(type) {
		case *ast.RestPatternNode:
			if restElementIsPresent {
				c.Errors.AddFailure(
					"there should be only a single rest element",
					element.Location(),
				)
			}
			restElementIsPresent = true
		default:
			subPatternElements = append(subPatternElements, e)
		}
	}
	c.enterPattern()

	c.emit(location.StartPos.Line, bytecode.DUP)
	c.emitValue(value.Ref(value.SetMixin), location)
	c.emit(location.StartPos.Line, bytecode.IS_A)

	jmp := c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NP)
	jumpsToPatch = append(jumpsToPatch, jmp)
	c.emit(location.StartPos.Line, bytecode.POP)

	c.emit(location.StartPos.Line, bytecode.DUP)
	callInfo := value.NewCallSiteInfo(symbol.L_length, 0)
	c.emitCallMethod(callInfo, location, false)

	if !restElementIsPresent {
		c.emitValue(value.SmallInt(len(subPatternElements)).ToValue(), location)
		c.emit(location.StartPos.Line, bytecode.EQUAL)
	} else {
		c.emitValue(value.SmallInt(len(subPatternElements)).ToValue(), location)
		c.emit(location.StartPos.Line, bytecode.GREATER_EQUAL)
	}

	jmp = c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NP)
	jumpsToPatch = append(jumpsToPatch, jmp)
	c.emit(location.StartPos.Line, bytecode.POP)

subPatternLoop:
	for _, element := range subPatternElements {
		switch element.(type) {
		case *ast.PrivateIdentifierNode:
			continue subPatternLoop
		}

		location := element.Location()
		c.emit(location.StartPos.Line, bytecode.DUP)
		c.compileNodeWithResult(element)
		callInfo := value.NewCallSiteInfo(symbol.L_contains, 1)
		c.emitCallMethod(callInfo, location, false)

		jmp := c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NP)
		jumpsToPatch = append(jumpsToPatch, jmp)
		c.emit(location.StartPos.Line, bytecode.POP)
	}

	// leave true as the result of the happy path
	c.emit(location.StartPos.Line, bytecode.TRUE)

	// leave false on the stack from the falsy if that jumped here
	for _, jmp := range jumpsToPatch {
		c.patchJump(jmp, location)
	}
	c.leavePattern()
}

func (c *Compiler) listOrTuplePattern(typ types.Type, location *position.Location, elements []ast.PatternNode, isList bool) {
	var jumpsToPatch []int

	var restVariableName string
	elementBeforeRestCount := -1
	for i, element := range elements {
		switch e := element.(type) {
		case *ast.RestPatternNode:
			if elementBeforeRestCount != -1 {
				c.Errors.AddFailure(
					"there should be only a single rest element",
					element.Location(),
				)
			}
			elementBeforeRestCount = i
			switch ident := e.Identifier.(type) {
			case *ast.PrivateIdentifierNode:
				restVariableName = ident.Value
			case *ast.PublicIdentifierNode:
				restVariableName = ident.Value
			case nil:
			default:
				return
			}
		}
	}
	elementAfterRestCount := len(elements) - 1 - elementBeforeRestCount
	var restListVar *local
	if restVariableName != "" {
		c.emit(location.StartPos.Line, bytecode.UNDEFINED)
		c.emit(location.StartPos.Line, bytecode.UNDEFINED)
		c.emitNewArrayList(0, location)
		restListVar = c.defineLocal(restVariableName, location)
		c.emitSetLocalNoPop(location.StartPos.Line, restListVar.index)
		c.emit(location.StartPos.Line, bytecode.POP)
	}
	c.enterPattern()

	c.emit(location.StartPos.Line, bytecode.DUP)
	if isList {
		c.emitValue(value.Ref(value.ListMixin), location)
	} else {
		c.emitValue(value.Ref(value.TupleMixin), location)
	}
	c.emit(location.StartPos.Line, bytecode.IS_A)

	jmp := c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NP)
	jumpsToPatch = append(jumpsToPatch, jmp)
	c.emit(location.StartPos.Line, bytecode.POP)

	c.emit(location.StartPos.Line, bytecode.DUP)
	callInfo := value.NewCallSiteInfo(symbol.L_length, 0)
	c.emitCallMethod(callInfo, location, false)
	var lengthVar *local
	if elementBeforeRestCount != -1 {
		lengthVar = c.defineLocal(fmt.Sprintf("#!listPatternLength%d", c.patternNesting), location)
		c.emitSetLocalNoPop(location.StartPos.Line, lengthVar.index)
	}

	if elementBeforeRestCount == -1 {
		c.emitValue(value.SmallInt(len(elements)).ToValue(), location)
		c.emit(location.StartPos.Line, bytecode.EQUAL_INT)
	} else {
		staticElementCount := elementBeforeRestCount + elementAfterRestCount
		c.emitValue(value.SmallInt(staticElementCount).ToValue(), location)
		c.emit(location.StartPos.Line, bytecode.GREATER_EQUAL_I)
	}

	jmp = c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NP)
	jumpsToPatch = append(jumpsToPatch, jmp)
	c.emit(location.StartPos.Line, bytecode.POP)

	elementsBeforeRest := elements
	if elementBeforeRestCount != -1 {
		elementsBeforeRest = elements[:elementBeforeRestCount]
	}
	for i, element := range elementsBeforeRest {
		location := element.Location()
		c.emit(location.StartPos.Line, bytecode.DUP)
		c.emitValue(value.SmallInt(i).ToValue(), element.Location())
		c.compileSubscript(typ, location)

		c.pattern(element)
		c.emit(location.StartPos.Line, bytecode.POP_SKIP_ONE)
		jmp := c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NP)
		jumpsToPatch = append(jumpsToPatch, jmp)
		c.emit(location.StartPos.Line, bytecode.POP)
	}

	if elementBeforeRestCount != -1 {
		iteratorVar := c.defineLocal(fmt.Sprintf("#!listPatternIterator%d", c.patternNesting), location)

		if restVariableName != "" {
			// adjust the length variable
			// length -= element_after_rest_count
			if elementAfterRestCount != 0 {
				c.emitGetLocal(location.StartPos.Line, lengthVar.index)
				c.emitValue(value.SmallInt(elementAfterRestCount).ToValue(), location)
				c.emit(location.StartPos.Line, bytecode.SUBTRACT_INT)
				c.emitSetLocalNoPop(location.StartPos.Line, lengthVar.index)
				c.emit(location.StartPos.Line, bytecode.POP)
			}

			// create the iterator variable
			// i := element_before_rest_count
			c.emitValue(value.SmallInt(elementBeforeRestCount).ToValue(), location)
			c.emitSetLocalNoPop(location.StartPos.Line, iteratorVar.index)
			c.emit(location.StartPos.Line, bytecode.POP)

			// loop header
			// i < length
			loopStartOffset := c.nextInstructionOffset()
			c.emitGetLocal(location.StartPos.Line, iteratorVar.index)
			c.emitGetLocal(location.StartPos.Line, lengthVar.index)
			loopEndJump := c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_ILT)

			// loop body
			c.emit(location.StartPos.Line, bytecode.DUP)
			c.emitGetLocal(location.StartPos.Line, iteratorVar.index)
			c.compileSubscript(typ, location)
			c.emitGetLocal(location.StartPos.Line, restListVar.index)
			c.emit(location.StartPos.Line, bytecode.SWAP)
			c.emit(location.StartPos.Line, bytecode.APPEND) // append to the list
			c.emit(location.StartPos.Line, bytecode.POP)
			// i++
			c.emitGetLocal(location.StartPos.Line, iteratorVar.index)
			c.compileIncrement(c.checker.StdInt(), location)
			c.emitSetLocalNoPop(location.StartPos.Line, iteratorVar.index)
			c.emit(location.StartPos.Line, bytecode.POP)

			c.emitLoop(location, loopStartOffset)
			// loop end
			c.patchJump(loopEndJump, location)
		} else {
			// create the iterator variable
			// i := length - element_after_rest_count
			if elementAfterRestCount == 0 {
				c.emitGetLocal(location.StartPos.Line, lengthVar.index)
				c.emitSetLocalNoPop(location.StartPos.Line, iteratorVar.index)
				c.emit(location.StartPos.Line, bytecode.POP)
			} else {
				c.emitGetLocal(location.StartPos.Line, lengthVar.index)
				c.emitValue(value.SmallInt(elementAfterRestCount).ToValue(), location)
				c.emit(location.StartPos.Line, bytecode.SUBTRACT_INT)
				c.emitSetLocalNoPop(location.StartPos.Line, iteratorVar.index)
				c.emit(location.StartPos.Line, bytecode.POP)
			}
		}

		elementsAfterRest := elements[elementBeforeRestCount+1:]
		for _, element := range elementsAfterRest {
			location := element.Location()
			c.emit(location.StartPos.Line, bytecode.DUP)
			c.emitGetLocal(location.StartPos.Line, iteratorVar.index)
			c.compileSubscript(typ, location)

			c.pattern(element)
			c.emit(location.StartPos.Line, bytecode.POP_SKIP_ONE)
			jmp := c.emitJump(location.StartPos.Line, bytecode.JUMP_UNLESS_NP)
			jumpsToPatch = append(jumpsToPatch, jmp)
			c.emit(location.StartPos.Line, bytecode.POP)

			// i++
			c.emitGetLocal(location.StartPos.Line, iteratorVar.index)
			c.compileIncrement(c.checker.StdInt(), location)
			c.emitSetLocalPop(location.StartPos.Line, iteratorVar.index)
		}
	}

	// leave true as the result of the happy path
	c.emit(location.StartPos.Line, bytecode.TRUE)

	// leave false on the stack from the falsy if that jumped here
	for _, jmp := range jumpsToPatch {
		c.patchJump(jmp, location)
	}
	c.leavePattern()
}

var matchesSymbol = value.ToSymbol("matches")

func (c *Compiler) enterPattern() {
	c.patternNesting++
}

func (c *Compiler) leavePattern() {
	c.patternNesting--
}

func (c *Compiler) compileSwitchExpressionNode(node *ast.SwitchExpressionNode, valueIsIgnored bool) expressionResult {
	location := node.Location()

	c.enterScope("", defaultScopeType)
	c.compileNodeWithResult(node.Value)

	var jumpToEndOffsets []int

	for _, caseNode := range node.Cases {
		c.enterScope("", defaultScopeType)

		caseSpan := caseNode.Location()
		c.pattern(caseNode.Pattern)

		jumpOverBodyOffset := c.emitJump(caseSpan.StartPos.Line, bytecode.JUMP_UNLESS)

		c.emit(caseSpan.StartPos.Line, bytecode.POP)

		c.compileStatements(caseNode.Body, caseSpan, valueIsIgnored)

		c.leaveScopeWithoutMutating(caseSpan.EndPos.Line)

		jumpToEndOffset := c.emitJump(caseSpan.EndPos.Line, bytecode.JUMP)
		jumpToEndOffsets = append(jumpToEndOffsets, jumpToEndOffset)

		c.patchJump(jumpOverBodyOffset, caseSpan)
		c.leaveScope(caseSpan.EndPos.Line)
	}

	c.emit(location.StartPos.Line, bytecode.POP)
	c.compileStatements(node.ElseBody, location, valueIsIgnored)

	for _, offset := range jumpToEndOffsets {
		c.patchJump(offset, location)
	}

	c.leaveScope(location.EndPos.Line)
	return valueIgnoredToResult(valueIsIgnored)
}

func (c *Compiler) compileSubscriptExpressionNode(node *ast.SubscriptExpressionNode) {
	if c.resolveAndEmit(node) {
		return
	}

	c.compileNodeWithResult(node.Receiver)
	c.compileNodeWithResult(node.Key)

	receiverType := c.typeOf(node.Receiver)
	c.compileSubscript(receiverType, node.Location())
}

func (c *Compiler) compileAttributeAccessNode(node *ast.AttributeAccessNode) {
	c.compileNodeWithResult(node.Receiver)

	name := identifierToName(node.AttributeName)
	nameSymbol := value.ToSymbol(name)
	callInfo := value.NewCallSiteInfo(nameSymbol, 0)
	if name == "call" {
		c.emitCall(callInfo, node.Location())
	} else {
		c.emitCallMethod(callInfo, node.Location(), false)
	}
}

func (c *Compiler) compileConstructorCallNode(node *ast.ConstructorCallNode) {
	c.compileConstructorCall(
		func() {
			c.compileNodeWithResult(node.ClassNode)
		},
		node.PositionalArguments,
		node.Location(),
	)
}

func (c *Compiler) compileNewExpressionNode(node *ast.NewExpressionNode) {
	c.compileConstructorCall(
		func() {
			c.emit(node.Location().StartPos.Line, bytecode.SELF)
		},
		node.PositionalArguments,
		node.Location(),
	)
}

func (c *Compiler) compileGenericConstructorCallNode(node *ast.GenericConstructorCallNode) {
	c.compileConstructorCall(
		func() {
			c.compileNodeWithResult(node.ClassNode)
		},
		node.PositionalArguments,
		node.Location(),
	)
}

func (c *Compiler) compileConstructorCall(class func(), args []ast.ExpressionNode, location *position.Location) {
	class()
	for _, posArg := range args {
		c.compileNodeWithResult(posArg)
	}

	c.emitInstantiate(len(args), location)
}

func (c *Compiler) compileMethodCallNode(node *ast.MethodCallNode) {
	c.compileMethodCall(
		node.Receiver,
		node.Op,
		node.MethodName,
		node.PositionalArguments,
		node.TailCall,
		node.Location(),
	)
}
func (c *Compiler) compileGenericMethodCallNode(node *ast.GenericMethodCallNode) {
	c.compileMethodCall(
		node.Receiver,
		node.Op,
		node.MethodName,
		node.PositionalArguments,
		node.TailCall,
		node.Location(),
	)
}

func (c *Compiler) compileMethodCall(receiver ast.ExpressionNode, op *token.Token, nameNode ast.IdentifierNode, args []ast.ExpressionNode, tailCall bool, location *position.Location) {
	_, onSelf := receiver.(*ast.SelfLiteralNode)
	name := identifierToName(nameNode)

	switch op.Type {
	case token.QUESTION_DOT:
		c.compileNodeWithResult(receiver)
		nilJump := c.emitJump(location.StartPos.Line, bytecode.JUMP_IF_NIL_NP)

		// if not nil
		// call the method
		c.compileInnerMethodCall(receiver, name, op, args, false, tailCall, location)

		// if nil
		// leave nil on the stack
		c.patchJump(nilJump, location)
	case token.QUESTION_DOT_DOT:
		c.compileNodeWithResult(receiver)
		c.emit(location.EndPos.Line, bytecode.DUP)
		nilJump := c.emitJump(location.StartPos.Line, bytecode.JUMP_IF_NIL_NP)

		// if not nil
		// call the method
		c.compileInnerMethodCall(receiver, name, op, args, false, tailCall, location)

		// if nil
		// leave nil on the stack
		c.patchJump(nilJump, location)
	case token.DOT_DOT:
		if !onSelf {
			c.compileNodeWithResult(receiver)
		}
		c.emit(location.EndPos.Line, bytecode.DUP)
		c.compileInnerMethodCall(receiver, name, op, args, onSelf, tailCall, location)
	case token.DOT:
		if !onSelf {
			c.compileNodeWithResult(receiver)
		}
		c.compileInnerMethodCall(receiver, name, op, args, onSelf, tailCall, location)
	default:
		panic(fmt.Sprintf("invalid method call operator: %#v", op))
	}
}

func (c *Compiler) compileInnerMethodCall(receiver ast.ExpressionNode, name string, op *token.Token, args []ast.ExpressionNode, onSelf bool, tailCall bool, location *position.Location) {
	for _, posArg := range args {
		c.compileNodeWithResult(posArg)
	}

	receiverType := c.typeOf(receiver)
	nameSym := value.ToSymbol(name)
	callInfo := value.NewCallSiteInfo(nameSym, len(args))
	if onSelf {
		c.emitCallSelf(callInfo, location, tailCall)
	} else {
		switch name {
		case "call":
			c.emitCall(callInfo, location)
		case "++":
			c.compileIncrement(receiverType, location)
		case "--":
			c.compileDecrement(receiverType, location)
		default:
			c.emitCallMethod(callInfo, location, tailCall)
		}
	}

	switch op.Type {
	case token.DOT_DOT, token.QUESTION_DOT_DOT:
		c.emit(location.EndPos.Line, bytecode.POP)
	case token.DOT, token.QUESTION_DOT:
	default:
		panic(fmt.Sprintf("invalid method call operator: %#v", op))
	}
}

func (c *Compiler) compileCallNode(node *ast.CallNode) {
	c.compileNodeWithResult(node.Receiver)

	if node.NilSafe {
		nilJump := c.emitJump(node.Location().StartPos.Line, bytecode.JUMP_IF_NIL_NP)

		// if not nil
		// call the method
		c.compileInnerCall(node)

		// if nil
		// leave nil on the stack
		c.patchJump(nilJump, node.Location())
		return
	}

	c.compileInnerCall(node)
}

func (c *Compiler) compileInnerCall(node *ast.CallNode) {
	for _, posArg := range node.PositionalArguments {
		c.compileNodeWithResult(posArg)
	}

	name := value.ToSymbol("call")
	callInfo := value.NewCallSiteInfo(name, len(node.PositionalArguments))
	c.emitCall(callInfo, node.Location())
}

func (c *Compiler) singletonBlockIsCompilable(node *ast.SingletonBlockExpressionNode) bool {
	if len(node.Body) <= 0 {
		return false
	}

	location := node.Location()
	singletonType := c.typeOf(node).(*types.SingletonClass)
	singletonName := singletonType.Name()

	singletonCompiler := New(fmt.Sprintf("<singleton_class: %s>", singletonName), namespaceMode, location, c.checker)
	singletonCompiler.Errors = c.Errors
	if !singletonCompiler.compileNamespace(node) {
		return false
	}

	node.Bytecode = singletonCompiler.Bytecode
	return true
}

func (c *Compiler) compileSingletonBlockExpressionNode(node *ast.SingletonBlockExpressionNode) expressionResult {
	if node.Bytecode == nil {
		return expressionIgnored
	}

	location := node.Location()
	c.emit(location.StartPos.Line, bytecode.SELF)
	c.emit(location.StartPos.Line, bytecode.GET_SINGLETON)

	c.emitValue(value.Ref(node.Bytecode), location)
	c.emit(location.StartPos.Line, bytecode.INIT_NAMESPACE)
	return expressionCompiled
}

func (c *Compiler) compileGoExpressionNode(node *ast.GoExpressionNode) {
	closureCompiler := New("<closure>", methodMode, node.Location(), c.checker)
	closureCompiler.Parent = c
	closureCompiler.Errors = c.Errors
	closureCompiler.compileFunction(node.Location(), nil, node.Body)

	result := closureCompiler.Bytecode
	c.emitValue(value.Ref(result), node.Location())

	c.emit(node.Location().StartPos.Line, bytecode.CLOSURE)

	for _, upvalue := range closureCompiler.upvalues {
		var flags bitfield.BitField8
		if upvalue.isLocal {
			flags.SetFlag(vm.UpvalueLocalFlag)
		}
		if upvalue.upIndex > math.MaxUint8 {
			flags.SetFlag(vm.UpvalueLongIndexFlag)
		}
		c.emitByte(flags.Byte())

		if flags.HasFlag(vm.UpvalueLongIndexFlag) {
			c.emitUint16(upvalue.upIndex)
		} else {
			c.emitByte(byte(upvalue.upIndex))
		}
	}

	c.emitByte(vm.ClosureTerminatorFlag)

	c.emit(node.Location().StartPos.Line, bytecode.GO)
}

func (c *Compiler) compileClosureLiteralNode(node *ast.ClosureLiteralNode) {
	closureCompiler := New("<closure>", methodMode, node.Location(), c.checker)
	closureCompiler.Parent = c
	closureCompiler.Errors = c.Errors
	closureCompiler.compileFunction(node.Location(), node.Parameters, node.Body)

	result := closureCompiler.Bytecode
	c.emitValue(value.Ref(result), node.Location())

	c.emit(node.Location().StartPos.Line, bytecode.CLOSURE)

	for _, upvalue := range closureCompiler.upvalues {
		var flags bitfield.BitField8
		if upvalue.isLocal {
			flags.SetFlag(vm.UpvalueLocalFlag)
		}
		if upvalue.upIndex > math.MaxUint8 {
			flags.SetFlag(vm.UpvalueLongIndexFlag)
		}
		c.emitByte(flags.Byte())

		if flags.HasFlag(vm.UpvalueLongIndexFlag) {
			c.emitUint16(upvalue.upIndex)
		} else {
			c.emitByte(byte(upvalue.upIndex))
		}
	}

	c.emitByte(vm.ClosureTerminatorFlag)
}

func (c *Compiler) mixinIsCompilable(node *ast.MixinDeclarationNode) bool {
	if len(node.Body) <= 0 {
		return false
	}

	mixinType := c.typeOf(node).(*types.Mixin)

	mixinCompiler := New(fmt.Sprintf("<mixin: %s>", mixinType.Name()), namespaceMode, node.Location(), c.checker)
	mixinCompiler.Errors = c.Errors
	if !mixinCompiler.compileNamespace(node) {
		return false
	}

	node.Bytecode = mixinCompiler.Bytecode
	return true
}

func (c *Compiler) compileMixinDeclarationNode(node *ast.MixinDeclarationNode) expressionResult {
	if node.Bytecode == nil {
		return expressionIgnored
	}

	location := node.Location()
	mixinType := c.typeOf(node).(*types.Mixin)
	mixinName := value.ToSymbol(mixinType.Name())

	c.emitGetConst(mixinName, node.Constant.Location())
	c.emitValue(value.Ref(node.Bytecode), location)
	c.emit(location.StartPos.Line, bytecode.INIT_NAMESPACE)
	return expressionCompiled
}

func (c *Compiler) moduleIsCompilable(node *ast.ModuleDeclarationNode) bool {
	if len(node.Body) <= 0 {
		return false
	}

	modType := c.typeOf(node).(*types.Module)
	modCompiler := New(fmt.Sprintf("<module: %s>", modType.Name()), namespaceMode, node.Location(), c.checker)
	modCompiler.Errors = c.Errors
	if !modCompiler.compileNamespace(node) {
		return false
	}
	node.Bytecode = modCompiler.Bytecode
	return true
}

func (c *Compiler) compileModuleDeclarationNode(node *ast.ModuleDeclarationNode) expressionResult {
	if node.Bytecode == nil {
		return expressionIgnored
	}

	location := node.Location()
	modType := c.typeOf(node).(*types.Module)
	modName := value.ToSymbol(modType.Name())

	c.emitGetConst(modName, node.Constant.Location())
	c.emitValue(value.Ref(node.Bytecode), location)
	c.emit(location.StartPos.Line, bytecode.INIT_NAMESPACE)
	return expressionCompiled
}

func (c *Compiler) interfaceIsCompilable(node *ast.InterfaceDeclarationNode) bool {
	if len(node.Body) <= 0 {
		return false
	}

	ifaceType := c.typeOf(node).(*types.Interface)

	ifaceCompiler := New(fmt.Sprintf("<interface: %s>", ifaceType.Name()), namespaceMode, node.Location(), c.checker)
	ifaceCompiler.Errors = c.Errors
	if !ifaceCompiler.compileNamespace(node) {
		return false
	}
	node.Bytecode = ifaceCompiler.Bytecode
	return true
}

func (c *Compiler) compileInterfaceDeclarationNode(node *ast.InterfaceDeclarationNode) expressionResult {
	if node.Bytecode == nil {
		return expressionIgnored
	}

	location := node.Location()
	ifaceType := c.typeOf(node).(*types.Interface)
	className := value.ToSymbol(ifaceType.Name())

	c.emitGetConst(className, node.Constant.Location())
	c.emitValue(value.Ref(node.Bytecode), location)
	c.emit(location.StartPos.Line, bytecode.INIT_NAMESPACE)
	return expressionCompiled
}

func (c *Compiler) classIsCompilable(node *ast.ClassDeclarationNode) bool {
	if len(node.Body) <= 0 {
		return false
	}

	classType := c.typeOf(node).(*types.Class)

	classCompiler := New(fmt.Sprintf("<class: %s>", classType.Name()), namespaceMode, node.Location(), c.checker)
	classCompiler.Errors = c.Errors
	if !classCompiler.compileNamespace(node) {
		return false
	}
	node.Bytecode = classCompiler.Bytecode
	return true
}

func (c *Compiler) compileClassDeclarationNode(node *ast.ClassDeclarationNode) expressionResult {
	if node.Bytecode == nil {
		return expressionIgnored
	}

	location := node.Location()
	classType := c.typeOf(node).(*types.Class)
	className := value.ToSymbol(classType.Name())

	c.emitGetConst(className, node.Constant.Location())
	c.emitValue(value.Ref(node.Bytecode), location)
	c.emit(location.StartPos.Line, bytecode.INIT_NAMESPACE)
	return expressionCompiled
}

func (c *Compiler) compileValuePatternDeclarationNode(node *ast.ValuePatternDeclarationNode) {
	previousMode := c.mode
	c.mode = valuePatternDeclarationNode
	defer func() {
		c.mode = previousMode
	}()

	location := node.Location()
	c.compileNodeWithResult(node.Initialiser)
	c.pattern(node.Pattern)

	jumpOverErrorOffset := c.emitJump(location.StartPos.Line, bytecode.JUMP_IF)

	c.emitValue(
		value.Ref(value.NewError(
			value.PatternNotMatchedErrorClass,
			"assigned value does not match the pattern defined in value declaration",
		)),
		location,
	)
	c.emit(location.EndPos.Line, bytecode.THROW)

	c.patchJump(jumpOverErrorOffset, location)
}

func (c *Compiler) compilerVariablePatternDeclarationNode(node *ast.VariablePatternDeclarationNode) {
	location := node.Location()
	c.compileNodeWithResult(node.Initialiser)
	c.pattern(node.Pattern)

	jumpOverErrorOffset := c.emitJump(location.StartPos.Line, bytecode.JUMP_IF)

	c.emitValue(
		value.Ref(value.NewError(
			value.PatternNotMatchedErrorClass,
			"assigned value does not match the pattern defined in variable declaration",
		)),
		location,
	)
	c.emit(location.EndPos.Line, bytecode.THROW)

	c.patchJump(jumpOverErrorOffset, location)
}

func (c *Compiler) compileVariableDeclarationNode(node *ast.VariableDeclarationNode, valueIsIgnored bool) expressionResult {
	initialised := node.Initialiser != nil

	local := c.defineLocal(identifierToName(node.Name), node.Location())
	if local == nil {
		return valueIgnoredToResult(valueIsIgnored)
	}
	if initialised {
		c.compileNodeWithResult(node.Initialiser)
	}

	if initialised {
		return c.emitSetLocal(node.Location().StartPos.Line, local.index, valueIsIgnored)
	}

	if !valueIsIgnored {
		c.emit(node.Location().StartPos.Line, bytecode.NIL)
		return expressionCompiled
	}

	return expressionCompiledWithoutResult
}

// Compile each element of a collection of statements.
func (c *Compiler) compileStatements(collection []ast.StatementNode, location *position.Location, valueIsIgnored bool) expressionResult {
	if valueIsIgnored {
		c.compileStatementsWithoutResult(collection)
		return expressionCompiledWithoutResult
	}

	c.compileStatementsWithResult(collection, location)
	return expressionCompiled
}

// Compiles a list of statements leaving no value on the stack
func (c *Compiler) compileStatementsWithoutResult(collection []ast.StatementNode) {
	for _, s := range collection {
		result := c.compileNode(s, true)
		switch result {
		case expressionCompiled:
			c.emit(s.Location().EndPos.Line, bytecode.POP)
		}
	}
}

// Compiles a list of statements leaving the value produced by the last statement on the stack
func (c *Compiler) compileStatementsWithResult(collection []ast.StatementNode, location *position.Location) {
	if !c.compileStatementsOk(collection) {
		c.emit(location.EndPos.Line, bytecode.NIL)
		return
	}
}

// Compiles a list of statements leaving the value produced by the last statement on the stack
func (c *Compiler) compileStatementsOk(collection []ast.StatementNode) bool {
	lastCompilableIndex := -1
	for i, s := range collection {
		if c.nodeIsCompilable(s) {
			lastCompilableIndex = i
		}
	}

	if lastCompilableIndex == -1 {
		return false
	}

	for i, s := range collection {
		isLast := lastCompilableIndex == i
		result := c.compileNode(s, !isLast)
		switch result {
		case expressionCompiled:
			if !isLast {
				c.emit(s.Location().EndPos.Line, bytecode.POP)
			}
		}
	}

	return true
}

func (c *Compiler) removeOpcode() {
	c.lastOpCode = c.secondToLastOpCode
	c.secondToLastOpCode = bytecode.NOOP
	c.Bytecode.RemoveByte()
}

func (c *Compiler) compileRangeLiteralNode(node *ast.RangeLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	location := node.Location()

	if node.Start == nil {
		c.compileNodeWithResult(node.End)

		switch node.Op.Type {
		case token.CLOSED_RANGE_OP, token.LEFT_OPEN_RANGE_OP:
			c.emit(location.StartPos.Line, bytecode.NEW_RANGE, bytecode.BEGINLESS_CLOSED_RANGE_FLAG)
		case token.RIGHT_OPEN_RANGE_OP, token.OPEN_RANGE_OP:
			c.emit(location.StartPos.Line, bytecode.NEW_RANGE, bytecode.BEGINLESS_OPEN_RANGE_FLAG)
		default:
			panic(fmt.Sprintf("invalid range operator: %#v", node.Op))
		}

		return
	}
	if node.End == nil {
		c.compileNodeWithResult(node.Start)

		switch node.Op.Type {
		case token.CLOSED_RANGE_OP, token.RIGHT_OPEN_RANGE_OP:
			c.emit(location.StartPos.Line, bytecode.NEW_RANGE, bytecode.ENDLESS_CLOSED_RANGE_FLAG)
		case token.LEFT_OPEN_RANGE_OP, token.OPEN_RANGE_OP:
			c.emit(location.StartPos.Line, bytecode.NEW_RANGE, bytecode.ENDLESS_OPEN_RANGE_FLAG)
		default:
			panic(fmt.Sprintf("invalid range operator: %#v", node.Op))
		}

		return
	}

	c.compileNodeWithResult(node.Start)
	c.compileNodeWithResult(node.End)
	switch node.Op.Type {
	case token.CLOSED_RANGE_OP:
		c.emit(location.StartPos.Line, bytecode.NEW_RANGE, bytecode.CLOSED_RANGE_FLAG)
	case token.OPEN_RANGE_OP:
		c.emit(location.StartPos.Line, bytecode.NEW_RANGE, bytecode.OPEN_RANGE_FLAG)
	case token.LEFT_OPEN_RANGE_OP:
		c.emit(location.StartPos.Line, bytecode.NEW_RANGE, bytecode.LEFT_OPEN_RANGE_FLAG)
	case token.RIGHT_OPEN_RANGE_OP:
		c.emit(location.StartPos.Line, bytecode.NEW_RANGE, bytecode.RIGHT_OPEN_RANGE_FLAG)
	default:
		panic(fmt.Sprintf("invalid range operator: %#v", node.Op))
	}
}

func (c *Compiler) compileHashSetLiteralNode(node *ast.HashSetLiteralNode) {
	location := node.Location()
	if c.resolveAndEmit(node) {
		return
	}

	baseSet := value.NewHashSet(len(node.Elements))
	firstDynamicIndex := -1

	for i, elementNode := range node.Elements {
		element := resolve(elementNode)
		if element.IsUndefined() || value.IsMutableCollection(element) {
			firstDynamicIndex = i
			break
		}

		vm.HashSetAppendWithMaxLoad(nil, baseSet, element, 1)
	}

	if node.Capacity == nil {
		c.emit(location.StartPos.Line, bytecode.UNDEFINED)
	} else {
		c.compileNodeWithResult(node.Capacity)
	}

	if baseSet.Length() == 0 && baseSet.Capacity() == 0 {
		c.emit(location.StartPos.Line, bytecode.UNDEFINED)
	} else {
		c.emitLoadValue(value.Ref(baseSet), location)
	}

	firstModifierElementIndex := -1
	var dynamicElementNodes []ast.ExpressionNode

	if firstDynamicIndex != -1 {
		dynamicElementNodes = node.Elements[firstDynamicIndex:]
	dynamicElementsLoop:
		for i, elementNode := range dynamicElementNodes {
			switch elementNode.(type) {
			case *ast.ModifierNode, *ast.ModifierForInNode, *ast.ModifierIfElseNode:
				if node.Capacity != nil {
					c.Errors.AddFailure(
						"capacity cannot be specified in collection literals with conditional elements or loops",
						node.Capacity.Location(),
					)
					return
				}
				c.emitNewHashSet(i, location)
				firstModifierElementIndex = i
				break dynamicElementsLoop
			default:
				c.compileNodeWithResult(elementNode)
			}
		}
	}

	if firstModifierElementIndex != -1 {
		modifierElementNodes := dynamicElementNodes[firstModifierElementIndex:]
		for _, elementNode := range modifierElementNodes {
			switch e := elementNode.(type) {
			case *ast.ModifierNode:
				var jumpOp bytecode.OpCode
				switch e.Modifier.Type {
				case token.IF:
					jumpOp = bytecode.JUMP_UNLESS
				case token.UNLESS:
					jumpOp = bytecode.JUMP_IF
				default:
					panic(fmt.Sprintf("invalid collection modifier: %#v", e.Modifier))
				}

				c.compileIfWithConditionExpression(
					jumpOp,
					e.Right,
					func() {
						c.compileNodeWithResult(e.Left)
						c.emit(e.Location().StartPos.Line, bytecode.APPEND)
					},
					func() {},
					e.Location(),
					false,
				)
			case *ast.ModifierIfElseNode:
				c.compileIfWithConditionExpression(
					bytecode.JUMP_UNLESS,
					e.Condition,
					func() {
						c.compileNodeWithResult(e.ThenExpression)
						c.emit(e.Location().StartPos.Line, bytecode.APPEND)
					},
					func() {
						c.compileNodeWithResult(e.ElseExpression)
						c.emit(e.Location().StartPos.Line, bytecode.APPEND)
					},
					e.Location(),
					false,
				)
			case *ast.ModifierForInNode:
				c.compileForIn(
					"",
					e.Pattern,
					e.InExpression,
					func() {
						c.compileNodeWithResult(e.ThenExpression)
						c.emit(e.ThenExpression.Location().EndPos.Line, bytecode.APPEND)
					},
					e.Location(),
					true,
				)
			default:
				c.compileNodeWithResult(elementNode)
				c.emit(e.Location().StartPos.Line, bytecode.APPEND)
			}
		}

		return
	}

	c.emitNewHashSet(len(dynamicElementNodes), location)
}

func (c *Compiler) compileHashMapLiteralNode(node *ast.HashMapLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	location := node.Location()
	baseMap := value.NewHashMap(len(node.Elements))
	firstDynamicIndex := -1

elementLoop:
	for i, elementNode := range node.Elements {
	elementSwitch:
		switch e := elementNode.(type) {
		case *ast.KeyValueExpressionNode:
			if !e.IsStatic() {
				break elementSwitch
			}
			key := resolve(e.Key)
			val := resolve(e.Value)
			if value.IsMutableCollection(key) || value.IsMutableCollection(val) {
				break elementSwitch
			}

			vm.HashMapSetWithMaxLoad(nil, baseMap, key, val, 1)
			continue elementLoop
		case *ast.SymbolKeyValueExpressionNode:
			if !e.IsStatic() {
				break elementSwitch
			}
			key := value.ToSymbol(identifierToName(e.Key))
			val := resolve(e.Value)
			if val.IsUndefined() || value.IsMutableCollection(val) {
				break elementSwitch
			}

			vm.HashMapSetWithMaxLoad(nil, baseMap, key.ToValue(), val, 1)
			continue elementLoop
		}

		firstDynamicIndex = i
		break elementLoop
	}

	if node.Capacity == nil {
		c.emit(location.StartPos.Line, bytecode.UNDEFINED)
	} else {
		c.compileNodeWithResult(node.Capacity)
	}

	if baseMap.Length() == 0 && baseMap.Capacity() == 0 {
		c.emit(location.StartPos.Line, bytecode.UNDEFINED)
	} else {
		c.emitLoadValue(value.Ref(baseMap), location)
	}

	firstModifierElementIndex := -1
	var dynamicElementNodes []ast.ExpressionNode

	if firstDynamicIndex != -1 {
		dynamicElementNodes = node.Elements[firstDynamicIndex:]
	dynamicElementsLoop:
		for i, elementNode := range dynamicElementNodes {
			switch element := elementNode.(type) {
			case *ast.ModifierNode, *ast.ModifierForInNode, *ast.ModifierIfElseNode:
				if node.Capacity != nil {
					c.Errors.AddFailure(
						"capacity cannot be specified in collection literals with conditional elements or loops",
						node.Capacity.Location(),
					)
					return
				}
				c.emitNewHashMap(i, location)
				firstModifierElementIndex = i
				break dynamicElementsLoop
			case *ast.KeyValueExpressionNode:
				c.compileNodeWithResult(element.Key)
				c.compileNodeWithResult(element.Value)
			case *ast.SymbolKeyValueExpressionNode:
				c.emitValue(value.ToSymbol(identifierToName(element.Key)).ToValue(), element.Location())
				c.compileNodeWithResult(element.Value)
			case *ast.PublicIdentifierNode:
				c.emitValue(value.ToSymbol(element.Value).ToValue(), element.Location())
				c.compileLocalVariableAccess(element.Value, element.Location())
			case *ast.PrivateIdentifierNode:
				c.emitValue(value.ToSymbol(element.Value).ToValue(), element.Location())
				c.compileLocalVariableAccess(element.Value, element.Location())
			default:
				panic(fmt.Sprintf("invalid element in hashmap literal: %#v", elementNode))
			}
		}
	}

	if firstModifierElementIndex != -1 {
		modifierElementNodes := dynamicElementNodes[firstModifierElementIndex:]
		for _, elementNode := range modifierElementNodes {
			switch e := elementNode.(type) {
			case *ast.KeyValueExpressionNode:
				c.compileNodeWithResult(e.Key)
				c.compileNodeWithResult(e.Value)
				c.emit(e.Location().StartPos.Line, bytecode.MAP_SET)
			case *ast.SymbolKeyValueExpressionNode:
				c.emitValue(value.ToSymbol(identifierToName(e.Key)).ToValue(), e.Location())
				c.compileNodeWithResult(e.Value)
				c.emit(e.Location().StartPos.Line, bytecode.MAP_SET)
			case *ast.ModifierNode:
				var jumpOp bytecode.OpCode
				switch e.Modifier.Type {
				case token.IF:
					jumpOp = bytecode.JUMP_UNLESS
				case token.UNLESS:
					jumpOp = bytecode.JUMP_IF
				default:
					panic(fmt.Sprintf("invalid collection modifier: %#v", e.Modifier))
				}

				c.compileIfWithConditionExpression(
					jumpOp,
					e.Right,
					func() {
						switch then := e.Left.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNodeWithResult(then.Key)
							c.compileNodeWithResult(then.Value)
							c.emit(then.Location().StartPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(identifierToName(then.Key)).ToValue(), then.Location())
							c.compileNodeWithResult(then.Value)
							c.emit(then.Location().StartPos.Line, bytecode.MAP_SET)
						default:
							panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
						}
					},
					func() {},
					e.Location(),
					false,
				)
			case *ast.ModifierIfElseNode:
				c.compileIfWithConditionExpression(
					bytecode.JUMP_UNLESS,
					e.Condition,
					func() {
						switch then := e.ThenExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNodeWithResult(then.Key)
							c.compileNodeWithResult(then.Value)
							c.emit(then.Location().StartPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(identifierToName(then.Key)).ToValue(), then.Location())
							c.compileNodeWithResult(then.Value)
							c.emit(then.Location().StartPos.Line, bytecode.MAP_SET)
						default:
							panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
						}
					},
					func() {
						switch els := e.ElseExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNodeWithResult(els.Key)
							c.compileNodeWithResult(els.Value)
							c.emit(els.Location().EndPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(identifierToName(els.Key)).ToValue(), els.Location())
							c.compileNodeWithResult(els.Value)
							c.emit(els.Location().EndPos.Line, bytecode.MAP_SET)
						default:
							panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
						}
					},
					e.Location(),
					false,
				)
			case *ast.ModifierForInNode:
				c.compileForIn(
					"",
					e.Pattern,
					e.InExpression,
					func() {
						switch then := e.ThenExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNodeWithResult(then.Key)
							c.compileNodeWithResult(then.Value)
							c.emit(then.Location().EndPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(identifierToName(then.Key)).ToValue(), then.Location())
							c.compileNodeWithResult(then.Value)
							c.emit(then.Location().EndPos.Line, bytecode.MAP_SET)
						default:
							panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
						}
					},
					e.Location(),
					true,
				)
			default:
				panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
			}
		}

		return
	}

	c.emitNewHashMap(len(dynamicElementNodes), location)
}

func (c *Compiler) compileHashRecordLiteralNode(node *ast.HashRecordLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	location := node.Location()
	baseRecord := value.NewHashRecord(len(node.Elements))
	firstDynamicIndex := -1

elementLoop:
	for i, elementNode := range node.Elements {
	elementSwitch:
		switch e := elementNode.(type) {
		case *ast.KeyValueExpressionNode:
			if !e.IsStatic() {
				break elementSwitch
			}
			key := resolve(e.Key)
			val := resolve(e.Value)
			if value.IsMutableCollection(key) || value.IsMutableCollection(val) {
				break elementSwitch
			}

			vm.HashRecordSetWithMaxLoad(nil, baseRecord, key, val, 1)
			continue elementLoop
		case *ast.SymbolKeyValueExpressionNode:
			if !e.IsStatic() {
				break elementSwitch
			}
			key := value.ToSymbol(identifierToName(e.Key))
			val := resolve(e.Value)
			if val.IsUndefined() || value.IsMutableCollection(val) {
				break elementSwitch
			}

			vm.HashRecordSetWithMaxLoad(nil, baseRecord, key.ToValue(), val, 1)
			continue elementLoop
		}

		firstDynamicIndex = i
		break elementLoop
	}

	if baseRecord.Length() == 0 {
		c.emit(location.StartPos.Line, bytecode.UNDEFINED)
	} else {
		c.emitLoadValue(value.Ref(baseRecord), location)
	}

	firstModifierElementIndex := -1
	var dynamicElementNodes []ast.ExpressionNode

	if firstDynamicIndex != -1 {
		dynamicElementNodes = node.Elements[firstDynamicIndex:]
	dynamicElementsLoop:
		for i, elementNode := range dynamicElementNodes {
			switch element := elementNode.(type) {
			case *ast.ModifierNode, *ast.ModifierForInNode, *ast.ModifierIfElseNode:
				c.emitNewHashRecord(i, location)
				firstModifierElementIndex = i
				break dynamicElementsLoop
			case *ast.KeyValueExpressionNode:
				c.compileNodeWithResult(element.Key)
				c.compileNodeWithResult(element.Value)
			case *ast.SymbolKeyValueExpressionNode:
				c.emitValue(value.ToSymbol(identifierToName(element.Key)).ToValue(), element.Location())
				c.compileNodeWithResult(element.Value)
			case *ast.PublicIdentifierNode:
				c.emitValue(value.ToSymbol(element.Value).ToValue(), element.Location())
				c.compileLocalVariableAccess(element.Value, element.Location())
			case *ast.PrivateIdentifierNode:
				c.emitValue(value.ToSymbol(element.Value).ToValue(), element.Location())
				c.compileLocalVariableAccess(element.Value, element.Location())
			default:
				panic(fmt.Sprintf("invalid element in hashmap literal: %#v", elementNode))
			}
		}
	}

	if firstModifierElementIndex != -1 {
		modifierElementNodes := dynamicElementNodes[firstModifierElementIndex:]
		for _, elementNode := range modifierElementNodes {
			switch e := elementNode.(type) {
			case *ast.KeyValueExpressionNode:
				c.compileNodeWithResult(e.Key)
				c.compileNodeWithResult(e.Value)
				c.emit(e.Location().StartPos.Line, bytecode.MAP_SET)
			case *ast.SymbolKeyValueExpressionNode:
				c.emitValue(value.ToSymbol(identifierToName(e.Key)).ToValue(), e.Location())
				c.compileNodeWithResult(e.Value)
				c.emit(e.Location().StartPos.Line, bytecode.MAP_SET)
			case *ast.ModifierNode:
				var jumpOp bytecode.OpCode
				switch e.Modifier.Type {
				case token.IF:
					jumpOp = bytecode.JUMP_UNLESS
				case token.UNLESS:
					jumpOp = bytecode.JUMP_IF
				default:
					panic(fmt.Sprintf("invalid collection modifier: %#v", e.Modifier))
				}

				c.compileIfWithConditionExpression(
					jumpOp,
					e.Right,
					func() {
						switch then := e.Left.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNodeWithResult(then.Key)
							c.compileNodeWithResult(then.Value)
							c.emit(then.Location().StartPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(identifierToName(then.Key)).ToValue(), then.Location())
							c.compileNodeWithResult(then.Value)
							c.emit(then.Location().StartPos.Line, bytecode.MAP_SET)
						default:
							panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
						}
					},
					func() {},
					e.Location(),
					false,
				)
			case *ast.ModifierIfElseNode:
				c.compileIfWithConditionExpression(
					bytecode.JUMP_UNLESS,
					e.Condition,
					func() {
						switch then := e.ThenExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNodeWithResult(then.Key)
							c.compileNodeWithResult(then.Value)
							c.emit(then.Location().StartPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(identifierToName(then.Key)).ToValue(), then.Location())
							c.compileNodeWithResult(then.Value)
							c.emit(then.Location().StartPos.Line, bytecode.MAP_SET)
						default:
							panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
						}
					},
					func() {
						switch els := e.ElseExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNodeWithResult(els.Key)
							c.compileNodeWithResult(els.Value)
							c.emit(els.Location().EndPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(identifierToName(els.Key)).ToValue(), els.Location())
							c.compileNodeWithResult(els.Value)
							c.emit(els.Location().EndPos.Line, bytecode.MAP_SET)
						default:
							panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
						}
					},
					e.Location(),
					false,
				)
			case *ast.ModifierForInNode:
				c.compileForIn(
					"",
					e.Pattern,
					e.InExpression,
					func() {
						switch then := e.ThenExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNodeWithResult(then.Key)
							c.compileNodeWithResult(then.Value)
							c.emit(then.Location().EndPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(identifierToName(then.Key)).ToValue(), then.Location())
							c.compileNodeWithResult(then.Value)
							c.emit(then.Location().EndPos.Line, bytecode.MAP_SET)
						default:
							panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
						}
					},
					e.Location(),
					true,
				)
			default:
				panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
			}
		}

		return
	}

	c.emitNewHashRecord(len(dynamicElementNodes), location)
}

func (c *Compiler) compileArrayListLiteralNode(node *ast.ArrayListLiteralNode) {
	location := node.Location()
	if c.resolveAndEmitList(node) {
		return
	}

	var keyValueCount int
	for _, elementNode := range node.Elements {
		switch elementNode.(type) {
		case *ast.KeyValueExpressionNode:
			keyValueCount++
		}
	}
	baseList := make(value.ArrayList, 0, len(node.Elements)-keyValueCount)
	firstDynamicIndex := -1

elementLoop:
	for i, elementNode := range node.Elements {
	elementSwitch:
		switch e := elementNode.(type) {
		case *ast.KeyValueExpressionNode:
			if !e.IsStatic() {
				break elementSwitch
			}
			key := resolve(e.Key)
			val := resolve(e.Value)
			index, ok := value.ToGoInt(key)
			if !ok {
				break elementSwitch
			}

			baseList.Expand((index + 1) - len(baseList))
			baseList[index] = val
			continue elementLoop
		}

		element := resolve(elementNode)
		if element.IsUndefined() || value.IsMutableCollection(element) {
			firstDynamicIndex = i
			break
		}

		baseList = append(baseList, element)
	}

	if node.Capacity == nil {
		c.emit(location.StartPos.Line, bytecode.UNDEFINED)
	} else {
		c.compileNodeWithResult(node.Capacity)
	}

	if len(baseList) == 0 && (keyValueCount == 0 || cap(baseList) == 0) {
		c.emit(location.StartPos.Line, bytecode.UNDEFINED)
	} else {
		c.emitLoadValue(value.Ref(&baseList), location)
	}

	firstModifierElementIndex := -1
	var dynamicElementNodes []ast.ExpressionNode

	if firstDynamicIndex != -1 {
		dynamicElementNodes = node.Elements[firstDynamicIndex:]
	dynamicElementsLoop:
		for i, elementNode := range dynamicElementNodes {
			switch elementNode.(type) {
			case *ast.ModifierNode, *ast.ModifierForInNode, *ast.ModifierIfElseNode:
				if node.Capacity != nil {
					c.Errors.AddFailure(
						"capacity cannot be specified in collection literals with conditional elements or loops",
						node.Capacity.Location(),
					)
					return
				}
				c.emitNewArrayList(i, location)
				firstModifierElementIndex = i
				break dynamicElementsLoop
			case *ast.KeyValueExpressionNode:
				c.emitNewArrayList(i, location)
				firstModifierElementIndex = i
				break dynamicElementsLoop
			default:
				c.compileNodeWithResult(elementNode)
			}
		}
	}

	if firstModifierElementIndex != -1 {
		modifierElementNodes := dynamicElementNodes[firstModifierElementIndex:]
		for _, elementNode := range modifierElementNodes {
			switch e := elementNode.(type) {
			case *ast.KeyValueExpressionNode:
				c.compileNodeWithResult(e.Key)
				c.compileNodeWithResult(e.Value)
				c.emit(e.Location().StartPos.Line, bytecode.APPEND_AT)
			case *ast.ModifierNode:
				var jumpOp bytecode.OpCode
				switch e.Modifier.Type {
				case token.IF:
					jumpOp = bytecode.JUMP_UNLESS
				case token.UNLESS:
					jumpOp = bytecode.JUMP_IF
				default:
					panic(fmt.Sprintf("invalid collection modifier: %#v", e.Modifier))
				}

				c.compileIfWithConditionExpression(
					jumpOp,
					e.Right,
					func() {
						switch then := e.Left.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNodeWithResult(then.Key)
							c.compileNodeWithResult(then.Value)
							c.emit(then.Location().StartPos.Line, bytecode.APPEND_AT)
						default:
							c.compileNodeWithResult(e.Left)
							c.emit(e.Location().StartPos.Line, bytecode.APPEND)
						}
					},
					func() {},
					e.Location(),
					false,
				)
			case *ast.ModifierIfElseNode:
				c.compileIfWithConditionExpression(
					bytecode.JUMP_UNLESS,
					e.Condition,
					func() {
						switch then := e.ThenExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNodeWithResult(then.Key)
							c.compileNodeWithResult(then.Value)
							c.emit(then.Location().StartPos.Line, bytecode.APPEND_AT)
						default:
							c.compileNodeWithResult(e.ThenExpression)
							c.emit(e.Location().StartPos.Line, bytecode.APPEND)
						}
					},
					func() {
						switch els := e.ElseExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNodeWithResult(els.Key)
							c.compileNodeWithResult(els.Value)
							c.emit(els.Location().StartPos.Line, bytecode.APPEND_AT)
						default:
							c.compileNodeWithResult(e.ElseExpression)
							c.emit(e.Location().StartPos.Line, bytecode.APPEND)
						}
					},
					e.Location(),
					false,
				)
			case *ast.ModifierForInNode:
				c.compileForIn(
					"",
					e.Pattern,
					e.InExpression,
					func() {
						switch then := e.ThenExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNodeWithResult(then.Key)
							c.compileNodeWithResult(then.Value)
							c.emit(then.Location().EndPos.Line, bytecode.APPEND_AT)
						default:
							c.compileNodeWithResult(e.ThenExpression)
							c.emit(then.Location().EndPos.Line, bytecode.APPEND)
						}
					},
					e.Location(),
					true,
				)
			default:
				c.compileNodeWithResult(elementNode)
				c.emit(e.Location().StartPos.Line, bytecode.APPEND)
			}
		}

		return
	}

	c.emitNewArrayList(len(dynamicElementNodes), location)
}

func (c *Compiler) compileArrayTupleLiteralNode(node *ast.ArrayTupleLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	location := node.Location()

	var baseArrayTuple value.ArrayTuple
	firstDynamicIndex := -1

elementLoop:
	for i, elementNode := range node.Elements {
	elementSwitch:
		switch e := elementNode.(type) {
		case *ast.KeyValueExpressionNode:
			if !e.IsStatic() {
				break elementSwitch
			}
			key := resolve(e.Key)
			val := resolve(e.Value)
			index, ok := value.ToGoInt(key)
			if !ok {
				break elementSwitch
			}

			baseArrayTuple.Expand((index + 1) - len(baseArrayTuple))
			baseArrayTuple[index] = val
			continue elementLoop
		}

		element := resolve(elementNode)
		if element.IsUndefined() {
			firstDynamicIndex = i
			break
		}

		baseArrayTuple = append(baseArrayTuple, element)
	}

	if len(baseArrayTuple) == 0 {
		c.emit(location.StartPos.Line, bytecode.UNDEFINED)
	} else {
		c.emitLoadValue(value.Ref(&baseArrayTuple), location)
	}

	firstModifierElementIndex := -1
	var dynamicElementNodes []ast.ExpressionNode

	if firstDynamicIndex != -1 {
		dynamicElementNodes = node.Elements[firstDynamicIndex:]
	dynamicElementsLoop:
		for i, elementNode := range dynamicElementNodes {
			switch e := elementNode.(type) {
			case *ast.ModifierNode, *ast.ModifierForInNode, *ast.ModifierIfElseNode, *ast.KeyValueExpressionNode:
				if i == 0 && firstDynamicIndex != 0 {
					c.emit(e.Location().StartPos.Line, bytecode.COPY)
				} else {
					c.emitNewArrayTuple(i, location)
				}
				firstModifierElementIndex = i
				break dynamicElementsLoop
			default:
				c.compileNodeWithResult(elementNode)
			}
		}
	}

	if firstModifierElementIndex != -1 {
		modifierElementNodes := dynamicElementNodes[firstModifierElementIndex:]
		for _, elementNode := range modifierElementNodes {
			switch e := elementNode.(type) {
			case *ast.KeyValueExpressionNode:
				c.compileNodeWithResult(e.Key)
				c.compileNodeWithResult(e.Value)
				c.emit(e.Location().StartPos.Line, bytecode.APPEND_AT)
			case *ast.ModifierNode:
				var jumpOp bytecode.OpCode
				switch e.Modifier.Type {
				case token.IF:
					jumpOp = bytecode.JUMP_UNLESS
				case token.UNLESS:
					jumpOp = bytecode.JUMP_IF
				default:
					panic(fmt.Sprintf("invalid collection modifier: %#v", e.Modifier))
				}

				c.compileIfWithConditionExpression(
					jumpOp,
					e.Right,
					func() {
						switch then := e.Left.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNodeWithResult(then.Key)
							c.compileNodeWithResult(then.Value)
							c.emit(then.Location().StartPos.Line, bytecode.APPEND_AT)
						default:
							c.compileNodeWithResult(e.Left)
							c.emit(e.Location().StartPos.Line, bytecode.APPEND)
						}
					},
					func() {},
					e.Location(),
					false,
				)
			case *ast.ModifierIfElseNode:
				c.compileIfWithConditionExpression(
					bytecode.JUMP_UNLESS,
					e.Condition,
					func() {
						switch then := e.ThenExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNodeWithResult(then.Key)
							c.compileNodeWithResult(then.Value)
							c.emit(then.Location().EndPos.Line, bytecode.APPEND_AT)
						default:
							c.compileNodeWithResult(e.ThenExpression)
							c.emit(e.Location().EndPos.Line, bytecode.APPEND)
						}
					},
					func() {
						switch els := e.ElseExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNodeWithResult(els.Key)
							c.compileNodeWithResult(els.Value)
							c.emit(els.Location().EndPos.Line, bytecode.APPEND_AT)
						default:
							c.compileNodeWithResult(e.ElseExpression)
							c.emit(e.Location().EndPos.Line, bytecode.APPEND)
						}
					},
					e.Location(),
					false,
				)
			case *ast.ModifierForInNode:
				c.compileForIn(
					"",
					e.Pattern,
					e.InExpression,
					func() {
						switch then := e.ThenExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNodeWithResult(then.Key)
							c.compileNodeWithResult(then.Value)
							c.emit(then.Location().EndPos.Line, bytecode.APPEND_AT)
						default:
							c.compileNodeWithResult(e.ThenExpression)
							c.emit(then.Location().EndPos.Line, bytecode.APPEND)
						}
					},
					e.Location(),
					true,
				)
			default:
				c.compileNodeWithResult(elementNode)
				c.emit(e.Location().StartPos.Line, bytecode.APPEND)
			}
		}

		return
	}

	c.emitNewArrayTuple(len(dynamicElementNodes), location)
}

func (c *Compiler) compileWordArrayTupleLiteralNode(node *ast.WordArrayTupleLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	c.Errors.AddFailure("invalid word arrayTuple literal", node.Location())
}

func (c *Compiler) compileBinArrayTupleLiteralNode(node *ast.BinArrayTupleLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	c.Errors.AddFailure("invalid binary arrayTuple literal", node.Location())
}

func (c *Compiler) compileSymbolArrayTupleLiteralNode(node *ast.SymbolArrayTupleLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	c.Errors.AddFailure("invalid symbol arrayTuple literal", node.Location())
}

func (c *Compiler) compileHexArrayTupleLiteralNode(node *ast.HexArrayTupleLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	c.Errors.AddFailure("invalid hex arrayTuple literal", node.Location())
}

func (c *Compiler) compileWordArrayListLiteralNode(node *ast.WordArrayListLiteralNode) {
	list := resolve(node)
	location := node.Location()
	if list.IsUndefined() {
		c.Errors.AddFailure("invalid word arrayList literal", location)
		return
	}

	if node.Capacity == nil {
		c.emitLoadValue(list, location)
		c.emit(location.EndPos.Line, bytecode.COPY)
	} else {
		c.compileNodeWithResult(node.Capacity)
		c.emitLoadValue(list, location)
		c.emitNewArrayList(0, location)
	}
}

func (c *Compiler) compileBinArrayListLiteralNode(node *ast.BinArrayListLiteralNode) {
	list := resolve(node)
	location := node.Location()
	if list.IsUndefined() {
		c.Errors.AddFailure("invalid bin arrayList literal", location)
		return
	}

	if node.Capacity == nil {
		c.emitLoadValue(list, location)
		c.emit(location.EndPos.Line, bytecode.COPY)
	} else {
		c.compileNodeWithResult(node.Capacity)
		c.emitLoadValue(list, location)
		c.emitNewArrayList(0, location)
	}
}

func (c *Compiler) compileSymbolArrayListLiteralNode(node *ast.SymbolArrayListLiteralNode) {
	list := resolve(node)
	location := node.Location()
	if list.IsUndefined() {
		c.Errors.AddFailure("invalid symbol arrayList literal", location)
		return
	}

	if node.Capacity == nil {
		c.emitLoadValue(list, location)
		c.emit(location.EndPos.Line, bytecode.COPY)
	} else {
		c.compileNodeWithResult(node.Capacity)
		c.emitLoadValue(list, location)
		c.emitNewArrayList(0, location)
	}
}

func (c *Compiler) compileHexArrayListLiteralNode(node *ast.HexArrayListLiteralNode) {
	list := resolve(node)
	location := node.Location()
	if list.IsUndefined() {
		c.Errors.AddFailure("invalid hex arrayList literal", location)
		return
	}

	if node.Capacity == nil {
		c.emitLoadValue(list, location)
		c.emit(location.EndPos.Line, bytecode.COPY)
	} else {
		c.compileNodeWithResult(node.Capacity)
		c.emitLoadValue(list, location)
		c.emitNewArrayList(0, location)
	}
}

func (c *Compiler) compileWordHashSetLiteralNode(node *ast.WordHashSetLiteralNode) {
	list := resolve(node)
	location := node.Location()
	if list.IsUndefined() {
		c.Errors.AddFailure("invalid word hashSet literal", location)
		return
	}

	if node.Capacity == nil {
		c.emitLoadValue(list, location)
		c.emit(location.EndPos.Line, bytecode.COPY)
	} else {
		c.compileNodeWithResult(node.Capacity)
		c.emitLoadValue(list, location)
		c.emitNewHashSet(0, location)
	}
}

func (c *Compiler) compileBinHashSetLiteralNode(node *ast.BinHashSetLiteralNode) {
	list := resolve(node)
	location := node.Location()
	if list.IsUndefined() {
		c.Errors.AddFailure("invalid bin hashSet literal", location)
		return
	}

	if node.Capacity == nil {
		c.emitLoadValue(list, location)
		c.emit(location.EndPos.Line, bytecode.COPY)
	} else {
		c.compileNodeWithResult(node.Capacity)
		c.emitLoadValue(list, location)
		c.emitNewHashSet(0, location)
	}
}

func (c *Compiler) compileSymbolHashSetLiteralNode(node *ast.SymbolHashSetLiteralNode) {
	list := resolve(node)
	location := node.Location()
	if list.IsUndefined() {
		c.Errors.AddFailure("invalid symbol hashSet literal", location)
		return
	}

	if node.Capacity == nil {
		c.emitLoadValue(list, location)
		c.emit(location.EndPos.Line, bytecode.COPY)
	} else {
		c.compileNodeWithResult(node.Capacity)
		c.emitLoadValue(list, location)
		c.emitNewHashSet(0, location)
	}
}

func (c *Compiler) compileHexHashSetLiteralNode(node *ast.HexHashSetLiteralNode) {
	list := resolve(node)
	location := node.Location()
	if list.IsUndefined() {
		c.Errors.AddFailure("invalid hex hashSet literal", location)
		return
	}

	if node.Capacity == nil {
		c.emitLoadValue(list, location)
		c.emit(location.EndPos.Line, bytecode.COPY)
	} else {
		c.compileNodeWithResult(node.Capacity)
		c.emitLoadValue(list, location)
		c.emitNewHashSet(0, location)
	}
}

func (c *Compiler) emitNewHashSet(size int, location *position.Location) {
	c.emitNewCollection(bytecode.NEW_HASH_SET8, bytecode.NEW_HASH_SET16, size, location)
}

func (c *Compiler) emitNewArrayTuple(size int, location *position.Location) {
	c.emitNewCollection(bytecode.NEW_ARRAY_TUPLE8, bytecode.NEW_ARRAY_TUPLE16, size, location)
}

func (c *Compiler) emitNewArrayList(size int, location *position.Location) {
	c.emitNewCollection(bytecode.NEW_ARRAY_LIST8, bytecode.NEW_ARRAY_LIST16, size, location)
}

func (c *Compiler) emitNewHashMap(size int, location *position.Location) {
	c.emitNewCollection(bytecode.NEW_HASH_MAP8, bytecode.NEW_HASH_MAP16, size, location)
}

func (c *Compiler) emitNewHashRecord(size int, location *position.Location) {
	c.emitNewCollection(bytecode.NEW_HASH_RECORD8, bytecode.NEW_HASH_RECORD16, size, location)
}

func (c *Compiler) emitNewRegex(flags bitfield.BitField8, size int, location *position.Location) {
	if size <= math.MaxUint8 {
		c.emit(location.EndPos.Line, bytecode.NEW_REGEX8, flags.Byte(), byte(size))
		return
	}

	if size <= math.MaxUint16 {
		c.emit(location.EndPos.Line, bytecode.NEW_REGEX16)
		c.emitByte(flags.Byte())
		c.emitUint16(uint16(size))
		return
	}

	c.Errors.AddFailure(
		fmt.Sprintf("max number of regex literal elements reached: %d", math.MaxUint16),
		location,
	)
}

func (c *Compiler) emitNewCollection(opcode8, opcode16 bytecode.OpCode, size int, location *position.Location) {
	if size <= math.MaxUint8 {
		c.emit(location.EndPos.Line, opcode8, byte(size))
		return
	}

	if size <= math.MaxUint16 {
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, uint16(size))
		c.emit(location.EndPos.Line, opcode16, bytes...)
		return
	}

	c.Errors.AddFailure(
		fmt.Sprintf("max number of collection literal elements reached: %d", math.MaxUint16),
		location,
	)
}

func (c *Compiler) compileUninterpolatedRegexLiteralNode(node *ast.UninterpolatedRegexLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	re, err := value.CompileRegex(node.Content, node.Flags)
	if errList, ok := err.(diagnostic.DiagnosticList); ok {
		regexStartPos := node.Location().StartPos
		for _, err := range errList {
			errStartPos := err.Span.StartPos
			errEndPos := err.Span.EndPos

			columnDifference := regexStartPos.Column - 1 + 2 // add 2 to account for `%/`
			byteDifference := regexStartPos.ByteOffset + 2   // add 2 to account for `%/`
			lineDifference := regexStartPos.Line - 1

			if errStartPos.Line == 1 {
				errStartPos.Column += columnDifference
			}
			errStartPos.Line += lineDifference
			errStartPos.ByteOffset += byteDifference

			if errEndPos != errStartPos {
				if errEndPos.Line == 1 {
					errEndPos.Column += columnDifference
				}
				errEndPos.Line += lineDifference
				errEndPos.ByteOffset += byteDifference
			}
			err.Location.FilePath = c.Bytecode.Location.FilePath

			c.Errors.Append(err)
		}
		return
	}

	if err != nil {
		c.Errors.AddFailure(err.Error(), node.Location())
		return
	}

	c.emitValue(value.Ref(re), node.Location())
}

func (c *Compiler) compileInterpolatedRegexLiteralNode(node *ast.InterpolatedRegexLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	for _, elementNode := range node.Content {
		switch element := elementNode.(type) {
		case *ast.RegexLiteralContentSectionNode:
			c.emitValue(value.Ref(value.String(element.Value)), element.Location())
		case *ast.RegexInterpolationNode:
			c.compileNodeWithResult(element.Expression)
		}
	}
	c.emitNewRegex(node.Flags, len(node.Content), node.Location())
}

var inspectSymbol = value.ToSymbol("inspect")

func (c *Compiler) compileInterpolatedStringLiteralNode(node *ast.InterpolatedStringLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	for _, elementNode := range node.Content {
		switch element := elementNode.(type) {
		case *ast.StringLiteralContentSectionNode:
			c.emitValue(value.Ref(value.String(element.Value)), element.Location())
		case *ast.StringInterpolationNode:
			c.compileNodeWithResult(element.Expression)
		case *ast.StringInspectInterpolationNode:
			c.compileNodeWithResult(element.Expression)
			callInfo := value.NewCallSiteInfo(inspectSymbol, 0)
			c.emitCallMethod(callInfo, element.Location(), false)
		}
	}

	c.emitNewCollection(bytecode.NEW_STRING8, bytecode.NEW_STRING16, len(node.Content), node.Location())
}

func (c *Compiler) compileInterpolatedSymbolLiteralNode(node *ast.InterpolatedSymbolLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	for _, elementNode := range node.Content.Content {
		switch element := elementNode.(type) {
		case *ast.StringLiteralContentSectionNode:
			c.emitValue(value.Ref(value.String(element.Value)), element.Location())
		case *ast.StringInterpolationNode:
			c.compileNodeWithResult(element.Expression)
		}
	}

	c.emitNewCollection(bytecode.NEW_SYMBOL8, bytecode.NEW_SYMBOL16, len(node.Content.Content), node.Location())
}

func (c *Compiler) compileIntLiteralNode(node *ast.IntLiteralNode) {
	i, err := value.ParseBigInt(node.Value, 0)
	if !err.IsUndefined() {
		c.Errors.AddFailure(err.Error(), node.Location())
		return
	}
	if i.IsSmallInt() {
		c.emitValue(i.ToSmallInt().ToValue(), node.Location())
		return
	}
	c.emitValue(value.Ref(i), node.Location())
}

// Compiles boolean binary operators
func (c *Compiler) compileLogicalExpressionNode(node *ast.LogicalExpressionNode, valueIsIgnored bool) expressionResult {
	if r := resolve(node); !r.IsUndefined() {
		if valueIsIgnored {
			return expressionCompiledWithoutResult
		}
		c.emitValue(r, node.Location())
		return expressionCompiled
	}

	switch node.Op.Type {
	case token.AND_AND:
		return c.logicalAnd(node, valueIsIgnored)
	case token.OR_OR:
		return c.logicalOr(node, valueIsIgnored)
	case token.QUESTION_QUESTION:
		return c.nilCoalescing(node, valueIsIgnored)
	default:
		c.Errors.AddFailure(fmt.Sprintf("unknown logical operator: %s", node.Op.String()), node.Location())
	}

	return expressionCompiled
}

// Compiles the `??` operator
func (c *Compiler) nilCoalescing(node *ast.LogicalExpressionNode, valueIsIgnored bool) expressionResult {
	c.compileNodeWithResult(node.Left)
	var jump int
	if valueIsIgnored {
		jump = c.emitJump(node.Location().StartPos.Line, bytecode.JUMP_UNLESS_NIL)
	} else {
		jump = c.emitJump(node.Location().StartPos.Line, bytecode.JUMP_UNLESS_NNP)
	}

	// if nil
	if !valueIsIgnored {
		c.emit(node.Location().StartPos.Line, bytecode.POP)
	}
	c.mustCompileNode(node.Right, valueIsIgnored)

	// if not nil
	c.patchJump(jump, node.Location())
	return valueIgnoredToResult(valueIsIgnored)
}

// Compiles the `||` operator
func (c *Compiler) logicalOr(node *ast.LogicalExpressionNode, valueIsIgnored bool) expressionResult {
	c.compileNodeWithResult(node.Left)
	var jump int
	if valueIsIgnored {
		jump = c.emitJump(node.Location().StartPos.Line, bytecode.JUMP_IF)
	} else {
		jump = c.emitJump(node.Location().StartPos.Line, bytecode.JUMP_IF_NP)
	}

	// if falsy
	if !valueIsIgnored {
		c.emit(node.Location().StartPos.Line, bytecode.POP)
	}
	c.mustCompileNode(node.Right, valueIsIgnored)

	// if truthy
	c.patchJump(jump, node.Location())
	return valueIgnoredToResult(valueIsIgnored)
}

// Compiles the `&&` operator
func (c *Compiler) logicalAnd(node *ast.LogicalExpressionNode, valueIsIgnored bool) expressionResult {
	c.compileNodeWithResult(node.Left)
	var jump int
	if valueIsIgnored {
		jump = c.emitJump(node.Location().StartPos.Line, bytecode.JUMP_UNLESS)
	} else {
		jump = c.emitJump(node.Location().StartPos.Line, bytecode.JUMP_UNLESS_NP)
	}

	// if truthy
	if !valueIsIgnored {
		c.emit(node.Location().StartPos.Line, bytecode.POP)
	}
	c.compileNode(node.Right, valueIsIgnored)

	// if falsy
	c.patchJump(jump, node.Location())
	return valueIgnoredToResult(valueIsIgnored)
}

func (c *Compiler) compileBinaryExpressionNode(node *ast.BinaryExpressionNode) {
	if c.resolveAndEmit(node) {
		return
	}
	c.compileNodeWithResult(node.Left)
	c.compileNodeWithResult(node.Right)
	c.emitBinaryOperation(c.typeOf(node.Left), node.Op, node.Location())
}

func (c *Compiler) emitBinaryOperation(typ types.Type, opToken *token.Token, location *position.Location) {
	line := location.StartPos.Line
	switch opToken.Type {
	case token.PLUS:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.ADD_INT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
			c.emit(line, bytecode.ADD_FLOAT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinAddable)) {
			c.emit(line, bytecode.ADD)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpAdd, 1), location, false)
	case token.MINUS:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.SUBTRACT_INT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
			c.emit(line, bytecode.SUBTRACT_FLOAT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinSubtractable)) {
			c.emit(line, bytecode.SUBTRACT)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpSubtract, 1), location, false)
	case token.STAR:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.MULTIPLY_INT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
			c.emit(line, bytecode.MULTIPLY_FLOAT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinMultipliable)) {
			c.emit(line, bytecode.MULTIPLY)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpMultiply, 1), location, false)
	case token.SLASH:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.DIVIDE_INT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
			c.emit(line, bytecode.DIVIDE_FLOAT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinMultipliable)) {
			c.emit(line, bytecode.DIVIDE)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpDivide, 1), location, false)
	case token.STAR_STAR:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.EXPONENTIATE_INT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinNumeric)) {
			c.emit(line, bytecode.EXPONENTIATE)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpExponentiate, 1), location, false)
	case token.LBITSHIFT:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.LBITSHIFT_INT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
			c.emit(line, bytecode.LBITSHIFT)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpLeftBitshift, 1), location, false)
	case token.LTRIPLE_BITSHIFT:
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinLogicBitshiftable)) {
			c.emit(line, bytecode.LOGIC_LBITSHIFT)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpLogicalLeftBitshift, 1), location, false)
	case token.RBITSHIFT:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.RBITSHIFT_INT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
			c.emit(line, bytecode.RBITSHIFT)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpRightBitshift, 1), location, false)
	case token.RTRIPLE_BITSHIFT:
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinLogicBitshiftable)) {
			c.emit(line, bytecode.LOGIC_RBITSHIFT)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpLogicalRightBitshift, 1), location, false)
	case token.AND:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.BITWISE_AND_INT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
			c.emit(line, bytecode.BITWISE_AND)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpAnd, 1), location, false)
	case token.AND_TILDE:
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
			c.emit(line, bytecode.BITWISE_AND_NOT)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpAndNot, 1), location, false)
	case token.OR:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.BITWISE_OR_INT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
			c.emit(line, bytecode.BITWISE_OR)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpOr, 1), location, false)
	case token.XOR:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.BITWISE_XOR_INT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
			c.emit(line, bytecode.BITWISE_XOR)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpXor, 1), location, false)
	case token.PERCENT:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.MODULO_INT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
			c.emit(line, bytecode.MODULO_FLOAT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinNumeric)) {
			c.emit(line, bytecode.MODULO)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpModulo, 1), location, false)
	case token.LAX_EQUAL:
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinEquatable)) {
			c.emit(line, bytecode.LAX_EQUAL)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpLaxEqual, 1), location, false)
	case token.LAX_NOT_EQUAL:
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinEquatable)) {
			c.emit(line, bytecode.LAX_NOT_EQUAL)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpLaxEqual, 1), location, false)
		c.emit(line, bytecode.NOT)
	case token.EQUAL_EQUAL:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.EQUAL_INT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
			c.emit(line, bytecode.EQUAL_INT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinEquatable)) {
			c.emit(line, bytecode.EQUAL)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpEqual, 1), location, false)
	case token.NOT_EQUAL:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.NOT_EQUAL_INT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.NOT_EQUAL_FLOAT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinEquatable)) {
			c.emit(line, bytecode.NOT_EQUAL)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpEqual, 1), location, false)
		c.emit(line, bytecode.NOT)
	case token.STRICT_EQUAL:
		c.emit(line, bytecode.STRICT_EQUAL)
	case token.STRICT_NOT_EQUAL:
		c.emit(line, bytecode.STRICT_NOT_EQUAL)
	case token.GREATER:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.GREATER_INT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
			c.emit(line, bytecode.GREATER_FLOAT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinNumeric)) {
			c.emit(line, bytecode.GREATER)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpGreaterThan, 1), location, false)
	case token.GREATER_EQUAL:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.GREATER_EQUAL_I)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
			c.emit(line, bytecode.GREATER_EQUAL_F)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinNumeric)) {
			c.emit(line, bytecode.GREATER_EQUAL)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpGreaterThanEqual, 1), location, false)
	case token.LESS:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.LESS_INT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
			c.emit(line, bytecode.LESS_FLOAT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinNumeric)) {
			c.emit(line, bytecode.LESS)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpLessThan, 1), location, false)
	case token.LESS_EQUAL:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.LESS_EQUAL_INT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
			c.emit(line, bytecode.LESS_EQUAL_FLOAT)
			return
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinNumeric)) {
			c.emit(line, bytecode.LESS_EQUAL)
			return
		}
		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpLessThanEqual, 1), location, false)
	case token.SPACESHIP_OP:
		c.emit(line, bytecode.COMPARE)
	case token.INSTANCE_OF_OP:
		c.emit(line, bytecode.INSTANCE_OF)
	case token.REVERSE_INSTANCE_OF_OP:
		c.emit(line, bytecode.INSTANCE_OF)
		c.emit(line, bytecode.NOT)
	case token.ISA_OP:
		c.emit(line, bytecode.IS_A)
	case token.REVERSE_ISA_OP:
		c.emit(line, bytecode.IS_A)
		c.emit(line, bytecode.NOT)
	default:
		c.Errors.AddFailure(fmt.Sprintf("unknown binary operator: %s", opToken.String()), location)
	}
}

// Resolves static AST expressions to Elk values
// and emits Bytecode that loads them.
// Returns false when the node cannot be optimised at compile-time
// and no Bytecode has been generated.
func (c *Compiler) resolveAndEmit(node ast.ExpressionNode) bool {
	result := resolve(node)
	if result.IsUndefined() {
		return false
	}

	c.emitValue(result, node.Location())
	return true
}

func (c *Compiler) resolveAndEmitList(node *ast.ArrayListLiteralNode) bool {
	result := resolveArrayListLiteral(node)
	if result.IsUndefined() {
		return false
	}

	c.emitValue(result, node.Location())
	return true
}

func (c *Compiler) emitValue(val value.Value, location *position.Location) {
	if val.IsReference() {
		switch v := val.AsReference().(type) {
		case *value.ArrayList:
			c.emitArrayList(v, location)
		case *value.ArrayTuple:
			c.emitArrayTuple(v, location)
		case *value.HashSet:
			c.emitHashSet(v, location)
		case *value.HashMap:
			c.emitHashMap(v, location)
		case *value.HashRecord:
			c.emitHashRecord(v, location)
		default:
			c.emitLoadValue(val, location)
		}
		return
	}

	switch val.ValueFlag() {
	case value.TRUE_FLAG:
		c.emit(location.StartPos.Line, bytecode.TRUE)
	case value.FALSE_FLAG:
		c.emit(location.StartPos.Line, bytecode.FALSE)
	case value.NIL_FLAG:
		c.emit(location.StartPos.Line, bytecode.NIL)
	case value.SMALL_INT_FLAG:
		c.emitSmallInt(val.AsSmallInt(), location)
	case value.INT64_FLAG:
		emitSignedInt(c, val, val.AsInlineInt64(), bytecode.LOAD_INT64_8, location)
	case value.UINT64_FLAG:
		emitUnsignedInt(c, val, val.AsInlineUInt64(), bytecode.LOAD_UINT64_8, location)
	case value.INT32_FLAG:
		emitSignedInt(c, val, val.AsInt32(), bytecode.LOAD_INT32_8, location)
	case value.UINT32_FLAG:
		emitUnsignedInt(c, val, val.AsUInt32(), bytecode.LOAD_UINT32_8, location)
	case value.INT16_FLAG:
		emitSignedInt(c, val, val.AsInt16(), bytecode.LOAD_INT16_8, location)
	case value.UINT16_FLAG:
		emitUnsignedInt(c, val, val.AsUInt16(), bytecode.LOAD_UINT16_8, location)
	case value.INT8_FLAG:
		emitSignedInt(c, val, val.AsInt8(), bytecode.LOAD_INT8, location)
	case value.UINT8_FLAG:
		emitUnsignedInt(c, val, val.AsUInt8(), bytecode.LOAD_UINT8, location)
	case value.CHAR_FLAG:
		c.emitChar(val.AsChar(), location)
	case value.FLOAT_FLAG:
		c.emitFloat(val.AsFloat(), location)
	default:
		c.emitLoadValue(val, location)
	}
}

func emitSignedInt[I value.SingedInt](c *Compiler, iVal value.Value, i I, opcodeLoad bytecode.OpCode, location *position.Location) {
	line := location.StartPos.Line
	if i >= math.MinInt8 && i <= math.MaxInt8 {
		c.emit(line, opcodeLoad, byte(i))
		return
	}

	c.emitLoadValue(iVal, location)
}

func emitUnsignedInt[I value.UnsignedInt](c *Compiler, iVal value.Value, i I, opcodeLoad bytecode.OpCode, location *position.Location) {
	line := location.StartPos.Line
	if i <= math.MaxUint8 {
		c.emit(line, opcodeLoad, byte(i))
		return
	}

	c.emitLoadValue(iVal, location)
}

func (c *Compiler) emitSmallInt(i value.SmallInt, location *position.Location) {
	line := location.StartPos.Line
	switch i {
	case -1:
		c.emit(line, bytecode.INT_M1)
		return
	case 0:
		c.emit(line, bytecode.INT_0)
		return
	case 1:
		c.emit(line, bytecode.INT_1)
		return
	case 2:
		c.emit(line, bytecode.INT_2)
		return
	case 3:
		c.emit(line, bytecode.INT_3)
		return
	case 4:
		c.emit(line, bytecode.INT_4)
		return
	case 5:
		c.emit(line, bytecode.INT_5)
		return
	}

	if i >= math.MinInt8 && i <= math.MaxInt8 {
		c.emit(line, bytecode.LOAD_INT_8, byte(i))
		return
	}
	if i >= math.MinInt16 && i <= math.MaxInt16 {
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, uint16(i))
		c.emit(line, bytecode.LOAD_INT_16, bytes...)
		return
	}

	c.emitLoadValue(i.ToValue(), location)
}

func (c *Compiler) emitChar(char value.Char, location *position.Location) {
	line := location.StartPos.Line

	if char >= math.MinInt8 && char <= math.MaxInt8 {
		c.emit(line, bytecode.LOAD_CHAR_8, byte(char))
		return
	}

	c.emitLoadValue(char.ToValue(), location)
}

func (c *Compiler) emitFloat(f value.Float, location *position.Location) {
	line := location.StartPos.Line
	switch f {
	case 0:
		c.emit(line, bytecode.FLOAT_0)
		return
	case 1:
		c.emit(line, bytecode.FLOAT_1)
		return
	case 2:
		c.emit(line, bytecode.FLOAT_2)
		return
	}

	c.emitLoadValue(f.ToValue(), location)
}

func (c *Compiler) emitHashSet(set *value.HashSet, location *position.Location) {
	baseSet := value.NewHashSet(set.Length())
	var mutableElements []value.Value

listLoop:
	for _, element := range set.Table {
		// skip if the bucket is empty or deleted
		if element.IsUndefined() || element == vm.DeletedHashSetValue {
			continue listLoop
		}

		if value.IsMutableCollection(element) {
			mutableElements = append(mutableElements, element)
			continue listLoop
		}

		vm.HashSetAppend(nil, baseSet, element)
	}

	if len(mutableElements) == 0 {
		c.emitLoadValue(value.Ref(baseSet), location)
		c.emit(location.EndPos.Line, bytecode.COPY)
		return
	}

	// capacity
	c.emit(location.StartPos.Line, bytecode.UNDEFINED)
	c.emitLoadValue(value.Ref(baseSet), location)

	for _, element := range mutableElements {
		c.emitValue(element, location)
	}

	c.emitNewHashMap(len(mutableElements), location)
}
func (c *Compiler) emitHashMap(hmap *value.HashMap, location *position.Location) {
	baseMap := value.NewHashMap(hmap.Length())
	var mutablePairs []value.Pair

listLoop:
	for _, element := range hmap.Table {
		// skip if the bucket is empty or deleted
		if element.Key.IsUndefined() {
			continue listLoop
		}

		if value.IsMutableCollection(element.Key) || value.IsMutableCollection(element.Value) {
			mutablePairs = append(mutablePairs, element)
			continue listLoop
		}

		vm.HashMapSet(nil, baseMap, element.Key, element.Value)
	}

	if len(mutablePairs) == 0 {
		c.emitLoadValue(value.Ref(baseMap), location)
		c.emit(location.EndPos.Line, bytecode.COPY)
		return
	}

	// capacity
	c.emit(location.StartPos.Line, bytecode.UNDEFINED)
	c.emitLoadValue(value.Ref(baseMap), location)

	for _, element := range mutablePairs {
		c.emitValue(element.Key, location)
		c.emitValue(element.Value, location)
	}

	c.emitNewHashMap(len(mutablePairs), location)
}

func (c *Compiler) emitHashRecord(hrec *value.HashRecord, location *position.Location) {
	baseRecord := value.NewHashRecord(hrec.Length())
	var mutablePairs []value.Pair

listLoop:
	for _, element := range hrec.Table {
		if element.Key.IsUndefined() {
			continue listLoop
		}

		if value.IsMutableCollection(element.Key) || value.IsMutableCollection(element.Value) {
			mutablePairs = append(mutablePairs, element)
			continue listLoop
		}

		vm.HashRecordSet(nil, baseRecord, element.Key, element.Value)
	}

	if len(mutablePairs) == 0 {
		c.emitLoadValue(value.Ref(baseRecord), location)
		return
	}

	c.emitLoadValue(value.Ref(baseRecord), location)

	for _, element := range mutablePairs {
		c.emitValue(element.Key, location)
		c.emitValue(element.Value, location)
	}

	c.emitNewHashRecord(len(mutablePairs), location)
}

func (c *Compiler) emitArrayList(list *value.ArrayList, location *position.Location) {
	firstMutableElementIndex := -1
	l := *list

listLoop:
	for i, element := range l {
		if value.IsMutableCollection(element) {
			firstMutableElementIndex = i
			break listLoop
		}
	}

	if firstMutableElementIndex == -1 {
		c.emitLoadValue(value.Ref(list), location)
		c.emit(location.EndPos.Line, bytecode.COPY)
		return
	}

	// capacity
	c.emit(location.StartPos.Line, bytecode.UNDEFINED)

	baseList := l[:firstMutableElementIndex]
	c.emitLoadValue(value.Ref(&baseList), location)

	rest := l[firstMutableElementIndex:]
	for _, element := range rest {
		c.emitValue(element, location)
	}

	c.emitNewArrayList(len(rest), location)
}

func (c *Compiler) emitArrayTuple(tuple *value.ArrayTuple, location *position.Location) {
	firstMutableElementIndex := -1
	t := *tuple

listLoop:
	for i, element := range t {
		if value.IsMutableCollection(element) {
			firstMutableElementIndex = i
			break listLoop
		}
	}

	if firstMutableElementIndex == -1 {
		c.emitLoadValue(value.Ref(tuple), location)
		return
	}

	baseTuple := t[:firstMutableElementIndex]
	c.emitLoadValue(value.Ref(&baseTuple), location)

	rest := t[firstMutableElementIndex:]
	for _, element := range rest {
		c.emitValue(element, location)
	}

	c.emitNewArrayList(len(rest), location)
}

func (c *Compiler) compileUnaryExpressionNode(node *ast.UnaryExpressionNode) {
	if c.resolveAndEmit(node) {
		return
	}
	c.compileNodeWithResult(node.Right)
	rightType := c.typeOf(node.Right)

	switch node.Op.Type {
	case token.PLUS:
		if c.checker.IsSubtype(rightType, c.checker.Std(symbol.S_BuiltinNumeric)) {
			c.emit(node.Location().StartPos.Line, bytecode.UNARY_PLUS)
			return
		}

		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpUnaryPlus, 0), node.Location(), false)
	case token.MINUS:
		if c.checker.IsSubtype(rightType, c.checker.StdInt()) {
			c.emit(node.Location().StartPos.Line, bytecode.NEGATE_INT)
			return
		}
		if c.checker.IsSubtype(rightType, c.checker.StdFloat()) {
			c.emit(node.Location().StartPos.Line, bytecode.NEGATE_FLOAT)
			return
		}
		if c.checker.IsSubtype(rightType, c.checker.Std(symbol.S_BuiltinNumeric)) {
			c.emit(node.Location().StartPos.Line, bytecode.NEGATE)
			return
		}

		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpNegate, 0), node.Location(), false)
	case token.BANG:
		// logical not
		c.emit(node.Location().StartPos.Line, bytecode.NOT)
	case token.TILDE:
		// binary negation
		if c.checker.IsSubtype(rightType, c.checker.Std(symbol.S_BuiltinNumeric)) {
			c.emit(node.Location().StartPos.Line, bytecode.BITWISE_NOT)
			return
		}

		c.emitCallMethod(value.NewCallSiteInfo(symbol.OpBitwiseNot, 0), node.Location(), false)
	case token.AND:
		// get singleton class
		c.emit(node.Location().StartPos.Line, bytecode.GET_SINGLETON)
	default:
		c.Errors.AddFailure(fmt.Sprintf("unknown unary operator: %s", node.Op.String()), node.Location())
	}
}

// Emit an instruction that jumps forward with a placeholder offset.
// Returns the offset of placeholder value that has to be patched.
func (c *Compiler) emitJump(line int, op bytecode.OpCode) int {
	c.emit(line, op, 0xff, 0xff)
	return c.nextInstructionOffset() - 2
}

// Emit an instruction that yields a value.
// Provide `nil` as the value when the yielded value is already
// on the stack.
func (c *Compiler) emitYield(location *position.Location, value ast.Node) {
	if value != nil {
		c.compileNodeWithResult(value)
	}

	c.emit(location.EndPos.Line, bytecode.YIELD)
}

// Emit an instruction that returns a value.
// Provide `nil` as the value when the returned value is already
// on the stack.
func (c *Compiler) emitReturn(location *position.Location, value ast.Node) {
	switch c.lastOpCode {
	case bytecode.RETURN, bytecode.RETURN_FIRST_ARG,
		bytecode.RETURN_SELF, bytecode.RETURN_FINALLY:
		return
	}

	if c.isGenerator {
		c.emitYield(location, value)
		c.emit(location.EndPos.Line, bytecode.STOP_ITERATION)
		return
	}

	switch c.mode {
	case setterMethodMode:
		if value != nil {
			c.compileNodeWithResult(value)
		}
		c.emit(location.EndPos.Line, bytecode.POP)
		if c.isNestedInFinally() {
			c.emit(location.EndPos.Line, bytecode.GET_LOCAL8, 1)
			c.emit(location.EndPos.Line, bytecode.RETURN_FINALLY)
		} else {
			c.emit(location.EndPos.Line, bytecode.RETURN_FIRST_ARG)
		}
	case initMethodMode:
		if value != nil {
			c.compileNodeWithResult(value)
		}
		c.emit(location.EndPos.Line, bytecode.POP)
		if c.isNestedInFinally() {
			c.emit(location.EndPos.Line, bytecode.SELF)
			c.emit(location.EndPos.Line, bytecode.RETURN_FINALLY)
		} else {
			c.emit(location.EndPos.Line, bytecode.RETURN_SELF)
		}
	case namespaceMode:
		if value != nil {
			c.compileNodeWithResult(value)
		}
		if c.lastOpCode != bytecode.NIL {
			c.emit(location.EndPos.Line, bytecode.POP)
			c.emit(location.EndPos.Line, bytecode.NIL)
		}
		if c.isNestedInFinally() {
			c.emit(location.EndPos.Line, bytecode.RETURN_FINALLY)
		} else {
			c.emit(location.EndPos.Line, bytecode.RETURN)
		}
	default:
		if value != nil {
			c.compileNodeWithResult(value)
		}
		if c.isNestedInFinally() {
			c.emit(location.EndPos.Line, bytecode.RETURN_FINALLY)
		} else {
			c.emit(location.EndPos.Line, bytecode.RETURN)
		}
	}
}

func (c *Compiler) emitFinalReturn(location *position.Location, value ast.Node) {
	if !c.isGenerator {
		c.emitReturn(location, value)
		return
	}

	if c.lastOpCode == bytecode.STOP_ITERATION {
		c.emit(location.StartPos.Line, bytecode.LOOP, 0, 4) // 3 bytes for LOOP + 1 byte for STOP_ITERATION
		return
	}

	c.emitYield(location, value)
	start := c.nextInstructionOffset()
	c.emit(location.EndPos.Line, bytecode.STOP_ITERATION)
	c.emitLoop(location, start)
}

// Emit an instruction that jumps back to the given Bytecode offset.
func (c *Compiler) emitLoop(location *position.Location, startOffset int) {
	c.emit(location.EndPos.Line, bytecode.LOOP)

	offset := c.nextInstructionOffset() - startOffset + 2
	if offset > math.MaxUint16 {
		c.Errors.AddFailure(
			fmt.Sprintf("too many bytes to jump backward: %d", math.MaxUint16),
			location,
		)
	}

	c.emitUint16(uint16(offset))
}

// Overwrite the placeholder operand of a jump instruction
func (c *Compiler) patchJumpWithTarget(target int, offset int, location *position.Location) {
	if target > math.MaxUint16 {
		c.Errors.AddFailure(
			fmt.Sprintf("too many bytes to jump over: %d", target),
			location,
		)
		return
	}

	c.Bytecode.Instructions[offset] = byte((target >> 8) & 0xff)
	c.Bytecode.Instructions[offset+1] = byte(target & 0xff)
}

// Overwrite the placeholder operand of a jump instruction
func (c *Compiler) patchJump(offset int, location *position.Location) {
	c.patchJumpWithTarget(c.nextInstructionOffset()-offset-2, offset, location)
}

// Emit an instruction that sets a local variable or value
func (c *Compiler) emitSetLocal(line int, index uint16, valueIsIgnored bool) expressionResult {
	if valueIsIgnored {
		c.emitSetLocalPop(line, index)
		return expressionCompiledWithoutResult
	} else {
		c.emitSetLocalNoPop(line, index)
		return expressionCompiled
	}
}

// Emit an instruction that sets a local variable or value without popping it
func (c *Compiler) emitSetLocalNoPop(line int, index uint16) {
	c.emit(line, bytecode.DUP)
	c.emitSetLocalPop(line, index)
}

// Emit an instruction that sets a local variable or value.
func (c *Compiler) emitSetLocalPop(line int, index uint16) {
	switch index {
	case 1:
		c.emit(line, bytecode.SET_LOCAL_1)
		return
	case 2:
		c.emit(line, bytecode.SET_LOCAL_2)
		return
	case 3:
		c.emit(line, bytecode.SET_LOCAL_3)
		return
	case 4:
		c.emit(line, bytecode.SET_LOCAL_4)
		return
	}

	if index > math.MaxUint8 {
		c.emit(line, bytecode.SET_LOCAL16)
		c.emitUint16(index)
		return
	}

	c.emit(line, bytecode.SET_LOCAL8, byte(index))
}

// Emit an instruction that gets the value of a local.
func (c *Compiler) emitGetLocal(line int, index uint16) {
	switch index {
	case 1:
		c.emit(line, bytecode.GET_LOCAL_1)
		return
	case 2:
		c.emit(line, bytecode.GET_LOCAL_2)
		return
	case 3:
		c.emit(line, bytecode.GET_LOCAL_3)
		return
	case 4:
		c.emit(line, bytecode.GET_LOCAL_4)
		return
	}

	if index > math.MaxUint8 {
		c.emit(line, bytecode.GET_LOCAL16)
		c.emitUint16(index)
		return
	}

	c.emit(line, bytecode.GET_LOCAL8, byte(index))
}

// Emit an instruction that sets an upvalue.
func (c *Compiler) emitSetUpvalue(line int, index uint16, valueIsIgnored bool) expressionResult {
	if valueIsIgnored {
		c.emitSetUpvaluePop(line, index)
		return expressionCompiledWithoutResult
	} else {
		c.emitSetUpvalueNoPop(line, index)
		return expressionCompiled
	}
}

// Emit an instruction that sets an upvalue.
func (c *Compiler) emitSetUpvaluePop(line int, index uint16) {
	switch index {
	case 0:
		c.emit(line, bytecode.SET_UPVALUE_0)
		return
	case 1:
		c.emit(line, bytecode.SET_UPVALUE_1)
		return
	}

	if index > math.MaxUint8 {
		c.emit(line, bytecode.SET_UPVALUE16)
		c.emitUint16(index)
		return
	}

	c.emit(line, bytecode.SET_UPVALUE8, byte(index))
}

// Emit an instruction that sets an upvalue without popping it.
func (c *Compiler) emitSetUpvalueNoPop(line int, index uint16) {
	c.emit(line, bytecode.DUP)
	c.emitSetUpvaluePop(line, index)
}

// Emit an instruction that gets the value of an upvalue.
func (c *Compiler) emitGetUpvalue(line int, index uint16) {
	switch index {
	case 0:
		c.emit(line, bytecode.GET_UPVALUE_0)
		return
	case 1:
		c.emit(line, bytecode.GET_UPVALUE_1)
		return
	}

	if index > math.MaxUint8 {
		c.emit(line, bytecode.GET_UPVALUE16)
		c.emitUint16(index)
		return
	}

	c.emit(line, bytecode.GET_UPVALUE8, byte(index))
}

// Emit an instruction that sets an upvalue.
func (c *Compiler) emitCloseUpvalue(line int, index uint16) {
	switch index {
	case 1:
		c.emit(line, bytecode.CLOSE_UPVALUE_1)
		return
	case 2:
		c.emit(line, bytecode.CLOSE_UPVALUE_2)
		return
	case 3:
		c.emit(line, bytecode.CLOSE_UPVALUE_3)
		return
	}

	if index > math.MaxUint8 {
		c.emit(line, bytecode.CLOSE_UPVALUE16)
		c.emitUint16(index)
		return
	}

	c.emit(line, bytecode.CLOSE_UPVALUE8, byte(index))
}

// Emit an instruction that loads a value from the pool
func (c *Compiler) emitAddValue(val value.Value, location *position.Location, opCode8, opCode16 bytecode.OpCode) int {
	id, size := c.Bytecode.AddValue(val)
	switch size {
	case bytecode.UINT8_SIZE:
		c.Bytecode.AddInstruction(location.StartPos.Line, opCode8, byte(id))
	case bytecode.UINT16_SIZE:
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, uint16(id))
		c.Bytecode.AddInstruction(location.StartPos.Line, opCode16, bytes...)
	default:
		c.Errors.AddFailure(
			fmt.Sprintf("value pool limit reached: %d", math.MaxUint16),
			location,
		)
		return -1
	}

	return id
}

// Emit an instruction that retrieves a constant
func (c *Compiler) emitGetConst(val value.Symbol, location *position.Location) int {
	return c.emitAddValue(
		val.ToValue(),
		location,
		bytecode.GET_CONST8,
		bytecode.GET_CONST16,
	)
}

// Add a value to the value pool and emit appropriate bytecode.
func (c *Compiler) emitLoadValue(val value.Value, location *position.Location) int {
	id, size := c.Bytecode.AddValue(val)

	switch id {
	case 0:
		c.emit(location.StartPos.Line, bytecode.LOAD_VALUE_0)
		return id
	case 1:
		c.emit(location.StartPos.Line, bytecode.LOAD_VALUE_1)
		return id
	case 2:
		c.emit(location.StartPos.Line, bytecode.LOAD_VALUE_2)
		return id
	case 3:
		c.emit(location.StartPos.Line, bytecode.LOAD_VALUE_3)
		return id
	}

	switch size {
	case bytecode.UINT8_SIZE:
		c.Bytecode.AddInstruction(location.StartPos.Line, bytecode.LOAD_VALUE8, byte(id))
	case bytecode.UINT16_SIZE:
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, uint16(id))
		c.Bytecode.AddInstruction(location.StartPos.Line, bytecode.LOAD_VALUE16, bytes...)
	default:
		c.Errors.AddFailure(
			fmt.Sprintf("value pool limit reached: %d", math.MaxUint16),
			location,
		)
		return -1
	}

	return id
}

// Emit an instruction that instantiates an object
func (c *Compiler) emitInstantiate(args int, location *position.Location) {
	if args <= math.MaxUint8 {
		c.Bytecode.AddInstruction(location.StartPos.Line, bytecode.INSTANTIATE8, byte(args))
		return
	}

	if args <= math.MaxUint16 {
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, uint16(args))
		c.Bytecode.AddInstruction(location.StartPos.Line, bytecode.INSTANTIATE8, bytes...)
		return
	}

	c.Errors.AddFailure(
		fmt.Sprintf("constructor argument limit reached: %d", math.MaxUint16),
		location,
	)
}

// Emit an instruction that sets the value of an instance variable
func (c *Compiler) emitSetInstanceVariable(name value.Symbol, location *position.Location, valueIsIgnored bool) {
	if valueIsIgnored {
		c.emitSetInstanceVariablePop(name, location)
		return
	}

	c.emitSetInstanceVariableNoPop(name, location)
}

// Emit an instruction that sets the value of an instance variable and pops it
func (c *Compiler) emitSetInstanceVariablePop(name value.Symbol, location *position.Location) {
	self := c.checker.SelfType()

	switch self := self.(type) {
	case types.NamespaceWithIvarIndices:
		index := self.IvarIndices().GetIndex(name)
		c.emitSetInstanceVariableByIndex(index, location)
	default:
		c.emitSetInstanceVariableByName(name, location)
	}
}

// Emit an instruction that sets the value of an instance variable by its index
func (c *Compiler) emitSetInstanceVariableByIndex(index int, location *position.Location) {
	line := location.StartPos.Line
	switch index {
	case 0:
		c.emit(line, bytecode.SET_IVAR_0)
		return
	case 1:
		c.emit(line, bytecode.SET_IVAR_1)
		return
	case 2:
		c.emit(line, bytecode.SET_IVAR_2)
		return
	}

	if index >= math.MinInt8 && index <= math.MaxInt8 {
		c.emit(line, bytecode.SET_IVAR8, byte(index))
		return
	}
	if index >= math.MinInt16 && index <= math.MaxInt16 {
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, uint16(index))
		c.emit(line, bytecode.SET_IVAR16, bytes...)
		return
	}

	c.Errors.AddFailure(
		fmt.Sprintf("too big instance variable index: %d", index),
		location,
	)
}

// Emit an instruction that sets the value of an instance variable by name
func (c *Compiler) emitSetInstanceVariableByName(name value.Symbol, location *position.Location) {
	id, size := c.Bytecode.AddValue(name.ToValue())
	switch size {
	case bytecode.UINT8_SIZE, bytecode.UINT16_SIZE:
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, uint16(id))
		c.Bytecode.AddInstruction(location.StartPos.Line, bytecode.SET_IVAR_NAME16, bytes...)
	default:
		c.Errors.AddFailure(
			fmt.Sprintf("value pool limit reached: %d", math.MaxUint16),
			location,
		)
	}
}

// Emit an instruction that sets the value of an instance variable without popping
func (c *Compiler) emitSetInstanceVariableNoPop(name value.Symbol, location *position.Location) {
	c.emit(location.StartPos.Line, bytecode.DUP)
	c.emitSetInstanceVariablePop(name, location)
}

func (c *Compiler) emitGetInstanceVariable(name value.Symbol, location *position.Location) {
	self := c.checker.SelfType()
	switch self := self.(type) {
	case types.NamespaceWithIvarIndices:
		ivarIndices := self.IvarIndices()
		index := ivarIndices.GetIndex(name)
		c.emitGetInstanceVariableByIndex(index, location)
	default:
		c.emitGetInstanceVariableByName(name, location)
	}
}

// Emit an instruction that reads the value of an instance variable by index.
func (c *Compiler) emitGetInstanceVariableByIndex(index int, location *position.Location) {
	line := location.StartPos.Line
	switch index {
	case 0:
		c.emit(line, bytecode.GET_IVAR_0)
		return
	case 1:
		c.emit(line, bytecode.GET_IVAR_1)
		return
	case 2:
		c.emit(line, bytecode.GET_IVAR_2)
		return
	}

	if index >= math.MinInt8 && index <= math.MaxInt8 {
		c.emit(line, bytecode.GET_IVAR8, byte(index))
		return
	}
	if index >= math.MinInt16 && index <= math.MaxInt16 {
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, uint16(index))
		c.emit(line, bytecode.GET_IVAR16, bytes...)
		return
	}

	c.Errors.AddFailure(
		fmt.Sprintf("too big instance variable index: %d", index),
		location,
	)
}

// Emit an instruction that reads the value of an instance variable by name.
func (c *Compiler) emitGetInstanceVariableByName(name value.Symbol, location *position.Location) int {
	id, size := c.Bytecode.AddValue(name.ToValue())
	switch size {
	case bytecode.UINT8_SIZE, bytecode.UINT16_SIZE:
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, uint16(id))
		c.Bytecode.AddInstruction(location.StartPos.Line, bytecode.GET_IVAR_NAME16, bytes...)
	default:
		c.Errors.AddFailure(
			fmt.Sprintf("value pool limit reached: %d", math.MaxUint16),
			location,
		)
		return -1
	}

	return id
}

// Emit an instruction that calls a method on self
func (c *Compiler) emitCallSelf(callInfo *value.CallSiteInfo, location *position.Location, tailCall bool) int {
	if tailCall {
		return c.emitAddValue(
			value.Ref(callInfo),
			location,
			bytecode.CALL_SELF_TCO8,
			bytecode.CALL_SELF_TCO16,
		)
	}

	return c.emitAddValue(
		value.Ref(callInfo),
		location,
		bytecode.CALL_SELF8,
		bytecode.CALL_SELF16,
	)
}

// Emit an instruction that calls the `call` method
func (c *Compiler) emitCall(callInfo *value.CallSiteInfo, location *position.Location) int {
	return c.emitAddValue(
		value.Ref(callInfo),
		location,
		bytecode.CALL8,
		bytecode.CALL16,
	)
}

// Emit an instruction that calls a method
func (c *Compiler) emitCallMethod(callInfo *value.CallSiteInfo, location *position.Location, tailCall bool) int {
	if tailCall {
		return c.emitAddValue(
			value.Ref(callInfo),
			location,
			bytecode.CALL_METHOD_TCO8,
			bytecode.CALL_METHOD_TCO16,
		)
	}

	return c.emitAddValue(
		value.Ref(callInfo),
		location,
		bytecode.CALL_METHOD8,
		bytecode.CALL_METHOD16,
	)
}

// Emit an instruction that calls the `next` method
func (c *Compiler) emitCallNext(callInfo *value.CallSiteInfo, location *position.Location) int {
	return c.emitAddValue(
		value.Ref(callInfo),
		location,
		bytecode.NEXT8,
		bytecode.NEXT16,
	)
}

// Emit an opcode with optional bytes.
func (c *Compiler) emit(line int, op bytecode.OpCode, bytes ...byte) {
	c.secondToLastOpCode = c.lastOpCode
	c.lastOpCode = op
	c.Bytecode.AddInstruction(line, op, bytes...)
}

func (c *Compiler) emitByte(byt byte) {
	c.Bytecode.AddBytes(byt)
}

func (c *Compiler) emitUint16(n uint16) {
	c.Bytecode.AppendUint16(n)
}

func (c *Compiler) emitUint32(n uint32) {
	c.Bytecode.AppendUint32(n)
}

func (c *Compiler) enterScope(label string, typ scopeType) {
	c.scopes = append(c.scopes, newScope(label, typ))
}

// Pop the values of local variables in the current scope
func (c *Compiler) leaveScope(line int) {
	varsToPop := c.leaveScopeWithoutMutating(line)

	currentDepth := len(c.scopes) - 1
	c.lastLocalIndex -= varsToPop
	c.scopes[currentDepth] = nil
	c.scopes = c.scopes[:currentDepth]
}

// Pop the values of local variables in the current scope.
// Allows you to emit the instructions to leave the same scope a few times,
// because it doesn't mutate the state of the compiler.
func (c *Compiler) leaveScopeWithoutMutating(line int) int {
	currentDepth := len(c.scopes) - 1

	c.closeUpvaluesInScope(line, c.scopes[currentDepth])

	varsToPop := len(c.scopes[currentDepth].localTable)
	c.emitLeaveScope(line, c.lastLocalIndex, varsToPop)
	return varsToPop
}

func (c *Compiler) closeUpvaluesInCurrentScope(line int) {
	currentDepth := len(c.scopes) - 1
	c.closeUpvaluesInScope(line, c.scopes[currentDepth])
}

func (c *Compiler) closeUpvaluesInScope(line int, scope *scope) {
	for _, local := range scope.localTable {
		if !local.hasUpvalue {
			continue
		}

		c.emitCloseUpvalue(line, local.index)
	}
}

func (c *Compiler) emitLeaveScope(line, maxLocalIndex, varsToPop int) {
	if varsToPop <= 0 {
		return
	}

	if maxLocalIndex > math.MaxUint8 || varsToPop > math.MaxUint8 {
		c.emit(line, bytecode.LEAVE_SCOPE32)
		c.emitUint16(uint16(maxLocalIndex))
		c.emitUint16(uint16(varsToPop))
	} else {
		c.emit(line, bytecode.LEAVE_SCOPE16, byte(maxLocalIndex), byte(varsToPop))
	}
}

// Register a local variable.
func (c *Compiler) defineLocal(name string, location *position.Location) *local {
	varScope := c.scopes.last()
	_, ok := varScope.localTable[name]
	if ok {
		c.Errors.AddFailure(
			fmt.Sprintf("a variable with this name has already been declared in this scope `%s`", name),
			location,
		)
		return nil
	}
	return c.defineVariableInScope(varScope, name, location)
}

// Register a local variable, reusing the variable with the same name that has already been defined in this scope.
func (c *Compiler) defineLocalOverrideCurrentScope(name string, location *position.Location) *local {
	varScope := c.scopes.last()
	if currentVar, ok := varScope.localTable[name]; ok {
		return currentVar
	}
	return c.defineVariableInScope(varScope, name, location)
}

func (c *Compiler) defineVariableInScope(scope *scope, name string, location *position.Location) *local {
	if c.lastLocalIndex == math.MaxUint16 {
		c.Errors.AddFailure(
			fmt.Sprintf("exceeded the maximum number of local variables (%d): %s", math.MaxUint16, name),
			location,
		)
		return nil
	}

	c.lastLocalIndex++
	if c.lastLocalIndex > c.maxLocalIndex {
		c.maxLocalIndex = c.lastLocalIndex
	}
	newVar := &local{
		index: uint16(c.lastLocalIndex),
	}
	scope.localTable[name] = newVar
	return newVar
}

// Resolve a local variable and get its index.
func (c *Compiler) resolveLocal(name string) (*local, bool) {
	var localVal *local
	var found bool
	for i := len(c.scopes) - 1; i >= 0; i-- {
		varScope := c.scopes[i]
		local, ok := varScope.localTable[name]
		if ok {
			localVal = local
			found = true
			break
		}
		if !c.unhygienic && varScope.typ == macroBoundaryScopeType {
			break
		}
	}

	if !found {
		return localVal, false
	}

	return localVal, true
}
