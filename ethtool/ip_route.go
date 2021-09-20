package ethtool

import(
    "fmt"
    "github.com/pigfall/tzzGoUtil/process"
)

const(
    IP_CMD="ip"
    IP_CMD_ROUTE="route"
)

func IpRouteAddWithGateway(dstIp string,devName string,gw string)(error){
    _,errOut,err := process.ExeOutput(
    IP_CMD,
    IP_CMD_ROUTE,
    dstIp,
    "via",
    gw,
    "dev",
        devName,
    )
    if err != nil{
        return fmt.Errorf("%s ,%w",errOut,err)
    }
    return nil
}

func IpRouteAdd(dstIp string,devName string)(error){
    _,errOut,err := process.ExeOutput(
    IP_CMD,
    IP_CMD_ROUTE,
    dstIp,
    "dev",
        devName,
    )
    if err != nil{
        return fmt.Errorf("%s ,%w",errOut,err)
    }
    return nil
}

func IpRouteAddWithSrc(dstIp string,srcIp string,gw string)(error){
    _,errOut,err := process.ExeOutput(
    IP_CMD,
    IP_CMD_ROUTE,
    dstIp,
    "via",
    gw,
    "src",
    srcIp,
    )
    if err != nil{
        return fmt.Errorf("%s ,%w",errOut,err)
    }
    return nil
}
