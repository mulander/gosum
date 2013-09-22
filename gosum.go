package gosum

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

type SumFile interface {
	// Returns a map of the entries
	Entries() map[string]string
	// Add the digest of src file as a new entry in the sum file
	Add(name string, src io.Reader) error
	// Checks if the file specified by name has the correct
	// sum against the sum computed from the Reader
	Check(name string, src io.Reader) (bool, error)
	// Writes a new entry from rfc 1321 formatted input
	Write(p []byte) (n int, err error)
	// Writes the file digests in the rfc 1321 format to the target writer
	WriteTo(w io.Writer) (n int64, err error)
	// Read entries in the rfc 1321 format
	Read(p []byte) (n int, err error)
	// Read the entries from a src reader
	ReadFrom(src io.Reader) (n int64, err error)
}

type MD5Sum struct {
	entries map[string]string
	reader  *bytes.Buffer
}

func NewMD5Sum() SumFile {
	return &MD5Sum{
		entries: make(map[string]string),
	}
}

func (m *MD5Sum) Entries() map[string]string {
	return m.entries
}

func (m *MD5Sum) digest(src io.Reader) string {
	hasher := md5.New()
	io.Copy(hasher, src)
	return hex.EncodeToString(hasher.Sum(nil))
}

func (m *MD5Sum) Write(p []byte) (n int, err error) {
	digests := bufio.NewScanner(bytes.NewReader(p))
	for digests.Scan() {
		n += len(digests.Text())
		entry := strings.Split(digests.Text(), "  ")
		digest, name := entry[0], entry[1]
		m.entries[name] = digest
	}
	if err := digests.Err(); err != nil {
		return n, err
	}
	return n, nil
}

func (m *MD5Sum) ReadFrom(src io.Reader) (n int64, err error) {
	digests := bufio.NewScanner(src)
	for digests.Scan() {
		n += int64(len(digests.Text()))
		entry := strings.Split(digests.Text(), "  ")
		digest, name := entry[0], entry[1]
		m.entries[name] = digest
	}

	if err := digests.Err(); err != nil {
		return n, err
	}
	return n, nil
}

func (m *MD5Sum) Add(name string, src io.Reader) error {
	md5sum := m.digest(src)
	m.entries[name] = md5sum
	return nil
}

func (m *MD5Sum) Read(p []byte) (n int, err error) {
	if m.reader == nil {
		var buffer bytes.Buffer
		for name, digest := range m.entries {
			buffer.Write([]byte(fmt.Sprintf("%s  %s\n", digest, name)))
		}
		m.reader = &buffer
	}
	n, err = m.reader.Read(p)
	if err == io.EOF {
		m.reader.Reset()
	}
	return n, err
}

func (m *MD5Sum) WriteTo(w io.Writer) (n int64, err error) {
	for name, digest := range m.entries {
		written, err := w.Write([]byte(fmt.Sprintf("%s  %s\n", digest, name)))
		if err != nil {
			return int64(written), err
		}
		n += int64(written)
	}
	return n, err
}

func (m *MD5Sum) Check(name string, src io.Reader) (bool, error) {
	digest, _ := m.entries[name]

	hash := m.digest(src)
	return hash == digest, nil
}
