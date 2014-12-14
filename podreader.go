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

	valid, err := podreader.VerifyMagic()
	if err != nil {
		return nil, err
	} else if valid == false {
		return nil, errors.New("Input file is not a POD file, does not start with POD magic.")
	}

	version, err := podreader.GetVersion()
	if err != nil {
		return nil, err
	} else if version != 5 {
		return nil, errors.New(fmt.Sprintf("Version %d of the POD file format not supported.", version))
	}

	podreader.version = version

	return podreader, nil
}

func (podreader *PodReader) ReadFile(file PodFile, out io.Writer) error {
	podstream := podreader.podstream

	if _, err := podstream.Seek(int64(file.offset), os.SEEK_SET); err != nil {
		return errors.New(fmt.Sprintf("Failed to seek to file offset %d: %s", file.offset, err.Error()))
	}

	size := file.size
	buffersize := int32(1024)
	buff := make([]byte, buffersize)

	for size > 0 {
		read, readerr := podstream.Read(buff)
		if read == 0 || readerr == io.EOF {
			return errors.New("Could not read all of file, ran into EOF of the POD!")
		} else if readerr != nil {
			return readerr
		}

		windowsize := buffersize
		if size < buffersize {
			windowsize = size
		}

		if _, err := out.Write(buff[:windowsize]); err != nil {
			return err
		}

		size = size - buffersize
	}

	return nil
}

func (podreader *PodReader) ReadFileTable() ([]PodFile, error) {
	podstream := podreader.podstream

	filecount, err := podreader.GetFileCount()
	if err != nil {
		return nil, err
	}

	// The TOC is stashed away at the end of the file.  Read the start address
	// from the POD file's header.
	tableaddress, err := podreader.GetFileTableAddress()
	if err != nil {
		return nil, err
	}

	if _, err := podstream.Seek(int64(tableaddress), os.SEEK_SET); err != nil {
		return nil, errors.New(fmt.Sprintf("Could not seek to file table: %s", err.Error()))
	}

	ret := make([]PodFile, filecount)
	nametablepositions := make([]int64, filecount)

	// First pass: read the data at the front of the TOC (position and size data).
	for i := 0; i < int(filecount); i++ {
		// Tuck away the seek offset to use in the next step.
		nameposition, err := podstream.ReadInt()
		if err != nil {
			return nil, err
		}

		nametablepositions[i] = int64(nameposition)

		// Read the rest of the data stored in the TOC.
		podfile := new(PodFile)

		podfile.size, err = podstream.ReadInt()
		if err != nil {
			return nil, err
		}

		podfile.offset, err = podstream.ReadInt()
		if err != nil {
			return nil, err
		}

		podfile.size2, err = podstream.ReadInt()
		if err != nil {
			return nil, err
		}

		podfile.unknown1, err = podstream.ReadBytes(4)
		if err != nil {
			return nil, err
		}

		podfile.unknown2, err = podstream.ReadBytes(4)
		if err != nil {
			return nil, err
		}

		podfile.unknown3, err = podstream.ReadBytes(4)
		if err != nil {
			return nil, err
		}

		ret[i] = *podfile
	}

	// Keep track of the start of the name table, since we're going to
	// use it to calculate the position of the name.
	namedictionarystart, err := podstream.Tell()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not determine the current cursor position in the file table: %s", err.Error()))
	}

	// Second pass: Read the data at the end of the TOC (file names)
	for i := 0; i < int(filecount); i++ {
		if _, err := podstream.Seek(namedictionarystart+nametablepositions[i], os.SEEK_SET); err != nil {
			return nil, errors.New(fmt.Sprintf("Could not seek to file name in table: %s", err.Error()))
		}

		ret[i].name, err = podstream.ReadNullTerminatedString()
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Could not read file name in table: %s", err.Error()))
		}
	}

	return ret, nil
}

func (podreader *PodReader) GetFileTableAddress() (int32, error) {
	podstream := podreader.podstream

	// Will this magic location be the same in all versions?
	if _, err := podstream.Seek(264, os.SEEK_SET); err != nil {
		return -1, errors.New(fmt.Sprintf("Could not seek to file table offset: %s", err.Error()))
	}

	result, err := podstream.ReadInt()
	if err != nil {
		return -1, errors.New(fmt.Sprintf("Could not read file table offset value: %s", err.Error()))
	}

	return result, nil
}

func (podreader *PodReader) GetVersion() (int, error) {
	podstream := podreader.podstream

	// Seek past "POD" to the version string byte.
	if _, err := podstream.Seek(3, os.SEEK_SET); err != nil {
		return -1, errors.New(fmt.Sprintf("Could not seek to version offset: %s", err.Error()))
	}

	value, err := podstream.ReadString(1)
	if err != nil {
		return -1, errors.New(fmt.Sprintf("Failed to read version value from POD: %s", err.Error()))
	}

	// Version is stored as a string and as part of the magic ("POD5"),
	// instead read it out as an int which is easier to deal with.
	ivalue, err := strconv.Atoi(value)
	if err != nil {
		return -1, errors.New(fmt.Sprintf("Failed to convert file version \"%s\" into number: %s", value, err.Error()))
	}

	return ivalue, nil
}

func (podreader *PodReader) GetFileCount() (int32, error) {
	podstream := podreader.podstream

	// Will this magic location be the same across all POD versions?
	if _, err := podstream.Seek(88, os.SEEK_SET); err != nil {
		return -1, errors.New(fmt.Sprintf("Could not seek to file count offset: %s", err.Error()))
	}

	value, err := podstream.ReadInt()
	if err != nil {
		return value, errors.New(fmt.Sprintf("Could not read file count value: %s", err.Error()))
	}

	return value, nil
}

func (podreader *PodReader) VerifyMagic() (bool, error) {
	podstream := podreader.podstream
	if _, err := podstream.Seek(0, os.SEEK_SET); err != nil {
		return false, errors.New(fmt.Sprintf("Could not seek to beginning of file to find magic: %s", err.Error()))
	}

	magic, err := podstream.ReadString(3)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Could not read file magic: %s", err.Error()))
	}

	return magic == "POD", nil
}
