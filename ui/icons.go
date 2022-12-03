package ui

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"time"

	"oddstream.games/gosol/util"
)

var IconMap = map[string]image.Image{}

// check is a helper function which streamlines error checking
func check(e error) {
	if e != nil {
		panic(e)
	}
}

type myCloser interface {
	Close() error
}

// closeFile is a helper function which streamlines closing
// with error checking on different file types.
func closeFile(f myCloser) {
	err := f.Close()
	check(err)
}

// readAll is a wrapper function for ioutil.ReadAll. It accepts a zip.File as
// its parameter, opens it, reads its content and returns it as a byte slice.
func readAll(file *zip.File) []byte {
	fc, err := file.Open()
	check(err)
	defer closeFile(fc)

	content, err := ioutil.ReadAll(fc)
	check(err)

	return content
}

func saveBytesToFile(bytes []byte, pngFname string) {

	// path, err := fullPath(pngFname)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// makeConfigDir()

	file, err := os.Create(pngFname)
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
}

func LoadIconMapFromZipFiles() {

	// temporary hack while figuring out size and type of icons

	// https://gist.github.com/madevelopers/40b269730df687cdcb8b

	// get the icon .zip files from here:
	// https://material.io/resources/icons/style=baseline
	// select Android, white and download the .zip file
	// edit /home/gilbert/ to match your folders

	println("loading ui icons from zip files")

	var gofile *os.File
	var err error
	gofile, err = os.Create("/media/gilbert/T7/gomps/5/ui/embeddedicons.go")
	if err != nil {
		log.Fatal(err)
	}
	defer closeFile(gofile)
	gofile.WriteString("package ui\n")
	gofile.WriteString("// Code automatically generated. DO NOT EDIT.\n\n")
	gofile.WriteString("//lint:file-ignore U1000,ST1003 Ignore unused code and underscores in generated code\n\n")
	gofile.WriteString("import (\n")
	gofile.WriteString("\t_ \"embed\" // go:embed only allowed in Go files that import \"embed\"\n")
	gofile.WriteString(")\n\n")

	iconNames := []string{"bookmark", "bookmark_add", "check_box", "check_box_outline_blank", "close", "done", "done_all", "info", "lightbulb", "list", "menu", "poll", "radio_button_checked", "radio_button_unchecked", "restore", "search", "settings", "star", "undo"}
	for _, iconName := range iconNames {
		zipFname := fmt.Sprintf("/home/gilbert/Downloads/%s-white-android.zip", iconName)
		zf, err := zip.OpenReader(zipFname)
		if err != nil {
			log.Fatal(err)
		}
		defer closeFile(zf)

		pngFname := fmt.Sprintf("res/drawable-hdpi/baseline_%s_white_24.png", iconName)
		for _, file := range zf.File {
			if file.Name == pngFname {
				imgBytes := readAll(file)
				img, _, err := image.Decode(bytes.NewReader(imgBytes))
				check(err)
				IconMap[iconName] = img

				saveBytesToFile(imgBytes, fmt.Sprintf("/home/gilbert/gosol/ui/icons/%s.png", iconName))

				gofile.WriteString(fmt.Sprintf("//go:embed icons/%s.png\n", iconName))
				gofile.WriteString(fmt.Sprintf("var %sIconBytes []byte\n\n", iconName))
			}
		}
	}

	gofile.WriteString("// LoadIconMapFromEmbedded loads icons from go:embed vars\n")
	gofile.WriteString("func LoadIconMapFromEmbedded() {\n")
	for _, iconName := range iconNames {
		gofile.WriteString(fmt.Sprintf("\tdecode(\"%s\", %sIconBytes)\n", iconName, iconName))
	}
	gofile.WriteString("}\n")

	// for _, file := range zf.File {
	// 	println(file.Name)
	// }
}

func decode(name string, variable []byte) {
	img, _, err := image.Decode(bytes.NewReader(variable))
	if err != nil {
		log.Panic(err)
	}
	IconMap[name] = img
}

func LoadIconMap() {
	defer util.Duration(time.Now(), "LoadIconMap")
	if GenerateIcons {
		LoadIconMapFromZipFiles()
	} else {
		LoadIconMapFromEmbedded()
	}
}
