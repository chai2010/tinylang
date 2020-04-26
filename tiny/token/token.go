// Copyright 2020 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TINY 记号类型
package token

import (
	"strconv"
	"unicode"
)

// TINY 语言的词法记号
type Token int

// 词法记号列表
const (
	// Special tokens
	EOF Token = iota // EOF 必须返回 0
	ILLEGAL
	COMMENT

	literal_beg
	IDENT  // main
	NUMBER // 123
	literal_end

	operator_beg
	ASSIGN // :=
	EQ     // =
	LT     // <
	PLUS   // +
	MINUS  // -
	TIMES  // *
	OVER   // /
	LPAREN // (
	RPAREN // )
	SEMI   // ;
	operator_end

	keyword_beg
	IF
	THEN
	ELSE
	END
	REPEAT
	UNTIL
	READ
	WRITE
	keyword_end
)

var tokens = [...]string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	NUMBER: "NUMBER",

	ASSIGN: ":=",
	EQ:     "=",
	LT:     "<",
	PLUS:   "+",
	MINUS:  "-",
	TIMES:  "*",
	OVER:   "/",
	LPAREN: "(",
	RPAREN: ")",
	SEMI:   ";",

	IF:     "if",
	THEN:   "then",
	ELSE:   "else",
	END:    "end",
	REPEAT: "repeat",
	UNTIL:  "until",
	READ:   "read",
	WRITE:  "write",
}

func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

// Lookup 查找 ident 字符串对应的关键字类型, 或者返回 IDENT 类型.
func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}

// LookupOperator 查找字符串对应的运算符.
func LookupOperator(s rune) (tok Token, suffix rune, ok bool) {
	if s == ':' {
		return ASSIGN, '=', true
	}
	for i := operator_beg + 1; i < operator_end; i++ {
		if tokens[i] != "" && tokens[i] == string(s) {
			return Token(i), 0, true
		}
	}

	return ILLEGAL, 0, false
}

// IsValid 判断 tok 是否有效.
func (tok Token) IsValid() bool {
	if 0 <= tok && tok < Token(len(tokens)) {
		return tokens[tok] != ""
	}
	return false
}

// IsOperator 判断记号是否为运算符.
func (tok Token) IsOperator() bool {
	return operator_beg < tok && tok < operator_end
}

// IsKeyword 判断记号是否为关键字.
func (tok Token) IsKeyword() bool {
	return keyword_beg < tok && tok < keyword_end
}

// IsKeyword 判断 name 字符串是否为关键字.
func IsKeyword(name string) bool {
	_, ok := keywords[name]
	return ok
}

// IsIdentifier 判断 name 字符串是否为标识符号.
func IsIdentifier(name string) bool {
	for i, c := range name {
		if !unicode.IsLetter(c) && c != '_' && (i == 0 || !unicode.IsDigit(c)) {
			return false
		}
	}
	return name != "" && !IsKeyword(name)
}
