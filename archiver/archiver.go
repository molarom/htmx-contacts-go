package archiver

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var a = new(Archiver)

type Archiver struct {
	status   string
	progress int
}

func Default() *Archiver {
	return a
}

func (a *Archiver) Status() string {
	return a.status
}

func (a *Archiver) Progress() int {
	return a.progress
}

func (a *Archiver) Run(ctx context.Context) {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * time.Duration(rand.Intn(2)))
		if a.status != "Running" {
			return
		}
		a.progress = (i + 1) / 10
		fmt.Print("Here..." + strconv.Itoa(a.progress))
	}
	time.Sleep(time.Second * 1)
	if a.status != "Running" {
		return
	}
	a.status = "Complete"
}

func (a *Archiver) Reset() {
	a.status = "Waiting"
}

func File() string {
	return "contacts.json"
}
