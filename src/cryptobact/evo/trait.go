package evo

import (
	"fmt"
	"strings"
)

type TraitMap map[string]*Trait

func (t TraitMap) String() string {
	result := make([]string, 0)

	for k, v := range t {
		result = append(result, fmt.Sprintf("[%s] %s", k, v))
	}

	return strings.Join(result, "\n")
}

type Trait struct {
	Pattern string
	Max     int
}

func (t *Trait) String() string {
	return fmt.Sprintf("{%s}: --> %d", t.Pattern, t.Max)
}
