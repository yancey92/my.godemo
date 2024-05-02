* init_demo 包中的测试验证了 **testing.Init() 中不会执行 flag.Parse()**
* flag_demo 包中测试验证了 **flag.Parse() 会被 go test 执行**
* 综上所述：go test 执行测试时，flag.Parse() 会被 go test 执行，但是不会在 testing.Init() 中被执行

具体 flag.Parse() 会被什么时候执行，那就要看代码运行过程中 flag.Parse() 出现的位置了。
测试时：
    先解析包中的init（解析顺序？？）
    再解析  testing.Init()
    再执行 flag.Parse()
    再执行具体的测试函数
编译时：
    从 main() 包开始，
    先解析包中的init（解析顺序？？）


golang 内置的 "flag" 包的特点：
		1.flag.Parse() 只能执行一次，重复执行则报错
		2.flag.Parse() 应该在获取所有参数之后调用；如果在其之后又执行类似 flag.Int() 这样的代码，会报错 "flag provided but not defined: -xxxxx"
		所以，建议使用 "github.com/spf13/pflag" 包。