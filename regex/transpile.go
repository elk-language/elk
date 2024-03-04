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

func (t *transpiler) transpileNode(node ast.Node) {
	switch n := node.(type) {
	case *ast.ConcatenationNode:
		t.concatenation(n)
	case *ast.ZeroOrOneQuantifierNode:
		t.transpileNode(n.Regex)
		if n.Alt {
			t.Buffer.WriteString(`??`)
		} else {
			t.Buffer.WriteRune('?')
		}
	case *ast.ZeroOrMoreQuantifierNode:
		t.transpileNode(n.Regex)
		if n.Alt {
			t.Buffer.WriteString(`*?`)
		} else {
			t.Buffer.WriteRune('*')
		}
	case *ast.OneOrMoreQuantifierNode:
		t.transpileNode(n.Regex)
		if n.Alt {
			t.Buffer.WriteString(`+?`)
		} else {
			t.Buffer.WriteRune('+')
		}
	case *ast.NQuantifierNode:
		t.transpileNode(n.Regex)
		t.Buffer.WriteRune('{')
		t.Buffer.WriteString(n.N)
		t.Buffer.WriteRune('}')
		if n.Alt {
			t.Buffer.WriteRune('?')
		}
	case *ast.NMQuantifierNode:
		t.transpileNode(n.Regex)
		t.Buffer.WriteRune('{')
		if n.N == "" {
			t.Buffer.WriteRune('0')
		} else {
			t.Buffer.WriteString(n.N)
		}
		t.Buffer.WriteRune(',')
		if n.M != "" {
			t.Buffer.WriteString(n.M)
		}
		t.Buffer.WriteRune('}')
		if n.Alt {
			t.Buffer.WriteRune('?')
		}
	case *ast.CharNode:
		t.Buffer.WriteRune(n.Value)
	case *ast.GroupNode:
		t.Buffer.WriteRune('(')
		t.transpileNode(n.Regex)
		t.Buffer.WriteRune(')')
	case *ast.MetaCharEscapeNode:
		t.Buffer.WriteRune('\\')
		t.Buffer.WriteRune(n.Value)
	case *ast.QuotedTextNode:
		t.Buffer.WriteString(`\Q`)
		t.Buffer.WriteString(n.Value)
		t.Buffer.WriteString(`\E`)
	case *ast.HexEscapeNode:
		t.Buffer.WriteString(`\x{`)
		t.Buffer.WriteString(n.Value)
		t.Buffer.WriteString(`}`)
	case *ast.OctalEscapeNode:
		t.Buffer.WriteRune('\\')
		t.Buffer.WriteString(fmt.Sprintf("%03s", n.Value))
	case *ast.UnicodeCharClassNode:
		t.Buffer.WriteRune('\\')
		if n.Negated {
			t.Buffer.WriteRune('P')
		} else {
			t.Buffer.WriteRune('p')
		}

		t.Buffer.WriteRune('{')
		t.Buffer.WriteString(n.Value)
		t.Buffer.WriteRune('}')
	case *ast.BellEscapeNode:
		t.Buffer.WriteString(`\a`)
	case *ast.FormFeedEscapeNode:
		t.Buffer.WriteString(`\f`)
	case *ast.TabEscapeNode:
		t.Buffer.WriteString(`\t`)
	case *ast.NewlineEscapeNode:
		t.Buffer.WriteString(`\n`)
	case *ast.CarriageReturnEscapeNode:
		t.Buffer.WriteString(`\r`)
	case *ast.StartOfStringAnchorNode:
		t.Buffer.WriteRune('^')
	case *ast.EndOfStringAnchorNode:
		t.Buffer.WriteRune('$')
	case *ast.AbsoluteStartOfStringAnchorNode:
		t.Buffer.WriteString(`\A`)
	case *ast.AbsoluteEndOfStringAnchorNode:
		t.Buffer.WriteString(`\z`)
	case *ast.WordBoundaryAnchorNode:
		t.Buffer.WriteString(`\b`)
	case *ast.NotWordBoundaryAnchorNode:
		t.Buffer.WriteString(`\B`)
	case *ast.WordCharClassNode:
		if t.AsciiMode {
			t.Buffer.WriteString(`\w`)
		} else {
			t.Buffer.WriteString(`[\p{L}\p{N}_]`)
		}
	case *ast.NotWordCharClassNode:
		if t.AsciiMode {
			t.Buffer.WriteString(`\W`)
		} else {
			t.Buffer.WriteString(`[^\p{L}\p{N}_]`)
		}
	case *ast.DigitCharClassNode:
		if t.AsciiMode {
			t.Buffer.WriteString(`\d`)
		} else {
			t.Buffer.WriteString(`\p{Nd}`)
		}
	case *ast.NotDigitCharClassNode:
		if t.AsciiMode {
			t.Buffer.WriteString(`\D`)
		} else {
			t.Buffer.WriteString(`\P{Nd}`)
		}
	case *ast.WhitespaceCharClassNode:
		if t.AsciiMode {
			t.Buffer.WriteString(`\s`)
		} else {
			t.Buffer.WriteString(`[\p{Z}\t\v]`)
		}
	case *ast.NotWhitespaceCharClassNode:
		if t.AsciiMode {
			t.Buffer.WriteString(`\S`)
		} else {
			t.Buffer.WriteString(`[^\p{Z}\t\v]`)
		}
	case *ast.HorizontalWhitespaceCharClassNode:
		if t.AsciiMode {
			t.Buffer.WriteString(`[\t ]`)
		} else {
			t.Buffer.WriteString(`[\t\p{Zs}]`)
		}
	case *ast.NotHorizontalWhitespaceCharClassNode:
		if t.AsciiMode {
			t.Buffer.WriteString(`[^\t ]`)
		} else {
			t.Buffer.WriteString(`[^\t\p{Zs}]`)
		}
	case *ast.VerticalWhitespaceCharClassNode:
		if t.AsciiMode {
			t.Buffer.WriteString(`[\n\x0b\f\r]`)
		} else {
			t.Buffer.WriteString(`[\n\x0b\f\r\x85\x{2028}\x{2029}]`)
		}
	case *ast.NotVerticalWhitespaceCharClassNode:
		if t.AsciiMode {
			t.Buffer.WriteString(`[^\n\x0b\f\r]`)
		} else {
			t.Buffer.WriteString(`[^\n\x0b\f\r\x85\x{2028}\x{2029}]`)
		}
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
