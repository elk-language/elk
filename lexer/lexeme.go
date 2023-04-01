package lexer

// Represents the type of lexeme
type LexemeType int

type Lexeme struct {
	Type       LexemeType
	Value      string // Literal value of the lexeme
	StartByte  int    // Index of the first byte of this lexeme
	ByteLength int    // Number of bytes of the lexeme
	Line       int    // Source line number where the lexeme starts
	Column     int    // Source column number where the lexeme starts
}

// Allocate a new End of File lexeme.
func newEOF() *Lexeme {
	return &Lexeme{Type: LexEOF}
}

const (
	LexEOF               = iota // End Of File has been reached
	LexSeparator                // Statement separator "\n", "\r\n" or ";"
	LexLParen                   // Left parenthesis "("
	LexRParen                   // Right parenthesis ")"
	LexLBrace                   // Left brace "{"
	LexRBrace                   // Right brace "}"
	LexLBracket                 // Left bracket "["
	LexRBracket                 // Right bracket "]"
	LexComma                    // Comma ","
	LexDot                      // Dot "."
	LexMinus                    // Minus "-"
	LexPlus                     // Plus "+"
	LexStar                     // Star "*"
	LexPower                    // Power "**"
	LexColon                    // Colon ":"
	LexLess                     // Less than "<"
	LexGreater                  // Greater than ">"
	LexAssign                   // Assign "="
	LexEqual                    // Equal "=="
	LexNotEqual                 // Not equal "!="
	LexRefEqual                 // Reference equality operator "=:="
	LexRefNotEqual              // Reference not equal operator "=!="
	LexStrictEqual              // Strict equal "==="
	LexStrictNotEqual           // Strict not equal "!=="
	LexLessEqual                // Less than or equal "<="
	LexGreaterEqual             // Greater than or equal ">="
	LexThickArrow               // Thick arrow "=>"
	LexThinArrow                // Thin arrow "->" (closure arrow)
	LexWigglyArrow              // Wiggly arrow "~>" (lambda arrow)
	LexWalrusAssign             // Walrus assign ":="
	LexAnd                      // Bitwise and "&"
	LexAndAnd                   // Logical and "&&"
	LexOr                       // Bitwise or "|"
	LexOrOr                     // Logical or "||"
	LexNilCoalesce              // Nil coalescing operator "??"
	LexBang                     // Logical not "!"
	LexSubtype                  // Subtype operator "<:"
	LexReverseSubtype           // Subtype operator ":>"
	LexInstanceOf               // Instance of operator ":>>"
	LexReverseInstanceOf        // Instance of operator "<<:"
	LexLBitShift                // Left bitwise shift "<<"
	LexRBitShift                // Right bitwise shift ">>"
	LexPercent                  // Percent "%"
	LexPercentW                 // Word collection literal prefix "%w"
	LexPercentS                 // Symbol collection literal prefix "%s"
	LexPercentI                 // Integer collection literal prefix "%i"
	LexPercentF                 // Float collection literal prefix "%f"
	LexSetLiteralBeg            // Set literal beginning "%{"
	LexTupleLiteralBeg          // Tuple literal beginning "%("
	LexPipeOperator             // Pipe operator "|>"
	LexScopeResOperator         // Scope resolution operator '::'
	LexKeyword                  // any types greater than this value can be considered keywords
)
