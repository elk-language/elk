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

	compileConstantContainer(buffer, env.Root, true)
}

func compileConstantContainer(buffer *bytes.Buffer, constContainer types.ConstantContainer, root bool) {
	// for name, constant := range constContainer.Constants() {
	// 	switch constant
	// }
}
