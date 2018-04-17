package node

import (
	"bufio"
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
			group.NewGroup(newGroupUUID, 20),
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
		if err != nil {
			// handle error
		}
		defer conn.Close()
		go node.handleConnection(conn)
	}
}

func (node Node) handleConnection(conn net.Conn) {
	conn.Write([]byte("Ack from " + node.nodeID.String()))
}

func (node Node) ping(existingGroup string) {
	var dest string
	nodeGroupView := node.groupView
	if nodeGroupView == nil || len(nodeGroupView.View) == 0 {
		if len(existingGroup) == 0 {
			fmt.Println("Found no groups to ping")
			return
		}
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
		log.Fatal("Error: ", err)
		return
	}
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	fmt.Println(scanner.Text())
	// log success
	// log.Println("Connected to ", conn.RemoteAddr())

}
