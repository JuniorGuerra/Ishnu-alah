package models

type PacketParser struct {
	packetChan chan *Packet
}

func NewPacketParser() *PacketParser {
	return &PacketParser{
		packetChan: make(chan *Packet),
	}
}

func (p *PacketParser) Handle(buff []byte) {
	// En Go, simplemente enviamos el paquete al channel.
	p.packetChan <- NewPacket(buff)
}
