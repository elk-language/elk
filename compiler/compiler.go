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
	"github.com/elk-language/elk/object"
	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/errors"

	"github.com/elk-language/elk/token"
)

// Compile the Elk source to a bytecode chunk.
func CompileSource(sourceName string, source string) (*bytecode.Chunk, errors.ErrorList) {
	ast, err := parser.Parse(sourceName, source)
	if err != nil {
		return nil, err
	}

	return CompileAST(sourceName, ast)
}

// Compile the AST node to a bytecode chunk.
func CompileAST(sourceName string, ast ast.Node) (*bytecode.Chunk, errors.ErrorList) {
	compiler := new(
		sourceName,
		position.NewLocationWithPosition(sourceName, ast.Pos()),
	)
	compiler.compile(ast)

	return compiler.bytecode, compiler.errors
}

// represents a local variable or value
type local struct {
	index            int16
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
	bytecode       *bytecode.Chunk
	errors         errors.ErrorList
	scopes         scopes
	lastLocalIndex int16 // index of the last local variable
}

// Instantiate a new compiler instance.
func new(sourceName string, loc *position.Location) *compiler {
	return &compiler{
		bytecode: bytecode.NewChunk(
			[]byte{},
			loc,
		),
		scopes:         scopes{localTable{}}, // start with an empty set for the 0th scope
		lastLocalIndex: -1,
		sourceName:     sourceName,
	}
}

// Create a new location struct with the given position.
func (c *compiler) newLocation(pos *position.Position) *position.Location {
	return position.NewLocationWithPosition(c.sourceName, pos)
}

func (c *compiler) compile(node ast.Node) {
	switch node := node.(type) {
	case *ast.ProgramNode:
		c.compileStatements(node.Body, node.Position)
	case *ast.ExpressionStatementNode:
		c.compile(node.Expression)
	case *ast.VariableDeclarationStatementNode:
		c.variableDeclaration(node)
	case *ast.ShortVariableDeclarationStatementNode:
		c.shortVariableDeclaration(node)
	case *ast.AssignmentExpressionNode:
		c.assignment(node)
	case *ast.PublicIdentifierNode:
		c.localVariableAccess(node.Value, node.Position)
	case *ast.PrivateIdentifierNode:
		c.localVariableAccess(node.Value, node.Position)
	case *ast.BinaryExpressionNode:
		c.binaryExpression(node)
	case *ast.UnaryExpressionNode:
		c.unaryExpression(node)
	case *ast.RawStringLiteralNode:
		c.emitConstant(object.String(node.Value), node.Position)
	case *ast.DoubleQuotedStringLiteralNode:
		c.emitConstant(object.String(node.Value), node.Position)
	case *ast.CharLiteralNode:
		c.emitConstant(object.Char(node.Value), node.Position)
	case *ast.RawCharLiteralNode:
		c.emitConstant(object.Char(node.Value), node.Position)
	case *ast.FalseLiteralNode:
		c.emit(node.Position.Line, bytecode.FALSE)
	case *ast.TrueLiteralNode:
		c.emit(node.Position.Line, bytecode.TRUE)
	case *ast.NilLiteralNode:
		c.emit(node.Position.Line, bytecode.NIL)
	case *ast.EmptyStatementNode:
	case *ast.DoExpressionNode:
		c.enterScope()
		c.compileStatements(node.Body, node.Position)
		c.leaveScope(node.Position.Line)
	case *ast.SimpleSymbolLiteralNode:
		c.emitConstant(object.SymbolTable.Add(node.Content), node.Position)
	case *ast.IntLiteralNode:
		c.intLiteral(node)
	case *ast.Int8LiteralNode:
		i, err := object.StrictParseInt(node.Value, 0, 8)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return
		}
		// BENCHMARK: Compare with storing
		// ints inline in bytecode instead of as constants.
		c.emitConstant(object.Int8(i), node.Position)
	case *ast.Int16LiteralNode:
		i, err := object.StrictParseInt(node.Value, 0, 16)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return
		}
		c.emitConstant(object.Int16(i), node.Position)
	case *ast.Int32LiteralNode:
		i, err := object.StrictParseInt(node.Value, 0, 32)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return
		}
		c.emitConstant(object.Int32(i), node.Position)
	case *ast.Int64LiteralNode:
		i, err := object.StrictParseInt(node.Value, 0, 64)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return
		}
		c.emitConstant(object.Int64(i), node.Position)
	case *ast.UInt8LiteralNode:
		i, err := object.StrictParseUint(node.Value, 0, 8)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return
		}
		c.emitConstant(object.UInt8(i), node.Position)
	case *ast.UInt16LiteralNode:
		i, err := object.StrictParseUint(node.Value, 0, 16)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return
		}
		c.emitConstant(object.UInt16(i), node.Position)
	case *ast.UInt32LiteralNode:
		i, err := object.StrictParseUint(node.Value, 0, 32)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return
		}
		c.emitConstant(object.UInt32(i), node.Position)
	case *ast.UInt64LiteralNode:
		i, err := object.StrictParseUint(node.Value, 0, 64)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return
		}
		c.emitConstant(object.UInt64(i), node.Position)
	case *ast.FloatLiteralNode:
		f, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return
		}
		c.emitConstant(object.Float(f), node.Position)
	case *ast.BigFloatLiteralNode:
		f, err := object.ParseBigFloat(node.Value)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return
		}
		c.emitConstant(f, node.Position)
	case *ast.Float64LiteralNode:
		f, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return
		}
		c.emitConstant(object.Float64(f), node.Position)
	case *ast.Float32LiteralNode:
		f, err := strconv.ParseFloat(node.Value, 32)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return
		}
		c.emitConstant(object.Float32(f), node.Position)

	default:
		c.errors.Add(
			fmt.Sprintf("compilation of this node has not been implemented: %T", node),
			c.newLocation(node.Pos()),
		)
	}
}

func (c *compiler) assignment(node *ast.AssignmentExpressionNode) {
	switch n := node.Left.(type) {
	case *ast.PublicIdentifierNode:
		c.localVariableAssignment(n.Value, node.Op, node.Right, node.Position)
	case *ast.PrivateIdentifierNode:
		c.localVariableAssignment(n.Value, node.Op, node.Right, node.Position)
	default:
		c.errors.Add(
			fmt.Sprintf("can't assign to: %T", node),
			c.newLocation(node.Pos()),
		)
	}
}

func (c *compiler) localVariableAssignment(name string, operator *token.Token, right ast.ExpressionNode, pos *position.Position) {
	c.compile(right)
	local, ok := c.resolveLocalVar(name, pos)
	if !ok {
		return
	}
	if local.initialised && local.singleAssignment {
		c.errors.Add(
			fmt.Sprintf("can't reassign a val: %s", name),
			c.newLocation(pos),
		)
	}
	local.initialised = true

	c.emit(pos.Line, bytecode.SET_LOCAL, byte(local.index))
}

func (c *compiler) localVariableAccess(name string, pos *position.Position) {
	local, ok := c.resolveLocalVar(name, pos)
	if !ok {
		return
	}
	if !local.initialised {
		c.errors.Add(
			fmt.Sprintf("can't access an uninitialised local: %s", name),
			c.newLocation(pos),
		)
		return
	}

	c.emit(pos.Line, bytecode.GET_LOCAL, byte(local.index))
}

func (c *compiler) shortVariableDeclaration(node *ast.ShortVariableDeclarationStatementNode) {
	c.compile(node.Initialiser)
	switch node.Name.Type {
	case token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER:
		c.defineLocal(node.Name.StringValue(), node.Position, false, true)
	default:
		c.errors.Add(
			fmt.Sprintf("can't compile a short variable declaration with: %s", node.Name.Type.String()),
			c.newLocation(node.Name.Pos()),
		)
	}
}

func (c *compiler) variableDeclaration(node *ast.VariableDeclarationStatementNode) {
	initialised := node.Initialiser != nil
	if initialised {
		c.compile(node.Initialiser)
	} else {
		// populate the variable slot with `nil` as a placeholder
		c.emit(node.Position.Line, bytecode.NIL)
	}

	switch node.Name.Type {
	case token.PUBLIC_IDENTIFIER, token.PRIVATE_IDENTIFIER:
		c.defineLocal(node.Name.StringValue(), node.Position, false, initialised)
	default:
		c.errors.Add(
			fmt.Sprintf("can't compile a variable declaration with: %s", node.Name.Type.String()),
			c.newLocation(node.Name.Pos()),
		)
	}
}

// Compile each element of a collection of statements.
func (c *compiler) compileStatements(collection []ast.StatementNode, pos *position.Position) {
	var nonEmptyStatements int
	for i, s := range collection {
		if _, ok := s.(*ast.EmptyStatementNode); ok {
			continue
		}
		c.compile(s)
		nonEmptyStatements++
		switch s.(type) {
		case *ast.VariableDeclarationStatementNode, *ast.ShortVariableDeclarationStatementNode:
			continue
		}
		if i != len(collection)-1 {
			c.emit(s.Pos().Line, bytecode.POP)
		}
	}

	if nonEmptyStatements == 0 {
		c.emit(pos.Line, bytecode.NIL)
	}
}

func (c *compiler) intLiteral(node *ast.IntLiteralNode) {
	i, err := object.ParseBigInt(node.Value, 0)
	if err != nil {
		c.errors.Add(err.Error(), c.newLocation(node.Position))
		return
	}
	if i.IsSmallInt() {
		c.emitConstant(i.ToSmallInt(), node.Position)
		return
	}
	c.emitConstant(i, node.Position)
}

func (c *compiler) binaryExpression(node *ast.BinaryExpressionNode) {
	c.compile(node.Left)
	c.compile(node.Right)
	switch node.Op.Type {
	case token.PLUS:
		c.emit(node.Position.Line, bytecode.ADD)
	case token.MINUS:
		c.emit(node.Position.Line, bytecode.SUBTRACT)
	case token.STAR:
		c.emit(node.Position.Line, bytecode.MULTIPLY)
	case token.SLASH:
		c.emit(node.Position.Line, bytecode.DIVIDE)
	case token.STAR_STAR:
		c.emit(node.Position.Line, bytecode.EXPONENTIATE)
	default:
		c.errors.Add(fmt.Sprintf("unknown binary operator: %s", node.Op.String()), c.newLocation(node.Position))
	}
}

func (c *compiler) unaryExpression(node *ast.UnaryExpressionNode) {
	c.compile(node.Right)
	switch node.Op.Type {
	case token.PLUS:
		// TODO: Implement unary plus
	case token.MINUS:
		c.emit(node.Position.Line, bytecode.NEGATE)
	case token.BANG:
		// logical not
		c.emit(node.Position.Line, bytecode.NOT)
	case token.TILDE:
		// binary negation
		c.emit(node.Position.Line, bytecode.BITWISE_NOT)
	default:
		c.errors.Add(fmt.Sprintf("unknown unary operator: %s", node.Op.String()), c.newLocation(node.Position))
	}
}

// Add a constant to the constant pool and emit appropriate bytecode.
func (c *compiler) emitConstant(val object.Value, pos *position.Position) {
	id, size := c.bytecode.AddConstant(val)
	switch size {
	case bytecode.UINT8_SIZE:
		c.bytecode.AddInstruction(pos.Line, bytecode.CONSTANT8, byte(id))
	case bytecode.UINT16_SIZE:
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, uint16(id))
		c.bytecode.AddInstruction(pos.Line, bytecode.CONSTANT16, bytes...)
	case bytecode.UINT32_SIZE:
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, uint32(id))
		c.bytecode.AddInstruction(pos.Line, bytecode.CONSTANT32, bytes...)
	default:
		c.errors.Add(
			fmt.Sprintf("constant pool limit reached: %d", math.MaxUint32),
			c.newLocation(pos),
		)
	}
}

// Emit an opcode with optional bytes.
func (c *compiler) emit(line int, op bytecode.OpCode, bytes ...byte) {
	c.bytecode.AddInstruction(line, op, bytes...)
}

func (c *compiler) enterScope() {
	c.scopes = append(c.scopes, localTable{})
}

func (c *compiler) leaveScope(line int) {
	currentDepth := len(c.scopes) - 1

	varsToPop := len(c.scopes[currentDepth])
	// TODO: fix
	if varsToPop > 0 {
		c.emit(line, bytecode.POP_N, byte(varsToPop))
	}

	c.lastLocalIndex -= int16(varsToPop)
	c.scopes[currentDepth] = nil
	c.scopes = c.scopes[:currentDepth]
}

// Register a local variable.
func (c *compiler) defineLocal(name string, pos *position.Position, singleAssignment, initialised bool) {
	varScope := c.scopes.last()
	_, ok := varScope[name]
	if ok {
		c.errors.Add(
			fmt.Sprintf("a variable with this name has already been declared in this scope: %s", name),
			c.newLocation(pos),
		)
	}
	if c.lastLocalIndex == math.MaxUint8 {
		c.errors.Add(
			fmt.Sprintf("exceeded the maximum number of local variables (%d): %s", math.MaxUint8, name),
			c.newLocation(pos),
		)
	}

	c.lastLocalIndex++
	varScope[name] = &local{
		index:            c.lastLocalIndex,
		initialised:      initialised,
		singleAssignment: singleAssignment,
	}
}

// Resolve a local variable and get its index.
func (c *compiler) resolveLocalVar(name string, pos *position.Position) (*local, bool) {
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
			c.newLocation(pos),
		)
		return localVal, false
	}

	return localVal, true
}
