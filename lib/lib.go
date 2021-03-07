package lib

import (
	"encoding/binary"
	"fmt"

	"github.com/hashicorp/memberlist"
)

// DeviceMeta holds the current meta data for a device
type DeviceMeta struct{ Voltage uint64 }

// NodeMeta is used to retrieve meta-data about the current node
// when broadcasting an alive message. It's length is limited to
// the given byte size. This metadata is available in the Node structure.
func (d *DeviceMeta) NodeMeta(limit int) []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, d.Voltage)
	return bs
}

// NotifyMsg is called when a user-data message is received.
// Care should be taken that this method does not block, since doing
// so would block the entire UDP packet receive loop. Additionally, the byte
// slice may be modified after the call returns, so it should be copied if needed
func (d *DeviceMeta) NotifyMsg(_ []byte) {

}

// GetBroadcasts is called when user data messages can be broadcast.
// It can return a list of buffers to send. Each buffer should assume an
// overhead as provided with a limit on the total byte size allowed.
// The total byte size of the resulting data to send must not exceed
// the limit. Care should be taken that this method does not block,
// since doing so would block the entire UDP packet receive loop.
func (d *DeviceMeta) GetBroadcasts(overhead int, limit int) [][]byte {
	return make([][]byte, 0)
}

// LocalState is used for a TCP Push/Pull. This is sent to
// the remote side in addition to the membership information. Any
// data can be sent here. See MergeRemoteState as well. The `join`
// boolean indicates this is for a join instead of a push/pull.
func (d *DeviceMeta) LocalState(join bool) []byte {
	return make([]byte, 0)
}

// MergeRemoteState is invoked after a TCP Push/Pull. This is the
// state received from the remote side and is the result of the
// remote side's LocalState call. The 'join'
// boolean indicates this is for a join instead of a push/pull.
func (d *DeviceMeta) MergeRemoteState(buf []byte, join bool) {
}

// PrintMembers prints all member of the list, including current Meta (voltage)
func PrintMembers(list *memberlist.Memberlist) {
	fmt.Println("Most Recent Memberlist:")
	for _, member := range list.Members() {
		meta := binary.LittleEndian.Uint64(member.Meta)
		fmt.Printf("Member: %s %s meta:%d\n", member.Name, member.Addr, meta)
	}
}
