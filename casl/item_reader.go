// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type itemReader struct {
	toks []Item // 文本
	pos  int    // 当前位置
}

func newItemReader(tokens []Item) *itemReader {
	return &itemReader{
		toks: tokens,
	}
}

func (p *itemReader) next() Item {
	if p.pos >= len(p.toks) {
		return Item{}
	}

	r := p.toks[p.pos]
	p.pos++

	return r
}
func (p *itemReader) peek() Item {
	if p.pos >= len(p.toks) {
		return Item{}
	}

	r := p.toks[p.pos]
	return r
}

// 只能在next之后后悔一次,
func (p *itemReader) backup() {
	p.pos--
}
