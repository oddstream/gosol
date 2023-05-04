//go:build linux || windows || android || darwin

package sol

import (
	"encoding/json"
	"log"

	"oddstream.games/gosol/util"
)

// Load an already existing Settings object from file
func (s *Settings) Load() {
	// defer util.Duration(time.Now(), "Settings.Load")

	bytes, count, err := util.LoadBytesFromFile("settings.json", false)
	if err != nil || count == 0 || bytes == nil {
		return
	}

	// golang gotcha reslice buffer to number of bytes actually read
	err = json.Unmarshal(bytes[:count], s)
	if err != nil {
		log.Panic("Settings.Load Unmarshal", err)
	}
}

// Save writes the Settings object to file
func (s *Settings) Save() {
	// defer util.Duration(time.Now(), "Settings.Save")

	s.LastVersionMajor = GosolVersionMajor
	s.LastVersionMinor = GosolVersionMinor
	// warning - calling ebiten function ouside RunGame loop will cause fatal panic
	bytes, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	util.SaveBytesToFile(bytes, "settings.json")
}

// Load statistics for all variants from JSON to an already-created Statistics object
func (s *Statistics) Load() {
	// defer util.Duration(time.Now(), "Statistics.Load")

	bytes, count, err := util.LoadBytesFromFile("statistics.json", false)
	if err != nil || count == 0 || bytes == nil {
		return
	}

	// golang gotcha reslice buffer to number of bytes actually read
	err = json.Unmarshal(bytes[:count], s)
	if err != nil {
		log.Fatal(err)
	}
}

// Save writes the Statistics object to file
func (s *Statistics) Save() {
	// defer util.Duration(time.Now(), "Statistics.Save")

	bytes, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	util.SaveBytesToFile(bytes, "statistics.json")
}

// Load an undo stack saved to json
func (b *Baize) Load() {
	bytes, count, err := util.LoadBytesFromFile("saved."+b.variant+".json", true)
	if err != nil || count == 0 || bytes == nil {
		return
	}
	var undoStack []*SavableBaize
	// golang gotcha reslice buffer to number of bytes actually read
	err = json.Unmarshal(bytes[:count], &undoStack)
	if err != nil {
		log.Fatal(err)
	}
	if !b.isSavableStackOk(undoStack) {
		log.Fatal("saved undo stack is not ok")
	}
	b.SetUndoStack(undoStack)
}

// Save the entire undo stack to file
func (b *Baize) Save() {
	// defer util.Duration(time.Now(), "Baize.Save")

	// do not bother to save virgin or completed games
	// if len(b.undoStack) < 2 || b.Complete() {
	// 	return
	// }

	bytes, err := json.MarshalIndent(b.undoStack, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	util.SaveBytesToFile(bytes, "saved."+b.variant+".json")
}
