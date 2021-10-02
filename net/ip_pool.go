package net

import(
	"fmt"
		"strings"
)



type IpPool struct{
	beenUsed map[string]*IpWithMask
	baseIpNet *IpWithMask
}

func NewIpPool(ipNet *IpWithMask,hasBeenUsed []*IpWithMask) (*IpPool,error){
	var invalidUsedIpNets = make([]*IpWithMask,0)
	for _,used := range hasBeenUsed{
		if !ipNet.Contains(used){
			invalidUsedIpNets  = append(invalidUsedIpNets,used)
		}
	}

	if len(invalidUsedIpNets) > 0 {
		var invalidUsesIpNetsStr = make([]string,0,len(invalidUsedIpNets))
		for _,v := range invalidUsedIpNets {
			invalidUsesIpNetsStr  = append(invalidUsesIpNetsStr, v.FormatAsIpSlashMask())
		}
		return nil,fmt.Errorf(" %s not in ipNet %s",strings.Join(invalidUsesIpNetsStr,","),ipNet.FormatAsIpSlashMask())
	}

	beenUsed := make(map[string]*IpWithMask)
	for _,v := range hasBeenUsed {
		beenUsed[v.FormatAsIpSlashMask()] = v
	}

	return &IpPool{
		beenUsed:beenUsed,
		baseIpNet:ipNet,
	},nil
}

