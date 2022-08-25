package main

import (
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
	"reflect"
	"runtime"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

type Player struct {
	Id   int64  `field:"id"`
	Name string `json:"name" field:"name"`
}

func getStructTag(f reflect.StructField, tagName string) string {
	return string(f.Tag.Get(tagName))
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Info("Hello World")

	pPlayer := new(Player)
	reflect.TypeOf(pPlayer).Elem().NumField()
}
