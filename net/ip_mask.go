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


// format as :   ip/mask    127.0.0.1/8
func (this *IpWithMask) FormatAsIpSlashMask()string{
	ipStr := this.Ip.String()
	ones,_ := this.Mask.Size()
	return fmt.Sprintf("%s/%d",ipStr,ones)
}

