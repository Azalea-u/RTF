package utils

import (
	"fmt"
	"log"
	"net"
)

// GeneratePort generates a random port number.
func GeneratePort() string {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("\033[31mError:\033[0m" + " Error listening on port: " + err.Error())
	}
	defer l.Close()
	return fmt.Sprintf("%d", l.Addr().(*net.TCPAddr).Port)
}
