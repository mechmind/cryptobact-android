package evo

import (
	"fmt"
	"io"
	"time"
	"math/big"
	"crypto/sha1"
)

var _ = time.Now

type Chromosome struct {
	CurrHash *big.Int
	PrevHash big.Int
	Author   uint64
	DNA      *DNA
	Nonce int
}

func (c *Chromosome) Hash(nonce int) *big.Int {
	h := sha1.New()
	io.WriteString(h, fmt.Sprintf("%x{%s}T%dN%dD%d",
		c.PrevHash,
		c.DNA,
		nonce))
	hash := big.NewInt(0)
	hash.SetBytes(h.Sum(nil))
	return hash
}

func (c *Chromosome) String() string {
	return fmt.Sprintf("A:%d [%d] %s", c.Author, c.Nonce, c.DNA)
}
