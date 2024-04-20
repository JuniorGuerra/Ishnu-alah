package models

type PacketParser struct {
	PacketChan chan *Packet
}

func NewPacketParser() *PacketParser {
	return &PacketParser{
		PacketChan: make(chan *Packet),
	}
}

func (p *PacketParser) Handle(buff []byte) {
	// En Go, simplemente enviamos el paquete al channel.
	p.PacketChan <- NewPacket(buff)
}
