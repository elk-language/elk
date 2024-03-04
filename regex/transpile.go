package regex

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/regex/parser"
	"github.com/elk-language/elk/regex/parser/ast"
)

// Transpile an Elk regex string to a Go regex string
func Transpile(elkRegex string) (string, errors.ErrorList) {
	ast, err := parser.Parse(elkRegex)
	if err != nil {
		return "", err
	}

	t := &transpiler{}
	t.transpileNode(ast)
	if t.Errors != nil {
		return "", t.Errors
	}
	return t.Buffer.String(), nil
}

// Holds the state of the Transpiler.
type transpiler struct {
	Errors    errors.ErrorList
	Buffer    strings.Builder
	AsciiMode bool
}

// Create a new location struct with the given position.
func (t *transpiler) newLocation(span *position.Span) *position.Location {
	return position.NewLocationWithSpan("regex", span)
}

func asciiLetterIndex(char rune) int {
	if char >= 'A' && char <= 'Z' {
		return int(char - 'A' + 1)
	}

	if char >= 'a' && char <= 'z' {
		return int(char - 'a' + 1)
	}

	panic(fmt.Sprintf("char is not an ASCII letter: %c", char))
}

func (t *transpiler) transpileNode(node ast.Node) {
	switch n := node.(type) {
	case *ast.ConcatenationNode:
		t.concatenation(n)
	case *ast.ZeroOrOneQuantifierNode:
		t.zeroOrOneQuantifier(n)
	case *ast.ZeroOrMoreQuantifierNode:
		t.zeroOrMoreQuantifier(n)
	case *ast.OneOrMoreQuantifierNode:
		t.oneOrMoreQuantifier(n)
	case *ast.NQuantifierNode:
		t.nQuantifier(n)
	case *ast.NMQuantifierNode:
		t.nmQuantifier(n)
	case *ast.MetaCharEscapeNode:
		t.metaCharEscape(n)
	case *ast.GroupNode:
		t.group(n)
	case *ast.CharClassNode:
		t.charClass(n)
	case *ast.QuotedTextNode:
		t.quotedText(n)
	case *ast.CharNode:
		t.char(n)
	case *ast.CaretEscapeNode:
		t.caretEscape(n)
	case *ast.UnicodeEscapeNode:
		t.unicodeEscape(n)
	case *ast.HexEscapeNode:
		t.hexEscape(n)
	case *ast.OctalEscapeNode:
		t.octalEscape(n)
	case *ast.UnicodeCharClassNode:
		t.unicodeCharClass(n)
	case *ast.BellEscapeNode:
		t.bellEscape()
	case *ast.FormFeedEscapeNode:
		t.formFeedEscape()
	case *ast.TabEscapeNode:
		t.tabEscape()
	case *ast.NewlineEscapeNode:
		t.newlineEscape()
	case *ast.CarriageReturnEscapeNode:
		t.carriageReturnEscape()
	case *ast.StartOfStringAnchorNode:
		t.startOfStringAnchor()
	case *ast.EndOfStringAnchorNode:
		t.endOfStringAnchor()
	case *ast.AbsoluteStartOfStringAnchorNode:
		t.absoluteStartOfStringAnchor()
	case *ast.AbsoluteEndOfStringAnchorNode:
		t.absoluteEndOfStringAnchor()
	case *ast.WordBoundaryAnchorNode:
		t.wordBoundaryAnchor()
	case *ast.NotWordBoundaryAnchorNode:
		t.notWordBoundaryAnchor()
	case *ast.WordCharClassNode:
		t.wordCharClass(false)
	case *ast.NotWordCharClassNode:
		t.notWordCharClass(false)
	case *ast.DigitCharClassNode:
		t.digitCharClass()
	case *ast.NotDigitCharClassNode:
		t.notDigitCharClass()
	case *ast.WhitespaceCharClassNode:
		t.whitespaceCharClass(false)
	case *ast.NotWhitespaceCharClassNode:
		t.notWhitespaceCharClass()
	case *ast.HorizontalWhitespaceCharClassNode:
		t.horizontalWhitespaceCharClass()
	case *ast.NotHorizontalWhitespaceCharClassNode:
		t.notHorizontalWhitespaceCharClass()
	case *ast.VerticalWhitespaceCharClassNode:
		t.verticalWhitespaceCharClass()
	case *ast.NotVerticalWhitespaceCharClassNode:
		t.notVerticalWhitespaceCharClass()
	case *ast.AnyCharClassNode:
		t.anyCharClass()
	case nil:
	default:
		t.Errors.Add(
			fmt.Sprintf("compilation of this node has not been implemented: %T", node),
			t.newLocation(node.Span()),
		)
	}
}

func (t *transpiler) concatenation(node *ast.ConcatenationNode) {
	for _, element := range node.Elements {
		t.transpileNode(element)
	}
}

func (t *transpiler) zeroOrOneQuantifier(node *ast.ZeroOrOneQuantifierNode) {
	t.transpileNode(node.Regex)
	if node.Alt {
		t.Buffer.WriteString(`??`)
	} else {
		t.Buffer.WriteRune('?')
	}
}

func (t *transpiler) zeroOrMoreQuantifier(node *ast.ZeroOrMoreQuantifierNode) {
	t.transpileNode(node.Regex)
	if node.Alt {
		t.Buffer.WriteString(`*?`)
	} else {
		t.Buffer.WriteRune('*')
	}
}

func (t *transpiler) oneOrMoreQuantifier(node *ast.OneOrMoreQuantifierNode) {
	t.transpileNode(node.Regex)
	if node.Alt {
		t.Buffer.WriteString(`+?`)
	} else {
		t.Buffer.WriteRune('+')
	}
}

func (t *transpiler) nQuantifier(node *ast.NQuantifierNode) {
	t.transpileNode(node.Regex)
	t.Buffer.WriteRune('{')
	t.Buffer.WriteString(node.N)
	t.Buffer.WriteRune('}')
	if node.Alt {
		t.Buffer.WriteRune('?')
	}
}

func (t *transpiler) nmQuantifier(node *ast.NMQuantifierNode) {
	t.transpileNode(node.Regex)
	t.Buffer.WriteRune('{')
	if node.N == "" {
		t.Buffer.WriteRune('0')
	} else {
		t.Buffer.WriteString(node.N)
	}
	t.Buffer.WriteRune(',')
	if node.M != "" {
		t.Buffer.WriteString(node.M)
	}
	t.Buffer.WriteRune('}')
	if node.Alt {
		t.Buffer.WriteRune('?')
	}
}

func (t *transpiler) char(node *ast.CharNode) {
	t.Buffer.WriteRune(node.Value)
}

func (t *transpiler) group(node *ast.GroupNode) {
	t.Buffer.WriteRune('(')
	if len(node.Flags) > 0 {
		// with flags
		t.Buffer.WriteRune('?')
		t.Buffer.WriteString(node.Flags)
		if node.Regex != nil {
			// with flags and content
			t.Buffer.WriteRune(':')
			t.transpileNode(node.Regex)
		}
		t.Buffer.WriteRune(')')
		return
	}

	if len(node.Name) > 0 {
		// named
		t.Buffer.WriteString(`?P<`)
		t.Buffer.WriteString(node.Name)
		t.Buffer.WriteRune('>')
	} else if node.NonCapturing {
		// non capturing
		t.Buffer.WriteString(`?:`)
	}

	t.transpileNode(node.Regex)
	t.Buffer.WriteRune(')')
}

func (t *transpiler) charClass(node *ast.CharClassNode) {
	t.Buffer.WriteRune('[')
	if node.Negated {
		t.Buffer.WriteRune('^')
	}

	for _, element := range node.Elements {
		t.charClassElement(element)
	}

	t.Buffer.WriteRune(']')
}

func (t *transpiler) charClassElement(node ast.CharClassElementNode) {
	switch n := node.(type) {
	case *ast.CharRangeNode:
		t.charRange(n)
	case *ast.NamedCharClassNode:
		t.namedCharClass(n)
	case *ast.MetaCharEscapeNode:
		t.metaCharEscape(n)
	case *ast.CharNode:
		t.char(n)
	case *ast.CaretEscapeNode:
		t.caretEscape(n)
	case *ast.UnicodeEscapeNode:
		t.unicodeEscape(n)
	case *ast.HexEscapeNode:
		t.hexEscape(n)
	case *ast.OctalEscapeNode:
		t.octalEscape(n)
	case *ast.UnicodeCharClassNode:
		t.unicodeCharClass(n)
	case *ast.BellEscapeNode:
		t.bellEscape()
	case *ast.FormFeedEscapeNode:
		t.formFeedEscape()
	case *ast.TabEscapeNode:
		t.tabEscape()
	case *ast.NewlineEscapeNode:
		t.newlineEscape()
	case *ast.CarriageReturnEscapeNode:
		t.carriageReturnEscape()
	case *ast.WordCharClassNode:
		t.wordCharClass(true)
	case *ast.NotWordCharClassNode:
		t.notWordCharClass(true)
	case *ast.DigitCharClassNode:
		t.digitCharClass()
	case *ast.NotDigitCharClassNode:
		t.notDigitCharClass()
	case *ast.WhitespaceCharClassNode:
		t.whitespaceCharClass(true)
	case *ast.NotWhitespaceCharClassNode:
		t.notWhitespaceCharClass()
	case *ast.HorizontalWhitespaceCharClassNode:
		t.horizontalWhitespaceCharClass()
	case *ast.NotHorizontalWhitespaceCharClassNode:
		t.notHorizontalWhitespaceCharClass()
	case *ast.VerticalWhitespaceCharClassNode:
		t.verticalWhitespaceCharClass()
	case *ast.NotVerticalWhitespaceCharClassNode:
		t.notVerticalWhitespaceCharClass()
	}
}

func (t *transpiler) namedCharClass(node *ast.NamedCharClassNode) {
	t.Buffer.WriteString(`[:`)
	if node.Negated {
		t.Buffer.WriteRune('^')
	}

	t.Buffer.WriteString(node.Name)
	t.Buffer.WriteString(`:]`)
}

func (t *transpiler) charRange(node *ast.CharRangeNode) {
	t.charClassElement(node.Left)
	t.Buffer.WriteRune('-')
	t.charClassElement(node.Right)
}

func (t *transpiler) metaCharEscape(node *ast.MetaCharEscapeNode) {
	t.Buffer.WriteRune('\\')
	t.Buffer.WriteRune(node.Value)
}

func (t *transpiler) quotedText(node *ast.QuotedTextNode) {
	t.Buffer.WriteString(`\Q`)
	t.Buffer.WriteString(node.Value)
	t.Buffer.WriteString(`\E`)
}

func (t *transpiler) caretEscape(node *ast.CaretEscapeNode) {
	t.Buffer.WriteString(`\x{`)
	fmt.Fprintf(&t.Buffer, "%x", asciiLetterIndex(node.Value))
	t.Buffer.WriteString(`}`)
}

func (t *transpiler) unicodeEscape(node *ast.UnicodeEscapeNode) {
	t.Buffer.WriteString(`\x{`)
	t.Buffer.WriteString(node.Value)
	t.Buffer.WriteString(`}`)
}

func (t *transpiler) hexEscape(node *ast.HexEscapeNode) {
	t.Buffer.WriteString(`\x{`)
	t.Buffer.WriteString(node.Value)
	t.Buffer.WriteString(`}`)
}

func (t *transpiler) octalEscape(node *ast.OctalEscapeNode) {
	t.Buffer.WriteRune('\\')
	t.Buffer.WriteString(fmt.Sprintf("%03s", node.Value))
}

func (t *transpiler) unicodeCharClass(node *ast.UnicodeCharClassNode) {
	t.Buffer.WriteRune('\\')
	if node.Negated {
		t.Buffer.WriteRune('P')
	} else {
		t.Buffer.WriteRune('p')
	}

	t.Buffer.WriteRune('{')
	t.Buffer.WriteString(node.Value)
	t.Buffer.WriteRune('}')
}

func (t *transpiler) bellEscape() {
	t.Buffer.WriteString(`\a`)
}

func (t *transpiler) formFeedEscape() {
	t.Buffer.WriteString(`\f`)
}

func (t *transpiler) tabEscape() {
	t.Buffer.WriteString(`\t`)
}

func (t *transpiler) newlineEscape() {
	t.Buffer.WriteString(`\n`)
}

func (t *transpiler) carriageReturnEscape() {
	t.Buffer.WriteString(`\r`)
}

func (t *transpiler) startOfStringAnchor() {
	t.Buffer.WriteRune('^')
}

func (t *transpiler) endOfStringAnchor() {
	t.Buffer.WriteRune('$')
}

func (t *transpiler) absoluteStartOfStringAnchor() {
	t.Buffer.WriteString(`\A`)
}

func (t *transpiler) absoluteEndOfStringAnchor() {
	t.Buffer.WriteString(`\z`)
}

func (t *transpiler) wordBoundaryAnchor() {
	t.Buffer.WriteString(`\b`)
}

func (t *transpiler) notWordBoundaryAnchor() {
	t.Buffer.WriteString(`\B`)
}

func (t *transpiler) wordCharClass(inCharClass bool) {
	if t.AsciiMode {
		t.Buffer.WriteString(`\w`)
		return
	}

	if inCharClass {
		t.Buffer.WriteString(`\p{L}\p{N}_`)
		return
	}
	t.Buffer.WriteString(`[\p{L}\p{N}_]`)
}

func (t *transpiler) notWordCharClass(inCharClass bool) {
	if t.AsciiMode {
		t.Buffer.WriteString(`\W`)
		return
	}

	if inCharClass {
		t.Buffer.WriteString(`[^\p{L}\p{N}_]`)
	}
	t.Buffer.WriteString(`^\p{L}\p{N}_`)
}

func (t *transpiler) digitCharClass() {
	if t.AsciiMode {
		t.Buffer.WriteString(`\d`)
	} else {
		t.Buffer.WriteString(`\p{Nd}`)
	}
}

func (t *transpiler) notDigitCharClass() {
	if t.AsciiMode {
		t.Buffer.WriteString(`\D`)
	} else {
		t.Buffer.WriteString(`\P{Nd}`)
	}
}

func (t *transpiler) whitespaceCharClass(inCharClass bool) {
	if t.AsciiMode {
		if inCharClass {
			t.Buffer.WriteString(`\s\v`)
			return
		}

		t.Buffer.WriteString(`[\s\v]`)
		return
	}

	if inCharClass {
		t.Buffer.WriteString(`\s\v\p{Z}\x85`)
		return
	}
	t.Buffer.WriteString(`[\s\v\p{Z}\x85]`)
}

func (t *transpiler) notWhitespaceCharClass() {
	if t.AsciiMode {
		t.Buffer.WriteString(`[^\s\v]`)
	} else {
		t.Buffer.WriteString(`[^\s\v\p{Z}\x85]`)
	}
}

func (t *transpiler) horizontalWhitespaceCharClass() {
	if t.AsciiMode {
		t.Buffer.WriteString(`[\t ]`)
	} else {
		t.Buffer.WriteString(`[\t\p{Zs}]`)
	}
}

func (t *transpiler) notHorizontalWhitespaceCharClass() {
	if t.AsciiMode {
		t.Buffer.WriteString(`[^\t ]`)
	} else {
		t.Buffer.WriteString(`[^\t\p{Zs}]`)
	}
}

func (t *transpiler) verticalWhitespaceCharClass() {
	if t.AsciiMode {
		t.Buffer.WriteString(`[\n\v\f\r]`)
	} else {
		t.Buffer.WriteString(`[\n\v\f\r\x85\x{2028}\x{2029}]`)
	}
}

func (t *transpiler) notVerticalWhitespaceCharClass() {
	if t.AsciiMode {
		t.Buffer.WriteString(`[^\n\x0b\f\r]`)
	} else {
		t.Buffer.WriteString(`[^\n\x0b\f\r\x85\x{2028}\x{2029}]`)
	}
}

func (t *transpiler) anyCharClass() {
	t.Buffer.WriteRune('.')
}
