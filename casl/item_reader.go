// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type itemReader struct {
	toks []Item // 文本
	pos  int    // 当前位置
}

func newItemReader(tokens []Item) *itemReader {
	if len(tokens) == 0 {
		tokens = []Item{{Typ: EOF}}
	}
	if tok := tokens[len(tokens)-1]; tok.Typ != EOF {
		tokens = append(tokens, Item{Typ: EOF})
	}
	return &itemReader{
		toks: tokens,
	}
}

func (p *itemReader) atEOF() bool {
	return p.pos >= len(p.toks)
}

func (p *itemReader) nextLine() (toks []Item) {
	for {
		// 读取一个记号
		tok := p.next()

		// 如果是行尾或文件结束, 则返回
		if tok.Typ == EOL || tok.Typ == EOF {
			return
		}

		// 记录记号到行
		toks = append(toks, tok)
	}
}

func (p *itemReader) next() Item {
	if p.pos >= len(p.toks) {
		return p.toks[len(p.toks)-1]
	}

	r := p.toks[p.pos]
	p.pos++

	return r
}
func (p *itemReader) peek() Item {
	if p.pos >= len(p.toks) {
		return p.toks[len(p.toks)-1]
	}

	r := p.toks[p.pos]
	return r
}

// 只能在next之后后悔一次,
func (p *itemReader) backup() {
	p.pos--
}
