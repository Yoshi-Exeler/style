package main

import "fmt"

// this is just a little playground to test examples

func main() {
	packet := []byte{1, 2, 3}

	connTCP := AbstractConnection{
		Sender: &TCPSender{},
		Closer: &QuickCloser{},
	}

	connTCP.Send(packet)
	connTCP.Close()

	connUDP := AbstractConnection{
		Sender: &UPDSender{},
		Closer: &SafeCloser{},
	}

	connUDP.Send(packet)
	connUDP.Close()

	grpcConn := &GRPCConnection{
		TCPSender: &TCPSender{},
		Closer:    &QuickCloser{},
	}

	closeConn(connTCP)

	closeConn(connUDP)

	closeConn(grpcConn)
}

func closeConn(c Connection) {
	c.Close()
}

type GRPCConnection struct {
	*TCPSender
	Closer
}

// trait definition
type Sender interface {
	Send(buff []byte) error
}

// trait definition
type Closer interface {
	Close() error
}

// trait implementation A
type TCPSender struct{}

func (tcps *TCPSender) Send(buff []byte) error {
	fmt.Println("tcp send")
	return nil
}

// trait implementation B
type UPDSender struct{}

func (upds *UPDSender) Send(buff []byte) error {
	fmt.Println("udp send")
	return nil
}

// trait implementation C
type SafeCloser struct{}

func (sc *SafeCloser) Close() error {
	fmt.Println("safe close")
	return nil
}

type QuickCloser struct{}

func (qc *QuickCloser) Close() error {
	fmt.Println("quick close")
	return nil
}

// concrete entity A
type Connection interface {
	Sender
	Closer
}

type AbstractConnection struct {
	Sender
	Closer
}
