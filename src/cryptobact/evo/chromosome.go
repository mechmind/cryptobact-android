package evo

import (
	"crypto/sha1"
	"fmt"
	"io"
	"math/big"
	"time"
)

var _ = time.Now

type DeviceId uint

type Chromosome struct {
	CurrHash *big.Int
	PrevHash big.Int
	Author   uint64
	DNA      *DNA
	//Time time.Time
	Nonce uint
	//Device DeviceId
}

func (c *Chromosome) Hash(nonce uint) *big.Int {
	h := sha1.New()
	io.WriteString(h, fmt.Sprintf("%x{%s}T%dN%dD%d",
		c.PrevHash,
		c.DNA,
		nonce))
	hash := big.NewInt(0)
	hash.SetBytes(h.Sum(nil))
	return hash
}
