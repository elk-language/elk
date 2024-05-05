// Package Compiler implements
// the Elk Bytecode Compiler.
// It takes in Elk source code and outputs
// Elk Bytecode that can be run the Elk VM.
package compiler

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"

	"github.com/elk-language/elk/token"
)

const MainName = "<main>"

// Compile the Elk source to a Bytecode chunk.
func CompileSource(sourceName string, source string) (*vm.BytecodeFunction, errors.ErrorList) {
	ast, err := parser.Parse(sourceName, source)
	if err != nil {
		return nil, err
	}

	return CompileAST(sourceName, ast)
}

// Compile the AST node to a Bytecode chunk.
func CompileAST(sourceName string, ast ast.Node) (*vm.BytecodeFunction, errors.ErrorList) {
	compiler := new(MainName, topLevelMode, position.NewLocationWithSpan(sourceName, ast.Span()))
	compiler.compileProgram(ast)

	return compiler.Bytecode, compiler.Errors
}

// Compile code for use in the REPL.
func CompileREPL(sourceName string, source string) (*Compiler, errors.ErrorList) {
	ast, err := parser.Parse(sourceName, source)
	if err != nil {
		return nil, err
	}

	compiler := new(MainName, topLevelMode, position.NewLocationWithSpan(sourceName, ast.Span()))
	compiler.compileProgram(ast)

	if compiler.Errors != nil {
		return nil, compiler.Errors
	}
	return compiler, nil
}

// Compiler mode
type mode uint8

const (
	topLevelMode mode = iota
	classMode
	mixinMode
	moduleMode
	functionMode
	setterMethodMode
	initMethodMode
	valuePatternDeclarationNode
)

// represents a local variable or value
type local struct {
	index            uint16
	singleAssignment bool
	initialised      bool
	hasUpvalue       bool // is captured by some upvalue in a closure
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
	Name             string
	Bytecode         *vm.BytecodeFunction
	Errors           errors.ErrorList
	scopes           scopes
	loopJumpSets     []*loopJumpSet
	offsetValueIds   []int // ids of integers in the value pool that represent bytecode offsets
	lastLocalIndex   int   // index of the last local variable
	maxLocalIndex    int   // max index of a local variable
	predefinedLocals int
	mode             mode
	lastOpCode       bytecode.OpCode
	patternNesting   int
	parent           *Compiler
	upvalues         []*upvalue
}

// Instantiate a new Compiler instance.
func new(name string, mode mode, loc *position.Location) *Compiler {
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
	}
	// reserve the first slot on the stack for `self`
	c.defineLocal("$self", &position.Span{}, true, true)
	switch mode {
	case topLevelMode, moduleMode, classMode, mixinMode:
		// reserve the second slot on the stack for the constant container
		c.defineLocal("$constant_container", &position.Span{}, true, true)
		// reserve the third slot on the stack for the method container
		c.defineLocal("$method_container", &position.Span{}, true, true)
		c.predefinedLocals = 3
	case functionMode, setterMethodMode, initMethodMode:
		c.predefinedLocals = 1
	}
	return c
}

// Create a new compiler based on the current compiler and compile
// new code using it.
// The new bytecode will be able to access variables defined in the previous
// chunk of bytecode produced by the previous compiler.
func (c *Compiler) CompileREPL(source string) (*Compiler, errors.ErrorList) {
	filename := c.Bytecode.Location.Filename
	ast, err := parser.Parse(filename, source)
	if err != nil {
		return nil, err
	}

	compiler := new(MainName, topLevelMode, position.NewLocationWithSpan(filename, ast.Span()))
	compiler.predefinedLocals = c.maxLocalIndex + 1
	compiler.scopes = c.scopes
	compiler.lastLocalIndex = c.lastLocalIndex
	compiler.maxLocalIndex = c.maxLocalIndex
	compiler.compileProgram(ast)

	if compiler.Errors != nil {
		return nil, compiler.Errors
	}
	return compiler, nil
}

// Entry point to the compilation process
func (c *Compiler) compileProgram(node ast.Node) {
	c.compileNode(node)
	c.emitReturn(node.Span(), nil)
	c.prepLocals()
}

// Entry point for compiling the body of a function.
func (c *Compiler) compileFunction(span *position.Span, parameters []ast.ParameterNode, body []ast.StatementNode) {
	if len(parameters) > 0 {
		c.Bytecode.SetParameters(make([]value.Symbol, 0, len(parameters)))
	}
	var positionalRestParamSeen bool

	for _, param := range parameters {
		p := param.(*ast.FormalParameterNode)
		pSpan := p.Span()

		switch p.Kind {
		case ast.NamedRestParameterKind:
			c.Bytecode.SetNamedRestParameter(true)
		case ast.PositionalRestParameterKind:
			positionalRestParamSeen = true
			c.Bytecode.IncrementPostRestParameterCount()
		default:
			if positionalRestParamSeen {
				c.Bytecode.IncrementPostRestParameterCount()
			}
		}

		local := c.defineLocal(p.Name, pSpan, false, true)
		if local == nil {
			return
		}
		c.Bytecode.AddParameterString(p.Name)
		c.predefinedLocals++

		if p.Initialiser != nil {
			c.Bytecode.IncrementOptionalParameterCount()

			c.emitGetLocal(span.StartPos.Line, local.index)
			jump := c.emitJump(pSpan.StartPos.Line, bytecode.JUMP_UNLESS_UNDEF)

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

// Entry point for compiling the body of a method.
func (c *Compiler) compileMethod(span *position.Span, parameters []ast.ParameterNode, body []ast.StatementNode) {
	if len(parameters) > 0 {
		c.Bytecode.SetParameters(make([]value.Symbol, 0, len(parameters)))
	}
	var positionalRestParamSeen bool

	for _, param := range parameters {
		p := param.(*ast.MethodParameterNode)
		pSpan := p.Span()

		switch p.Kind {
		case ast.NamedRestParameterKind:
			c.Bytecode.SetNamedRestParameter(true)
		case ast.PositionalRestParameterKind:
			positionalRestParamSeen = true
			c.Bytecode.IncrementPostRestParameterCount()
		default:
			if positionalRestParamSeen {
				c.Bytecode.IncrementPostRestParameterCount()
			}
		}

		local := c.defineLocal(p.Name, pSpan, false, true)
		if local == nil {
			return
		}
		c.Bytecode.AddParameterString(p.Name)
		c.predefinedLocals++

		if p.Initialiser != nil {
			c.Bytecode.IncrementOptionalParameterCount()

			c.emitGetLocal(span.StartPos.Line, local.index)
			jump := c.emitJump(pSpan.StartPos.Line, bytecode.JUMP_UNLESS_UNDEF)

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

// Entry point for compiling the body of a Module, Class, Mixin, Struct.
func (c *Compiler) compileModule(node ast.Node) {
	span := node.Span()
	switch n := node.(type) {
	case *ast.ClassDeclarationNode:
		c.compileStatements(n.Body, span)
	case *ast.ModuleDeclarationNode:
		c.compileStatements(n.Body, span)
	case *ast.MixinDeclarationNode:
		c.compileStatements(n.Body, span)
	case *ast.SingletonBlockExpressionNode:
		c.compileStatements(n.Body, span)
	default:
		c.Errors.Add(
			fmt.Sprintf("incorrect module type %#v", n),
			c.newLocation(span),
		)
		return
	}
	c.emitReturn(span, nil)
	c.prepLocals()
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
		currentValue := c.Bytecode.Values[id].(value.SmallInt)
		c.Bytecode.Values[id] = currentValue + value.SmallInt(len(newInstructions))
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
		c.Errors.Add(
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

	c.Errors.Add(
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

func (c *Compiler) compileNode(node ast.Node) {
	switch node := node.(type) {
	case *ast.ProgramNode:
		c.compileStatements(node.Body, node.Span())
	case *ast.ExpressionStatementNode:
		c.compileNode(node.Expression)
	case *ast.LabeledExpressionNode:
		c.labeledExpression(node)
	case *ast.ConstantLookupNode:
		c.constantLookup(node)
	case *ast.GenericConstantNode:
		c.compileNode(node.Constant)
	case *ast.SelfLiteralNode:
		c.emit(node.Span().StartPos.Line, bytecode.SELF)
	case *ast.AssignmentExpressionNode:
		c.assignment(node)
	case *ast.AliasDeclarationNode:
		c.aliasDeclaration(node)
	case *ast.GetterDeclarationNode:
		c.getterDeclaration(node)
	case *ast.SetterDeclarationNode:
		c.setterDeclaration(node)
	case *ast.AccessorDeclarationNode:
		c.accessorDeclaration(node)
	case *ast.ClassDeclarationNode:
		c.classDeclaration(node)
	case *ast.ModuleDeclarationNode:
		c.moduleDeclaration(node)
	case *ast.MixinDeclarationNode:
		c.mixinDeclaration(node)
	case *ast.MethodDefinitionNode:
		c.methodDefinition(node)
	case *ast.FunctionLiteralNode:
		c.functionLiteral(node)
	case *ast.InitDefinitionNode:
		c.initDefinition(node)
	case *ast.IncludeExpressionNode:
		c.includeExpression(node)
	case *ast.ExtendExpressionNode:
		c.extendExpression(node)
	case *ast.SingletonBlockExpressionNode:
		c.singletonBlock(node)
	case *ast.SwitchExpressionNode:
		c.switchExpression(node)
	case *ast.SubscriptExpressionNode:
		c.subscriptExpression(node)
	case *ast.NilSafeSubscriptExpressionNode:
		c.nilSafeSubscriptExpression(node)
	case *ast.AttributeAccessNode:
		c.attributeAccess(node)
	case *ast.ConstructorCallNode:
		c.constructorCall(node)
	case *ast.MethodCallNode:
		c.methodCall(node)
	case *ast.CallNode:
		c.call(node)
	case *ast.ReceiverlessMethodCallNode:
		c.receiverlessMethodCall(node)
	case *ast.ReturnExpressionNode:
		c.returnExpression(node)
	case *ast.VariablePatternDeclarationNode:
		c.variablePatternDeclaration(node)
	case *ast.VariableDeclarationNode:
		c.variableDeclaration(node)
	case *ast.ValuePatternDeclarationNode:
		c.valuePatternDeclaration(node)
	case *ast.ValueDeclarationNode:
		c.valueDeclaration(node)
	case *ast.PublicIdentifierNode:
		c.localVariableAccess(node.Value, node.Span())
	case *ast.PrivateIdentifierNode:
		c.localVariableAccess(node.Value, node.Span())
	case *ast.InstanceVariableNode:
		c.instanceVariableAccess(node.Value, node.Span())
	case *ast.BinaryExpressionNode:
		c.binaryExpression(node)
	case *ast.LogicalExpressionNode:
		c.logicalExpression(node)
	case *ast.UnaryExpressionNode:
		c.unaryExpression(node)
	case *ast.RangeLiteralNode:
		c.rangeLiteral(node)
	case *ast.HashSetLiteralNode:
		c.hashSetLiteral(node)
	case *ast.HashMapLiteralNode:
		c.hashMapLiteral(node)
	case *ast.HashRecordLiteralNode:
		c.hashRecordLiteral(node)
	case *ast.ArrayTupleLiteralNode:
		c.arrayTupleLiteral(node)
	case *ast.WordArrayTupleLiteralNode:
		c.wordArrayTupleLiteral(node)
	case *ast.SymbolArrayTupleLiteralNode:
		c.symbolArrayTupleLiteral(node)
	case *ast.BinArrayTupleLiteralNode:
		c.binArrayTupleLiteral(node)
	case *ast.HexArrayTupleLiteralNode:
		c.hexArrayTupleLiteral(node)
	case *ast.ArrayListLiteralNode:
		c.arrayListLiteral(node)
	case *ast.WordArrayListLiteralNode:
		c.wordArrayListLiteral(node)
	case *ast.SymbolArrayListLiteralNode:
		c.symbolArrayListLiteral(node)
	case *ast.BinArrayListLiteralNode:
		c.binArrayListLiteral(node)
	case *ast.HexArrayListLiteralNode:
		c.hexArrayListLiteral(node)
	case *ast.WordHashSetLiteralNode:
		c.wordHashSetLiteral(node)
	case *ast.SymbolHashSetLiteralNode:
		c.symbolHashSetLiteral(node)
	case *ast.BinHashSetLiteralNode:
		c.binHashSetLiteral(node)
	case *ast.HexHashSetLiteralNode:
		c.hexHashSetLiteral(node)
	case *ast.UninterpolatedRegexLiteralNode:
		c.uninterpolatedRegexLiteral(node)
	case *ast.InterpolatedRegexLiteralNode:
		c.interpolatedRegexLiteral(node)
	case *ast.RawStringLiteralNode:
		c.emitValue(value.String(node.Value), node.Span())
	case *ast.DoubleQuotedStringLiteralNode:
		c.emitValue(value.String(node.Value), node.Span())
	case *ast.InterpolatedStringLiteralNode:
		c.interpolatedStringLiteral(node)
	case *ast.InterpolatedSymbolLiteralNode:
		c.interpolatedSymbolLiteral(node)
	case *ast.CharLiteralNode:
		c.emitValue(value.Char(node.Value), node.Span())
	case *ast.RawCharLiteralNode:
		c.emitValue(value.Char(node.Value), node.Span())
	case *ast.FalseLiteralNode:
		c.emit(node.Span().StartPos.Line, bytecode.FALSE)
	case *ast.TrueLiteralNode:
		c.emit(node.Span().StartPos.Line, bytecode.TRUE)
	case *ast.NilLiteralNode:
		c.emit(node.Span().StartPos.Line, bytecode.NIL)
	case *ast.EmptyStatementNode:
	case *ast.ThrowExpressionNode:
		c.throwExpression(node)
	case *ast.DoExpressionNode:
		c.doExpression(node)
	case *ast.IfExpressionNode:
		c.ifExpression(false, node.Condition, node.ThenBody, node.ElseBody, node.Span())
	case *ast.UnlessExpressionNode:
		c.ifExpression(true, node.Condition, node.ThenBody, node.ElseBody, node.Span())
	case *ast.ModifierIfElseNode:
		c.modifierIfExpression(false, node.Condition, node.ThenExpression, node.ElseExpression, node.Span())
	case *ast.ModifierNode:
		c.modifierExpression("", node)
	case *ast.DocCommentNode:
		c.docComment(node)
	case *ast.BreakExpressionNode:
		c.breakExpression(node)
	case *ast.ContinueExpressionNode:
		c.continueExpression(node)
	case *ast.LoopExpressionNode:
		c.loopExpression("", node.ThenBody, node.Span())
	case *ast.WhileExpressionNode:
		c.whileExpression("", node)
	case *ast.UntilExpressionNode:
		c.untilExpression("", node)
	case *ast.NumericForExpressionNode:
		c.numericForExpression("", node)
	case *ast.ForInExpressionNode:
		c.forInExpression("", node)
	case *ast.ModifierForInNode:
		c.modifierForIn("", node)
	case *ast.PostfixExpressionNode:
		c.postfixExpression(node)
	case *ast.SimpleSymbolLiteralNode:
		c.emitValue(value.ToSymbol(node.Content), node.Span())
	case *ast.IntLiteralNode:
		c.intLiteral(node)
	case *ast.Int8LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 8)
		if err != nil {
			c.Errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		// BENCHMARK: Compare with storing
		// ints inline in Bytecode instead of as constants.
		c.emitValue(value.Int8(i), node.Span())
	case *ast.Int16LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 16)
		if err != nil {
			c.Errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitValue(value.Int16(i), node.Span())
	case *ast.Int32LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 32)
		if err != nil {
			c.Errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitValue(value.Int32(i), node.Span())
	case *ast.Int64LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 64)
		if err != nil {
			c.Errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitValue(value.Int64(i), node.Span())
	case *ast.UInt8LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 8)
		if err != nil {
			c.Errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitValue(value.UInt8(i), node.Span())
	case *ast.UInt16LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 16)
		if err != nil {
			c.Errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitValue(value.UInt16(i), node.Span())
	case *ast.UInt32LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 32)
		if err != nil {
			c.Errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitValue(value.UInt32(i), node.Span())
	case *ast.UInt64LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 64)
		if err != nil {
			c.Errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitValue(value.UInt64(i), node.Span())
	case *ast.FloatLiteralNode:
		f, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			c.Errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitValue(value.Float(f), node.Span())
	case *ast.BigFloatLiteralNode:
		f, err := value.ParseBigFloat(node.Value)
		if err != nil {
			c.Errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitValue(f, node.Span())
	case *ast.Float64LiteralNode:
		f, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			c.Errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitValue(value.Float64(f), node.Span())
	case *ast.Float32LiteralNode:
		f, err := strconv.ParseFloat(node.Value, 32)
		if err != nil {
			c.Errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitValue(value.Float32(f), node.Span())

	case nil:
	default:
		c.Errors.Add(
			fmt.Sprintf("compilation of this node has not been implemented: %T", node),
			c.newLocation(node.Span()),
		)
	}
}

func (c *Compiler) throwExpression(node *ast.ThrowExpressionNode) {
	span := node.Span()
	if node.Value != nil {
		c.compileNode(node.Value)
	} else {
		c.emitValue(value.NewError(value.ErrorClass, "error"), span)
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

func (c *Compiler) doExpression(node *ast.DoExpressionNode) {
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
		// pop the boolean return value of the pattern
		c.emit(span.StartPos.Line, bytecode.POP)

		c.compileStatements(catchNode.Body, catchNode.Span())

		if len(node.Finally) < 1 {
			// pop the thrown value and the stack trace, leaving the return value of the catch
			c.emit(span.EndPos.Line, bytecode.POP_N_SKIP_ONE, 2)
		}
		jump := c.emitJump(span.EndPos.Line, bytecode.JUMP)
		jumpsToEndOfCatch = append(jumpsToEndOfCatch, jump)

		c.patchJump(jumpOverCatchBody, span)
		// pop the boolean return value of the pattern after jump
		c.emit(span.EndPos.Line, bytecode.POP)
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
		jumpOverFinallyBreakOrContinueOffset := c.emitJump(span.EndPos.Line, bytecode.JUMP_UNLESS_UNDEF)
		c.emit(span.EndPos.Line, bytecode.POP_N, 2)
		c.emit(span.EndPos.Line, bytecode.JUMP_TO_FINALLY)
		c.patchJump(jumpOverFinallyBreakOrContinueOffset, span)

		jumpToRethrowOffset := c.emitJump(span.EndPos.Line, bytecode.JUMP_IF)
		jumpToFinallyReturnOffset := c.emitJump(span.EndPos.Line, bytecode.JUMP_IF_NIL)
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

func (c *Compiler) breakExpression(node *ast.BreakExpressionNode) {
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

	c.emitValue(value.SmallInt(finallyCount), span)
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

func (c *Compiler) continueExpression(node *ast.ContinueExpressionNode) {
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

	c.emitValue(value.SmallInt(finallyCount), span)
	c.emit(span.StartPos.Line, bytecode.JUMP_TO_FINALLY)
}

// Patch loop jump addresses for `break` and `continue` expressions.
func (c *Compiler) patchLoopJumps(continueOffset int) {
	lastLoopJumpSet := c.loopJumpSets[len(c.loopJumpSets)-1]
	for _, loopJump := range lastLoopJumpSet.loopJumps {
		switch loopJump.typ {
		case breakFinallyLoopJump:
			c.Bytecode.Values[loopJump.offset] = value.SmallInt(c.nextInstructionOffset())
		case continueFinallyLoopJump:
			c.Bytecode.Values[loopJump.offset] = value.SmallInt(continueOffset)
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

func (c *Compiler) loopExpression(label string, body []ast.StatementNode, span *position.Span) {
	c.enterScope(label, loopScopeType)
	c.initLoopJumpSet(label, false)

	start := c.nextInstructionOffset()
	if c.compileStatementsOk(body, span) {
		c.emit(span.EndPos.Line, bytecode.POP)
	}
	c.emitLoop(span, start)

	c.leaveScope(span.EndPos.Line)
	c.patchLoopJumps(start)
}

func (c *Compiler) whileExpression(label string, node *ast.WhileExpressionNode) {
	span := node.Span()

	if result := resolve(node.Condition); result != nil {
		if value.Falsy(result) {
			// the loop won't run at all
			// it can be optimised into a simple NIL operation
			c.emit(span.StartPos.Line, bytecode.NIL)
			return
		}

		// the loop is endless
		c.loopExpression(label, node.ThenBody, span)
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
	// pop the condition value
	// and the return value of the last iteration
	c.emit(span.StartPos.Line, bytecode.POP_N, 2)

	// loop body
	c.compileStatements(node.ThenBody, span)

	c.leaveScope(span.EndPos.Line)
	// jump to loop condition
	c.emitLoop(span, start)

	// after loop
	c.patchJump(loopBodyOffset, span)
	// pop the condition value
	c.emit(span.EndPos.Line, bytecode.POP)

	c.leaveScope(span.EndPos.Line)
	c.patchLoopJumps(start)
}

func (c *Compiler) modifierWhileExpression(label string, node *ast.ModifierNode) {
	span := node.Span()

	body := node.Left
	condition := node.Right

	var conditionIsStaticFalsy bool

	if result := resolve(condition); result != nil {
		if value.Truthy(result) {
			// the loop is endless
			c.loopExpression(label, ast.ExpressionToStatements(body), span)
			return
		}

		// the loop will only iterate once
		conditionIsStaticFalsy = true
	}

	c.enterScope(label, loopScopeType)
	c.initLoopJumpSet(label, true)

	// loop start
	start := c.nextInstructionOffset()
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
	// pop the condition value
	// and the return value of the last iteration
	c.emit(span.StartPos.Line, bytecode.POP_N, 2)

	// jump to loop start
	c.emitLoop(span, start)

	// after loop
	c.patchJump(loopBodyOffset, span)
	// pop the condition value
	c.emit(span.EndPos.Line, bytecode.POP)

	c.leaveScope(span.EndPos.Line)
	c.patchLoopJumps(continueOffset)
}

func (c *Compiler) modifierUntilExpression(label string, node *ast.ModifierNode) {
	span := node.Span()

	body := node.Left
	condition := node.Right

	var conditionIsStaticTruthy bool

	if result := resolve(condition); result != nil {
		if value.Falsy(result) {
			// the loop is endless
			c.loopExpression(label, ast.ExpressionToStatements(body), span)
			return
		}

		// the loop will only iterate once
		conditionIsStaticTruthy = true
	}

	c.enterScope(label, loopScopeType)
	c.initLoopJumpSet(label, true)

	// loop start
	start := c.nextInstructionOffset()
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
	// pop the condition value
	// and the return value of the last iteration
	c.emit(span.StartPos.Line, bytecode.POP_N, 2)

	// jump to loop start
	c.emitLoop(span, start)

	// after loop
	c.patchJump(loopBodyOffset, span)
	// pop the condition value
	c.emit(span.EndPos.Line, bytecode.POP)

	c.leaveScope(span.EndPos.Line)
	c.patchLoopJumps(continueOffset)
}

func (c *Compiler) untilExpression(label string, node *ast.UntilExpressionNode) {
	span := node.Span()

	if result := resolve(node.Condition); result != nil {
		if value.Falsy(result) {
			// the loop is endless
			c.loopExpression(label, node.ThenBody, span)
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
	var loopBodyOffset int

	// loop condition eg. `i > 5`
	c.compileNode(node.Condition)
	// jump past the loop if the condition is truthy
	loopBodyOffset = c.emitJump(span.StartPos.Line, bytecode.JUMP_IF)
	// pop the condition value
	// and the return value of the last iteration
	c.emit(span.StartPos.Line, bytecode.POP_N, 2)

	// loop body
	c.compileStatements(node.ThenBody, span)

	// jump to loop condition
	c.emitLoop(span, start)

	// after loop
	c.patchJump(loopBodyOffset, span)
	// pop the condition value
	c.emit(span.EndPos.Line, bytecode.POP)

	c.leaveScope(span.EndPos.Line)
	c.patchLoopJumps(start)
}

// Compile a labeled expression eg. `$foo: println("bar")`
func (c *Compiler) labeledExpression(node *ast.LabeledExpressionNode) {
	switch expr := node.Expression.(type) {
	case *ast.WhileExpressionNode:
		c.whileExpression(node.Label, expr)
	case *ast.UntilExpressionNode:
		c.untilExpression(node.Label, expr)
	case *ast.LoopExpressionNode:
		c.loopExpression(node.Label, expr.ThenBody, expr.Span())
	case *ast.NumericForExpressionNode:
		c.numericForExpression(node.Label, expr)
	case *ast.ForInExpressionNode:
		c.forInExpression(node.Label, expr)
	case *ast.ModifierForInNode:
		c.modifierForIn(node.Label, expr)
	case *ast.ModifierNode:
		c.modifierExpression(node.Label, expr)
	default:
		c.compileNode(node.Expression)
	}
}

// Compile a constant lookup expressions eg. `Foo::Bar`
func (c *Compiler) constantLookup(node *ast.ConstantLookupNode) {
	if node.Left == nil {
		c.emit(node.Span().StartPos.Line, bytecode.ROOT)
	} else {
		c.compileNode(node.Left)
	}

	switch r := node.Right.(type) {
	case *ast.PublicConstantNode:
		c.emitGetModConst(value.ToSymbol(r.Value), node.Span())
	default:
		c.Errors.Add(
			fmt.Sprintf("incorrect right side of constant lookup: %T", node.Right),
			c.newLocation(node.Span()),
		)
	}
}

// Compile a for in loop eg. `for i in [1, 2] then println(i)`
func (c *Compiler) forInExpression(label string, node *ast.ForInExpressionNode) {
	c.compileForIn(
		label,
		node.Parameter,
		node.InExpression,
		func() {
			c.compileStatements(node.ThenBody, node.Span())
		},
		node.Span(),
		false,
	)
}

// Compile a for in loop eg. `println(i) for i in [1, 2]`
func (c *Compiler) modifierForIn(label string, node *ast.ModifierForInNode) {
	c.compileForIn(
		label,
		node.Parameter,
		node.InExpression,
		func() {
			c.compileNode(node.ThenExpression)
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
	iteratorVar := c.defineLocal(iteratorVarName, span, true, true)
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
		paramVar := c.defineLocal(p.Value, param.Span(), true, true)
		c.emitSetLocal(param.Span().StartPos.Line, paramVar.index)
		c.emit(param.Span().EndPos.Line, bytecode.POP)
	case *ast.PublicIdentifierNode:
		paramVar := c.defineLocal(p.Value, param.Span(), true, true)
		c.emitSetLocal(param.Span().StartPos.Line, paramVar.index)
		c.emit(param.Span().EndPos.Line, bytecode.POP)
	default:
		c.pattern(param)
		jumpOverErrorOffset := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF)
		c.emit(span.EndPos.Line, bytecode.POP)

		c.emitValue(
			value.NewError(
				value.PatternNotMatchedErrorClass,
				"assigned value does not match the pattern defined in for in loop",
			),
			span,
		)
		c.emit(span.EndPos.Line, bytecode.THROW)

		c.patchJump(jumpOverErrorOffset, span)
		c.emit(span.EndPos.Line, bytecode.POP)
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
func (c *Compiler) numericForExpression(label string, node *ast.NumericForExpressionNode) {
	span := node.Span()

	if node.Initialiser == nil && node.Condition == nil && node.Increment == nil {
		// the loop is endless
		c.loopExpression(label, node.ThenBody, span)
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
		// pop the condition value
		// and the return value of the last iteration
		c.emit(span.EndPos.Line, bytecode.POP_N, 2)
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
		// pop the condition value
		c.emit(span.EndPos.Line, bytecode.POP)
	}

	c.leaveScope(span.EndPos.Line)
	c.patchLoopJumps(continueOffset)
}

func (c *Compiler) complexSetterCall(opCode bytecode.OpCode, node *ast.AttributeAccessNode, val ast.ExpressionNode, span *position.Span) {
	c.compileNode(node.Receiver)
	name := value.ToSymbol(node.AttributeName)
	callInfo := value.NewCallSiteInfo(name, 0, nil)
	c.emitCallMethod(callInfo, node.Span())

	c.compileNode(val)
	c.emit(span.StartPos.Line, opCode)

	c.emitSetterCall(node.AttributeName, node.Span())
}

func (c *Compiler) emitSetterCall(name string, span *position.Span) {
	nameSymbol := value.ToSymbol(name + "=")
	callInfo := value.NewCallSiteInfo(nameSymbol, 1, nil)
	c.emitCallMethod(callInfo, span)
}

func (c *Compiler) emitGetterCall(name string, span *position.Span) {
	nameSymbol := value.ToSymbol(name)
	callInfo := value.NewCallSiteInfo(nameSymbol, 0, nil)
	c.emitCallMethod(callInfo, span)
}

func (c *Compiler) postfixExpression(node *ast.PostfixExpressionNode) {
	switch n := node.Expression.(type) {
	case *ast.PublicIdentifierNode:
		// get variable value
		c.localVariableAccess(n.Value, n.Span())

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
		c.localVariableAccess(n.Value, n.Span())

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
			c.Errors.Add(
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
		callInfo := value.NewCallSiteInfo(name, 0, nil)
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
		c.Errors.Add(
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

		jump := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF)

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

		jump := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)

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

		nilJump := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF_NIL)
		nonNilJump := c.emitJump(span.StartPos.Line, bytecode.JUMP)

		// if nil
		c.patchJump(nilJump, span)
		c.emit(span.StartPos.Line, bytecode.POP)
		c.compileNode(node.Right)
		c.emitSetterCall(attr.AttributeName, span)

		// if not nil
		c.patchJump(nonNilJump, span)
	case token.PLUS_EQUAL:
		c.complexSetterCall(bytecode.ADD, attr, node.Right, node.Span())
	case token.MINUS_EQUAL:
		c.complexSetterCall(bytecode.SUBTRACT, attr, node.Right, node.Span())
	case token.STAR_EQUAL:
		c.complexSetterCall(bytecode.MULTIPLY, attr, node.Right, node.Span())
	case token.SLASH_EQUAL:
		c.complexSetterCall(bytecode.DIVIDE, attr, node.Right, node.Span())
	case token.STAR_STAR_EQUAL:
		c.complexSetterCall(bytecode.EXPONENTIATE, attr, node.Right, node.Span())
	case token.PERCENT_EQUAL:
		c.complexSetterCall(bytecode.MODULO, attr, node.Right, node.Span())
	case token.LBITSHIFT_EQUAL:
		c.complexSetterCall(bytecode.LBITSHIFT, attr, node.Right, node.Span())
	case token.LTRIPLE_BITSHIFT_EQUAL:
		c.complexSetterCall(bytecode.LOGIC_LBITSHIFT, attr, node.Right, node.Span())
	case token.RBITSHIFT_EQUAL:
		c.complexSetterCall(bytecode.RBITSHIFT, attr, node.Right, node.Span())
	case token.RTRIPLE_BITSHIFT_EQUAL:
		c.complexSetterCall(bytecode.LOGIC_RBITSHIFT, attr, node.Right, node.Span())
	case token.AND_EQUAL:
		c.complexSetterCall(bytecode.BITWISE_AND, attr, node.Right, node.Span())
	case token.OR_EQUAL:
		c.complexSetterCall(bytecode.BITWISE_OR, attr, node.Right, node.Span())
	case token.XOR_EQUAL:
		c.complexSetterCall(bytecode.BITWISE_XOR, attr, node.Right, node.Span())
	default:
		c.Errors.Add(fmt.Sprintf("unknown binary operator: %s", node.Op.String()), c.newLocation(node.Span()))
	}
}

func (c *Compiler) complexInstanceVariableAssignment(ivarSymbol value.Symbol, valueNode ast.ExpressionNode, opcode bytecode.OpCode, span *position.Span) {
	c.emitGetInstanceVariable(ivarSymbol, span)
	c.compileNode(valueNode)
	c.emit(span.StartPos.Line, opcode)
	c.emitSetInstanceVariable(ivarSymbol, span)
}

func (c *Compiler) instanceVariableAssignment(node *ast.AssignmentExpressionNode, ivar *ast.InstanceVariableNode) {
	switch c.mode {
	case topLevelMode:
		c.Errors.Add(
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

		jump := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF)

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

		jump := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)

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

		nilJump := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF_NIL)
		nonNilJump := c.emitJump(span.StartPos.Line, bytecode.JUMP)

		// if nil
		c.patchJump(nilJump, span)
		c.emit(span.StartPos.Line, bytecode.POP)
		c.compileNode(node.Right)
		c.emitSetInstanceVariable(ivarSymbol, ivar.Span())

		// if not nil
		c.patchJump(nonNilJump, span)
	case token.PLUS_EQUAL:
		c.complexInstanceVariableAssignment(ivarSymbol, node.Right, bytecode.ADD, node.Span())
	case token.MINUS_EQUAL:
		c.complexInstanceVariableAssignment(ivarSymbol, node.Right, bytecode.SUBTRACT, node.Span())
	case token.STAR_EQUAL:
		c.complexInstanceVariableAssignment(ivarSymbol, node.Right, bytecode.MULTIPLY, node.Span())
	case token.SLASH_EQUAL:
		c.complexInstanceVariableAssignment(ivarSymbol, node.Right, bytecode.DIVIDE, node.Span())
	case token.STAR_STAR_EQUAL:
		c.complexInstanceVariableAssignment(ivarSymbol, node.Right, bytecode.EXPONENTIATE, node.Span())
	case token.PERCENT_EQUAL:
		c.complexInstanceVariableAssignment(ivarSymbol, node.Right, bytecode.MODULO, node.Span())
	case token.LBITSHIFT_EQUAL:
		c.complexInstanceVariableAssignment(ivarSymbol, node.Right, bytecode.LBITSHIFT, node.Span())
	case token.LTRIPLE_BITSHIFT_EQUAL:
		c.complexInstanceVariableAssignment(ivarSymbol, node.Right, bytecode.LOGIC_LBITSHIFT, node.Span())
	case token.RBITSHIFT_EQUAL:
		c.complexInstanceVariableAssignment(ivarSymbol, node.Right, bytecode.RBITSHIFT, node.Span())
	case token.RTRIPLE_BITSHIFT_EQUAL:
		c.complexInstanceVariableAssignment(ivarSymbol, node.Right, bytecode.LOGIC_RBITSHIFT, node.Span())
	case token.AND_EQUAL:
		c.complexInstanceVariableAssignment(ivarSymbol, node.Right, bytecode.BITWISE_AND, node.Span())
	case token.OR_EQUAL:
		c.complexInstanceVariableAssignment(ivarSymbol, node.Right, bytecode.BITWISE_OR, node.Span())
	case token.XOR_EQUAL:
		c.complexInstanceVariableAssignment(ivarSymbol, node.Right, bytecode.BITWISE_XOR, node.Span())
	default:
		c.Errors.Add(fmt.Sprintf("unknown binary operator: %s", node.Op.String()), c.newLocation(node.Span()))
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

		jump := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF)

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

		jump := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)

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

		nilJump := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF_NIL)
		nonNilJump := c.emitJump(span.StartPos.Line, bytecode.JUMP)

		// if nil
		c.patchJump(nilJump, span)
		c.emit(span.StartPos.Line, bytecode.POP)
		c.compileNode(node.Right)
		c.emit(node.Span().EndPos.Line, bytecode.SUBSCRIPT_SET)

		// if not nil
		c.patchJump(nonNilJump, span)
		c.emit(span.StartPos.Line, bytecode.POP_N_SKIP_ONE, 2)
	case token.PLUS_EQUAL:
		c.complexSubscriptAssignment(subscript, node.Right, bytecode.ADD, node.Span())
	case token.MINUS_EQUAL:
		c.complexSubscriptAssignment(subscript, node.Right, bytecode.SUBTRACT, node.Span())
	case token.STAR_EQUAL:
		c.complexSubscriptAssignment(subscript, node.Right, bytecode.MULTIPLY, node.Span())
	case token.SLASH_EQUAL:
		c.complexSubscriptAssignment(subscript, node.Right, bytecode.DIVIDE, node.Span())
	case token.STAR_STAR_EQUAL:
		c.complexSubscriptAssignment(subscript, node.Right, bytecode.EXPONENTIATE, node.Span())
	case token.PERCENT_EQUAL:
		c.complexSubscriptAssignment(subscript, node.Right, bytecode.MODULO, node.Span())
	case token.LBITSHIFT_EQUAL:
		c.complexSubscriptAssignment(subscript, node.Right, bytecode.LBITSHIFT, node.Span())
	case token.LTRIPLE_BITSHIFT_EQUAL:
		c.complexSubscriptAssignment(subscript, node.Right, bytecode.LOGIC_LBITSHIFT, node.Span())
	case token.RBITSHIFT_EQUAL:
		c.complexSubscriptAssignment(subscript, node.Right, bytecode.RBITSHIFT, node.Span())
	case token.RTRIPLE_BITSHIFT_EQUAL:
		c.complexSubscriptAssignment(subscript, node.Right, bytecode.LOGIC_RBITSHIFT, node.Span())
	case token.AND_EQUAL:
		c.complexSubscriptAssignment(subscript, node.Right, bytecode.BITWISE_AND, node.Span())
	case token.OR_EQUAL:
		c.complexSubscriptAssignment(subscript, node.Right, bytecode.BITWISE_OR, node.Span())
	case token.XOR_EQUAL:
		c.complexSubscriptAssignment(subscript, node.Right, bytecode.BITWISE_XOR, node.Span())
	default:
		c.Errors.Add(fmt.Sprintf("unknown binary operator: %s", node.Op.String()), c.newLocation(node.Span()))
	}
}

func (c *Compiler) assignment(node *ast.AssignmentExpressionNode) {
	switch n := node.Left.(type) {
	case *ast.PublicIdentifierNode:
		c.localVariableAssignment(n.Value, node.Op, node.Right, node.Span())
	case *ast.PrivateIdentifierNode:
		c.localVariableAssignment(n.Value, node.Op, node.Right, node.Span())
	case *ast.SubscriptExpressionNode:
		c.subscriptAssignment(node, n)
	case *ast.ConstantLookupNode:
		if node.Op.Type != token.COLON_EQUAL {
			c.Errors.Add(
				fmt.Sprintf("cannot assign constants using `%s`", node.Op.StringValue()),
				c.newLocation(node.Span()),
			)
		}
		c.compileNode(node.Right)
		if n.Left == nil {
			c.emit(node.Span().StartPos.Line, bytecode.ROOT)
		} else {
			c.compileNode(n.Left)
		}

		switch r := n.Right.(type) {
		case *ast.PublicConstantNode:
			c.emitDefModConst(value.ToSymbol(r.Value), n.Span())
		default:
			c.Errors.Add(
				fmt.Sprintf("incorrect right side of constant lookup: %T", n.Right),
				c.newLocation(n.Right.Span()),
			)
		}
	case *ast.PublicConstantNode:
		c.compileSimpleConstantAssignment(n.Value, node.Op, node.Right, node.Span())
	case *ast.PrivateConstantNode:
		c.compileSimpleConstantAssignment(n.Value, node.Op, node.Right, node.Span())
	case *ast.InstanceVariableNode:
		c.instanceVariableAssignment(node, n)
	case *ast.AttributeAccessNode:
		c.attributeAssignment(node, n)
	default:
		c.Errors.Add(
			fmt.Sprintf("cannot assign to: %T", node.Left),
			c.newLocation(node.Span()),
		)
	}
}

func (c *Compiler) compileSimpleConstantAssignment(name string, op *token.Token, right ast.ExpressionNode, span *position.Span) {
	if op.Type != token.COLON_EQUAL {
		c.Errors.Add(
			fmt.Sprintf("cannot assign constants using `%s`", op.StringValue()),
			c.newLocation(span),
		)
	}
	c.compileNode(right)
	c.emit(span.StartPos.Line, bytecode.CONSTANT_CONTAINER)
	c.emitDefModConst(value.ToSymbol(name), span)
}

func (c *Compiler) complexAssignment(name string, valueNode ast.ExpressionNode, opcode bytecode.OpCode, span *position.Span) {
	local, upvalue, ok := c.localVariableAccess(name, span)
	if !ok {
		return
	}
	c.compileNode(valueNode)
	c.emit(span.StartPos.Line, opcode)

	if local.initialised && local.singleAssignment {
		c.Errors.Add(
			fmt.Sprintf("cannot reassign a val: %s", name),
			c.newLocation(span),
		)
	}
	local.initialised = true
	if upvalue != nil {
		c.emitSetUpvalue(span.StartPos.Line, upvalue.index)
	} else {
		c.emitSetLocal(span.StartPos.Line, local.index)
	}
}

// Return the offset of the Bytecode next instruction.
func (c *Compiler) nextInstructionOffset() int {
	return len(c.Bytecode.Instructions)
}

func (c *Compiler) setLocalWithoutValue(name string, span *position.Span) {
	if local, ok := c.resolveLocal(name, span); ok {
		if local.initialised && local.singleAssignment {
			c.Errors.Add(
				fmt.Sprintf("cannot reassign a val: %s", name),
				c.newLocation(span),
			)
		}
		local.initialised = true
		c.emitSetLocal(span.StartPos.Line, local.index)
	} else if upvalue, ok := c.resolveUpvalue(name, span); ok {
		local := upvalue.local
		if local.initialised && local.singleAssignment {
			c.Errors.Add(
				fmt.Sprintf("cannot reassign a val: %s", name),
				c.newLocation(span),
			)
		}
		local.initialised = true
		c.emitSetUpvalue(span.StartPos.Line, upvalue.index)
	} else {
		c.Errors.Add(
			fmt.Sprintf("undeclared variable: %s", name),
			c.newLocation(span),
		)
	}
}

func (c *Compiler) setLocal(name string, valueNode ast.ExpressionNode, span *position.Span) {
	c.compileNode(valueNode)
	c.setLocalWithoutValue(name, span)
}

func (c *Compiler) localVariableAssignment(name string, operator *token.Token, right ast.ExpressionNode, span *position.Span) {
	switch operator.Type {
	case token.OR_OR_EQUAL:
		c.localVariableAccess(name, span)
		jump := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF)

		// if falsy
		c.emit(span.StartPos.Line, bytecode.POP)
		c.setLocal(name, right, span)

		// if truthy
		c.patchJump(jump, span)
	case token.AND_AND_EQUAL:
		c.localVariableAccess(name, span)
		jump := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)

		// if truthy
		c.emit(span.StartPos.Line, bytecode.POP)
		c.setLocal(name, right, span)

		// if falsy
		c.patchJump(jump, span)
	case token.QUESTION_QUESTION_EQUAL:
		c.localVariableAccess(name, span)
		nilJump := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF_NIL)
		nonNilJump := c.emitJump(span.StartPos.Line, bytecode.JUMP)

		// if nil
		c.patchJump(nilJump, span)
		c.emit(span.StartPos.Line, bytecode.POP)
		c.setLocal(name, right, span)

		// if not nil
		c.patchJump(nonNilJump, span)
	case token.PLUS_EQUAL:
		c.complexAssignment(name, right, bytecode.ADD, span)
	case token.MINUS_EQUAL:
		c.complexAssignment(name, right, bytecode.SUBTRACT, span)
	case token.STAR_EQUAL:
		c.complexAssignment(name, right, bytecode.MULTIPLY, span)
	case token.SLASH_EQUAL:
		c.complexAssignment(name, right, bytecode.DIVIDE, span)
	case token.STAR_STAR_EQUAL:
		c.complexAssignment(name, right, bytecode.EXPONENTIATE, span)
	case token.PERCENT_EQUAL:
		c.complexAssignment(name, right, bytecode.MODULO, span)
	case token.AND_EQUAL:
		c.complexAssignment(name, right, bytecode.BITWISE_AND, span)
	case token.OR_EQUAL:
		c.complexAssignment(name, right, bytecode.BITWISE_OR, span)
	case token.XOR_EQUAL:
		c.complexAssignment(name, right, bytecode.BITWISE_XOR, span)
	case token.LBITSHIFT_EQUAL:
		c.complexAssignment(name, right, bytecode.LBITSHIFT, span)
	case token.LTRIPLE_BITSHIFT_EQUAL:
		c.complexAssignment(name, right, bytecode.LOGIC_LBITSHIFT, span)
	case token.RBITSHIFT_EQUAL:
		c.complexAssignment(name, right, bytecode.RBITSHIFT, span)
	case token.RTRIPLE_BITSHIFT_EQUAL:
		c.complexAssignment(name, right, bytecode.LOGIC_RBITSHIFT, span)
	case token.EQUAL_OP:
		c.setLocal(name, right, span)
	case token.COLON_EQUAL:
		c.compileNode(right)
		local := c.defineLocal(name, span, false, true)
		if local == nil {
			return
		}
		c.emitSetLocal(span.StartPos.Line, local.index)
	default:
		c.Errors.Add(
			fmt.Sprintf("assignment using this operator has not been implemented: %s", operator.Type.String()),
			c.newLocation(span),
		)
		return
	}
}

func (c *Compiler) instanceVariableAccess(name string, span *position.Span) {
	switch c.mode {
	case topLevelMode:
		c.Errors.Add(
			"cannot read instance variables in the top level",
			c.newLocation(span),
		)
		return
	}

	c.emitGetInstanceVariable(value.ToSymbol(name), span)
}

func (c *Compiler) localVariableAccess(name string, span *position.Span) (*local, *upvalue, bool) {
	if local, ok := c.resolveLocal(name, span); ok {
		if !local.initialised {
			c.Errors.Add(
				fmt.Sprintf("cannot access an uninitialised local: %s", name),
				c.newLocation(span),
			)
			return nil, nil, false
		}

		c.emitGetLocal(span.StartPos.Line, local.index)
		return local, nil, true
	} else if upvalue, ok := c.resolveUpvalue(name, span); ok {
		local := upvalue.local
		if !local.initialised {
			c.Errors.Add(
				fmt.Sprintf("cannot access an uninitialised local: %s", name),
				c.newLocation(span),
			)
			return nil, nil, false
		}

		c.emitGetUpvalue(span.StartPos.Line, upvalue.index)
		return local, upvalue, true
	}

	c.Errors.Add(
		fmt.Sprintf("undeclared variable: %s", name),
		c.newLocation(span),
	)
	return nil, nil, false
}

// Resolve an upvalue from an outer context and get its index.
func (c *Compiler) resolveUpvalue(name string, span *position.Span) (*upvalue, bool) {
	parent := c.parent
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
		c.Errors.Add(
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

func (c *Compiler) docComment(node *ast.DocCommentNode) {
	c.emitValue(value.String(node.Comment), node.Span())
	c.compileNode(node.Expression)
	c.emit(node.Span().EndPos.Line, bytecode.DOC_COMMENT)
}

func (c *Compiler) modifierExpression(label string, node *ast.ModifierNode) {
	switch node.Modifier.Type {
	case token.IF:
		c.modifierIfExpression(false, node.Right, node.Left, nil, node.Span())
	case token.UNLESS:
		c.modifierIfExpression(true, node.Right, node.Left, nil, node.Span())
	case token.WHILE:
		c.modifierWhileExpression(label, node)
	case token.UNTIL:
		c.modifierUntilExpression(label, node)
	default:
		c.Errors.Add(
			fmt.Sprintf("illegal modifier: %s", node.Modifier.StringValue()),
			c.newLocation(node.Span()),
		)
	}
}

func (c *Compiler) modifierIfExpression(unless bool, condition, then, els ast.ExpressionNode, span *position.Span) {
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
	c.compileIf(
		jumpOp,
		condition,
		func() {
			c.compileNode(then)
		},
		elsFunc,
		span,
	)
}

func (c *Compiler) ifExpression(unless bool, condition ast.ExpressionNode, then, els []ast.StatementNode, span *position.Span) {
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

	c.compileIf(
		jumpOp,
		condition,
		func() {
			c.compileStatements(then, span)
		},
		elsFunc,
		span,
	)
}

func (c *Compiler) compileIf(jumpOp bytecode.OpCode, condition ast.ExpressionNode, then, els func(), span *position.Span) {
	if result := resolve(condition); result != nil {
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

	c.enterScope("", defaultScopeType)
	c.compileNode(condition)

	thenJumpOffset := c.emitJump(span.StartPos.Line, jumpOp)

	c.emit(span.StartPos.Line, bytecode.POP)

	then()
	c.leaveScope(span.StartPos.Line)

	elseJumpOffset := c.emitJump(span.StartPos.Line, bytecode.JUMP)

	c.patchJump(thenJumpOffset, span)
	c.emit(span.StartPos.Line, bytecode.POP)

	if els != nil {
		c.enterScope("", defaultScopeType)
		els()
		c.leaveScope(span.StartPos.Line)
	} else {
		c.emit(span.StartPos.Line, bytecode.NIL)
	}
	c.patchJump(elseJumpOffset, span)
}

func (c *Compiler) valueDeclaration(node *ast.ValueDeclarationNode) {
	initialised := node.Initialiser != nil

	switch node.Name.Type {
	case token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER:
		if initialised {
			c.compileNode(node.Initialiser)
		}
		local := c.defineLocal(node.Name.StringValue(), node.Span(), true, initialised)
		if local == nil {
			return
		}
		if initialised {
			c.emitSetLocal(node.Span().StartPos.Line, local.index)
		} else {
			c.emit(node.Span().StartPos.Line, bytecode.NIL)
		}
	default:
		c.Errors.Add(
			fmt.Sprintf("cannot compile a value declaration with: %s", node.Name.Type.String()),
			c.newLocation(node.Name.Span()),
		)
	}
}

func (c *Compiler) returnExpression(node *ast.ReturnExpressionNode) {
	span := node.Span()
	if node.Value != nil {
		c.emitReturn(span, node.Value)
	} else {
		c.emit(span.StartPos.Line, bytecode.NIL)
		c.emitReturn(span, nil)
	}
}

func (c *Compiler) receiverlessMethodCall(node *ast.ReceiverlessMethodCallNode) {
	for _, posArg := range node.PositionalArguments {
		c.compileNode(posArg)
	}

	var namedArgs []value.Symbol
namedArgNodeLoop:
	for _, namedArgVal := range node.NamedArguments {
		namedArg := namedArgVal.(*ast.NamedCallArgumentNode)
		namedArgName := value.ToSymbol(namedArg.Name)
		for _, argName := range namedArgs {
			if argName == namedArgName {
				c.Errors.Add(
					fmt.Sprintf("duplicated named argument in call: %s", argName.Inspect()),
					c.newLocation(namedArg.Span()),
				)
				continue namedArgNodeLoop
			}
		}
		namedArgs = append(namedArgs, namedArgName)
		c.compileNode(namedArg.Value)
	}

	name := value.ToSymbol(node.MethodName)
	argumentCount := len(node.PositionalArguments) + len(node.NamedArguments)
	callInfo := value.NewCallSiteInfo(name, argumentCount, namedArgs)
	c.emitCallFunction(callInfo, node.Span())
}

func (c *Compiler) nilSafeSubscriptExpression(node *ast.NilSafeSubscriptExpressionNode) {
	if c.resolveAndEmit(node) {
		return
	}

	c.compileIf(
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

func (c *Compiler) literalPattern(callInfo *value.CallSiteInfo, pattern ast.Node) {
	span := pattern.Span()
	c.emit(span.StartPos.Line, bytecode.DUP)
	c.compileNode(pattern)
	c.emitCallPattern(callInfo, span)
}

var (
	containsSymbol       = value.ToSymbol("contains")
	lengthSymbol         = value.ToSymbol("length")
	equalSymbol          = value.ToSymbol("==")
	notEqualSymbol       = value.ToSymbol("!=")
	laxEqualSymbol       = value.ToSymbol("=~")
	laxNotEqualSymbol    = value.ToSymbol("!~")
	strictEqualSymbol    = value.ToSymbol("===")
	strictNotEqualSymbol = value.ToSymbol("!==")
	lessSymbol           = value.ToSymbol("<")
	lessEqualSymbol      = value.ToSymbol("<=")
	greaterSymbol        = value.ToSymbol(">")
	greaterEqualSymbol   = value.ToSymbol(">=")
)

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
			value.NewCallSiteInfo(equalSymbol, 1, nil),
			pat,
		)
	case *ast.RangeLiteralNode:
		c.emit(span.StartPos.Line, bytecode.DUP)
		c.rangeLiteral(pat)
		c.emit(span.StartPos.Line, bytecode.SWAP)
		callInfo := value.NewCallSiteInfo(containsSymbol, 1, nil)
		c.emitCallMethod(callInfo, span)
	case *ast.PublicIdentifierNode:
		switch c.mode {
		case valuePatternDeclarationNode:
			c.defineLocal(pat.Value, span, true, false)
			c.setLocalWithoutValue(pat.Value, span)
			c.emit(span.StartPos.Line, bytecode.TRUE)
		default:
			c.defineLocalOverrideCurrentScope(pat.Value, span, false)
			c.setLocalWithoutValue(pat.Value, span)
			c.emit(span.StartPos.Line, bytecode.TRUE)
		}
	case *ast.PrivateIdentifierNode:
		switch c.mode {
		case valuePatternDeclarationNode:
			c.defineLocal(pat.Value, span, true, false)
			c.setLocalWithoutValue(pat.Value, span)
			c.emit(span.StartPos.Line, bytecode.TRUE)
		default:
			c.defineLocalOverrideCurrentScope(pat.Value, span, false)
			c.setLocalWithoutValue(pat.Value, span)
			c.emit(span.StartPos.Line, bytecode.TRUE)
		}
	case *ast.ObjectPatternNode:
		c.objectPattern(pat)
	case *ast.AsPatternNode:
		c.asPattern(pat)
	case *ast.UninterpolatedRegexLiteralNode, *ast.InterpolatedRegexLiteralNode:
		c.emit(span.StartPos.Line, bytecode.DUP)
		c.compileNode(pat)
		c.emit(span.StartPos.Line, bytecode.SWAP)
		callInfo := value.NewCallSiteInfo(matchesSymbol, 1, nil)
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
		c.Errors.Add(
			fmt.Sprintf("compilation of this pattern has not been implemented: %T", pattern),
			c.newLocation(span),
		)
	}
}

func (c *Compiler) unaryPattern(pat *ast.UnaryExpressionNode) {
	var methodName value.Symbol
	switch pat.Op.Type {
	case token.EQUAL_EQUAL:
		methodName = equalSymbol
	case token.NOT_EQUAL:
		methodName = notEqualSymbol
	case token.LAX_EQUAL:
		methodName = laxEqualSymbol
	case token.LAX_NOT_EQUAL:
		methodName = laxNotEqualSymbol
	case token.STRICT_EQUAL:
		methodName = strictEqualSymbol
	case token.STRICT_NOT_EQUAL:
		methodName = strictNotEqualSymbol
	case token.LESS:
		methodName = lessSymbol
	case token.LESS_EQUAL:
		methodName = lessEqualSymbol
	case token.GREATER:
		methodName = greaterSymbol
	case token.GREATER_EQUAL:
		methodName = greaterEqualSymbol
	default:
		c.literalPattern(
			value.NewCallSiteInfo(equalSymbol, 1, nil),
			pat,
		)
		return
	}

	c.literalPattern(
		value.NewCallSiteInfo(methodName, 1, nil),
		pat.Right,
	)
}

func (c *Compiler) binaryPattern(pat *ast.BinaryPatternNode) {
	span := pat.Span()
	var op bytecode.OpCode
	switch pat.Op.Type {
	case token.OR_OR:
		op = bytecode.JUMP_IF
	case token.AND_AND:
		op = bytecode.JUMP_UNLESS
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

	c.defineLocal(varName, span, true, false)
	c.setLocalWithoutValue(varName, span)
	c.pattern(node.Pattern)
}

func (c *Compiler) identifierObjectPatternAttribute(name string, span *position.Span) {
	c.emit(span.StartPos.Line, bytecode.DUP)
	callInfo := value.NewCallSiteInfo(value.ToSymbol(name), 0, nil)
	c.emitCallMethod(callInfo, span)

	identVar := c.defineLocal(name, span, true, true)
	c.emitSetLocal(span.StartPos.Line, identVar.index)
	c.emit(span.StartPos.Line, bytecode.POP)
}

func (c *Compiler) objectPattern(node *ast.ObjectPatternNode) {
	var jumpsToPatch []int
	c.enterPattern()

	span := node.Span()
	c.emit(node.Class.Span().StartPos.Line, bytecode.DUP)
	c.compileNode(node.Class)
	c.emit(node.Class.Span().StartPos.Line, bytecode.IS_A)

	jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)
	jumpsToPatch = append(jumpsToPatch, jmp)
	c.emit(span.StartPos.Line, bytecode.POP)

	for _, attr := range node.Attributes {
		span := attr.Span()
		switch e := attr.(type) {
		case *ast.SymbolKeyValuePatternNode:
			c.emit(span.StartPos.Line, bytecode.DUP)
			callInfo := value.NewCallSiteInfo(value.ToSymbol(e.Key), 0, nil)
			c.emitCallMethod(callInfo, span)

			c.pattern(e.Value)
			c.emit(span.StartPos.Line, bytecode.POP_SKIP_ONE)
			jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)
			jumpsToPatch = append(jumpsToPatch, jmp)
			c.emit(span.StartPos.Line, bytecode.POP)
		case *ast.PublicIdentifierNode:
			c.identifierObjectPatternAttribute(e.Value, span)
		case *ast.PrivateIdentifierNode:
			c.identifierObjectPatternAttribute(e.Value, span)
		default:
			c.Errors.Add(
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
		c.emitValue(value.ListMixin, span)
	case *ast.WordArrayTupleLiteralNode, *ast.SymbolArrayTupleLiteralNode, *ast.BinArrayTupleLiteralNode, *ast.HexArrayTupleLiteralNode:
		c.emitValue(value.TupleMixin, span)
	case *ast.WordHashSetLiteralNode, *ast.SymbolHashSetLiteralNode, *ast.BinHashSetLiteralNode, *ast.HexHashSetLiteralNode:
		c.emitValue(value.SetMixin, span)
	default:
		panic(fmt.Sprintf("invalid special collection pattern node: %#v", node))
	}
	c.emit(span.StartPos.Line, bytecode.IS_A)

	jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)
	c.emit(span.StartPos.Line, bytecode.POP)

	c.emit(span.StartPos.Line, bytecode.DUP)
	c.compileNode(node)
	c.emit(span.StartPos.Line, bytecode.LAX_EQUAL)

	// leave false on the stack from the falsy if that jumped here
	c.patchJump(jmp, span)
}

func (c *Compiler) identifierMapPatternElement(name string, span *position.Span) {
	c.emit(span.StartPos.Line, bytecode.DUP)
	c.emitValue(value.ToSymbol(name), span)
	c.emit(span.StartPos.Line, bytecode.SUBSCRIPT)
	identVar := c.defineLocal(name, span, false, true)
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
		c.emitValue(value.MapMixin, span)
	} else {
		c.emitValue(value.RecordMixin, span)
	}
	c.emit(span.StartPos.Line, bytecode.IS_A)

	jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)
	jumpsToPatch = append(jumpsToPatch, jmp)
	c.emit(span.StartPos.Line, bytecode.POP)

	for _, element := range elements {
		span := element.Span()
		switch e := element.(type) {
		case *ast.SymbolKeyValuePatternNode:
			c.emit(span.StartPos.Line, bytecode.DUP)
			c.emitValue(value.ToSymbol(e.Key), span)
			c.emit(span.StartPos.Line, bytecode.SUBSCRIPT)

			c.pattern(e.Value)
			c.emit(span.StartPos.Line, bytecode.POP_SKIP_ONE)
			jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)
			jumpsToPatch = append(jumpsToPatch, jmp)
			c.emit(span.StartPos.Line, bytecode.POP)
		case *ast.KeyValuePatternNode:
			c.emit(span.StartPos.Line, bytecode.DUP)
			c.compileNode(e.Key)
			c.emit(span.StartPos.Line, bytecode.SUBSCRIPT)

			c.pattern(e.Value)
			c.emit(span.StartPos.Line, bytecode.POP_SKIP_ONE)
			jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)
			jumpsToPatch = append(jumpsToPatch, jmp)
			c.emit(span.StartPos.Line, bytecode.POP)
		case *ast.PublicIdentifierNode:
			c.identifierMapPatternElement(e.Value, span)
		case *ast.PrivateIdentifierNode:
			c.identifierMapPatternElement(e.Value, span)
		default:
			c.Errors.Add(
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
				c.Errors.Add(
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
	c.emitValue(value.SetMixin, span)
	c.emit(span.StartPos.Line, bytecode.IS_A)

	jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)
	jumpsToPatch = append(jumpsToPatch, jmp)
	c.emit(span.StartPos.Line, bytecode.POP)

	c.emit(span.StartPos.Line, bytecode.DUP)
	callInfo := value.NewCallSiteInfo(lengthSymbol, 0, nil)
	c.emitCallMethod(callInfo, span)

	if !restElementIsPresent {
		c.emitValue(value.SmallInt(len(subPatternElements)), span)
		c.emit(span.StartPos.Line, bytecode.EQUAL)
	} else {
		c.emitValue(value.SmallInt(len(subPatternElements)), span)
		c.emit(span.StartPos.Line, bytecode.GREATER_EQUAL)
	}

	jmp = c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)
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
		callInfo := value.NewCallSiteInfo(containsSymbol, 1, nil)
		c.emitCallMethod(callInfo, span)

		jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)
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
				c.Errors.Add(
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
		restListVar = c.defineLocal(restVariableName, span, true, true)
		c.emitSetLocal(span.StartPos.Line, restListVar.index)
		c.emit(span.StartPos.Line, bytecode.POP)
	}
	c.enterPattern()

	c.emit(span.StartPos.Line, bytecode.DUP)
	if isList {
		c.emitValue(value.ListMixin, span)
	} else {
		c.emitValue(value.TupleMixin, span)
	}
	c.emit(span.StartPos.Line, bytecode.IS_A)

	jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)
	jumpsToPatch = append(jumpsToPatch, jmp)
	c.emit(span.StartPos.Line, bytecode.POP)

	c.emit(span.StartPos.Line, bytecode.DUP)
	callInfo := value.NewCallSiteInfo(lengthSymbol, 0, nil)
	c.emitCallMethod(callInfo, span)
	var lengthVar *local
	if elementBeforeRestCount != -1 {
		lengthVar = c.defineLocal(fmt.Sprintf("#!listPatternLength%d", c.patternNesting), span, true, true)
		c.emitSetLocal(span.StartPos.Line, lengthVar.index)
	}

	if elementBeforeRestCount == -1 {
		c.emitValue(value.SmallInt(len(elements)), span)
		c.emit(span.StartPos.Line, bytecode.EQUAL)
	} else {
		staticElementCount := elementBeforeRestCount + elementAfterRestCount
		c.emitValue(value.SmallInt(staticElementCount), span)
		c.emit(span.StartPos.Line, bytecode.GREATER_EQUAL)
	}

	jmp = c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)
	jumpsToPatch = append(jumpsToPatch, jmp)
	c.emit(span.StartPos.Line, bytecode.POP)

	elementsBeforeRest := elements
	if elementBeforeRestCount != -1 {
		elementsBeforeRest = elements[:elementBeforeRestCount]
	}
	for i, element := range elementsBeforeRest {
		span := element.Span()
		c.emit(span.StartPos.Line, bytecode.DUP)
		c.emitValue(value.SmallInt(i), element.Span())
		c.emit(span.StartPos.Line, bytecode.SUBSCRIPT)

		c.pattern(element)
		c.emit(span.StartPos.Line, bytecode.POP_SKIP_ONE)
		jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)
		jumpsToPatch = append(jumpsToPatch, jmp)
		c.emit(span.StartPos.Line, bytecode.POP)
	}

	if elementBeforeRestCount != -1 {
		iteratorVar := c.defineLocal(fmt.Sprintf("#!listPatternIterator%d", c.patternNesting), span, true, true)

		if restVariableName != "" {
			// adjust the length variable
			// length -= element_after_rest_count
			if elementAfterRestCount != 0 {
				c.emitGetLocal(span.StartPos.Line, lengthVar.index)
				c.emitValue(value.SmallInt(elementAfterRestCount), span)
				c.emit(span.StartPos.Line, bytecode.SUBTRACT)
				c.emitSetLocal(span.StartPos.Line, lengthVar.index)
				c.emit(span.StartPos.Line, bytecode.POP)
			}

			// create the iterator variable
			// i := element_before_rest_count
			c.emitValue(value.SmallInt(elementBeforeRestCount), span)
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
			c.emit(span.StartPos.Line, bytecode.POP)
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
			c.emit(span.StartPos.Line, bytecode.POP)
		} else {
			// create the iterator variable
			// i := length - element_after_rest_count
			if elementAfterRestCount == 0 {
				c.emitGetLocal(span.StartPos.Line, lengthVar.index)
				c.emitSetLocal(span.StartPos.Line, iteratorVar.index)
				c.emit(span.StartPos.Line, bytecode.POP)
			} else {
				c.emitGetLocal(span.StartPos.Line, lengthVar.index)
				c.emitValue(value.SmallInt(elementAfterRestCount), span)
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
			jmp := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)
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

func (c *Compiler) switchExpression(node *ast.SwitchExpressionNode) {
	span := node.Span()

	c.enterScope("", defaultScopeType)
	c.compileNode(node.Value)

	var jumpToEndOffsets []int

	for _, caseNode := range node.Cases {
		c.enterScope("", defaultScopeType)

		caseSpan := caseNode.Span()
		c.pattern(caseNode.Pattern)

		jumpOverBodyOffset := c.emitJump(caseSpan.StartPos.Line, bytecode.JUMP_UNLESS)

		c.emit(caseSpan.StartPos.Line, bytecode.POP_N, 2)

		c.compileStatements(caseNode.Body, caseSpan)

		c.leaveScopeWithoutMutating(caseSpan.EndPos.Line)

		jumpToEndOffset := c.emitJump(caseSpan.EndPos.Line, bytecode.JUMP)
		jumpToEndOffsets = append(jumpToEndOffsets, jumpToEndOffset)

		c.patchJump(jumpOverBodyOffset, caseSpan)
		c.leaveScope(caseSpan.EndPos.Line)
		c.emit(caseSpan.EndPos.Line, bytecode.POP)
	}

	c.emit(span.StartPos.Line, bytecode.POP)
	c.compileStatements(node.ElseBody, span)

	for _, offset := range jumpToEndOffsets {
		c.patchJump(offset, span)
	}

	c.leaveScope(span.EndPos.Line)
}

func (c *Compiler) subscriptExpression(node *ast.SubscriptExpressionNode) {
	if c.resolveAndEmit(node) {
		return
	}

	c.compileNode(node.Receiver)
	c.compileNode(node.Key)
	c.emit(node.Span().EndPos.Line, bytecode.SUBSCRIPT)
}

func (c *Compiler) attributeAccess(node *ast.AttributeAccessNode) {
	c.compileNode(node.Receiver)

	name := value.ToSymbol(node.AttributeName)
	callInfo := value.NewCallSiteInfo(name, 0, nil)
	if node.AttributeName == "call" {
		c.emitCall(callInfo, node.Span())
	} else {
		c.emitCallMethod(callInfo, node.Span())
	}
}

func (c *Compiler) constructorCall(node *ast.ConstructorCallNode) {
	c.compileNode(node.Class)
	for _, posArg := range node.PositionalArguments {
		c.compileNode(posArg)
	}

	var namedArgs []value.Symbol
namedArgNodeLoop:
	for _, namedArgVal := range node.NamedArguments {
		namedArg := namedArgVal.(*ast.NamedCallArgumentNode)
		namedArgName := value.ToSymbol(namedArg.Name)
		for _, argName := range namedArgs {
			if argName == namedArgName {
				c.Errors.Add(
					fmt.Sprintf("duplicated named argument in call: %s", argName.Inspect()),
					c.newLocation(namedArg.Span()),
				)
				continue namedArgNodeLoop
			}
		}
		namedArgs = append(namedArgs, namedArgName)
		c.compileNode(namedArg.Value)
	}

	name := value.ToSymbol("#init")
	argumentCount := len(node.PositionalArguments) + len(node.NamedArguments)
	callInfo := value.NewCallSiteInfo(name, argumentCount, namedArgs)
	c.emitInstantiate(callInfo, node.Span())
}

func (c *Compiler) methodCall(node *ast.MethodCallNode) {
	c.compileNode(node.Receiver)

	if node.NilSafe {
		nilJump := c.emitJump(node.Span().StartPos.Line, bytecode.JUMP_IF_NIL)

		// if not nil
		// call the method
		c.innerMethodCall(node)

		// if nil
		// leave nil on the stack
		c.patchJump(nilJump, node.Span())
		return
	}

	c.innerMethodCall(node)
}

func (c *Compiler) innerMethodCall(node *ast.MethodCallNode) {
	for _, posArg := range node.PositionalArguments {
		c.compileNode(posArg)
	}

	var namedArgs []value.Symbol
namedArgNodeLoop:
	for _, namedArgVal := range node.NamedArguments {
		namedArg := namedArgVal.(*ast.NamedCallArgumentNode)
		namedArgName := value.ToSymbol(namedArg.Name)
		for _, argName := range namedArgs {
			if argName == namedArgName {
				c.Errors.Add(
					fmt.Sprintf("duplicated named argument in call: %s", argName.Inspect()),
					c.newLocation(namedArg.Span()),
				)
				continue namedArgNodeLoop
			}
		}
		namedArgs = append(namedArgs, namedArgName)
		c.compileNode(namedArg.Value)
	}

	name := value.ToSymbol(node.MethodName)
	argumentCount := len(node.PositionalArguments) + len(node.NamedArguments)
	callInfo := value.NewCallSiteInfo(name, argumentCount, namedArgs)
	if node.MethodName == "call" {
		c.emitCall(callInfo, node.Span())
	} else {
		c.emitCallMethod(callInfo, node.Span())
	}
}

func (c *Compiler) call(node *ast.CallNode) {
	c.compileNode(node.Receiver)

	if node.NilSafe {
		nilJump := c.emitJump(node.Span().StartPos.Line, bytecode.JUMP_IF_NIL)

		// if not nil
		// call the method
		c.innerCall(node)

		// if nil
		// leave nil on the stack
		c.patchJump(nilJump, node.Span())
		return
	}

	c.innerCall(node)
}

func (c *Compiler) innerCall(node *ast.CallNode) {
	for _, posArg := range node.PositionalArguments {
		c.compileNode(posArg)
	}

	var namedArgs []value.Symbol
namedArgNodeLoop:
	for _, namedArgVal := range node.NamedArguments {
		namedArg := namedArgVal.(*ast.NamedCallArgumentNode)
		namedArgName := value.ToSymbol(namedArg.Name)
		for _, argName := range namedArgs {
			if argName == namedArgName {
				c.Errors.Add(
					fmt.Sprintf("duplicated named argument in call: %s", argName.Inspect()),
					c.newLocation(namedArg.Span()),
				)
				continue namedArgNodeLoop
			}
		}
		namedArgs = append(namedArgs, namedArgName)
		c.compileNode(namedArg.Value)
	}

	name := value.ToSymbol("call")
	argumentCount := len(node.PositionalArguments) + len(node.NamedArguments)
	callInfo := value.NewCallSiteInfo(name, argumentCount, namedArgs)
	c.emitCall(callInfo, node.Span())
}

func (c *Compiler) singletonBlock(node *ast.SingletonBlockExpressionNode) {
	span := node.Span()
	switch c.mode {
	case classMode, mixinMode, moduleMode:
	case topLevelMode:
		c.Errors.Add(
			"cannot open a singleton class in the top level",
			c.newLocation(span),
		)
		return
	case functionMode, setterMethodMode, initMethodMode:
		c.Errors.Add(
			"cannot open a singleton class in a method",
			c.newLocation(span),
		)
		return
	default:
		c.Errors.Add(
			"cannot open a singleton class in this context",
			c.newLocation(span),
		)
		return
	}

	if len(node.Body) == 0 {
		c.emit(span.StartPos.Line, bytecode.UNDEFINED)
	} else {
		singletonCompiler := new("<singleton_class>", classMode, c.newLocation(span))
		singletonCompiler.Errors = c.Errors
		singletonCompiler.compileModule(node)
		c.Errors = singletonCompiler.Errors

		result := singletonCompiler.Bytecode
		c.emitValue(result, span)
	}
	c.emit(span.StartPos.Line, bytecode.SELF)
	c.emit(span.StartPos.Line, bytecode.DEF_SINGLETON)
}

func (c *Compiler) methodDefinition(node *ast.MethodDefinitionNode) {
	switch c.mode {
	case functionMode, setterMethodMode, initMethodMode:
		c.Errors.Add(
			fmt.Sprintf("methods cannot be nested: %s", node.Name),
			c.newLocation(node.Span()),
		)
		return
	}

	var mode mode
	if node.IsSetter() {
		mode = setterMethodMode
	} else {
		mode = functionMode
	}

	methodCompiler := new(node.Name, mode, c.newLocation(node.Span()))
	methodCompiler.Errors = c.Errors
	methodCompiler.compileMethod(node.Span(), node.Parameters, node.Body)
	c.Errors = methodCompiler.Errors

	result := methodCompiler.Bytecode
	c.emitValue(result, node.Span())

	c.emitValue(value.ToSymbol(node.Name), node.Span())

	c.emit(node.Span().StartPos.Line, bytecode.DEF_METHOD)
}

func (c *Compiler) functionLiteral(node *ast.FunctionLiteralNode) {
	functionCompiler := new("<function>", functionMode, c.newLocation(node.Span()))
	functionCompiler.parent = c
	functionCompiler.Errors = c.Errors
	functionCompiler.compileFunction(node.Span(), node.Parameters, node.Body)
	c.Errors = functionCompiler.Errors

	result := functionCompiler.Bytecode
	c.emitValue(result, node.Span())

	c.emit(node.Span().StartPos.Line, bytecode.CLOSURE)

	for _, upvalue := range functionCompiler.upvalues {
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

func (c *Compiler) initDefinition(node *ast.InitDefinitionNode) {
	switch c.mode {
	case functionMode, setterMethodMode, initMethodMode:
		c.Errors.Add(
			"methods cannot be nested: #init",
			c.newLocation(node.Span()),
		)
		return
	case topLevelMode:
		c.Errors.Add(
			"init cannot be defined in the top level",
			c.newLocation(node.Span()),
		)
		return
	case moduleMode:
		c.Errors.Add(
			"modules cannot have initializers",
			c.newLocation(node.Span()),
		)
		return
	}

	methodCompiler := new("#init", initMethodMode, c.newLocation(node.Span()))
	methodCompiler.Errors = c.Errors
	methodCompiler.compileMethod(node.Span(), node.Parameters, node.Body)
	c.Errors = methodCompiler.Errors

	result := methodCompiler.Bytecode
	c.emitValue(result, node.Span())

	c.emitValue(value.ToSymbol("#init"), node.Span())

	c.emit(node.Span().StartPos.Line, bytecode.DEF_METHOD)
}

func (c *Compiler) extendExpression(node *ast.ExtendExpressionNode) {
	switch c.mode {
	case classMode, mixinMode, moduleMode:
	case topLevelMode:
		c.Errors.Add(
			"cannot extend mixins in the top level",
			c.newLocation(node.Span()),
		)
		return
	case functionMode:
		c.Errors.Add(
			"cannot extend mixins in a method",
			c.newLocation(node.Span()),
		)
		return
	default:
		c.Errors.Add(
			"cannot extend mixins in this context",
			c.newLocation(node.Span()),
		)
		return
	}

	span := node.Span()
	for _, constant := range node.Constants {
		c.compileNode(constant)
		c.emit(span.StartPos.Line, bytecode.SELF)
		c.emit(span.StartPos.Line, bytecode.GET_SINGLETON)
		c.emit(span.StartPos.Line, bytecode.INCLUDE)
	}

	c.emit(span.EndPos.Line, bytecode.NIL)
}

func (c *Compiler) includeExpression(node *ast.IncludeExpressionNode) {
	switch c.mode {
	case classMode, mixinMode:
	case topLevelMode:
		c.Errors.Add(
			"cannot include mixins in the top level",
			c.newLocation(node.Span()),
		)
		return
	case moduleMode:
		c.Errors.Add(
			"cannot include mixins in a module",
			c.newLocation(node.Span()),
		)
		return
	case functionMode:
		c.Errors.Add(
			"cannot include mixins in a method",
			c.newLocation(node.Span()),
		)
		return
	default:
		c.Errors.Add(
			"cannot include mixins in this context",
			c.newLocation(node.Span()),
		)
		return
	}

	span := node.Span()
	for _, constant := range node.Constants {
		c.compileNode(constant)
		c.emit(span.StartPos.Line, bytecode.SELF)
		c.emit(span.StartPos.Line, bytecode.INCLUDE)
	}

	c.emit(span.EndPos.Line, bytecode.NIL)
}

func (c *Compiler) mixinDeclaration(node *ast.MixinDeclarationNode) {
	switch c.mode {
	case functionMode:
		if node.Constant != nil {
			c.Errors.Add(
				fmt.Sprintf("cannot define named mixins inside of a method: %s", c.Bytecode.Name().ToString()),
				c.newLocation(node.Span()),
			)
			return
		}
	}

	if len(node.Body) == 0 {
		c.emit(node.Span().StartPos.Line, bytecode.UNDEFINED)
	} else {
		mixinCompiler := new("<mixin>", mixinMode, c.newLocation(node.Span()))
		mixinCompiler.Errors = c.Errors
		mixinCompiler.compileModule(node)
		c.Errors = mixinCompiler.Errors

		result := mixinCompiler.Bytecode
		c.emitValue(result, node.Span())
	}

	switch constant := node.Constant.(type) {
	case *ast.ConstantLookupNode:
		if constant.Left != nil {
			c.compileNode(constant.Left)
		} else {
			c.emit(constant.Span().StartPos.Line, bytecode.ROOT)
		}
		switch r := constant.Right.(type) {
		case *ast.PublicConstantNode:
			c.emitValue(value.ToSymbol(r.Value), r.Span())
		case *ast.PrivateConstantNode:
			c.emitValue(value.ToSymbol(r.Value), r.Span())
		default:
			c.Errors.Add(
				fmt.Sprintf("incorrect right side of constant lookup: %T", constant.Right),
				c.newLocation(constant.Right.Span()),
			)
			return
		}
	case *ast.PublicConstantNode:
		c.emit(constant.Span().StartPos.Line, bytecode.CONSTANT_CONTAINER)
		c.emitValue(value.ToSymbol(constant.Value), constant.Span())
	case *ast.PrivateConstantNode:
		c.emit(constant.Span().StartPos.Line, bytecode.CONSTANT_CONTAINER)
		c.emitValue(value.ToSymbol(constant.Value), constant.Span())
	case nil:
		c.emit(node.Span().StartPos.Line, bytecode.DEF_ANON_MIXIN)
		return
	default:
		c.Errors.Add(
			fmt.Sprintf("incorrect mixin name: %T", constant),
			c.newLocation(constant.Span()),
		)
		return
	}

	c.emit(node.Span().StartPos.Line, bytecode.DEF_MIXIN)
}

func (c *Compiler) moduleDeclaration(node *ast.ModuleDeclarationNode) {
	switch c.mode {
	case functionMode:
		if node.Constant != nil {
			c.Errors.Add(
				fmt.Sprintf("cannot define named modules inside of a method: %s", c.Bytecode.Name().ToString()),
				c.newLocation(node.Span()),
			)
			return
		}
	}

	if len(node.Body) == 0 {
		c.emit(node.Span().StartPos.Line, bytecode.UNDEFINED)
	} else {
		modCompiler := new("<module>", moduleMode, c.newLocation(node.Span()))
		modCompiler.Errors = c.Errors
		modCompiler.compileModule(node)
		c.Errors = modCompiler.Errors

		result := modCompiler.Bytecode
		c.emitValue(result, node.Span())
	}

	switch constant := node.Constant.(type) {
	case *ast.ConstantLookupNode:
		if constant.Left != nil {
			c.compileNode(constant.Left)
		} else {
			c.emit(constant.Span().StartPos.Line, bytecode.ROOT)
		}
		switch r := constant.Right.(type) {
		case *ast.PublicConstantNode:
			c.emitValue(value.ToSymbol(r.Value), r.Span())
		case *ast.PrivateConstantNode:
			c.emitValue(value.ToSymbol(r.Value), r.Span())
		default:
			c.Errors.Add(
				fmt.Sprintf("incorrect right side of constant lookup: %T", constant.Right),
				c.newLocation(constant.Right.Span()),
			)
			return
		}
	case *ast.PublicConstantNode:
		c.emit(constant.Span().StartPos.Line, bytecode.CONSTANT_CONTAINER)
		c.emitValue(value.ToSymbol(constant.Value), constant.Span())
	case *ast.PrivateConstantNode:
		c.emit(constant.Span().StartPos.Line, bytecode.CONSTANT_CONTAINER)
		c.emitValue(value.ToSymbol(constant.Value), constant.Span())
	case nil:
		c.emit(node.Span().StartPos.Line, bytecode.DEF_ANON_MODULE)
		return
	default:
		c.Errors.Add(
			fmt.Sprintf("incorrect module name: %T", constant),
			c.newLocation(constant.Span()),
		)
		return
	}

	c.emit(node.Span().StartPos.Line, bytecode.DEF_MODULE)
}

func (c *Compiler) getterDeclaration(node *ast.GetterDeclarationNode) {
	switch c.mode {
	case functionMode:
		c.Errors.Add(
			"cannot define getters in this context",
			c.newLocation(node.Span()),
		)
		return
	}

	for _, entry := range node.Entries {
		getterEntry := entry.(*ast.AttributeParameterNode)
		c.emitValue(value.ToSymbol(getterEntry.Name), node.Span())
		c.emit(node.Span().StartPos.Line, bytecode.DEF_GETTER)
	}
	c.emit(node.Span().EndPos.Line, bytecode.NIL)
}

func (c *Compiler) setterDeclaration(node *ast.SetterDeclarationNode) {
	switch c.mode {
	case functionMode:
		c.Errors.Add(
			"cannot define setters in this context",
			c.newLocation(node.Span()),
		)
		return
	}

	for _, entry := range node.Entries {
		getterEntry := entry.(*ast.AttributeParameterNode)
		c.emitValue(value.ToSymbol(getterEntry.Name), node.Span())
		c.emit(node.Span().StartPos.Line, bytecode.DEF_SETTER)
	}
	c.emit(node.Span().EndPos.Line, bytecode.NIL)
}

func (c *Compiler) accessorDeclaration(node *ast.AccessorDeclarationNode) {
	switch c.mode {
	case functionMode:
		c.Errors.Add(
			"cannot define accessors in this context",
			c.newLocation(node.Span()),
		)
		return
	}

	for _, entry := range node.Entries {
		getterEntry := entry.(*ast.AttributeParameterNode)
		c.emitValue(value.ToSymbol(getterEntry.Name), node.Span())
		c.emit(node.Span().StartPos.Line, bytecode.DEF_GETTER)

		c.emitValue(value.ToSymbol(getterEntry.Name), node.Span())
		c.emit(node.Span().StartPos.Line, bytecode.DEF_SETTER)
	}
	c.emit(node.Span().EndPos.Line, bytecode.NIL)
}

func (c *Compiler) aliasDeclaration(node *ast.AliasDeclarationNode) {
	switch c.mode {
	case functionMode:
		c.Errors.Add(
			"cannot define aliases in this context",
			c.newLocation(node.Span()),
		)
		return
	}

	for _, aliasEntry := range node.Entries {
		c.emitValue(value.ToSymbol(aliasEntry.OldName), node.Span())
		c.emitValue(value.ToSymbol(aliasEntry.NewName), node.Span())
		c.emit(node.Span().StartPos.Line, bytecode.DEF_ALIAS)
	}

	c.emit(node.Span().EndPos.Line, bytecode.NIL)
}

func (c *Compiler) classDeclaration(node *ast.ClassDeclarationNode) {
	switch c.mode {
	case functionMode:
		if node.Constant != nil {
			c.Errors.Add(
				fmt.Sprintf("cannot define named classes inside of a method: %s", c.Bytecode.Name().ToString()),
				c.newLocation(node.Span()),
			)
			return
		}
	}

	if len(node.Body) == 0 {
		c.emit(node.Span().StartPos.Line, bytecode.UNDEFINED)
	} else {
		modCompiler := new("<class>", classMode, c.newLocation(node.Span()))
		modCompiler.Errors = c.Errors
		modCompiler.compileModule(node)
		c.Errors = modCompiler.Errors

		result := modCompiler.Bytecode
		c.emitValue(result, node.Span())
	}

	switch constant := node.Constant.(type) {
	case *ast.ConstantLookupNode:
		if constant.Left != nil {
			c.compileNode(constant.Left)
		} else {
			c.emit(constant.Span().StartPos.Line, bytecode.ROOT)
		}
		switch r := constant.Right.(type) {
		case *ast.PublicConstantNode:
			c.emitValue(value.ToSymbol(r.Value), r.Span())
		case *ast.PrivateConstantNode:
			c.emitValue(value.ToSymbol(r.Value), r.Span())
		default:
			c.Errors.Add(
				fmt.Sprintf("incorrect right side of constant lookup: %T", constant.Right),
				c.newLocation(constant.Right.Span()),
			)
			return
		}
	case *ast.PublicConstantNode:
		c.emit(constant.Span().StartPos.Line, bytecode.CONSTANT_CONTAINER)
		c.emitValue(value.ToSymbol(constant.Value), constant.Span())
	case *ast.PrivateConstantNode:
		c.emit(constant.Span().StartPos.Line, bytecode.CONSTANT_CONTAINER)
		c.emitValue(value.ToSymbol(constant.Value), constant.Span())
	case nil:
		c.compileClassSuperclass(node)
		c.emit(node.Span().StartPos.Line, bytecode.DEF_ANON_CLASS)
		return
	default:
		c.Errors.Add(
			fmt.Sprintf("incorrect class name: %T", constant),
			c.newLocation(constant.Span()),
		)
		return
	}
	c.compileClassSuperclass(node)

	var flags bitfield.BitFlag8
	if node.Abstract {
		flags |= value.CLASS_ABSTRACT_FLAG
	}
	if node.Sealed {
		flags |= value.CLASS_SEALED_FLAG
	}

	c.emit(node.Span().StartPos.Line, bytecode.DEF_CLASS, byte(flags))
}

func (c *Compiler) compileClassSuperclass(node *ast.ClassDeclarationNode) {
	if node.Superclass != nil {
		c.compileNode(node.Superclass)
	} else {
		c.emit(node.Span().StartPos.Line, bytecode.UNDEFINED)
	}
}

func (c *Compiler) valuePatternDeclaration(node *ast.ValuePatternDeclarationNode) {
	previousMode := c.mode
	c.mode = valuePatternDeclarationNode
	defer func() {
		c.mode = previousMode
	}()

	span := node.Span()
	c.compileNode(node.Initialiser)
	c.pattern(node.Pattern)

	jumpOverErrorOffset := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF)
	c.emit(span.EndPos.Line, bytecode.POP)

	c.emitValue(
		value.NewError(
			value.PatternNotMatchedErrorClass,
			"assigned value does not match the pattern defined in value declaration",
		),
		span,
	)
	c.emit(span.EndPos.Line, bytecode.THROW)

	c.patchJump(jumpOverErrorOffset, span)
	c.emit(span.EndPos.Line, bytecode.POP)
}

func (c *Compiler) variablePatternDeclaration(node *ast.VariablePatternDeclarationNode) {
	span := node.Span()
	c.compileNode(node.Initialiser)
	c.pattern(node.Pattern)

	jumpOverErrorOffset := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF)
	c.emit(span.EndPos.Line, bytecode.POP)

	c.emitValue(
		value.NewError(
			value.PatternNotMatchedErrorClass,
			"assigned value does not match the pattern defined in variable declaration",
		),
		span,
	)
	c.emit(span.EndPos.Line, bytecode.THROW)

	c.patchJump(jumpOverErrorOffset, span)
	c.emit(span.EndPos.Line, bytecode.POP)
}

func (c *Compiler) variableDeclaration(node *ast.VariableDeclarationNode) {
	initialised := node.Initialiser != nil

	switch node.Name.Type {
	case token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER:
		if initialised {
			c.compileNode(node.Initialiser)
		}
		local := c.defineLocal(node.Name.StringValue(), node.Span(), false, initialised)
		if local == nil {
			return
		}
		if initialised {
			c.emitSetLocal(node.Span().StartPos.Line, local.index)
		} else {
			c.emit(node.Span().StartPos.Line, bytecode.NIL)
		}
	case token.INSTANCE_VARIABLE:
		switch c.mode {
		case classMode, mixinMode, moduleMode:
		default:
			c.Errors.Add(
				"instance variables can only be declared in class, module, mixin bodies",
				c.newLocation(node.Span()),
			)
		}
		c.emit(node.Span().StartPos.Line, bytecode.NIL)
	default:
		c.Errors.Add(
			fmt.Sprintf("cannot compile a variable declaration with: %s", node.Name.Type.String()),
			c.newLocation(node.Name.Span()),
		)
	}
}

// Compile each element of a collection of statements.
func (c *Compiler) compileStatements(collection []ast.StatementNode, span *position.Span) {
	if !c.compileStatementsOk(collection, span) {
		c.emit(span.EndPos.Line, bytecode.NIL)
	}
}

// Same as [compileStatements] but returns false when no instructions were emitted instead
// emitting a `nil` value.
func (c *Compiler) compileStatementsOk(collection []ast.StatementNode, span *position.Span) bool {
	var nonEmptyStatements []ast.StatementNode
	for _, s := range collection {
		if _, ok := s.(*ast.EmptyStatementNode); ok {
			continue
		}
		nonEmptyStatements = append(nonEmptyStatements, s)
	}

	for i, s := range nonEmptyStatements {
		isLast := i == len(nonEmptyStatements)-1
		if !isLast && s.IsStatic() {
			continue
		}
		c.compileNode(s)
		if !isLast {
			c.emit(s.Span().EndPos.Line, bytecode.POP)
		}
	}

	return len(nonEmptyStatements) != 0
}

func (c *Compiler) rangeLiteral(node *ast.RangeLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	span := node.Span()

	if node.From == nil {
		c.compileNode(node.To)

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
	if node.To == nil {
		c.compileNode(node.From)

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

	c.compileNode(node.From)
	c.compileNode(node.To)
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

func (c *Compiler) hashSetLiteral(node *ast.HashSetLiteralNode) {
	span := node.Span()
	if c.resolveAndEmit(node) {
		return
	}

	baseSet := value.NewHashSet(len(node.Elements))
	firstDynamicIndex := -1

	for i, elementNode := range node.Elements {
		element := resolve(elementNode)
		if element == nil || value.IsMutableCollection(element) {
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
		c.emitLoadValue(baseSet, span)
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
					c.Errors.Add(
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

				c.compileIf(
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
				c.compileIf(
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
					e.Parameter,
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

func (c *Compiler) hashMapLiteral(node *ast.HashMapLiteralNode) {
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
			if val == nil || value.IsMutableCollection(val) {
				break elementSwitch
			}

			vm.HashMapSetWithMaxLoad(nil, baseMap, key, val, 1)
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
		c.emitLoadValue(baseMap, span)
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
					c.Errors.Add(
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
				c.emitValue(value.ToSymbol(element.Key), element.Span())
				c.compileNode(element.Value)
			case *ast.PublicIdentifierNode:
				c.emitValue(value.ToSymbol(element.Value), element.Span())
				c.localVariableAccess(element.Value, element.Span())
			case *ast.PrivateIdentifierNode:
				c.emitValue(value.ToSymbol(element.Value), element.Span())
				c.localVariableAccess(element.Value, element.Span())
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
				c.emitValue(value.ToSymbol(e.Key), e.Span())
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

				c.compileIf(
					jumpOp,
					e.Right,
					func() {
						switch then := e.Left.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNode(then.Key)
							c.compileNode(then.Value)
							c.emit(then.Span().StartPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(then.Key), then.Span())
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
				c.compileIf(
					bytecode.JUMP_UNLESS,
					e.Condition,
					func() {
						switch then := e.ThenExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNode(then.Key)
							c.compileNode(then.Value)
							c.emit(then.Span().StartPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(then.Key), then.Span())
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
							c.emitValue(value.ToSymbol(els.Key), els.Span())
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
					e.Parameter,
					e.InExpression,
					func() {
						switch then := e.ThenExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNode(then.Key)
							c.compileNode(then.Value)
							c.emit(then.Span().EndPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(then.Key), then.Span())
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

func (c *Compiler) hashRecordLiteral(node *ast.HashRecordLiteralNode) {
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
			if val == nil || value.IsMutableCollection(val) {
				break elementSwitch
			}

			vm.HashRecordSetWithMaxLoad(nil, baseRecord, key, val, 1)
			continue elementLoop
		}

		firstDynamicIndex = i
		break elementLoop
	}

	if baseRecord.Length() == 0 {
		c.emit(span.StartPos.Line, bytecode.UNDEFINED)
	} else {
		c.emitLoadValue(baseRecord, span)
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
				c.emitValue(value.ToSymbol(element.Key), element.Span())
				c.compileNode(element.Value)
			case *ast.PublicIdentifierNode:
				c.emitValue(value.ToSymbol(element.Value), element.Span())
				c.localVariableAccess(element.Value, element.Span())
			case *ast.PrivateIdentifierNode:
				c.emitValue(value.ToSymbol(element.Value), element.Span())
				c.localVariableAccess(element.Value, element.Span())
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
				c.emitValue(value.ToSymbol(e.Key), e.Span())
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

				c.compileIf(
					jumpOp,
					e.Right,
					func() {
						switch then := e.Left.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNode(then.Key)
							c.compileNode(then.Value)
							c.emit(then.Span().StartPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(then.Key), then.Span())
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
				c.compileIf(
					bytecode.JUMP_UNLESS,
					e.Condition,
					func() {
						switch then := e.ThenExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNode(then.Key)
							c.compileNode(then.Value)
							c.emit(then.Span().StartPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(then.Key), then.Span())
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
							c.emitValue(value.ToSymbol(els.Key), els.Span())
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
					e.Parameter,
					e.InExpression,
					func() {
						switch then := e.ThenExpression.(type) {
						case *ast.KeyValueExpressionNode:
							c.compileNode(then.Key)
							c.compileNode(then.Value)
							c.emit(then.Span().EndPos.Line, bytecode.MAP_SET)
						case *ast.SymbolKeyValueExpressionNode:
							c.emitValue(value.ToSymbol(then.Key), then.Span())
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

func (c *Compiler) arrayListLiteral(node *ast.ArrayListLiteralNode) {
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
		if element == nil || value.IsMutableCollection(element) {
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
		c.emitLoadValue(&baseList, span)
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
					c.Errors.Add(
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

				c.compileIf(
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
				c.compileIf(
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
					e.Parameter,
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

func (c *Compiler) arrayTupleLiteral(node *ast.ArrayTupleLiteralNode) {
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
		if element == nil {
			firstDynamicIndex = i
			break
		}

		baseArrayTuple = append(baseArrayTuple, element)
	}

	if len(baseArrayTuple) == 0 {
		c.emit(span.StartPos.Line, bytecode.UNDEFINED)
	} else {
		c.emitLoadValue(&baseArrayTuple, span)
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

				c.compileIf(
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
				c.compileIf(
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
					e.Parameter,
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

func (c *Compiler) wordArrayTupleLiteral(node *ast.WordArrayTupleLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	c.Errors.Add("invalid word arrayTuple literal", c.newLocation(node.Span()))
}

func (c *Compiler) binArrayTupleLiteral(node *ast.BinArrayTupleLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	c.Errors.Add("invalid binary arrayTuple literal", c.newLocation(node.Span()))
}

func (c *Compiler) symbolArrayTupleLiteral(node *ast.SymbolArrayTupleLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	c.Errors.Add("invalid symbol arrayTuple literal", c.newLocation(node.Span()))
}

func (c *Compiler) hexArrayTupleLiteral(node *ast.HexArrayTupleLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	c.Errors.Add("invalid hex arrayTuple literal", c.newLocation(node.Span()))
}

func (c *Compiler) wordArrayListLiteral(node *ast.WordArrayListLiteralNode) {
	list := resolve(node)
	span := node.Span()
	if list == nil {
		c.Errors.Add("invalid word arrayList literal", c.newLocation(span))
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

func (c *Compiler) binArrayListLiteral(node *ast.BinArrayListLiteralNode) {
	list := resolve(node)
	span := node.Span()
	if list == nil {
		c.Errors.Add("invalid bin arrayList literal", c.newLocation(span))
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

func (c *Compiler) symbolArrayListLiteral(node *ast.SymbolArrayListLiteralNode) {
	list := resolve(node)
	span := node.Span()
	if list == nil {
		c.Errors.Add("invalid symbol arrayList literal", c.newLocation(span))
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

func (c *Compiler) hexArrayListLiteral(node *ast.HexArrayListLiteralNode) {
	list := resolve(node)
	span := node.Span()
	if list == nil {
		c.Errors.Add("invalid hex arrayList literal", c.newLocation(span))
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

func (c *Compiler) wordHashSetLiteral(node *ast.WordHashSetLiteralNode) {
	list := resolve(node)
	span := node.Span()
	if list == nil {
		c.Errors.Add("invalid word hashSet literal", c.newLocation(span))
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

func (c *Compiler) binHashSetLiteral(node *ast.BinHashSetLiteralNode) {
	list := resolve(node)
	span := node.Span()
	if list == nil {
		c.Errors.Add("invalid bin hashSet literal", c.newLocation(span))
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

func (c *Compiler) symbolHashSetLiteral(node *ast.SymbolHashSetLiteralNode) {
	list := resolve(node)
	span := node.Span()
	if list == nil {
		c.Errors.Add("invalid symbol hashSet literal", c.newLocation(span))
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

func (c *Compiler) hexHashSetLiteral(node *ast.HexHashSetLiteralNode) {
	list := resolve(node)
	span := node.Span()
	if list == nil {
		c.Errors.Add("invalid hex hashSet literal", c.newLocation(span))
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
	c.emitNewCollection(bytecode.NEW_HASH_SET8, bytecode.NEW_HASH_SET32, size, span)
}

func (c *Compiler) emitNewArrayTuple(size int, span *position.Span) {
	c.emitNewCollection(bytecode.NEW_ARRAY_TUPLE8, bytecode.NEW_ARRAY_TUPLE32, size, span)
}

func (c *Compiler) emitNewArrayList(size int, span *position.Span) {
	c.emitNewCollection(bytecode.NEW_ARRAY_LIST8, bytecode.NEW_ARRAY_LIST32, size, span)
}

func (c *Compiler) emitNewHashMap(size int, span *position.Span) {
	c.emitNewCollection(bytecode.NEW_HASH_MAP8, bytecode.NEW_HASH_MAP32, size, span)
}

func (c *Compiler) emitNewHashRecord(size int, span *position.Span) {
	c.emitNewCollection(bytecode.NEW_HASH_RECORD8, bytecode.NEW_HASH_RECORD32, size, span)
}

func (c *Compiler) emitNewRegex(flags bitfield.BitField8, size int, span *position.Span) {
	if size <= math.MaxUint8 {
		c.emit(span.EndPos.Line, bytecode.NEW_REGEX8, flags.Byte(), byte(size))
		return
	}

	if size <= math.MaxUint32 {
		c.emit(span.EndPos.Line, bytecode.NEW_REGEX32)
		c.emitByte(flags.Byte())
		c.emitUint32(uint32(size))
		return
	}

	c.Errors.Add(
		fmt.Sprintf("max number of regex literal elements reached: %d", math.MaxUint32),
		c.newLocation(span),
	)
}

func (c *Compiler) emitNewCollection(opcode8, opcode32 bytecode.OpCode, size int, span *position.Span) {
	if size <= math.MaxUint8 {
		c.emit(span.EndPos.Line, opcode8, byte(size))
		return
	}

	if size <= math.MaxUint32 {
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, uint32(size))
		c.emit(span.EndPos.Line, opcode32, bytes...)
		return
	}

	c.Errors.Add(
		fmt.Sprintf("max number of collection literal elements reached: %d", math.MaxUint32),
		c.newLocation(span),
	)
}

func (c *Compiler) uninterpolatedRegexLiteral(node *ast.UninterpolatedRegexLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	re, err := value.CompileRegex(node.Content, node.Flags)
	if errList, ok := err.(errors.ErrorList); ok {
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
		c.Errors.Add(err.Error(), c.newLocation(node.Span()))
		return
	}

	c.emitValue(re, node.Span())
}

func (c *Compiler) interpolatedRegexLiteral(node *ast.InterpolatedRegexLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	for _, elementNode := range node.Content {
		switch element := elementNode.(type) {
		case *ast.RegexLiteralContentSectionNode:
			c.emitValue(value.String(element.Value), element.Span())
		case *ast.RegexInterpolationNode:
			c.compileNode(element.Expression)
		}
	}
	c.emitNewRegex(node.Flags, len(node.Content), node.Span())
}

func (c *Compiler) interpolatedStringLiteral(node *ast.InterpolatedStringLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	for _, elementNode := range node.Content {
		switch element := elementNode.(type) {
		case *ast.StringLiteralContentSectionNode:
			c.emitValue(value.String(element.Value), element.Span())
		case *ast.StringInterpolationNode:
			c.compileNode(element.Expression)
		}
	}

	c.emitNewCollection(bytecode.NEW_STRING8, bytecode.NEW_STRING32, len(node.Content), node.Span())
}

func (c *Compiler) interpolatedSymbolLiteral(node *ast.InterpolatedSymbolLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	for _, elementNode := range node.Content.Content {
		switch element := elementNode.(type) {
		case *ast.StringLiteralContentSectionNode:
			c.emitValue(value.String(element.Value), element.Span())
		case *ast.StringInterpolationNode:
			c.compileNode(element.Expression)
		}
	}

	c.emitNewCollection(bytecode.NEW_SYMBOL8, bytecode.NEW_SYMBOL32, len(node.Content.Content), node.Span())
}

func (c *Compiler) intLiteral(node *ast.IntLiteralNode) {
	i, err := value.ParseBigInt(node.Value, 0)
	if err != nil {
		c.Errors.Add(err.Error(), c.newLocation(node.Span()))
		return
	}
	if i.IsSmallInt() {
		c.emitValue(i.ToSmallInt(), node.Span())
		return
	}
	c.emitValue(i, node.Span())
}

// Compiles boolean binary operators
func (c *Compiler) logicalExpression(node *ast.LogicalExpressionNode) {
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
		c.Errors.Add(fmt.Sprintf("unknown logical operator: %s", node.Op.String()), c.newLocation(node.Span()))
	}
}

// Compiles the `??` operator
func (c *Compiler) nilCoalescing(node *ast.LogicalExpressionNode) {
	c.compileNode(node.Left)
	nilJump := c.emitJump(node.Span().StartPos.Line, bytecode.JUMP_IF_NIL)
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
	jump := c.emitJump(node.Span().StartPos.Line, bytecode.JUMP_IF)

	// if falsy
	c.emit(node.Span().StartPos.Line, bytecode.POP)
	c.compileNode(node.Right)

	// if truthy
	c.patchJump(jump, node.Span())
}

// Compiles the `&&` operator
func (c *Compiler) logicalAnd(node *ast.LogicalExpressionNode) {
	c.compileNode(node.Left)
	jump := c.emitJump(node.Span().StartPos.Line, bytecode.JUMP_UNLESS)

	// if truthy
	c.emit(node.Span().StartPos.Line, bytecode.POP)
	c.compileNode(node.Right)

	// if falsy
	c.patchJump(jump, node.Span())
}

func (c *Compiler) binaryExpression(node *ast.BinaryExpressionNode) {
	if c.resolveAndEmit(node) {
		return
	}
	c.compileNode(node.Left)
	c.compileNode(node.Right)
	c.emitBinaryOperation(node.Op, node.Span())
}

func (c *Compiler) emitBinaryOperation(opToken *token.Token, span *position.Span) {
	line := span.StartPos.Line
	switch opToken.Type {
	case token.PLUS:
		c.emit(line, bytecode.ADD)
	case token.MINUS:
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
		c.emit(line, bytecode.LESS_EQUAL)
	case token.SPACESHIP_OP:
		c.emit(line, bytecode.COMPARE)
	case token.INSTANCE_OF_OP:
		c.emit(line, bytecode.INSTANCE_OF)
	case token.ISA_OP:
		c.emit(line, bytecode.IS_A)
	default:
		c.Errors.Add(fmt.Sprintf("unknown binary operator: %s", opToken.String()), c.newLocation(span))
	}
}

// Resolves static AST expressions to Elk values
// and emits Bytecode that loads them.
// Returns false when the node cannot be optimised at compile-time
// and no Bytecode has been generated.
func (c *Compiler) resolveAndEmit(node ast.ExpressionNode) bool {
	result := resolve(node)
	if result == nil {
		return false
	}

	c.emitValue(result, node.Span())
	return true
}

func (c *Compiler) resolveAndEmitList(node *ast.ArrayListLiteralNode) bool {
	result := resolveArrayListLiteral(node)
	if result == nil {
		return false
	}

	c.emitValue(result, node.Span())
	return true
}

func (c *Compiler) emitValue(val value.Value, span *position.Span) {
	switch v := val.(type) {
	case value.TrueType:
		c.emit(span.StartPos.Line, bytecode.TRUE)
	case value.FalseType:
		c.emit(span.StartPos.Line, bytecode.FALSE)
	case value.NilType:
		c.emit(span.StartPos.Line, bytecode.NIL)
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
}

func (c *Compiler) emitHashSet(set *value.HashSet, span *position.Span) {
	baseSet := value.NewHashSet(set.Length())
	var mutableElements []value.Value

listLoop:
	for _, element := range set.Table {
		// skip if the bucket is empty or deleted
		if element == nil || element == value.Undefined {
			continue listLoop
		}

		if value.IsMutableCollection(element) {
			mutableElements = append(mutableElements, element)
			continue listLoop
		}

		vm.HashSetAppend(nil, baseSet, element)
	}

	if len(mutableElements) == 0 {
		c.emitLoadValue(baseSet, span)
		c.emit(span.EndPos.Line, bytecode.COPY)
		return
	}

	// capacity
	c.emit(span.StartPos.Line, bytecode.UNDEFINED)
	c.emitLoadValue(baseSet, span)

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
		if element.Key == nil {
			continue listLoop
		}

		if value.IsMutableCollection(element.Key) || value.IsMutableCollection(element.Value) {
			mutablePairs = append(mutablePairs, element)
			continue listLoop
		}

		vm.HashMapSet(nil, baseMap, element.Key, element.Value)
	}

	if len(mutablePairs) == 0 {
		c.emitLoadValue(baseMap, span)
		c.emit(span.EndPos.Line, bytecode.COPY)
		return
	}

	// capacity
	c.emit(span.StartPos.Line, bytecode.UNDEFINED)
	c.emitLoadValue(baseMap, span)

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
		if element.Key == nil {
			continue listLoop
		}

		if value.IsMutableCollection(element.Key) || value.IsMutableCollection(element.Value) {
			mutablePairs = append(mutablePairs, element)
			continue listLoop
		}

		vm.HashRecordSet(nil, baseRecord, element.Key, element.Value)
	}

	if len(mutablePairs) == 0 {
		c.emitLoadValue(baseRecord, span)
		return
	}

	c.emitLoadValue(baseRecord, span)

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
		c.emitLoadValue(list, span)
		c.emit(span.EndPos.Line, bytecode.COPY)
		return
	}

	// capacity
	c.emit(span.StartPos.Line, bytecode.UNDEFINED)

	baseList := l[:firstMutableElementIndex]
	c.emitLoadValue(&baseList, span)

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
		c.emitLoadValue(tuple, span)
		return
	}

	baseTuple := t[:firstMutableElementIndex]
	c.emitLoadValue(&baseTuple, span)

	rest := t[firstMutableElementIndex:]
	for _, element := range rest {
		c.emitValue(element, span)
	}

	c.emitNewArrayList(len(rest), span)
}

func (c *Compiler) unaryExpression(node *ast.UnaryExpressionNode) {
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
		c.Errors.Add(fmt.Sprintf("unknown unary operator: %s", node.Op.String()), c.newLocation(node.Span()))
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
	case moduleMode, mixinMode, classMode, initMethodMode:
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
		c.Errors.Add(
			fmt.Sprintf("too many bytes to jump backward: %d", math.MaxUint16),
			c.newLocation(span),
		)
	}

	c.emitUint16(uint16(offset))
}

// Overwrite the placeholder operand of a jump instruction
func (c *Compiler) patchJumpWithTarget(target int, offset int, span *position.Span) {
	if target > math.MaxUint16 {
		c.Errors.Add(
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
	if index > math.MaxUint8 {
		c.emit(line, bytecode.SET_LOCAL16)
		c.emitUint16(index)
		return
	}

	c.emit(line, bytecode.SET_LOCAL8, byte(index))
}

// Emit an instruction that gets the value of a local.
func (c *Compiler) emitGetLocal(line int, index uint16) {
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

// Emit an instruction that calls a function
func (c *Compiler) emitAddValue(val value.Value, span *position.Span, opCode8, opCode16, opCode32 bytecode.OpCode) int {
	id, size := c.Bytecode.AddValue(val)
	switch size {
	case bytecode.UINT8_SIZE:
		c.Bytecode.AddInstruction(span.StartPos.Line, opCode8, byte(id))
	case bytecode.UINT16_SIZE:
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, uint16(id))
		c.Bytecode.AddInstruction(span.StartPos.Line, opCode16, bytes...)
	case bytecode.UINT32_SIZE:
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, uint32(id))
		c.Bytecode.AddInstruction(span.StartPos.Line, opCode32, bytes...)
	default:
		c.Errors.Add(
			fmt.Sprintf("value pool limit reached: %d", math.MaxUint32),
			c.newLocation(span),
		)
		return -1
	}

	return id
}

// Add a value to the value pool and emit appropriate bytecode.
func (c *Compiler) emitLoadValue(val value.Value, span *position.Span) int {
	return c.emitAddValue(
		val,
		span,
		bytecode.LOAD_VALUE8,
		bytecode.LOAD_VALUE16,
		bytecode.LOAD_VALUE32,
	)
}

// Emit an instruction that instantiates an object
func (c *Compiler) emitInstantiate(callInfo *value.CallSiteInfo, span *position.Span) int {
	return c.emitAddValue(
		callInfo,
		span,
		bytecode.INSTANTIATE8,
		bytecode.INSTANTIATE16,
		bytecode.INSTANTIATE32,
	)
}

// Emit an instruction that sets the value of an instance variable.
func (c *Compiler) emitSetInstanceVariable(name value.Symbol, span *position.Span) int {
	return c.emitAddValue(
		name,
		span,
		bytecode.SET_IVAR8,
		bytecode.SET_IVAR16,
		bytecode.SET_IVAR32,
	)
}

// Emit an instruction that reads the value of an instance variable.
func (c *Compiler) emitGetInstanceVariable(name value.Symbol, span *position.Span) int {
	return c.emitAddValue(
		name,
		span,
		bytecode.GET_IVAR8,
		bytecode.GET_IVAR16,
		bytecode.GET_IVAR32,
	)
}

// Emit an instruction that calls a function
func (c *Compiler) emitCallFunction(callInfo *value.CallSiteInfo, span *position.Span) int {
	return c.emitAddValue(
		callInfo,
		span,
		bytecode.CALL_SELF8,
		bytecode.CALL_SELF16,
		bytecode.CALL_SELF32,
	)
}

// Emit an instruction that calls a method in a pattern
func (c *Compiler) emitCallPattern(callInfo *value.CallSiteInfo, span *position.Span) int {
	return c.emitAddValue(
		callInfo,
		span,
		bytecode.CALL_PATTERN8,
		bytecode.CALL_PATTERN16,
		bytecode.CALL_PATTERN32,
	)
}

// Emit an instruction that calls the `call` method
func (c *Compiler) emitCall(callInfo *value.CallSiteInfo, span *position.Span) int {
	return c.emitAddValue(
		callInfo,
		span,
		bytecode.CALL8,
		bytecode.CALL16,
		bytecode.CALL32,
	)
}

// Emit an instruction that calls a method
func (c *Compiler) emitCallMethod(callInfo *value.CallSiteInfo, span *position.Span) int {
	return c.emitAddValue(
		callInfo,
		span,
		bytecode.CALL_METHOD8,
		bytecode.CALL_METHOD16,
		bytecode.CALL_METHOD32,
	)
}

// Emit an instruction that gets the value of a module constant.
func (c *Compiler) emitGetModConst(name value.Symbol, span *position.Span) int {
	return c.emitAddValue(
		name,
		span,
		bytecode.GET_MOD_CONST8,
		bytecode.GET_MOD_CONST16,
		bytecode.GET_MOD_CONST32,
	)
}

// Emit an instruction that defines a module constant.
func (c *Compiler) emitDefModConst(name value.Symbol, span *position.Span) int {
	return c.emitAddValue(
		name,
		span,
		bytecode.DEF_MOD_CONST8,
		bytecode.DEF_MOD_CONST16,
		bytecode.DEF_MOD_CONST32,
	)
}

// Emit an opcode with optional bytes.
func (c *Compiler) emit(line int, op bytecode.OpCode, bytes ...byte) {
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
func (c *Compiler) defineLocal(name string, span *position.Span, singleAssignment, initialised bool) *local {
	varScope := c.scopes.last()
	_, ok := varScope.localTable[name]
	if ok {
		c.Errors.Add(
			fmt.Sprintf("a variable with this name has already been declared in this scope: %s", name),
			c.newLocation(span),
		)
		return nil
	}
	return c.defineVariableInScope(varScope, name, span, singleAssignment, initialised)
}

// Register a local variable, reusing the variable with the same name that has already been defined in this scope.
func (c *Compiler) defineLocalOverrideCurrentScope(name string, span *position.Span, initialised bool) *local {
	varScope := c.scopes.last()
	if currentVar, ok := varScope.localTable[name]; ok {
		return currentVar
	}
	return c.defineVariableInScope(varScope, name, span, false, initialised)
}

func (c *Compiler) defineVariableInScope(scope *scope, name string, span *position.Span, singleAssignment, initialised bool) *local {
	if c.lastLocalIndex == math.MaxUint16 {
		c.Errors.Add(
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
		index:            uint16(c.lastLocalIndex),
		initialised:      initialised,
		singleAssignment: singleAssignment,
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
