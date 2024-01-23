package main

import (
	//"crypto/md5"
	"encoding/hex"
	_ "image/jpeg"
	_ "image/png"
	"log"

	"github.com/CrowdSurge/banner"
	//metrics "github.com/armon/go-metrics"
)

// AsciiHandler msr broker
func AsciiHandler(r *Req, i int) {
	//var t messageText

	//json.Unmarshal(m.Data, &t)
	hexDecodedStr, _ := hex.DecodeString(r.Hextr)

	hexEncodedStr := hex.EncodeToString([]byte(banner.PrintS(string(hexDecodedStr))))

	if err := EC.Publish("data.json.hash", &Req{Token: r.Token, Hextr: hexEncodedStr, Reply: r.Reply, Db: r.Db}); err != nil {
		log.Print(err)
	}
	REQ0 = REQ0 + 1

}
