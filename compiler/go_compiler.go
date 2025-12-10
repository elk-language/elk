package compiler

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

func CreateGoCompiler(parent *GoCompiler, checker types.Checker, loc *position.Location, errors *diagnostic.SyncDiagnosticList) *GoCompiler {
	compiler := NewGoCompiler("main", topLevelGoCompilerMode, loc, checker)
	compiler.Errors = errors
	compiler.Parent = parent
	return compiler
}

func (c *GoCompiler) InitGlobalEnv() *GoCompiler {
	envCompiler := NewGoCompiler("initGlobalEnv", topLevelGoCompilerMode, c.loc, c.checker)
	envCompiler.Parent = c
	envCompiler.Errors = c.Errors
	envCompiler.compileGlobalEnv()
	return envCompiler
}

func (c *GoCompiler) EmitExecInParent() {
	parent := c.Parent
	parent.emit("%s()\n", c.Name)
}

// Compiler mode
type goMode uint8

const (
	topLevelGoCompilerMode goMode = iota
	methodGoCompilerMode
)

// Compiles Elk source code to Go source code.
type GoCompiler struct {
	Errors           *diagnostic.SyncDiagnosticList
	Parent           *GoCompiler
	Name             string
	buff             strings.Builder
	outerBuff        strings.Builder
	checker          types.Checker
	loc              *position.Location
	mode             goMode
	tmpLocalCounter  int
	callCacheCounter int
	bigIntCounter    int
}

type goValue struct {
	tmpLocal string
	inline   string
	typ      types.Type
}

func (v *goValue) IsInline() bool {
	return v.inline != ""
}

func (v *goValue) IsTmp() bool {
	return v.tmpLocal != ""
}

func (v *goValue) Value() string {
	if v.inline != "" {
		return v.inline
	}

	return v.tmpLocal
}

func newInlineValue(v string, typ types.Type) *goValue {
	return &goValue{
		inline: v,
		typ:    typ,
	}
}

func newTmpValue(v string, typ types.Type) *goValue {
	return &goValue{
		tmpLocal: v,
		typ:      typ,
	}
}

func NewGoCompiler(name string, mode goMode, loc *position.Location, checker types.Checker) *GoCompiler {
	return &GoCompiler{
		Name:    name,
		mode:    mode,
		Errors:  diagnostic.NewSyncDiagnosticList(),
		checker: checker,
		loc:     loc,
	}
}

func (c *GoCompiler) compileGlobalEnv() {
	c.emit("package main\n\n")
	c.emit("func initGlobalEnv() {\n")

	env := c.checker.Env()
	c.emit("var parentNamespace, namespace value.Value\n")
	c.compileModuleDefinition(env.Root, env.Root, value.ToSymbol("Root"))

	c.emit("}\n")
}

func (c *GoCompiler) compileNamespaceDefinition(parentNamespace, namespace types.Namespace, constName value.Symbol) {
	if !namespace.IsDefined() && !namespace.IsNative() {
		switch p := parentNamespace.(type) {
		case *types.SingletonClass:
			c.emit("parentNamespace = value.RootModule.Constants.Get(value.ToSymbol(%q)).SingletonClass()\n", p.AttachedObject.Name())
		default:
			c.emit("parentNamespace = value.RootModule.Constants.Get(value.ToSymbol(%q))\n", p.Name())
		}

		switch namespace.(type) {
		case *types.Module:
			c.emit("namespace = value.Ref(value.NewModule())\n")
		case *types.Class:
			c.emit("namespace = value.Ref(value.NewClassWithOptions(value.ClassWithSuperclass(nil)))\n")
		case *types.Mixin:
			c.emit("namespace = value.Ref(value.NewMixin())\n")
		case *types.Interface:
			c.emit("namespace = value.Ref(value.NewInterface())\n")
		}
		c.emit("value.AddConstant(parentNamespace, value.ToSymbol(%q), namespace)\n\n", constName)
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

// Emit code inside of a func
func (c *GoCompiler) emit(format string, a ...any) {
	fmt.Fprintf(&c.buff, format, a...)
}

// Emit package level code
func (c *GoCompiler) emitOuter(format string, a ...any) {
	fmt.Fprintf(&c.outerBuff, format, a...)
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

func (c *GoCompiler) CompileClassInheritance(class *types.Class, location *position.Location) {
	if class.IsCompiled() {
		return
	}
	superclass := class.Superclass()
	if superclass == nil {
		return
	}

	class.SetCompiled(true)

	c.emit("class = (*value.Class)(value.RootModule.Constants.Get(value.ToSymbol(%q)).Pointer())", class.Name())
	c.emit("superclass = (*value.Class)(value.RootModule.Constants.Get(value.ToSymbol(%q)).Pointer())", superclass.Name())
	c.emit("class.SetSuperclass(superclass)")
}

func (c *GoCompiler) CompileIvarIndices(target types.NamespaceWithIvarIndices, location *position.Location) {
	switch target := target.(type) {
	case *types.SingletonClass:
		c.emit("class = value.RootModule.Constants.Get(value.ToSymbol(%q)).SingletonClass()\n", target.AttachedObject.Name())
	case *types.Module:
		c.emit("class = value.RootModule.Constants.Get(value.ToSymbol(%q)).SingletonClass()\n", target.Name())
	default:
		c.emit("class = (*value.Class)(value.RootModule.Constants.Get(value.ToSymbol(%q)).Pointer())\n", target.Name())
	}

	c.emit("class.IvarIndices = %s\n", target.IvarIndices().ToGoSource())
}

func (c *GoCompiler) CompileInclude(target types.Namespace, mixin *types.Mixin, location *position.Location) {
	switch t := target.(type) {
	case *types.SingletonClass:
		c.emit("class = value.RootModule.Constants.Get(value.ToSymbol(%q)).SingletonClass()\n", t.AttachedObject.Name())
	default:
		c.emit("class = (*value.Class)(value.RootModule.Constants.Get(value.ToSymbol(%q)).Pointer())\n", target.Name())
	}

	c.emit("mixin = (*value.Mixin)(value.RootModule.Constants.Get(value.ToSymbol(%q)).Pointer())\n", mixin.Name())
	c.emit("class.IncludeMixin(mixin)\n")
}

func mangleFileName(name string) string {
	var b strings.Builder

	b.WriteString("__file_")

	for i, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			// Identifiers cannot start with a digit
			if i == 0 && unicode.IsDigit(r) {
				b.WriteByte('_')
			}
			b.WriteRune(r)
			continue
		}

		// Everything else becomes underscore
		b.WriteByte('_')
	}

	return b.String()
}

func (c *GoCompiler) InitExpressionCompiler(location *position.Location) *GoCompiler {
	name := mangleFileName(location.FilePath)
	exprCompiler := NewGoCompiler(name, topLevelGoCompilerMode, location, c.checker)
	exprCompiler.Errors = c.Errors

	c.emit("%s()\n", name)

	return exprCompiler
}

func (c *GoCompiler) CompileExpressionsInFile(node *ast.ProgramNode) {
	c.emit("func %s() {\n", c.Name)

	c.compileProgram(node)

	c.emit("}\n")
}

// Entry point to the compilation process
func (c *GoCompiler) compileProgram(node *ast.ProgramNode) {
	for _, stmt := range node.Body {
		c.compileStatement(stmt)
	}
}

func (c *GoCompiler) compileStatement(node ast.StatementNode) {
	switch node := node.(type) {
	case *ast.ExpressionStatementNode:
		c.compileExpression(node.Expression)
	}
}

func (c *GoCompiler) compileExpression(node ast.ExpressionNode) *goValue {
	switch node := node.(type) {
	case *ast.IntLiteralNode:
		i, err := value.ParseBigInt(node.Value, 0)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		if i.IsSmallInt() {
			return newInlineValue(
				fmt.Sprintf("value.SmallInt(%d).ToValue()", i.ToSmallInt()),
				c.checker.StdValue(),
			)
		}
		bigIntVar := c.emitBigIntVar(node.Value)
		return newTmpValue(
			bigIntVar,
			c.checker.StdValue(),
		)
	case *ast.Int8LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 8)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineValue(
			fmt.Sprintf("value.Int8(%d)", i),
			c.checker.Std(symbol.Int8),
		)
	case *ast.Int16LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 16)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineValue(
			fmt.Sprintf("value.Int16(%d)", i),
			c.checker.Std(symbol.Int16),
		)
	case *ast.Int32LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 32)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineValue(
			fmt.Sprintf("value.Int32(%d)", i),
			c.checker.Std(symbol.Int32),
		)
	case *ast.Int64LiteralNode:
		i, err := value.StrictParseInt(node.Value, 0, 64)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineValue(
			fmt.Sprintf("value.Int64(%d)", i),
			c.checker.Std(symbol.Int64),
		)
	case *ast.UInt8LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 8)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineValue(
			fmt.Sprintf("value.UInt8(%d)", i),
			c.checker.Std(symbol.UInt8),
		)
	case *ast.UInt16LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 16)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineValue(
			fmt.Sprintf("value.UInt16(%d)", i),
			c.checker.Std(symbol.UInt16),
		)
	case *ast.UInt32LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 32)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineValue(
			fmt.Sprintf("value.UInt32(%d)", i),
			c.checker.Std(symbol.UInt32),
		)
	case *ast.UInt64LiteralNode:
		i, err := value.StrictParseUint(node.Value, 0, 64)
		if !err.IsUndefined() {
			c.Errors.AddFailure(err.Error(), node.Location())
			return nil
		}
		return newInlineValue(
			fmt.Sprintf("value.UInt64(%d)", i),
			c.checker.Std(symbol.UInt64),
		)
	case *ast.BinaryExpressionNode:
		return c.compileBinaryExpressionNode(node)
	default:
		panic(fmt.Sprintf("invalid expression node: %T", node))
	}
}

func (c *GoCompiler) getTmpIdent() string {
	c.tmpLocalCounter++
	return fmt.Sprintf("t%d", c.tmpLocalCounter)
}

func (c *GoCompiler) emitErrorPropagation() {
	switch c.mode {
	case topLevelGoCompilerMode:
		c.emit("if err.IsNotUndefined() { thread.PrintErrorValue(err); os.Exit(1) }\n")
	default:
		c.emit("if err.IsNotUndefined() { return value.Undefined, err }\n")
	}
}

func (c *GoCompiler) emitCallCache() string {
	c.callCacheCounter++
	callCacheName := fmt.Sprintf("cc_%s_%d", c.Name, c.callCacheCounter)
	c.emitOuter("var %s *value.CallCache", callCacheName)

	return callCacheName
}

func (c *GoCompiler) emitBigIntVar(val string) string {
	c.bigIntCounter++
	varName := fmt.Sprintf("bi_%s_%d", c.Name, c.bigIntCounter)
	c.emitOuter("var %s value.Value = value.MustParseInt(%q, 0)", varName, val)

	return varName
}

func (c *GoCompiler) emitAddCallFrame(loc *position.Location) {
	c.emit(
		"thread.AddCallFrame(value.CallFrame{FuncName: %q, FileName: %q, LineNumber: %d})\n",
		c.Name,
		loc.FilePath,
		loc.StartPos.Line,
	)
}
func (c *GoCompiler) emitPopCallFrame() {
	c.emit("thread.PopCallFrame()\n")
}

func (c *GoCompiler) compileBinaryExpressionNode(node *ast.BinaryExpressionNode) *goValue {
	if resolved := c.resolve(node); resolved != nil {
		return resolved
	}

	left := c.compileExpression(node.Left)
	right := c.compileExpression(node.Right)
	typ := left.typ

	switch node.Op.Type {
	case token.PLUS:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			tmp := c.getTmpIdent()
			c.emit("%s, err := value.AddInt(%s, %s)\n", tmp, left.Value(), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
			tmp := c.getTmpIdent()
			c.emit("%s, err := %s.AddVal(%s)\n", tmp, left.Value(), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinAddable)) {
			tmp := c.getTmpIdent()
			c.emit("%s, err := value.AddVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}

		callCache := c.emitCallCache()
		tmp := c.getTmpIdent()
		c.emitAddCallFrame(node.Location())
		c.emit("%s, err := thread.CallMethodByNameWithCache(symbol.OpAdd, %s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
		c.emitPopCallFrame()
		c.emitErrorPropagation()
		return newTmpValue(
			tmp,
			c.checker.Std(symbol.Value),
		)
	case token.MINUS:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			tmp := c.getTmpIdent()
			c.emit("%s, err := value.SubtractInt(%s, %s)\n", tmp, left.Value(), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
			tmp := c.getTmpIdent()
			c.emit("%s, err := %s.SubtractVal(%s)\n", tmp, left.Value(), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinAddable)) {
			tmp := c.getTmpIdent()
			c.emit("%s, err := value.SubtractVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}

		callCache := c.emitCallCache()
		tmp := c.getTmpIdent()
		c.emitAddCallFrame(node.Location())
		c.emit("%s, err := thread.CallMethodByNameWithCache(symbol.OpSubtract, %s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
		c.emitPopCallFrame()
		c.emitErrorPropagation()
		return newTmpValue(
			tmp,
			c.checker.Std(symbol.Value),
		)
	case token.STAR:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			tmp := c.getTmpIdent()
			c.emit("%s, err := value.MultiplyInt(%s, %s)\n", tmp, left.Value(), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
			tmp := c.getTmpIdent()
			c.emit("%s, err := %s.MultiplyVal(%s)\n", tmp, left.Value(), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinAddable)) {
			tmp := c.getTmpIdent()
			c.emit("%s, err := value.MultiplyVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}

		callCache := c.emitCallCache()
		tmp := c.getTmpIdent()
		c.emitAddCallFrame(node.Location())
		c.emit("%s, err := thread.CallMethodByNameWithCache(symbol.OpMultiply, %s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
		c.emitPopCallFrame()
		c.emitErrorPropagation()
		return newTmpValue(
			tmp,
			c.checker.Std(symbol.Value),
		)
	case token.SLASH:
		c.emitAddCallFrame(node.Location())

		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			tmp := c.getTmpIdent()
			c.emit("%s, err := value.DivideInt(%s, %s)\n", tmp, left.Value(), c.convertToValue(right))
			c.emitPopCallFrame()
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
			tmp := c.getTmpIdent()
			c.emit("%s, err := %s.DivideVal(%s)\n", tmp, left.Value(), c.convertToValue(right))
			c.emitPopCallFrame()
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinAddable)) {
			tmp := c.getTmpIdent()
			c.emit("%s, err := value.DivideVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
			c.emitPopCallFrame()
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}

		callCache := c.emitCallCache()
		tmp := c.getTmpIdent()
		c.emit("%s, err := thread.CallMethodByNameWithCache(symbol.OpMultiply, %s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
		c.emitPopCallFrame()
		c.emitErrorPropagation()
		return newTmpValue(
			tmp,
			c.checker.Std(symbol.Value),
		)
	case token.STAR_STAR:
		if c.checker.IsSubtype(typ, c.checker.StdInt()) {
			tmp := c.getTmpIdent()
			c.emit("%s, err := value.ExponentiateInt(%s, %s)\n", tmp, left.Value(), c.convertToValue(right))
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinNumeric)) {
			tmp := c.getTmpIdent()
			c.emit("%s, err := value.ExponentiateVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
			c.emitPopCallFrame()
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}

		callCache := c.emitCallCache()
		tmp := c.getTmpIdent()
		c.emitAddCallFrame(node.Location())
		c.emit("%s, err := thread.CallMethodByNameWithCache(symbol.OpMultiply, %s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
		c.emitPopCallFrame()
		c.emitErrorPropagation()
		return newTmpValue(
			tmp,
			c.checker.Std(symbol.Value),
		)
	case token.LBITSHIFT:
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
			tmp := c.getTmpIdent()
			c.emit("%s, err := value.LeftBitshiftVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
			c.emitPopCallFrame()
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}

		callCache := c.emitCallCache()
		tmp := c.getTmpIdent()
		c.emitAddCallFrame(node.Location())
		c.emit("%s, err := thread.CallMethodByNameWithCache(symbol.OpLeftBitshift, %s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
		c.emitPopCallFrame()
		c.emitErrorPropagation()
		return newTmpValue(
			tmp,
			c.checker.Std(symbol.Value),
		)
	case token.LTRIPLE_BITSHIFT:
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinLogicBitshiftable)) {
			tmp := c.getTmpIdent()
			c.emit("%s, err := value.LogicalLeftBitshiftVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
			c.emitPopCallFrame()
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}

		callCache := c.emitCallCache()
		tmp := c.getTmpIdent()
		c.emitAddCallFrame(node.Location())
		c.emit("%s, err := thread.CallMethodByNameWithCache(symbol.OpLogicalLeftBitshift, %s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
		c.emitPopCallFrame()
		c.emitErrorPropagation()
		return newTmpValue(
			tmp,
			c.checker.Std(symbol.Value),
		)
	case token.RBITSHIFT:
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
			tmp := c.getTmpIdent()
			c.emit("%s, err := value.RightBitshiftVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
			c.emitPopCallFrame()
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}

		callCache := c.emitCallCache()
		tmp := c.getTmpIdent()
		c.emitAddCallFrame(node.Location())
		c.emit("%s, err := thread.CallMethodByNameWithCache(symbol.OpRightBitshift, %s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
		c.emitPopCallFrame()
		c.emitErrorPropagation()
		return newTmpValue(
			tmp,
			c.checker.Std(symbol.Value),
		)
	case token.RTRIPLE_BITSHIFT:
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinLogicBitshiftable)) {
			tmp := c.getTmpIdent()
			c.emit("%s, err := value.LogicalRightBitshiftVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
			c.emitPopCallFrame()
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}

		callCache := c.emitCallCache()
		tmp := c.getTmpIdent()
		c.emitAddCallFrame(node.Location())
		c.emit("%s, err := thread.CallMethodByNameWithCache(symbol.OpLogicalRightBitshift, %s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
		c.emitPopCallFrame()
		c.emitErrorPropagation()
		return newTmpValue(
			tmp,
			c.checker.Std(symbol.Value),
		)
	case token.AND:
		if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
			tmp := c.getTmpIdent()
			c.emit("%s, err := value.BitwiseAndVal(%s, %s)\n", tmp, c.convertToValue(left), c.convertToValue(right))
			c.emitPopCallFrame()
			c.emitErrorPropagation()

			return newTmpValue(
				tmp,
				c.checker.Std(symbol.Value),
			)
		}

		callCache := c.emitCallCache()
		tmp := c.getTmpIdent()
		c.emitAddCallFrame(node.Location())
		c.emit("%s, err := thread.CallMethodByNameWithCache(symbol.OpAnd, %s, %s, %s)", tmp, callCache, c.convertToValue(left), c.convertToValue(right))
		c.emitPopCallFrame()
		c.emitErrorPropagation()
		return newTmpValue(
			tmp,
			c.checker.Std(symbol.Value),
		)
	// case token.AND_TILDE:
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
	// 		c.emit(line, bytecode.BITWISE_AND_NOT)
	// 		return
	// 	}
	// 	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpAndNot, 1), location, false)
	// case token.OR:
	// 	if c.checker.IsSubtype(typ, c.checker.StdInt()) {
	// 		c.emit(line, bytecode.BITWISE_OR_INT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
	// 		c.emit(line, bytecode.BITWISE_OR)
	// 		return
	// 	}
	// 	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpOr, 1), location, false)
	// case token.XOR:
	// 	if c.checker.IsSubtype(typ, c.checker.StdInt()) {
	// 		c.emit(line, bytecode.BITWISE_XOR_INT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinInt)) {
	// 		c.emit(line, bytecode.BITWISE_XOR)
	// 		return
	// 	}
	// 	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpXor, 1), location, false)
	// case token.PERCENT:
	// 	if c.checker.IsSubtype(typ, c.checker.StdInt()) {
	// 		c.emit(line, bytecode.MODULO_INT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
	// 		c.emit(line, bytecode.MODULO_FLOAT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinNumeric)) {
	// 		c.emit(line, bytecode.MODULO)
	// 		return
	// 	}
	// 	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpModulo, 1), location, false)
	// case token.LAX_EQUAL:
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinEquatable)) {
	// 		c.emit(line, bytecode.LAX_EQUAL)
	// 		return
	// 	}
	// 	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpLaxEqual, 1), location, false)
	// case token.LAX_NOT_EQUAL:
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinEquatable)) {
	// 		c.emit(line, bytecode.LAX_NOT_EQUAL)
	// 		return
	// 	}
	// 	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpLaxEqual, 1), location, false)
	// 	c.emit(line, bytecode.NOT)
	// case token.EQUAL_EQUAL:
	// 	if c.checker.IsSubtype(typ, c.checker.StdInt()) {
	// 		c.emit(line, bytecode.EQUAL_INT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
	// 		c.emit(line, bytecode.EQUAL_INT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinEquatable)) {
	// 		c.emit(line, bytecode.EQUAL)
	// 		return
	// 	}
	// 	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpEqual, 1), location, false)
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
	// case token.LESS_EQUAL:
	// 	if c.checker.IsSubtype(typ, c.checker.StdInt()) {
	// 		c.emit(line, bytecode.LESS_EQUAL_INT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.StdFloat()) {
	// 		c.emit(line, bytecode.LESS_EQUAL_FLOAT)
	// 		return
	// 	}
	// 	if c.checker.IsSubtype(typ, c.checker.Std(symbol.S_BuiltinNumeric)) {
	// 		c.emit(line, bytecode.LESS_EQUAL)
	// 		return
	// 	}
	// 	c.emitCallMethod(value.NewCallSiteInfo(symbol.OpLessThanEqual, 1), location, false)
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
		return nil
	}
}

func (c *GoCompiler) resolve(node ast.ExpressionNode) *goValue {
	result := resolve(node)
	if result.IsUndefined() {
		return nil
	}

	return c.valueToGoSource(result)
}

func (c *GoCompiler) valueToGoSource(val value.Value) *goValue {
	if val.IsReference() {
		switch v := val.AsReference().(type) {
		case *value.ArrayList:
			return c.arrayListToGoSource(v)
		case *value.ArrayTuple:
			return c.arrayTupleToGoSource(v)
		case *value.HashSet:
			return c.hashSetToGoSource(v)
		case *value.HashMap:
			return c.hashMapToGoSource(v)
		case *value.HashRecord:
			return c.hashRecordToGoSource(v)
		case value.Int64:
			return newInlineValue(
				fmt.Sprintf("value.Int64(%d)", v),
				c.checker.Std(symbol.Int64),
			)
		case value.UInt64:
			return newInlineValue(
				fmt.Sprintf("value.UInt64(%d)", v),
				c.checker.Std(symbol.UInt64),
			)
		default:
			return nil
		}
	}

	switch val.ValueFlag() {
	case value.TRUE_FLAG:
		return newInlineValue(
			"value.True",
			c.checker.StdValue(),
		)
	case value.FALSE_FLAG:
		return newInlineValue(
			"value.False",
			c.checker.StdValue(),
		)
	case value.NIL_FLAG:
		return newInlineValue(
			"value.Nil",
			c.checker.StdValue(),
		)
	case value.SMALL_INT_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.SmallInt(%d).ToValue()", val.AsSmallInt()),
			c.checker.StdValue(),
		)
	case value.INT64_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.Int64(%d)", val.AsInt64()),
			c.checker.Std(symbol.Int64),
		)
	case value.UINT64_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.UInt64(%d)", val.AsUInt64()),
			c.checker.Std(symbol.UInt64),
		)
	case value.INT32_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.Int32(%d)", val.AsInt32()),
			c.checker.Std(symbol.Int32),
		)
	case value.UINT32_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.UInt32(%d)", val.AsUInt32()),
			c.checker.Std(symbol.UInt32),
		)
	case value.INT16_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.Int16(%d)", val.AsInt16()),
			c.checker.Std(symbol.Int16),
		)
	case value.UINT16_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.UInt16(%d)", val.AsUInt16()),
			c.checker.Std(symbol.UInt16),
		)
	case value.INT8_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.Int8(%d)", val.AsInt8()),
			c.checker.Std(symbol.Int8),
		)
	case value.UINT8_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.UInt8(%d)", val.AsUInt8()),
			c.checker.Std(symbol.UInt8),
		)
	case value.CHAR_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.Char(%q)", val.AsChar()),
			c.checker.Std(symbol.Char),
		)
	case value.FLOAT_FLAG:
		return newInlineValue(
			fmt.Sprintf("value.Float(%f)", val.AsFloat()),
			c.checker.Std(symbol.Float),
		)
	}

	return nil
}

func (c *GoCompiler) arrayListToGoSource(v *value.ArrayList) *goValue {
	var buff strings.Builder

	fmt.Fprintf(&buff, "value.NewArrayListWithElements(%d, ", v.LeftCapacity())

	for _, element := range *v {
		el := c.valueToGoSource(element)
		if el == nil {
			return nil
		}

		fmt.Fprintf(&buff, "%s, ", c.convertToValue(el))
	}

	buff.WriteString(")")
	return newInlineValue(buff.String(), c.checker.Std(symbol.ArrayList))
}

func (c *GoCompiler) convertToValue(v *goValue) string {
	if c.checker.IsTheSameType(v.typ, c.checker.Std(symbol.Value)) {
		return v.Value()
	}
	if c.checker.IsSubtype(v.typ, c.checker.Std(symbol.Object)) {
		return fmt.Sprintf("value.Ref(%s)", v.Value())
	}

	return fmt.Sprintf("%s.ToValue()", v.Value())
}

func (c *GoCompiler) arrayTupleToGoSource(v *value.ArrayTuple) *goValue {
	var buff strings.Builder

	buff.WriteString("value.NewArrayTupleWithElements(0, ")

	for _, element := range *v {
		el := c.valueToGoSource(element)
		if el == nil {
			return nil
		}

		fmt.Fprintf(&buff, "%s, ", c.convertToValue(el))
	}

	buff.WriteString(")")
	return newInlineValue(buff.String(), c.checker.Std(symbol.ArrayTuple))
}

func (c *GoCompiler) hashSetToGoSource(v *value.HashSet) *goValue {
	var buff strings.Builder

	fmt.Fprintf(&buff, "vm.MustNewHashSetWithCapacityAndElements(nil, %d, ", v.LeftCapacity())

	for element := range v.All() {
		el := c.valueToGoSource(element)
		if el == nil {
			return nil
		}

		fmt.Fprintf(&buff, "%s, ", c.convertToValue(el))
	}

	buff.WriteString(")")
	return newInlineValue(buff.String(), c.checker.Std(symbol.HashSet))
}

func (c *GoCompiler) hashMapToGoSource(v *value.HashMap) *goValue {
	var buff strings.Builder

	fmt.Fprintf(&buff, "vm.MustNewHashMapWithCapacityAndElements(nil, %d, ", v.LeftCapacity())

	for pair := range v.All() {
		p := c.valuePairToGoSource(pair)
		if p == nil {
			return nil
		}

		fmt.Fprintf(&buff, "%s, ", p)
	}

	buff.WriteString(")")
	return newInlineValue(buff.String(), c.checker.Std(symbol.HashMap))
}

func (c *GoCompiler) hashRecordToGoSource(v *value.HashRecord) *goValue {
	var buff strings.Builder

	buff.WriteString("vm.MustNewHashRecordWithElements(nil, ")

	for pair := range v.All() {
		p := c.valuePairToGoSource(pair)
		if p == nil {
			return nil
		}

		fmt.Fprintf(&buff, "%s, ", p)
	}

	buff.WriteString(")")
	return newInlineValue(buff.String(), c.checker.Std(symbol.HashRecord))
}

func (c *GoCompiler) valuePairToGoSource(p value.Pair) *goValue {
	k := c.valueToGoSource(p.Key)
	if k == nil {
		return nil
	}
	v := c.valueToGoSource(p.Key)
	if v == nil {
		return nil
	}

	return newInlineValue(
		fmt.Sprintf("Pair{Key: %s, Value: %s}", k.inline, v.inline),
		c.checker.Std(symbol.Pair),
	)
}
