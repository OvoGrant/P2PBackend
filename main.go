package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const UDP_SERVER_PORT = "4911"
const TCP_SERVER_PORT = ":4912"

func main() {
	initClient()

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

		log.Println("Incoming request from ", conn)
		go handleConnection(conn)
	}

}
func RunUDP() {

	address, err := net.ResolveUDPAddr("udp4", ":"+UDP_SERVER_PORT)

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

		buffer := make([]byte, 1024)

		_, addr, err := conn.ReadFromUDP(buffer)

		filenames := strings.Split(string(buffer), "\n")

		//set the peer that made the connection as active
		SetPeerActive(addr.IP.String())

		//add the peer to file's peer set
		for _, filename := range filenames {
			AddPeerToFile(addr.IP.String(), filename)
		}

		conn.WriteToUDP([]byte("ACK\n"), addr)

		if err != nil {
			log.Println(err)
			return
		}

	}

}

// this function handles all the tcp logic
func handleConnection(conn net.Conn) {

	reader := bufio.NewReader(conn)

	line, err := reader.ReadString('\n')

	wt := bufio.NewWriter(conn)

	if err != nil {
		fmt.Println("Error: ", err)
		wt.Write([]byte("NACK\n"))
		wt.Write([]byte("Format Error\n"))
		return
	}

	peers, err := GetPeersWithFile(line)

	//write that the file does not exist
	if err != nil {
		wt.Write([]byte("NACK\n"))
		wt.Write([]byte("File Does Not Exist\n"))
		return
	}

	//write that there are no users with the file
	if len(peers) == 0 {
		wt.Write([]byte("NACK\n"))
		wt.Write([]byte("No Peers With File\n"))
		return
	}

	wt.Write([]byte("ACK\n"))
	wt.Write([]byte(strings.Join(peers, "\n")))

}
