package net

import(
)



func ForEachAddr(
	do func(dev DeviceIfce,addr IpWithMask)error,
)(error){
	devs,err := ListDevices()
	if err != nil{
		return err
	}
	for _,dev := range devs{
		addrs,err := dev.Addrs()
		if err != nil{
			return err
		}
		for _,addr := range addrs{
			err = do(dev,addr)
			if err != nil{
				return err
			}
		}
	}
	return nil
}

func ListIpV4Addrs()([]IpWithMask,error){
	var addrs = make([]IpWithMask,0)
	err := ForEachAddr(
		func(dev DeviceIfce,addr IpWithMask)error{
			if addr.IsIpV4(){
				addrs = append(addrs,addr)
			}
			return nil
		},
	)
	if err != nil{
		return nil,err
	}
	return addrs,nil
}


func IpSubnetCoincideOrCoinCided(ip *IpWithMask,toCompares []IpWithMask)bool{
	var concide bool
	for _,toCompare := range toCompares {
		concide = (ip.Contains(&toCompare) || toCompare.Contains(ip))
		if concide {
			return true
		}
	}
	return false
}
