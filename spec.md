Tiny玩具语言
==========


## COMET虚拟计算机说明

COMET是一台字长为16位的定点计算机，主存储器的容量是65536字节，按编号`0000－FFFF`(十六进制)编址。一个字的位数为`0 1 2 ⋯ 15`。

COMET机有5个通用寄存器GR(16位)，一个指令计数器PC(16位)和一个标志寄存器FR(2位)。其中GR1，GR2，GR3，GR4通用寄存器兼作变址寄存器。另外，GR4还兼作栈指针(SP)用，栈指针是存放栈顶地址用的寄存器。PC(指令寄存器)　在执行指令的过程中，PC中存放着正在执行的指令的第一个字的地址(一条指令占两个字)。当指令执行结束时，一般是把PC的内容加2，只有在执行转移指令且条件成立时，才将转移指令地址置入PC中。FR(标志寄存器)　在ADD，SUB，MUL，DIV，MOD，AND，OR，EOR，CPA，CPL，SLA，SRA，SLL，SRL，LEA等指令执行结束时，根据执行的结果，将FR置成00，01或10(大于、等于、小于；或负数、零、正数)。它不会因其它指令的执行而改变。

COMET指令格式： OP GR，ADR[，XR]，其中OP对应第一个字的高四位(0-7位)，GR为第一个字的(8-11位)，XR为第一个字的(12-15位)，ADR对应第二个字；即一个指令为两个字长。如果为直接寻址，即无XR，则第一个字的8-11为全部为0(GR0不能用作变址寻址)！

编码，E为地址(等于ADR+[XR])，在LEA等指令中为直接数，即称取址

```
    HALT,   // 0X0  停机
    LD,     // 0X1  取数，GR = (E) 
    ST,     // 0X2  存数，E = (GR)
    LEA,    // 0X3  取地址，GR = E
    
    ADD,    // 0X4  相加，GR = (GR)+(E)
    SUB,    // 0X5  相减，GR = (GR)-(E)
    MUL,    // 0X6  相乘，GR = (GR)*(E)
    DIV,    // 0X7  相除，GR = (GR)/(E)
    MOD,    // 0X8  取模，GR = (GR)%(E)

    AND,    // 0X9  与，GR = (GR)&(E)
    OR,     // 0XA  或，GR = (GR)|(E)
    EOR,    // 0XB  异或，GR = (GR)^(E)
    
    CPA,    // 0XC  算术比较，(GR)-(E)，有符号数，设置FR
    CPL,    // 0XD  逻辑比较，(GR)-(E)，无符号数，设置FR
    
    SLA,    // 0XE   算术左移，空出的的位置补0
    SRA,    // 0XF   算术右移，空出的的位置被置成第0位的值
    SLL,    // 0X10  逻辑左移，空出的的位置补0
    SRL,    // 0X11  逻辑右移，空出的的位置被置0
    
    JMP,    // 0X12 无条件跳转，PC = E 
    JPZ,    // 0X13 不小于跳转，PC = E
    JMI,    // 0X14 小于跳转，  PC = E
    JNE,    // 0X15 不等于跳转，PC = E
    JZE,    // 0X16 等于跳转，  PC = E
    
    PUSH,   // 0X17 进栈，SP = (SP)-1，(SP) = E
    POP,    // 0X18 出栈，GR = ((SP))，SP = (SP)+1
    
    CALL,   // 0X19 调用，SP = (SP)-1，(SP) = (PC)+2，PC = E
    RET     // 0X1A 返回，SP = (SP)+1
 ```
 
 外设备，用户可以自己配置

   a) 输出输出设备(键盘和显示器)。

   有两个设备寄存器，IO_ADDR、IO_FLAG。IO_ADDR保存要传输数据的内存地址；

   IO_FLAG为IO的标志位，其8-15位是要传输数据的个数(0表示无IO)，7位表示

   输入或输出(1表示输入，0为输出)，6位在出现IO错误时设置，3-5位为传输的

   类型(有字符、八进制、十进制、十六进制等)，0-2保留(可能用于表示IO设备)

   b) 相关的值

```
   IO_ADDR  = 0xFD10  // 数据地址
   IO_FLAG  = 0xFD11  // 标志位
   IO_FIO   = 0x0100  // 输入输出
   IO_TYPE  = 0x1C00  // 传输类型
   IO_MAX   = 0x00FF  // 最大数目
   IO_ERROR = 0x0200  // 错误位
   IO_IN    = 0x0000  // 输入
   IO_OUT   = 0x0100  // 输出
   IO_CHR   = 0x0400  // 字符
   IO_OCT   = 0x0800  // 八进制
   IO_DEC   = 0x0C00  // 十进制
   IO_HEX   = 0x1000  // 十六进制
```

内存使用约定: comet计算机有64k字的内存，默认程序从0地址装如，栈FC00向下增长，FC00-FCFF的526字空间机器保留，FD00-FDFF的256字为外设备寄存器区(如IO设备)FE00-FEFF的256字为系统使用的临时数据区，FF00-FEFF为系统使用的临时数据区

comet计算机集成调试功能，用户可以在使用时参考帮助。comet计算机目前还不完全支持unicode字符。


## CASL汇编语言说明

1. CASL由4种伪指令(START、END、DS、DC)，5种宏指令(READ、WRITE、IN、OUT、EXIT)和27种

符号指令(COMET的指令)组成。CASL的每条指令书写在一行内(最多不超过72个字符)，它的书写

格式如下：

```
    标号        指令码     操作数              注释

   [LABEL]      START      [LABEL]
   [LABEL]      END        空白
   [LABEL]      DC         常数
   [LABEL]      DS         区域的字数
   [LABEL]      READ       ALABEL
   [LABEL]      WRITE      ALABEL
   [LABEL]      IN         ALABEL,NLABEL
   [LABEL]      OUT        ALABEL,NLABEL
   [LABEL]      EXIT       空白
   [LABEL]      符号指令参照comet计算机说明
```

由上表可知，CASL每条指令由标号(可缺省)、指令码、操作数(可缺省)4栏构成，每一栏的书写

规则如下：

```
标号栏    从第一个字符开始，最多不超过6个字符位置。

指令栏    在无标号时，从第二个字符位置以后的任意字符位置开始。有标号时，标号后面至少
         有一个空白从其后的任意位置开始。

操作数栏　指令码后至少有一个空白，其后到72个字符位置。不能继续到下一行。

注释栏　  行里有分号(;)，其后直到终了作为注释处理(但DC指令里的字符串中的分号除外)。
         此外在第一字符位置为分号或在分号前只有空白的情况下，该行全部作为注释处理，
         在注释栏里，可以书写任何字符。

LABEL    泛指标号，标号最多不超过6个字符，开头必须是英文大写字母，以后可为英文字母或
         数字。
```

用空白表示的栏目里不得写入字符。

2. 伪指令

   a) [LABEL] START [LABEL]

   表示程序的开头，即在程序的开始必须书写。

   操作数栏中的标号是这个程序中定义的标号，它指出该程序的启动地址。在省略的情况下，

   程序从开始执行。标号栏中的标号可以作为其它程序进入该程序的入口。

   b) [LABEL] END

   表示程序的终止，在程序的未尾必须书写。

   c) [LABEL]　DC　常数

   用来指定和存储常数。常数分十进制常数，十六进制常数，地址常数和字符串常数四种。

   标号栏中的标号是代表被指定的十进制常数、十六进制常数、地址常数的存储地址或代表被

   指定的字符串常数的存储区域的第一字的地址。

   十进制常数：DC　n

   用n指定一个十进制(－32768 < n < 65536)，并将n转换成二进制数存储在一个字中。如果n

   超出规定的范围，则将其16位存储起来。对32768－65535的十进制数也可以用负的十进制常数

   表示。

   十六进制常数：DC　#h

   用h指定一个4位十六进制数(0000－FFFF)，并将h对应的二进制数存储在一个字中(在h的前面

   必须写上#)。

   字符串常数：DC　'字符串'

   将字符串中从左开始的每个字符转换成字符数据，并依次把字符数据存储在连续的各字中。在

   字符串中出现('\0', '\n', '\t', '\'', '\\')时，在在前面加转义字符'\\'。
　
   d) [LABEL] DS　区域的字数

   用来保留指定的字数的存储区域的第一个字的地址。区域的字数为零时，存储区域不存在，但

   是标号栏中的标号仍有效，即代表下一字的地址。

3. 宏指令

宏指令是根据事先定义的指令串和操作的信息，生成指定功能的指令串。CASL中有进行输入、输出

及结束程序等宏指令，而没有定义输入，输出符号指令，这类处理由操作系统完成。程序中出现宏

指令时，CASL生成调用操作系统的指令串，但是，生成的指令串的字数不定。执行宏指令时，GR的

内容保持不变而FR的内容不确定。

   a) [LABEL] READ ALABEL

   输入一个十进制数在ALABEL位置的内存中。

   b) [LABEL] WRITE ALABEL

   将内存ALABEL位置的数以十进制形式输出。

   c) [LABEL] IN ALABEL, NLABEL

   从输入装置上输入一个记录，记录中的信息(字符)依次按字符数据的形式被顺序存放在标号为

   ALABEL开始的区域内，已输入的字符个数以二进制的形式存放在标号为NLABEL的字中。NLABEL

   大小不得超过256。

   d) [LABEL] OUT ALABEL, NLABEL

   将存放在从标号ALABEL开始的区域中的字符数据，作为一个记录向输出装置输出，输出的字符

   个数由标号为NLABEL的字中的内容确定。NLABEL大小不得超过256。

   e) [LABEL] EXIT

   序执行的终止，控制返回操作系统。


三、TINY高级语言说明

1. 简要

tiny程序结构很简单：仅是由分号分隔开的语句序列，并且也没有过程声明。所有的变量都是整型

变量，通过对其赋值可以方便地声明变量。它只有两个控制语句：if语句和repeat语句，这两个控

制语句本身也可以包含语句序列。if语句有一个可选的else部分且必须由关键字end结束。除此之外，

read语句和write语句完成输入输出。在花括弧中可以有注释，但注释不能嵌套。

tiny的表达式也局限于布尔表达式和整型算术表达式。布尔表达式由对两个算术表达式的比较组成，

该比较使用 < 和 = 比较算符。算术表达式可以包括整型常数、变量以及4个整型算符+、-、*、/，

此外还有一般的数学属性。布尔表达式只能作为测试出现在控制语句中——而没有布尔型变量。

tiny的记号分为3个典型类型：保留字、特殊符号和“其他”记号。保留字一共有8个，它们的含义类

似。特殊符号共有10种：分别是4种基本的整数运算符号、2种比较符号(小于、等于)，以及括号、分

号和赋值符号。除了赋值符号是两个字符长度外，其余均为一个字符。

其他记号就是数了，它们是一个或者多个数字以及标识符的序列，而标识符又是一个或多个字母序列。

除了记号外，tiny还要遵循以下的词法惯例：注释应放在花括弧{...}中，且不能嵌套；代码是自由格

式；空白格有空格、制表符和新行；最长子串原则后须接识别符号。

2. tiny的EBNF文法

```
   proram        -> stmt-sequence
   stmt-sequence -> statement { ; statement }
   statement     -> if-stmt | repeat-stmt | assign-stmt | read-stmt | write-stmt
   if-stmt       -> if exp then stmt-sequence [ else stmt-sequence ] end
   repeat-stmt   -> repeat stmt-sequence until exp
   assign-stmt   -> identifier := exp
   read-stmt     -> read identifier
   write-stmt    -> write exp
   exp           -> simple-exp [ comparison-op simple-exp ]
   comparison-op -> < | =
   simple-exp    -> term { addop term }
   addop         -> + | -
   term          -> factor { mulop factor }
   mulop         -> * | /
   factor        -> (exp) | number | identifier
```

四、应用实例

1. 一个简单的tiny程序

```
{ sum.tiny 计算 1 + 2 + ... + n 的和 }

read n; { 输入一个整数 }
if 0 < n then { 如果 0 < n 则执行 }
  sum := 0; { 赋值同时声明变量sum }
  repeat { repeat循环 }
    sum := sum + n;
    n := n - 1
  until n = 0; { 当 n = 0 时循环结束 }
  write sum { 输出sum的值 }
end
```

2. 编译tiny程序到casl汇编程序

输入命令：tiny sum
=====================
TINY编译器 到CASL语言
=====================

编译文件 sum.tiny

编译中...

分析结束:)

D:\work\ting-lang>

生成两个文件：sum.list、sum.casl。

sum.list是编译过程中产生的信息。如果没有错误，则sum.casl是得到的汇编程序。

查看sum.list信息

```
TINY编译器: sum.tiny

   1: 
   2: { 计算 1 + 2 + ... + n 的和 }
   3: 
   4: read n; { 输入一个整数 }
	4: 关键字: read
	4: 变量, 名称= n
	4: ;
   5: if 0 < n then { 如果 0 < n 则执行 }
	5: 关键字: if
	5: 数值, 值= 0
	5: <
	5: 变量, 名称= n
	5: 关键字: then
   6:   sum := 0; { 赋值同时声明变量sum }
	6: 变量, 名称= sum
	6: :=
	6: 数值, 值= 0
	6: ;
   7:   repeat { repeat循环 }
	7: 关键字: repeat
   8:     sum := sum + n;
	8: 变量, 名称= sum
	8: :=
	8: 变量, 名称= sum
	8: +
	8: 变量, 名称= n
	8: ;
   9:     n := n - 1
	9: 变量, 名称= n
	9: :=
	9: 变量, 名称= n
	9: -
	9: 数值, 值= 1
  10:   until n = 0; { 当 n = 0 时循环结束 }
	10: 关键字: until
	10: 变量, 名称= n
	10: =
	10: 数值, 值= 0
	10: ;
  11:   write sum { 输出sum的值 }
	11: 关键字: write
	11: 变量, 名称= sum
  12: end
	12: 关键字: end
	12: 文件结束

语法树:

Read读: n
If判断
	运算符: <
		常数: 0
		标号: n
	Assign赋值: sum
		常数: 0
	Repeat循环
		Assign赋值: sum
			运算符: +
				标号: sum
				标号: n
		Assign赋值: n
			运算符: -
				标号: n
				常数: 1
		运算符: =
			标号: n
			常数: 0
	Write写
		标号: sum

符号表：

变量名称  对应标号  初始行号
--------  --------  --------
sum       ABBAAA    6       
n         ABAAAA    4       
```

查看sum.casl信息

```
; ============
; CASL汇编程序
; ============

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
```

3. 将得到的sum.casl汇编程序翻译为comet计算机机器语言

输入命令：casl sum

```
==================
CASL汇编语言编译器
==================

输入文件 sum.casl
输出文件 sum.comet
```

生成sum.comet二进制文件。

4. 运行sum.comet程序

输入命令：comet sum

```
===============
COMET虚拟计算机
===============

100
5050
```

计算得到 1 + 2 + ... + 100 = 5050

5. 调试sum.comet程序

输入命令：comet -debug sum

```
===============
COMET虚拟计算机
===============

调试 （帮助输入 help）...

输入命令: help
命令列表:
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
输入命令: i
显示内存指令
mem[0   ]: JMP  5
输入命令: r
显示寄存器数据
GR[0] =    0    PC =    0
GR[1] =    0    SP = fb00
GR[2] =    0    FR =   01
GR[3] =    0
输入命令: q
退出调试...
```

具体细节用户可以自己尝试 :)
