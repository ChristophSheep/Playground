package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/mysheep/cell"
	"github.com/mysheep/cell/brain"
)

func getCount(bits *[]bool) int {
	count := 0
	for _, bit := range *bits {
		if bit {
			count = count + 1
		}
	}
	return count
}

func getPixelsByName(name string) ([]bool, error) {

	fileName := getFilename(name)

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

func getWeights(name string) ([]float64, error) {

	bits, err := getPixelsByName(name)

	if err != nil {
		return nil, err
	}

	count := getCount(&bits)
	weights := make([]float64, len(bits))
	weight := float64(len(bits)) / float64(count)

	//
	// TODO : negative Weight values to REPRESENT NOT
	//

	for i, bit := range bits {
		if bit {
			weights[i] = weight
		}
	}

	return weights, nil

}

func getAllWeights(names []string) ([][]float64, error) {

	wweights := make([][]float64, len(names))

	i := 0
	for j, name := range names {
		fmt.Printf("%d:%25s ", j, name)
		if i == 4 {
			fmt.Println()
			i = 0
		}

		weights, err := getWeights(name)
		if err != nil {
			return nil, err
		}
		wweights[j] = weights
		i = i + 1
	}
	fmt.Println()

	return wweights, nil
}

func getNow() string {
	return time.Now().Format(brain.TIME_FORMAT)
}

/*
	Create retina cells
*/
func createRetinaCells(retinaCells []*brain.EmitterCell) {
	for i, _ := range retinaCells {
		retinaCells[i] = brain.MakeEmitterCell(fmt.Sprintf("retina%2d", i))
	}
}

/*
	Create objects (recognition) cells
*/
func createObjectCells(objectCells []*brain.MultiCell, files []string, THRESHOLD float64) {
	for j, _ := range objectCells {
		objectCells[j] = brain.MakeMultiCell(files[j], THRESHOLD)
	}
}

/*
	Create display cells
*/
func createDisplayCells(displayCells []*brain.DisplayCell, files []string) {
	for j, _ := range displayCells {
		displayCells[j] = brain.MakeDisplayCell(files[j])
	}
}

/*
	Connect object with display cells
*/
func connectObjectWithDisplayCells(objectCells []*brain.MultiCell, displayCells []*brain.DisplayCell) {
	for j, _ := range objectCells {
		brain.ConnectBy(objectCells[j], displayCells[j], float64(1.0))
	}
}

/*
	Connect retina with objects cells
*/
func connectRetinaWithObjectCells(retinaCells []*brain.EmitterCell, objectCells []*brain.MultiCell, wweights [][]float64) {
	for o, _ := range objectCells {

		if math.Mod(float64(o), float64(200)) == 0.0 {
			fmt.Println(fmt.Sprintf("Connect %d of %d", o, len(objectCells)))
		}

		for r, _ := range retinaCells {
			// TODO: MassConnect without append
			weight := wweights[o][r]
			brain.ConnectBy(retinaCells[r], objectCells[o], weight)
		}
	}
}

func seePixel(name string, retinaCells []*brain.EmitterCell) {

	pixels, err := getPixelsByName(name)

	if err != nil {
		return
	}

	YSIZE := SIZE
	XSIZE := SIZE

	var sendPixelRow = func(row []bool, rowIndex int, time time.Time) {
		for colIndex, bit := range row {
			if bit {
				retinaCells[rowIndex*YSIZE+colIndex].EmitOne(time)
			} else {
				retinaCells[rowIndex*YSIZE+colIndex].EmitZero(time)
			}
		}
	}

	startTime := time.Now()
	for rowIndex := 0; rowIndex < YSIZE; rowIndex = rowIndex + 1 {
		i := rowIndex * XSIZE
		row := pixels[i : i+XSIZE]
		go sendPixelRow(row, rowIndex, startTime)
	}
	elapsed := time.Now().Sub(startTime)
	fmt.Printf("Send duration was %v\n", elapsed)
}

func main() {

	done := make(chan bool)
	waitUntilDone := func() { <-done }

	//
	// Setup Network
	//

	fmt.Printf("SIZE is set to %d\n", SIZE)

	files, err := getFiles(getFolder())

	if err != nil {
		return
	}

	var countObjects = len(files)
	const THRESHOLD = SIZE*SIZE - 2 // TODO:???

	fmt.Printf("%d objects found\n", countObjects)
	fmt.Printf("Cell threshold is set to %d\n", THRESHOLD)

	retinaCells := make([]*brain.EmitterCell, SIZE*SIZE)
	objectCells := make([]*brain.MultiCell, countObjects)
	displayCells := make([]*brain.DisplayCell, countObjects)

	allWeights, err := getAllWeights(files)

	if err != nil {
		fmt.Println(err)
		return
	}

	createRetinaCells(retinaCells)
	createObjectCells(objectCells, files, THRESHOLD)
	createDisplayCells(displayCells, files)

	connectObjectWithDisplayCells(objectCells, displayCells)
	connectRetinaWithObjectCells(retinaCells, objectCells, allWeights)

	//
	// Console Commands
	//
	cmds := map[string]func([]string){
		"quit": func(params []string) { done <- true },
		"exit": func(params []string) { done <- true },
		"q":    func(params []string) { done <- true },
		"see": func(params []string) {
			objIndex, err := strconv.Atoi(params[0])
			name := files[objIndex]
			if err == nil {
				fmt.Println(getNow(), "-", fmt.Sprintf("Retina cells see now '%s'", name))
				fmt.Println(getNow(), "-", "Waiting for answer ...")
				seePixel(name, retinaCells)
			}
		},
		"ws": func(params []string) {
			i, err := strconv.Atoi(params[0])
			if err == nil {
				fmt.Printf("%v\n", objectCells[i].Weights())
			}
		},
	}

	go cell.Console(cmds)

	// Wait until Done
	//
	waitUntilDone()
	//
	// Wait until Done

	fmt.Println("BYE")
}
