package bytecode

// The maximum number of bytes a single
// instruction can take up.
const maxInstructionByteLength = 5

// Represents Operation Codes
// used by the Elk Virtual Machine.
type OpCode uint8

func (o OpCode) String() string {
	if int(o) > len(opCodeNames) {
		return "UNKNOWN"
	}

	return opCodeNames[o]
}

const (
	RETURN     OpCode = iota // Return from the current frame
	CONSTANT8                // Push a constant with a single byte index onto the value stack
	CONSTANT16               // Push a constant with a two byte index onto the value stack
	CONSTANT32               // Push a constant with a four byte index onto the value stack
)

var opCodeNames = [...]string{
	RETURN:     "RETURN",
	CONSTANT8:  "CONSTANT8",
	CONSTANT16: "CONSTANT16",
	CONSTANT32: "CONSTANT32",
}
