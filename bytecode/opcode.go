package bytecode

// The maximum number of bytes a single
// instruction can take up.
const MaxInstructionByteCount = 5

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
	ZERO_VALUE         OpCode = iota // Zero value
	RETURN                           // Return from the current frame
	LOAD_VALUE8                      // Push a value with a single byte index onto the value stack
	LOAD_VALUE16                     // Push a value with a two byte index onto the value stack
	LOAD_VALUE32                     // Push a value with a four byte index onto the value stack
	ADD                              // Take two values from the stack, add them together (or call the + method) and push the result
	SUBTRACT                         // Take two values from the stack, subtract them (or call the - method) and push the result
	MULTIPLY                         // Take two values from the stack, multiply them (or call the * method) and push the result
	DIVIDE                           // Take two values from the stack, divide them (or call the / method) and push the result
	EXPONENTIATE                     // Take two values from the stack, raise one to the power signified by the other
	NEGATE                           // Take a value off the stack and negate it
	NOT                              // Take a value off the stack and perform boolean negation (converting it to a Bool)
	BITWISE_NOT                      // Take a value off the stack and perform bitwise negation
	TRUE                             // Push true onto the stack
	FALSE                            // Push false onto the stack
	NIL                              // Push nil onto the stack
	POP                              // Pop an element off the stack.
	POP_N                            // Pop n elements off the stack.
	LEAVE_SCOPE16                    // Leave a scope and pop off any local variables (16 bit operand)
	LEAVE_SCOPE32                    // Leave a scope and pop off any local variables (32 bit operand)
	PREP_LOCALS8                     // Prepare slots for local variables and values (8 bit operand)
	PREP_LOCALS16                    // Prepare slots for local variables and values (16 bit operand)
	SET_LOCAL8                       // Assign the value on top of the stack to the local variable with the given index (8 bit operand)
	SET_LOCAL16                      // Assign the value on top of the stack to the local variable with the given index (16 bit operand)
	GET_LOCAL8                       // Push the value of the local variable with the given index onto the stack (8 bit operand)
	GET_LOCAL16                      // Push the value of the local variable with the given index onto the stack (16 bit operand)
	JUMP_UNLESS                      // Jump n bytes forward if the value on the stack is falsy
	JUMP                             // Jump n bytes forward
	JUMP_IF                          // Jump n bytes forward if the value on the stack is truthy
	LOOP                             // Jump n bytes backward
	JUMP_IF_NIL                      // Jump n bytes forward if the value on the stack is nil
	RBITSHIFT                        // Take two values from the stack, perform a right bitshift and push the result
	LOGIC_RBITSHIFT                  // Take two values from the stack, perform a logical right bitshift and push the result
	LBITSHIFT                        // Take two values from the stack, perform a left bitshift and push the result
	LOGIC_LBITSHIFT                  // Take two values from the stack, perform a logical left bitshift and push the result
	BITWISE_AND                      // Take two values from the stack, perform a bitwise AND and push the result
	BITWISE_OR                       // Take two values from the stack, perform a bitwise OR and push the result
	BITWISE_XOR                      // Take two values from the stack, perform a bitwise XOR and push the result
	MODULO                           // Take two values from the stack, perform modulo and push the result
	EQUAL                            // Take two values from the stack, check if they're equal and push the result
	STRICT_EQUAL                     // Take two values from the stack, check if they're strictly equal and push the result
	GREATER                          // Take two values from the stack, check if the first value is greater than the second and push the result
	GREATER_EQUAL                    // Take two values from the stack, check if the first value is greater than or equal to the second and push the result
	LESS                             // Take two values from the stack, check if the first value is less than the second and push the result
	LESS_EQUAL                       // Take two values from the stack, check if the first value is less than or equal to the second and push the result
	GET_MOD_CONST8                   // Pop one value off the stack (module) and get the value of the constant with the name stored under the given index in the constant pool (8 bit operand)
	GET_MOD_CONST16                  // Pop one value off the stack (module) and get the value of the constant with the name stored under the given index in the constant pool (16 bit operand)
	GET_MOD_CONST32                  // Pop one value off the stack (module) and get the value of the constant with the name stored under the given index in the constant pool (32 bit operand)
	ROOT                             // Push `Std::Root` onto the stack.
	NOT_EQUAL                        // Take two values from the stack, check if they're not equal and push the result
	STRICT_NOT_EQUAL                 // Take two values from the stack, check if they're strictly not equal and push the result
	DEF_MOD_CONST8                   // Pop one value off the stack (module) and define a new constant under it (8 bit operand)
	DEF_MOD_CONST16                  // Pop one value off the stack (module) and define a new constant under it (16 bit operand)
	DEF_MOD_CONST32                  // Pop one value off the stack (module) and define a new constant under it (32 bit operand)
	CONSTANT_CONTAINER               // Push the module/class/mixin that will hold constants defined in this context
	DEF_CLASS                        // Define a new class
	SELF                             // Push `self` onto the stack
	DEF_MODULE                       // Define a new module
	CALL_METHOD8                     // Call a method with an explicit receiver eg. `foo.bar(2)` (8 bit operand)
	CALL_METHOD16                    // Call a method with an explicit receiver eg. `foo.bar(2)` (16 bit operand)
	CALL_METHOD32                    // Call a method with an explicit receiver eg. `foo.bar(2)` (32 bit operand)
	DEF_METHOD                       // Define a new method
	UNDEFINED                        // Push the undefined value onto the stack
	DEF_ANON_CLASS                   // Define a new anonymous class
	DEF_ANON_MODULE                  // Define a new anonymous module
	CALL_FUNCTION8                   // Call a method with an implicit receiver eg. `bar(2)` (8 bit operand)
	CALL_FUNCTION16                  // Call a method with an implicit receiver eg. `bar(2)` (16 bit operand)
	CALL_FUNCTION32                  // Call a method with an implicit receiver eg. `bar(2)` (32 bit operand)
	DEF_MIXIN                        // Define a new mixin
	DEF_ANON_MIXIN                   // Define a new anonymous mixin
	INCLUDE                          // Include a mixin in a class/mixin
	GET_SINGLETON                    // Pop one value off the stack push its singleton class
	JUMP_UNLESS_UNDEF                // Jump n bytes forward unless the value on the stack is undefined
	DEF_ALIAS                        // Define a method alias
	METHOD_CONTAINER                 // Push the class/mixin that will hold methods defined in this context
	COMPARE                          // Pop two values, compare them using `<=>` and push the result
	DOC_COMMENT                      // Attach a doc comment to an Elk object
	DEF_GETTER                       // Define a getter method
	DEF_SETTER                       // Define a setter method
	DEF_SINGLETON                    // Open the definition of a singleton class of the given object
	RETURN_FIRST_ARG                 // Push the first given argument (constant container for modules, classes etc) and return
	INSTANTIATE8                     // Create a new instance of a class (8 bit operand)
	INSTANTIATE16                    // Create a new instance of a class (16 bit operand)
	INSTANTIATE32                    // Create a new instance of a class (32 bit operand)
	RETURN_SELF                      // Push self and return
	GET_IVAR8                        // Get the value of an instance variable (8 bit operand)
	GET_IVAR16                       // Get the value of an instance variable (16 bit operand)
	GET_IVAR32                       // Get the value of an instance variable (32 bit operand)
	SET_IVAR8                        // Set the value of an instance variable (8 bit operand)
	SET_IVAR16                       // Set the value of an instance variable (16 bit operand)
	SET_IVAR32                       // Set the value of an instance variable (32 bit operand)
	NEW_TUPLE8                       // Create a new tuple (8 bit operand)
	NEW_TUPLE32                      // Create a new tuple (32 bit operand)
	APPEND_TUPLE                     // Append an element to a tuple, pops the element and leaves the tuple on the stack
)

var opCodeNames = [...]string{
	ZERO_VALUE:         "ZERO_VALUE",
	RETURN:             "RETURN",
	LOAD_VALUE8:        "LOAD_VALUE8",
	LOAD_VALUE16:       "LOAD_VALUE16",
	LOAD_VALUE32:       "LOAD_VALUE32",
	POP:                "POP",
	POP_N:              "POP_N",
	ADD:                "ADD",
	SUBTRACT:           "SUBTRACT",
	MULTIPLY:           "MULTIPLY",
	DIVIDE:             "DIVIDE",
	EXPONENTIATE:       "EXPONENTIATE",
	NEGATE:             "NEGATE",
	NOT:                "NOT",
	BITWISE_NOT:        "BITWISE_NOT",
	TRUE:               "TRUE",
	FALSE:              "FALSE",
	NIL:                "NIL",
	LEAVE_SCOPE16:      "LEAVE_SCOPE16",
	LEAVE_SCOPE32:      "LEAVE_SCOPE32",
	PREP_LOCALS8:       "PREP_LOCALS8",
	PREP_LOCALS16:      "PREP_LOCALS16",
	SET_LOCAL8:         "SET_LOCAL8",
	SET_LOCAL16:        "SET_LOCAL16",
	GET_LOCAL8:         "GET_LOCAL8",
	GET_LOCAL16:        "GET_LOCAL16",
	JUMP_UNLESS:        "JUMP_UNLESS",
	JUMP:               "JUMP",
	JUMP_IF:            "JUMP_IF",
	LOOP:               "LOOP",
	JUMP_IF_NIL:        "JUMP_IF_NIL",
	RBITSHIFT:          "RBITSHIFT",
	LOGIC_RBITSHIFT:    "LOGIC_RBITSHIFT",
	LBITSHIFT:          "LBITSHIFT",
	LOGIC_LBITSHIFT:    "LOGIC_LBITSHIFT",
	BITWISE_AND:        "BITWISE_AND",
	BITWISE_OR:         "BITWISE_OR",
	BITWISE_XOR:        "BITWISE_XOR",
	MODULO:             "MODULO",
	EQUAL:              "EQUAL",
	STRICT_EQUAL:       "STRICT_EQUAL",
	GREATER:            "GREATER",
	GREATER_EQUAL:      "GREATER_EQUAL",
	LESS:               "LESS",
	LESS_EQUAL:         "LESS_EQUAL",
	GET_MOD_CONST8:     "GET_MOD_CONST8",
	GET_MOD_CONST16:    "GET_MOD_CONST16",
	GET_MOD_CONST32:    "GET_MOD_CONST32",
	ROOT:               "ROOT",
	NOT_EQUAL:          "NOT_EQUAL",
	STRICT_NOT_EQUAL:   "STRICT_NOT_EQUAL",
	DEF_MOD_CONST8:     "DEF_MOD_CONST8",
	DEF_MOD_CONST16:    "DEF_MOD_CONST16",
	DEF_MOD_CONST32:    "DEF_MOD_CONST32",
	CONSTANT_CONTAINER: "CONSTANT_CONTAINER",
	DEF_CLASS:          "DEF_CLASS",
	SELF:               "SELF",
	DEF_MODULE:         "DEF_MODULE",
	CALL_METHOD8:       "CALL_METHOD8",
	CALL_METHOD16:      "CALL_METHOD16",
	CALL_METHOD32:      "CALL_METHOD32",
	DEF_METHOD:         "DEF_METHOD",
	UNDEFINED:          "UNDEFINED",
	DEF_ANON_CLASS:     "DEF_ANON_CLASS",
	DEF_ANON_MODULE:    "DEF_ANON_MODULE",
	CALL_FUNCTION8:     "CALL_FUNCTION8",
	CALL_FUNCTION16:    "CALL_FUNCTION16",
	CALL_FUNCTION32:    "CALL_FUNCTION32",
	DEF_MIXIN:          "DEF_MIXIN",
	DEF_ANON_MIXIN:     "DEF_ANON_MIXIN",
	INCLUDE:            "INCLUDE",
	GET_SINGLETON:      "GET_SINGLETON",
	JUMP_UNLESS_UNDEF:  "JUMP_UNLESS_UNDEF",
	DEF_ALIAS:          "DEF_ALIAS",
	METHOD_CONTAINER:   "METHOD_CONTAINER",
	COMPARE:            "COMPARE",
	DOC_COMMENT:        "DOC_COMMENT",
	DEF_GETTER:         "DEF_GETTER",
	DEF_SETTER:         "DEF_SETTER",
	DEF_SINGLETON:      "DEF_SINGLETON",
	RETURN_FIRST_ARG:   "RETURN_FIRST_ARG",
	INSTANTIATE8:       "INSTANTIATE8",
	INSTANTIATE16:      "INSTANTIATE16",
	INSTANTIATE32:      "INSTANTIATE32",
	RETURN_SELF:        "RETURN_SELF",
	GET_IVAR8:          "GET_IVAR8",
	GET_IVAR16:         "GET_IVAR16",
	GET_IVAR32:         "GET_IVAR32",
	SET_IVAR8:          "SET_IVAR8",
	SET_IVAR16:         "SET_IVAR16",
	SET_IVAR32:         "SET_IVAR32",
	NEW_TUPLE8:         "NEW_TUPLE8",
	NEW_TUPLE32:        "NEW_TUPLE32",
	APPEND_TUPLE:       "APPEND_TUPLE",
}
