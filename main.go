package main

import (
	"os"
	"fmt"
)

func main() {
	podfile := os.Stdin

	podstream := NewPodStream(podfile)

	ispod, err := VerifyMagic(podstream)
	if err != nil {
		panic(err)
	}

	if ispod == false {
		fmt.Println("Input is not a POD file!")
		os.Exit(-1)
	}

	version, err := GetVersion(podstream)
	if err != nil {
		panic(err)
	}

	if *version != "5" {
		fmt.Println(fmt.Sprintf("Can't read version %s of the POD file format.", *version))
		os.Exit(-1)
	}

	fmt.Println(fmt.Sprintf("Input is a POD file, version %s.", *version))

	filecount, err := GetFileCount(podstream)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Found %d files.", *filecount))

	os.Exit(0)
}

func GetVersion(podstream *PodStream) (*string, error) {
	return podstream.ReadString(3, 1)
}

func GetFileCount(podstream *PodStream) (*uint32, error) {
	return podstream.ReadUInt(88)
}

func VerifyMagic(podstream *PodStream) (bool, error) {
	magic, err := ReadMagic(podstream)

	if err != nil {
		return false, err
	}

	if *magic != "POD" {
		return false, err
	}

	return true, nil
}

func ReadMagic(podstream *PodStream) (*string, error) {
	return podstream.ReadString(0, 3)
}
