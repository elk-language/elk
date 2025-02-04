package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/types/checker"
	"github.com/elk-language/elk/value/symbol"
)

// Compiles Elk headers

func main() {
	env := types.NewGlobalEnvironmentWithoutHeaders()
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	pathToMainFile := filepath.Join(workingDir, "headers", "main.elh")
	_, errList := checker.CheckFile(pathToMainFile, env, true)
	if len(errList) > 0 {
		fmt.Println(errList.HumanString(true))
		os.Exit(1)
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
				var mixin *Mixin
				mixin.IsLiteral() // noop - avoid unused variable error
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
	_, isSingleton := namespace.(*types.SingletonClass)
	return len(namespace.Constants()) > 0 ||
		len(namespace.Methods()) > 0 ||
		len(namespace.InstanceVariables()) > 0 ||
		!isSingleton && namespace.Parent() != nil && namespace.Parent() != objectClass
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
					namespace := namespace.MustSubtype(%q).(*%s)
			`,
			types.GetConstantName(namespace.Name()),
			namespaceType,
		)
	}
	buffer.WriteString("\nnamespace.Name() // noop - avoid unused variable error\n")

	setTypeParameters(buffer, namespace)

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

	includeMixinsAndImplementInterfaces(buffer, namespace, env)
	defineMethods(buffer, namespace)
	defineConstants(buffer, namespace)
	defineInstanceVariables(buffer, namespace)
	singleton := namespace.Singleton()
	if singleton != nil {
		defineMethodsWithinSingleton(buffer, singleton, env)
	}
	defineMethodsWithinSubtypes(buffer, namespace, env)

	buffer.WriteString("}")
}

func defineMethodsWithinSingleton(buffer *bytes.Buffer, namespace types.Namespace, env *types.GlobalEnvironment) {
	hasContent := namespaceHasContent(namespace, env)

	if !hasContent {
		return
	}

	fmt.Fprint(
		buffer,
		`
				{
					namespace := namespace.Singleton()
			`,
	)
	buffer.WriteString("\nnamespace.Name() // noop - avoid unused variable error\n")

	setTypeParameters(buffer, namespace)

	includeMixinsAndImplementInterfaces(buffer, namespace, env)
	defineMethods(buffer, namespace)
	defineConstants(buffer, namespace)
	defineInstanceVariables(buffer, namespace)
	singleton := namespace.Singleton()
	if singleton != nil {
		defineMethodsWithinNamespace(buffer, singleton, env, false)
	}
	defineMethodsWithinSubtypes(buffer, namespace, env)

	buffer.WriteString("}")
}

func includeMixinsAndImplementInterfaces(buffer *bytes.Buffer, namespace types.Namespace, env *types.GlobalEnvironment) {
	buffer.WriteString("\n// Include mixins and implement interfaces\n")
	for parent := range types.Backward(types.DirectlyIncludedAndImplemented(namespace)) {
		switch p := parent.(type) {
		case *types.MixinWithWhere:
			fmt.Fprintf(
				buffer,
				`
					// %s
					mixin = NewMixin("", false, "", env)
					{
						namespace := mixin
						namespace.Name() // noop - avoid unused variable error
				`,
				p.InspectExtend(),
			)

			defineTypeParametersInSubtypes(buffer, p)
			defineMethods(buffer, p)

			buffer.WriteString("}\n")

			buffer.WriteString("IncludeMixinWithWhere(namespace, mixin, ")
			createTypeParametersForMixinWithWhere(buffer, p.Where)
			buffer.WriteString(")\n")
		case *types.MixinProxy:
			fmt.Fprintf(
				buffer,
				`IncludeMixin(namespace, %s)
				`,
				namespaceToCode(p),
			)
		case *types.InterfaceProxy:
			fmt.Fprintf(
				buffer,
				`ImplementInterface(namespace, %s)
				`,
				namespaceToCode(p),
			)
		case *types.Generic:
			switch p.Namespace.(type) {
			case *types.MixinProxy:
				fmt.Fprintf(
					buffer,
					`IncludeMixin(namespace, %s)
					`,
					namespaceToCode(p),
				)
			case *types.InterfaceProxy:
				fmt.Fprintf(
					buffer,
					`ImplementInterface(namespace, %s)
					`,
					namespaceToCode(p),
				)
			}
		}
	}
}

func setTypeParameters(buffer *bytes.Buffer, namespace types.Namespace) {
	if !namespace.IsGeneric() {
		return
	}

	buffer.WriteString("\n// Set up type parameters\nvar typeParam *TypeParameter\n")
	fmt.Fprintf(
		buffer,
		"typeParams := make([]*TypeParameter, %d)\n",
		len(namespace.TypeParameters()),
	)

	for i, param := range namespace.TypeParameters() {
		fmt.Fprintf(
			buffer,
			`
				typeParam = NewTypeParameter(value.ToSymbol(%[1]q), namespace, Never{}, Any{}, nil, %[2]s)
				typeParams[%[3]d] = typeParam
				namespace.DefineSubtype(value.ToSymbol(%[1]q), typeParam)
				namespace.DefineConstant(value.ToSymbol(%[1]q), NoValue{})
			`,
			param.Name.String(),
			param.Variance.String(),
			i,
		)

		if !types.IsNever(param.LowerBound) {
			fmt.Fprintf(
				buffer,
				"typeParam.LowerBound = %s\n",
				typeToCode(param.LowerBound, false),
			)
		}
		if !types.IsAny(param.UpperBound) {
			fmt.Fprintf(
				buffer,
				"typeParam.UpperBound = %s\n",
				typeToCode(param.UpperBound, false),
			)
		}
		if param.Default != nil {
			fmt.Fprintf(
				buffer,
				"typeParam.Default = %s\n",
				typeToCode(param.Default, false),
			)
		}
	}

	buffer.WriteString("\nnamespace.SetTypeParameters(typeParams)\n\n")
}

func createTypeParametersForMixinWithWhere(buffer *bytes.Buffer, typeParams []*types.TypeParameter) {
	buffer.WriteString(`[]*TypeParameter{`)
	for _, param := range typeParams {
		fmt.Fprintf(
			buffer,
			"NewTypeParameter(value.ToSymbol(%q), mixin, %s, %s, %s, %s)",
			param.Name.String(),
			typeToCode(param.LowerBound, false),
			typeToCode(param.UpperBound, false),
			typeToCode(param.Default, false),
			param.Variance.String(),
		)
	}
	buffer.WriteString("}")
}

func createTypeParameters(buffer *bytes.Buffer, typeParams []*types.TypeParameter) {
	buffer.WriteString(`[]*TypeParameter{`)
	for _, param := range typeParams {
		fmt.Fprintf(
			buffer,
			"%s,",
			typeToCode(param, true),
		)
	}
	buffer.WriteString("}")
}

func defineTypeParametersInSubtypes(buffer *bytes.Buffer, namespace types.Namespace) {
	for _, subtype := range types.SortedSubtypes(namespace) {
		param, ok := subtype.Type.(*types.TypeParameter)
		if !ok {
			continue
		}

		defineSubtype(buffer, param, param.Name.String())
	}
}

func defineConstants(buffer *bytes.Buffer, namespace types.Namespace) {
	buffer.WriteString("\n// Define constants\n")
	for name, typ := range types.SortedConstants(namespace) {
		fmt.Fprintf(
			buffer,
			"namespace.DefineConstant(value.ToSymbol(%q), %s)\n",
			name.String(),
			typeToCode(typ.Type, false),
		)
	}
}

func defineInstanceVariables(buffer *bytes.Buffer, namespace types.Namespace) {
	buffer.WriteString("\n// Define instance variables\n")
	for name, typ := range types.SortedOwnInstanceVariables(namespace) {
		fmt.Fprintf(
			buffer,
			"namespace.DefineInstanceVariable(value.ToSymbol(%q), %s)\n",
			name,
			typeToCode(typ, false),
		)
	}
}

func defineMethodsWithinSubtypes(buffer *bytes.Buffer, namespace types.Namespace, env *types.GlobalEnvironment) {
	for _, subtype := range types.SortedSubtypes(namespace) {
		subtypeNamespace, ok := subtype.Type.(types.Namespace)
		if !ok {
			continue
		}
		if subtypeNamespace == namespace {
			continue
		}

		defineMethodsWithinNamespace(buffer, subtypeNamespace, env, false)
	}
}

func defineMethods(buffer *bytes.Buffer, namespace types.Namespace) {
	buffer.WriteString("\n// Define methods\n")

	for methodName, method := range types.SortedOwnMethods(namespace) {
		fmt.Fprintf(
			buffer,
			"namespace.DefineMethod(%q, 0",
			method.DocComment,
		)

		if method.IsAbstract() {
			buffer.WriteString("| METHOD_ABSTRACT_FLAG")
		}
		if method.IsSealed() {
			buffer.WriteString("| METHOD_SEALED_FLAG")
		}
		if method.IsNative() {
			buffer.WriteString("| METHOD_NATIVE_FLAG")
		}
		if method.IsGenerator() {
			buffer.WriteString("| METHOD_GENERATOR_FLAG")
		}
		if method.IsAsync() {
			buffer.WriteString("| METHOD_ASYNC_FLAG")
		}

		fmt.Fprintf(
			buffer,
			", value.ToSymbol(%q), ",
			methodName.String(),
		)

		if len(method.TypeParameters) > 0 {
			buffer.WriteString("[]*TypeParameter{")
			for _, param := range method.TypeParameters {
				fmt.Fprintf(
					buffer,
					"%s,",
					typeToCode(param, true),
				)
			}
			buffer.WriteString("}, ")
		} else {
			buffer.WriteString("nil, ")
		}

		if len(method.Params) > 0 {
			buffer.WriteString("[]*Parameter{")
			for _, param := range method.Params {
				fmt.Fprintf(
					buffer,
					"NewParameter(value.ToSymbol(%q), %s, %s, %t),",
					param.Name.String(),
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
	}
}

func defineSubtypesWithinNamespace(buffer *bytes.Buffer, namespace types.Namespace) {
	for name, subtype := range types.SortedSubtypes(namespace) {
		if subtype.Type == namespace {
			continue
		}
		switch s := subtype.Type.(type) {
		case *types.Class:
			defineClass(buffer, s, name.String())
		case *types.Mixin:
			defineMixin(buffer, s, name.String())
		case *types.Module:
			defineModule(buffer, s, name.String())
		case *types.Interface:
			defineInterface(buffer, s, name.String())
		case *types.TypeParameter:
		default:
			defineSubtype(buffer, s, name.String())
		}
	}
}

func defineSubtype(buffer *bytes.Buffer, subtype types.Type, name string) {
	fmt.Fprintf(
		buffer,
		"namespace.DefineSubtype(value.ToSymbol(%q), %s)\n",
		name,
		typeToCode(subtype, true),
	)
}

func defineClass(buffer *bytes.Buffer, class *types.Class, constantName string) {
	hasSubtypes := len(class.Subtypes()) > 0
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
		`namespace.TryDefineClass(%q, %t, %t, %t, %t, value.ToSymbol(%q), %s, env)
		`,
		class.DocComment(),
		class.IsAbstract(),
		class.IsSealed(),
		class.IsPrimitive(),
		class.IsNoInit(),
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
	hasSubtypes := len(mixin.Subtypes()) > 0
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
	hasSubtypes := len(module.Subtypes()) > 0
	if hasSubtypes {
		buffer.WriteString(`{ namespace :=`)
	}

	fmt.Fprintf(
		buffer,
		`namespace.TryDefineModule(%q, value.ToSymbol(%q), env)
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
	hasSubtypes := len(iface.Subtypes()) > 0
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
func namespaceToCode(typ types.Namespace) string {
	switch t := typ.(type) {
	case nil:
		return "nil"
	case *types.Class:
		return fmt.Sprintf(
			"NameToType(%q, env).(*Class)",
			t.Name(),
		)
	case *types.SingletonClass:
		return fmt.Sprintf(
			"NameToNamespace(%q, env).Singleton()",
			t.AttachedObject.Name(),
		)
	case *types.Mixin:
		return fmt.Sprintf(
			"NameToType(%q, env).(*Mixin)",
			t.Name(),
		)
	case *types.MixinProxy:
		return fmt.Sprintf(
			"NameToType(%q, env).(*Mixin)",
			t.Name(),
		)
	case *types.Module:
		return fmt.Sprintf(
			"NameToType(%q, env).(*Module)",
			t.Name(),
		)
	case *types.Interface:
		return fmt.Sprintf(
			"NameToType(%q, env).(*Interface)",
			t.Name(),
		)
	case *types.InterfaceProxy:
		return fmt.Sprintf(
			"NameToType(%q, env).(*Interface)",
			t.Name(),
		)
	case *types.Generic:
		return typeToCode(typ, false)
	case *types.TypeParamNamespace:
		return fmt.Sprintf("NewTypeParamNamespace(%q, %t)", t.DocComment(), t.ForMethod)
	case *types.Closure:
		return typeToCode(typ, false)
	default:
		panic(
			fmt.Sprintf("invalid type: %T", typ),
		)
	}
}

// Serialise the type to Go code
func typeToCode(typ types.Type, init bool) string {
	switch t := typ.(type) {
	case nil:
		return "nil"
	case types.Any:
		return "Any{}"
	case types.Self:
		return "Self{}"
	case types.Void:
		return "Void{}"
	case types.Nil:
		return "Nil{}"
	case types.Bool:
		return "Bool{}"
	case types.True:
		return "True{}"
	case types.False:
		return "False{}"
	case types.Never:
		return "Never{}"
	case types.Untyped:
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
	case *types.TypeParameter:
		namespaceName := t.Namespace.Name()
		if init || len(namespaceName) == 0 {
			return fmt.Sprintf(
				"NewTypeParameter(value.ToSymbol(%q), %s, %s, %s, %s, %s)",
				t.Name.String(),
				namespaceToCode(t.Namespace),
				typeToCode(t.LowerBound, false),
				typeToCode(t.UpperBound, false),
				typeToCode(t.Default, false),
				t.Variance.String(),
			)
		}

		return fmt.Sprintf(
			`NameToType("%s::%s", env)`,
			namespaceName,
			t.Name.String(),
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
	case *types.SingletonOf:
		return fmt.Sprintf("NewSingletonOf(%s)", typeToCode(t.Type, init))
	case *types.InstanceOf:
		return fmt.Sprintf("NewInstanceOf(%s)", typeToCode(t.Type, init))
	case *types.Generic:
		buff := new(strings.Builder)
		fmt.Fprintf(
			buff,
			"NewGeneric(%s, NewTypeArguments(TypeArgumentMap{",
			namespaceToCode(t.Namespace),
		)
		for name, arg := range t.TypeArguments.AllArguments() {
			fmt.Fprintf(
				buff,
				"value.ToSymbol(%q): NewTypeArgument(%s, %s),",
				name.String(),
				typeToCode(arg.Type, init),
				arg.Variance.String(),
			)
		}
		buff.WriteString("}, []value.Symbol{")

		for _, name := range t.TypeArguments.ArgumentOrder {
			fmt.Fprintf(
				buff,
				"value.ToSymbol(%q),",
				name.String(),
			)
		}
		buff.WriteString("}))")

		return buff.String()
	case *types.TypeParamNamespace:
		return fmt.Sprintf("NewTypeParamNamespace(%q, %t)", t.DocComment(), t.ForMethod)
	case *types.Closure:
		buff := new(strings.Builder)
		fmt.Fprintf(
			buff,
			"NewClosureWithMethod(%q, 0",
			t.Body.DocComment,
		)

		if t.Body.IsAbstract() {
			buff.WriteString("| METHOD_ABSTRACT_FLAG")
		}
		if t.Body.IsSealed() {
			buff.WriteString("| METHOD_SEALED_FLAG")
		}
		if t.Body.IsNative() {
			buff.WriteString("| METHOD_NATIVE_FLAG")
		}
		if t.Body.IsGenerator() {
			buff.WriteString("| METHOD_GENERATOR_FLAG")
		}
		if t.Body.IsAsync() {
			buff.WriteString("| METHOD_ASYNC_FLAG")
		}

		fmt.Fprintf(
			buff,
			", value.ToSymbol(%q), ",
			t.Body.Name.String(),
		)

		if len(t.Body.TypeParameters) > 0 {
			buff.WriteString("[]*TypeParameter{")
			for _, param := range t.Body.TypeParameters {
				fmt.Fprintf(
					buff,
					"%s.(*TypeParameter),",
					typeToCode(param, false),
				)
			}
			buff.WriteString("}, ")
		} else {
			buff.WriteString("nil, ")
		}

		if len(t.Body.Params) > 0 {
			buff.WriteString("[]*Parameter{")
			for _, param := range t.Body.Params {
				fmt.Fprintf(
					buff,
					"NewParameter(value.ToSymbol(%q), %s, %s, %t),",
					param.Name.String(),
					typeToCode(param.Type, false),
					param.Kind,
					param.InstanceVariable,
				)
			}
			buff.WriteString("}, ")
		} else {
			buff.WriteString("nil, ")
		}

		fmt.Fprintf(
			buff,
			"%s, %s)",
			typeToCode(t.Body.ReturnType, false),
			typeToCode(t.Body.ThrowType, false),
		)
		return buff.String()
	default:
		panic(
			fmt.Sprintf("invalid type: %T", typ),
		)
	}
}
