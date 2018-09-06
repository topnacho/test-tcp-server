package main

import (
	"encoding/hex"
	"fmt"
	"github.com/rs/zerolog/log"
	"net"
	"os"
	"time"
)

const (
	ConnPortStandard = "3333"
	ConnPortEcho = "6666"
	ConnType = "tcp"
)

func main() {
	//Listen for incoming connections standard.
	l, err := net.Listen(ConnType, ":"+ConnPortStandard)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + ":" + ConnPortStandard)

	// Listen for incoming connections standard.
	l1, err := net.Listen(ConnType, ":"+ConnPortEcho)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l1.Close()
	fmt.Println("Listening on " + ":" + ConnPortEcho)


	go acceptLoop(l, ConnPortStandard)
	acceptLoop(l1, ConnPortEcho)


}

func acceptLoop(l net.Listener, port string) {
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error().Err(err)
			return
		}

		log.Info().Msgf("******** Handling new connection from %s on port %s", conn.RemoteAddr(), port)

		var echo bool
		switch port {
		case ConnPortEcho:
			echo = true
		case ConnPortStandard:
			echo = false
		default:
			log.Error().Msgf("unknown port %s", port)
		}

		go handleRequest(conn, echo)

	}
}


// Handles incoming requests and returns a simple ACK
func handleRequest(conn net.Conn, echo bool) {

	var ACK = byte(0xAA)

	// close the connection when this function ends
	defer func() {
		log.Info().Msgf("******** Closing connection from %s", conn.RemoteAddr())
		conn.Close()
	}()

	// set a timeout
	timeoutDuration := 5 * time.Second
	buf := make([]byte, 1024)

	for {
		// Set a deadline for reading. Read operation will fail if no data is received after deadline
		conn.SetReadDeadline(time.Now().Add(timeoutDuration))

		bytesRead, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				return
			}
			log.Error().Msgf("******** error while reading from %s [%s]", conn.RemoteAddr(), err.Error())
			return
		}
		
		print(hex.Dump(buf[0:bytesRead]))

		if echo {
			conn.Write(buf[0:bytesRead])

		} else {
			// Send ACK
			ack := []byte{ACK}
			conn.Write(ack)
		}

	}
}


