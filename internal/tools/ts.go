package tools

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"reflect"
	"strings"
)

type TsField struct {
	name      string
	fieldType string
	nullable  bool
}

type TsType struct {
	name   string
	fields []TsField
}

type Parser struct {
	Files    []string
	Prefix   string
	TsPrefix string
	TsFile   string
}

func NewGoToTs(files []string, prefix, tsPrefix, tsFile string) *Parser {
	return &Parser{
		Files:    files,
		Prefix:   prefix,
		TsPrefix: tsPrefix,
		TsFile:   tsFile,
	}
}

func (p *Parser) ToTs() {
	fset := token.NewFileSet()

	tsDefs := ""

	for _, file := range p.Files {
		ts, err := p.parseFile(fset, file)
		if err != nil {
			log.Fatalf("error parsing file %s: %s", file, err)
		}

		for _, tsType := range ts {
			tsDefs += fmt.Sprintf("%s\n", generateTsStruct(tsType))
		}
	}

	f, err := os.Create(p.TsFile)
	if err != nil {
		log.Fatalln("error creating ts file: ", err)
	}
	_, err = f.WriteString(tsDefs)
	if err != nil {
		log.Fatalln("error writing generated ts: ", err)
		f.Close()
	}
}

func (p *Parser) parseFile(fset *token.FileSet, name string) ([]TsType, error) {
	f, err := parser.ParseFile(fset, name, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var types []TsType

	tsTypes := getTsType(f, p.Prefix, p.TsPrefix)

	for _, ty := range tsTypes {
		if ty != nil {
			types = append(types, *ty)
		}
	}

	return types, nil
}

func getTsType(f *ast.File, structPrefix, tsPrefix string) []*TsType {
	var results []*TsType

	ast.Inspect(f, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		structName := typeSpec.Name.Name

		if !strings.HasPrefix(structName, structPrefix) {
			return true
		}

		trimmedName := strings.TrimPrefix(structName, structPrefix)

		tsType := TsType{
			name: fmt.Sprintf("%s%s", tsPrefix, trimmedName),
		}

		for _, field := range structType.Fields.List {
			tsField := TsField{}

			if len(field.Names) > 0 {
				tsField.name = field.Names[0].Name
			}

			tsField.fieldType = goToTSType(exprToString(field.Type))

			if field.Tag != nil {
				tag, nullable := parseJSONTag(field.Tag.Value)

				if tag == "" {
					log.Printf("field %s in struct %s does not have json struct field\n", tsField.name, tsType.name)
					return true
				}
				tsField.name = tag
				tsField.nullable = nullable
			} else {
				log.Printf("field %s in struct %s does not have struct field\n", tsField.name, tsType.name)
				return true
			}

			tsType.fields = append(tsType.fields, tsField)
		}

		results = append(results, &tsType)
		return false
	})

	return results
}

func generateTsStruct(tsType TsType) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "export type %s = {\n", tsType.name)

	for _, field := range tsType.fields {
		fieldName := field.name
		fieldType := field.fieldType

		if field.nullable {
			fieldType += " | null"
		}

		fmt.Fprintf(&sb, "  %s: %s;\n", fieldName, fieldType)
	}

	sb.WriteString("};\n")

	return sb.String()
}

func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", exprToString(t.X), t.Sel.Name)
	case *ast.StarExpr:
		return "*" + exprToString(t.X)
	case *ast.ArrayType:
		return "[]" + exprToString(t.Elt)
	default:
		return fmt.Sprintf("%T", expr)
	}
}

func goToTSType(goType string) string {
	goType = strings.Replace(goType, "*", "", -1)

	if strings.HasPrefix(goType, "int") || strings.HasPrefix(goType, "float") {
		return "number"
	}

	if strings.HasPrefix(goType, "[]int") || strings.HasPrefix(goType, "[]float") {
		return "number[]"
	}
	if strings.HasPrefix(goType, "[]string") {
		return "string[]"
	}

	if goType == "string" {
		return "string"
	}
	if goType == "bool" {
		return "boolean"
	}
	if goType == "[]byte" {
		return "Uint8Array"
	}
	if goType == "map[string]interface{}" {
		return "{ [key: string]: any }"
	}
	if goType == "interface{}" {
		return "any"
	}

	return "unknown"
}

func parseJSONTag(tag string) (string, bool) {
	tag = strings.Replace(tag, "`", "", -1)
	tag = strings.Replace(tag, "`", "", -1)

	structTag := reflect.StructTag(tag)

	jsonTag := structTag.Get("json")

	if jsonTag == "" {
		return "", false
	}

	parts := strings.Split(jsonTag, ",")

	omitempty := false
	for _, part := range parts[1:] {
		if part == "omitempty" {
			omitempty = true
			break
		}
	}

	return parts[0], omitempty
}
