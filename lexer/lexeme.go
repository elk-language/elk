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
	LexSeparator                // Statement separator `\n`, `\r\n` or `;`
	LexLParen                   // Left parenthesis `(`
	LexRParen                   // Right parenthesis `)`
	LexLBrace                   // Left brace `{`
	LexRBrace                   // Right brace `}`
	LexLBracket                 // Left bracket `[`
	LexRBracket                 // Right bracket `]`
	LexComma                    // Comma `,`
	LexDot                      // Dot `.`
	LexMinus                    // Minus `-`
	LexMinusEqual               // Minus equal `-=`
	LexPlus                     // Plus `+`
	LexPlusEqual                // Plus equal `+=`
	LexStar                     // Star `*`
	LexStarEqual                // Star equal `*=`
	LexPower                    // Power `**`
	LexPowerEqual               // Power equal `**=`
	LexColon                    // Colon `:`
	LexColonEqual               // Colon equal `:=`
	LexLess                     // Less than `<`
	LexLessEqual                // Less than or equal `<=`
	LexGreater                  // Greater than `>`
	LexGreaterEqual             // Greater than or equal `>=`
	LexAssign                   // Assign `=`
	LexEqual                    // Equal `==`
	LexNotEqual                 // Not equal `!=`
	LexRefEqual                 // Reference equality operator `=:=`
	LexRefNotEqual              // Reference not equal operator `=!=`
	LexStrictEqual              // Strict equal `===`
	LexStrictNotEqual           // Strict not equal `!==`
	LexThickArrow               // Thick arrow `=>`
	LexThinArrow                // Thin arrow `->` (closure arrow)
	LexTilde                    // Tilde `~`
	LexTildeEqual               // Tilde equal `~=`
	LexMatchOperator            // Match operator `=~`
	LexWigglyArrow              // Wiggly arrow `~>` (lambda arrow)
	LexAnd                      // Bitwise and `&`
	LexAndEqual                 // Bitwise and equal `&=`
	LexAndAnd                   // Logical and `&&`
	LexAndAndEqual              // Logical and equal `&&=`
	LexOr                       // Bitwise or `|`
	LexOrEqual                  // Bitwise or equal `|=`
	LexOrOr                     // Logical or `||`
	LexOrOrEqual                // Logical or `||=`
	LexNilCoalesce              // Nil coalescing operator `??`
	LexNilCoalesceEqual         // Nil coalescing equal operator `??=`
	LexBang                     // Logical not `!`
	LexQuestionMark             // Question mark `?`
	LexSubtype                  // Subtype operator `<:`
	LexReverseSubtype           // Reverse subtype operator `:>`
	LexInstanceOf               // Instance of operator `<<:`
	LexReverseInstanceOf        // Reverse instance of operator `:>>`
	LexLBitShift                // Left bitwise shift `<<`
	LexLBitShiftEqual           // Left bitwise shift equal `<<=`
	LexRBitShift                // Right bitwise shift `>>`
	LexRBitShiftEqual           // Right bitwise shift equal `>>=`
	LexPercent                  // Percent `%`
	LexPercentEqual             // Percent equal `%=`
	LexPercentW                 // Word collection literal prefix `%w`
	LexPercentS                 // Symbol collection literal prefix `%s`
	LexPercentI                 // Integer collection literal prefix `%i`
	LexPercentF                 // Float collection literal prefix `%f`
	LexSetLiteralBeg            // Set literal beginning `%{`
	LexTupleLiteralBeg          // Tuple literal beginning `%(`
	LexPipeOperator             // Pipe operator `|>`
	LexScopeResOperator         // Scope resolution operator `::`
	LexDocComment               // Documentation comment `##[` ... `]##`
	LexRawString                // Raw String literal delimited by single quotes `'` ... `'`
	LexStringContent            // String literal content
	LexStringBeg                // Beginning delimiter of String literals `"`
	LexStringEnd                // Ending delimiter of String literals `"`
	LexInt                      // Int literal
	LexFloat                    // Float literal
	LexIdentifier               // Identifier
	LexPrivateIdentifier        // Identifier with a initial underscore
	LexConstant                 // Constant (identifier with an initial capital letter)
	LexPrivateConstant          // Constant with an initial underscore
	// Keywords start here
	LexKeyword   // any types greater than this value can be considered keywords
	LexNil       // Keyword `nil`
	LexFalse     // Keyword `false`
	LexTrue      // Keyword `true`
	LexIf        // Keyword `if`
	LexElse      // Keyword `else`
	LexElsif     // Keyword `elsif`
	LexUnless    // Keyword `unless`
	LexWhile     // Keyword `while`
	LexUntil     // Keyword `until`
	LexLoop      // Keyword `loop`
	LexBreak     // Keyword `break`
	LexReturn    // Keyword `return`
	LexDef       // Keyword `def`
	LexEnd       // Keyword `end`
	LexThen      // Keyword `then`
	LexClass     // Keyword `class`
	LexModule    // Keyword `module`
	LexMixin     // Keyword `mixin`
	LexInterface // Keyword `interface`
	LexType      // Keyword `type`
	LexVar       // Keyword `var`
	LexThrow     // Keyword `throw`
	LexCatch     // Keyword `catch`
	LexDo        // Keyword `do`
	LexEnsure    // Keyword `ensure`
	LexAlias     // Keyword `alias`
	LexSelf      // Keyword `self`
	LexSuper     // Keyword `super`
)

// Maps keywords to their Lexeme Type.
var keywords = map[string]LexemeType{
	"nil":       LexNil,
	"false":     LexFalse,
	"true":      LexTrue,
	"if":        LexIf,
	"else":      LexElse,
	"elsif":     LexElsif,
	"unless":    LexUnless,
	"while":     LexWhile,
	"until":     LexUntil,
	"loop":      LexLoop,
	"break":     LexBreak,
	"return":    LexReturn,
	"def":       LexDef,
	"end":       LexEnd,
	"then":      LexThen,
	"class":     LexClass,
	"module":    LexModule,
	"mixin":     LexMixin,
	"interface": LexInterface,
	"type":      LexType,
	"var":       LexVar,
	"throw":     LexThrow,
	"catch":     LexCatch,
	"do":        LexDo,
	"ensure":    LexEnsure,
	"alias":     LexAlias,
	"self":      LexSelf,
	"super":     LexSuper,
}
