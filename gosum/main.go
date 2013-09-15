package main

import (
	"flag"
	"github.com/mulander/gosum"
	"log"
	"os"
)

func main() {
	flag.Parse()
	md5sum := gosum.NewMD5Sum()
	err := md5sum.Open("test.md5sum")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%q", flag.Args())
	for _, file := range flag.Args() {
		log.Println(file)
		src, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer src.Close()
		md5sum.Write(src)
	}
	err = md5sum.Close()
	if err != nil {
		log.Fatal(err)
	}
}
