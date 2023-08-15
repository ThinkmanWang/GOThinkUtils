package main

import (
	"fmt"
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/qqwry"
)

var g_qqwry *qqwry.QQwry

func main() {
	g_qqwry = qqwry.NewQQwry("./qqwry.dat")

	ret := g_qqwry.Find("182.139.183.98")
	fmt.Println(ret)

	ret = g_qqwry.Find("121.13.218.229")
	fmt.Println(ret)

	ret = g_qqwry.Find("223.104.3.15")
	fmt.Println(ret)

	ret = g_qqwry.Find("116.233.109.130")
	fmt.Println(ret)

	ret = g_qqwry.Find("61.128.128.68")
	fmt.Println(ret)

	ret = g_qqwry.Find("223.104.7.98")
	fmt.Println(ret)
}
