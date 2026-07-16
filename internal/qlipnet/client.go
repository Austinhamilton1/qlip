package qlipnet

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/Austinhamilton1/qlip/internal/qlipboard"
)

type Client struct {
	conn net.Conn
	quit chan bool
}

func NewClient() *Client {
	return &Client{
		conn: nil,
		quit: make(chan bool),
	}
}

// Connect to the qlipnet server and initalize a handshake.
func (c *Client) Connect(ipAddr string, port int) error {
	connStr := fmt.Sprintf("%v:%v", ipAddr, port)
	conn, err := net.Dial("tcp", connStr)
	if err != nil {
		return err
	}

	p := Packet{
		Type: InitMsgType,
		Len:  0,
		Data: nil,
	}

	err = writePacket(conn, p)
	if err != nil {
		return err
	}

	p, err = readPacket(conn)
	if err != nil {
		return err
	}

	if p.Type != InitMsgType {
		errorMsg := fmt.Sprintf("invalid response from server: %v", p.Type)
		return errors.New(errorMsg)
	}

	c.conn = conn

	return nil
}

// Main event loop. This event loop works two fold. First,
// the loop watches for changes on the clipboard. If any are
// found, then it sends the data to the serve. Second, the loop
// listens for messages coming from the server and synchronizes
// the incoming messages to the local clipboard.
func (c *Client) Loop() {
	qlipboard.Init()

	update := qlipboard.Watch(c.quit)
	synchronizer := qlipboard.Sync(c.quit)

	go func(conn net.Conn) {
		for {
			p, err := readPacket(conn)
			if err != nil {
				log.Printf("could not read packet: %v", err)
				continue
			}

			switch p.Type {
			case CloseMsgType:
				close(c.quit)
				return
			case QlipBoardTextMsgType:
				synchronizer <- qlipboard.Data{
					DataType: qlipboard.Text,
					Bytes:    p.Data,
				}
			case QlipBoardImgMsgType:
				synchronizer <- qlipboard.Data{
					DataType: qlipboard.Image,
					Bytes:    p.Data,
				}
			default:
				log.Printf("invalid packet type %d", p.Type)
			}
		}
	}(c.conn)

	// Watch clipboard and if there is a change, send the
	// data to the server
	for data := range update {
		switch data.DataType {
		case qlipboard.Text:
			err := writePacket(
				c.conn,
				Packet{
					Type: QlipBoardTextMsgType,
					Len:  uint64(len(data.Bytes)),
					Data: data.Bytes,
				},
			)
			if err != nil {
				log.Printf("missed clipboard text data: %v", err)
			}
		case qlipboard.Image:
			err := writePacket(
				c.conn,
				Packet{
					Type: QlipBoardImgMsgType,
					Len:  uint64(len(data.Bytes)),
					Data: data.Bytes,
				},
			)
			if err != nil {
				log.Printf("missed clipboard image data: %v", err)
			}
		}
	}

	err := writePacket(
		c.conn,
		Packet{
			Type: CloseMsgType,
			Len:  0,
			Data: nil,
		},
	)
	if err != nil {
		log.Printf("could not send close message: %v", err)
	}

	c.conn.Close()
}

// Gracefully close the client.
func (c *Client) Close() {
	close(c.quit)
}
