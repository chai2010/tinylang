// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import "log"

// 内置的系统调用
const (
	SYSCALL_READ      = 0 // 读取N个字节
	SYSCALL_WRITE     = 1 // 写N个字节
	SYSCALL_READ_INT  = 2 // 读一个十进制整数
	SYSCALL_WRITE_INT = 3 // 写一个十进制整数
	SYSCALL_READ_STR  = 4 // 读一个字符串
	SYSCALL_WRITE_STR = 5 // 写一个字符串

	SYSCALL_USER_START = 64 // 用户的系统调号从此开始
)

// 注册内置的系统调用
func init() {
	RegisterSyscall(SYSCALL_READ, builtinSyscall_read)
	RegisterSyscall(SYSCALL_WRITE, builtinSyscall_write)

	RegisterSyscall(SYSCALL_READ_INT, builtinSyscall_readInt)
	RegisterSyscall(SYSCALL_WRITE_INT, builtinSyscall_writeInt)

	RegisterSyscall(SYSCALL_READ_STR, builtinSyscall_readStr)
	RegisterSyscall(SYSCALL_WRITE_STR, builtinSyscall_writeStr)
}

// 系统调用表格
var syscallTable [256]func(ctx *Comet, id uint8) uint16

// 系统调用
func Syscall(ctx *Comet, id uint8) uint16 {
	if fn := syscallTable[id]; fn != nil {
		return fn(ctx, id)
	}
	return 0
}

// 注册系统调用(会覆盖之前的系统调用)
func RegisterSyscall(id uint8, syscall func(ctx *Comet, id uint8) uint16) error {
	if syscallTable[id] != nil {
		log.Printf("COMET: 系统调用 [%d] 被覆盖\n", id)
	}
	syscallTable[id] = syscall
	return nil
}

func builtinSyscall_read(ctx *Comet, id uint8) uint16 {
	return 0 // TODO
}
func builtinSyscall_write(ctx *Comet, id uint8) uint16 {
	return 0 // TODO
}

func builtinSyscall_readInt(ctx *Comet, id uint8) uint16 {
	return 0 // TODO
}
func builtinSyscall_writeInt(ctx *Comet, id uint8) uint16 {
	return 0 // TODO
}

func builtinSyscall_readStr(ctx *Comet, id uint8) uint16 {
	return 0 // TODO
}
func builtinSyscall_writeStr(ctx *Comet, id uint8) uint16 {
	return 0 // TODO
}
