// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

// 解析全部记号
func LexAll(caslCode string) ([]Item, error) {
	return newLexer(caslCode).lexAll()
}

// 词法分析器
type lexer struct {
	r *txtReader
}

func newLexer(code string) *lexer {
	return &lexer{r: newTxtReader(code)}
}

// 解析全部
func (l *lexer) lexAll() ([]Item, error) {
	var tokens []Item
Loop:
	for {
		// 跳过开头空白
		l.skipSpace()

		// 解析下个记号
		switch r := l.r.peek(); true {
		case r == eof: // 文件结束
			break Loop
		case r == '\n' || r == '\r': // 一行结束
		case r == ';' || r == '#': // 行注释, # 是扩展语法
		case r >= '0' && r <= '9': // 数字
		case r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'): // 标识符 或 关键字
		case r == ',': // 逗号
		default: // 错误
			return nil, fmt.Errorf("未知记号: %q, at %d", r, l.r.pos)
		}
	}

	return tokens, nil
}

// 跳过空白
func (l *lexer) skipSpace() {
	//
}

// 解析注释
func (l *lexer) skipComment() {
	//
}

// 解析数字
func (l *lexer) skipNumber() {
	//
}

// 解析字符串
func (l *lexer) skipString() {
	//
}
