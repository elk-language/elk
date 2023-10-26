// Package compiler implements
// the Elk bytecode compiler.
// It takes in Elk source code and outputs
// Elk bytecode that can be run the Elk VM.
package compiler

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/value"

	"github.com/elk-language/elk/token"
)

// Compile the Elk source to a bytecode chunk.
func CompileSource(sourceName string, source string) (*value.BytecodeFunction, errors.ErrorList) {
	ast, err := parser.Parse(sourceName, source)
	if err != nil {
		return nil, err
	}

	return CompileAST(sourceName, ast)
}

// Compile the AST node to a bytecode chunk.
func CompileAST(sourceName string, ast ast.Node) (*value.BytecodeFunction, errors.ErrorList) {
	compiler := new(
		sourceName,
		position.NewLocationWithSpan(sourceName, ast.Span()),
	)
	compiler.compile(ast)

	return compiler.function, compiler.errors
}

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

// Holds the state of the compiler.
type compiler struct {
	sourceName     string
	function       *value.BytecodeFunction
	errors         errors.ErrorList
	scopes         scopes
	lastLocalIndex int // index of the last local variable
	maxLocalIndex  int // max index of a local variable
}

// Instantiate a new compiler instance.
func new(sourceName string, loc *position.Location) *compiler {
	return &compiler{
		function: value.NewBytecodeFunction(
			[]byte{},
			loc,
		),
		scopes:         scopes{localTable{}}, // start with an empty set for the 0th scope
		lastLocalIndex: -1,
		maxLocalIndex:  -1,
		sourceName:     sourceName,
	}
}

// Entry point to the compilation process
func (c *compiler) compile(node ast.Node) {
	c.compileNode(node)
	c.prepLocals()
	c.emit(node.Span().EndPos.Line, bytecode.RETURN)
}

// Create a new location struct with the given position.
func (c *compiler) newLocation(span *position.Span) *position.Location {
	return position.NewLocationWithSpan(c.sourceName, span)
}

func (c *compiler) prepLocals() {
	if c.maxLocalIndex < 0 {
		return
	}

	var newInstructions []byte
	if c.maxLocalIndex >= math.MaxUint8 {
		newInstructions = make([]byte, 0, len(c.function.Instructions)+5)
		newInstructions = append(newInstructions, byte(bytecode.PREP_LOCALS16))
		newInstructions = binary.BigEndian.AppendUint16(newInstructions, uint16(c.maxLocalIndex+1))
	} else {
		newInstructions = make([]byte, 0, len(c.function.Instructions)+2)
		newInstructions = append(newInstructions, byte(bytecode.PREP_LOCALS8), byte(c.maxLocalIndex+1))
	}

	c.function.Instructions = append(
		newInstructions,
		c.function.Instructions...,
	)
	lineInfo := c.function.LineInfoList.First()
	lineInfo.InstructionCount++
}

func (c *compiler) compileNode(node ast.Node) {
	switch node := node.(type) {
	case *ast.ProgramNode:
		c.compileStatements(node.Body, node.Span())
	case *ast.ExpressionStatementNode:
		c.compileNode(node.Expression)
	case *ast.ConstantLookupNode:
		c.constantLookup(node)
	case *ast.AssignmentExpressionNode:
		c.assignment(node)
	case *ast.VariableDeclarationNode:
		c.variableDeclaration(node)
	case *ast.ValueDeclarationNode:
		c.valueDeclaration(node)
	case *ast.PublicIdentifierNode:
		c.localVariableAccess(node.Value, node.Span())
	case *ast.PrivateIdentifierNode:
		c.localVariableAccess(node.Value, node.Span())
	case *ast.BinaryExpressionNode:
		c.binaryExpression(node)
	case *ast.LogicalExpressionNode:
		c.logicalExpression(node)
	case *ast.UnaryExpressionNode:
		c.unaryExpression(node)
	case *ast.RawStringLiteralNode:
		c.emitConstant(value.String(node.Value), node.Span())
	case *ast.DoubleQuotedStringLiteralNode:
		c.emitConstant(value.String(node.Value), node.Span())
	case *ast.CharLiteralNode:
		c.emitConstant(value.Char(node.Value), node.Span())
	case *ast.RawCharLiteralNode:
		c.emitConstant(value.Char(node.Value), node.Span())
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
		c.modifierExpression(node)
	case *ast.LoopExpressionNode:
		c.loopExpression(node)
	case *ast.NumericForExpressionNode:
		c.numericForExpression(node)
	case *ast.SimpleSymbolLiteralNode:
		c.emitConstant(value.SymbolTable.Add(node.Content), node.Span())
	case *ast.IntLiteralNode:
		c.intLiteral(node)
	case *ast.Int8LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 8)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		// BENCHMARK: Compare with storing
		// ints inline in bytecode instead of as constants.
		c.emitConstant(value.Int8(i), node.Span())
	case *ast.Int16LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 16)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitConstant(value.Int16(i), node.Span())
	case *ast.Int32LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 32)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitConstant(value.Int32(i), node.Span())
	case *ast.Int64LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 64)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitConstant(value.Int64(i), node.Span())
	case *ast.UInt8LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 8)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitConstant(value.UInt8(i), node.Span())
	case *ast.UInt16LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 16)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitConstant(value.UInt16(i), node.Span())
	case *ast.UInt32LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 32)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitConstant(value.UInt32(i), node.Span())
	case *ast.UInt64LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 64)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitConstant(value.UInt64(i), node.Span())
	case *ast.FloatLiteralNode:
		f, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitConstant(value.Float(f), node.Span())
	case *ast.BigFloatLiteralNode:
		f, err := value.ParseBigFloat(node.Value)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitConstant(f, node.Span())
	case *ast.Float64LiteralNode:
		f, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitConstant(value.Float64(f), node.Span())
	case *ast.Float32LiteralNode:
		f, err := strconv.ParseFloat(node.Value, 32)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Span()))
			return
		}
		c.emitConstant(value.Float32(f), node.Span())

	case nil:
	default:
		c.errors.Add(
			fmt.Sprintf("compilation of this node has not been implemented: %T", node),
			c.newLocation(node.Span()),
		)
	}
}

func (c *compiler) loopExpression(node *ast.LoopExpressionNode) {
	c.enterScope()
	defer c.leaveScope(node.Span().EndPos.Line)

	start := c.nextInstructionOffset()
	if c.compileStatementsOk(node.ThenBody, node.Span()) {
		c.emit(node.Span().EndPos.Line, bytecode.POP)
	}
	c.emitLoop(node.Span(), start)
}

// Compile a constant lookup expressions eg. `Foo::Bar`
func (c *compiler) constantLookup(node *ast.ConstantLookupNode) {
	if node.Left == nil {
		c.emit(node.Span().StartPos.Line, bytecode.ROOT)
	} else {
		c.compileNode(node.Left)
	}

	switch r := node.Right.(type) {
	case *ast.PublicConstantNode:
		c.emitGetModConst(value.SymbolTable.Add(r.Value), node.Span())
	default:
		c.errors.Add(
			fmt.Sprintf("incorrect right side of constant lookup: %T", node),
			c.newLocation(node.Span()),
		)
	}
}

// Compile a numeric for loop eg. `for i := 0; i < 5; i += 1 then println(i)`
func (c *compiler) numericForExpression(node *ast.NumericForExpressionNode) {
	span := node.Span()
	c.enterScope()
	defer c.leaveScope(span.EndPos.Line)

	// loop initialiser eg. `i := 0`
	if node.Initialiser != nil {
		c.compileNode(node.Initialiser)
		c.emit(span.EndPos.Line, bytecode.POP)
	}

	// loop start
	start := c.nextInstructionOffset()

	var loopBodyOffset int
	// loop condition eg. `i < 5`
	if node.Condition != nil {
		c.compileNode(node.Condition)
		// jump past the loop if the condition is falsy
		loopBodyOffset = c.emitJump(span.StartPos.Line, bytecode.JUMP_UNLESS)
		c.emit(span.EndPos.Line, bytecode.POP)
	}

	// loop body
	if c.compileStatementsOk(node.ThenBody, span) {
		c.emit(span.EndPos.Line, bytecode.POP)
	}

	if node.Increment != nil {
		// increment step eg. `i += 1`
		c.compileNode(node.Increment)
	}

	// jump to loop condition
	c.emitLoop(span, start)

	// after loop
	if node.Condition != nil {
		c.patchJump(loopBodyOffset, span)
		// pop the condition value
		c.emit(span.EndPos.Line, bytecode.POP)
	}
	c.emit(span.EndPos.Line, bytecode.NIL)
}

func (c *compiler) assignment(node *ast.AssignmentExpressionNode) {
	switch n := node.Left.(type) {
	case *ast.PublicIdentifierNode:
		c.localVariableAssignment(n.Value, node.Op, node.Right, node.Span())
	case *ast.PrivateIdentifierNode:
		c.localVariableAssignment(n.Value, node.Op, node.Right, node.Span())
	default:
		c.errors.Add(
			fmt.Sprintf("can't assign to: %T", node),
			c.newLocation(node.Span()),
		)
	}
}

func (c *compiler) complexAssignment(name string, valueNode ast.ExpressionNode, opcode bytecode.OpCode, span *position.Span) {
	local, ok := c.localVariableAccess(name, span)
	if !ok {
		return
	}
	c.compileNode(valueNode)
	c.emit(span.StartPos.Line, opcode)

	if local.initialised && local.singleAssignment {
		c.errors.Add(
			fmt.Sprintf("can't reassign a val: %s", name),
			c.newLocation(span),
		)
	}
	local.initialised = true
	c.emitSetLocal(span.StartPos.Line, local.index)
}

// Return the offset of the bytecode next instruction.
func (c *compiler) nextInstructionOffset() int {
	return len(c.function.Instructions)
}

func (c *compiler) setLocal(name string, valueNode ast.ExpressionNode, span *position.Span) {
	c.compileNode(valueNode)
	local, ok := c.resolveLocal(name, span)
	if !ok {
		return
	}
	if local.initialised && local.singleAssignment {
		c.errors.Add(
			fmt.Sprintf("can't reassign a val: %s", name),
			c.newLocation(span),
		)
	}
	local.initialised = true
	c.emitSetLocal(span.StartPos.Line, local.index)
}

func (c *compiler) localVariableAssignment(name string, operator *token.Token, right ast.ExpressionNode, span *position.Span) {
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
		c.emitSetLocal(span.StartPos.Line, local.index)
	default:
		c.errors.Add(
			fmt.Sprintf("assignment using this operator has not been implemented: %s", operator.Type.String()),
			c.newLocation(span),
		)
		return
	}
}

func (c *compiler) localVariableAccess(name string, span *position.Span) (*local, bool) {
	local, ok := c.resolveLocal(name, span)
	if !ok {
		return nil, false
	}
	if !local.initialised {
		c.errors.Add(
			fmt.Sprintf("can't access an uninitialised local: %s", name),
			c.newLocation(span),
		)
		return nil, false
	}

	c.emitGetLocal(span.StartPos.Line, local.index)
	return local, true
}

func (c *compiler) modifierExpression(node *ast.ModifierNode) {
	switch node.Modifier.Type {
	case token.IF:
		c.modifierIfExpression(false, node.Right, node.Left, nil, node.Span())
	case token.UNLESS:
		c.modifierIfExpression(true, node.Right, node.Left, nil, node.Span())
	default:
		c.errors.Add(
			fmt.Sprintf("illegal modifier: %s", node.Modifier.StringValue()),
			c.newLocation(node.Span()),
		)
	}
}

func (c *compiler) modifierIfExpression(unless bool, condition, then, els ast.ExpressionNode, span *position.Span) {
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

func (c *compiler) ifExpression(unless bool, condition ast.ExpressionNode, then, els []ast.StatementNode, span *position.Span) {
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

func (c *compiler) valueDeclaration(node *ast.ValueDeclarationNode) {
	initialised := node.Initialiser != nil

	switch node.Name.Type {
	case token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER:
		if initialised {
			c.compileNode(node.Initialiser)
		}
		local := c.defineLocal(node.Name.StringValue(), node.Span(), true, initialised)
		if initialised {
			c.emitSetLocal(node.Span().StartPos.Line, local.index)
		} else {
			c.emit(node.Span().StartPos.Line, bytecode.NIL)
		}
	default:
		c.errors.Add(
			fmt.Sprintf("can't compile a value declaration with: %s", node.Name.Type.String()),
			c.newLocation(node.Name.Span()),
		)
	}
}

func (c *compiler) variableDeclaration(node *ast.VariableDeclarationNode) {
	initialised := node.Initialiser != nil

	switch node.Name.Type {
	case token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER:
		if initialised {
			c.compileNode(node.Initialiser)
		}
		local := c.defineLocal(node.Name.StringValue(), node.Span(), false, initialised)
		if initialised {
			c.emitSetLocal(node.Span().StartPos.Line, local.index)
		} else {
			c.emit(node.Span().StartPos.Line, bytecode.NIL)
		}
	default:
		c.errors.Add(
			fmt.Sprintf("can't compile a variable declaration with: %s", node.Name.Type.String()),
			c.newLocation(node.Name.Span()),
		)
	}
}

// Compile each element of a collection of statements.
func (c *compiler) compileStatements(collection []ast.StatementNode, span *position.Span) {
	if !c.compileStatementsOk(collection, span) {
		c.emit(span.EndPos.Line, bytecode.NIL)
	}
}

// Same as [compileStatements] but returns false when no instructions were emitted instead
// emitting a `nil` value.
func (c *compiler) compileStatementsOk(collection []ast.StatementNode, span *position.Span) bool {
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

func (c *compiler) intLiteral(node *ast.IntLiteralNode) {
	i, err := value.ParseBigInt(node.Value, 0)
	if err != nil {
		c.errors.Add(err.Error(), c.newLocation(node.Span()))
		return
	}
	if i.IsSmallInt() {
		c.emitConstant(i.ToSmallInt(), node.Span())
		return
	}
	c.emitConstant(i, node.Span())
}

// Compiles boolean binary operators
func (c *compiler) logicalExpression(node *ast.LogicalExpressionNode) {
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
		c.errors.Add(fmt.Sprintf("unknown logical operator: %s", node.Op.String()), c.newLocation(node.Span()))
	}
}

// Compiles the `??` operator
func (c *compiler) nilCoalescing(node *ast.LogicalExpressionNode) {
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
func (c *compiler) logicalOr(node *ast.LogicalExpressionNode) {
	c.compileNode(node.Left)
	jump := c.emitJump(node.Span().StartPos.Line, bytecode.JUMP_IF)

	// if falsy
	c.emit(node.Span().StartPos.Line, bytecode.POP)
	c.compileNode(node.Right)

	// if truthy
	c.patchJump(jump, node.Span())
}

// Compiles the `&&` operator
func (c *compiler) logicalAnd(node *ast.LogicalExpressionNode) {
	c.compileNode(node.Left)
	jump := c.emitJump(node.Span().StartPos.Line, bytecode.JUMP_UNLESS)

	// if truthy
	c.emit(node.Span().StartPos.Line, bytecode.POP)
	c.compileNode(node.Right)

	// if falsy
	c.patchJump(jump, node.Span())
}

func (c *compiler) binaryExpression(node *ast.BinaryExpressionNode) {
	if c.resolveAndEmit(node) {
		return
	}
	c.compileNode(node.Left)
	c.compileNode(node.Right)
	switch node.Op.Type {
	case token.PLUS:
		c.emit(node.Span().StartPos.Line, bytecode.ADD)
	case token.MINUS:
		c.emit(node.Span().StartPos.Line, bytecode.SUBTRACT)
	case token.STAR:
		c.emit(node.Span().StartPos.Line, bytecode.MULTIPLY)
	case token.SLASH:
		c.emit(node.Span().StartPos.Line, bytecode.DIVIDE)
	case token.STAR_STAR:
		c.emit(node.Span().StartPos.Line, bytecode.EXPONENTIATE)
	case token.LBITSHIFT:
		c.emit(node.Span().StartPos.Line, bytecode.LBITSHIFT)
	case token.LTRIPLE_BITSHIFT:
		c.emit(node.Span().StartPos.Line, bytecode.LOGIC_LBITSHIFT)
	case token.RBITSHIFT:
		c.emit(node.Span().StartPos.Line, bytecode.RBITSHIFT)
	case token.RTRIPLE_BITSHIFT:
		c.emit(node.Span().StartPos.Line, bytecode.LOGIC_RBITSHIFT)
	case token.AND:
		c.emit(node.Span().StartPos.Line, bytecode.BITWISE_AND)
	case token.OR:
		c.emit(node.Span().StartPos.Line, bytecode.BITWISE_OR)
	case token.XOR:
		c.emit(node.Span().StartPos.Line, bytecode.BITWISE_XOR)
	case token.PERCENT:
		c.emit(node.Span().StartPos.Line, bytecode.MODULO)
	case token.EQUAL_EQUAL:
		c.emit(node.Span().StartPos.Line, bytecode.EQUAL)
	case token.NOT_EQUAL:
		c.emit(node.Span().StartPos.Line, bytecode.NOT_EQUAL)
	case token.STRICT_EQUAL:
		c.emit(node.Span().StartPos.Line, bytecode.STRICT_EQUAL)
	case token.STRICT_NOT_EQUAL:
		c.emit(node.Span().StartPos.Line, bytecode.STRICT_NOT_EQUAL)
	case token.GREATER:
		c.emit(node.Span().StartPos.Line, bytecode.GREATER)
	case token.GREATER_EQUAL:
		c.emit(node.Span().StartPos.Line, bytecode.GREATER_EQUAL)
	case token.LESS:
		c.emit(node.Span().StartPos.Line, bytecode.LESS)
	case token.LESS_EQUAL:
		c.emit(node.Span().StartPos.Line, bytecode.LESS_EQUAL)
	default:
		c.errors.Add(fmt.Sprintf("unknown binary operator: %s", node.Op.String()), c.newLocation(node.Span()))
	}
}

// Resolves static AST expressions to Elk values
// and emits bytecode that loads them.
// Returns false when the node can't be optimised at compile-time
// and no bytecode has been generated.
func (c *compiler) resolveAndEmit(node ast.ExpressionNode) bool {
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
		c.emitConstant(result, node.Span())
	}
	return true
}

func (c *compiler) unaryExpression(node *ast.UnaryExpressionNode) {
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
	default:
		c.errors.Add(fmt.Sprintf("unknown unary operator: %s", node.Op.String()), c.newLocation(node.Span()))
	}
}

// Emit an instruction that jumps forward with a placeholder offset.
// Returns the offset of placeholder value that has to be patched.
func (c *compiler) emitJump(line int, op bytecode.OpCode) int {
	c.emit(line, op, 0xff, 0xff)
	return c.nextInstructionOffset() - 2
}

// Emit an instruction that jumps back to the given bytecode offset.
func (c *compiler) emitLoop(span *position.Span, startOffset int) {
	c.emit(span.EndPos.Line, bytecode.LOOP)

	offset := c.nextInstructionOffset() - startOffset + 2
	if offset > math.MaxUint16 {
		c.errors.Add(
			fmt.Sprintf("too many bytes to jumbytep backward: %d", math.MaxUint16),
			c.newLocation(span),
		)
	}

	c.function.AppendUint16(uint16(offset))
}

// Overwrite the placeholder operand of a jump instruction
func (c *compiler) patchJump(offset int, span *position.Span) {
	jump := c.nextInstructionOffset() - offset - 2

	if jump > math.MaxUint16 {
		c.errors.Add(
			fmt.Sprintf("too many bytes to jump over: %d", jump),
			c.newLocation(span),
		)
		return
	}

	c.function.Instructions[offset] = byte((jump >> 8) & 0xff)
	c.function.Instructions[offset+1] = byte(jump & 0xff)
}

// Emit an instruction that sets a local variable or value.
func (c *compiler) emitSetLocal(line int, index uint16) {
	if index > math.MaxUint8 {
		c.emit(line, bytecode.SET_LOCAL16)
		c.function.AppendUint16(index)
		return
	}

	c.emit(line, bytecode.SET_LOCAL8, byte(index))
}

// Emit an instruction that gets the value of a local.
func (c *compiler) emitGetLocal(line int, index uint16) {
	if index > math.MaxUint8 {
		c.emit(line, bytecode.GET_LOCAL16)
		c.function.AppendUint16(index)
		return
	}

	c.emit(line, bytecode.GET_LOCAL8, byte(index))
}

// Add a constant to the constant pool and emit appropriate bytecode.
func (c *compiler) emitConstant(val value.Value, span *position.Span) {
	id, size := c.function.AddConstant(val)
	switch size {
	case bytecode.UINT8_SIZE:
		c.function.AddInstruction(span.StartPos.Line, bytecode.CONSTANT8, byte(id))
	case bytecode.UINT16_SIZE:
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, uint16(id))
		c.function.AddInstruction(span.StartPos.Line, bytecode.CONSTANT16, bytes...)
	case bytecode.UINT32_SIZE:
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, uint32(id))
		c.function.AddInstruction(span.StartPos.Line, bytecode.CONSTANT32, bytes...)
	default:
		c.errors.Add(
			fmt.Sprintf("constant pool limit reached: %d", math.MaxUint32),
			c.newLocation(span),
		)
	}
}

// Emit an instruction that gets the value of a module constant.
func (c *compiler) emitGetModConst(name value.Symbol, span *position.Span) {
	id, size := c.function.AddConstant(name)
	switch size {
	case bytecode.UINT8_SIZE:
		c.function.AddInstruction(span.StartPos.Line, bytecode.GET_MOD_CONST8, byte(id))
	case bytecode.UINT16_SIZE:
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, uint16(id))
		c.function.AddInstruction(span.StartPos.Line, bytecode.GET_MOD_CONST16, bytes...)
	case bytecode.UINT32_SIZE:
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, uint32(id))
		c.function.AddInstruction(span.StartPos.Line, bytecode.GET_MOD_CONST32, bytes...)
	default:
		c.errors.Add(
			fmt.Sprintf("constant pool limit reached: %d", math.MaxUint32),
			c.newLocation(span),
		)
	}
}

// Emit an opcode with optional bytes.
func (c *compiler) emit(line int, op bytecode.OpCode, bytes ...byte) {
	c.function.AddInstruction(line, op, bytes...)
}

func (c *compiler) enterScope() {
	c.scopes = append(c.scopes, localTable{})
}

func (c *compiler) leaveScope(line int) {
	currentDepth := len(c.scopes) - 1

	varsToPop := len(c.scopes[currentDepth])
	if varsToPop > 0 {
		if c.lastLocalIndex > math.MaxUint8 || varsToPop > math.MaxUint8 {
			c.emit(line, bytecode.LEAVE_SCOPE32)
			c.function.AppendUint16(uint16(c.lastLocalIndex))
			c.function.AppendUint16(uint16(varsToPop))
		} else {
			c.emit(line, bytecode.LEAVE_SCOPE16, byte(c.lastLocalIndex), byte(varsToPop))
		}
	}

	c.lastLocalIndex -= varsToPop
	c.scopes[currentDepth] = nil
	c.scopes = c.scopes[:currentDepth]
}

// Register a local variable.
func (c *compiler) defineLocal(name string, span *position.Span, singleAssignment, initialised bool) *local {
	varScope := c.scopes.last()
	_, ok := varScope[name]
	if ok {
		c.errors.Add(
			fmt.Sprintf("a variable with this name has already been declared in this scope: %s", name),
			c.newLocation(span),
		)
	}
	if c.lastLocalIndex == math.MaxUint16 {
		c.errors.Add(
			fmt.Sprintf("exceeded the maximum number of local variables (%d): %s", math.MaxUint16, name),
			c.newLocation(span),
		)
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
func (c *compiler) resolveLocal(name string, span *position.Span) (*local, bool) {
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
		c.errors.Add(
			fmt.Sprintf("undeclared variable: %s", name),
			c.newLocation(span),
		)
		return localVal, false
	}

	return localVal, true
}
