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

	// open file and copy to stdout
	f, err := sftp.Open(os.Getenv("GET_FILE"))
	if err != nil {
		log.Fatalf("unable to open file: %v", err)
	}
	defer f.Close()

	_, err = io.Copy(os.Stdout, f)
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
