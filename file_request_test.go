package main

import (
	"io"
	"net"
	"strings"
	"testing"
)

func TestHandleConnection(t *testing.T) {

	initClient()
	initLoggers()

	tt := []struct {
		name     string
		filename string
		expected string
	}{
		{"file exists", "castigate", "ACK\\npeerOne\\npeerTwo"},
		{"file not exist", "does_not_exist", "NACK\\nfile not found"},
		{"file exists no peers", "masticate", "NACK\\nNo Peers With File"},
	}

	SetPeerActive("peerOne")
	SetPeerActive("peerTwo")
	AddPeerToFile("peerOne", "castigate")
	AddPeerToFile("peerTwo", "castigate")
	AddPeerToFile("peerThree", "masticate")
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			one, two := net.Pipe()

			go handleConnection(two)

			one.Write([]byte(tc.filename + "\n"))

			response, err := io.ReadAll(one)

			formatted := strings.Replace(string(response), "\n", "\\n", -1)

			if err != nil {
				t.Errorf(" Want '%s', got '%s'", tc.expected, err.Error())
			}

			if formatted != tc.expected {
				t.Errorf("Want '%s', got '%s'", tc.expected, formatted)
			}

		})
	}
}

