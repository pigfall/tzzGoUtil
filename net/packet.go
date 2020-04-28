package net

import(
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcap"
    "github.com/Peanuttown/tzzGoUtil/log"
)

func Listen(){
    devices,err := pcap.FindAllDevs()
    if err != nil{
        log.Debug(err)
        return
    }
    log.Debugf("find devices:%d\n",len(devices))
    for device := range devices{
        log.Debugf("%#v\n",device)
    }
    handler,err:=pcap.OpenLive("ethernet_32768",65536,true,pcap.BlockForever)
    if err != nil{
        log.Debug(err)
        return
    }
    src := gopacket.NewPacketSource(handler,layers.LayerTypeEthernet)
    packetCh := src.Packets()
    log.Debug("over")
    for{
        log.Debug((<-packetCh).String())
    }
}
