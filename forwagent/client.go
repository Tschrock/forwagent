package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/dainnilsson/forwagent/common"
)

type dialServer = func() (net.Conn, error)

func createListener(name string) (net.Listener, error) {
	os.Remove(name)
	time.Sleep(100 * time.Millisecond)
	return net.Listen("unix", name)
}

func handleConnections(listener net.Listener, kind string, dial dialServer) {
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
		} else {
			go handleConnection(dial, conn, kind)
		}
	}

}

func handleConnection(dial dialServer, conn net.Conn, connType string) {
	defer conn.Close()

	serverConn, err := dial()
	if err != nil {
		fmt.Println("Error connecting to server:", err.Error())
		return
	}

	io.WriteString(serverConn, connType)

	common.ProxyConnections(conn, serverConn)
}

func runClient(gpgSocketPath string, sshSocketPath string, dial dialServer) error {
	gpgSocket, err := createListener(gpgSocketPath)
	if err != nil {
		return err
	}
	sshSocket, err := createListener(sshSocketPath)
	if err != nil {
		return err
	}

	go handleConnections(gpgSocket, "GPG", dial)
	handleConnections(sshSocket, "SSH", dial)

	return nil
}
