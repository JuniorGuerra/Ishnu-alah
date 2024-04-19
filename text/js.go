package text

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// PhotonPacketParser es la estructura que hereda de EventEmitter en JS.
type PhotonPacketParser struct {
	// En Go, usamos channels en lugar de EventEmitter.
	packetChan chan *PhotonPacket
}

// NewPhotonPacketParser crea una nueva instancia de PhotonPacketParser.
func NewPhotonPacketParser() *PhotonPacketParser {
	return &PhotonPacketParser{
		packetChan: make(chan *PhotonPacket),
	}
}

// Handle toma un buffer y emite un paquete Photon.
func (p *PhotonPacketParser) Handle(buff []byte) {
	// En Go, simplemente enviamos el paquete al channel.
	p.packetChan <- NewPhotonPacket(buff)
}

// PhotonPacket representa un paquete Photon.
type PhotonPacket struct {
	parent       *PhotonPacketParser
	payload      *bytes.Buffer
	peerId       uint16
	flags        uint8
	commandCount uint8
	timestamp    uint32
	challenge    uint32
	commands     []*PhotonCommand
}

// NewPhotonPacket crea una nueva instancia de PhotonPacket.
func NewPhotonPacket(buff []byte) *PhotonPacket {
	packet := &PhotonPacket{
		payload: bytes.NewBuffer(buff),
	}
	packet.parsePacket()
	return packet
}

// parsePhotonHeader analiza la cabecera del paquete Photon.
func (p *PhotonPacket) parsePhotonHeader() {
	binary.Read(p.payload, binary.BigEndian, &p.peerId)
	binary.Read(p.payload, binary.BigEndian, &p.flags)
	binary.Read(p.payload, binary.BigEndian, &p.commandCount)
	binary.Read(p.payload, binary.BigEndian, &p.timestamp)
	binary.Read(p.payload, binary.BigEndian, &p.challenge)
}

// parsePacket analiza el paquete completo.
func (p *PhotonPacket) parsePacket() {
	p.parsePhotonHeader()
	for i := 0; i < int(p.commandCount); i++ {
		command := NewPhotonCommand(p.payload)
		p.commands = append(p.commands, command)
	}
}

// PhotonCommand representa un comando Photon.
type PhotonCommand struct {
	payload        *bytes.Buffer
	commandType    uint8
	channelId      uint8
	commandFlags   uint8
	commandLength  uint32
	sequenceNumber uint32
	messageType    uint8
	data           interface{}
}

// NewPhotonCommand crea una nueva instancia de PhotonCommand.
func NewPhotonCommand(payload *bytes.Buffer) *PhotonCommand {
	command := &PhotonCommand{
		payload: payload,
	}
	command.parseCommand()
	return command
}

// parseCommandHeader analiza la cabecera del comando.
func (c *PhotonCommand) parseCommandHeader() {
	binary.Read(c.payload, binary.BigEndian, &c.commandType)
	binary.Read(c.payload, binary.BigEndian, &c.channelId)
	binary.Read(c.payload, binary.BigEndian, &c.commandFlags)
	c.payload.Next(1) // Saltar 1 byte.
	binary.Read(c.payload, binary.BigEndian, &c.commandLength)
	binary.Read(c.payload, binary.BigEndian, &c.sequenceNumber)
}

// parseCommand analiza el comando completo.
func (c *PhotonCommand) parseCommand() {
	c.parseCommandHeader()
	// Aquí iría la lógica para analizar los diferentes tipos de comandos.
}

func main() {
	// Ejemplo de cómo se usaría.
	parser := NewPhotonPacketParser()
	go func() {
		for packet := range parser.packetChan {
			fmt.Println(packet)
			// Procesar el paquete recibido.
		}
	}()

	// Simular la recepción de un paquete.
	parser.Handle([]byte{0x01, 0x02, 0x03})
}
