// Copyright 2021 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/chai2010/tinylang/tiny/parser"
	"github.com/chai2010/tinylang/wasm/pkg/compiler"
	"github.com/chai2010/tinylang/wasm/pkg/vm"
	"github.com/chai2010/tinylang/wasm/pkg/wasm/encoding"
)

var (
	flagAst    = flag.Bool("ast", false, "print ast and exit")
	flagDebug  = flag.Bool("debug", false, "run with debug mode")
	flagOutput = flag.String("o", "", "set output file")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stdout, `usage: tiny-wasm [flags] file.tiny`)
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
	m, err := c.Compile(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *flagDebug {
		fmt.Println(c.ModuleString())
	}

	var buf bytes.Buffer
	if err := encoding.WriteModule(&buf, m); err != nil {
		panic(err)
	}

	if *flagOutput != "" {
		if err := os.WriteFile(*flagOutput, buf.Bytes(), 0666); err != nil {
			panic(err)
		}
	} else {
		if err := vm.Run(buf.Bytes()); err != nil {
			panic(err)
		}
	}
}
