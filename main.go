package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"strings"
	"time"
)

// UDP_SERVER_PORT represents used to listen for udp connections
const UDP_SERVER_PORT = ":4911"

// TCP_SERVER_PORT represents the port used to listen for incoming TCP connections
const TCP_SERVER_PORT = ":4912"

// PEER_TTL represents the maximum time a peer has to contact the server before it is deemed inactive
const PEER_TTL = 60

func main() {

	initClient()
	initLoggers()

	go ResetPeerSet()
	go RunUDP()
	runTCP()
}

// runTCP is a function that runs a TCP server capable of handling requests for a file
func runTCP() {

	//create the listener on the port
	listener, err := net.Listen("tcp", TCP_SERVER_PORT)

	defer func() {

		err = listener.Close()

	}()

	//if there is an error starting the server log the error and shutdown
	if err != nil {
		log.Fatal(err)
	}

	//create infinite loop to listen for incoming connections
	for {

		conn, err := listener.Accept()

		//if there is an error accepting the connection print the error and log it
		if err != nil {
			ErrorLogger.Println(err)
			continue
		}

		//handle the connection in a separate go routine
		go handleConnection(conn)
	}

}

// ResetPeerSet is a function that runs continuously in a go routine and periodically sets peers as inactive
func ResetPeerSet() {

	for {
		time.Sleep(PEER_TTL * time.Second)
		ClearPeerSet()
		DeletionLogger.Println("Deleted all peers from peer set")
	}

}

// RunUDP runs continuously in a go routine and handles incoming UDP messages
func RunUDP() {

	//get the udp address
	address, err := net.ResolveUDPAddr("udp4", UDP_SERVER_PORT)

	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp4", address)

	if err != nil {
		log.Fatal(err)
	}

	//for every incoming connection handle it
	for {
		handleUDPConnection(conn)
	}

}

// handleUDPConnection handles the incoming UDP messages
func handleUDPConnection(conn *net.UDPConn) {

	//make a buffer to store the message
	buffer := make([]byte, 1024)

	//read the incoming message into the buffer
	n, addr, err := conn.ReadFromUDP(buffer)

	ConnectionLogger.Printf("Incoming message from %s\n", addr.IP.String())

	//extract the body data by converting bytes to a string and splitting at line separators
	bodyData := strings.Split(string(buffer[:n]), "\n")

	//extract peer data stored on the first line
	peerInfo := bodyData[0]

	//create peer struct
	var peer Peer

	//marshal json data into the peer structure
	json.Unmarshal([]byte(peerInfo), &peer)

	//add address information
	peer.Address = addr.IP.String()

	res, err := json.Marshal(peer)

	//set the peer that made the connection as active
	SetPeerActive(string(res))

	if len(bodyData) == 1 {
		conn.WriteToUDP([]byte("ACK"), addr)
		return
	}

	//filenames should be stored in all the subsequenc lines
	filenames := bodyData[1:]

	//add the peers to the set
	for _, filename := range filenames {
		IndexingLogger.Printf("Adding %s to %s set\n", string(res), filename)
		AddPeerToFile(string(res), filename)
	}

	_, err = conn.WriteToUDP([]byte("ACK"), addr)

	if err != nil {
		ErrorLogger.Println(err)
		return
	}
}

// handleConnection is a function that handles in coming TCP requests for a file
func handleConnection(conn net.Conn) {

	defer func() {

		conn.Close()

	}()

	reader := bufio.NewScanner(conn)

	if !reader.Scan() {
		ErrorLogger.Println(reader.Err())
		conn.Write([]byte("NACK\nFormat Error"))
		return
	}

	line := reader.Text()

	ConnectionLogger.Printf("Received request from %s for %s\n", conn.RemoteAddr().String(), line)

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
