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

// BENCHMARK: compare with a dynamically allocated array
const MAX_LOCAL_COUNT = math.MaxUint8 + 1

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

// set of local variable names
type varNameSet map[string]bool

// indices represent scope depths
// and elements are sets of local variable names in a particular scope
type scopes []varNameSet

// Get the last local variable scope.
func (s scopes) last() varNameSet {
	return s[len(s)-1]
}

// Holds the state of the compiler.
type compiler struct {
	sourceName string
	bytecode   *bytecode.Chunk
	errors     errors.ErrorList
	scopes     scopes
}

// Instantiate a new compiler instance.
func new(sourceName string, loc *position.Location) *compiler {
	return &compiler{
		bytecode: bytecode.NewChunk(
			[]byte{},
			loc,
		),
		scopes:     scopes{varNameSet{}}, // start with an empty set for the 0th scope
		sourceName: sourceName,
	}
}

// Create a new location struct with the given position.
func (c *compiler) newLocation(pos *position.Position) *position.Location {
	return position.NewLocationWithPosition(c.sourceName, pos)
}

func (c *compiler) compile(node ast.Node) bool {
	switch node := node.(type) {
	case *ast.ProgramNode:
		return c.compileStatements(node.Body, node.Position)
	case *ast.ExpressionStatementNode:
		return c.compile(node.Expression)
	case *ast.BinaryExpressionNode:
		return c.binaryExpression(node)
	case *ast.UnaryExpressionNode:
		return c.unaryExpression(node)
	case *ast.RawStringLiteralNode:
		return c.emitConstant(object.String(node.Value), node.Position)
	case *ast.DoubleQuotedStringLiteralNode:
		return c.emitConstant(object.String(node.Value), node.Position)
	case *ast.CharLiteralNode:
		return c.emitConstant(object.Char(node.Value), node.Position)
	case *ast.RawCharLiteralNode:
		return c.emitConstant(object.Char(node.Value), node.Position)
	case *ast.FalseLiteralNode:
		c.emit(node.Line, bytecode.FALSE)
	case *ast.TrueLiteralNode:
		c.emit(node.Line, bytecode.TRUE)
	case *ast.NilLiteralNode:
		c.emit(node.Line, bytecode.NIL)
	case *ast.EmptyStatementNode:
	case *ast.DoExpressionNode:
		c.enterScope()
		result := c.compileStatements(node.Body, node.Position)
		c.leaveScope(node.Line)
		return result
	case *ast.VariableDeclarationNode:
		return c.addLocalVar(node.Name.StringValue(), node.Position)
	case *ast.SimpleSymbolLiteralNode:
		return c.emitConstant(object.SymbolTable.Add(node.Content), node.Position)
	case *ast.IntLiteralNode:
		return c.intLiteral(node)
	case *ast.Int8LiteralNode:
		i, err := object.StrictParseInt(node.Value, 0, 8)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		// BENCHMARK: Compare with storing
		// ints inline in bytecode instead of as constants.
		return c.emitConstant(object.Int8(i), node.Position)
	case *ast.Int16LiteralNode:
		i, err := object.StrictParseInt(node.Value, 0, 16)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.Int16(i), node.Position)
	case *ast.Int32LiteralNode:
		i, err := object.StrictParseInt(node.Value, 0, 32)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.Int32(i), node.Position)
	case *ast.Int64LiteralNode:
		i, err := object.StrictParseInt(node.Value, 0, 64)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.Int64(i), node.Position)
	case *ast.UInt8LiteralNode:
		i, err := object.StrictParseUint(node.Value, 0, 8)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.UInt8(i), node.Position)
	case *ast.UInt16LiteralNode:
		i, err := object.StrictParseUint(node.Value, 0, 16)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.UInt16(i), node.Position)
	case *ast.UInt32LiteralNode:
		i, err := object.StrictParseUint(node.Value, 0, 32)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.UInt32(i), node.Position)
	case *ast.UInt64LiteralNode:
		i, err := object.StrictParseUint(node.Value, 0, 64)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.UInt64(i), node.Position)
	case *ast.FloatLiteralNode:
		f, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.Float(f), node.Position)
	case *ast.BigFloatLiteralNode:
		f, err := object.ParseBigFloat(node.Value)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(f, node.Position)
	case *ast.Float64LiteralNode:
		f, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.Float64(f), node.Position)
	case *ast.Float32LiteralNode:
		f, err := strconv.ParseFloat(node.Value, 32)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.Float32(f), node.Position)

	default:
		c.errors.Add(
			fmt.Sprintf("compilation of this node has not been implemented: %T", node),
			c.newLocation(node.Pos()),
		)
		return false
	}

	return true
}

// Register a local variable.
func (c *compiler) addLocalVar(name string, pos *position.Position) bool {
	varScope := c.scopes.last()
	varExists := varScope[name]
	if varExists {
		c.errors.Add(
			fmt.Sprintf("a variable with this name has already been declared in this scope: %s", name),
			c.newLocation(pos),
		)
		return false
	}

	varScope[name] = true
	return true
}

// Compile each element of a collection of statements.
func (c *compiler) compileStatements(collection []ast.StatementNode, pos *position.Position) bool {
	if len(collection) == 0 {
		c.emit(pos.Line, bytecode.NIL)
		return true
	}

	for i, s := range collection {
		if !c.compile(s) {
			return false
		}
		if i != len(collection)-1 {
			c.emit(s.Pos().Line, bytecode.POP)
		}
	}

	return true
}

func (c *compiler) intLiteral(node *ast.IntLiteralNode) bool {
	i, err := object.ParseBigInt(node.Value, 0)
	if err != nil {
		c.errors.Add(err.Error(), c.newLocation(node.Position))
		return false
	}
	if i.IsSmallInt() {
		return c.emitConstant(i.ToSmallInt(), node.Position)
	}
	return c.emitConstant(i, node.Position)
}

func (c *compiler) binaryExpression(node *ast.BinaryExpressionNode) bool {
	if !c.compile(node.Left) {
		return false
	}
	if !c.compile(node.Right) {
		return false
	}
	switch node.Op.Type {
	case token.PLUS:
		c.emit(node.Line, bytecode.ADD)
	case token.MINUS:
		c.emit(node.Line, bytecode.SUBTRACT)
	case token.STAR:
		c.emit(node.Line, bytecode.MULTIPLY)
	case token.SLASH:
		c.emit(node.Line, bytecode.DIVIDE)
	case token.STAR_STAR:
		c.emit(node.Line, bytecode.EXPONENTIATE)
	default:
		c.errors.Add(fmt.Sprintf("unknown binary operator: %s", node.Op.String()), c.newLocation(node.Position))
		return false
	}

	return true
}

func (c *compiler) unaryExpression(node *ast.UnaryExpressionNode) bool {
	if !c.compile(node.Right) {
		return false
	}
	switch node.Op.Type {
	case token.PLUS:
		// TODO: Implement unary plus
	case token.MINUS:
		c.emit(node.Line, bytecode.NEGATE)
	case token.BANG:
		// logical not
		c.emit(node.Line, bytecode.NOT)
	case token.TILDE:
		// binary negation
		c.emit(node.Line, bytecode.BITWISE_NOT)
	default:
		c.errors.Add(fmt.Sprintf("unknown unary operator: %s", node.Op.String()), c.newLocation(node.Position))
		return false
	}

	return true
}

// Add a constant to the constant pool and emit appropriate bytecode.
func (c *compiler) emitConstant(val object.Value, pos *position.Position) bool {
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
		return false
	}
	return true
}

// Emit an opcode with optional bytes.
func (c *compiler) emit(line int, op bytecode.OpCode, bytes ...byte) {
	c.bytecode.AddInstruction(line, op, bytes...)
}

func (c *compiler) enterScope() {
	c.scopes = append(c.scopes, varNameSet{})
}

func (c *compiler) leaveScope(line int) {
	currentDepth := len(c.scopes) - 1

	varsToPop := len(c.scopes[currentDepth])
	// TODO: fix
	if varsToPop > 0 {
		c.emit(line, bytecode.POP_N, byte(varsToPop))
	}

	c.scopes[currentDepth] = nil
	c.scopes = c.scopes[:currentDepth]
}
