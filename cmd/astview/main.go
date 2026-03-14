package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	// Создаем набор файлов
	fset := token.NewFileSet()

	// Парсим файл
	node, err := parser.ParseFile(fset, "testdata/valid/main.go", nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	// Обходим AST
	ast.Inspect(node, func(n ast.Node) bool {
		// Если это вызов функции
		if call, ok := n.(*ast.CallExpr); ok {
			fmt.Printf("Найден вызов функции: %T\n", call.Fun)

			// Смотрим на функцию
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if pkg, ok := sel.X.(*ast.Ident); ok {
					fmt.Printf("  Пакет: %s\n", pkg.Name)
					fmt.Printf("  Функция: %s\n", sel.Sel.Name)
				}
			}

			// Смотрим на аргументы
			for i, arg := range call.Args {
				fmt.Printf("  Аргумент %d: %T\n", i, arg)
				if lit, ok := arg.(*ast.BasicLit); ok {
					fmt.Printf("    Значение: %s (тип: %s)\n", lit.Value, lit.Kind)
				}
			}
		}
		return true
	})
}
