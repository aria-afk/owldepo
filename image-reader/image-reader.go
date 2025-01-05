package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image-reader/pg"
	"image/png"
	"os"
	"os/exec"

	_ "github.com/joho/godotenv/autoload"
)

// TODO:
// Change panic(fmt.Errorf) to just log.Fatalf

func main() {
	sourceDir := "../scrapper/images/"
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		panic(fmt.Errorf("Could not read source images\n%s", err))
	}
	db := pg.NewPG()
	for i, file := range files {
		// TESTING: remove after
		if i > 10 {
			break
		}
		splitImg, err := splitImage(sourceDir + file.Name())
		if err != nil {
			panic(fmt.Errorf("Error splitting image\n%s", err))
		}
		sir, _ := readSplitImage(splitImg)
		saveResultsToDB(sir, db)
	}
}

// TODO: Add all row areas once testing is done
type SplitImage struct {
	FileName   string
	SearchArea image.Image
	R1UserID   image.Image
	R1Quantity image.Image
	R1Price    image.Image
}

// splitImage takes an uncropped owl screenshot and crops it into all
// possible sections for data gathering
func splitImage(filename string) (SplitImage, error) {
	si := SplitImage{
		FileName: filename,
	}
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
	FileName   string
	SearchArea string
	R1UserID   string
	R1Quantity string
	R1Price    string
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

func readSplitImage(si SplitImage) (SplitImageResults, error) {
	sir := SplitImageResults{
		FileName: si.FileName,
	}

	tmpPath := "./tmp/"
	if err := os.Mkdir(tmpPath, os.ModePerm); err != nil && !errors.Is(err, os.ErrExist) {
		return sir, err
	}

	// Create temp cropped images to be read by ocr
	writeImage(si.SearchArea, tmpPath+"SearchArea.png")
	writeImage(si.R1UserID, tmpPath+"R1UserID.png")
	writeImage(si.R1Quantity, tmpPath+"R1Quantity.png")
	writeImage(si.R1Price, tmpPath+"R1Price.png")

	// Read in image data
	files, err := os.ReadDir(tmpPath)
	if err != nil {
		panic(fmt.Errorf("Issue reading tmp path\n%s", err))
	}
	for _, file := range files {
		cmd := exec.Command(
			"tesseract",
			tmpPath+file.Name(),
			"stdout",
			"--tessdata-dir",
			"../trainer/tess_files/tesstrain/data/",
			"--psm",
			"7",
			"-l",
			"maplestory",
			"--loglevel",
			"ALL",
		)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "TESSDATA_PREFIX=../trainer/tess_files/tesseract/tessdata/")
		if err := cmd.Run(); err != nil {
			panic(fmt.Errorf("Error setting tessdata prefix \n%s", err))
		}
		switch file.Name() {
		case "SearchArea.png":
			sir.SearchArea = out.String()
		case "R1UserID.png":
			sir.R1UserID = out.String()
		case "R1Quantity.png":
			sir.R1Quantity = out.String()
		case "R1Price.png":
			sir.R1Price = out.String()
		default:
			fmt.Println("missed data:", file.Name(), out.String())
		}
	}

	// remove tmp dir
	os.RemoveAll(tmpPath)

	return sir, nil
}

func saveResultsToDB(sir SplitImageResults, db *pg.PG) {
	fmt.Println(sir.FileName)
}
