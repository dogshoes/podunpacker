package main

import (
	"os"
	"fmt"
	"errors"
	"encoding/binary"
)

func main() {
	podfile := os.Stdin

	ispod, err := VerifyMagic(podfile)
	if err != nil {
		panic(err)
	}

	if ispod == false {
		fmt.Println("Input is not a POD file!")
		os.Exit(-1)
	}

	version, err := GetVersion(podfile)
	if err != nil {
		panic(err)
	}

	if *version != "5" {
		fmt.Println(fmt.Sprintf("Can't read version %s of the POD file format.", *version))
		os.Exit(-1)
	}

	fmt.Println(fmt.Sprintf("Input is a POD file, version %s.", *version))

	filecount, err := GetFileCount(podfile)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Found %d files.", *filecount))

	os.Exit(0)
}

func GetVersion(in *os.File) (*string, error) {
	return ReadString(in, 3, 1)
}

func GetFileCount(in *os.File) (*uint32, error) {
	return ReadUInt(in, 88)
}

func VerifyMagic(in *os.File) (bool, error) {
	magic, err := ReadMagic(in)

	if err != nil {
		return false, err
	}

	if *magic != "POD" {
		return false, err
	}

	return true, nil
}

func ReadMagic(in *os.File) (*string, error) {
	return ReadString(in, 0, 3)
}

func ReadUInt(in *os.File, position int64) (*uint32, error) {
	value, err := ReadBytes(in, position, 4)
	if err != nil {
		return nil, err
	}

	result := binary.LittleEndian.Uint32(value)
	return &result, nil
}

func ReadString(in *os.File, position int64, length int) (*string, error) {
	value, err := ReadBytes(in, position, length)
	if err != nil {
		return nil, err
	}

	result := string(value[:])

	return &result, nil
}

func ReadBytes(in *os.File, position int64, length int) ([]byte, error) {
	in.Seek(position, 0)

	value := make([]byte, length)

	bytecount, err := in.Read(value)

	if err != nil {
		return value, err
	} else if bytecount != length {
		return value, errors.New("Unable to read to length expected.")
	}

	return value, nil
}
