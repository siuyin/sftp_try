package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"siuyin/sftp_try"

	"github.com/pkg/sftp"
)

func main() {
	client, err := sftp_try.GetClient()
	if err != nil {
		log.Fatalf("unable to get client connection: %v", err)
	}

	// open an SFTP session over an existing ssh connection.
	sftp, err := sftp.NewClient(client)
	if err != nil {
		log.Fatal(err)
	}
	defer sftp.Close()

	// walk a directory
	w := sftp.Walk(os.Getenv("WALK_DIR"))
	for w.Step() {
		if w.Err() != nil {
			continue
		}
		fmt.Println(w.Path())
	}

	// read directory
	ents, err := sftp.ReadDir(".")
	if err != nil {
		log.Fatalf("could not read dir: %v", err)
	}
	for _, v := range ents {
		fmt.Printf("%v%v: %v %s\n", v.Name(), dirFlag(v.IsDir()), v.Size(), v.ModTime())
	}

	// open file for reading
	f, err := sftp.Open(os.Getenv("GET_FILE"))
	if err != nil {
		log.Fatalf("unable to open file: %v", err)
	}
	defer f.Close()

	// create output file
	of, err := os.Create(os.Getenv("GET_FILE"))
	if err != nil {
		log.Fatalf("unable to create file: %v", err)
	}
	defer of.Close()

	// create a tee reader
	tee := io.TeeReader(f, os.Stdout)
	_, err = io.Copy(of, tee)
	if err != nil {
		log.Printf("copy error: %v", err)
	}
}

func dirFlag(isDir bool) string {
	if isDir {
		return "/"
	}
	return ""
}
