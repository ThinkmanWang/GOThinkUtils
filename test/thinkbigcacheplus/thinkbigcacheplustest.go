package main

import (
	"encoding/gob"
	"strings"

	"github.com/ThinkmanWang/GOThinkUtils/thinkutils"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

// MyStruct 是示例数据类型。gob 只编码导出字段，所以字段必须大写。
type MyStruct struct {
	Name string
	Age  int
	Tags []string
}

func init() {
	gob.Register(MyStruct{})
}

// 各分区的类型化句柄。通常在各模块 init 里注册好，业务代码直接用句柄读写。
// RegType[T] 内部用 gob 自动生成编解码，无需手写 ToByte/FromByte。
var (
	userPartition = thinkutils.ThinkBigCachePlusRegType[MyStruct](thinkutils.ThinkBigCachePlusInstance(), "user")
	strPartition  = thinkutils.ThinkBigCachePlusRegType[string](thinkutils.ThinkBigCachePlusInstance(), "commonStr")
)

func testSave() {
	_ = thinkutils.ThinkBigCachePlusInstance().Start()

	// Set 编译期锁类型：传错类型直接编译不过。
	_ = userPartition.Set("u1", MyStruct{Name: "Thinkman", Age: 18, Tags: []string{"go", "cache"}})
	_ = strPartition.Set("hello", "world123")

	_ = thinkutils.ThinkBigCachePlusInstance().SaveToDisk()
}

func testLoad() {
	_ = thinkutils.ThinkBigCachePlusInstance().Start()

	// Get 直接返回具体类型 T，无需 .(T) 断言。
	if data, err := userPartition.Get("u1"); err != nil {
		log.Error("Failed to get u1 from cache : %s", err.Error())
	} else {
		log.Info("user: name=%s age=%d tags=[%s]", data.Name, data.Age, strings.Join(data.Tags, ","))
	}

	if data, err := strPartition.Get("hello"); err != nil {
		log.Error("Failed to get hello from cache : %s", err.Error())
	} else {
		log.Info("str: %s", data)
	}
}

func main() {
	testSave()
	//testLoad()
}
