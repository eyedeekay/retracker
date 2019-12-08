package tracker

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/vvampirius/retracker/bittorrent/common"
	"log"
	"strconv"
	"time"
)

var i2pB64enc *base64.Encoding = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-~")

type Request struct {
	timestamp  time.Time
	remoteAddr common.Address
	InfoHash   common.InfoHash
	PeerID     common.PeerID
	Port       int
	Uploaded   uint64
	Downloaded uint64
	Left       uint64
	IP         common.Address
	NumWant    uint64
	Event      string
}

func (self *Request) Peer() common.Peer {
	peer := common.Peer{
		PeerID: self.PeerID,
		IP:     self.IP,
		Port:   self.Port,
	}
	if !peer.IP.Valid() {
		peer.IP = self.remoteAddr
	}
	return peer
}

// get base64 representation of i2p dest sha256 hash(the 44-character one)
func (p *Request) I2PHash() common.Address {
	hash := sha256.New()
	hash.Write([]byte(p.IP))
	digest := hash.Sum(nil)
	buf := make([]byte, 44)
	i2pB64enc.Encode(buf, digest)
	return common.Address(buf)
}

func (self *Request) Compact() common.Address {
	if len(self.IP) > 350 {
		return self.I2PHash()
	}
	return self.IP
}

func (self *Request) CompactPeer() common.Peer {
	peer := common.Peer{
		IP:   self.Compact(),
		Port: self.Port,
	}
	if !peer.IP.Valid() {
		peer.IP = self.remoteAddr
	}
	return peer
}

func (self *Request) String() string {
	return fmt.Sprintf("%s info_hash:%x peer_id:%x port:%d ip:%s numwant:%d event:%s", self.remoteAddr, self.InfoHash, self.PeerID, self.Port, self.IP, self.NumWant, self.Event)
}

func (self *Request) TimeStampDelta() float64 {
	return time.Now().Sub(self.timestamp).Minutes()
}

func MakeRequest(remoteAddr, infoHash, peerID, port, uploaded, downloaded, left, ip, numwant,
	event string, logger *log.Logger) (*Request, error) {
	request := Request{timestamp: time.Now(), remoteAddr: common.Address(remoteAddr)}

	if v := common.InfoHash(infoHash); v.Valid() {
		request.InfoHash = v
	} else {
		return nil, errors.New(`info_hash is not valid`)
	}

	if v := common.PeerID(peerID); v.Valid() {
		request.PeerID = v
	} else {
		return nil, errors.New(`peer_id is not valid`)
	}

	if d, err := strconv.Atoi(port); err == nil {
		request.Port = d
	} else {
		return nil, errors.New(fmt.Sprintf("Can't parse 'port' value: '%s'", err.Error()))
	}

	if d, err := strconv.ParseUint(uploaded, 10, 64); err == nil {
		request.Uploaded = d
	} else {
		return nil, errors.New(fmt.Sprintf("Can't parse 'uploaded' value: '%s'", err.Error()))
	}

	if d, err := strconv.ParseUint(downloaded, 10, 64); err == nil {
		request.Downloaded = d
	} else {
		return nil, errors.New(fmt.Sprintf("Can't parse 'downloaded' value: '%s'", err.Error()))
	}

	if d, err := strconv.ParseUint(left, 10, 64); err == nil {
		request.Left = d
	} else {
		return nil, errors.New(fmt.Sprintf("Can't parse 'left' value: '%s'", err.Error()))
	}

	request.IP = common.Address(ip)

	if d, err := strconv.ParseUint(numwant, 10, 64); err == nil {
		request.NumWant = d
	}

	if event := event; event == `` || event == `started` || event == `stopped` || event == `completed` {
		request.Event = event
	} else {
		if logger != nil {
			logger.Printf("WARNING! Got '%s' event in announce.\n", event)
		}
	}

	return &request, nil
}
