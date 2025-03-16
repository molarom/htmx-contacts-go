package archiver

import (
	"fmt"
	"math/rand"
	"time"

	"gitlab.com/romalor/htmx-contacts/pkg/atomic"
)

var a = &Archiver{
	status:   new(atomic.String),
	progress: new(atomic.Float64),
}

type Archiver struct {
	status   *atomic.String
	progress *atomic.Float64
}

func Default() *Archiver {
	if a.status.Value() != "" {
		return a
	}
	a.status.Set("Waiting")
	return a
}

func (a *Archiver) Status() string {
	return a.status.Value()
}

func (a *Archiver) Progress() float64 {
	return a.progress.Value()
}

func (a *Archiver) Run() {
	if a.status.Value() == "Waiting" {
		a.status.Set("Running")
		a.progress.Set(0)
		go a.run()
	}
}

func (a *Archiver) run() {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * time.Duration(rand.Intn(2)))
		if a.status.Value() != "Running" {
			return
		}
		a.progress.Set(float64(i) / float64(10))
		fmt.Println("Here...", a.progress.Value())
	}
	time.Sleep(time.Second * 1)
	if a.status.Value() != "Running" {
		return
	}
	a.status.Set("Complete")
}

func (a *Archiver) Reset() {
	a.status.Set("Waiting")
}

func (a *Archiver) File() string {
	return "contacts.json"
}
