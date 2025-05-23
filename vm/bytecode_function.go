package vm

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strings"
	"unsafe"

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
	parameterCount         int
	optionalParameterCount int
}

func (b *BytecodeFunction) ipAddRaw(n uintptr) uintptr {
	inst := b.Instructions
	return uintptr(unsafe.Pointer(&inst[0])) + n
}

func (b *BytecodeFunction) Name() value.Symbol {
	return b.name
}

func (b *BytecodeFunction) ParameterCount() int {
	return b.parameterCount
}

func (b *BytecodeFunction) SetParameterCount(paramCount int) {
	b.parameterCount = paramCount
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

func (*BytecodeFunction) Class() *value.Class {
	return value.MethodClass
}

func (*BytecodeFunction) DirectClass() *value.Class {
	return value.MethodClass
}

func (*BytecodeFunction) SingletonClass() *value.Class {
	return nil
}

func (b *BytecodeFunction) Copy() value.Reference {
	return b
}

func (b *BytecodeFunction) Inspect() string {
	return fmt.Sprintf("Method{name: %s, type: :bytecode, location: %s}", b.name.Inspect(), b.Location.String())
}

func (b *BytecodeFunction) Error() string {
	return b.Inspect()
}

func (*BytecodeFunction) InstanceVariables() value.SymbolMap {
	return nil
}

func (b *BytecodeFunction) FileName() string {
	if b.Location == nil {
		return ""
	}
	return b.Location.FilePath
}

func (b *BytecodeFunction) GetLineNumber(ip int) int {
	return b.LineInfoList.GetLineNumber(ip)
}

// Create a new bytecode method.
func NewBytecodeFunctionSimple(name value.Symbol, instruct []byte, loc *position.Location) *BytecodeFunction {
	return &BytecodeFunction{
		Instructions: instruct,
		Location:     loc,
		name:         name,
	}
}

// Create a new bytecode method.
func NewBytecodeFunction(
	name value.Symbol,
	instruct []byte,
	loc *position.Location,
	lineInfo bytecode.LineInfoList,
	paramCount int,
	optParamCount int,
	values []value.Value,
) *BytecodeFunction {
	return &BytecodeFunction{
		name:                   name,
		Instructions:           instruct,
		Location:               loc,
		LineInfoList:           lineInfo,
		parameterCount:         paramCount,
		optionalParameterCount: optParamCount,
		Values:                 values,
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

func BytecodeFunctionWithParameters(params int) BytecodeFunctionOption {
	return func(b *BytecodeFunction) {
		b.parameterCount = params
	}
}

func BytecodeFunctionWithOptionalParameters(optParams int) BytecodeFunctionOption {
	return func(b *BytecodeFunction) {
		b.optionalParameterCount = optParams
	}
}

// Create a new bytecode method with options.
func NewBytecodeFunctionWithOptions(opts ...BytecodeFunctionOption) *BytecodeFunction {
	b := &BytecodeFunction{}

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
	return NewBytecodeFunction(name, instruct, loc, lineInfo, 0, 0, values)
}

// Create a new bytecode method.
func NewBytecodeFunctionWithCatchEntries(
	name value.Symbol,
	instruct []byte,
	loc *position.Location,
	lineInfo bytecode.LineInfoList,
	params int,
	optParamCount int,
	values []value.Value,
	catchEntries []*CatchEntry,
) *BytecodeFunction {
	return &BytecodeFunction{
		name:                   name,
		Instructions:           instruct,
		Location:               loc,
		LineInfoList:           lineInfo,
		CatchEntries:           catchEntries,
		parameterCount:         params,
		optionalParameterCount: optParamCount,
		Values:                 values,
	}
}

// Create a new bytecode method.
func NewBytecodeFunctionWithUpvalues(
	name value.Symbol,
	instruct []byte,
	loc *position.Location,
	lineInfo bytecode.LineInfoList,
	params int,
	optParamCount int,
	values []value.Value,
	upvalueCount int,
) *BytecodeFunction {
	return &BytecodeFunction{
		name:                   name,
		Instructions:           instruct,
		Location:               loc,
		LineInfoList:           lineInfo,
		UpvalueCount:           upvalueCount,
		parameterCount:         params,
		optionalParameterCount: optParamCount,
		Values:                 values,
	}
}

// Add a parameter to the method.
func (f *BytecodeFunction) AddParameter() {
	f.parameterCount++
}

// Add an instruction to the bytecode chunk.
func (f *BytecodeFunction) AddInstruction(lineNumber int, op bytecode.OpCode, bytes ...byte) {
	f.LineInfoList.AddLineNumber(lineNumber, len(bytes)+1)
	f.Instructions = append(f.Instructions, byte(op))
	f.Instructions = append(f.Instructions, bytes...)
}

// Add an instruction to the bytecode chunk.
func (f *BytecodeFunction) RemoveByte() {
	if len(f.Instructions) < 1 {
		panic("cannot remove a byte from an empty bytecode function")
	}

	f.LineInfoList.RemoveByte()
	f.Instructions = f.Instructions[:len(f.Instructions)-1]
}

// Add bytes to the bytecode chunk.
func (f *BytecodeFunction) AddBytes(bytes ...byte) {
	f.LineInfoList.AddBytesToLastLine(len(bytes))
	f.Instructions = append(f.Instructions, bytes...)
}

// Append two bytes to the bytecode chunk.
func (f *BytecodeFunction) AppendUint16(n uint16) {
	f.LineInfoList.AddBytesToLastLine(2)
	f.Instructions = binary.BigEndian.AppendUint16(f.Instructions, n)
}

// Append four bytes to the bytecode chunk.
func (f *BytecodeFunction) AppendUint32(n uint32) {
	f.LineInfoList.AddBytesToLastLine(4)
	f.Instructions = binary.BigEndian.AppendUint32(f.Instructions, n)
}

// Size of an integer.
type IntSize uint8

// Add a value to the value pool.
// Returns the index of the constant.
func (f *BytecodeFunction) AddValue(obj value.Value) (int, IntSize) {
	var id int
	if obj.IsReference() {
		objRef := obj.AsReference()
		i := -1
		for j, value := range f.Values {
			if !value.IsReference() {
				continue
			}

			if value.AsReference() == objRef {
				i = j
				id = j
				break
			}
		}
		if i == -1 {
			id = len(f.Values)
			f.Values = append(f.Values, obj)
		}
	} else {
		switch obj.ValueFlag() {
		case value.SYMBOL_FLAG, value.SMALL_INT_FLAG,
			value.INT64_FLAG, value.INT32_FLAG, value.INT16_FLAG, value.INT8_FLAG,
			value.UINT64_FLAG, value.UINT32_FLAG, value.UINT16_FLAG, value.UINT8_FLAG,
			value.FLOAT_FLAG, value.FLOAT32_FLAG, value.FLOAT64_FLAG:
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
	}

	if id <= math.MaxUint8 {
		return id, bytecode.UINT8_SIZE
	}

	if id <= math.MaxUint16 {
		return id, bytecode.UINT16_SIZE
	}

	if int64(id) <= math.MaxUint32 {
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
	// pp.Fprintln(output, f)
	fmt.Fprintf(output, "== Disassembly of %s at: %s ==\n\n", f.name.String(), f.Location.String())

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
		fn, ok := constant.SafeAsReference().(*BytecodeFunction)
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
		bytecode.RBITSHIFT, bytecode.LBITSHIFT, bytecode.NOOP,
		bytecode.LOGIC_RBITSHIFT, bytecode.LOGIC_LBITSHIFT,
		bytecode.BITWISE_AND, bytecode.BITWISE_OR, bytecode.BITWISE_XOR, bytecode.MODULO,
		bytecode.EQUAL, bytecode.STRICT_EQUAL, bytecode.GREATER, bytecode.GREATER_EQUAL, bytecode.LESS, bytecode.LESS_EQUAL,
		bytecode.NOT_EQUAL, bytecode.STRICT_NOT_EQUAL, bytecode.SELF, bytecode.INIT_NAMESPACE, bytecode.DEF_METHOD,
		bytecode.UNDEFINED, bytecode.INCLUDE, bytecode.GET_SINGLETON, bytecode.GET_CLASS, bytecode.COMPARE, bytecode.DOC_COMMENT,
		bytecode.DEF_GETTER, bytecode.DEF_SETTER, bytecode.RETURN_FIRST_ARG,
		bytecode.RETURN_SELF, bytecode.APPEND, bytecode.COPY, bytecode.SUBSCRIPT, bytecode.SUBSCRIPT_SET,
		bytecode.APPEND_AT, bytecode.GET_ITERATOR, bytecode.MAP_SET, bytecode.LAX_EQUAL, bytecode.LAX_NOT_EQUAL,
		bytecode.BITWISE_AND_NOT, bytecode.UNARY_PLUS, bytecode.INCREMENT, bytecode.DECREMENT, bytecode.DUP,
		bytecode.SWAP, bytecode.INSTANCE_OF, bytecode.IS_A, bytecode.POP_SKIP_ONE, bytecode.INSPECT_STACK,
		bytecode.THROW, bytecode.RETHROW, bytecode.RETURN_FINALLY, bytecode.JUMP_TO_FINALLY,
		bytecode.MUST, bytecode.AS, bytecode.SET_SUPERCLASS, bytecode.DEF_CONST, bytecode.EXEC, bytecode.DEF_METHOD_ALIAS,
		bytecode.INT_M1, bytecode.INT_0, bytecode.INT_1, bytecode.INT_2, bytecode.INT_3, bytecode.INT_4, bytecode.INT_5,
		bytecode.FLOAT_0, bytecode.FLOAT_1, bytecode.FLOAT_2,
		bytecode.LESS_EQUAL_INT, bytecode.ADD_INT, bytecode.SUBTRACT_INT,
		bytecode.GET_LOCAL_1, bytecode.GET_LOCAL_2, bytecode.GET_LOCAL_3, bytecode.GET_LOCAL_4,
		bytecode.SET_LOCAL_1, bytecode.SET_LOCAL_2, bytecode.SET_LOCAL_3, bytecode.SET_LOCAL_4,
		bytecode.GET_UPVALUE_0, bytecode.GET_UPVALUE_1,
		bytecode.SET_UPVALUE_0, bytecode.SET_UPVALUE_1, bytecode.SET_UPVALUE8, bytecode.SET_UPVALUE16,
		bytecode.POP_2, bytecode.POP_2_SKIP_ONE, bytecode.DUP_2,
		bytecode.ADD_FLOAT, bytecode.SUBTRACT_FLOAT, bytecode.MULTIPLY_INT, bytecode.MULTIPLY_FLOAT,
		bytecode.DIVIDE_INT, bytecode.DIVIDE_FLOAT, bytecode.EXPONENTIATE_INT, bytecode.NEGATE_INT, bytecode.NEGATE_FLOAT,
		bytecode.RBITSHIFT_INT, bytecode.LBITSHIFT_INT, bytecode.BITWISE_AND_INT, bytecode.BITWISE_OR_INT,
		bytecode.BITWISE_XOR_INT, bytecode.MODULO_INT, bytecode.MODULO_FLOAT, bytecode.EQUAL_INT, bytecode.EQUAL_FLOAT,
		bytecode.GREATER_INT, bytecode.GREATER_FLOAT, bytecode.GREATER_EQUAL_I, bytecode.GREATER_EQUAL_F,
		bytecode.LESS_INT, bytecode.LESS_FLOAT, bytecode.LESS_EQUAL_FLOAT, bytecode.NOT_EQUAL_INT, bytecode.NOT_EQUAL_FLOAT,
		bytecode.INCREMENT_INT, bytecode.DECREMENT_INT, bytecode.CLOSE_UPVALUE_1, bytecode.CLOSE_UPVALUE_2, bytecode.CLOSE_UPVALUE_3,
		bytecode.GENERATOR, bytecode.YIELD, bytecode.STOP_ITERATION, bytecode.GO, bytecode.DUP_SECOND,
		bytecode.PROMISE, bytecode.AWAIT, bytecode.AWAIT_RESULT:
		return f.disassembleOneByteInstruction(output, opcode.String(), offset), nil
	case bytecode.SET_LOCAL8, bytecode.GET_LOCAL8, bytecode.PREP_LOCALS8,
		bytecode.NEW_ARRAY_TUPLE8, bytecode.NEW_ARRAY_LIST8, bytecode.NEW_STRING8,
		bytecode.NEW_HASH_MAP8, bytecode.NEW_HASH_RECORD8, bytecode.NEW_SYMBOL8,
		bytecode.NEW_HASH_SET8, bytecode.GET_UPVALUE8, bytecode.CLOSE_UPVALUE8,
		bytecode.INSTANTIATE8, bytecode.LOAD_UINT64_8,
		bytecode.LOAD_UINT32_8, bytecode.LOAD_UINT16_8,
		bytecode.LOAD_UINT8:
		return f.disassembleUnsignedNumericOperands(output, 1, 1, offset)
	case bytecode.LOAD_INT_8, bytecode.LOAD_INT64_8,
		bytecode.LOAD_INT32_8, bytecode.LOAD_INT16_8, bytecode.LOAD_INT8:
		return f.disassembleSignedNumericOperands(output, 1, 1, offset)
	case bytecode.LOAD_CHAR_8:
		return f.disassembleChar(output, offset)
	case bytecode.PREP_LOCALS16, bytecode.SET_LOCAL16, bytecode.GET_LOCAL16, bytecode.JUMP_UNLESS, bytecode.JUMP,
		bytecode.JUMP_IF, bytecode.LOOP, bytecode.JUMP_IF_NIL, bytecode.JUMP_UNLESS_UNP, bytecode.FOR_IN_BUILTIN,
		bytecode.FOR_IN, bytecode.GET_UPVALUE16, bytecode.CLOSE_UPVALUE16,
		bytecode.INSTANTIATE16, bytecode.NEW_ARRAY_TUPLE16, bytecode.NEW_ARRAY_LIST16, bytecode.NEW_STRING16,
		bytecode.NEW_HASH_MAP16, bytecode.NEW_HASH_RECORD16, bytecode.NEW_SYMBOL16,
		bytecode.NEW_HASH_SET16, bytecode.JUMP_IF_IEQ, bytecode.JUMP_UNLESS_IEQ, bytecode.JUMP_UNLESS_IGE,
		bytecode.JUMP_UNLESS_IGT, bytecode.JUMP_UNLESS_ILT, bytecode.JUMP_UNLESS_ILE, bytecode.JUMP_UNLESS_NIL,
		bytecode.JUMP_IF_NP, bytecode.JUMP_UNLESS_NP, bytecode.JUMP_IF_NIL_NP, bytecode.JUMP_UNLESS_NNP,
		bytecode.JUMP_IF_EQ, bytecode.JUMP_UNLESS_EQ, bytecode.JUMP_UNLESS_GE,
		bytecode.JUMP_UNLESS_GT, bytecode.JUMP_UNLESS_LT, bytecode.JUMP_UNLESS_LE, bytecode.JUMP_UNLESS_UNDEF:
		return f.disassembleUnsignedNumericOperands(output, 1, 2, offset)
	case bytecode.LOAD_INT_16:
		return f.disassembleSignedNumericOperands(output, 1, 2, offset)
	case bytecode.LEAVE_SCOPE16:
		return f.disassembleUnsignedNumericOperands(output, 2, 1, offset)
	case bytecode.LEAVE_SCOPE32:
		return f.disassembleUnsignedNumericOperands(output, 2, 2, offset)
	case bytecode.DEF_NAMESPACE:
		return f.disassembleDefNamespace(output, offset)
	case bytecode.NEW_REGEX8:
		return f.disassembleNewRegex(output, 1, offset)
	case bytecode.NEW_REGEX16:
		return f.disassembleNewRegex(output, 2, offset)
	case bytecode.CLOSURE:
		return f.disassembleClosure(output, offset)
	case bytecode.NEW_RANGE:
		return f.disassembleNewRange(output, offset)
	case bytecode.LOAD_VALUE8, bytecode.CALL_METHOD8, bytecode.CALL_METHOD_TCO8,
		bytecode.CALL_SELF8, bytecode.CALL_SELF_TCO8,
		bytecode.GET_IVAR8, bytecode.SET_IVAR8,
		bytecode.CALL8, bytecode.GET_CONST8, bytecode.NEXT8:
		return f.disassembleValue(output, 2, offset)
	case bytecode.LOAD_VALUE_0:
		return f._disassembleValue(output, 1, 0, offset)
	case bytecode.LOAD_VALUE_1:
		return f._disassembleValue(output, 1, 1, offset)
	case bytecode.LOAD_VALUE_2:
		return f._disassembleValue(output, 1, 2, offset)
	case bytecode.LOAD_VALUE_3:
		return f._disassembleValue(output, 1, 3, offset)
	case bytecode.LOAD_VALUE16, bytecode.CALL_METHOD16, bytecode.CALL_METHOD_TCO16,
		bytecode.CALL_SELF16, bytecode.CALL_SELF_TCO16,
		bytecode.GET_IVAR16, bytecode.SET_IVAR16,
		bytecode.CALL16, bytecode.GET_CONST16, bytecode.NEXT16:
		return f.disassembleValue(output, 3, offset)
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

func (f *BytecodeFunction) disassembleChar(output io.Writer, offset int) (int, error) {
	operandBytes := 1
	bytes := 1 + operandBytes
	if result, err := f.checkBytes(output, offset, bytes); err != nil {
		return result, err
	}

	opcode := bytecode.OpCode(f.Instructions[offset])

	f.printLineNumber(output, offset)
	f.dumpBytes(output, offset, bytes)
	f.printOpCode(output, opcode)

	a := rune(f.Instructions[offset+1])
	fmt.Fprintf(output, "%-16c", a)
	fmt.Fprintln(output)

	return offset + bytes, nil
}

func (f *BytecodeFunction) disassembleUnsignedNumericOperands(output io.Writer, operands, operandBytes, offset int) (int, error) {
	bytes := 1 + operands*operandBytes
	if result, err := f.checkBytes(output, offset, bytes); err != nil {
		return result, err
	}

	opcode := bytecode.OpCode(f.Instructions[offset])

	f.printLineNumber(output, offset)
	f.dumpBytes(output, offset, bytes)
	f.printOpCode(output, opcode)

	readFunc := readFuncForUnsignedBytes(operandBytes)

	for i := 0; i < operands; i++ {
		a := readFunc(f.Instructions[offset+1+i*operandBytes : offset+1+(i+1)*operandBytes])
		printNumField(output, a)
	}
	fmt.Fprintln(output)

	return offset + bytes, nil
}

func (f *BytecodeFunction) disassembleSignedNumericOperands(output io.Writer, operands, operandBytes, offset int) (int, error) {
	bytes := 1 + operands*operandBytes
	if result, err := f.checkBytes(output, offset, bytes); err != nil {
		return result, err
	}

	opcode := bytecode.OpCode(f.Instructions[offset])

	f.printLineNumber(output, offset)
	f.dumpBytes(output, offset, bytes)
	f.printOpCode(output, opcode)

	readFunc := readFuncForSignedBytes(operandBytes)

	for i := 0; i < operands; i++ {
		a := readFunc(f.Instructions[offset+1+i*operandBytes : offset+1+(i+1)*operandBytes])
		printNumField(output, a)
	}
	fmt.Fprintln(output)

	return offset + bytes, nil
}

func readFuncForUnsignedBytes(bytes int) unsignedIntReadFunc {
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

func readFuncForSignedBytes(bytes int) signedIntReadFunc {
	switch bytes {
	case 8:
		return readInt64
	case 4:
		return readInt32
	case 2:
		return readInt16
	case 1:
		return readInt8
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
		startIndex := offset + bytes
		if len(f.Instructions)-1 < startIndex {
			break
		}
		flagsByte := f.Instructions[startIndex]
		fmt.Fprintf(output, "%04d  ", startIndex)
		if flagsByte == ClosureTerminatorFlag {
			f.printLineNumber(output, startIndex)
			f.dumpBytes(output, startIndex, 1)
			fmt.Fprintf(output, "%-18s", "|")
			fmt.Fprintln(output, "terminator")
			bytes++
			break
		}
		flags := bitfield.BitField8FromInt(flagsByte)
		var upIndex int
		var upvalueBytes int
		if flags.HasFlag(UpvalueLongIndexFlag) {
			upIndex = int(readUint16(f.Instructions[offset+bytes+1:]))
			bytes += 2
			upvalueBytes = 3
		} else {
			upIndex = int(readUint8(f.Instructions[offset+bytes+1:]))
			bytes++
			upvalueBytes = 2
		}

		var title string
		if flags.HasFlag(UpvalueLocalFlag) {
			title = "local"
		} else {
			title = "upvalue"
		}

		f.printLineNumber(output, startIndex)
		f.dumpBytes(output, startIndex, upvalueBytes)
		fmt.Fprintf(output, "%-18s", "|")
		fmt.Fprintf(output, "%s %d", title, upIndex)
		bytes++
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

func (f *BytecodeFunction) disassembleDefNamespace(output io.Writer, offset int) (int, error) {
	bytes := 2
	if result, err := f.checkBytes(output, offset, bytes); err != nil {
		return result, err
	}

	opcode := bytecode.OpCode(f.Instructions[offset])

	f.printLineNumber(output, offset)
	f.dumpBytes(output, offset, bytes)
	f.printOpCode(output, opcode)

	typeByte := readUint8(f.Instructions[offset+1 : offset+2])
	var namespaceType string
	switch typeByte {
	case 0:
		namespaceType = "module"
	case 1:
		namespaceType = "class"
	case 2:
		namespaceType = "mixin"
	case 3:
		namespaceType = "interface"
	default:
		return offset + bytes, fmt.Errorf("invalid namespace byte %d", typeByte)
	}
	fmt.Fprintf(output, "%-16s\n", fmt.Sprintf("%d (%s)", typeByte, namespaceType))

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

	sizeReadFunc := readFuncForUnsignedBytes(sizeBytes)
	size := sizeReadFunc(f.Instructions[offset+1+flagBytes : offset+1+flagBytes+sizeBytes])
	printNumField(output, size)
	fmt.Fprintln(output)

	return offset + bytes, nil
}

type unsignedIntReadFunc func([]byte) uint64

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

type signedIntReadFunc func([]byte) int64

func readInt64(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}

func readInt32(b []byte) int64 {
	return int64(int32(binary.BigEndian.Uint32(b)))
}

func readInt16(b []byte) int64 {
	return int64(int16(binary.BigEndian.Uint16(b)))
}

func readInt8(b []byte) int64 {
	return int64(int8(b[0]))
}

func (f *BytecodeFunction) disassembleValue(output io.Writer, byteLength, offset int) (int, error) {
	return f._disassembleValue(output, byteLength, -1, offset)
}

func (f *BytecodeFunction) _disassembleValue(output io.Writer, byteLength, constantIndex, offset int) (int, error) {
	opcode := bytecode.OpCode(f.Instructions[offset])

	if result, err := f.checkBytes(output, offset, byteLength); err != nil {
		return result, err
	}

	switch byteLength {
	case 1:
	case 2:
		constantIndex = int(f.Instructions[offset+1])
	case 3:
		constantIndex = int(binary.BigEndian.Uint16(f.Instructions[offset+1 : offset+3]))
	case 5:
		constantIndex = int(binary.BigEndian.Uint32(f.Instructions[offset+1 : offset+5]))
	default:
		panic(fmt.Sprintf("%d is not a valid byteLength for a value opcode!", byteLength))
	}

	f.printLineNumber(output, offset)
	f.dumpBytes(output, offset, byteLength)
	f.printOpCode(output, opcode)

	if constantIndex >= len(f.Values) {
		msg := fmt.Sprintf("invalid value index %d (0x%X)", constantIndex, constantIndex)
		fmt.Fprintln(output, msg)
		return offset + byteLength, errors.New(msg)
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
	return len(f.Instructions) - 1, errors.New(msg)
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

func printNumField[I int64 | uint64](output io.Writer, n I) {
	fmt.Fprintf(output, "%-16d", n)
}
