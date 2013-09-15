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
	// Returns a map of the entries
	Entries() map[string]string
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
	sumFile string
	entries map[string]string
}

func NewMD5Sum() SumFile {
	return &MD5Sum{
		sumFile: "",
		entries: make(map[string]string),
	}
}

func (m *MD5Sum) Open(name string) error {
	if _, err := os.Stat(name); err == nil {
		src, err := os.Open(name)
		if err != nil {
			return err
		}
		defer src.Close()

		digests := bufio.NewScanner(src)
		for digests.Scan() {
			entry := strings.Split(digests.Text(), "  ")
			digest, name := entry[0], entry[1]
			m.entries[name] = digest
		}
		if err := digests.Err(); err != nil {
			return err
		}
	}
	m.sumFile = name
	return nil
}

func (m *MD5Sum) Entries() map[string]string {
	return m.entries
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
	m.entries[name] = md5sum
	return nil
}

func (m *MD5Sum) Close() error {
	dst, err := os.Create(m.sumFile)
	if err != nil {
		return err
	}
	defer dst.Close()
	for name, digest := range m.entries {
		dst.Write([]byte(fmt.Sprintf("%s  %s\n", digest, name)))
	}
	return nil
}

func (m *MD5Sum) Check(name string) (bool, error) {
	src, err := os.Open(name)
	if err != nil {
		return false, err
	}

	digest, _ := m.entries[name]

	hash := m.digest(src)
	return hash == digest, err
}
