// Copyright 2020 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

var (
	_ Node = (*File)(nil)
	_ Node = (Stmt)(nil)
	_ Node = (Expr)(nil)
)

var (
	_ Stmt = (*BlockStmt)(nil)
	_ Stmt = (*IfStmt)(nil)
	_ Stmt = (*RepeatStmt)(nil)
	_ Stmt = (*AssignStmt)(nil)
	_ Stmt = (*ReadStmt)(nil)
	_ Stmt = (*WriteStmt)(nil)
)

var (
	_ Expr = (*Ident)(nil)
	_ Expr = (*Number)(nil)
	_ Expr = (*BinaryExpr)(nil)
	_ Expr = (*ParenExpr)(nil)
)

func (p *File) node_private() {}

func (p *BlockStmt) node_private()  {}
func (p *IfStmt) node_private()     {}
func (p *RepeatStmt) node_private() {}
func (p *AssignStmt) node_private() {}
func (p *ReadStmt) node_private()   {}
func (p *WriteStmt) node_private()  {}

func (p *Ident) node_private()      {}
func (p *Number) node_private()     {}
func (p *BinaryExpr) node_private() {}
func (p *ParenExpr) node_private()  {}

func (p *Comment) node_private()      {}
func (p *CommentGroup) node_private() {}

func (p *BlockStmt) stmt_private()  {}
func (p *IfStmt) stmt_private()     {}
func (p *RepeatStmt) stmt_private() {}
func (p *AssignStmt) stmt_private() {}
func (p *ReadStmt) stmt_private()   {}
func (p *WriteStmt) stmt_private()  {}

func (p *Ident) expr_private()      {}
func (p *Number) expr_private()     {}
func (p *BinaryExpr) expr_private() {}
func (p *ParenExpr) expr_private()  {}
