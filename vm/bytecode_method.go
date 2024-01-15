package vm

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strings"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// A single unit of Elk bytecode.
type BytecodeMethod struct {
	Instructions []byte
	Values       []value.Value // The value pool
	LineInfoList bytecode.LineInfoList
	Location     *position.Location
	Doc          value.Value

	name                   value.Symbol
	parameters             []value.Symbol
	optionalParameterCount int
	postRestParameterCount int
	namedRestParameter     bool
	sealed                 bool
}

func (b *BytecodeMethod) Name() value.Symbol {
	return b.name
}

func (b *BytecodeMethod) ParameterCount() int {
	return len(b.parameters)
}

func (b *BytecodeMethod) SetParameters(params []value.Symbol) {
	b.parameters = params
}

func (b *BytecodeMethod) Parameters() []value.Symbol {
	return b.parameters
}

func (b *BytecodeMethod) SetOptionalParameterCount(optParamCount int) {
	b.optionalParameterCount = optParamCount
}

func (b *BytecodeMethod) IncrementOptionalParameterCount() {
	b.optionalParameterCount++
}

func (b *BytecodeMethod) OptionalParameterCount() int {
	return b.optionalParameterCount
}

func (b *BytecodeMethod) SetPostRestParameterCount(postParamCount int) {
	b.postRestParameterCount = postParamCount
}

func (b *BytecodeMethod) IncrementPostRestParameterCount() {
	b.postRestParameterCount++
}

func (b *BytecodeMethod) PostRestParameterCount() int {
	return b.postRestParameterCount
}

func (b *BytecodeMethod) SetNamedRestParameter(present bool) {
	b.namedRestParameter = present
}

func (b *BytecodeMethod) NamedRestParameter() bool {
	return b.namedRestParameter
}

func (*BytecodeMethod) Class() *value.Class {
	return value.MethodClass
}

func (*BytecodeMethod) DirectClass() *value.Class {
	return value.MethodClass
}

func (*BytecodeMethod) SingletonClass() *value.Class {
	return nil
}

func (b *BytecodeMethod) IsSealed() bool {
	return b.sealed
}

func (b *BytecodeMethod) SetSealed() {
	b.sealed = true
}

func (b *BytecodeMethod) Copy() value.Value {
	return b
}

func (b *BytecodeMethod) Inspect() string {
	return fmt.Sprintf("Method{name: %s, type: :bytecode, location: %s}", b.name.Inspect(), b.Location.String())
}

func (*BytecodeMethod) InstanceVariables() value.SymbolMap {
	return nil
}

// Create a new bytecode method.
func NewBytecodeMethodSimple(name value.Symbol, instruct []byte, loc *position.Location) *BytecodeMethod {
	return &BytecodeMethod{
		Instructions:           instruct,
		Location:               loc,
		name:                   name,
		postRestParameterCount: -1,
	}
}

// Create a new bytecode method.
func NewBytecodeMethod(
	name value.Symbol,
	instruct []byte,
	loc *position.Location,
	lineInfo bytecode.LineInfoList,
	params []value.Symbol,
	optParamCount int,
	postRestParamCount int,
	namedRestParam bool,
	sealed bool,
	values []value.Value,
) *BytecodeMethod {
	return &BytecodeMethod{
		name:                   name,
		Instructions:           instruct,
		Location:               loc,
		LineInfoList:           lineInfo,
		parameters:             params,
		optionalParameterCount: optParamCount,
		postRestParameterCount: postRestParamCount,
		namedRestParameter:     namedRestParam,
		Values:                 values,
		sealed:                 sealed,
	}
}

type BytecodeMethodOption func(*BytecodeMethod)

func BytecodeMethodWithName(name value.Symbol) BytecodeMethodOption {
	return func(b *BytecodeMethod) {
		b.name = name
	}
}

func BytecodeMethodWithStringName(name string) BytecodeMethodOption {
	return func(b *BytecodeMethod) {
		b.name = value.ToSymbol(name)
	}
}

func BytecodeMethodWithInstructions(instructs []byte) BytecodeMethodOption {
	return func(b *BytecodeMethod) {
		b.Instructions = instructs
	}
}

func BytecodeMethodWithLocation(loc *position.Location) BytecodeMethodOption {
	return func(b *BytecodeMethod) {
		b.Location = loc
	}
}

func BytecodeMethodWithLineInfoList(lineInfo bytecode.LineInfoList) BytecodeMethodOption {
	return func(b *BytecodeMethod) {
		b.LineInfoList = lineInfo
	}
}

func BytecodeMethodWithParameters(params []value.Symbol) BytecodeMethodOption {
	return func(b *BytecodeMethod) {
		b.parameters = params
	}
}

func BytecodeMethodWithOptionalParameters(optParams int) BytecodeMethodOption {
	return func(b *BytecodeMethod) {
		b.optionalParameterCount = optParams
	}
}

func BytecodeMethodWithPostParameters(postParams int) BytecodeMethodOption {
	return func(b *BytecodeMethod) {
		b.postRestParameterCount = postParams
	}
}

func BytecodeMethodWithPositionalRestParameter() BytecodeMethodOption {
	return func(b *BytecodeMethod) {
		b.postRestParameterCount = 0
	}
}

func BytecodeMethodWithNamedRestParameter() BytecodeMethodOption {
	return func(b *BytecodeMethod) {
		b.namedRestParameter = true
	}
}

func BytecodeMethodWithSealed() BytecodeMethodOption {
	return func(b *BytecodeMethod) {
		b.sealed = true
	}
}

// Create a new bytecode method with options.
func NewBytecodeMethodWithOptions(opts ...BytecodeMethodOption) *BytecodeMethod {
	b := &BytecodeMethod{
		postRestParameterCount: -1,
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

// Create a new bytecode method.
func NewBytecodeMethodNoParams(
	name value.Symbol,
	instruct []byte,
	loc *position.Location,
	lineInfo bytecode.LineInfoList,
	values []value.Value,
) *BytecodeMethod {
	return NewBytecodeMethod(name, instruct, loc, lineInfo, nil, 0, -1, false, false, values)
}

// Add a parameter to the method.
func (f *BytecodeMethod) AddParameter(name value.Symbol) {
	f.parameters = append(f.parameters, name)
}

// Add a parameter to the method.
func (f *BytecodeMethod) AddParameterString(name string) {
	f.parameters = append(f.parameters, value.ToSymbol(name))
}

// Add an instruction to the bytecode chunk.
func (f *BytecodeMethod) AddInstruction(lineNumber int, op bytecode.OpCode, bytes ...byte) {
	f.LineInfoList.AddLineNumber(lineNumber)
	f.Instructions = append(f.Instructions, byte(op))
	f.Instructions = append(f.Instructions, bytes...)
}

// Add bytes to the bytecode chunk.
func (f *BytecodeMethod) AddBytes(bytes ...byte) {
	f.Instructions = append(f.Instructions, bytes...)
}

// Append two bytes to the bytecode chunk.
func (f *BytecodeMethod) AppendUint16(n uint16) {
	f.Instructions = binary.BigEndian.AppendUint16(f.Instructions, n)
}

// Append four bytes to the bytecode chunk.
func (f *BytecodeMethod) AppendUint32(n uint32) {
	f.Instructions = binary.BigEndian.AppendUint32(f.Instructions, n)
}

// Size of an integer.
type IntSize uint8

// Add a constant to the constant pool.
// Returns the index of the constant.
func (f *BytecodeMethod) AddValue(obj value.Value) (int, IntSize) {
	var id int
	switch obj.(type) {
	case value.String, value.Symbol, value.SmallInt, value.Int64, value.Int32, value.Int16,
		value.Int8, value.UInt64, value.UInt32, value.UInt16, value.UInt8,
		value.Float, value.Float32, value.Float64:
		if i := slices.Index(f.Values, obj); i != -1 {
			id = i
			break
		}
		id = len(f.Values)
		f.Values = append(f.Values, obj)
	default:
		id = len(f.Values)
		f.Values = append(f.Values, obj)
	}

	if id <= math.MaxUint8 {
		return id, bytecode.UINT8_SIZE
	}

	if id <= math.MaxUint16 {
		return id, bytecode.UINT16_SIZE
	}

	if id <= math.MaxUint32 {
		return id, bytecode.UINT32_SIZE
	}

	return id, bytecode.UINT64_SIZE
}

// Disassemble the bytecode chunk and write the
// output to stdout.
func (f *BytecodeMethod) DisassembleStdout() {
	f.Disassemble(os.Stdout)
}

// Disassemble the bytecode chunk and return a string
// containing the result.
func (f *BytecodeMethod) DisassembleString() (string, error) {
	var buffer strings.Builder
	err := f.Disassemble(&buffer)
	if err != nil {
		return buffer.String(), err
	}

	return buffer.String(), nil
}

// Disassemble the bytecode chunk and write the
// output to a writer.
func (f *BytecodeMethod) Disassemble(output io.Writer) error {
	fmt.Fprintf(output, "== Disassembly of %s at: %s ==\n\n", f.name.ToString(), f.Location.String())

	if len(f.Instructions) == 0 {
		return nil
	}

	var offset int
	var instructionIndex int
	for {
		result, err := f.DisassembleInstruction(output, offset, instructionIndex)
		if err != nil {
			return err
		}
		offset = result
		instructionIndex++
		if offset >= len(f.Instructions) {
			break
		}
	}

	for _, constant := range f.Values {
		fn, ok := constant.(*BytecodeMethod)
		if !ok {
			continue
		}
		fmt.Fprintln(output)
		fn.Disassemble(output)
	}

	return nil
}

func (f *BytecodeMethod) DisassembleInstruction(output io.Writer, offset, instructionIndex int) (int, error) {
	fmt.Fprintf(output, "%04d  ", offset)
	opcodeByte := f.Instructions[offset]
	opcode := bytecode.OpCode(opcodeByte)
	switch opcode {
	case bytecode.RETURN, bytecode.ADD, bytecode.SUBTRACT,
		bytecode.MULTIPLY, bytecode.DIVIDE, bytecode.EXPONENTIATE,
		bytecode.NEGATE, bytecode.NOT, bytecode.BITWISE_NOT,
		bytecode.TRUE, bytecode.FALSE, bytecode.NIL, bytecode.POP,
		bytecode.RBITSHIFT, bytecode.LBITSHIFT,
		bytecode.LOGIC_RBITSHIFT, bytecode.LOGIC_LBITSHIFT,
		bytecode.BITWISE_AND, bytecode.BITWISE_OR, bytecode.BITWISE_XOR, bytecode.MODULO,
		bytecode.EQUAL, bytecode.STRICT_EQUAL, bytecode.GREATER, bytecode.GREATER_EQUAL, bytecode.LESS, bytecode.LESS_EQUAL,
		bytecode.ROOT, bytecode.NOT_EQUAL, bytecode.STRICT_NOT_EQUAL,
		bytecode.CONSTANT_CONTAINER, bytecode.SELF, bytecode.DEF_MODULE, bytecode.DEF_METHOD,
		bytecode.UNDEFINED, bytecode.DEF_ANON_CLASS, bytecode.DEF_ANON_MODULE,
		bytecode.DEF_MIXIN, bytecode.DEF_ANON_MIXIN, bytecode.INCLUDE, bytecode.GET_SINGLETON,
		bytecode.DEF_ALIAS, bytecode.METHOD_CONTAINER, bytecode.COMPARE, bytecode.DOC_COMMENT,
		bytecode.DEF_GETTER, bytecode.DEF_SETTER, bytecode.DEF_SINGLETON, bytecode.RETURN_FIRST_ARG,
		bytecode.RETURN_SELF, bytecode.APPEND, bytecode.COPY, bytecode.SUBSCRIPT, bytecode.SUBSCRIPT_SET,
		bytecode.APPEND_AT, bytecode.GET_ITERATOR:
		return f.disassembleOneByteInstruction(output, opcode.String(), offset, instructionIndex), nil
	case bytecode.POP_N, bytecode.SET_LOCAL8, bytecode.GET_LOCAL8, bytecode.PREP_LOCALS8,
		bytecode.DEF_CLASS, bytecode.NEW_ARRAY_TUPLE8, bytecode.NEW_ARRAY_LIST8, bytecode.NEW_STRING8:
		return f.disassembleNumericOperands(output, 1, 1, offset, instructionIndex)
	case bytecode.PREP_LOCALS16, bytecode.SET_LOCAL16, bytecode.GET_LOCAL16, bytecode.JUMP_UNLESS, bytecode.JUMP,
		bytecode.JUMP_IF, bytecode.LOOP, bytecode.JUMP_IF_NIL, bytecode.JUMP_UNLESS_UNDEF, bytecode.FOR_IN:
		return f.disassembleNumericOperands(output, 1, 2, offset, instructionIndex)
	case bytecode.NEW_ARRAY_TUPLE32, bytecode.NEW_ARRAY_LIST32, bytecode.NEW_STRING32:
		return f.disassembleNumericOperands(output, 1, 4, offset, instructionIndex)
	case bytecode.LEAVE_SCOPE16:
		return f.disassembleNumericOperands(output, 2, 1, offset, instructionIndex)
	case bytecode.LEAVE_SCOPE32:
		return f.disassembleNumericOperands(output, 2, 2, offset, instructionIndex)
	case bytecode.LOAD_VALUE8, bytecode.GET_MOD_CONST8,
		bytecode.DEF_MOD_CONST8, bytecode.CALL_METHOD8,
		bytecode.CALL_FUNCTION8, bytecode.INSTANTIATE8,
		bytecode.GET_IVAR8, bytecode.SET_IVAR8:
		return f.disassembleConstant(output, 2, offset, instructionIndex)
	case bytecode.LOAD_VALUE16, bytecode.GET_MOD_CONST16,
		bytecode.DEF_MOD_CONST16, bytecode.CALL_METHOD16,
		bytecode.CALL_FUNCTION16, bytecode.INSTANTIATE16,
		bytecode.GET_IVAR16, bytecode.SET_IVAR16:
		return f.disassembleConstant(output, 3, offset, instructionIndex)
	case bytecode.LOAD_VALUE32, bytecode.GET_MOD_CONST32,
		bytecode.DEF_MOD_CONST32, bytecode.CALL_METHOD32,
		bytecode.CALL_FUNCTION32, bytecode.INSTANTIATE32,
		bytecode.GET_IVAR32, bytecode.SET_IVAR32:
		return f.disassembleConstant(output, 5, offset, instructionIndex)
	default:
		f.printLineNumber(output, instructionIndex)
		f.dumpBytes(output, offset, 1)
		fmt.Fprintf(output, "unknown operation %d (0x%X)\n", opcodeByte, opcodeByte)
		return offset + 1, fmt.Errorf("unknown operation %d (0x%X) at offset %d (0x%X)", opcodeByte, opcodeByte, offset, offset)
	}
}

func (f *BytecodeMethod) dumpBytes(output io.Writer, offset, count int) {
	for i := offset; i < offset+count; i++ {
		fmt.Fprintf(output, "%02X ", f.Instructions[i])
	}

	for i := count; i < bytecode.MaxInstructionByteCount; i++ {
		fmt.Fprint(output, "   ")
	}
}

func (f *BytecodeMethod) disassembleOneByteInstruction(output io.Writer, name string, offset, instructionIndex int) int {
	f.printLineNumber(output, instructionIndex)
	f.dumpBytes(output, offset, 1)
	fmt.Fprintln(output, name)
	return offset + 1
}

func (f *BytecodeMethod) disassembleNumericOperands(output io.Writer, operands, operandBytes, offset, instructionIndex int) (int, error) {
	bytes := 1 + operands*operandBytes
	if result, err := f.checkBytes(output, offset, instructionIndex, bytes); err != nil {
		return result, err
	}

	opcode := bytecode.OpCode(f.Instructions[offset])

	f.printLineNumber(output, instructionIndex)
	f.dumpBytes(output, offset, bytes)
	f.printOpCode(output, opcode)

	var readFunc intReadFunc
	switch operandBytes {
	case 8:
		readFunc = readUint64
	case 4:
		readFunc = readUint32
	case 2:
		readFunc = readUint16
	case 1:
		readFunc = readUint8
	default:
		panic(fmt.Sprintf("incorrect bytesize of operands: %d", operandBytes))
	}

	for i := 0; i < operands; i++ {
		a := readFunc(f.Instructions[offset+1+i*operandBytes : offset+1+(i+1)*operandBytes])
		f.printNumField(output, a)
	}
	fmt.Fprintln(output)

	return offset + bytes, nil
}

type intReadFunc func([]byte) uint64

func readUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func readUint32(b []byte) uint64 {
	return uint64(binary.BigEndian.Uint32(b))
}

func readUint16(b []byte) uint64 {
	return uint64(binary.BigEndian.Uint16(b))
}

func readUint8(b []byte) uint64 {
	return uint64(b[0])
}

func (f *BytecodeMethod) disassembleConstant(output io.Writer, byteLength, offset, instructionIndex int) (int, error) {
	opcode := bytecode.OpCode(f.Instructions[offset])

	if result, err := f.checkBytes(output, offset, instructionIndex, byteLength); err != nil {
		return result, err
	}

	var constantIndex int
	if byteLength == 2 {
		constantIndex = int(f.Instructions[offset+1])
	} else if byteLength == 3 {
		constantIndex = int(binary.BigEndian.Uint16(f.Instructions[offset+1 : offset+3]))
	} else if byteLength == 5 {
		constantIndex = int(binary.BigEndian.Uint32(f.Instructions[offset+1 : offset+5]))
	} else {
		panic(fmt.Sprintf("%d is not a valid byteLength for a value opcode!", byteLength))
	}

	f.printLineNumber(output, instructionIndex)
	f.dumpBytes(output, offset, byteLength)
	f.printOpCode(output, opcode)

	if constantIndex >= len(f.Values) {
		msg := fmt.Sprintf("invalid value index %d (0x%X)", constantIndex, constantIndex)
		fmt.Fprintln(output, msg)
		return offset + byteLength, fmt.Errorf(msg)
	}
	constant := f.Values[constantIndex]
	fmt.Fprintln(output, constant.Inspect())

	return offset + byteLength, nil
}

func (f *BytecodeMethod) checkBytes(output io.Writer, offset, instructionIndex, byteLength int) (int, error) {
	opcode := bytecode.OpCode(f.Instructions[offset])
	if len(f.Instructions)-offset >= byteLength {
		return 0, nil
	}
	f.printLineNumber(output, instructionIndex)
	f.dumpBytes(output, offset, len(f.Instructions)-offset)
	f.printOpCode(output, opcode)
	msg := "not enough bytes"
	fmt.Fprintln(output, msg)
	return len(f.Instructions) - 1, fmt.Errorf(msg)
}

func (f *BytecodeMethod) printLineNumber(output io.Writer, instructionIndex int) {
	fmt.Fprintf(output, "%-8s", f.getLineNumberString(instructionIndex))
}

func (f *BytecodeMethod) getLineNumberString(instructionIndex int) string {
	currentLineNumber := f.LineInfoList.GetLineNumber(instructionIndex)
	if instructionIndex == 0 {
		return fmt.Sprintf("%d", currentLineNumber)
	}

	previousLineNumber := f.LineInfoList.GetLineNumber(instructionIndex - 1)
	if previousLineNumber == currentLineNumber {
		return "|"
	}

	return fmt.Sprintf("%d", currentLineNumber)
}

func (f *BytecodeMethod) printOpCode(output io.Writer, opcode bytecode.OpCode) {
	fmt.Fprintf(output, "%-18s", opcode.String())
}

func (f *BytecodeMethod) printNumField(output io.Writer, n uint64) {
	fmt.Fprintf(output, "%-16d", n)
}
