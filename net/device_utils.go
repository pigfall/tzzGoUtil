package net



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


func IpSubnetCoincide(ip IpWithMask,toCompares []IpWithMask)bool{
	panic("UNIMPL")
}
