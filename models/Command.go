package models

import (
	"bytes"
	"encoding/binary"
)

type Command struct {
	parent         *PacketParser
	payload        *bytes.Buffer
	commandType    uint8
	channelId      uint8
	commandFlags   uint8
	commandLength  uint32
	sequenceNumber uint32
	messageType    uint8
	data           interface{}
}

func NewCommand(payload *bytes.Buffer) *Command {
	command := &Command{
		payload: payload,
	}
	command.parseCommand()
	return command
}

func (c *Command) parseCommand() {
	c.parseCommandHeader()
	switch c.commandType {
	case 7:
		c.payload.Next(4)
		c.payload = bytes.NewBuffer(c.payload.Bytes())
	case 6:
		c.ParseReliableCommand()

	case 4:
		break
	}
}

func (c *Command) parseCommandHeader() {
	binary.Read(c.payload, binary.BigEndian, &c.commandType)
	binary.Read(c.payload, binary.BigEndian, &c.channelId)
	binary.Read(c.payload, binary.BigEndian, &c.commandFlags)
	c.payload.Next(1) // Saltar 1 byte.
	binary.Read(c.payload, binary.BigEndian, &c.commandLength)
	binary.Read(c.payload, binary.BigEndian, &c.sequenceNumber)
}

func (c *Command) ParseReliableCommand() error {

	c.payload.Next(1)

	messageType, err := c.payload.ReadByte()
	if err != nil {
		return err
	}
	c.messageType = messageType

	newPayload := bytes.NewBuffer(c.payload.Bytes()[2:])

	switch c.messageType {
	case 2:
		data, err := DeserializeOperationRequest(newPayload)
		if err != nil {
			return err
		}
		c.data = data

		c.parent.Emit("request", data)
	case 3:
		data, err := DeserializeOperationResponse(newPayload)
		if err != nil {
			return err
		}
		c.data = data

		c.parent.Emit("response", data)
	case 4:
		data, err := DeserializeEventData(newPayload)
		if err != nil {
			return err
		}
		c.data = data

		c.parent.Emit("event", data)
	}

	return nil
}

func (p *PacketParser) Emit(eventType string, data interface{}) {

}
