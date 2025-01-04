package main

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
)

func main() {
	sourceDir := "../scrapper/images/"
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		panic(fmt.Errorf("Could not read source images\n%s", err))
	}
	for i, file := range files {
		// TESTING: remove after
		if i > 0 {
			break
		}
		splitImg, err := splitImage(sourceDir + file.Name())
		if err != nil {
			panic(fmt.Errorf("Error splitting image\n%s", err))
		}
		readSplitImage(splitImg)
	}
}

type SplitImage struct {
	SearchArea image.Image
	R1UserID   image.Image
	R1Quantity image.Image
	R1Price    image.Image
}

// splitImage takes an uncropped owl screenshot and crops it into all
// possible sections for data gathering
func splitImage(filename string) (SplitImage, error) {
	file, err := os.Open(filename)
	if err != nil {
		return SplitImage{}, err
	}
	img, _, err := image.Decode(file)
	file.Close()
	if err != nil {
		return SplitImage{}, err
	}

	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	simg, ok := img.(subImager)
	if !ok {
		return SplitImage{}, fmt.Errorf("Image does not support cropping")
	}

	areasToCrop := map[string]image.Rectangle{
		"SearchArea": image.Rect(14, 39, 410, 54),
		"R1UserID":   image.Rect(14, 129, 82, 146),
		"R1Quantity": image.Rect(265, 128, 299, 147),
		"R1Price":    image.Rect(181, 128, 260, 147),
	}

	si := SplitImage{}
	for key, area := range areasToCrop {
		switch key {
		case "SearchArea":
			si.SearchArea = simg.SubImage(area)
		case "R1UserID":
			si.R1UserID = simg.SubImage(area)
		case "R1Quantity":
			si.R1Quantity = simg.SubImage(area)
		case "R1Price":
			si.R1Price = simg.SubImage(area)
		}
	}

	return si, nil
}

type SplitImageResults struct {
	SearchArea string
	R1UserID   string
	R1Quantity string
	R1Price    string
}

func readSplitImage(si SplitImage) (SplitImageResults, error) {
	sir := SplitImageResults{}

	tmpPath := "./tmp/"
	if err := os.Mkdir(tmpPath, os.ModePerm); err != nil && !errors.Is(err, os.ErrExist) {
		return sir, err
	}

	writeImage(si.SearchArea, tmpPath+"searchArea.png")
	writeImage(si.R1UserID, tmpPath+"R1UserID.png")
	writeImage(si.R1Quantity, tmpPath+"R1Quantity.png")
	writeImage(si.R1Price, tmpPath+"R1Price.png")

	return sir, nil
}

func writeImage(img image.Image, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(fmt.Errorf("Issue writing image\n%s", err))
	}
	defer file.Close()
	err = png.Encode(file, img)
	if err != nil {
		panic(fmt.Errorf("Issue writing image\n%s", err))
	}
}
