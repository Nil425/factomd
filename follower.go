// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/log"
)

var _ = fmt.Print
var _ = log.Printf

func Follower(state interfaces.IState) {

	for {
		msg := <-state.FollowerInMsgQueue()
		if state.PrintType(msg.Type()) {
			fmt.Printf("%20s %s\n", "Follower:", msg.String())
		}
		msg.FollowerExecute(state)
		state.UpdateState()	
	}

}
