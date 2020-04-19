package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/dainnilsson/forwagent/common"
	"github.com/go-noisesocket/noisesocket"
)

func verifyCallback(publicKey []byte, data []byte) error {
	keys, err := common.ReadKeyList("servers")
	if err != nil {
		return err
	}
	for _, key := range keys {
		if bytes.Equal(key, publicKey) {
			return nil
		}
	}

	publicB64 := base64.StdEncoding.EncodeToString(publicKey)
	fmt.Println("Unknown server key:", publicB64)
	fmt.Println("To allow:")
	fmt.Println("\necho '" + publicB64 + "' >> ~/.forwagent/servers.allowed\n")
	return errors.New("connection closed, unknown public key")
}

func main() {
	keys, err := common.GetKeyPair("client")
	if err != nil {
		fmt.Println("Couldn't read or generate key pair!", err.Error())
		os.Exit(1)
	}

	var host string
	var gpgSocket string
	var sshSocket string
	flag.StringVar(&host, "host", "127.0.0.1:4711", "The host to connect to.")
	flag.StringVar(&gpgSocket, "gpgsocket", filepath.Join(common.GetHomeDir(), ".gnupg", "S.gpg-agent"), "The gpg socket location.")
	flag.StringVar(&sshSocket, "sshsocket", filepath.Join(common.GetHomeDir(), ".gnupg", "S.gpg-agent.ssh"), "The ssh socket location.")
	flag.Parse()

	fmt.Println("Using server:", host)
	fmt.Println("Client key:", base64.StdEncoding.EncodeToString(keys.Public))

	config := noisesocket.ConnectionConfig{
		StaticKey:      keys,
		VerifyCallback: verifyCallback,
	}

	err = runClient(gpgSocket, sshSocket, func() (net.Conn, error) {
		return noisesocket.Dial(host, &config)
	})
	if err != nil {
		fmt.Println("Couldn't start client:", err.Error())
		os.Exit(1)
	}
}
