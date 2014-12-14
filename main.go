package main

import (
	"os"
	"fmt"
	"path"
	"strings"
	"errors"
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
		file := files[i]

		fmt.Println(fmt.Sprintf("* %s (%d bytes / %d bytes), u1: %x, u2: %x, u3: %x", file.name, file.size, file.size2, file.unknown1, file.unknown2, file.unknown3))

		outpath := NormalizePodPath(file.name)
		if patherr := FilePathIsValid(outpath); patherr != nil {
			panic(patherr)
		}

		handle, err := CreateFile(outpath)
		if err != nil {
			panic(err)
		}
		
		podreader.ReadFile(file, handle)

		handle.Sync()
		handle.Close()
	}

	os.Exit(0)
}

func NormalizePodPath(podpath string) (string) {
	pathparts := strings.Split(podpath, "\\")
	return strings.Join(pathparts, "/")
}

func FilePathIsValid(normalizedpath string) (error) {
	if path.IsAbs(normalizedpath) {
		return errors.New(fmt.Sprintf("Cannot unpack file with absolute file path: %s.", normalizedpath))
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

	handle, err := os.OpenFile(normalizedpath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0664)
	if err != nil {
		return nil, err
	}

	return handle, nil
}
