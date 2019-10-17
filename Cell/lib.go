package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

const SIZE = 32
const FOLDER = "./Font-Awesome-SVG-PNG/black/png/%d/"

func getFolder() string {
	dir := fmt.Sprintf(FOLDER, SIZE)
	return dir
}

func getFilename(name string) string {
	return getFolder() + name
}

func getFiles(folder string) ([]string, error) {

	file, err := os.Open(folder)
	if err != nil {
		return nil, err
	}

	fileinfos, err := file.Readdir(-1 /*all files*/)
	if err != nil {
		return nil, err
	}
	filenames := make([]string, 0)
	for _, f := range fileinfos {
		filenames = append(filenames, f.Name())
	}

	return filenames, nil
}

func getImage(fileName string) (image.Image, error) {
	r, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(r)
	if err != nil {
		return nil, err
	}
	return img, err
}

func getBit(x, y int, img *image.Image) bool {
	_, _, _, a := (*img).At(x, y).RGBA()
	if a > 128 {
		return true
	}
	return false
}

func printBit(bit bool) {
	if bit {
		fmt.Printf("1")
	} else {
		fmt.Print(".")
	}
}

func printPixels(pixels []bool) {
	for _, b := range pixels {
		printBit(b)

	}
}

func getPixels(img image.Image) ([]bool, error) {

	pixels := make([]bool, SIZE*SIZE)

	for y := 0; y < SIZE; y++ {
		for x := 0; x < SIZE; x++ {
			i := x + y*SIZE
			pixels[i] = getBit(x, y, &img)
		}
	}

	return pixels, nil
}

func test() {

	files, err := getFiles(getFolder())
	if err != nil {
		return
	}

	for i, file := range files {
		fmt.Println(i, "-", file)
	}

	img, err := getImage(getFilename(files[0]))

	if err != nil {
		log.Println(err)
		return
	}

	pixels, err := getPixels(img)
	if err == nil {
		printPixels(pixels)
	}

	fmt.Println()
	fmt.Println(len(pixels))

}
