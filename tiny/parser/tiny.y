// Copyright 2020 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/****************************************************/
/* File: tiny.y                                     */
/* The TINY Yacc/Bison specification file           */
/* Compiler Construction: Principles and Practice   */
/* Kenneth C. Louden                                */
/****************************************************/
%{
package parser

import (
	"fmt"
	"strconv"

	"github.com/chai2010/tinylang/tiny/ast"
	"github.com/chai2010/tinylang/tiny/token"
)

var _ = fmt.Sprint
%}

%union {
	ctx struct {
		file     *ast.File
		expr     ast.Expr
		stmt     ast.Stmt
		stmtList []ast.Stmt
		node     ast.Node

		pos token.Pos
		tok token.Token
		lit string
	}
}

// ctx.ident
%token <ctx> _IDENT

// ctx.lit
%token <ctx> _NUMBER

// ctx.tok
%token <ctx> _EQ _LT _PLUS _MINUS _TIMES _OVER
%token <ctx> _ASSIGN _IF _THEN _ELSE _END _REPEAT _UNTIL _READ _WRITE
%token <ctx> _LPAREN _RPAREN _SEMI

// ctx.file
%type  <ctx>  program

// ctx.expr
%type  <ctx>  factor exp term simple_exp

// ctx.stmt
%type  <ctx>  assign_stmt
%type  <ctx>  if_stmt
%type  <ctx>  repeat_stmt
%type  <ctx>  read_stmt
%type  <ctx>  write_stmt
%type  <ctx>  stmt

// ctx.stmtList
%type  <ctx>  stmt_seq

%% /* Grammar for TINY */

program
	: stmt_seq {
		yyDebugln("program.stmt_seq end")
		$$.file = &ast.File {
			List: $1.stmtList,
		}

		// https://github.com/golang/tools/commit/a965a571dd795205ab6582d15ea92e1350374b58
		//
		// Code inside the grammar actions may refer to yyrcvr,
		// which holds the yyParser.
		yyrcvr.lval.ctx.file = $$.file
	}
	| {
		yyDebugln("program.<empty>")
		$$.file = &ast.File {}
	}

stmt_seq
	: stmt_seq _SEMI stmt {
		yyDebugln("stmt_seq.stmt_seq SEMI stmt")
		$$.stmtList = append($$.stmtList, $3.stmt)
	}
	| stmt {
		yyDebugln("stmt_seq.stmt")
		if $1.stmt != nil {
			$$.stmtList = []ast.Stmt{$1.stmt}
		}
	}

stmt
	: if_stmt {
		yyDebugln("stmt.if_stmt")
		$$.stmt = $1.stmt;
	}
	| repeat_stmt {
		yyDebugln("stmt.repeat_stmt")
		$$.stmt = $1.stmt;
	}
	| assign_stmt {
		yyDebugln("stmt.assign_stmt")
		$$.stmt = $1.stmt;
	}
	| read_stmt {
		yyDebugln("stmt.read_stmt end")
		$$.stmt = $1.stmt;
	}
	| write_stmt {
		yyDebugln("stmt.write_stmt")
		$$.stmt = $1.stmt;
	}
	| {
		$$.stmt = nil
	}

if_stmt
	: _IF exp _THEN stmt_seq _END {
		yyDebugln("if_stmt.IF exp THEN stmt_seq END")
		$$.stmt = &ast.IfStmt {
			If:   $1.pos,
			Cond: $2.expr,
			Body: &ast.BlockStmt {
				List: $4.stmtList,
			},
		}
	}
	| _IF exp _THEN stmt_seq _ELSE stmt_seq _END {
		yyDebugln("if_stmt.IF exp THEN stmt_seq ELSE stmt_seq END")
		$$.stmt = &ast.IfStmt {
			If:   $1.pos,
			Cond: $2.expr,
			Body: &ast.BlockStmt {
				List: $4.stmtList,
			},
			Else: &ast.BlockStmt {
				List: $6.stmtList,
			},
		}
	}

repeat_stmt
	: _REPEAT stmt_seq _UNTIL exp {
		yyDebugln("repeat_stmt.REPEAT stmt_seq UNTIL exp")
		$$.stmt = &ast.RepeatStmt {
			Repeat: $1.pos,
			Until:  $4.expr,
			Body:   &ast.BlockStmt {
				List: $2.stmtList,
			},
		}
	}

assign_stmt
	: _IDENT _ASSIGN exp {
		yyDebugln("assign_stmt.ASSIGN exp")
		$$.stmt = &ast.AssignStmt {
			Target: &ast.Ident{
				NamePos: $1.pos,
				Name:    $1.lit,
			},
			TokPos: $2.pos,
			Value:  $3.expr,
		}
	}

read_stmt
	: _READ _IDENT {
		yyDebugln("read_stmt.READ ID")
		$$.stmt = &ast.ReadStmt {
			Read:   $1.pos,
			Target: &ast.Ident{
				NamePos: $2.pos,
				Name:    $2.lit,
			},
		}
	}

write_stmt
	: _WRITE exp {
		yyDebugln("write_stmt.WRITE exp")
		$$.stmt = &ast.WriteStmt {
			Write: $1.pos,
			Value: $2.expr,
		}
	}

exp
	: simple_exp _LT simple_exp {
		yyDebugln("exp.simple_exp LT simple_exp")
		$$.expr = &ast.BinaryExpr {
			X:     $1.expr,
			OpPos: $2.pos,
			Op:    $2.tok,
			Y:     $3.expr,
		}
	}
	| simple_exp _EQ simple_exp {
		yyDebugln("exp.simple_exp EQ simple_exp")
		$$.expr = &ast.BinaryExpr {
			X:     $1.expr,
			OpPos: $2.pos,
			Op:    $2.tok,
			Y:     $3.expr,
		}
	}
	| simple_exp {
		yyDebugln("exp.simple_exp")
		$$.expr = $1.expr;
	}

simple_exp
	: simple_exp _PLUS term {
		yyDebugln("simple_exp.simple_exp PLUS term")
		$$.expr = &ast.BinaryExpr {
			X:     $1.expr,
			OpPos: $2.pos,
			Op:    $2.tok,
			Y:     $3.expr,
		}
	}
	| simple_exp _MINUS term {
		yyDebugln("simple_exp.simple_exp MINUS term")
		$$.expr = &ast.BinaryExpr {
			X:     $1.expr,
			OpPos: $2.pos,
			Op:    $2.tok,
			Y:     $3.expr,
		}
	} 
	| term {
		yyDebugln("simple_exp.term")
		$$.expr = $1.expr;
	}

term
	: term _TIMES factor {
		yyDebugln("term.term TIMES factor")
		$$.expr = &ast.BinaryExpr {
			X:     $1.expr,
			OpPos: $2.pos,
			Op:    $2.tok,
			Y:     $3.expr,
		}
	}
	| term _OVER factor {
		yyDebugln("term.term OVER factor")
		$$.expr = &ast.BinaryExpr {
			X:     $1.expr,
			OpPos: $2.pos,
			Op:    $2.tok,
			Y:     $3.expr,
		}
	}
	| factor {
		yyDebugln("term.factor")
		$$.expr = $1.expr;
	}

factor
	: _LPAREN exp _RPAREN {
		yyDebugln("factor.LPAREN exp RPAREN")
		$$.expr = &ast.ParenExpr{
			Lparen: $1.pos,
			X:      $2.expr,
			Rparen: $3.pos,
		}
	}
	| _NUMBER {
		yyDebugln("factor.NUM")
		v, err := strconv.Atoi($1.lit)
		if err != nil {
			yylex.Error("aaa")
			return 0
		}
		$$.expr = &ast.Number{
			ValuePos: $1.pos,
			ValueEnd: $1.pos + token.Pos(len($1.lit)),
			Value:    v,
		}
	}
	| _IDENT {
		yyDebugln("factor.ID")
		$$.expr = &ast.Ident{
			NamePos: $1.pos,
			Name:    $1.lit,
		}
	}

%%
