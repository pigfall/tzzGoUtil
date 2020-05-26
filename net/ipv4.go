package net

import(
    "math"
    "strings"
    "net"
    "fmt"
    "strconv"
)


func ParseIPv4(ip string)(net.IP,error){
    elems := strings.Split(ip,".")
    if len(elems) != 4{
        return nil,fmt.Errorf("parse failed,elem num not 4:%s",ip)
    }
    ret := make([]byte,4)
    for i,v:= range elems{
        value,err := strconv.ParseInt(v,10,64)
        if err != nil{
            return nil,fmt.Errorf("parse ipv4 elem to number failed:%s",ip)
        }
        if math.MaxUint8 < value{
            return nil,fmt.Errorf("parse ipv4 failed,elem value > 255:%s",ip)
        }
        ret[i] = byte(value)
    }
    return ret,nil
}
