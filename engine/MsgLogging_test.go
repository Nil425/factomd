package engine_test

import (
	"testing"

	. "github.com/FactomProject/factomd/engine"
	"github.com/FactomProject/factomd/state"
	"github.com/FactomProject/factomd/common/messages"
	"github.com/FactomProject/factomd/common/primitives"
)



func TestMessageLoging(t *testing.T) {
	msgLog := new(MsgLog)
	msgLog.Init(true, 0)

	fnode := new(FactomNode)
	fnode.State = new(state.State)
	s := fnode.State

	msg := new(messages.Bounce)
	msg.Name = "bob"
	msg.Timestamp = primitives.NewTimestampNow()
	msg.Data = []byte("here is some data")
	msg.Stamps = append(msg.Stamps, primitives.NewTimestampNow())

	msgLog.PrtMsgs(s)

	msgLog.Add2(fnode, true, "peer","where",true,msg)
	msgLog.Add2(fnode, false, "peer","where",true,msg)

	msgLog.Startp = primitives.NewTimestampFromMilliseconds(0)
	msgLog.Add2(fnode, false, "peer","where",true,msg)

	if len(msgLog.MsgList) != 3 {
		t.Error("Should have three entries")
	}
	msgLog.PrtMsgs(s)
	
	msgLog.Last = primitives.NewTimestampFromMilliseconds(0)
	msgLog.Add2(fnode, false, "peer","where",true,msg)


	if len(msgLog.MsgList) != 1 {
		t.Error("Should have one message")
	}
	msgLog.PrtMsgs(s)
}
