package eventmessages

import (
	"encoding/binary"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/messages"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/proto"
	"time"
)

type SourceEvent struct {
	eventSource    EventSource
	messagePayload interfaces.IMsg
}

func (srcEvent *SourceEvent) GetEventSource() EventSource {
	return srcEvent.eventSource
}

func (srcEvent *SourceEvent) GetMessagePayload() interfaces.IMsg {
	return srcEvent.messagePayload
}

func SourceEventFromMessage(eventSource EventSource, msg interfaces.IMsg) *SourceEvent {
	return &SourceEvent{
		eventSource:    eventSource,
		messagePayload: msg}
}

func MapToFactomEvent(sourceEvent SourceEvent) *FactomEvent {
	if sourceEvent.messagePayload != nil {
		return msgToFactomEvent(sourceEvent.GetEventSource(), sourceEvent.messagePayload)
	}
	panic("No payload found in source event.")
}

func msgToFactomEvent(eventSource EventSource, msg interfaces.IMsg) *FactomEvent {
	event := &FactomEvent{}
	event.EventSource = eventSource
	switch msg.(type) {
	case *messages.DBStateMsg:
		event.Value = mapDBState(msg.(*messages.DBStateMsg))
	case *messages.CommitChainMsg:
		event.Value = mapCommitChain(msg)
	case *messages.CommitEntryMsg:
		event.Value = mapCommitEvent(msg)
	default:
		return nil
	}
	return event
}

func mapDBState(dbStateMessage *messages.DBStateMsg) *FactomEvent_AnchorEvent {
	event := &FactomEvent_AnchorEvent{AnchorEvent: &AnchoredEvent{
		DirectoryBlock: mapDirBlock(dbStateMessage.DirectoryBlock),
	}}
	return event
}

func mapDirBlock(block interfaces.IDirectoryBlock) *DirectoryBlock {
	result := &DirectoryBlock{}
	result.Header = mapDirHeader(block.GetHeader())
	result.Entries = mapDirEntries(block.GetDBEntries())
	return result
}

func mapDirHeader(header interfaces.IDirectoryBlockHeader) *DirectoryBlockHeader {

	time := header.GetTimestamp().GetTime()
	result := &DirectoryBlockHeader{
		BodyMerkleRoot: &Hash{
			HashValue: header.GetBodyMR().Bytes(),
		},
		PreviousKeyMerkleRoot: &Hash{
			HashValue: header.GetPrevKeyMR().Bytes(),
		},
		PreviousFullHash: &Hash{
			HashValue: header.GetPrevFullHash().Bytes(),
		},
		Timestamp:  &types.Timestamp{Seconds: int64(time.Second()), Nanos: int32(time.Nanosecond())},
		DBHeight:   header.GetDBHeight(),
		BlockCount: header.GetBlockCount(),
	}

	return result
}

func mapDirEntries(entries []interfaces.IDBEntry) []*Entry {
	result := make([]*Entry, len(entries))
	for i, entry := range entries {
		result[i] = mapDirEntry(entry)
	}
	return result
}

func mapDirEntry(entry interfaces.IDBEntry) *Entry {
	result := &Entry{
		ChainID: &Hash{
			HashValue: entry.GetChainID().Bytes(),
		},
		KeyMerkleRoot: &Hash{
			HashValue: entry.GetKeyMR().Bytes(),
		},
	}
	return result
}

func mapCommitChain(msg interfaces.IMsg) *FactomEvent_CommitChain {
	commitChain := msg.(*messages.CommitChainMsg).CommitChain
	ecPubKey := commitChain.ECPubKey.Fixed()
	sig := commitChain.Sig

	result := &FactomEvent_CommitChain{
		CommitChain: &CommitChain{
			ChainIDHash: &Hash{
				HashValue: commitChain.ChainIDHash.Bytes()},
			EntryHash: &Hash{
				HashValue: commitChain.EntryHash.Bytes(),
			},
			Timestamp: convertToTimestamp(commitChain.MilliTime),
			Credits:   uint32(commitChain.Credits),
			ECPubKey:  ecPubKey[:],
			Sig:       sig[:],
		}}
	return result
}

func mapCommitEvent(msg interfaces.IMsg) *FactomEvent_CommitEntry {
	commitEntry := msg.(*messages.CommitEntryMsg).CommitEntry
	ecPubKey := commitEntry.ECPubKey.Fixed()
	sig := commitEntry.Sig

	result := &FactomEvent_CommitEntry{
		CommitEntry: &CommitEntry{
			EntryHash: &Hash{
				HashValue: commitEntry.EntryHash.Bytes(),
			},
			Timestamp: convertToTimestamp(commitEntry.MilliTime),
			Credits:   uint32(commitEntry.Credits),
			ECPubKey:  ecPubKey[:],
			Sig:       sig[:],
		}}
	return result
}

func convertToTimestamp(milliTime *primitives.ByteSlice6) *types.Timestamp {
	// TODO Is there an easier way to do this?
	slice8 := make([]byte, 8)
	copy(slice8[2:], milliTime[:])
	millis := int64(binary.BigEndian.Uint64(slice8))
	time := time.Unix(0, millis*1000000)
	return &types.Timestamp{Seconds: int64(time.Second()), Nanos: int32(time.Nanosecond())}
}

type Event interface {
	Reset()
	String() string
	ProtoMessage()
	XXX_Unmarshal(b []byte) error
	XXX_Marshal(b []byte, deterministic bool) ([]byte, error)
	XXX_Merge(src proto.Message)
	XXX_Size() int
	XXX_DiscardUnknown()
}
