# YaccGo
Through Google has tool about yacc named goyacc, But it generate go code can not debug! see the issue 

https://github.com/golang/vscode-go/issues/1674#event-5612030543 

I tried to modify the codes, It is not good readable, and I have no patient to do it, So I write a **YaccGo**

YaccGo is an unstantable and debugable Yacc in Go. . It is written in Go and generates parsers written in Go ,typescript, rust ...

# Quick Start
clone the code
```
make all
# generate the typescript parser code
bin/yaccgo generate typescript examples/exprts.y expts.ts

```
at your `y `file  You should do as follower

# Design

LALR1 Algorithm Base on

https://hassan-ait-kaci.net/pdf/others/p615-deremer.pdf

Lexer inspired from 

https://www.youtube.com/watch?v=HxaD_trXwRE

# RoadMap

1. support language:

- [x] go
- [x] typescript

- [ ] rust


2. DotGraph

   will support Dot Graph by svg

### Contributing

Welcome to contributing, We appreciate your help! please make sure 

* `staticcheck` no error
* codes should has test 

## License


[MIT](LICENSE)
