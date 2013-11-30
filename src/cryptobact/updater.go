package main

import (
    "cryptobact/engine"
)

const UPDATER_QUEUE = 1

type Updater struct {
    reqs chan *updateRequest
    render *Render
    done, updateReady chan struct{}
}

type updateRequest struct {
    w *engine.World
    done chan struct{}
}

func newUpdater(r *Render) *Updater {
    return &Updater{make(chan *updateRequest, UPDATER_QUEUE), r, make(chan struct{}),
        make(chan struct{})}
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
            // send ok to engine
            req.done <- struct{}{}
            // send ping to main render loop
            r.updateReady <- struct{}{}
        case <- r.done:
            break
        }
    }
}

func (r *Updater) isWorldUpdated() bool {
    select {
    case <-r.updateReady:
        return true
    default:
        return false
    }
}
