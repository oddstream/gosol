package ui

import (
	"archive/zip"
	"bytes"
	_ "embed" // go:embed only allowed in Go files that import "embed"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"time"

	"oddstream.games/gosol/util"
)

/*
//go:embed res/drawable-xhdpi/baseline_undo_white_24.png
var undoBytes []byte
*/

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

	// https://material.io/resources/icons/style=baseline

	println("loading ui icons from zip files")

	var gofile *os.File
	var err error
	gofile, err = os.Create("/home/gilbert/gosol/ui/embeddedicons.go")
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

	iconNames := []string{"bookmark", "bookmark_add", "close", "done", "done_all", "info", "list", "menu", "restore", "search", "settings", "star", "undo"}
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

func LoadIconMapFromEmbedded() {
	println("loading ui icons go:embed")
	decode("bookmark", bookmarkIconBytes)
	decode("bookmark_add", bookmarkIconBytes)
	decode("close", closeIconBytes)
	decode("done", doneIconBytes)
	decode("done_all", done_allIconBytes)
	decode("info", infoIconBytes)
	decode("list", listIconBytes)
	decode("menu", menuIconBytes)
	decode("restore", restoreIconBytes)
	decode("search", searchIconBytes)
	decode("settings", settingsIconBytes)
	decode("star", starIconBytes)
	decode("undo", undoIconBytes)
}

func LoadIconMap() {
	defer util.Duration(time.Now(), "LoadIconMap")
	if GenerateIcons {
		LoadIconMapFromZipFiles()
	} else {
		LoadIconMapFromEmbedded()
	}
}
