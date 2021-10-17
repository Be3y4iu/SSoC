package client

import (
	"bufio"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

// BUFFERSIZE is default size of chunk of file that we send
var BUFFERSIZE int64 = 1024

// Connect func connects to the server and provide list of commands
func Connect(network, address string) error {
	conn, err := connect(network, address)
	if err != nil {
		return fmt.Errorf("cannot connect to the server %q", err)
	}

	text, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return fmt.Errorf("cannor read from server")
	}
	fmt.Println(text)
	for {
		command, err := getCommand()
		if err != nil {
			logrus.Warn(err)
			return err
		}

		logrus.Info("command is ", command)
		if err != nil {
			return err
		}
		logrus.Info("trying to write")
		_, err = conn.Write([]byte(command))
		if err != nil {
			logrus.Warn("trouble")
			return fmt.Errorf("cannot write to the connection %q", err)
		}

		provideCommand(command, conn)
	}
	return nil
}

func connect(network, address string) (connection net.Conn, err error) {
	connection, err = net.Dial(network, address)
	return
}
