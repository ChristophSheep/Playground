package main

import (
	"fmt"
	"math"

	"github.com/mysheep/cell"
	"github.com/mysheep/cell/brain"
)

func main() {

	done := make(chan bool)
	waitUntilDone := func() { <-done }

	//
	// Setup Network
	//

	const SIZE = 64
	const OBJECTS = 700
	const THRESHOLD = 32

	retinaCells := make([]*brain.EmitterCell, SIZE*SIZE)
	objectCells := make([]*brain.Cell, OBJECTS)
	displayCells := make([]*brain.DisplayCell, OBJECTS)

	fmt.Printf("Create %d retina cells\n", len(retinaCells))

	// Create retina cells
	//
	for i, _ := range retinaCells {
		retinaCells[i] = brain.MakeEmitterCell(fmt.Sprintf("retina%d", i))
	}

	fmt.Printf("Create %d object and display cells\n", len(objectCells))

	// Create object and display cells
	for j, _ := range objectCells {
		objectCells[j] = brain.MakeMultiCell(fmt.Sprintf("cell%d", j), THRESHOLD) // TODO: Object Name
		displayCells[j] = brain.MakeDisplayCell(fmt.Sprintf("display%d", j))      // TODO: Object Name
		brain.ConnectBy(objectCells[j], displayCells[j], 1)
	}

	fmt.Printf("Connect %d object cells with display cells\n", len(objectCells))

	for j, _ := range objectCells {
		brain.ConnectBy(objectCells[j], displayCells[j], 1)
	}

	fmt.Printf("Connect %d retina with objects cells \n", len(retinaCells)*len(objectCells))

	// Connect retina cells with object cells
	//
	for r, _ := range retinaCells {

		if int(math.Mod(float64(r), 100)) == 0 {
			fmt.Printf("%d ", r)
		}

		for o, _ := range objectCells {
			brain.ConnectBy(retinaCells[r], objectCells[o], 0)
		}
	}

	// Set weights of object cells
	//
	for o, _ := range objectCells {
		o = o // TODO
	}
	//
	// Console Commands
	//
	cmds := map[string]func(){
		"quit": func() { done <- true },
		"exit": func() { done <- true },
		"q":    func() { done <- true },
		"show": func() {
			for i, _ := range retinaCells {
				retinaCells[i].EmitOne()
				if int(math.Mod(float64(i), 100)) == 0 {
					fmt.Printf("%d ", i)
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
