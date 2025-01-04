package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/exp/rand"
)

func main() {
	iterations := flag.Int("iterations", 10, "The amount of images to mark")
	flag.Parse()

	db, err := leveldb.OpenFile("./db", nil)
	if err != nil {
		log.Fatalf("Could not open DB\n%s", err)
	}

	rawImagesPath := "../../scrapper/images/"
	outPath := "./maplestory-ground-truth/"
	err = os.Mkdir(outPath, os.ModePerm)
	if !errors.Is(err, os.ErrExist) && err != nil {
		log.Fatalf("Could not create the outpath\n%s", err)
	}

	files, err := os.ReadDir(rawImagesPath)
	if err != nil {
		log.Fatalf("Could not open images path to get training data\n%s", err)
	}

	i := 0
	for _, file := range files {
		if i >= *iterations {
			break
		}

		// Skip files we already marked
		if fileAlreadyMarked(db, file.Name()) {
			continue
		}

		img, err := readImage(rawImagesPath + file.Name())
		if err != nil {
			log.Fatalf("Error opening image file\n%s", err)
		}
		croppedImg, err := cropImage(img, getCropShape())
		if err != nil {
			log.Fatalf("Error cropping image\n%s", err)
		}
		err = writeImage(croppedImg, outPath+file.Name())
		if err != nil {
			log.Fatalf("Error writting cropped image\n%s", err)
		}

		imgNameSplit := strings.Split(file.Name(), ".png")
		vimTextFilePath := outPath + imgNameSplit[0] + ".gt.txt"
		// Print vim command out
		fmt.Printf("\r      vim %s", vimTextFilePath)
		fmt.Printf("\r")

		cmd := exec.Command("xdg-open", outPath+file.Name())
		if err := cmd.Run(); err != nil {
			log.Fatalf("Error running xdg-open\n%s", err)
			break
		}

		err = db.Put([]byte(file.Name()), []byte(file.Name()), nil)
		if err != nil {
			// TODO: Handle this err better
			log.Println(err)
		}

		// Prune any files that dont have a .gt.txt
		// (any "skipped" files)
		if _, err := os.Stat(vimTextFilePath); errors.Is(err, os.ErrNotExist) {
			os.Remove(outPath + file.Name())
		}

		i += 1
	}
}

// ---------------------------------- DB Utils  ----------------------------
func fileAlreadyMarked(db *leveldb.DB, filename string) bool {
	existing, err := db.Get([]byte(filename), nil)
	if !errors.Is(err, leveldb.ErrNotFound) && err != nil {
		log.Fatalf("Error looking up existing filenmae in db\n%s", err)
	} else if err != nil && errors.Is(err, leveldb.ErrNotFound) {
		return false
	}
	return len(string(existing)) > 4
}

// -------------------------------------------------------------------------

// ---------------------------------- Image Utils  _------------------------
func readImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func cropImage(img image.Image, crop image.Rectangle) (image.Image, error) {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	simg, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("Image does not support cropping")
	}
	return simg.SubImage(crop), nil
}

func writeImage(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, img)
}

// https://pixspy.com/
func getCropShape() image.Rectangle {
	// For training data we want to limit this to the top search
	// and then the first rows results
	// training data MUST BE 1 line of text.
	dims := []image.Rectangle{
		// Top of image "Search results for {ItemName} that you entered."
		image.Rect(14, 39, 410, 54),
		// Users ID (first row)
		image.Rect(14, 129, 82, 146),
		// Quantity (first row)
		image.Rect(265, 128, 299, 147),
		// Image price (first row)
		image.Rect(181, 128, 260, 147),
	}
	return dims[rand.Intn(len(dims))]
}

// -------------------------------------------------------------------------
