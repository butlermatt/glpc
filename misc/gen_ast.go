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

	expressions := []string{
		"Assign   : Name *lexer.Token, Value Expr",
		"Binary   : Left Expr, Operator *lexer.Token, Right Expr",
		"Boolean  : Token *lexer.Token, Value bool",
		"Grouping : Expression Expr",
		"Index    : Left Expr, Operator *lexer.Token, Right Expr",
		"List     : Values []Expr",
		"Logical  : Left Expr, Operator *lexer.Token, Right Expr",
		"Number   : Token *lexer.Token, Float float64, Int int",
		"Null     : Token *lexer.Token, Value interface{}",
		"String   : Token *lexer.Token, Value string",
		"Unary    : Operator *lexer.Token, Right Expr",
		"Variable : Name *lexer.Token",
	}

	statements := []string{
		"Expression : Expression Expr",
		"Var        : Name *lexer.Token, Value Expr",
	}

	err := defineAst(outDir, expressions, statements)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v", err)
		os.Exit(1)
	}
}

func defineAst(outDir string, exprs []string, stmts []string) error {
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
	for _, definition := range exprs {
		parts := strings.Split(definition, " : ")
		name := strings.TrimSpace(parts[0])
		exprNames = append(exprNames, name)
		err = defineType(file, "Expr", name, strings.TrimSpace(parts[1]))
		if err != nil {
			return err
		}
	}

	err = defineVisitor(file, "Expr", exprNames)
	if err != nil {
		return err
	}

	var stmtNames []string
	for _, definition := range stmts {
		parts := strings.Split(definition, " : ")
		name := strings.TrimSpace(parts[0])
		stmtNames = append(stmtNames, name)
		err = defineType(file, "Stmt", name, strings.TrimSpace(parts[1]))
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
