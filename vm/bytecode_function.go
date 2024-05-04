package vm

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strings"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/regex/flag"
	"github.com/elk-language/elk/value"
)

// A single unit of Elk bytecode.
type BytecodeFunction struct {
	Instructions []byte
	Values       []value.Value // The value pool
	LineInfoList bytecode.LineInfoList
	Location     *position.Location
	Doc          value.Value
	CatchEntries []*CatchEntry
	UpvalueCount int

	name                   value.Symbol
	parameters             []value.Symbol
	optionalParameterCount int
	postRestParameterCount int
	namedRestParameter     bool
	sealed                 bool
}

func (b *BytecodeFunction) Name() value.Symbol {
	return b.name
}

func (b *BytecodeFunction) ParameterCount() int {
	return len(b.parameters)
}

func (b *BytecodeFunction) SetParameters(params []value.Symbol) {
	b.parameters = params
}

func (b *BytecodeFunction) Parameters() []value.Symbol {
	return b.parameters
}

func (b *BytecodeFunction) SetOptionalParameterCount(optParamCount int) {
	b.optionalParameterCount = optParamCount
}

func (b *BytecodeFunction) IncrementOptionalParameterCount() {
	b.optionalParameterCount++
}

func (b *BytecodeFunction) OptionalParameterCount() int {
	return b.optionalParameterCount
}

func (b *BytecodeFunction) SetPostRestParameterCount(postParamCount int) {
	b.postRestParameterCount = postParamCount
}

func (b *BytecodeFunction) IncrementPostRestParameterCount() {
	b.postRestParameterCount++
}

func (b *BytecodeFunction) PostRestParameterCount() int {
	return b.postRestParameterCount
}

func (b *BytecodeFunction) SetNamedRestParameter(present bool) {
	b.namedRestParameter = present
}

func (b *BytecodeFunction) NamedRestParameter() bool {
	return b.namedRestParameter
}

func (*BytecodeFunction) Class() *value.Class {
	return value.MethodClass
}

func (*BytecodeFunction) DirectClass() *value.Class {
	return value.MethodClass
}

func (*BytecodeFunction) SingletonClass() *value.Class {
	return nil
}

func (b *BytecodeFunction) IsSealed() bool {
	return b.sealed
}

func (b *BytecodeFunction) SetSealed() {
	b.sealed = true
}

func (b *BytecodeFunction) Copy() value.Value {
	return b
}

func (b *BytecodeFunction) Inspect() string {
	return fmt.Sprintf("Method{name: %s, type: :bytecode, location: %s}", b.name.Inspect(), b.Location.String())
}

func (*BytecodeFunction) InstanceVariables() value.SymbolMap {
	return nil
}

func (b *BytecodeFunction) FileName() string {
	if b.Location == nil {
		return ""
	}
	return b.Location.Filename
}

func (b *BytecodeFunction) GetLineNumber(ip int) int {
	return b.LineInfoList.GetLineNumber(ip)
}

// Create a new bytecode method.
func NewBytecodeFunctionSimple(name value.Symbol, instruct []byte, loc *position.Location) *BytecodeFunction {
	return &BytecodeFunction{
		Instructions:           instruct,
		Location:               loc,
		name:                   name,
		postRestParameterCount: -1,
	}
}

// Create a new bytecode method.
func NewBytecodeFunction(
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
) *BytecodeFunction {
	return &BytecodeFunction{
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

type BytecodeFunctionOption func(*BytecodeFunction)

func BytecodeFunctionWithName(name value.Symbol) BytecodeFunctionOption {
	return func(b *BytecodeFunction) {
		b.name = name
	}
}

func BytecodeFunctionWithStringName(name string) BytecodeFunctionOption {
	return func(b *BytecodeFunction) {
		b.name = value.ToSymbol(name)
	}
}

func BytecodeFunctionWithInstructions(instructs []byte) BytecodeFunctionOption {
	return func(b *BytecodeFunction) {
		b.Instructions = instructs
	}
}

func BytecodeFunctionWithLocation(loc *position.Location) BytecodeFunctionOption {
	return func(b *BytecodeFunction) {
		b.Location = loc
	}
}

func BytecodeFunctionWithLineInfoList(lineInfo bytecode.LineInfoList) BytecodeFunctionOption {
	return func(b *BytecodeFunction) {
		b.LineInfoList = lineInfo
	}
}

func BytecodeFunctionWithParameters(params []value.Symbol) BytecodeFunctionOption {
	return func(b *BytecodeFunction) {
		b.parameters = params
	}
}

func BytecodeFunctionWithOptionalParameters(optParams int) BytecodeFunctionOption {
	return func(b *BytecodeFunction) {
		b.optionalParameterCount = optParams
	}
}

func BytecodeFunctionWithPostParameters(postParams int) BytecodeFunctionOption {
	return func(b *BytecodeFunction) {
		b.postRestParameterCount = postParams
	}
}

func BytecodeFunctionWithPositionalRestParameter() BytecodeFunctionOption {
	return func(b *BytecodeFunction) {
		b.postRestParameterCount = 0
	}
}

func BytecodeFunctionWithNamedRestParameter() BytecodeFunctionOption {
	return func(b *BytecodeFunction) {
		b.namedRestParameter = true
	}
}

func BytecodeFunctionWithSealed() BytecodeFunctionOption {
	return func(b *BytecodeFunction) {
		b.sealed = true
	}
}

// Create a new bytecode method with options.
func NewBytecodeFunctionWithOptions(opts ...BytecodeFunctionOption) *BytecodeFunction {
	b := &BytecodeFunction{
		postRestParameterCount: -1,
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

// Create a new bytecode method.
func NewBytecodeFunctionNoParams(
	name value.Symbol,
	instruct []byte,
	loc *position.Location,
	lineInfo bytecode.LineInfoList,
	values []value.Value,
) *BytecodeFunction {
	return NewBytecodeFunction(name, instruct, loc, lineInfo, nil, 0, -1, false, false, values)
}

// Create a new bytecode method.
func NewBytecodeFunctionWithCatchEntries(
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
	catchEntries []*CatchEntry,
) *BytecodeFunction {
	return &BytecodeFunction{
		name:                   name,
		Instructions:           instruct,
		Location:               loc,
		LineInfoList:           lineInfo,
		CatchEntries:           catchEntries,
		parameters:             params,
		optionalParameterCount: optParamCount,
		postRestParameterCount: postRestParamCount,
		namedRestParameter:     namedRestParam,
		Values:                 values,
		sealed:                 sealed,
	}
}

// Add a parameter to the method.
func (f *BytecodeFunction) AddParameter(name value.Symbol) {
	f.parameters = append(f.parameters, name)
}

// Add a parameter to the method.
func (f *BytecodeFunction) AddParameterString(name string) {
	f.parameters = append(f.parameters, value.ToSymbol(name))
}

// Add an instruction to the bytecode chunk.
func (f *BytecodeFunction) AddInstruction(lineNumber int, op bytecode.OpCode, bytes ...byte) {
	f.LineInfoList.AddLineNumber(lineNumber, len(bytes)+1)
	f.Instructions = append(f.Instructions, byte(op))
	f.Instructions = append(f.Instructions, bytes...)
}

// Add bytes to the bytecode chunk.
func (f *BytecodeFunction) AddBytes(bytes ...byte) {
	f.Instructions = append(f.Instructions, bytes...)
}

// Append two bytes to the bytecode chunk.
func (f *BytecodeFunction) AppendUint16(n uint16) {
	f.Instructions = binary.BigEndian.AppendUint16(f.Instructions, n)
}

// Append four bytes to the bytecode chunk.
func (f *BytecodeFunction) AppendUint32(n uint32) {
	f.Instructions = binary.BigEndian.AppendUint32(f.Instructions, n)
}

// Size of an integer.
type IntSize uint8

// Add a constant to the constant pool.
// Returns the index of the constant.
func (f *BytecodeFunction) AddValue(obj value.Value) (int, IntSize) {
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
func (f *BytecodeFunction) DisassembleStdout() {
	f.Disassemble(os.Stdout)
}

// Disassemble the bytecode chunk and return a string
// containing the result.
func (f *BytecodeFunction) DisassembleString() (string, error) {
	var buffer strings.Builder
	err := f.Disassemble(&buffer)
	if err != nil {
		return buffer.String(), err
	}

	return buffer.String(), nil
}

// Disassemble the bytecode chunk and write the
// output to a writer.
func (f *BytecodeFunction) Disassemble(output io.Writer) error {
	fmt.Fprintf(output, "== Disassembly of %s at: %s ==\n\n", f.name.ToString(), f.Location.String())

	if len(f.CatchEntries) > 0 {
		fmt.Fprintln(output, "-- Catch entries --")
		for _, catchEntry := range f.CatchEntries {
			fmt.Fprintf(output, "%04d:%04d -> %04d", catchEntry.From, catchEntry.To, catchEntry.JumpAddress)
			if catchEntry.Finally {
				fmt.Fprint(output, " (finally)")
			}
			fmt.Fprintln(output)
		}
		fmt.Fprintln(output)
	}

	if len(f.Instructions) == 0 {
		return nil
	}

	var offset int
	for {
		result, err := f.DisassembleInstruction(output, offset)
		if err != nil {
			return err
		}
		offset = result
		if offset >= len(f.Instructions) {
			break
		}
	}

	for _, constant := range f.Values {
		fn, ok := constant.(*BytecodeFunction)
		if !ok {
			continue
		}
		fmt.Fprintln(output)
		fn.Disassemble(output)
	}

	return nil
}

func (f *BytecodeFunction) DisassembleInstruction(output io.Writer, offset int) (int, error) {
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
		bytecode.APPEND_AT, bytecode.GET_ITERATOR, bytecode.MAP_SET, bytecode.LAX_EQUAL, bytecode.LAX_NOT_EQUAL,
		bytecode.BITWISE_AND_NOT, bytecode.UNARY_PLUS, bytecode.INCREMENT, bytecode.DECREMENT, bytecode.DUP,
		bytecode.SWAP, bytecode.INSTANCE_OF, bytecode.IS_A, bytecode.POP_SKIP_ONE, bytecode.INSPECT_STACK,
		bytecode.THROW, bytecode.RETHROW, bytecode.POP_ALL, bytecode.RETURN_FINALLY, bytecode.JUMP_TO_FINALLY:
		return f.disassembleOneByteInstruction(output, opcode.String(), offset), nil
	case bytecode.POP_N, bytecode.SET_LOCAL8, bytecode.GET_LOCAL8, bytecode.PREP_LOCALS8,
		bytecode.DEF_CLASS, bytecode.NEW_ARRAY_TUPLE8, bytecode.NEW_ARRAY_LIST8, bytecode.NEW_STRING8,
		bytecode.NEW_HASH_MAP8, bytecode.NEW_HASH_RECORD8, bytecode.DUP_N, bytecode.POP_N_SKIP_ONE, bytecode.NEW_SYMBOL8,
		bytecode.NEW_HASH_SET8, bytecode.SET_UPVALUE8, bytecode.GET_UPVALUE8:
		return f.disassembleNumericOperands(output, 1, 1, offset)
	case bytecode.PREP_LOCALS16, bytecode.SET_LOCAL16, bytecode.GET_LOCAL16, bytecode.JUMP_UNLESS, bytecode.JUMP,
		bytecode.JUMP_IF, bytecode.LOOP, bytecode.JUMP_IF_NIL, bytecode.JUMP_UNLESS_UNDEF, bytecode.FOR_IN,
		bytecode.SET_UPVALUE16, bytecode.GET_UPVALUE16:
		return f.disassembleNumericOperands(output, 1, 2, offset)
	case bytecode.NEW_ARRAY_TUPLE32, bytecode.NEW_ARRAY_LIST32, bytecode.NEW_STRING32,
		bytecode.NEW_HASH_MAP32, bytecode.NEW_HASH_RECORD32, bytecode.NEW_SYMBOL32,
		bytecode.NEW_HASH_SET32:
		return f.disassembleNumericOperands(output, 1, 4, offset)
	case bytecode.LEAVE_SCOPE16:
		return f.disassembleNumericOperands(output, 2, 1, offset)
	case bytecode.LEAVE_SCOPE32:
		return f.disassembleNumericOperands(output, 2, 2, offset)
	case bytecode.NEW_REGEX8:
		return f.disassembleNewRegex(output, 1, offset)
	case bytecode.NEW_REGEX32:
		return f.disassembleNewRegex(output, 4, offset)
	case bytecode.CLOSURE:
		return f.disassembleClosure(output, offset)
	case bytecode.NEW_RANGE:
		return f.disassembleNewRange(output, offset)
	case bytecode.LOAD_VALUE8, bytecode.GET_MOD_CONST8,
		bytecode.DEF_MOD_CONST8, bytecode.CALL_METHOD8,
		bytecode.CALL_SELF8, bytecode.INSTANTIATE8,
		bytecode.GET_IVAR8, bytecode.SET_IVAR8, bytecode.CALL_PATTERN8,
		bytecode.CALL8:
		return f.disassembleConstant(output, 2, offset)
	case bytecode.LOAD_VALUE16, bytecode.GET_MOD_CONST16,
		bytecode.DEF_MOD_CONST16, bytecode.CALL_METHOD16,
		bytecode.CALL_SELF16, bytecode.INSTANTIATE16,
		bytecode.GET_IVAR16, bytecode.SET_IVAR16, bytecode.CALL_PATTERN16,
		bytecode.CALL16:
		return f.disassembleConstant(output, 3, offset)
	case bytecode.LOAD_VALUE32, bytecode.GET_MOD_CONST32,
		bytecode.DEF_MOD_CONST32, bytecode.CALL_METHOD32,
		bytecode.CALL_SELF32, bytecode.INSTANTIATE32,
		bytecode.GET_IVAR32, bytecode.SET_IVAR32, bytecode.CALL_PATTERN32,
		bytecode.CALL32:
		return f.disassembleConstant(output, 5, offset)
	default:
		f.printLineNumber(output, offset)
		f.dumpBytes(output, offset, 1)
		fmt.Fprintf(output, "unknown operation %d (0x%X)\n", opcodeByte, opcodeByte)
		return offset + 1, fmt.Errorf("unknown operation %d (0x%X) at offset %d (0x%X)", opcodeByte, opcodeByte, offset, offset)
	}
}

func (f *BytecodeFunction) dumpBytes(output io.Writer, offset, count int) {
	for i := offset; i < offset+count; i++ {
		fmt.Fprintf(output, "%02X ", f.Instructions[i])
	}

	for i := count; i < bytecode.MaxInstructionByteCount; i++ {
		fmt.Fprint(output, "   ")
	}
}

func (f *BytecodeFunction) disassembleOneByteInstruction(output io.Writer, name string, offset int) int {
	f.printLineNumber(output, offset)
	f.dumpBytes(output, offset, 1)
	fmt.Fprintln(output, name)
	return offset + 1
}

func (f *BytecodeFunction) disassembleNumericOperands(output io.Writer, operands, operandBytes, offset int) (int, error) {
	bytes := 1 + operands*operandBytes
	if result, err := f.checkBytes(output, offset, bytes); err != nil {
		return result, err
	}

	opcode := bytecode.OpCode(f.Instructions[offset])

	f.printLineNumber(output, offset)
	f.dumpBytes(output, offset, bytes)
	f.printOpCode(output, opcode)

	readFunc := readFuncForBytes(operandBytes)

	for i := 0; i < operands; i++ {
		a := readFunc(f.Instructions[offset+1+i*operandBytes : offset+1+(i+1)*operandBytes])
		f.printNumField(output, a)
	}
	fmt.Fprintln(output)

	return offset + bytes, nil
}

func readFuncForBytes(bytes int) intReadFunc {
	switch bytes {
	case 8:
		return readUint64
	case 4:
		return readUint32
	case 2:
		return readUint16
	case 1:
		return readUint8
	default:
		panic(fmt.Sprintf("incorrect bytesize of operands: %d", bytes))
	}
}

func (f *BytecodeFunction) disassembleClosure(output io.Writer, offset int) (int, error) {
	bytes := 1
	if result, err := f.checkBytes(output, offset, bytes); err != nil {
		return result, err
	}

	opcode := bytecode.OpCode(f.Instructions[offset])

	f.printLineNumber(output, offset)
	f.dumpBytes(output, offset, bytes)
	f.printOpCode(output, opcode)

	for {
		fmt.Fprintln(output)
		if len(f.Instructions)-1 < offset+bytes {
			break
		}
		flags := bitfield.BitField8FromInt(f.Instructions[offset+bytes])
		fmt.Fprintf(output, "%04d  ", offset+bytes)
		var upIndex int
		if flags.HasFlag(UpvalueLongIndexFlag) {
			upIndex = readUint16(f.Instructions[offset+bytes+1:])
			bytes += 2
		} else {
			upIndex = readUint8(f.Instructions[offset+bytes+1:])
			bytes++
		}

		f.printLineNumber(output, offset)
		f.dumpBytes(output, offset+bytes, 2)
		fmt.Fprintf(output, "%-18s", "|")
		bytes += 2
	}
	return offset + bytes, nil
}

func (f *BytecodeFunction) disassembleNewRange(output io.Writer, offset int) (int, error) {
	bytes := 2
	if result, err := f.checkBytes(output, offset, bytes); err != nil {
		return result, err
	}

	opcode := bytecode.OpCode(f.Instructions[offset])

	f.printLineNumber(output, offset)
	f.dumpBytes(output, offset, bytes)
	f.printOpCode(output, opcode)

	flagByte := f.Instructions[offset+1]
	var argString string

	switch flagByte {
	case bytecode.CLOSED_RANGE_FLAG:
		argString = "x...y"
	case bytecode.LEFT_OPEN_RANGE_FLAG:
		argString = "x<..y"
	case bytecode.RIGHT_OPEN_RANGE_FLAG:
		argString = "x..<y"
	case bytecode.OPEN_RANGE_FLAG:
		argString = "x<.<y"
	case bytecode.BEGINLESS_CLOSED_RANGE_FLAG:
		argString = "...x"
	case bytecode.BEGINLESS_OPEN_RANGE_FLAG:
		argString = "..<x"
	case bytecode.ENDLESS_CLOSED_RANGE_FLAG:
		argString = "x..."
	case bytecode.ENDLESS_OPEN_RANGE_FLAG:
		argString = "x<.."
	}
	fmt.Fprintf(output, "%d (%s)", flagByte, argString)
	fmt.Fprintln(output)

	return offset + bytes, nil
}

func (f *BytecodeFunction) disassembleNewRegex(output io.Writer, sizeBytes, offset int) (int, error) {
	flagBytes := 1
	bytes := 1 + flagBytes + sizeBytes
	if result, err := f.checkBytes(output, offset, bytes); err != nil {
		return result, err
	}

	opcode := bytecode.OpCode(f.Instructions[offset])

	f.printLineNumber(output, offset)
	f.dumpBytes(output, offset, bytes)
	f.printOpCode(output, opcode)

	flagByte := readUint8(f.Instructions[offset+1 : offset+1+flagBytes])
	flags := bitfield.BitField8FromInt(flagByte)
	fmt.Fprintf(output, "%-16s", fmt.Sprintf("%d (%s)", flagByte, flag.ToStringWithDisabledFlags(flags)))

	sizeReadFunc := readFuncForBytes(sizeBytes)
	size := sizeReadFunc(f.Instructions[offset+1+flagBytes : offset+1+flagBytes+sizeBytes])
	f.printNumField(output, size)
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

func (f *BytecodeFunction) disassembleConstant(output io.Writer, byteLength, offset int) (int, error) {
	opcode := bytecode.OpCode(f.Instructions[offset])

	if result, err := f.checkBytes(output, offset, byteLength); err != nil {
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

	f.printLineNumber(output, offset)
	f.dumpBytes(output, offset, byteLength)
	f.printOpCode(output, opcode)

	if constantIndex >= len(f.Values) {
		msg := fmt.Sprintf("invalid value index %d (0x%X)", constantIndex, constantIndex)
		fmt.Fprintln(output, msg)
		return offset + byteLength, fmt.Errorf(msg)
	}
	constant := f.Values[constantIndex]
	fmt.Fprintf(output, "%d (%s)\n", constantIndex, constant.Inspect())

	return offset + byteLength, nil
}

func (f *BytecodeFunction) checkBytes(output io.Writer, offset, byteLength int) (int, error) {
	opcode := bytecode.OpCode(f.Instructions[offset])
	if len(f.Instructions)-offset >= byteLength {
		return 0, nil
	}
	f.printLineNumber(output, offset)
	f.dumpBytes(output, offset, len(f.Instructions)-offset)
	f.printOpCode(output, opcode)
	msg := "not enough bytes"
	fmt.Fprintln(output, msg)
	return len(f.Instructions) - 1, fmt.Errorf(msg)
}

func (f *BytecodeFunction) printLineNumber(output io.Writer, offset int) {
	fmt.Fprintf(output, "%-8s", f.getLineNumberString(offset))
}

func (f *BytecodeFunction) getLineNumberString(offset int) string {
	currentLineNumber := f.LineInfoList.GetLineNumber(offset)
	if offset == 0 {
		return fmt.Sprintf("%d", currentLineNumber)
	}

	previousLineNumber := f.LineInfoList.GetLineNumber(offset - 1)
	if previousLineNumber == currentLineNumber {
		return "|"
	}

	return fmt.Sprintf("%d", currentLineNumber)
}

func (f *BytecodeFunction) printOpCode(output io.Writer, opcode bytecode.OpCode) {
	fmt.Fprintf(output, "%-18s", opcode.String())
}

func (f *BytecodeFunction) printNumField(output io.Writer, n uint64) {
	fmt.Fprintf(output, "%-16d", n)
}
