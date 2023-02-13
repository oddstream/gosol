// https://github.com/golang/go/wiki/WebAssembly
// https://pkg.go.dev/syscall/js
// https://github.com/dennwc/dom
// "You cannot import "syscall/js" without GOOS=js/GOARCH=wasm"
// https://github.com/golang/tools/blob/master/gopls/doc/settings.md

package sol

import (
	"encoding/json"
	"fmt"
	"log"
	"syscall/js"
)

const keyPrefix = "gosol/"

func loadBytesFromLocalStorage(key string, leaveNoTrace bool) ([]byte, error) {
	keyName := keyPrefix + key
	localStorage := js.Global().Get("window").Get("localStorage")
	v := localStorage.Get(keyName)
	if v.IsUndefined() {
		return nil, fmt.Errorf("%s undefined", keyName)
	}
	if leaveNoTrace {
		// https://developer.mozilla.org/en-US/docs/Web/API/Storage/removeItem
		localStorage.Call("removeItem", keyName)
	}
	// if v.String() != "<undefined>" {
	bytes := []byte(v.String())
	return bytes, nil
	// }
}

func saveBytesToLocalStorage(bytes []byte, key string) {
	keyName := keyPrefix + key
	js.Global().Get("window").Get("localStorage").Set(keyName, string(bytes))
}

// Load an already existing Settings object from browser localStorage
func (ud *Settings) Load() {

	bytes, err := loadBytesFromLocalStorage("preferences", false)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(bytes, ud)
	if err != nil {
		log.Println("Settings.Load().Unmarshal() error", err)
	}

}

// Save writes the Settings object to localStorage
func (ud *Settings) Save() {

	bytes, err := json.Marshal(ud)
	if err != nil {
		log.Println("Settings.Save().Marshal() error", err)
	} else {
		saveBytesToLocalStorage(bytes, "preferences")
	}

}

// Load statistics for all variants from JSON to an already-created Statistics object
func (s *Statistics) Load() {

	bytes, err := loadBytesFromLocalStorage("statistics", false)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(bytes, s)
	if err != nil {
		log.Println("Statistics.Load().Unmarshal() error", err)
	}

}

// Save writes the Statistics object to file
func (s *Statistics) Save() {

	bytes, err := json.Marshal(s)
	if err != nil {
		log.Println("Statistics.Save().Marshal() error", err)
	} else {
		saveBytesToLocalStorage(bytes, "statistics")
	}

}

// Save the entire undo stack to file
func (b *Baize) Save() {

	// do not bother to save virgin or completed games
	if len(b.undoStack) < 2 || b.Complete() {
		return
	}

	bytes, err := json.Marshal(b.undoStack)
	if err != nil {
		log.Println("Baize.Save().Marshal() error", err)
	} else {
		saveBytesToLocalStorage(bytes, "saved")
	}

}

// Load the entire undo stack from file
// func (b *Baize) Load(v string) bool {

// 	if runtime.GOARCH != "wasm" {
// 		log.Fatal("GOOS=js GOARCH=wasm required")
// 	}

// 	bytes, err := loadBytesFromLocalStorage(v, true)
// 	if err != nil {
// 		log.Println(err)
// 		return false
// 	}
// 	err = json.Unmarshal(bytes, &b.UndoStack)
// 	if err != nil {
// 		log.Println("%s.Load().Unmarshal() error", v, err)
// 		return false
// 	}
// 	return b.UndoStack != nil && len(b.UndoStack) > 0

// }

func LoadUndoStack() []*SavableBaize {
	// defer util.Duration(time.Now(), "LoadUndoStack")

	bytes, err := loadBytesFromLocalStorage("saved", true)
	if err != nil {
		log.Println(err)
		return nil
	}

	var undoStack []*SavableBaize
	err = json.Unmarshal(bytes, &undoStack)
	if err != nil {
		log.Println("LoadUndoStack().Unmarshal() error", err)
		// log.Fatal(err)
	}

	if len(undoStack) > 0 {
		return undoStack
	}
	return nil
}
