// This is just temp for now
// it seems some of the file names in the training data
// `maplestory-ground-truth` make the tesstrain makefile mad
// this just gives them a number to each file name and its matching image
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	basePath := "../tesstrain/data/maplestory-ground-truth/"
	files, err := os.ReadDir(basePath)
	if err != nil {
		panic(err)
	}
	for i, file := range files {
		if strings.Contains(file.Name(), ".gt.txt") {
			splitNameTxt := strings.Split(file.Name(), ".gt.txt")
			imagePath := basePath + splitNameTxt[0] + ".png"
			txtPath := basePath + file.Name()

			err := os.Rename(imagePath, fmt.Sprintf("%s%d.png", basePath, i))
			if err != nil {
				panic(err)
			}
			err = os.Rename(txtPath, fmt.Sprintf("%s%d.gt.txt", basePath, i))
			if err != nil {
				panic(err)
			}
		}
	}
}
