// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strings"
	"unicode/utf8"
)

// 文件结尾rune值
const eof = -1

// 读文本
type txtReader struct {
	txt   string // 文本
	pos   int    // 当前位置
	width int    // 刚读rune的宽度
}

// 构造
func newTxtReader(txt string) *txtReader {
	return &txtReader{txt: txt}
}

// 读一个字符
func (p *txtReader) next() rune {
	if p.pos >= len(p.txt) {
		return eof
	}

	r, w := utf8.DecodeRuneInString(p.txt[p.pos:])
	p.width = w
	p.pos += w
	return r
}

// 查看下一个字符
func (p *txtReader) peek() rune {
	if p.pos >= len(p.txt) {
		return eof
	}

	r, w := utf8.DecodeRuneInString(p.txt[p.pos:])
	p.width = w
	return r
}

// 只能在next之后后悔一次,
func (p *txtReader) backup() {
	p.pos -= p.width
}

// 计算pos对应的行列位置(行列号从1开始)
func (p *txtReader) position(pos int) (line, column int) {
	if pos < 0 {
		return 0, 0
	}
	if pos > len(p.txt) {
		pos = len(p.txt)
	}
	line = strings.Count(p.txt[:pos], "\n") + 1
	if i := strings.LastIndexByte(p.txt[:pos], '\n'); i >= 0 {
		column = pos - i
	} else {
		column = 1
	}
	return
}
