package main

import "github.com/FactomProject/electiontesting/ET2/dive"

func main() {
	dive.Main()
}

//package main
//
//import (
//	"bytes"
//	"encoding/gob"
//	"fmt"
//
//	"github.com/FactomProject/electiontesting/controller"
//	"github.com/FactomProject/electiontesting/election"
//	"github.com/FactomProject/electiontesting/imessage"
//
//	"crypto/sha256"
//
//	"github.com/FactomProject/electiontesting/messages"
//	"github.com/FactomProject/electiontesting/primitives"
//	"github.com/dustin/go-humanize"
//)
//
//var mirrorMap map[[32]byte][]byte
//
//var solutions = 0
//var breadth = 0
//var loops = 0
//var mirrors = 0
//var depths []int
//var solutionsAt []int
//var mirrorsAt []int
//var deadMessagesAt []int
//var failuresAt []int
//var hitlimit int
//var maxdepth int
//var failure int
//var errConclusions int
//
//var globalRunNumber = 0
//
//var leadersMap = make(map[primitives.Identity]int)
//var audsMap = make(map[primitives.Identity]int)
//
//var extraPrints = true
//var extraPrints1 = true
//var extraPrints2 = true
//var insanePrints = false
//
////================ main =================
//func main() {
//	recurse(2, 5, 100)
//}
//
//// newElections will return an array of elections (1 per leader) and an array
//// of volunteers messages to kick things off.
////		Params:
////			feds int   Number of Federated Nodes
////			auds int   Number of Volunteers
////			noDisplay  Passing a true here will reduce memory consumption, as it is a debugging tool
////
////		Returns:
////			controller *Controller  This can used for debugging (Printing votes)
////			elections []*election   Nodes you can execute on (returns msg, statchange)
////			volmsgs   []*VoluntMsg	Volunteer msgs you can start things with
//func newElections(feds, auds int, noDisplay bool) (*controller.Controller, []*election.Election, []*controller.DirectedMsg) {
//	con := controller.NewController(feds, auds)
//
//	if noDisplay {
//		for _, e := range con.Elections {
//			e.Display = nil
//		}
//		con.GlobalDisplay = nil
//	}
//	var msgs []*controller.DirectedMsg
//	fmt.Println("Starting")
//	for _, v := range con.Volunteers {
//		for i, _ := range con.Elections {
//			my := new(controller.DirectedMsg)
//			my.LeaderIdx = i
//			my.Msg = v
//			msgs = append(msgs, my)
//			fmt.Println(my.Msg.String(), my.LeaderIdx)
//		}
//	}
//
//	global := con.Elections[0].Display.Global
//	for i, ldr := range con.Elections {
//		con.Elections[i] = CloneElection(ldr)
//		con.Elections[i].Display.Global = global
//	}
//
//	for _, l := range con.Elections {
//		leadersMap[l.Self] = con.Elections[0].FedIDtoIndex(l.Self)
//	}
//
//	for _, a := range con.AuthSet.GetAuds() {
//		audsMap[a] = con.Elections[0].GetVolunteerPriority(a)
//	}
//
//	return con, con.Elections, msgs
//}
//
//type mymsg struct {
//	leaderIdx int
//	msg       imessage.IMessage
//}
//
//var cnt = 0
//
//// dive
//// Pass a list of messages to process, to a set of leaders, at a current depth, with a particular limit.
//// Provided a msgPath, and updated for recording purposes.
//// Returns
//// limitHit -- path hit the limit, and recursed.  Can happen on loops
//// leaf -- All messages were processed, and no message resulted in a change of state.
//// seeSuccess -- Some path below this dive produced a solution
//// Note that we actually dive 100 levels beyond our limit, and declare seeSuccess past our limit as proof we are
//// in a loop.
//// Hitting the limit and seeSuccess is proof of a loop that none the less can resolve.
//func dive(msgs []*controller.DirectedMsg, leaders []*election.Election, depth int, limit int, msgPath []*controller.DirectedMsg) (limitHit bool, leaf bool, seeSuccess bool) {
//	depths = incCounter(depths, depth)
//	depth++
//
//	if globalRunNumber < 1000 || globalRunNumber%50000 == 0 {
//		extraPrints = true
//		extraPrints1 = true
//		extraPrints2 = true
//	}
//
//	printState := func() {
//		fmt.Printf("%s%s%4d%s%4d %s %12s %s%12s %s%3d %s%12s %12s  %s %12s %s %12s %s %12s %s %12s %s %12s", "=============== ",
//			" Depth=", depth, "/", maxdepth,
//			"| Multiple Conclusions", humanize.Comma(int64(errConclusions)),
//			"| Failures=", humanize.Comma(int64(failure)),
//			"| MsgQ=", len(msgs),
//			"| Mirrors=", humanize.Comma(int64(mirrors)), humanize.Comma(int64(len(mirrorMap))),
//			"| Hit the Limits=", humanize.Comma(int64(hitlimit)),
//			"| Breadth=", humanize.Comma(int64(breadth)),
//			"| solutions so far =", humanize.Comma(int64(solutions)),
//			"| global count= ", humanize.Comma(int64(globalRunNumber)),
//			"| loops detected=", humanize.Comma(int64(loops)))
//
//		prt := func(counter []int, msg string) {
//			fmt.Printf("\n=%20s", msg)
//			if len(counter) == 0 {
//				fmt.Println("\n=     None Found\n=")
//			}
//			for i, v := range counter {
//				if i%16 == 0 {
//					fmt.Println("")
//					fmt.Print("=")
//				}
//				str := fmt.Sprintf("%s[%3d]", humanize.Comma(int64(v)), i)
//				fmt.Printf("%12s ", str)
//			}
//		}
//		prt(deadMessagesAt, "Dead Messages")
//		prt(mirrorsAt, "Mirrors")
//		prt(solutionsAt, "Solutions")
//		prt(failuresAt, "Failures")
//		prt(depths, "Depths")
//		fmt.Println()
//
//		// Lots of printing... Not necessary....
//		fmt.Println(leaders[0].Display.Global.String())
//
//		for _, ldr := range leaders {
//			fmt.Println(ldr.Display.String())
//		}
//
//		if insanePrints {
//			// Example of a run that has a werid msg state
//			if globalRunNumber > -1 {
//				fmt.Println("Leader 0")
//				fmt.Println(leaders[0].PrintMessages())
//				fmt.Println("Leader 1")
//				fmt.Println(leaders[1].PrintMessages())
//				fmt.Println("Leader 2")
//				fmt.Println(leaders[2].PrintMessages())
//			}
//		}
//		fmt.Printf("%d %d setcon\n", len(leadersMap), len(audsMap))
//		for i, v := range msgPath {
//			fmt.Println(formatForInterpreter(v), "#", i, v.LeaderIdx, "<==", leaders[0].Display.FormatMessage(v.Msg))
//		}
//		fmt.Println("<b> # Pending:")
//		for i, v := range msgs {
//			fmt.Println(formatForInterpreter(v), "#", i, v.LeaderIdx, "<==", leaders[0].Display.FormatMessage(v.Msg))
//		}
//
//	}
//
//	if depth > limit {
//		if extraPrints {
//			fmt.Println(">>>>>>>>>>>>>>>>>> Hit Limit <<<<<<<<<<<<<<<<<")
//			printState()
//			extraPrints = false
//		}
//		breadth++
//		hitlimit++
//		return true, false, false
//	}
//
//	if depth > maxdepth {
//		maxdepth = depth
//	}
//
//	if depth < 4 {
//		fmt.Println("==================== Depth %d ====================")
//		printState()
//	}
//
//	//done := 0
//	//for _, ldr := range leaders {
//	//	if ldr.Committed {
//	//		done++
//	//	}
//	//}
//	if complete, err := nodesCompleted(leaders); complete { // done == len(leaders)/2+1 {
//		solutionsAt = incCounter(solutionsAt, depth)
//		if extraPrints {
//			fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>> Solution Found @ ", depth)
//		}
//		breadth++
//		solutions++
//		if extraPrints2 {
//			fmt.Println("!!!!!!!!!!!!!!!!!!  Success!")
//			printState()
//			extraPrints2 = false
//		}
//		return false, true, true
//
//	} else if err != nil {
//		// Bad! This means the algorithm is broken
//		errConclusions++
//		fmt.Println("^&**(^%$& Broken! @ depth:", depth)
//		fmt.Println(err.Error())
//		printState()
//	}
//
//	// Look for mirrorMap, but only after we have been going a bit.
//	if depth > 0 {
//		var hashes [][32]byte
//		var strings []string
//		for _, ldr := range leaders {
//			bits := ldr.NormalizedString()
//			if bits != nil {
//				strings = append(strings, string(bits))
//				h := Sha(bits)
//				hashes = append(hashes, h)
//			} else {
//				panic("shouldn't happen")
//			}
//		}
//		for i := 0; i < len(hashes)-1; i++ {
//			for j := 0; j < len(hashes)-1-i; j++ {
//				if bytes.Compare(hashes[j][:], hashes[j+1][:]) > 0 {
//					hashes[j], hashes[j+1] = hashes[j+1], hashes[j+1]
//					strings[j], strings[j+1] = strings[j+1], strings[j]
//				}
//			}
//		}
//		var all []byte
//		var alls string
//		for i, h := range hashes {
//			all = append(all, h[:]...)
//			alls += strings[i]
//		}
//		mh := Sha(all)
//		if mirrorMap[mh] != nil {
//			mirrors++
//			breadth++
//			mirrorsAt = incCounter(mirrorsAt, depth)
//			return false, false, true
//		}
//		mirrorMap[mh] = mh[:]
//	}
//
//	leaf = true
//	d := 3
//	for range msgs {
//
//		d += 3
//		d = d % len(msgs)
//
//		v := msgs[d]
//
//		var msgs2 []*controller.DirectedMsg
//		msgs2 = append(msgs2, msgs[0:d]...)
//		msgs2 = append(msgs2, msgs[d+1:]...)
//		ml2 := len(msgs2)
//
//		globalRunNumber++
//
//		cl := CloneElection(leaders[v.LeaderIdx])
//
//		//if !spewSame(cl, leaders[v.leaderIdx]) {
//		//	fmt.Println("Clone Failed")
//		//	debugClone(cl, leaders[v.leaderIdx])
//		//	os.Exit(0)
//		//}
//
//		msg, changed := leaders[v.LeaderIdx].Execute(v.Msg, depth)
//
//		msgPath2 := append(msgPath, v)
//
//		if changed {
//			leaf = false
//			if msg != nil {
//				for i, _ := range leaders {
//					if i != v.LeaderIdx {
//						my := new(controller.DirectedMsg)
//						my.LeaderIdx = i
//						my.Msg = msg
//						msgs2 = append(msgs2, my)
//					}
//				}
//			}
//			gl := leaders[v.LeaderIdx].Display.Global
//			for _, ldr := range leaders {
//				ldr.Display.Global = gl
//			}
//			// Recursive Dive
//			lim, _, ss := dive(msgs2, leaders, depth, limit, msgPath2)
//			_ = lim || ss
//			seeSuccess = seeSuccess || ss
//			limitHit = limitHit || lim
//			for _, ldr := range leaders {
//				ldr.Display.Global = cl.Display.Global
//			}
//			msgs2 = msgs2[:ml2]
//		} else {
//			deadMessagesAt = incCounter(deadMessagesAt, depth)
//		}
//		leaders[v.LeaderIdx] = cl
//	}
//	if limitHit {
//		leaf = false
//	}
//	if limitHit {
//		if depth == 9 {
//
//			if seeSuccess {
//				loops++
//			} else {
//				failure++
//				if extraPrints1 {
//					extraPrints1 = false
//					fmt.Println("/////////////// Loops Fail //////////////////////")
//					fmt.Printf("%d %d setcon\n", len(leadersMap), len(audsMap))
//					printState()
//				}
//			}
//			limitHit = false
//		}
//	} else {
//		if leaf {
//			incCounter(failuresAt, depth)
//			failure++
//			leaf = false
//
//			if extraPrints1 {
//				extraPrints1 = false
//				fmt.Println("/////////////// Fail //////////////////////")
//				fmt.Printf("%d %d setcon\n", len(leadersMap), len(audsMap))
//				printState()
//			}
//
//		}
//	}
//
//	return limitHit, leaf, seeSuccess
//}
//
//func nodesCompleted(nodes []*election.Election) (bool, error) {
//	done := 0
//	prev := -1
//	for _, n := range nodes {
//		if n.Committed {
//			done++
//			if prev != -1 && n.CurrentVote.VolunteerPriority != prev {
//				return false, fmt.Errorf("2 nodes committed on different results. %d and %d", prev, n.CurrentVote.VolunteerPriority)
//			}
//			prev = n.CurrentVote.VolunteerPriority
//		}
//	}
//
//	return done >= (len(nodes)/2)+1, nil
//}
//
//func formatForInterpreter(my *controller.DirectedMsg) string {
//	msg := my.Msg
//	switch msg.(type) {
//	case *messages.LeaderLevelMessage:
//		l := msg.(*messages.LeaderLevelMessage)
//		from := leadersMap[l.Signer]
//
//		return fmt.Sprintf("{ %d } %d { %d } <-l", from, l.Level, my.LeaderIdx)
//	case *messages.VolunteerMessage:
//		a := msg.(*messages.VolunteerMessage)
//		from := audsMap[a.Signer]
//
//		return fmt.Sprintf("%d { %d } <-v", from, my.LeaderIdx)
//	case *messages.VoteMessage:
//		l := msg.(*messages.VoteMessage)
//		from := leadersMap[l.Signer]
//		vol := audsMap[l.Volunteer.Signer]
//
//		return fmt.Sprintf("{ %d } %d { %d } <-o", from, vol, my.LeaderIdx)
//	}
//	return "NA"
//}
//
//func incCounter(counter []int, depth int) []int {
//	for len(counter) <= depth {
//		counter = append(counter, 0)
//	}
//	counter[depth]++
//	return counter
//}
//
//func recurse(auds int, feds int, limit int) {
//
//	_, leaders, msgs := newElections(feds, auds, false)
//	var msgpath []*controller.DirectedMsg
//	dive(msgs, leaders, 0, limit, msgpath)
//}
//
//// reuse encoder/decoder so we don't recompile the struct definition
//var enc *gob.Encoder
//var dec *gob.Decoder
//
//// LoopingDetected will the number of looping leaders
//func LoopingDetected(global *election.Display) int {
//	return global.DetectLoops()
//}
//
//func init() {
//	buff := new(bytes.Buffer)
//	enc = gob.NewEncoder(buff)
//	dec = gob.NewDecoder(buff)
//	mirrorMap = make(map[[32]byte][]byte, 10000)
//}
//
//func CloneElection(src *election.Election) *election.Election {
//	return src.Copy()
//	dst := new(election.Election)
//	err := enc.Encode(src)
//	if err != nil {
//		errConclusions++
//	}
//	err = dec.Decode(dst)
//	if err != nil {
//		errConclusions++
//	}
//	return dst
//}
//
//// Create a Sha256 Hash from a byte array
//func Sha(p []byte) [32]byte {
//	b := sha256.Sum256(p)
//	return b
//}
