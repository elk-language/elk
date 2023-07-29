package repl

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/go-prompt"
	pstrings "github.com/elk-language/go-prompt/strings"
)

type Lexer struct {
	lexer.Lexer
}

func (l *Lexer) Init(input string) {
	l.Lexer = *lexer.New([]byte(input))
}

func (l *Lexer) Next() (prompt.Token, bool) {
	t := l.Lexer.Next()
	if t.Type == token.END_OF_FILE {
		return nil, false
	}

	return t, true
}

func Log(format string, a ...any) {
	f, err := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	fmt.Fprintf(f, format+"\n", a...)
}

func executeOnEnter(pr *prompt.Prompt, indentSize int) (indent int, execute bool) {
	doc := pr.Buffer().Document()
	if doc.OnLastLine() {
		input := doc.Text
		p := parser.New("(eval)", []byte(input))
		p.Parse()
		// Log(pp.Sprint(ast))
		// Log(pp.Sprint(LastBlockNodePosition(ast)))

		prevIndent := doc.PreviousLineIndentSpaces()
		baseIndent := doc.LastLineIndentSpaces()
		if len(input) >= 3 && input[len(input)-3:] == "end" && baseIndent != 0 || prevIndent != 0 {
			if baseIndent >= prevIndent {
				var indentDiff int
				if baseIndent != prevIndent {
					indentDiff = baseIndent - prevIndent + indentSize
					if indentDiff > baseIndent {
						indentDiff = baseIndent
					}
				}
				pr.CursorLeftRunes(pstrings.RuneNumber(indentDiff + 3))
				pr.InsertTextMoveCursor("end", true)
				pr.DeleteRunes(pstrings.RuneNumber(indentDiff))
				baseIndent -= indentSize
			} else if prevIndent > baseIndent {
				indentDiff := prevIndent - baseIndent - indentSize
				if indentDiff < 0 {
					indentDiff = 0
				}
				pr.CursorLeftRunes(3)
				pr.InsertTextMoveCursor(strings.Repeat(" ", indentDiff), false)
				pr.CursorRightRunes(3)
				baseIndent = prevIndent - 1
			}
		}

		if p.ShouldIndent() {
			return baseIndent/indentSize + 1, false
		}
		if p.IsIncomplete() {
			return baseIndent / indentSize, false
		}

		return 0, true
	}

	input := pr.Buffer().Document().TextBeforeCursor()
	p := parser.New("(eval)", []byte(input))
	p.Parse()

	baseIndent := pr.Buffer().Document().PreviousLineIndentLevel(indentSize)
	if len(input) > 3 && baseIndent > 0 && input[len(input)-3:] == "end" {
		pr.CursorLeftRunes(pstrings.RuneNumber(indentSize + 3))
		pr.InsertTextMoveCursor("end", true)
		pr.DeleteRunes(pstrings.RuneNumber(indentSize))
		baseIndent--
	}

	if p.ShouldIndent() {
		return baseIndent + 1, false
	}

	return baseIndent, false
}

func LastBlockNodePosition(node ast.Node) *position.Position {
	switch n := node.(type) {
	case *ast.ProgramNode:
		if len(n.Body) == 0 {
			return n.Position
		}
		return LastBlockNodePosition(n.Body[len(n.Body)-1])
	case *ast.ExpressionStatementNode:
		return LastBlockNodePosition(n.Expression)
	case *ast.MethodDefinitionNode:
		if len(n.Body) == 0 {
			return n.Position
		}
		return LastBlockNodePosition(n.Body[len(n.Body)-1])
	case *ast.InitDefinitionNode:
		if len(n.Body) == 0 {
			return n.Position
		}
		return LastBlockNodePosition(n.Body[len(n.Body)-1])
	case *ast.ClosureLiteralNode:
		if len(n.Body) == 0 {
			return n.Position
		}
		return LastBlockNodePosition(n.Body[len(n.Body)-1])
	case *ast.ClassDeclarationNode:
		if len(n.Body) == 0 {
			return n.Position
		}
		return LastBlockNodePosition(n.Body[len(n.Body)-1])
	case *ast.ModuleDeclarationNode:
		if len(n.Body) == 0 {
			return n.Position
		}
		return LastBlockNodePosition(n.Body[len(n.Body)-1])
	case *ast.MixinDeclarationNode:
		if len(n.Body) == 0 {
			return n.Position
		}
		return LastBlockNodePosition(n.Body[len(n.Body)-1])
	case *ast.InterfaceDeclarationNode:
		if len(n.Body) == 0 {
			return n.Position
		}
		return LastBlockNodePosition(n.Body[len(n.Body)-1])
	case *ast.StructDeclarationNode:
		if len(n.Body) == 0 {
			return n.Position
		}
		return LastBlockNodePosition(n.Body[len(n.Body)-1])
	case *ast.LoopExpressionNode:
		if len(n.ThenBody) == 0 {
			return n.Position
		}
		return LastBlockNodePosition(n.ThenBody[len(n.ThenBody)-1])
	case *ast.WhileExpressionNode:
		if len(n.ThenBody) == 0 {
			return n.Position
		}
		return LastBlockNodePosition(n.ThenBody[len(n.ThenBody)-1])
	case *ast.UntilExpressionNode:
		if len(n.ThenBody) == 0 {
			return n.Position
		}
		return LastBlockNodePosition(n.ThenBody[len(n.ThenBody)-1])
	case *ast.ForExpressionNode:
		if len(n.ThenBody) == 0 {
			return n.Position
		}
		return LastBlockNodePosition(n.ThenBody[len(n.ThenBody)-1])
	case *ast.IfExpressionNode:
		if len(n.ElseBody) > 0 {
			return LastBlockNodePosition(n.ElseBody[len(n.ElseBody)-1])
		}
		if len(n.ThenBody) > 0 {
			return LastBlockNodePosition(n.ThenBody[len(n.ThenBody)-1])
		}

		return n.Position
	default:
		return nil
	}
}

func Run() {
	p := prompt.New(
		executor,
		prompt.WithLexer(&Lexer{}),
		prompt.WithExecuteOnEnterCallback(executeOnEnter),
	)
	p.Run()
}

func executor(input string) {
}
