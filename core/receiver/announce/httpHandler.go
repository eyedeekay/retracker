package announce

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func (self *Announce) HttpHandler(w http.ResponseWriter, r *http.Request) {
	xrealip := r.Header.Get(`X-Real-IP`)
	xi2pdest := r.Header.Get(`X-I2p-Dest-Base64`)
	compact := r.URL.Query().Get(`compact`)
	// if x-i2p-dest-base64 header is present enforce it's usage for announces.
	// if it isn't, assume it's on purpose because by default it will be
	// present on I2P clients.
	if xi2pdest != "" {
		if !strings.HasPrefix(r.URL.Query().Get(`ip`), xi2pdest) {
			return
		}
	}
	if self.Logger != nil {
		self.Logger.Printf("%s %s %s %s\n", r.RemoteAddr, xrealip, r.RequestURI, r.UserAgent())
	}
	rr := self.ProcessAnnounce(
		self.getRemoteAddr(r, xrealip, xi2pdest),
		r.URL.Query().Get(`info_hash`),
		r.URL.Query().Get(`peer_id`),
		r.URL.Query().Get(`port`),
		r.URL.Query().Get(`uploaded`),
		r.URL.Query().Get(`downloaded`),
		r.URL.Query().Get(`left`),
		r.URL.Query().Get(`ip`),
		r.URL.Query().Get(`numwant`),
		r.URL.Query().Get(`event`),
		compact,
	)
	if d, err := rr.Bencode(); err == nil {
		fmt.Fprint(w, d)
		if self.Logger != nil && self.Config.Debug {
			self.Logger.Printf("Bencode: %s\n", d)
		}
	} else {
		self.Logger.Println(err.Error())
	}
}

func (self *Announce) getRemoteAddr(r *http.Request, xrealip, xi2pdest string) string {
	// if we're given an I2P base64, always return it. Append .i2p to the dest
	// if it's not present.
	if xi2pdest != `` {
		if strings.HasSuffix(xi2pdest, ".i2p") {
			return xi2pdest
		}
		return xi2pdest + ".i2p"
	}
	if self.Config.XRealIP && xrealip != `` {
		return xrealip
	}
	return self.parseRemoteAddr(r.RemoteAddr, `127.0.0.1`)
}

func (self *Announce) parseRemoteAddr(in, def string) string {
	address := def
	r := regexp.MustCompile(`(.*):\d+$`)
	if match := r.FindStringSubmatch(in); len(match) == 2 {
		address = match[1]
	}
	return address
}
