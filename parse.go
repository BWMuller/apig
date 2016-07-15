package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func parseField(field *ast.Field) map[string]string {
	fields := make(map[string]string)
	fieldNames := []string{}

	for _, name := range field.Names {
		fieldNames = append(fieldNames, name.Name)
	}

	var fieldType string

	switch x := field.Type.(type) {
	case *ast.Ident: // e.g. string
		fieldType = x.Name
	case *ast.StarExpr: // e.g. *time.Time
		switch x2 := x.X.(type) {
		case *ast.SelectorExpr:
			switch x3 := x2.X.(type) {
			case *ast.Ident:
				fieldType = x3.Name + "." + x2.Sel.Name
			}
		}
	}

	for _, name := range fieldNames {
		fields[name] = fieldType
	}

	return fields
}

func parseFile(path string) ([]*Model, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)

	if err != nil {
		return nil, err
	}

	models := []*Model{}

	ast.Inspect(f, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.GenDecl:
			if x.Tok != token.TYPE {
				break
			}

			for _, spec := range x.Specs {
				fields := make(map[string]string)

				var modelName string

				switch x2 := spec.(type) {
				case *ast.TypeSpec:
					modelName = x2.Name.Name

					switch x3 := x2.Type.(type) {
					case *ast.StructType:
						for _, field := range x3.Fields.List {
							fs := parseField(field)

							for k, v := range fs {
								fields[k] = v
							}
						}
					}

					models = append(models, &Model{
						Name:   modelName,
						Fields: fields,
					})
				}
			}
		}

		return true
	})

	return models, nil
}
