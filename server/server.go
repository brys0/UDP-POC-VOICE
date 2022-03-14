package main

//
//import (
//	"Stimme/cmd"
//	"bytes"
//	"encoding/gob"
//	"fmt"
//	"github.com/gen2brain/malgo"
//	"log"
//	"net"
//	"os"
//	"strconv"
//)
//
//func main() {
//	serverPort := 8000
//	if len(os.Args) > 1 {
//		if v, err := strconv.Atoi(os.Args[1]); err != nil {
//			fmt.Printf("Invalid port %v, err %v", os.Args[1], err)
//			os.Exit(-1)
//		} else {
//			serverPort = v
//		}
//	}
//
//	addr := net.UDPAddr{
//		Port: serverPort,
//		IP:   net.ParseIP("0.0.0.0"),
//	}
//	server, err := net.ListenUDP("udp", &addr)
//	if err != nil {
//		fmt.Printf("Listen err %v\n", err)
//		os.Exit(-1)
//	}
//	fmt.Printf("Listen at %v\n", addr.String())
//	audioSystem, ctx := makeAudioStream()
//	var decoded cmd.Packet
//	onSendFrames := func(pSample, nil []byte, framecount uint32) {
//		copy(pSample, decoded.Frame.Data)
//	}
//	playbackCallback := malgo.DeviceCallbacks{
//		Data: onSendFrames,
//	}
//	device, err := malgo.InitDevice(ctx.Context, audioSystem, playbackCallback)
//	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//
//	err = device.Start()
//	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//	lastIteration := 0
//	for {
//		p := make([]byte, 2048)
//		nn, raddr, err := server.ReadFromUDP(p)
//		if err != nil {
//			fmt.Printf("Read err  %v", err)
//			continue
//		}
//
//		bytes := p[:nn]
//		decoded = DecodePacket(bytes)
//		//fmt.Printf("\r Received Packet size : %v", decoded.BinarySize())
//
//		fmt.Printf("\r Falling behind! %v ticks behind client total ticks %v \n", decoded.Iteration-lastIteration, decoded.Iteration)
//		lastIteration = decoded.Iteration
//		go func(conn *net.UDPConn, raddr *net.UDPAddr, packet cmd.Packet) {
//			_, err := conn.WriteToUDP([]byte(fmt.Sprintf("Pong: %v", packet)), raddr)
//			if err != nil {
//				fmt.Printf("Response err %v", err)
//			}
//		}(server, raddr, decoded)
//	}
//}
//func makeAudioStream() (malgo.DeviceConfig, *malgo.AllocatedContext) {
//	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
//		fmt.Printf("LOG <%v>", message)
//	})
//	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//	defer func() {
//		_ = ctx.Uninit()
//		ctx.Free()
//	}()
//	deviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)
//	deviceConfig.Capture.Format = malgo.FormatS16
//	deviceConfig.Capture.Channels = 1
//	deviceConfig.Playback.Format = malgo.FormatS16
//	deviceConfig.Playback.Channels = 1
//	deviceConfig.SampleRate = 44000
//	deviceConfig.Alsa.NoMMap = 1
//	return deviceConfig, ctx
//}
//func EncodeToBytes(p interface{}) []byte {
//
//	buf := bytes.Buffer{}
//	enc := gob.NewEncoder(&buf)
//	err := enc.Encode(p)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return buf.Bytes()
//}
//func DecodePacket(s []byte) cmd.Packet {
//	p := cmd.Packet{}
//	dec := gob.NewDecoder(bytes.NewReader(s))
//	err := dec.Decode(&p)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return p
//}
//
////func main() {
//
////	_, err = fmt.Scanln()
////	if err != nil {
////		println(err)
////		os.Exit(1)
////	}
////}
