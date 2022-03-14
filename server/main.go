package main

import (
	"Stimme/cmd"
	"fmt"
	"github.com/asaskevich/EventBus"
	"os"
	"sync"
)

var wg sync.WaitGroup

func main() {
	bus := EventBus.New()
	bus.Subscribe("server:socket:online", func(status bool) {
		fmt.Printf("\n Server is online")
	})
	server, err := cmd.NewServer(8080, "localhost", bus)

	if err != nil {
		fmt.Printf("\n Error from server %v", err)
		os.Exit(6)
	}
	fmt.Printf("\n Server created %v", server)

	defer server.UDP.Close()
	wg.Add(1)
	go server.Messages()
	wg.Wait()
}
func KeepProcessAlive() {
	var ch chan int
	<-ch
}
