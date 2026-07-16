package qlipnet

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
)

// Type of message being sent (metadata).
type MsgType byte

const InitMsgType MsgType = 1          // Initialization handshake
const QlipBoardTextMsgType MsgType = 2 // Clipboard text data broadcast
const QlipBoardImgMsgType MsgType = 3  // Clipboard image data broadcast
const CloseMsgType MsgType = 4         // Close gracefully handshake

// Communication is sent in a standardized Packet
// protocol. Each packet will contain the type header.
// If the type is QlipBoardMsgType then the packet will
// have a length and data bytes.
type Packet struct {
	Type MsgType
	Len  uint64
	Data []byte
}

// Serialize a packet into bytes so that it
// can be sent over the a network connection.
func (p Packet) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.WriteByte(byte(p.Type))

	err := binary.Write(buf, binary.BigEndian, p.Len)
	if err != nil {
		return nil, err
	}

	buf.Write(p.Data)

	return buf.Bytes(), nil
}

// Read a full packet from a connection.
func readPacket(conn net.Conn) (Packet, error) {
	var p Packet

	// Read packet type
	var typ [1]byte
	if _, err := io.ReadFull(conn, typ[:]); err != nil {
		return p, err
	}
	p.Type = MsgType(typ[0])

	// Read payload length
	if err := binary.Read(conn, binary.BigEndian, &p.Len); err != nil {
		return p, err
	}

	// Read payload
	p.Data = make([]byte, p.Len)
	if _, err := io.ReadFull(conn, p.Data); err != nil {
		return p, err
	}

	return p, nil
}

// Serialize and then write a packet to a connection.
func writePacket(conn net.Conn, p Packet) error {
	buf, err := p.Serialize()
	if err != nil {
		return err
	}

	_, err = conn.Write(buf)
	return err
}
