package scripting

import (
	"fmt"
	"io/ioutil"
	"net"
)

type TCPClient struct {
	hostname string
	port     uint16
	client   *net.TCPConn
}

func NewTCPClient(hostname string, port uint16) TCPClient {
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", hostname, port))
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	return TCPClient{hostname, port, conn}
}

func (client TCPClient) Send(buffer []byte) int {
	if client.client != nil {
		written, _ := client.client.Write(buffer)
		return written
	}
	return -1
}

func (client TCPClient) Recv() []byte {
	if client.client != nil {
		result, _ := ioutil.ReadAll(client.client)
		return result
	}
	return make([]byte, 0)
}
