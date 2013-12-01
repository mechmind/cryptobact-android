package main

import (
    "cryptobact/engine"
    "log"
)

const UPDATER_QUEUE = 1

type Updater struct {
    reqs chan *updateRequest
    render *Render
    done chan struct{}
    updateReady chan chan struct{}
}

type updateRequest struct {
    w *engine.World
    done chan struct{}
}

func newUpdater() *Updater {
    return &Updater{make(chan *updateRequest, UPDATER_QUEUE), nil, make(chan struct{}),
        make(chan chan struct{})}
}

func (r *Updater) AttachRender(rn *Render) {
    r.render = rn
}

func (r *Updater) Update(w *engine.World) {
    req := &updateRequest{w, make(chan struct{})}
    select {
    case r.reqs <- req:
        <-req.done
    default:
    }
}

func (r *Updater) fetchUpdates() {
    for {
        select {
        case req := <-r.reqs:
            // update render's bb
            //log.Println("handling world update")
            if r.render == nil {
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
        case <- r.done:
            break
        }
    }
}

func (r *Updater) handleUpdate(w *engine.World) {
    var foodCount int
    for _, f := range w.Food {
        if f != nil {
            r.render.UpdateSet(ID_FOOD, float32(f.X), float32(f.Y), 1.0)
            foodCount++
        }
    }

    //log.Println("handled", foodCount, "food")

    var bactCount int
    for _, p := range w.Populations {
        for _, b := range p.GetBacts() {
            //if b != nil && b.Born {
            if b != nil {
                r.render.UpdateSet(ID_BACTERIA, float32(b.X), float32(b.Y), 1.0)
                bactCount++
            }
        }
    }
    log.Println("handled", bactCount, "bacts")
}

func (r *Updater) isWorldUpdated() chan struct{} {
    select {
    case status := <-r.updateReady:
        return status
    default:
        return nil
    }
}
