package prompt

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types/checker"
	goprompt "github.com/elk-language/go-prompt"
	pstrings "github.com/elk-language/go-prompt/strings"
)

type Command struct {
	Names       []string
	Description string
}

var Commands = []Command{
	{Names: []string{":quit", ":q", "exit"}, Description: "Exit the program"},
	{Names: []string{":help", ":h", "help"}, Description: "Print the help message"},
	{Names: []string{":variables", ":v"}, Description: "Inspect the currently available variables"},
	{Names: []string{":stack", ":s"}, Description: "Inspect the values on the VM's stack"},
	{Names: []string{":disassemble", ":disasm"}, Description: "Disassemble the bytecode of the previous input and print it"},
	{Names: []string{":clear", "clear"}, Description: "Clear the screen"},
}

func RenderCommandsHelp(commands []Command) string {
	nameLists := make([]string, len(commands))
	width := 0
	for i, cmd := range commands {
		nameLists[i] = strings.Join(cmd.Names, ", ")
		if n := len(nameLists[i]); n > width {
			width = n
		}
	}

	var buf strings.Builder
	buf.WriteString("\nAvailable Commands:\n")
	for i, cmd := range commands {
		fmt.Fprintf(&buf, "  %-*s   %s\n", width, nameLists[i], cmd.Description)
	}
	return buf.String()
}

func PrintHelp(banner string, commands []Command) {
	fmt.Print(banner)
	fmt.Print(KeyBindsHelp)
	fmt.Println()
	fmt.Print(RenderCommandsHelp(commands))
	fmt.Println()
}

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

func exitChecker(input string, breakline bool) bool {
	if !breakline {
		return false
	}

	switch input {
	case ":c", ":continue":
		return true
	case ":q", ":quit", "exit":
		os.Exit(0)
		return true
	default:
		return false
	}
}

func New(evaluator *Evaluator, opts ...goprompt.Option) *goprompt.Prompt {
	opts = append(
		[]goprompt.Option{
			goprompt.WithLexer(&Lexer{}),
			goprompt.WithExecuteOnEnterCallback(executeOnEnter),
			goprompt.WithPrefix(">> "),
			goprompt.WithCompleter(evaluator.Complete),
			goprompt.WithExitChecker(exitChecker),
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
		evaluator.Executor,
		opts...,
	)
}

// Start the REPL.
func Run(evaluator *Evaluator, opts ...goprompt.Option) {
	prompt := New(evaluator, opts...)
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
		if strings.HasSuffix(doc.Text, ";;") {
			return 0, true
		}
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

type Evaluator struct {
	Commands                         []Command
	Executor                         goprompt.Executor
	SourceMap                        map[string]string
	InputIndex                       int
	ElkTypechecker                   *checker.Checker
	Prepend                          []ast.StatementNode
	Prompt                           *goprompt.Prompt
	Completer                        goprompt.Completer
	LastCompletionSource             string
	LastCompletionParsedProgramNode  *ast.ProgramNode
	LastCompletionCheckedProgramNode *ast.ProgramNode
}

const SourceNamePrefix = "<repl:"

func (e *Evaluator) SourceName() string {
	return fmt.Sprintf("%s%d>", SourceNamePrefix, e.InputIndex)
}

func (e *Evaluator) DeleteSource(sourceName string) {
	if e.ElkTypechecker != nil && !e.ElkTypechecker.DefinedMacros() {
		delete(e.SourceMap, sourceName)
	}
}

func (e *Evaluator) AddSource(input string) string {
	sourceName := e.SourceName()
	e.InputIndex++
	e.SourceMap[sourceName] = input
	return sourceName
}

func (e *Evaluator) DeepCopyPrepend() []ast.StatementNode {
	if len(e.Prepend) < 1 {
		return nil
	}

	result := make([]ast.StatementNode, len(e.Prepend))

	for i, statement := range e.Prepend {
		stmt := ast.DeepCopyGeneric(statement)
		stmt.SetPrepended(true)
		result[i] = stmt
	}

	return result
}

func (e *Evaluator) UsingToPrepend(using []*ast.UsingExpressionNode) []ast.StatementNode {
	if len(using) < 1 {
		return nil
	}

	result := make([]ast.StatementNode, len(using))

	for i, usingElement := range using {
		copy := ast.DeepCopy(usingElement).(*ast.UsingExpressionNode)
		copy.SetType(nil)

		result[i] = ast.NewExpressionStatementNode(
			copy.Location(),
			copy,
		)
	}

	return result
}

func normaliseSuggestionDescription(desc string) string {
	var buf bytes.Buffer

	var lastChar rune
	for _, char := range desc {
		if unicode.IsSpace(char) {
			goto cont
		}
		if unicode.IsSpace(lastChar) {
			buf.WriteRune(' ')
		}

		buf.WriteRune(char)
	cont:
		lastChar = char
	}

	return buf.String()
}

type completionTextField uint8

const (
	completionTextFieldValue completionTextField = iota
	completionTextFieldLabel
	completionTextFieldIdentifier
)

func appendCompletionsToSuggestions(target *[]goprompt.Suggest, source []*checker.Completion, completionTextField completionTextField) {
	for _, completion := range source {
		var text string
		switch completionTextField {
		case completionTextFieldValue:
			text = completion.Value
		case completionTextFieldLabel:
			text = completion.Label
		case completionTextFieldIdentifier:
			text = completion.Identifier
		}

		alreadyPresent := slices.ContainsFunc(*target, func(suggest goprompt.Suggest) bool {
			return suggest.Label == completion.Label
		})

		if alreadyPresent {
			continue
		}

		*target = append(*target, goprompt.Suggest{
			Label:       completion.Label,
			Text:        text,
			Description: normaliseSuggestionDescription(completion.Description),
		})
	}
}

const completionSourceName = "<completion>"

func (e *Evaluator) completeCommands(source string, suggestions *[]goprompt.Suggest) (startChar, endChar pstrings.RuneNumber) {
	if len(source) > 0 && len(source) < 15 && !strings.ContainsAny(source, " \n\r\t\v") {
		lowerSource := strings.ToLower(source)

		var found bool
	commandLoop:
		for _, command := range e.Commands {
			for _, commandName := range command.Names {
				lowerCommandName := strings.ToLower(commandName)
				if strings.Contains(lowerCommandName, lowerSource) {
					*suggestions = append(*suggestions, goprompt.Suggest{
						Text:        commandName,
						Description: normaliseSuggestionDescription(command.Description),
					})
					found = true
					continue commandLoop
				}
			}
		}

		if found {
			return 0, pstrings.RuneCountInString(source)
		}
	}

	return 0, 0
}

func isLookupNode(node ast.Node) bool {
	switch node.(type) {
	case *ast.MethodCallNode, *ast.ReceiverlessMethodCallNode,
		*ast.AttributeAccessNode, *ast.ConstantLookupNode, *ast.ScopedMacroCallNode,
		*ast.MethodLookupNode:
		return true
	default:
		return false
	}
}

// Handles autocomplete
func (e *Evaluator) Complete(doc goprompt.Document) (suggestions []goprompt.Suggest, startChar, endChar pstrings.RuneNumber) {
	source := doc.Text

	startChar, endChar = e.completeCommands(doc.Text, &suggestions)

	var parsedProgramNode *ast.ProgramNode
	_ = parsedProgramNode
	var checkedProgramNode *ast.ProgramNode
	if e.LastCompletionSource == source && e.LastCompletionParsedProgramNode != nil {
		parsedProgramNode = e.LastCompletionParsedProgramNode
		checkedProgramNode = e.LastCompletionCheckedProgramNode
	} else {
		parsedProgramNode, checkedProgramNode = e.ElkTypechecker.CheckSourceREPLAutocomplete(completionSourceName, source, e.DeepCopyPrepend())
		e.LastCompletionSource = source
		e.LastCompletionParsedProgramNode = parsedProgramNode
		e.LastCompletionCheckedProgramNode = checkedProgramNode
	}
	e.ElkTypechecker.ClearErrors()

	cursorByteOffset := doc.CursorByteOffset()

	var narrowestParsedNode ast.Node
	var narrowestParsedLookupNode ast.Node
	ast.Traverse(
		parsedProgramNode,
		func(node, parent ast.Node) ast.TraverseOption {
			loc := node.Location()
			if loc == nil || cursorByteOffset < pstrings.ByteNumber(loc.StartPos.ByteOffset) || cursorByteOffset-1 > pstrings.ByteNumber(loc.EndPos.ByteOffset) {
				return ast.TraverseSkip
			}

			if narrowestParsedNode == nil {
				narrowestParsedNode = node
				if isLookupNode(node) {
					narrowestParsedLookupNode = node
				}
				return ast.TraverseContinue
			}

			if node.Location().ByteCount() <= narrowestParsedNode.Location().ByteCount() {
				narrowestParsedNode = node
				if isLookupNode(node) {
					narrowestParsedLookupNode = node
				}
			}

			return ast.TraverseContinue
		},
		nil,
	)

	var narrowestCheckedNode ast.Node
	var narrowestCheckedLookupNode ast.Node
	ast.Traverse(
		checkedProgramNode,
		func(node, parent ast.Node) ast.TraverseOption {
			if stmt, ok := node.(ast.StatementNode); ok && stmt.Prepended() {
				return ast.TraverseSkip
			}

			loc := node.Location()
			if loc == nil || cursorByteOffset < pstrings.ByteNumber(loc.StartPos.ByteOffset) || cursorByteOffset-1 > pstrings.ByteNumber(loc.EndPos.ByteOffset) {
				return ast.TraverseSkip
			}

			if narrowestCheckedNode == nil {
				narrowestCheckedNode = node
				if isLookupNode(node) {
					narrowestCheckedLookupNode = node
				}
				return ast.TraverseContinue
			}

			if node.Location().ByteCount() <= narrowestCheckedNode.Location().ByteCount() {
				narrowestCheckedNode = node
				if isLookupNode(node) {
					narrowestCheckedLookupNode = node
				}
			}

			return ast.TraverseContinue
		},
		nil,
	)

	if narrowestCheckedNode == nil {
		return
	}

	switch n := narrowestParsedLookupNode.(type) {
	case *ast.ConstantLookupNode:
		return e.completeConstantLookup(doc, n, &suggestions)
	case *ast.MethodLookupNode:
		return e.completeMethodLookup(doc, n, &suggestions)
	case *ast.ScopedMacroCallNode:
		return e.completeScopedMacroCall(doc, n, &suggestions)
	case *ast.AttributeAccessNode:
		receiver := getReceiverOfMethodCall(narrowestCheckedLookupNode)
		if receiver == nil {
			return
		}
		return e.completeMethodCall(doc, receiver, n.AttributeName, &suggestions)
	case *ast.MethodCallNode:
		receiver := getReceiverOfMethodCall(narrowestCheckedLookupNode)
		if receiver == nil {
			return
		}
		return e.completeMethodCall(doc, receiver, n.MethodName, &suggestions)
	case *ast.ReceiverlessMethodCallNode:
		return e.completeFunction(doc, n.MethodName, &suggestions)
	}

	switch n := narrowestParsedNode.(type) {
	case *ast.PublicConstantNode:
		return e.completeConstant(doc, n.Value, n.Location(), &suggestions)
	case *ast.PrivateConstantNode:
		return e.completeConstant(doc, n.Value, n.Location(), &suggestions)
	case *ast.ReceiverlessMethodCallNode:
		return e.completeFunction(doc, n.MethodName, &suggestions)
	case *ast.PublicIdentifierNode:
		return e.completeLocalOrFunction(doc, n.Value, n.Location(), &suggestions)
	case *ast.PrivateIdentifierNode:
		return e.completeLocalOrFunction(doc, n.Value, n.Location(), &suggestions)
	default:
		return
	}
}

func getReceiverOfMethodCall(node ast.Node) ast.ExpressionNode {
	switch n := node.(type) {
	case *ast.AttributeAccessNode:
		return n.Receiver
	case *ast.MethodCallNode:
		return n.Receiver
	default:
		return nil
	}
}

func (e *Evaluator) completeConstant(doc goprompt.Document, name string, loc *position.Location, suggestions *[]goprompt.Suggest) (s []goprompt.Suggest, startChar, endChar pstrings.RuneNumber) {
	cursorByteOffset := doc.CursorByteOffset()
	var currentName string
	currentName, startChar, endChar, _ = getCompletionStartAndEnd(doc, name, loc)

	node, _ := e.ElkTypechecker.ASTCache.GetUnsafe(completionSourceName)
	lspData := node.LSPData.(*checker.LSPData)
	if lspData == nil {
		return
	}

	completions := lspData.GetConstantCompletions(e.ElkTypechecker, currentName, int(cursorByteOffset))
	appendCompletionsToSuggestions(suggestions, completions, completionTextFieldValue)

	return *suggestions, startChar, endChar
}

func constantToString(constant ast.ComplexConstantNode) string {
	switch constant := constant.(type) {
	case *ast.PublicConstantNode:
		return constant.Value
	case *ast.PrivateConstantNode:
		return constant.Value
	}

	return ""
}

func (e *Evaluator) completeConstantLookup(doc goprompt.Document, node *ast.ConstantLookupNode, suggestions *[]goprompt.Suggest) (s []goprompt.Suggest, startChar, endChar pstrings.RuneNumber) {
	name := constantToString(node.Right)

	var currentName string
	var cursorIsWithin bool
	currentName, startChar, endChar, cursorIsWithin = getCompletionStartAndEnd(doc, name, node.Right.Location())
	if !cursorIsWithin {
		return
	}

	program, _ := e.ElkTypechecker.ASTCache.GetUnsafe(completionSourceName)
	lspData := program.LSPData.(*checker.LSPData)
	if lspData == nil {
		return
	}

	completions := e.ElkTypechecker.GetConstantLookupCompletions(node.Left, currentName)
	appendCompletionsToSuggestions(suggestions, completions, completionTextFieldValue)

	if currentName == "" {
		completions := e.ElkTypechecker.GetScopedMacroCompletions(node.Left, currentName)
		appendCompletionsToSuggestions(suggestions, completions, completionTextFieldIdentifier)
	}

	return *suggestions, startChar, endChar
}

func (e *Evaluator) completeMethodLookup(doc goprompt.Document, node *ast.MethodLookupNode, suggestions *[]goprompt.Suggest) (s []goprompt.Suggest, startChar, endChar pstrings.RuneNumber) {
	name := ast.IdentifierToString(node.Name)

	var currentName string
	var cursorIsWithin bool
	currentName, startChar, endChar, cursorIsWithin = getCompletionStartAndEnd(doc, name, node.Name.Location())
	if !cursorIsWithin {
		return
	}

	program, _ := e.ElkTypechecker.ASTCache.GetUnsafe(completionSourceName)
	lspData := program.LSPData.(*checker.LSPData)
	if lspData == nil {
		return
	}

	completions := e.ElkTypechecker.GetScopedMacroCompletions(node.Receiver, currentName)
	appendCompletionsToSuggestions(suggestions, completions, completionTextFieldIdentifier)
	return *suggestions, startChar, endChar
}

func (e *Evaluator) completeScopedMacroCall(doc goprompt.Document, node *ast.ScopedMacroCallNode, suggestions *[]goprompt.Suggest) (s []goprompt.Suggest, startChar, endChar pstrings.RuneNumber) {
	name := ast.IdentifierToString(node.MacroName)

	var currentName string
	var cursorIsWithin bool
	currentName, startChar, endChar, cursorIsWithin = getCompletionStartAndEnd(doc, name, node.MacroName.Location())
	if !cursorIsWithin {
		return
	}

	program, _ := e.ElkTypechecker.ASTCache.GetUnsafe(completionSourceName)
	lspData := program.LSPData.(*checker.LSPData)
	if lspData == nil {
		return
	}

	completions := e.ElkTypechecker.GetScopedMacroCompletions(node.Receiver, currentName)
	appendCompletionsToSuggestions(suggestions, completions, completionTextFieldValue)

	return *suggestions, startChar, endChar
}

func (e *Evaluator) completeFunction(doc goprompt.Document, nameNode ast.IdentifierNode, suggestions *[]goprompt.Suggest) (s []goprompt.Suggest, startChar, endChar pstrings.RuneNumber) {
	name := ast.IdentifierToString(nameNode)

	var currentName string
	var cursorIsWithin bool
	currentName, startChar, endChar, cursorIsWithin = getCompletionStartAndEnd(doc, name, nameNode.Location())
	if !cursorIsWithin {
		return
	}

	node, _ := e.ElkTypechecker.ASTCache.GetUnsafe(completionSourceName)
	lspData := node.LSPData.(*checker.LSPData)
	if lspData == nil {
		return
	}

	completions := lspData.GetFunctionCompletions(e.ElkTypechecker, currentName, 0)
	appendCompletionsToSuggestions(suggestions, completions, completionTextFieldLabel)

	return *suggestions, startChar, endChar
}

func (e *Evaluator) completeLocalOrFunction(doc goprompt.Document, name string, loc *position.Location, suggestions *[]goprompt.Suggest) (s []goprompt.Suggest, startChar, endChar pstrings.RuneNumber) {
	cursorByteOffset := doc.CursorByteOffset()
	var currentName string
	var cursorIsWithin bool
	currentName, startChar, endChar, cursorIsWithin = getCompletionStartAndEnd(doc, name, loc)
	if !cursorIsWithin {
		return
	}

	node, _ := e.ElkTypechecker.ASTCache.GetUnsafe(completionSourceName)
	lspData := node.LSPData.(*checker.LSPData)
	if lspData == nil {
		return
	}

	completions := lspData.GetLocalCompletions(currentName, int(cursorByteOffset))
	appendCompletionsToSuggestions(suggestions, completions, completionTextFieldValue)

	completions = lspData.GetFunctionCompletions(e.ElkTypechecker, currentName, 0)
	appendCompletionsToSuggestions(suggestions, completions, completionTextFieldLabel)

	return *suggestions, startChar, endChar
}

func (e *Evaluator) completeMethodCall(doc goprompt.Document, receiverNode, nameNode ast.Node, suggestions *[]goprompt.Suggest) (s []goprompt.Suggest, startChar, endChar pstrings.RuneNumber) {
	receiverType := e.ElkTypechecker.TypeOf(receiverNode)

	var currentName string
	switch n := nameNode.(type) {
	case *ast.PublicIdentifierNode:
		currentName = n.Value
	case *ast.PrivateIdentifierNode:
		currentName = n.Value
	}

	var cursorIsWithin bool
	currentName, startChar, endChar, cursorIsWithin = getCompletionStartAndEnd(doc, currentName, nameNode.Location())
	if !cursorIsWithin {
		return
	}

	completions := e.ElkTypechecker.GetMethodCompletions(receiverType, currentName)
	appendCompletionsToSuggestions(suggestions, completions, completionTextFieldValue)

	return *suggestions, startChar, endChar
}

func getCompletionStartAndEnd(doc goprompt.Document, name string, loc *position.Location) (currentName string, startChar, endChar pstrings.RuneNumber, cursorIsWithin bool) {
	currentName = name
	source := doc.Text
	cursorRuneOffset := doc.CursorRuneOffset()
	cursorByteOffset := doc.CursorByteOffset()

	if name == "" {
		cursorIsWithin = true
	}

	if loc.ByteIsWithin(int(cursorByteOffset - 1)) {
		startChar = pstrings.RuneNumber(loc.StartPos.RuneOffset(source))
		endChar = pstrings.RuneNumber(loc.EndPos.RuneOffset(source)) + 1
		truncate := loc.EndPos.ByteOffset - int(cursorByteOffset) + 1
		if truncate > 0 {
			currentName = currentName[:len(currentName)-truncate]
		}
		cursorIsWithin = true
	} else if cursorByteOffset > pstrings.ByteNumber(loc.EndPos.ByteOffset) {
		startChar = pstrings.RuneNumber(loc.StartPos.RuneOffset(source))
		endChar = pstrings.RuneNumber(loc.EndPos.RuneOffset(source)) + 1
	} else {
		startChar = cursorRuneOffset
		endChar = cursorRuneOffset
	}

	return
}
