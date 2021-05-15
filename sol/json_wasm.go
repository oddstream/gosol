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
	"runtime"
	"syscall/js"
	"time"

	"oddstream.games/gosol/util"
)

const keyPrefix = "Gosol/"

func loadBytesFromLocalStorage(key string) ([]byte, error) {
	localStorage := js.Global().Get("window").Get("localStorage")
	v := localStorage.Get(keyPrefix + key)
	if v.String() != "<undefined>" {
		bytes := []byte(v.String())
		return bytes, nil
	}
	return nil, fmt.Errorf("%s undefined", keyPrefix+key)
}

func saveBytesToLocalStorage(bytes []byte, key string) {
	js.Global().Get("window").Get("localStorage").Set(keyPrefix+key, string(bytes))
}

// Load an already existing UserData object from browser localStorage
func (ud *UserData) Load() {

	if runtime.GOARCH != "wasm" {
		log.Fatal("GOOS=js GOARCH=wasm required")
	}

	bytes, err := loadBytesFromLocalStorage("UserData")
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(bytes, ud)
	if err != nil {
		log.Println("UserData.Load().Unmarshal() error", err)
	}

}

// Save writes the UserData object to localStorage
func (ud *UserData) Save() {

	if runtime.GOARCH != "wasm" {
		log.Fatal("GOOS=js GOARCH=wasm required")
	}

	bytes, err := json.Marshal(ud)
	if err != nil {
		log.Println("UserData.Save().Marshal() error", err)
	} else {
		saveBytesToLocalStorage(bytes, "UserData")
	}

}

// Load statistics for all variants from JSON to an already-created Statistics object
func (s *Statistics) Load() {

	if runtime.GOARCH != "wasm" {
		log.Fatal("GOOS=js GOARCH=wasm required")
	}

	bytes, err := loadBytesFromLocalStorage("Statistics")
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

	if runtime.GOARCH != "wasm" {
		log.Fatal("GOOS=js GOARCH=wasm required")
	}

	bytes, err := json.Marshal(s)
	if err != nil {
		log.Println("Statistics.Save().Marshal() error", err)
	} else {
		saveBytesToLocalStorage(bytes, "Statistics")
	}

}

// Save the entire undo stack to file
func (b *Baize) Save() {

	if runtime.GOARCH != "wasm" {
		log.Fatal("GOOS=js GOARCH=wasm required")
	}

	bytes, err := json.Marshal(b.UndoStack)
	if err != nil {
		log.Println("%s.Save().Marshal() error", b.Variant, err)
	} else {
		saveBytesToLocalStorage(bytes, b.Variant)
	}

}

// Load the entire undo stack from file
// func (b *Baize) Load(v string) bool {

// 	if runtime.GOARCH != "wasm" {
// 		log.Fatal("GOOS=js GOARCH=wasm required")
// 	}

// 	bytes, err := loadBytesFromLocalStorage(v)
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

func LoadUndoStack(v string) []SaveableBaize {
	defer util.Duration(time.Now(), "LoadUndoStack")

	if runtime.GOARCH != "wasm" {
		log.Fatal("GOOS=js GOARCH=wasm required")
	}

	bytes, err := loadBytesFromLocalStorage(v)
	if err != nil {
		log.Println(err)
		return nil
	}

	var undoStack []SaveableBaize
	// golang gotcha reslice buffer to number of bytes actually read
	err = json.Unmarshal(bytes, &undoStack)
	if err != nil {
		log.Println("%s.Load().Unmarshal() error", v, err)
		// log.Fatal(err)
	}

	if len(undoStack) > 0 {
		return undoStack
	}
	return nil
}
