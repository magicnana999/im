package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

func Start() {

	wg.Add(1)

	brokerAddr := `127.0.0.1:7539`

	tcpAddr, err := net.ResolveTCPAddr("tcp4", brokerAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	go heartbeat(conn)
	wg.Wait()
}

func heartbeat(c *net.TCPConn) {

	timer := time.NewTimer(time.Second * 30)
	defer timer.Stop()

	buffer := new(bytes.Buffer)

	binary.Write(buffer, binary.BigEndian, 1)

	binary.Write(buffer, binary.BigEndian, 12)

	c.Write(buffer.Bytes())
}
