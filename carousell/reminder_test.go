package carousell

import (
	"carousell-gobot/data/config"
	"carousell-gobot/data/state"
	"strconv"
	"testing"
	"time"
)

func TestReminders(t *testing.T) {
	// setup temporary state
	tmp, err := state.CreateTmp()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err := state.RemoveTmp()
		if err != nil {
			t.Error(err)
		}
	}()
	err = state.Load(tmp.Name())
	if err != nil {
		t.Error(err)
	}

	// set reminder config: 1h & 4h
	config.Config.Reminders = []int8{1, 4}

	// create 2 fake states (deal on 3hr & 6hr later)
	now := time.Now()
	for i := 1; i <= 2; i++ {
		id := strconv.Itoa(i)
		cState, _ := state.Get(id)
		cState.ID = id
		cState.DealOn = now.Add(time.Hour * time.Duration(3*i))
	}

	// init reminders
	InitReminders()

	// should have 2 sets of reminders (2hr & 5hr later)
	if len(reminders) != 2 {
		t.Errorf("%d vs 2", len(reminders))
	}

	// check 2 hours later should have 2 reminders
	three, ok := reminders[now.Add(time.Hour*2)]
	if !ok {
		t.Errorf("missing 2 hours later reminder")
	}
	if len(three) != 2 {
		t.Errorf("%d vs 2", len(three))
	}

	// check 5 hours later should have 1 reminders
	six, ok := reminders[now.Add(time.Hour*5)]
	if !ok {
		t.Errorf("missing 5 hours later reminder")
	}
	if len(six) != 1 {
		t.Errorf("%d vs 1", len(six))
	}
}
