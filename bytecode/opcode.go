package bytecode

// The maximum number of bytes a single
// instruction can take up.
const maxInstructionByteLength = 5

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
	RETURN       OpCode = iota // Return from the current frame
	CONSTANT8                  // Push a constant with a single byte index onto the value stack
	CONSTANT16                 // Push a constant with a two byte index onto the value stack
	CONSTANT32                 // Push a constant with a four byte index onto the value stack
	ADD                        // Take two values from the stack, add them together (or call the + method) and push the result
	SUBTRACT                   // Take two values from the stack, subtract them (or call the - method) and push the result
	MULTIPLY                   // Take two values from the stack, multiply them (or call the * method) and push the result
	DIVIDE                     // Take two values from the stack, divide them (or call the / method) and push the result
	EXPONENTIATE               // Take two values from the stack, raise one to the power signified by the other
	NEGATE                     // Take a value off the stack and negate it
	NOT                        // Take a value off the stack and perform boolean negation (converting it to a Bool)
	BITWISE_NOT                // Take a value off the stack and perform bitwise negation
	TRUE                       // Push true onto the stack
	FALSE                      // Push false onto the stack
	NIL                        // Push nil onto the stack
)

var opCodeNames = [...]string{
	RETURN:       "RETURN",
	CONSTANT8:    "CONSTANT8",
	CONSTANT16:   "CONSTANT16",
	CONSTANT32:   "CONSTANT32",
	ADD:          "ADD",
	SUBTRACT:     "SUBTRACT",
	MULTIPLY:     "MULTIPLY",
	DIVIDE:       "DIVIDE",
	EXPONENTIATE: "EXPONENTIATE",
	NEGATE:       "NEGATE",
	NOT:          "NOT",
	BITWISE_NOT:  "BITWISE_NOT",
	TRUE:         "TRUE",
	FALSE:        "FALSE",
	NIL:          "NIL",
}
