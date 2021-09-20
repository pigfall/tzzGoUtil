package net

import(
		"fmt"
		"net"
)


type IpWithMask struct{
	Ip net.IP
	Mask net.IPMask
}

// format as :   ip/mask    127.0.0.1/8
func (this *IpWithMask) FormatAsIpSlashMask()string{
	ipStr := this.Ip.String()
	ones,_ := this.Mask.Size()
	return fmt.Sprintf("%s/%s",ipStr,ones)
}

