package evo

import "time"

type DeviceId uint
type ChromosomeHash uint

type Chromosome struct {
    PrevHash ChromosomeHash
    //author PubKey
    DNA *DNA
    Date time.Time
    Nonce uint
    Device DeviceId
}
