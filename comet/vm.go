// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

const (
	MEM_SIZE = 1 << 16 // 内存大小

	SP_START = 0xFC00 // SP栈开始地址
	PC_START = 0x0000 // PC默认开始地址
	PC_MAX   = 0xFC00 // PC最大地址
)

type Comet struct {
	CPU
	Stdin    *bufio.Reader              // 标准输入输出(VM自身使用)
	Stdout   io.Writer                  // 标准输入输出(VM自身使用)
	Shutdown bool                       // 已经关机
	Syscall  func(ctx *Comet, id uint8) // 系统调用(GR0是返回值)
}

type CPU struct {
	PC  uint16          // 指令计数器
	FR  int16           // 标志寄存器
	GR  [5]uint16       // 通用寄存器
	Mem [1 << 16]uint16 // 64KB内存
}

type stdReadWriter struct{}

func (*stdReadWriter) Read(p []byte) (n int, err error) {
	return os.Stdin.Read(p)
}
func (*stdReadWriter) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func NewComet(prog []uint16, pc int) *Comet {
	p := new(Comet)
	copy(p.Mem[:], prog)

	p.PC = uint16(pc)
	p.GR[4] = SP_START

	p.Stdin = bufio.NewReader(os.Stdin)
	p.Stdout = os.Stdout

	return p
}

func (p *Comet) Run() {
	if p.Shutdown {
		return
	}
	for !p.Shutdown {
		p.StepRun()
	}
}

func (p *Comet) StepRun() {
	if p.Shutdown {
		return
	}

	var op = OpType(p.Mem[p.PC] / 0x100)
	var gr = (p.Mem[p.PC] % 0x100) / 0x10
	var xr = p.Mem[p.PC] % 0x10
	var adr = p.Mem[p.PC+1]
	var syscalId = uint8(p.Mem[p.PC] % 0x100)

	if gr < 0 || gr > 4 {
		fmt.Printf("非法指令：mem[%x] = %x\n", p.PC, p.Mem[p.PC])
		p.Shutdown = true
		return
	}
	if xr < 0 || xr > 4 {
		fmt.Printf("非法指令：mem[%x] = %x\n", p.PC, p.Mem[p.PC])
		p.Shutdown = true
		return
	}
	if xr != 0 {
		adr = uint16(int32(adr) + int32(p.GR[xr]))
	}

	// 临时: 处理IO
	p.io()

	// 指令解码
	switch op {
	case HALT:
		p.PC += 1
		p.Shutdown = true
	case LD:
		p.PC += 2
		p.GR[gr] = p.Mem[adr]
	case ST:
		p.PC += 2
		p.Mem[adr] = p.GR[gr]
	case LEA:
		p.PC += 2
		p.GR[gr] = adr
		p.FR = int16(p.GR[gr])
	case ADD:
		p.PC += 2
		p.GR[gr] = uint16(int16(p.GR[gr]) + int16(p.Mem[adr]))
		p.FR = int16(p.GR[gr])
	case SUB:
		p.PC += 2
		p.GR[gr] = uint16(int16(p.GR[gr]) - int16(p.Mem[adr]))
		p.FR = int16(p.GR[gr])
	case MUL:
		p.PC += 2
		p.GR[gr] = uint16(int16(p.GR[gr]) * int16(p.Mem[adr]))
		p.FR = int16(p.GR[gr])
	case DIV:
		p.PC += 2
		p.GR[gr] = uint16(int16(p.GR[gr]) / int16(p.Mem[adr]))
		p.FR = int16(p.GR[gr])
	case MOD:
		p.PC += 2
		p.GR[gr] = uint16(int16(p.GR[gr]) % int16(p.Mem[adr]))
		p.FR = int16(p.GR[gr])
	case AND:
		p.PC += 2
		p.GR[gr] &= p.Mem[adr]
		p.FR = int16(p.GR[gr])
	case OR:
		p.PC += 2
		p.GR[gr] |= p.Mem[adr]
		p.FR = int16(p.GR[gr])
	case EOR:
		p.PC += 2
		p.GR[gr] ^= p.Mem[adr]
		p.FR = int16(p.GR[gr])
	case SLA:
		p.PC += 2
		p.GR[gr] = uint16(int16(p.GR[gr]) << int16(p.Mem[adr]))
		p.FR = int16(p.GR[gr])
	case SRA:
		p.PC += 2
		p.GR[gr] = uint16(int16(p.GR[gr]) >> int16(p.Mem[adr]))
		p.FR = int16(p.GR[gr])
	case SLL:
		p.PC += 2
		p.GR[gr] = p.GR[gr] << p.Mem[adr]
		p.FR = int16(p.GR[gr])
	case SRL:
		p.PC += 2
		p.GR[gr] = p.GR[gr] >> p.Mem[adr]
		p.FR = int16(p.GR[gr])
	case CPA:
		p.PC += 2
		p.FR = int16(p.GR[gr]) - int16(p.Mem[adr])
	case CPL:
		p.PC += 2
		p.FR = int16(p.GR[gr] - p.Mem[adr])
	case JMP:
		p.PC += 2
		p.PC = adr
	case JPZ:
		p.PC += 2
		if p.FR >= 0 {
			p.PC = adr
		}
	case JMI:
		p.PC += 2
		if p.FR < 0 {
			p.PC = adr
		}
	case JNZ:
		p.PC += 2
		if p.FR != 0 {
			p.PC = adr
		}
	case JZE:
		p.PC += 2
		if p.FR == 0 {
			p.PC = adr
		}
	case PUSH:
		p.PC += 2
		p.Mem[p.GR[4]-1] = p.Mem[adr]
		p.GR[4]--
	case POP:
		p.PC += 1
		p.GR[gr] = p.Mem[p.GR[4]]
		p.GR[4]++
	case CALL:
		p.PC += 2
		p.Mem[p.GR[4]-1] = p.PC
		p.PC = p.Mem[adr]
		p.GR[4]--
	case RET:
		p.PC += 1
		p.PC = p.Mem[p.GR[4]]
		p.GR[4]++

	case SYSCALL:
		p.PC += 1
		p.Syscall(p, syscalId)

	default:
		p.Shutdown = true
		fmt.Printf("非法指令：mem[%x] = %x\n", p.PC, p.Mem[p.PC])
	}
}

func (p *Comet) DebugRun() {
	var (
		backup  = *p
		stepcnt int
		pntflag bool
		traflag bool
	)

	fmt.Println("调试 （帮助输入 help）...")
	fmt.Println()

	for {
		fmt.Print("输入命令: ")
		line, _, _ := p.Stdin.ReadLine()

		// 删除空白字符
		line = bytes.TrimSpace(line)

		// 跳过空白行
		if string(line) == "" {
			fmt.Fprintln(p.Stdout)
			continue
		}

		var cmd, x1, x2 = "", 0, 0
		n, _ := fmt.Fscanf(bytes.NewBuffer(line), "%s%x%x", &cmd, &x1, &x2)

		switch cmd {
		case "help", "h":
			fmt.Println(p.DebugHelp())
		case "go", "g":
			if p.Shutdown {
				fmt.Println("已经停机, 输入 `clear` 指令重置机器")
				continue
			}
			stepcnt = 0
			for !p.Shutdown {
				stepcnt++
				if traflag {
					fmt.Print(p.FormatInstruction(p.PC, 1))
				}

				// 单步执行(可能执行HALT关机指令)
				p.StepRun()
			}
			if pntflag {
				fmt.Printf("执行指令数目 = %d\n", stepcnt)
			}

		case "step", "s":
			if p.Shutdown {
				fmt.Println("已经停机, 输入 `clear` 指令重置机器")
				continue
			}

			if n >= 2 {
				stepcnt = x1
			} else {
				stepcnt = 1
			}

			var i int
			for i = 0; i < stepcnt && !p.Shutdown; i++ {
				if traflag {
					fmt.Print(p.FormatInstruction(p.PC, 1))
				}

				// 单步执行(可能执行HALT关机指令)
				p.StepRun()
			}
			if pntflag {
				fmt.Printf("执行指令数目 = %d\n", i)
			}

		case "jump", "j":
			if n >= 2 {
				fmt.Printf("指令跳转到 %x\n", x1)
				p.PC = uint16(x1)
			} else {
				fmt.Println("错误: 缺少跳转地址")
			}

		case "regs", "r":
			fmt.Println("显示寄存器数据")

			switch {
			case p.FR > 0:
				fmt.Printf("GR[0] = %04x\tPC = %04x\n", p.GR[0], p.PC)
				fmt.Printf("GR[1] = %04x\tSP = %04x\n", p.GR[1], uint16(p.GR[4]))
				fmt.Printf("GR[2] = %04x\tFR = ..00\n", p.GR[2])
				fmt.Printf("GR[3] = %04x\n", p.GR[3])
			case p.FR < 0:
				fmt.Printf("GR[0] = %04x\tPC = %04x\n", p.GR[0], p.PC)
				fmt.Printf("GR[1] = %04x\tSP = %04x\n", p.GR[1], uint16(p.GR[4]))
				fmt.Printf("GR[2] = %04x\tFR = ..10\n", p.GR[2])
				fmt.Printf("GR[3] = %04x\n", p.GR[3])
			default:
				fmt.Printf("GR[0] = %04x\tPC = %04x\n", p.GR[0], p.PC)
				fmt.Printf("GR[1] = %04x\tSP = %04x\n", p.GR[1], uint16(p.GR[4]))
				fmt.Printf("GR[2] = %04x\tFR = ..01\n", p.GR[2])
				fmt.Printf("GR[3] = %04x\n", p.GR[3])
			}

		case "iMem", "imem", "i":
			fmt.Println("显示内存指令")

			x1 := uint16(x1)
			if n < 2 {
				x1 = p.PC
			}
			if n < 3 {
				x2 = 1
			}

			fmt.Print(p.FormatInstruction(x1, x2))

		case "dMem", "dmem", "d":
			x1 := uint16(x1)
			if n < 2 {
				x1 = p.PC
			}
			if n < 3 {
				x2 = 1
			}

			for i := 0; i < x2 && i < len(p.Mem); i++ {
				fmt.Printf("mem[%04x] = %04x\n", x1, uint16(p.Mem[x1]))
				x1++
			}

		case "alter", "a":
			if n == 3 {
				fmt.Printf("修改内存数据  mem[%x] = %x\n", x1, x2)
				p.Mem[x1] = uint16(x2)
			} else {
				fmt.Println("修改内存数据 失败！")
			}

		case "trace", "t":
			traflag = !traflag
			if traflag {
				fmt.Println("指令显示功能 打开")
			} else {
				fmt.Println("指令显示功能 关闭")
			}

		case "print", "p":
			pntflag = !pntflag
			if pntflag {
				fmt.Println("指令计数功能 打开")
			} else {
				fmt.Println("指令计数功能 关闭")
			}

		case "clear", "c":
			fmt.Println("程序重新载入内存")
			*p = backup
			stepcnt = 0

		case "quit", "q":
			fmt.Println("退出调试...")
			return

		default:
			fmt.Println("未知命令", cmd)
		}
	}
}

func (p *Comet) DebugHelp() string {
	return `命令列表:
  h)elp           显示本命令列表
  g)o             运行程序直到停止
  s)tep  <n>      执行 n 条指令 （默认为 1 ）
  j)ump  <b>      跳转到 b 地址 （默认为当前地址）
  r)egs           显示寄存器内容
  i)Mem  <b <n>>  显示从 b 开始 n 个内存数据
  d)Mem  <b <n>>  显示从 b 开始 n 个内存指令
  a(lter <b <v>>  修改 b 位置的内存数据为 v 值
  t)race          开关指令显示功能
  p)rint          开关指令计数功能
  c)lear          重置模拟器内容
  q)uit           终止模拟器
`
}

// 格式化pc开始的n个指令
func (p *Comet) FormatInstruction(pc uint16, n int) string {
	var buf bytes.Buffer

	for i := 0; i < n; i++ {
		ins, ok := p.ParseInstruction(pc)
		if !ok {
			fmt.Fprintf(&buf, "mem[%04x]: 未知\n", pc)
			break
		}

		fmt.Fprintf(&buf, "mem[%04x]: %v\n", pc, ins)
		pc += ins.Op.Size()
	}

	return buf.String()
}

func (p *Comet) io() {
	cnt := p.Mem[IO_FLAG] & IO_MAX
	if cnt == 0 {
		return
	}

	fio := p.Mem[IO_FLAG] & IO_FIO
	typ := p.Mem[IO_FLAG] & IO_TYPE
	adr := p.Mem[IO_ADDR]

	var format string
	switch {
	case typ == IO_CHR:
		format = "%c"
	case typ == IO_OCT:
		format = "%o"
	case typ == IO_DEC:
		format = "%d"
	case typ == IO_HEX:
		format = "%x"
	default:
		p.Mem[IO_FLAG] |= IO_ERROR
		p.Mem[IO_FLAG] &= ^uint16(IO_MAX)
		return
	}

	for i := 0; i < int(cnt); i++ {
		if fio == IO_IN {
			fmt.Fscanf(p.Stdin, format, &p.Mem[adr])
			adr++
		} else {
			fmt.Fprintf(p.Stdout, format, p.Mem[adr])
			adr++
		}
	}

	p.Mem[IO_FLAG] &= ^uint16(IO_MAX)
}
