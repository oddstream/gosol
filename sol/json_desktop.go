// +build desktop

package sol

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"runtime"
)

func fullPath(jsonFname string) (string, error) {
	// os.Getenv("HOME") == "" on WASM
	// could use something like errors.New("math: square root of negative number")
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		println(err)
		return "", err
	}
	// println("UserConfigDir", userConfigDir) // /home/gilbert/.config
	return path.Join(userConfigDir, "oddstream.games", "solitaire", jsonFname), nil
}

func makeConfigDir() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	dir := path.Join(userConfigDir, "oddstream.games", "solitaire")
	err = os.MkdirAll(dir, 0755) // https://stackoverflow.com/questions/14249467/os-mkdir-and-os-mkdirall-permission-value
	if err != nil {
		log.Fatal(err)
	}
	// if path is already a directory, MkdirAll does nothing and returns nil
}

func loadBytesFromFile(jsonFname string) ([]byte, int, error) {

	if runtime.GOARCH == "wasm" {
		log.Fatal("WASM detected")
	}

	path, err := fullPath(jsonFname)
	if err != nil {
		return nil, 0, err
	}

	file, err := os.Open(path)
	if err == nil && file != nil {

		fi, err := file.Stat()
		if err != nil {
			log.Fatal("error getting FileInfo for ", path)
		}
		bytes := make([]byte, fi.Size()+8)

		var count int
		count, err = file.Read(bytes)
		if err != nil {
			log.Fatal(err)
		}
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
		println("loaded", path)
		return bytes, count, nil
	}
	println(path, "does not exist")
	return nil, 0, nil // file does not exist (which is ok)
}

func saveBytesToFile(bytes []byte, jsonFname string) {

	if runtime.GOARCH == "wasm" {
		log.Fatal("WASM detected")
	}

	path, err := fullPath(jsonFname)
	if err != nil {
		log.Fatal(err)
	}

	makeConfigDir()

	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}

	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}

	println("saved", path)
}

// Load an already existing UserData object from file
func (ud *UserData) Load() {

	bytes, count, err := loadBytesFromFile("userdata.json")
	if err != nil || count == 0 {
		return
	}

	// golang gotcha reslice buffer to number of bytes actually read
	err = json.Unmarshal(bytes[:count], ud)
	if err != nil {
		log.Fatal(err)
	}

}

// Save writes the UserData object to file
func (ud *UserData) Save() {

	bytes, err := json.Marshal(ud)
	if err != nil {
		log.Fatal(err)
	}

	saveBytesToFile(bytes, "userdata.json")

}

// Load statistics for all variants from JSON to an already-created Statistics object
func (s *Statistics) Load() {

	bytes, count, err := loadBytesFromFile("statistics.json")
	if err != nil || count == 0 {
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

	bytes, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}

	saveBytesToFile(bytes, "statistics.json")

}

// Save the entire undo stack to file
func (b *Baize) Save() {

	if len(b.UndoStack) == 0 {
		return
	}

	// push an extra state; this will be popped after a Load() and used to populate the baize

	bytes, err := json.MarshalIndent(b.UndoStack, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	saveBytesToFile(bytes, b.Variant+".json")

}

// Load the entire undo stack from file
func (b *Baize) Load(v string) bool {

	bytes, count, err := loadBytesFromFile(v + ".json")
	if err != nil || count == 0 {
		return false
	}

	// golang gotcha reslice buffer to number of bytes actually read
	err = json.Unmarshal(bytes[:count], &b.UndoStack)
	if err != nil {
		log.Fatal(err)
	}
	return b.UndoStack != nil && len(b.UndoStack) > 0
}
