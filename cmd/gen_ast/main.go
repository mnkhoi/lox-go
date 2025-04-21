package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args[1:]) != 1 {
		fmt.Println("Usage: generate_ast <output directory>")
		os.Exit(64)
	}

	outputDir := os.Args[1]

	defineAst(outputDir, "Expr", []string{
		"Binary   : Left Expr, Operator Token, Right Expr",
		"Grouping : Expression Expr",
		"Literal  : Value any",
		"Unary    : Operator Token, Right Expr",
	})
}

func defineAst(outdir, basename string, types []string) error {
	path := fmt.Sprintf("%s/%s.go", outdir, strings.ToLower(basename))
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprintf(f, "package lox\n\n")
	fmt.Fprintf(f, "type %s interface {\n", basename)
	fmt.Fprintf(f, "\tAccept(visitor Visitor) any\n")
	fmt.Fprintf(f, "}\n")

	defineVisitor(f, basename, types)

	for _, ttype := range types {
		splitted := strings.Split(ttype, ":")
		className := strings.Trim(splitted[0], " ")
		fields := strings.Trim(splitted[1], " ")
		defineType(f, basename, className, fields)
	}

	return nil
}

func defineType(f *os.File, basename, classname, fieldlist string) error {
	fmt.Fprintf(f, "\ttype %s struct {\n", classname)
	for field := range strings.SplitSeq(fieldlist, ", ") {
		name := strings.Split(field, " ")[0]
		ttype := strings.Split(field, " ")[1]
		fmt.Fprintf(f, "\t\t%s %s\n", name, ttype)
	}
	fmt.Fprintf(f, "}\n")
	fmt.Fprintf(f, "func (t %s) Accept(visitor Visitor) any {\n", classname)
	fmt.Fprintf(f, "\treturn visitor.Visit%s%s(t)\n", basename, classname)
	fmt.Fprintf(f, "}\n")

	// fmt.Fprintf(f, "func New%s(", classname)
	// for field := range strings.SplitSeq(fieldlist, ", ") {
	// 	name := strings.ToLower(strings.Split(field, " ")[0])
	// 	ttype := strings.Split(field, " ")[1]
	// 	fmt.Fprintf(f, "%s %s, ", name, ttype)
	// }
	// fmt.Fprintf(f, "\treturn &%s{}")

	return nil
}

func defineVisitor(f *os.File, basename string, types []string) {
	fmt.Fprintf(f, "type Visitor interface {\n")

	for _, t := range types {
		classname := strings.Trim(strings.Split(t, ":")[0], " ")
		fmt.Fprintf(f, "\tVisit%s%s(%s %s) any\n", basename, classname, strings.ToLower(classname), classname)
	}
	fmt.Fprintf(f, "}\n")
}
