package tests

import (
	"fmt"
	"github.com/edgexfoundry/powerflow/blockchainPowerflow/data/peerList"
	s "github.com/edgexfoundry/powerflow/commonPowerFlow/identity"
	"reflect"
	"testing"
)

func TestPeerListRebalance(t *testing.T) {
	peers := peerList.NewPeerList(5, s.Identity{}, 4)
	peers.Add("1111", 1)
	peers.Add("4444", 4)
	peers.Add("-1-1", -1)
	peers.Add("0000", 0)
	peers.Add("2121", 21)
	peers.Rebalance()
	expected := peerList.NewPeerList(5, s.Identity{}, 4)
	expected.Add("1111", 1)
	expected.Add("4444", 4)
	expected.Add("2121", 21)
	expected.Add("-1-1", -1)
	fmt.Println(reflect.DeepEqual(peers, expected))

	peers = peerList.NewPeerList(5, s.Identity{}, 2)
	peers.Add("1111", 1)
	peers.Add("4444", 4)
	peers.Add("-1-1", -1)
	peers.Add("0000", 0)
	peers.Add("2121", 21)
	peers.Rebalance()
	expected = peerList.NewPeerList(5, s.Identity{}, 2)
	expected.Add("4444", 4)
	expected.Add("2121", 21)
	fmt.Println(reflect.DeepEqual(peers, expected))

	peers = peerList.NewPeerList(5, s.Identity{}, 4)
	peers.Add("1111", 1)
	peers.Add("7777", 7)
	peers.Add("9999", 9)
	peers.Add("11111111", 11)
	peers.Add("2020", 20)
	peers.Rebalance()
	expected = peerList.NewPeerList(5, s.Identity{}, 4)
	expected.Add("1111", 1)
	expected.Add("7777", 7)
	expected.Add("9999", 9)
	expected.Add("2020", 20)
	fmt.Println(reflect.DeepEqual(peers, expected))
}
