package net

import(
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
    "math"
    "strings"
    "net"
    "fmt"
    "strconv"
)


func IsIpv4(ip net.IP)(bool){
	return ip.To4()!= nil
}

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

type Ipv4Packet struct{
    packet  gopacket.Packet
    ipLayer *layers.IPv4
}

func ParseIpv4Packet(raw []byte)(*Ipv4Packet,error){
    ipPacket := &Ipv4Packet{}
    pac:= gopacket.NewPacket(raw,layers.LayerTypeIPv4,gopacket.Default)
    if pac.ErrorLayer() != nil{
        return nil,pac.ErrorLayer().Error()
    }
    ipPacket.packet = pac
    linkLayer := pac.LinkLayer()
    if linkLayer.LayerType() != layers.LayerTypeIPv4{
        return nil,fmt.Errorf("not a ipv4 packet")
    }
    
    ipLayer := pac.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
    ipPacket.ipLayer = ipLayer
    return ipPacket,nil
}
