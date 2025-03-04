package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/iancoleman/strcase"
)

// Parsers go source files looking for
// Go struct definitions and generates
// Go code that declared Elk constructors,
// getters and header files

type structDefinition struct {
	name   string
	fields []*field
	doc    string
}

type field struct {
	name          string
	doc           string
	inConstructor bool
	fieldType     *fieldType
}

type fieldType struct {
	name      string
	pkg       string
	isSlice   bool
	isPointer bool
}

type structMap map[string]*structDefinition

func main() {
	sourcePath := "parser/ast"
	module := "Std::Elk::AST"

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, sourcePath, nil, parser.ParseComments)
	if err != nil {
		panic(fmt.Sprintf("parsing error: %s", err))
	}

	cache := make(map[string]*structDefinition)

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				decl, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}

				structDoc := decl.Doc.Text()
				for _, spec := range decl.Specs {
					spec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					structType, ok := spec.Type.(*ast.StructType)
					if !ok {
						continue
					}

					analyseStruct(cache, spec.Name.Name, structDoc, structType)
				}
			}
		}
	}

	err = os.MkdirAll("tmp/headers", os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("could not create tmp directory: %s", err))
	}
	err = os.MkdirAll("tmp/methods", os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("could not create tmp directory: %s", err))
	}

	for _, structDef := range cache {
		generateStruct(structDef, module)
	}

	cmd := exec.Command("go", "fmt", "./tmp/...")
	cmd.Run()
}

func generateStruct(structDef *structDefinition, module string) {
	generateHeaderForStruct(structDef, module)
	generateMethodsForStruct(structDef, module)
}

func simpleGoTypeToElkType(goType string) string {
	switch goType {
	case "bool":
		return "Bool"
	case "string":
		return "String"
	case "int":
		return "Int"
	case "uint8", "byte":
		return "UInt8"
	default:
		return goType
	}
}

func goTypeToElkType(goType string, isSlice bool) string {
	simpleElkType := simpleGoTypeToElkType(goType)
	if !isSlice {
		return simpleElkType
	}

	return fmt.Sprintf("ArrayTuple[%s]", simpleElkType)
}

const indentUnit = "  "

func generateHeaderForStruct(structDef *structDefinition, module string) {
	buffer := new(bytes.Buffer)

	baseIndentLevel := 0
	if len(module) != 0 {
		fmt.Fprintf(buffer, "module %s\n", module)
		baseIndentLevel++
	}

	baseIndent := strings.Repeat(indentUnit, baseIndentLevel)
	indentOne := strings.Repeat(indentUnit, baseIndentLevel+1)

	if structDef.doc != "" {
		fmt.Fprintf(
			buffer,
			"%s##[\n%s%s%[1]s]##",
			baseIndent,
			indentOne,
			structDef.doc,
		)
	}

	fmt.Fprintf(
		buffer,
		"\n%ssealed primitive class %s\n",
		baseIndent,
		structDef.name,
	)

	fmt.Fprintf(buffer, "%sconstructor(", indentOne)
	for i, field := range structDef.fields {
		if !field.inConstructor {
			continue
		}
		elkFieldName := strcase.ToSnake(field.name)
		elkType := goTypeToElkType(field.fieldType.name, field.fieldType.isSlice)

		if i != 0 {
			buffer.WriteString(", ")
		}
		fmt.Fprintf(
			buffer,
			"%s: %s",
			elkFieldName,
			elkType,
		)
	}
	buffer.WriteString("); end\n")

	for _, field := range structDef.fields {
		elkFieldName := strcase.ToSnake(field.name)
		elkType := goTypeToElkType(field.fieldType.name, field.fieldType.isSlice)

		if field.doc != "" {
			fmt.Fprintf(buffer, `
%[1]s##[
%s%s
%[1]s]##
			`,
				indentOne,
				strings.Repeat(indentUnit, baseIndentLevel+3),
				field.doc,
			)
		}
		fmt.Fprintf(buffer, "%sdef %s: %s; end\n",
			indentOne,
			elkFieldName,
			elkType,
		)
	}

	fmt.Fprintf(buffer, "%send\n", baseIndent)

	if len(module) != 0 {
		buffer.WriteString("end\n")
	}

	fileName := strcase.ToSnake(structDef.name)
	filePath := fmt.Sprintf("./tmp/headers/%s.elh", fileName)
	err := os.WriteFile(filePath, buffer.Bytes(), os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("could not write elk header file: %s", err))
	}
}

func generateMethodsForStruct(structDef *structDefinition, module string) {
	buffer := new(bytes.Buffer)

	buffer.WriteString(
		`package ast

import (
	"github.com/elk-language/elk/vm"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
)
`,
	)

	fmt.Fprintf(buffer, "func init%s() {\n", structDef.name)
	fmt.Fprintf(buffer, "c := &value.%sClass.MethodContainer", structDef.name)

	generateConstructorForStruct(buffer, structDef)

	buffer.WriteString("\n}\n")

	fileName := strcase.ToSnake(structDef.name)
	filePath := fmt.Sprintf("./tmp/methods/%s.go", fileName)
	err := os.WriteFile(filePath, buffer.Bytes(), os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("could not write elk method file: %s", err))
	}
}

func generateConstructorForStruct(buffer *bytes.Buffer, structDef *structDefinition) {
	buffer.WriteString(
		`
		vm.Def(
			c,
			"#init",
			func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
		`,
	)

	i := 0
	for _, field := range structDef.fields {
		if !field.inConstructor {
			continue
		}

		if field.fieldType.isSlice {
			fmt.Fprintf(
				buffer,
				`
					arg%[1]dTuple := args[%[1]d].MustReference().(*value.ArrayTuple)
					arg%[1]d := make([]%[2]s, arg%[1]dTuple.Length())
					for _, el := range *arg%[1]dTuple {
						arg%[1]d = append(arg%[1]d, %[3]s)
					}
				`,
				i,
				valueTypeName(field.fieldType),
				elkTypeToGoTypeAssertion("el", field.fieldType),
			)
		} else {
			fmt.Fprintf(
				buffer,
				"arg%[1]d := %s\n",
				i,
				elkTypeToGoTypeAssertion(fmt.Sprintf("args[%d]", i), field.fieldType),
			)
		}

		i++
	}

	fmt.Fprintf(buffer, "self := ast.New%s(\nposition.DefaultSpan,\n", structDef.name)
	i = 0
	for _, field := range structDef.fields {
		if !field.inConstructor {
			continue
		}

		fmt.Fprintf(buffer, "arg%[1]d,\n", i)
		i++
	}
	buffer.WriteString("\n)\n")
	buffer.WriteString("return value.Ref(self), value.Undefined\n")

	buffer.WriteString("\n},\n")
	fmt.Fprintf(buffer, "vm.DefWithParameters(%d),\n", len(structDef.fields))
	buffer.WriteString("\n)\n")
}

func valueTypeName(fldType *fieldType) string {
	if fldType.isPointer {
		return fmt.Sprintf("*%s.%s", fldType.pkg, fldType.name)
	}
	return fmt.Sprintf("%s.%s", fldType.pkg, fldType.name)
}

func elkTypeToGoTypeAssertion(val string, fldType *fieldType) string {
	switch fldType.name {
	case "bool":
		return fmt.Sprintf("value.Truthy(%s)", val)
	case "string":
		return fmt.Sprintf("(string)(%s.MustReference().(value.String))", val)
	case "int":
		return fmt.Sprintf("(int)(%s.AsInt())", val)
	case "uint8", "byte":
		return fmt.Sprintf("(uint8)(%s.AsUInt8())", val)
	default:
		return fmt.Sprintf("%s.MustReference().(%s)", val, valueTypeName(fldType))
	}
}

func analyseStruct(cache structMap, typeName string, doc string, node *ast.StructType) {
	structDefinition := &structDefinition{
		name: typeName,
		doc:  doc,
	}
	cache[typeName] = structDefinition

	for _, fld := range node.Fields.List {

		fldType := getFieldType(fld.Type)
		for _, fieldIdent := range fld.Names {
			if len(fieldIdent.Name) == 0 {
				continue
			}

			firstChar, _ := utf8.DecodeRuneInString(fieldIdent.Name)
			if !unicode.IsUpper(firstChar) {
				continue
			}

			structField := &field{
				name:          fieldIdent.Name,
				fieldType:     fldType,
				doc:           fld.Doc.Text(),
				inConstructor: true,
			}
			structDefinition.fields = append(
				structDefinition.fields,
				structField,
			)
		}
	}

	structDefinition.fields = append(
		structDefinition.fields,
		&field{
			name: "Span",
			fieldType: &fieldType{
				name:      "Span",
				pkg:       "position",
				isPointer: true,
			},
		},
	)
}

func getFieldType(typeNode ast.Expr) *fieldType {
	switch typeNode := typeNode.(type) {
	case *ast.Ident:
		return &fieldType{
			name: typeNode.Name,
			pkg:  "ast",
		}
	case *ast.SelectorExpr:
		pkgNode, ok := typeNode.X.(*ast.Ident)
		if !ok {
			return nil
		}

		return &fieldType{
			name: typeNode.Sel.Name,
			pkg:  pkgNode.Name,
		}
	case *ast.StarExpr:
		typ := getFieldType(typeNode.X)
		if typ != nil {
			typ.isPointer = true
		}
		return typ
	case *ast.ArrayType:
		typ := getFieldType(typeNode.Elt)
		if typ != nil {
			typ.isSlice = true
		}
		return typ
	default:
		return nil
	}
}
