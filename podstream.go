package main

import (
	"os"
	"bufio"
	"errors"
	"bytes"
	"encoding/binary"
)

type PodStream struct {
	podfile *os.File
}

func NewPodStream(podfile *os.File) (*PodStream) {
	podstream := new(PodStream)
	podstream.podfile = podfile

	return podstream
}

func (podstream *PodStream) Seek(offset int64, whence int) (int64) {
	ret, err := podstream.podfile.Seek(offset, whence)
	if err != nil {
		panic(err)
	}

	return ret
}

func (podstream *PodStream) ReadInt() (int32) {
	value, err := podstream.ReadBytes(4)
	if err != nil {
		panic(err)
	}

	var result int32
	err = binary.Read(bytes.NewBuffer(value), binary.LittleEndian, &result)
	if err != nil {
		panic(err)
	}

	return result
}

func (podstream *PodStream) ReadString(length int) (string) {
	value, err := podstream.ReadBytes(length)
	if err != nil {
		panic(err)
	}

	result := string(value[:])

	return result
}

func (podstream *PodStream) ReadNullTerminatedString() (string) {
	value, err := podstream.ReadUntil(0x00)
	if err != nil {
		panic(err)
	}

	result := string(value[:])

	return result
}

func (podstream *PodStream) Tell() (int64) {
	ret, err := podstream.podfile.Seek(0, os.SEEK_CUR)
	if err != nil {
		panic(err)
	}

	return ret
}

func (podstream *PodStream) ReadUntil(delim byte) ([]byte, error) {
	reader := bufio.NewReader(podstream.podfile)

	result, err := reader.ReadBytes(delim)
	if err != nil {
		return nil, err
	}

	// Trim off the last byte we receive which will be our delim.
	length := len(result)
	if length > 0 {
		length = length - 1
	}

	return result[:length], nil
}

func (podstream *PodStream) ReadBytes(length int) ([]byte, error) {
	value := make([]byte, length)

	bytecount, err := podstream.podfile.Read(value)

	if err != nil {
		return value, err
	} else if bytecount != length {
		return value, errors.New("Unable to read to length expected.")
	}

	return value, nil
}

func (podstream *PodStream) Read(p []byte) (int, error) {
	return podstream.podfile.Read(p)
}
