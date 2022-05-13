//go:build linux || windows || android || darwin

package sol

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"runtime"
	"time"

	"oddstream.games/gosol/util"
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
	return path.Join(userConfigDir, "oddstream.games", "gosol", jsonFname), nil
}

func makeConfigDir() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	dir := path.Join(userConfigDir, "oddstream.games", "gosol")
	err = os.MkdirAll(dir, 0755) // https://stackoverflow.com/questions/14249467/os-mkdir-and-os-mkdirall-permission-value
	if err != nil {
		log.Fatal(err)
	}
	// if path is already a directory, MkdirAll does nothing and returns nil
}

func loadBytesFromFile(jsonFname string, leaveNoTrace bool) ([]byte, int, error) {

	if runtime.GOARCH == "wasm" {
		log.Fatal("WASM detected")
	}

	path, err := fullPath(jsonFname)
	if err != nil {
		return nil, 0, err
	}

	file, err := os.Open(path)
	if err == nil && file != nil {
		var bytes []byte
		var count int
		fi, err := file.Stat()
		if err != nil {
			log.Fatal(err, " getting FileInfo ", path)
		}
		if fi.Size() == 0 {
			log.Print("empty file ", path)
		} else {
			bytes = make([]byte, fi.Size()+8)
			count, err = file.Read(bytes)
			if err != nil {
				log.Fatal(err, " reading ", path)
			}
		}
		err = file.Close()
		if err != nil {
			log.Fatal(err, " closing ", path)
		}
		println("loaded", path)
		if leaveNoTrace {
			os.Remove(path)
		}
		return bytes, count, nil
	}
	// log.Print(err, path)
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

// Load an already existing Preferences object from file
func (prefs *Preferences) Load() {
	if DebugMode {
		defer util.Duration(time.Now(), "Preferences.Load")
	}
	bytes, count, err := loadBytesFromFile("preferences.json", false)
	if err != nil || count == 0 || bytes == nil {
		return
	}

	// golang gotcha reslice buffer to number of bytes actually read
	err = json.Unmarshal(bytes[:count], prefs)
	if err != nil {
		log.Panic("Preferences.Load Unmarshal", err)
	}
}

// Save writes the Preferences object to file
func (prefs *Preferences) Save() {
	if DebugMode {
		defer util.Duration(time.Now(), "Preferences.Save")
	}
	// warning - calling ebiten function ouside RunGame loop will cause fatal panic
	bytes, err := json.MarshalIndent(prefs, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	saveBytesToFile(bytes, "preferences.json")
}

// Load statistics for all variants from JSON to an already-created Statistics object
func (s *Statistics) Load() {
	if DebugMode {
		defer util.Duration(time.Now(), "Statistics.Load")
	}
	bytes, count, err := loadBytesFromFile("statistics.json", false)
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
	if DebugMode {
		defer util.Duration(time.Now(), "Statistics.Save")
	}
	bytes, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	saveBytesToFile(bytes, "statistics.json")
}

// Save the entire undo stack to file
func (b *Baize) Save() {
	if DebugMode {
		defer util.Duration(time.Now(), "Baize.Save")
	}
	// do not bother to save virgin or completed games
	// if len(b.undoStack) < 2 || b.Complete() {
	// 	return
	// }

	bytes, err := json.MarshalIndent(b.undoStack, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	saveBytesToFile(bytes, "saved.json")
}

func LoadUndoStack() []*SavableBaize {
	if DebugMode {
		defer util.Duration(time.Now(), "LoadUndoStack")
	}
	bytes, count, err := loadBytesFromFile("saved.json", true)
	if err != nil || count == 0 || bytes == nil {
		return nil
	}

	var undoStack []*SavableBaize
	// golang gotcha reslice buffer to number of bytes actually read
	err = json.Unmarshal(bytes[:count], &undoStack)
	if err != nil {
		log.Fatal(err)
	}

	if len(undoStack) > 0 {
		return undoStack
	}
	return nil
}
