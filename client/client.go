package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/gen2brain/malgo"
	"log"
	"net"
	"os"
	"strings"
)

type PacketType uint8

type Packet struct {
	Type     PacketType
	ClientID string
	Frame    AudioFrame
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
func main() {
	serverEP := "127.0.0.1"
	if len(os.Args) > 1 {
		serverEP = os.Args[1]
	}
	if !strings.Contains(serverEP, ":") {
		serverEP = fmt.Sprintf("%v:8000", serverEP)
	}

	conn, err := net.Dial("udp", serverEP)
	if err != nil {
		fmt.Printf("Dial err %v", err)
		os.Exit(-1)
	}
	defer conn.Close()

	createAudioStream(conn)

	p := make([]byte, 1024)
	nn, err := conn.Read(p)
	if err != nil {
		fmt.Printf("Read err %v\n", err)
		os.Exit(-1)
	}

	fmt.Printf("%v\n", string(p[:nn]))
}

func createAudioStream(conn net.Conn) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		fmt.Printf("LOG <%v>", message)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()
	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.Playback.Format = malgo.FormatS16
	deviceConfig.Playback.Channels = 1
	deviceConfig.SampleRate = 44000
	deviceConfig.Alsa.NoMMap = 1
	//sizeInBytes := uint32(malgo.SampleSizeInBytes(deviceConfig.Capture.Format))
	onRecvFrames := func(pSample2, pSample []byte, framecount uint32) {
		p := Packet{
			Type:     1,
			ClientID: "newClient",
			Frame: AudioFrame{
				Data: pSample,
			},
		}
		conn.Write(EncodeToBytes(p))
		fmt.Printf("\r Output : %v", p.BinarySize())
	}
	fmt.Println("Hear you speak!")
	captureCallbacks := malgo.DeviceCallbacks{
		Data: onRecvFrames,
	}
	device, err := malgo.InitDevice(ctx.Context, deviceConfig, captureCallbacks)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = device.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Scanln()
}

func EncodeToBytes(p interface{}) []byte {

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
