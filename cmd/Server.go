package cmd

import (
	"fmt"
	"github.com/asaskevich/EventBus"
	"net"
)

type Server struct {
	Port    int
	Host    string
	UDP     *net.UDPConn
	Clients []Client
	Bus     *EventBus.EventBus
}
type Client struct {
	ClientID  int64
	Connected bool
	Addr      *net.UDPAddr
}

func (server *Server) Messages() {
	clients := make(map[string]Client)
	for {
		packetBuf := make([]byte, 0)
		nn, raddr, err := server.UDP.ReadFromUDP(packetBuf)
		if err != nil {
			fmt.Printf("\r Read err %v", err)
			continue
		}
		packetBytes := packetBuf[:nn]
		packet := DecodePacket(packetBytes)
		server.Bus.Publish("server:socket:message", packet)
		switch packet.Type {
		case PacketJoin:
			{
				clients[fmt.Sprintf("%v:%v", raddr.IP, raddr.Port)] = Client{
					ClientID:  199299299292,
					Connected: true,
					Addr:      raddr,
				}
				server.Bus.Publish("server:socket:join", packet)
				break
			}
		case PacketLeave:
			{
				delete(clients, fmt.Sprintf("%v:%v", raddr.IP, raddr.Port))
				server.Bus.Publish("server:socket:leave", packet)
				break
			}
		case PacketAudio:
			{
				server.Bus.Publish("server:socket:audio", packet)
				break
			}
		default:
			{
				break
			}
		}

		if packet.Type == PacketJoin {
			clients[fmt.Sprintf("%v:%v", raddr.IP, raddr.Port)] = Client{
				ClientID:  199299299292,
				Connected: true,
				Addr:      raddr,
			}
		}
	}
}
func NewServer(port int, host string, bus *EventBus.EventBus) (*Server, error) {
	server := Server{
		Port: port,
		Host: host,
		Bus:  bus,
	}
	addr := net.UDPAddr{
		Port: server.Port,
		IP:   net.ParseIP(server.Host),
	}
	serverUDP, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("\n Something went wrong when trying to listen to UDP socket %v", err)
		server.Bus.Publish("server:socket:failed", true)
		return nil, err
	} else {
		server.Bus.Publish("server:socket:online", true)
	}
	server.UDP = serverUDP
	server.Messages()
	return &server, nil
}
