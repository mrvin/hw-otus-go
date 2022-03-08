package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

type field struct {
	Name      string
	TypeField string
	Tags      map[string]string
}

type validStruct struct {
	Name   string
	Fields []field
}

type validStructs struct {
	PackageName string
	Items       []validStruct
}

var vStructs validStructs

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("go-validate: not enough argument")
	}

	fset := token.NewFileSet()
	astTree, err := parser.ParseFile(fset, os.Args[1], nil, 0)
	if err != nil {
		log.Fatalf("go-validate: %v", err)
	}

	vStructs.PackageName = os.Getenv("GOPACKAGE")
	ast.Inspect(astTree, InspectNode)

	fileName := outFileName(os.Args[1])

	if err := genOutFile(fileName, vStructs); err != nil {
		log.Fatalf("go-validate: %v", err)
	}
}

func outFileName(srcFileName string) string {
	const suffixFileName = "validation_generated.go"

	dir, fileName := filepath.Split(srcFileName)
	fileName = fileName[:len(fileName)-len(filepath.Ext(fileName))]

	return filepath.Join(dir, fmt.Sprintf("%s_%s", fileName, suffixFileName))
}

func InspectNode(n ast.Node) bool {
	specNode, ok := n.(*ast.TypeSpec)
	if ok {
		vStruct := validStruct{specNode.Name.Name, make([]field, 0)}
		structType, ok := specNode.Type.(*ast.StructType)
		if ok {
			for _, f := range structType.Fields.List {
				if tag := tagParse(f.Tag); tag != nil {
					// todo:f.Names[0].Name
					vField := field{f.Names[0].Name, types.ExprString(f.Type), tag}
					vStruct.Fields = append(vStruct.Fields, vField)
				}
			}
		}
		vStructs.Items = append(vStructs.Items, vStruct)
	}

	return true
}

func tagParse(t *ast.BasicLit) map[string]string {
	if t == nil {
		return nil
	}
	mTag := make(map[string]string)

	tag, err := strconv.Unquote(t.Value)
	if err != nil {
		log.Printf("Invalid field tags: %s", t.Value)
		return nil
	}

	slValid := strings.Split(tag, " ")
	for _, valid := range slValid {
		if strings.HasPrefix(valid, "validate:") {
			vTag := strings.TrimPrefix(valid, "validate:")
			//todo: vTag[1:len(vTag)-1
			sl := strings.Split(vTag[1:len(vTag)-1], "|")
			for _, str := range sl {
				m := strings.SplitN(str, ":", 2)
				if len(m) > 1 {
					mTag[m[0]] = m[1]
				}
			}
		}
	}

	return mTag
}

func receiverName(name string) string {
	return strings.ToLower(name[:1])
}

func genOutFile(fileName string, structs validStructs) error {
	var buf bytes.Buffer

	outFileTemp, err := template.New("outfiletemp").
		Funcs(template.FuncMap{"receiverName": receiverName}).
		Parse(templ)
	if err != nil {
		return fmt.Errorf("сan't create template: %w", err)
	}

	if err := outFileTemp.Execute(&buf, structs); err != nil {
		return fmt.Errorf("сan't execute template: %w", err)
	}

	out, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("сan't format: %w", err)
	}

	outFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("сan't create file: %w", err)
	}
	_, err = outFile.Write(out)

	if closeErr := outFile.Close(); err == nil {
		err = closeErr
	}

	return err
}
