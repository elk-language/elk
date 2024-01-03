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

// Compile the Elk source to a Bytecode chunk.
func CompileSource(sourceName string, source string) (*vm.BytecodeMethod, errors.ErrorList) {
	ast, err := parser.Parse(sourceName, source)
	if err != nil {
		return nil, err
	}

	return CompileAST(sourceName, ast)
}

// Compile the AST node to a Bytecode chunk.
func CompileAST(sourceName string, ast ast.Node) (*vm.BytecodeMethod, errors.ErrorList) {
	compiler := new("main", topLevelMode, position.NewLocationWithSpan(sourceName, ast.Span()))
	compiler.compileProgram(ast)

	return compiler.Bytecode, compiler.Errors
}

// Compile code for use in the REPL.
func CompileREPL(sourceName string, source string) (*Compiler, errors.ErrorList) {
	ast, err := parser.Parse(sourceName, source)
	if err != nil {
		return nil, err
	}

	compiler := new("main", topLevelMode, position.NewLocationWithSpan(sourceName, ast.Span()))
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
	methodMode
	setterMethodMode
	initMethodMode
)

// represents a local variable or value
type local struct {
	index            uint16
	singleAssignment bool
	initialised      bool
}

// set of local variables
type localTable map[string]*local

// indices represent scope depths
// and elements are sets of local variable names in a particular scope
type scopes []localTable

// Get the last local variable scope.
func (s scopes) last() localTable {
	return s[len(s)-1]
}

type loopJumpInfoType uint8

const (
	breakLoopJump loopJumpInfoType = iota
	continueLoopJump
)

type loopJumpInfo struct {
	typ    loopJumpInfoType
	offset int
	span   *position.Span
}

type loopJumpSet struct {
	label     string
	infinite  bool
	loopJumps []*loopJumpInfo
}

// Holds the state of the Compiler.
type Compiler struct {
	Name             string
	Bytecode         *vm.BytecodeMethod
	Errors           errors.ErrorList
	scopes           scopes
	loopJumpSets     []*loopJumpSet
	lastLocalIndex   int // index of the last local variable
	maxLocalIndex    int // max index of a local variable
	predefinedLocals int
	mode             mode
	lastOpCode       bytecode.OpCode
}

// Instantiate a new Compiler instance.
func new(name string, mode mode, loc *position.Location) *Compiler {
	c := &Compiler{
		Bytecode: vm.NewBytecodeMethodSimple(
			value.ToSymbol(name),
			[]byte{},
			loc,
		),
		scopes:         scopes{localTable{}}, // start with an empty set for the 0th scope
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
	case methodMode, setterMethodMode, initMethodMode:
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

	compiler := new("main", topLevelMode, position.NewLocationWithSpan(filename, ast.Span()))
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
	c.prepLocals()
	c.emitReturn(node.Span(), nil)
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
		case ast.PositionalRestParameterKind:
			positionalRestParamSeen = true
		case ast.NamedRestParameterKind:
			c.Errors.Add(
				fmt.Sprintf("named rest parameters are not supported yet: %s", p.Name),
				c.newLocation(pSpan),
			)
			continue
		}

		if positionalRestParamSeen {
			c.Bytecode.IncrementPostRestParameterCount()
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
	c.prepLocals()

	c.emitReturn(span, nil)
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
	c.prepLocals()

	c.emitReturn(span, nil)
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
	if c.maxLocalIndex >= math.MaxUint8 {
		newInstructions = make([]byte, 0, len(c.Bytecode.Instructions)+3)
		newInstructions = append(newInstructions, byte(bytecode.PREP_LOCALS16))
		newInstructions = binary.BigEndian.AppendUint16(newInstructions, uint16(localCount))
	} else {
		newInstructions = make([]byte, 0, len(c.Bytecode.Instructions)+2)
		newInstructions = append(newInstructions, byte(bytecode.PREP_LOCALS8), byte(localCount))
	}

	c.Bytecode.Instructions = append(
		newInstructions,
		c.Bytecode.Instructions...,
	)
	lineInfo := c.Bytecode.LineInfoList.First()
	if lineInfo != nil {
		lineInfo.InstructionCount++
	}
}

func (c *Compiler) initLoopJumpSet(label string, infinite bool) {
	c.loopJumpSets = append(
		c.loopJumpSets,
		&loopJumpSet{
			label:    label,
			infinite: infinite,
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
	case *ast.InitDefinitionNode:
		c.initDefinition(node)
	case *ast.IncludeExpressionNode:
		c.includeExpression(node)
	case *ast.ExtendExpressionNode:
		c.extendExpression(node)
	case *ast.SingletonBlockExpressionNode:
		c.singletonBlock(node)
	case *ast.AttributeAccessNode:
		c.attributeAccess(node)
	case *ast.ConstructorCallNode:
		c.constructorCall(node)
	case *ast.MethodCallNode:
		c.methodCall(node)
	case *ast.FunctionCallNode:
		c.functionCall(node)
	case *ast.ReturnExpressionNode:
		c.returnExpression(node)
	case *ast.VariableDeclarationNode:
		c.variableDeclaration(node)
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
	case *ast.TupleLiteralNode:
		c.tupleLiteral(node)
	case *ast.RawStringLiteralNode:
		c.emitValue(value.String(node.Value), node.Span())
	case *ast.DoubleQuotedStringLiteralNode:
		c.emitValue(value.String(node.Value), node.Span())
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
	case *ast.DoExpressionNode:
		c.enterScope()
		c.compileStatements(node.Body, node.Span())
		c.leaveScope(node.Span().EndPos.Line)
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

func (c *Compiler) breakExpression(node *ast.BreakExpressionNode) {
	span := node.Span()
	if node.Value == nil {
		c.emit(span.StartPos.Line, bytecode.NIL)
	} else {
		c.compileNode(node.Value)
	}

	breakJumpOffset := c.emitJump(span.StartPos.Line, bytecode.JUMP)
	c.addLoopJump(node.Label, breakLoopJump, breakJumpOffset, span)
}

func (c *Compiler) continueExpression(node *ast.ContinueExpressionNode) {
	span := node.Span()
	loop := c.findLoopJumpSet(node.Label, span)
	if loop == nil {
		return
	}

	if loop.infinite {
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

	continueJumpOffset := c.emitJump(span.StartPos.Line, bytecode.LOOP)
	c.addLoopJumpTo(loop, continueLoopJump, continueJumpOffset)
}

// Patch loop jump addresses for `break` and `continue` expressions.
func (c *Compiler) patchLoopJumps(continueOffset int) {
	lastLoopJumpSet := c.loopJumpSets[len(c.loopJumpSets)-1]
	for _, loopJump := range lastLoopJumpSet.loopJumps {
		switch loopJump.typ {
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
	c.enterScope()
	defer c.leaveScope(span.EndPos.Line)
	c.initLoopJumpSet(label, true)

	start := c.nextInstructionOffset()
	if c.compileStatementsOk(body, span) {
		c.emit(span.EndPos.Line, bytecode.POP)
	}
	c.emitLoop(span, start)

	c.patchLoopJumps(start)
}

func (c *Compiler) whileExpression(label string, node *ast.WhileExpressionNode) {
	span := node.Span()

	if result, ok := resolve(node.Condition); ok {
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

	c.enterScope()
	defer c.leaveScope(span.EndPos.Line)
	c.initLoopJumpSet(label, false)

	c.emit(span.StartPos.Line, bytecode.NIL)
	// loop start
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

	// jump to loop condition
	c.emitLoop(span, start)

	// after loop
	c.patchJump(loopBodyOffset, span)
	// pop the condition value
	c.emit(span.EndPos.Line, bytecode.POP)
	c.patchLoopJumps(start)
}

func (c *Compiler) modifierWhileExpression(label string, node *ast.ModifierNode) {
	span := node.Span()

	body := node.Left
	condition := node.Right

	var conditionIsStaticFalsy bool

	if result, ok := resolve(condition); ok {
		if value.Truthy(result) {
			// the loop is endless
			c.loopExpression(label, ast.ExpressionToStatements(body), span)
			return
		}

		// the loop will only iterate once
		conditionIsStaticFalsy = true
	}

	c.enterScope()
	defer c.leaveScope(span.EndPos.Line)
	c.initLoopJumpSet(label, false)

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
	c.patchLoopJumps(continueOffset)
}

func (c *Compiler) modifierUntilExpression(label string, node *ast.ModifierNode) {
	span := node.Span()

	body := node.Left
	condition := node.Right

	var conditionIsStaticTruthy bool

	if result, ok := resolve(condition); ok {
		if value.Falsy(result) {
			// the loop is endless
			c.loopExpression(label, ast.ExpressionToStatements(body), span)
			return
		}

		// the loop will only iterate once
		conditionIsStaticTruthy = true
	}

	c.enterScope()
	defer c.leaveScope(span.EndPos.Line)
	c.initLoopJumpSet(label, false)

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
	c.patchLoopJumps(continueOffset)
}

func (c *Compiler) untilExpression(label string, node *ast.UntilExpressionNode) {
	span := node.Span()

	if result, ok := resolve(node.Condition); ok {
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

	c.enterScope()
	defer c.leaveScope(span.EndPos.Line)
	c.initLoopJumpSet(label, false)

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

// Compile a numeric for loop eg. `for i := 0; i < 5; i += 1 then println(i)`
func (c *Compiler) numericForExpression(label string, node *ast.NumericForExpressionNode) {
	span := node.Span()

	if node.Initialiser == nil && node.Condition == nil && node.Increment == nil {
		// the loop is endless
		c.loopExpression(label, node.ThenBody, span)
		return
	}

	c.enterScope()
	defer c.leaveScope(span.EndPos.Line)
	c.initLoopJumpSet(label, false)

	// loop initialiser eg. `i := 0`
	if node.Initialiser != nil {
		c.compileNode(node.Initialiser)
		c.emit(span.EndPos.Line, bytecode.POP)
	}

	c.emit(span.StartPos.Line, bytecode.NIL)
	// loop start
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

	// jump to loop condition
	c.emitLoop(span, start)

	// after loop
	if node.Condition != nil {
		c.patchJump(loopBodyOffset, span)
		// pop the condition value
		c.emit(span.EndPos.Line, bytecode.POP)
	}

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

func (c *Compiler) assignment(node *ast.AssignmentExpressionNode) {
	switch n := node.Left.(type) {
	case *ast.PublicIdentifierNode:
		c.localVariableAssignment(n.Value, node.Op, node.Right, node.Span())
	case *ast.PrivateIdentifierNode:
		c.localVariableAssignment(n.Value, node.Op, node.Right, node.Span())
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
		switch c.mode {
		case topLevelMode:
			c.Errors.Add(
				"instance variables cannot be set in the top level",
				c.newLocation(node.Span()),
			)
		}
		c.compileNode(node.Right)
		c.emitSetInstanceVariable(value.ToSymbol(n.Value), n.Span())
	case *ast.AttributeAccessNode:
		// compile the argument
		switch node.Op.Type {
		case token.EQUAL_OP:
			c.compileNode(n.Receiver)
			c.compileNode(node.Right)
			c.emitSetterCall(n.AttributeName, node.Span())
		case token.OR_OR_EQUAL:
			span := node.Span()
			// Read the current value
			c.compileNode(n.Receiver)
			c.emitGetterCall(n.AttributeName, span)

			jump := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF)

			// if falsy
			c.emit(span.StartPos.Line, bytecode.POP)
			c.compileNode(node.Right)
			c.emitSetterCall(n.AttributeName, span)

			// if truthy
			c.patchJump(jump, span)
		case token.AND_AND_EQUAL:
			span := node.Span()
			// Read the current value
			c.compileNode(n.Receiver)
			c.emitGetterCall(n.AttributeName, span)

			jump := c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)

			// if truthy
			c.emit(span.StartPos.Line, bytecode.POP)
			c.compileNode(node.Right)
			c.emitSetterCall(n.AttributeName, span)

			// if falsy
			c.patchJump(jump, span)
		case token.QUESTION_QUESTION_EQUAL:
			span := node.Span()
			// Read the current value
			c.compileNode(n.Receiver)
			c.emitGetterCall(n.AttributeName, span)

			nilJump := c.emitJump(span.StartPos.Line, bytecode.JUMP_IF_NIL)
			nonNilJump := c.emitJump(span.StartPos.Line, bytecode.JUMP)

			// if nil
			c.patchJump(nilJump, span)
			c.emit(span.StartPos.Line, bytecode.POP)
			c.compileNode(node.Right)
			c.emitSetterCall(n.AttributeName, span)

			// if not nil
			c.patchJump(nonNilJump, span)
		case token.PLUS_EQUAL:
			c.complexSetterCall(bytecode.ADD, n, node.Right, node.Span())
		case token.MINUS_EQUAL:
			c.complexSetterCall(bytecode.SUBTRACT, n, node.Right, node.Span())
		case token.STAR_EQUAL:
			c.complexSetterCall(bytecode.MULTIPLY, n, node.Right, node.Span())
		case token.SLASH_EQUAL:
			c.complexSetterCall(bytecode.DIVIDE, n, node.Right, node.Span())
		case token.STAR_STAR_EQUAL:
			c.complexSetterCall(bytecode.EXPONENTIATE, n, node.Right, node.Span())
		case token.PERCENT_EQUAL:
			c.complexSetterCall(bytecode.MODULO, n, node.Right, node.Span())
		case token.LBITSHIFT_EQUAL:
			c.complexSetterCall(bytecode.LBITSHIFT, n, node.Right, node.Span())
		case token.LTRIPLE_BITSHIFT_EQUAL:
			c.complexSetterCall(bytecode.LOGIC_LBITSHIFT, n, node.Right, node.Span())
		case token.RBITSHIFT_EQUAL:
			c.complexSetterCall(bytecode.RBITSHIFT, n, node.Right, node.Span())
		case token.RTRIPLE_BITSHIFT_EQUAL:
			c.complexSetterCall(bytecode.LOGIC_RBITSHIFT, n, node.Right, node.Span())
		case token.AND_EQUAL:
			c.complexSetterCall(bytecode.BITWISE_AND, n, node.Right, node.Span())
		case token.OR_EQUAL:
			c.complexSetterCall(bytecode.BITWISE_OR, n, node.Right, node.Span())
		case token.XOR_EQUAL:
			c.complexSetterCall(bytecode.BITWISE_XOR, n, node.Right, node.Span())
		default:
			c.Errors.Add(fmt.Sprintf("unknown binary operator: %s", node.Op.String()), c.newLocation(node.Span()))
		}

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
	local, ok := c.localVariableAccess(name, span)
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
	c.emitSetLocal(span.StartPos.Line, local.index)
}

// Return the offset of the Bytecode next instruction.
func (c *Compiler) nextInstructionOffset() int {
	return len(c.Bytecode.Instructions)
}

func (c *Compiler) setLocal(name string, valueNode ast.ExpressionNode, span *position.Span) {
	c.compileNode(valueNode)
	local, ok := c.resolveLocal(name, span)
	if !ok {
		return
	}
	if local.initialised && local.singleAssignment {
		c.Errors.Add(
			fmt.Sprintf("cannot reassign a val: %s", name),
			c.newLocation(span),
		)
	}
	local.initialised = true
	c.emitSetLocal(span.StartPos.Line, local.index)
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

func (c *Compiler) localVariableAccess(name string, span *position.Span) (*local, bool) {
	local, ok := c.resolveLocal(name, span)
	if !ok {
		return nil, false
	}
	if !local.initialised {
		c.Errors.Add(
			fmt.Sprintf("cannot access an uninitialised local: %s", name),
			c.newLocation(span),
		)
		return nil, false
	}

	c.emitGetLocal(span.StartPos.Line, local.index)
	return local, true
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
	if result, ok := resolve(condition); ok {
		// if gets optimised away
		c.enterScope()
		defer c.leaveScope(span.StartPos.Line)

		var truthyBody, falsyBody ast.ExpressionNode
		if unless {
			truthyBody = els
			falsyBody = then
		} else {
			truthyBody = then
			falsyBody = els
		}
		if value.Truthy(result) {
			if truthyBody == nil {
				c.emit(span.StartPos.Line, bytecode.NIL)
				return
			}
			c.compileNode(truthyBody)
			return
		}

		if falsyBody == nil {
			c.emit(span.StartPos.Line, bytecode.NIL)
			return
		}
		c.compileNode(falsyBody)
		return
	}

	c.enterScope()
	c.compileNode(condition)
	var jumpOp bytecode.OpCode
	if unless {
		jumpOp = bytecode.JUMP_IF
	} else {
		jumpOp = bytecode.JUMP_UNLESS
	}
	thenJumpOffset := c.emitJump(span.StartPos.Line, jumpOp)
	c.emit(span.StartPos.Line, bytecode.POP)

	c.compileNode(then)
	c.leaveScope(span.StartPos.Line)

	elseJumpOffset := c.emitJump(span.StartPos.Line, bytecode.JUMP)

	c.patchJump(thenJumpOffset, span)
	c.emit(span.StartPos.Line, bytecode.POP)

	if els != nil {
		c.enterScope()
		c.compileNode(els)
		c.leaveScope(span.StartPos.Line)
	} else {
		c.emit(span.StartPos.Line, bytecode.NIL)
	}
	c.patchJump(elseJumpOffset, span)
}

func (c *Compiler) ifExpression(unless bool, condition ast.ExpressionNode, then, els []ast.StatementNode, span *position.Span) {
	if result, ok := resolve(condition); ok {
		// if gets optimised away
		c.enterScope()
		defer c.leaveScope(span.StartPos.Line)

		var truthyBody, falsyBody []ast.StatementNode
		if unless {
			truthyBody = els
			falsyBody = then
		} else {
			truthyBody = then
			falsyBody = els
		}
		if value.Truthy(result) {
			c.compileStatements(truthyBody, span)
			return
		}

		c.compileStatements(falsyBody, span)
		return
	}

	c.enterScope()
	c.compileNode(condition)

	var jumpOp bytecode.OpCode
	if unless {
		jumpOp = bytecode.JUMP_IF
	} else {
		jumpOp = bytecode.JUMP_UNLESS
	}
	thenJumpOffset := c.emitJump(span.StartPos.Line, jumpOp)

	c.emit(span.StartPos.Line, bytecode.POP)

	c.compileStatements(then, span)
	c.leaveScope(span.StartPos.Line)

	elseJumpOffset := c.emitJump(span.StartPos.Line, bytecode.JUMP)

	c.patchJump(thenJumpOffset, span)
	c.emit(span.StartPos.Line, bytecode.POP)

	if els != nil {
		c.enterScope()
		c.compileStatements(els, span)
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

func (c *Compiler) functionCall(node *ast.FunctionCallNode) {
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

func (c *Compiler) attributeAccess(node *ast.AttributeAccessNode) {
	c.compileNode(node.Receiver)

	name := value.ToSymbol(node.AttributeName)
	callInfo := value.NewCallSiteInfo(name, 0, nil)
	c.emitCallMethod(callInfo, node.Span())
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

	if node.NilSafe {
		c.Errors.Add(
			fmt.Sprintf("nil safe method calls are not supported yet: %s", node.MethodName),
			c.newLocation(node.Span()),
		)
		return
	}

	name := value.ToSymbol(node.MethodName)
	argumentCount := len(node.PositionalArguments) + len(node.NamedArguments)
	callInfo := value.NewCallSiteInfo(name, argumentCount, namedArgs)
	c.emitCallMethod(callInfo, node.Span())
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
	case methodMode, setterMethodMode, initMethodMode:
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
		singletonCompiler := new("singleton_class", classMode, c.newLocation(span))
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
	case methodMode, setterMethodMode, initMethodMode:
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
		mode = methodMode
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

func (c *Compiler) initDefinition(node *ast.InitDefinitionNode) {
	switch c.mode {
	case methodMode, setterMethodMode, initMethodMode:
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
	case methodMode:
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
	case methodMode:
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
	case methodMode:
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
		mixinCompiler := new("mixin", mixinMode, c.newLocation(node.Span()))
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
	case methodMode:
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
		modCompiler := new("module", moduleMode, c.newLocation(node.Span()))
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
	case methodMode:
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
	case methodMode:
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
	case methodMode:
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
	case methodMode:
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
	case methodMode:
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
		modCompiler := new("class", classMode, c.newLocation(node.Span()))
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
		if sExpr, ok := s.(*ast.ExpressionStatementNode); ok && isLast && c.resolveAndEmit(sExpr.Expression) {
			continue
		}
		c.compileNode(s)
		if !isLast {
			c.emit(s.Span().EndPos.Line, bytecode.POP)
		}
	}

	return len(nonEmptyStatements) != 0
}

func (c *Compiler) tupleLiteral(node *ast.TupleLiteralNode) {
	if c.resolveAndEmit(node) {
		return
	}

	panic("non-static tuple literal are not supported yet")
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
	case token.OR:
		c.emit(line, bytecode.BITWISE_OR)
	case token.XOR:
		c.emit(line, bytecode.BITWISE_XOR)
	case token.PERCENT:
		c.emit(line, bytecode.MODULO)
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
	default:
		c.Errors.Add(fmt.Sprintf("unknown binary operator: %s", opToken.String()), c.newLocation(span))
	}
}

// Resolves static AST expressions to Elk values
// and emits Bytecode that loads them.
// Returns false when the node cannot be optimised at compile-time
// and no Bytecode has been generated.
func (c *Compiler) resolveAndEmit(node ast.ExpressionNode) bool {
	result, ok := resolve(node)
	if !ok {
		return false
	}

	switch result.(type) {
	case value.TrueType:
		c.emit(node.Span().StartPos.Line, bytecode.TRUE)
	case value.FalseType:
		c.emit(node.Span().StartPos.Line, bytecode.FALSE)
	case value.NilType:
		c.emit(node.Span().StartPos.Line, bytecode.NIL)
	default:
		c.emitValue(result, node.Span())
	}
	return true
}

func (c *Compiler) unaryExpression(node *ast.UnaryExpressionNode) {
	if c.resolveAndEmit(node) {
		return
	}
	c.compileNode(node.Right)
	switch node.Op.Type {
	case token.PLUS:
		// TODO: Implement unary plus
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
	case bytecode.RETURN, bytecode.RETURN_FIRST_ARG:
		return
	}

	switch c.mode {
	case setterMethodMode:
		if value == nil {
			c.emit(span.EndPos.Line, bytecode.POP)
		}
		c.emit(span.EndPos.Line, bytecode.RETURN_FIRST_ARG)
	case moduleMode, mixinMode, classMode, initMethodMode:
		if value == nil {
			c.emit(span.EndPos.Line, bytecode.POP)
		}
		c.emit(span.EndPos.Line, bytecode.RETURN_SELF)
	default:
		if value != nil {
			c.compileNode(value)
		}
		c.emit(span.EndPos.Line, bytecode.RETURN)
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

	c.Bytecode.AppendUint16(uint16(offset))
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
		c.Bytecode.AppendUint16(index)
		return
	}

	c.emit(line, bytecode.SET_LOCAL8, byte(index))
}

// Emit an instruction that gets the value of a local.
func (c *Compiler) emitGetLocal(line int, index uint16) {
	if index > math.MaxUint8 {
		c.emit(line, bytecode.GET_LOCAL16)
		c.Bytecode.AppendUint16(index)
		return
	}

	c.emit(line, bytecode.GET_LOCAL8, byte(index))
}

// Emit an instruction that calls a function
func (c *Compiler) emitAddValue(val value.Value, span *position.Span, opCode8, opCode16, opCode32 bytecode.OpCode) {
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
	}
}

// Add a value to the value pool and emit appropriate bytecode.
func (c *Compiler) emitValue(val value.Value, span *position.Span) {
	c.emitAddValue(
		val,
		span,
		bytecode.LOAD_VALUE8,
		bytecode.LOAD_VALUE16,
		bytecode.LOAD_VALUE32,
	)
}

// Emit an instruction that instantiates an object
func (c *Compiler) emitInstantiate(callInfo *value.CallSiteInfo, span *position.Span) {
	c.emitAddValue(
		callInfo,
		span,
		bytecode.INSTANTIATE8,
		bytecode.INSTANTIATE16,
		bytecode.INSTANTIATE32,
	)
}

// Emit an instruction that sets the value of an instance variable.
func (c *Compiler) emitSetInstanceVariable(name value.Symbol, span *position.Span) {
	c.emitAddValue(
		name,
		span,
		bytecode.SET_IVAR8,
		bytecode.SET_IVAR16,
		bytecode.SET_IVAR32,
	)
}

// Emit an instruction that reads the value of an instance variable.
func (c *Compiler) emitGetInstanceVariable(name value.Symbol, span *position.Span) {
	c.emitAddValue(
		name,
		span,
		bytecode.GET_IVAR8,
		bytecode.GET_IVAR16,
		bytecode.GET_IVAR32,
	)
}

// Emit an instruction that calls a function
func (c *Compiler) emitCallFunction(callInfo *value.CallSiteInfo, span *position.Span) {
	c.emitAddValue(
		callInfo,
		span,
		bytecode.CALL_FUNCTION8,
		bytecode.CALL_FUNCTION16,
		bytecode.CALL_FUNCTION32,
	)
}

// Emit an instruction that calls a method
func (c *Compiler) emitCallMethod(callInfo *value.CallSiteInfo, span *position.Span) {
	c.emitAddValue(
		callInfo,
		span,
		bytecode.CALL_METHOD8,
		bytecode.CALL_METHOD16,
		bytecode.CALL_METHOD32,
	)
}

// Emit an instruction that gets the value of a module constant.
func (c *Compiler) emitGetModConst(name value.Symbol, span *position.Span) {
	c.emitAddValue(
		name,
		span,
		bytecode.GET_MOD_CONST8,
		bytecode.GET_MOD_CONST16,
		bytecode.GET_MOD_CONST32,
	)
}

// Emit an instruction that defines a module constant.
func (c *Compiler) emitDefModConst(name value.Symbol, span *position.Span) {
	c.emitAddValue(
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

func (c *Compiler) enterScope() {
	c.scopes = append(c.scopes, localTable{})
}

func (c *Compiler) leaveScope(line int) {
	currentDepth := len(c.scopes) - 1

	varsToPop := len(c.scopes[currentDepth])
	if varsToPop > 0 {
		if c.lastLocalIndex > math.MaxUint8 || varsToPop > math.MaxUint8 {
			c.emit(line, bytecode.LEAVE_SCOPE32)
			c.Bytecode.AppendUint16(uint16(c.lastLocalIndex))
			c.Bytecode.AppendUint16(uint16(varsToPop))
		} else {
			c.emit(line, bytecode.LEAVE_SCOPE16, byte(c.lastLocalIndex), byte(varsToPop))
		}
	}

	c.lastLocalIndex -= varsToPop
	c.scopes[currentDepth] = nil
	c.scopes = c.scopes[:currentDepth]
}

// Register a local variable.
func (c *Compiler) defineLocal(name string, span *position.Span, singleAssignment, initialised bool) *local {
	varScope := c.scopes.last()
	_, ok := varScope[name]
	if ok {
		c.Errors.Add(
			fmt.Sprintf("a variable with this name has already been declared in this scope: %s", name),
			c.newLocation(span),
		)
		return nil
	}
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
	varScope[name] = newVar
	return newVar
}

// Resolve a local variable and get its index.
func (c *Compiler) resolveLocal(name string, span *position.Span) (*local, bool) {
	var localVal *local
	var found bool
	for i := len(c.scopes) - 1; i >= 0; i-- {
		varScope := c.scopes[i]
		local, ok := varScope[name]
		if !ok {
			continue
		}
		localVal = local
		found = true
		break
	}

	if !found {
		c.Errors.Add(
			fmt.Sprintf("undeclared variable: %s", name),
			c.newLocation(span),
		)
		return localVal, false
	}

	return localVal, true
}
