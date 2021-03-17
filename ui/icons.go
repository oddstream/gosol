package ui

import (
	"archive/zip"
	"bytes"
	_ "embed" // go:embed only allowed in Go files that import "embed"
	"fmt"
	"image"
	"io/ioutil"
	"log"
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

func LoadIconMap() {
	// img, _, err := image.Decode(bytes.NewReader(undoBytes))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// IconMap["undo"] = img

	// dci := gg.NewContextForImage(img)
	// w := dci.Width()
	// h := dci.Height()
	// println("button image dimensions", w, h)
	// dci.Scale(float64(rb.width)/float64(w), float64(rb.height)/float64(h))

	iconNames := []string{"help_outline", "menu", "undo"}
	for _, iconName := range iconNames {
		zipFname := fmt.Sprintf("/home/gilbert/gosol/ui/%s-white-android.zip", iconName)
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
			}
		}
	}
	// for _, file := range zf.File {
	// 	println(file.Name)
	// }
}
