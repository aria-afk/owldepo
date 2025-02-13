package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image-reader/pg"
	"image/png"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

// TODO:
// Change panic(fmt.Errorf) to just log.Fatalf

type ItemMap struct {
	Entries map[string]ItemMapEntry `json:"entries"`
}

type ItemMapEntry struct {
	LibHref string `json:"libHref"`
	Name    string `json:"name"`
	Type    string `json:"type"`
}

func loadItemMap() ItemMap {
	var im ItemMap
	itemsPath := "../trainer/item-lib-scraper/mapleitems.json"
	itemFile, err := os.Open(itemsPath)
	if err != nil {
		panic(fmt.Sprintf("Could not load the item map from ../trainer/item-lib-scraper/mapleitems.json\n%s", err))
	}

	if err = json.NewDecoder(itemFile).Decode(&im); err != nil {
		panic(fmt.Sprintf("Could not decode json into struct from ../trainer/item-lib-scraper/mapleitems.json\n%s", err))
	}

	return im
}

func main() {
	itemMap := loadItemMap()
	sourceDir := "../scrapper/images/"
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		panic(fmt.Errorf("Could not read source images\n%s", err))
	}
	db := pg.NewPG()
	for i, file := range files {
		var exists bool
		err := db.Conn.QueryRow(`SELECT EXISTS(SELECT imagekey FROM images_ingested WHERE imagekey = $1)`, file.Name()).Scan(&exists)
		if err != nil {
			log.Fatalf("Error performing retrival query on images_ingested \n%s", err)
		}
		if exists {
			continue
		}

		fmt.Printf("\r processing images [%d/%d]", i, len(files))
		splitImg, err := splitImage(sourceDir + file.Name())
		if err != nil {
			panic(fmt.Errorf("Error splitting image\n%s", err))
		}
		sir, _ := readSplitImage(splitImg)
		saveResultsToDB(sir, file.Name(), itemMap, db)
	}
}

// TODO: Add all row areas once testing is done
type SplitImage struct {
	SearchArea image.Image
	R1UserID   image.Image
	R1Quantity image.Image
	R1Price    image.Image
	FileName   string
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
			sir.SearchArea = strings.ReplaceAll(out.String(), "\n", "")
		case "R1UserID.png":
			sir.R1UserID = strings.ReplaceAll(out.String(), "\n", "")
		case "R1Quantity.png":
			sir.R1Quantity = strings.ReplaceAll(out.String(), "\n", "")
		case "R1Price.png":
			sir.R1Price = strings.ReplaceAll(out.String(), "\n", "")
		default:
			fmt.Println("missed data:", file.Name(), out.String())
		}
	}

	// remove tmp dir
	os.RemoveAll(tmpPath)

	return sir, nil
}

func saveResultsToDB(sir SplitImageResults, imagekey string, itemMap ItemMap, db *pg.PG) {
	pfn := sir.parseFileName()
	ok, itemName := sir.parseItemName(itemMap)
	if !ok {
		_, err := db.Conn.Exec(`INSERT INTO failed_images (imagekey) VALUES ($1) ON CONFLICT DO NOTHING`, imagekey)
		if err != nil {
			log.Fatalf("Error performing insert into failed_images query\n%s", err)
		}
		return
	}

	rowResults := []RowResults{}
	rowCount := 1
	for rowCount <= 8 {
		ok, parsedRow := sir.parseRow(rowCount)
		if ok {
			rowResults = append(rowResults, parsedRow)
		}
		rowCount++
	}

	_, err := db.Conn.Exec(`INSERT INTO ITEMS (id) VALUES ($1) ON CONFLICT DO NOTHING`, itemName)
	if err != nil {
		log.Fatalf("Error performing insert into items query\n%s", err)
	}

	if len(rowResults) > 0 {
		var itemEntryUuid string
		err := db.Conn.QueryRow(`INSERT INTO item_entries (item_id, time) VALUES ($1, $2) RETURNING item_entries.uuid`, itemName, pfn.Timestamp).Scan(&itemEntryUuid)
		if err != nil {
			log.Fatalf("Error performing insert into item_entries query\n%s", err)
		}
		for _, r := range rowResults {
			_, err := db.Conn.Exec(`INSERT INTO item_entry_info (item_entry_uuid, seller_id, quantity, price) VALUES ($1, $2, $3, $4);`,
				itemEntryUuid,
				r.UserID,
				r.Quantity,
				r.Price,
			)
			if err != nil {
				log.Fatalf("Error performing insert into item_entries query\n%s", err)
			}
		}
	}

	_, err = db.Conn.Exec(`INSERT INTO images_ingested (imagekey) VALUES ($1)`, imagekey)
	if err != nil {
		log.Fatalf("Error performing insert into images_ingested query\n%s", err)
	}
}

type ParsedFileName struct {
	TaskId    string
	Timestamp string
	FileName  string
}

func (sir *SplitImageResults) parseFileName() ParsedFileName {
	// NOTE: THIS NEEDS TO BE UPDATED IF WE CHANGE PATHS
	prepend := "../scrapper/images/"
	fileNameNoPrepend := strings.Split(sir.FileName, prepend)
	splitFileName := strings.Split(fileNameNoPrepend[1], "~")
	return ParsedFileName{
		TaskId:    splitFileName[0],
		Timestamp: splitFileName[1],
		FileName:  splitFileName[2],
	}
}

func (sir *SplitImageResults) parseItemName(itemMap ItemMap) (bool, string) {
	endingTrim := strings.Split(sir.SearchArea, "that you")
	prependTrim := strings.Split(endingTrim[0], "Search results for")
	if len(prependTrim) <= 1 {
		return false, ""
	}
	parsedItemName := prependTrim[1][1 : len(prependTrim[1])-1]

	_, ok := itemMap.Entries[parsedItemName]
	if !ok {
		// if an item was not recognized as an item we can preform a simalarity search
		// on each item to attempt to locate the proper one.
		// for now I want to collect the errors and review them.
		return false, ""
	}
	return true, parsedItemName
}

type RowResults struct {
	UserID   string
	Quantity int
	Price    int
}

func (sir *SplitImageResults) parseRow(rowIdx int) (bool, RowResults) {
	var userId string
	var quantity string
	var price string
	switch rowIdx {
	case 1:
		userId = strings.ReplaceAll(sir.R1UserID, " ", "")
		quantity = sir.R1Quantity
		price = strings.ReplaceAll(sir.R1Price, ",", "")
	default:
		return false, RowResults{}
	}

	priceInt, err := strconv.Atoi(price)
	if err != nil {
		return false, RowResults{}
	}
	quantityInt, err := strconv.Atoi(quantity)
	if err != nil {
		return false, RowResults{}
	}

	return true, RowResults{
		UserID:   userId,
		Quantity: quantityInt,
		Price:    priceInt,
	}
}
