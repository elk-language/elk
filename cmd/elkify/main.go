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
	isGetter      bool
	fieldType     *fieldType
}

func (f *field) elkName() string {
	return strcase.ToSnake(f.name)
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
	case "uint16":
		return "UInt16"
	case "uint32":
		return "UInt32"
	case "uint64":
		return "UInt64"
	case "int8":
		return "Int8"
	case "int16":
		return "Int16"
	case "int32":
		return "Int32"
	case "int64":
		return "Int64"
	case "float32":
		return "Float"
	case "float64":
		return "Float"
	case "rune":
		return "Char"
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

func simpleGoTypeToElkTypeConversion(value string, fldType *fieldType) string {
	switch fldType.name {
	case "bool":
		return fmt.Sprintf("value.ToElkBool(%s)", value)
	case "string":
		return fmt.Sprintf("value.Ref(value.String(%s))", value)
	case "int":
		return fmt.Sprintf("value.SmallInt(%s).ToValue()", value)
	case "uint8", "byte":
		return fmt.Sprintf("value.UInt8(%s).ToValue()", value)
	case "uint16":
		return fmt.Sprintf("value.UInt16(%s).ToValue()", value)
	case "uint32":
		return fmt.Sprintf("value.UInt32(%s).ToValue()", value)
	case "uint64":
		return fmt.Sprintf("value.UInt32(%s).ToValue()", value)
	case "int8":
		return fmt.Sprintf("value.Int8(%s).ToValue()", value)
	case "int16":
		return fmt.Sprintf("value.Int16(%s).ToValue()", value)
	case "int32":
		return fmt.Sprintf("value.Int32(%s).ToValue()", value)
	case "int64":
		return fmt.Sprintf("value.Int32(%s).ToValue()", value)
	case "float32":
		return fmt.Sprintf("value.Float(%s).ToValue()", value)
	case "float64":
		return fmt.Sprintf("value.Float(%s).ToValue()", value)
	case "rune":
		return fmt.Sprintf("value.Char(%s).ToValue()", value)
	default:
		if fldType.pkg == "position" {
			switch fldType.name {
			case "Span", "Position":
				return fmt.Sprintf(
					"value.Ref((*value.%s)(%s))",
					fldType.name,
					value,
				)
			}
		}
		return fmt.Sprintf("value.Ref(%s)", value)
	}
}

func generateGetterConversionToElkType(buffer *bytes.Buffer, value string, fldType *fieldType, isSlice bool) {
	if !isSlice {
		typ := simpleGoTypeToElkTypeConversion(value, fldType)
		fmt.Fprintf(buffer, "result := %s\n", typ)
		return
	}

	fmt.Fprintf(
		buffer,
		`
			collection := %[1]s
			arrayTuple := value.NewArrayTupleWithLength(len(collection))
			for _, el := range collection {
				arrayTuple.Append(%s)
			}
			result := value.Ref(arrayTuple)
		`,
		value,
		simpleGoTypeToElkTypeConversion("el", fldType),
	)
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
		elkFieldName := field.elkName()
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
	generateInstanceMethodsForStruct(buffer, structDef)

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

	for i, fld := range structDef.fields {
		if i == len(structDef.fields)-1 {
			continue
		}

		if fld.fieldType.isSlice {
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
				valueTypeName(fld.fieldType),
				elkTypeToGoTypeAssertion("el", fld.fieldType),
			)
		} else {
			fmt.Fprintf(
				buffer,
				"arg%[1]d := %s\n",
				i,
				elkTypeToGoTypeAssertion(fmt.Sprintf("args[%d]", i), fld.fieldType),
			)
		}
	}

	fmt.Fprintf(
		buffer,
		`
			var argSpan *position.Span
			if args[%[1]d].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[%[1]d].Pointer())
			}
		`,
		len(structDef.fields)-1,
	)

	fmt.Fprintf(buffer, "self := ast.New%s(\nargSpan,\n", structDef.name)
	for i := range len(structDef.fields) - 1 {
		fmt.Fprintf(buffer, "arg%[1]d,\n", i)
	}
	buffer.WriteString("\n)\n")
	buffer.WriteString("return value.Ref(self), value.Undefined\n")

	buffer.WriteString("\n},\n")
	fmt.Fprintf(buffer, "vm.DefWithParameters(%d),\n", len(structDef.fields))
	buffer.WriteString("\n)\n")
}

func generateInstanceMethodsForStruct(buffer *bytes.Buffer, structDef *structDefinition) {
	for _, field := range structDef.fields {
		elkFieldName := field.elkName()
		fmt.Fprintf(
			buffer,
			`
			vm.Def(
				c,
				"%s",
				func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			`,
			elkFieldName,
		)

		fmt.Fprintf(
			buffer,
			"self := args[0].MustReference().(*ast.%s)\n",
			structDef.name,
		)

		var getter string
		if field.isGetter {
			getter = fmt.Sprintf("self.%s()", field.name)
		} else {
			getter = fmt.Sprintf("self.%s", field.name)
		}
		generateGetterConversionToElkType(
			buffer,
			getter,
			field.fieldType,
			field.fieldType.isSlice,
		)
		fmt.Fprintf(buffer, "return result, value.Undefined\n")

		buffer.WriteString("\n},\n")
		buffer.WriteString("\n)\n")
	}
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
	case "uint16":
		return fmt.Sprintf("(uint16)(%s.AsUInt16())", val)
	case "uint32":
		return fmt.Sprintf("(uint32)(%s.AsUInt32())", val)
	case "uint64":
		return fmt.Sprintf("(uint64)(%s.AsUInt64())", val)
	case "int8":
		return fmt.Sprintf("(int8)(%s.AsInt8())", val)
	case "int16":
		return fmt.Sprintf("(int16)(%s.AsInt16())", val)
	case "int32":
		return fmt.Sprintf("(int32)(%s.AsInt32())", val)
	case "int64":
		return fmt.Sprintf("(int64)(%s.AsInt64())", val)
	case "float64":
		return fmt.Sprintf("(float64)(%s.AsFloat())", val)
	case "float32":
		return fmt.Sprintf("(float32)(%s.AsFloat())", val)
	case "rune":
		return fmt.Sprintf("(rune)(%s.AsChar())", val)
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

			if !fieldIdent.IsExported() {
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
			name:     "Span",
			isGetter: true,
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
