// Copyright 2020 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"os"
	"text/template"

	"github.com/chai2010/tinylang/tiny/token"
)

func main() {
	var toks []TokenInfo
	for i := 0; i < 100; i++ {
		if tok := token.Token(i); tok.IsValid() {
			switch tok {
			case token.EOF, token.ILLEGAL, token.COMMENT:
				// skip
			default:
				toks = append(toks, TokenInfo{
					Name:  tok.Name(),
					Value: tok.String(),
				})
			}
		}
	}

	tmpl, err := template.New("").Parse(tmpl_Code)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, toks)
	if err != nil {
		panic(err)
	}
}

type TokenInfo struct {
	Name  string
	Value string
}

const tmpl_Code = `
{{- $tokInfoList := . -}}

// Copyright 2020 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Auto generated. DO NOT ENIT!!!

package parser

import (
	"github.com/chai2010/tinylang/tiny/token"
)

func yyTok2tok(x int) token.Token {
	switch x {
	{{- range $i, $v := $tokInfoList}}
	case _{{$v.Name}}:
		return token.{{$v.Name}}
	{{- end}}
	}
	return token.ILLEGAL
}


func tok2yyTok(x token.Token) int {
	switch x {
	{{- range $i, $v := $tokInfoList}}
	case token.{{$v.Name}}:
		return _{{$v.Name}}
	{{- end}}
	}
	return 0
}
`
