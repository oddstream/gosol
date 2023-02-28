//go:build linux || windows || android

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
		log.Println(err)
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
		log.Println("loaded", path)
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

	log.Println("saved", path)
}

// Load an already existing Settings object from file
func (s *Settings) Load() {
	// defer util.Duration(time.Now(), "Settings.Load")

	bytes, count, err := loadBytesFromFile("preferences.json", false)
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
	saveBytesToFile(bytes, "preferences.json")
}

// Load statistics for all variants from JSON to an already-created Statistics object
func (s *Statistics) Load() {
	// defer util.Duration(time.Now(), "Statistics.Load")

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
	// defer util.Duration(time.Now(), "Statistics.Save")

	bytes, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	saveBytesToFile(bytes, "statistics.json")
}

// Load an undo stack saved to json
func (b *Baize) Load() {
	bytes, count, err := loadBytesFromFile("saved."+b.variant+".json", true)
	if err != nil || count == 0 || bytes == nil {
		return
	}
	var undoStack []*SavableBaize
	// golang gotcha reslice buffer to number of bytes actually read
	err = json.Unmarshal(bytes[:count], &undoStack)
	if err != nil {
		log.Fatal(err)
	}
	if !b.IsSavableStackOk(undoStack) {
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

	saveBytesToFile(bytes, "saved."+b.variant+".json")
}
