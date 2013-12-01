package infektor

import "fmt"
import "net"
import "time"
import "log"

var _ = time.Now
var _ = fmt.Println

type Infektor struct {
    ports []uint
    sockets []*net.UDPConn
}

func NewInfektor(ports []uint) *Infektor {
    return &Infektor{ports: ports}
}

func (ifk *Infektor) Listen() bool {
    ifk.sockets = make([]*net.UDPConn, 0)
    for _, num := range ifk.ports {
        server, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", num))
        log.Printf("!!! listen on %s\n", server)
        sock, err := net.ListenUDP("udp", server)
        if err == nil {
            ifk.sockets = append(ifk.sockets, sock)
        } else {
            log.Printf("!!! %s\n", err)
            continue
        }

        log.Printf("!!! success without a error")

        go func(sock *net.UDPConn) {
            for {
                buf := make([]byte, 512)
                rlen, _, err := sock.ReadFromUDP(buf)
                if rlen <= 0 || err != nil {
                    continue
                }

                log.Printf("!!! %s\n", buf)
            }
        }(sock)
    }

    return len(ifk.sockets) > 0
}



func (ifk *Infektor) TransmitDisease() bool {
    addrs := GetBroadcastAddrs()
    if len(addrs) == 0 {
        return false
    }

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

            sock.Write([]byte{'h', 'u', 'i', '\n'})
            log.Printf("sent at sock %s:%d\n", a, p)
        }
    }

    return true
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
