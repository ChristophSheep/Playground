package main

import "github.com/mysheep/playground/Cell/example2"

func main() {

	//example1.Run()

	spec2 := example2.Spec{
		Size:           32,
		FolderTemplate: "./images/Some-Characters/%d/",
	}

	example2.Run(spec2)
}
