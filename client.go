package sshhelper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"code.google.com/p/go.crypto/ssh"
)

const (
	user    = "root"
	goodkey = "/Users/name/.ssh/id_rsa"
	server  = "my.lan:22"
)

func main() {
	good()
}

func good() {
	RunRemote("ls", goodkey, user, server)
}

func parsekey(file string) ssh.Signer {
	privateBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal(err)
	}
	return private
}

func RunRemote(cmd, key, user, server string) {
	pkey := parsekey(key)

	auths := []ssh.AuthMethod{ssh.PublicKeys(pkey)}

	cfg := &ssh.ClientConfig{
		User: user,
		Auth: auths,
	}
	cfg.SetDefaults()

	client, err := ssh.Dial("tcp", server, cfg)
	if err != nil {
		log.Fatal(err)
	}

	var session *ssh.Session

	session, err = client.NewSession()
	if err != nil {
		log.Fatal(err)
		//fmt.Println(err)
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(cmd)

	fmt.Println(stdoutBuf.String())
}
