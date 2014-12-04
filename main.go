package main

import (
	"os"
	"fmt"
)

func main() {
	podfile := os.Stdin

	podreader, err := NewPodReader(podfile)

	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Input is a POD file, version %d.", podreader.GetVersion()))

	filecount := podreader.GetFileCount()

	fmt.Println(fmt.Sprintf("Found %d files.", filecount))

	files := podreader.ReadFileTable()

	for i := 0; i < len(files); i++ {
		fmt.Println(fmt.Sprintf("* %s (%d bytes)", files[i].name, files[i].size))
	}

	os.Exit(0)
}
