package main

import (
	"fmt"

	"github.com/shrubins/swim-protocol/node"
)

func main() {
	nodes := []node.Node{
		node.NewLeaderNode(),
		node.NewLeaderNode(),
		node.NewLeaderNode(),
		node.NewLeaderNode(),
	}

	startServers(nodes)
}

func startServers(nodes []node.Node) {
	done := make(chan bool, len(nodes))
	for _, n := range nodes {
		go n.StartNode(done)
	}
	for i := 0; i < len(nodes); i++ {
		select {
		case <-done:
			fmt.Println("Node", i+1, "Terminated")
		}
	}
}
