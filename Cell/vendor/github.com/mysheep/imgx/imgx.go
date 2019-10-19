package imgx

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func GetPixelsByName(folderTemplate string, size int, name string) ([]bool, error) {

	fileName := getFilename(folderTemplate, size, name)

	img, err := getImage(fileName)

	if err != nil {
		return nil, err
	}

	bits, err := getPixels(img)

	if err != nil {
		return nil, err
	}

	return bits, err
}

func GetImgFolder(folderTemplate string, size int) string {
	dir := fmt.Sprintf(folderTemplate, size)
	return dir
}

func getFilename(folderTemplate string, size int, name string) string {
	return GetImgFolder(folderTemplate, size) + name
}

func GetImgFiles(folder string) ([]string, error) {

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
		if f.IsDir() == false {
			filenames = append(filenames, f.Name())
		}
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

	xSize := img.Bounds().Size().X
	ySize := img.Bounds().Size().Y

	pixels := make([]bool, xSize*ySize)

	for y := 0; y < ySize; y++ {
		for x := 0; x < xSize; x++ {
			i := x + y*xSize
			pixels[i] = getBit(x, y, &img)
		}
	}

	return pixels, nil
}
