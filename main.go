package main

import (
	"ishnu-alah/models"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const (
	//	device = `\Device\NPF_{3B4C1643-F01E-4542-9FC8-CE8A27A755D2}`
	device = "wlo1"
	filter = "udp and (dst port 5056 or src port 5056)"
)

func main() {

	handle, err := pcap.OpenLive(device, 1600, true, pcap.BlockForever)

	if err != nil {
		panic(err)
	}

	if err := handle.SetBPFFilter(filter); err != nil {
		panic(err)
	}
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {

		// b := packet.ApplicationLayer()
		processPacket(packet)

	}

}

func handlePayload(payload []byte) {

	packetPharser := models.NewPacketParser()

	go func() {
		for packet := range packetPharser.PacketChan {
			// fmt.Println("canal del packete", packet)
			if packet != nil {
				continue
			}
		}
	}()

	packetPharser.Handle(payload)

}

func processPacket(packet gopacket.Packet) {
	// if ipv4Layer := packet.Layer(layers.LayerTypeIPv4); ipv4Layer != nil {
	// ipv4, _ := ipv4Layer.(*layers.IPv4)
	if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {

		udp, _ := udpLayer.(*layers.UDP)
		handlePayload(udp.Payload)
	}

}
