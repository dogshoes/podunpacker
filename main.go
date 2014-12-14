package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

func main() {
	podfile := os.Stdin

	podreader, err := NewPodReader(podfile)
	if err != nil {
		panic(err)
	}

	version, err := podreader.GetVersion()
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Input is a POD file, version %d.", version))

	filecount, err := podreader.GetFileCount()
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Found %d files.", filecount))

	files, err := podreader.ReadFileTable()
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(files); i++ {
		file := files[i]

		fmt.Println(fmt.Sprintf("* %s (%d bytes / %d bytes), u1: %x, u2: %x, u3: %x", file.name, file.size, file.size2, file.unknown1, file.unknown2, file.unknown3))

		ExtractFile(podreader, file)
	}

	os.Exit(0)
}

func ExtractFile(podreader *PodReader, file PodFile) {
	outpath := NormalizePodPath(file.name)

	if err := FilePathIsValid(outpath); err != nil {
		panic(err)
	}

	handle, err := CreateFile(outpath)
	if err != nil {
		panic(err)
	}

	defer handle.Close()

	if err := podreader.ReadFile(file, handle); err != nil {
		panic(err)
	}

	handle.Sync()
}

func NormalizePodPath(podpath string) string {
	pathparts := strings.Split(podpath, "\\")
	return strings.Join(pathparts, "/")
}

func FilePathIsValid(normalizedpath string) error {
	if path.IsAbs(normalizedpath) {
		return errors.New(fmt.Sprintf("Cannot unpack file with embedded absolute file path: %s", normalizedpath))
	}

	return nil
}

func CreateFile(normalizedpath string) (*os.File, error) {
	rootpath := path.Dir(normalizedpath)
	if rootpath != "." {
		if err := os.MkdirAll(rootpath, 0777); err != nil {
			return nil, err
		}
	}

	handle, err := os.OpenFile(normalizedpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return nil, err
	}

	return handle, nil
}
