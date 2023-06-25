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
	"github.com/elk-language/elk/token"
)

// Compile the Elk source to a bytecode chunk.
func CompileSource(sourceName string, source []byte) (*bytecode.Chunk, position.ErrorList) {
	ast, err := parser.Parse(sourceName, source)
	if err != nil {
		return nil, err
	}

	return CompileAST(sourceName, ast)
}

// Compile the AST node to a bytecode chunk.
func CompileAST(sourceName string, ast ast.Node) (*bytecode.Chunk, position.ErrorList) {
	compiler := new(
		sourceName,
		position.NewLocationWithPosition(sourceName, ast.Pos()),
	)
	compiler.compile(ast)

	return compiler.bytecode, compiler.errors
}

// Holds the state of the compiler.
type compiler struct {
	sourceName string
	bytecode   *bytecode.Chunk
	errors     position.ErrorList
}

// Instantiate a new compiler instance.
func new(sourceName string, loc *position.Location) *compiler {
	return &compiler{
		bytecode: bytecode.NewChunk(
			[]byte{},
			loc,
		),
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
		for _, s := range node.Body {
			if !c.compile(s) {
				return false
			}
		}
	case *ast.ExpressionStatementNode:
		return c.compile(node.Expression)
	case *ast.BinaryExpressionNode:
		return c.binaryExpression(node)
	case *ast.UnaryExpressionNode:
		return c.unaryExpression(node)
	case *ast.IntLiteralNode:
		// TODO: Implement BigInt compilation
		i, err := object.StrictParseInt(node.Value, 0, 64)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.SmallInt(i), node)
	case *ast.Int8LiteralNode:
		i, err := object.StrictParseInt(node.Value, 0, 8)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.Int8(i), node)
	case *ast.Int16LiteralNode:
		i, err := object.StrictParseInt(node.Value, 0, 16)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.Int16(i), node)
	case *ast.Int32LiteralNode:
		i, err := object.StrictParseInt(node.Value, 0, 32)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.Int32(i), node)
	case *ast.Int64LiteralNode:
		i, err := object.StrictParseInt(node.Value, 0, 64)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.Int64(i), node)
	case *ast.UInt8LiteralNode:
		i, err := object.StrictParseUint(node.Value, 0, 8)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.UInt8(i), node)
	case *ast.UInt16LiteralNode:
		i, err := object.StrictParseUint(node.Value, 0, 16)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.UInt16(i), node)
	case *ast.UInt32LiteralNode:
		i, err := object.StrictParseUint(node.Value, 0, 32)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.UInt32(i), node)
	case *ast.UInt64LiteralNode:
		i, err := object.StrictParseUint(node.Value, 0, 64)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.UInt64(i), node)
	case *ast.FloatLiteralNode:
		f, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.Float(f), node)
	case *ast.Float64LiteralNode:
		f, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.Float64(f), node)
	case *ast.Float32LiteralNode:
		f, err := strconv.ParseFloat(node.Value, 32)
		if err != nil {
			c.errors.Add(err.Error(), c.newLocation(node.Position))
			return false
		}
		return c.emitConstant(object.Float32(f), node)
	}

	return true
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
		c.emit(node, bytecode.ADD)
	case token.MINUS:
		c.emit(node, bytecode.SUBTRACT)
	case token.STAR:
		c.emit(node, bytecode.MULTIPLY)
	case token.SLASH:
		c.emit(node, bytecode.DIVIDE)
	case token.STAR_STAR:
		c.emit(node, bytecode.EXPONENTIATE)
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
		c.emit(node, bytecode.NEGATE)
	case token.BANG:
		// logical not
		c.emit(node, bytecode.NOT)
	case token.TILDE:
		// binary negation
		c.emit(node, bytecode.BITWISE_NOT)
	default:
		c.errors.Add(fmt.Sprintf("unknown unary operator: %s", node.Op.String()), c.newLocation(node.Position))
		return false
	}

	return true
}

// Add a constant to the constant pool and emit appropriate bytecode.
func (c *compiler) emitConstant(val object.Value, node ast.Node) bool {
	id, size := c.bytecode.AddConstant(val)
	switch size {
	case bytecode.UINT8_SIZE:
		c.bytecode.AddInstruction(node.Pos().Line, bytecode.CONSTANT8, byte(id))
	case bytecode.UINT16_SIZE:
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, uint16(id))
		c.bytecode.AddInstruction(node.Pos().Line, bytecode.CONSTANT16, bytes...)
	case bytecode.UINT32_SIZE:
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, uint32(id))
		c.bytecode.AddInstruction(node.Pos().Line, bytecode.CONSTANT32, bytes...)
	default:
		c.errors.Add(
			fmt.Sprintf("constant pool limit reached: %d", math.MaxUint32),
			c.newLocation(node.Pos()),
		)
		return false
	}
	return true
}

// Add a constant to the constant pool and emit appropriate bytecode.
func (c *compiler) emit(node ast.Node, op bytecode.OpCode, bytes ...byte) {
	c.bytecode.AddInstruction(node.Pos().Line, op, bytes...)
}
