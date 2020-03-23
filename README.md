- *赞助 BTC: 1Cbd6oGAUUyBi7X7MaR4np4nTmQZXVgkCW*
- *赞助 ETH: 0x623A3C3a72186A6336C79b18Ac1eD36e1c71A8a6*

----

# Tiny玩具语言

Tiny语言是[《编译原理及实践》](https://book.douban.com/subject/1088057/)书中定义的玩具语言。

实现原理：

- [COMET虚拟计算机说明](./comet/README.md)
- [CASL汇编器的设计与实现.pdf](./docs/casl-assembler.pdf)
- [COMET虚拟机的设计与实现.pdf](./docs/comet-vm.pdf)

## 例子

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

第一步：编译tiny到casl汇编程序：

```
$ go run ./tiny sum.tiny
=====================
TINY编译器 到CASL语言
=====================

编译文件 sum.tiny

编译中...

分析结束 :)
```

第二步：编译casl汇编到comet二进制格式：

```
$ go run ./casl sum.casl
==================
CASL汇编语言编译器
==================

输入文件 sum.casl
输出文件 sum.comet
```

第三步：虚拟机运行程序：

```
$ go run . sum.comet
100
5050
```

## 补充说明

目前的实现版本是基于2005年实现的C语言版本。当时版本的C语言组织方式并不合理，稍后会逐步改造为Go语言实现。

## 版权

保留所有权利。
