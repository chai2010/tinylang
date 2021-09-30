
# Tiny玩具语言(Go语言版)

Tiny语言是[《编译原理及实践》](https://book.douban.com/subject/1088057/)书中定义的玩具语言。

这里是Go语言实现(注释采用Go语言风格)。

实现原理：

- [COMET虚拟计算机说明](./comet/README.md)
- [COMET虚拟机的设计与实现.pdf](./docs/comet-vm.pdf)

## 例子

以下的例子计算1到n之和：

```
// sum = 1 + 2 + ... + n

read n;
if 0 < n then
  sum := 0;
  repeat
    sum := sum + n;
    n := n - 1
  until n = 0;
  write sum
end
```

运行tiny程序:

```
$ tinylang sum.tiny 
100
5050
```

也可以通过`-ast`查看生成的语法树，通过`-casl`查看生成的CASL汇编程序，或通过`-debug`调试执行：

![](./_docs/images/tiny-demo.cast.gif)

## WebAssembly 支持

切换到 `wasm` 目录, 通过以下命令执行 `hello.tiny`:

```
$ go run main.go hello.tiny
READ: 100
5050
```

`READ` 表示需要输入一个整数, 这里的 100 表示计算 1 到 100 的和, 结果是 5050.

通过 `-o` 参数输出 WebAssembly 二进制模块:

```
$ go run main.go -o hello.wasm hello.tiny
```

实用 `wasm2wat` 将二进制格式转换为文本格式:

```
$ wasm2wat hello.wasm > hello.wast
```

输出的 `hello.wast` 内容如下:

```wasm
(module $tinylang
  (type (;0;) (func (result i32)))
  (type (;1;) (func (param i32)))
  (type (;2;) (func))
  (import "env" "__tiny_read" (func (;0;) (type 0)))
  (import "env" "__tiny_write" (func (;1;) (type 1)))
  (func $_start (type 2)
    call 0
    set_global 1
    i32.const 0
    get_global 1
    i32.lt_s
    if  ;; label = @1
      i32.const 0
      set_global 2
      loop  ;; label = @2
        get_global 2
        get_global 1
        i32.add
        set_global 2
        get_global 1
        i32.const 1
        i32.sub
        set_global 1
        get_global 1
        i32.const 0
        i32.eq
        i32.eqz
        br_if 0 (;@2;)
      end
      get_global 2
      call 1
    end)
  (memory (;0;) 1 1)
  (global (;0;) i32 (i32.const 0))
  (global (;1;) (mut i32) (i32.const 0))
  (global (;2;) (mut i32) (i32.const 0))
  (export "memory" (memory 0))
  (export "_start" (func $_start)))
```

如果要单独执行输出的 WebAssembly 模块, 需要注入 `__tiny_read` 和 `__tiny_write` 函数。

## 版权

保留所有权利。
