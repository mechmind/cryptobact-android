package transport

import (
	"cryptobact/evo"

	"encoding/json"
	"log"
	"net"
)

var _ = log.Println

type InfectoChan chan *evo.Population

type Transporter interface {
	Infect(*evo.Population)
	Catch() InfectoChan
}

func GetBroadcastIPAddrs() []net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}

	broadcast := make([]net.IP, 0)
	for _, addr := range addrs {
		ipaddr, ipnet, err := net.ParseCIDR(addr.String())

		if len(ipaddr.To4()) != 4 {
			// @TODO support ipv6
			continue
		}

		if err != nil {
			continue
		}

		ip := make(net.IP, len(ipnet.Mask))
		for i, b := range ipnet.Mask {
			ip[i] = ipnet.IP[i] | (^b)
		}

		broadcast = append(broadcast, ip)
	}

	return broadcast
}

func PackPopulation(pop *evo.Population) []byte {
	buffer, err := json.Marshal(*pop)
	if err != nil {
		log.Println("can no pack population", err)
		return nil
	}

	return buffer
}
