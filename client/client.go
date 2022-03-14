package main

import (
	"Stimme/cmd"
	"fmt"
	"github.com/gen2brain/malgo"
	"net"
	"os"
	"strings"
)

func main() {
	serverEP := "localhost"
	if len(os.Args) > 1 {
		serverEP = os.Args[1]
	}
	if !strings.Contains(serverEP, ":") {
		serverEP = fmt.Sprintf("%v:8080", serverEP)
	}

	conn, err := net.Dial("udp", serverEP)
	fmt.Printf("\n Connection %v", conn.LocalAddr())
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
	packetIteration := 0
	onRecvFrames := func(pSample2, pSample []byte, framecount uint32) {
		packetIteration += 1
		p := cmd.Packet{
			Iteration: packetIteration,
			Type:      3,
			ClientID:  "newClient",
			Frame: cmd.AudioFrame{
				Data: pSample,
			},
		}
		fmt.Printf("\r Total iteration: %v", packetIteration)
		conn.Write(p.Encode())
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
