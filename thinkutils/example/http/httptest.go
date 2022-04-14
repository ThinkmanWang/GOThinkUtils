package main

import (
	"github.com/ThinkmanWang/github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"
	"github.com/levigross/grequests"
	"runtime"
)

var (
	log *logger.LocalLogger = logger.DefaultLogger()
)

func simpleGet() {
	opt := grequests.RequestOptions{Headers: map[string]string{
		"a": "1",
	}, Params: map[string]string{
		"aaa": "2",
	}}
	resp, err := grequests.Get("http://httpbin.org/get", &opt)
	if nil != err {
		return
	}

	log.Info(resp.String())
}

func simplePost() {
	opt := grequests.RequestOptions{Headers: map[string]string{
		"a": "1",
	}, Params: map[string]string{
		"aaa": "2",
	}, Data: map[string]string{
		"bbb": "3",
	}}
	resp, err := grequests.Post("http://httpbin.org/post?b=2", &opt)
	if nil != err {
		return
	}

	log.Info(resp.String())
}

type User struct {
	Name string `json:"name"`
	Age  uint8  `json:"age"`
}

func postJson() {
	opt := grequests.RequestOptions{
		JSON: User{Name: "abc", Age: 1},
	}
	resp, err := grequests.Post("http://httpbin.org/post", &opt)
	if nil != err {
		return
	}

	log.Info(resp.String())
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	simpleGet()
	simplePost()
	postJson()
}
