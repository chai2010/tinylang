// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/chai2010/tinylang/comet"
)

// TODO: 返回指令列表, 然后再生成指令

// 解析CASL程序, 返回机器码和对应的可读格式指令
func ParseCASL(caslCode string) (prog []comet.Instruction, err error) {
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

func (p *parser) paseAll() (prog []comet.Instruction, err error) {
	// CASL字符串解析为记号列表
	toks, err := LexAll(p.caslCode)
	if err != nil {
		return nil, err
	}

	// 行记号读接口
	r := newItemReader(toks)

	// 依次处理每行的记号
Loop:
	for !r.atEOF() {
		// 读取一行
		toks := r.nextLine()

		// 跳过空行(行尾记号已经被丢弃)
		if len(toks) == 0 {
			continue Loop
		}

		// 要解析的指令
		var ins comet.Instruction

		// 解析标号
		if len(toks) > 0 {
			if tok := toks[0]; tok.Typ == ID {
				ins.Label = tok.Val
				toks = toks[1:]
			}
		}

		// 解析指令
		var op = Lookup(toks[0].Val)

		// 解析宏指令
		if op.IsMACRO() {
			// TODO
			continue
		}

		// 解析系统指令
		if op.IsMACRO_SYSCALL() {
			// TODO
			continue
		}

		// 解析COMET指令
		if op.IsCOMET_INS() {
			// TODO
			continue
		}

		// 其它未知指令
		if tok := toks[0]; tok.Typ.IsKeyword() {
			//ins.Op = Lookup(tok.Val)
			toks = toks[1:]
		} else {
			err = fmt.Errorf("非法记号: %v", tok)
			return
		}

		// 判断行内第一个记号

		switch {
		case toks[0].Typ == ID: // 标号

		default:

		}
	}

	panic("TODO")
}

// 解析一行
func (p *parser) paseLine() {
	//var toks []
	//
}
