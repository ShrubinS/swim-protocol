package main

import (
	"flag"

	"github.com/shrubins/swim-protocol/node"
)

var (
	host          = flag.String("host", "localhost", "The hostname or IP to connect to; defaults to \"localhost\".")
	port          = flag.Int("port", 8000, "The port to connect to; defaults to 8000.")
	creategroup   = flag.Bool("new", false, "Create a new group")
	existingGroup = flag.String("group", "", "Existing group to join")
)

func main() {
	flag.Parse()
	member := node.NewNode(*host, *port, *creategroup)
	done := make(chan bool, 1)
	member.StartNode(done, *existingGroup)
	<-done
}

// Start multiple nodes
// func startServers(nodes []node.Node) {
// 	done := make(chan bool, len(nodes))
// 	for _, n := range nodes {
// 		go n.StartNode(done)
// 	}
// 	for i := 0; i < len(nodes); i++ {
// 		select {
// 		case <-done:
// 			fmt.Println("Node", i+1, "Terminated")
// 		}
// 	}
// }
