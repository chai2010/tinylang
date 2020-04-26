// Copyright 2020 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"bytes"
	"text/scanner"

	"github.com/chai2010/tinylang/tiny/token"
)

type Scanner struct {
	filename string
	src      []byte
	scanner  scanner.Scanner
	lastErr  error
}

func NewScanner(filename string, src []byte) *Scanner {
	p := &Scanner{
		filename: filename,
		src:      src,
	}
	p.scanner.Init(bytes.NewReader(src))
	return p
}

func (p *Scanner) PosString(pos token.Pos) string {
	return token.PosString(p.filename, p.src, pos)
}

func (p *Scanner) Scan() (pos token.Pos, tok token.Token, lit string) {
	if p.lastErr != nil {
		return
	}

	x := p.scanner.Scan()
	lit = p.scanner.TokenText()
	pos = token.Pos(p.scanner.Pos().Offset + 1)

	switch x {
	case scanner.EOF:
		tok = token.EOF
		return
	case scanner.Int:
		tok = token.NUMBER
		return

	case scanner.Float:
		tok = token.ILLEGAL
		return

	case scanner.Ident:
		tok = token.Lookup(lit)
		return

	case scanner.Char:
		tok = token.ILLEGAL
		return
	case scanner.String:
		tok = token.ILLEGAL
		return
	case scanner.RawString:
		tok = token.ILLEGAL
		return
	case scanner.Comment:
		tok = token.ILLEGAL
		return

	default:
		opTok, opSuffix, ok := token.LookupOperator(x)
		if !ok {
			tok = token.ILLEGAL
			return
		}
		if opSuffix != 0 {
			if next := p.scanner.Scan(); next == opSuffix {
				lit += string(next)
				tok = opTok
				return
			} else {
				tok = token.ILLEGAL
				return
			}
		}
		tok = opTok
		return
	}
}
