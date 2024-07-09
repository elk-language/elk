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
	return namespace.Constants().Len() > 0 || namespace.Methods().Len() > 0 || namespace.Parent() != nil && namespace.Parent() != objectClass
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
	types.ForeachConstant(namespace, func(name string, typ types.Type) {
		fmt.Fprintf(
			buffer,
			`namespace.DefineConstant(%q, %s)
			`,
			name,
			typeToCode(typ),
		)
	})

	for _, subtype := range namespace.Subtypes().Map {
		subtypeNamespace, ok := subtype.(types.Namespace)
		if !ok {
			continue
		}

		defineMethodsWithinNamespace(buffer, subtypeNamespace, env, false)
	}

	buffer.WriteString("}")
}

func defineMethods(buffer *bytes.Buffer, namespace types.Namespace) {
	buffer.WriteString("\n// Define methods\n")
	for _, method := range namespace.Methods().Map {
		fmt.Fprintf(
			buffer,
			"namespace.DefineMethod(%q, %t, %t, %t, %q, ",
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
					"NewParameter(value.ToSymbol(%q), %s, %s, %t)",
					param.Name,
					typeToCode(param.Type),
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
			typeToCode(method.ReturnType),
			typeToCode(method.ThrowType),
		)
	}
}

func defineSubtypesWithinNamespace(buffer *bytes.Buffer, namespace types.Namespace) {
	for name, subtype := range namespace.Subtypes().Map {
		switch s := subtype.(type) {
		case *types.Class:
			defineClass(buffer, s, name)
		case *types.Mixin:
			defineMixin(buffer, s, name)
		case *types.Module:
			defineModule(buffer, s, name)
		case *types.Interface:
			defineInterface(buffer, s, name)
		}
	}
}

func defineClass(buffer *bytes.Buffer, class *types.Class, constantName value.Symbol) {
	hasSubtypes := class.Subtypes().Len() > 0
	if hasSubtypes {
		buffer.WriteString(`{ namespace :=`)
	}

	fmt.Fprintf(
		buffer,
		`namespace.TryDefineClass(%q, %t, %t, %t, %q, objectClass, env)
		`,
		class.DocComment(),
		class.IsAbstract(),
		class.IsSealed(),
		class.IsPrimitive(),
		constantName,
	)

	defineSubtypesWithinNamespace(buffer, class)
	if hasSubtypes {
		buffer.WriteString("}\n")
	}
}

func defineMixin(buffer *bytes.Buffer, mixin *types.Mixin, constantName value.Symbol) {
	hasSubtypes := mixin.Subtypes().Len() > 0
	if hasSubtypes {
		buffer.WriteString(`{ namespace :=`)
	}

	fmt.Fprintf(
		buffer,
		`namespace.TryDefineMixin(%q, %t, %q, env)
		`,
		mixin.DocComment(),
		mixin.IsAbstract(),
		constantName,
	)

	defineSubtypesWithinNamespace(buffer, mixin)
	if hasSubtypes {
		buffer.WriteString("}\n")
	}
}

func defineModule(buffer *bytes.Buffer, module *types.Module, constantName value.Symbol) {
	hasSubtypes := module.Subtypes().Len() > 0
	if hasSubtypes {
		buffer.WriteString(`{ namespace :=`)
	}

	fmt.Fprintf(
		buffer,
		`namespace.TryDefineModule(%q, %q)
		`,
		module.DocComment(),
		constantName,
	)

	defineSubtypesWithinNamespace(buffer, module)
	if hasSubtypes {
		buffer.WriteString("}\n")
	}
}

func defineInterface(buffer *bytes.Buffer, iface *types.Interface, constantName value.Symbol) {
	hasSubtypes := iface.Subtypes().Len() > 0
	if hasSubtypes {
		buffer.WriteString(`{ namespace :=`)
	}

	fmt.Fprintf(
		buffer,
		`namespace.TryDefineInterface(%q, %q, env)
		`,
		iface.DocComment(),
		constantName,
	)

	defineSubtypesWithinNamespace(buffer, iface)
	if hasSubtypes {
		buffer.WriteString("}\n")
	}
}

// Serialise the type to Go code
func typeToCode(typ types.Type) string {
	switch t := typ.(type) {
	case nil:
		return "nil"
	case types.Any:
		return "Any{}"
	case types.Void:
		return "Void{}"
	case types.Never:
		return "Never{}"
	case *types.NamedType:
		return fmt.Sprintf(
			"NewNamedType(%q, %s)",
			t.Name,
			typeToCode(t.Type),
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
			typeToCode(t.Type),
		)
	case *types.Union:
		buff := new(strings.Builder)
		buff.WriteString("NewUnion(")
		for _, element := range t.Elements {
			fmt.Fprintf(
				buff,
				"%s, ",
				typeToCode(element),
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
				typeToCode(element),
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