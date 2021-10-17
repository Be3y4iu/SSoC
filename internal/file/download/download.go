package download

import (
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// BUFFERSIZE is default size of chunk of file that we send
const (
	BUFFERSIZE = 1024
	Megabyte   = 1000000.0
)

func Download(connection net.Conn) error {
	logrus.Info("Connected to server, start receiving the file name and file size")
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	newFile, err := os.Create(fileName)

	if err != nil {
		return err
	}
	defer newFile.Close()
	var receivedBytes int64
	startTime := time.Now()
	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			_, err := io.CopyN(newFile, connection, fileSize-receivedBytes)
			if err != nil {
				return err
			}
			connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}

		n, err := io.CopyN(newFile, connection, BUFFERSIZE)
		if err != nil {
			return err
		}
		receivedBytes += n
	}
	downloadTime := time.Since(startTime).Seconds()
	logrus.Infof("speed is %.2f Mb/seconds", float64(receivedBytes)/downloadTime/Megabyte)

	logrus.Info("Received file completely!")

	return nil
}
