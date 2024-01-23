package main

import (
	_ "image/jpeg"
	_ "image/png"
	"log"
	"strconv"
	"time"
	//metrics "github.com/armon/go-metrics"
)

//DataHandler export broker msg func
func DataHandler(r *Req, i int) {

	REQ0 = REQ0 + 1
	var err error
	var Payload string = r.Hextr

	tokenStr := strconv.FormatUint(uint64(r.Token), 10)

	if r.Db == "write" || r.Db == "rw" {

		_, err = STMTIns.Exec(r.Token, r.Hextr)

		if err != nil {
			log.Print(err)
		}
	}

	if r.Db == "read" || r.Db == "rw" {
		// additional iteration
		err = STMTSel.QueryRow(r.Token).Scan(&Payload) // WHERE number = 13

		if err != nil {
			Payload = r.Hextr
			log.Printf("QueryRowErr: %s", err) // proper error handling instead of panic in your app
		}

	}

	sec, _ := time.ParseDuration(AppCacheExpire)

	err = CACHE.Set(tokenStr, Payload, sec).Err()

	if err != nil {
		log.Print(err)
	}

	err = NC.Publish(r.Reply, []byte(tokenStr))

	if err != nil {
		log.Print(err)
	}

}
