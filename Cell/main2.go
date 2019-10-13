package main

import (
	"fmt"

	"github.com/mysheep/cell"
	"github.com/mysheep/cell/brain"
)

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
	const THRESHOLD = 256

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
		brain.ConnectBy(objectCells[j], displayCells[j], 1)
	}

	fmt.Printf("Connect %d object cells with display cells\n", len(objectCells))

	for j, _ := range objectCells {
		brain.ConnectBy(objectCells[j], displayCells[j], 1)
	}

	fmt.Printf("Connect retina cells with countObjects cells - %d connections\n", len(retinaCells)*len(objectCells))

	// Connect retina cells with object cells
	//
	for r, _ := range retinaCells {
		for o, _ := range objectCells {
			brain.ConnectBy(retinaCells[r], objectCells[o], 1)
		}
	}

	// Set weights of object cells
	//
	for i, _ := range objectCells {

		name := files[i]

		fileName := getFilename(files[i])

		fmt.Println("name", name)

		img, err := getImage(fileName)

		if err == nil {
			bits, err := getPixels(img)
			if err == nil {
				for j, bit := range bits {
					if bit {
						objectCells[i].SetWeight(j, 1)
					} else {
						objectCells[i].SetWeight(j, 0)
					}
				}
			}
		}
	}

	//
	// Console Commands
	//
	cmds := map[string]func(){
		"quit": func() { done <- true },
		"exit": func() { done <- true },
		"q":    func() { done <- true },
		"obj1": func() {
			for i, w := range objectCells[0].Weights() {

				if w > 0 {
					retinaCells[i].EmitOne()
				}
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
