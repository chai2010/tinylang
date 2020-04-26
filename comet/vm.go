// Copyright 2019 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// COMET 虚拟机
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
)

type Comet struct {
	CPU
	Stdin    *bufio.Reader                  // 标准输入输出(VM自身使用)
	Stdout   io.Writer                      // 标准输入输出(VM自身使用)
	Shutdown bool                           // 已经关机
	Syscall  func(ctx *Comet, id SyscallId) // 系统调用(GR0是返回值)
}

type CPU struct {
	PC  uint16          // 指令计数器
	FR  int16           // 标志寄存器
	GR  [5]uint16       // 通用寄存器
	Mem [1 << 16]uint16 // 64KB内存
}

// COMET 程序
type Program struct {
	PC  uint16
	Bin []uint16
}

func NewComet(prog *Program) *Comet {
	p := &Comet{
		Syscall: func(ctx *Comet, id SyscallId) {
			Syscall(ctx, id)
		},
	}
	copy(p.Mem[:], prog.Bin)

	p.PC = prog.PC
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
	var syscalId = SyscallId(adr)

	// 非法指令
	if !op.Valid() {
		fmt.Printf("Illegal code: mem[%x] = %x\n", p.PC, p.Mem[p.PC])
		p.Shutdown = true
		return
	}

	// 判断 GR
	if !op.UseGR() && gr != 0 {
		fmt.Printf("Illegal code: mem[%x] = %x; %s\n", p.PC, p.Mem[p.PC], "invalid GR")
		p.Shutdown = true
		return
	}
	if gr < 0 || gr > 4 {
		fmt.Printf("Illegal code: mem[%x] = %x； %s\n", p.PC, p.Mem[p.PC], "invalid GR")
		p.Shutdown = true
		return
	}

	// 判断 xr
	if !op.UseADR() && xr != 0 {
		fmt.Printf("Illegal code: mem[%x] = %x; %s\n", p.PC, p.Mem[p.PC], "invalid XR")
		p.Shutdown = true
		return
	}
	if xr < 0 || xr > 4 {
		fmt.Printf("Illegal code: mem[%x] = %x; %s\n", p.PC, p.Mem[p.PC], "invalid XR")
		p.Shutdown = true
		return
	}
	if xr != 0 {
		adr = uint16(int32(adr) + int32(p.GR[xr]))
	}

	// 调整 PC
	p.PC += op.Size()

	// 执行指令
	switch op {
	case HALT:
		p.Shutdown = true
	case LD:
		p.GR[gr] = p.Mem[adr]
	case ST:
		p.Mem[adr] = p.GR[gr]
	case LEA:
		p.GR[gr] = adr
		p.FR = int16(p.GR[gr])
	case ADD:
		p.GR[gr] = uint16(int16(p.GR[gr]) + int16(p.Mem[adr]))
		p.FR = int16(p.GR[gr])
	case SUB:
		p.GR[gr] = uint16(int16(p.GR[gr]) - int16(p.Mem[adr]))
		p.FR = int16(p.GR[gr])
	case MUL:
		p.GR[gr] = uint16(int16(p.GR[gr]) * int16(p.Mem[adr]))
		p.FR = int16(p.GR[gr])
	case DIV:
		p.GR[gr] = uint16(int16(p.GR[gr]) / int16(p.Mem[adr]))
		p.FR = int16(p.GR[gr])
	case MOD:
		p.GR[gr] = uint16(int16(p.GR[gr]) % int16(p.Mem[adr]))
		p.FR = int16(p.GR[gr])
	case AND:
		p.GR[gr] &= p.Mem[adr]
		p.FR = int16(p.GR[gr])
	case OR:
		p.GR[gr] |= p.Mem[adr]
		p.FR = int16(p.GR[gr])
	case EOR:
		p.GR[gr] ^= p.Mem[adr]
		p.FR = int16(p.GR[gr])
	case SLA:
		p.GR[gr] = uint16(int16(p.GR[gr]) << int16(p.Mem[adr]))
		p.FR = int16(p.GR[gr])
	case SRA:
		p.GR[gr] = uint16(int16(p.GR[gr]) >> int16(p.Mem[adr]))
		p.FR = int16(p.GR[gr])
	case SLL:
		p.GR[gr] = p.GR[gr] << p.Mem[adr]
		p.FR = int16(p.GR[gr])
	case SRL:
		p.GR[gr] = p.GR[gr] >> p.Mem[adr]
		p.FR = int16(p.GR[gr])
	case CPA:
		p.FR = int16(p.GR[gr]) - int16(p.Mem[adr])
	case CPL:
		p.FR = int16(p.GR[gr] - p.Mem[adr])
	case JMP:
		p.PC = adr
	case JPZ:
		if p.FR >= 0 {
			p.PC = adr
		}
	case JMI:
		if p.FR < 0 {
			p.PC = adr
		}
	case JNZ:
		if p.FR != 0 {
			p.PC = adr
		}
	case JZE:
		if p.FR == 0 {
			p.PC = adr
		}
	case PUSH:
		p.Mem[p.GR[4]-1] = p.Mem[adr]
		p.GR[4]--
	case POP:
		p.GR[gr] = p.Mem[p.GR[4]]
		p.GR[4]++
	case CALL:
		p.Mem[p.GR[4]-1] = p.PC
		p.PC = p.Mem[adr]
		p.GR[4]--
	case RET:
		p.PC = p.Mem[p.GR[4]]
		p.GR[4]++
	case NOP:
		// empty
	case SYSCALL:
		p.Syscall(p, syscalId)

	default:
		panic("unreachable")
	}
}

func (p *Comet) DebugRun() {
	var (
		backup  = *p
		stepcnt int
		pntflag bool
		traflag bool
	)

	fmt.Println("Debug (enter h for help)...")
	fmt.Println()

	for {
		fmt.Print("Enter command: ")
		line, _, _ := p.Stdin.ReadLine()

		// 删除空白字符
		line = bytes.TrimSpace(line)

		// 跳过空白行
		if string(line) == "" {
			fmt.Println()
			continue
		}

		var cmd, x1, x2 = "", 0, 0
		n, _ := fmt.Fscanf(bytes.NewBuffer(line), "%s%x%x", &cmd, &x1, &x2)

		switch cmd {
		case "help", "h":
			fmt.Println(p.DebugHelp())
		case "go", "g":
			if p.Shutdown {
				fmt.Println("halted, enter `clear` to reset VM")
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
				fmt.Printf("step count = %d\n", stepcnt)
			}

		case "step", "s":
			if p.Shutdown {
				fmt.Println("halted, enter `clear` to reset VM")
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
				fmt.Printf("step count = %d\n", i)
			}

		case "jump", "j":
			if n >= 2 {
				p.PC = uint16(x1)
				fmt.Printf("PC = %x\n", x1)
			} else {
				fmt.Println("invalid command")
			}

		case "regs", "r":
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
				fmt.Printf("mem[%x] = %x\n", x1, x2)
				p.Mem[x1] = uint16(x2)
			} else {
				fmt.Println("invalid command")
			}

		case "trace", "t":
			traflag = !traflag
			if traflag {
				fmt.Println("trace instruction now on")
			} else {
				fmt.Println("trace instruction now off")
			}

		case "print", "p":
			pntflag = !pntflag
			if pntflag {
				fmt.Println("Printing instruction count now on")
			} else {
				fmt.Println("Printing instruction count now off")
			}

		case "clear", "c":
			*p = backup
			stepcnt = 0

		case "quit", "q":
			return

		default:
			fmt.Println("unknown command:", cmd)
		}
	}
}

func (p *Comet) DebugHelp() string {
	return `Commands are:
  h)elp           show help command list
  g)o             run instructions until HALT
  s)tep  <n>      run n (default 1) instructions
  j)ump  <b>      jump to the b (default is current location)
  r)egs           print the contents of the registers
  i)Mem  <b <n>>  print n iMem locations starting at b
  d)Mem  <b <n>>  print n dMem locations starting at b
  a(lter <b <v>>  change the memory value at v
  t)race          toggle instruction trace
  p)rint          toggle print of total instructions executed
  c)lear          reset comet VM
  q)uit           exit
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
