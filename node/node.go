package node

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/shrubins/swim-protocol/message"

	"github.com/shrubins/swim-protocol/group"

	"github.com/google/uuid"
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

	member := message.Member{
		NodeID:  newNodeUUID,
		Address: fmt.Sprintf("%s:%s", host, strconv.Itoa(port)),
	}

	if createGroup {

		return Node{
			newNodeUUID,
			host,
			port,
			group.NewGroup(newGroupUUID, member),
		}
	}

	return Node{
		newNodeUUID,
		host,
		port,
		group.NotGroup(member),
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
	// fmt.Println("ping from", conn.RemoteAddr())
	decoder := json.NewDecoder(conn)
	var receivedMessage message.Message
	err := decoder.Decode(&receivedMessage)
	if err != nil {
		fmt.Println("Server says: No message recieved", err)
	}
	fmt.Println("ping received")

	receivedMessage.String()

	for _, m := range receivedMessage.Members {
		node.groupView.AddToGroup(m)
	}

	sendMessage := message.Message{
		GroupID: node.groupView.GroupID,
		Members: node.groupView.View,
	}
	write, err := json.Marshal(sendMessage)
	if err != nil {
		log.Println("JSON Marshal error:", err)
	}
	conn.Write(write)
}

func (node Node) ping(existingGroup string) {
	var dest string
	updatedGroupView := make(chan *group.Group, 1)
	nodeGroupView := node.groupView
	if nodeGroupView == nil || len(nodeGroupView.View) == 0 || !nodeGroupView.RealGroup {
		if len(existingGroup) == 0 {
			// fmt.Println("Found no groups to ping")
			return
		}
		fmt.Println("joining existing group", existingGroup)
		dest = existingGroup
	} else {
		fmt.Println("node group view has ", len(nodeGroupView.View), "members")
		index := rand.Intn(len(nodeGroupView.View))
		dest = nodeGroupView.View[index].Address
	}
	conn, err := net.DialTimeout("tcp", dest, timeout)
	if err != nil {
		// log error
		log.Println("Error: ", err)
		return
	}

	go node.readConnection(conn, updatedGroupView)

	update := <-updatedGroupView
	*node.groupView = *update
}

func (node Node) readConnection(conn net.Conn, groupView chan *group.Group) {
	defer func() {
		fmt.Println("Closing connection...")
		conn.Close()
	}()

	sendMessage := message.Message{
		GroupID: node.groupView.GroupID,
		Members: node.groupView.View,
	}

	fmt.Println("Pinging message...")
	sendMessage.String()

	write, err := json.Marshal(sendMessage)
	if err != nil {
		log.Println("JSON Marshal error:", err)
	}
	conn.Write(write)

	decoder := json.NewDecoder(conn)
	var receivedMessage message.Message
	err = decoder.Decode(&receivedMessage)
	if err != nil {
		fmt.Println(err)
	}
	gv := group.MakeGroup(receivedMessage.GroupID, receivedMessage.Members)
	groupView <- gv
	fmt.Println(receivedMessage.Members)
}

// func handleConnection(conn net.Conn) {
// 	fmt.Println("Handling new connection...")

// 	// Close connection when this function ends
// 	defer func() {
// 		fmt.Println("...Closing connection...")
// 		conn.Close()
// 	}()

// 	timeoutDuration := 5 * time.Second
// 	bufReader := bufio.NewReader(conn)

// 	for {
// 		// Set a deadline for reading. Read operation will fail if no data
// 		// is received after deadline.
// 		conn.SetReadDeadline(time.Now().Add(timeoutDuration))

// 		// Read tokens delimited by newline
// 		bytes, err := bufReader.ReadBytes('\n')
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}

// 		fmt.Printf("%s", bytes)
// 	}
// }
