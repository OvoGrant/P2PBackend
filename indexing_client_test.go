package main

import (
	"net"
	"testing"
)

func TestUDPConnection(t *testing.T) {

	tt := []struct {
		name     string
		request  []byte
		expected []byte
	}{
		{"test for ack", []byte("123.0.0.1"), []byte("ACK")},
	}

	protocol := "udp"
	port := ":5234"

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			addr, err := net.ResolveUDPAddr(protocol, port)

			conn, err := net.DialUDP("udp", nil, addr)

			if err != nil {

			}

			conn.Write(tc.request)

			defer conn.Close()

			handleUDPConnection(conn)

		})
	}
}
