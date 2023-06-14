package bytecode

// The maximum number of bytes a single
// instruction can take up.
const maxInstructionByteLength = 4

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
	RETURN   OpCode = iota // Return from the current frame
	CONSTANT               // Push a constant onto the value stack
)

var opCodeNames = [...]string{
	RETURN:   "RETURN",
	CONSTANT: "CONSTANT",
}
