package main

import "GOThinkUtils/thinkutils/logger"

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

type MyStruct struct {
	Name string
}

func (this *MyStruct) Func1() {
	log.Info("%p %s", this, this.Name)
}

func (this *MyStruct) Func2() {
	this.Name = "456"
}

func (this *MyStruct) Func3() {
	this.Name = "789"
}

func main() {
	log.Info("FXXK")

	p := MyStruct{Name: "123"}
	log.Info("%p", &p)

	p.Func2()
	log.Info(p.Name)

	p.Func3()
	log.Info(p.Name)

}
