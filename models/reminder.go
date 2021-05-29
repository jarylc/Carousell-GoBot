package models

import (
	"time"
)

type Reminder struct {
	Time       time.Time
	ChanCancel chan bool
}

func (r Reminder) Cancel() {
	defer r.Close()
	r.ChanCancel <- true
}

func (r Reminder) Close() {
	close(r.ChanCancel)
}

func NewReminder(time time.Time) *Reminder {
	return &Reminder{
		Time:       time,
		ChanCancel: make(chan bool),
	}
}
