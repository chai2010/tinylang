// Copyright 2020 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

// Walk 遍历每个节点, 如果 fn 返回 false 则返回.
func Walk(node Node, fn func(node Node) bool) {
	walk(node, fn)
}

func walk(n Node, fn func(n Node) bool) {
	if n == nil || !fn(n) {
		return
	}
	switch n := n.(type) {
	case *File:
		for _, stmt := range n.List {
			walk(stmt, fn)
		}

	case *IfStmt:
		walk(n.Cond, fn)
		for _, stmt := range n.Body.List {
			walk(stmt, fn)
		}
		walk(n.Else, fn)

	case *RepeatStmt:
		for _, stmt := range n.Body.List {
			walk(stmt, fn)
		}
		walk(n.Until, fn)

	case *AssignStmt:
		walk(n.Target, fn)
		walk(n.Value, fn)

	case *ReadStmt:
		walk(n.Target, fn)

	case *WriteStmt:
		walk(n.Value, fn)

	case *Ident:
		return
	case *Number:
		return

	case *BinaryExpr:
		walk(n.X, fn)
		walk(n.Y, fn)

	case *ParenExpr:
		walk(n.X, fn)

	default:
		panic("unreachable")
	}
}
