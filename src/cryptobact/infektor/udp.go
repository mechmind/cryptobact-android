package infektor

// @TBD make transport agnostic

import "fmt"
import "net"
import "time"
//import "log"
//import "bytes"
import "encoding/json"

import "cryptobact/evo"

var _ = time.Now
var _ = fmt.Println

type Infektor struct {
    ports []uint
    sockets []*net.UDPConn
}

func NewInfektor(ports []uint) *Infektor {
    return &Infektor{ports: ports}
}

func (ifk *Infektor) Listen() chan *evo.Chromosome {
    infections := make(chan *evo.Chromosome, 255)
    ifk.sockets = make([]*net.UDPConn, 0)
    for _, num := range ifk.ports {
        server, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", num))
        sock, err := net.ListenUDP("udp", server)
        if err == nil {
            ifk.sockets = append(ifk.sockets, sock)
        } else {
            continue
        }

        go func(sock *net.UDPConn) {
            for {
                buf := make([]byte, 512)
                rlen, _, err := sock.ReadFromUDP(buf)
                if rlen <= 0 || err != nil {
                    continue
                }

                new_chromo := &evo.Chromosome{}
                err = json.Unmarshal(buf[:rlen], new_chromo)
                if err != nil {
                    continue
                }

                infections <- new_chromo
            }
        }(sock)
    }

    return infections
}


func (ifk *Infektor) Spread(population *evo.Population, d time.Duration) {
    ticker := time.NewTicker(d)
    go func (ifk *Infektor, population *evo.Population, ticker *time.Ticker) {
        for {
            <- ticker.C
            ifk.TransmitDisease(population)
        }
    }(ifk, population, ticker)
}


// @TBD send Chromochain instead of evo.Population
func (ifk *Infektor) TransmitDisease(population *evo.Population) bool {
    addrs := GetBroadcastAddrs()
    if len(addrs) == 0 {
        return false
    }

    packet := CreateDiseasePacket(population)

    for _, p := range ifk.ports {
        for _, a := range addrs {
            server, err := net.ResolveUDPAddr("udp",
                fmt.Sprintf("%s:%d", a, p))
            if err != nil {
                continue
            }
            sock, err := net.DialUDP("udp", nil, server)
            if err != nil {
                continue
            }

            sock.Write(packet)
        }
    }

    return true
}

func CreateDiseasePacket(population *evo.Population) []byte {
    //encoder := json.NewEncoder(buffer)
    bacts := population.GetBacts()
    if len(bacts) > 0 {
        buffer, _ := json.Marshal(bacts[0].Chromosome)
        return buffer
    } else {
        return nil
    }
}

func GetBroadcastAddrs() []net.IP {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return nil
    }

    broadcast := make([]net.IP, 0)
    for _, addr := range addrs {
        _, ipnet, err := net.ParseCIDR(addr.String())
        if err != nil {
            continue
        }

        ip := net.IP(make(net.IP, len(ipnet.Mask)))
        for i, b := range ipnet.Mask {
            ip[i] = ipnet.IP[i] | (^b)
        }

        broadcast = append(broadcast, ip)
    }

    return broadcast
}
