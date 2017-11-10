package main

import (
	"fmt"
	"log"

	"siuyin/sftp_try"

	"golang.org/x/crypto/ssh"
)

func main() {
	conn, err := sftp_try.GetClient()
	if err != nil {
		log.Fatalf("cannot get connection: %v", err)
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
