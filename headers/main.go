package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/types/checker"
)

// Compiles Elk headers

func main() {
	env := types.NewGlobalEnvironment()
	items, _ := os.ReadDir(".")
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		if !strings.HasSuffix(item.Name(), ".elh") {
			continue
		}

		bytes, err := os.ReadFile(item.Name())
		if err != nil {
			panic(err)
		}
		source := string(bytes)
		_, errList := checker.CheckSource(item.Name(), source, env, true)
		if len(errList) > 0 {
			fmt.Println(errList.HumanStringWithSource(source, true))
			os.Exit(1)
		}
	}
	buffer := new(bytes.Buffer)
	buffer.WriteString(
		`
			package headers

			import (
				"github.com/elk-language/elk/types"
			)

			func SetupGlobalEnvironment(env *types.GlobalEnvironment) {

		`,
	)

	compileNamespace(buffer, env.Root, true)

	buffer.WriteString(
		`
			}
		`,
	)
}

func compileNamespace(buffer *bytes.Buffer, namespace types.Namespace, root bool) {
	// for name, subtype := range namespace.Subtypes().Map {
	// 	switch s := subtype.(type) {
	// 	case *types.Class:
	// 		compileClass(buffer, s)
	// 	}
	// }
}

func compileClass(buffer *bytes.Buffer, class *types.Class) {

}
