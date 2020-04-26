// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"fmt"
	"go/token"
	"strings"
)

// 全局变量数据(CASL汇编定义数据用)
type DataStore struct {
	Name    string    // 数据名字
	Size    uint16    // 数据大小
	Comment string    // 注释
	Pos     token.Pos // 在源文件的位置
}

// 对应 CASL 的汇编格式
func (p *DataStore) String() string {
	name := p.Name

	// CASL 名字最多6个字母
	// 因此默认采用 6 个字母对齐
	if len(name) < 6 {
		name = (name + "      ")[:6]
	}

	var s = fmt.Sprintf("%s DS %d", name, p.Size)
	if p.Comment != "" {
		s = (s + strings.Repeat(" ", 24+7))[:24+7] + p.Comment
	}
	return s
}

// 字节码
func (p *DataStore) Bytes() []byte {
	return make([]byte, p.Size*2)
}
