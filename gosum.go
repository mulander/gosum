package gosum

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

type SumFile interface {
	// Opens an existing sum file and reads it's contents
	// into memory. If the sum file doesn't exist it will
	// be the target for saving.
	Open(name string) error
	// Add the src file as a new entry in the sum file
	Write(src *os.File) error
	// Checks if the specified by name file has the correct
	// sum against the sum from the sum file
	Check(name string) (bool, error)
	// Writes the content of the sum file back to disk.
	// Will overwrite any existing content.
	Close() error
}

type MD5Sum struct {
	SumFile string
	Entries map[string]string
}

func NewMD5Sum() SumFile {
	return &MD5Sum{
		SumFile: "",
		Entries: make(map[string]string),
	}
}

func (m *MD5Sum) Open(name string) error {
	if _, err := os.Stat(name); err == nil {
		src, err := os.Open(name)
		if err != nil {
			return err
		}
		src.Close()

		digests := bufio.NewScanner(src)
		for digests.Scan() {
			entry := strings.Split(digests.Text(), "  ")
			digest, name := entry[0], entry[1]
			m.Entries[name] = digest
		}
		if err := digests.Err(); err != nil {
			return err
		}
	}
	m.SumFile = name
	return nil
}

func (m *MD5Sum) digest(src *os.File) string {
	hasher := md5.New()
	io.Copy(hasher, src)
	return hex.EncodeToString(hasher.Sum(nil))
}

func (m *MD5Sum) Write(src *os.File) error {
	fileinfo, err := src.Stat()
	if err != nil {
		return err
	}

	md5sum := m.digest(src)

	name := fileinfo.Name()
	if name == "stdin" {
		name = "-"
	}
	m.Entries[name] = md5sum
	return nil
}

func (m *MD5Sum) Close() error {
	dst, err := os.Create(m.SumFile)
	if err != nil {
		return err
	}
	defer dst.Close()
	for name, digest := range m.Entries {
		dst.Write([]byte(fmt.Sprintf("%s  %s\n", digest, name)))
	}
	return nil
}

func (m *MD5Sum) Check(name string) (bool, error) {
	src, err := os.Open(name)
	if err != nil {
		return false, err
	}

	digest, _ := m.Entries[name]

	hash := m.digest(src)
	return hash == digest, err
}
