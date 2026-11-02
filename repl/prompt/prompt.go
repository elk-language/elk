package prompt

import (
	"strings"
	"unicode/utf8"

	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/token"
	goprompt "github.com/elk-language/go-prompt"
	pstrings "github.com/elk-language/go-prompt/strings"
)

// Adapter for `lexer.Lexer` that
// provides an interface compatible with go-prompt.
type Lexer struct {
	lexer.Lexer
}

func (l *Lexer) Init(input string) {
	l.Lexer = *lexer.New(input)
}

func (l *Lexer) Next() (goprompt.Token, bool) {
	t := l.Lexer.Next()
	if t.Type == token.END_OF_FILE {
		return nil, false
	}

	return t, true
}

func New(executor goprompt.Executor, opts ...goprompt.Option) *goprompt.Prompt {
	opts = append(
		[]goprompt.Option{
			goprompt.WithLexer(&Lexer{}),
			goprompt.WithExecuteOnEnterCallback(executeOnEnter),
			goprompt.WithPrefix(">> "),
			goprompt.WithKeyBindMode(goprompt.CustomOverrideEmacsKeyBind),
			goprompt.WithKeyBind(
				goprompt.KeyBind{
					Key: goprompt.AltRight,
					Fn:  cursorTokenRight,
				},
				goprompt.KeyBind{
					Key: goprompt.EscapeRight,
					Fn:  cursorTokenRight,
				},
				goprompt.KeyBind{
					Key: goprompt.ControlRight,
					Fn:  cursorTokenRight,
				},
				goprompt.KeyBind{
					Key: goprompt.EscapeF,
					Fn:  cursorTokenRight,
				},
				goprompt.KeyBind{
					Key: goprompt.AltLeft,
					Fn:  cursorTokenLeft,
				},
				goprompt.KeyBind{
					Key: goprompt.EscapeLeft,
					Fn:  cursorTokenLeft,
				},
				goprompt.KeyBind{
					Key: goprompt.ControlLeft,
					Fn:  cursorTokenLeft,
				},
				goprompt.KeyBind{
					Key: goprompt.EscapeB,
					Fn:  cursorTokenLeft,
				},
				goprompt.KeyBind{
					Key: goprompt.AltBackspace,
					Fn:  deleteTokenLeft,
				},
				goprompt.KeyBind{
					Key: goprompt.ControlBackspace,
					Fn:  deleteTokenLeft,
				},
				goprompt.KeyBind{
					Key: goprompt.ControlW,
					Fn:  deleteTokenLeft,
				},
				goprompt.KeyBind{
					Key: goprompt.EscapeD,
					Fn:  deleteTokenRight,
				},
				goprompt.KeyBind{
					Key: goprompt.ControlDelete,
					Fn:  deleteTokenRight,
				},
				goprompt.KeyBind{
					Key: goprompt.AltDelete,
					Fn:  deleteTokenRight,
				},
				goprompt.KeyBind{
					Key: goprompt.EscapeU,
					Fn:  upcaseTokenAfter,
				},
				goprompt.KeyBind{
					Key: goprompt.EscapeL,
					Fn:  downcaseTokenAfter,
				},
				goprompt.KeyBind{
					Key: goprompt.EscapeC,
					Fn:  capitalizeTokenAfter,
				},
				goprompt.KeyBind{
					Key: goprompt.EscapeT,
					Fn:  swapTokensBefore,
				},
			),
		},
		opts...,
	)

	return goprompt.New(
		executor,
		opts...,
	)
}

// Start the REPL.
func Run(executor goprompt.Executor, opts ...goprompt.Option) {
	prompt := New(executor, opts...)
	prompt.Run()
}

var KeyBindsHelp = `
Key binds:
  Keep in mind that some key binds may only work
  on specific systems or terminal emulators.

  Enter                        Run the input, or insert an indented newline when it is incomplete
  Backspace                    Delete the character before the cursor
  Up, Ctrl+p                   Move up a line, or show the older history entry on the first line
  Down, Ctrl+n                 Move down a line, or show the newer history entry on the last line
  Left, Ctrl+b                 Move one character left
  Right, Ctrl+f                Move one character right
  Home, Ctrl+a                 Move to the start of the line
  End, Ctrl+e                  Move to the end of the line
  Alt/Ctrl/Esc+Left, Esc+b     Move to the start of the token before the cursor
  Alt/Ctrl/Esc+Right, Esc+f    Move to the end of the token at or after the cursor
  Alt/Ctrl+Backspace, Ctrl+w   Delete back to the start of the token before the cursor
  Esc+d, Alt/Ctrl+Delete       Delete forward to the end of the token at or after the cursor
  Ctrl+c                       Discard the current input
  Ctrl+d                       Exit the REPL (if the input is blank)
  Ctrl+l                       Clear the screen
  Tab                          Insert an indent
  Shift+Tab                    Remove an indent from leading whitespace
  Ctrl+k, Meta+Delete          Delete from the cursor to the end of the line
  Ctrl+u, Meta+Backspace       Delete from the cursor to the start of the line
  Delete, Ctrl+d               Delete the character at the cursor
  Ctrl+o, Esc/Shift+Enter      Insert a newline without running the input
  Esc+<                        Show the oldest history entry
  Esc+>                        Show the newest history entry
  Ctrl+h, Esc+Backspace        Delete the word before the cursor
  Esc+Delete                   Delete the word after the cursor
  Esc+u                        Uppercase the token at or after the cursor
  Esc+l                        Lowercase the token at or after the cursor
  Esc+c                        Capitalize the token at or after the cursor
  Ctrl+t                       Swap the two characters before the cursor
  Esc+t                        Swap the two tokens before the cursor
`

func cursorTokenRight(p *goprompt.Prompt) bool {
	doc := p.Buffer().Document()
	cursorByteOffset := doc.CursorByteOffset()
	line := doc.CurrentLine()

	var rightBytes pstrings.ByteNumber
	lex := lexer.New(line)
	for {
		token := lex.Next()
		if token.IsEndOfFile() {
			break
		}
		endPos := token.Location().EndPos
		if endPos.ByteOffset >= int(cursorByteOffset) {
			rightBytes = pstrings.ByteNumber(endPos.ByteOffset) - cursorByteOffset
			break
		}
	}

	return p.CursorRightBytes(rightBytes + 1)
}

func cursorTokenLeft(p *goprompt.Prompt) bool {
	doc := p.Buffer().Document()
	cursorByteOffset := doc.CursorByteOffset()
	line := doc.CurrentLine()

	lex := lexer.New(line)
	var prevToken *token.Token
	var cursorToken *token.Token

	for {
		token := lex.Next()
		cursorToken = token
		if token.IsEndOfFile() {
			break
		}
		endByteOffset := pstrings.ByteNumber(token.Location().EndPos.ByteOffset)
		if endByteOffset >= cursorByteOffset {
			break
		}
		prevToken = token
	}

	if prevToken != nil {
		startOfPreviousToken := pstrings.ByteNumber(prevToken.Location().StartPos.ByteOffset)
		return p.CursorLeftBytes(cursorByteOffset - startOfPreviousToken)
	}

	startOfCurrentToken := pstrings.ByteNumber(cursorToken.Location().StartPos.ByteOffset)
	return p.CursorLeftBytes(cursorByteOffset - startOfCurrentToken)
}

func deleteTokenLeft(p *goprompt.Prompt) bool {
	buf := p.Buffer()
	line, _, cursor := cursorOnLine(buf.Document())

	lex := lexer.New(line)
	start := -1
	for {
		tok := lex.Next()
		if tok.IsEndOfFile() {
			break
		}
		tokStart := tok.Location().StartPos.ByteOffset
		if tokStart >= cursor {
			break
		}
		start = tokStart
	}
	if start < 0 {
		return false
	}

	buf.DeleteBeforeCursorRunes(
		pstrings.RuneCountInString(line[start:cursor]),
		p.TerminalColumns(),
		p.TerminalRows(),
	)
	return true
}

func deleteTokenRight(p *goprompt.Prompt) bool {
	buf := p.Buffer()
	line, _, cursor := cursorOnLine(buf.Document())
	_, end, ok := tokenAtOrAfter(line, cursor)
	if !ok {
		return false
	}

	runeCount := pstrings.RuneCountInString(line[cursor:end])
	if runeCount == 0 {
		return false
	}
	buf.DeleteRunes(
		runeCount,
		p.TerminalColumns(),
		p.TerminalRows(),
	)
	return true
}

func upcaseTokenAfter(p *goprompt.Prompt) bool {
	return transformTokenAfter(p, strings.ToUpper)
}

func downcaseTokenAfter(p *goprompt.Prompt) bool {
	return transformTokenAfter(p, strings.ToLower)
}

func capitalizeTokenAfter(p *goprompt.Prompt) bool {
	return transformTokenAfter(p, capitalizeToken)
}

func capitalizeToken(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if size == 0 {
		return s
	}
	return strings.ToUpper(string(r)) + strings.ToLower(s[size:])
}

func swapTokensBefore(p *goprompt.Prompt) bool {
	doc := p.Buffer().Document()
	line, lineStart, cursor := cursorOnLine(doc)

	type tokenSpan struct {
		start int
		end   int
	}
	lex := lexer.New(line)
	var before []tokenSpan
	for {
		tok := lex.Next()
		if tok.IsEndOfFile() {
			break
		}
		start := tok.Location().StartPos.ByteOffset
		end := tok.Location().EndPos.ByteOffset + 1
		if start >= cursor {
			break
		}
		before = append(before, tokenSpan{start, end})
	}
	if len(before) < 2 {
		return false
	}

	left := before[len(before)-2]
	right := before[len(before)-1]
	middle := line[left.end:right.start]
	replacement := line[right.start:right.end] + middle + line[left.start:left.end]
	replaceSpan(p, lineStart+left.start, lineStart+right.end, replacement)
	// A cursor that started after the span stays there. Move it back to the
	// end of the word now in second place.
	if cursor > right.end {
		p.CursorLeftBytes(pstrings.ByteNumber(cursor - right.end))
	}
	return true
}

func transformTokenAfter(p *goprompt.Prompt, transform func(string) string) bool {
	doc := p.Buffer().Document()
	line, lineStart, cursor := cursorOnLine(doc)
	start, end, ok := tokenAtOrAfter(line, cursor)
	if !ok || start >= end {
		return false
	}
	replaceSpan(p, lineStart+start, lineStart+end, transform(line[start:end]))
	return true
}

// tokenAtOrAfter returns the byte range of the token that contains the cursor
// or, when the cursor is between tokens, the next one. end is exclusive.
func tokenAtOrAfter(line string, cursor int) (start, end int, ok bool) {
	lex := lexer.New(line)
	for {
		tok := lex.Next()
		if tok.IsEndOfFile() {
			return 0, 0, false
		}
		start = tok.Location().StartPos.ByteOffset
		end = tok.Location().EndPos.ByteOffset + 1
		if end > cursor {
			return start, end, true
		}
	}
}

func cursorOnLine(doc *goprompt.Document) (line string, lineStart, cursor int) {
	line = doc.CurrentLine()
	cursor = len(doc.CurrentLineBeforeCursor())
	lineStart = int(doc.CursorByteOffset()) - cursor
	return line, lineStart, cursor
}

func replaceSpan(p *goprompt.Prompt, docStart, docEnd int, replacement string) {
	cursor := int(p.Buffer().Document().CursorByteOffset())
	if cursor > docStart {
		p.CursorLeftBytes(pstrings.ByteNumber(cursor - docStart))
	} else if cursor < docStart {
		p.CursorRightBytes(pstrings.ByteNumber(docStart - cursor))
	}

	span := p.Buffer().Document().Text[docStart:docEnd]
	p.DeleteRunes(pstrings.RuneCountInString(span))
	p.InsertTextMoveCursor(replacement, false)
	if cursor > docEnd {
		p.CursorRightBytes(pstrings.ByteNumber(cursor - docEnd))
	}
}

// A Set of keywords that end a block of code
var blockEndKeywords = ds.MakeSet(
	"end",
)

// A Set of keywords that separate multiple blocks of code
var blockSeparatorKeywords = ds.MakeSet(
	"else",
	"elsif",
	"case",
	"catch",
	"finally",
)

// Callback triggered when the Enter key is pressed.
// Decides whether the input is complete and should be executed
// or whether a newline with indentation should be added to the buffer.
func executeOnEnter(pr *goprompt.Prompt, indentSize int) (indent int, execute bool) {
	doc := pr.Buffer().Document()
	var input string
	if doc.OnLastLine() {
		input = doc.Text
	} else {
		input = doc.TextBeforeCursor()
	}

	p := parser.New(sourceName, input)
	p.Parse()

	baseIndent := doc.CurrentLineIndentSpaces()
	currentLine := doc.CurrentLine()
	lex := lexer.New(currentLine)
	firstToken := lex.Next()
	firstWord := firstToken.FetchValue()
	isBlockEnd := blockEndKeywords.Contains(firstWord)
	isBlockSeparator := blockSeparatorKeywords.Contains(firstWord)

	var movedBack bool
	if isBlockEnd || isBlockSeparator {
		var indentDiff int
		var nextIndentDiff int

		indentDiff = min(indentSize-(baseIndent%indentSize), baseIndent)
		if isBlockEnd {
			nextIndentDiff = indentDiff
		}

		if indentDiff != 0 {
			movedBack = true
		}
		toLeft := pstrings.RuneNumber(utf8.RuneCountInString(currentLine) - baseIndent + indentDiff)
		pr.CursorLeftRunes(toLeft)
		pr.InsertTextMoveCursor(currentLine[baseIndent:], false)
		pr.DeleteRunes(toLeft)
		baseIndent -= nextIndentDiff
	}

	indent = baseIndent / indentSize
	if doc.OnLastLine() {
		if p.ShouldIndent() && !movedBack {
			return indent + 1, false
		}
		if p.IsIncomplete() {
			return indent, false
		}

		lastChar, _ := utf8.DecodeLastRuneInString(currentLine)
		switch lastChar {
		case ';', '\\':
			return 0, false
		}

		return 0, true
	}

	if p.ShouldIndent() && !movedBack {
		return indent + 1, false
	}

	return indent, false
}

const (
	sourceName = "REPL"
)
