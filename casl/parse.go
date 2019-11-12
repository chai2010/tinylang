// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// 解析CASL程序, 返回机器码和对应的可读格式指令
func ParseCASL(caslCode string) (prog []uint16, ins []string, err error) {
	return newParser(caslCode).paseAll()
}

// 语法解析器
type parser struct {
	caslCode string
	macroTok Token // ILLEGAL/START/END
}

// 构建新的语法解析器
func newParser(caslCode string) *parser {
	return &parser{
		caslCode: caslCode,
	}
}

func (p *parser) paseAll() (prog []uint16, ins []string, err error) {
	// CASL字符串解析为记号列表
	toks, err := LexAll(p.caslCode)
	if err != nil {
		return nil, nil, err
	}

	// 依次处理每个记号
	r := newItemReader(toks)
Loop:
	for {
		switch tok := r.peek(); true {
		case tok.Typ == EOL || tok.Typ == EOF: // 一行只有一个语句
			continue Loop
		case tok.Typ == ID: // 标号

		default:
			_ = tok
		}
	}

	panic("TODO")
}

// 解析一行
func (p *parser) paseLine() {
	//var toks []
	//
}
