package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <output directory>", os.Args[0])
		os.Exit(1)
	}

	outDir := os.Args[1]

	expressions := map[string]string{
		"Boolean": "Token *lexer.Token, Value bool",
		"List":    "Values []Expr",
		"Number":  "Token *lexer.Token, Float float64, Int int",
		"Null":    "Token *lexer.Token, Value interface{}",
		"String":  "Token *lexer.Token, Value string",
		"Unary":   "Operator *lexer.Token, Right Expr",
	}

	statements := map[string]string{
		"Expression": "Expression Expr",
	}

	err := defineAst(outDir, expressions, statements)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v", err)
		os.Exit(1)
	}
}

func defineAst(outDir string, exprs map[string]string, stmts map[string]string) error {
	file, err := os.Create(outDir + "/ast.go")
	if err != nil {
		return err
	}
	defer file.Close()

	checkWrite(file, "package %s\n", outDir)
	checkWrite(file, "import \"github.com/butlermatt/glpc/lexer\"\n")
	checkWrite(file,
		`// Expr is an AST expression which returns a value of type Object or an error.
type Expr interface {
	Accept(ExprVisitor) (Object, error)
}

// Stmt is an AST statement which returns no value but may produce an error.
type Stmt interface {
	Accept(StmtVisitor) error
}
`)

	var exprNames []string
	for name, definition := range exprs {
		exprNames = append(exprNames, name)
		err = defineType(file, "Expr", name, definition)
		if err != nil {
			return err
		}
	}

	err = defineVisitor(file, "Expr", exprNames)
	if err != nil {
		return err
	}

	var stmtNames []string
	for name, definition := range stmts {
		stmtNames = append(stmtNames, name)
		err = defineType(file, "Stmt", name, definition)
		if err != nil {
			return err
		}
	}

	return defineVisitor(file, "Stmt", stmtNames)
}

func defineVisitor(file *os.File, nodeType string, names []string) error {
	checkWrite(file, "\n// %sVisitor will visit %[1]s objects and must receive calls to their applicable methods.", nodeType)
	checkWrite(file, "type %sVisitor interface {", nodeType)

	lower := strings.ToLower(nodeType)
	for _, name := range names {
		if lower == "expr" {
			checkWrite(file, "\tVisit%s(%s *%[1]s) (Object, error)", name+nodeType, lower)
		} else {
			checkWrite(file, "\tVisit%s(%s *%[1]s) error", name+nodeType, lower)
		}
	}

	return checkWrite(file, "}")
}

func defineType(file *os.File, nodeType, name, definition string) error {
	checkWrite(file, "// %s is a %s of a %s", name+nodeType, nodeType, name)
	checkWrite(file, "type %s struct {", name+nodeType)
	fields := strings.Split(definition, ", ")
	for _, field := range fields {
		checkWrite(file, "\t%s", strings.TrimSpace(field))
	}

	checkWrite(file, "}\n")
	ptr := strings.ToLower(name[0:1])
	checkWrite(file, "// Accept calls the correct visit method on %sVisitor, passing a reference to itself as a value", nodeType)
	if nodeType == "Expr" {
		return checkWrite(file, "func (%s *%s) Accept(visitor %sVisitor) (Object, error) { return visitor.Visit%[2]s(%[1]s) }\n", ptr, name+nodeType, nodeType)
	}

	return checkWrite(file, "func (%s *%s) Accept(visitor %sVisitor) error { return visitor.Visit%[2]s(%[1]s) }\n", ptr, name+nodeType, nodeType)
}

var writeError error = nil

func checkWrite(f *os.File, str string, args ...interface{}) error {
	if writeError != nil {
		return writeError
	}

	_, writeError = fmt.Fprintf(f, str+"\n", args...)
	return writeError
}
