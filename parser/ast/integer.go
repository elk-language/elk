package ast

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// Int literal eg. `5`, `125_355`, `0xff`
type IntLiteralNode struct {
	TypedNodeBase
	Value string
}

func (n *IntLiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &IntLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *IntLiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::IntLiteralNode", env)
}

func (n *IntLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *IntLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*IntLiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.loc.Equal(o.loc)
}

func (n *IntLiteralNode) String() string {
	return n.Value
}

func (*IntLiteralNode) IsStatic() bool {
	return true
}

func (*IntLiteralNode) Class() *value.Class {
	return value.IntLiteralNodeClass
}

func (*IntLiteralNode) DirectClass() *value.Class {
	return value.IntLiteralNodeClass
}

func (n *IntLiteralNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::IntLiteralNode{location: %s, value: %s}",
		(*value.Location)(n.loc).Inspect(), n.Value,
	)
}

func (n *IntLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new int literal node eg. `5`, `125_355`, `0xff`
func NewIntLiteralNode(loc *position.Location, val string) *IntLiteralNode {
	return &IntLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}

// Int64 literal eg. `5i64`, `125_355i64`, `0xffi64`
type Int64LiteralNode struct {
	TypedNodeBase
	Value string
}

func (n *Int64LiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &Int64LiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *Int64LiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::Int64LiteralNode", env)
}

func (n *Int64LiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *Int64LiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*Int64LiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.loc.Equal(o.loc)
}

func (n *Int64LiteralNode) String() string {
	return fmt.Sprintf("%si64", n.Value)
}

func (*Int64LiteralNode) IsStatic() bool {
	return true
}

func (*Int64LiteralNode) Class() *value.Class {
	return value.Int64LiteralNodeClass
}

func (*Int64LiteralNode) DirectClass() *value.Class {
	return value.Int64LiteralNodeClass
}

func (n *Int64LiteralNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::Int64LiteralNode{location: %s, value: %s}",
		(*value.Location)(n.loc).Inspect(),
		n.Value,
	)
}

func (n *Int64LiteralNode) Error() string {
	return n.Inspect()
}

// Create a new Int64 literal node eg. `5i64`, `125_355i64`, `0xffi64`
func NewInt64LiteralNode(loc *position.Location, val string) *Int64LiteralNode {
	return &Int64LiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}

// UInt literal eg. `5u`, `125_355u`, `0xffu`
type UIntLiteralNode struct {
	TypedNodeBase
	Value string
}

func (n *UIntLiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::UIntLiteralNode", env)
}

func (n *UIntLiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &UIntLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *UIntLiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *UIntLiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UIntLiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.loc.Equal(o.loc)
}

func (n *UIntLiteralNode) String() string {
	return fmt.Sprintf("%su", n.Value)
}

func (*UIntLiteralNode) IsStatic() bool {
	return true
}

func (*UIntLiteralNode) Class() *value.Class {
	return value.UIntLiteralNodeClass
}

func (*UIntLiteralNode) DirectClass() *value.Class {
	return value.UIntLiteralNodeClass
}

func (n *UIntLiteralNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::UIntLiteralNode{location: %s, value: %s}",
		(*value.Location)(n.loc).Inspect(),
		n.Value,
	)
}

func (n *UIntLiteralNode) Error() string {
	return n.Inspect()
}

// Create a new UInt literal node eg. `5u`, `125_355u`, `0xffu`
func NewUIntLiteralNode(loc *position.Location, val string) *UIntLiteralNode {
	return &UIntLiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}

// UInt64 literal eg. `5u64`, `125_355u64`, `0xffu64`
type UInt64LiteralNode struct {
	TypedNodeBase
	Value string
}

func (n *UInt64LiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::UInt64LiteralNode", env)
}

func (n *UInt64LiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &UInt64LiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *UInt64LiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *UInt64LiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UInt64LiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.loc.Equal(o.loc)
}

func (n *UInt64LiteralNode) String() string {
	return fmt.Sprintf("%su64", n.Value)
}

func (*UInt64LiteralNode) IsStatic() bool {
	return true
}

func (*UInt64LiteralNode) Class() *value.Class {
	return value.UInt64LiteralNodeClass
}

func (*UInt64LiteralNode) DirectClass() *value.Class {
	return value.UInt64LiteralNodeClass
}

func (n *UInt64LiteralNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::UInt64LiteralNode{location: %s, value: %s}",
		(*value.Location)(n.loc).Inspect(),
		n.Value,
	)
}

func (n *UInt64LiteralNode) Error() string {
	return n.Inspect()
}

// Create a new UInt64 literal node eg. `5u64`, `125_355u64`, `0xffu64`
func NewUInt64LiteralNode(loc *position.Location, val string) *UInt64LiteralNode {
	return &UInt64LiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}

// Int32 literal eg. `5i32`, `1_20i32`, `0xffi32`
type Int32LiteralNode struct {
	TypedNodeBase
	Value string
}

func (n *Int32LiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &Int32LiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *Int32LiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::Int32LiteralNode", env)
}

func (n *Int32LiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *Int32LiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*Int32LiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.loc.Equal(o.loc)
}

func (n *Int32LiteralNode) String() string {
	return fmt.Sprintf("%si32", n.Value)
}

func (*Int32LiteralNode) IsStatic() bool {
	return true
}

func (*Int32LiteralNode) Class() *value.Class {
	return value.Int32LiteralNodeClass
}

func (*Int32LiteralNode) DirectClass() *value.Class {
	return value.Int32LiteralNodeClass
}

func (n *Int32LiteralNode) Inspect() string {
	return fmt.Sprintf("Std::Elk::AST::Int32LiteralNode{location: %s, value: %s}", (*value.Location)(n.loc).Inspect(), n.Value)
}

func (n *Int32LiteralNode) Error() string {
	return n.Inspect()
}

// Create a new Int32 literal node eg. `5i32`, `1_20i32`, `0xffi32`
func NewInt32LiteralNode(loc *position.Location, val string) *Int32LiteralNode {
	return &Int32LiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}

// UInt32 literal eg. `5u32`, `1_20u32`, `0xffu32`
type UInt32LiteralNode struct {
	TypedNodeBase
	Value string
}

func (n *UInt32LiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &UInt32LiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *UInt32LiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::UInt32LiteralNode", env)
}

func (n *UInt32LiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *UInt32LiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UInt32LiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.loc.Equal(o.loc)
}

func (n *UInt32LiteralNode) String() string {
	return fmt.Sprintf("%su32", n.Value)
}

func (*UInt32LiteralNode) IsStatic() bool {
	return true
}

func (*UInt32LiteralNode) Class() *value.Class {
	return value.UInt32LiteralNodeClass
}

func (*UInt32LiteralNode) DirectClass() *value.Class {
	return value.UInt32LiteralNodeClass
}

func (n *UInt32LiteralNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::UInt32LiteralNode{location: %s, value: %s}",
		(*value.Location)(n.loc).Inspect(),
		n.Value,
	)
}

func (n *UInt32LiteralNode) Error() string {
	return n.Inspect()
}

// Create a new UInt32 literal node eg. `5u32`, `1_20u32`, `0xffu32`
func NewUInt32LiteralNode(loc *position.Location, val string) *UInt32LiteralNode {
	return &UInt32LiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}

// Int16 literal eg. `5i16`, `1_20i16`, `0xffi16`
type Int16LiteralNode struct {
	TypedNodeBase
	Value string
}

func (n *Int16LiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &Int16LiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *Int16LiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::Int16LiteralNode", env)
}

func (n *Int16LiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *Int16LiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*Int16LiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.loc.Equal(o.loc)
}

func (n *Int16LiteralNode) String() string {
	return fmt.Sprintf("%si16", n.Value)
}

func (*Int16LiteralNode) IsStatic() bool {
	return true
}

func (*Int16LiteralNode) Class() *value.Class {
	return value.Int16LiteralNodeClass
}

func (*Int16LiteralNode) DirectClass() *value.Class {
	return value.Int16LiteralNodeClass
}

func (n *Int16LiteralNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::Int16LiteralNode{location: %s, value: %s}",
		(*value.Location)(n.loc).Inspect(),
		n.Value,
	)
}

func (n *Int16LiteralNode) Error() string {
	return n.Inspect()
}

// Create a new Int16 literal node eg. `5i16`, `1_20i16`, `0xffi16`
func NewInt16LiteralNode(loc *position.Location, val string) *Int16LiteralNode {
	return &Int16LiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}

// UInt16 literal eg. `5u16`, `1_20u16`, `0xffu16`
type UInt16LiteralNode struct {
	TypedNodeBase
	Value string
}

func (n *UInt16LiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &UInt16LiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *UInt16LiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::UInt16LiteralNode", env)
}

func (n *UInt16LiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *UInt16LiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UInt16LiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.loc.Equal(o.loc)
}

func (n *UInt16LiteralNode) String() string {
	return fmt.Sprintf("%su16", n.Value)
}

func (*UInt16LiteralNode) IsStatic() bool {
	return true
}

func (*UInt16LiteralNode) Class() *value.Class {
	return value.UInt16LiteralNodeClass
}

func (*UInt16LiteralNode) DirectClass() *value.Class {
	return value.UInt16LiteralNodeClass
}

func (n *UInt16LiteralNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::UInt16LiteralNode{location: %s, value: %s}",
		(*value.Location)(n.loc).Inspect(),
		n.Value,
	)
}

func (n *UInt16LiteralNode) Error() string {
	return n.Inspect()
}

// Create a new UInt16 literal node eg. `5u16`, `1_20u16`, `0xffu16`
func NewUInt16LiteralNode(loc *position.Location, val string) *UInt16LiteralNode {
	return &UInt16LiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}

// Int8 literal eg. `5i8`, `1_20i8`, `0xffi8`
type Int8LiteralNode struct {
	TypedNodeBase
	Value string
}

func (n *Int8LiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &Int8LiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *Int8LiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::Int8LiteralNode", env)
}

func (n *Int8LiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *Int8LiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*Int8LiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.loc.Equal(o.loc)
}

func (n *Int8LiteralNode) String() string {
	return fmt.Sprintf("%si8", n.Value)
}

func (*Int8LiteralNode) IsStatic() bool {
	return true
}

func (*Int8LiteralNode) Class() *value.Class {
	return value.Int8LiteralNodeClass
}

func (*Int8LiteralNode) DirectClass() *value.Class {
	return value.Int8LiteralNodeClass
}

func (n *Int8LiteralNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::Int8LiteralNode{location: %s, value: %s}",
		(*value.Location)(n.loc).Inspect(),
		n.Value,
	)
}

func (n *Int8LiteralNode) Error() string {
	return n.Inspect()
}

// Create a new Int8 literal node eg. `5i8`, `1_20i8`, `0xffi8`
func NewInt8LiteralNode(loc *position.Location, val string) *Int8LiteralNode {
	return &Int8LiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}

// UInt8 literal eg. `5u8`, `1_20u8`, `0xffu8`
type UInt8LiteralNode struct {
	TypedNodeBase
	Value string
}

func (n *UInt8LiteralNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return &UInt8LiteralNode{
		TypedNodeBase: TypedNodeBase{loc: position.SpliceLocation(loc, n.loc, unquote), typ: n.typ},
		Value:         n.Value,
	}
}

func (n *UInt8LiteralNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::UInt8LiteralNode", env)
}

func (n *UInt8LiteralNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n *UInt8LiteralNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*UInt8LiteralNode)
	if !ok {
		return false
	}

	return n.Value == o.Value &&
		n.loc.Equal(o.loc)
}

func (n *UInt8LiteralNode) String() string {
	return fmt.Sprintf("%su8", n.Value)
}

func (*UInt8LiteralNode) IsStatic() bool {
	return true
}

func (*UInt8LiteralNode) Class() *value.Class {
	return value.UInt8LiteralNodeClass
}

func (*UInt8LiteralNode) DirectClass() *value.Class {
	return value.UInt8LiteralNodeClass
}

func (n *UInt8LiteralNode) Inspect() string {
	return fmt.Sprintf(
		"Std::Elk::AST::UInt8LiteralNode{location: %s, value: %s}",
		(*value.Location)(n.loc).Inspect(),
		n.Value,
	)
}

func (n *UInt8LiteralNode) Error() string {
	return n.Inspect()
}

// Create a new UInt8 literal node eg. `5u8`, `1_20u8`, `0xffu8`
func NewUInt8LiteralNode(loc *position.Location, val string) *UInt8LiteralNode {
	return &UInt8LiteralNode{
		TypedNodeBase: TypedNodeBase{loc: loc},
		Value:         val,
	}
}
