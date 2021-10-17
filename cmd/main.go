package main

import (
	"client/internal/client"

	"github.com/sirupsen/logrus"
)

func main() {
	err := client.Connect("tcp", "localhost:27001")
	if err != nil {
		logrus.Errorf("connection problem: %q", err)
		return
	}
}
