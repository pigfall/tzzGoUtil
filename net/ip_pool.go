package net

import(
	"sync"
	"fmt"
		"strings"
)



type IpPool struct{
	beenUsed map[string]*IpWithMask
	baseIpNet *IpWithMask
	l sync.Mutex
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


func (this *IpPool) Take()(*IpWithMask,error){
	this.l.Lock()
	defer this.l.Unlock()
	var take *IpWithMask
	this.baseIpNet.ForEachIpInThisCidr(
		func(ipNet *IpWithMask)(stop bool,err error){
			if used  := this.beenUsed[ipNet.FormatAsIpSlashMask()];used != nil{
				return false,nil
			}
			take =ipNet
			return true,nil
		},
	)
	if take == nil{
		return nil,fmt.Errorf("No Aviable ip is this ip pool")
	}

	return take,nil
}

func (this *IpPool) Release(ipNet *IpWithMask){
	this.l.Lock()
	defer this.l.Unlock()
	delete(this.beenUsed,ipNet.FormatAsIpSlashMask())
}
