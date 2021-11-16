package net

import(
		"fmt"
		"strconv"
		"net"
		"strings"
)

//eg: 127.0.0.1:80
// {
type IpPortFormat string

func IpPortFormatFromIpPort(ip net.IP,port int)IpPortFormat{
	return IpPortFormat(fmt.Sprintf("%s:%d",ip.String(),port))
}
// }

type IpPort struct{
	IP net.IP
	Port int
}

// format x.x.x.x:8
func IpPortFromString(format string)(*IpPort,error){
	lastColumnIndex := strings.LastIndex(format,":")
	if lastColumnIndex<0{
		return nil,fmt.Errorf("Invalid ipPort format %s",format)
	}
	ip := format[:lastColumnIndex]
	ipObj:=net.ParseIP(ip)
	if ipObj == nil{
		return nil,fmt.Errorf("Invalid ipPort format %s",format)
	}
	port := format[lastColumnIndex+1:]
	portInt,err := strconv.ParseUint(port,10,64)
	if err != nil{
		return nil,fmt.Errorf("Invalid ipPort format %s",format)
	}

	return &IpPort{
		IP:ipObj,
		Port:int(portInt),
	},nil


}

func (this *IpPort) ToIpPortFormat()IpPortFormat{
	return IpPortFormat(fmt.Sprintf("%s:%d",this.IP.String(),this.Port))
}

func (this *IpPort) ToString() string{
	return fmt.Sprintf("%s:%v",this.IP.String(),this.Port)
}


//eg: 127.0.0.1
type IpFormat string

//eg: 127.0.0.1/8
type IpNetFormat string


type IpWithMask struct{
	Ip net.IP
	Mask net.IPMask
}


func (this *IpWithMask) ToString()(string){
	return string(this.ToIpNetFormat())

}

func (this *IpWithMask) ToIpNetFormat()(IpNetFormat){
	return IpNetFormat(this.FormatAsIpSlashMask())
}

func (this *IpWithMask) IsIpV4()bool{
	return IsIpv4(this.Ip)
}

func (this *IpWithMask) IpFormat()IpFormat{
	return IpFormat(fmt.Sprintf("%s",this.Ip.String()))
}

func (this *IpWithMask) ToIpNet()(*net.IPNet){
	return &net.IPNet{
		IP:this.Ip,
		Mask:this.Mask,
	}
}

func (this *IpWithMask) Contains(other *IpWithMask)bool{
	return this.ToIpNet().Contains(other.Ip)
}

func (this *IpWithMask) String() (string){
	onesCount,_:= this.Mask.Size()
	return fmt.Sprintf("%s/%d",this.Ip.String(),onesCount)
}

func FromIpColonMask(format string)(*IpWithMask,error){
	return  fromIpSplitSymbolMask(format,":")
}

func fromIpSplitSymbolMask(format,splitSymbol string)(*IpWithMask,error){
	slashIndex := strings.Index(format,splitSymbol)
	if slashIndex  < 0 { // not found slash
		return nil,fmt.Errorf("Not found spilit symbol %s",splitSymbol)
	}
	ip := net.ParseIP(format[:slashIndex])
	if ip == nil {
		return nil,fmt.Errorf("Parse ip failed: ip format is invalid")
	}
	if slashIndex+1 == len(format){
		return nil, fmt.Errorf("Invalid format, no mask after %s",splitSymbol)
	}
	maskInt,err := strconv.ParseInt(format[slashIndex+1:],10,64)
	if err != nil{
		return nil,fmt.Errorf("Parse mask to int failed, invalid mask")
	}
	// < TODO
	if len(ip)== 4{
		mask := net.CIDRMask(int(maskInt),32)
		return &IpWithMask{
			Ip:ip,
			Mask:mask,
		},nil
	}
	// >
	panic("only supported ipv4")
}

// from format as '127.0.0.1/8' => IpWithMask
func FromIpSlashMask(format string)(*IpWithMask,error){
	slashIndex := strings.Index(format,"/")
	if slashIndex  < 0 { // not found slash
		return nil,fmt.Errorf("Not found slash")
	}
	ip := net.ParseIP(format[:slashIndex])
	if ip == nil {
		return nil,fmt.Errorf("Parse ip failed: ip format is invalid")
	}
	if slashIndex+1 == len(format){
		return nil, fmt.Errorf("Invalid format, no mask after slash")
	}
	maskInt,err := strconv.ParseInt(format[slashIndex+1:],10,64)
	if err != nil{
		return nil,fmt.Errorf("Parse mask to int failed, invalid mask")
	}
	// < TODO
	if IsIpv4(ip){

		mask := net.CIDRMask(int(maskInt),32)
		return &IpWithMask{
			Ip:ip,
			Mask:mask,
		},nil
	}
	// >
	panic("only supported ipv4")
}


// format as :   ip/mask    127.0.0.1/8
func (this *IpWithMask) FormatAsIpSlashMask()string{
	ipStr := this.Ip.String()
	ones,_ := this.Mask.Size()
	return fmt.Sprintf("%s/%d",ipStr,ones)
}


func (this *IpWithMask) BaseIpNet() *IpWithMask{
	// { TODO
	if !IsIpv4(this.Ip){
		panic("TODO, only supported ipv4")
	}
	_,cidr,err := net.ParseCIDR(this.FormatAsIpSlashMask())
	if err != nil{
		panic(err)
	}
	return IpWithMaskFromIpNet(cidr)
}

func IpWithMaskFromIpNet(ipNet *net.IPNet)*IpWithMask{
	if ipNet == nil{
		return nil
	}
	return &IpWithMask{
		Ip:ipNet.IP,
		Mask:ipNet.Mask,
	}
}

func (this *IpWithMask) ForEachIpInThisCidr(do func(ipWithMask *IpWithMask)(stop bool,err error))error{
	_ , ipnet, err := net.ParseCIDR(this.FormatAsIpSlashMask())
	if err != nil{
		return err
	}
	baseIp := make([]byte,len(ipnet.IP))
	copy(baseIp,ipnet.IP)
	baseIpNet := &net.IPNet{
		IP:baseIp,
		Mask:ipnet.Mask,
	}
	incIp(ipnet.IP)
	if !baseIpNet.Contains(ipnet.IP){
		return nil
	}

	for {
		ipCopy := make([]byte,len(ipnet.IP))
		copy(ipCopy,ipnet.IP)
		incIp(ipnet.IP)
		if !baseIpNet.Contains(ipnet.IP){
			break
		}
		stop,err := do(IpWithMaskFromIpNet(&net.IPNet{
			IP:ipCopy,
			Mask:ipnet.Mask,
		}))
		if err != nil{
			return err
		}
		if stop {
			break
		}
	}

	return nil
}

//  http://play.golang.org/p/m8TNTtygK0
func incIp(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}



// 255.255.0.0
func MaskFormatTo255(mask net.IPMask)string{
	var elems = make([]string,0,len(mask))
	for _,m := range mask{
		elems = append(elems,fmt.Sprintf("%d",int(m)))
	}
	return strings.Join(elems,".")
}

