// +build browser

// https://github.com/golang/go/wiki/WebAssembly
// https://pkg.go.dev/syscall/js
// https://github.com/dennwc/dom
// "You cannot import "syscall/js" without GOOS=js/GOARCH=wasm"
// https://github.com/golang/tools/blob/master/gopls/doc/settings.md

package maze

import (
	"encoding/json"
	"log"
	"runtime"
	"syscall/js"
)

// Load an already existing UserData object from browser localStorage
func (ud *UserData) Load() {

	if runtime.GOARCH != "wasm" {
		log.Fatal("GOOS=js GOARCH=wasm required")
	}

	localStorage := js.Global().Get("window").Get("localStorage")
	v := localStorage.Get("gomaze")
	// println(v.String())
	bytes := []byte(v.String())
	if len(bytes) > 0 {
		err := json.Unmarshal(bytes[:], ud)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Save writes the UserData object to localStorage
func (ud *UserData) Save() {

	if runtime.GOARCH != "wasm" {
		log.Fatal("GOOS=js GOARCH=wasm required")
	}

	bytes, err := json.Marshal(ud)
	if err != nil {
		log.Fatal(err)
	}
	str := string(bytes[:])
	localStorage := js.Global().Get("window").Get("localStorage")
	// localStorage.Set("CompletedLevels", strconv.Itoa(ud.CompletedLevels))
	localStorage.Set("gomaze", str)
}
