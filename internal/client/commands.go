package client

import (
	"bufio"
	"client/internal/file/download"
	"client/internal/file/upload"
	"fmt"
	"net"
	"os"
	"strings"
)

func provideCommand(command string, conn net.Conn) error {
	var err error
	switch {
	case strings.Contains(command, "echo"):
		echoCommand(conn)
	case strings.Contains(command, "time"):
		timeCommand(conn)
	case strings.Contains(command, "close"):
	case strings.Contains(command, "download"):
		err = download.Download(conn)
	case strings.Contains(command, "upload"):
		err = upload.UploadFile(conn, strings.TrimLeft(strings.ToLower(command), "upload"))
	}
	return err
}

func timeCommand(conn net.Conn) {
	bytes := make([]byte, 128)
	conn.Read(bytes)
	fmt.Printf("%s\n", bytes)
}

func echoCommand(conn net.Conn) {
	bytes := make([]byte, 128)
	conn.Read(bytes)
	fmt.Printf("%s\n", bytes)
}

func getCommand() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	var command string
	var err error
loop:
	for {
		fmt.Print("Enter command: ")
		command, err = reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("cannot read string from stdin")
		}
		err = checkCommand(command)
		if err != nil {
			fmt.Print("Unknown command. Try this commands: \n1.echo 'text'" +
				"\n2.time" +
				"\n3.download + 'filename'" +
				"\n4.upload + 'filename'" +
				"\n5.close\n")
		} else {
			break loop
		}
	}

	return command, nil
}

func checkCommand(command string) error {
	switch {
	case strings.Contains(command, "echo"):
		return nil
	case strings.EqualFold(command, "time\n"):
		return nil
	case strings.Contains(command, "download"):
		return nil
	case strings.Contains(command, "upload"):
		return nil
	case strings.EqualFold(command, "close\n"):
		return nil
	default:
		return fmt.Errorf("unknown command")
	}
}
