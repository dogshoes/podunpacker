package main

import (
	"os"
	"fmt"
	"strconv"
)

func main() {
	podfile := os.Stdin

	podstream := NewPodStream(podfile)

	if VerifyMagic(podstream) == false {
		fmt.Println("Input is not a POD file!")
		os.Exit(-1)
	}

	version := GetVersion(podstream)
	if version != 5 {
		fmt.Println(fmt.Sprintf("Can't read version %d of the POD file format.", version))
		os.Exit(-1)
	}

	fmt.Println(fmt.Sprintf("Input is a POD file, version %d.", version))

	filecount := GetFileCount(podstream)

	fmt.Println(fmt.Sprintf("Found %d files.", filecount))

	files := ReadFileTable(podstream)

	for i := 0; i < len(files); i++ {
		fmt.Println(fmt.Sprintf("* %s (%d bytes)", files[i].name, files[i].size))
	}

	os.Exit(0)
}

type PodFile struct {
	size, offset, size2, unknown1, unknown2, unknown3 int32
	name string
}

func ReadFileTable(podstream *PodStream) ([]PodFile) {
	filecount := GetFileCount(podstream)

	// The TOC is stashed away at the end of the file.  Read the start address
	// from the POD file's header.
	tableaddress := GetFileTableAddress(podstream)
	podstream.Seek(int64(tableaddress), os.SEEK_SET)

	ret := make([]PodFile, filecount)
	nametablepositions := make([]int64, filecount)

	// First pass: read the data at the front of the TOC (position and size data).
	for i := 0; i < int(filecount); i++ {
		// Tuck away the seek offset to use in the next step.
		nametablepositions[i] = int64(podstream.ReadInt())

		// Read the rest of the data stored in the TOC.
		podfile := new(PodFile)
		podfile.size = podstream.ReadInt()
		podfile.offset = podstream.ReadInt()
		podfile.size2 = podstream.ReadInt()
		podfile.unknown1 = podstream.ReadInt()
		podfile.unknown2 = podstream.ReadInt()
		podfile.unknown3 = podstream.ReadInt()

		ret[i] = *podfile
	}

	// Keep track of the start of the name table, since we're going to
	// use it to calculate the position of the name.
	namedictionarystart := podstream.Tell()

	// Second pass: Read the data at the end of the TOC (file names)
	for i := 0; i < int(filecount); i++ {
		podstream.Seek(namedictionarystart + nametablepositions[i], os.SEEK_SET)
		ret[i].name = podstream.ReadNullTerminatedString()
	}

	return ret
}

func GetFileTableAddress(podstream *PodStream) (int32) {
	podstream.Seek(264, os.SEEK_SET)

	return podstream.ReadInt()
}

func GetVersion(podstream *PodStream) (int) {
	podstream.Seek(3, os.SEEK_SET)

	// Version is stored as a string and as part of the magic ("POD5"),
	// instead read it out as an int which is easier to deal with.
	value := podstream.ReadString(1)
	ivalue, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}

	return ivalue
}

func GetFileCount(podstream *PodStream) (int32) {
	podstream.Seek(88, os.SEEK_SET)

	return podstream.ReadInt()
}

func VerifyMagic(podstream *PodStream) (bool) {
	podstream.Seek(0, os.SEEK_SET)

	magic := podstream.ReadString(3)
	return magic == "POD"
}
