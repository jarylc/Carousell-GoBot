package state

import (
	"carousell-gobot/data/config"
	"carousell-gobot/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

var debug = os.Getenv("DEBUG") == "1"

var states = map[string]*models.State{}
var mutex sync.Mutex

var finalPath string

// Load - load state file
func Load(path string) error {
	err := getPath(path)
	if err != nil {
		return err
	}

	if _, err := os.Stat(finalPath); os.IsNotExist(err) {
		err := Save()
		if err != nil {
			return err
		}
	}

	raw, err := ioutil.ReadFile(finalPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw, &states)
	if err != nil {
		return err
	}

	loadPruner()

	return nil
}

// Save - save state file
func Save() error {
	save, err := json.Marshal(states)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(finalPath, save, 0600)
	if err != nil {
		return err
	}
	return nil
}

// ListIDs - list all state IDs
func ListIDs() []string {
	ids := make([]string, 0, len(states))
	for k := range states {
		ids = append(ids, k)
	}
	return ids
}

// Get - get a specific state from state ID
func Get(id string) (*models.State, bool) {
	state, ok := states[id]
	if !ok {
		state = models.NewState(id)
		states[id] = state
		return state, true
	}
	return state, false
}

func getPath(path string) error {
	if strings.HasPrefix(path, "/") {
		finalPath = path
		return nil
	}
	var err error
	finalPath, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	finalPath += "/" + path
	return nil
}

func loadPruner() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			now := time.Now()
			select {
			case <-interrupt:
				return
			case <-time.After(now.Truncate(24*time.Hour).AddDate(0, 0, 1).Sub(now)):
				prune()
			}
		}
	}()
	prune()
}
func prune() {
	if debug {
		log.Println("Running scheduled state pruning")
	}
	mutex.Lock()
	defer mutex.Unlock()
	for id, cState := range states {
		if -time.Until(cState.LastActivity).Hours() >= float64(config.Config.StatePrune*24) { // past pruning date
			if debug {
				log.Printf("\t%s state pruned", id)
			}
			delete(states, id)
		}
	}
	err := Save()
	if err != nil {
		log.Println(err)
	}
}

// TESTING UTILS
var tmp *os.File

func CreateTmp() (*os.File, error) {
	var err error
	tmp, err = ioutil.TempFile("", "state")
	if err != nil {
		return nil, err
	}
	if _, err = tmp.Write([]byte("{}")); err != nil {
		return nil, err
	}
	return tmp, nil
}
func RemoveTmp() error {
	err := os.Remove(tmp.Name())
	if err != nil {
		return err
	}
	return nil
}
