package example2

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/mysheep/cell/brain"
	"github.com/mysheep/console"
	"github.com/mysheep/imgx"
)

type Spec struct {
	Size           int
	FolderTemplate string
}

func getCountBlack(bits *[]bool) int {
	count := 0
	for _, bit := range *bits {
		if bit {
			count = count + 1
		}
	}
	return count
}

func getWeights(folderTemplate string, size int, name string, threshold float64) ([]float64, error) {

	bits, err := imgx.GetPixelsByName(folderTemplate, size, name)

	if err != nil {
		return nil, err
	}

	countAll := len(bits)
	countBlack := getCountBlack(&bits)
	weights := make([]float64, countAll)

	weightBlack := threshold / float64(countBlack)
	weightWhite := -1.0 * threshold / float64(countAll)

	for i, bit := range bits {
		if bit {
			weights[i] = weightBlack
		} else {
			weights[i] = weightWhite
		}
	}

	return weights, nil

}

func fillSpaces(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s = s + " "
	}
	return s
}

func getAllWeights(folderTemplate string, size int, names []string, threshold float64) ([][]float64, error) {

	wweights := make([][]float64, len(names))

	i := 0
	fmt.Println()
	for j, name := range names {
		fmt.Printf("[%2d] %s%s", j, name, fillSpaces(20-len(name)))

		weights, err := getWeights(folderTemplate, size, name, threshold)
		if err != nil {
			return nil, err
		}
		wweights[j] = weights

		i = i + 1
		if i == 4 {
			i = 0
			fmt.Println()
		}
	}
	fmt.Println()
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

/*
  Let retina see an image with given name
*/
func seePixel(spec Spec, name string, retinaCells []*brain.EmitterCell) {

	pixels, err := imgx.GetPixelsByName(spec.FolderTemplate, spec.Size, name)

	if err != nil {
		return
	}

	YSIZE := spec.Size
	XSIZE := spec.Size

	var sendPixelRow = func(row []bool, rowIndex int, time time.Time) {
		for colIndex, bit := range row {
			if bit {
				retinaCells[rowIndex*YSIZE+colIndex].EmitOne(time)
			} // Do not emit ZERO because cell is in rest and is not firing
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

func Run(spec Spec) {

	done := make(chan bool)
	waitUntilDone := func() { <-done }

	// -------------
	// Setup Network
	// -------------
	//
	//           ObjectCell
	// Retina    +-------+      DisplayCell
	// [o]---w0->|       |
	// [o]---w1->|   A   o----->[Display "A"]
	// [o]---wn->|       |
	//           +-------+
	//

	fmt.Printf("SIZE is set to %d\n", spec.Size)

	imgFiles, err := imgx.GetImgFiles(imgx.GetImgFolder(spec.FolderTemplate, spec.Size))

	if err != nil {
		fmt.Println(err)
		return
	}

	var countObjects = len(imgFiles)
	var THRESHOLD = float64(spec.Size * spec.Size)

	fmt.Printf("%d objects found\n", countObjects)
	fmt.Printf("Cell threshold is set to %f\n", THRESHOLD)

	retinaCells := make([]*brain.EmitterCell, spec.Size*spec.Size)
	objectCells := make([]*brain.MultiCell, countObjects)
	displayCells := make([]*brain.DisplayCell, countObjects)

	allWeights, err := getAllWeights(spec.FolderTemplate, spec.Size, imgFiles, THRESHOLD)

	if err != nil {
		fmt.Println(err)
		return
	}

	createRetinaCells(retinaCells)
	createObjectCells(objectCells, imgFiles, THRESHOLD)
	createDisplayCells(displayCells, imgFiles)

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
			if objIndex >= len(imgFiles) {
				return
			}
			name := imgFiles[objIndex]
			if err == nil {
				fmt.Println(getNow(), "-", fmt.Sprintf("Retina cells see now '%s'", name))
				fmt.Println(getNow(), "-", "Waiting for answer ...")
				seePixel(spec, name, retinaCells)
			}
		},
		"ws": func(params []string) {
			if len(params) == 0 {
				return
			}
			i, err := strconv.Atoi(params[0])
			if err == nil {
				fmt.Printf("%v\n", objectCells[i].Weights())
			}
		},
	}

	go console.Go(cmds)

	// Wait until Done
	//
	waitUntilDone()
	//
	// Wait until Done

	fmt.Println("BYE")
}
