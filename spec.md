## Tiny语言的BNF语法

```bnf
program  ->  stmt-sequence
stmt-sequence  ->  stmt-sequence; statement | statement
statement  ->  if-stmt | repeat-stmt | assign-stmt | read-stmt | write-stmt
if-stmt  ->  if exp then stmt-sequence end
             | if exp then stmt-sequence else stmt-sequence end
repeat-stmt  ->  repeat stmt-sequence until exp
assign-stmt  ->  identifier := exp
exp  ->  simple-exp | simple-exp comparison-op simple-exp
comparison-op  ->  < | =
simple-exp  ->  simple-exp addop term | term
addop  ->  + | -
term  ->  term mulop factor | factor
mulop  ->  * | /
factor  ->  (exp) | numver | identifier
```

注释采用Go语言风格的注释。
