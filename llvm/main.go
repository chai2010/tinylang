// Copyright 2021 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/chai2010/tinylang/llvm/pkg/compiler"
	"github.com/chai2010/tinylang/tiny/parser"
)

var (
	flagOutput = flag.String("o", "", "set output file")
	flagAst    = flag.Bool("ast", false, "print ast and exit")
	flagLLIR   = flag.Bool("llir", false, "print llvm-ir and exit")
	flagDebug  = flag.Bool("debug", false, "run with debug mode")
)

//go:embed tiny-lib/tiny_lib.ll
var tiny_lib_ll string

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stdout, `usage: tiny-llvm [flags] file.tiny`)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		return
	}

	filename := flag.Arg(0)
	f, err := parser.ParseFile(filename, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *flagAst {
		fmt.Print(f)
		return
	}

	c := compiler.NewCompiler(&compiler.Options{
		ModulePath:  "tiny.ll",
		ModuleBytes: []byte(tiny_lib_ll),
	})
	module, err := c.Compile(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *flagOutput == "" {
		fmt.Println(module)
		return
	}

	outfile, err := os.Create(*flagOutput)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer outfile.Close()

	_, err = module.WriteTo(outfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
