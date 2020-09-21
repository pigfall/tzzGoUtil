package net

import (
    "fmt"
    "testing"
    "github.com/Peanuttown/tzzGoUtil/process"
)

// 测试 当 tun,tap 设备消失时 ，与之相关的 路由规则也会消失
func TestTunRoute(t *testing.T) {
    var devName = "testDev"
    _,err  :=NewTun(devName)
    if err != nil{
        t.Fatal(err)
    }
    err = DevUp(devName)
    if err != nil{
        t.Fatal(err)
    }
    _,errOut,err := process.ExeOutput("ip","route","add","192.16.2.0/24","dev",devName)
    if err != nil{
        t.Fatal(fmt.Errorf("%s, %w",errOut,err))
    }
}
