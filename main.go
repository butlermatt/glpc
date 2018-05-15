package main

import (
	"fmt"
	"github.com/butlermatt/glpc/interpreter"
	"github.com/butlermatt/glpc/lexer"
	"github.com/butlermatt/glpc/parser"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [script]", os.Args[0])
	}

	runFile(os.Args[1])
}

func runFile(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %+v", err)
		os.Exit(1)
	}

	err = run(data, path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

func run(input []byte, filename string) error {
	l := lexer.New(input, filename)
	p := parser.New(l)
	interp := interpreter.New()
	env, err := interp.Interpret(p, filename)
	if err != nil {
		return err
	}

	return interp.RunMain(env)
}
