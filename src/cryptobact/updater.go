package main

import (
	"cryptobact/engine"
	"cryptobact/ui"

	"log"
)

var _ = log.Println

const UPDATER_QUEUE = 1

type Updater struct {
	reqs        chan *updateRequest
	field       *ui.Field
	done        chan struct{}
	updateReady chan chan struct{}
}

type updateRequest struct {
	w    *engine.World
	done chan struct{}
}

func newUpdater() *Updater {
	return &Updater{
		make(chan *updateRequest, UPDATER_QUEUE),
		nil,
		make(chan struct{}),
		make(chan chan struct{})}
}

func (r *Updater) AttachField(f *ui.Field) {
	r.field = f
	r.done <- struct{}{}
}

func (r *Updater) Update(w *engine.World) {
	select {
	case <-r.done:
		req := &updateRequest{w.Snapshot(), make(chan struct{})}
		r.reqs <- req
		r.done = req.done
	default:
	}
}

func (r *Updater) fetchUpdates() {
	for {
		select {
		case req := <-r.reqs:
			// update render's bb
			//log.Println("handling world update")
			if r.field == nil {
				// discard update silently
				req.done <- struct{}{}
			} else {
				r.handleUpdate(req.w)
				// send ping to main render loop
				status := make(chan struct{})
				r.updateReady <- status
				<-status
				// send ok to engine
				req.done <- struct{}{}
			}
		}
	}
}

func (r *Updater) handleUpdate(w *engine.World) {
	var foodCount int
	for _, f := range w.Food {
		if f != nil {
			r.field.UpdateFood(float32(f.X), float32(f.Y))
			foodCount++
		}
	}

	//log.Println("handled", foodCount, "food")

	var bactCount int
	for _, p := range w.Populations {
		for _, b := range p.Bacts {
			if b != nil {
				if b.Born {
					r.field.UpdateBact(float32(b.X), float32(b.Y), float32(b.Angle),
						[3]byte{1, 1, 1})
					bactCount++
				} else {
					r.field.UpdateEgg(float32(b.X), float32(b.Y), [3]byte{1, 1, 1})
				}
			}
		}
	}
}

func (r *Updater) isWorldUpdated() chan struct{} {
	select {
	case status := <-r.updateReady:
		return status
	default:
		return nil
	}
}
