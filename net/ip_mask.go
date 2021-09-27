package net

import(
		"fmt"
		"strconv"
		"net"
		"strings"
)


type IpWithMask struct{
	Ip net.IP
	Mask net.IPMask
}

func (this *IpWithMask) IsIpV4()bool{
	return IsIpv4(this.Ip)
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

