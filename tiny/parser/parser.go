// Copyright 2020 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TINY 语言词法和语法分析器
package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/chai2010/tinylang/tiny/ast"
	"github.com/chai2010/tinylang/tiny/token"
)

var _ = fmt.Sprint

// GoyaccOption 设置 goyacc 代码的调试选项.
func GoyaccOption(debug int, errVerbose bool) {
	yyDebug = debug
	yyErrorVerbose = errVerbose
}

func yyDebugf(format string, a ...interface{}) {
	if yyDebug > 0 {
		fmt.Printf("DEBUG: "+format, a...)
	}
}
func yyDebugln(a ...interface{}) {
	if yyDebug > 0 {
		fmt.Print("DEBUG: ")
		fmt.Println(a...)
	}
}

type yyLexerImpl struct {
	scanner     *Scanner
	lastSymType *yySymType
	lastErr     error
}

var _ yyLexer = (*yyLexerImpl)(nil)

func yyNewLexer(scanner *Scanner) *yyLexerImpl {
	p := &yyLexerImpl{scanner: scanner}
	return p
}

func (p *yyLexerImpl) Lex(lval *yySymType) int {
	p.lastSymType = lval

	lval.ctx.pos, lval.ctx.tok, lval.ctx.lit = p.scanner.Scan()
	lval.ctx.pos -= token.Pos(len(lval.ctx.lit))

	if lval.ctx.tok == token.EOF {
		return 0
	}
	if lval.ctx.tok == token.ILLEGAL {
		p.lastErr = fmt.Errorf("%s: illegal token: %v",
			p.scanner.PosString(lval.ctx.pos),
			lval.ctx.lit,
		)
		return 0
	}

	if lval.ctx.tok == token.NUMBER {
		_, err := strconv.Atoi(lval.ctx.lit)
		if err != nil {
			p.lastErr = err
			return 0
		}
	}

	return tok2yyTok(lval.ctx.tok)
}

func (p *yyLexerImpl) Error(s string) {
	p.lastErr = errors.New("goyacc: " + s)
	return
}

// ParseFile 将 TINY 文件解析为语法树.
func ParseFile(filename string, src interface{}) (f *ast.File, err error) {
	text, err := readSource(filename, src)
	if err != nil {
		return nil, err
	}

	scanner := NewScanner(filename, text)
	lexer := yyNewLexer(scanner)

	parser := yyNewParser()
	parser.Parse(lexer)

	f = lexer.lastSymType.ctx.file
	if f != nil {
		f.Name = filename
		f.Data = text
	}

	err = lexer.lastErr
	return
}

func readSource(filename string, src interface{}) ([]byte, error) {
	if src != nil {
		switch s := src.(type) {
		case string:
			return []byte(s), nil
		case []byte:
			return s, nil
		case *bytes.Buffer:
			// is io.Reader, but src is already available in []byte form
			if s != nil {
				return s.Bytes(), nil
			}
		case io.Reader:
			return io.ReadAll(s)
		}
		return nil, errors.New("invalid source")
	}
	return os.ReadFile(filename)
}
