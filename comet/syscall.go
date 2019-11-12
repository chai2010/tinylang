// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"fmt"
	"log"
)

// 内置的系统调用
const (
	SYSCALL_READ  = 1 // 读一个十进制整数到GR0
	SYSCALL_WRITE = 2 // 输出GR10十进制格式整数

	SYSCALL_IN   = 3 // 读N个字符, GR0是地址, GR1是N
	SYSCALL_OUT  = 4 // 写N个字符, GR0是地址, GR1是N
	SYSCALL_EXIT = 5 // 结束程序

	SYSCALL_USER_START = 64 // 用户的系统调号从此开始
)

// 注册内置的系统调用
func init() {
	RegisterSyscall(SYSCALL_READ, builtinSyscall_readInt)
	RegisterSyscall(SYSCALL_WRITE, builtinSyscall_writeInt)

	RegisterSyscall(SYSCALL_IN, builtinSyscall_readStr)
	RegisterSyscall(SYSCALL_OUT, builtinSyscall_writeStr)

	RegisterSyscall(SYSCALL_EXIT, builtinSyscall_exit)
}

// 系统调用表格
var syscallTable [256]func(ctx *Comet)

// 系统调用
func Syscall(ctx *Comet, id uint8) {
	if fn := syscallTable[id]; fn != nil {
		fn(ctx)
	}
}

// 注册系统调用(会覆盖之前的系统调用)
func RegisterSyscall(id uint8, syscall func(ctx *Comet)) error {
	if syscallTable[id] != nil {
		log.Printf("COMET: 系统调用 [%d] 被覆盖\n", id)
	}
	syscallTable[id] = syscall
	return nil
}

// 读一个十进制整数到GR0
func builtinSyscall_readInt(ctx *Comet) {
	fmt.Fscanf(ctx.Stdin, "%d", ctx.GR[0])
}

// 输出GR10十进制格式整数
func builtinSyscall_writeInt(ctx *Comet) {
	fmt.Fprintln(ctx.Stdout, ctx.GR[0])
}

// 读N个字符, GR0是地址, GR1是N
func builtinSyscall_readStr(ctx *Comet) {
	var adr = ctx.GR[0]
	var cnt = ctx.GR[1]
	for i := uint16(0); i < cnt; i++ {
		var c rune
		fmt.Fscanf(ctx.Stdin, "%c", &c)
		ctx.Mem[adr+i] = uint16(c)
	}
}

// 写N个字符, GR0是地址, GR1是N
func builtinSyscall_writeStr(ctx *Comet) {
	var adr = ctx.GR[0]
	var cnt = ctx.GR[1]
	for i := uint16(0); i < cnt; i++ {
		fmt.Fprint(ctx.Stdout, rune(ctx.Mem[adr+i]))
	}
}

// 退出程序(和停机类似)
func builtinSyscall_exit(ctx *Comet) {
	ctx.Shutdown = true
}
