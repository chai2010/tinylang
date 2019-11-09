# Tiny玩具语言

Tiny语言是[《编译原理及实践》](https://book.douban.com/subject/1088057/)书中定义的玩具语言。本项目基于[`goyacc`](https://github.com/golang/tools/tree/master/cmd/goyacc)工具重新实现Tiny语言。

废话少说，先给个例子：

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

## 语言说明

Tiny程序结构很简单：仅是由分号分隔开的语句序列，并且也没有过程声明。所有的变量都是整型变量，通过对其赋值可以方便地声明变量。它只有两个控制语句：`if`语句和`repeat`语句，这两个控制语句本身也可以包含语句序列。`if`语句有一个可选的`else`部分且必须由关键字`end`结束。除此之外，`read`语句和`write`语句完成输入输出。在花括弧中可以有注释，但注释不能嵌套。

Tiny的表达式也局限于布尔表达式和整型算术表达式。布尔表达式由对两个算术表达式的比较组成，该比较使用`<`和`=`比较算符。算术表达式可以包括整型常数、变量以及4个整型算符`+`、`-`、`*`、`/`，此外还有一般的数学属性。布尔表达式只能作为测试出现在控制语句中——而没有布尔型变量。

Tiny的记号分为3个典型类型：保留字、特殊符号和“其他”记号。保留字一共有8个，它们的含义类似。特殊符号共有10种：分别是4种基本的整数运算符号、2种比较符号(小于、等于)，以及括号、分号和赋值符号。除了赋值符号是两个字符长度外，其余均为一个字符。

其他记号就是数了，它们是一个或者多个数字以及标识符的序列，而标识符又是一个或多个字母序列。除了记号外，tiny还要遵循以下的词法惯例：注释应放在花括弧`{...}`中，且不能嵌套；代码是自由格式；空白格有空格、制表符和新行；最长子串原则后须接识别符号。

## EBNF文法

- 语法图：https://chai2010.cn/tinylang/spec

```ebnf
proram        ::= stmt_sequence
stmt_sequence ::= statement | statement ';' statement
statement     ::= if_stmt | repeat_stmt | assign_stmt | read_stmt | write_stmt
if-stmt       ::= if exp then stmt-sequence | if exp then stmt-sequence else stmt-sequence end
repeat-stmt   ::= repeat stmt-sequence until exp
assign_stmt   ::= identifier ':=' exp
read_stmt     ::= read identifier
write_stmt    ::= write exp
exp           ::= simple_exp | simple_exp '<' simple_exp | simple_exp '=' simple_exp
simple_exp    ::= term | term '+' term | term '-' term
term          ::= factor | factor '*' factor | factor '/' factor
factor        ::= (exp) | number | identifier
```

- [https://en.wikipedia.org/wiki/Extended_Backus-Naur_form](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form)
- https://www.bottlecaps.de/rr/ui

