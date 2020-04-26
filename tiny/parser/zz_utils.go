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
	case _IDENT:
		return token.IDENT
	case _NUMBER:
		return token.NUMBER
	case _ASSIGN:
		return token.ASSIGN
	case _EQ:
		return token.EQ
	case _LT:
		return token.LT
	case _PLUS:
		return token.PLUS
	case _MINUS:
		return token.MINUS
	case _TIMES:
		return token.TIMES
	case _OVER:
		return token.OVER
	case _LPAREN:
		return token.LPAREN
	case _RPAREN:
		return token.RPAREN
	case _SEMI:
		return token.SEMI
	case _IF:
		return token.IF
	case _THEN:
		return token.THEN
	case _ELSE:
		return token.ELSE
	case _END:
		return token.END
	case _REPEAT:
		return token.REPEAT
	case _UNTIL:
		return token.UNTIL
	case _READ:
		return token.READ
	case _WRITE:
		return token.WRITE
	}
	return token.ILLEGAL
}

func tok2yyTok(x token.Token) int {
	switch x {
	case token.IDENT:
		return _IDENT
	case token.NUMBER:
		return _NUMBER
	case token.ASSIGN:
		return _ASSIGN
	case token.EQ:
		return _EQ
	case token.LT:
		return _LT
	case token.PLUS:
		return _PLUS
	case token.MINUS:
		return _MINUS
	case token.TIMES:
		return _TIMES
	case token.OVER:
		return _OVER
	case token.LPAREN:
		return _LPAREN
	case token.RPAREN:
		return _RPAREN
	case token.SEMI:
		return _SEMI
	case token.IF:
		return _IF
	case token.THEN:
		return _THEN
	case token.ELSE:
		return _ELSE
	case token.END:
		return _END
	case token.REPEAT:
		return _REPEAT
	case token.UNTIL:
		return _UNTIL
	case token.READ:
		return _READ
	case token.WRITE:
		return _WRITE
	}
	return 0
}
