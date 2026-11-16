package checker

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/parser/ast"
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
	Label       string
	Identifier  string
	Value       string
	Description string
}

func (l *LSPData) GetLocalCompletions(name string, byteOffset int) []*Completion {
	var result []*Completion

	seenLocals := ds.MakeSet[symbol.Symbol]()
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
			if seenLocals.Contains(localName) {
				continue
			}

			seenLocals.Add(localName)
			localNameString := localName.String()
			if strings.HasPrefix(localNameString, "#") {
				continue
			}

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
	seenMethods := ds.MakeSet[symbol.Symbol]()
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
			if seenMethods.Contains(methodName) {
				continue
			}

			seenMethods.Add(methodName)
			methodNameString := methodName.String()
			if strings.HasPrefix(methodNameString, "#") {
				continue
			}

			lowerMethodName := strings.ToLower(methodNameString)
			if !strings.Contains(lowerMethodName, lowerName) {
				continue
			}

			var label string
			var value string
			var identifier string
			if method.IsMacro() {
				label = methodNameString + "()"
				identifier = methodNameString
				value = methodNameString[:len(methodNameString)-1]
			} else {
				label = methodNameString + "()"
				identifier = methodNameString
				value = methodNameString
			}
			doc := fmt.Sprintf("`%s` %s", types.Inspect(method), method.DocComment)
			result = append(result, &Completion{
				Label:       label,
				Identifier:  identifier,
				Value:       value,
				Description: doc,
			})
		}
	}

	sortCompletions(result, lowerName)

	return result
}

func (l *LSPData) GetConstantCompletions(checker *Checker, name string, byteOffset int) []*Completion {
	var result []*Completion

	lowerName := strings.ToLower(name)

	seenConstants := ds.MakeSet[symbol.Symbol]()

	for _, constScope := range slices.Backward(l.constantScopes.Slice) {
		switch constScope.kind {
		case scopeUsingBufferKind, scopeUsingKind:
		default:
			if byteOffset > 0 && constScope.span != nil && !constScope.span.ByteIsWithin(byteOffset) {
				continue
			}
		}

		for constantName, constant := range constScope.container.Subtypes() {
			if seenConstants.Contains(constantName) {
				continue
			}

			seenConstants.Add(constantName)
			constantNameString := constantName.String()
			if strings.HasPrefix(constantNameString, "#") {
				continue
			}

			lowerConstantName := strings.ToLower(constantNameString)
			if !strings.Contains(lowerConstantName, lowerName) {
				continue
			}

			result = append(result, &Completion{
				Value:       constantNameString,
				Description: getDocForSubtype(constant),
			})
		}

		for constantName, constant := range constScope.container.Constants() {
			if seenConstants.Contains(constantName) {
				continue
			}

			seenConstants.Add(constantName)
			constantNameString := constantName.String()
			if strings.HasPrefix(constantNameString, "#") {
				continue
			}

			lowerConstantName := strings.ToLower(constantNameString)
			if !strings.Contains(lowerConstantName, lowerName) {
				continue
			}

			doc := fmt.Sprintf("`%s` %s", constant.Inspect(), constant.DocComment())
			result = append(result, &Completion{
				Value:       constantNameString,
				Description: doc,
			})
		}
	}

	sortCompletions(result, lowerName)

	return result
}

func getDocForSubtype(constant *types.Constant) string {
	typ := constant.Type
	switch t := typ.(type) {
	case *types.Mixin:
		return fmt.Sprintf("Mixin `%s` %s", types.Inspect(typ), t.DocComment())
	case *types.Interface:
		return fmt.Sprintf("Interface `%s` %s", types.Inspect(typ), t.DocComment())
	case *types.Module:
		return fmt.Sprintf("Module `%s` %s", types.Inspect(typ), t.DocComment())
	case *types.Class:
		return fmt.Sprintf("Class `%s` %s", types.Inspect(typ), t.DocComment())
	case *types.TypeParameter:
		return fmt.Sprintf("Type Parameter `%s` of `%s`", t.InspectSignature(), types.Inspect(t.Namespace))
	case types.TypeWithDocComment:
		return fmt.Sprintf("Type `%s` %s", types.Inspect(typ), t.DocComment())
	default:
		return fmt.Sprintf("Type `%s` %s", types.Inspect(typ), constant.DocComment())
	}
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
		if !strings.Contains(strings.ToLower(name.String()), lowerName) {
			continue
		}

		methodNameString := name.String()
		doc := fmt.Sprintf("`%s` %s", types.Inspect(method), method.DocComment)

		var label string
		var value string
		var identifier string
		if method.IsMacro() {
			label = methodNameString + "()"
			identifier = methodNameString
			value = methodNameString[:len(methodNameString)-1]
		} else {
			label = methodNameString + "()"
			identifier = methodNameString
			value = methodNameString
		}
		result = append(
			result,
			&Completion{
				Label:       label,
				Value:       value,
				Identifier:  identifier,
				Description: doc,
			},
		)
	}

	sortCompletions(result, lowerName)

	return result
}

func (c *Checker) GetScopedMacroCompletions(namespaceNode ast.ExpressionNode, name string) []*Completion {
	var result []*Completion

	lowerName := strings.ToLower(name)
	namespaceNode = ast.DeepCopyGeneric(namespaceNode)

	namespaceConstant, _ := c.resolveConstant(namespaceNode)
	defer c.ClearErrors()

	if c.Errors.IsFailure() {
		return nil
	}
	if namespaceConstant == nil {
		return nil
	}

	namespace, ok := namespaceConstant.Type.(types.Namespace)
	if !ok {
		return nil
	}

	seenMacros := ds.MakeSet[symbol.Symbol]()
	for name, method := range c.MethodsInNamespace(namespace) {
		if !method.IsMacro() {
			continue
		}
		if seenMacros.Contains(name) {
			continue
		}

		seenMacros.Add(name)
		if !strings.Contains(strings.ToLower(name.String()), lowerName) {
			continue
		}

		methodNameString := name.String()
		doc := fmt.Sprintf("`%s` %s", types.Inspect(method), method.DocComment)
		label := methodNameString + "()"
		identifier := methodNameString
		value := methodNameString[:len(methodNameString)-1]
		result = append(
			result,
			&Completion{
				Label:       label,
				Value:       value,
				Identifier:  identifier,
				Description: doc,
			},
		)
	}

	sortCompletions(result, lowerName)

	return result
}

func (c *Checker) GetConstantLookupCompletions(namespaceNode ast.ExpressionNode, name string) []*Completion {
	var result []*Completion

	lowerName := strings.ToLower(name)
	namespaceNode = ast.DeepCopyGeneric(namespaceNode)

	namespaceConstant, _ := c.resolveConstantSubtype(namespaceNode)
	defer c.ClearErrors()

	if c.Errors.IsFailure() {
		return nil
	}
	if namespaceConstant == nil {
		return nil
	}

	namespace, ok := namespaceConstant.Type.(types.Namespace)
	if !ok {
		return nil
	}

	seenConstants := ds.MakeSet[symbol.Symbol]()
	for constantName, constant := range namespace.Subtypes() {
		if seenConstants.Contains(constantName) {
			continue
		}

		seenConstants.Add(constantName)
		constantNameString := constantName.String()
		if strings.HasPrefix(constantNameString, "#") {
			continue
		}

		lowerConstantName := strings.ToLower(constantNameString)
		if !strings.Contains(lowerConstantName, lowerName) {
			continue
		}

		result = append(result, &Completion{
			Value:       constantNameString,
			Description: getDocForSubtype(constant),
		})
	}

	for constantName, constant := range namespace.Constants() {
		if seenConstants.Contains(constantName) {
			continue
		}

		seenConstants.Add(constantName)
		constantNameString := constantName.String()
		if strings.HasPrefix(constantNameString, "#") {
			continue
		}

		lowerConstantName := strings.ToLower(constantNameString)
		if !strings.Contains(lowerConstantName, lowerName) {
			continue
		}

		doc := fmt.Sprintf("`%s` %s", constant.Inspect(), constant.DocComment())
		result = append(result, &Completion{
			Value:       constantNameString,
			Description: doc,
		})
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
