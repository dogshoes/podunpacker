package main

import (
	"os"
	"errors"
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

func (podstream *PodStream) ReadUInt(position int64) (*uint32, error) {
	value, err := podstream.ReadBytes(position, 4)
	if err != nil {
		return nil, err
	}

	result := binary.LittleEndian.Uint32(value)
	return &result, nil
}

func (podstream *PodStream) ReadString(position int64, length int) (*string, error) {
	value, err := podstream.ReadBytes(position, length)
	if err != nil {
		return nil, err
	}

	result := string(value[:])

	return &result, nil
}

func (podstream *PodStream) ReadBytes(position int64, length int) ([]byte, error) {
	podstream.podfile.Seek(position, 0)

	value := make([]byte, length)

	bytecount, err := podstream.podfile.Read(value)

	if err != nil {
		return value, err
	} else if bytecount != length {
		return value, errors.New("Unable to read to length expected.")
	}

	return value, nil
}
