package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/types/checker"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// Compiles Elk headers

const pathToHeaders = "./headers/"

func main() {
	env := types.NewGlobalEnvironmentWithoutHeaders()
	items, _ := os.ReadDir(pathToHeaders)
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		if !strings.HasSuffix(item.Name(), ".elh") {
			continue
		}

		pathToFile := pathToHeaders + item.Name()
		bytes, err := os.ReadFile(pathToFile)
		if err != nil {
			panic(err)
		}
		source := string(bytes)
		_, errList := checker.CheckSource(pathToFile, source, env, true)
		if len(errList) > 0 {
			fmt.Println(errList.HumanStringWithSource(source, true))
			os.Exit(1)
		}
	}
	buffer := new(bytes.Buffer)
	buffer.WriteString(
		`
			package types

			// This file is auto-generated, please do not edit it manually

			import (
				"github.com/elk-language/elk/value"
				"github.com/elk-language/elk/value/symbol"
			)

			func setupGlobalEnvironmentFromHeaders(env *GlobalEnvironment) {
				objectClass := env.StdSubtypeClass(symbol.Object)
				namespace := env.Root
		`,
	)

	buffer.WriteString("\n// Define all namespaces\n")
	defineSubtypesWithinNamespace(buffer, env.Root)

	buffer.WriteString("\n// Define methods, constants\n")
	defineMethodsWithinNamespace(buffer, env.Root, env, true)

	buffer.WriteString(
		`
			}
		`,
	)

	os.WriteFile("types/headers.go", buffer.Bytes(), 0666)
}

func namespaceHasContent(namespace types.Namespace, env *types.GlobalEnvironment) bool {
	objectClass := env.StdSubtypeClass(symbol.Object)
	return namespace.Constants().Len() > 0 ||
		namespace.Methods().Len() > 0 ||
		namespace.InstanceVariables().Len() > 0 ||
		namespace.Parent() != nil && namespace.Parent() != objectClass
}

func defineMethodsWithinNamespace(buffer *bytes.Buffer, namespace types.Namespace, env *types.GlobalEnvironment, root bool) {
	namespaceClass, namespaceIsClass := namespace.(*types.Class)
	hasContent := namespaceHasContent(namespace, env)
	objectClass := env.StdSubtypeClass(symbol.Object)

	if !hasContent {
		return
	}

	if root {
		buffer.WriteString(
			`
				{
					namespace := env.Root
			`,
		)
	} else {
		var namespaceType string

		switch namespace.(type) {
		case *types.Class:
			namespaceType = "Class"
		case *types.Mixin:
			namespaceType = "Mixin"
		case *types.Interface:
			namespaceType = "Interface"
		case *types.Module:
			namespaceType = "Module"
		}
		fmt.Fprintf(
			buffer,
			`
				{
					namespace := namespace.SubtypeString(%q).(*%s)
			`,
			types.GetConstantName(namespace.Name()),
			namespaceType,
		)
	}
	buffer.WriteString("\nnamespace.Name() // noop - avoid unused variable error\n")
	if namespaceIsClass {
		superclass := namespaceClass.Superclass()
		if superclass != nil && superclass != objectClass {

			fmt.Fprintf(
				buffer,
				`namespace.SetParent(NameToNamespace(%q, env))
				`,
				namespaceClass.Superclass().Name(),
			)
		}
	}

	buffer.WriteString("\n// Include mixins\n")
	types.ForeachIncludedMixin(namespace, func(m *types.Mixin) {
		fmt.Fprintf(
			buffer,
			`namespace.IncludeMixin(NameToType(%q, env).(*Mixin))
			`,
			m.Name(),
		)
	})

	buffer.WriteString("\n// Implement interfaces\n")
	types.ForeachImplementedInterface(namespace, func(i *types.Interface) {
		fmt.Fprintf(
			buffer,
			`namespace.ImplementInterface(NameToType(%q, env).(*Interface))
			`,
			i.Name(),
		)
	})

	defineMethods(buffer, namespace)

	buffer.WriteString("\n// Define constants\n")
	types.ForeachConstantSorted(namespace, func(name value.Symbol, typ types.Type) {
		fmt.Fprintf(
			buffer,
			`namespace.DefineConstant(value.ToSymbol(%q), %s)
			`,
			name,
			typeToCode(typ, false),
		)
	})

	buffer.WriteString("\n// Define instance variables\n")
	types.ForeachOwnInstanceVariableSorted(namespace, func(name value.Symbol, typ types.Type) {
		fmt.Fprintf(
			buffer,
			`namespace.DefineInstanceVariable(value.ToSymbol(%q), %s)
			`,
			name,
			typeToCode(typ, false),
		)
	})

	types.ForeachSubtypeSorted(namespace, func(name value.Symbol, subtype types.Type) {
		subtypeNamespace, ok := subtype.(types.Namespace)
		if !ok {
			return
		}

		defineMethodsWithinNamespace(buffer, subtypeNamespace, env, false)
	})

	buffer.WriteString("}")
}

func defineMethods(buffer *bytes.Buffer, namespace types.Namespace) {
	buffer.WriteString("\n// Define methods\n")

	types.ForeachOwnMethodSorted(namespace, func(name value.Symbol, method *types.Method) {
		fmt.Fprintf(
			buffer,
			"namespace.DefineMethod(%q, %t, %t, %t, value.ToSymbol(%q), ",
			method.DocComment,
			method.IsAbstract(),
			method.IsSealed(),
			method.IsNative(),
			method.Name,
		)
		if len(method.Params) > 0 {
			buffer.WriteString("[]*Parameter{")
			for _, param := range method.Params {
				fmt.Fprintf(
					buffer,
					"NewParameter(value.ToSymbol(%q), %s, %s, %t),",
					param.Name,
					typeToCode(param.Type, false),
					param.Kind,
					param.InstanceVariable,
				)
			}
			buffer.WriteString("}, ")
		} else {
			buffer.WriteString("nil, ")
		}

		fmt.Fprintf(
			buffer,
			"%s, %s)\n",
			typeToCode(method.ReturnType, false),
			typeToCode(method.ThrowType, false),
		)
	})
}

func defineSubtypesWithinNamespace(buffer *bytes.Buffer, namespace types.Namespace) {
	types.ForeachSubtypeSorted(namespace, func(name value.Symbol, subtype types.Type) {
		switch s := subtype.(type) {
		case *types.Class:
			defineClass(buffer, s, name.String())
		case *types.Mixin:
			defineMixin(buffer, s, name.String())
		case *types.Module:
			defineModule(buffer, s, name.String())
		case *types.Interface:
			defineInterface(buffer, s, name.String())
		default:
			defineSubtype(buffer, s, name.String())
		}
	})
}

func defineSubtype(buffer *bytes.Buffer, subtype types.Type, name string) {
	fmt.Fprintf(
		buffer,
		`namespace.DefineSubtype(value.ToSymbol(%q), %s)
		`,
		name,
		typeToCode(subtype, true),
	)
}

func defineClass(buffer *bytes.Buffer, class *types.Class, constantName string) {
	hasSubtypes := class.Subtypes().Len() > 0
	if hasSubtypes {
		buffer.WriteString(`{ namespace :=`)
	}

	var initialSuperclass string
	if class.Superclass() == nil {
		initialSuperclass = "nil"
	} else {
		initialSuperclass = "objectClass"
	}

	fmt.Fprintf(
		buffer,
		`namespace.TryDefineClass(%q, %t, %t, %t, value.ToSymbol(%q), %s, env)
		`,
		class.DocComment(),
		class.IsAbstract(),
		class.IsSealed(),
		class.IsPrimitive(),
		constantName,
		initialSuperclass,
	)

	defineSubtypesWithinNamespace(buffer, class)
	if hasSubtypes {
		buffer.WriteString("namespace.Name() // noop - avoid unused variable error\n")
		buffer.WriteString("}\n")
	}
}

func defineMixin(buffer *bytes.Buffer, mixin *types.Mixin, constantName string) {
	hasSubtypes := mixin.Subtypes().Len() > 0
	if hasSubtypes {
		buffer.WriteString(`{ namespace :=`)
	}

	fmt.Fprintf(
		buffer,
		`namespace.TryDefineMixin(%q, %t, value.ToSymbol(%q), env)
		`,
		mixin.DocComment(),
		mixin.IsAbstract(),
		constantName,
	)

	defineSubtypesWithinNamespace(buffer, mixin)
	if hasSubtypes {
		buffer.WriteString("namespace.Name() // noop - avoid unused variable error\n")
		buffer.WriteString("}\n")
	}
}

func defineModule(buffer *bytes.Buffer, module *types.Module, constantName string) {
	hasSubtypes := module.Subtypes().Len() > 0
	if hasSubtypes {
		buffer.WriteString(`{ namespace :=`)
	}

	fmt.Fprintf(
		buffer,
		`namespace.TryDefineModule(%q, value.ToSymbol(%q))
		`,
		module.DocComment(),
		constantName,
	)

	defineSubtypesWithinNamespace(buffer, module)
	if hasSubtypes {
		buffer.WriteString("namespace.Name() // noop - avoid unused variable error\n")
		buffer.WriteString("}\n")
	}
}

func defineInterface(buffer *bytes.Buffer, iface *types.Interface, constantName string) {
	hasSubtypes := iface.Subtypes().Len() > 0
	if hasSubtypes {
		buffer.WriteString(`{ namespace :=`)
	}

	fmt.Fprintf(
		buffer,
		`namespace.TryDefineInterface(%q, value.ToSymbol(%q), env)
		`,
		iface.DocComment(),
		constantName,
	)

	defineSubtypesWithinNamespace(buffer, iface)
	if hasSubtypes {
		buffer.WriteString("namespace.Name() // noop - avoid unused variable error\n")
		buffer.WriteString("}\n")
	}
}

// Serialise the type to Go code
func typeToCode(typ types.Type, init bool) string {
	switch t := typ.(type) {
	case nil:
		return "nil"
	case types.Any:
		return "Any{}"
	case types.Void:
		return "Void{}"
	case types.Nil:
		return "Nil{}"
	case types.True:
		return "True{}"
	case types.False:
		return "False{}"
	case types.Never:
		return "Never{}"
	case types.Nothing:
		return "Nothing{}"
	case *types.Not:
		return fmt.Sprintf(
			"NewNot(%s)",
			typeToCode(t.Type, init),
		)
	case *types.NamedType:
		if init {
			return fmt.Sprintf(
				"NewNamedType(%q, %s)",
				t.Name,
				typeToCode(t.Type, init),
			)
		}

		return fmt.Sprintf(
			"NameToType(%q, env)",
			t.Name,
		)
	case *types.Class:
		return fmt.Sprintf(
			"NameToType(%q, env)",
			t.Name(),
		)
	case *types.SingletonClass:
		return fmt.Sprintf(
			"NameToNamespace(%q, env).Singleton()",
			t.AttachedObject.Name(),
		)
	case *types.Mixin:
		return fmt.Sprintf(
			"NameToType(%q, env)",
			t.Name(),
		)
	case *types.Module:
		return fmt.Sprintf(
			"NameToType(%q, env)",
			t.Name(),
		)
	case *types.Interface:
		return fmt.Sprintf(
			"NameToType(%q, env)",
			t.Name(),
		)
	case *types.Nilable:
		return fmt.Sprintf(
			"NewNilable(%s)",
			typeToCode(t.Type, init),
		)
	case *types.Union:
		buff := new(strings.Builder)
		buff.WriteString("NewUnion(")
		for _, element := range t.Elements {
			fmt.Fprintf(
				buff,
				"%s, ",
				typeToCode(element, init),
			)
		}
		buff.WriteRune(')')
		return buff.String()
	case *types.Intersection:
		buff := new(strings.Builder)
		buff.WriteString("NewIntersection(")
		for _, element := range t.Elements {
			fmt.Fprintf(
				buff,
				"%s, ",
				typeToCode(element, init),
			)
		}
		buff.WriteRune(')')
		return buff.String()
	case *types.SymbolLiteral:
		return fmt.Sprintf("NewSymbolLiteral(%q)", t.Value)
	case *types.StringLiteral:
		return fmt.Sprintf("NewStringLiteral(%q)", t.Value)
	case *types.CharLiteral:
		return fmt.Sprintf("NewCharLiteral(%q)", t.Value)
	case *types.FloatLiteral:
		return fmt.Sprintf("NewFloatLiteral(%q)", t.Value)
	case *types.Float32Literal:
		return fmt.Sprintf("NewFloat32Literal(%q)", t.Value)
	case *types.Float64Literal:
		return fmt.Sprintf("NewFloat64Literal(%q)", t.Value)
	case *types.IntLiteral:
		return fmt.Sprintf("NewIntLiteral(%q)", t.Value)
	case *types.Int64Literal:
		return fmt.Sprintf("NewInt64Literal(%q)", t.Value)
	case *types.Int32Literal:
		return fmt.Sprintf("NewInt32Literal(%q)", t.Value)
	case *types.Int16Literal:
		return fmt.Sprintf("NewInt16Literal(%q)", t.Value)
	case *types.Int8Literal:
		return fmt.Sprintf("NewInt8Literal(%q)", t.Value)
	case *types.UInt64Literal:
		return fmt.Sprintf("NewUInt64Literal(%q)", t.Value)
	case *types.UInt32Literal:
		return fmt.Sprintf("NewUInt32Literal(%q)", t.Value)
	case *types.UInt16Literal:
		return fmt.Sprintf("NewUInt16Literal(%q)", t.Value)
	case *types.UInt8Literal:
		return fmt.Sprintf("NewUInt8Literal(%q)", t.Value)
	default:
		panic(
			fmt.Sprintf("invalid type: %T", typ),
		)
	}
}
