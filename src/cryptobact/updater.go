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
		//time.Sleep(time.Second / 5)
		req := &updateRequest{w.Snapshot(), nil}
		r.reqs <- req
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
				r.done <- struct{}{}
			} else {
				r.handleUpdate(req.w)
				// send ping to main render loop
				status := make(chan struct{})
				r.updateReady <- status
				<-status
				// send ok to engine
				r.done <- struct{}{}
			}
		}
	}
}

type bactinfo struct{ x, y float32 }

var lastlist []bactinfo
var thresh float32 = 5.0

func abs(v1, v2 float32) float32 {
	d := v1 - v2
	if d < 0 {
		return -d
	}
	return d
}

func logdiff(b []bactinfo) {
	if len(b) != len(lastlist) {
		lastlist = b
		return
	}
	for idx, i := range b {
		ii := lastlist[idx]
		if abs(ii.x, i.x) > thresh || abs(ii.y, i.y) > thresh {
			log.Println("updater: found teleport at ", idx, ":", ii, i)
			log.Println("updater: all of new them", b)
			log.Println("updater: all of old them", lastlist)
		}
	}
	lastlist = b
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
	var bacts []bactinfo
	for _, p := range w.Populations {
		for _, b := range p.Bacts {
			if b != nil {
				if b.Born {
					r.field.UpdateBact(float32(b.X), float32(b.Y), float32(b.Angle),
						b.GetColor())
					bactCount++
					bacts = append(bacts, bactinfo{float32(b.X), float32(b.Y)})
				} else {
					r.field.UpdateEgg(float32(b.X), float32(b.Y), b.GetColor())
				}
			}
		}
	}
	//log.Println("updater: all bacts", bacts)
}

func (r *Updater) IsWorldUpdated() chan struct{} {
	select {
	case status := <-r.updateReady:
		return status
	default:
		return nil
	}
}
