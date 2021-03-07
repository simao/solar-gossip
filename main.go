package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/simao/solar-gossip/lib"
)

func findBestNode(current *memberlist.Node, list *memberlist.Memberlist) *memberlist.Node {
	var max uint64
	var best *memberlist.Node
	for _, member := range list.Members() {
		voltage := binary.LittleEndian.Uint64(member.Meta)

		if voltage > max {
			best = member
			max = voltage
		}
	}

	if current == nil || best != current {
		if best == list.LocalNode() {
			fmt.Println("I am node with most voltage, update dns")
		} else {
			fmt.Printf("%s has most voltage\n", best.Name)
		}
	}

	return best
}

func main() {
	fmt.Println(len(os.Args))

	var port int
	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1])
	} else {
		port = 7000
	}

	var remotePort int
	if len(os.Args) > 2 {
		remotePort, _ = strconv.Atoi(os.Args[2])
	} else {
		remotePort = 7000
	}

	var mymeta lib.DeviceMeta
	if len(os.Args) > 3 {
		arg, _ := strconv.Atoi(os.Args[3])
		mymeta = lib.DeviceMeta{Voltage: uint64(arg)}
	} else {
		mymeta = lib.DeviceMeta{Voltage: uint64(rand.Int31n(1000))}
	}

	listEventsCh := make(chan memberlist.NodeEvent, 10)

	config := memberlist.DefaultLocalConfig()
	config.BindPort = port
	config.Name = fmt.Sprintf("asterix:%d", port)
	config.Events = &memberlist.ChannelEventDelegate{Ch: listEventsCh}
	config.Delegate = &mymeta

	fmt.Printf("port=%d, remote_port=%d", port, remotePort)

	list, err := memberlist.Create(config)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	if port != remotePort { // I am not first node
		addr := fmt.Sprintf("0.0.0.0:%d", remotePort)
		fmt.Println("Connecting to cluster using ", addr)

		_, err := list.Join([]string{addr})
		if err != nil {
			panic("Failed to join cluster: " + err.Error())
		}
	} else { // I am first node
		fmt.Println("Waiting for clients on port", port)
	}

	tick := time.Tick(10 * time.Second)
	var currentBest *memberlist.Node

	for {
		select {
		case <-tick:
			lib.PrintMembers(list)

		case e := <-listEventsCh:
			fmt.Println("")
			fmt.Println("Received", e)

			currentBest = findBestNode(currentBest, list)

			lib.PrintMembers(list)

		default:
			fmt.Print(".")
			time.Sleep(time.Second)
		}
	}
}
