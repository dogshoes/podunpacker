package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type PodReader struct {
	podstream *PodStream
	version   int
}

type PodFile struct {
	size, offset, size2          int32
	name                         string
	unknown1, unknown2, unknown3 []byte
}

func NewPodReader(podfile *os.File) (*PodReader, error) {
	podreader := new(PodReader)
	podreader.podstream = NewPodStream(podfile)

	if podreader.VerifyMagic() == false {
		return nil, errors.New("Input file is not a POD file.")
	}

	version := podreader.GetVersion()
	if version != 5 {
		return nil, errors.New(fmt.Sprintf("Can't read version %d of the POD file format.", version))
	}

	podreader.version = version

	return podreader, nil
}

func (podreader *PodReader) ReadFile(file PodFile, out io.Writer) error {
	podstream := podreader.podstream

	podstream.Seek(int64(file.offset), os.SEEK_SET)

	size := file.size
	buffersize := int32(1024)
	buff := make([]byte, buffersize)

	for size > 0 {
		read, readerr := podstream.Read(buff)
		if read == 0 {
			// EOF
			break
		} else if readerr != nil {
			return readerr
		}

		windowsize := buffersize
		if size < buffersize {
			windowsize = size
		}

		_, err := out.Write(buff[:windowsize])
		if err != nil {
			return err
		}

		size = size - buffersize
	}

	return nil
}

func (podreader *PodReader) ReadFileTable() []PodFile {
	podstream := podreader.podstream

	filecount := podreader.GetFileCount()

	// The TOC is stashed away at the end of the file.  Read the start address
	// from the POD file's header.
	tableaddress := podreader.GetFileTableAddress()
	podstream.Seek(int64(tableaddress), os.SEEK_SET)

	ret := make([]PodFile, filecount)
	nametablepositions := make([]int64, filecount)

	// First pass: read the data at the front of the TOC (position and size data).
	for i := 0; i < int(filecount); i++ {
		var err error
		// Tuck away the seek offset to use in the next step.
		nametablepositions[i] = int64(podstream.ReadInt())

		// Read the rest of the data stored in the TOC.
		podfile := new(PodFile)
		podfile.size = podstream.ReadInt()
		podfile.offset = podstream.ReadInt()
		podfile.size2 = podstream.ReadInt()
		if podfile.unknown1, err = podstream.ReadBytes(4); err != nil {
			panic(err)
		}

		if podfile.unknown2, err = podstream.ReadBytes(4); err != nil {
			panic(err)
		}

		if podfile.unknown3, err = podstream.ReadBytes(4); err != nil {
			panic(err)
		}

		ret[i] = *podfile
	}

	// Keep track of the start of the name table, since we're going to
	// use it to calculate the position of the name.
	namedictionarystart := podstream.Tell()

	// Second pass: Read the data at the end of the TOC (file names)
	for i := 0; i < int(filecount); i++ {
		podstream.Seek(namedictionarystart+nametablepositions[i], os.SEEK_SET)
		ret[i].name = podstream.ReadNullTerminatedString()
	}

	return ret
}

func (podreader *PodReader) GetFileTableAddress() int32 {
	podstream := podreader.podstream
	podstream.Seek(264, os.SEEK_SET)

	return podstream.ReadInt()
}

func (podreader *PodReader) GetVersion() int {
	podstream := podreader.podstream
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

func (podreader *PodReader) GetFileCount() int32 {
	podstream := podreader.podstream
	podstream.Seek(88, os.SEEK_SET)

	return podstream.ReadInt()
}

func (podreader *PodReader) VerifyMagic() bool {
	podstream := podreader.podstream
	podstream.Seek(0, os.SEEK_SET)

	magic := podstream.ReadString(3)
	return magic == "POD"
}
