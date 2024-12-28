// markdata.go tool to help with marking training data for maplestory data
// since we need a lot of it.
package main

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"markdata/lvldb"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// TODO: Make this stuff env args after testing is done
	dataPath := "../../scrapper/images/"
	outPath := "./out/"
	maxIters := 5

	err := os.Mkdir(outPath, os.ModePerm)
	if !errors.Is(err, os.ErrExist) && err != nil {
		panic(fmt.Sprintf("Could not create out dir\n%s", err))
	}

	db := lvldb.NewLvlDB()

	files, err := os.ReadDir(dataPath)
	panicf("Could not open data path for read", err)

	for i, file := range files {
		if i > maxIters {
			break
		}
		name := file.Name()
		exists, err := db.Get(name)
		if err != nil {
			fmt.Println("1", err)
			continue
		}
		if exists != "" {
			continue
		}

		img, err := readImage(dataPath + name)
		if err != nil {
			fmt.Println("2", err)
			continue
		}

		// TODO: We want to map different crop locations
		// to get different samples from the owl printout
		img, err = cropImage(img, image.Rect(0, 0, 100, 100))
		if err != nil {
			fmt.Println("3", err)
			continue
		}

		writeImage(img, outPath+name)
		// TODO: Update db that this image is done

		imageNameNoExtension := strings.Split(name, ".png")
		textFilePath := outPath + imageNameNoExtension[0] + ".gt.txt"
		fmt.Printf("\r vim %s ", textFilePath)
		fmt.Printf("\r")

		cmd := exec.Command("xdg-open", outPath+name)
		if err := cmd.Run(); err != nil {
			fmt.Println("4", err)
			break
		}
	}
}

func readImage(name string) (image.Image, error) {
	fd, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	img, _, err := image.Decode(fd)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func cropImage(img image.Image, crop image.Rectangle) (image.Image, error) {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	// img is an Image interface. This checks if the underlying value has a
	// method called SubImage. If it does, then we can use SubImage to crop the
	// image.
	simg, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}

	return simg.SubImage(crop), nil
}

func writeImage(img image.Image, name string) error {
	fd, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fd.Close()

	return png.Encode(fd, img)
}

func panicf(message string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s\n%s", message, err))
	}
}
