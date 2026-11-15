package checker

import (
	"fmt"
	"slices"
	"strings"

	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/position"
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
			if !strings.Contains(localNameString, name) {
				continue
			}

			desc := fmt.Sprintf("`%s` %s", local.Inspect(), local.docComment)
			result = append(result, &Completion{
				Value:       localNameString,
				Description: desc,
			})
		}
	}

	return result
}

func (l *LSPData) GetFunctionCompletions(name string, byteOffset int) []*Completion {
	var result []*Completion

	// TODO

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
