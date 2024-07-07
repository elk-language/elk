package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/types/checker"
	"github.com/elk-language/elk/value"
)

// Compiles Elk headers

const pathToHeaders = "./headers/"

func main() {
	env := types.NewGlobalEnvironment()
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
			package headers

			// This file is auto-generated, please do not edit it manually

			import (
				"github.com/elk-language/elk/types"
				"github.com/elk-language/elk/value/symbol"
			)

			func SetupGlobalEnvironment(env *types.GlobalEnvironment) {
				objectClass := env.StdSubtypeClass(symbol.Object)
				namespace := env.Root
		`,
	)

	defineSubtypesWithinNamespace(buffer, env.Root)

	buffer.WriteString(
		`
			}
		`,
	)

	os.WriteFile("headers/headers.go", buffer.Bytes(), 0666)
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
