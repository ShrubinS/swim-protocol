package node

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/shrubins/swim-protocol/group"
)

const (
	protocolPeriod = 1000 * time.Millisecond
	timeout        = 300 * time.Millisecond
)

type Node struct {
	nodeID    uuid.UUID
	host      string
	port      int
	groupView *group.Group
}

func NewNode(host string, port int, createGroup bool) Node {
	newNodeUUID, err := uuid.NewRandom()
	if err != nil {
		log.Fatal("Error generating Node UUID:\n" + err.Error())
	}
	newGroupUUID, err := uuid.NewRandom()
	if err != nil {
		log.Fatal("Error generating Group UUID:\n" + err.Error())
	}

	if createGroup {
		return Node{
			newNodeUUID,
			host,
			port,
			group.NewGroup(newGroupUUID, fmt.Sprintf("%s:%s", host, strconv.Itoa(port))),
		}
	}

	return Node{
		newNodeUUID,
		host,
		port,
		nil,
	}
}

func (node Node) StartNode(done chan bool, existingGroup string) {
	ticker := time.NewTicker(protocolPeriod)
	go func() {
		for range ticker.C {
			node.ping(existingGroup)
		}
	}()
	node.listen()
	time.Sleep(2 * time.Second)
	ticker.Stop()
	fmt.Println("Terminating Node...")
	done <- true
}

func (node Node) processHeartBeat() {
	fmt.Println("Heartbeat", node.nodeID.String())
}

func (node Node) listen() {
	host := node.host + ":" + strconv.Itoa(node.port)
	ln, err := net.Listen("tcp", host)
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		// fmt.Println("conn", conn.RemoteAddr())
		if err != nil {
			// handle error
		}
		defer conn.Close()
		go node.handleConnection(conn)
	}
}

func (node Node) handleConnection(conn net.Conn) {
	fmt.Println("ping from", conn.RemoteAddr())
	sendMessage := Message{
		node.groupView.GroupID,
		node.nodeID,
		node.groupView.View,
	}
	write, err := json.Marshal(sendMessage)
	if err != nil {
		log.Println("JSON Marshal error:", err)
	}
	conn.Write(write)
}

func (node Node) ping(existingGroup string) {
	var dest string
	nodeGroupView := node.groupView
	fmt.Println("node groupview before", nodeGroupView)
	if nodeGroupView == nil || len(nodeGroupView.View) == 0 {
		if len(existingGroup) == 0 {
			// fmt.Println("Found no groups to ping")
			return
		}
		fmt.Println("joining existing group")
		dest = existingGroup
	} else {
		fmt.Println("existing group is", existingGroup)
		fmt.Println("node group view has ", len(nodeGroupView.View))
		index := rand.Intn(len(nodeGroupView.View))
		dest = nodeGroupView.View[index]
	}
	fmt.Println("Pinging", dest)
	conn, err := net.DialTimeout("tcp", dest, timeout)
	if err != nil {
		// log error
		log.Println("Error: ", err)
		return
	}

	go readConnection(conn, &nodeGroupView)
}

func readConnection(conn net.Conn, groupView **group.Group) {
	defer func() {
		fmt.Println("Closing connection...")
		conn.Close()
	}()
	decoder := json.NewDecoder(conn)
	var receivedMessage Message
	err := decoder.Decode(&receivedMessage)
	if err != nil {
		fmt.Println(err)
	}
	gv := group.MakeGroup(receivedMessage.GroupID, receivedMessage.Members)
	*groupView = gv
	fmt.Println(receivedMessage.Members)
}

func handleConnection(conn net.Conn) {
	fmt.Println("Handling new connection...")

	// Close connection when this function ends
	defer func() {
		fmt.Println("Closing connection...")
		conn.Close()
	}()

	timeoutDuration := 5 * time.Second
	bufReader := bufio.NewReader(conn)

	for {
		// Set a deadline for reading. Read operation will fail if no data
		// is received after deadline.
		conn.SetReadDeadline(time.Now().Add(timeoutDuration))

		// Read tokens delimited by newline
		bytes, err := bufReader.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%s", bytes)
	}
}
