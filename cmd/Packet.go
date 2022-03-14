package cmd

import (
	"bytes"
	"encoding/gob"
	"log"
)

type PacketType uint8

const (
	PacketLeave = iota
	PacketJoin  = iota + 1
	PacketAudio = iota + 2
)

type Packet struct {
	Iteration int
	Type      PacketType
	ClientID  string
	Frame     AudioFrame
}
type AudioFrame struct {
	Data []byte
}

func (p *Packet) BinarySize() int {
	return 1 + len(p.ClientID) + p.audioFramesSize()
}
func (p *Packet) audioFramesSize() int {
	return len(p.Frame.Data)
}
func (f *AudioFrame) BinarySize() int {
	return 2 + len(f.Data)
}
func (p *Packet) Encode() []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}
func DecodePacket(s []byte) Packet {
	p := Packet{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil {
		log.Fatal(err)
	}
	return p
}
