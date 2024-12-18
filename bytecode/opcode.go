package bytecode

// The maximum number of bytes a single
// instruction can take up.
const MaxInstructionByteCount = 6

const (
	UINT8_SIZE  = iota // The integer fits in a uint8
	UINT16_SIZE        // The integer fits in a uint16
	UINT32_SIZE        // The integer fits in a uint32
	UINT64_SIZE        // The integer fits in a uint64
)

const (
	CLOSED_RANGE_FLAG byte = iota
	OPEN_RANGE_FLAG
	LEFT_OPEN_RANGE_FLAG
	RIGHT_OPEN_RANGE_FLAG
	BEGINLESS_CLOSED_RANGE_FLAG
	BEGINLESS_OPEN_RANGE_FLAG
	ENDLESS_CLOSED_RANGE_FLAG
	ENDLESS_OPEN_RANGE_FLAG
)

const (
	DEF_MODULE_FLAG byte = iota
	DEF_CLASS_FLAG
	DEF_MIXIN_FLAG
	DEF_INTERFACE_FLAG
)

// Represents Operation Codes
// used by the Elk Virtual Machine.
type OpCode byte

func (o OpCode) String() string {
	if int(o) > len(opCodeNames) {
		return "UNKNOWN"
	}

	return opCodeNames[o]
}

const (
	ZERO_VALUE        OpCode = iota // Zero value
	RETURN                          // Return from the current frame
	LOAD_VALUE8                     // Push a value with a single byte index onto the value stack
	LOAD_VALUE16                    // Push a value with a two byte index onto the value stack
	ADD                             // Take two values from the stack, add them together (or call the + method) and push the result
	ADD_INT                         // Take two values from the stack, add them together (or call the + method) and push the result
	SUBTRACT                        // Take two values from the stack, subtract them (or call the - method) and push the result
	SUBTRACT_INT                    // Take two values from the stack, subtract them (or call the - method) and push the result
	MULTIPLY                        // Take two values from the stack, multiply them (or call the * method) and push the result
	DIVIDE                          // Take two values from the stack, divide them (or call the / method) and push the result
	EXPONENTIATE                    // Take two values from the stack, raise one to the power signified by the other
	NEGATE                          // Take a value off the stack and negate it
	NOT                             // Take a value off the stack and perform boolean negation (converting it to a Bool)
	BITWISE_NOT                     // Take a value off the stack and perform bitwise negation
	TRUE                            // Push true onto the stack
	FALSE                           // Push false onto the stack
	NIL                             // Push nil onto the stack
	POP                             // Pop an element off the stack.
	POP_N                           // Pop n elements off the stack.
	LEAVE_SCOPE16                   // Leave a scope and pop off any local variables (16 bit operand)
	LEAVE_SCOPE32                   // Leave a scope and pop off any local variables (32 bit operand)
	PREP_LOCALS8                    // Prepare slots for local variables and values (8 bit operand)
	PREP_LOCALS16                   // Prepare slots for local variables and values (16 bit operand)
	SET_LOCAL_1                     // Assign the value on top of the stack to the local variable with index 1
	SET_LOCAL_2                     // Assign the value on top of the stack to the local variable with index 2
	SET_LOCAL_3                     // Assign the value on top of the stack to the local variable with index 3
	SET_LOCAL_4                     // Assign the value on top of the stack to the local variable with index 4
	SET_LOCAL8                      // Assign the value on top of the stack to the local variable with the given index (8 bit operand)
	SET_LOCAL16                     // Assign the value on top of the stack to the local variable with the given index (16 bit operand)
	GET_LOCAL_1                     // Push the value of the local variable with index 1
	GET_LOCAL_2                     // Push the value of the local variable with index 2
	GET_LOCAL_3                     // Push the value of the local variable with index 3
	GET_LOCAL_4                     // Push the value of the local variable with index 4
	GET_LOCAL8                      // Push the value of the local variable with the given index onto the stack (8 bit operand)
	GET_LOCAL16                     // Push the value of the local variable with the given index onto the stack (16 bit operand)
	JUMP_UNLESS_ILE                 // Jump n bytes forward if the the Int on the stack is not less than or equal to the second value
	JUMP_UNLESS_ILT                 // Jump n bytes forward if the the Int on the stack is not less than the second value
	JUMP_UNLESS_IGE                 // Jump n bytes forward if the the Int on the stack is not greater than or equal to the second value
	JUMP_UNLESS_IGT                 // Jump n bytes forward if the the Int on the stack is not greater than the second value
	JUMP_UNLESS_IEQ                 // Jump n bytes forward if the the Int on the stack is not equal to the second value
	JUMP_UNLESS                     // Jump n bytes forward if the value on the stack is falsy
	JUMP_UNLESS_NIL                 // Jump n bytes forward if the value on the stack is not nil
	JUMP                            // Jump n bytes forward
	JUMP_IF                         // Jump n bytes forward if the value on the stack is truthy
	JUMP_IF_IEQ                     // Jump n bytes forward if the the Int on the stack is equal to the second value
	LOOP                            // Jump n bytes backward
	JUMP_IF_NIL                     // Jump n bytes forward if the value on the stack is nil
	RBITSHIFT                       // Take two values from the stack, perform a right bitshift and push the result
	LOGIC_RBITSHIFT                 // Take two values from the stack, perform a logical right bitshift and push the result
	LBITSHIFT                       // Take two values from the stack, perform a left bitshift and push the result
	LOGIC_LBITSHIFT                 // Take two values from the stack, perform a logical left bitshift and push the result
	BITWISE_AND                     // Take two values from the stack, perform a bitwise AND and push the result
	BITWISE_OR                      // Take two values from the stack, perform a bitwise OR and push the result
	BITWISE_XOR                     // Take two values from the stack, perform a bitwise XOR and push the result
	MODULO                          // Take two values from the stack, perform modulo and push the result
	EQUAL                           // Take two values from the stack, check if they're equal and push the result
	STRICT_EQUAL                    // Take two values from the stack, check if they're strictly equal and push the result
	GREATER                         // Take two values from the stack, check if the first value is greater than the second and push the result
	GREATER_EQUAL                   // Take two values from the stack, check if the first value is greater than or equal to the second and push the result
	LESS                            // Take two values from the stack, check if the first value is less than the second and push the result
	LESS_EQUAL                      // Take two values from the stack, check if the first value is less than or equal to the second and push the result
	LESS_EQUAL_INT                  // Take two values from the stack, check if the first value is less than or equal to the second and push the result
	NOOP                            // Does not perform any operation, placeholder.
	ROOT                            // Push `Std::Root` onto the stack.
	NOT_EQUAL                       // Take two values from the stack, check if they're not equal and push the result
	STRICT_NOT_EQUAL                // Take two values from the stack, check if they're strictly not equal and push the result
	INIT_NAMESPACE                  // Initialise a namespace
	SELF                            // Push `self` onto the stack
	CALL_METHOD8                    // Call a method with an explicit receiver eg. `foo.bar(2)` (8 bit operand)
	CALL_METHOD16                   // Call a method with an explicit receiver eg. `foo.bar(2)` (16 bit operand)
	DEF_METHOD                      // Define a new method
	UNDEFINED                       // Push the undefined value onto the stack
	GET_CLASS                       // Pop one value off the stack push its class
	CALL_SELF8                      // Call a method with an implicit receiver eg. `bar(2)` (8 bit operand)
	CALL_SELF16                     // Call a method with an implicit receiver eg. `bar(2)` (16 bit operand)
	INCLUDE                         // Include a mixin in a class/mixin
	GET_SINGLETON                   // Pop one value off the stack push its singleton class
	JUMP_UNLESS_UNDEF               // Jump n bytes forward unless the value on the stack is undefined
	COMPARE                         // Pop two values, compare them using `<=>` and push the result
	DOC_COMMENT                     // Attach a doc comment to an Elk object
	DEF_GETTER                      // Define a getter method
	DEF_SETTER                      // Define a setter method
	RETURN_FIRST_ARG                // Push the first given argument (constant container for modules, classes etc) and return
	INSTANTIATE8                    // Create a new instance of a class (8 bit operand)
	INSTANTIATE16                   // Create a new instance of a class (16 bit operand)
	RETURN_SELF                     // Push self and return
	GET_IVAR8                       // Get the value of an instance variable (8 bit operand)
	GET_IVAR16                      // Get the value of an instance variable (16 bit operand)
	SET_IVAR8                       // Set the value of an instance variable (8 bit operand)
	SET_IVAR16                      // Set the value of an instance variable (16 bit operand)
	NEW_ARRAY_TUPLE8                // Create a new arrayTuple (8 bit operand)
	NEW_ARRAY_TUPLE16               // Create a new arrayTuple (16 bit operand)
	APPEND                          // Append an element to a list or arrayTuple, pops the element and leaves the collection on the stack
	COPY                            // Create a copy of the value on top of the stack and replace it on the stack.
	SUBSCRIPT                       // Pops 2 values off the stack. Get the element in a ArrayList, ArrayTuple or HashMap under the given key.
	SUBSCRIPT_SET                   // Pops 3 values off the stack. Set the element in a ArrayList, ArrayTuple or HashMap under the given key.
	APPEND_AT                       // Set an element at the given index in the ArrayTuple or ArrayList, if the index is out of range, expand the collection, filling the empty slots with `nil`
	NEW_ARRAY_LIST8                 // Create a new list (8 bit operand)
	NEW_ARRAY_LIST16                // Create a new list (16 bit operand)
	GET_ITERATOR                    // Get the iterator of the value on top of the stack.
	FOR_IN                          // Drives the for..in loop
	NEW_STRING8                     // Create a new string (8 bit operand)
	NEW_STRING16                    // Create a new string (16 bit operand)
	NEW_HASH_MAP8                   // Create a new hashmap (8 bit operand)
	NEW_HASH_MAP16                  // Create a new hashmap (16 bit operand)
	MAP_SET                         // Set a value under the given key in a hash record or hashmap, pops the key and value and leaves the collection on the stack
	NEW_HASH_RECORD8                // Create a new hash record (8 bit operand)
	NEW_HASH_RECORD16               // Create a new hash record (16 bit operand)
	LAX_EQUAL                       // Take two values from the stack, check if they are equal and push the result
	LAX_NOT_EQUAL                   // Take two values from the stack, check if they are not equal and push the result
	NEW_REGEX8                      // Create a new regex (8 bit operand)
	NEW_REGEX16                     // Create a new regex (16 bit operand)
	BITWISE_AND_NOT                 // Take two values from the stack, perform a bitwise AND NOT and push the result
	UNARY_PLUS                      // Perform unary plus on the value on top of the stack like `+a`
	INCREMENT                       // Increment the value on top of the stack
	DECREMENT                       // Decrement the value on top of the stack
	DUP                             // Duplicate the value on top of the stack
	DUP_N                           // Duplicate the top N values on top of the stack
	POP_N_SKIP_ONE                  // Pop the top N values on top of the stack skipping the first one
	NEW_SYMBOL8                     // Create a new symbol (8 bit operand)
	NEW_SYMBOL16                    // Create a new symbol (16 bit operand)
	SWAP                            // Swap the top two values on the stack
	NEW_RANGE                       // Create a new range
	SET_SUPERCLASS                  // Sets the superclass/parent of a class
	AS                              // Throw an error if the second value on the stack is not an instance of the class/mixin on top of the stack
	MUST                            // Throw an error if the value on top of the stack is `nil`
	INSTANCE_OF                     // Pop two values of the stack, check whether one is an instance of the other
	IS_A                            // Pop two values of the stack, check whether one is an instance of the subclass of the other
	POP_SKIP_ONE                    // Pop the value on top of the stack skipping the first one
	INSPECT_STACK                   // Prints the stack, for debugging
	NEW_HASH_SET8                   // Create a new hashset (8 bit operand)
	NEW_HASH_SET16                  // Create a new hashset (16 bit operand)
	THROW                           // Throw a value/error
	RETHROW                         // Rethrow a value/error
	POP_ALL                         // Pop all values on the stack, leaving only the slots reserved for locals
	RETURN_FINALLY                  // Execute all finally blocks this line is nested in and return from the current frame
	JUMP_TO_FINALLY                 // Jump to the specified instruction after executing finally blocks
	CLOSURE                         // Wrap the function on top of the stack in a closure
	CALL8                           // Call the `call` method with an explicit receiver eg. `foo.call(2)` (8 bit operand)
	CALL16                          // Call the `call` method with an explicit receiver eg. `foo.call(2)` (16 bit operand)
	SET_UPVALUE8                    // Assign the value on top of the stack to the upvalue with the given index (8 bit operand)
	SET_UPVALUE16                   // Assign the value on top of the stack to the upvalue with the given index (16 bit operand)
	GET_UPVALUE8                    // Push the value of the upvalue with the given index onto the stack (8 bit operand)
	GET_UPVALUE16                   // Push the value of the upvalue with the given index onto the stack (16 bit operand)
	CLOSE_UPVALUE8                  // Close an upvalue with the given index, moving it from the stack to the heap (8 bit operand)
	CLOSE_UPVALUE16                 // Close an upvalue with the given index, moving it from the stack to the heap (16 bit operand)
	DEF_NAMESPACE                   // Define a new namespace
	DEF_METHOD_ALIAS                // Define a new method alias
	GET_CONST8                      // Get the value of the constant with the name stored under the given index in the value pool (8 bit operand)
	GET_CONST16                     // Get the value of the constant with the name stored under the given index in the value pool (16 bit operand)
	DEF_CONST                       // Define a new constant
	EXEC                            // Execute a chunk of bytecode
	INT_M1                          // Push -1 onto the stack
	INT_0                           // Push 0 onto the stack
	INT_1                           // Push 1 onto the stack
	INT_2                           // Push 2 onto the stack
	INT_3                           // Push 3 onto the stack
	INT_4                           // Push 4 onto the stack
	INT_5                           // Push 5 onto the stack
	LOAD_INT_8                      // Push an 8 bit Int onto the stack
	LOAD_INT_16                     // Push a 16 bit Int onto the stack
	INT64_0                         // Push 0i64 onto the stack
	INT64_1                         // Push 1i64 onto the stack
	LOAD_INT64_8                    // Push an 8 bit Int64 onto the stack
	UINT64_0                        // Push 0u64 onto the stack
	UINT64_1                        // Push 1u64 onto the stack
	LOAD_UINT64_8                   // Push an 8 bit UInt64 onto the stack
	INT32_0                         // Push 0i32 onto the stack
	INT32_1                         // Push 1i32 onto the stack
	LOAD_INT32_8                    // Push an 8 bit Int32 onto the stack
	UINT32_0                        // Push 0u32 onto the stack
	UINT32_1                        // Push 1u32 onto the stack
	LOAD_UINT32_8                   // Push an 8 bit UInt32 onto the stack
	INT16_0                         // Push 0i16 onto the stack
	INT16_1                         // Push 1i16 onto the stack
	LOAD_INT16_8                    // Push an 8 bit Int16 onto the stack
	UINT16_0                        // Push 0u16 onto the stack
	UINT16_1                        // Push 1u16 onto the stack
	LOAD_UINT16_8                   // Push an 8 bit UInt16 onto the stack
	INT8_0                          // Push 0i8 onto the stack
	INT8_1                          // Push 1i8 onto the stack
	LOAD_INT8                       // Push an Int8 onto the stack
	UINT8_0                         // Push 0u8 onto the stack
	UINT8_1                         // Push 1u8 onto the stack
	LOAD_UINT8                      // Push a UInt8 onto the stack
	LOAD_CHAR_8                     // Push an 8 bit Char onto the stack
	FLOAT_0                         // Push 0.0 onto the stack
	FLOAT_1                         // Push 1.0 onto the stack
	FLOAT_2                         // Push 1.0 onto the stack
)

var opCodeNames = [...]string{
	ZERO_VALUE:        "ZERO_VALUE",
	RETURN:            "RETURN",
	LOAD_VALUE8:       "LOAD_VALUE8",
	LOAD_VALUE16:      "LOAD_VALUE16",
	POP:               "POP",
	POP_N:             "POP_N",
	ADD:               "ADD",
	ADD_INT:           "ADD_INT",
	SUBTRACT:          "SUBTRACT",
	SUBTRACT_INT:      "SUBTRACT_INT",
	MULTIPLY:          "MULTIPLY",
	DIVIDE:            "DIVIDE",
	EXPONENTIATE:      "EXPONENTIATE",
	NEGATE:            "NEGATE",
	NOT:               "NOT",
	BITWISE_NOT:       "BITWISE_NOT",
	TRUE:              "TRUE",
	FALSE:             "FALSE",
	NIL:               "NIL",
	LEAVE_SCOPE16:     "LEAVE_SCOPE16",
	LEAVE_SCOPE32:     "LEAVE_SCOPE32",
	PREP_LOCALS8:      "PREP_LOCALS8",
	PREP_LOCALS16:     "PREP_LOCALS16",
	SET_LOCAL_1:       "SET_LOCAL_1",
	SET_LOCAL_2:       "SET_LOCAL_2",
	SET_LOCAL_3:       "SET_LOCAL_3",
	SET_LOCAL_4:       "SET_LOCAL_4",
	SET_LOCAL8:        "SET_LOCAL8",
	SET_LOCAL16:       "SET_LOCAL16",
	GET_LOCAL_1:       "GET_LOCAL_1",
	GET_LOCAL_2:       "GET_LOCAL_2",
	GET_LOCAL_3:       "GET_LOCAL_3",
	GET_LOCAL_4:       "GET_LOCAL_4",
	GET_LOCAL8:        "GET_LOCAL8",
	GET_LOCAL16:       "GET_LOCAL16",
	JUMP_UNLESS_ILE:   "JUMP_UNLESS_ILE",
	JUMP_UNLESS_ILT:   "JUMP_UNLESS_ILT",
	JUMP_UNLESS_IGE:   "JUMP_UNLESS_IGE",
	JUMP_UNLESS_IGT:   "JUMP_UNLESS_IGT",
	JUMP_UNLESS_IEQ:   "JUMP_UNLESS_IEQ",
	JUMP_UNLESS:       "JUMP_UNLESS",
	JUMP_UNLESS_NIL:   "JUMP_UNLESS_NIL",
	JUMP:              "JUMP",
	JUMP_IF:           "JUMP_IF",
	JUMP_IF_IEQ:       "JUMP_IF_IEQ",
	LOOP:              "LOOP",
	JUMP_IF_NIL:       "JUMP_IF_NIL",
	RBITSHIFT:         "RBITSHIFT",
	LOGIC_RBITSHIFT:   "LOGIC_RBITSHIFT",
	LBITSHIFT:         "LBITSHIFT",
	LOGIC_LBITSHIFT:   "LOGIC_LBITSHIFT",
	BITWISE_AND:       "BITWISE_AND",
	BITWISE_OR:        "BITWISE_OR",
	BITWISE_XOR:       "BITWISE_XOR",
	MODULO:            "MODULO",
	EQUAL:             "EQUAL",
	STRICT_EQUAL:      "STRICT_EQUAL",
	GREATER:           "GREATER",
	GREATER_EQUAL:     "GREATER_EQUAL",
	LESS:              "LESS",
	LESS_EQUAL:        "LESS_EQUAL",
	LESS_EQUAL_INT:    "LESS_EQUAL_INT",
	NOOP:              "NOOP",
	ROOT:              "ROOT",
	NOT_EQUAL:         "NOT_EQUAL",
	STRICT_NOT_EQUAL:  "STRICT_NOT_EQUAL",
	INIT_NAMESPACE:    "INIT_NAMESPACE",
	SELF:              "SELF",
	CALL_METHOD8:      "CALL_METHOD8",
	CALL_METHOD16:     "CALL_METHOD16",
	DEF_METHOD:        "DEF_METHOD",
	UNDEFINED:         "UNDEFINED",
	GET_CLASS:         "GET_CLASS",
	CALL_SELF8:        "CALL_SELF8",
	CALL_SELF16:       "CALL_SELF16",
	INCLUDE:           "INCLUDE",
	GET_SINGLETON:     "GET_SINGLETON",
	JUMP_UNLESS_UNDEF: "JUMP_UNLESS_UNDEF",
	COMPARE:           "COMPARE",
	DOC_COMMENT:       "DOC_COMMENT",
	DEF_GETTER:        "DEF_GETTER",
	DEF_SETTER:        "DEF_SETTER",
	RETURN_FIRST_ARG:  "RETURN_FIRST_ARG",
	INSTANTIATE8:      "INSTANTIATE8",
	INSTANTIATE16:     "INSTANTIATE16",
	RETURN_SELF:       "RETURN_SELF",
	GET_IVAR8:         "GET_IVAR8",
	GET_IVAR16:        "GET_IVAR16",
	SET_IVAR8:         "SET_IVAR8",
	SET_IVAR16:        "SET_IVAR16",
	NEW_ARRAY_TUPLE8:  "NEW_ARRAY_TUPLE8",
	NEW_ARRAY_TUPLE16: "NEW_ARRAY_TUPLE16",
	APPEND:            "APPEND",
	COPY:              "COPY",
	SUBSCRIPT:         "SUBSCRIPT",
	SUBSCRIPT_SET:     "SUBSCRIPT_SET",
	APPEND_AT:         "APPEND_AT",
	NEW_ARRAY_LIST8:   "NEW_ARRAY_LIST8",
	NEW_ARRAY_LIST16:  "NEW_ARRAY_LIST16",
	GET_ITERATOR:      "GET_ITERATOR",
	FOR_IN:            "FOR_IN",
	NEW_STRING8:       "NEW_STRING8",
	NEW_STRING16:      "NEW_STRING16",
	NEW_HASH_MAP8:     "NEW_HASH_MAP8",
	NEW_HASH_MAP16:    "NEW_HASH_MAP16",
	MAP_SET:           "MAP_SET",
	NEW_HASH_RECORD8:  "NEW_HASH_RECORD8",
	NEW_HASH_RECORD16: "NEW_HASH_RECORD16",
	LAX_EQUAL:         "LAX_EQUAL",
	LAX_NOT_EQUAL:     "LAX_NOT_EQUAL",
	NEW_REGEX8:        "NEW_REGEX8",
	NEW_REGEX16:       "NEW_REGEX16",
	BITWISE_AND_NOT:   "BITWISE_AND_NOT",
	UNARY_PLUS:        "UNARY_PLUS",
	INCREMENT:         "INCREMENT",
	DECREMENT:         "DECREMENT",
	DUP:               "DUP",
	DUP_N:             "DUP_N",
	POP_N_SKIP_ONE:    "POP_N_SKIP_ONE",
	NEW_SYMBOL8:       "NEW_SYMBOL8",
	NEW_SYMBOL16:      "NEW_SYMBOL16",
	SWAP:              "SWAP",
	NEW_RANGE:         "NEW_RANGE",
	SET_SUPERCLASS:    "SET_SUPERCLASS",
	AS:                "AS",
	MUST:              "MUST",
	INSTANCE_OF:       "INSTANCE_OF",
	IS_A:              "IS_A",
	POP_SKIP_ONE:      "POP_SKIP_ONE",
	INSPECT_STACK:     "INSPECT_STACK",
	NEW_HASH_SET8:     "NEW_HASH_SET8",
	NEW_HASH_SET16:    "NEW_HASH_SET16",
	THROW:             "THROW",
	RETHROW:           "RETHROW",
	POP_ALL:           "POP_ALL",
	RETURN_FINALLY:    "RETURN_FINALLY",
	JUMP_TO_FINALLY:   "JUMP_TO_FINALLY",
	CLOSURE:           "CLOSURE",
	CALL8:             "CALL8",
	CALL16:            "CALL16",
	SET_UPVALUE8:      "SET_UPVALUE8",
	SET_UPVALUE16:     "SET_UPVALUE16",
	GET_UPVALUE8:      "GET_UPVALUE8",
	GET_UPVALUE16:     "GET_UPVALUE16",
	CLOSE_UPVALUE8:    "CLOSE_UPVALUE8",
	CLOSE_UPVALUE16:   "CLOSE_UPVALUE16",
	DEF_NAMESPACE:     "DEF_NAMESPACE",
	DEF_METHOD_ALIAS:  "DEF_METHOD_ALIAS",
	GET_CONST8:        "GET_CONST8",
	GET_CONST16:       "GET_CONST16",
	DEF_CONST:         "DEF_CONST",
	EXEC:              "EXEC",
	INT_M1:            "INT_M1",
	INT_0:             "INT_0",
	INT_1:             "INT_1",
	INT_2:             "INT_2",
	INT_3:             "INT_3",
	INT_4:             "INT_4",
	INT_5:             "INT_5",
	LOAD_INT_8:        "LOAD_INT_8",
	LOAD_INT_16:       "LOAD_INT_16",
	INT64_0:           "INT64_0",
	INT64_1:           "INT64_1",
	LOAD_INT64_8:      "LOAD_INT64_8",
	UINT64_0:          "UINT64_0",
	UINT64_1:          "UINT64_1",
	LOAD_UINT64_8:     "LOAD_UINT64_8",
	INT32_0:           "INT32_0",
	INT32_1:           "INT32_1",
	LOAD_INT32_8:      "LOAD_INT32_8",
	UINT32_0:          "UINT32_0",
	UINT32_1:          "UINT32_1",
	LOAD_UINT32_8:     "LOAD_UINT32_8",
	INT16_0:           "INT16_0",
	INT16_1:           "INT16_1",
	LOAD_INT16_8:      "LOAD_INT16_8",
	UINT16_0:          "UINT16_0",
	UINT16_1:          "UINT16_1",
	LOAD_UINT16_8:     "LOAD_UINT16_8",
	INT8_0:            "INT8_0",
	INT8_1:            "INT8_1",
	LOAD_INT8:         "LOAD_INT8",
	UINT8_0:           "UINT8_0",
	UINT8_1:           "UINT8_1",
	LOAD_UINT8:        "LOAD_UINT8",
	LOAD_CHAR_8:       "LOAD_CHAR_8",
	FLOAT_0:           "FLOAT_0",
	FLOAT_1:           "FLOAT_1",
	FLOAT_2:           "FLOAT_2",
}
