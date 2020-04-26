// Copyright 2020 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

import (
	"strconv"
)

func (tok Token) Name() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens_name)) {
		s = tokens_name[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

var tokens_name = [...]string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	NUMBER: "NUMBER",

	ASSIGN: "ASSIGN",
	EQ:     "EQ",
	LT:     "LT",
	PLUS:   "PLUS",
	MINUS:  "MINUS",
	TIMES:  "TIMES",
	OVER:   "OVER",
	LPAREN: "LPAREN",
	RPAREN: "RPAREN",
	SEMI:   "SEMI",

	IF:     "IF",
	THEN:   "THEN",
	ELSE:   "ELSE",
	END:    "END",
	REPEAT: "REPEAT",
	UNTIL:  "UNTIL",
	READ:   "READ",
	WRITE:  "WRITE",
}
