package state

import (
	"testing"
)

func TestState(t *testing.T) {
	tmp, err := CreateTmp()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err := RemoveTmp()
		if err != nil {
			t.Error(err)
		}
	}()
	err = Load(tmp.Name())
	if err != nil {
		t.Error(err)
	}
	err = Save()
	if err != nil {
		t.Error(err)
	}
}
