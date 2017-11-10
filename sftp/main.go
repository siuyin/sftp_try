package main

import (
	"fmt"
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

	if err := sftp_try.Get(sftp, os.Getenv("GET_FILE")); err != nil {
		log.Printf("copy error: %v", err)
	}

}

func dirFlag(isDir bool) string {
	if isDir {
		return "/"
	}
	return ""
}
