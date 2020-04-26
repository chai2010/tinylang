// Copyright 2020 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chai2010/tinylang/comet"
	"github.com/chai2010/tinylang/tiny/compiler"
	"github.com/chai2010/tinylang/tiny/parser"
)

var (
	flagAst   = flag.Bool("ast", false, "print ast and exit")
	flagCASL  = flag.Bool("casl", false, "print casl assembly and exit")
	flagDebug = flag.Bool("debug", false, "run with debug mode")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stdout, `usage: tiny [flags] file.tiny`)
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

	c := compiler.NewCompiler()
	bytecode, err := c.Compile(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *flagCASL {
		fmt.Print(c.CASLString())
		return
	}

	prog, err := comet.LoadProgram(filename, bytecode)
	if err != nil {
		fmt.Println(err)
		return
	}

	vm := comet.NewComet(prog)
	if *flagDebug {
		vm.DebugRun()
	} else {
		vm.Run()
	}
}
