package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const UDP_SERVER_PORT = ":4911"
const TCP_SERVER_PORT = ":4912"

func main() {

	initClient()
	initLoggers()

	go ResetPeerSet()
	go RunUDP()
	runTCP()
}

func runTCP() {

	//create the tcp listener
	listener, err := net.Listen("tcp", TCP_SERVER_PORT)

	defer listener.Close()

	if err != nil {
		fmt.Println("Error: ", err)
	}

	for {

		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error: err ")
			continue
		}

		go handleConnection(conn)
	}

}

func ResetPeerSet() {

	for {
		time.Sleep(60 * time.Second)
		ClearPeerSet()
		DeletionLogger.Println("Deleted all peers from peer set")
	}

}

func RunUDP() {

	address, err := net.ResolveUDPAddr("udp4", UDP_SERVER_PORT)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp4", address)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	for {
		handleUDPConnection(conn)
	}

}

func handleUDPConnection(conn *net.UDPConn) {
	buffer := make([]byte, 1024)

	n, addr, err := conn.ReadFromUDP(buffer)

	filenames := strings.Split(string(buffer[:n]), "\n")

	ConnectionLogger.Printf("Incoming message from %s\n", addr.IP.String())

	//set the peer that made the connection as active
	SetPeerActive(addr.IP.String())

	if len(filenames) == 1 {
		conn.WriteToUDP([]byte("ACK"), addr)
		return
	}

	//add the peer to files peer set
	for _, filename := range filenames {
		IndexingLogger.Printf("Adding %s to %s set\n", addr.IP.String(), filename)
		AddPeerToFile(addr.IP.String(), filename)
	}

	conn.WriteToUDP([]byte("ACK"), addr)

	if err != nil {
		log.Println(err)
		return
	}
}

// this function handles all the tcp logic
func handleConnection(conn net.Conn) {

	defer conn.Close()

	reader := bufio.NewScanner(conn)

	if !reader.Scan() {
		fmt.Println("Error: ", reader.Err())
		conn.Write([]byte("NACK\nFormat Error"))
		return
	}

	line := reader.Text()

	ConnectionLogger.Printf("Received request from %s for %s\n", conn.LocalAddr().String(), line)

	peers, err := GetPeersWithFile(line)

	//write that the file does not exist
	if err != nil {
		conn.Write([]byte("NACK\n" + err.Error()))
		return
	}

	//write that there are no users with the file
	if len(peers) == 0 {
		conn.Write([]byte("NACK\nNo Peers With File"))
		return
	}

	conn.Write([]byte("ACK\n"))
	conn.Write([]byte(strings.Join(peers, "\n")))

}
