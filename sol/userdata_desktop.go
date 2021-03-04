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

// Load an already existing UserData object from file
func (ud *UserData) Load() {

	if runtime.GOARCH == "wasm" {
		log.Fatal("WASM detected")
	}

	path, err := fullPath("userdata.json")
	if err != nil {
		return
	}
	file, err := os.Open(path)
	if err == nil && file != nil {
		defer file.Close()

		bytes := make([]byte, 256)
		var count int
		count, err = file.Read(bytes)
		if err != nil {
			log.Fatal(err)
		}
		if count > 0 {
			// golang gotcha reslice buffer to number of bytes actually read
			err = json.Unmarshal(bytes[:count], ud)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// Save writes the UserData object to file
func (ud *UserData) Save() {

	if runtime.GOARCH == "wasm" {
		log.Fatal("WASM detected")
	}

	bytes, err := json.Marshal(ud)
	if err != nil {
		log.Fatal(err)
	}

	path, err := fullPath("userdata.json")
	if err != nil {
		return
	}

	makeConfigDir()

	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}

}

// Load statistics for all variants from JSON to an already-created Statistics object
func (s *Statistics) Load() {

	if runtime.GOARCH == "wasm" {
		log.Fatal("WASM detected")
	}

	path, err := fullPath("statistics.json")
	if err != nil {
		return
	}
	file, err := os.Open(path)
	if err == nil && file != nil {
		defer file.Close()
		fi, err := file.Stat()
		if err != nil {
			log.Fatal("error getting FileInfo for ", path)
		}
		bytes := make([]byte, fi.Size()+8)
		var count int
		count, err = file.Read(bytes)
		if err != nil {
			log.Fatal(path, err)
		}
		if count > 0 {
			// golang gotcha reslice buffer to number of bytes actually read
			err = json.Unmarshal(bytes[:count], s)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// Save writes the Statistics object to file
func (s *Statistics) Save() {

	if runtime.GOARCH == "wasm" {
		log.Fatal("WASM detected")
	}

	bytes, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}

	path, err := fullPath("statistics.json")
	if err != nil {
		return
	}

	makeConfigDir()

	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}

}
