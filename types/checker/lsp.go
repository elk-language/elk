package checker

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value/symbol"
)

// LSPData contains additional typechecker
// data used by the LSP for a particular file/program node.
type LSPData struct {
	localEnvs      concurrent.Slice[*localEnvironment]
	methodScopes   concurrent.Slice[*methodScope]
	constantScopes concurrent.Slice[*constantScope]
}

func (l *LSPData) LSPData() {}

func NewLSPData() *LSPData {
	return &LSPData{}
}

type Completion struct {
	Value       string
	Description string
}

func (l *LSPData) GetLocalCompletions(name string, byteOffset int) []*Completion {
	var result []*Completion

	lowerName := strings.ToLower(name)
	for _, localEnv := range slices.Backward(l.localEnvs.Slice) {
		if byteOffset > 0 && localEnv.span != nil && !localEnv.span.ByteIsWithin(byteOffset) {
			continue
		}

		for localName, local := range symbol.SortKeysIter(localEnv.locals) {
			// declaration of the variable is after the current cursor
			// so it is not available
			if byteOffset > 0 && local.declaration != nil && byteOffset < local.declaration.Span.EndPos.ByteOffset {
				continue
			}
			localNameString := localName.String()
			lowerLocalName := strings.ToLower(localNameString)
			if !strings.Contains(lowerLocalName, lowerName) {
				continue
			}

			desc := fmt.Sprintf("`%s` %s", local.Inspect(), local.docComment)
			result = append(result, &Completion{
				Value:       localNameString,
				Description: desc,
			})
		}
	}

	sortCompletions(result, lowerName)

	return result
}

func sortCompletions(completions []*Completion, query string) {
	slices.SortFunc(completions, func(a, b *Completion) int {
		if c := cmp.Compare(strings.Index(strings.ToLower(a.Value), query), strings.Index(strings.ToLower(b.Value), query)); c != 0 {
			return c
		}
		return cmp.Compare(len(a.Value), len(b.Value))
	})
}

func (l *LSPData) GetFunctionCompletions(checker *Checker, name string, byteOffset int) []*Completion {
	var result []*Completion

	lowerName := strings.ToLower(name)
	for _, methodScope := range slices.Backward(l.methodScopes.Slice) {
		if byteOffset > 0 && methodScope.span != nil && !methodScope.span.ByteIsWithin(byteOffset) {
			continue
		}
		switch methodScope.kind {
		case scopeUsingBufferKind, scopeUsingKind:
		default:
			continue
		}

		for methodName, method := range checker.MethodsInNamespace(methodScope.container) {
			if method.IsRegisteredOverload() {
				continue
			}

			methodNameString := methodName.String()
			lowerMethodName := strings.ToLower(methodNameString)
			if !strings.Contains(lowerMethodName, lowerName) {
				continue
			}

			doc := fmt.Sprintf("`%s` %s", types.Inspect(method), method.DocComment)
			result = append(result, &Completion{
				Value:       methodNameString,
				Description: doc,
			})
		}
	}

	sortCompletions(result, lowerName)

	return result
}

func (c *Checker) GetMethodCompletions(receiverType types.Type, name string) []*Completion {
	var result []*Completion

	lowerName := strings.ToLower(name)

	receiverType = c.ToNonLiteral(receiverType, true)
	receiverNamespace, ok := receiverType.(types.Namespace)
	if !ok {
		return nil
	}

	for name, method := range c.MethodsInNamespace(receiverNamespace) {
		if method.IsRegisteredOverload() {
			continue
		}
		if strings.Contains(strings.ToLower(name.String()), lowerName) {
			doc := fmt.Sprintf("`%s` %s", types.Inspect(method), method.DocComment)
			result = append(
				result,
				&Completion{
					Value:       name.String(),
					Description: doc,
				},
			)
		}
	}

	sortCompletions(result, lowerName)

	return result
}

func (l *LSPData) finalise() {
	l.sortByPosition()
}

func (l *LSPData) sortByPosition() {
	slices.SortFunc(l.localEnvs.Slice, func(a, b *localEnvironment) int {
		return sortBySpan(a.span, b.span)
	})
	slices.SortFunc(l.methodScopes.Slice, func(a, b *methodScope) int {
		return sortBySpan(a.span, b.span)
	})
	slices.SortFunc(l.constantScopes.Slice, func(a, b *constantScope) int {
		return sortBySpan(a.span, b.span)
	})
}

func sortBySpan(a, b *position.Span) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return 1
	}
	if b == nil {
		return -1
	}
	if a.StartPos.ByteOffset > b.StartPos.ByteOffset {
		return 1
	}
	if b.StartPos.ByteOffset > a.StartPos.ByteOffset {
		return -1
	}

	aByteCount := a.ByteCount()
	bByteCount := b.ByteCount()
	if aByteCount > bByteCount {
		return -1
	}
	if bByteCount > aByteCount {
		return 1
	}
	return 0
}
