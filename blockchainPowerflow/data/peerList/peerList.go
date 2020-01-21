package peerList

import (
	"bytes"
	"container/ring"
	"encoding/json"
	s "github.com/edgexfoundry/powerflow/commonPowerFlow/identity"
	"log"
	"sort"
	"strconv"
	"sync"
)

//Peer
type Peer struct {
	Id     int32
	IpPort string
	Pid    s.PublicIdentity
}

//PeerList contains numId, peerMap, max length, and a mutex
type PeerList struct {
	numId    int32
	ipPort   string
	secureId s.Identity
	//peerMap    map[string]int32
	//peerMapPid map[string]s.PublicIdentity
	pairMap   map[string]int32
	peers     map[string]Peer
	maxLength int32

	mux sync.Mutex
}

//NewPeerList func creates a New PeerList for a id and maxLength
func NewPeerList(numId int32, selfIpPort string, secureId s.Identity, maxLength int32) PeerList {

	return PeerList{
		numId:     numId,
		ipPort:    selfIpPort,
		secureId:  secureId,
		pairMap:   make(map[string]int32),
		peers:     make(map[string]Peer),
		maxLength: maxLength,
	}
}

// ONLY FOR TEST PURPOSES
func TestNewPeerList(id int32 /*sid s.Identity,*/, maxLength int32) PeerList {

	return PeerList{
		numId: id,
		/*secureId:   sid,*/
		pairMap:   make(map[string]int32),
		peers:     make(map[string]Peer),
		maxLength: maxLength,
	}
}

///////////
// Pair - data structure to hold a key/value pair - addr/id.
type Pair struct {
	addr string
	id   int32
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].id < p[j].id }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[string]int32) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		//fmt.Println("in sortMapByValue : k, v :", k, v)
		p[i] = Pair{
			addr: k,
			id:   int32(v),
		}
		//fmt.Println("in sortMapByValue : p[i] :", p[i])
		i++
	}
	//fmt.Println("in sortMapByValue : p :", p)
	sort.Sort(p)
	//fmt.Println("in sortMapByValue : sorted p :", p)
	return p
}

///////////

//Add func adds a peer with addr and id to peerMap
func (peers *PeerList) Add(addr string, id int32) {
	peers.mux.Lock()
	defer peers.mux.Unlock()

	peers.pairMap[addr] = id
}

//Add func adds a peer with addr and id to peerMap
func (peers *PeerList) AddPid(ipPort string, p Peer) {
	peers.mux.Lock()
	defer peers.mux.Unlock()

	peers.peers[ipPort] = p
}

//Delete func deletes a peer with specific addr
func (peers *PeerList) Delete(ipPort string) {
	peers.mux.Lock()
	defer peers.mux.Unlock()

	delete(peers.peers, ipPort)
}

////Delete func deletes a peer with specific addr
//func (peers *PeerList) DeletePid(addr string) {
//	peers.mux.Lock()
//	defer peers.mux.Unlock()
//
//	delete(peers.peerMapPid, addr)
//}

//Rebalance func changes the PeerMap to contain take maxLength(32) closest peers (by Id)
func (peers *PeerList) Rebalance() {

	if int32(len(peers.pairMap)) > peers.maxLength {

		peers.mux.Lock()
		defer peers.mux.Unlock()

		//fmt.Println("in Rebalance")
		//fmt.Println("in Rebalance : original peerMap length : ", len(peers.peerMap))
		peers.pairMap[peers.ipPort] = peers.numId //adding self id to peerMap
		sortedAddrIDList := sortMapByValue(peers.pairMap)
		//fmt.Println("in Rebalance : sortedAddrIDList : ", sortedAddrIDList)
		sortedAddrIDListLength := len(sortedAddrIDList)
		//fmt.Println("in Rebalance : sortedAddrIDListLength : ", sortedAddrIDListLength)

		peers.pairMap = peers.getBalancedPeerMap(sortedAddrIDListLength, sortedAddrIDList)

	}
}

func (peers *PeerList) getBalancedPeerMap(sortedAddrIDListLength int, sortedAddrIDList PairList) map[string]int32 {
	r := ring.New(sortedAddrIDListLength) // new ring
	useRingPtr := r

	//initialize ring with sortedAddrIDList values
	for i := 0; i < sortedAddrIDListLength; i++ {
		r.Value = sortedAddrIDList[i]
		//fmt.Println("in Rebalance : r.Value : ", r.Value)
		if sortedAddrIDList[i].id == peers.numId {
			useRingPtr = r
			//fmt.Println("in Rebalance : useRingPtr : ", useRingPtr)
		}
		r = r.Next()
	}
	newPeerMap := make(map[string]int32)
	r = useRingPtr
	//fmt.Println("in Rebalance : useRingPtr : ", useRingPtr)
	for i := 1; i <= int(peers.maxLength/2); i++ {
		r = r.Prev()
		pair := r.Value.(Pair)
		newPeerMap[pair.addr] = pair.id
	}
	r = useRingPtr
	for i := 1; i <= int(peers.maxLength/2); i++ {
		r = r.Next()
		pair := r.Value.(Pair)
		newPeerMap[pair.addr] = pair.id
	}

	return newPeerMap
}

//Show func returns PeerMap string
func (peers *PeerList) Show() string {
	var buffer bytes.Buffer

	buffer.WriteString("This is PeerMap:\n")
	for k := range peers.pairMap {
		buffer.WriteString("Addr:" + k + " Id:" + strconv.Itoa(int(peers.pairMap[k])) + "\n")
	}
	return buffer.String()
}

//Show func returns PeerMap string
func (peers *PeerList) ShowPids() string {
	var buffer bytes.Buffer

	buffer.WriteString("This is PeerMapPid:\n")
	for k := range peers.peers {
		thisPeerPid := peers.peers[k].Pid
		buffer.WriteString("Addr:" + k + " Pid PubK : " + thisPeerPid.PublicIdentityToJson() + "\n" +
			"ipPort : " + peers.peers[k].IpPort + "\n" +
			"labelt : " + peers.peers[k].Pid.Label + "\n")
	}
	return buffer.String()
}

////Register func assigns a value to numId
//func (peers *PeerList) Register(id int32) {
//	peers.numId = id
//	fmt.Printf("NumId=%v\n", id)
//}

//Copy func returns a copy of the peerMap
func (peers *PeerList) Copy() map[string]int32 {

	peers.mux.Lock()
	defer peers.mux.Unlock()

	copyOfPeerMap := make(map[string]int32)
	for k := range peers.pairMap {
		copyOfPeerMap[k] = peers.pairMap[k]
	}

	return copyOfPeerMap
}

//Copy func returns a copy of the peerMap
func (peers *PeerList) CopyPids() map[string]Peer {

	peers.mux.Lock()
	defer peers.mux.Unlock()

	copyOfPeers := make(map[string]Peer)
	for k := range peers.peers {
		copyOfPeers[k] = peers.peers[k]
	}

	return copyOfPeers
}

//GetNumId func returns numId of Peer
func (peers *PeerList) GetNumId() int32 {
	return peers.numId
}

//PeerMapToJson func returns a json string of PeerMap or an error
func (peers *PeerList) PairMapToJson() (string, error) {
	peers.mux.Lock()

	jsonOfPeerMap, err := json.Marshal(peers.pairMap)

	peers.mux.Unlock()

	return string(jsonOfPeerMap), err
}

////PeerMapToJson func returns a json string of PeerMap or an error
//func PairMapToJson(pairMap map[string]int32) (string, error) {
//
//	jsonOfPeerMap, err := json.Marshal(pairMap)
//
//	return string(jsonOfPeerMap), err
//}

//PeerMapToJson func returns a json string of PeerMap or an error
func (peers *PeerList) PeersToJson() (string, error) {
	peers.mux.Lock()

	jsonOfPeerMapPid, err := json.Marshal(peers.peers)

	peers.mux.Unlock()

	return string(jsonOfPeerMapPid), err
}

//PeerMapSIDToJson func returns a json string of PeerMap or an error
func PeersToJson(peers map[string]Peer) (string, error) {

	jsonOfPeerMapPid, err := json.Marshal(peers)

	return string(jsonOfPeerMapPid), err
}

//InjectPeerMapJson func injects the new PeerMap into existing PeerMap, except for the entry corresponding to self
func (peers *PeerList) InjectPairMapJson(peerMapJsonStr string, selfAddr string) {

	var newPeerMap map[string]int32
	err := json.Unmarshal([]byte(peerMapJsonStr), &newPeerMap)
	if err == nil {
		peers.mux.Lock()

		for k := range newPeerMap {
			if /*_, ok := peers.peerMap[k]; !ok &&*/ k != selfAddr {
				peers.pairMap[k] = newPeerMap[k]
			}
		}

		peers.mux.Unlock()
	}
}

//InjectPeerMapJson func injects the new PeerMap into existing PeerMap, except for the entry corresponding to self
func (peers *PeerList) InjectPeersJson(peerMapPidJsonStr string, selfAddr string) {

	var recvPeerMapPid map[string]Peer
	err := json.Unmarshal([]byte(peerMapPidJsonStr), &recvPeerMapPid)
	if err == nil {
		//peerMapPidCopy := peers.CopyPids()
		for addr, pid := range recvPeerMapPid {
			if _, ok := peers.peers[addr]; !ok {
				peers.mux.Lock()
				peers.peers[addr] = pid
				peers.mux.Unlock()
			}
		}
	} else {
		log.Println("Error in Inject PeersJson : err : ", err)
	}
}

//
//func TestPeerListRebalance() {
//	peers := NewPeerList(5, s.Identity{}, 4)
//	peers.Add("1111", 1)
//	peers.Add("4444", 4)
//	peers.Add("-1-1", -1)
//	peers.Add("0000", 0)
//	peers.Add("2121", 21)
//	peers.Rebalance()
//	expected := NewPeerList(5, s.Identity{}, 4)
//	expected.Add("1111", 1)
//	expected.Add("4444", 4)
//	expected.Add("2121", 21)
//	expected.Add("-1-1", -1)
//	fmt.Println(reflect.DeepEqual(peers, expected))
//
//	peers = NewPeerList(5, s.Identity{}, 2)
//	peers.Add("1111", 1)
//	peers.Add("4444", 4)
//	peers.Add("-1-1", -1)
//	peers.Add("0000", 0)
//	peers.Add("2121", 21)
//	peers.Rebalance()
//	expected = NewPeerList(5, s.Identity{}, 2)
//	expected.Add("4444", 4)
//	expected.Add("2121", 21)
//	fmt.Println(reflect.DeepEqual(peers, expected))
//
//	peers = NewPeerList(5, s.Identity{}, 4)
//	peers.Add("1111", 1)
//	peers.Add("7777", 7)
//	peers.Add("9999", 9)
//	peers.Add("11111111", 11)
//	peers.Add("2020", 20)
//	peers.Rebalance()
//	expected = NewPeerList(5, s.Identity{}, 4)
//	expected.Add("1111", 1)
//	expected.Add("7777", 7)
//	expected.Add("9999", 9)
//	expected.Add("2020", 20)
//	fmt.Println(reflect.DeepEqual(peers, expected))
//}
