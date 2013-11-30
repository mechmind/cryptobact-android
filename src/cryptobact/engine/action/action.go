package action

import(
	"cryptobact/engine/world"
)

type ActionSubject interface {
}

type ActionObject interface {
}

type Action struct {
}

func (a *Action) Apply(w *World) {
	return
}
