package main

import (
	"flag"
	"github.com/mulander/gosum"
	"io"
	"log"
	"os"
)

func main() {
	flag.Parse()
	md5sum := gosum.NewMD5Sum()

	log.Println("Current contents of test.md5sum:")
	current := gosum.NewMD5Sum()
	fileSum, err := os.Open("test.md5sum")
	if err != nil {
		log.Println("Can't open test.md5sum")
	} else {
		io.Copy(current, fileSum)
		io.Copy(os.Stdout, current)
		fileSum.Close()
	}

	log.Printf("%q", flag.Args())
	for _, file := range flag.Args() {
		log.Println(file)
		src, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer src.Close()
		fileinfo, err := src.Stat()
		name := fileinfo.Name()
		if name == "stdin" {
			name = "-"
		}
		md5sum.Add(name, src)
	}

	io.Copy(os.Stdout, md5sum)
	dst, err := os.Create("test.md5sum")
	if err != nil {
		log.Fatal(err)
	}
	defer dst.Close()
	io.Copy(dst, md5sum)

}
