package gosum

import (
	"bytes"
	"io"
	"testing"
)

func TestAdd(t *testing.T) {
	md5sum := NewMD5Sum()
	md5sum.Add("a.txt", bytes.NewReader([]byte("test")))
	expected := "098f6bcd4621d373cade4e832627b4f6  a.txt\n"

	var result bytes.Buffer
	io.Copy(&result, md5sum)

	if result.String() != expected {
		t.Error("md5sum.Add checksum mismatch")
	}
}

func TestCheck(t *testing.T) {
	md5sum := NewMD5Sum()
	md5sum.Add("a.txt", bytes.NewReader([]byte("test")))

	ok, err := md5sum.Check("a.txt", bytes.NewReader([]byte("test")))
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("md5sum.Add checksum mismatch")
	}
}
