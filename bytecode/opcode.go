package bytecode

// The maximum number of bytes a single
// instruction can take up.
const MaxInstructionByteLength = 5

const (
	UINT8_SIZE  = iota // The integer fits in a uint8
	UINT16_SIZE        // The integer fits in a uint16
	UINT32_SIZE        // The integer fits in a uint32
	UINT64_SIZE        // The integer fits in a uint64
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
	RETURN           OpCode = iota // Return from the current frame
	CONSTANT8                      // Push a constant with a single byte index onto the value stack
	CONSTANT16                     // Push a constant with a two byte index onto the value stack
	CONSTANT32                     // Push a constant with a four byte index onto the value stack
	ADD                            // Take two values from the stack, add them together (or call the + method) and push the result
	SUBTRACT                       // Take two values from the stack, subtract them (or call the - method) and push the result
	MULTIPLY                       // Take two values from the stack, multiply them (or call the * method) and push the result
	DIVIDE                         // Take two values from the stack, divide them (or call the / method) and push the result
	EXPONENTIATE                   // Take two values from the stack, raise one to the power signified by the other
	NEGATE                         // Take a value off the stack and negate it
	NOT                            // Take a value off the stack and perform boolean negation (converting it to a Bool)
	BITWISE_NOT                    // Take a value off the stack and perform bitwise negation
	TRUE                           // Push true onto the stack
	FALSE                          // Push false onto the stack
	NIL                            // Push nil onto the stack
	POP                            // Pop an element off the stack.
	POP_N                          // Pop n elements off the stack.
	LEAVE_SCOPE16                  // Leave a scope and pop off any local variables (16 bit operand)
	LEAVE_SCOPE32                  // Leave a scope and pop off any local variables (32 bit operand)
	PREP_LOCALS8                   // Prepare slots for local variables and values (8 bit operand)
	PREP_LOCALS16                  // Prepare slots for local variables and values (16 bit operand)
	SET_LOCAL8                     // Assign the value on top of the stack to the local variable with the given index (8 bit operand)
	SET_LOCAL16                    // Assign the value on top of the stack to the local variable with the given index (16 bit operand)
	GET_LOCAL8                     // Push the value of the local variable with the given index onto the stack (8 bit operand)
	GET_LOCAL16                    // Push the value of the local variable with the given index onto the stack (16 bit operand)
	JUMP_UNLESS                    // Jump n bytes forward if the value on the stack is falsy
	JUMP                           // Jump n bytes forward
	JUMP_IF                        // Jump n bytes forward if the value on the stack is truthy
	LOOP                           // Jump n bytes backward
	JUMP_IF_NIL                    // Jump n bytes forward if the value on the stack is nil
	RBITSHIFT                      // Take two values from the stack, perform a right bitshift and push the result
	LOGIC_RBITSHIFT                // Take two values from the stack, perform a logical right bitshift and push the result
	LBITSHIFT                      // Take two values from the stack, perform a left bitshift and push the result
	LOGIC_LBITSHIFT                // Take two values from the stack, perform a logical left bitshift and push the result
	BITWISE_AND                    // Take two values from the stack, perform a bitwise AND and push the result
	BITWISE_OR                     // Take two values from the stack, perform a bitwise OR and push the result
	BITWISE_XOR                    // Take two values from the stack, perform a bitwise XOR and push the result
	MODULO                         // Take two values from the stack, perform modulo and push the result
	EQUAL                          // Take two values from the stack, check if they're equal and push the result
	STRICT_EQUAL                   // Take two values from the stack, check if they're strictly equal and push the result
	GREATER                        // Take two values from the stack, check if the first value is greater than the second and push the result
	GREATER_EQUAL                  // Take two values from the stack, check if the first value is greater than or equal to the second and push the result
	LESS                           // Take two values from the stack, check if the first value is less than the second and push the result
	LESS_EQUAL                     // Take two values from the stack, check if the first value is less than or equal to the second and push the result
	GET_MOD_CONST8                 // Pop one value off the stack (module) and get the value of the constant with the name stored under the given index in the constant pool (8 bit operand)
	GET_MOD_CONST16                // Pop one value off the stack (module) and get the value of the constant with the name stored under the given index in the constant pool (16 bit operand)
	GET_MOD_CONST32                // Pop one value off the stack (module) and get the value of the constant with the name stored under the given index in the constant pool (32 bit operand)
	ROOT                           // Push `Std::Root` onto the stack.
	NOT_EQUAL                      // Take two values from the stack, check if they're not equal and push the result
	STRICT_NOT_EQUAL               // Take two values from the stack, check if they're strictly not equal and push the result
	DEF_MOD_CONST8                 // Pop one value off the stack (module) and define a new constant under it (8 bit operand)
	DEF_MOD_CONST16                // Pop one value off the stack (module) and define a new constant under it (16 bit operand)
	DEF_MOD_CONST32                // Pop one value off the stack (module) and define a new constant under it (32 bit operand)
	CONSTANT_BASE                  // Push the module/class/mixin that will hold constants defined in this context
	DEF_CLASS                      // Define a new class
	SELF                           // Push `self` onto the stack
	DEF_MODULE                     // Define a new module
)

var opCodeNames = [...]string{
	RETURN:           "RETURN",
	CONSTANT8:        "CONSTANT8",
	CONSTANT16:       "CONSTANT16",
	CONSTANT32:       "CONSTANT32",
	POP:              "POP",
	POP_N:            "POP_N",
	ADD:              "ADD",
	SUBTRACT:         "SUBTRACT",
	MULTIPLY:         "MULTIPLY",
	DIVIDE:           "DIVIDE",
	EXPONENTIATE:     "EXPONENTIATE",
	NEGATE:           "NEGATE",
	NOT:              "NOT",
	BITWISE_NOT:      "BITWISE_NOT",
	TRUE:             "TRUE",
	FALSE:            "FALSE",
	NIL:              "NIL",
	LEAVE_SCOPE16:    "LEAVE_SCOPE16",
	LEAVE_SCOPE32:    "LEAVE_SCOPE32",
	PREP_LOCALS8:     "PREP_LOCALS8",
	PREP_LOCALS16:    "PREP_LOCALS16",
	SET_LOCAL8:       "SET_LOCAL8",
	SET_LOCAL16:      "SET_LOCAL16",
	GET_LOCAL8:       "GET_LOCAL8",
	GET_LOCAL16:      "GET_LOCAL16",
	JUMP_UNLESS:      "JUMP_UNLESS",
	JUMP:             "JUMP",
	JUMP_IF:          "JUMP_IF",
	LOOP:             "LOOP",
	JUMP_IF_NIL:      "JUMP_IF_NIL",
	RBITSHIFT:        "RBITSHIFT",
	LOGIC_RBITSHIFT:  "LOGIC_RBITSHIFT",
	LBITSHIFT:        "LBITSHIFT",
	LOGIC_LBITSHIFT:  "LOGIC_LBITSHIFT",
	BITWISE_AND:      "BITWISE_AND",
	BITWISE_OR:       "BITWISE_OR",
	BITWISE_XOR:      "BITWISE_XOR",
	MODULO:           "MODULO",
	EQUAL:            "EQUAL",
	STRICT_EQUAL:     "STRICT_EQUAL",
	GREATER:          "GREATER",
	GREATER_EQUAL:    "GREATER_EQUAL",
	LESS:             "LESS",
	LESS_EQUAL:       "LESS_EQUAL",
	GET_MOD_CONST8:   "GET_MOD_CONST8",
	GET_MOD_CONST16:  "GET_MOD_CONST16",
	GET_MOD_CONST32:  "GET_MOD_CONST32",
	ROOT:             "ROOT",
	NOT_EQUAL:        "NOT_EQUAL",
	STRICT_NOT_EQUAL: "STRICT_NOT_EQUAL",
	DEF_MOD_CONST8:   "DEF_MOD_CONST8",
	DEF_MOD_CONST16:  "DEF_MOD_CONST16",
	DEF_MOD_CONST32:  "DEF_MOD_CONST32",
	CONSTANT_BASE:    "CONSTANT_BASE",
	DEF_CLASS:        "DEF_CLASS",
	SELF:             "SELF",
	DEF_MODULE:       "DEF_MODULE",
}
