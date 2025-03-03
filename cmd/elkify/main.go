package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
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
	name        string
	typeName    string
	typeIsSlice bool
	doc         string
}

type structMap map[string]*structDefinition

func main() {
	sourcePath := "parser/ast"

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

	err = os.MkdirAll("tmp", os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("could not create tmp directory: %s", err))
	}

	for _, structDef := range cache {
		generateStruct(structDef)
	}

}

func generateStruct(structDef *structDefinition) {
	generateHeaderForStruct(structDef)
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
		return "Int8"
	default:
		return goType
	}
}

func goTypeToElkType(goType string, isSlice bool) string {
	simpleElkType := simpleGoTypeToElkType(goType)
	if !isSlice {
		return simpleElkType
	}

	return fmt.Sprintf("Tuple[%s]", simpleElkType)
}

func generateHeaderForStruct(structDef *structDefinition) {
	buffer := new(bytes.Buffer)

	if structDef.doc != "" {
		fmt.Fprintf(
			buffer,
			`
##[
	%s]##`,
			structDef.doc,
		)
	}

	fmt.Fprintf(
		buffer,
		"\nsealed primitive class %s\n",
		structDef.name,
	)

	buffer.WriteString("  constructor(")
	for i, field := range structDef.fields {
		elkFieldName := strcase.ToSnake(field.name)
		elkType := goTypeToElkType(field.typeName, field.typeIsSlice)

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
		elkType := goTypeToElkType(field.typeName, field.typeIsSlice)

		if field.doc != "" {
			fmt.Fprintf(buffer, `
		##[
			%s
		]##
			`,
				field.doc,
			)
		}
		fmt.Fprintf(buffer, "  getter %s: %s\n",
			elkFieldName,
			elkType,
		)
	}

	buffer.WriteString("end\n")

	fileName := strcase.ToSnake(structDef.name)
	filePath := fmt.Sprintf("./tmp/%s.elh", fileName)
	err := os.WriteFile(filePath, buffer.Bytes(), os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("could not write elk header file: %s", err))
	}
}

func analyseStruct(cache structMap, typeName string, doc string, node *ast.StructType) {
	structDefinition := &structDefinition{
		name: typeName,
		doc:  doc,
	}
	cache[typeName] = structDefinition

	for _, fld := range node.Fields.List {
		var fieldTypeName string
		var fieldTypeIsSlice bool

		switch fld := fld.Type.(type) {
		case *ast.Ident:
			fieldTypeName = fld.Name
		case *ast.StarExpr:
			ident, ok := fld.X.(*ast.Ident)
			if !ok {
				continue
			}

			fieldTypeName = ident.Name
		case *ast.ArrayType:
			ident, ok := fld.Elt.(*ast.Ident)
			if !ok {
				continue
			}

			fieldTypeName = ident.Name
		default:
			continue
		}

		for _, fieldIdent := range fld.Names {
			if len(fieldIdent.Name) == 0 {
				continue
			}

			firstChar, _ := utf8.DecodeRuneInString(fieldIdent.Name)
			if !unicode.IsUpper(firstChar) {
				continue
			}

			structField := &field{
				typeName:    fieldTypeName,
				name:        fieldIdent.Name,
				typeIsSlice: fieldTypeIsSlice,
				doc:         fld.Doc.Text(),
			}
			structDefinition.fields = append(
				structDefinition.fields,
				structField,
			)
		}
	}
}
