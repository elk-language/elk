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
	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"

	"github.com/elk-language/elk/token"
)

const MainName = "<main>"

func CreateMainCompiler(checker types.Checker, loc *position.Location, errors *error.SyncErrorList) *Compiler {
	compiler := New(loc.Filename, topLevelMode, loc, checker)
	compiler.Errors = errors
	return compiler
}

func (c *Compiler) CreateMainCompiler(checker types.Checker, loc *position.Location, errors *error.SyncErrorList) *Compiler {
	compiler := New(loc.Filename, topLevelMode, loc, checker)
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
	span := &parent.Bytecode.Location.Span
	parent.emitValue(value.Ref(c.Bytecode), span)
	parent.emit(span.StartPos.Line, bytecode.EXEC)
	parent.emit(span.StartPos.Line, bytecode.POP)
}

// Compiler mode
type mode uint8

const (
	topLevelMode mode = iota
	namespaceMode
	methodMode
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
	defaultScopeType   scopeType = iota
	loopScopeType                // this scope is a loop
	doFinallyScopeType           // this scope is inside do with a finally block
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
	typ    loopJumpInfoType
	offset int
	span   *position.Span
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
	Errors             *error.SyncErrorList
	scopes             scopes
	loopJumpSets       []*loopJumpSet
	offsetValueIds     []int // ids of integers in the value pool that represent bytecode offsets
	lastLocalIndex     int   // index of the last local variable
	maxLocalIndex      int   // max index of a local variable
	predefinedLocals   int
	mode               mode
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
		Errors:         error.NewSyncErrorList(),
	}
	// reserve the first slot on the stack for `self`
	c.defineLocal("$self", &position.Span{})
	switch mode {
	case topLevelMode, namespaceMode,
		methodMode, setterMethodMode, initMethodMode:
		c.predefinedLocals = 1
	}
	return c
}

func (c *Compiler) EmitReturnNil() {
	span := &c.Bytecode.Location.Span
	c.emit(span.EndPos.Line, bytecode.NIL)
	c.emit(span.EndPos.Line, bytecode.RETURN)
}

func (c *Compiler) EmitReturn() {
	span := &c.Bytecode.Location.Span
	if c.lastOpCode != bytecode.RETURN {
		c.emit(span.EndPos.Line, bytecode.RETURN)
	}
	c.prepLocals()
}

func (c *Compiler) typeOf(node ast.Node) types.Type {
	return node.Type(c.checker.Env())
}

func (c *Compiler) compileGlobalEnv() {
	span := &c.Bytecode.Location.Span
	env := c.checker.Env()
	c.compileModuleDefinition(env.Root, env.Root, value.ToSymbol("Root"), span)
}

func (c *Compiler) compileNamespaceDefinition(parentNamespace, namespace types.Namespace, namespaceType byte, constName value.Symbol, span *position.Span) {
	if !namespace.IsDefined() {
		switch p := parentNamespace.(type) {
		case *types.SingletonClass:
			c.emitGetConst(value.ToSymbol(p.AttachedObject.Name()), span)
			c.emit(span.StartPos.Line, bytecode.GET_SINGLETON)
		default:
			c.emitGetConst(value.ToSymbol(p.Name()), span)
		}
		c.emitValue(constName.ToValue(), span)
		c.emit(span.StartPos.Line, bytecode.DEF_NAMESPACE, namespaceType)
		namespace.SetDefined(true)
	}

	for name, subtype := range types.SortedSubtypes(namespace) {
		if subtype.Type == namespace {
			continue
		}
		c.compileSubtypeDefinition(namespace, subtype.Type, name, span)
	}
}

func (c *Compiler) compileModuleDefinition(parentNamespace types.Namespace, module *types.Module, constName value.Symbol, span *position.Span) {
	c.compileNamespaceDefinition(parentNamespace, module, bytecode.DEF_MODULE_FLAG, constName, span)
}

func (c *Compiler) compileClassDefinition(parentNamespace types.Namespace, class *types.Class, constName value.Symbol, span *position.Span) {
	c.compileNamespaceDefinition(parentNamespace, class, bytecode.DEF_CLASS_FLAG, constName, span)
}

func (c *Compiler) compileMixinDefinition(parentNamespace types.Namespace, mixin *types.Mixin, constName value.Symbol, span *position.Span) {
	c.compileNamespaceDefinition(parentNamespace, mixin, bytecode.DEF_MIXIN_FLAG, constName, span)
}

func (c *Compiler) compileInterfaceDefinition(parentNamespace types.Namespace, iface *types.Interface, constName value.Symbol, span *position.Span) {
	c.compileNamespaceDefinition(parentNamespace, iface, bytecode.DEF_INTERFACE_FLAG, constName, span)
}

func (c *Compiler) compileSubtypeDefinition(parentNamespace types.Namespace, typ types.Type, constName value.Symbol, span *position.Span) {
	switch t := typ.(type) {
	case *types.Module:
		c.compileModuleDefinition(parentNamespace, t, constName, span)
	case *types.Class:
		c.compileClassDefinition(parentNamespace, t, constName, span)
	case *types.Mixin:
		c.compileMixinDefinition(parentNamespace, t, constName, span)
	case *types.Interface:
		c.compileInterfaceDefinition(parentNamespace, t, constName, span)
	}
}

func (c *Compiler) CompileClassInheritance(class *types.Class, span *position.Span) {
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
	c.emitGetConst(name, span)

	superclassName := value.ToSymbol(superclass.Name())
	c.emitGetConst(superclassName, span)

	c.emit(span.StartPos.Line, bytecode.SET_SUPERCLASS)
}

func (c *Compiler) CompileInclude(target types.Namespace, mixin *types.Mixin, span *position.Span) {
	switch t := target.(type) {
	case *types.SingletonClass:
		targetName := value.ToSymbol(t.AttachedObject.Name())
		c.emitGetConst(targetName, span)
		c.emit(span.StartPos.Line, bytecode.GET_SINGLETON)
	default:
		targetName := value.ToSymbol(t.Name())
		c.emitGetConst(targetName, span)
	}

	mixinName := value.ToSymbol(mixin.Name())
	c.emitGetConst(mixinName, span)

	c.emit(span.StartPos.Line, bytecode.INCLUDE)
}

func (c *Compiler) InitExpressionCompiler(filename string, span *position.Span) *Compiler {
	exprCompiler := New(filename, topLevelMode, c.Bytecode.Location, c.checker)
	exprCompiler.Errors = c.Errors

	c.emitValue(value.Ref(exprCompiler.Bytecode), span)
	c.emit(span.StartPos.Line, bytecode.EXEC)
	c.emit(span.StartPos.Line, bytecode.POP)

	return exprCompiler
}

func (c *Compiler) CompileExpressionsInFile(node *ast.ProgramNode) {
	c.compileNode(node)
}

// Entry point to the compilation process
func (c *Compiler) compileProgram(node ast.Node) {
	c.compileNode(node)
	c.emitReturn(node.Span(), nil)
	c.prepLocals()
}

// Entry point for compiling the body of a function.
func (c *Compiler) compileFunction(span *position.Span, parameters []ast.ParameterNode, body []ast.StatementNode) {
	c.Bytecode.SetParameterCount(len(parameters))

	for _, param := range parameters {
		p := param.(*ast.FormalParameterNode)
		pSpan := p.Span()

		local := c.defineLocal(p.Name, pSpan)
		if local == nil {
			return
		}
		c.predefinedLocals++

		if p.Initialiser != nil {
			c.Bytecode.IncrementOptionalParameterCount()

			c.emitGetLocal(span.StartPos.Line, local.index)
			jump := c.emitJump(pSpan.StartPos.Line, bytecode.JUMP_UNLESS_UNP)

			c.emit(pSpan.StartPos.Line, bytecode.POP)
			c.compileNode(p.Initialiser)
			c.emitSetLocal(pSpan.StartPos.Line, local.index)

			c.patchJump(jump, pSpan)
			// pops the value after SET_LOCAL when the argument was missing
			// or pops the condition value used for jump when the argument was present
			c.emit(pSpan.StartPos.Line, bytecode.POP)
		}
	}
	c.compileStatements(body, span)

	c.emitReturn(span, nil)
	c.prepLocals()
}

func (c *Compiler) InitMethodCompiler(span *position.Span) (*Compiler, int) {
	methodCompiler := New("<methodDefinitions>", topLevelMode, c.Bytecode.Location, c.checker)
	methodCompiler.Errors = c.Errors
	methodCompiler.Parent = c

	offset := c.nextInstructionOffset()
	c.emitValue(value.Ref(methodCompiler.Bytecode), span)
	c.emit(span.StartPos.Line, bytecode.EXEC)
	c.emit(span.StartPos.Line, bytecode.POP)

	return methodCompiler, offset
}

func (c *Compiler) CompileMethods(span *position.Span, execOffset int) {
	c.compileMethodsWithinModule(c.checker.Env().Root, span)
	if len(c.Bytecode.Instructions) > 0 {
		c.emit(span.EndPos.Line, bytecode.NIL)
		c.emit(span.EndPos.Line, bytecode.RETURN)
		return
	}

	// If no instructions were emitted, remove the EXEC instruction block
	c.Parent.removeBytes(execOffset, 4)
	c.Parent.removeMethodDefinitionsBytecodeFunction()
}

func (c *Compiler) removeBytes(offset int, count int) {
	c.Bytecode.Instructions = slices.Concat(c.Bytecode.Instructions[:offset], c.Bytecode.Instructions[offset+count:])
	lineInfo := c.Bytecode.LineInfoList.GetLineInfo(offset)
	lineInfo.InstructionCount -= count
}

var methodDefinitionsSymbol = value.ToSymbol("<methodDefinitions>")

func (c *Compiler) removeMethodDefinitionsBytecodeFunction() {
	for i, val := range c.Bytecode.Values {
		val, ok := val.SafeAsReference().(*vm.BytecodeFunction)
		if !ok {
			continue
		}

		if val.Name() == methodDefinitionsSymbol {
			c.Bytecode.Values[i] = value.Undefined
			break
		}
	}
}

func (c *Compiler) compileMethodsWithinModule(module *types.Module, span *position.Span) {
	if types.NamespaceHasAnyDefinableMethods(module) {
		c.emitGetConst(value.ToSymbol(module.Name()), span)
		c.emit(span.StartPos.Line, bytecode.GET_SINGLETON)

		for methodName, method := range types.SortedOwnMethods(module) {
			c.compileMethodDefinition(methodName, method, span)
		}

		for aliasName, alias := range types.SortedOwnMethodAliases(module) {
			c.compileMethodAliasDefinition(aliasName, alias, span)
		}

		c.emit(span.StartPos.Line, bytecode.POP)
	}

	for _, subtype := range types.SortedSubtypes(module) {
		if subtype.Type == module {
			continue
		}
		c.compileMethodsWithinType(subtype.Type, span)
	}
}

func (c *Compiler) compileMethodAliasDefinition(aliasName value.Symbol, alias *types.MethodAlias, span *position.Span) {
	if !alias.IsDefinable() {
		return
	}

	c.emitValue(alias.Method.Name.ToValue(), span)
	c.emitValue(aliasName.ToValue(), span)
	c.emit(span.StartPos.Line, bytecode.DEF_METHOD_ALIAS)
	alias.Compiled = true
}

func (c *Compiler) compileMethodDefinition(name value.Symbol, method *types.Method, span *position.Span) {
	if !method.IsDefinable() {
		return
	}

	if method.IsAttribute() {
		if method.IsSetter() {
			nameStr := name.String()
			c.emitValue(value.ToSymbol(nameStr[:len(nameStr)-1]).ToValue(), span)
			c.emit(span.StartPos.Line, bytecode.DEF_SETTER)
			method.SetCompiled(true)
			return
		}

		c.emitValue(name.ToValue(), span)
		c.emit(span.StartPos.Line, bytecode.DEF_GETTER)
		method.SetCompiled(true)
		return
	}

	c.emitValue(value.Ref(method.Bytecode), span)
	c.emitValue(name.ToValue(), span)
	c.emit(span.StartPos.Line, bytecode.DEF_METHOD)
	method.SetCompiled(true)
}

func (c *Compiler) compileMethodsWithinNamespace(namespace types.Namespace, span *position.Span) {
	namespaceHasCompiledMethods := types.NamespaceHasAnyDefinableMethods(namespace)

	singleton := namespace.Singleton()
	singletonHasCompiledMethods := types.NamespaceHasAnyDefinableMethods(singleton)

	if namespaceHasCompiledMethods || singletonHasCompiledMethods {
		c.emitGetConst(value.ToSymbol(namespace.Name()), span)

		for methodName, method := range types.SortedOwnMethods(namespace) {
			c.compileMethodDefinition(methodName, method, span)
		}

		for aliasName, alias := range types.SortedOwnMethodAliases(namespace) {
			c.compileMethodAliasDefinition(aliasName, alias, span)
		}

		if singletonHasCompiledMethods {
			c.emit(span.StartPos.Line, bytecode.GET_SINGLETON)

			for methodName, method := range types.SortedOwnMethods(singleton) {
				c.compileMethodDefinition(methodName, method, span)
			}

			for aliasName, alias := range types.SortedOwnMethodAliases(singleton) {
				c.compileMethodAliasDefinition(aliasName, alias, span)
			}
		}

		c.emit(span.StartPos.Line, bytecode.POP)
	}

	for _, subtype := range types.SortedSubtypes(namespace) {
		if subtype.Type == namespace {
			continue
		}
		c.compileMethodsWithinType(subtype.Type, span)
	}
}

func (c *Compiler) compileMethodsWithinType(typ types.Type, span *position.Span) {
	switch t := typ.(type) {
	case *types.Module:
		c.compileMethodsWithinModule(t, span)
	case *types.Class:
		c.compileMethodsWithinNamespace(t, span)
	case *types.Mixin:
		c.compileMethodsWithinNamespace(t, span)
	case *types.Interface:
		c.compileMethodsWithinNamespace(t, span)
	}
}

func (c *Compiler) CompileMethodBody(node *ast.MethodDefinitionNode, name value.Symbol) *vm.BytecodeFunction {
	var mode mode
	if node.IsSetter() {
		mode = setterMethodMode
	} else if node.Name == "#init" {
		mode = initMethodMode
	} else {
		mode = methodMode
	}

	methodCompiler := New(name.String(), mode, c.newLocation(node.Span()), c.checker)
	methodCompiler.Errors = c.Errors
	methodCompiler.compileMethodBody(node.Span(), node.Parameters, node.Body)

	return methodCompiler.Bytecode
}

// Entry point for compiling the body of a method.
func (c *Compiler) compileMethodBody(span *position.Span, parameters []ast.ParameterNode, body []ast.StatementNode) {
	c.Bytecode.SetParameterCount(len(parameters))

	for _, param := range parameters {
		p := param.(*ast.MethodParameterNode)
		pSpan := p.Span()

		local := c.defineLocal(p.Name, pSpan)
		if local == nil {
			return
		}
		c.predefinedLocals++

		if p.Initialiser != nil {
			c.Bytecode.IncrementOptionalParameterCount()

			c.emitGetLocal(span.StartPos.Line, local.index)
			jump := c.emitJump(pSpan.StartPos.Line, bytecode.JUMP_UNLESS_UNP)

			c.emit(pSpan.StartPos.Line, bytecode.POP)
			c.compileNode(p.Initialiser)
			c.emitSetLocal(pSpan.StartPos.Line, local.index)

			c.patchJump(jump, pSpan)
			// pops the value after SET_LOCAL when the argument was missing
			// or pops the condition value used for jump when the argument was present
			c.emit(pSpan.StartPos.Line, bytecode.POP)
		}

		if p.SetInstanceVariable {
			c.emitGetLocal(span.StartPos.Line, local.index)
			c.emitSetInstanceVariable(value.ToSymbol(p.Name), pSpan)
			// pop the value after setting it
			c.emit(pSpan.StartPos.Line, bytecode.POP)
		}
	}
	c.compileStatements(body, span)

	c.emitReturn(span, nil)
	c.prepLocals()
}

// Entry point for compiling the body of a namespace eg. Module, Class, Mixin, Struct, Interface.
func (c *Compiler) compileNamespace(node ast.Node) bool {
	span := node.Span()
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
			c.newLocation(span),
		)
		return false
	}

	c.emitReturn(span, nil)
	c.prepLocals()
	return true
}

// Create a new location struct with the given position.
func (c *Compiler) newLocation(span *position.Span) *position.Location {
	return position.NewLocationWithSpan(c.Bytecode.Location.Filename, span)
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

func (c *Compiler) findLoopJumpSet(label string, span *position.Span) *loopJumpSet {
	if len(c.loopJumpSets) < 1 {
		c.Errors.AddFailure(
			"cannot jump with `break` or `continue` outside of a loop",
			c.newLocation(span),
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
		c.newLocation(span),
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

func (c *Compiler) addLoopJump(label string, typ loopJumpInfoType, offset int, span *position.Span) {
	jumpSet := c.findLoopJumpSet(label, span)
	if jumpSet == nil {
		return
	}

	c.addLoopJumpTo(jumpSet, typ, offset)
}

func (c *Compiler) compilePublicConstantNode(node *ast.PublicConstantNode) {
	c.emitGetConst(value.ToSymbol(node.Value), node.Span())
}

func (c *Compiler) compilePrivateConstantNode(node *ast.PrivateConstantNode) {
	c.emitGetConst(value.ToSymbol(node.Value), node.Span())
}

func (c *Compiler) compileNode(node ast.Node) bool {
	switch node := node.(type) {
	case nil, *ast.AliasDeclarationNode, *ast.IncludeExpressionNode,
		*ast.EmptyStatementNode, *ast.MethodDefinitionNode, *ast.UsingExpressionNode,
		*ast.ConstantDeclarationNode, *ast.TypeDefinitionNode, *ast.GenericTypeDefinitionNode,
		*ast.MethodSignatureDefinitionNode, *ast.ImplementExpressionNode,
		*ast.StructDeclarationNode, *ast.GenericReceiverlessMethodCallNode,
		*ast.ReceiverlessMethodCallNode, *ast.AttrDeclarationNode,
		*ast.SetterDeclarationNode, *ast.GetterDeclarationNode, *ast.InitDefinitionNode,
		*ast.InstanceVariableDeclarationNode:
		return false
	case *ast.ProgramNode:
		c.compileStatements(node.Body, node.Span())
	case *ast.ExtendWhereBlockExpressionNode:
		c.compileStatements(node.Body, node.Span())
	case *ast.ExpressionStatementNode:
		return c.compileNode(node.Expression)
	case *ast.LabeledExpressionNode:
		c.compileLabeledExpressionNode(node)
	case *ast.UndefinedLiteralNode:
		c.emit(node.Span().StartPos.Line, bytecode.UNDEFINED)
	case *ast.PublicConstantNode:
		c.compilePublicConstantNode(node)
	case *ast.PrivateConstantNode:
		c.compilePrivateConstantNode(node)
	case *ast.GenericConstantNode:
		c.compileNode(node.Constant)
	case *ast.SelfLiteralNode:
		c.emit(node.Span().StartPos.Line, bytecode.SELF)
	case *ast.AssignmentExpressionNode:
		c.compileAssignmentExpressionNode(node)
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
	case *ast.ClosureLiteralNode:
		c.compileClosureLiteralNode(node)
	case *ast.SwitchExpressionNode:
		c.compileSwitchExpressionNode(node)
	case *ast.SubscriptExpressionNode:
		c.compileSubscriptExpressionNode(node)
	case *ast.NilSafeSubscriptExpressionNode:
		c.compileNilSafeSubscriptExpressionNode(node)
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
	case *ast.VariablePatternDeclarationNode:
		c.compilerVariablePatternDeclarationNode(node)
	case *ast.VariableDeclarationNode:
		c.compileVariableDeclarationNode(node)
	case *ast.ValuePatternDeclarationNode:
		c.compileValuePatternDeclarationNode(node)
	case *ast.ValueDeclarationNode:
		c.compileValueDeclarationNode(node)
	case *ast.PublicIdentifierNode:
		c.compileLocalVariableAccess(node.Value, node.Span())
	case *ast.PrivateIdentifierNode:
		c.compileLocalVariableAccess(node.Value, node.Span())
	case *ast.InstanceVariableNode:
		c.compileInstanceVariableAccess(node.Value, node.Span())
	case *ast.BinaryExpressionNode:
		c.compileBinaryExpressionNode(node)
	case *ast.LogicalExpressionNode:
		c.compileLogicalExpressionNode(node)
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
		c.emitValue(value.Ref(value.String(node.Value)), node.Span())
	case *ast.DoubleQuotedStringLiteralNode:
		c.emitValue(value.Ref(value.String(node.Value)), node.Span())
	case *ast.InterpolatedStringLiteralNode:
		c.compileInterpolatedStringLiteralNode(node)
	case *ast.InterpolatedSymbolLiteralNode:
		c.compileInterpolatedSymbolLiteralNode(node)
	case *ast.CharLiteralNode:
		c.emitValue(value.Char(node.Value).ToValue(), node.Span())
	case *ast.RawCharLiteralNode:
		c.emitValue(value.Char(node.Value).ToValue(), node.Span())
	case *ast.FalseLiteralNode:
		c.emit(node.Span().StartPos.Line, bytecode.FALSE)
	case *ast.TrueLiteralNode:
		c.emit(node.Span().StartPos.Line, bytecode.TRUE)
	case *ast.NilLiteralNode:
		c.emit(node.Span().StartPos.Line, bytecode.NIL)
	case *ast.ThrowExpressionNode:
		c.compileThrowExpressionNode(node)
	case *ast.MustExpressionNode:
		c.compileMustExpressionNode(node)
	case *ast.TryExpressionNode:
		c.compileTryExpressionNode(node)
	case *ast.AsExpressionNode:
		c.compileAsExpressionNode(node)
	case *ast.TypeofExpressionNode:
		c.compileTypeofExpressionNode(node)
	case *ast.DoExpressionNode:
		c.compileDoExpressionNode(node)
	case *ast.IfExpressionNode:
		c.compileIfExpression(false, node.Condition, node.ThenBody, node.ElseBody, node.Span())
	case *ast.UnlessExpressionNode:
		c.compileIfExpression(true, node.Condition, node.ThenBody, node.ElseBody, node.Span())
	case *ast.ModifierIfElseNode:
		c.compileModifierIfExpression(false, node.Condition, node.ThenExpression, node.ElseExpression, node.Span())
	case *ast.ModifierNode:
		c.compileModifierExpressionNode("", node)
	case *ast.BreakExpressionNode:
		c.compileBreakExpressionNode(node)
	case *ast.ContinueExpressionNode:
		c.compileContinueExpressionNode(node)
	case *ast.LoopExpressionNode:
		c.compileLoopExpressionNode("", node.ThenBody, node.Span())
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
		c.compilePostfixExpressionNode(node)
	case *ast.SimpleSymbolLiteralNode:
		c.emitValue(value.ToSymbol(node.Content).ToValue(), node.Span())
	case *ast.IntLiteralNode:
		c.compileIntLiteralNode(node)
	case *ast.Int8LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 8)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), c.newLocation(node.Span()))
			return true
		}
		// BENCHMARK: Compare with storing
		// ints inline in Bytecode instead of as constants.
		c.emitValue(value.Int8(i).ToValue(), node.Span())
	case *ast.Int16LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 16)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), c.newLocation(node.Span()))
			return true
		}
		c.emitValue(value.Int16(i).ToValue(), node.Span())
	case *ast.Int32LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 32)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), c.newLocation(node.Span()))
			return true
		}
		c.emitValue(value.Int32(i).ToValue(), node.Span())
	case *ast.Int64LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 64)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), c.newLocation(node.Span()))
			return true
		}
		c.emitValue(value.Int64(i).ToValue(), node.Span())
	case *ast.UInt8LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 8)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), c.newLocation(node.Span()))
			return true
		}
		c.emitValue(value.UInt8(i).ToValue(), node.Span())
	case *ast.UInt16LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 16)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), c.newLocation(node.Span()))
			return true
		}
		c.emitValue(value.UInt16(i).ToValue(), node.Span())
	case *ast.UInt32LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 32)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), c.newLocation(node.Span()))
			return true
		}
		c.emitValue(value.UInt32(i).ToValue(), node.Span())
	case *ast.UInt64LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 64)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), c.newLocation(node.Span()))
			return true
		}
		c.emitValue(value.UInt64(i).ToValue(), node.Span())
	case *ast.FloatLiteralNode:
		f, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			c.Errors.AddFailure(err.Error(), c.newLocation(node.Span()))
			return true
		}
		c.emitValue(value.Float(f).ToValue(), node.Span())
	case *ast.BigFloatLiteralNode:
		f, err := value.ParseBigFloat(node.Value)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), c.newLocation(node.Span()))
			return true
		}
		c.emitValue(value.Ref(f), node.Span())
	case *ast.Float64LiteralNode:
		f, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			c.Errors.AddFailure(err.Error(), c.newLocation(node.Span()))
			return true
		}
		c.emitValue(value.Float64(f).ToValue(), node.Span())
	case *ast.Float32LiteralNode:
		f, err := strconv.ParseFloat(node.Value, 32)
		if err != nil {
			c.Errors.AddFailure(err.Error(), c.newLocation(node.Span()))
			return true
		}
		c.emitValue(value.Float32(f).ToValue(), node.Span())
	default:
		c.Errors.AddFailure(
			fmt.Sprintf("compilation of this node has not been implemented: %T", node),
			c.newLocation(node.Span()),
		)
	}

	return true
}

func (c *Compiler) compileTypeofExpressionNode(node *ast.TypeofExpressionNode) {
	c.compileNode(node.Value)
}

func (c *Compiler) compileTryExpressionNode(node *ast.TryExpressionNode) {
	c.compileNode(node.Value)
}

func (c *Compiler) compileMustExpressionNode(node *ast.MustExpressionNode) {
	span := node.Span()
	c.compileNode(node.Value)
	c.emit(span.StartPos.Line, bytecode.MUST)
}

func (c *Compiler) compileAsExpressionNode(node *ast.AsExpressionNode) {
	span := node.Span()
	c.compileNode(node.Value)
	c.compileNode(node.RuntimeType)
	c.emit(span.StartPos.Line, bytecode.AS)
}

func (c *Compiler) compileThrowExpressionNode(node *ast.ThrowExpressionNode) {
	span := node.Span()
	if node.Value != nil {
		c.compileNode(node.Value)
	} else {
		c.emitValue(value.Ref(value.NewError(value.ErrorClass, "error")), span)
	}

	c.emit(span.StartPos.Line, bytecode.THROW)
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
	span := node.Span()
	switch n := namespace.(type) {
	case *types.SingletonClass:
		namespaceName := value.ToSymbol(n.AttachedObject.Name())
		c.emitGetConst(namespaceName, node.Constant.Span())
		c.emit(span.StartPos.Line, bytecode.GET_SINGLETON)
	default:
		namespaceName := value.ToSymbol(n.Name())
		c.emitGetConst(namespaceName, node.Constant.Span())
	}
	c.emitValue(constName.ToValue(), span)
	c.compileNode(node.Initialiser)
	c.emit(span.StartPos.Line, bytecode.DEF_CONST)
}

func (c *Compiler) compileDoExpressionNode(node *ast.DoExpressionNode) {
	span := node.Span()

	doStartOffset := c.nextInstructionOffset()

	var scopeType scopeType
	if len(node.Finally) > 0 {
		scopeType = doFinallyScopeType
	} else {
		scopeType = defaultScopeType
	}

	c.enterScope("", scopeType)
	c.compileStatements(node.Body, span)
	c.leaveScope(span.EndPos.Line)

	doEndOffset := c.nextInstructionOffset()

	if len(node.Finally) > 0 {
		c.enterScope("", defaultScopeType)
		c.compileStatements(node.Finally, span)
		// pop the return value of finally leaving the return value of do
		c.emit(span.StartPos.Line, bytecode.POP)
		c.leaveScope(span.EndPos.Line)
	}

	if len(node.Catches) <= 0 && len(node.Finally) <= 0 {
		return
	}

	jumpOverCatchOffset := c.emitJump(span.StartPos.Line, bytecode.JUMP)

	var jumpsToEndOfCatch []int
	catchStartOffset := c.nextInstructionOffset()

	c.registerCatch(doStartOffset, doEndOffset, catchStartOffset, false)

	for _, catchNode := range node.Catches {
		span := catchNode.Span()
		c.pattern(catchNode.Pattern)
		jumpOverCatchBody := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)

		c.compileStatements(catchNode.Body, catchNode.Span())

		if len(node.Finally) < 1 {
			// pop the thrown value and the stack trace, leaving the return value of the catch
			c.emit(span.EndPos.Line, bytecode.POP_N_SKIP_ONE, 2)
		}
		jump := c.emitJump(span.EndPos.Line, bytecode.JUMP)
		jumpsToEndOfCatch = append(jumpsToEndOfCatch, jump)

		c.patchJump(jumpOverCatchBody, span)
	}

	if len(node.Finally) > 0 {
		c.emit(span.EndPos.Line, bytecode.TRUE)
	} else {
		c.emit(span.EndPos.Line, bytecode.RETHROW)
	}

	var jumpOverFalseOffset int
	if len(node.Finally) > 0 {

		jumpOverFalseOffset = c.emitJump(span.EndPos.Line, bytecode.JUMP)
	}
	for _, jump := range jumpsToEndOfCatch {
		c.patchJump(jump, span)
	}
	if len(node.Finally) > 0 {
		c.emit(span.EndPos.Line, bytecode.FALSE)
		c.patchJump(jumpOverFalseOffset, span)

		jumpOverReturnBreakOrContinueEntryOffset := c.emitJump(span.EndPos.Line, bytecode.JUMP)
		finallyEntryOffset := c.nextInstructionOffset()
		c.registerCatch(doStartOffset, doEndOffset, finallyEntryOffset, true)
		// entry point for return when executing finally
		c.emit(span.EndPos.Line, bytecode.NIL)

		jumpOverBreakOrContinueEntryOffset := c.emitJump(span.EndPos.Line, bytecode.JUMP)
		// entry point for break or continue when executing finally
		c.emit(span.EndPos.Line, bytecode.UNDEFINED)

		c.patchJump(jumpOverBreakOrContinueEntryOffset, span)
		c.patchJump(jumpOverReturnBreakOrContinueEntryOffset, span)

		c.compileStatements(node.Finally, span)

		c.emit(span.EndPos.Line, bytecode.SWAP)
		jumpOverFinallyBreakOrContinueOffset := c.emitJump(span.EndPos.Line, bytecode.JUMP_UNLESS_UNP)
		c.emit(span.EndPos.Line, bytecode.POP_N, 2)
		c.emit(span.EndPos.Line, bytecode.JUMP_TO_FINALLY)
		c.patchJump(jumpOverFinallyBreakOrContinueOffset, span)

		jumpToRethrowOffset := c.emitJump(span.EndPos.Line, bytecode.JUMP_IF_NP)
		jumpToFinallyReturnOffset := c.emitJump(span.EndPos.Line, bytecode.JUMP_IF_NIL_NP)
		// FALSE
		c.emit(span.EndPos.Line, bytecode.POP_N, 2)          // pop the flag and return value of finally
		c.emit(span.EndPos.Line, bytecode.POP_N_SKIP_ONE, 2) // pop the thrown value and the stack trace leaving the return value of catch
		jumpToEndOffset := c.emitJump(span.EndPos.Line, bytecode.JUMP)

		c.patchJump(jumpToFinallyReturnOffset, span)
		// return with finally
		c.emit(span.EndPos.Line, bytecode.POP_N, 2) // pop the flag and return value of finally
		c.emit(span.EndPos.Line, bytecode.RETURN_FINALLY)

		c.patchJump(jumpToRethrowOffset, span)
		// pop the flag and the return value of finally
		c.emit(span.EndPos.Line, bytecode.POP_N, 2)
		c.emit(span.EndPos.Line, bytecode.RETHROW)

		c.patchJump(jumpToEndOffset, span)
	}

	c.patchJump(jumpOverCatchOffset, span)
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
	span := node.Span()
	if node.Value == nil {
		c.emit(span.StartPos.Line, bytecode.NIL)
	} else {
		c.compileNode(node.Value)
	}

	finallyCount := c.countFinallyInLoop(node.Label)
	if finallyCount <= 0 {
		c.leaveScopeOnBreak(span.StartPos.Line, node.Label)

		breakJumpOffset := c.emitJump(span.StartPos.Line, bytecode.JUMP)
		c.addLoopJump(node.Label, breakLoopJump, breakJumpOffset, span)
		return
	}

	jumpOffsetId := c.emitLoadValue(value.Undefined, span)
	c.offsetValueIds = append(c.offsetValueIds, jumpOffsetId)
	c.addLoopJump(node.Label, breakFinallyLoopJump, jumpOffsetId, span)

	c.emitValue(value.SmallInt(finallyCount).ToValue(), span)
	c.emit(span.StartPos.Line, bytecode.JUMP_TO_FINALLY)
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
	span := node.Span()
	loop := c.findLoopJumpSet(node.Label, span)
	if loop == nil {
		return
	}

	if !loop.returnsValueFromLastIteration {
		if node.Value != nil {
			c.compileNode(node.Value)
			c.emit(span.StartPos.Line, bytecode.POP)
		}
	} else {
		if node.Value == nil {
			c.emit(span.StartPos.Line, bytecode.NIL)
		} else {
			c.compileNode(node.Value)
		}
	}

	finallyCount := c.countFinallyInLoop(node.Label)
	if finallyCount <= 0 {
		c.leaveScopeOnContinue(span.StartPos.Line, node.Label)

		continueJumpOffset := c.emitJump(span.StartPos.Line, bytecode.LOOP)
		c.addLoopJumpTo(loop, continueLoopJump, continueJumpOffset)
		return
	}

	jumpOffsetId := c.emitLoadValue(value.Undefined, span)
	c.offsetValueIds = append(c.offsetValueIds, jumpOffsetId)
	c.addLoopJump(node.Label, continueFinallyLoopJump, jumpOffsetId, span)

	c.emitValue(value.SmallInt(finallyCount).ToValue(), span)
	c.emit(span.StartPos.Line, bytecode.JUMP_TO_FINALLY)
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
			c.patchJump(loopJump.offset, loopJump.span)
		case continueLoopJump:
			target := continueOffset - loopJump.offset
			if target >= 0 {
				// jump forward
				// override the opcode to JUMP
				c.Bytecode.Instructions[loopJump.offset-1] = byte(bytecode.JUMP)
				c.patchJumpWithTarget(target-2, loopJump.offset, loopJump.span)
			} else {
				// jump backward
				// override the opcode to LOOP
				c.Bytecode.Instructions[loopJump.offset-1] = byte(bytecode.LOOP)
				c.patchJumpWithTarget((-target)+2, loopJump.offset, loopJump.span)
			}
		default:
			panic(fmt.Sprintf("invalid loop jump info: %#v", loopJump))
		}
	}
	c.loopJumpSets = c.loopJumpSets[:len(c.loopJumpSets)-1]
}

func (c *Compiler) compileLoopExpressionNode(label string, body []ast.StatementNode, span *position.Span) {
	c.enterScope(label, loopScopeType)
	c.initLoopJumpSet(label, false)

	start := c.nextInstructionOffset()
	c.enterScope("", defaultScopeType)
	if c.compileStatementsOk(body) {
		c.emit(span.EndPos.Line, bytecode.POP)
	}
	c.leaveScope(span.EndPos.Line)
	c.emitLoop(span, start)

	c.leaveScope(span.EndPos.Line)
	c.patchLoopJumps(start)
}

func (c *Compiler) compileWhileExpressionNode(label string, node *ast.WhileExpressionNode) {
	span := node.Span()

	if result := resolve(node.Condition); !result.IsUndefined() {
		if value.Falsy(result) {
			// the loop won't run at all
			// it can be optimised into a simple NIL operation
			c.emit(span.StartPos.Line, bytecode.NIL)
			return
		}

		// the loop is endless
		c.compileLoopExpressionNode(label, node.ThenBody, span)
		return
	}

	c.enterScope(label, loopScopeType)
	c.initLoopJumpSet(label, true)

	c.emit(span.StartPos.Line, bytecode.NIL)
	// loop start
	c.enterScope("", defaultScopeType)
	start := c.nextInstructionOffset()
	var loopBodyOffset int

	// loop condition eg. `i < 5`
	c.compileNode(node.Condition)
	// jump past the loop if the condition is falsy
	loopBodyOffset = c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)
	// pop the return value of the last iteration
	c.emit(span.StartPos.Line, bytecode.POP)

	// loop body
	c.compileStatements(node.ThenBody, span)

	c.leaveScope(span.EndPos.Line)
	// jump to loop condition
	c.emitLoop(span, start)

	// after loop
	c.patchJump(loopBodyOffset, span)

	c.leaveScope(span.EndPos.Line)
	c.patchLoopJumps(start)
}

func (c *Compiler) modifierWhileExpression(label string, node *ast.ModifierNode) {
	span := node.Span()

	body := node.Left
	condition := node.Right

	var conditionIsStaticFalsy bool

	if result := resolve(condition); !result.IsUndefined() {
		if value.Truthy(result) {
			// the loop is endless
			c.compileLoopExpressionNode(label, ast.ExpressionToStatements(body), span)
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
	c.compileNode(body)
	// continue
	continueOffset := c.nextInstructionOffset()
	if conditionIsStaticFalsy {
		// the loop has a static falsy condition
		// it will only finish one iteration
		c.leaveScope(span.EndPos.Line)
		c.patchLoopJumps(continueOffset)
		return
	}

	// loop condition eg. `i < 5`
	c.compileNode(condition)
	// jump past the loop if the condition is falsy
	loopBodyOffset = c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)
	// pop the return value of the last iteration
	c.emit(span.StartPos.Line, bytecode.POP)

	c.leaveScope(span.EndPos.Line)
	// jump to loop start
	c.emitLoop(span, start)

	// after loop
	c.patchJump(loopBodyOffset, span)

	c.leaveScope(span.EndPos.Line)
	c.patchLoopJumps(continueOffset)
}

func (c *Compiler) modifierUntilExpression(label string, node *ast.ModifierNode) {
	span := node.Span()

	body := node.Left
	condition := node.Right

	var conditionIsStaticTruthy bool

	if result := resolve(condition); !result.IsUndefined() {
		if value.Falsy(result) {
			// the loop is endless
			c.compileLoopExpressionNode(label, ast.ExpressionToStatements(body), span)
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
	c.compileNode(body)
	// continue
	continueOffset := c.nextInstructionOffset()
	if conditionIsStaticTruthy {
		// the loop has a static truthy condition
		// it will only finish one iteration
		c.leaveScope(span.EndPos.Line)
		c.patchLoopJumps(continueOffset)
		return
	}

	// loop condition eg. `i > 5`
	c.compileNode(condition)
	// jump past the loop if the condition is truthy
	loopBodyOffset = c.emitJump(span.StartPos.Line, bytecode.JUMP_IF)
	// pop the return value of the last iteration
	c.emit(span.StartPos.Line, bytecode.POP)

	c.leaveScope(span.EndPos.Line)
	// jump to loop start
	c.emitLoop(span, start)

	// after loop
	c.patchJump(loopBodyOffset, span)

	c.leaveScope(span.EndPos.Line)
	c.patchLoopJumps(continueOffset)
}

func (c *Compiler) compileUntilExpressionNode(label string, node *ast.UntilExpressionNode) {
	span := node.Span()

	if result := resolve(node.Condition); !result.IsUndefined() {
		if value.Falsy(result) {
			// the loop is endless
			c.compileLoopExpressionNode(label, node.ThenBody, span)
			return
		}

		// the loop won't run at all
		// it can be optimised into a simple NIL operation
		c.emit(span.StartPos.Line, bytecode.NIL)
		return
	}

	c.enterScope(label, loopScopeType)
	c.initLoopJumpSet(label, true)

	c.emit(span.StartPos.Line, bytecode.NIL)
	// loop start
	start := c.nextInstructionOffset()
	c.enterScope("", defaultScopeType)
	var loopBodyOffset int

	// loop condition eg. `i > 5`
	c.compileNode(node.Condition)
	// jump past the loop if the condition is truthy
	loopBodyOffset = c.emitJump(span.StartPos.Line, bytecode.JUMP_IF)
	// pop the return value of the last iteration
	c.emit(span.StartPos.Line, bytecode.POP)

	// loop body
	c.compileStatements(node.ThenBody, span)

	c.leaveScope(span.EndPos.Line)
	// jump to loop condition
	c.emitLoop(span, start)

	// after loop
	c.patchJump(loopBodyOffset, span)

	c.leaveScope(span.EndPos.Line)
	c.patchLoopJumps(start)
}

// Compile a labeled expression eg. `$foo: println("bar")`
func (c *Compiler) compileLabeledExpressionNode(node *ast.LabeledExpressionNode) {
	switch expr := node.Expression.(type) {
	case *ast.WhileExpressionNode:
		c.compileWhileExpressionNode(node.Label, expr)
	case *ast.UntilExpressionNode:
		c.compileUntilExpressionNode(node.Label, expr)
	case *ast.LoopExpressionNode:
		c.compileLoopExpressionNode(node.Label, expr.ThenBody, expr.Span())
	case *ast.NumericForExpressionNode:
		c.compileNumericForExpressionNode(node.Label, expr)
	case *ast.ForInExpressionNode:
		c.compileForInExpressionNode(node.Label, expr)
	case *ast.ModifierForInNode:
		c.compileModifierForInNode(node.Label, expr)
	case *ast.ModifierNode:
		c.compileModifierExpressionNode(node.Label, expr)
	default:
		c.compileNode(node.Expression)
	}
}

// Compile a for in loop eg. `for i in [1, 2] then println(i)`
func (c *Compiler) compileForInExpressionNode(label string, node *ast.ForInExpressionNode) {
	c.compileForIn(
		label,
		node.Pattern,
		node.InExpression,
		func() {
			c.compileStatements(node.ThenBody, node.Span())
		},
		node.Span(),
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
			if !c.compileNode(node.ThenExpression) {
				c.emit(node.ThenExpression.Span().StartPos.Line, bytecode.NIL)
			}
		},
		node.Span(),
		false,
	)
}

func (c *Compiler) compileForIn(
	label string,
	param ast.PatternNode,
	inExpression ast.ExpressionNode,
	then func(),
	span *position.Span,
	collectionLiteral bool,
) {
	c.enterScope(label, loopScopeType)
	c.initLoopJumpSet(label, false)

	c.compileNode(inExpression)
	c.emit(span.StartPos.Line, bytecode.GET_ITERATOR)

	iteratorVarName := fmt.Sprintf("#!forIn%d", len(c.scopes))
	iteratorVar := c.defineLocal(iteratorVarName, span)
	c.emitSetLocal(span.StartPos.Line, iteratorVar.index)
	c.emit(span.EndPos.Line, bytecode.POP)

	// loop start
	start := c.nextInstructionOffset()
	continueOffset := start
	c.enterScope("", defaultScopeType)

	c.emitGetLocal(span.StartPos.Line, iteratorVar.index)
	loopBodyOffset := c.emitJump(span.StartPos.Line, bytecode.FOR_IN)

	switch p := param.(type) {
	case *ast.PrivateIdentifierNode:
		paramVar := c.defineLocal(p.Value, param.Span())
		c.emitSetLocal(param.Span().StartPos.Line, paramVar.index)
		c.emit(param.Span().EndPos.Line, bytecode.POP)
	case *ast.PublicIdentifierNode:
		paramVar := c.defineLocal(p.Value, param.Span())
		c.emitSetLocal(param.Span().StartPos.Line, paramVar.index)
		c.emit(param.Span().EndPos.Line, bytecode.POP)
	default:
		c.pattern(param)
		jumpOverErrorOffset := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF)

		c.emitValue(
			value.Ref(value.NewError(
				value.PatternNotMatchedErrorClass,
				"assigned value does not match the pattern defined in for in loop",
			)),
			span,
		)
		c.emit(span.EndPos.Line, bytecode.THROW)

		c.patchJump(jumpOverErrorOffset, span)
	}

	// loop body
	then()

	// pop the return value of the block
	if !collectionLiteral {
		c.emit(span.EndPos.Line, bytecode.POP)
	}
	c.leaveScope(span.EndPos.Line)
	// jump to loop condition
	c.emitLoop(span, start)

	// after loop
	c.patchJump(loopBodyOffset, span)
	if !collectionLiteral {
		c.emit(span.EndPos.Line, bytecode.NIL)
	}

	c.leaveScope(span.EndPos.Line)
	c.patchLoopJumps(continueOffset)
}

// Compile a numeric for loop eg. `for i := 0; i < 5; i += 1 then println(i)`
func (c *Compiler) compileNumericForExpressionNode(label string, node *ast.NumericForExpressionNode) {
	span := node.Span()

	if node.Initialiser == nil && node.Condition == nil && node.Increment == nil {
		// the loop is endless
		c.compileLoopExpressionNode(label, node.ThenBody, span)
		return
	}

	c.enterScope(label, loopScopeType)
	c.initLoopJumpSet(label, true)

	// loop initialiser eg. `i := 0`
	if node.Initialiser != nil {
		c.compileNode(node.Initialiser)
		c.emit(span.EndPos.Line, bytecode.POP)
	}

	c.emit(span.StartPos.Line, bytecode.NIL)
	// loop start
	c.enterScope("", defaultScopeType)
	start := c.nextInstructionOffset()
	continueOffset := start

	var loopBodyOffset int
	// loop condition eg. `i < 5`
	if node.Condition != nil {
		c.compileNode(node.Condition)
		// jump past the loop if the condition is falsy
		loopBodyOffset = c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)
		// pop the return value of the last iteration
		c.emit(span.EndPos.Line, bytecode.POP)
	} else {
		// pop the return value of the last iteration
		c.emit(span.EndPos.Line, bytecode.POP)
	}

	// loop body
	c.compileStatements(node.ThenBody, span)

	if node.Increment != nil {
		continueOffset = c.nextInstructionOffset()
		// increment step eg. `i += 1`
		c.compileNode(node.Increment)
		c.emit(span.EndPos.Line, bytecode.POP)
	}

	c.leaveScope(span.EndPos.Line)
	// jump to loop condition
	c.emitLoop(span, start)

	// after loop
	if node.Condition != nil {
		c.patchJump(loopBodyOffset, span)
	}

	c.leaveScope(span.EndPos.Line)
	c.patchLoopJumps(continueOffset)
}

func (c *Compiler) emitSetterCall(name string, span *position.Span) {
	nameSymbol := value.ToSymbol(name + "=")
	callInfo := value.NewCallSiteInfo(nameSymbol, 1)
	c.emitCallMethod(callInfo, span)
}

func (c *Compiler) emitGetterCall(name string, span *position.Span) {
	nameSymbol := value.ToSymbol(name)
	callInfo := value.NewCallSiteInfo(nameSymbol, 0)
	c.emitCallMethod(callInfo, span)
}

func (c *Compiler) compilePostfixExpressionNode(node *ast.PostfixExpressionNode) {
	switch n := node.Expression.(type) {
	case *ast.PublicIdentifierNode:
		// get variable value
		c.compileLocalVariableAccess(n.Value, n.Span())

		switch node.Op.Type {
		case token.PLUS_PLUS:
			c.emit(node.Span().EndPos.Line, bytecode.INCREMENT)
		case token.MINUS_MINUS:
			c.emit(node.Span().EndPos.Line, bytecode.DECREMENT)
		default:
			panic(fmt.Sprintf("invalid postfix operator: %#v", node.Op))
		}

		// set variable
		c.setLocalWithoutValue(n.Value, n.Span())
	case *ast.PrivateIdentifierNode:
		// get variable value
		c.compileLocalVariableAccess(n.Value, n.Span())

		switch node.Op.Type {
		case token.PLUS_PLUS:
			c.emit(node.Span().EndPos.Line, bytecode.INCREMENT)
		case token.MINUS_MINUS:
			c.emit(node.Span().EndPos.Line, bytecode.DECREMENT)
		default:
			panic(fmt.Sprintf("invalid postfix operator: %#v", node.Op))
		}

		// set variable
		c.setLocalWithoutValue(n.Value, n.Span())
	case *ast.SubscriptExpressionNode:
		// get value
		c.compileNode(n.Receiver)
		c.compileNode(n.Key)
		c.emit(node.Span().EndPos.Line, bytecode.DUP_N, 2)
		c.emit(node.Span().EndPos.Line, bytecode.SUBSCRIPT)

		switch node.Op.Type {
		case token.PLUS_PLUS:
			c.emit(node.Span().EndPos.Line, bytecode.INCREMENT)
		case token.MINUS_MINUS:
			c.emit(node.Span().EndPos.Line, bytecode.DECREMENT)
		default:
			panic(fmt.Sprintf("invalid postfix operator: %#v", node.Op))
		}

		// set value
		c.emit(node.Span().EndPos.Line, bytecode.SUBSCRIPT_SET)
	case *ast.InstanceVariableNode:
		switch c.mode {
		case topLevelMode:
			c.Errors.AddFailure(
				"instance variables cannot be set in the top level",
				c.newLocation(node.Span()),
			)
		}
		// get value
		ivarSymbol := value.ToSymbol(n.Value)
		c.emitGetInstanceVariable(ivarSymbol, n.Span())

		switch node.Op.Type {
		case token.PLUS_PLUS:
			c.emit(node.Span().EndPos.Line, bytecode.INCREMENT)
		case token.MINUS_MINUS:
			c.emit(node.Span().EndPos.Line, bytecode.DECREMENT)
		default:
			panic(fmt.Sprintf("invalid postfix operator: %#v", node.Op))
		}

		// set instance variable
		c.emitSetInstanceVariable(ivarSymbol, node.Span())
	case *ast.AttributeAccessNode:
		// get value
		c.compileNode(n.Receiver)
		name := value.ToSymbol(n.AttributeName)
		callInfo := value.NewCallSiteInfo(name, 0)
		c.emitCallMethod(callInfo, node.Span())

		switch node.Op.Type {
		case token.PLUS_PLUS:
			c.emit(node.Span().EndPos.Line, bytecode.INCREMENT)
		case token.MINUS_MINUS:
			c.emit(node.Span().EndPos.Line, bytecode.DECREMENT)
		default:
			panic(fmt.Sprintf("invalid postfix operator: %#v", node.Op))
		}

		// set attribute
		c.emitSetterCall(n.AttributeName, node.Span())
	default:
		c.Errors.AddFailure(
			fmt.Sprintf("cannot assign to: %T", node.Expression),
			c.newLocation(node.Span()),
		)
	}
}

func (c *Compiler) attributeAssignment(node *ast.AssignmentExpressionNode, attr *ast.AttributeAccessNode) {
	// compile the argument
	switch node.Op.Type {
	case token.EQUAL_OP:
		c.compileNode(attr.Receiver)
		c.compileNode(node.Right)
		c.emitSetterCall(attr.AttributeName, node.Span())
	case token.OR_OR_EQUAL:
		span := node.Span()
		// Read the current value
		c.compileNode(attr.Receiver)
		c.emitGetterCall(attr.AttributeName, span)

		jump := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF_NP)

		// if falsy
		c.emit(span.StartPos.Line, bytecode.POP)
		c.compileNode(node.Right)
		c.emitSetterCall(attr.AttributeName, span)

		// if truthy
		c.patchJump(jump, span)
	case token.AND_AND_EQUAL:
		span := node.Span()
		// Read the current value
		c.compileNode(attr.Receiver)
		c.emitGetterCall(attr.AttributeName, span)

		jump := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS_NP)

		// if truthy
		c.emit(span.StartPos.Line, bytecode.POP)
		c.compileNode(node.Right)
		c.emitSetterCall(attr.AttributeName, span)

		// if falsy
		c.patchJump(jump, span)
	case token.QUESTION_QUESTION_EQUAL:
		span := node.Span()
		// Read the current value
		c.compileNode(attr.Receiver)
		c.emitGetterCall(attr.AttributeName, span)

		nilJump := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF_NIL_NP)
		nonNilJump := c.emitJump(span.StartPos.Line, bytecode.JUMP)

		// if nil
		c.patchJump(nilJump, span)
		c.emit(span.StartPos.Line, bytecode.POP)
		c.compileNode(node.Right)
		c.emitSetterCall(attr.AttributeName, span)

		// if not nil
		c.patchJump(nonNilJump, span)
	default:
		c.Errors.AddFailure(fmt.Sprintf("unknown binary operator: %s", node.Op.String()), c.newLocation(node.Span()))
	}
}

func (c *Compiler) instanceVariableAssignment(node *ast.AssignmentExpressionNode, ivar *ast.InstanceVariableNode) {
	switch c.mode {
	case topLevelMode:
		c.Errors.AddFailure(
			"instance variables cannot be set in the top level",
			c.newLocation(node.Span()),
		)
	}

	ivarSymbol := value.ToSymbol(ivar.Value)
	switch node.Op.Type {
	case token.EQUAL_OP:
		c.compileNode(node.Right)
		c.emitSetInstanceVariable(ivarSymbol, ivar.Span())
	case token.OR_OR_EQUAL:
		span := node.Span()
		// Read the current value
		c.emitGetInstanceVariable(ivarSymbol, span)

		jump := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF_NP)

		// if falsy
		c.emit(span.StartPos.Line, bytecode.POP)
		c.compileNode(node.Right)
		c.emitSetInstanceVariable(ivarSymbol, ivar.Span())

		// if truthy
		c.patchJump(jump, span)
	case token.AND_AND_EQUAL:
		span := node.Span()
		// Read the current value
		c.emitGetInstanceVariable(ivarSymbol, span)

		jump := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS_NP)

		// if truthy
		c.emit(span.StartPos.Line, bytecode.POP)
		c.compileNode(node.Right)
		c.emitSetInstanceVariable(ivarSymbol, ivar.Span())

		// if falsy
		c.patchJump(jump, span)
	case token.QUESTION_QUESTION_EQUAL:
		span := node.Span()
		// Read the current value
		c.emitGetInstanceVariable(ivarSymbol, span)

		nilJump := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF_NIL_NP)
		nonNilJump := c.emitJump(span.StartPos.Line, bytecode.JUMP)

		// if nil
		c.patchJump(nilJump, span)
		c.emit(span.StartPos.Line, bytecode.POP)
		c.compileNode(node.Right)
		c.emitSetInstanceVariable(ivarSymbol, ivar.Span())

		// if not nil
		c.patchJump(nonNilJump, span)
	default:
		c.Errors.AddFailure(fmt.Sprintf("unknown binary operator: %s", node.Op.String()), c.newLocation(node.Span()))
	}
}

func (c *Compiler) complexSubscriptAssignment(subscript *ast.SubscriptExpressionNode, valueNode ast.ExpressionNode, opcode bytecode.OpCode, span *position.Span) {
	c.compileNode(subscript.Receiver)
	c.compileNode(subscript.Key)
	c.emit(span.EndPos.Line, bytecode.DUP_N, 2)
	c.emit(span.EndPos.Line, bytecode.SUBSCRIPT)

	c.compileNode(valueNode)
	c.emit(span.StartPos.Line, opcode)
	c.emit(span.EndPos.Line, bytecode.SUBSCRIPT_SET)
}

func (c *Compiler) subscriptAssignment(node *ast.AssignmentExpressionNode, subscript *ast.SubscriptExpressionNode) {
	switch node.Op.Type {
	case token.EQUAL_OP:
		c.compileNode(subscript.Receiver)
		c.compileNode(subscript.Key)
		c.compileNode(node.Right)
		c.emit(node.Span().EndPos.Line, bytecode.SUBSCRIPT_SET)
	case token.OR_OR_EQUAL:
		span := node.Span()
		// Read the current value
		c.compileNode(subscript.Receiver)
		c.compileNode(subscript.Key)
		c.emit(node.Span().EndPos.Line, bytecode.DUP_N, 2)
		c.emit(node.Span().EndPos.Line, bytecode.SUBSCRIPT)

		jump := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF_NP)

		// if falsy
		c.emit(span.StartPos.Line, bytecode.POP)
		c.compileNode(node.Right)
		c.emit(node.Span().EndPos.Line, bytecode.SUBSCRIPT_SET)

		// if truthy
		c.patchJump(jump, span)
		c.emit(span.StartPos.Line, bytecode.POP_N_SKIP_ONE, 2)
	case token.AND_AND_EQUAL:
		span := node.Span()
		// Read the current value
		c.compileNode(subscript.Receiver)
		c.compileNode(subscript.Key)
		c.emit(node.Span().EndPos.Line, bytecode.DUP_N, 2)
		c.emit(node.Span().EndPos.Line, bytecode.SUBSCRIPT)

		jump := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS_NP)

		// if truthy
		c.emit(span.StartPos.Line, bytecode.POP)
		c.compileNode(node.Right)
		c.emit(node.Span().EndPos.Line, bytecode.SUBSCRIPT_SET)

		// if falsy
		c.patchJump(jump, span)
		c.emit(span.StartPos.Line, bytecode.POP_N_SKIP_ONE, 2)
	case token.QUESTION_QUESTION_EQUAL:
		span := node.Span()
		// Read the current value
		c.compileNode(subscript.Receiver)
		c.compileNode(subscript.Key)
		c.emit(node.Span().EndPos.Line, bytecode.DUP_N, 2)
		c.emit(node.Span().EndPos.Line, bytecode.SUBSCRIPT)

		nilJump := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF_NIL_NP)
		nonNilJump := c.emitJump(span.StartPos.Line, bytecode.JUMP)

		// if nil
		c.patchJump(nilJump, span)
		c.emit(span.StartPos.Line, bytecode.POP)
		c.compileNode(node.Right)
		c.emit(node.Span().EndPos.Line, bytecode.SUBSCRIPT_SET)

		// if not nil
		c.patchJump(nonNilJump, span)
		c.emit(span.StartPos.Line, bytecode.POP_N_SKIP_ONE, 2)
	default:
		c.Errors.AddFailure(fmt.Sprintf("unknown binary operator: %s", node.Op.String()), c.newLocation(node.Span()))
	}
}

func (c *Compiler) compileAssignmentExpressionNode(node *ast.AssignmentExpressionNode) {
	switch n := node.Left.(type) {
	case *ast.PublicIdentifierNode:
		c.localVariableAssignment(n.Value, node.Op, node.Right, node.Span())
	case *ast.PrivateIdentifierNode:
		c.localVariableAssignment(n.Value, node.Op, node.Right, node.Span())
	case *ast.SubscriptExpressionNode:
		c.subscriptAssignment(node, n)
	case *ast.InstanceVariableNode:
		c.instanceVariableAssignment(node, n)
	case *ast.AttributeAccessNode:
		c.attributeAssignment(node, n)
	default:
		c.Errors.AddFailure(
			fmt.Sprintf("cannot assign to: %T", node.Left),
			c.newLocation(node.Span()),
		)
	}
}

// Return the offset of the next instruction.
func (c *Compiler) nextInstructionOffset() int {
	return len(c.Bytecode.Instructions)
}

func (c *Compiler) setLocalWithoutValue(name string, span *position.Span) {
	if local, ok := c.resolveLocal(name, span); ok {
		c.emitSetLocal(span.StartPos.Line, local.index)
	} else if upvalue, ok := c.resolveUpvalue(name, span); ok {
		c.emitSetUpvalue(span.StartPos.Line, upvalue.index)
	}
}

func (c *Compiler) setLocal(name string, valueNode ast.ExpressionNode, span *position.Span) {
	c.compileNode(valueNode)
	c.setLocalWithoutValue(name, span)
}

func (c *Compiler) localVariableAssignment(name string, operator *token.Token, right ast.ExpressionNode, span *position.Span) {
	switch operator.Type {
	case token.OR_OR_EQUAL:
		c.compileLocalVariableAccess(name, span)
		jump := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF_NP)

		// if falsy
		c.emit(span.StartPos.Line, bytecode.POP)
		c.setLocal(name, right, span)

		// if truthy
		c.patchJump(jump, span)
	case token.AND_AND_EQUAL:
		c.compileLocalVariableAccess(name, span)
		jump := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS_NP)

		// if truthy
		c.emit(span.StartPos.Line, bytecode.POP)
		c.setLocal(name, right, span)

		// if falsy
		c.patchJump(jump, span)
	case token.QUESTION_QUESTION_EQUAL:
		c.compileLocalVariableAccess(name, span)
		nilJump := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF_NIL_NP)
		nonNilJump := c.emitJump(span.StartPos.Line, bytecode.JUMP)

		// if nil
		c.patchJump(nilJump, span)
		c.emit(span.StartPos.Line, bytecode.POP)
		c.setLocal(name, right, span)

		// if not nil
		c.patchJump(nonNilJump, span)
	case token.EQUAL_OP:
		c.setLocal(name, right, span)
	case token.COLON_EQUAL:
		c.compileNode(right)
		local := c.defineLocal(name, span)
		if local == nil {
			return
		}
		c.emitSetLocal(span.StartPos.Line, local.index)
	default:
		c.Errors.AddFailure(
			fmt.Sprintf("assignment using this operator has not been implemented: %s", operator.Type.String()),
			c.newLocation(span),
		)
		return
	}
}

func (c *Compiler) compileInstanceVariableAccess(name string, span *position.Span) {
	c.emitGetInstanceVariable(value.ToSymbol(name), span)
}

func (c *Compiler) compileLocalVariableAccess(name string, span *position.Span) (*local, *upvalue, bool) {
	if local, ok := c.resolveLocal(name, span); ok {
		c.emitGetLocal(span.StartPos.Line, local.index)
		return local, nil, true
	} else if upvalue, ok := c.resolveUpvalue(name, span); ok {
		local := upvalue.local
		c.emitGetUpvalue(span.StartPos.Line, upvalue.index)
		return local, upvalue, true
	}

	return nil, nil, false
}

// Resolve an upvalue from an outer context and get its index.
func (c *Compiler) resolveUpvalue(name string, span *position.Span) (*upvalue, bool) {
	parent := c.Parent
	if parent == nil {
		return nil, false
	}
	local, ok := parent.resolveLocal(name, span)
	if ok {
		return c.addUpvalue(local, local.index, true, span), true
	}

	upvalue, ok := parent.resolveUpvalue(name, span)
	if ok {
		return c.addUpvalue(upvalue.local, upvalue.index, false, span), true
	}

	return nil, false
}

func (c *Compiler) addUpvalue(local *local, upIndex uint16, isLocal bool, span *position.Span) *upvalue {
	for _, upvalue := range c.upvalues {
		if upvalue.upIndex == upIndex && upvalue.isLocal == isLocal {
			return upvalue
		}
	}

	if len(c.upvalues) > math.MaxUint16 {
		c.Errors.AddFailure(
			fmt.Sprintf("upvalue limit reached: %d", math.MaxUint16),
			c.newLocation(span),
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

func (c *Compiler) compileModifierExpressionNode(label string, node *ast.ModifierNode) {
	switch node.Modifier.Type {
	case token.IF:
		c.compileModifierIfExpression(false, node.Right, node.Left, nil, node.Span())
	case token.UNLESS:
		c.compileModifierIfExpression(true, node.Right, node.Left, nil, node.Span())
	case token.WHILE:
		c.modifierWhileExpression(label, node)
	case token.UNTIL:
		c.modifierUntilExpression(label, node)
	default:
		c.Errors.AddFailure(
			fmt.Sprintf("illegal modifier: %s", node.Modifier.StringValue()),
			c.newLocation(node.Span()),
		)
	}
}

func (c *Compiler) compileModifierIfExpression(unless bool, condition, then, els ast.ExpressionNode, span *position.Span) {
	var elsFunc func()
	if els != nil {
		elsFunc = func() {
			c.compileNode(els)
		}
	}
	var jumpOp bytecode.OpCode
	if unless {
		jumpOp = bytecode.JUMP_IF
	} else {
		jumpOp = bytecode.JUMP_UNLESS
	}
	c.compileIfWithConditionExpression(
		jumpOp,
		condition,
		func() {
			c.compileNode(then)
		},
		elsFunc,
		span,
	)
}

func (c *Compiler) compileIfExpression(unless bool, condition ast.ExpressionNode, then, els []ast.StatementNode, span *position.Span) {
	var elsFunc func()
	if els != nil {
		elsFunc = func() {
			c.compileStatements(els, span)
		}
	}

	var jumpOp bytecode.OpCode
	if unless {
		jumpOp = bytecode.JUMP_IF
	} else {
		jumpOp = bytecode.JUMP_UNLESS
	}

	c.compileIfWithConditionExpression(
		jumpOp,
		condition,
		func() {
			c.compileStatements(then, span)
		},
		elsFunc,
		span,
	)
}

func (c *Compiler) compileIf(jumpOp bytecode.OpCode, condition, then, els func(), span *position.Span) {
	c.enterScope("", defaultScopeType)
	condition()

	thenJumpOffset := c.emitJump(span.StartPos.Line, jumpOp)

	then()
	c.leaveScope(span.StartPos.Line)

	elseJumpOffset := c.emitJump(span.StartPos.Line, bytecode.JUMP)

	c.patchJump(thenJumpOffset, span)

	if els != nil {
		c.enterScope("", defaultScopeType)
		els()
		c.leaveScope(span.StartPos.Line)
	} else {
		c.emit(span.StartPos.Line, bytecode.NIL)
	}
	c.patchJump(elseJumpOffset, span)
}

func (c *Compiler) compileIfWithConditionExpression(jumpOp bytecode.OpCode, condition ast.ExpressionNode, then, els func(), span *position.Span) {
	if result := resolve(condition); !result.IsUndefined() {
		// if gets optimised away
		c.enterScope("", defaultScopeType)
		defer c.leaveScope(span.StartPos.Line)

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
				c.emit(span.StartPos.Line, bytecode.NIL)
				return
			}
			then()
			return
		}

		if els == nil {
			c.emit(span.StartPos.Line, bytecode.NIL)
			return
		}
		els()
		return
	}

	if c.optimiseIf(jumpOp, condition, then, els, span) {
		return
	}

	c.compileIf(
		jumpOp,
		func() {
			c.compileNode(condition)
		},
		then,
		els,
		span,
	)
}

func (c *Compiler) optimiseIf(jumpOp bytecode.OpCode, condition ast.ExpressionNode, then, els func(), span *position.Span) bool {
	cond, ok := condition.(*ast.BinaryExpressionNode)
	if !ok {
		return false
	}

	switch cond.Op.Type {
	case token.LESS_EQUAL:
		return c.optimiseIfLessEqual(jumpOp, cond, then, els, span)
	case token.LESS:
		return c.optimiseIfLess(jumpOp, cond, then, els, span)
	case token.GREATER_EQUAL:
		return c.optimiseIfGreaterEqual(jumpOp, cond, then, els, span)
	case token.GREATER:
		return c.optimiseIfGreater(jumpOp, cond, then, els, span)
	case token.EQUAL_EQUAL:
		return c.optimiseIfEqual(jumpOp, cond, then, els, span)
	case token.NOT_EQUAL:
		return c.optimiseIfNotEqual(jumpOp, cond, then, els, span)
	}

	return false
}

func (c *Compiler) optimiseIfNotEqual(jumpOp bytecode.OpCode, condition *ast.BinaryExpressionNode, then, els func(), span *position.Span) bool {
	leftType := c.typeOf(condition.Left)
	rightType := c.typeOf(condition.Right)

	if jumpOp == bytecode.JUMP_UNLESS {
		if c.checker.IsSubtype(leftType, c.checker.StdNil()) {
			c.compileIf(
				bytecode.JUMP_IF_NIL,
				func() {
					c.compileNode(condition.Right)
				},
				then,
				els,
				span,
			)
			return true
		}
		if c.checker.IsSubtype(rightType, c.checker.StdNil()) {
			c.compileIf(
				bytecode.JUMP_IF_NIL,
				func() {
					c.compileNode(condition.Left)
				},
				then,
				els,
				span,
			)
			return true
		}
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			c.compileIf(
				bytecode.JUMP_IF_IEQ,
				func() {
					c.compileNode(condition.Left)
					c.compileNode(condition.Right)
				},
				then,
				els,
				span,
			)
			return true
		}
		if c.checker.IsSubtype(rightType, c.checker.StdInt()) {
			c.compileIf(
				bytecode.JUMP_IF_IEQ,
				func() {
					c.compileNode(condition.Right)
					c.compileNode(condition.Left)
				},
				then,
				els,
				span,
			)
			return true
		}
	}
	if jumpOp == bytecode.JUMP_IF {
		if c.checker.IsSubtype(leftType, c.checker.StdNil()) {
			c.compileIf(
				bytecode.JUMP_UNLESS_NIL,
				func() {
					c.compileNode(condition.Right)
				},
				then,
				els,
				span,
			)
			return true
		}
		if c.checker.IsSubtype(rightType, c.checker.StdNil()) {
			c.compileIf(
				bytecode.JUMP_UNLESS_NIL,
				func() {
					c.compileNode(condition.Left)
				},
				then,
				els,
				span,
			)
			return true
		}
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			c.compileIf(
				bytecode.JUMP_UNLESS_IEQ,
				func() {
					c.compileNode(condition.Left)
					c.compileNode(condition.Right)
				},
				then,
				els,
				span,
			)
			return true
		}
		if c.checker.IsSubtype(rightType, c.checker.StdInt()) {
			c.compileIf(
				bytecode.JUMP_UNLESS_IEQ,
				func() {
					c.compileNode(condition.Right)
					c.compileNode(condition.Left)
				},
				then,
				els,
				span,
			)
			return true
		}
	}

	return false
}

func (c *Compiler) optimiseIfEqual(jumpOp bytecode.OpCode, condition *ast.BinaryExpressionNode, then, els func(), span *position.Span) bool {
	leftType := c.typeOf(condition.Left)
	rightType := c.typeOf(condition.Right)

	if jumpOp == bytecode.JUMP_UNLESS {
		if c.checker.IsSubtype(leftType, c.checker.StdNil()) {
			c.compileIf(
				bytecode.JUMP_UNLESS_NIL,
				func() {
					c.compileNode(condition.Right)
				},
				then,
				els,
				span,
			)
			return true
		}
		if c.checker.IsSubtype(rightType, c.checker.StdNil()) {
			c.compileIf(
				bytecode.JUMP_UNLESS_NIL,
				func() {
					c.compileNode(condition.Left)
				},
				then,
				els,
				span,
			)
			return true
		}
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			c.compileIf(
				bytecode.JUMP_UNLESS_IEQ,
				func() {
					c.compileNode(condition.Left)
					c.compileNode(condition.Right)
				},
				then,
				els,
				span,
			)
			return true
		}
		if c.checker.IsSubtype(rightType, c.checker.StdInt()) {
			c.compileIf(
				bytecode.JUMP_UNLESS_IEQ,
				func() {
					c.compileNode(condition.Right)
					c.compileNode(condition.Left)
				},
				then,
				els,
				span,
			)
			return true
		}
	}
	if jumpOp == bytecode.JUMP_IF {
		if c.checker.IsSubtype(leftType, c.checker.StdNil()) {
			c.compileIf(
				bytecode.JUMP_IF_NIL,
				func() {
					c.compileNode(condition.Right)
				},
				then,
				els,
				span,
			)
			return true
		}
		if c.checker.IsSubtype(rightType, c.checker.StdNil()) {
			c.compileIf(
				bytecode.JUMP_IF_NIL,
				func() {
					c.compileNode(condition.Left)
				},
				then,
				els,
				span,
			)
			return true
		}
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			c.compileIf(
				bytecode.JUMP_IF_IEQ,
				func() {
					c.compileNode(condition.Left)
					c.compileNode(condition.Right)
				},
				then,
				els,
				span,
			)
			return true
		}
		if c.checker.IsSubtype(rightType, c.checker.StdInt()) {
			c.compileIf(
				bytecode.JUMP_IF_IEQ,
				func() {
					c.compileNode(condition.Right)
					c.compileNode(condition.Left)
				},
				then,
				els,
				span,
			)
			return true
		}
	}

	return false
}

func (c *Compiler) optimiseIfGreater(jumpOp bytecode.OpCode, condition *ast.BinaryExpressionNode, then, els func(), span *position.Span) bool {
	leftType := c.typeOf(condition.Left)
	rightType := c.typeOf(condition.Right)

	if jumpOp == bytecode.JUMP_UNLESS {
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			c.compileIf(
				bytecode.JUMP_UNLESS_IGT,
				func() {
					c.compileNode(condition.Left)
					c.compileNode(condition.Right)
				},
				then,
				els,
				span,
			)
			return true
		}
		if c.checker.IsSubtype(rightType, c.checker.StdInt()) {
			c.compileIf(
				bytecode.JUMP_UNLESS_ILT,
				func() {
					c.compileNode(condition.Right)
					c.compileNode(condition.Left)
				},
				then,
				els,
				span,
			)
			return true
		}
	}
	if jumpOp == bytecode.JUMP_IF {
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			c.compileIf(
				bytecode.JUMP_UNLESS_ILE,
				func() {
					c.compileNode(condition.Left)
					c.compileNode(condition.Right)
				},
				then,
				els,
				span,
			)
			return true
		}
	}

	return false
}

func (c *Compiler) optimiseIfGreaterEqual(jumpOp bytecode.OpCode, condition *ast.BinaryExpressionNode, then, els func(), span *position.Span) bool {
	leftType := c.typeOf(condition.Left)
	rightType := c.typeOf(condition.Right)

	if jumpOp == bytecode.JUMP_UNLESS {
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			c.compileIf(
				bytecode.JUMP_UNLESS_IGE,
				func() {
					c.compileNode(condition.Left)
					c.compileNode(condition.Right)
				},
				then,
				els,
				span,
			)
			return true
		}
		// Reverse only when leftType is subtype of BuiltinComparable
		if c.checker.IsSubtype(rightType, c.checker.StdInt()) {
			c.compileIf(
				bytecode.JUMP_UNLESS_ILE,
				func() {
					c.compileNode(condition.Right)
					c.compileNode(condition.Left)
				},
				then,
				els,
				span,
			)
			return true
		}
	}
	if jumpOp == bytecode.JUMP_IF {
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			c.compileIf(
				bytecode.JUMP_UNLESS_ILT,
				func() {
					c.compileNode(condition.Left)
					c.compileNode(condition.Right)
				},
				then,
				els,
				span,
			)
			return true
		}
	}

	return false
}

func (c *Compiler) optimiseIfLess(jumpOp bytecode.OpCode, condition *ast.BinaryExpressionNode, then, els func(), span *position.Span) bool {
	leftType := c.typeOf(condition.Left)
	rightType := c.typeOf(condition.Right)

	if jumpOp == bytecode.JUMP_UNLESS {
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			c.compileIf(
				bytecode.JUMP_UNLESS_ILT,
				func() {
					c.compileNode(condition.Left)
					c.compileNode(condition.Right)
				},
				then,
				els,
				span,
			)
			return true
		}
		if c.checker.IsSubtype(rightType, c.checker.StdInt()) {
			c.compileIf(
				bytecode.JUMP_UNLESS_IGT,
				func() {
					c.compileNode(condition.Right)
					c.compileNode(condition.Left)
				},
				then,
				els,
				span,
			)
			return true
		}
	}
	if jumpOp == bytecode.JUMP_IF {
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			c.compileIf(
				bytecode.JUMP_UNLESS_IGE,
				func() {
					c.compileNode(condition.Left)
					c.compileNode(condition.Right)
				},
				then,
				els,
				span,
			)
			return true
		}
	}

	return false
}

func (c *Compiler) optimiseIfLessEqual(jumpOp bytecode.OpCode, condition *ast.BinaryExpressionNode, then, els func(), span *position.Span) bool {
	leftType := c.typeOf(condition.Left)
	rightType := c.typeOf(condition.Right)

	if jumpOp == bytecode.JUMP_UNLESS {
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			c.compileIf(
				bytecode.JUMP_UNLESS_ILE,
				func() {
					c.compileNode(condition.Left)
					c.compileNode(condition.Right)
				},
				then,
				els,
				span,
			)
			return true
		}
		if c.checker.IsSubtype(rightType, c.checker.StdInt()) {
			c.compileIf(
				bytecode.JUMP_UNLESS_IGE,
				func() {
					c.compileNode(condition.Right)
					c.compileNode(condition.Left)
				},
				then,
				els,
				span,
			)
			return true
		}
	}
	if jumpOp == bytecode.JUMP_IF {
		if c.checker.IsSubtype(leftType, c.checker.StdInt()) {
			c.compileIf(
				bytecode.JUMP_UNLESS_IGT,
				func() {
					c.compileNode(condition.Left)
					c.compileNode(condition.Right)
				},
				then,
				els,
				span,
			)
			return true
		}
	}

	return false
}

func (c *Compiler) compileValueDeclarationNode(node *ast.ValueDeclarationNode) {
	initialised := node.Initialiser != nil

	if initialised {
		c.compileNode(node.Initialiser)
	}
	local := c.defineLocal(node.Name, node.Span())
	if local == nil {
		return
	}
	if initialised {
		c.emitSetLocal(node.Span().StartPos.Line, local.index)
	} else {
		c.emit(node.Span().StartPos.Line, bytecode.NIL)
	}
}

func (c *Compiler) compileReturnExpressionNode(node *ast.ReturnExpressionNode) {
	span := node.Span()
	if node.Value != nil {
		c.emitReturn(span, node.Value)
	} else {
		c.emit(span.StartPos.Line, bytecode.NIL)
		c.emitReturn(span, nil)
	}
}

func (c *Compiler) compileNilSafeSubscriptExpressionNode(node *ast.NilSafeSubscriptExpressionNode) {
	if c.resolveAndEmit(node) {
		return
	}

	c.compileIfWithConditionExpression(
		bytecode.JUMP_IF_NIL,
		node.Receiver,
		func() {
			c.compileNode(node.Receiver)
			c.compileNode(node.Key)
			c.emit(node.Span().EndPos.Line, bytecode.SUBSCRIPT)
		},
		func() {
			c.emit(node.Span().EndPos.Line, bytecode.NIL)
		},
		node.Span(),
	)
}

func (c *Compiler) relationalPattern(pattern ast.Node, opcode bytecode.OpCode) {
	span := pattern.Span()

	c.compileIf(
		bytecode.JUMP_UNLESS,
		func() {
			c.emit(span.StartPos.Line, bytecode.DUP)
			c.compileNode(pattern)
			c.emit(span.StartPos.Line, bytecode.DUP_N, 2)
			c.emit(span.StartPos.Line, bytecode.SWAP)
			c.emit(span.StartPos.Line, bytecode.GET_CLASS)
			c.emit(span.StartPos.Line, bytecode.IS_A)
		},
		func() {
			c.emit(span.StartPos.Line, opcode)
		},
		func() {
			c.emit(span.StartPos.Line, bytecode.POP_N, 2)
			c.emit(span.StartPos.Line, bytecode.FALSE)
		},
		span,
	)
}

func (c *Compiler) literalPattern(pattern ast.Node, opcode bytecode.OpCode) {
	span := pattern.Span()
	c.emit(span.StartPos.Line, bytecode.DUP)
	c.compileNode(pattern)
	c.emit(span.StartPos.Line, opcode)
}

func (c *Compiler) pattern(pattern ast.PatternNode) {
	span := pattern.Span()
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
		c.emit(span.StartPos.Line, bytecode.DUP)
		c.compileRangeLiteralNode(pat)
		c.emit(span.StartPos.Line, bytecode.SWAP)
		callInfo := value.NewCallSiteInfo(symbol.S_contains, 1)
		c.emitCallMethod(callInfo, span)
	case *ast.PublicIdentifierNode:
		switch c.mode {
		case valuePatternDeclarationNode:
			c.defineLocal(pat.Value, span)
		default:
			c.defineLocalOverrideCurrentScope(pat.Value, span)
		}
		c.setLocalWithoutValue(pat.Value, span)
		c.emit(span.StartPos.Line, bytecode.TRUE)
	case *ast.PrivateIdentifierNode:
		switch c.mode {
		case valuePatternDeclarationNode:
			c.defineLocal(pat.Value, span)
		default:
			c.defineLocalOverrideCurrentScope(pat.Value, span)
		}
		c.setLocalWithoutValue(pat.Value, span)
		c.emit(span.StartPos.Line, bytecode.TRUE)
	case *ast.ObjectPatternNode:
		c.objectPattern(pat)
	case *ast.AsPatternNode:
		c.asPattern(pat)
	case *ast.UninterpolatedRegexLiteralNode, *ast.InterpolatedRegexLiteralNode:
		c.emit(span.StartPos.Line, bytecode.DUP)
		c.compileNode(pat)
		c.emit(span.StartPos.Line, bytecode.SWAP)
		callInfo := value.NewCallSiteInfo(matchesSymbol, 1)
		c.emitCallMethod(callInfo, span)
	case *ast.UnaryExpressionNode:
		c.unaryPattern(pat)
	case *ast.BinaryPatternNode:
		c.binaryPattern(pat)
	case *ast.MapPatternNode:
		c.mapOrRecordPattern(pat.Span(), pat.Elements, true)
	case *ast.RecordPatternNode:
		c.mapOrRecordPattern(pat.Span(), pat.Elements, false)
	case *ast.SetPatternNode:
		c.setPattern(pat.Span(), pat.Elements)
	case *ast.ListPatternNode:
		c.listOrTuplePattern(pat.Span(), pat.Elements, true)
	case *ast.TuplePatternNode:
		c.listOrTuplePattern(pat.Span(), pat.Elements, false)
	case *ast.WordArrayListLiteralNode, *ast.SymbolArrayListLiteralNode, *ast.BinArrayListLiteralNode, *ast.HexArrayListLiteralNode,
		*ast.WordArrayTupleLiteralNode, *ast.SymbolArrayTupleLiteralNode, *ast.BinArrayTupleLiteralNode, *ast.HexArrayTupleLiteralNode,
		*ast.WordHashSetLiteralNode, *ast.SymbolHashSetLiteralNode, *ast.BinHashSetLiteralNode, *ast.HexHashSetLiteralNode:
		c.specialCollectionPattern(pat)
	default:
		c.Errors.AddFailure(
			fmt.Sprintf("compilation of this pattern has not been implemented: %T", pattern),
			c.newLocation(span),
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
	span := pat.Span()
	var op bytecode.OpCode
	switch pat.Op.Type {
	case token.OR_OR:
		op = bytecode.JUMP_IF_NP
	case token.AND_AND:
		op = bytecode.JUMP_UNLESS_NP
	default:
		panic(fmt.Sprintf("invalid binary pattern operator: %s", pat.Op.Type.String()))
	}

	c.pattern(pat.Left)
	jump := c.emitJump(span.StartPos.Line, op)

	// branch one
	c.emit(span.StartPos.Line, bytecode.POP)
	c.pattern(pat.Right)

	// branch two
	c.patchJump(jump, span)
}

func (c *Compiler) asPattern(node *ast.AsPatternNode) {
	span := node.Span()
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
		c.defineLocal(varName, span)
	default:
		c.defineLocalOverrideCurrentScope(varName, span)
	}
	c.setLocalWithoutValue(varName, span)
	c.pattern(node.Pattern)
}

func (c *Compiler) identifierObjectPatternAttribute(name string, span *position.Span) {
	c.emit(span.StartPos.Line, bytecode.DUP)
	callInfo := value.NewCallSiteInfo(value.ToSymbol(name), 0)
	c.emitCallMethod(callInfo, span)

	var identVar *local
	switch c.mode {
	case valuePatternDeclarationNode:
		identVar = c.defineLocal(name, span)
	default:
		identVar = c.defineLocalOverrideCurrentScope(name, span)
	}
	c.emitSetLocal(span.StartPos.Line, identVar.index)
	c.emit(span.StartPos.Line, bytecode.POP)
}

func (c *Compiler) objectPattern(node *ast.ObjectPatternNode) {
	var jumpsToPatch []int
	c.enterPattern()

	span := node.Span()
	c.emit(node.ObjectType.Span().StartPos.Line, bytecode.DUP)
	c.compileNode(node.ObjectType)
	c.emit(node.ObjectType.Span().StartPos.Line, bytecode.IS_A)

	jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS_NP)
	jumpsToPatch = append(jumpsToPatch, jmp)
	c.emit(span.StartPos.Line, bytecode.POP)

	for _, attr := range node.Attributes {
		span := attr.Span()
		switch e := attr.(type) {
		case *ast.SymbolKeyValuePatternNode:
			c.emit(span.StartPos.Line, bytecode.DUP)
			callInfo := value.NewCallSiteInfo(value.ToSymbol(e.Key), 0)
			c.emitCallMethod(callInfo, span)

			c.pattern(e.Value)
			c.emit(span.StartPos.Line, bytecode.POP_SKIP_ONE)
			jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS_NP)
			jumpsToPatch = append(jumpsToPatch, jmp)
			c.emit(span.StartPos.Line, bytecode.POP)
		case *ast.PublicIdentifierNode:
			c.identifierObjectPatternAttribute(e.Value, span)
		case *ast.PrivateIdentifierNode:
			c.identifierObjectPatternAttribute(e.Value, span)
		default:
			c.Errors.AddFailure(
				fmt.Sprintf("invalid object pattern attribute: %T", attr),
				c.newLocation(span),
			)
		}
	}

	// leave true as the result of the happy path
	c.emit(span.StartPos.Line, bytecode.TRUE)

	// leave false on the stack from the falsy if that jumped here
	for _, jmp := range jumpsToPatch {
		c.patchJump(jmp, span)
	}
	c.leavePattern()
}

func (c *Compiler) specialCollectionPattern(node ast.PatternNode) {
	span := node.Span()
	c.emit(span.StartPos.Line, bytecode.DUP)
	switch node.(type) {
	case *ast.WordArrayListLiteralNode, *ast.SymbolArrayListLiteralNode, *ast.BinArrayListLiteralNode, *ast.HexArrayListLiteralNode:
		c.emitValue(value.Ref(value.ListMixin), span)
	case *ast.WordArrayTupleLiteralNode, *ast.SymbolArrayTupleLiteralNode, *ast.BinArrayTupleLiteralNode, *ast.HexArrayTupleLiteralNode:
		c.emitValue(value.Ref(value.TupleMixin), span)
	case *ast.WordHashSetLiteralNode, *ast.SymbolHashSetLiteralNode, *ast.BinHashSetLiteralNode, *ast.HexHashSetLiteralNode:
		c.emitValue(value.Ref(value.SetMixin), span)
	default:
		panic(fmt.Sprintf("invalid special collection pattern node: %#v", node))
	}
	c.emit(span.StartPos.Line, bytecode.IS_A)

	jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS_NP)
	c.emit(span.StartPos.Line, bytecode.POP)

	c.emit(span.StartPos.Line, bytecode.DUP)
	c.compileNode(node)
	c.emit(span.StartPos.Line, bytecode.LAX_EQUAL)

	// leave false on the stack from the falsy if that jumped here
	c.patchJump(jmp, span)
}

func (c *Compiler) identifierMapPatternElement(name string, span *position.Span) {
	c.emit(span.StartPos.Line, bytecode.DUP)
	c.emitValue(value.ToSymbol(name).ToValue(), span)
	c.emit(span.StartPos.Line, bytecode.SUBSCRIPT)
	var identVar *local
	switch c.mode {
	case valuePatternDeclarationNode:
		identVar = c.defineLocal(name, span)
	default:
		identVar = c.defineLocalOverrideCurrentScope(name, span)
	}
	if identVar == nil {
		return
	}
	c.emitSetLocal(span.StartPos.Line, identVar.index)
	c.emit(span.StartPos.Line, bytecode.POP)
}

func (c *Compiler) mapOrRecordPattern(span *position.Span, elements []ast.PatternNode, isMap bool) {
	var jumpsToPatch []int
	c.enterPattern()

	c.emit(span.StartPos.Line, bytecode.DUP)
	if isMap {
		c.emitValue(value.Ref(value.MapMixin), span)
	} else {
		c.emitValue(value.Ref(value.RecordMixin), span)
	}
	c.emit(span.StartPos.Line, bytecode.IS_A)

	jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS_NP)
	jumpsToPatch = append(jumpsToPatch, jmp)
	c.emit(span.StartPos.Line, bytecode.POP)

	for _, element := range elements {
		span := element.Span()
		switch e := element.(type) {
		case *ast.SymbolKeyValuePatternNode:
			c.emit(span.StartPos.Line, bytecode.DUP)
			c.emitValue(value.ToSymbol(e.Key).ToValue(), span)
			c.emit(span.StartPos.Line, bytecode.SUBSCRIPT)

			c.pattern(e.Value)
			c.emit(span.StartPos.Line, bytecode.POP_SKIP_ONE)
			jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS_NP)
			jumpsToPatch = append(jumpsToPatch, jmp)
			c.emit(span.StartPos.Line, bytecode.POP)
		case *ast.KeyValuePatternNode:
			c.emit(span.StartPos.Line, bytecode.DUP)
			c.compileNode(e.Key)
			c.emit(span.StartPos.Line, bytecode.SUBSCRIPT)

			c.pattern(e.Value)
			c.emit(span.StartPos.Line, bytecode.POP_SKIP_ONE)
			jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS_NP)
			jumpsToPatch = append(jumpsToPatch, jmp)
			c.emit(span.StartPos.Line, bytecode.POP)
		case *ast.PublicIdentifierNode:
			c.identifierMapPatternElement(e.Value, span)
		case *ast.PrivateIdentifierNode:
			c.identifierMapPatternElement(e.Value, span)
		default:
			c.Errors.AddFailure(
				fmt.Sprintf("invalid map pattern element: %T", element),
				c.newLocation(span),
			)
		}
	}

	// leave true as the result of the happy path
	c.emit(span.StartPos.Line, bytecode.TRUE)

	// leave false on the stack from the falsy if that jumped here
	for _, jmp := range jumpsToPatch {
		c.patchJump(jmp, span)
	}
	c.leavePattern()
}

func (c *Compiler) setPattern(span *position.Span, elements []ast.PatternNode) {
	var jumpsToPatch []int
	var subPatternElements []ast.PatternNode

	var restElementIsPresent bool
	for _, element := range elements {
		switch e := element.(type) {
		case *ast.RestPatternNode:
			if restElementIsPresent {
				c.Errors.AddFailure(
					"there should be only a single rest element",
					c.newLocation(element.Span()),
				)
			}
			restElementIsPresent = true
		default:
			subPatternElements = append(subPatternElements, e)
		}
	}
	c.enterPattern()

	c.emit(span.StartPos.Line, bytecode.DUP)
	c.emitValue(value.Ref(value.SetMixin), span)
	c.emit(span.StartPos.Line, bytecode.IS_A)

	jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS_NP)
	jumpsToPatch = append(jumpsToPatch, jmp)
	c.emit(span.StartPos.Line, bytecode.POP)

	c.emit(span.StartPos.Line, bytecode.DUP)
	callInfo := value.NewCallSiteInfo(symbol.L_length, 0)
	c.emitCallMethod(callInfo, span)

	if !restElementIsPresent {
		c.emitValue(value.SmallInt(len(subPatternElements)).ToValue(), span)
		c.emit(span.StartPos.Line, bytecode.EQUAL)
	} else {
		c.emitValue(value.SmallInt(len(subPatternElements)).ToValue(), span)
		c.emit(span.StartPos.Line, bytecode.GREATER_EQUAL)
	}

	jmp = c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS_NP)
	jumpsToPatch = append(jumpsToPatch, jmp)
	c.emit(span.StartPos.Line, bytecode.POP)

subPatternLoop:
	for _, element := range subPatternElements {
		switch element.(type) {
		case *ast.PrivateIdentifierNode:
			continue subPatternLoop
		}

		span := element.Span()
		c.emit(span.StartPos.Line, bytecode.DUP)
		c.compileNode(element)
		callInfo := value.NewCallSiteInfo(symbol.L_contains, 1)
		c.emitCallMethod(callInfo, span)

		jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS_NP)
		jumpsToPatch = append(jumpsToPatch, jmp)
		c.emit(span.StartPos.Line, bytecode.POP)
	}

	// leave true as the result of the happy path
	c.emit(span.StartPos.Line, bytecode.TRUE)

	// leave false on the stack from the falsy if that jumped here
	for _, jmp := range jumpsToPatch {
		c.patchJump(jmp, span)
	}
	c.leavePattern()
}

func (c *Compiler) listOrTuplePattern(span *position.Span, elements []ast.PatternNode, isList bool) {
	var jumpsToPatch []int

	var restVariableName string
	elementBeforeRestCount := -1
	for i, element := range elements {
		switch e := element.(type) {
		case *ast.RestPatternNode:
			if elementBeforeRestCount != -1 {
				c.Errors.AddFailure(
					"there should be only a single rest element",
					c.newLocation(element.Span()),
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
		c.emit(span.StartPos.Line, bytecode.UNDEFINED)
		c.emit(span.StartPos.Line, bytecode.UNDEFINED)
		c.emitNewArrayList(0, span)
		restListVar = c.defineLocal(restVariableName, span)
		c.emitSetLocal(span.StartPos.Line, restListVar.index)
		c.emit(span.StartPos.Line, bytecode.POP)
	}
	c.enterPattern()

	c.emit(span.StartPos.Line, bytecode.DUP)
	if isList {
		c.emitValue(value.Ref(value.ListMixin), span)
	} else {
		c.emitValue(value.Ref(value.TupleMixin), span)
	}
	c.emit(span.StartPos.Line, bytecode.IS_A)

	jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS_NP)
	jumpsToPatch = append(jumpsToPatch, jmp)
	c.emit(span.StartPos.Line, bytecode.POP)

	c.emit(span.StartPos.Line, bytecode.DUP)
	callInfo := value.NewCallSiteInfo(symbol.L_length, 0)
	c.emitCallMethod(callInfo, span)
	var lengthVar *local
	if elementBeforeRestCount != -1 {
		lengthVar = c.defineLocal(fmt.Sprintf("#!listPatternLength%d", c.patternNesting), span)
		c.emitSetLocal(span.StartPos.Line, lengthVar.index)
	}

	if elementBeforeRestCount == -1 {
		c.emitValue(value.SmallInt(len(elements)).ToValue(), span)
		c.emit(span.StartPos.Line, bytecode.EQUAL)
	} else {
		staticElementCount := elementBeforeRestCount + elementAfterRestCount
		c.emitValue(value.SmallInt(staticElementCount).ToValue(), span)
		c.emit(span.StartPos.Line, bytecode.GREATER_EQUAL)
	}

	jmp = c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS_NP)
	jumpsToPatch = append(jumpsToPatch, jmp)
	c.emit(span.StartPos.Line, bytecode.POP)

	elementsBeforeRest := elements
	if elementBeforeRestCount != -1 {
		elementsBeforeRest = elements[:elementBeforeRestCount]
	}
	for i, element := range elementsBeforeRest {
		span := element.Span()
		c.emit(span.StartPos.Line, bytecode.DUP)
		c.emitValue(value.SmallInt(i).ToValue(), element.Span())
		c.emit(span.StartPos.Line, bytecode.SUBSCRIPT)

		c.pattern(element)
		c.emit(span.StartPos.Line, bytecode.POP_SKIP_ONE)
		jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS_NP)
		jumpsToPatch = append(jumpsToPatch, jmp)
		c.emit(span.StartPos.Line, bytecode.POP)
	}

	if elementBeforeRestCount != -1 {
		iteratorVar := c.defineLocal(fmt.Sprintf("#!listPatternIterator%d", c.patternNesting), span)

		if restVariableName != "" {
			// adjust the length variable
			// length -= element_after_rest_count
			if elementAfterRestCount != 0 {
				c.emitGetLocal(span.StartPos.Line, lengthVar.index)
				c.emitValue(value.SmallInt(elementAfterRestCount).ToValue(), span)
				c.emit(span.StartPos.Line, bytecode.SUBTRACT)
				c.emitSetLocal(span.StartPos.Line, lengthVar.index)
				c.emit(span.StartPos.Line, bytecode.POP)
			}

			// create the iterator variable
			// i := element_before_rest_count
			c.emitValue(value.SmallInt(elementBeforeRestCount).ToValue(), span)
			c.emitSetLocal(span.StartPos.Line, iteratorVar.index)
			c.emit(span.StartPos.Line, bytecode.POP)

			// loop header
			// i < length
			loopStartOffset := c.nextInstructionOffset()
			c.emitGetLocal(span.StartPos.Line, iteratorVar.index)
			c.emitGetLocal(span.StartPos.Line, lengthVar.index)
			c.emit(span.StartPos.Line, bytecode.LESS)
			loopEndJump := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)

			// loop body
			c.emit(span.StartPos.Line, bytecode.DUP)
			c.emitGetLocal(span.StartPos.Line, iteratorVar.index)
			c.emit(span.StartPos.Line, bytecode.SUBSCRIPT)
			c.emitGetLocal(span.StartPos.Line, restListVar.index)
			c.emit(span.StartPos.Line, bytecode.SWAP)
			c.emit(span.StartPos.Line, bytecode.APPEND) // append to the list
			c.emit(span.StartPos.Line, bytecode.POP)
			// i++
			c.emitGetLocal(span.StartPos.Line, iteratorVar.index)
			c.emit(span.StartPos.Line, bytecode.INCREMENT)
			c.emitSetLocal(span.StartPos.Line, iteratorVar.index)
			c.emit(span.StartPos.Line, bytecode.POP)

			c.emitLoop(span, loopStartOffset)
			// loop end
			c.patchJump(loopEndJump, span)
		} else {
			// create the iterator variable
			// i := length - element_after_rest_count
			if elementAfterRestCount == 0 {
				c.emitGetLocal(span.StartPos.Line, lengthVar.index)
				c.emitSetLocal(span.StartPos.Line, iteratorVar.index)
				c.emit(span.StartPos.Line, bytecode.POP)
			} else {
				c.emitGetLocal(span.StartPos.Line, lengthVar.index)
				c.emitValue(value.SmallInt(elementAfterRestCount).ToValue(), span)
				c.emit(span.StartPos.Line, bytecode.SUBTRACT)
				c.emitSetLocal(span.StartPos.Line, iteratorVar.index)
				c.emit(span.StartPos.Line, bytecode.POP)
			}
		}

		elementsAfterRest := elements[elementBeforeRestCount+1:]
		for _, element := range elementsAfterRest {
			span := element.Span()
			c.emit(span.StartPos.Line, bytecode.DUP)
			c.emitGetLocal(span.StartPos.Line, iteratorVar.index)
			c.emit(span.StartPos.Line, bytecode.SUBSCRIPT)

			c.pattern(element)
			c.emit(span.StartPos.Line, bytecode.POP_SKIP_ONE)
			jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS_NP)
			jumpsToPatch = append(jumpsToPatch, jmp)
			c.emit(span.StartPos.Line, bytecode.POP)

			// i++
			c.emitGetLocal(span.StartPos.Line, iteratorVar.index)
			c.emit(span.StartPos.Line, bytecode.INCREMENT)
			c.emitSetLocal(span.StartPos.Line, iteratorVar.index)
			c.emit(span.StartPos.Line, bytecode.POP)
		}
	}

	// leave true as the result of the happy path
	c.emit(span.StartPos.Line, bytecode.TRUE)

	// leave false on the stack from the falsy if that jumped here
	for _, jmp := range jumpsToPatch {
		c.patchJump(jmp, span)
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

func (c *Compiler) compileSwitchExpressionNode(node *ast.SwitchExpressionNode) {
	span := node.Span()

	c.enterScope("", defaultScopeType)
	c.compileNode(node.Value)

	var jumpToEndOffsets []int

	for _, caseNode := range node.Cases {
		c.enterScope("", defaultScopeType)

		caseSpan := caseNode.Span()
		c.pattern(caseNode.Pattern)

		jumpOverBodyOffset := c.emitJump(caseSpan.StartPos.Line, bytecode.JUMP_UNLESS)

		c.emit(caseSpan.StartPos.Line, bytecode.POP)

		c.compileStatements(caseNode.Body, caseSpan)

		c.leaveScopeWithoutMutating(caseSpan.EndPos.Line)

		jumpToEndOffset := c.emitJump(caseSpan.EndPos.Line, bytecode.JUMP)
		jumpToEndOffsets = append(jumpToEndOffsets, jumpToEndOffset)

		c.patchJump(jumpOverBodyOffset, caseSpan)
		c.leaveScope(caseSpan.EndPos.Line)
	}

	c.emit(span.StartPos.Line, bytecode.POP)
	c.compileStatements(node.ElseBody, span)

	for _, offset := range jumpToEndOffsets {
		c.patchJump(offset, span)
	}

	c.leaveScope(span.EndPos.Line)
}

func (c *Compiler) compileSubscriptExpressionNode(node *ast.SubscriptExpressionNode) {
	if c.resolveAndEmit(node) {
		return
	}

	c.compileNode(node.Receiver)
	c.compileNode(node.Key)
	c.emit(node.Span().EndPos.Line, bytecode.SUBSCRIPT)
}

func (c *Compiler) compileAttributeAccessNode(node *ast.AttributeAccessNode) {
	c.compileNode(node.Receiver)

	name := value.ToSymbol(node.AttributeName)
	callInfo := value.NewCallSiteInfo(name, 0)
	if node.AttributeName == "call" {
		c.emitCall(callInfo, node.Span())
	} else {
		c.emitCallMethod(callInfo, node.Span())
	}
}

func (c *Compiler) compileConstructorCallNode(node *ast.ConstructorCallNode) {
	c.compileConstructorCall(
		func() {
			c.compileNode(node.Class)
		},
		node.PositionalArguments,
		node.Span(),
	)
}

func (c *Compiler) compileNewExpressionNode(node *ast.NewExpressionNode) {
	c.compileConstructorCall(
		func() {
			c.emit(node.Span().StartPos.Line, bytecode.SELF)
		},
		node.PositionalArguments,
		node.Span(),
	)
}

func (c *Compiler) compileGenericConstructorCallNode(node *ast.GenericConstructorCallNode) {
	c.compileConstructorCall(
		func() {
			c.compileNode(node.Class)
		},
		node.PositionalArguments,
		node.Span(),
	)
}

func (c *Compiler) compileConstructorCall(class func(), args []ast.ExpressionNode, span *position.Span) {
	class()
	for _, posArg := range args {
		c.compileNode(posArg)
	}

	c.emitInstantiate(len(args), span)
}

func (c *Compiler) compileMethodCallNode(node *ast.MethodCallNode) {
	c.compileMethodCall(
		node.Receiver,
		node.Op,
		node.MethodName,
		node.PositionalArguments,
		node.Span(),
	)
}
func (c *Compiler) compileGenericMethodCallNode(node *ast.GenericMethodCallNode) {
	c.compileMethodCall(
		node.Receiver,
		node.Op,
		node.MethodName,
		node.PositionalArguments,
		node.Span(),
	)
}

func (c *Compiler) compileMethodCall(receiver ast.ExpressionNode, op *token.Token, name string, args []ast.ExpressionNode, span *position.Span) {
	_, onSelf := receiver.(*ast.SelfLiteralNode)

	switch op.Type {
	case token.QUESTION_DOT:
		c.compileNode(receiver)
		nilJump := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF_NIL_NP)

		// if not nil
		// call the method
		c.compileInnerMethodCall(name, op, args, false, span)

		// if nil
		// leave nil on the stack
		c.patchJump(nilJump, span)
	case token.QUESTION_DOT_DOT:
		c.compileNode(receiver)
		c.emit(span.EndPos.Line, bytecode.DUP)
		nilJump := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF_NIL_NP)

		// if not nil
		// call the method
		c.compileInnerMethodCall(name, op, args, false, span)

		// if nil
		// leave nil on the stack
		c.patchJump(nilJump, span)
	case token.DOT_DOT:
		if !onSelf {
			c.compileNode(receiver)
		}
		c.emit(span.EndPos.Line, bytecode.DUP)
		c.compileInnerMethodCall(name, op, args, onSelf, span)
	case token.DOT:
		if !onSelf {
			c.compileNode(receiver)
		}
		c.compileInnerMethodCall(name, op, args, onSelf, span)
	default:
		panic(fmt.Sprintf("invalid method call operator: %#v", op))
	}
}

func (c *Compiler) compileInnerMethodCall(name string, op *token.Token, args []ast.ExpressionNode, onSelf bool, span *position.Span) {
	for _, posArg := range args {
		c.compileNode(posArg)
	}

	nameSym := value.ToSymbol(name)
	callInfo := value.NewCallSiteInfo(nameSym, len(args))
	if onSelf {
		c.emitCallSelf(callInfo, span)
	} else if name == "call" {
		c.emitCall(callInfo, span)
	} else {
		c.emitCallMethod(callInfo, span)
	}

	switch op.Type {
	case token.DOT_DOT, token.QUESTION_DOT_DOT:
		c.emit(span.EndPos.Line, bytecode.POP)
	case token.DOT, token.QUESTION_DOT:
	default:
		panic(fmt.Sprintf("invalid method call operator: %#v", op))
	}
}

func (c *Compiler) compileCallNode(node *ast.CallNode) {
	c.compileNode(node.Receiver)

	if node.NilSafe {
		nilJump := c.emitJump(node.Span().StartPos.Line, bytecode.JUMP_IF_NIL_NP)

		// if not nil
		// call the method
		c.compileInnerCall(node)

		// if nil
		// leave nil on the stack
		c.patchJump(nilJump, node.Span())
		return
	}

	c.compileInnerCall(node)
}

func (c *Compiler) compileInnerCall(node *ast.CallNode) {
	for _, posArg := range node.PositionalArguments {
		c.compileNode(posArg)
	}

	name := value.ToSymbol("call")
	callInfo := value.NewCallSiteInfo(name, len(node.PositionalArguments))
	c.emitCall(callInfo, node.Span())
}

func (c *Compiler) compileSingletonBlockExpressionNode(node *ast.SingletonBlockExpressionNode) bool {
	if len(node.Body) <= 0 {
		return false
	}

	span := node.Span()
	singletonType := c.typeOf(node).(*types.SingletonClass)
	singletonName := singletonType.Name()

	singletonCompiler := New(fmt.Sprintf("<singleton_class: %s>", singletonName), namespaceMode, c.newLocation(span), c.checker)
	singletonCompiler.Errors = c.Errors
	if !singletonCompiler.compileNamespace(node) {
		return false
	}

	c.emit(span.StartPos.Line, bytecode.SELF)
	c.emit(span.StartPos.Line, bytecode.GET_SINGLETON)

	result := singletonCompiler.Bytecode
	c.emitValue(value.Ref(result), span)
	c.emit(span.StartPos.Line, bytecode.INIT_NAMESPACE)
	return true
}

func (c *Compiler) compileClosureLiteralNode(node *ast.ClosureLiteralNode) {
	closureCompiler := New("<closure>", methodMode, c.newLocation(node.Span()), c.checker)
	closureCompiler.Parent = c
	closureCompiler.Errors = c.Errors
	closureCompiler.compileFunction(node.Span(), node.Parameters, node.Body)

	result := closureCompiler.Bytecode
	c.emitValue(value.Ref(result), node.Span())

	c.emit(node.Span().StartPos.Line, bytecode.CLOSURE)

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

func (c *Compiler) compileMixinDeclarationNode(node *ast.MixinDeclarationNode) bool {
	if len(node.Body) <= 0 {
		return false
	}

	span := node.Span()
	mixinType := c.typeOf(node).(*types.Mixin)
	mixinName := value.ToSymbol(mixinType.Name())

	mixinCompiler := New(fmt.Sprintf("<mixin: %s>", mixinType.Name()), namespaceMode, c.newLocation(node.Span()), c.checker)
	mixinCompiler.Errors = c.Errors
	if !mixinCompiler.compileNamespace(node) {
		return false
	}

	c.emitGetConst(mixinName, node.Constant.Span())
	result := mixinCompiler.Bytecode
	c.emitValue(value.Ref(result), span)
	c.emit(span.StartPos.Line, bytecode.INIT_NAMESPACE)
	return true
}

func (c *Compiler) compileModuleDeclarationNode(node *ast.ModuleDeclarationNode) bool {
	if len(node.Body) <= 0 {
		return false
	}

	span := node.Span()
	modType := c.typeOf(node).(*types.Module)
	modName := value.ToSymbol(modType.Name())

	modCompiler := New(fmt.Sprintf("<module: %s>", modType.Name()), namespaceMode, c.newLocation(node.Span()), c.checker)
	modCompiler.Errors = c.Errors
	if !modCompiler.compileNamespace(node) {
		return false
	}

	c.emitGetConst(modName, node.Constant.Span())
	result := modCompiler.Bytecode
	c.emitValue(value.Ref(result), span)
	c.emit(span.StartPos.Line, bytecode.INIT_NAMESPACE)
	return true
}

func (c *Compiler) compileInterfaceDeclarationNode(node *ast.InterfaceDeclarationNode) bool {
	if len(node.Body) <= 0 {
		return false
	}

	span := node.Span()
	ifaceType := c.typeOf(node).(*types.Interface)
	className := value.ToSymbol(ifaceType.Name())

	ifaceCompiler := New(fmt.Sprintf("<interface: %s>", ifaceType.Name()), namespaceMode, c.newLocation(node.Span()), c.checker)
	ifaceCompiler.Errors = c.Errors
	if !ifaceCompiler.compileNamespace(node) {
		return false
	}

	c.emitGetConst(className, node.Constant.Span())
	result := ifaceCompiler.Bytecode
	c.emitValue(value.Ref(result), span)
	c.emit(span.StartPos.Line, bytecode.INIT_NAMESPACE)
	return true
}

func (c *Compiler) compileClassDeclarationNode(node *ast.ClassDeclarationNode) bool {
	if len(node.Body) <= 0 {
		return false
	}

	span := node.Span()
	classType := c.typeOf(node).(*types.Class)
	className := value.ToSymbol(classType.Name())

	classCompiler := New(fmt.Sprintf("<class: %s>", classType.Name()), namespaceMode, c.newLocation(node.Span()), c.checker)
	classCompiler.Errors = c.Errors
	if !classCompiler.compileNamespace(node) {
		return false
	}

	c.emitGetConst(className, node.Constant.Span())
	result := classCompiler.Bytecode
	c.emitValue(value.Ref(result), span)
	c.emit(span.StartPos.Line, bytecode.INIT_NAMESPACE)
	return true
}

func (c *Compiler) compileValuePatternDeclarationNode(node *ast.ValuePatternDeclarationNode) {
	previousMode := c.mode
	c.mode = valuePatternDeclarationNode
	defer func() {
		c.mode = previousMode
	}()

	span := node.Span()
	c.compileNode(node.Initialiser)
	c.pattern(node.Pattern)

	jumpOverErrorOffset := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF)

	c.emitValue(
		value.Ref(value.NewError(
			value.PatternNotMatchedErrorClass,
			"assigned value does not match the pattern defined in value declaration",
		)),
		span,
	)
	c.emit(span.EndPos.Line, bytecode.THROW)

	c.patchJump(jumpOverErrorOffset, span)
}

func (c *Compiler) compilerVariablePatternDeclarationNode(node *ast.VariablePatternDeclarationNode) {
	span := node.Span()
	c.compileNode(node.Initialiser)
	c.pattern(node.Pattern)

	jumpOverErrorOffset := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF)

	c.emitValue(
		value.Ref(value.NewError(
			value.PatternNotMatchedErrorClass,
			"assigned value does not match the pattern defined in variable declaration",
		)),
		span,
	)
	c.emit(span.EndPos.Line, bytecode.THROW)

	c.patchJump(jumpOverErrorOffset, span)
}

func (c *Compiler) compileVariableDeclarationNode(node *ast.VariableDeclarationNode) {
	initialised := node.Initialiser != nil

	if initialised {
		c.compileNode(node.Initialiser)
	}
	local := c.defineLocal(node.Name, node.Span())
	if local == nil {
		return
	}
	if initialised {
		c.emitSetLocal(node.Span().StartPos.Line, local.index)
	} else {
		c.emit(node.Span().StartPos.Line, bytecode.NIL)
	}
}

// Compile each element of a collection of statements.
func (c *Compiler) compileStatements(collection []ast.StatementNode, span *position.Span) {
	if !c.compileStatementsOk(collection) {
		c.emit(span.EndPos.Line, bytecode.NIL)
	}
}

// Same as [compileStatements] but returns false when no instructions were emitted instead of
// emitting a `nil` value.
func (c *Compiler) compileStatementsOk(collection []ast.StatementNode) bool {
	var isPresent bool
	for _, s := range collection {
		exprIsPresent := c.compileNode(s)
		if exprIsPresent {
			isPresent = true
			c.emit(s.Span().EndPos.Line, bytecode.POP)
		}
	}

	if isPresent {
		// remove the last POP instruction
		c.removeOpcode()
		return true
	}

	return false
}

func (c *Compiler) removeOpcode() {
	c.lastOpCode = c.secondToLastOpCode
	c.secondToLastOpCode = bytecode.ZERO_VALUE
	c.Bytecode.RemoveByte()
}

func (c *Compiler) compileRangeLiteralNode(node *ast.RangeLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	span := node.Span()

	if node.Start == nil {
		c.compileNode(node.End)

		switch node.Op.Type {
		case token.CLOSED_RANGE_OP, token.LEFT_OPEN_RANGE_OP:
			c.emit(span.StartPos.Line, bytecode.NEW_RANGE, bytecode.BEGINLESS_CLOSED_RANGE_FLAG)
		case token.RIGHT_OPEN_RANGE_OP, token.OPEN_RANGE_OP:
			c.emit(span.StartPos.Line, bytecode.NEW_RANGE, bytecode.BEGINLESS_OPEN_RANGE_FLAG)
		default:
			panic(fmt.Sprintf("invalid range operator: %#v", node.Op))
		}

		return
	}
	if node.End == nil {
		c.compileNode(node.Start)

		switch node.Op.Type {
		case token.CLOSED_RANGE_OP, token.RIGHT_OPEN_RANGE_OP:
			c.emit(span.StartPos.Line, bytecode.NEW_RANGE, bytecode.ENDLESS_CLOSED_RANGE_FLAG)
		case token.LEFT_OPEN_RANGE_OP, token.OPEN_RANGE_OP:
			c.emit(span.StartPos.Line, bytecode.NEW_RANGE, bytecode.ENDLESS_OPEN_RANGE_FLAG)
		default:
			panic(fmt.Sprintf("invalid range operator: %#v", node.Op))
		}

		return
	}

	c.compileNode(node.Start)
	c.compileNode(node.End)
	switch node.Op.Type {
	case token.CLOSED_RANGE_OP:
		c.emit(span.StartPos.Line, bytecode.NEW_RANGE, bytecode.CLOSED_RANGE_FLAG)
	case token.OPEN_RANGE_OP:
		c.emit(span.StartPos.Line, bytecode.NEW_RANGE, bytecode.OPEN_RANGE_FLAG)
	case token.LEFT_OPEN_RANGE_OP:
		c.emit(span.StartPos.Line, bytecode.NEW_RANGE, bytecode.LEFT_OPEN_RANGE_FLAG)
	case token.RIGHT_OPEN_RANGE_OP:
		c.emit(span.StartPos.Line, bytecode.NEW_RANGE, bytecode.RIGHT_OPEN_RANGE_FLAG)
	default:
		panic(fmt.Sprintf("invalid range operator: %#v", node.Op))
	}
}

func (c *Compiler) compileHashSetLiteralNode(node *ast.HashSetLiteralNode) {
	span := node.Span()
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
		c.emit(span.StartPos.Line, bytecode.UNDEFINED)
	} else {
		c.compileNode(node.Capacity)
	}

	if baseSet.Length() == 0 && baseSet.Capacity() == 0 {
		c.emit(span.StartPos.Line, bytecode.UNDEFINED)
	} else {
		c.emitLoadValue(value.Ref(baseSet), span)
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
						c.newLocation(node.Capacity.Span()),
					)
					return
				}
				c.emitNewHashSet(i, span)
				firstModifierElementIndex = i
				break dynamicElementsLoop
			default:
				c.compileNode(elementNode)
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
						c.compileNode(e.Left)
						c.emit(e.Span().StartPos.Line, bytecode.APPEND)
					},
					func() {},
					e.Span(),
				)
			case *ast.ModifierIfElseNode:
				c.compileIfWithConditionExpression(
					bytecode.JUMP_UNLESS,
					e.Condition,
					func() {
						c.compileNode(e.ThenExpression)
						c.emit(e.Span().StartPos.Line, bytecode.APPEND)
					},
					func() {
						c.compileNode(e.ElseExpression)
						c.emit(e.Span().StartPos.Line, bytecode.APPEND)
					},
					e.Span(),
				)
			case *ast.ModifierForInNode:
				c.compileForIn(
					"",
					e.Pattern,
					e.InExpression,
					func() {
						c.compileNode(e.ThenExpression)
						c.emit(e.ThenExpression.Span().EndPos.Line, bytecode.APPEND)
					},
					e.Span(),
					true,
				)
			default:
				c.compileNode(elementNode)
				c.emit(e.Span().StartPos.Line, bytecode.APPEND)
			}
		}

		return
	}

	c.emitNewHashSet(len(dynamicElementNodes), span)
}

func (c *Compiler) compileHashMapLiteralNode(node *ast.HashMapLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	span := node.Span()
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
			key := value.ToSymbol(e.Key)
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
		c.emit(span.StartPos.Line, bytecode.UNDEFINED)
	} else {
		c.compileNode(node.Capacity)
	}

	if baseMap.Length() == 0 && baseMap.Capacity() == 0 {
		c.emit(span.StartPos.Line, bytecode.UNDEFINED)
	} else {
		c.emitLoadValue(value.Ref(baseMap), span)
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
						c.newLocation(node.Capacity.Span()),
					)
					return
				}
				c.emitNewHashMap(i, span)
				firstModifierElementIndex = i
				break dynamicElementsLoop
			case *ast.KeyValueExpressionNode:
				c.compileNode(element.Key)
				c.compileNode(element.Value)
			case *ast.SymbolKeyValueExpressionNode:
				c.emitValue(value.ToSymbol(element.Key).ToValue(), element.Span())
				c.compileNode(element.Value)
			case *ast.PublicIdentifierNode:
				c.emitValue(value.ToSymbol(element.Value).ToValue(), element.Span())
				c.compileLocalVariableAccess(element.Value, element.Span())
			case *ast.PrivateIdentifierNode:
				c.emitValue(value.ToSymbol(element.Value).ToValue(), element.Span())
				c.compileLocalVariableAccess(element.Value, element.Span())
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
				c.compileNode(e.Key)
				c.compileNode(e.Value)
				c.emit(e.Span().StartPos.Line, bytecode.MAP_SET)
			case *ast.SymbolKeyValueExpressionNode:
				c.emitValue(value.ToSymbol(e.Key).ToValue(), e.Span())
				c.compileNode(e.Value)
				c.emit(e.Span().StartPos.Line, bytecode.MAP_SET)
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
							c.compileNode(then.Key)
							c.compileNode(then.Value)
							c.emit(then.Span().StartPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(then.Key).ToValue(), then.Span())
							c.compileNode(then.Value)
							c.emit(then.Span().StartPos.Line, bytecode.MAP_SET)
						default:
							panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
						}
					},
					func() {},
					e.Span(),
				)
			case *ast.ModifierIfElseNode:
				c.compileIfWithConditionExpression(
					bytecode.JUMP_UNLESS,
					e.Condition,
					func() {
						switch then := e.ThenExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNode(then.Key)
							c.compileNode(then.Value)
							c.emit(then.Span().StartPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(then.Key).ToValue(), then.Span())
							c.compileNode(then.Value)
							c.emit(then.Span().StartPos.Line, bytecode.MAP_SET)
						default:
							panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
						}
					},
					func() {
						switch els := e.ElseExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNode(els.Key)
							c.compileNode(els.Value)
							c.emit(els.Span().EndPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(els.Key).ToValue(), els.Span())
							c.compileNode(els.Value)
							c.emit(els.Span().EndPos.Line, bytecode.MAP_SET)
						default:
							panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
						}
					},
					e.Span(),
				)
			case *ast.ModifierForInNode:
				c.compileForIn(
					"",
					e.Pattern,
					e.InExpression,
					func() {
						switch then := e.ThenExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNode(then.Key)
							c.compileNode(then.Value)
							c.emit(then.Span().EndPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(then.Key).ToValue(), then.Span())
							c.compileNode(then.Value)
							c.emit(then.Span().EndPos.Line, bytecode.MAP_SET)
						default:
							panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
						}
					},
					e.Span(),
					true,
				)
			default:
				panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
			}
		}

		return
	}

	c.emitNewHashMap(len(dynamicElementNodes), span)
}

func (c *Compiler) compileHashRecordLiteralNode(node *ast.HashRecordLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	span := node.Span()
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
			key := value.ToSymbol(e.Key)
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
		c.emit(span.StartPos.Line, bytecode.UNDEFINED)
	} else {
		c.emitLoadValue(value.Ref(baseRecord), span)
	}

	firstModifierElementIndex := -1
	var dynamicElementNodes []ast.ExpressionNode

	if firstDynamicIndex != -1 {
		dynamicElementNodes = node.Elements[firstDynamicIndex:]
	dynamicElementsLoop:
		for i, elementNode := range dynamicElementNodes {
			switch element := elementNode.(type) {
			case *ast.ModifierNode, *ast.ModifierForInNode, *ast.ModifierIfElseNode:
				c.emitNewHashRecord(i, span)
				firstModifierElementIndex = i
				break dynamicElementsLoop
			case *ast.KeyValueExpressionNode:
				c.compileNode(element.Key)
				c.compileNode(element.Value)
			case *ast.SymbolKeyValueExpressionNode:
				c.emitValue(value.ToSymbol(element.Key).ToValue(), element.Span())
				c.compileNode(element.Value)
			case *ast.PublicIdentifierNode:
				c.emitValue(value.ToSymbol(element.Value).ToValue(), element.Span())
				c.compileLocalVariableAccess(element.Value, element.Span())
			case *ast.PrivateIdentifierNode:
				c.emitValue(value.ToSymbol(element.Value).ToValue(), element.Span())
				c.compileLocalVariableAccess(element.Value, element.Span())
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
				c.compileNode(e.Key)
				c.compileNode(e.Value)
				c.emit(e.Span().StartPos.Line, bytecode.MAP_SET)
			case *ast.SymbolKeyValueExpressionNode:
				c.emitValue(value.ToSymbol(e.Key).ToValue(), e.Span())
				c.compileNode(e.Value)
				c.emit(e.Span().StartPos.Line, bytecode.MAP_SET)
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
							c.compileNode(then.Key)
							c.compileNode(then.Value)
							c.emit(then.Span().StartPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(then.Key).ToValue(), then.Span())
							c.compileNode(then.Value)
							c.emit(then.Span().StartPos.Line, bytecode.MAP_SET)
						default:
							panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
						}
					},
					func() {},
					e.Span(),
				)
			case *ast.ModifierIfElseNode:
				c.compileIfWithConditionExpression(
					bytecode.JUMP_UNLESS,
					e.Condition,
					func() {
						switch then := e.ThenExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNode(then.Key)
							c.compileNode(then.Value)
							c.emit(then.Span().StartPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(then.Key).ToValue(), then.Span())
							c.compileNode(then.Value)
							c.emit(then.Span().StartPos.Line, bytecode.MAP_SET)
						default:
							panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
						}
					},
					func() {
						switch els := e.ElseExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNode(els.Key)
							c.compileNode(els.Value)
							c.emit(els.Span().EndPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(els.Key).ToValue(), els.Span())
							c.compileNode(els.Value)
							c.emit(els.Span().EndPos.Line, bytecode.MAP_SET)
						default:
							panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
						}
					},
					e.Span(),
				)
			case *ast.ModifierForInNode:
				c.compileForIn(
					"",
					e.Pattern,
					e.InExpression,
					func() {
						switch then := e.ThenExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNode(then.Key)
							c.compileNode(then.Value)
							c.emit(then.Span().EndPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(then.Key).ToValue(), then.Span())
							c.compileNode(then.Value)
							c.emit(then.Span().EndPos.Line, bytecode.MAP_SET)
						default:
							panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
						}
					},
					e.Span(),
					true,
				)
			default:
				panic(fmt.Sprintf("invalid hash map element: %#v", elementNode))
			}
		}

		return
	}

	c.emitNewHashRecord(len(dynamicElementNodes), span)
}

func (c *Compiler) compileArrayListLiteralNode(node *ast.ArrayListLiteralNode) {
	span := node.Span()
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
		c.emit(span.StartPos.Line, bytecode.UNDEFINED)
	} else {
		c.compileNode(node.Capacity)
	}

	if len(baseList) == 0 && (keyValueCount == 0 || cap(baseList) == 0) {
		c.emit(span.StartPos.Line, bytecode.UNDEFINED)
	} else {
		c.emitLoadValue(value.Ref(&baseList), span)
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
						c.newLocation(node.Capacity.Span()),
					)
					return
				}
				c.emitNewArrayList(i, span)
				firstModifierElementIndex = i
				break dynamicElementsLoop
			case *ast.KeyValueExpressionNode:
				c.emitNewArrayList(i, span)
				firstModifierElementIndex = i
				break dynamicElementsLoop
			default:
				c.compileNode(elementNode)
			}
		}
	}

	if firstModifierElementIndex != -1 {
		modifierElementNodes := dynamicElementNodes[firstModifierElementIndex:]
		for _, elementNode := range modifierElementNodes {
			switch e := elementNode.(type) {
			case *ast.KeyValueExpressionNode:
				c.compileNode(e.Key)
				c.compileNode(e.Value)
				c.emit(e.Span().StartPos.Line, bytecode.APPEND_AT)
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
							c.compileNode(then.Key)
							c.compileNode(then.Value)
							c.emit(then.Span().StartPos.Line, bytecode.APPEND_AT)
						default:
							c.compileNode(e.Left)
							c.emit(e.Span().StartPos.Line, bytecode.APPEND)
						}
					},
					func() {},
					e.Span(),
				)
			case *ast.ModifierIfElseNode:
				c.compileIfWithConditionExpression(
					bytecode.JUMP_UNLESS,
					e.Condition,
					func() {
						switch then := e.ThenExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNode(then.Key)
							c.compileNode(then.Value)
							c.emit(then.Span().StartPos.Line, bytecode.APPEND_AT)
						default:
							c.compileNode(e.ThenExpression)
							c.emit(e.Span().StartPos.Line, bytecode.APPEND)
						}
					},
					func() {
						switch els := e.ElseExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNode(els.Key)
							c.compileNode(els.Value)
							c.emit(els.Span().StartPos.Line, bytecode.APPEND_AT)
						default:
							c.compileNode(e.ElseExpression)
							c.emit(e.Span().StartPos.Line, bytecode.APPEND)
						}
					},
					e.Span(),
				)
			case *ast.ModifierForInNode:
				c.compileForIn(
					"",
					e.Pattern,
					e.InExpression,
					func() {
						switch then := e.ThenExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNode(then.Key)
							c.compileNode(then.Value)
							c.emit(then.Span().EndPos.Line, bytecode.APPEND_AT)
						default:
							c.compileNode(e.ThenExpression)
							c.emit(then.Span().EndPos.Line, bytecode.APPEND)
						}
					},
					e.Span(),
					true,
				)
			default:
				c.compileNode(elementNode)
				c.emit(e.Span().StartPos.Line, bytecode.APPEND)
			}
		}

		return
	}

	c.emitNewArrayList(len(dynamicElementNodes), span)
}

func (c *Compiler) compileArrayTupleLiteralNode(node *ast.ArrayTupleLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	span := node.Span()

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
		c.emit(span.StartPos.Line, bytecode.UNDEFINED)
	} else {
		c.emitLoadValue(value.Ref(&baseArrayTuple), span)
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
					c.emit(e.Span().StartPos.Line, bytecode.COPY)
				} else {
					c.emitNewArrayTuple(i, span)
				}
				firstModifierElementIndex = i
				break dynamicElementsLoop
			default:
				c.compileNode(elementNode)
			}
		}
	}

	if firstModifierElementIndex != -1 {
		modifierElementNodes := dynamicElementNodes[firstModifierElementIndex:]
		for _, elementNode := range modifierElementNodes {
			switch e := elementNode.(type) {
			case *ast.KeyValueExpressionNode:
				c.compileNode(e.Key)
				c.compileNode(e.Value)
				c.emit(e.Span().StartPos.Line, bytecode.APPEND_AT)
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
							c.compileNode(then.Key)
							c.compileNode(then.Value)
							c.emit(then.Span().StartPos.Line, bytecode.APPEND_AT)
						default:
							c.compileNode(e.Left)
							c.emit(e.Span().StartPos.Line, bytecode.APPEND)
						}
					},
					func() {},
					e.Span(),
				)
			case *ast.ModifierIfElseNode:
				c.compileIfWithConditionExpression(
					bytecode.JUMP_UNLESS,
					e.Condition,
					func() {
						switch then := e.ThenExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNode(then.Key)
							c.compileNode(then.Value)
							c.emit(then.Span().EndPos.Line, bytecode.APPEND_AT)
						default:
							c.compileNode(e.ThenExpression)
							c.emit(e.Span().EndPos.Line, bytecode.APPEND)
						}
					},
					func() {
						switch els := e.ElseExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNode(els.Key)
							c.compileNode(els.Value)
							c.emit(els.Span().EndPos.Line, bytecode.APPEND_AT)
						default:
							c.compileNode(e.ElseExpression)
							c.emit(e.Span().EndPos.Line, bytecode.APPEND)
						}
					},
					e.Span(),
				)
			case *ast.ModifierForInNode:
				c.compileForIn(
					"",
					e.Pattern,
					e.InExpression,
					func() {
						switch then := e.ThenExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNode(then.Key)
							c.compileNode(then.Value)
							c.emit(then.Span().EndPos.Line, bytecode.APPEND_AT)
						default:
							c.compileNode(e.ThenExpression)
							c.emit(then.Span().EndPos.Line, bytecode.APPEND)
						}
					},
					e.Span(),
					true,
				)
			default:
				c.compileNode(elementNode)
				c.emit(e.Span().StartPos.Line, bytecode.APPEND)
			}
		}

		return
	}

	c.emitNewArrayTuple(len(dynamicElementNodes), span)
}

func (c *Compiler) compileWordArrayTupleLiteralNode(node *ast.WordArrayTupleLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	c.Errors.AddFailure("invalid word arrayTuple literal", c.newLocation(node.Span()))
}

func (c *Compiler) compileBinArrayTupleLiteralNode(node *ast.BinArrayTupleLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	c.Errors.AddFailure("invalid binary arrayTuple literal", c.newLocation(node.Span()))
}

func (c *Compiler) compileSymbolArrayTupleLiteralNode(node *ast.SymbolArrayTupleLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	c.Errors.AddFailure("invalid symbol arrayTuple literal", c.newLocation(node.Span()))
}

func (c *Compiler) compileHexArrayTupleLiteralNode(node *ast.HexArrayTupleLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	c.Errors.AddFailure("invalid hex arrayTuple literal", c.newLocation(node.Span()))
}

func (c *Compiler) compileWordArrayListLiteralNode(node *ast.WordArrayListLiteralNode) {
	list := resolve(node)
	span := node.Span()
	if list.IsUndefined() {
		c.Errors.AddFailure("invalid word arrayList literal", c.newLocation(span))
		return
	}

	if node.Capacity == nil {
		c.emitLoadValue(list, span)
		c.emit(span.EndPos.Line, bytecode.COPY)
	} else {
		c.compileNode(node.Capacity)
		c.emitLoadValue(list, span)
		c.emitNewArrayList(0, span)
	}
}

func (c *Compiler) compileBinArrayListLiteralNode(node *ast.BinArrayListLiteralNode) {
	list := resolve(node)
	span := node.Span()
	if list.IsUndefined() {
		c.Errors.AddFailure("invalid bin arrayList literal", c.newLocation(span))
		return
	}

	if node.Capacity == nil {
		c.emitLoadValue(list, span)
		c.emit(span.EndPos.Line, bytecode.COPY)
	} else {
		c.compileNode(node.Capacity)
		c.emitLoadValue(list, span)
		c.emitNewArrayList(0, span)
	}
}

func (c *Compiler) compileSymbolArrayListLiteralNode(node *ast.SymbolArrayListLiteralNode) {
	list := resolve(node)
	span := node.Span()
	if list.IsUndefined() {
		c.Errors.AddFailure("invalid symbol arrayList literal", c.newLocation(span))
		return
	}

	if node.Capacity == nil {
		c.emitLoadValue(list, span)
		c.emit(span.EndPos.Line, bytecode.COPY)
	} else {
		c.compileNode(node.Capacity)
		c.emitLoadValue(list, span)
		c.emitNewArrayList(0, span)
	}
}

func (c *Compiler) compileHexArrayListLiteralNode(node *ast.HexArrayListLiteralNode) {
	list := resolve(node)
	span := node.Span()
	if list.IsUndefined() {
		c.Errors.AddFailure("invalid hex arrayList literal", c.newLocation(span))
		return
	}

	if node.Capacity == nil {
		c.emitLoadValue(list, span)
		c.emit(span.EndPos.Line, bytecode.COPY)
	} else {
		c.compileNode(node.Capacity)
		c.emitLoadValue(list, span)
		c.emitNewArrayList(0, span)
	}
}

func (c *Compiler) compileWordHashSetLiteralNode(node *ast.WordHashSetLiteralNode) {
	list := resolve(node)
	span := node.Span()
	if list.IsUndefined() {
		c.Errors.AddFailure("invalid word hashSet literal", c.newLocation(span))
		return
	}

	if node.Capacity == nil {
		c.emitLoadValue(list, span)
		c.emit(span.EndPos.Line, bytecode.COPY)
	} else {
		c.compileNode(node.Capacity)
		c.emitLoadValue(list, span)
		c.emitNewHashSet(0, span)
	}
}

func (c *Compiler) compileBinHashSetLiteralNode(node *ast.BinHashSetLiteralNode) {
	list := resolve(node)
	span := node.Span()
	if list.IsUndefined() {
		c.Errors.AddFailure("invalid bin hashSet literal", c.newLocation(span))
		return
	}

	if node.Capacity == nil {
		c.emitLoadValue(list, span)
		c.emit(span.EndPos.Line, bytecode.COPY)
	} else {
		c.compileNode(node.Capacity)
		c.emitLoadValue(list, span)
		c.emitNewHashSet(0, span)
	}
}

func (c *Compiler) compileSymbolHashSetLiteralNode(node *ast.SymbolHashSetLiteralNode) {
	list := resolve(node)
	span := node.Span()
	if list.IsUndefined() {
		c.Errors.AddFailure("invalid symbol hashSet literal", c.newLocation(span))
		return
	}

	if node.Capacity == nil {
		c.emitLoadValue(list, span)
		c.emit(span.EndPos.Line, bytecode.COPY)
	} else {
		c.compileNode(node.Capacity)
		c.emitLoadValue(list, span)
		c.emitNewHashSet(0, span)
	}
}

func (c *Compiler) compileHexHashSetLiteralNode(node *ast.HexHashSetLiteralNode) {
	list := resolve(node)
	span := node.Span()
	if list.IsUndefined() {
		c.Errors.AddFailure("invalid hex hashSet literal", c.newLocation(span))
		return
	}

	if node.Capacity == nil {
		c.emitLoadValue(list, span)
		c.emit(span.EndPos.Line, bytecode.COPY)
	} else {
		c.compileNode(node.Capacity)
		c.emitLoadValue(list, span)
		c.emitNewHashSet(0, span)
	}
}

func (c *Compiler) emitNewHashSet(size int, span *position.Span) {
	c.emitNewCollection(bytecode.NEW_HASH_SET8, bytecode.NEW_HASH_SET16, size, span)
}

func (c *Compiler) emitNewArrayTuple(size int, span *position.Span) {
	c.emitNewCollection(bytecode.NEW_ARRAY_TUPLE8, bytecode.NEW_ARRAY_TUPLE16, size, span)
}

func (c *Compiler) emitNewArrayList(size int, span *position.Span) {
	c.emitNewCollection(bytecode.NEW_ARRAY_LIST8, bytecode.NEW_ARRAY_LIST16, size, span)
}

func (c *Compiler) emitNewHashMap(size int, span *position.Span) {
	c.emitNewCollection(bytecode.NEW_HASH_MAP8, bytecode.NEW_HASH_MAP16, size, span)
}

func (c *Compiler) emitNewHashRecord(size int, span *position.Span) {
	c.emitNewCollection(bytecode.NEW_HASH_RECORD8, bytecode.NEW_HASH_RECORD16, size, span)
}

func (c *Compiler) emitNewRegex(flags bitfield.BitField8, size int, span *position.Span) {
	if size <= math.MaxUint8 {
		c.emit(span.EndPos.Line, bytecode.NEW_REGEX8, flags.Byte(), byte(size))
		return
	}

	if size <= math.MaxUint16 {
		c.emit(span.EndPos.Line, bytecode.NEW_REGEX16)
		c.emitByte(flags.Byte())
		c.emitUint16(uint16(size))
		return
	}

	c.Errors.AddFailure(
		fmt.Sprintf("max number of regex literal elements reached: %d", math.MaxUint16),
		c.newLocation(span),
	)
}

func (c *Compiler) emitNewCollection(opcode8, opcode16 bytecode.OpCode, size int, span *position.Span) {
	if size <= math.MaxUint8 {
		c.emit(span.EndPos.Line, opcode8, byte(size))
		return
	}

	if size <= math.MaxUint16 {
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, uint16(size))
		c.emit(span.EndPos.Line, opcode16, bytes...)
		return
	}

	c.Errors.AddFailure(
		fmt.Sprintf("max number of collection literal elements reached: %d", math.MaxUint16),
		c.newLocation(span),
	)
}

func (c *Compiler) compileUninterpolatedRegexLiteralNode(node *ast.UninterpolatedRegexLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	re, err := value.CompileRegex(node.Content, node.Flags)
	if errList, ok := err.(error.ErrorList); ok {
		regexStartPos := node.Span().StartPos
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
			err.Location.Filename = c.Bytecode.Location.Filename

			c.Errors.Append(err)
		}
		return
	}

	if err != nil {
		c.Errors.AddFailure(err.Error(), c.newLocation(node.Span()))
		return
	}

	c.emitValue(value.Ref(re), node.Span())
}

func (c *Compiler) compileInterpolatedRegexLiteralNode(node *ast.InterpolatedRegexLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	for _, elementNode := range node.Content {
		switch element := elementNode.(type) {
		case *ast.RegexLiteralContentSectionNode:
			c.emitValue(value.Ref(value.String(element.Value)), element.Span())
		case *ast.RegexInterpolationNode:
			c.compileNode(element.Expression)
		}
	}
	c.emitNewRegex(node.Flags, len(node.Content), node.Span())
}

var inspectSymbol = value.ToSymbol("inspect")

func (c *Compiler) compileInterpolatedStringLiteralNode(node *ast.InterpolatedStringLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	for _, elementNode := range node.Content {
		switch element := elementNode.(type) {
		case *ast.StringLiteralContentSectionNode:
			c.emitValue(value.Ref(value.String(element.Value)), element.Span())
		case *ast.StringInterpolationNode:
			c.compileNode(element.Expression)
		case *ast.StringInspectInterpolationNode:
			c.compileNode(element.Expression)
			callInfo := value.NewCallSiteInfo(inspectSymbol, 0)
			c.emitCallMethod(callInfo, element.Span())
		}
	}

	c.emitNewCollection(bytecode.NEW_STRING8, bytecode.NEW_STRING16, len(node.Content), node.Span())
}

func (c *Compiler) compileInterpolatedSymbolLiteralNode(node *ast.InterpolatedSymbolLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	for _, elementNode := range node.Content.Content {
		switch element := elementNode.(type) {
		case *ast.StringLiteralContentSectionNode:
			c.emitValue(value.Ref(value.String(element.Value)), element.Span())
		case *ast.StringInterpolationNode:
			c.compileNode(element.Expression)
		}
	}

	c.emitNewCollection(bytecode.NEW_SYMBOL8, bytecode.NEW_SYMBOL16, len(node.Content.Content), node.Span())
}

func (c *Compiler) compileIntLiteralNode(node *ast.IntLiteralNode) {
	i, err := value.ParseBigInt(node.Value, 0)
	if !err.IsUndefined() {
		c.Errors.AddFailure(err.Error(), c.newLocation(node.Span()))
		return
	}
	if i.IsSmallInt() {
		c.emitValue(i.ToSmallInt().ToValue(), node.Span())
		return
	}
	c.emitValue(value.Ref(i), node.Span())
}

// Compiles boolean binary operators
func (c *Compiler) compileLogicalExpressionNode(node *ast.LogicalExpressionNode) {
	if c.resolveAndEmit(node) {
		return
	}
	switch node.Op.Type {
	case token.AND_AND:
		c.logicalAnd(node)
	case token.OR_OR:
		c.logicalOr(node)
	case token.QUESTION_QUESTION:
		c.nilCoalescing(node)
	default:
		c.Errors.AddFailure(fmt.Sprintf("unknown logical operator: %s", node.Op.String()), c.newLocation(node.Span()))
	}
}

// Compiles the `??` operator
func (c *Compiler) nilCoalescing(node *ast.LogicalExpressionNode) {
	c.compileNode(node.Left)
	nilJump := c.emitJump(node.Span().StartPos.Line, bytecode.JUMP_IF_NIL_NP)
	nonNilJump := c.emitJump(node.Span().StartPos.Line, bytecode.JUMP)

	// if nil
	c.patchJump(nilJump, node.Span())
	c.emit(node.Span().StartPos.Line, bytecode.POP)
	c.compileNode(node.Right)

	// if not nil
	c.patchJump(nonNilJump, node.Span())
}

// Compiles the `||` operator
func (c *Compiler) logicalOr(node *ast.LogicalExpressionNode) {
	c.compileNode(node.Left)
	jump := c.emitJump(node.Span().StartPos.Line, bytecode.JUMP_IF_NP)

	// if falsy
	c.emit(node.Span().StartPos.Line, bytecode.POP)
	c.compileNode(node.Right)

	// if truthy
	c.patchJump(jump, node.Span())
}

// Compiles the `&&` operator
func (c *Compiler) logicalAnd(node *ast.LogicalExpressionNode) {
	c.compileNode(node.Left)
	jump := c.emitJump(node.Span().StartPos.Line, bytecode.JUMP_UNLESS_NP)

	// if truthy
	c.emit(node.Span().StartPos.Line, bytecode.POP)
	c.compileNode(node.Right)

	// if falsy
	c.patchJump(jump, node.Span())
}

func (c *Compiler) compileBinaryExpressionNode(node *ast.BinaryExpressionNode) {
	if c.resolveAndEmit(node) {
		return
	}
	c.compileNode(node.Left)
	c.compileNode(node.Right)
	c.emitBinaryOperation(c.typeOf(node.Left), node.Op, node.Span())
}

func (c *Compiler) emitBinaryOperation(typ types.Type, opToken *token.Token, span *position.Span) {
	line := span.StartPos.Line
	switch opToken.Type {
	case token.PLUS:
		// c.emitCallMethod(value.NewCallSiteInfo(symbol.OpAdd, 1), span)
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.ADD_INT)
			return
		}
		c.emit(line, bytecode.ADD)
	case token.MINUS:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.SUBTRACT_INT)
			return
		}
		c.emit(line, bytecode.SUBTRACT)
	case token.STAR:
		c.emit(line, bytecode.MULTIPLY)
	case token.SLASH:
		c.emit(line, bytecode.DIVIDE)
	case token.STAR_STAR:
		c.emit(line, bytecode.EXPONENTIATE)
	case token.LBITSHIFT:
		c.emit(line, bytecode.LBITSHIFT)
	case token.LTRIPLE_BITSHIFT:
		c.emit(line, bytecode.LOGIC_LBITSHIFT)
	case token.RBITSHIFT:
		c.emit(line, bytecode.RBITSHIFT)
	case token.RTRIPLE_BITSHIFT:
		c.emit(line, bytecode.LOGIC_RBITSHIFT)
	case token.AND:
		c.emit(line, bytecode.BITWISE_AND)
	case token.AND_TILDE:
		c.emit(line, bytecode.BITWISE_AND_NOT)
	case token.OR:
		c.emit(line, bytecode.BITWISE_OR)
	case token.XOR:
		c.emit(line, bytecode.BITWISE_XOR)
	case token.PERCENT:
		c.emit(line, bytecode.MODULO)
	case token.LAX_EQUAL:
		c.emit(line, bytecode.LAX_EQUAL)
	case token.LAX_NOT_EQUAL:
		c.emit(line, bytecode.LAX_NOT_EQUAL)
	case token.EQUAL_EQUAL:
		c.emit(line, bytecode.EQUAL)
	case token.NOT_EQUAL:
		c.emit(line, bytecode.NOT_EQUAL)
	case token.STRICT_EQUAL:
		c.emit(line, bytecode.STRICT_EQUAL)
	case token.STRICT_NOT_EQUAL:
		c.emit(line, bytecode.STRICT_NOT_EQUAL)
	case token.GREATER:
		c.emit(line, bytecode.GREATER)
	case token.GREATER_EQUAL:
		c.emit(line, bytecode.GREATER_EQUAL)
	case token.LESS:
		c.emit(line, bytecode.LESS)
	case token.LESS_EQUAL:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			c.emit(line, bytecode.LESS_EQUAL_INT)
			return
		}
		c.emit(line, bytecode.LESS_EQUAL)
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
		c.Errors.AddFailure(fmt.Sprintf("unknown binary operator: %s", opToken.String()), c.newLocation(span))
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

	c.emitValue(result, node.Span())
	return true
}

func (c *Compiler) resolveAndEmitList(node *ast.ArrayListLiteralNode) bool {
	result := resolveArrayListLiteral(node)
	if result.IsUndefined() {
		return false
	}

	c.emitValue(result, node.Span())
	return true
}

func (c *Compiler) emitValue(val value.Value, span *position.Span) {
	if val.IsReference() {
		switch v := val.AsReference().(type) {
		case *value.ArrayList:
			c.emitArrayList(v, span)
		case *value.ArrayTuple:
			c.emitArrayTuple(v, span)
		case *value.HashSet:
			c.emitHashSet(v, span)
		case *value.HashMap:
			c.emitHashMap(v, span)
		case *value.HashRecord:
			c.emitHashRecord(v, span)
		default:
			c.emitLoadValue(val, span)
		}
		return
	}

	switch val.ValueFlag() {
	case value.TRUE_FLAG:
		c.emit(span.StartPos.Line, bytecode.TRUE)
	case value.FALSE_FLAG:
		c.emit(span.StartPos.Line, bytecode.FALSE)
	case value.NIL_FLAG:
		c.emit(span.StartPos.Line, bytecode.NIL)
	case value.SMALL_INT_FLAG:
		c.emitSmallInt(val.AsSmallInt(), span)
	case value.INT64_FLAG:
		emitSignedInt(c, val, val.AsInt64(), bytecode.INT64_0, bytecode.INT64_1, bytecode.LOAD_INT64_8, span)
	case value.UINT64_FLAG:
		emitUnsignedInt(c, val, val.AsUInt64(), bytecode.UINT64_0, bytecode.UINT64_1, bytecode.LOAD_UINT64_8, span)
	case value.INT32_FLAG:
		emitSignedInt(c, val, val.AsInt32(), bytecode.INT32_0, bytecode.INT32_1, bytecode.LOAD_INT32_8, span)
	case value.UINT32_FLAG:
		emitUnsignedInt(c, val, val.AsUInt32(), bytecode.UINT32_0, bytecode.UINT32_1, bytecode.LOAD_UINT32_8, span)
	case value.INT16_FLAG:
		emitSignedInt(c, val, val.AsInt16(), bytecode.INT16_0, bytecode.INT16_1, bytecode.LOAD_INT16_8, span)
	case value.UINT16_FLAG:
		emitUnsignedInt(c, val, val.AsUInt16(), bytecode.UINT16_0, bytecode.UINT16_1, bytecode.LOAD_UINT16_8, span)
	case value.INT8_FLAG:
		emitSignedInt(c, val, val.AsInt8(), bytecode.INT8_0, bytecode.INT8_1, bytecode.LOAD_INT8, span)
	case value.UINT8_FLAG:
		emitUnsignedInt(c, val, val.AsUInt8(), bytecode.UINT8_0, bytecode.UINT8_1, bytecode.LOAD_UINT8, span)
	case value.CHAR_FLAG:
		c.emitChar(val.AsChar(), span)
	case value.FLOAT_FLAG:
		c.emitFloat(val.AsFloat(), span)
	default:
		c.emitLoadValue(val, span)
	}
}

func emitSignedInt[I value.SingedInt](c *Compiler, iVal value.Value, i I, opcode0, opcode1, opcodeLoad bytecode.OpCode, span *position.Span) {
	line := span.StartPos.Line
	switch i {
	case 0:
		c.emit(line, opcode0)
		return
	case 1:
		c.emit(line, opcode1)
		return
	}

	if i >= math.MinInt8 && i <= math.MaxInt8 {
		c.emit(line, opcodeLoad, byte(i))
		return
	}

	c.emitLoadValue(iVal, span)
}

func emitUnsignedInt[I value.UnsignedInt](c *Compiler, iVal value.Value, i I, opcode0, opcode1, opcodeLoad bytecode.OpCode, span *position.Span) {
	line := span.StartPos.Line
	switch i {
	case 0:
		c.emit(line, opcode0)
		return
	case 1:
		c.emit(line, opcode1)
		return
	}

	if i <= math.MaxUint8 {
		c.emit(line, opcodeLoad, byte(i))
		return
	}

	c.emitLoadValue(iVal, span)
}

func (c *Compiler) emitSmallInt(i value.SmallInt, span *position.Span) {
	line := span.StartPos.Line
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

	c.emitLoadValue(i.ToValue(), span)
}

func (c *Compiler) emitChar(char value.Char, span *position.Span) {
	line := span.StartPos.Line

	if char >= math.MinInt8 && char <= math.MaxInt8 {
		c.emit(line, bytecode.LOAD_CHAR_8, byte(char))
		return
	}

	c.emitLoadValue(char.ToValue(), span)
}

func (c *Compiler) emitFloat(f value.Float, span *position.Span) {
	line := span.StartPos.Line
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

	c.emitLoadValue(f.ToValue(), span)
}

func (c *Compiler) emitHashSet(set *value.HashSet, span *position.Span) {
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
		c.emitLoadValue(value.Ref(baseSet), span)
		c.emit(span.EndPos.Line, bytecode.COPY)
		return
	}

	// capacity
	c.emit(span.StartPos.Line, bytecode.UNDEFINED)
	c.emitLoadValue(value.Ref(baseSet), span)

	for _, element := range mutableElements {
		c.emitValue(element, span)
	}

	c.emitNewHashMap(len(mutableElements), span)
}
func (c *Compiler) emitHashMap(hmap *value.HashMap, span *position.Span) {
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
		c.emitLoadValue(value.Ref(baseMap), span)
		c.emit(span.EndPos.Line, bytecode.COPY)
		return
	}

	// capacity
	c.emit(span.StartPos.Line, bytecode.UNDEFINED)
	c.emitLoadValue(value.Ref(baseMap), span)

	for _, element := range mutablePairs {
		c.emitValue(element.Key, span)
		c.emitValue(element.Value, span)
	}

	c.emitNewHashMap(len(mutablePairs), span)
}

func (c *Compiler) emitHashRecord(hrec *value.HashRecord, span *position.Span) {
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
		c.emitLoadValue(value.Ref(baseRecord), span)
		return
	}

	c.emitLoadValue(value.Ref(baseRecord), span)

	for _, element := range mutablePairs {
		c.emitValue(element.Key, span)
		c.emitValue(element.Value, span)
	}

	c.emitNewHashRecord(len(mutablePairs), span)
}

func (c *Compiler) emitArrayList(list *value.ArrayList, span *position.Span) {
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
		c.emitLoadValue(value.Ref(list), span)
		c.emit(span.EndPos.Line, bytecode.COPY)
		return
	}

	// capacity
	c.emit(span.StartPos.Line, bytecode.UNDEFINED)

	baseList := l[:firstMutableElementIndex]
	c.emitLoadValue(value.Ref(&baseList), span)

	rest := l[firstMutableElementIndex:]
	for _, element := range rest {
		c.emitValue(element, span)
	}

	c.emitNewArrayList(len(rest), span)
}

func (c *Compiler) emitArrayTuple(tuple *value.ArrayTuple, span *position.Span) {
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
		c.emitLoadValue(value.Ref(tuple), span)
		return
	}

	baseTuple := t[:firstMutableElementIndex]
	c.emitLoadValue(value.Ref(&baseTuple), span)

	rest := t[firstMutableElementIndex:]
	for _, element := range rest {
		c.emitValue(element, span)
	}

	c.emitNewArrayList(len(rest), span)
}

func (c *Compiler) compileUnaryExpressionNode(node *ast.UnaryExpressionNode) {
	if c.resolveAndEmit(node) {
		return
	}
	c.compileNode(node.Right)
	switch node.Op.Type {
	case token.PLUS:
		c.emit(node.Span().StartPos.Line, bytecode.UNARY_PLUS)
	case token.MINUS:
		c.emit(node.Span().StartPos.Line, bytecode.NEGATE)
	case token.BANG:
		// logical not
		c.emit(node.Span().StartPos.Line, bytecode.NOT)
	case token.TILDE:
		// binary negation
		c.emit(node.Span().StartPos.Line, bytecode.BITWISE_NOT)
	case token.AND:
		// get singleton class
		c.emit(node.Span().StartPos.Line, bytecode.GET_SINGLETON)
	default:
		c.Errors.AddFailure(fmt.Sprintf("unknown unary operator: %s", node.Op.String()), c.newLocation(node.Span()))
	}
}

// Emit an instruction that jumps forward with a placeholder offset.
// Returns the offset of placeholder value that has to be patched.
func (c *Compiler) emitJump(line int, op bytecode.OpCode) int {
	c.emit(line, op, 0xff, 0xff)
	return c.nextInstructionOffset() - 2
}

// Emit an instruction that returns a value.
// Provide `nil` as the value when the returned value is already
// on the stack.
func (c *Compiler) emitReturn(span *position.Span, value ast.Node) {
	switch c.lastOpCode {
	case bytecode.RETURN, bytecode.RETURN_FIRST_ARG,
		bytecode.RETURN_SELF, bytecode.RETURN_FINALLY:
		return
	}

	switch c.mode {
	case setterMethodMode:
		if value != nil {
			c.compileNode(value)
		}
		c.emit(span.EndPos.Line, bytecode.POP)
		if c.isNestedInFinally() {
			c.emit(span.EndPos.Line, bytecode.GET_LOCAL8, 1)
			c.emit(span.EndPos.Line, bytecode.RETURN_FINALLY)
		} else {
			c.emit(span.EndPos.Line, bytecode.RETURN_FIRST_ARG)
		}
	case initMethodMode:
		if value != nil {
			c.compileNode(value)
		}
		c.emit(span.EndPos.Line, bytecode.POP)
		if c.isNestedInFinally() {
			c.emit(span.EndPos.Line, bytecode.SELF)
			c.emit(span.EndPos.Line, bytecode.RETURN_FINALLY)
		} else {
			c.emit(span.EndPos.Line, bytecode.RETURN_SELF)
		}
	case namespaceMode:
		if value != nil {
			c.compileNode(value)
		}
		if c.lastOpCode != bytecode.NIL {
			c.emit(span.EndPos.Line, bytecode.POP)
			c.emit(span.EndPos.Line, bytecode.NIL)
		}
		if c.isNestedInFinally() {
			c.emit(span.EndPos.Line, bytecode.RETURN_FINALLY)
		} else {
			c.emit(span.EndPos.Line, bytecode.RETURN)
		}
	default:
		if value != nil {
			c.compileNode(value)
		}
		if c.isNestedInFinally() {
			c.emit(span.EndPos.Line, bytecode.RETURN_FINALLY)
		} else {
			c.emit(span.EndPos.Line, bytecode.RETURN)
		}
	}
}

// Emit an instruction that jumps back to the given Bytecode offset.
func (c *Compiler) emitLoop(span *position.Span, startOffset int) {
	c.emit(span.EndPos.Line, bytecode.LOOP)

	offset := c.nextInstructionOffset() - startOffset + 2
	if offset > math.MaxUint16 {
		c.Errors.AddFailure(
			fmt.Sprintf("too many bytes to jump backward: %d", math.MaxUint16),
			c.newLocation(span),
		)
	}

	c.emitUint16(uint16(offset))
}

// Overwrite the placeholder operand of a jump instruction
func (c *Compiler) patchJumpWithTarget(target int, offset int, span *position.Span) {
	if target > math.MaxUint16 {
		c.Errors.AddFailure(
			fmt.Sprintf("too many bytes to jump over: %d", target),
			c.newLocation(span),
		)
		return
	}

	c.Bytecode.Instructions[offset] = byte((target >> 8) & 0xff)
	c.Bytecode.Instructions[offset+1] = byte(target & 0xff)
}

// Overwrite the placeholder operand of a jump instruction
func (c *Compiler) patchJump(offset int, span *position.Span) {
	c.patchJumpWithTarget(c.nextInstructionOffset()-offset-2, offset, span)
}

// Emit an instruction that sets a local variable or value.
func (c *Compiler) emitSetLocal(line int, index uint16) {
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
func (c *Compiler) emitSetUpvalue(line int, index uint16) {
	if index > math.MaxUint8 {
		c.emit(line, bytecode.SET_UPVALUE16)
		c.emitUint16(index)
		return
	}

	c.emit(line, bytecode.SET_UPVALUE8, byte(index))
}

// Emit an instruction that gets the value of an upvalue.
func (c *Compiler) emitGetUpvalue(line int, index uint16) {
	if index > math.MaxUint8 {
		c.emit(line, bytecode.GET_UPVALUE16)
		c.emitUint16(index)
		return
	}

	c.emit(line, bytecode.GET_UPVALUE8, byte(index))
}

// Emit an instruction that sets an upvalue.
func (c *Compiler) emitCloseUpvalue(line int, index uint16) {
	if index > math.MaxUint8 {
		c.emit(line, bytecode.CLOSE_UPVALUE16)
		c.emitUint16(index)
		return
	}

	c.emit(line, bytecode.CLOSE_UPVALUE8, byte(index))
}

// Emit an instruction that loads a value from the pool
func (c *Compiler) emitAddValue(val value.Value, span *position.Span, opCode8, opCode16 bytecode.OpCode) int {
	id, size := c.Bytecode.AddValue(val)
	switch size {
	case bytecode.UINT8_SIZE:
		c.Bytecode.AddInstruction(span.StartPos.Line, opCode8, byte(id))
	case bytecode.UINT16_SIZE:
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, uint16(id))
		c.Bytecode.AddInstruction(span.StartPos.Line, opCode16, bytes...)
	default:
		c.Errors.AddFailure(
			fmt.Sprintf("value pool limit reached: %d", math.MaxUint16),
			c.newLocation(span),
		)
		return -1
	}

	return id
}

// Emit an instruction that retrieves a constant
func (c *Compiler) emitGetConst(val value.Symbol, span *position.Span) int {
	return c.emitAddValue(
		val.ToValue(),
		span,
		bytecode.GET_CONST8,
		bytecode.GET_CONST16,
	)
}

// Add a value to the value pool and emit appropriate bytecode.
func (c *Compiler) emitLoadValue(val value.Value, span *position.Span) int {
	return c.emitAddValue(
		val,
		span,
		bytecode.LOAD_VALUE8,
		bytecode.LOAD_VALUE16,
	)
}

// Emit an instruction that instantiates an object
func (c *Compiler) emitInstantiate(args int, span *position.Span) {
	if args <= math.MaxUint8 {
		c.Bytecode.AddInstruction(span.StartPos.Line, bytecode.INSTANTIATE8, byte(args))
		return
	}

	if args <= math.MaxUint16 {
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, uint16(args))
		c.Bytecode.AddInstruction(span.StartPos.Line, bytecode.INSTANTIATE8, bytes...)
		return
	}

	c.Errors.AddFailure(
		fmt.Sprintf("constructor argument limit reached: %d", math.MaxUint16),
		c.newLocation(span),
	)
}

// Emit an instruction that sets the value of an instance variable.
func (c *Compiler) emitSetInstanceVariable(name value.Symbol, span *position.Span) int {
	return c.emitAddValue(
		name.ToValue(),
		span,
		bytecode.SET_IVAR8,
		bytecode.SET_IVAR16,
	)
}

// Emit an instruction that reads the value of an instance variable.
func (c *Compiler) emitGetInstanceVariable(name value.Symbol, span *position.Span) int {
	return c.emitAddValue(
		name.ToValue(),
		span,
		bytecode.GET_IVAR8,
		bytecode.GET_IVAR16,
	)
}

// Emit an instruction that calls a method on self
func (c *Compiler) emitCallSelf(callInfo *value.CallSiteInfo, span *position.Span) int {
	return c.emitAddValue(
		value.Ref(callInfo),
		span,
		bytecode.CALL_SELF8,
		bytecode.CALL_SELF16,
	)
}

// Emit an instruction that calls the `call` method
func (c *Compiler) emitCall(callInfo *value.CallSiteInfo, span *position.Span) int {
	return c.emitAddValue(
		value.Ref(callInfo),
		span,
		bytecode.CALL8,
		bytecode.CALL16,
	)
}

// Emit an instruction that calls a method
func (c *Compiler) emitCallMethod(callInfo *value.CallSiteInfo, span *position.Span) int {
	return c.emitAddValue(
		value.Ref(callInfo),
		span,
		bytecode.CALL_METHOD8,
		bytecode.CALL_METHOD16,
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
func (c *Compiler) defineLocal(name string, span *position.Span) *local {
	varScope := c.scopes.last()
	_, ok := varScope.localTable[name]
	if ok {
		c.Errors.AddFailure(
			fmt.Sprintf("a variable with this name has already been declared in this scope `%s`", name),
			c.newLocation(span),
		)
		return nil
	}
	return c.defineVariableInScope(varScope, name, span)
}

// Register a local variable, reusing the variable with the same name that has already been defined in this scope.
func (c *Compiler) defineLocalOverrideCurrentScope(name string, span *position.Span) *local {
	varScope := c.scopes.last()
	if currentVar, ok := varScope.localTable[name]; ok {
		return currentVar
	}
	return c.defineVariableInScope(varScope, name, span)
}

func (c *Compiler) defineVariableInScope(scope *scope, name string, span *position.Span) *local {
	if c.lastLocalIndex == math.MaxUint16 {
		c.Errors.AddFailure(
			fmt.Sprintf("exceeded the maximum number of local variables (%d): %s", math.MaxUint16, name),
			c.newLocation(span),
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
func (c *Compiler) resolveLocal(name string, span *position.Span) (*local, bool) {
	var localVal *local
	var found bool
	for i := len(c.scopes) - 1; i >= 0; i-- {
		varScope := c.scopes[i]
		local, ok := varScope.localTable[name]
		if !ok {
			continue
		}
		localVal = local
		found = true
		break
	}

	if !found {
		return localVal, false
	}

	return localVal, true
}
