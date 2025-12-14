package value

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/elk-language/elk/indent"
)

// Maps instance variable names to their indices
type IvarIndices map[Symbol]int

func (n IvarIndices) SetIndex(name Symbol, i int) {
	n[name] = i
}

func (n IvarIndices) GetIndex(name Symbol) int {
	return n[name]
}

func (n IvarIndices) GetIndexOk(name Symbol) (int, bool) {
	val, ok := n[name]
	return val, ok
}

func (n IvarIndices) ToGoSource() string {
	var buff strings.Builder

	buff.WriteString("IvarIndices{")
	for key, val := range n {
		fmt.Fprintf(&buff, "value.ToSymbol(%q): %d,", key.String(), val)
	}
	buff.WriteString("}")

	return buff.String()
}

func (n IvarIndices) GetName(index int) Symbol {
	name, _ := n.GetNameOk(index)
	return name
}

func (n IvarIndices) GetNameOk(index int) (Symbol, bool) {
	for ivarName, i := range n {
		if i == index {
			return ivarName, true
		}
	}

	return 0, false
}

func (n *IvarIndices) Copy() Reference {
	return n
}

func (*IvarIndices) Class() *Class {
	return nil
}

func (*IvarIndices) DirectClass() *Class {
	return nil
}

func (*IvarIndices) SingletonClass() *Class {
	return nil
}

func (in *IvarIndices) Length() int {
	return len(*in)
}

const MAX_IVAR_INDICES_ELEMENTS_IN_INSPECT = 300

func (in *IvarIndices) Inspect() string {
	var hasMultilineElements bool
	keyStrings := make(
		[]string,
		0,
		min(MAX_IVAR_INDICES_ELEMENTS_IN_INSPECT, in.Length()),
	)
	valStrings := make(
		[]string,
		0,
		min(MAX_IVAR_INDICES_ELEMENTS_IN_INSPECT, in.Length()),
	)

	i := 0
	for key, val := range *in {
		keyString := key.Inspect()
		keyStrings = append(keyStrings, keyString)

		valString := strconv.Itoa(val)
		valStrings = append(valStrings, valString)

		if strings.ContainsRune(keyString, '\n') ||
			strings.ContainsRune(valString, '\n') {
			hasMultilineElements = true
		}

		if i >= MAX_IVAR_INDICES_ELEMENTS_IN_INSPECT-1 {
			break
		}
		i++
	}

	var buff strings.Builder

	buff.WriteString("IvarIndices{")
	if hasMultilineElements || in.Length() > 15 {
		buff.WriteRune('\n')
		for i := range len(keyStrings) {
			keyString := keyStrings[i]
			valString := valStrings[i]

			if i != 0 {
				buff.WriteString(",\n")
			}
			indent.IndentString(&buff, keyString, 1)
			buff.WriteString(" => ")
			indent.IndentStringFromSecondLine(&buff, valString, 1)

			if i >= MAX_IVAR_INDICES_ELEMENTS_IN_INSPECT-1 {
				buff.WriteString(",\n  ...")
				break
			}
		}
		buff.WriteRune('\n')
	} else {
		for i := range len(keyStrings) {
			keyString := keyStrings[i]
			valString := valStrings[i]

			if i != 0 {
				buff.WriteString(", ")
			}
			buff.WriteString(keyString)
			buff.WriteString(" => ")
			buff.WriteString(valString)

			if i >= MAX_IVAR_INDICES_ELEMENTS_IN_INSPECT-1 {
				buff.WriteString(", ...")
				break
			}
		}
	}
	buff.WriteRune('}')

	return buff.String()
}

func (n *IvarIndices) Error() string {
	return n.Inspect()
}

func (*IvarIndices) InstanceVariables() *InstanceVariables {
	return nil
}
