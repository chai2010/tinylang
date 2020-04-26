// Copyright 2020 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import (
	"bytes"
	goast "go/ast"
	gotoken "go/token"
	"io"
	"os"
)

// Print 打印语法树到 stdout
func Print(node Node) {
	Fprint(os.Stdout, node)
}

// Fprint 打印语法树到指定目标
func Fprint(w io.Writer, node Node) {
	fset := gotoken.NewFileSet()

	if f, ok := node.(*File); ok {
		file := *f
		if len(file.Data) > 0 {
			fset.AddFile(f.Name, 1, len(f.Data)).SetLinesForContent(f.Data)
			file.Data = nil
		}
		node = &file
	}

	goast.Fprint(w, fset, node, goast.NotNilFilter)
}

func (p *File) String() string {
	var buf bytes.Buffer
	Fprint(&buf, p)
	return buf.String()
}
