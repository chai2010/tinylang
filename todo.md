# TODO列表

1. 改造为Go语言版本
2. 包装为库, 可以从外部导入
3. 重新设计虚拟机, 支持扩展函数
4. 用json表示二进制格式, 并保留对应的CASL指令
5. 先转到`<stdint.h>`类型，然后逐步重构
6. COMET增加一个`SYSCALL`指令，GR0对应调用编号，GR1、GR2、GR3为参数，返回结果在GR0～GR4中
7. 增加单元测试

系统调用函数格式：`func syscall(gr0, gr1, gr2, gr3 uint16) (gr0, gr1, gr2, gr3 uint16)`
