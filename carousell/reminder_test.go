package carousell

import (
	"carousell-gobot/data/config"
	"carousell-gobot/data/state"
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

	// create base time
	now := time.Now()

	// create state #1 (3 hours later)
	state1, _ := state.Get("1")
	state1.DealOn = now.Add(time.Hour * time.Duration(3))

	// create state #2 (6 hours later)
	state2, _ := state.Get("2")
	state2.DealOn = now.Add(time.Hour * time.Duration(6))

	// init reminders
	InitReminders()

	// should have 2 sets of reminders (states 1 & 2)
	if len(reminders) != 2 {
		t.Errorf("%d vs 2", len(reminders))
	}

	// check state #1 to have 1 reminders and correct time (2 hours later)
	rState1, ok := reminders[state1]
	if !ok {
		t.Errorf("missing state #1 reminder")
	}
	if len(rState1) != 1 {
		t.Errorf("%d vs 1", len(rState1))
	}
	add2 := now.Add(time.Hour * 2)
	if rState1[0].Time != add2 {
		t.Errorf("%s vs %s", rState1[0].Time.String(), add2)
	}
	rState1[0].Cancel()

	// check state #2 to have 2 reminders and correct times (2 hours & 5 hours later)
	rState2, ok := reminders[state2]
	if !ok {
		t.Errorf("missing state #1 reminder")
	}
	if len(rState2) != 2 {
		t.Errorf("%d vs 2", len(rState2))
	}
	add5 := now.Add(time.Hour * 5)
	if rState2[0].Time != add5 {
		t.Errorf("%s vs %s", rState2[0].Time.String(), add5)
	}
	rState2[0].Cancel()
	if rState2[1].Time != add2 {
		t.Errorf("%s vs %s", rState2[1].Time.String(), add2)
	}
	rState2[1].Cancel()
}
