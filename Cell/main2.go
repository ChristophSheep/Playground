package main

import (
	"fmt"
	"strconv"

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

func getWeights(name string) ([]float64, error) {

	fileName := getFilename(name)

	//	fmt.Println("name", name)

	img, err := getImage(fileName)

	if err != nil {
		return nil, err
	}

	bits, err := getPixels(img)

	if err != nil {
		return nil, err
	}

	count := getCount(&bits)
	weights := make([]float64, len(bits))
	weight := float64(len(bits)) / float64(count)

	//	fmt.Printf("%s: count %d weight %d\n", name, count, weight)

	for i, bit := range bits {
		if bit {
			weights[i] = weight
		}
	}

	return weights, nil

}

func main() {

	done := make(chan bool)
	waitUntilDone := func() { <-done }

	//
	// Setup Network
	//

	fmt.Printf("size is set to %d\n", size)

	files, err := getFiles(getFolder())

	if err != nil {
		return
	}

	var countObjects = len(files)
	const THRESHOLD = size*size - 2 // TODO:???

	fmt.Printf("%d objects found\n", countObjects)
	fmt.Printf("Cell threshold is set to %d\n", THRESHOLD)

	retinaCells := make([]*brain.EmitterCell, size*size)
	objectCells := make([]*brain.Cell, countObjects)
	displayCells := make([]*brain.DisplayCell, countObjects)

	fmt.Printf("Create %d retina cells\n", len(retinaCells))

	// Create retina cells
	//
	for i, _ := range retinaCells {
		retinaCells[i] = brain.MakeEmitterCell(fmt.Sprintf("retina%2d", i))
	}

	fmt.Printf("Create %d object and display cells\n", len(objectCells))

	// Create object and display cells
	for j, _ := range objectCells {
		objectCells[j] = brain.MakeMultiCell(files[j], THRESHOLD)
		displayCells[j] = brain.MakeDisplayCell(files[j])
		brain.ConnectBy(objectCells[j], displayCells[j], float64(1.0))
	}

	fmt.Printf("Connect %d object cells with display cells\n", len(objectCells))

	for j, _ := range objectCells {
		// TODO: MassConnect or BulkConnect(0..j)
		// TODO: without append
		brain.ConnectBy(objectCells[j], displayCells[j], float64(1.0))
	}

	fmt.Printf("Connect retina cells with countObjects cells - %d connections\n", len(retinaCells)*len(objectCells))

	wweights := make([][]float64, len(objectCells))

	for j, _ := range objectCells {
		name := files[j]
		weights, err := getWeights(name)
		if err != nil {
			panic(fmt.Sprint("Count not get weights from file ", name))
		}
		wweights[j] = weights
	}

	// Connect retina cells with object cells
	//
	for r, _ := range retinaCells {
		for o, _ := range objectCells {
			// TODO: MassConnect without append
			weight := wweights[o][r]
			brain.ConnectBy(retinaCells[r], objectCells[o], weight)
		}
	}

	//
	// Console Commands
	//
	cmds := map[string]func([]string){
		"quit": func(params []string) { done <- true },
		"exit": func(params []string) { done <- true },
		"q":    func(params []string) { done <- true },
		"obj": func(params []string) {
			i, err := strconv.Atoi(params[0])
			if err == nil {

				for j, w := range objectCells[i].Weights() {
					if w > 0 {
						retinaCells[j].EmitOne()
					}
				}
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
