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

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/value"

	"github.com/elk-language/elk/token"
)

// Compile the Elk source to a Bytecode chunk.
func CompileSource(sourceName string, source string) (*value.BytecodeFunction, errors.ErrorList) {
	ast, err := parser.Parse(sourceName, source)
	if err != nil {
		return nil, err
	}

	return CompileAST(sourceName, ast)
}

// Compile the AST node to a Bytecode chunk.
func CompileAST(sourceName string, ast ast.Node) (*value.BytecodeFunction, errors.ErrorList) {
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

// Holds the state of the Compiler.
type Compiler struct {
	Name             string
	Bytecode         *value.BytecodeFunction
	Errors           errors.ErrorList
	scopes           scopes
	lastLocalIndex   int // index of the last local variable
	maxLocalIndex    int // max index of a local variable
	predefinedLocals int
	mode             mode
}

// Instantiate a new Compiler instance.
func new(name string, mode mode, loc *position.Location) *Compiler {
	c := &Compiler{
		Bytecode: value.NewBytecodeFunction(
			value.SymbolTable.Add(name),
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
	case methodMode:
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
	compiler.predefinedLocals = len(c.scopes.last())
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
	c.emit(node.Span().EndPos.Line, bytecode.RETURN)
}

// Entry point for compiling the body of a method.
func (c *Compiler) compileMethod(node *ast.MethodDefinitionNode) {
	span := node.Span()
	if len(node.Parameters) > 0 {
		c.Bytecode.Parameters = make([]value.Symbol, 0, len(node.Parameters))
	}
	for _, param := range node.Parameters {
		p := param.(*ast.MethodParameterNode)
		if p.SetInstanceVariable {
			c.Errors.Add(
				fmt.Sprintf("instance variable parameters are not supported yet: %s", p.Name),
				c.newLocation(p.Span()),
			)
			continue
		}

		if p.Kind != ast.NormalParameterKind {
			c.Errors.Add(
				fmt.Sprintf("splat parameters are not supported yet: %s", p.Name),
				c.newLocation(p.Span()),
			)
			continue
		}

		if p.Initialiser != nil {
			c.Errors.Add(
				fmt.Sprintf("optional parameters are not supported yet: %s", p.Name),
				c.newLocation(p.Span()),
			)
			continue
		}

		c.defineLocal(p.Name, p.Span(), false, true)
		c.Bytecode.Parameters = append(c.Bytecode.Parameters, value.SymbolTable.Add(p.Name))
		c.predefinedLocals++
	}
	c.compileStatements(node.Body, span)
	c.prepLocals()
	c.emit(span.EndPos.Line, bytecode.RETURN)
}

// Entry point for compiling the body of a Module, Class, Mixin, Struct.
func (c *Compiler) compileModule(node ast.Node) {
	span := node.Span()
	switch n := node.(type) {
	case *ast.ClassDeclarationNode:
		c.compileModuleStatements(n.Body, span)
	case *ast.ModuleDeclarationNode:
		c.compileModuleStatements(n.Body, span)
	case *ast.MixinDeclarationNode:
		c.compileModuleStatements(n.Body, span)
	}
	c.prepLocals()
	c.emit(span.EndPos.Line, bytecode.RETURN)
}

// Compile the top level statements of a module body.
func (c *Compiler) compileModuleStatements(collection []ast.StatementNode, span *position.Span) {
	if c.compileStatementsOk(collection, span) {
		c.emit(span.EndPos.Line, bytecode.POP)
	}
	// return the module itself
	c.emit(span.EndPos.Line, bytecode.SELF)
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
	lineInfo.InstructionCount++
}

func (c *Compiler) compileNode(node ast.Node) {
	switch node := node.(type) {
	case *ast.ProgramNode:
		c.compileStatements(node.Body, node.Span())
	case *ast.ExpressionStatementNode:
		c.compileNode(node.Expression)
	case *ast.ConstantLookupNode:
		c.constantLookup(node)
	case *ast.GenericConstantNode:
		c.compileNode(node.Constant)
	case *ast.SelfLiteralNode:
		c.emit(node.Span().StartPos.Line, bytecode.SELF)
	case *ast.AssignmentExpressionNode:
		c.assignment(node)
	case *ast.ClassDeclarationNode:
		c.classDeclaration(node)
	case *ast.ModuleDeclarationNode:
		c.moduleDeclaration(node)
	case *ast.MixinDeclarationNode:
		c.mixinDeclaration(node)
	case *ast.MethodDefinitionNode:
		c.methodDefinition(node)
	case *ast.IncludeExpressionNode:
		c.includeExpression(node)
	case *ast.ExtendExpressionNode:
		c.extendExpression(node)
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
	case *ast.BinaryExpressionNode:
		c.binaryExpression(node)
	case *ast.LogicalExpressionNode:
		c.logicalExpression(node)
	case *ast.UnaryExpressionNode:
		c.unaryExpression(node)
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
		c.modifierExpression(node)
	case *ast.LoopExpressionNode:
		c.loopExpression(node)
	case *ast.NumericForExpressionNode:
		c.numericForExpression(node)
	case *ast.SimpleSymbolLiteralNode:
		c.emitValue(value.SymbolTable.Add(node.Content), node.Span())
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

func (c *Compiler) loopExpression(node *ast.LoopExpressionNode) {
	c.enterScope()
	defer c.leaveScope(node.Span().EndPos.Line)

	start := c.nextInstructionOffset()
	if c.compileStatementsOk(node.ThenBody, node.Span()) {
		c.emit(node.Span().EndPos.Line, bytecode.POP)
	}
	c.emitLoop(node.Span(), start)
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
		c.emitGetModConst(value.SymbolTable.Add(r.Value), node.Span())
	default:
		c.Errors.Add(
			fmt.Sprintf("incorrect right side of constant lookup: %T", node.Right),
			c.newLocation(node.Span()),
		)
	}
}

// Compile a numeric for loop eg. `for i := 0; i < 5; i += 1 then println(i)`
func (c *Compiler) numericForExpression(node *ast.NumericForExpressionNode) {
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

func (c *Compiler) assignment(node *ast.AssignmentExpressionNode) {
	switch n := node.Left.(type) {
	case *ast.PublicIdentifierNode:
		c.localVariableAssignment(n.Value, node.Op, node.Right, node.Span())
	case *ast.PrivateIdentifierNode:
		c.localVariableAssignment(n.Value, node.Op, node.Right, node.Span())
	case *ast.ConstantLookupNode:
		if node.Op.Type != token.COLON_EQUAL {
			c.Errors.Add(
				fmt.Sprintf("can't assign constants using `%s`", node.Op.StringValue()),
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
			c.emitDefModConst(value.SymbolTable.Add(r.Value), n.Span())
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
	default:
		c.Errors.Add(
			fmt.Sprintf("can't assign to: %T", node.Left),
			c.newLocation(node.Span()),
		)
	}
}

func (c *Compiler) compileSimpleConstantAssignment(name string, op *token.Token, right ast.ExpressionNode, span *position.Span) {
	if op.Type != token.COLON_EQUAL {
		c.Errors.Add(
			fmt.Sprintf("can't assign constants using `%s`", op.StringValue()),
			c.newLocation(span),
		)
	}
	c.compileNode(right)
	c.emit(span.StartPos.Line, bytecode.CONSTANT_CONTAINER)
	c.emitDefModConst(value.SymbolTable.Add(name), span)
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
			fmt.Sprintf("can't reassign a val: %s", name),
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
			fmt.Sprintf("can't reassign a val: %s", name),
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
		c.emitSetLocal(span.StartPos.Line, local.index)
	default:
		c.Errors.Add(
			fmt.Sprintf("assignment using this operator has not been implemented: %s", operator.Type.String()),
			c.newLocation(span),
		)
		return
	}
}

func (c *Compiler) localVariableAccess(name string, span *position.Span) (*local, bool) {
	local, ok := c.resolveLocal(name, span)
	if !ok {
		return nil, false
	}
	if !local.initialised {
		c.Errors.Add(
			fmt.Sprintf("can't access an uninitialised local: %s", name),
			c.newLocation(span),
		)
		return nil, false
	}

	c.emitGetLocal(span.StartPos.Line, local.index)
	return local, true
}

func (c *Compiler) modifierExpression(node *ast.ModifierNode) {
	switch node.Modifier.Type {
	case token.IF:
		c.modifierIfExpression(false, node.Right, node.Left, nil, node.Span())
	case token.UNLESS:
		c.modifierIfExpression(true, node.Right, node.Left, nil, node.Span())
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
		if initialised {
			c.emitSetLocal(node.Span().StartPos.Line, local.index)
		} else {
			c.emit(node.Span().StartPos.Line, bytecode.NIL)
		}
	default:
		c.Errors.Add(
			fmt.Sprintf("can't compile a value declaration with: %s", node.Name.Type.String()),
			c.newLocation(node.Name.Span()),
		)
	}
}

func (c *Compiler) returnExpression(node *ast.ReturnExpressionNode) {
	span := node.Span()
	if node.Value != nil {
		c.compileNode(node.Value)
	} else {
		c.emit(span.StartPos.Line, bytecode.NIL)
	}

	c.emit(span.StartPos.Line, bytecode.RETURN)
}

func (c *Compiler) functionCall(node *ast.FunctionCallNode) {
	for _, posArg := range node.PositionalArguments {
		c.compileNode(posArg)
	}

	if node.NamedArguments != nil {
		c.Errors.Add(
			fmt.Sprintf("named arguments are not supported yet: %s", node.MethodName),
			c.newLocation(node.Span()),
		)
		return
	}

	name := value.SymbolTable.Add(node.MethodName)
	callInfo := value.NewCallSiteInfo(name, len(node.PositionalArguments))
	c.emitCallFunction(callInfo, node.Span())
}

func (c *Compiler) methodCall(node *ast.MethodCallNode) {
	c.compileNode(node.Receiver)
	for _, posArg := range node.PositionalArguments {
		c.compileNode(posArg)
	}

	if node.NamedArguments != nil {
		c.Errors.Add(
			fmt.Sprintf("named arguments are not supported yet: %s", node.MethodName),
			c.newLocation(node.Span()),
		)
		return
	}

	if node.NilSafe {
		c.Errors.Add(
			fmt.Sprintf("nil safe method calls are not supported yet: %s", node.MethodName),
			c.newLocation(node.Span()),
		)
		return
	}

	name := value.SymbolTable.Add(node.MethodName)
	callInfo := value.NewCallSiteInfo(name, len(node.PositionalArguments))
	c.emitCallMethod(callInfo, node.Span())
}

func (c *Compiler) methodDefinition(node *ast.MethodDefinitionNode) {
	if c.mode == methodMode {
		c.Errors.Add(
			fmt.Sprintf("methods can't be nested: %s", node.Name),
			c.newLocation(node.Span()),
		)
		return
	}
	methodCompiler := new(node.Name, methodMode, c.newLocation(node.Span()))
	methodCompiler.Errors = c.Errors
	methodCompiler.compileMethod(node)
	c.Errors = methodCompiler.Errors

	result := methodCompiler.Bytecode
	c.emitValue(result, node.Span())

	c.emitValue(value.SymbolTable.Add(node.Name), node.Span())

	c.emit(node.Span().StartPos.Line, bytecode.DEF_METHOD)
}

func (c *Compiler) extendExpression(node *ast.ExtendExpressionNode) {
	switch c.mode {
	case classMode, mixinMode, moduleMode:
	case topLevelMode:
		c.Errors.Add(
			"can't extend mixins in the top level",
			c.newLocation(node.Span()),
		)
		return
	case methodMode:
		c.Errors.Add(
			"can't extend mixins in a method",
			c.newLocation(node.Span()),
		)
		return
	default:
		c.Errors.Add(
			"can't extend mixins in this context",
			c.newLocation(node.Span()),
		)
		return
	}

	span := node.Span()
	for _, constant := range node.Constants {
		c.compileNode(constant)
		c.emit(span.StartPos.Line, bytecode.SELF)
		c.emit(span.StartPos.Line, bytecode.GET_SINGLETON_CLASS)
		c.emit(span.StartPos.Line, bytecode.INCLUDE)
	}

	c.emit(span.EndPos.Line, bytecode.NIL)
}

func (c *Compiler) includeExpression(node *ast.IncludeExpressionNode) {
	switch c.mode {
	case classMode, mixinMode:
	case topLevelMode:
		c.Errors.Add(
			"can't include mixins in the top level",
			c.newLocation(node.Span()),
		)
		return
	case moduleMode:
		c.Errors.Add(
			"can't include mixins in a module",
			c.newLocation(node.Span()),
		)
		return
	case methodMode:
		c.Errors.Add(
			"can't include mixins in a method",
			c.newLocation(node.Span()),
		)
		return
	default:
		c.Errors.Add(
			"can't include mixins in this context",
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
			c.emitValue(value.SymbolTable.Add(r.Value), r.Span())
		case *ast.PrivateConstantNode:
			c.emitValue(value.SymbolTable.Add(r.Value), r.Span())
		default:
			c.Errors.Add(
				fmt.Sprintf("incorrect right side of constant lookup: %T", constant.Right),
				c.newLocation(constant.Right.Span()),
			)
			return
		}
	case *ast.PublicConstantNode:
		c.emit(constant.Span().StartPos.Line, bytecode.CONSTANT_CONTAINER)
		c.emitValue(value.SymbolTable.Add(constant.Value), constant.Span())
	case *ast.PrivateConstantNode:
		c.emit(constant.Span().StartPos.Line, bytecode.CONSTANT_CONTAINER)
		c.emitValue(value.SymbolTable.Add(constant.Value), constant.Span())
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
			c.emitValue(value.SymbolTable.Add(r.Value), r.Span())
		case *ast.PrivateConstantNode:
			c.emitValue(value.SymbolTable.Add(r.Value), r.Span())
		default:
			c.Errors.Add(
				fmt.Sprintf("incorrect right side of constant lookup: %T", constant.Right),
				c.newLocation(constant.Right.Span()),
			)
			return
		}
	case *ast.PublicConstantNode:
		c.emit(constant.Span().StartPos.Line, bytecode.CONSTANT_CONTAINER)
		c.emitValue(value.SymbolTable.Add(constant.Value), constant.Span())
	case *ast.PrivateConstantNode:
		c.emit(constant.Span().StartPos.Line, bytecode.CONSTANT_CONTAINER)
		c.emitValue(value.SymbolTable.Add(constant.Value), constant.Span())
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

func (c *Compiler) classDeclaration(node *ast.ClassDeclarationNode) {
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
			c.emitValue(value.SymbolTable.Add(r.Value), r.Span())
		case *ast.PrivateConstantNode:
			c.emitValue(value.SymbolTable.Add(r.Value), r.Span())
		default:
			c.Errors.Add(
				fmt.Sprintf("incorrect right side of constant lookup: %T", constant.Right),
				c.newLocation(constant.Right.Span()),
			)
			return
		}
	case *ast.PublicConstantNode:
		c.emit(constant.Span().StartPos.Line, bytecode.CONSTANT_CONTAINER)
		c.emitValue(value.SymbolTable.Add(constant.Value), constant.Span())
	case *ast.PrivateConstantNode:
		c.emit(constant.Span().StartPos.Line, bytecode.CONSTANT_CONTAINER)
		c.emitValue(value.SymbolTable.Add(constant.Value), constant.Span())
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

	c.emit(node.Span().StartPos.Line, bytecode.DEF_CLASS)
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
		if initialised {
			c.emitSetLocal(node.Span().StartPos.Line, local.index)
		} else {
			c.emit(node.Span().StartPos.Line, bytecode.NIL)
		}
	default:
		c.Errors.Add(
			fmt.Sprintf("can't compile a variable declaration with: %s", node.Name.Type.String()),
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
		c.Errors.Add(fmt.Sprintf("unknown binary operator: %s", node.Op.String()), c.newLocation(node.Span()))
	}
}

// Resolves static AST expressions to Elk values
// and emits Bytecode that loads them.
// Returns false when the node can't be optimised at compile-time
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

// Emit an instruction that jumps back to the given Bytecode offset.
func (c *Compiler) emitLoop(span *position.Span, startOffset int) {
	c.emit(span.EndPos.Line, bytecode.LOOP)

	offset := c.nextInstructionOffset() - startOffset + 2
	if offset > math.MaxUint16 {
		c.Errors.Add(
			fmt.Sprintf("too many bytes to jumbytep backward: %d", math.MaxUint16),
			c.newLocation(span),
		)
	}

	c.Bytecode.AppendUint16(uint16(offset))
}

// Overwrite the placeholder operand of a jump instruction
func (c *Compiler) patchJump(offset int, span *position.Span) {
	jump := c.nextInstructionOffset() - offset - 2

	if jump > math.MaxUint16 {
		c.Errors.Add(
			fmt.Sprintf("too many bytes to jump over: %d", jump),
			c.newLocation(span),
		)
		return
	}

	c.Bytecode.Instructions[offset] = byte((jump >> 8) & 0xff)
	c.Bytecode.Instructions[offset+1] = byte(jump & 0xff)
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
	}
	if c.lastLocalIndex == math.MaxUint16 {
		c.Errors.Add(
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
