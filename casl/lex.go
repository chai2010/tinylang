// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

// 解析全部记号
func LexAll(caslCode string) ([]Item, error) {
	return newLexer(caslCode).lexAll()
}

// 词法分析器
type lexer struct {
	r *txtReader
}

// 构造词法解析器
func newLexer(code string) *lexer {
	return &lexer{r: newTxtReader(code)}
}

// 解析全部
func (l *lexer) lexAll() (tokens []Item, err error) {
Loop:
	for {
		// 跳过开头空白
		l.skipSpace()

		// 解析下个记号
		switch r := l.r.peek(); true {
		case r == eof: // 文件结束
			break Loop
		case r == '\n' || r == '\r': // 一行结束
			continue Loop

		case r == ';' || r == '#': // 行注释, # 是扩展语法
			tok, err := l.lexComment()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, tok)

		case r == '+' || r == '-' || (r >= '0' && r <= '9'): // 数字
			tok, err := l.lexNumber()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, tok)

		case r == '\'': // 字符串
			tok, err := l.lexString()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, tok)

		case l.isAlphaNumer(r): // 标识符 或 关键字
			tok, err := l.lexIdent()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, tok)

		case r == ',': // 逗号
			tok := Item{
				Typ: COMMA,
				Val: ",",
				Pos: l.r.pos,
				End: l.r.pos + 1,
			}
			tokens = append(tokens, tok)

		default: // 错误
			return nil, fmt.Errorf("未知记号: %q, at %d", r, l.r.pos)
		}
	}

	return tokens, nil
}

// 跳过空白
func (l *lexer) skipSpace() {
	for {
		switch r := l.r.peek(); true {
		case l.isEneOfLine(r) || l.isEOF(r):
			l.r.next()
		case !l.isSpace(r):
			l.r.next()
		default:
			return
		}
	}
}

// 解析注释
func (l *lexer) lexComment() (tok Item, err error) {
	tok.Typ = COMMENT
	tok.Pos = l.r.pos
	for {
		switch r := l.r.peek(); true {
		case l.isEneOfLine(r) || l.isEOF(r):
			tok.Val = l.r.txt[tok.Pos:l.r.pos]
			tok.End = l.r.pos
			return
		default:
			l.r.next()
		}
	}
}

// 解析数字
func (l *lexer) lexNumber() (Item, error) {
	// 1. 跳过符号
	// 2. 支持十六进制
	panic("todo")
}

// 解析字符串
func (l *lexer) lexString() (Item, error) {
	// 1. 支持转义字符
	panic("todo")
}

// 解析标识符
func (l *lexer) lexIdent() (Item, error) {
	panic("todo")
}

// 空白字符(不含换行符号)
func (l *lexer) isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// 行尾符号
func (l *lexer) isEneOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

// 文件结束
func (l *lexer) isEOF(r rune) bool {
	return r == eof
}

// 是否为字面或数字(包含下划线, 不支持中文字符)
func (l *lexer) isAlphaNumer(r rune) bool {
	if r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' || r <= 'Z') {
		return true
	}
	if r >= '0' && r <= '9' {
		return true
	}
	return false
}
