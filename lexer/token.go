package lexer

// Represents the type of token
type TokenType int

type Token struct {
	Type       TokenType
	Value      string // Literal value of the token
	StartByte  int    // Index of the first byte of this token
	ByteLength int    // Number of bytes of the token
	Line       int    // Source line number where the token starts
	Column     int    // Source column number where the token starts
}

// Allocate a new End of File token.
func newEOF() *Token {
	return &Token{Type: EOFToken}
}

const (
	EOFToken               TokenType = iota // End Of File has been reached
	SeparatorToken                          // Statement separator `\n`, `\r\n` or `;`
	LParenToken                             // Left parenthesis `(`
	RParenToken                             // Right parenthesis `)`
	LBraceToken                             // Left brace `{`
	RBraceToken                             // Right brace `}`
	LBracketToken                           // Left bracket `[`
	RBracketToken                           // Right bracket `]`
	CommaToken                              // Comma `,`
	DotToken                                // Dot `.`
	MinusToken                              // Minus `-`
	MinusEqualToken                         // Minus equal `-=`
	PlusToken                               // Plus `+`
	PlusEqualToken                          // Plus equal `+=`
	StarToken                               // Star `*`
	StarEqualToken                          // Star equal `*=`
	PowerToken                              // Power `**`
	PowerEqualToken                         // Power equal `**=`
	ColonToken                              // Colon `:`
	ColonEqualToken                         // Colon equal `:=`
	LessToken                               // Less than `<`
	LessEqualToken                          // Less than or equal `<=`
	GreaterToken                            // Greater than `>`
	GreaterEqualToken                       // Greater than or equal `>=`
	AssignToken                             // Assign `=`
	EqualToken                              // Equal `==`
	NotEqualToken                           // Not equal `!=`
	RefEqualToken                           // Reference equality operator `=:=`
	RefNotEqualToken                        // Reference not equal operator `=!=`
	StrictEqualToken                        // Strict equal `===`
	StrictNotEqualToken                     // Strict not equal `!==`
	ThickArrowToken                         // Thick arrow `=>`
	ThinArrowToken                          // Thin arrow `->` (closure arrow)
	TildeToken                              // Tilde `~`
	TildeEqualToken                         // Tilde equal `~=`
	MatchOperatorToken                      // Match operator `=~`
	WigglyArrowToken                        // Wiggly arrow `~>` (lambda arrow)
	AndToken                                // Bitwise and `&`
	AndEqualToken                           // Bitwise and equal `&=`
	AndAndToken                             // Logical and `&&`
	AndAndEqualToken                        // Logical and equal `&&=`
	OrToken                                 // Bitwise or `|`
	TheAnswerToken                          // The answer to the great question of life, the universe, and everything.
	OrEqualToken                            // Bitwise or equal `|=`
	OrOrToken                               // Logical or `||`
	OrOrEqualToken                          // Logical or `||=`
	NilCoalesceToken                        // Nil coalescing operator `??`
	NilCoalesceEqualToken                   // Nil coalescing equal operator `??=`
	BangToken                               // Logical not `!`
	QuestionMarkToken                       // Question mark `?`
	SubtypeToken                            // Subtype operator `<:`
	ReverseSubtypeToken                     // Reverse subtype operator `:>`
	InstanceOfToken                         // Instance of operator `<<:`
	ReverseInstanceOfToken                  // Reverse instance of operator `:>>`
	LBitShiftToken                          // Left bitwise shift `<<`
	LBitShiftEqualToken                     // Left bitwise shift equal `<<=`
	RBitShiftToken                          // Right bitwise shift `>>`
	RBitShiftEqualToken                     // Right bitwise shift equal `>>=`
	PercentToken                            // Percent `%`
	PercentEqualToken                       // Percent equal `%=`
	PercentWToken                           // Word collection literal prefix `%w`
	PercentSToken                           // Symbol collection literal prefix `%s`
	PercentIToken                           // Integer collection literal prefix `%i`
	PercentFToken                           // Float collection literal prefix `%f`
	SetLiteralBegToken                      // Set literal beginning `%{`
	TupleLiteralBegToken                    // Tuple literal beginning `%(`
	PipeOperatorToken                       // Pipe operator `|>`
	ScopeResOperatorToken                   // Scope resolution operator `::`
	DocCommentToken                         // Documentation comment `##[` ... `]##`
	RawStringToken                          // Raw String literal delimited by single quotes `'` ... `'`
	StringBegToken                          // Beginning delimiter of String literals `"`
	StringContentToken                      // String literal content
	StringInterpBegToken                    // Beginning of string interpolation `${`
	StringInterpEndToken                    // End of string interpolation `}`
	StringEndToken                          // Ending delimiter of String literals `"`
	IntToken                                // Int literal
	FloatToken                              // Float literal
	IdentifierToken                         // Identifier
	PrivateIdentifierToken                  // Identifier with a initial underscore
	ConstantToken                           // Constant (identifier with an initial capital letter)
	PrivateConstantToken                    // Constant with an initial underscore
	// Keywords start here
	KeywordBegToken // any types greater than this value can be considered keywords
	NilToken        // Keyword `nil`
	FalseToken      // Keyword `false`
	TrueToken       // Keyword `true`
	IfToken         // Keyword `if`
	ElseToken       // Keyword `else`
	ElsifToken      // Keyword `elsif`
	UnlessToken     // Keyword `unless`
	WhileToken      // Keyword `while`
	UntilToken      // Keyword `until`
	LoopToken       // Keyword `loop`
	BreakToken      // Keyword `break`
	ReturnToken     // Keyword `return`
	DefToken        // Keyword `def`
	EndToken        // Keyword `end`
	ThenToken       // Keyword `then`
	ClassToken      // Keyword `class`
	ModuleToken     // Keyword `module`
	MixinToken      // Keyword `mixin`
	InterfaceToken  // Keyword `interface`
	TypeToken       // Keyword `type`
	VarToken        // Keyword `var`
	ThrowToken      // Keyword `throw`
	CatchToken      // Keyword `catch`
	DoToken         // Keyword `do`
	EnsureToken     // Keyword `ensure`
	AliasToken      // Keyword `alias`
	SelfToken       // Keyword `self`
	SuperToken      // Keyword `super`
	SwitchToken     // Keyword `switch`
	CaseToken       // Keyword `case`
	KeywordEndToken // any types lesser than this value can be considered keywords
)

// Maps keywords to their Token Type.
var keywords = map[string]TokenType{
	"nil":       NilToken,
	"false":     FalseToken,
	"true":      TrueToken,
	"if":        IfToken,
	"else":      ElseToken,
	"elsif":     ElsifToken,
	"unless":    UnlessToken,
	"while":     WhileToken,
	"until":     UntilToken,
	"loop":      LoopToken,
	"break":     BreakToken,
	"return":    ReturnToken,
	"def":       DefToken,
	"end":       EndToken,
	"then":      ThenToken,
	"class":     ClassToken,
	"module":    ModuleToken,
	"mixin":     MixinToken,
	"interface": InterfaceToken,
	"type":      TypeToken,
	"var":       VarToken,
	"throw":     ThrowToken,
	"catch":     CatchToken,
	"do":        DoToken,
	"ensure":    EnsureToken,
	"alias":     AliasToken,
	"self":      SelfToken,
	"super":     SuperToken,
	"switch":    SwitchToken,
	"case":      CaseToken,
}
