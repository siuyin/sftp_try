package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func main() {
	// A public key may be used to authenticate against the remote
	// server by using an unencrypted PEM-encoded private key file.
	//
	// If you have an encrypted private key, the crypto/x509 package
	// can be used to decrypt it.
	key, err := ioutil.ReadFile(os.Getenv("KEYFILE"))
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKeyWithPassphrase(key, []byte(os.Getenv("PASSWD")))
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	// An SSH client is represented with a ClientConn.
	//
	// To authenticate with the remote server you must pass at least one
	// implementation of AuthMethod via the Auth field in ClientConfig,
	// and provide a HostKeyCallback.
	config := &ssh.ClientConfig{
		User: os.Getenv("USER"),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	conn, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer conn.Close()

	sshCmd(conn, "ls")
	fmt.Println("runing next command ------------")
	sshCmd(conn, "cat junk.md")

}

func sshCmd(c *ssh.Client, cmd string) {
	sess, err := c.NewSession()
	if err != nil {
		log.Fatalf("unable to create session: %v", err)
	}
	defer sess.Close()

	op, err := sess.CombinedOutput(cmd)
	if err != nil {

		log.Fatalf("unable to run ls in session: %v", err)
	}
	fmt.Printf("%s", op)
}
