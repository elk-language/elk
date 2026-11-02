package prompt

import (
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

// Start the REPL.
func Run(executor goprompt.Executor, opts ...goprompt.Option) {
	opts = append(
		[]goprompt.Option{
			goprompt.WithLexer(&Lexer{}),
			goprompt.WithExecuteOnEnterCallback(executeOnEnter),
			goprompt.WithPrefix(">> "),
			goprompt.WithKeyBindMode(goprompt.CustomOverrideEmacsKeyBind),
			goprompt.WithKeyBind(
				goprompt.KeyBind{
					Key: goprompt.ControlO,
					Fn:  insertNewline,
				},
				goprompt.KeyBind{
					Key: goprompt.AltRight,
					Fn:  cursorRight,
				},
				goprompt.KeyBind{
					Key: goprompt.AltLeft,
					Fn:  cursorLeft,
				},
				goprompt.KeyBind{
					Key: goprompt.AltBackspace,
					Fn:  deleteWord,
				},
				goprompt.KeyBind{
					Key: goprompt.ControlW,
					Fn:  deleteWord,
				},
			),
		},
		opts...,
	)

	p := goprompt.New(
		executor,
		opts...,
	)
	p.Run()
}

func insertNewline(p *goprompt.Prompt) bool {
	p.Buffer().NewLine(p.UserInputColumns(), p.TerminalRows(), false)
	return true
}

func cursorRight(p *goprompt.Prompt) bool {
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

func cursorLeft(p *goprompt.Prompt) bool {
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

func deleteWord(p *goprompt.Prompt) bool {
	buf := p.Buffer()
	doc := buf.Document()
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

	var token *token.Token
	if prevToken != nil {
		token = prevToken
	} else {
		token = cursorToken
	}

	startByte := token.Location().StartPos.ByteOffset
	endByte := token.Location().EndPos.ByteOffset
	runeCount := pstrings.RuneCountInString(doc.Text[startByte : endByte+1])
	buf.DeleteBeforeCursorRunes(
		runeCount,
		p.TerminalColumns(),
		p.TerminalRows(),
	)
	return true
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
