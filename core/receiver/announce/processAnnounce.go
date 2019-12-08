package announce

import (
	"github.com/vvampirius/retracker/bittorrent/tracker"
)

func (self *Announce) ProcessAnnounce(remoteAddr, infoHash, peerID, port, uploaded, downloaded, left, ip, numwant,
	event, compact string) *tracker.Response {
	if request, err := tracker.MakeRequest(remoteAddr, infoHash, peerID, port, uploaded, downloaded, left, ip, numwant,
		event, self.Logger); err == nil {
		if self.Logger != nil {
			self.Logger.Println(request.String())
		}

		response := tracker.Response{
			Interval: 30,
		}

		if request.Event != `stopped` {
			self.Storage.Update(*request)
			response.Peers = self.Storage.GetPeers(request.InfoHash, compact)
		} else {
			self.Storage.Delete(*request)
			//TODO: make another response ?
		}

		return &response
	} else {
		if self.Logger != nil {
			self.Logger.Println(err.Error())
		}
	}

	return nil
}
