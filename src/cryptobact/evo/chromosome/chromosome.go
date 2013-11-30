package chromosome

import "time"

import . "cryptobact/evo/dna"

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
