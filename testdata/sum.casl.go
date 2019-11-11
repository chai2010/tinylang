package main

const GR0 = 0

var AC = newInt16(1)     // AC 存放临时变量
var ABBAAA = newInt16(1) // 对应变量 sum
var ABAAAA = newInt16(1) // 对应变量 n
var CASL00 = newInt16(0) // CASL00 为程序的启动地址

func newInt16(v int16) int16 {
	return 0
}

func READ(adr int16) {
	return
}

func LD(gr int16, dar int16)  {}
func ST(gr int16, dar int16)  {}
func LEA(gr int16, dar int16) {}
func CPA(gr int16, dar int16) {}
func JPZ(dar int16)           {}

func start(pc uint16) {
	READ(ABBAAA)    // read 语句
	LD(GR0, ABAAAA) // if 判断语句开始
	_ = 0           // 对应变量; 运算符右边的值
	ST(GR0, AC)     // 保存运算符右边的值
	LEA(GR0, 0)     // 对应常量; 运算符左边的值
	CPA(GR0, AC)    // 计算表达式的值
	//JPZ(ABBBAA)     // if 语句，跳转到else部分
	//LEA	GR0,	0		; 对应常量
	//ST	GR0,	ABBAAA		; assign 语句
	//				; repeat 循环语句开始
}

func main() {}

/*
	START	CASL00			; 程序入口
AC    	DS	1			; AC 存放临时变量
ABBAAA	DS	1			; 对应变量 sum
ABAAAA	DS	1			; 对应变量 n
CASL00	DS	0			; CASL00 为程序的启动地址
	READ	ABAAAA			; read 语句
					; if 判断语句开始
	LD	GR0,	ABAAAA		; 对应变量; 运算符右边的值
	ST	GR0,	AC		; 保存运算符右边的值
	LEA	GR0,	0		; 对应常量; 运算符左边的值
	CPA	GR0,	AC		; 计算表达式的值
	JPZ	ABBBAA			; if 语句，跳转到else部分
	LEA	GR0,	0		; 对应常量
	ST	GR0,	ABBAAA		; assign 语句
					; repeat 循环语句开始
ABBBBB	DS	0			; 对应 repeat 语句开始地址
	LD	GR0,	ABAAAA		; 对应变量; 运算符右边的值
	ST	GR0,	AC		; 保存运算符右边的值
	LD	GR0,	ABBAAA		; 对应变量; 运算符左边的值
	ADD	GR0,	AC		; 计算表达式的值
	ST	GR0,	ABBAAA		; assign 语句
	LEA	GR0,	1		; 对应常量; 运算符右边的值
	ST	GR0,	AC		; 保存运算符右边的值
	LD	GR0,	ABAAAA		; 对应变量; 运算符左边的值
	SUB	GR0,	AC		; 计算表达式的值
	ST	GR0,	ABAAAA		; assign 语句
	LEA	GR0,	0		; 对应常量; 运算符右边的值
	ST	GR0,	AC		; 保存运算符右边的值
	LD	GR0,	ABAAAA		; 对应变量; 运算符左边的值
	CPA	GR0,	AC		; 计算表达式的值
	JNZ	ABBBBB			; 跳转到 repeat 循环语句开始
					; repeat 循环语句结束
	LD	GR0,	ABBAAA		; 对应变量
	ST	GR0,	AC		; 保存 write 值在 AC 中
	WRITE	AC			; write 语句，输出 AC 值
	JMP	ABBBBA			; if 语句，跳转到end部分
ABBBAA	DS	0			; 对应 if 语句的 else 地址
ABBBBA	DS	0			; 对应 if 语句的 end 地址
					; if 语句结束
	HALT				; 停机
	END				; 程序结束
*/
