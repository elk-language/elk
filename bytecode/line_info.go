package bytecode

// Contains a source code line number
// with the count of bytecode instructions that were generated from that line
type LineInfo struct {
	LineNumber       int // Number of the line of source code that the instructions were generated from
	InstructionCount int // Number of consecutive bytecode instructions that share a single line of source code
}

// Create a new LineInfo.
// LineNumber and InstructionCount should be greater than 0.
func NewLineInfo(lineNum, instructCount int) *LineInfo {
	return &LineInfo{
		LineNumber:       lineNum,
		InstructionCount: instructCount,
	}
}

type LineInfoList []*LineInfo

// Retrieve the last LineInfo.
func (l LineInfoList) Last() *LineInfo {
	if len(l) == 0 {
		return nil
	}
	return l[len(l)-1]
}

// Retrieve the first LineInfo.
func (l LineInfoList) First() *LineInfo {
	if len(l) == 0 {
		return nil
	}
	return l[0]
}

// Get the source code line number for the given
// bytecode instruction index.
// Returns -1 when the line number couldn't be found.
func (l LineInfoList) GetLineNumber(instructionIndex int) int {
	lineInfo := l.GetLineInfo(instructionIndex)
	if lineInfo == nil {
		return -1
	}
	return lineInfo.LineNumber
}

func (l LineInfoList) GetLineInfo(instructionIndex int) *LineInfo {
	currentBytecodeOffset := 0
	for _, lineInfo := range l {
		currentBytecodeOffset += lineInfo.InstructionCount
		if currentBytecodeOffset-1 >= instructionIndex {
			return lineInfo
		}
	}

	return nil
}

func (l *LineInfoList) AddBytesToLastLine(bytes int) {
	lastLineInfo := l.Last()
	lastLineInfo.InstructionCount += bytes
}

// Set the source code line number for the next bytecode instruction.
func (l *LineInfoList) AddLineNumber(lineNumber int, bytes int) {
	lastLineInfo := l.Last()
	if lastLineInfo != nil && lastLineInfo.LineNumber == lineNumber {
		lastLineInfo.InstructionCount += bytes
		return
	}

	*l = append(*l, NewLineInfo(lineNumber, bytes))
}

func (l *LineInfoList) RemoveByte() {
	lastLineInfo := l.Last()
	if lastLineInfo == nil {
		panic("cannot remove a byte from an empty line info list")
	}

	if lastLineInfo.InstructionCount == 1 {
		*l = (*l)[:len(*l)-1]
		return
	}

	lastLineInfo.InstructionCount--
}

func (l *LineInfoList) RemoveBytes(count int) {
	for range count {
		l.RemoveByte()
	}
}
