package sftp_try

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// GetClient read parameters from the environment and returns
// an ssh client connection.
func GetClient() (*ssh.Client, error) {
	// A public key may be used to authenticate against the remote
	// server by using an unencrypted PEM-encoded private key file.
	//
	// If you have an encrypted private key, the crypto/x509 package
	// can be used to decrypt it.
	key, err := ioutil.ReadFile(os.Getenv("KEYFILE"))
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKeyWithPassphrase(key, []byte(os.Getenv("PASSWD")))
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
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
	client, err := ssh.Dial("tcp", host+":"+port, config)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial: %v", err)
	}
	return client, nil
}

// Get retrieves a file with sftp.
func Get(sc *sftp.Client, fn string) error {
	// open file for reading
	f, err := sc.Open(fn)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer f.Close()

	// create output file
	of, err := os.Create(fn)
	if err != nil {
		return fmt.Errorf("unable to create file: %v", err)
	}
	defer of.Close()

	_, err = io.Copy(of, f)
	if err != nil {
		return fmt.Errorf("copy error: %v", err)
	}
	return nil
}
