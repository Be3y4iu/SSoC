package upload

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	BUFFERSIZE = 1024
	Megabyte   = 1000000.0
)

func UploadFile(conn net.Conn, fileName string) error {
	logrus.Info(strings.TrimSpace(fileName))
	file, err := os.Open(strings.TrimSpace(fileName))
	if err != nil {
		return fmt.Errorf("cannot open file %q", err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("cannot get file info %q", err)
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName = fillString(fileInfo.Name(), 64)
	fmt.Println("Sending filename and filesize")
	conn.Write([]byte(fileSize))
	conn.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file!")
	startTime := time.Now()
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		conn.Write(sendBuffer)
	}
	downloadTime := time.Since(startTime).Nanoseconds()
	// todo fix speed
	logrus.Infof("speed is %.5f Mb/seconds", float64(fileInfo.Size())/float64(downloadTime)/Megabyte*1000000000)
	logrus.Info("file has been sent, closing connection!")
	return nil
}

func fillString(returnString string, toLength int) string {
	for {
		stringLen := len(returnString)
		if stringLen < toLength {
			returnString = returnString + ":"
			continue
		}
		break
	}
	return returnString
}
