package regex

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/regex/flag"
	"github.com/elk-language/elk/regex/parser"
	"github.com/elk-language/elk/regex/parser/ast"
)

// Transpile an Elk regex string to a Go regex string
func Transpile(elkRegex string, flags bitfield.BitField8) (string, error.ErrorList) {
	ast, err := parser.Parse(elkRegex)
	if err != nil {
		return "", err
	}

	t := &transpiler{Flags: flags}
	t.transpileNode(ast)
	if t.Errors != nil {
		return "", t.Errors
	}
	return t.Buffer.String(), nil
}

// Transpiler mode
type mode uint8

const (
	topLevelMode mode = iota
	charClassMode
	negatedCharClassMode
)

// Holds the state of the Transpiler.
type transpiler struct {
	Errors error.ErrorList
	Buffer strings.Builder
	Mode   mode
	Flags  bitfield.BitField8
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

func (t *transpiler) nodeHasToBeSplitInCharacterClasses(node ast.Node) bool {
	switch node.(type) {
	case *ast.NotHorizontalWhitespaceCharClassNode, *ast.NotVerticalWhitespaceCharClassNode:
		switch t.Mode {
		case topLevelMode, negatedCharClassMode:
			return false
		case charClassMode:
			return true
		}
	case *ast.NotWhitespaceCharClassNode, *ast.NotWordCharClassNode:
		if t.Flags.HasFlag(flag.ASCIIFlag) {
			return false
		}
		switch t.Mode {
		case topLevelMode, negatedCharClassMode:
			return false
		case charClassMode:
			return true
		}
	}

	return false
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
	case *ast.UnionNode:
		t.union(n)
	case *ast.CharClassNode:
		t.charClass(n)
	case *ast.QuotedTextNode:
		t.quotedText(n)
	case *ast.CharNode:
		if t.Flags.HasFlag(flag.ExtendedFlag) && unicode.IsSpace(n.Value) {
			return
		}
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
		t.wordCharClass()
	case *ast.NotWordCharClassNode:
		t.notWordCharClass(n)
	case *ast.DigitCharClassNode:
		t.digitCharClass()
	case *ast.NotDigitCharClassNode:
		t.notDigitCharClass()
	case *ast.WhitespaceCharClassNode:
		t.whitespaceCharClass()
	case *ast.NotWhitespaceCharClassNode:
		t.notWhitespaceCharClass(n)
	case *ast.HorizontalWhitespaceCharClassNode:
		t.horizontalWhitespaceCharClass()
	case *ast.NotHorizontalWhitespaceCharClassNode:
		t.notHorizontalWhitespaceCharClass(n)
	case *ast.VerticalWhitespaceCharClassNode:
		t.verticalWhitespaceCharClass()
	case *ast.NotVerticalWhitespaceCharClassNode:
		t.notVerticalWhitespaceCharClass(n)
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
	var inComment bool

	for _, element := range node.Elements {
		if t.Flags.HasFlag(flag.ExtendedFlag) {
			if inComment {
				if ch, ok := element.(*ast.CharNode); ok && ch.Value == '\n' {
					inComment = false
				}
				continue
			}
			if ch, ok := element.(*ast.CharNode); ok && ch.Value == '#' {
				inComment = true
				continue
			}
		}
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
	originalFlags := t.Flags
	var visibleSetFlags, visibleUnsetFlags bitfield.BitField8

	for _, fl := range flag.Flags {
		if node.SetFlags.HasFlag(fl) {
			t.Flags.SetFlag(fl)
			if flag.IsSupportedByGo(fl) {
				visibleSetFlags.SetFlag(fl)
			}
		}

		if node.UnsetFlags.HasFlag(fl) {
			t.Flags.UnsetFlag(fl)
			if flag.IsSupportedByGo(fl) {
				visibleUnsetFlags.SetFlag(fl)
			}
		}
	}

	var hasVisibleContent bool
	if node.Regex != nil || (visibleSetFlags.IsAnyFlagSet() || visibleUnsetFlags.IsAnyFlagSet()) {
		t.Buffer.WriteRune('(')
		hasVisibleContent = true
	}

	if visibleSetFlags.IsAnyFlagSet() || visibleUnsetFlags.IsAnyFlagSet() {
		// with flags
		t.Buffer.WriteRune('?')
		for _, fl := range flag.Flags {
			if !visibleSetFlags.HasFlag(fl) {
				continue
			}

			char := flag.ToChar(fl)
			t.Buffer.WriteRune(char)
		}
		if visibleUnsetFlags.IsAnyFlagSet() {
			t.Buffer.WriteRune('-')
			for _, fl := range flag.Flags {
				if !visibleUnsetFlags.HasFlag(fl) {
					continue
				}

				char := flag.ToChar(fl)
				t.Buffer.WriteRune(char)
			}
		}
		if node.Regex != nil {
			// with flags and content
			t.Buffer.WriteRune(':')
			t.transpileNode(node.Regex)
			t.Flags = originalFlags
		}
		t.Buffer.WriteRune(')')
		return
	} else if hasVisibleContent && node.IsAnyFlagSet() {
		t.Buffer.WriteString(`?:`)
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

	if node.Regex != nil {
		t.Flags = originalFlags
	}
	if hasVisibleContent {
		t.Buffer.WriteRune(')')
	}
}

func (t *transpiler) union(node *ast.UnionNode) {
	t.transpileNode(node.Left)
	t.Buffer.WriteRune('|')
	t.transpileNode(node.Right)
}

func (t *transpiler) charClass(node *ast.CharClassNode) {
	var internalNodes []ast.CharClassElementNode
	var nodesToSplit []ast.CharClassElementNode

	if node.Negated {
		t.Mode = negatedCharClassMode
	} else {
		t.Mode = charClassMode
	}

	for _, element := range node.Elements {
		if t.nodeHasToBeSplitInCharacterClasses(element) {
			nodesToSplit = append(nodesToSplit, element)
			continue
		}

		internalNodes = append(internalNodes, element)
	}
	if len(nodesToSplit) > 0 {
		t.Buffer.WriteString(`(?:[`)
	} else {
		t.Buffer.WriteRune('[')
	}
	if node.Negated {
		t.Buffer.WriteRune('^')
	}

	for _, element := range internalNodes {
		t.charClassElement(element)
	}

	t.Buffer.WriteRune(']')
	t.Mode = topLevelMode
	if len(nodesToSplit) > 0 {
		for _, element := range nodesToSplit {
			t.Buffer.WriteRune('|')
			t.charClassElement(element)
		}
		t.Buffer.WriteRune(')')
	}
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
		t.wordCharClass()
	case *ast.NotWordCharClassNode:
		t.notWordCharClass(n)
	case *ast.DigitCharClassNode:
		t.digitCharClass()
	case *ast.NotDigitCharClassNode:
		t.notDigitCharClass()
	case *ast.WhitespaceCharClassNode:
		t.whitespaceCharClass()
	case *ast.NotWhitespaceCharClassNode:
		t.notWhitespaceCharClass(n)
	case *ast.HorizontalWhitespaceCharClassNode:
		t.horizontalWhitespaceCharClass()
	case *ast.NotHorizontalWhitespaceCharClassNode:
		t.notHorizontalWhitespaceCharClass(n)
	case *ast.VerticalWhitespaceCharClassNode:
		t.verticalWhitespaceCharClass()
	case *ast.NotVerticalWhitespaceCharClassNode:
		t.notVerticalWhitespaceCharClass(n)
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

func (t *transpiler) wordCharClass() {
	if t.Flags.HasFlag(flag.ASCIIFlag) {
		t.Buffer.WriteString(`\w`)
		return
	}

	switch t.Mode {
	case topLevelMode:
		t.Buffer.WriteString(`[\p{L}\p{Mn}\p{Nd}\p{Pc}]`)
	case charClassMode, negatedCharClassMode:
		t.Buffer.WriteString(`\p{L}\p{Mn}\p{Nd}\p{Pc}`)
	}
}

func (t *transpiler) notWordCharClass(node *ast.NotWordCharClassNode) {
	if t.Flags.HasFlag(flag.ASCIIFlag) {
		t.Buffer.WriteString(`\W`)
		return
	}

	switch t.Mode {
	case topLevelMode:
		t.Buffer.WriteString(`[^\p{L}\p{Mn}\p{Nd}\p{Pc}]`)
	case charClassMode:
		t.Errors.Add(
			`unicode-aware \W in char classes is not supported`,
			t.newLocation(node.Span()),
		)
	case negatedCharClassMode:
		t.Errors.Add(
			`double negation of unicode-aware \W is not supported`,
			t.newLocation(node.Span()),
		)
	}
}

func (t *transpiler) digitCharClass() {
	if t.Flags.HasFlag(flag.ASCIIFlag) {
		t.Buffer.WriteString(`\d`)
	} else {
		t.Buffer.WriteString(`\p{Nd}`)
	}
}

func (t *transpiler) notDigitCharClass() {
	if t.Flags.HasFlag(flag.ASCIIFlag) {
		t.Buffer.WriteString(`\D`)
	} else {
		t.Buffer.WriteString(`\P{Nd}`)
	}
}

func (t *transpiler) whitespaceCharClass() {
	if t.Flags.HasFlag(flag.ASCIIFlag) {
		t.Buffer.WriteString(`\s`)
		return
	}

	switch t.Mode {
	case topLevelMode:
		t.Buffer.WriteString(`[\s\v\p{Z}\x85]`)
	case charClassMode, negatedCharClassMode:
		t.Buffer.WriteString(`\s\v\p{Z}\x85`)
	}
}

func (t *transpiler) notWhitespaceCharClass(node *ast.NotWhitespaceCharClassNode) {
	if t.Flags.HasFlag(flag.ASCIIFlag) {
		t.Buffer.WriteString(`\S`)
		return
	}

	switch t.Mode {
	case topLevelMode:
		t.Buffer.WriteString(`[^\s\v\p{Z}\x85]`)
	case charClassMode:
		t.Errors.Add(
			`unicode-aware \S in char classes is not supported`,
			t.newLocation(node.Span()),
		)
	case negatedCharClassMode:
		t.Errors.Add(
			`double negation of unicode-aware \S is not supported`,
			t.newLocation(node.Span()),
		)
	}
}

func (t *transpiler) horizontalWhitespaceCharClass() {
	switch t.Mode {
	case topLevelMode:
		if t.Flags.HasFlag(flag.ASCIIFlag) {
			t.Buffer.WriteString(`[\t ]`)
		} else {
			t.Buffer.WriteString(`[\t\p{Zs}]`)
		}
	case charClassMode, negatedCharClassMode:
		if t.Flags.HasFlag(flag.ASCIIFlag) {
			t.Buffer.WriteString(`\t `)
		} else {
			t.Buffer.WriteString(`\t\p{Zs}`)
		}
	}
}

func (t *transpiler) notHorizontalWhitespaceCharClass(node *ast.NotHorizontalWhitespaceCharClassNode) {
	switch t.Mode {
	case topLevelMode:
		if t.Flags.HasFlag(flag.ASCIIFlag) {
			t.Buffer.WriteString(`[^\t ]`)
		} else {
			t.Buffer.WriteString(`[^\t\p{Zs}]`)
		}
	case charClassMode:
		t.Errors.Add(
			`unicode-aware \H in char classes is not supported`,
			t.newLocation(node.Span()),
		)
	case negatedCharClassMode:
		t.Errors.Add(
			`double negation of unicode-aware \H is not supported`,
			t.newLocation(node.Span()),
		)
	}
}

func (t *transpiler) verticalWhitespaceCharClass() {
	switch t.Mode {
	case topLevelMode:
		if t.Flags.HasFlag(flag.ASCIIFlag) {
			t.Buffer.WriteString(`[\n\v\f\r]`)
		} else {
			t.Buffer.WriteString(`[\n\v\f\r\x85\x{2028}\x{2029}]`)
		}
	case charClassMode, negatedCharClassMode:
		if t.Flags.HasFlag(flag.ASCIIFlag) {
			t.Buffer.WriteString(`\n\v\f\r`)
		} else {
			t.Buffer.WriteString(`\n\v\f\r\x85\x{2028}\x{2029}`)
		}
	}
}

func (t *transpiler) notVerticalWhitespaceCharClass(node *ast.NotVerticalWhitespaceCharClassNode) {
	switch t.Mode {
	case topLevelMode:
		if t.Flags.HasFlag(flag.ASCIIFlag) {
			t.Buffer.WriteString(`[^\n\v\f\r]`)
		} else {
			t.Buffer.WriteString(`[^\n\v\f\r\x85\x{2028}\x{2029}]`)
		}
	case charClassMode:
		t.Errors.Add(
			`unicode-aware \V in char classes is not supported`,
			t.newLocation(node.Span()),
		)
	case negatedCharClassMode:
		t.Errors.Add(
			`double negation of unicode-aware \V is not supported`,
			t.newLocation(node.Span()),
		)
	}
}

func (t *transpiler) anyCharClass() {
	t.Buffer.WriteRune('.')
}
