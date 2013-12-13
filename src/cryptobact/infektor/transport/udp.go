package transport

import (
	"cryptobact/evo"

	"encoding/json"
	"fmt"
	"log"
	"net"
)

type UDP struct {
	ports   []int
	bufSize int
}

func NewUDP(ports []int) *UDP {
	return &UDP{ports: ports, bufSize: 4096}
}

func (u *UDP) Infect(pop *evo.Population) {
	addrs := GetBroadcastIPAddrs()
	if len(addrs) == 0 {
		return
	}

	packet := PackPopulation(pop)

	for _, p := range u.ports {
		for _, a := range addrs {
			server, err := net.ResolveUDPAddr("udp",
				fmt.Sprintf("%s:%d", a, p))
			if err != nil {
				log.Println("can not resolve udp addr", err)
				continue
			}

			sock, err := net.DialUDP("udp", nil, server)
			if err != nil {
				log.Println("can not dial udp", err)
				continue
			}

			sock.Write(packet)
		}
	}
}

func (u *UDP) Catch() InfectoChan {
	ch := make(InfectoChan)
	for _, num := range u.ports {
		server, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", num))
		sock, err := net.ListenUDP("udp", server)
		if err != nil {
			log.Println("can not open port", num, err)
			continue
		}

		go receiveLoop(sock, ch, u.bufSize)
	}

	return ch
}

func receiveLoop(sock *net.UDPConn, ch InfectoChan, bufSize int) {
	for {
		buf := make([]byte, bufSize)
		rlen, _, err := sock.ReadFromUDP(buf)
		if rlen <= 0 {
			continue
		}

		if err != nil {
			log.Println("failed to read from infekted socket", err)
		}

		infection := evo.Population{}
		err = json.Unmarshal(buf[:rlen], &infection)
		if err != nil {
			log.Printf("can not unmarshal json %s\n", buf[:rlen], err)
			continue
		}

		ch <- &infection
	}
}
