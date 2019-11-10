// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

// 一个词法元素
type Item struct {
	Typ Token  // 记号类型
	Val string // 字符串值
	Num int    // 数字值
	Pos int    // 开始位置
	End int    // 结束位置
}

func (i Item) String() string {
	switch {
	case i.Typ == ILLEGAL:
		return i.Val // 错误
	case i.Typ == EOF:
		return "EOF"
	case i.Typ == EOL:
		return "EOL"
	case i.Typ.IsKeyword():
		return fmt.Sprintf("<%s>", i.Val)
	case len(i.Val) > 10:
		return fmt.Sprintf("%.10q...", i.Val)
	}
	return fmt.Sprintf("%q", i.Val)
}
