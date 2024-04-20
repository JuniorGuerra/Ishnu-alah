package models

import (
	"bytes"
	"encoding/binary"
)

type Packet struct {
	// parent       *PacketParser
	payload      *bytes.Buffer
	peerId       uint16
	flags        uint8
	commandCount uint8
	timestamp    uint32
	challenge    uint32
	commands     []*Command
}

func NewPacket(buff []byte) *Packet {
	packet := &Packet{
		payload: bytes.NewBuffer(buff),
	}

	packet.parsePacket()
	return packet
}

func (p *Packet) parsePhotonHeader() {
	binary.Read(p.payload, binary.BigEndian, &p.peerId)
	binary.Read(p.payload, binary.BigEndian, &p.flags)
	binary.Read(p.payload, binary.BigEndian, &p.commandCount)
	binary.Read(p.payload, binary.BigEndian, &p.timestamp)
	binary.Read(p.payload, binary.BigEndian, &p.challenge)

}

func (p *Packet) parsePacket() {
	p.parsePhotonHeader()
	for i := 0; i < int(p.commandCount); i++ {
		command := NewCommand(p.payload)
		p.commands = append(p.commands, command)
	}

}
