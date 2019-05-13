package discovery

import (
	"reflect"
	"strings"
	"testing"

	"github.com/harmony-one/bls/ffi/go/bls"
	"github.com/harmony-one/harmony/api/proto"
	"github.com/harmony-one/harmony/api/proto/node"
	"github.com/harmony-one/harmony/crypto/pki"
	"github.com/harmony-one/harmony/p2p"
)

var (
	pubKey1 = pki.GetBLSPrivateKeyFromInt(333).GetPublicKey()
	p1      = p2p.Peer{
		IP:              "127.0.0.1",
		Port:            "9999",
		ConsensusPubKey: pubKey1,
	}
	e1 = `ssW`
	e3 = "sssss"

	pubKey2 = pki.GetBLSPrivateKeyFromInt(999).GetPublicKey()

	p2 = []p2p.Peer{
		{
			IP:              "127.0.0.1",
			Port:            "8888",
			ConsensusPubKey: pubKey1,
		},
		{
			IP:              "127.0.0.1",
			Port:            "9999",
			ConsensusPubKey: pubKey2,
		},
	}
	e2 = "pong:1=>length:2"

	leaderPubKey = pki.GetBLSPrivateKeyFromInt(888).GetPublicKey()

	pubKeys = []*bls.PublicKey{pubKey1, pubKey2}

	buf1 []byte
	buf2 []byte
)

func TestString(test *testing.T) {
	ping1 := NewPingMessage(p1, false)

	if strings.Compare(ping1.String(), e1) != 0 {
		test.Errorf("expect: %v, got: %v", e1, ping1.String())
	}

	ping1.Node.Role = uint32(node.ClientRole)

	if strings.Compare(ping1.String(), e3) != 0 {
		test.Errorf("expect: %v, got: %v", e3, ping1.String())
	}

	pong1 := NewPongMessage(p2, pubKeys, leaderPubKey, 0)

	if !strings.HasPrefix(pong1.String(), e2) {
		test.Errorf("expect: %v, got: %v", e2, pong1.String())
	}
}

func TestSerialize(test *testing.T) {
	ping1 := NewPingMessage(p1, true)
	buf1 = ping1.ConstructPingMessage()
	msg1, err := proto.GetMessagePayload(buf1)
	if err != nil {
		test.Error("GetMessagePayload Failed!")
	}
	ping, err := GetPingMessage(msg1)
	if err != nil {
		test.Error("Ping failed!")
	}
	if !reflect.DeepEqual(ping, ping1) {
		test.Error("Serialize/Deserialze Ping Message Failed")
	}

	pong1 := NewPongMessage(p2, pubKeys, leaderPubKey, 0)
	buf2 = pong1.ConstructPongMessage()

	msg2, err := proto.GetMessagePayload(buf2)
	if err != nil {
		test.Error("GetMessagePayload Failed!")
	}
	pong, err := GetPongMessage(msg2)
	if err != nil {
		test.Error("Pong failed!")
	}

	if !reflect.DeepEqual(pong, pong1) {
		test.Error("Serialize/Deserialze Pong Message Failed")
	}
}
