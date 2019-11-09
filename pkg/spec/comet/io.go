// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

// IO外设相关
const (
	IO_ADDR  = 0xFD10 // 数据地址
	IO_FLAG  = 0xFD11 // 标志位
	IO_FIO   = 0x0100 // 输入输出
	IO_TYPE  = 0x1C00 // 传输类型
	IO_MAX   = 0x00FF // 最大数目
	IO_ERROR = 0x0200 // 错误位
	IO_IN    = 0x0000 // 输入
	IO_OUT   = 0x0100 // 输出
	IO_CHR   = 0x0400 // 字符
	IO_OCT   = 0x0800 // 八进制
	IO_DEC   = 0x0C00 // 十进制
	IO_HEX   = 0x1000 // 十六进制
)
