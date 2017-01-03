// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.
package main

import (
	"flag"
	"fmt"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/messages"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/engine"
	"github.com/FactomProject/factomd/p2p"
	"math/rand"
	"strings"
	"sync"
	"time"
	"os"
)

var p2pProxy *engine.P2PProxy

var old map[[32]byte]interfaces.IMsg
var oldsync sync.Mutex

var oldcnt int

var broadcastSent int
var broadcastReceived int
var p2pSent int
var p2pReceived int
var p2pRequestSent int
var p2pRequestReceived int

var name string
var isp2p bool
var numStamps int
var numReplies int
var size int

func InitNetwork() {

	go engine.StartProfiler()

	namePtr := flag.String("name", fmt.Sprintf("%d", rand.Int()), "Name for this node")
	networkPortOverridePtr := flag.String("networkPort", "8108", "Address for p2p network to listen on.")
	peersPtr := flag.String("peers", "", "Array of peer addresses. ")
	netdebugPtr := flag.Int("netdebug", 0, "0-5: 0 = quiet, >0 = increasing levels of logging")
	exclusivePtr := flag.Bool("exclusive", false, "If true, we only dial out to special/trusted peers.")
	deadlinePtr := flag.Int64("deadline", 1, "Deadline for Reads and Writes to conn.")
	p2pPtr := flag.Bool("p2p", false, "Test p2p messages (default to false)")
	numStampsPtr := flag.Int("numstamps", 1, "Number of timestamps per reply on p2p test. (makes messages big)")
	numReplysPtr := flag.Int("numreplies", 1, "Number of replies to any request")
	sizePtr := flag.Int("size", 0, "size.  We will add a payload of random data of this many K.")

	flag.Parse()

	numReplies = *numReplysPtr
	numStamps = *numStampsPtr
	name = *namePtr
	port := *networkPortOverridePtr
	peers := *peersPtr
	netdebug := *netdebugPtr
	exclusive := *exclusivePtr
	p2p.NetworkDeadline = time.Duration(*deadlinePtr) * time.Millisecond
	isp2p = *p2pPtr
	size = *sizePtr*1024


	os.Stderr.WriteString("\nnetTest is a standalone program that generates factomd messages (bounce and bounceResponse)\n ")
	os.Stderr.WriteString("        and sends them to other nodes on the network.  This allows testing of the network\n ")
	os.Stderr.WriteString("        without running all of factomd.  Note you can control the size of messages and other\n ")
	os.Stderr.WriteString("        variables like the deadline used in the network, and p2p testing.\n\n")

	os.Stderr.WriteString("Settings\n")
	os.Stderr.WriteString(fmt.Sprintf("%20s -- %s\n","name",name))
	os.Stderr.WriteString(fmt.Sprintf("%20s -- %s\n","networkPort",port))
	os.Stderr.WriteString(fmt.Sprintf("%20s -- %s\n","peers",peers))
	os.Stderr.WriteString(fmt.Sprintf("%20s -- %d\n","netdebug",netdebug))
	os.Stderr.WriteString(fmt.Sprintf("%20s -- %v\n","exclusive",exclusive))
	os.Stderr.WriteString(fmt.Sprintf("%20s -- %d\n","deadline",p2p.NetworkDeadline.Seconds()))
	os.Stderr.WriteString(fmt.Sprintf("%20s -- %v\n","p2p",isp2p))
	os.Stderr.WriteString(fmt.Sprintf("%20s -- %dk\n\n","size", size))


	old = make(map[[32]byte]interfaces.IMsg, 0)
	connectionMetricsChannel := make(chan interface{}, p2p.StandardChannelSize)
	ci := p2p.ControllerInit{
		Port:                     port,
		PeersFile:                "peers.json",
		Network:                  1,
		Exclusive:                exclusive,
		SeedURL:                  "",
		SpecialPeers:             peers,
		ConnectionMetricsChannel: connectionMetricsChannel,
	}
	p2pNetwork := new(p2p.Controller).Init(ci)
	p2pNetwork.StartNetwork()
	// Setup the proxy (Which translates from network parcels to factom messages, handling addressing for directed messages)
	p2pProxy = new(engine.P2PProxy).Init("testnode", "P2P Network").(*engine.P2PProxy)
	p2pProxy.FromNetwork = p2pNetwork.FromNetwork
	p2pProxy.ToNetwork = p2pNetwork.ToNetwork
	p2pProxy.SetDebugMode(netdebug)

	if netdebug > 0 {
		p2pNetwork.StartLogging(uint8(netdebug))
	} else {
		p2pNetwork.StartLogging(uint8(0))
	}
	p2pProxy.StartProxy()
	// Command line peers lets us manually set special peers
	p2pNetwork.DialSpecialPeersString("")
}

var cntreq int32
var cntreply int32

// Returns true if message is new
func MsgIsNew(msg interfaces.IMsg) bool {
	oldsync.Lock()
	defer oldsync.Unlock()
	return old[msg.GetHash().Fixed()] == nil
}

func SetMsg(msg interfaces.IMsg) {
	oldsync.Lock()
	old[msg.GetHash().Fixed()] = msg
	oldsync.Unlock()
}

func listen() {

	for {
		msg, err := p2pProxy.Recieve()
		if err != nil || msg == nil {
			time.Sleep(1 * time.Millisecond)
			continue
		}
		time.Sleep(1 * time.Millisecond)

		bounce, ok1 := msg.(*messages.Bounce)
		bounceReply, ok2 := msg.(*messages.BounceReply)

		if MsgIsNew(msg) {
			SetMsg(msg)

			fmt.Println("    ", msg.String())

			if ok1 && len(bounce.Stamps) < 5 {
				if isp2p {
					for i := 0; i < numReplies; i++ {
						bounceReply = new(messages.BounceReply)
						bounceReply.SetPeer2Peer(true)

						bounceReply.Number = cntreply
						cntreply++
						bounceReply.Name = name + "->" + strings.TrimSpace(bounce.Name)

						bounceReply.Timestamp = bounce.Timestamp
						bounceReply.Stamps = append(bounceReply.Stamps, bounce.Stamps...)

						for j := 0; j < numStamps; j++ {
							bounceReply.Stamps = append(bounceReply.Stamps, primitives.NewTimestampNow())
						}

						bounceReply.SetOrigin(bounce.GetOrigin())
						bounceReply.SetNetworkOrigin(bounce.GetNetworkOrigin())

						SetMsg(msg)
						p2pProxy.Send(bounceReply)

						p2pSent++
					}
					p2pRequestReceived++
				} else {
					bounce.Stamps = append(bounce.Stamps, primitives.NewTimestampNow())
					bounce.Number = cntreq
					bounce.Name = strings.TrimSpace(bounce.Name) + "-" + name
					cntreq++

					SetMsg(msg)
					p2pProxy.Send(msg)

					broadcastReceived++
					broadcastSent++
				}
			}
			if false && ok2 && len(bounceReply.Stamps) < 5 {
				bounceReply.Stamps = append(bounceReply.Stamps, primitives.NewTimestampNow())

				SetMsg(msg)
				p2pProxy.Send(msg)

				p2pReceived++
				p2pSent++
			}

		} else {
			oldcnt++
			fmt.Println("OLD:", msg.String())
		}

	}
}

func main() {
	InitNetwork()

	time.Sleep(1 * time.Second)
	fmt.Println("Starting...")

	go listen()

	for {
		bounce := new(messages.Bounce)
		bounce.Number = cntreq
		cntreq++
		bounce.Name = name
		bounce.AddData(size)
		bounce.Timestamp = primitives.NewTimestampNow()
		bounce.Stamps = append(bounce.Stamps, primitives.NewTimestampNow())
		if isp2p {
			bounce.SetPeer2Peer(true)
			p2pRequestSent++
		} else {
			broadcastSent++
		}
		p2pProxy.Send(bounce)
		SetMsg(bounce)

		if isp2p {
			fmt.Printf("netTest(%s):  ::p2p:: request sent: %d request recieved %d sent: %d received: %d\n",
				name,
				p2pRequestSent, p2pRequestReceived,
				p2pSent, p2pReceived)

		} else {
			fmt.Printf("netTest(%s):  ::: broadcast sent: %d broadcast received: %d\n",
				name,

				broadcastSent, broadcastReceived)
		}
		time.Sleep(8 * time.Second)
	}
}

// if isp2p {
// 	fmt.Printf("netTest(%s): Reads: %d errs %d Writes %d errs %d  ::p2p:: request sent: %d request recieved %d sent: %d received: %d\n",
// 		name,
// 		p2p.Reads, p2p.ReadsErr,
// 		p2p.Writes, p2p.WritesErr,
// 		p2pRequestSent, p2pRequestReceived,
// 		p2pSent, p2pReceived)

// } else {
// 	fmt.Printf("netTest(%s): Reads: %d errs %d Writes %d errs %d  ::: broadcast sent: %d broadcast received: %d\n",
// 		name,
// 		p2p.Reads, p2p.ReadsErr,
// 		p2p.Writes, p2p.WritesErr,
// 		broadcastSent, broadcastReceived)
// }
